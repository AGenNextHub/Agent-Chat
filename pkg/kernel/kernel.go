// Package kernel is the kernel core: the control-plane admission/reconcile loop
// that turns desired Capability contracts into an admitted Registry. It is pure
// Go; the Kubernetes operator (controller-runtime) wraps this core and is gated
// by the supply-chain policy (docs/SUPPLY_CHAIN.md). This is the "kernel control
// loop" expressed without a cluster dependency.
package kernel

import "github.com/agennext/agent-chat/pkg/capability"

// Core is the kernel's admission engine over the capability registry.
type Core struct {
	Registry *capability.Registry
}

// New builds a kernel core over a fresh registry.
func New() *Core { return &Core{Registry: capability.NewRegistry()} }

// Reconcile drives a desired set of contracts into the admitted set — the
// control-loop step the operator invokes from CRD state. It admits what it can
// and reports per-contract failures; one bad contract never blocks the others
// (design for failure). The result is deterministic given the input.
func (k *Core) Reconcile(desired []capability.Contract) (admitted []string, failures map[string]error) {
	failures = make(map[string]error)
	for _, c := range desired {
		if err := k.Registry.Register(c); err != nil {
			failures[c.Name] = err
			continue
		}
		admitted = append(admitted, c.Name)
	}
	return admitted, failures
}
