# AGenNext Chat — Specification (Define → Design → Align → Act)

This is the canonical spec. It is organized by the **Define → Design → Align → Act**
method: *Define* what must be true, *Design* how, ***Align* on sign-off before acting**,
then *Act* in ordered milestones. No phase proceeds until the prior one is agreed. It is
built strictly from the recorded inputs and instructions (see "Provenance" at the end).

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

### Constraints / non-goals
No mandatory hyperscale or dedicated GPU; no high-stakes autonomy; no model retraining;
no new wire format where a CNCF standard fits. (See [`SCOPE.md`](SCOPE.md).)

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

## ALIGN (sign-off gate — current phase)

Nothing in ACT begins until this gate is explicitly cleared. Align is where the holder of
the spec and the holder of the work agree, on the record, that Define + Design are correct.

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
- "Act" is unblocked only by an explicit go from the spec owner.

---

## ACT (proposed milestones — gated by ALIGN, not yet started)

| M | Milestone | Exit criterion |
|---|---|---|
| M0 | **Foundation** (this PR): scope, spec, architecture, protocol, stack, principles | Docs reviewed & agreed |
| M1 | **Contracts**: machine-readable Capability schema + Kernel CRDs | `Capability`/`Tenant`/`AgentSession`/`Policy` CRDs validate |
| M2 | **Edge Gate (protocol-first)**: gate processing order, mTLS, OPA+OpenFGA stubs | A scoped request is admitted/denied per `PROTOCOL.md` |
| M3 | **Runtime Core**: RAG + inference + agent loop over CloudEvents/NATS | One tenant chatbot answers a query end-to-end |
| M4 | **Security**: Guard Prompts + detector + PII screening | Benchmark meets F4/F5 targets |
| M5 | **Deploy**: k3s + Cilium overlay + Argo CD GitOps | Multi-node mesh reconciles from Git |
| M6 | **Governance gates**: inspectability + fairness CI | N5 checks enforced in CI |

### Open decisions (gate the move from M0→M1)
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
landscape with no vendor lock and no license lock; method = Define → Design → Align → Act.
