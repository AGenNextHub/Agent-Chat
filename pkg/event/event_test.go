package event

import (
	"context"
	"errors"
	"testing"
)

func TestEventValidate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		ev      Event
		wantErr bool
	}{
		{"valid", New("1", "web", "chat.message.v1", "sess-1", "acme", "user-1", "chat.answer", []byte("hi")), false},
		{"missing id", New("", "web", "chat.message.v1", "sess-1", "acme", "user-1", "chat.answer", nil), true},
		{"missing tenant", New("1", "web", "chat.message.v1", "sess-1", "", "user-1", "chat.answer", nil), true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.ev.Validate()
			if tc.wantErr && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tc.wantErr && !errors.Is(err, ErrInvalidEvent) {
				t.Fatalf("expected ErrInvalidEvent, got %v", err)
			}
		})
	}
}

func TestMemBusPubSub(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	bus := NewMemBus()
	ch, err := bus.Subscribe(ctx, "chat.message.v1")
	if err != nil {
		t.Fatalf("subscribe: %v", err)
	}
	want := New("e1", "web", "chat.message.v1", "s1", "acme", "u1", "chat.answer", []byte("hello"))
	if err := bus.Publish(ctx, want); err != nil {
		t.Fatalf("publish: %v", err)
	}
	got := <-ch
	if got.ID != want.ID || got.Tenant() != "acme" {
		t.Fatalf("got %+v, want id=e1 tenant=acme", got)
	}
}

func TestMemBusRejectsInvalid(t *testing.T) {
	t.Parallel()
	bus := NewMemBus()
	if err := bus.Publish(context.Background(), Event{}); err == nil {
		t.Fatal("expected publish of invalid event to fail")
	}
}
