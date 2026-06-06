// Package edge implements the Edge Protocol Gate: the only boundary that
// terminates external traffic and the encrypted overlay. The gate is itself an
// agent at the edge — it reasons about admission in a fixed, normative order
// (authenticate → authorize → scope → screen) and admits only requests that pass
// every step. Nothing speaks raw protocol to the Runtime Core. See
// docs/PROTOCOL.md and docs/THREAT_MODEL.md.
package edge

import (
	"context"
	"errors"
	"fmt"

	"github.com/agennext/agent-chat/pkg/capability"
	"github.com/agennext/agent-chat/pkg/event"
	"github.com/agennext/agent-chat/pkg/guard"
	"github.com/agennext/agent-chat/pkg/loop"
)

// Admission failure reasons. Each maps to a normative gate step.
var (
	// ErrUnauthenticated is returned when the caller cannot be authenticated.
	ErrUnauthenticated = errors.New("unauthenticated")
	// ErrUnauthorized is returned when authorization or policy denies the call.
	ErrUnauthorized = errors.New("unauthorized")
	// ErrOutOfScope is returned when the requested scope exceeds the contract.
	ErrOutOfScope = errors.New("out of scope")
	// ErrBlocked is returned when L1 screening flags the payload.
	ErrBlocked = errors.New("blocked by screen")
)

// Authenticator establishes the calling principal from an event. The production
// binding verifies mTLS/SPIFFE identity terminated from the overlay.
type Authenticator interface {
	Authenticate(ctx context.Context, e event.Event) (string, error)
}

// Authorizer answers whether a principal may invoke a capability for a tenant.
// The production binding is OpenFGA (relationship-based authorization).
type Authorizer interface {
	Authorize(ctx context.Context, principal, tenant, capName string) (bool, error)
}

// ScopeAuthority returns the scope a principal is GRANTED for a capability and
// tenant. Scope is derived here, from the authority — never supplied by the
// caller. The production binding derives grants from OpenFGA relations.
type ScopeAuthority interface {
	Granted(ctx context.Context, principal, tenant, capName string) (capability.Scope, error)
}

// Gate admits events into the Runtime Core after the full admission sequence.
type Gate struct {
	// Authn authenticates the caller.
	Authn Authenticator
	// Authz authorizes the caller (OpenFGA-shaped).
	Authz Authorizer
	// Grants derives the principal's granted scope (never caller-supplied).
	Grants ScopeAuthority
	// Decider applies in-process policy (OPA-shaped).
	Decider guard.Decider
	// Screener performs L1 prompt-injection screening on the payload.
	Screener guard.Screener
	// Prompt supplies the Guard Prompt prepended for generation.
	Prompt guard.PromptProvider
	// Registry holds capability contracts (for scope validation).
	Registry *capability.Registry
}

// Admit runs the normative gate sequence and returns an AdmittedEvent on
// success. The admitted scope is DERIVED — the principal's granted scope
// intersected with the capability's contract scope — never taken from the
// caller. A caller cannot widen its own authority by asking.
func (g *Gate) Admit(ctx context.Context, e event.Event) (loop.AdmittedEvent, error) {
	var zero loop.AdmittedEvent

	// Step 1: validate the envelope.
	if err := e.Validate(); err != nil {
		return zero, err
	}

	// Step 2: authenticate.
	principal, err := g.Authn.Authenticate(ctx, e)
	if err != nil {
		return zero, fmt.Errorf("%w: %w", ErrUnauthenticated, err)
	}

	// Step 3: authorize (relationship check + policy), default deny.
	ok, err := g.Authz.Authorize(ctx, principal, e.Tenant(), e.Capability())
	if err != nil {
		return zero, fmt.Errorf("%w: %w", ErrUnauthorized, err)
	}
	if !ok {
		return zero, fmt.Errorf("%w: relationship check failed", ErrUnauthorized)
	}
	if dec := g.Decider.Decide(ctx, guard.Request{
		Tenant: e.Tenant(), Principal: principal, Capability: e.Capability(), Action: "invoke",
	}); !dec.Allow {
		return zero, fmt.Errorf("%w: %s", ErrUnauthorized, dec.Reason)
	}

	// Step 4: DERIVE the effective scope = principal's grant ∩ contract scope.
	// Scope is never taken from the caller, so a request cannot self-authorize.
	contract, found := g.Registry.Get(e.Capability())
	if !found {
		return zero, fmt.Errorf("%w: unknown capability %q", ErrUnauthorized, e.Capability())
	}
	granted, err := g.Grants.Granted(ctx, principal, e.Tenant(), e.Capability())
	if err != nil {
		return zero, fmt.Errorf("%w: %w", ErrUnauthorized, err)
	}
	effective := granted.Intersect(contract.Scope)
	if effective.Empty() {
		return zero, fmt.Errorf("%w: principal has no granted scope for %q", ErrOutOfScope, e.Capability())
	}

	// Step 5: L1 screen the (untrusted) payload before it reaches the runtime.
	if v := g.Screener.Screen(ctx, string(e.Data), guard.OriginUser); v.Malicious {
		return zero, fmt.Errorf("%w: %s", ErrBlocked, v.Reason)
	}

	// Step 6: forward — scoped (derived), screened, with the Guard Prompt attached.
	return loop.AdmittedEvent{
		Event:       e,
		Principal:   principal,
		Scope:       effective,
		GuardPrompt: g.Prompt.GuardPrompt(),
	}, nil
}
