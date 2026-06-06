# References

## Grounding paper
1. J. Xie, B. Li, H. Fu, C. Gao, Z. Xu, F. Han. *Securing LLM-as-a-Service for Small
   Businesses: An Industry Case Study of a Distributed Chatbot Deployment Platform.*
   AISC 2026. arXiv:2601.15528v1 [cs.DC]. <https://arxiv.org/abs/2601.15528>
   — distributed k3s edge cloud, multi-tenant isolation, layered prompt-injection defence,
   PII screening. Source: <https://aisuko.github.io/secure_llm/>

## Security / prompt injection
2. R. Li et al. *GenTel-Safe: A Unified Benchmark and Shielding Framework for Defending
   Against Prompt Injection Attacks.* arXiv:2409.19521. <https://arxiv.org/abs/2409.19521>
3. K. Greshake et al. *Not what you've signed up for: Compromising Real-World
   LLM-Integrated Applications with Indirect Prompt Injection.* arXiv:2302.12173.
4. F. Perez, I. Ribeiro. *Ignore Previous Prompt: Attack Techniques For Language Models.*
   arXiv:2211.09527.
5. S. Willison. *Delimiters Won't Save You from Prompt Injection.* 2023.
   <https://simonwillison.net/2023/May/11/delimiters-wont-save-you/>

## Method / RAG / agents
6. P. Lewis et al. *Retrieval-Augmented Generation for Knowledge-Intensive NLP Tasks.*
   NeurIPS 2020.
7. S. Yao et al. *ReAct: Synergizing Reasoning and Acting in Language Models.* ICLR 2023.
   <https://openreview.net/forum?id=WE_vluYUL-X> — the reason→act loop.

## Cloud-native building blocks
8. CNCF Landscape. <https://landscape.cncf.io/>
9. Kubernetes. <https://kubernetes.io> · k3s. <https://k3s.io/>
10. CloudEvents 1.0 specification. <https://cloudevents.io/>
11. Open Policy Agent (OPA). <https://www.openpolicyagent.org/>
12. OpenFGA. <https://openfga.dev/>
13. SPIFFE/SPIRE. <https://spiffe.io/>
14. Cilium. <https://cilium.io/> · Linkerd. <https://linkerd.io/> · Envoy Gateway.
    <https://gateway.envoyproxy.io/>
15. NATS. <https://nats.io/> · OpenTelemetry. <https://opentelemetry.io/>
16. Argo CD. <https://argo-cd.readthedocs.io/> · Helm. <https://helm.sh/> · ORAS.
    <https://oras.land/>

## Supply chain
17. OpenSSF SLSA. <https://slsa.dev/>
18. Sigstore / cosign. <https://www.sigstore.dev/>
19. Go `govulncheck`. <https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck>
20. CycloneDX SBOM. <https://cyclonedx.org/>

## Regulatory (respect laws of the land)
21. Privacy Act 1988 (Australia). <https://www.legislation.gov.au/>
22. Australian Privacy Principles (APPs), OAIC.
    <https://www.oaic.gov.au/privacy/australian-privacy-principles>
