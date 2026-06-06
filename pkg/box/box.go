// Package box is the platform's primitive noun: the content-addressed Box.
//
// A Box's identity IS its content (digest = sha256 of its canonical encoding),
// so it is self-grounding, immutable, and tamper-evident. Boxes reference other
// boxes by digest, forming a Merkle DAG — the composition graph. The Graph
// resolver answers the only question that matters of a building block: does the
// graph resolve (every ref present, no cycles)?
//
// Pure standard library (crypto/sha256, encoding/json) — zero dependencies.
package box

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
)

// Box is a bounded unit of content. Field order is fixed so the JSON encoding —
// and therefore the digest — is deterministic. Content is the ground (opaque,
// undefined here); Refs are child digests (composition); Lineage are parent
// digests (history).
type Box struct {
	MediaType string   `json:"mediaType"`
	Schema    string   `json:"schema,omitempty"`
	Content   []byte   `json:"content,omitempty"`
	Refs      []string `json:"refs,omitempty"`
	Lineage   []string `json:"lineage,omitempty"`
}

// Definition is the FORM of a box — its media type and schema (contract) — with
// no content yet. "Define Box" declares the form; "Build Box" fills it.
type Definition struct {
	MediaType string
	Schema    string
}

// Define declares a box's form (the contract), never its content. The ground
// stays opaque and self-grounding. — Define Box
func Define(mediaType, schema string) Definition {
	return Definition{MediaType: mediaType, Schema: schema}
}

// Build instantiates the definition with content (the ground) plus optional
// child box digests (composition). — Build Box
func (d Definition) Build(content []byte, refs ...string) Box {
	return Box{MediaType: d.MediaType, Schema: d.Schema, Content: content, Refs: refs}
}

// Encode returns the canonical, deterministic byte encoding of the box (its
// `.box` form). The digest is taken over exactly these bytes.
func (b Box) Encode() ([]byte, error) { return json.Marshal(b) }

// Digest returns the box's self-assigned content identity, "sha256:<hex>".
// Because Refs are themselves digests, a box's identity transitively depends on
// its children — the Merkle property.
func (b Box) Digest() (string, error) {
	enc, err := b.Encode()
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(enc)
	return "sha256:" + hex.EncodeToString(sum[:]), nil
}

// Graph is a content-addressed set of boxes (digest → box).
type Graph struct {
	boxes map[string]Box
}

// Errors returned by Resolve.
var (
	// ErrMissingRef means a referenced box is not in the graph (unresolvable).
	ErrMissingRef = errors.New("missing ref")
	// ErrCycle means the reference graph is not a DAG.
	ErrCycle = errors.New("cycle")
)

// NewGraph returns an empty graph.
func NewGraph() *Graph { return &Graph{boxes: make(map[string]Box)} }

// Add stores a box under its content digest and returns that digest.
func (g *Graph) Add(b Box) (string, error) {
	d, err := b.Digest()
	if err != nil {
		return "", err
	}
	g.boxes[d] = b
	return d, nil
}

// Get returns the box for a digest and whether it is present.
func (g *Graph) Get(digest string) (Box, bool) {
	b, ok := g.boxes[digest]
	return b, ok
}

// Resolve walks the DAG rooted at digest and returns a dependency-first
// (topological) ordering. It fails closed: a missing ref or a cycle is an error,
// so "the graph resolves" is a checked property, not an assumption.
func (g *Graph) Resolve(root string) ([]string, error) {
	const (
		gray  = 1 // on the current DFS path
		black = 2 // fully resolved
	)
	state := make(map[string]int)
	var order []string

	var visit func(d string) error
	visit = func(d string) error {
		switch state[d] {
		case black:
			return nil
		case gray:
			return fmt.Errorf("%w: at %s", ErrCycle, short(d))
		}
		b, ok := g.boxes[d]
		if !ok {
			return fmt.Errorf("%w: %s", ErrMissingRef, short(d))
		}
		state[d] = gray
		for _, ref := range b.Refs {
			if err := visit(ref); err != nil {
				return err
			}
		}
		state[d] = black
		order = append(order, d) // dependency-first
		return nil
	}

	if err := visit(root); err != nil {
		return nil, err
	}
	return order, nil
}

// Verify recomputes a box's digest and reports whether it matches the claimed
// identity (tamper check).
func Verify(digest string, b Box) (bool, error) {
	d, err := b.Digest()
	if err != nil {
		return false, err
	}
	return d == digest, nil
}

func short(d string) string {
	if len(d) > 14 {
		return d[:14]
	}
	return d
}
