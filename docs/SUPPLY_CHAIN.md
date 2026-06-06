# Supply Chain & SBOM Policy

> **No component is adopted without an SBOM and a supply-chain risk check.
> Only trusted sources.** This policy gates every dependency and image.

## Posture

1. **Minimal dependencies.** The core (`pkg/`, `cmd/`) is implemented in **pure Go
   standard library — zero third-party modules.** `go.mod` has no `require` block. This
   makes the SBOM trivially auditable and removes transitive-dependency risk by
   construction.
2. **Trusted sources only.** Permitted sources, in order of preference:
   - The Go standard library and official `golang.org/x` modules.
   - CNCF **Graduated/Incubating** projects (see STACK.md), pinned by digest.
   - Official upstream container images, pinned by `@sha256:` digest — never floating tags.
3. **Maturity-ranked.** New components are chosen Graduated > Incubating > Sandbox, with
   an ADR recording the decision (see `docs/adr/`).
4. **No license lock.** Permissive/OSI licenses only (Apache-2.0, BSD, MIT, MPL-2.0).
   AGPL/SSPL/BUSL/RSAL-relicensed projects are rejected — see STACK.md.

## Checks (run in CI)

| Check | Tool (trusted) | Gate |
|---|---|---|
| Vulnerabilities (Go) | `govulncheck` (official Go) | fail on known vulns |
| SBOM generation | Syft (CycloneDX/SPDX) | artifact published per build |
| Image/SBOM vuln scan | Grype | fail on fixable High/Critical |
| Image provenance | cosign / SLSA provenance | signed, verifiable images |
| Dependency review | CI dependency-review | block disallowed licenses |

## Adding a dependency (required procedure)

Before any `go get` or new image:

1. Generate/inspect its **SBOM**.
2. Run a **vulnerability scan** and review transitive deps.
3. Verify **license** is permitted (no license lock).
4. Confirm **source trust** (official/CNCF, signed, digest-pinnable).
5. Record an **ADR** with the maturity rationale.
6. Pin by **digest** and commit the SBOM delta.

A dependency that fails any step is not adopted. This procedure is itself a governance
gate (see PRINCIPLES.md).
