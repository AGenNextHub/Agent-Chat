# AGenNext Chat

A **composable, Kubernetes-native, multi-tenant agent-chat platform** for deploying
secure RAG-based chatbots on a distributed mesh of low-cost edge nodes — without
hyperscale cloud, dedicated GPUs, or in-house ML engineering.

Grounded in *"Securing LLM-as-a-Service for Small Businesses"* (Xie et al., RMIT
University, AISC 2026, arXiv:2601.15528v1) and generalized into a layered, CNCF-aligned,
AI-native platform. See [`docs/REFERENCES.md`](docs/REFERENCES.md).

> **Status — v0.1.x (the loop spine):** the agent-loop core is implemented in **pure Go
> standard library (zero third-party dependencies)** with tests, an end-to-end demo, and
> full GitOps/OSS/supply-chain scaffolding. Control-plane (CRDs/operator) and production
> bindings (NATS, OPA, OpenFGA, KServe, PostgreSQL) are next — see
> [`docs/ROADMAP.md`](docs/ROADMAP.md).

## The model

```
        external clients + federating edge nodes (encrypted overlay = "the matrix")
                                  │
                ┌─────────────────▼─────────────────┐
                │  EDGE — Protocol Gate (pkg/edge)   │  authN/Z · scope · L1 screen
                └─────────────────┬─────────────────┘
                ┌─────────────────▼─────────────────┐
                │  RUNTIME CORE — the loop (pkg/loop)│  reason→guard→act→screen→persist
                └─────────────────┬─────────────────┘
                ┌─────────────────▼─────────────────┐
                │  KERNEL — k3s/k8s control loop     │  scheduling · isolation · policy
                └────────────────────────────────────┘
```

The conceptual spine is the **Trinity**: *Context → Capability → Contract*. See
[`docs/CONCEPTS.md`](docs/CONCEPTS.md).

## Quickstart

```bash
go run ./cmd/agennextd     # run one turn end-to-end; prints the inspectable trace
```

The demo admits a chat event through the Edge Gate, runs the bounded agent loop
(retrieve → guard → screen → answer), and emits a JSON trace where every step is
attributable — no hidden logic.

## Verify (make verification accessible)

```bash
make check     # tidy + go vet + golangci-lint + go test -race   (mirrors CI)
make cover     # coverage summary
make run       # the end-to-end demo
```

CI ([`.github/workflows/ci.yaml`](.github/workflows/ci.yaml)) runs the same checks plus
`govulncheck`, SBOM generation/scan, and `helm lint`.

## What's implemented

| Package | Role |
|---|---|
| `pkg/event` | CloudEvents 1.0 envelope (content-type agnostic — all I/O types) + bus |
| `pkg/capability` | capability-as-contract: `Contract`, `Scope`, `Registry` |
| `pkg/store` | working `Context` + durable `Memory` (load/clear) |
| `pkg/guard` | injection `Screener`, Guard Prompts, in-process policy `Decider` |
| `pkg/loop` | the agent loop: bounded, guard-before-act, screen-on-output, idempotent |
| `pkg/edge` | the Protocol Gate: authn → authz → scope → L1 screen |
| `cmd/agennextd` | single-node end-to-end demo |

## Documentation

| Doc | Purpose |
|---|---|
| [`docs/SPEC.md`](docs/SPEC.md) | the canonical spec + delivery pipeline + milestones |
| [`docs/CONCEPTS.md`](docs/CONCEPTS.md) | the Trinity, primitives, agent-at-all-gates |
| [`docs/LOOP.md`](docs/LOOP.md) | the agent loop in detail |
| [`docs/ARCHITECTURE.md`](docs/ARCHITECTURE.md) | the three-layer model |
| [`docs/PROTOCOL.md`](docs/PROTOCOL.md) | overlay + edge-gate contract |
| [`docs/CAPABILITIES.md`](docs/CAPABILITIES.md) | capability-as-contract model |
| [`docs/STACK.md`](docs/STACK.md) | maturity-ranked CNCF stack (no vendor/license lock) |
| [`docs/SUPPLY_CHAIN.md`](docs/SUPPLY_CHAIN.md) | SBOM + supply-chain policy |
| [`docs/THREAT_MODEL.md`](docs/THREAT_MODEL.md) | threats + layered defences |
| [`docs/PRINCIPLES.md`](docs/PRINCIPLES.md) | testable governance tenets |
| [`docs/SCOPE.md`](docs/SCOPE.md) · [`docs/ROADMAP.md`](docs/ROADMAP.md) | scope & roadmap |
| [`docs/adr/`](docs/adr/) | architecture decision records |

## Contributing & security

See [`CONTRIBUTING.md`](CONTRIBUTING.md), [`CODE_OF_CONDUCT.md`](CODE_OF_CONDUCT.md), and
[`SECURITY.md`](SECURITY.md).

## License

Apache-2.0 (CNCF-standard). See [`LICENSE`](LICENSE) and [`NOTICE`](NOTICE).
