// Command agennextd is a single-node demonstration of the AGenNext Chat spine:
// it wires the Edge Gate and the agent loop with in-memory bindings, admits a
// chat event, runs one turn end-to-end, and prints the inspectable trace.
//
// This is the end-to-end check for the M0 build pass: a real, dependency-free
// path from an inbound event to a screened, scoped, governed answer. Production
// bindings (NATS, OPA, OpenFGA, KServe, PostgreSQL) replace the in-memory stubs
// behind the same interfaces.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/agennext/agent-chat/pkg/capability"
	"github.com/agennext/agent-chat/pkg/edge"
	"github.com/agennext/agent-chat/pkg/event"
	"github.com/agennext/agent-chat/pkg/guard"
	"github.com/agennext/agent-chat/pkg/loop"
	"github.com/agennext/agent-chat/pkg/store"
)

// demoAuthn accepts the principal carried on the event (single-node stand-in for
// mTLS/SPIFFE identity termination at the overlay).
type demoAuthn struct{}

func (demoAuthn) Authenticate(_ context.Context, e event.Event) (string, error) {
	return e.Principal(), nil
}

// demoAuthz permits the demo tenant/principal (stand-in for OpenFGA).
type demoAuthz struct{}

func (demoAuthz) Authorize(_ context.Context, _, _, _ string) (bool, error) { return true, nil }

// demoReasoner: retrieve once, then answer from the observation. Deterministic
// stand-in for the model so the demo is reproducible.
type demoReasoner struct{ used bool }

func (r *demoReasoner) Reason(_ context.Context, s loop.State) (loop.Action, error) {
	if !r.used {
		r.used = true
		return loop.Action{
			Kind:       loop.ActionInvoke,
			Capability: "rag.retrieve",
			Input:      []byte("return policy"),
			Scope:      capability.Scope{Tenants: []string{"acme"}, Data: []string{"tenant://acme/kb/faq"}},
		}, nil
	}
	answer := "I could not find an answer."
	for _, o := range s.Scratch {
		if !o.Blocked && o.Output != "" {
			answer = "Based on our knowledge base: " + o.Output
		}
	}
	return loop.Action{Kind: loop.ActionAnswer, Answer: answer}, nil
}

// demoInvoker returns a benign knowledge-base passage.
type demoInvoker struct{}

func (demoInvoker) Invoke(_ context.Context, _ string, _ []byte) (loop.Output, error) {
	return loop.Output{
		Data:   []byte("returns are accepted within 30 days of delivery with proof of purchase."),
		Origin: guard.OriginRetrieved,
	}, nil
}

func main() {
	ctx := context.Background()

	reg := capability.NewRegistry()
	must(reg.Register(capability.Contract{
		Name:       "rag.retrieve",
		Version:    "0.1.0",
		Provides:   []string{"retrieve(query) -> passages"},
		Scope:      capability.Scope{Tenants: []string{"acme"}, Data: []string{"tenant://acme/kb/*"}},
		Policy:     "opa://policies/rag-retrieve.rego",
		AuthZ:      "openfga://type/capability/relation/can_invoke",
		Artifact:   "oci://registry/agennext/rag-retrieve@sha256:demo",
		Sandbox:    capability.SandboxIsolated,
		Idempotent: true,
	}))

	gate := &edge.Gate{
		Authn:    demoAuthn{},
		Authz:    demoAuthz{},
		Decider:  guard.NewStaticDecider("rag.retrieve"),
		Screener: guard.NewHeuristicScreener(),
		Prompt:   guard.NewStaticPrompt(),
		Registry: reg,
	}

	engine := &loop.Engine{
		Reasoner: &demoReasoner{},
		Invoker:  demoInvoker{},
		Registry: reg,
		Screener: guard.NewHeuristicScreener(),
		Decider:  guard.NewStaticDecider("rag.retrieve"),
		Ctx:      store.NewMemContextStore(),
		Mem:      store.NewMemMemoryStore(),
		Dedupe:   loop.NewMemDeduper(),
		Budget:   loop.DefaultBudget(),
	}

	in := event.New("evt-1", "channel/web", "chat.message.v1", "session-acme-1",
		"acme", "user-1", "rag.retrieve", []byte("What is your return policy?"))
	requested := capability.Scope{Tenants: []string{"acme"}, Data: []string{"tenant://acme/kb/faq"}}

	admitted, err := gate.Admit(ctx, in, requested)
	if err != nil {
		fmt.Fprintln(os.Stderr, "gate denied:", err)
		os.Exit(1)
	}

	res, err := engine.Run(ctx, admitted)
	if err != nil {
		fmt.Fprintln(os.Stderr, "loop error:", err)
		os.Exit(1)
	}

	out := map[string]any{
		"session":    res.SessionID,
		"answer":     res.Answer,
		"iterations": res.Iterations,
		"stopped_by": res.StoppedBy,
		"trace":      res.Trace,
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	must(enc.Encode(out))
}

func must(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "fatal:", err)
		os.Exit(1)
	}
}
