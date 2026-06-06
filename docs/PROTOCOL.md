# The Protocol — overlay & edge-gate contract ("the matrix")

Protocol first. "The matrix" is the **encrypted overlay mesh of heterogeneous edge
nodes**; "the protocol" is the contract those nodes and clients speak, and the **Edge
Protocol Gate** is the only place it terminates.

## Goals

- Terminate untrusted external traffic at a single, inspectable boundary.
- Carry tenant identity and capability scope with every request.
- Survive 60–200 ms inter-node latency and heterogeneous, low-cost nodes.
- Reuse CNCF standards rather than invent a new wire format where one fits.

## Planes

| Plane | Concern | Building block |
|---|---|---|
| **Transport** | confidentiality, node identity | encrypted overlay + **mTLS** |
| **AuthN** | who is calling | mTLS identity / token |
| **AuthZ** | may they do this | **OpenFGA** relations + **OPA** policy |
| **Admission** | is the capability contract satisfied | Kernel (see `CAPABILITIES.md`) |
| **Application** | the chat/RAG request itself | request envelope below |

## Request envelope (conceptual)

Every admitted request carries:

```
overlay frame  ──▶  mTLS(node identity)
  └─ envelope
       tenant:      <tenant-id>           # required; scopes everything downstream
       principal:   <caller identity>     # for OpenFGA checks
       capability:  <name@version>        # what is being invoked
       scope:       <declared least-priv> # must be ⊆ capability contract scope
       payload:     <query / documents>   # untrusted content
       trace:       <OpenTelemetry ctx>
```

## Gate processing order (normative)

1. **Terminate overlay + verify mTLS** node/peer identity.
2. **Authenticate** the principal.
3. **Authorize** — OpenFGA relation check, gated by OPA policy. Deny ⇒ stop.
4. **Derive effective scope** = principal's grant ∩ capability contract scope.
   Scope is never taken from the caller; an empty result ⇒ stop.
5. **Security filter (L1):** inject Guard Prompts; run the pre-generation injection
   detector over `payload` (query *and* any retrieved content). Flagged ⇒ refuse.
6. **Forward** the now-scoped, screened request to the Runtime Core.

A request that fails any step **never reaches inference**. This ordering is part of the
threat model — see [`THREAT_MODEL.md`](THREAT_MODEL.md).

## Federation note

Multiple edge gates may peer over the overlay to pool heterogeneous nodes. Peers
authenticate by mTLS and are subject to the same admission and authZ planes — a peer is
not more trusted than a client.

> This document defines the **contract**, not a frozen byte layout. The concrete encoding
> (e.g. gRPC/HTTP + protobuf) is an open decision tracked in `SCOPE.md §5`.
