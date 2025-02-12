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
	mux         *http.ServeMux
	counter     uniq.Counter
	promCounter prometheus.Counter
}

// NewLogServer creates a new LogServer with the provided InMemory storage.
func NewLogServer(c uniq.Counter, promc prometheus.Counter) *LogServer {
	srv := &LogServer{
		counter:     c,
		promCounter: promc,
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
		if jsonRecord.IPAddress != "" && !l.counter.Add(jsonRecord.IPAddress) {
			l.promCounter.Inc()
		}

		w.WriteHeader(http.StatusOK)
	}
}
