package box

import (
	"errors"
	"testing"
)

func TestDigestIsContentAddressed(t *testing.T) {
	t.Parallel()
	a := Box{MediaType: "application/vnd.agennext.tokens+json", Content: []byte(`{"primary":"#4338CA"}`)}
	da, _ := a.Digest()
	db, _ := Box{MediaType: a.MediaType, Content: a.Content}.Digest()
	if da != db {
		t.Fatal("identical content must yield identical digest")
	}
	c := a
	c.Content = []byte(`{"primary":"#000000"}`)
	dc, _ := c.Digest()
	if dc == da {
		t.Fatal("different content must yield a different digest")
	}
	ok, _ := Verify(da, a)
	if !ok {
		t.Fatal("Verify should accept the matching digest")
	}
	if ok, _ := Verify(da, c); ok {
		t.Fatal("Verify must reject tampered content")
	}
}

// buildDAG: leaf b, then c→b, then a→[b,c]. Returns the graph and root digest.
func buildDAG(t *testing.T) (*Graph, string) {
	t.Helper()
	g := NewGraph()
	bD, _ := g.Add(Box{MediaType: "leaf", Content: []byte("b")})
	cD, _ := g.Add(Box{MediaType: "node", Content: []byte("c"), Refs: []string{bD}})
	aD, _ := g.Add(Box{MediaType: "root", Content: []byte("a"), Refs: []string{bD, cD}})
	return g, aD
}

func TestGraphResolvesWithBoxAsBuildingBlock(t *testing.T) {
	t.Parallel()
	g, root := buildDAG(t)
	order, err := g.Resolve(root)
	if err != nil {
		t.Fatalf("graph must resolve: %v", err)
	}
	if len(order) != 3 {
		t.Fatalf("expected 3 boxes, got %d", len(order))
	}
	// Dependency-first: the root is last, and every ref precedes its parent.
	if order[len(order)-1] != root {
		t.Fatal("root must resolve last (deps first)")
	}
	pos := map[string]int{}
	for i, d := range order {
		pos[d] = i
	}
	for _, d := range order {
		b, _ := g.Get(d)
		for _, ref := range b.Refs {
			if pos[ref] > pos[d] {
				t.Fatalf("ref %s resolved after its parent %s", short(ref), short(d))
			}
		}
	}
}

func TestGraphRejectsCycle(t *testing.T) {
	t.Parallel()
	// Hand-build a cycle by inserting boxes whose refs point at each other's
	// digests (digests computed without the back-ref, then patched in).
	g := NewGraph()
	a := Box{MediaType: "a", Content: []byte("a")}
	b := Box{MediaType: "b", Content: []byte("b")}
	aD, _ := a.Digest()
	bD, _ := b.Digest()
	a.Refs = []string{bD}
	b.Refs = []string{aD}
	// Store under the digests the refs name (pre-patch identities).
	g.boxes[aD] = a
	g.boxes[bD] = b
	if _, err := g.Resolve(aD); !errors.Is(err, ErrCycle) {
		t.Fatalf("expected ErrCycle, got %v", err)
	}
}

func TestGraphRejectsMissingRef(t *testing.T) {
	t.Parallel()
	g := NewGraph()
	root, _ := g.Add(Box{MediaType: "root", Content: []byte("a"), Refs: []string{"sha256:deadbeef"}})
	if _, err := g.Resolve(root); !errors.Is(err, ErrMissingRef) {
		t.Fatalf("expected ErrMissingRef, got %v", err)
	}
}

// TestTamperBreaksResolution shows the Merkle property: mutating a child changes
// its digest, so the parent's ref no longer resolves — the graph fails closed.
func TestTamperBreaksResolution(t *testing.T) {
	t.Parallel()
	g, root := buildDAG(t)
	// Replace child "b" content; its digest changes, leaving the old ref dangling.
	for d, b := range g.boxes {
		if string(b.Content) == "b" {
			delete(g.boxes, d)
			tampered := b
			tampered.Content = []byte("b-tampered")
			_, _ = g.Add(tampered) // stored under a NEW digest
		}
	}
	if _, err := g.Resolve(root); !errors.Is(err, ErrMissingRef) {
		t.Fatalf("tamper must break resolution, got %v", err)
	}
}
