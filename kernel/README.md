# Kernel — kernel-native control plane

The Kubernetes / k3s control plane **is** the platform kernel. See
[`../docs/ARCHITECTURE.md`](../docs/ARCHITECTURE.md) and
[`../docs/STACK.md`](../docs/STACK.md).

## Responsibilities
- Scheduling & placement across heterogeneous edge nodes.
- Tenant lifecycle and container-based isolation (blast-radius containment).
- Capability admission (contract + OPA policy + OpenFGA relations).
- The **control loop** that reconciles desired state.

## Planned CRDs
| CRD | Purpose |
|---|---|
| `Tenant` | isolation boundary, data scope, quotas |
| `Capability` | contract-bound, OCI-referenced module (see `CAPABILITIES.md`) |
| `AgentSession` | a conversational session with context/memory controls |
| `Policy` | OPA bundle + OpenFGA model binding |

## Build
- controller-runtime / Kubebuilder operator (Apache-2.0).
- State in etcd (via the API server). Delivery reconciled by Argo CD.

> Design only — no operator code yet. Open decision: implementation language (Go default).
