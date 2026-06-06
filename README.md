# AGenNext Chat

A **composable, Kubernetes-native, multi-tenant agent-chat platform** for deploying
secure RAG-based chatbots on a distributed mesh of low-cost edge nodes — without
hyperscale cloud, dedicated GPUs, or in-house ML engineering.

AGenNext Chat is grounded in the industry case study *"Securing LLM-as-a-Service for
Small Businesses: An Industry Case Study of a Distributed Chatbot Deployment Platform"*
(Xie, Li, Fu, Gao, Xu, Han — RMIT University, AISC 2026, arXiv:2601.15528v1) and
generalizes its findings into a layered, CNCF-aligned platform.

> **Status:** founding scaffold. This repository currently defines the **scope,
> architecture, protocol, threat model, and governance principles**. Component
> directories contain design READMEs, not yet runnable services.

---

## The model in one picture

```
        external clients + federating edge nodes (the "matrix" overlay)
                                  │
                ┌─────────────────▼─────────────────┐
                │  EDGE  —  the Protocol Gate        │  ingress · authN/Z · load balance
                │  (Client + overlay termination)    │  TLS · guard prompts · injection filter
                └─────────────────┬─────────────────┘
                ┌─────────────────▼─────────────────┐
                │  RUNTIME CORE                      │  RAG · inference · agent sessions
                │  (composable modules)              │  context + memory controls · sandbox
                └─────────────────┬─────────────────┘
                ┌─────────────────▼─────────────────┐
                │  KERNEL  —  kernel-native control  │  k3s/k8s · scheduling · tenant isolation
                │  (control loop = reconciliation)   │  lifecycle · policy · blast-radius limits
                └────────────────────────────────────┘
```

- **Kernel (kernel-native):** the Kubernetes (k3s) control plane *is* the platform
  kernel. Scheduling, per-tenant isolation, lifecycle, and policy are expressed as the
  cluster's reconciliation loop.
- **Runtime Core:** composable execution layer — RAG pipeline, LLM inference, agent
  sessions, and the **context / memory** control surface. Every tenant runs sandboxed.
- **Edge Protocol Gate:** the only place external traffic and the encrypted overlay
  ("the matrix") terminate — ingress, auth, load balancing, and the first layer of
  prompt-injection defence.

## Governance tenets (non-negotiable)

- **No hidden logic** — every guard prompt, routing rule, and policy is inspectable.
- **No bias** — fairness is an explicit, testable property of model behaviour.
- **Clear (inspectable) context** — session context and memory are observable and
  operator-controllable (`load`/`clear`).
- **Privacy by default** — PII screening at ingestion; per-tenant data isolation.
- **Composable** — layers and modules are independently deployable and swappable.

## Documentation

| Doc | Purpose |
|---|---|
| [`docs/SCOPE.md`](docs/SCOPE.md) | What is and is **not** in scope (the founding definition) |
| [`docs/ARCHITECTURE.md`](docs/ARCHITECTURE.md) | The three-layer model in detail |
| [`docs/PROTOCOL.md`](docs/PROTOCOL.md) | The overlay/edge-gate protocol ("the matrix") |
| [`docs/THREAT_MODEL.md`](docs/THREAT_MODEL.md) | Prompt injection & multi-tenant risks + mitigations |
| [`docs/PRINCIPLES.md`](docs/PRINCIPLES.md) | Governance tenets, expanded and made testable |

## Repository layout

```
kernel/        control-plane / operator design (k3s-native)
runtime-core/  RAG + inference + agent-session modules
edge-gate/     protocol gate: ingress, auth, security filtering
deploy/        k3s / Helm / overlay deployment scaffolding
docs/          scope, architecture, protocol, threat model, principles
```

## Attribution

The security design (layered Guard Prompts + GenTel-Shield detection), the distributed
k3s edge-cloud topology, and the multi-tenant isolation model derive from
arXiv:2601.15528v1. See [`docs/THREAT_MODEL.md`](docs/THREAT_MODEL.md) for the mapping.

## License

Apache-2.0 (CNCF-standard). See [`LICENSE`](LICENSE).
