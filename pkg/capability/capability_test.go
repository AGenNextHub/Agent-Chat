package capability

import (
	"errors"
	"testing"
)

func validContract() Contract {
	return Contract{
		Name:     "rag.retrieve",
		Version:  "0.1.0",
		Provides: []string{"retrieve(query) -> passages"},
		Scope:    Scope{Tenants: []string{"acme"}, Data: []string{"tenant://acme/kb/*"}},
		Sandbox:  SandboxIsolated,
	}
}

func TestContractValidate(t *testing.T) {
	t.Parallel()
	if err := validContract().Validate(); err != nil {
		t.Fatalf("valid contract rejected: %v", err)
	}
	bad := []Contract{
		{Name: "Bad", Version: "0.1.0", Provides: []string{"x"}, Sandbox: SandboxIsolated},
		{Name: "rag.retrieve", Version: "0.1", Provides: []string{"x"}, Sandbox: SandboxIsolated},
		{Name: "rag.retrieve", Version: "0.1.0", Sandbox: SandboxIsolated},
		{Name: "rag.retrieve", Version: "0.1.0", Provides: []string{"x"}, Sandbox: "weird"},
	}
	for i, c := range bad {
		if err := c.Validate(); !errors.Is(err, ErrInvalidContract) {
			t.Fatalf("case %d: expected ErrInvalidContract, got %v", i, err)
		}
	}
}

func TestScopeContains(t *testing.T) {
	t.Parallel()
	granted := Scope{Tenants: []string{"acme"}, Data: []string{"tenant://acme/kb/*"}, Tools: []string{"index.search"}}
	tests := []struct {
		name string
		req  Scope
		want bool
	}{
		{"in scope", Scope{Tenants: []string{"acme"}, Data: []string{"tenant://acme/kb/faq"}}, true},
		{"cross tenant denied", Scope{Tenants: []string{"evil"}}, false},
		{"data outside glob denied", Scope{Data: []string{"tenant://acme/secrets/x"}}, false},
		{"tool not granted denied", Scope{Tools: []string{"net.fetch"}}, false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := granted.Contains(tc.req); got != tc.want {
				t.Fatalf("Contains(%+v) = %v, want %v", tc.req, got, tc.want)
			}
		})
	}
}

func TestScopeWildcard(t *testing.T) {
	t.Parallel()
	any := Scope{Tenants: []string{"*"}, Data: []string{"*"}, Tools: []string{"*"}}
	if !any.Contains(Scope{Tenants: []string{"anyone"}, Data: []string{"anything"}, Tools: []string{"whatever"}}) {
		t.Fatal("wildcard scope should contain anything")
	}
}

func TestRegistry(t *testing.T) {
	t.Parallel()
	r := NewRegistry()
	if err := r.Register(validContract()); err != nil {
		t.Fatalf("register: %v", err)
	}
	if _, ok := r.Get("rag.retrieve"); !ok {
		t.Fatal("expected capability to be registered")
	}
	// Dependency on an unregistered capability must be rejected.
	dep := validContract()
	dep.Name = "rag.answer"
	dep.Requires = []string{"missing.cap"}
	if err := r.Register(dep); !errors.Is(err, ErrInvalidContract) {
		t.Fatalf("expected dependency rejection, got %v", err)
	}
	if len(r.List()) != 1 {
		t.Fatalf("expected 1 registered contract, got %d", len(r.List()))
	}
}
