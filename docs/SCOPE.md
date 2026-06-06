# AGenNext Chat — Scope Definition

This is the founding scope document. It consolidates the platform's intent, boundaries,
and control surface. Everything built in this repository must trace back to a line here.

## 1. Mission

Enable resource-constrained organizations to deploy **secure, multi-tenant, RAG-based
agent chatbots** via a no-code workflow, on a **distributed mesh of low-cost edge
nodes** — without hyperscale cloud, dedicated GPUs, or in-house ML/security teams.

## 2. In scope

### 2.1 Platform (non-functional properties)
- **CNCF- / Kubernetes-native** control plane, targeting lightweight **k3s** for
  edge and resource-constrained deployments.
- **Distributed overlay mesh** ("the matrix") of heterogeneous nodes joined by an
  **encrypted overlay network**, tolerant of 60–200 ms inter-node latency.
- **Composable**: every layer and module is independently deployable and swappable.
- **Multi-tenant** with container-based isolation and per-tenant data-access controls;
  designed to limit the blast radius of a compromised component.

### 2.2 Architecture layers
- **Kernel (kernel-native control plane):** scheduling, tenant lifecycle, isolation,
  policy — expressed as the k8s/k3s reconciliation loop.
- **Runtime Core:** RAG pipeline, LLM inference, and **agent sessions**, including the
  context/memory control surface and per-session sandboxing.
- **Edge Protocol Gate:** ingress, authN/Z, load balancing, overlay termination, and
  first-layer security filtering. The only boundary that speaks the overlay protocol.

### 2.3 Agent / session control surface
These are **platform features**, exposed to operators and (scoped) to tenants:
- `load context` / `clear context` — manage a session's working context window.
- `load memory` / `clear memory` — manage a session's persistent memory store.
- `sandbox mode` — run an agent/tenant in an isolated, restricted-capability execution
  environment (default for untrusted workloads).
- `define scope` — bound an agent's capabilities, tools, and data access (least
  privilege).

### 2.4 Security & privacy (from arXiv:2601.15528v1)
- **Layered prompt-injection defence:** system-level **Guard Prompts** + a pre-generation
  **GenTel-Shield**-style detector applied to user queries *and* retrieved content.
- **PII screening & de-identification** at document-ingestion time.
- **Tenant isolation** sufficient to prevent cross-tenant data leakage.

### 2.5 Governance tenets (testable)
- **No hidden logic** — prompts, guardrails, routing, and policy are inspectable.
- **No bias** — fairness is an explicit, measured property; see `docs/PRINCIPLES.md`.
- **Clear context** — session context/memory are observable and operator-controllable.

## 3. Out of scope (explicit non-goals)

- **Hyperscale-cloud dependence** or mandatory dedicated-GPU infrastructure.
- **High-stakes autonomous decision-making** or unrestricted generative tasks — the
  platform is for *constrained, domain-bounded* question answering under operator policy.
- **Model training / fine-tuning** — defences are model-agnostic and require no retraining.
- **Building a new wire protocol from scratch** where a CNCF standard fits (prefer
  existing overlay/mTLS/OCI/OpenTelemetry building blocks).
- **A consumer end-user identity system** — tenancy and operator identity first.

## 4. Success criteria

1. A tenant can deploy a domain-bounded RAG chatbot via no-code config on a k3s mesh.
2. Prompt-injection defences measurably reduce attack success (target: layered config
   ≈100% recall, ≥99% F1, per the case study) with bounded latency overhead.
3. The edge private cloud introduces **no net inference-latency penalty** versus
   bare-metal for API-served LLM workloads.
4. Every active guardrail and policy is inspectable at runtime (no hidden logic).

## 5. Open decisions (to confirm before implementation)

- **Primary language/stack** for kernel operator and runtime modules (Go is the
  CNCF-native default; not yet chosen).
- **Inference backend(s)**: API-served models, local small models (e.g. 3B-class), or both.
- **Overlay technology** for the encrypted mesh (e.g. WireGuard-based vs. service mesh).
- **Memory store** backing `load/clear memory`.

These are tracked here rather than guessed; they gate the move from scaffold to code.
