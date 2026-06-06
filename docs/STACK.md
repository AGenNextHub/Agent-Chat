# Technology Stack — CNCF-native, maturity-ranked, no lock-in

Selection rules (from requirements):
1. **CNCF / cloud-native first** — prefer the [CNCF landscape](https://landscape.cncf.io/).
2. **Pick by maturity** — prefer **Graduated** > **Incubating** > **Sandbox**.
3. **No vendor lock** — open governance, self-hostable, swappable behind a contract.
4. **No license lock** — permissive/OSI (Apache-2.0, BSD, MIT, MPL-2.0). **Reject**
   AGPL/SSPL/BUSL/RSAL re-licensed projects.

> Every choice sits behind a **capability contract** (`CAPABILITIES.md`), so any tool here
> can be swapped for a contract-compatible alternative. The list is a default, not a cage.

## Selected stack

| Layer | Need | Choice | Maturity | License |
|---|---|---|---|---|
| Kernel | Orchestration | **Kubernetes** | Graduated | Apache-2.0 |
| Kernel | Edge distro | **k3s** | CNCF (Sandbox) | Apache-2.0 |
| Kernel | Container runtime | **containerd** | Graduated | Apache-2.0 |
| Kernel | Operator framework | **controller-runtime / Kubebuilder** | k8s SIG | Apache-2.0 |
| Kernel | Cluster state | **etcd** | Graduated | Apache-2.0 |
| Delivery | GitOps loop | **Argo CD** | Graduated | Apache-2.0 |
| Delivery | Event/pipeline loop | **Argo Events + Workflows** | Graduated/Incubating | Apache-2.0 |
| Packaging | Charts/release | **Helm** | Graduated | Apache-2.0 |
| Packaging | OCI artifacts (capabilities) | **ORAS** | CNCF | Apache-2.0 |
| Packaging | Registry | **Harbor** | Graduated | Apache-2.0 |
| Identity | Workload identity / mTLS | **SPIFFE / SPIRE** | Graduated | Apache-2.0 |
| Policy | Admission/policy | **OPA / Gatekeeper** | Graduated | Apache-2.0 |
| AuthZ | Fine-grained relations | **OpenFGA** | Incubating | Apache-2.0 |
| Network | CNI + encrypted overlay | **Cilium** (eBPF + WireGuard) | Graduated | Apache-2.0 |
| Network | Service mesh / mTLS | **Linkerd** | Graduated | Apache-2.0 |
| Edge gate | Ingress / API gateway | **Envoy + Gateway API** | Graduated | Apache-2.0 |
| Eventing | Event format | **CloudEvents** | Graduated | Apache-2.0 |
| Eventing | Messaging bus (multi-channel) | **NATS** | Incubating | Apache-2.0 |
| Observability | Telemetry | **OpenTelemetry** | Graduated | Apache-2.0 |
| Observability | Metrics | **Prometheus** | Graduated | Apache-2.0 |
| Observability | Edge logs | **Fluent Bit** | Graduated | Apache-2.0 |
| Observability | Tracing backend | **Jaeger** | Graduated | Apache-2.0 |
| Observability | Dashboards | **Perses** | CNCF (Sandbox) | Apache-2.0 |
| Secrets | Secret sync | **External Secrets Operator** | Incubating | Apache-2.0 |
| Secrets | Secret store | **OpenBao** | LF Edge / OpenSSF | MPL-2.0 |
| State | Persistent memory + vectors | **PostgreSQL + pgvector** | (de-facto std) | PostgreSQL / permissive |
| State | Ephemeral context cache | **Valkey** | LF | BSD-3 |
| State | Scale-out vector DB (opt.) | **Qdrant** | — | Apache-2.0 |
| AI | Model serving | **KServe** | Incubating | Apache-2.0 |
| AI | Local inference | **vLLM** / **Ollama** | — | Apache-2.0 / MIT |
| AI | Provider-agnostic gateway | **LiteLLM** | — | MIT |

## Explicitly rejected (license lock)

| Tool | Why rejected | Chosen instead |
|---|---|---|
| **Grafana** | AGPLv3 | Perses (Apache-2.0) |
| **Redis** | RSALv2 / SSPL | Valkey (BSD-3) |
| **HashiCorp Vault / Consul** | BUSL | OpenBao (MPL-2.0) / Linkerd |
| **Elasticsearch** (classic) | SSPL | OpenSearch (Apache-2.0) if needed |

## "Build a loop" — the three reconciliation loops

AGenNext Chat is loop-driven at every layer; nothing is imperative one-shot.

1. **Kernel control loop** — controller-runtime reconciles `Tenant` / `Capability` /
   `AgentSession` / `Policy` CRDs to desired state. *This is the kernel.*
2. **Delivery loop (GitOps)** — Argo CD continuously reconciles cluster state to Git;
   capabilities are OCI artifacts referenced from Git. Distributed edge clusters each
   pull their slice.
3. **Agent loop** — within the Runtime Core, a reason→act→observe loop drives RAG +
   tool use, **event-driven via CloudEvents over NATS**. Multi-channel adapters publish
   inbound events; the loop consumes, reasons, and emits responses back to the channel.

## AI-native, multi-channel, multi-device, distributed

- **AI-native:** inference (KServe/vLLM), RAG, and the agent loop are first-class kernel
  workloads, not bolt-ons; model choice is provider-agnostic (LiteLLM) — *no vendor lock*.
- **Multi-channel:** each channel (web, chat apps, email, voice) is a capability that
  speaks **CloudEvents** on **NATS**; adding a channel = adding a contract-bound adapter.
- **Multi-device:** clients are stateless against the Edge Gate; session context/memory
  live server-side and are addressable across devices via the agent-session identity.
- **Distributed:** heterogeneous low-cost nodes joined by the Cilium/WireGuard overlay
  ("the matrix"), tolerant of 60–200 ms latency, scheduled by the kernel.

> Open decisions (encoding, primary inference backend, overlay specifics) remain tracked
> in [`SCOPE.md §5`](SCOPE.md) — chosen here as defaults, confirmable before code.
