package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

type jobResponse struct {
	Name        string     `json:"name"`
	Schedule    string     `json:"schedule"`
	LastSuccess *time.Time `json:"last_success,omitempty"`
	LastFailure *time.Time `json:"last_failure,omitempty"`
	MissCount   int        `json:"miss_count"`
	Status      string     `json:"status"`
}

func (s *Server) handleListJobs(w http.ResponseWriter, r *http.Request) {
	jobs := s.cfg.Jobs
	resp := make([]jobResponse, 0, len(jobs))
	for _, j := range jobs {
		state, err := s.tracker.Get(j.Name)
		if err != nil {
			continue
		}
		resp = append(resp, jobResponse{
			Name:        j.Name,
			Schedule:    j.Schedule,
			LastSuccess: nilIfZero(state.LastSuccess),
			LastFailure: nilIfZero(state.LastFailure),
			MissCount:   state.MissCount,
			Status:      deriveStatus(state),
		})
	}
	json.NewEncoder(w).Encode(resp)
}

func (s *Server) handleGetJob(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/api/v1/jobs/")
	if name == "" {
		http.Error(w, "missing job name", http.StatusBadRequest)
		return
	}
	state, err := s.tracker.Get(name)
	if err != nil {
		http.Error(w, "job not found", http.StatusNotFound)
		return
	}
	var schedule string
	for _, j := range s.cfg.Jobs {
		if j.Name == name {
			schedule = j.Schedule
			break
		}
	}
	resp := jobResponse{
		Name:        name,
		Schedule:    schedule,
		LastSuccess: nilIfZero(state.LastSuccess),
		LastFailure: nilIfZero(state.LastFailure),
		MissCount:   state.MissCount,
		Status:      deriveStatus(state),
	}
	json.NewEncoder(w).Encode(resp)
}

func nilIfZero(t time.Time) *time.Time {
	if t.IsZero() {
		return nil
	}
	return &t
}

func deriveStatus(state interface{ GetMissCount() int; IsHealthy() bool }) string {
	if state.IsHealthy() {
		return "ok"
	}
	if state.GetMissCount() > 0 {
		return "missing"
	}
	return "failed"
}
