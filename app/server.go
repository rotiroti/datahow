package app

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/rotiroti/datahow/uniq"
)

type MetricCollector interface {
	Inc()
}

// Server represents an HTTP server that handles incoming log records.
type Server struct {
	mux     *http.ServeMux
	counter uniq.Counter
	metrics MetricCollector
}

// NewServer creates a new Server with the provided Counter and Metrics.
func NewServer(c uniq.Counter, m MetricCollector) *Server {
	srv := &Server{
		counter: c,
		metrics: m,
	}
	srv.routes()

	return srv
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) routes() {
	s.mux = http.NewServeMux()
	s.mux.HandleFunc("POST /logs", s.handleLog())
}

func (s *Server) handleLog() http.HandlerFunc {
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
		if jsonRecord.IPAddress != "" && !s.counter.Add(jsonRecord.IPAddress) {
			s.metrics.Inc()
		}

		w.WriteHeader(http.StatusOK)
	}
}
