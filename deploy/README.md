# Deploy — k3s / Helm / overlay scaffolding

Deployment foundation for the distributed edge mesh. See
[`../docs/STACK.md`](../docs/STACK.md).

## Topology
- **k3s** clusters across heterogeneous, low-cost nodes.
- **Cilium** CNI with **WireGuard** transparent encryption = the overlay ("the matrix"),
  tolerant of 60–200 ms inter-node latency.
- **Argo CD** reconciles each cluster to Git (the delivery loop).

## Layout (planned)
```
deploy/
  helm/        umbrella + per-component charts
  k3s/         cluster bootstrap (control-plane + worker roles)
  argocd/      app-of-apps GitOps definitions
  overlay/     Cilium/WireGuard config
```

## Principles
- Capabilities pulled as **OCI** artifacts (ORAS) from **Harbor**.
- Secrets via External Secrets Operator + **OpenBao** (no BUSL lock).
- Telemetry via OpenTelemetry → Prometheus / Jaeger / Perses (no AGPL lock).

> Scaffolding outline only — manifests/charts to follow.
