package api

import (
	"encoding/json"
	"net/http"

	"github.com/cronwatch/internal/metrics"
)

// metricsResponse is the JSON shape returned by GET /api/metrics.
type metricsResponse struct {
	Jobs map[string]metrics.Stats `json:"jobs"`
}

// RegisterMetricsRoute adds the /api/metrics endpoint to the given mux.
func (s *Server) registerMetricsRoute(mux *http.ServeMux) {
	mux.HandleFunc("/api/metrics", s.getMetrics)
}

func (s *Server) getMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	all := s.metrics.All()
	resp := metricsResponse{Jobs: all}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "encode error", http.StatusInternalServerError)
	}
}
