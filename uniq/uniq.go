package uniq

import "sync"

var _ Counter = (*HSet)(nil)

// Counter defines the interface for counting unique elements.
type Counter interface {
	Add(ip string) bool
	Count() int64
}

// HSet represents a concurrent-safe set of IP addresses.
// It uses a map internally to store unique IPs and a separate mutex for synchronization.
type HSet struct {
	u  map[string]struct{}
	mu sync.Mutex
}

// NewHSet initializes a new HSet.
func NewHSet() *HSet {
	return &HSet{u: make(map[string]struct{})}
}

// Add checks if an `ip` address exists in the set and adds it if not present.
// Returns true if the IP was already in the set, false if it was added.
func (h *HSet) Add(ip string) bool {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.u[ip]; ok {
		return true
	}

	h.u[ip] = struct{}{}

	return false
}

// Count returns the number of unique IP addresses in the set.
func (h *HSet) Count() int64 {
	h.mu.Lock()
	defer h.mu.Unlock()

	return int64(len(h.u))
}
