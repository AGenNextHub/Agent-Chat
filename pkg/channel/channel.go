// Package channel defines the composable channel contract.
//
// The core is headless and has one contract — CloudEvents. A channel is an
// adapter that connects a touch point (web, Slack, Mattermost, Matrix, WhatsApp)
// to the core: it maps an inbound message to a CloudEvent and an outbound result
// back to its medium. Add a channel by implementing Adapter; the core is
// untouched. Payloads are content-type agnostic (multimodal) via the event
// envelope (`event.Event.Data` + `DataContentType`).
package channel

import (
	"context"

	"github.com/agennext/agent-chat/pkg/event"
	"github.com/agennext/agent-chat/pkg/loop"
)

// Adapter connects one touch point to the core, both directions. It speaks the
// channel's medium on one side and the platform's CloudEvent contract on the
// other; it holds no business logic — translation only.
type Adapter interface {
	// Name identifies the channel, e.g. "web", "slack", "matrix".
	Name() string
	// Inbound translates a raw channel message into a CloudEvent for the core.
	Inbound(ctx context.Context, raw []byte) (event.Event, error)
	// Outbound renders a turn result back into the channel's format (multimodal).
	Outbound(ctx context.Context, res loop.Result) ([]byte, error)
}
