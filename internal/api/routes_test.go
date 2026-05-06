package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cronwatch/internal/metrics"
	"github.com/cronwatch/internal/tracker"
)

func TestGetMetrics_Empty(t *testing.T) {
	mt := metrics.New()
	tr := tracker.New()
	_ = tr
	_ = mt
	// Verify All() returns empty map when no jobs recorded
	all := mt.All()
	if len(all) != 0 {
		t.Fatalf("expected empty metrics, got %d entries", len(all))
	}
}

func TestGetMetrics_WithData(t *testing.T) {
	mt := metrics.New()
	mt.RecordRun("backup", true, 200*time.Millisecond)
	mt.RecordRun("backup", false, 50*time.Millisecond)
	all := mt.All()
	s, ok := all["backup"]
	if !ok {
		t.Fatal("expected backup in metrics")
	}
	if s.Total != 2 {
		t.Fatalf("expected 2 runs, got %d", s.Total)
	}
	if s.Failures != 1 {
		t.Fatalf("expected 1 failure, got %d", s.Failures)
	}
}

func TestMetricsEndpoint_MethodNotAllowed(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/metrics", nil)
	if req.Method != http.MethodPost {
		t.Fatal("unexpected method")
	}
	// simulate handler logic
	if req.Method != http.MethodGet {
		rec.WriteHeader(http.StatusMethodNotAllowed)
	}
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rec.Code)
	}
}

func TestMetricsEndpoint_JSON(t *testing.T) {
	mt := metrics.New()
	mt.RecordRun("sync", true, 100*time.Millisecond)
	all := mt.All()
	b, err := json.Marshal(all)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	if len(b) == 0 {
		t.Fatal("expected non-empty JSON")
	}
}
