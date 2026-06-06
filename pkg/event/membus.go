package event

import (
	"context"
	"sync"
)

// Bus is a minimal publish/subscribe transport for events. The production
// binding is NATS (CloudEvents over NATS); MemBus is an in-process
// implementation used for tests and single-node runs.
type Bus interface {
	// Publish delivers an event to all subscribers whose filter matches.
	Publish(ctx context.Context, e Event) error
	// Subscribe returns a channel of events whose Type equals typeFilter,
	// or all events when typeFilter is empty.
	Subscribe(ctx context.Context, typeFilter string) (<-chan Event, error)
}

// MemBus is an in-memory, goroutine-safe Bus implementation.
type MemBus struct {
	mu   sync.RWMutex
	subs []subscription
}

type subscription struct {
	filter string
	ch     chan Event
}

// NewMemBus returns an empty in-memory bus.
func NewMemBus() *MemBus { return &MemBus{} }

// Publish fans out e to every matching subscriber without blocking.
func (b *MemBus) Publish(_ context.Context, e Event) error {
	if err := e.Validate(); err != nil {
		return err
	}
	b.mu.RLock()
	defer b.mu.RUnlock()
	for _, s := range b.subs {
		if s.filter == "" || s.filter == e.Type {
			select {
			case s.ch <- e:
			default: // drop for slow consumers; at-least-once is the bus's job in prod
			}
		}
	}
	return nil
}

// Subscribe registers a buffered subscriber. The channel is closed when ctx ends.
func (b *MemBus) Subscribe(ctx context.Context, typeFilter string) (<-chan Event, error) {
	ch := make(chan Event, 64)
	b.mu.Lock()
	b.subs = append(b.subs, subscription{filter: typeFilter, ch: ch})
	b.mu.Unlock()
	go func() {
		<-ctx.Done()
		b.mu.Lock()
		defer b.mu.Unlock()
		for i, s := range b.subs {
			if s.ch == ch {
				b.subs = append(b.subs[:i], b.subs[i+1:]...)
				close(ch)
				break
			}
		}
	}()
	return ch, nil
}
