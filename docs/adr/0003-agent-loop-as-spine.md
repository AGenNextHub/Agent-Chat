# 3. The agent loop is the spine; security lives inside the loop

Date: 2026-06-06

## Status
Accepted

## Context
AGenNext Chat is AI-native: the conversation/turn is the core unit. Prompt injection can
arrive directly (user input) or indirectly (retrieved/tool output). A gate-only defence
misses indirect injection, which enters *after* admission.

## Decision
We implement the agent loop (`pkg/loop`) as a bounded reason→act cycle and place the
platform's invariants *inside* the iteration:
- **Guard before Act** — every proposed action is checked against its capability contract
  scope and an in-process policy decision before it runs.
- **Screen tool output** — untrusted ACT output is screened for injection before it can
  re-enter REASON (closes threat T2).
- **Transactional turn** — context/memory deltas are buffered and committed once.
- **Bounded** — max iterations / token / wall-clock budgets guarantee termination.
- **Idempotent intake** — at-least-once replay returns the prior result.

## Consequences
- Security is not a wrapper around the loop; it is a property of each iteration.
- The same guard primitives are reusable at the Edge Gate ("agent at all gates").
- Reasoner and Invoker are interfaces, so models and tools are swappable without touching
  the invariants.
