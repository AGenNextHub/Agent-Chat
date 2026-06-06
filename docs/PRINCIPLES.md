# Governance Principles (made testable)

Principles are only real if they can fail a test. Each tenet below has an enforcement
point and an observable check.

## No hidden logic

> Every guard prompt, routing rule, policy, and capability is inspectable at runtime.

- **Enforcement:** all behaviour flows from declared **capability contracts** + **OPA**
  policies; nothing runs outside an admitted contract.
- **Check:** for any response, the full set of active guard prompts, policies, and invoked
  capabilities is retrievable and attributable (audit trail via OpenTelemetry).
- **Fails if:** any decision cannot be traced to a declared contract/policy.

## No bias

> Fairness is an explicit, measured property of model behaviour — not an assumption.

- **Enforcement:** evaluation gates in CI over fairness test sets; configurable refusal
  and guardrail behaviour is uniform across protected attributes.
- **Check:** disparity metrics across demographic slices stay within an agreed threshold
  on the benchmark suite; results published, not hidden.
- **Fails if:** measured disparity exceeds threshold, or no fairness evaluation exists.

## Clear (inspectable) context

> Session context and memory are observable and operator-controllable.

- **Enforcement:** the `load context` / `clear context` / `load memory` / `clear memory`
  control surface is first-class and audited.
- **Check:** an operator can view and clear exactly what a session holds; clears are
  verifiable.
- **Fails if:** a session retains context/memory that cannot be inspected or cleared.

## Privacy by default

> PII is screened at ingestion; tenant data never crosses tenant boundaries.

- **Enforcement:** ingestion-time PII screening + de-identification; capability `scope`
  pinned to `$caller.tenant`; OpenFGA relations.
- **Check:** cross-tenant access attempts are structurally denied (tested in CI).
- **Fails if:** any path allows one tenant to read another's data.

## Least privilege (define scope)

> A capability runs only with the authority its contract declares.

- **Enforcement:** Kernel admission rejects requested scope ⊄ contract scope.
- **Check:** out-of-scope tool/data access is denied and logged.
- **Fails if:** a capability acts beyond its declared scope.

## Composable

> Layers and capabilities are independently deployable and swappable.

- **Enforcement:** OCI-packaged capabilities; contract-bound interfaces.
- **Check:** a capability can be replaced by another satisfying the same contract with no
  changes to its consumers.
- **Fails if:** swapping a contract-compatible capability breaks consumers.
