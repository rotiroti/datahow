package datahow

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rotiroti/datahow/uniq"
)

// LogServer represents an HTTP server that handles incoming log records.
type LogServer struct {
	mux    *http.ServeMux
	store  *uniq.InMemory
	metric prometheus.Counter
}

// NewLogServer creates a new LogServer with the provided InMemory storage.
func NewLogServer(im *uniq.InMemory, c prometheus.Counter) *LogServer {
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
