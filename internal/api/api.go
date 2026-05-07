package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/cronwatch/internal/metrics"
	"github.com/cronwatch/internal/tracker"
)

// Server exposes a simple HTTP API for job status and metrics.
type Server struct {
	tracker *tracker.Tracker
	metrics *metrics.Metrics
	server  *http.Server
}

// New creates a new API server bound to addr.
func New(addr string, t *tracker.Tracker, m *metrics.Metrics) *Server {
	s := &Server{tracker: t, metrics: m}
	mux := http.NewServeMux()
	mux.HandleFunc("/api/jobs", s.listJobs)
	mux.HandleFunc("/api/jobs/", s.getJob)
	s.server = &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	return s
}

// Start begins listening in a goroutine.
func (s *Server) Start() error {
	go s.server.ListenAndServe() //nolint:errcheck
	return nil
}

// Stop gracefully shuts down the server with a timeout to allow in-flight
// requests to complete before forcefully closing connections.
func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.server.Shutdown(ctx)
}

func (s *Server) listJobs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	states := s.tracker.All()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(states) //nolint:errcheck
}

func (s *Server) getJob(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	name := r.URL.Path[len("/api/jobs/"):]
	if name == "" {
		http.Error(w, "job name required", http.StatusBadRequest)
		return
	}
	state, ok := s.tracker.Get(name)
	if !ok {
		http.Error(w, "job not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(state) //nolint:errcheck
}
