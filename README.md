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

## The contract of this repo — Agent Language

Plain words: the Agent Language has one format — a **command set**, where each command is
**structured instructions**. The format is general: it can define anything. A person reads
the instructions; a machine runs the commands — same meaning on both sides. This is the
contract.

A word may carry **multiple definitions**, but every meaning must be **clear**,
**executable** (a machine can run it), **explainable** (a person can understand it), and
**defendable** (you can justify it) — each grounded in an **existing source of authority**
(a dictionary, a standard, a law). Nothing is invented; nothing is left to interpretation.

**The context — the Unboxd Dictionary.** Every meaning comes from a context: the **Unboxd
Dictionary**. This contract is **owned by its context** — change a meaning there and every
contract bound to it changes; the context's maintainer is the authority. The dictionary
must be:
- **public** — anyone can read it
- **versioned** — every change is a new version; nothing is silently edited
- **signed** — each version is signed, so you know who stands behind it
- **community-maintained** — kept by a community of open-source contributors

**Execution.** Every command and every code push runs **under this contract** — validated
against the Unboxd Dictionary. A command whose words are not in the context, or whose
meaning cannot be defended, does not execute. The contract is the gate.

**Run the chat** — `go run ./cmd/agennextd`
- starts the chat service on port 8080
- checks who is calling and what they may do (the gate)
- runs the agent one bounded step at a time
- returns an answer with a full trace of every step — nothing hidden

**Build a box** (a sealed unit of content) — `boxctl build -t <type> -f <file> -o my.box`
- wraps your file into a box
- gives it a fingerprint (digest) computed from its content
- writes the box file

**Check a box** — `boxctl verify <sealed.box>`
- confirms the content was not changed
- confirms it was signed by the stated person

**Resolve a box's parts** — `boxctl resolve -dir <folder> <fingerprint>`
- walks every part the box depends on
- confirms each part is present
- refuses if a part is missing or the parts loop

**Sign a box** (a person, never the agent) — `cosign sign <registry>/<name>:<tag>`
- you sign it with your own identity
- nobody else can claim it; the agent never signs for you

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
go run ./cmd/agennextd                 # serve the headless API on :8080
curl localhost:8080/healthz            # {"status":"ok"}
curl -XPOST localhost:8080/v1/chat -H 'content-type: application/json' \
  -d '{"id":"1","session":"s1","tenant":"acme","principal":"u1","capability":"rag.retrieve","message":"return policy?"}'

go run ./cmd/agennextd -demo           # or: run one turn and print the inspectable trace
```

The daemon admits a chat event through the Edge Gate, runs the bounded agent loop
(retrieve → guard → screen → answer), and returns a JSON trace where every step is
attributable — no hidden logic. The request carries **no scope**: authority is derived at
the gate. Admission failures map to safe HTTP codes (401/403/422) that bleed no detail.

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
