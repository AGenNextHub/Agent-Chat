# AGenNext Chat — Specification

**Delivery pipeline (every stage gated by the one before):**

> **Define → Design → Align → Take Approval → Act (one pass) → Self-review → Sanity → End-to-End Check → Confirm**

- **Define** — what must be true (requirements, acceptance criteria).
- **Design** — how (architecture, protocol, stack, contracts).
- **Align** — agree Define + Design are correct.
- **Take Approval** — explicit go from the spec owner. Nothing is built before this.
- **Act (one pass)** — implement the whole agreed scope in a single coherent pass.
- **Self-review** — the builder critiques its own output against the spec.
- **Sanity** — quick checks: it builds, lints, obvious correctness holds.
- **End-to-End Check** — exercise the full path against the acceptance criteria.
- **Confirm** — report results on the record; only then is the work done.

This is the canonical spec, built strictly from the recorded inputs and instructions
(see "Provenance" at the end).

---

## DEFINE

### Problem
Small / resource-constrained organizations cannot affordably and safely run RAG-based
chatbots: hyperscale cloud is costly, prompt injection is a real threat, and they lack
ML/security teams.

### Users
- **Operators** — deploy and govern the platform across an edge mesh.
- **Tenants** — deploy a domain-bounded chatbot via no-code config.
- **End users** — converse with a tenant's chatbot across channels/devices.

### Functional requirements
| ID | Requirement |
|---|---|
| F1 | Multi-tenant, no-code deployment of RAG chatbots on a distributed k3s mesh |
| F2 | Agent sessions with `load/clear context` and `load/clear memory` controls |
| F3 | `sandbox mode` isolated execution; `define scope` least-privilege per capability |
| F4 | Layered prompt-injection defence: Guard Prompts + pre-generation detector |
| F5 | PII screening + de-identification at ingestion |
| F6 | Multi-channel (web/chat/email/voice) via CloudEvents over NATS |
| F7 | Multi-device: server-side session state addressable across devices |
| F8 | Composability: every feature is a contract-bound, OCI-packaged Capability |

### Non-functional requirements
| ID | Requirement |
|---|---|
| N1 | CNCF/Kubernetes-native; k3s for edge; loop-driven (no imperative one-shots) |
| N2 | AI-native: inference/RAG/agent-loop are first-class kernel workloads |
| N3 | Distributed over an encrypted overlay; tolerant of 60–200 ms latency |
| N4 | No vendor lock; no license lock (permissive/OSI only) |
| N5 | No hidden logic; no bias; clear (inspectable) context — all testable |
| N6 | Tenant isolation prevents cross-tenant leakage by construction |
| N7 | Edge private cloud adds no net inference-latency penalty vs bare-metal |
| N8 | **Lightweight**: ships as a single static binary; runs with in-process defaults and **no mandatory external services**; resource floor ~50m CPU / 64Mi |
| N9 | **Headless**: no bundled UI; all interaction via the Edge Gate API + CloudEvents channels; any UI is an out-of-scope channel adapter |
| N10 | **Single open stack**: one architecture that scales progressively (in-process → externalized) by config, not by forking; **maximally open-source / open standards** |

### Deployment posture — one stack, progressive scale
There are **not** two stacks. The same components and seams run everywhere; scale is a
configuration choice, not a different architecture:

- **Single-binary (default, lightweight):** in-process `Decider` (OPA-shaped),
  `Screener`, and in-memory `Bus`/`Store` — the binary runs a node with no external
  dependencies. This is the "cheap edge box" mode.
- **Externalized (scale-out):** the *same* seams bind to NATS, OPA, OpenFGA, and
  PostgreSQL+pgvector. No code or architecture change — only configuration.

Every component is open-source under a permissive license (see STACK.md, SUPPLY_CHAIN.md);
openness is maximized end to end.

### Constraints / non-goals
No mandatory hyperscale or dedicated GPU; no high-stakes autonomy; no model retraining;
no new wire format where a CNCF standard fits; **no bundled UI; no proprietary
components.** (See [`SCOPE.md`](SCOPE.md).)

### Acceptance criteria
1. Tenant deploys a domain-bounded RAG chatbot via no-code on a k3s mesh.
2. Layered defence ≈100% recall / ≥99% F1 on the injection benchmark, bounded latency.
3. No net latency penalty vs bare-metal for API-served models.
4. Every active guardrail/policy is runtime-inspectable; fairness is measured in CI.

---

## DESIGN

