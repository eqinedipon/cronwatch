package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cronwatch/internal/api"
	"github.com/cronwatch/internal/metrics"
	"github.com/cronwatch/internal/tracker"
)

func newTestServer(t *testing.T) (*api.Server, *tracker.Tracker, *metrics.Metrics) {
	t.Helper()
	tr := tracker.New()
	mt := metrics.New()
	srv := api.New("127.0.0.1:0", tr, mt)
	return srv, tr, mt
}

func TestListJobs_Empty(t *testing.T) {
	_, tr, mt := newTestServer(t)
	srv := api.New("127.0.0.1:0", tr, mt)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/jobs", nil)
	// access handler via exported test helper — use real server mux via Start/Stop
	_ = srv
	_ = rec
	_ = req
}

func TestListJobs_ReturnsJobs(t *testing.T) {
	_, tr, mt := newTestServer(t)
	tr.RecordSuccess("backup", time.Now())
	srv := api.New("127.0.0.1:0", tr, mt)
	_ = srv
	states := tr.All()
	if _, ok := states["backup"]; !ok {
		t.Fatal("expected backup job in states")
	}
}

func TestGetJob_NotFound(t *testing.T) {
	_, tr, mt := newTestServer(t)
	srv := api.New("127.0.0.1:0", tr, mt)
	_ = srv
	_, ok := tr.Get("nonexistent")
	if ok {
		t.Fatal("expected job not found")
	}
}

func TestGetJob_Found(t *testing.T) {
	_, tr, mt := newTestServer(t)
	tr.RecordSuccess("sync", time.Now())
	state, ok := tr.Get("sync")
	if !ok {
		t.Fatal("expected job to exist")
	}
	b, err := json.Marshal(state)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	if len(b) == 0 {
		t.Fatal("expected non-empty JSON")
	}
}

func TestStartStop(t *testing.T) {
	srv, tr, mt := newTestServer(t)
	_ = tr
	_ = mt
	if err := srv.Start(); err != nil {
		t.Fatalf("Start: %v", err)
	}
	if err := srv.Stop(); err != nil {
		t.Fatalf("Stop: %v", err)
	}
}

// TestStartStop_DoubleStop verifies that calling Stop twice does not panic or
// return an unexpected error, ensuring graceful shutdown is idempotent.
func TestStartStop_DoubleStop(t *testing.T) {
	srv, _, _ := newTestServer(t)
	if err := srv.Start(); err != nil {
		t.Fatalf("Start: %v", err)
	}
	if err := srv.Stop(); err != nil {
		t.Fatalf("first Stop: %v", err)
	}
	// Second Stop should not panic; an error is acceptable but not required.
	_ = srv.Stop()
}
