# Concepts — the conceptual model of AGenNext Chat

This document fixes the vocabulary. Every other doc and every package name traces
back to a term defined here.

## The Trinity

Three irreducible elements recur at every layer of the platform:

> **Context → Capability → Contract**

- **Context** — the *demand*. What the situation requires (a turn, a query, a state).
- **Capability** — the *delivery*. What fulfils the demand. *Context demands capability;
  capability must be delivered.*
- **Contract** — the *binding*. What bounds and governs the delivery. *A capability is a
  contract; a service is a contract.*

The Trinity is not an analogy bolted on top — it is the literal control flow of the agent
loop: a turn's **context** produces a demand, a **contract** gates it, the loop **delivers**
the capability. It also projects onto structure:

| Trinity element | Layer | Agent face |
|---|---|---|
| Context (sense the demand) | Edge Protocol Gate | **Interface** |
| Capability (deliver) | Runtime Core | **Orchestrator** |
| Contract (govern) | Kernel | **Operator** |

## The platform is the agent

The same **agent** abstraction recurs at every layer rather than living in one place:

- **Orchestrator** — the agent loop orchestrating a turn (Runtime Core).
- **Operator** — the agent as the Kubernetes operator reconciling desired state (Kernel).
- **Interface** — the agent as the conversational surface and the gate (Edge).
- **Agent at all gates / edges** — every boundary hosts an agent that *reasons* about
  admission; a gate is not a static filter. The `loop` engine is reusable as a gate's
  decision engine, which is why `edge.Gate` runs the same guard primitives the loop does.
- **Platform is the edge / the surface** — what every stakeholder *touches* is the boundary
  (the edge/surface); the cores live edge-free behind it. The platform's outward identity is
  all edge, even though its cores have none. Inside: core. Outside: edge.

## The loop is the blueprint

The loop is not merely the data-plane spine — it is the **blueprint every layer
instantiates**. The same shape recurs at every altitude:

> perceive → decide-under-contract → act → observe → repeat, **bounded**

| Instance | Altitude | "Perceive → … → act" |
|---|---|---|
| **Agent loop** | a turn | event → reason → guard → act → screen → observe |
| **Kernel control loop** | desired state | watch → diff → admit (contract) → reconcile |
| **Gate as agent** | a boundary | request → authn/z → scope → screen → admit |
| **GitOps loop** | cluster ⇆ Git | observe Git → diff → apply → observe |

This is why *agent at all gates* and *the platform is the agent* hold: each layer is the
same loop, governed by a contract, bounded so it always terminates. Implement the loop
once (`pkg/loop`) and the pattern — and its invariants — are reusable everywhere.

## Composition hierarchy

> **capability ⊂ service ⊂ solution**

- **Capability** — the atomic, contract-bound unit (`pkg/capability`). OCI-packaged.
- **Service** — a composition of capabilities exposed under one contract.
  *Service is a contract.*
- **Solution** — what the tenant actually deploys. *The solution is the service.*

## The core

- **Core is the contract; the contract is the core.** The capability core
  (`pkg/capability`) is the **canonical** definition; every other core is an expression of
  it. The core *defines* capability — not the transport, not the edge.
- **The core has no edges.** A core package never imports the boundary (`pkg/edge`); the
  **edge composes around the core**. This is enforced (CI checks no core imports edge).
- **The core expands in the loop.** The core grows only by adding capabilities the loop can
  invoke — never by absorbing the boundary.
- **The core delivers capability through tangibles.** Capability reaches the real world only
  via the loop's ACT step (the `Invoker`) — real tools, OCI artifacts, and devices. Each
  delivery is bounded by the contract's scope and sandbox: delivery is where the contract
  meets the world, never unbounded.

| Core | Package | Notes |
|---|---|---|
| Capability / canonical | `pkg/capability` | the contract — the canonical core |
| Runtime / Agent | `pkg/loop` | the agent loop (the spine) |
| Chat | `pkg/chat` | runs admitted turns; edge-free |
| Kernel | `pkg/kernel` | admits/reconciles contracts (control loop) |
| Cloud / Kubernetes | *(gated)* | operator/client-go — needs the release-gated deps |

The edge (`pkg/edge`) and transports (gRPC/HTTP) wrap these cores from outside.

## Primitives

The smallest building blocks, each a Go type in `pkg/`:

| Primitive | Package | Meaning |
|---|---|---|
| **Event** | `event` | CloudEvents 1.0 envelope; content-type agnostic (all I/O types) |
| **Contract** | `capability` | the complete, inspectable description of a capability |
| **Scope** | `capability` | least-privilege authority (tenants/data/tools) |
| **Context** | `store` | a session's working window |
| **Memory** | `store` | a session's durable facts |
| **Screener / Verdict** | `guard` | injection detection over untrusted text |
| **Decider / Decision** | `guard` | in-process (OPA-shaped) policy decision |
| **Action / Observation** | `loop` | a proposed step and its sanitized result |
| **AdmittedEvent** | `loop` | an event the gate has cleared |
| **Result / TraceEntry** | `loop` | the turn outcome and its inspectable trace |

## Tenets (testable — see PRINCIPLES.md)

- **No hidden logic** · **No bias** · **Clear context** · **Operate with clarity** —
  every decision is attributable to a contract and visible in the turn trace.
- **Secure contract · protect privacy · compliant at all points · respect laws of the
  land** — PII screened at ingestion, tenant isolation by construction, per-jurisdiction
  compliance and governance gates.
- **Build for billions** — stateless workers, idempotent at-least-once intake, distributed
  edge mesh; turn state lives in stores, not workers, so the loop scales horizontally.
- **No vendor lock · no license lock · only trusted sources · SBOM-gated** — see
  STACK.md and SUPPLY_CHAIN.md.
