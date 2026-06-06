// Package store defines the session state surfaces the agent loop reads and
// writes: working Context and persistent Memory. Both expose load/clear so the
// platform's "clear context" / "no hidden logic" tenets are first-class and
// auditable. The production bindings are Valkey (context cache) and
// PostgreSQL+pgvector (memory); the in-memory implementations here back tests
// and single-node runs. See docs/PRINCIPLES.md.
package store

import (
	"context"
	"sync"
	"time"
)

// Message is one turn in a session's working context.
type Message struct {
	// Role is "user" or "assistant".
	Role string
	// Content is the message text.
	Content string
	// At is when the message was recorded.
	At time.Time
}

// Context is a session's working context window.
type Context struct {
	// Messages is the ordered conversation history.
	Messages []Message
	// Updated is the last write time.
	Updated time.Time
}

// MemoryItem is one durable fact retained across turns.
type MemoryItem struct {
	// Key categorizes the item (e.g. "retrieval", "summary").
	Key string
	// Value is the retained content.
	Value string
	// At is when the item was stored.
	At time.Time
}

// ContextStore persists working context per session.
type ContextStore interface {
	Load(ctx context.Context, sessionID string) (Context, error)
	Save(ctx context.Context, sessionID string, c Context) error
	Clear(ctx context.Context, sessionID string) error
}

// MemoryStore persists durable memory per session.
type MemoryStore interface {
	Load(ctx context.Context, sessionID string) ([]MemoryItem, error)
	Append(ctx context.Context, sessionID string, item MemoryItem) error
	Clear(ctx context.Context, sessionID string) error
}

// MemContextStore is an in-memory ContextStore.
type MemContextStore struct {
	mu sync.RWMutex
	m  map[string]Context
}

// NewMemContextStore returns an empty in-memory context store.
func NewMemContextStore() *MemContextStore { return &MemContextStore{m: make(map[string]Context)} }

// Load returns the stored context for sessionID (zero value if absent).
func (s *MemContextStore) Load(_ context.Context, sessionID string) (Context, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.m[sessionID], nil
}

// Save replaces the stored context for sessionID.
func (s *MemContextStore) Save(_ context.Context, sessionID string, c Context) error {
	c.Updated = time.Now().UTC()
	s.mu.Lock()
	defer s.mu.Unlock()
	s.m[sessionID] = c
	return nil
}

// Clear removes the stored context for sessionID.
func (s *MemContextStore) Clear(_ context.Context, sessionID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.m, sessionID)
	return nil
}

// MemMemoryStore is an in-memory MemoryStore.
type MemMemoryStore struct {
	mu sync.RWMutex
	m  map[string][]MemoryItem
}

// NewMemMemoryStore returns an empty in-memory memory store.
func NewMemMemoryStore() *MemMemoryStore { return &MemMemoryStore{m: make(map[string][]MemoryItem)} }

// Load returns the durable memory items for sessionID.
func (s *MemMemoryStore) Load(_ context.Context, sessionID string) ([]MemoryItem, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	items := s.m[sessionID]
	out := make([]MemoryItem, len(items))
	copy(out, items)
	return out, nil
}

// Append adds one durable memory item for sessionID.
func (s *MemMemoryStore) Append(_ context.Context, sessionID string, item MemoryItem) error {
	if item.At.IsZero() {
		item.At = time.Now().UTC()
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.m[sessionID] = append(s.m[sessionID], item)
	return nil
}

// Clear removes all durable memory for sessionID.
func (s *MemMemoryStore) Clear(_ context.Context, sessionID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.m, sessionID)
	return nil
}
