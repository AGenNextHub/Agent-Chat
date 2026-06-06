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

func TestScopeEmpty(t *testing.T) {
	t.Parallel()
	if !(Scope{}).Empty() {
		t.Fatal("zero scope must be empty")
	}
	if (Scope{Tenants: []string{"acme"}}).Empty() {
		t.Fatal("non-empty scope must not report empty")
	}
}

func TestScopeIntersect(t *testing.T) {
	t.Parallel()
	contract := Scope{Tenants: []string{"acme"}, Data: []string{"tenant://acme/kb/*"}, Tools: []string{"index.search"}}
	tests := []struct {
		name      string
		grant     Scope
		wantEmpty bool
		check     Scope // a sub-scope the result must contain (if !wantEmpty)
	}{
		{"grant within contract", Scope{Tenants: []string{"acme"}, Data: []string{"tenant://acme/kb/faq"}}, false, Scope{Tenants: []string{"acme"}, Data: []string{"tenant://acme/kb/faq"}}},
		{"wildcard grant -> contract", Scope{Tenants: []string{"*"}, Data: []string{"*"}, Tools: []string{"*"}}, false, Scope{Tenants: []string{"acme"}, Data: []string{"tenant://acme/kb/x"}, Tools: []string{"index.search"}}},
		{"disjoint tenant -> empty", Scope{Tenants: []string{"evil"}}, true, Scope{}},
		{"disjoint data only -> empty", Scope{Data: []string{"tenant://other/*"}}, true, Scope{}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.grant.Intersect(contract)
			if got.Empty() != tc.wantEmpty {
				t.Fatalf("Intersect empty=%v, want %v (got %+v)", got.Empty(), tc.wantEmpty, got)
			}
			if !tc.wantEmpty && !got.Contains(tc.check) {
				t.Fatalf("effective scope %+v does not contain %+v", got, tc.check)
			}
		})
	}
}

func TestScopeIntersectCommutativeOnTenants(t *testing.T) {
	t.Parallel()
	a := Scope{Tenants: []string{"acme", "globex"}}
	b := Scope{Tenants: []string{"globex"}}
	if !a.Intersect(b).Contains(Scope{Tenants: []string{"globex"}}) {
		t.Fatal("intersection should keep the common tenant")
	}
	if a.Intersect(b).Contains(Scope{Tenants: []string{"acme"}}) {
		t.Fatal("intersection must drop the non-common tenant")
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
