// Package capability implements the platform's atomic, composable unit:
// the capability, which IS a contract. A capability declares exactly what it
// provides, requires, and may touch; nothing runs outside its contract.
// Services and solutions are compositions of capability contracts.
// See docs/CAPABILITIES.md and docs/CONCEPTS.md.
package capability

import "path"

// Scope is the least-privilege authority a capability may exercise: the
// tenants, data globs, and tools it is allowed to access. An empty list means
// "none"; the wildcard "*" means "any".
type Scope struct {
	// Tenants this capability may act for ("*" = any).
	Tenants []string
	// Data glob patterns (path.Match syntax) it may read, e.g. "tenant://acme/kb/*".
	Data []string
	// Tools (capability names) it may invoke.
	Tools []string
}

// Contains reports whether req is a subset of s (req ⊆ s): every tenant, data
// pattern, and tool requested must be permitted by s. This is the predicate the
// Edge Gate and the agent loop's GUARD step enforce.
func (s Scope) Contains(req Scope) bool {
	return subsetExact(s.Tenants, req.Tenants) &&
		subsetGlob(s.Data, req.Data) &&
		subsetExact(s.Tools, req.Tools)
}

// subsetExact reports whether every element of req is permitted by granted,
// where granted may contain the wildcard "*".
func subsetExact(granted, req []string) bool {
	for _, r := range req {
		if !containsExact(granted, r) {
			return false
		}
	}
	return true
}

func containsExact(granted []string, r string) bool {
	for _, g := range granted {
		if g == "*" || g == r {
			return true
		}
	}
	return false
}

// subsetGlob reports whether every requested pattern is covered by a granted
// glob (path.Match). A granted "*" covers anything.
func subsetGlob(granted, req []string) bool {
	for _, r := range req {
		if !coveredByGlob(granted, r) {
			return false
		}
	}
	return true
}

func coveredByGlob(granted []string, r string) bool {
	for _, g := range granted {
		if g == "*" || g == r {
			return true
		}
		if ok, err := path.Match(g, r); err == nil && ok {
			return true
		}
	}
	return false
}

// Empty reports whether the scope grants no authority at all. An empty scope is
// the safe default-deny: it permits nothing.
func (s Scope) Empty() bool {
	return len(s.Tenants) == 0 && len(s.Data) == 0 && len(s.Tools) == 0
}

// Intersect returns the least-privilege scope permitted by BOTH s and o — the
// effective authority when a principal's grant meets a capability's contract.
// It is the mechanism that lets the Edge Gate DERIVE an admitted scope instead
// of trusting a caller-supplied one.
func (s Scope) Intersect(o Scope) Scope {
	return Scope{
		Tenants: intersectExact(s.Tenants, o.Tenants),
		Data:    intersectGlob(s.Data, o.Data),
		Tools:   intersectExact(s.Tools, o.Tools),
	}
}

func hasWildcard(list []string) bool {
	for _, x := range list {
		if x == "*" {
			return true
		}
	}
	return false
}

// intersectExact intersects two exact-match lists, honoring "*" as "any".
func intersectExact(a, b []string) []string {
	if hasWildcard(a) {
		return dedup(b)
	}
	if hasWildcard(b) {
		return dedup(a)
	}
	var out []string
	for _, x := range a {
		if containsExact(b, x) {
			out = append(out, x)
		}
	}
	return dedup(out)
}

// intersectGlob intersects two glob lists, keeping patterns covered by both
// sides (the tighter, mutually-permitted set). "*" means "any".
func intersectGlob(a, b []string) []string {
	if hasWildcard(a) {
		return dedup(b)
	}
	if hasWildcard(b) {
		return dedup(a)
	}
	var out []string
	for _, x := range a {
		if coveredByGlob(b, x) {
			out = append(out, x)
		}
	}
	for _, y := range b {
		if coveredByGlob(a, y) {
			out = append(out, y)
		}
	}
	return dedup(out)
}

func dedup(in []string) []string {
	if len(in) == 0 {
		return nil
	}
	seen := make(map[string]bool, len(in))
	var out []string
	for _, x := range in {
		if !seen[x] {
			seen[x] = true
			out = append(out, x)
		}
	}
	return out
}
