package store

import (
	"context"
	"testing"
)

func TestMemContextStore(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	s := NewMemContextStore()
	if c, _ := s.Load(ctx, "s1"); len(c.Messages) != 0 {
		t.Fatal("expected empty context for unknown session")
	}
	if err := s.Save(ctx, "s1", Context{Messages: []Message{{Role: "user", Content: "hi"}}}); err != nil {
		t.Fatalf("save: %v", err)
	}
	c, _ := s.Load(ctx, "s1")
	if len(c.Messages) != 1 || c.Updated.IsZero() {
		t.Fatalf("unexpected context: %+v", c)
	}
	if err := s.Clear(ctx, "s1"); err != nil {
		t.Fatalf("clear: %v", err)
	}
	if c, _ := s.Load(ctx, "s1"); len(c.Messages) != 0 {
		t.Fatal("expected context cleared")
	}
}

func TestMemMemoryStore(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	s := NewMemMemoryStore()
	if err := s.Append(ctx, "s1", MemoryItem{Key: "retrieval", Value: "v"}); err != nil {
		t.Fatalf("append: %v", err)
	}
	items, _ := s.Load(ctx, "s1")
	if len(items) != 1 || items[0].At.IsZero() {
		t.Fatalf("unexpected items: %+v", items)
	}
	if err := s.Clear(ctx, "s1"); err != nil {
		t.Fatalf("clear: %v", err)
	}
	if items, _ := s.Load(ctx, "s1"); len(items) != 0 {
		t.Fatal("expected memory cleared")
	}
}
