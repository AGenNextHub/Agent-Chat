# 2. Pure standard-library core (zero third-party dependencies)

Date: 2026-06-06

## Status
Accepted

## Context
Requirements mandate an SBOM and supply-chain risk check before adopting any component, and
"only trusted sources." Third-party dependencies introduce transitive risk and SBOM noise.

## Decision
The core (`pkg/`, `cmd/`) is implemented using only the Go standard library. `go.mod` has no
`require` block. CloudEvents, the policy/authz/screening seams, and stores are modeled as our
own types and interfaces, with production systems (NATS, OPA, OpenFGA, KServe, PostgreSQL)
bound behind those interfaces later.

## Consequences
- The SBOM is trivially auditable; transitive-dependency risk is eliminated by construction.
- Adding any dependency requires the procedure in `docs/SUPPLY_CHAIN.md` and an ADR.
- Some functionality (e.g. a trained injection detector) arrives later, behind an interface,
  rather than via an immediate heavyweight dependency.
