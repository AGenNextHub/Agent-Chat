// Package guard implements the platform's in-loop and at-gate defences:
//
//   - Screener: a model-agnostic prompt-injection detector applied to untrusted
//     text (user input AND retrieved/tool output). This is the GenTel-Shield
//     slot; HeuristicScreener is a deterministic baseline meant to be swapped
//     for a trained detector behind the same interface.
//   - PromptProvider: the system-level Guard Prompts that establish irrecoverable
//     behavioural constraints.
//   - Decider: an in-process, OPA-shaped policy decision point used by the GUARD
//     step and the Edge Gate (default deny-unknown).
//
// Defences derive from arXiv:2601.15528v1 §4. See docs/THREAT_MODEL.md.
package guard

import (
	"context"
	"regexp"
)

// Origin labels where a piece of text came from, so screening can be applied
// with the right suspicion. Retrieved and tool output are the indirect-injection
// surface (threat T2) and must be screened before re-entering the loop.
type Origin string

const (
	// OriginUser is text supplied directly by the end user.
	OriginUser Origin = "user"
	// OriginRetrieved is text pulled from a knowledge base during RAG.
	OriginRetrieved Origin = "retrieved"
	// OriginTool is output returned by an invoked capability/tool.
	OriginTool Origin = "tool"
)

// Verdict is the result of screening a piece of text.
type Verdict struct {
	// Malicious is true when injection characteristics are detected.
	Malicious bool
	// Reason is a short, human-readable explanation (inspectability).
	Reason string
	// Score is a coarse confidence in [0,1].
	Score float64
}

// Screener classifies untrusted text for prompt-injection intent.
type Screener interface {
	Screen(ctx context.Context, text string, origin Origin) Verdict
}

// injectionPatterns are deliberately conservative, documented signatures. They
// are a transparent baseline, not a substitute for a trained detector.
var injectionPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)ignore\s+(all\s+|the\s+|your\s+)?(previous|prior|above)\s+(instruction|prompt|rule)`),
	regexp.MustCompile(`(?i)disregard\s+(all\s+|the\s+|your\s+)?(previous|prior|above)?\s*(instruction|rule|prompt)`),
	regexp.MustCompile(`(?i)(reveal|show|print|leak)\s+(the\s+|your\s+)?(system\s+prompt|hidden|internal|secret|api\s*key|password)`),
	regexp.MustCompile(`(?i)you\s+are\s+now\s+`),
	regexp.MustCompile(`(?i)pretend\s+(you\s+are|to\s+be)\s+`),
	regexp.MustCompile(`(?i)(new|updated)\s+(rule|policy|system\s+patch|instruction)s?\b`),
	regexp.MustCompile(`(?i)\bexfiltrat`),
	regexp.MustCompile(`(?i)override\s+(the\s+|your\s+)?(instruction|policy|guard|safety)`),
}

// HeuristicScreener is a deterministic, dependency-free baseline detector.
type HeuristicScreener struct{}

// NewHeuristicScreener returns a baseline screener.
func NewHeuristicScreener() HeuristicScreener { return HeuristicScreener{} }

// Screen flags text matching any known injection signature. Retrieved and tool
// origins are treated as higher risk and reported with their origin.
func (HeuristicScreener) Screen(_ context.Context, text string, origin Origin) Verdict {
	matches := 0
	for _, re := range injectionPatterns {
		if re.MatchString(text) {
			matches++
		}
	}
	if matches == 0 {
		return Verdict{Malicious: false, Score: 0}
	}
	score := float64(matches) / float64(len(injectionPatterns))
	if score > 1 {
		score = 1
	}
	return Verdict{
		Malicious: true,
		Reason:    "matched " + string(origin) + " injection signature",
		Score:     score,
	}
}
