# 4. Protobuf-first contracts (common language / communication layer)

Date: 2026-06-06

## Status
Proposed — the IDL is adopted now; the code-generation step and its runtime
dependency are **gated** (supply-chain check + human release gate).

## Context
"Protocol first," "service is a contract," "common language / universal framework / common
communication layer," and "fine-grained access control" all point at one decision: a single,
language-neutral **contract IDL**. Protocol Buffers (+ gRPC) is the natural fit — it is the
canonical contract for the headless API and for the capability/event types, generated into
every language a channel or service is written in.

This collides with ADR-0002 (pure-stdlib, zero third-party dependencies): generated Go
requires `google.golang.org/protobuf`, and the gRPC service requires `google.golang.org/grpc`.
Both are trusted (Google) and permissively licensed (BSD-3 / Apache-2.0) — they pass
*only-trusted-sources* and *no-license-lock* — but they end "zero third-party deps."

## Decision
1. **Adopt `.proto` as the canonical contract language now.** `proto/agennext/v1/` holds the
   source of truth: `Contract` (capability), `Event` (CloudEvents), and the `Chat` service
   (the headless gate API). A `.proto` is plain text and adds no dependency.
2. **Defer code generation and the runtime dependency.** `buf.gen.yaml` is provided but not
   run. Adopting `protobuf`/`grpc` requires: an SBOM + supply-chain check (per
   SUPPLY_CHAIN.md), digest pinning, and explicit **release-gate** approval.
3. **Until codegen lands,** the hand-written Go types in `pkg/` are the implementation and
   MUST stay in sync with the `.proto`. The proto lint/breaking checks run in CI via buf.
4. The headless `Chat.Send` request **carries no scope** — authority is derived at the gate
   (fine-grained access control via grants ∩ contract), reinforcing ADR's least-privilege.

## Consequences
- One contract, many languages: channels/services interoperate over a common communication
  layer without re-specifying types.
- The zero-dep property is preserved for now and will be traded *consciously*, through a
  gate, when codegen is adopted — not silently.
- A drift risk exists between `.proto` and hand-written Go until codegen replaces the latter;
  CI proto checks and review mitigate it.
