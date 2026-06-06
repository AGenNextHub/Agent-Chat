package capability

import (
	"fmt"
	"sort"
	"sync"
)

// Registry is the admitted set of capability contracts. Registration rejects
// invalid contracts, so the registry only ever holds runnable, contract-valid
// capabilities. It is goroutine-safe.
type Registry struct {
	mu sync.RWMutex
	m  map[string]Contract
}

// NewRegistry returns an empty registry.
func NewRegistry() *Registry { return &Registry{m: make(map[string]Contract)} }

// Register validates and admits a contract. A contract whose invariants or
// dependencies are unmet is rejected, never partially admitted.
func (r *Registry) Register(c Contract) error {
	if err := c.Validate(); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, dep := range c.Requires {
		if _, ok := r.m[dep]; !ok {
			return fmt.Errorf("%w: %q requires unregistered capability %q", ErrInvalidContract, c.Name, dep)
		}
	}
	r.m[c.Name] = c
	return nil
}

// Get returns the contract for name and whether it is registered.
func (r *Registry) Get(name string) (Contract, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	c, ok := r.m[name]
	return c, ok
}

// List returns all registered contracts sorted by name (for inspectability).
func (r *Registry) List() []Contract {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]Contract, 0, len(r.m))
	for _, c := range r.m {
		out = append(out, c)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}
