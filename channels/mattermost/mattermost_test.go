package mattermost

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/agennext/agent-chat/pkg/capability"
	"github.com/agennext/agent-chat/pkg/chat"
	"github.com/agennext/agent-chat/pkg/edge"
	"github.com/agennext/agent-chat/pkg/event"
	"github.com/agennext/agent-chat/pkg/guard"
	"github.com/agennext/agent-chat/pkg/loop"
	"github.com/agennext/agent-chat/pkg/store"
)

const token = "mm-token"

func form(values map[string]string) []byte {
	v := url.Values{}
	for k, val := range values {
		v.Set(k, val)
	}
	return []byte(v.Encode())
}

func TestInbound(t *testing.T) {
	t.Parallel()
	a := Adapter{Tenant: "acme", Capability: "rag.retrieve"}
	raw := form(map[string]string{
		"token": token, "post_id": "P1", "channel_id": "C1",
		"user_id": "U1", "text": "return policy?",
	})
	ev, err := a.Inbound(context.Background(), raw)
	if err != nil {
		t.Fatalf("inbound: %v", err)
	}
	if ev.Subject != "C1" || ev.Principal() != "U1" || ev.Tenant() != "acme" || string(ev.Data) != "return policy?" {
		t.Fatalf("unexpected event: %+v", ev)
	}
	// A post with no text/user is not actionable.
	if _, err := a.Inbound(context.Background(), form(map[string]string{"channel_id": "C1"})); err == nil {
		t.Fatal("empty post should not be actionable")
	}
}

func TestOutbound(t *testing.T) {
	t.Parallel()
	b, _ := Adapter{}.Outbound(context.Background(), loop.Result{Answer: "hello"})
	if string(b) != `{"text":"hello"}` {
		t.Fatalf("outbound: %s", b)
	}
}

func TestTokenOK(t *testing.T) {
	t.Parallel()
	if !tokenOK(token, token) {
		t.Fatal("matching token rejected")
	}
	if tokenOK(token, "wrong") {
		t.Fatal("wrong token accepted")
	}
	if tokenOK("", "") {
		t.Fatal("empty configured token must never pass")
	}
}

// --- core wiring for the hook flow ---

type okAuthn struct{}

func (okAuthn) Authenticate(_ context.Context, e event.Event) (string, error) { return e.Principal(), nil }

type okAuthz struct{}

func (okAuthz) Authorize(_ context.Context, _, _, _ string) (bool, error) { return true, nil }

type okGrants struct{}

func (okGrants) Granted(_ context.Context, _, _, _ string) (capability.Scope, error) {
	return capability.Scope{Tenants: []string{"acme"}, Data: []string{"tenant://acme/kb/*"}}, nil
}

type answerReasoner struct{ text string }

func (a answerReasoner) Reason(_ context.Context, _ loop.State) (loop.Action, error) {
	return loop.Action{Kind: loop.ActionAnswer, Answer: a.text}, nil
}

func newBridge(t *testing.T) *Bridge {
	t.Helper()
	reg := capability.NewRegistry()
	if err := reg.Register(capability.Contract{
		Name: "rag.retrieve", Version: "0.1.0", Provides: []string{"retrieve"},
		Scope:   capability.Scope{Tenants: []string{"acme"}, Data: []string{"tenant://acme/kb/*"}},
		Sandbox: capability.SandboxIsolated,
	}); err != nil {
		t.Fatalf("register: %v", err)
	}
	gate := &edge.Gate{
		Authn: okAuthn{}, Authz: okAuthz{}, Grants: okGrants{},
		Decider: guard.NewStaticDecider("rag.retrieve"), Screener: guard.NewHeuristicScreener(),
		Prompt: guard.NewStaticPrompt(), Registry: reg,
	}
	eng := &loop.Engine{
		Reasoner: answerReasoner{"hi from mattermost"}, Registry: reg,
		Screener: guard.NewHeuristicScreener(), Decider: guard.NewStaticDecider("rag.retrieve"),
		Ctx: store.NewMemContextStore(), Mem: store.NewMemMemoryStore(),
		Dedupe: loop.NewMemDeduper(), Budget: loop.DefaultBudget(),
	}
	return &Bridge{
		Adapter: Adapter{Token: token, Tenant: "acme", Capability: "rag.retrieve"},
		Gate:    gate, Core: chat.New(eng),
	}
}

func hook(t *testing.T, h http.Handler, body []byte) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, "/mattermost/hooks", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func TestHookFlow(t *testing.T) {
	t.Parallel()
	h := newBridge(t).Handler()
	body := form(map[string]string{
		"token": token, "post_id": "P9", "channel_id": "C1",
		"user_id": "U1", "text": "hi",
	})
	rec := hook(t, h, body)
	if rec.Code != http.StatusOK {
		t.Fatalf("code %d body=%s", rec.Code, rec.Body)
	}
	var out map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatalf("decode reply: %v", err)
	}
	if out["text"] != "hi from mattermost" {
		t.Fatalf("expected answer in webhook reply, got %q", out["text"])
	}
}

func TestBadTokenRejected(t *testing.T) {
	t.Parallel()
	h := newBridge(t).Handler()
	body := form(map[string]string{
		"token": "wrong", "post_id": "P9", "channel_id": "C1",
		"user_id": "U1", "text": "hi",
	})
	rec := hook(t, h, body)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}

func TestNonActionablePostAckedSilently(t *testing.T) {
	t.Parallel()
	h := newBridge(t).Handler()
	// Valid token but no text/user — Mattermost system messages etc.
	body := form(map[string]string{"token": token, "channel_id": "C1"})
	rec := hook(t, h, body)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 ack, got %d", rec.Code)
	}
	if strings.TrimSpace(rec.Body.String()) != `{"text":""}` {
		t.Fatalf("expected empty reply, got %s", rec.Body)
	}
}
