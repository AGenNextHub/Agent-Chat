package box

import (
	"crypto/ed25519"
	"crypto/rand"
	"errors"
	"testing"
)

// TestBoxLifecycle exercises the canonical verbs: Define → Build → Sign → Share.
func TestBoxLifecycle(t *testing.T) {
	t.Parallel()

	// Define Box (form only) → Build Box (form + content/ground).
	b := Define("application/vnd.agennext.tokens+json", "oci://schema/tokens").
		Build([]byte(`{"primary":"#4338CA"}`))

	// Sign Box — the HUMAN signs (the test plays the key holder).
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("keygen: %v", err)
	}
	signed, err := Sign(b, priv)
	if err != nil {
		t.Fatalf("sign: %v", err)
	}
	if ok, err := signed.Verify(); err != nil || !ok {
		t.Fatalf("signed box must verify (ok=%v err=%v)", ok, err)
	}

	// Agent path: Attach a human-produced signature WITHOUT holding the key.
	digest, _ := b.Digest()
	humanSig := ed25519.Sign(priv, []byte(digest))
	att, err := Attach(b, pub, humanSig)
	if err != nil {
		t.Fatalf("attach: %v", err)
	}
	if att.Digest != signed.Digest {
		t.Fatal("attach and sign must agree on the digest")
	}

	// Share Box — fail-closed store; only verified boxes are admitted.
	st := NewCollab()
	d, err := st.Publish(signed)
	if err != nil {
		t.Fatalf("share: %v", err)
	}
	if _, ok := st.Fetch(d); !ok {
		t.Fatal("shared box must be fetchable by digest")
	}
}

func TestVerifyRejectsTamperAfterSign(t *testing.T) {
	t.Parallel()
	_, priv, _ := ed25519.GenerateKey(rand.Reader)
	signed, _ := Sign(Define("m", "s").Build([]byte("real")), priv)
	signed.Box.Content = []byte("tampered") // mutate after signing
	if ok, _ := signed.Verify(); ok {
		t.Fatal("tampering content after signing must fail verification")
	}
}

func TestPublishRejectsUnsigned(t *testing.T) {
	t.Parallel()
	st := NewCollab()
	b := Define("m", "s").Build([]byte("x"))
	bad := Sealed{Box: b}
	bad.Digest, _ = b.Digest() // right digest, no/garbage signature
	if _, err := st.Publish(bad); !errors.Is(err, ErrUnsealed) {
		t.Fatalf("unsigned box must be rejected, got %v", err)
	}
}

func TestAttachRejectsBadSignature(t *testing.T) {
	t.Parallel()
	pub, _, _ := ed25519.GenerateKey(rand.Reader)
	b := Define("m", "s").Build([]byte("x"))
	if _, err := Attach(b, pub, []byte("not-a-real-signature")); !errors.Is(err, ErrUnsealed) {
		t.Fatalf("bad signature must be rejected, got %v", err)
	}
}
