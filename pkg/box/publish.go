package box

// Collab is the collaboration surface — where the work of many hands meets
// (Define → Build → Sign → Publish into a Collab). Sharing is not accountable;
// publishing into a Collab is, because every box admitted is SIGNED by an
// identified human key holder. Collab is the tool; Publish is the verb.
type Collab struct {
	m map[string]Sealed
}

// NewCollab returns an empty collaboration surface.
func NewCollab() *Collab { return &Collab{m: make(map[string]Sealed)} }

// Publish admits a signed box into the collaboration: it verifies the seal, then
// stores the box under its content digest and returns that digest. Only verified,
// human-signed boxes are admitted (fail closed) — accountability, not sharing.
func (c *Collab) Publish(sb Sealed) (string, error) {
	ok, err := sb.Verify()
	if err != nil {
		return "", err
	}
	if !ok {
		return "", ErrUnsealed
	}
	c.m[sb.Digest] = sb
	return sb.Digest, nil
}

// Fetch returns the published box for a digest and whether it is present.
func (c *Collab) Fetch(digest string) (Sealed, bool) {
	sb, ok := c.m[digest]
	return sb, ok
}
