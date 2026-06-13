# Deploy runbook — k3s

Step-by-step to run AGenNext Chat on a k3s cluster. Every command here was
checked against the chart, CRDs, and the `agennextd` binary in this repo. The
order is: **gate → publish → install → verify**.

## 0. Prerequisites

- `kubectl` pointed at your k3s cluster (`kubectl config current-context`).
- `helm` v3.
- An image published by the **release** workflow (`.github/workflows/release.yaml`):
  it builds from `Dockerfile` and pushes `ghcr.io/agennext/agent-chat` on a
  human-cut semver tag. The image is `distroless static + nonroot`, port `8080`.

This runbook deploys **reviewed, merged** code. The merge (PR review gate) and the
release tag (publish gate) are human steps — do them first.

## 1. Publish the image (CI, once per release)

The image is built and pushed by CI when you cut a tag — this is the release gate:

```bash
# on the merged main branch, at the version you are releasing
git tag v0.1.0
git push origin v0.1.0
```

The `release` workflow logs in to GHCR, builds with SLSA provenance + an SBOM
attestation, signs the digest with cosign (keyless), and prints the
`…@sha256:` digest to the run summary. For an immutable deploy, copy that digest.

If the GHCR package is private, make it public, or create a pull secret in the
namespace (step 2) and reference it.

## 2. Install the CRDs

The platform's contracts are CRDs (Tenant, Capability, Policy, AgentSession):

```bash
kubectl apply -k deploy/crds/
kubectl get crds | grep agennext
```

## 3. Install the chart

```bash
helm install agennext-chat deploy/helm/agennext-chat \
  --namespace agennext --create-namespace \
  --set image.repository=ghcr.io/agennext/agent-chat \
  --set image.tag=0.1.0
```

Production: pin the **digest** instead of the tag (immutable):

```bash
helm upgrade --install agennext-chat deploy/helm/agennext-chat \
  --namespace agennext \
  --set image.repository=ghcr.io/agennext/agent-chat \
  --set image.tag=0.1.0@sha256:<digest-from-step-1>
```

The chart defaults are HA-minded: 3 stateless replicas, a PodDisruptionBudget
(`minAvailable: 2`), `maxUnavailable: 0` rollouts, topology spread across nodes,
and a hardened pod (`runAsNonRoot`, read-only rootfs, all caps dropped,
`RuntimeDefault` seccomp). On a single-node k3s, relax spread/PDB:

```bash
  --set replicaCount=1 --set podDisruptionBudget.enabled=false
```

### Private GHCR (only if the package is not public)

```bash
kubectl create secret docker-registry ghcr \
  --namespace agennext \
  --docker-server=ghcr.io \
  --docker-username=<github-user> \
  --docker-password=<github-token-with-read:packages>
# then add imagePullSecrets to the pod (values override or patch the ServiceAccount)
```

## 4. Verify the rollout

```bash
kubectl -n agennext rollout status deploy/agennext-chat --timeout=120s
kubectl -n agennext get pods -l app.kubernetes.io/name=agennext-chat
```

Pods become Ready when `GET /readyz` returns 200 (served by `agennextd`, mapped in
`pkg/server/server.go`). Liveness is `GET /healthz`.

## 5. Smoke test

```bash
kubectl -n agennext port-forward svc/agennext-chat 8080:8080 &
curl -fsS http://127.0.0.1:8080/healthz && echo " healthz ok"
curl -fsS -X POST http://127.0.0.1:8080/v1/chat \
  -H 'content-type: application/json' \
  -d '{"id":"evt-1","source":"channel/web","type":"chat.message.v1","session":"session-1","tenant":"acme","principal":"user-1","capability":"rag.retrieve","message":"What is your return policy?"}'
```

Expect a JSON turn result with `answer`, `iterations`, `stopped_by`, and the
inspectable `trace`.

## 6. Channels (optional)

Peer-platform channels mount at the edge in `agennextd` behind env vars
(`SLACK_SIGNING_SECRET`, `SLACK_BOT_TOKEN`, `MATTERMOST_TOKEN`). The current chart
does not yet template pod env/secrets — wiring those into `values.yaml` is the
next deploy increment. Until then, the web/API surface (`/v1/chat`) is live and
the channels run via the daemon's env when set.

## Rollback

```bash
helm rollback agennext-chat            # previous release
helm uninstall agennext-chat -n agennext   # remove (CRDs persist; delete separately if intended)
```

> GitOps alternative: `deploy/argocd/agennext-chat.yaml` reconciles this chart
> from Git instead of imperative `helm install` — the delivery loop.
