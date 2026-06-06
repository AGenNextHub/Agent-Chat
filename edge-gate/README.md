# Edge Protocol Gate

The only boundary that terminates the encrypted overlay ("the matrix") and faces clients.
See [`../docs/PROTOCOL.md`](../docs/PROTOCOL.md) and
[`../docs/THREAT_MODEL.md`](../docs/THREAT_MODEL.md).

## Processing order (normative)
1. Terminate overlay + verify mTLS (Cilium/WireGuard + SPIRE identity).
2. Authenticate principal.
3. Authorize — OpenFGA relations gated by OPA policy.
4. Validate capability contract scope ⊇ requested scope.
5. Security filter (L1): inject Guard Prompts + run pre-generation injection detector
   over query *and* retrieved content.
6. Forward scoped, screened request to the Runtime Core.

A request failing any step **never reaches inference**.

## Build
- Envoy + Gateway API for ingress (Apache-2.0).
- Multi-channel adapters publish/consume CloudEvents over NATS.

> Design only — no gate code yet.
