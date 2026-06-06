package box

import (
	"crypto/ed25519"
	"encoding/json"
	"errors"
)

// Sealed is a signed, finalized box: it binds the box's content (via its digest)
// to a signer (an Ed25519 public key) with a signature over that digest.
// Verifying a Sealed box checks BOTH integrity (content still hashes to digest)
// and authenticity (the signature is valid for that digest under the key).
// Ed25519 is standard-library — no third-party crypto in the supply chain.
type Sealed struct {
	Box       Box    `json:"box"`
	Digest    string `json:"digest"`
	PublicKey []byte `json:"publicKey"` // ed25519 public key
	Signature []byte `json:"signature"` // sign(digest)
}

// ErrUnsealed is returned when a sealed box fails verification.
var ErrUnsealed = errors.New("seal invalid")

// Sign computes the box's content digest, signs it with priv, and returns the
// sealed envelope. The seal is over the digest, so it transitively covers all
// content and (via the Merkle digest) every referenced child box. — Sign Box
//
// IMPORTANT: Sign holds a PRIVATE KEY and therefore belongs to a HUMAN's signing
// tool (or tests) — never to the agent or the daemon. A signature the agent
// makes over its own output establishes no real-world accountability; realness
// is rooted in a human key holder. The agent uses Attach + Verify only, and must
// never possess a private key. (human is real; if the agent signs, it isn't.)
func Sign(b Box, priv ed25519.PrivateKey) (Sealed, error) {
	d, err := b.Digest()
	if err != nil {
		return Sealed{}, err
	}
	pub, ok := priv.Public().(ed25519.PublicKey)
	if !ok {
		return Sealed{}, errors.New("box: invalid private key")
	}
	return Sealed{
		Box:       b,
		Digest:    d,
		PublicKey: pub,
		Signature: ed25519.Sign(priv, []byte(d)),
	}, nil
}

// Verify checks integrity then authenticity. It returns true only if the box's
// content still hashes to the sealed digest AND the signature is valid for that
// digest under the sealed public key.
func (s Sealed) Verify() (bool, error) {
	ok, err := Verify(s.Digest, s.Box) // integrity: content ↔ digest
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	if len(s.PublicKey) != ed25519.PublicKeySize {
		return false, nil // a malformed seal is invalid, not an error
	}
	return ed25519.Verify(ed25519.PublicKey(s.PublicKey), []byte(s.Digest), s.Signature), nil
}

// Attach assembles a Sealed box from a signature produced OUT OF BAND by a human
// (or a real-world key holder). The agent computes the digest and attaches the
// human's signature — it never holds a private key and never signs. The seal is
// verified before return, so a wrong or tampered signature is rejected here.
// This is the agent's path; Sign is the human's.
func Attach(b Box, pub ed25519.PublicKey, sig []byte) (Sealed, error) {
	d, err := b.Digest()
	if err != nil {
		return Sealed{}, err
	}
	s := Sealed{Box: b, Digest: d, PublicKey: pub, Signature: sig}
	ok, err := s.Verify()
	if err != nil {
		return Sealed{}, err
	}
	if !ok {
		return Sealed{}, ErrUnsealed
	}
	return s, nil
}

// Encode returns the canonical `.box` envelope bytes (a sealed box on the wire).
func (s Sealed) Encode() ([]byte, error) { return json.Marshal(s) }
