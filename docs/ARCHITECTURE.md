# AGenNext Chat — Architecture

Foundation-first, protocol-first. Three layers, one direction of trust: external traffic
is untrusted until the **Edge Protocol Gate** admits it; the **Kernel** is the root of
trust; the **Runtime Core** executes only what the Kernel scheduled and the Gate admitted.

```
        external clients + federating edge nodes (encrypted overlay = "the matrix")
                                  │
                ┌─────────────────▼─────────────────┐
   admit/deny   │  EDGE — Protocol Gate              │  overlay termination, mTLS
   ◀────────────│  ingress · authN/Z · LB · filter   │  OpenFGA authZ · OPA policy
                └─────────────────┬─────────────────┘  Guard Prompts (L1) + detector
                                  │ admitted, scoped request
                ┌─────────────────▼─────────────────┐
   schedule     │  RUNTIME CORE                      │  RAG · inference · agent sessions
   ◀────────────│  composable Capability modules     │  context/memory · sandbox mode
                └─────────────────┬─────────────────┘  (each module = OCI artifact)
                                  │ reconcile
                ┌─────────────────▼─────────────────┐
   root of      │  KERNEL — kernel-native control    │  k3s/k8s control loop
   trust        │  scheduling · isolation · policy   │  tenant lifecycle · CRDs
                └────────────────────────────────────┘
```

## Layer 1 — Kernel (kernel-native control plane)

The Kubernetes / k3s control plane **is** the kernel. There is no separate orchestrator.

Responsibilities:
- **Scheduling & placement** of Capability modules across heterogeneous edge nodes,
  accounting for overlay latency (60–200 ms) and node resources (CPU-only included).
- **Tenant lifecycle & isolation** — each tenant is a namespace boundary with
  container-based isolation and per-tenant data-access controls.
- **Policy root** — admission control and reconciliation enforce platform invariants.
- **Blast-radius containment** — a compromised module cannot reach another tenant's data.

Expressed as **CRDs + an operator**: `Tenant`, `Capability`, `AgentSession`, `Policy`.
The reconciliation loop is the kernel's scheduler.

## Layer 2 — Runtime Core

The composable execution layer. Everything here is a **Capability module** (see
[`CAPABILITIES.md`](CAPABILITIES.md)) packaged as an **OCI artifact** and bound by a
**contract**.

Core modules:
- **RAG pipeline** — ingestion (with PII screening/de-identification), embedding,
  retrieval.
- **Inference** — model-agnostic; API-served and/or local small models.
- **Agent sessions** — the unit a user converses with. Owns the **context** and
  **memory** control surface:
  - `load context` / `clear context` — working context window.
  - `load memory` / `clear memory` — persistent per-session memory store.
- **Sandbox mode** — default isolated, restricted-capability execution for untrusted
  workloads.

## Layer 3 — Edge Protocol Gate

The only component that terminates the encrypted overlay ("the matrix") and faces
clients. Nothing speaks raw protocol to the Runtime Core.

Responsibilities:
- **Overlay termination & mTLS**, signature verification of federating nodes.
- **Authentication** of callers and **authorization** via **OpenFGA** (relationship-based,
  fine-grained) gated by **OPA** policy.
- **Load balancing** across runtime replicas.
- **Security filtering (L1)** — Guard Prompts injected and the pre-generation
  injection detector applied before any request reaches inference.

See [`PROTOCOL.md`](PROTOCOL.md) for the gate's wire contract.

## Cross-cutting

| Concern | Building block |
|---|---|
| Packaging & distribution | **OCI** artifacts |
| Policy | **OPA** (Open Policy Agent) |
| Authorization | **OpenFGA** |
| Orchestration | **Kubernetes / k3s** |
| Observability | OpenTelemetry |
| Transport security | mTLS over the encrypted overlay |

## Trust & data-flow invariants

1. A request reaches inference **only** after Edge admission + Guard Prompt + detector.
2. A Capability runs **only** with the scope its contract declares (least privilege).
3. Cross-tenant data access is structurally impossible, not merely policy-discouraged.
4. Every active guardrail/policy is inspectable at runtime — **no hidden logic**.
