# The Loop — the spine of AGenNext Chat

Defined first, on purpose. AGenNext Chat is AI-native, which means the **agent loop** is
the product; the Kernel, Edge Gate, and stack exist to serve it. Everything else in the
spec wraps around what is defined here.

There are three loops in the platform; only the **agent loop** is the data-plane spine.
The other two are control-plane and exist to keep the spine running.

| Loop | Plane | Owner | Drives |
|---|---|---|---|
| **Agent loop** | data | Runtime Core | a single conversation/turn to a result |
| Kernel control loop | control | Kernel (operator) | desired state of CRDs |
| GitOps loop | control | Argo CD | cluster ⇆ Git |

This document defines the **agent loop**.

## Definition

The agent loop is a **bounded, event-driven, reason→act cycle** executed per
`AgentSession`. It consumes one admitted inbound event and runs until it produces a
response or hits a budget bound.

```
        inbound CloudEvent (admitted by Edge Gate)
                     │
            ┌────────▼─────────┐
            │ 1. PERCEIVE      │  parse event: tenant, principal, capability, scope, payload
            └────────┬─────────┘
            ┌────────▼─────────┐
            │ 2. LOAD          │  hydrate context (Valkey) + memory (pgvector retrieval)
            └────────┬─────────┘
            ┌────────▼─────────┐
            │ 3. REASON        │  model proposes next action: answer | invoke capability
            └────────┬─────────┘
            ┌────────▼─────────┐
            │ 4. GUARD         │  re-check proposed action vs contract scope + OPA policy
            └────────┬─────────┘   (deny ⇒ refuse/repair, never silently proceed)
            ┌────────▼─────────┐
            │ 5. ACT           │  run the capability in its sandbox, within scope
            └────────┬─────────┘
            ┌────────▼─────────┐
            │ 6. OBSERVE       │  capture result + emit OpenTelemetry span
            └────────┬─────────┘
                     │  more steps needed and budget remains?
             ┌───────┴────────┐
          yes│ (back to 3)    │no
             ▲                ▼
            ┌──────────────────────┐
            │ 7. PERSIST           │  update context + memory (honor load/clear controls)
            └────────┬─────────────┘
            ┌────────▼─────────────┐
            │ 8. EMIT              │  outbound CloudEvent → channel/device
            └──────────────────────┘
```

## States

`PERCEIVE → LOAD → REASON → GUARD → (ACT → OBSERVE → REASON)* → PERSIST → EMIT`

The starred segment is the iteration. The loop **never** transitions ACT before GUARD.

## Inputs / outputs

- **In:** one admitted `CloudEvent` (the Edge Gate has already done authN/Z, scope
  validation, and the L1 injection filter — the loop trusts only admitted events).
- **Out:** one or more outbound `CloudEvent`s carrying the response, plus telemetry and a
  persisted context/memory delta.

## Invariants (non-negotiable)

1. **Bounded** — every loop has a max-iteration count and a token/time budget. No
   unbounded loops. Exhausting the budget yields a graceful, explained stop.
2. **Guard before Act** — security is *inside* the loop, not a wrapper. Each proposed
   action is checked against its capability contract scope and OPA policy before it runs.
3. **Attributable** — every ACT is traceable to a capability contract (no hidden logic).
4. **Stateless workers** — loop state lives in context/memory stores, not the worker. This
   is what makes the platform multi-device and distributable across the edge mesh.
5. **Idempotent intake** — events may be redelivered (at-least-once); a turn must tolerate
   replay without duplicate side effects.
6. **Inspectable** — at any point the active guards, the chosen actions, and the loaded
   context are retrievable for audit (`clear context` / `no hidden logic`).

## Budgets & termination

| Bound | Default (proposed) | Effect on hit |
|---|---|---|
| max iterations | (open decision) | stop, return best answer + reason |
| token budget | (open decision) | stop, return partial + reason |
| wall-clock | (open decision) | stop, return timeout notice |

## Control surface mapping

| Control | Where it acts in the loop |
|---|---|
| `load context` / `clear context` | step 2 (LOAD) and step 7 (PERSIST) |
| `load memory` / `clear memory` | step 2 (LOAD) and step 7 (PERSIST) |
| `sandbox mode` | step 5 (ACT) — execution profile |
| `define scope` | step 4 (GUARD) — the scope checked against |

## Relationship to the control-plane loops

- The **Kernel control loop** schedules/sandboxes the workers that run the agent loop and
  reconciles `AgentSession` lifecycle.
- The **GitOps loop** delivers the capability artifacts (OCI) the agent loop invokes.

Neither runs a conversation; both exist so the agent loop can run safely and everywhere.

> Open decisions to resolve at Take Approval: the three budget defaults above, and whether
> REASON is a single model call or itself a planner/executor split.
