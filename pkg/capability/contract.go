package capability

import (
	"errors"
	"fmt"
	"regexp"
)

// SandboxProfile is the execution isolation a capability runs under.
type SandboxProfile string

const (
	// SandboxIsolated is the default: no network/tool access beyond declared scope.
	SandboxIsolated SandboxProfile = "isolated"
	// SandboxRestricted permits a constrained set of host interactions.
	SandboxRestricted SandboxProfile = "restricted"
	// SandboxTrusted is for first-party platform capabilities only.
	SandboxTrusted SandboxProfile = "trusted"
)

// nameRe constrains capability names to a dotted, lowercase identifier form.
var nameRe = regexp.MustCompile(`^[a-z][a-z0-9]*(\.[a-z][a-z0-9]*)+$`)

// semverRe is a minimal semantic-version matcher (major.minor.patch).
var semverRe = regexp.MustCompile(`^\d+\.\d+\.\d+$`)

// Contract is the complete, inspectable description of a capability. A feature
// is a contract; if it is not declared here, the platform will not run it.
type Contract struct {
	// Name is a dotted identifier, e.g. "rag.retrieve".
	Name string
	// Version is a semantic version, e.g. "0.1.0".
	Version string
	// Provides lists the interfaces this capability exposes.
	Provides []string
	// Requires lists capabilities/services this one depends on.
	Requires []string
	// Scope is the least-privilege authority granted to this capability.
	Scope Scope
	// Policy references an OPA bundle that must pass for admission/execution.
	Policy string
	// AuthZ references the OpenFGA relation governing who may invoke it.
	AuthZ string
	// Artifact is the digest-pinned OCI reference to the implementation.
	Artifact string
	// Sandbox is the execution profile (default isolated).
	Sandbox SandboxProfile
	// Idempotent declares the capability safe to retry without duplicate side
	// effects. Non-idempotent (side-effecting) capabilities must be guarded by
	// the loop's at-least-once dedupe before ACT. Defaults to false (unsafe).
	Idempotent bool
}

// ErrInvalidContract is returned by Validate for a malformed contract.
var ErrInvalidContract = errors.New("invalid contract")

// Validate enforces the structural invariants a contract must satisfy before
// the Kernel will admit it.
func (c Contract) Validate() error {
	if !nameRe.MatchString(c.Name) {
		return fmt.Errorf("%w: name %q must be a dotted lowercase identifier", ErrInvalidContract, c.Name)
	}
	if !semverRe.MatchString(c.Version) {
		return fmt.Errorf("%w: version %q must be semantic (major.minor.patch)", ErrInvalidContract, c.Version)
	}
	if len(c.Provides) == 0 {
		return fmt.Errorf("%w: at least one provided interface is required", ErrInvalidContract)
	}
	switch c.Sandbox {
	case SandboxIsolated, SandboxRestricted, SandboxTrusted:
	case "":
		return fmt.Errorf("%w: sandbox profile is required", ErrInvalidContract)
	default:
		return fmt.Errorf("%w: unknown sandbox profile %q", ErrInvalidContract, c.Sandbox)
	}
	return nil
}
