package replay_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/your-org/cronwatch/internal/config"
	"github.com/your-org/cronwatch/internal/replay"
	"github.com/your-org/cronwatch/internal/tracker"
)

func TestHandler_MethodNotAllowed(t *testing.T) {
	tr := tracker.New()
	disp := &mockDispatcher{}
	r := replay.New(&config.Config{}, tr, disp)

	req := httptest.NewRequest(http.MethodGet, "/replay", nil)
	w := httptest.NewRecorder()
	replay.Handler(r)(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", w.Code)
	}
}

func TestHandler_NoMissedJobs(t *testing.T) {
	tr := tracker.New()
	disp := &mockDispatcher{}
	cfg := &config.Config{Jobs: []config.Job{}}
	r := replay.New(cfg, tr, disp)

	req := httptest.NewRequest(http.MethodPost, "/replay", nil)
	w := httptest.NewRecorder()
	replay.Handler(r)(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var resp struct {
		Dispatched []string          `json:"dispatched"`
		Errors     map[string]string `json:"errors"`
	}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(resp.Dispatched) != 0 {
		t.Errorf("expected empty dispatched list, got %v", resp.Dispatched)
	}
}

func TestHandler_MissedJobDispatched(t *testing.T) {
	tr := tracker.New()
	disp := &mockDispatcher{}
	cfg := &config.Config{
		Jobs: []config.Job{
			{Name: "sync", Schedule: "* * * * *", Command: "sync.sh"},
		},
	}
	r := replay.New(cfg, tr, disp)

	req := httptest.NewRequest(http.MethodPost, "/replay", nil)
	w := httptest.NewRecorder()
	replay.Handler(r)(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "sync") {
		t.Errorf("expected 'sync' in response body, got: %s", body)
	}
}

func TestHandler_ContentTypeIsJSON(t *testing.T) {
	tr := tracker.New()
	disp := &mockDispatcher{}
	r := replay.New(&config.Config{}, tr, disp)

	req := httptest.NewRequest(http.MethodPost, "/replay", nil)
	w := httptest.NewRecorder()
	replay.Handler(r)(w, req)

	ct := w.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected application/json, got %q", ct)
	}
}

// ensure mockDispatcher satisfies replay.Dispatcher at compile time
var _ replay.Dispatcher = (*mockDispatcher)(nil)

func init() {
	// verify context plumbing compiles
	_ = context.Background
}
