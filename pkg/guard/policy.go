package guard

import "context"

// Request is the input to a policy decision: who is asking to invoke what.
// It is intentionally free of domain types so guard stays dependency-light.
type Request struct {
	// Tenant the call is scoped to.
	Tenant string
	// Principal making the call.
	Principal string
	// Capability being invoked.
	Capability string
	// Action being requested (e.g. "invoke").
	Action string
}

// Decision is the outcome of a policy evaluation.
type Decision struct {
	// Allow is true when the request is permitted.
	Allow bool
	// Reason explains the decision (inspectability / no hidden logic).
	Reason string
}

// Decider is an in-process policy decision point. The production binding is OPA
// evaluated in-process (a compiled wasm bundle) so the agent loop's GUARD step
// stays local and cheap on edge nodes; this interface is the seam.
type Decider interface {
	Decide(ctx context.Context, req Request) Decision
}

// StaticDecider is a deny-by-default decider that permits only an explicit
// allowlist of capabilities. It models OPA's default-deny posture for tests and
// single-node runs.
type StaticDecider struct {
	// Allowed is the set of capability names permitted to be invoked.
	Allowed map[string]bool
}

// NewStaticDecider builds a StaticDecider from the given capability names.
func NewStaticDecider(allowed ...string) StaticDecider {
	m := make(map[string]bool, len(allowed))
	for _, a := range allowed {
		m[a] = true
	}
	return StaticDecider{Allowed: m}
}

// Decide permits a request only if its capability is on the allowlist.
func (d StaticDecider) Decide(_ context.Context, req Request) Decision {
	if d.Allowed[req.Capability] {
		return Decision{Allow: true, Reason: "allowlisted"}
	}
	return Decision{Allow: false, Reason: "default-deny: capability not allowlisted"}
}
