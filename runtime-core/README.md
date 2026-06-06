# Runtime Core — composable execution layer

Hosts RAG, inference, and agent sessions as **contract-bound, OCI-packaged Capabilities**.
See [`../docs/ARCHITECTURE.md`](../docs/ARCHITECTURE.md) and
[`../docs/CAPABILITIES.md`](../docs/CAPABILITIES.md).

## Modules (capabilities)
- **RAG**: ingestion (PII screening + de-identification), embedding, retrieval.
- **Inference**: model-agnostic via KServe/vLLM, routed by LiteLLM (no vendor lock).
- **Agent session**: the reason→act→observe **loop**, event-driven over
  CloudEvents/NATS.

## Context & memory control surface
| Control | Effect |
|---|---|
| `load context` / `clear context` | working context window (Valkey-backed, ephemeral) |
| `load memory` / `clear memory` | persistent memory + vectors (PostgreSQL + pgvector) |
| `sandbox mode` | isolated, restricted-capability execution (default for untrusted) |
| `define scope` | least-privilege capability scope, enforced at admission |

All controls are audited — see the *clear context* and *no hidden logic* tenets in
[`../docs/PRINCIPLES.md`](../docs/PRINCIPLES.md).

> Design only — no module code yet.
