package kernel

import (
	"testing"

	"github.com/agennext/agent-chat/pkg/capability"
)

func goodContract() capability.Contract {
	return capability.Contract{
		Name:     "rag.retrieve",
		Version:  "0.1.0",
		Provides: []string{"retrieve"},
		Sandbox:  capability.SandboxIsolated,
	}
}

func TestKernelReconcile(t *testing.T) {
	t.Parallel()
	k := New()
	bad := capability.Contract{Name: "Bad"} // invalid: fails validation
	admitted, failures := k.Reconcile([]capability.Contract{goodContract(), bad})

	if len(admitted) != 1 || admitted[0] != "rag.retrieve" {
		t.Fatalf("admitted = %v, want [rag.retrieve]", admitted)
	}
	if _, ok := failures["Bad"]; !ok {
		t.Fatal("invalid contract must be reported as a failure")
	}
	if _, ok := k.Registry.Get("rag.retrieve"); !ok {
		t.Fatal("valid contract must be admitted into the registry")
	}
	// Design for failure: one bad contract does not block the good one.
}
