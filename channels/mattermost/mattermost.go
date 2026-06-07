// Package mattermost is the AGenNext channel for Mattermost — open-source, the
// on-brand peer platform. It implements pkg/channel.Adapter over Mattermost
// outgoing webhooks (form-encoded, token-verified). Pure stdlib; the reply is
// the synchronous webhook response, so no outbound API call is needed.
package mattermost

import (
	"context"
	"crypto/hmac"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/agennext/agent-chat/pkg/chat"
	"github.com/agennext/agent-chat/pkg/edge"
	"github.com/agennext/agent-chat/pkg/event"
	"github.com/agennext/agent-chat/pkg/loop"
)

// Adapter maps Mattermost outgoing-webhook posts to/from the CloudEvent contract.
type Adapter struct {
	Token      string // the outgoing-webhook token Mattermost sends (verify it)
	Tenant     string
	Capability string
}

// Name identifies the channel.
func (Adapter) Name() string { return "mattermost" }

// Inbound parses a Mattermost outgoing-webhook (form-encoded) body into a
// CloudEvent. Session is the channel id; principal is the user id.
func (a Adapter) Inbound(_ context.Context, raw []byte) (event.Event, error) {
	v, err := url.ParseQuery(string(raw))
	if err != nil {
		return event.Event{}, err
	}
	if v.Get("user_id") == "" || v.Get("text") == "" {
		return event.Event{}, fmt.Errorf("not an actionable post")
	}
	return event.New(
		v.Get("post_id"), "channel/mattermost", "chat.message.v1",
		v.Get("channel_id"), a.Tenant, v.Get("user_id"), a.Capability,
		[]byte(v.Get("text")),
	), nil
}

// Outbound renders a result as the webhook response body.
func (Adapter) Outbound(_ context.Context, res loop.Result) ([]byte, error) {
	return json.Marshal(map[string]string{"text": res.Answer})
}

func tokenOK(want, got string) bool {
	return want != "" && hmac.Equal([]byte(want), []byte(got))
}

// Bridge serves the Mattermost outgoing-webhook endpoint at the edge.
type Bridge struct {
	Adapter Adapter
	Gate    *edge.Gate
	Core    *chat.Core
}

// Handler returns the webhook endpoint.
func (b *Bridge) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /mattermost/hooks", b.Hook)
	return mux
}

// Hook verifies the token, maps the post to a CloudEvent, runs it through the
// gate and core, and returns the answer as the webhook reply.
func (b *Bridge) Hook(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 1<<20))
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	v, _ := url.ParseQuery(string(body))
	if !tokenOK(b.Adapter.Token, v.Get("token")) {
		http.Error(w, "bad token", http.StatusUnauthorized)
		return
	}
	reply := func(text string) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"text": text})
	}
	ev, err := b.Adapter.Inbound(r.Context(), body)
	if err != nil {
		reply("") // non-actionable; ack with no message
		return
	}
	admitted, err := b.Gate.Admit(r.Context(), ev)
	if err != nil {
		reply("") // denied at the gate; ack silently (no detail bled)
		return
	}
	res, err := b.Core.Run(r.Context(), admitted)
	if err != nil {
		http.Error(w, "internal", http.StatusInternalServerError)
		return
	}
	reply(res.Answer)
}
