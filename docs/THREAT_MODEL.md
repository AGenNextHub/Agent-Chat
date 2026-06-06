# Threat Model

Derived from arXiv:2601.15528v1 §4 and generalized to the AGenNext Chat layering.

## Assets

- Tenant knowledge bases (may contain PII despite screening).
- System prompts, guard prompts, policies (must remain non-exfiltrable).
- Cross-tenant boundary (its integrity is the platform's core promise).

## Adversaries & attack surfaces

| # | Threat | Surface | Primary mitigation |
|---|---|---|---|
| T1 | **Direct prompt injection** — user input overrides instructions | Edge → Runtime | Guard Prompts (L1) + injection detector before inference |
| T2 | **Indirect prompt injection** — malicious instructions in retrieved docs/web | RAG retrieval | Detector applied to *retrieved content*, not just queries |
| T3 | **Prompt / policy leakage** — probing for system rules | Inference | Guard Prompts prohibit disclosure of internal rules |
| T4 | **Cross-tenant data leakage** | Runtime / storage | Namespace isolation + capability `scope` ⊆ tenant; OpenFGA |
| T5 | **Privilege / capability escalation** | Capability invocation | Contract admission; OPA policy; least-privilege scope |
| T6 | **Compromised module blast radius** | Any Runtime module | Container isolation; sandbox mode; kernel-enforced boundaries |
| T7 | **PII exposure at ingestion** | Document upload | PII screening + de-identification before indexing |
| T8 | **Untrusted peer / overlay MITM** | Federation overlay | mTLS, signature verification, peers ≤ client trust |

## Layered prompt-injection defence (T1–T3)

Two complementary layers, per the case study:

1. **Guard Prompts (rule-based, model-agnostic):** system-level constraints that prohibit
   role switching, permission escalation, execution of instructions embedded in retrieved
   content, and disclosure of internal prompts/rules. Near-zero runtime overhead.
   *Limitation:* static; weaker against obfuscated/indirect injection.
2. **Pre-generation detector (learned, GenTel-Shield-style):** classifies queries *and*
   retrieved content; malicious inputs are blocked before reaching the LLM. Model-agnostic,
   no retraining.

Case-study evidence (for calibration of success criteria):

| Config | Recall | F1 |
|---|---|---|
| Pure LLM | ~0.4–1.2% | ~1–3% |
| Guard Prompts | 99.6–100% | ≈99.8–100% |
| Detector only | ~81.6% | ~89.7% |
| **Guard Prompts + Detector** | **100%** | **≈99.8%** |

Conclusion carried into design: **layer both**. Guard Prompts give a high ceiling but need
tuning and may not generalize; the learned detector gives stable, model-agnostic recall.

## Non-goals (see SCOPE §3)

High-stakes autonomous action, unrestricted generation, and model retraining are out of
scope; this narrows the attack surface by construction.
