// Package health exposes a simple HTTP endpoint reporting daemon status.
package health

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"
)

// Status holds a snapshot of daemon health.
type Status struct {
	OK        bool      `json:"ok"`
	Uptime    string    `json:"uptime"`
	JobsTotal int       `json:"jobs_total"`
	CheckedAt time.Time `json:"checked_at"`
}

// Server is a lightweight HTTP health server.
type Server struct {
	port      int
	start     time.Time
	jobsTotal int
	running   atomic.Bool
	server    *http.Server
}

// New creates a new health Server.
func New(port, jobsTotal int) *Server {
	return &Server{
		port:      port,
		start:     time.Now(),
		jobsTotal: jobsTotal,
	}
}

// Start begins serving health checks in the background.
func (s *Server) Start() {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", s.handleHealth)

	s.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.port),
		Handler: mux,
	}

	s.running.Store(true)
	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			// non-fatal: health endpoint is best-effort
			_ = err
		}
	}()
}

// Stop shuts down the health server.
func (s *Server) Stop() {
	if s.running.Load() {
		_ = s.server.Close()
		s.running.Store(false)
	}
}

func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	status := Status{
		OK:        true,
		Uptime:    time.Since(s.start).Round(time.Second).String(),
		JobsTotal: s.jobsTotal,
		CheckedAt: time.Now().UTC(),
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(status)
}
