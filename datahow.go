package datahow

import "sync"

// A InMemory represents a concurrent-safe set of IP addresses.
// It uses a map internally to store unique IPs.
type InMemory struct {
	u  map[string]struct{}
	mu sync.Mutex
}

// New initializes a new InMemory set.
func New() *InMemory {
	return &InMemory{u: make(map[string]struct{})}
}

// ExistOrAdd checks if an `ip` address exists in the set and adds it if not present.
// Returns true if the IP was already in the set, false if it was added.
func (i *InMemory) ExistOrAdd(ip string) bool {
	i.mu.Lock()
	defer i.mu.Unlock()

	if _, ok := i.u[ip]; ok {
		return true
	}

	i.u[ip] = struct{}{}

	return false
}
