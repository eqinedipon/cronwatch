package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cronwatch/internal/config"
	"github.com/cronwatch/internal/tracker"
)

func TestHandleListJobs_Empty(t *testing.T) {
	srv := newTestServer(t)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/jobs", nil)
	srv.router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var jobs []jobResponse
	if err := json.NewDecoder(rec.Body).Decode(&jobs); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(jobs) != 0 {
		t.Errorf("expected empty list, got %d jobs", len(jobs))
	}
}

func TestHandleListJobs_WithJobs(t *testing.T) {
	srv := newTestServer(t)
	srv.cfg.Jobs = []config.Job{
		{Name: "backup", Schedule: "0 2 * * *"},
	}
	srv.tracker.RecordSuccess("backup", time.Now())
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/jobs", nil)
	srv.router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var jobs []jobResponse
	if err := json.NewDecoder(rec.Body).Decode(&jobs); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(jobs) != 1 {
		t.Fatalf("expected 1 job, got %d", len(jobs))
	}
	if jobs[0].Name != "backup" {
		t.Errorf("expected name backup, got %s", jobs[0].Name)
	}
	if jobs[0].Status != "ok" {
		t.Errorf("expected status ok, got %s", jobs[0].Status)
	}
}

func TestHandleGetJob_NotFound(t *testing.T) {
	srv := newTestServer(t)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/jobs/ghost", nil)
	srv.router.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestHandleGetJob_Found(t *testing.T) {
	srv := newTestServer(t)
	srv.cfg.Jobs = []config.Job{
		{Name: "sync", Schedule: "*/5 * * * *"},
	}
	srv.tracker.RecordMiss("sync")
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/jobs/sync", nil)
	srv.router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var job jobResponse
	if err := json.NewDecoder(rec.Body).Decode(&job); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if job.MissCount != 1 {
		t.Errorf("expected miss_count 1, got %d", job.MissCount)
	}
	_ = tracker.New // ensure import used
}