- **Architecture** — three layers (Kernel / Runtime Core / Edge Protocol Gate):
  [`ARCHITECTURE.md`](ARCHITECTURE.md).
- **Protocol** — overlay + edge-gate contract, normative gate ordering:
  [`PROTOCOL.md`](PROTOCOL.md).
- **Composability** — capability-as-contract (provides/requires/scope/policy/authz/OCI):
  [`CAPABILITIES.md`](CAPABILITIES.md).
- **Stack** — maturity-ranked CNCF selection, license-lock rejections, the three loops:
  [`STACK.md`](STACK.md).
- **Security** — threat model + layered defences: [`THREAT_MODEL.md`](THREAT_MODEL.md).
- **Governance** — testable tenets: [`PRINCIPLES.md`](PRINCIPLES.md).

---

## ALIGN → TAKE APPROVAL (sign-off gate — current phase)

Nothing in ACT begins until this gate is explicitly cleared. Align is where the holder of
the spec and the holder of the work agree, on the record, that Define + Design are correct;
**Take Approval** is the explicit "go" that immediately precedes the single Act pass.

Alignment checklist:
- [ ] Scope (in/out, non-goals) is correct — `SCOPE.md`
- [ ] Define requirements (F*/N*) and acceptance criteria are correct
- [ ] Architecture (three layers, trust direction) is correct — `ARCHITECTURE.md`
- [ ] Protocol gate ordering is correct — `PROTOCOL.md`
- [ ] Capability-as-contract model is correct — `CAPABILITIES.md`
- [ ] Stack selections + license-lock rejections are accepted — `STACK.md`
- [ ] The four open decisions below are resolved

Rules of this phase:
- Changes happen locally and are **not pushed** until alignment is granted.
- "Act" is unblocked only by an explicit **Take Approval** go from the spec owner.

---

## ACT (one pass — gated by approval, not yet started)

On approval, the agreed scope is implemented in a **single coherent pass** (not dribbled
pushes), covering the milestones below as one body of work, then run through the closing
gates: **Self-review → Sanity → End-to-End Check → Confirm**.

| M | Milestone (within the one pass) | Exit criterion |
|---|---|---|
| M0 | **Foundation**: scope, spec, architecture, protocol, stack, principles | Docs reviewed & agreed |
| M1 | **Contracts**: machine-readable Capability schema + Kernel CRDs | `Capability`/`Tenant`/`AgentSession`/`Policy` CRDs validate |
| M2 | **Edge Gate (protocol-first)**: gate processing order, mTLS, OPA+OpenFGA stubs | A scoped request is admitted/denied per `PROTOCOL.md` |
| M3 | **Runtime Core**: RAG + inference + agent loop over CloudEvents/NATS | One tenant chatbot answers a query end-to-end |
| M4 | **Security**: Guard Prompts + detector + PII screening | Benchmark meets F4/F5 targets |
| M5 | **Deploy**: k3s + Cilium overlay + Argo CD GitOps | Multi-node mesh reconciles from Git |
| M6 | **Governance gates**: inspectability + fairness CI | N5 checks enforced in CI |

### Closing gates (after the one pass)
- **Self-review** — critique the output against this spec; list gaps and fix them.
- **Sanity** — build, lint, obvious-correctness checks pass.
- **End-to-End Check** — exercise the full path against the acceptance criteria.
- **Confirm** — report results on the record; only then is the work done.

### Open decisions (must be resolved at Take Approval)
1. Primary language/stack for kernel operator + runtime (Go is the CNCF-native default).
2. Inference backend: API-served, local small models, or both.
3. Overlay specifics (Cilium+WireGuard assumed).
4. Memory store confirmation (PostgreSQL+pgvector + Valkey assumed).

---

## Provenance (inputs & instructions this spec is built from)
Platform name *AGenNext Chat*; theory arXiv:2601.15528v1 (distributed k3s chatbot,
Guard Prompts + GenTel-Shield, PII screening); kernel-native / runtime-core /
edge-as-protocol-gate; CNCF / k8s-native / composable; no hidden logic / no bias /
clear context; OCI / OPA / OpenFGA / CloudEvents; multi-channel / multi-device /
AI-native / distributed; build a loop; pick the complete stack by maturity from the CNCF
landscape with no vendor lock and no license lock; SBOM + supply-chain check, only trusted
sources; lightweight; headless; single open stack (open source as much as possible);
framework is the input; respect laws of the land / ensure governance / compliant at all
points; build for billions; method = Define → Design → Align → Take Approval → Act (one
pass) → Self-review → Sanity → End-to-End Check → Confirm.
