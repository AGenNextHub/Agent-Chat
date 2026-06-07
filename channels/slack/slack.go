// Package slack is the AGenNext SlackBot channel — a pure-stdlib adapter mapping
// Slack's Events API to the core's CloudEvent contract. No SDK: net/http +
// encoding/json + crypto/hmac. The Adapter is translation only; the Bridge wires
// it to the gate and the chat core.
package slack

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/agennext/agent-chat/pkg/chat"
	"github.com/agennext/agent-chat/pkg/edge"
	"github.com/agennext/agent-chat/pkg/event"
	"github.com/agennext/agent-chat/pkg/loop"
)

// Adapter maps Slack messages to/from the platform's CloudEvent contract.
type Adapter struct {
	SigningSecret string // verifies inbound Slack request signatures
	Tenant        string // tenant these messages belong to
	Capability    string // capability they invoke, e.g. "rag.retrieve"
}

// Name identifies the channel.
func (Adapter) Name() string { return "slack" }

// slackEnvelope is the subset of Slack's Events API payload we read.
type slackEnvelope struct {
	Type      string `json:"type"`      // "url_verification" | "event_callback"
	Challenge string `json:"challenge"` // present on url_verification
	EventID   string `json:"event_id"`
	Event     struct {
		Type    string `json:"type"` // "message" | "app_mention"
		Text    string `json:"text"`
		User    string `json:"user"`
		Channel string `json:"channel"`
		BotID   string `json:"bot_id"` // set on bot messages (ignore to avoid loops)
	} `json:"event"`
}

// Inbound translates a Slack event payload into a CloudEvent. The session is the
// Slack channel id; the principal is the Slack user id.
func (a Adapter) Inbound(_ context.Context, raw []byte) (event.Event, error) {
	var e slackEnvelope
	if err := json.Unmarshal(raw, &e); err != nil {
		return event.Event{}, err
	}
	if e.Event.BotID != "" {
		return event.Event{}, fmt.Errorf("ignoring bot message")
	}
	return event.New(
		e.EventID, "channel/slack", "chat.message.v1",
		e.Event.Channel, a.Tenant, e.Event.User, a.Capability,
		[]byte(e.Event.Text),
	), nil
}

// Outbound renders a turn result into a Slack chat.postMessage body.
func (Adapter) Outbound(_ context.Context, res loop.Result) ([]byte, error) {
	return json.Marshal(map[string]string{"text": res.Answer})
}

// verifySignature checks Slack's v0 request signature (HMAC-SHA256 over
// "v0:<timestamp>:<body>") in constant time.
func verifySignature(secret, timestamp, signature string, body []byte) bool {
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = io.WriteString(mac, "v0:"+timestamp+":")
	mac.Write(body)
	want := "v0=" + hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(want), []byte(signature))
}

// Bridge serves Slack's Events API and connects it to the core.
type Bridge struct {
	Adapter Adapter
	Gate    *edge.Gate
	Core    *chat.Core
	// Post delivers the answer back to a Slack channel. Injectable for tests; the
	// production value posts to chat.postMessage with the bot token.
	Post func(ctx context.Context, channel, text string) error
	// Now is injectable for tests (timestamp-freshness check).
	Now func() time.Time
}

func (b *Bridge) now() time.Time {
	if b.Now != nil {
		return b.Now()
	}
	return time.Now()
}

// Handler returns the Slack events endpoint.
func (b *Bridge) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /slack/events", b.Events)
	return mux
}

// Events handles a Slack Events API request: verify signature + freshness,
// answer url_verification, else map inbound -> gate -> core -> post back.
func (b *Bridge) Events(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 1<<20))
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	// Reject stale requests (>5 min) and bad signatures (replay + authenticity).
	ts := r.Header.Get("X-Slack-Request-Timestamp")
	if n, err := strconv.ParseInt(ts, 10, 64); err != nil || absDuration(b.now().Sub(time.Unix(n, 0))) > 5*time.Minute {
		http.Error(w, "stale", http.StatusUnauthorized)
		return
	}
	if !verifySignature(b.Adapter.SigningSecret, ts, r.Header.Get("X-Slack-Signature"), body) {
		http.Error(w, "bad signature", http.StatusUnauthorized)
		return
	}

	// URL verification handshake.
	var probe slackEnvelope
	_ = json.Unmarshal(body, &probe)
	if probe.Type == "url_verification" {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"challenge": probe.Challenge})
		return
	}

	ev, err := b.Adapter.Inbound(r.Context(), body)
	if err != nil {
		w.WriteHeader(http.StatusOK) // ack non-actionable events (e.g. bot echoes)
		return
	}
	admitted, err := b.Gate.Admit(r.Context(), ev)
	if err != nil {
		w.WriteHeader(http.StatusOK) // denied at the gate; ack so Slack does not retry
		return
	}
	res, err := b.Core.Run(r.Context(), admitted)
	if err != nil {
		http.Error(w, "internal", http.StatusInternalServerError)
		return
	}
	if b.Post != nil {
		_ = b.Post(r.Context(), ev.Subject, res.Answer)
	}
	w.WriteHeader(http.StatusOK)
}

func absDuration(d time.Duration) time.Duration {
	if d < 0 {
		return -d
	}
	return d
}

// PostMessage returns a Bridge.Post bound to a bot token; it delivers the answer
// to a Slack channel via chat.postMessage. Pure net/http — no SDK.
func PostMessage(botToken string) func(ctx context.Context, channel, text string) error {
	return func(ctx context.Context, channel, text string) error {
		payload, err := json.Marshal(map[string]string{"channel": channel, "text": text})
		if err != nil {
			return err
		}
		req, err := http.NewRequestWithContext(ctx, http.MethodPost,
			"https://slack.com/api/chat.postMessage", bytes.NewReader(payload))
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", "Bearer "+botToken)
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		defer func() { _ = res.Body.Close() }()
		var out struct {
			OK    bool   `json:"ok"`
			Error string `json:"error"`
		}
		if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
			return err
		}
		if !out.OK {
			return fmt.Errorf("slack postMessage: %s", out.Error)
		}
		return nil
	}
}
