# Roadmap

Managed under the delivery pipeline **Define → Design → Align → Take Approval → Act (one
pass) → Self-review → Sanity → End-to-End Check → Confirm** (see SPEC.md). Versioning is
SemVer; releases are cut from `main` via signed, SBOM-attached artifacts.

## Now — v0.1.x (this pass): the loop spine
- [x] Core primitives: `event`, `capability`, `store`, `guard`, `loop`, `edge` (pure stdlib)
- [x] Agent loop: bounded, guard-before-act, screen-on-tool-output, transactional turn,
      idempotent replay
- [x] Edge Gate: authn → authz → scope → L1 screen admission sequence
- [x] End-to-end runnable demo (`cmd/agennextd`) with inspectable trace
- [x] GitOps + OSS scaffolding (CI, Helm, Argo CD, supply-chain checks)

## Next — v0.2.x: contracts & control plane
- [ ] Machine-readable Capability schema + Kernel CRDs (`Tenant`/`Capability`/
      `AgentSession`/`Policy`)
- [ ] controller-runtime operator (the kernel control loop)
- [ ] OPA-in-process Decider binding; OpenFGA Authorizer binding

## Later — v0.3.x: runtime & delivery
- [ ] NATS Bus binding (CloudEvents over NATS); streaming EMIT
- [ ] RAG ingestion with PII screening; PostgreSQL+pgvector memory; Valkey context
- [ ] KServe/vLLM inference; LiteLLM provider-agnostic gateway
- [ ] k3s + Cilium overlay + Argo CD app-of-apps

## Hardening — v0.4.x: security & governance gates
- [ ] Trained injection detector (GenTel-Shield) behind `guard.Screener`
- [ ] Fairness evaluation in CI (no bias); inspectability/audit endpoints
- [ ] Multi-channel/multi-device adapters; data-residency controls

## Product management
- Issues use the templates in `.github/ISSUE_TEMPLATE/`; changes land via PRs against the
  acceptance criteria in SPEC.md.
- Each significant decision is recorded as an ADR in `docs/adr/`.
- Definition of Done = the closing gates (Self-review → Sanity → End-to-End → Confirm) pass
  and the relevant acceptance criterion is demonstrably met.
