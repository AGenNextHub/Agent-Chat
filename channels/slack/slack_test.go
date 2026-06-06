package slack

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/agennext/agent-chat/pkg/capability"
	"github.com/agennext/agent-chat/pkg/chat"
	"github.com/agennext/agent-chat/pkg/edge"
	"github.com/agennext/agent-chat/pkg/event"
	"github.com/agennext/agent-chat/pkg/guard"
	"github.com/agennext/agent-chat/pkg/loop"
	"github.com/agennext/agent-chat/pkg/store"
)

const secret = "shhh"

func sign(ts string, body []byte) string {
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = io.WriteString(mac, "v0:"+ts+":")
	mac.Write(body)
	return "v0=" + hex.EncodeToString(mac.Sum(nil))
}

func TestVerifySignature(t *testing.T) {
	t.Parallel()
	body := []byte(`{"hi":1}`)
	ts := "1700000000"
	if !verifySignature(secret, ts, sign(ts, body), body) {
		t.Fatal("valid signature rejected")
	}
	if verifySignature(secret, ts, sign(ts, body), []byte("tampered")) {
		t.Fatal("tampered body accepted")
	}
}

func TestInbound(t *testing.T) {
	t.Parallel()
	a := Adapter{Tenant: "acme", Capability: "rag.retrieve"}
	raw := []byte(`{"type":"event_callback","event_id":"Ev1","event":{"type":"message","text":"return policy?","user":"U1","channel":"C1"}}`)
	ev, err := a.Inbound(context.Background(), raw)
	if err != nil {
		t.Fatalf("inbound: %v", err)
	}
	if ev.Subject != "C1" || ev.Principal() != "U1" || ev.Tenant() != "acme" || string(ev.Data) != "return policy?" {
		t.Fatalf("unexpected event: %+v", ev)
	}
	// Bot messages are ignored to avoid loops.
	if _, err := a.Inbound(context.Background(), []byte(`{"event":{"bot_id":"B1","text":"x"}}`)); err == nil {
		t.Fatal("bot message should be ignored")
	}
}

func TestOutbound(t *testing.T) {
	t.Parallel()
	b, _ := Adapter{}.Outbound(context.Background(), loop.Result{Answer: "hello"})
	if string(b) != `{"text":"hello"}` {
		t.Fatalf("outbound: %s", b)
	}
}

// --- core wiring for the events flow ---

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

func newBridge(t *testing.T, capture *string) *Bridge {
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
		Reasoner: answerReasoner{"hi from slack"}, Registry: reg,
		Screener: guard.NewHeuristicScreener(), Decider: guard.NewStaticDecider("rag.retrieve"),
		Ctx: store.NewMemContextStore(), Mem: store.NewMemMemoryStore(),
		Dedupe: loop.NewMemDeduper(), Budget: loop.DefaultBudget(),
	}
	fixed := time.Unix(1700000000, 0)
	return &Bridge{
		Adapter: Adapter{SigningSecret: secret, Tenant: "acme", Capability: "rag.retrieve"},
		Gate:    gate, Core: chat.New(eng),
		Post: func(_ context.Context, _, text string) error { *capture = text; return nil },
		Now:  func() time.Time { return fixed },
	}
}

func post(t *testing.T, h http.Handler, body []byte) *httptest.ResponseRecorder {
	t.Helper()
	ts := "1700000000"
	req := httptest.NewRequest(http.MethodPost, "/slack/events", bytes.NewReader(body))
	req.Header.Set("X-Slack-Request-Timestamp", ts)
	req.Header.Set("X-Slack-Signature", sign(ts, body))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func TestURLVerification(t *testing.T) {
	t.Parallel()
	var got string
	h := newBridge(t, &got).Handler()
	rec := post(t, h, []byte(`{"type":"url_verification","challenge":"abc123"}`))
	if rec.Code != http.StatusOK {
		t.Fatalf("code %d", rec.Code)
	}
	var out map[string]string
	_ = json.Unmarshal(rec.Body.Bytes(), &out)
	if out["challenge"] != "abc123" {
		t.Fatalf("challenge not echoed: %s", rec.Body)
	}
}

func TestEventsFlow(t *testing.T) {
	t.Parallel()
	var posted string
	h := newBridge(t, &posted).Handler()
	body := []byte(`{"type":"event_callback","event_id":"Ev9","event":{"type":"message","text":"hi","user":"U1","channel":"C1"}}`)
	rec := post(t, h, body)
	if rec.Code != http.StatusOK {
		t.Fatalf("code %d body=%s", rec.Code, rec.Body)
	}
	if posted != "hi from slack" {
		t.Fatalf("expected answer posted to Slack, got %q", posted)
	}
}

func TestBadSignatureRejected(t *testing.T) {
	t.Parallel()
	var got string
	h := newBridge(t, &got).Handler()
	body := []byte(`{"type":"event_callback"}`)
	req := httptest.NewRequest(http.MethodPost, "/slack/events", bytes.NewReader(body))
	req.Header.Set("X-Slack-Request-Timestamp", strconv.FormatInt(time.Unix(1700000000, 0).Unix(), 10))
	req.Header.Set("X-Slack-Signature", "v0=deadbeef")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}
