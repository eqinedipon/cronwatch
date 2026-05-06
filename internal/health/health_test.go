package health_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/example/cronwatch/internal/health"
)

func TestHealthEndpoint_ReturnsOK(t *testing.T) {
	srv := health.New(19876, 5)
	srv.Start()
	defer srv.Stop()

	// give the server a moment to bind
	time.Sleep(50 * time.Millisecond)

	resp, err := http.Get(fmt.Sprintf("http://localhost:%d/healthz", 19876))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var status health.Status
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		t.Fatalf("decode failed: %v", err)
	}

	if !status.OK {
		t.Error("expected ok=true")
	}
	if status.JobsTotal != 5 {
		t.Errorf("expected jobs_total=5, got %d", status.JobsTotal)
	}
	if status.Uptime == "" {
		t.Error("expected non-empty uptime")
	}
	if status.CheckedAt.IsZero() {
		t.Error("expected non-zero checked_at")
	}
}

func TestHealthEndpoint_StopPreventsRequests(t *testing.T) {
	srv := health.New(19877, 2)
	srv.Start()
	time.Sleep(50 * time.Millisecond)
	srv.Stop()
	time.Sleep(50 * time.Millisecond)

	_, err := http.Get(fmt.Sprintf("http://localhost:%d/healthz", 19877))
	if err == nil {
		t.Fatal("expected connection refused after Stop")
	}
}

func TestNew_ZeroJobs(t *testing.T) {
	srv := health.New(19878, 0)
	srv.Start()
	defer srv.Stop()

	time.Sleep(50 * time.Millisecond)

	resp, err := http.Get(fmt.Sprintf("http://localhost:%d/healthz", 19878))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	var status health.Status
	_ = json.NewDecoder(resp.Body).Decode(&status)

	if status.JobsTotal != 0 {
		t.Errorf("expected 0 jobs, got %d", status.JobsTotal)
	}
}
