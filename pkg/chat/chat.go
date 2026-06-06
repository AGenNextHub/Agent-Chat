// Package chat is the Chat core: the conversational runtime that runs admitted
// turns through the agent loop.
//
// The core has NO edges. Admission is the edge boundary's job, performed before
// the core is called; the core never imports or embeds the gate. The core
// expands only through the loop — adding capabilities widens its reach — never by
// absorbing the boundary. The edge composes around this core (see cmd/agennextd),
// not inside it.
package chat

import (
	"context"

	"github.com/agennext/agent-chat/pkg/loop"
)

// Core is the Chat runtime core. It owns no state and no edge; all turn state
// lives in the engine's stores, so the core scales horizontally (build for
// billions).
type Core struct {
	Engine *loop.Engine
}

// New builds a Chat core over an agent-loop engine.
func New(e *loop.Engine) *Core { return &Core{Engine: e} }

// Run executes one already-admitted turn through the loop. The caller (the edge
// boundary) admits first; the core trusts only admitted events. An admitted turn
// always resolves — a clean answer or a human escalation.
func (c *Core) Run(ctx context.Context, admitted loop.AdmittedEvent) (loop.Result, error) {
	return c.Engine.Run(ctx, admitted)
}
