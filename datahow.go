package datahow

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// A InMemory represents a concurrent-safe set of IP addresses.
// It uses a map internally to store unique IPs.
type InMemory struct {
	u  map[string]struct{}
	mu sync.Mutex
}

// New initializes a new InMemory set.
func NewInMemory() *InMemory {
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

// LogServer represents an HTTP server that handles incoming log records.
type LogServer struct {
	mux    *http.ServeMux
	store  *InMemory
	metric prometheus.Counter
}

// NewLogServer creates a new LogServer with the provided InMemory storage.
func NewLogServer(im *InMemory, c prometheus.Counter) *LogServer {
	srv := &LogServer{
		store:  im,
		metric: c,
	}
	srv.routes()

	return srv
}

func (l *LogServer) routes() {
	l.mux = http.NewServeMux()
	l.mux.HandleFunc("POST /logs", l.handleLog())
}

func (l *LogServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	l.mux.ServeHTTP(w, r)
}

func (l *LogServer) handleLog() http.HandlerFunc {
	type record struct {
		Timestamp time.Time `json:"timestamp"`
		IPAddress string    `json:"ip"`
		URL       string    `json:"url"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var jsonRecord record

		if err := json.NewDecoder(r.Body).Decode(&jsonRecord); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}

		// Increment Prometheus metric only if the IP is new
		if jsonRecord.IPAddress != "" && !l.store.ExistOrAdd(jsonRecord.IPAddress) {
			l.metric.Inc()
		}

		w.WriteHeader(http.StatusOK)
	}
}
