package replay_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/your-org/cronwatch/internal/config"
	"github.com/your-org/cronwatch/internal/replay"
	"github.com/your-org/cronwatch/internal/tracker"
)

type mockDispatcher struct {
	called []string
	err    error
}

func (m *mockDispatcher) Dispatch(_ context.Context, name, _ string) error {
	m.called = append(m.called, name)
	return m.err
}

func newReplayer(jobs []config.Job, t *tracker.Tracker, d replay.Dispatcher) *replay.Replayer {
	cfg := &config.Config{Jobs: jobs}
	r := replay.New(cfg, t, d)
	return r
}

func TestReplay_DispatchesMissedJob(t *testing.T) {
	tr := tracker.New()
	disp := &mockDispatcher{}

	jobs := []config.Job{{Name: "backup", Schedule: "* * * * *", Command: "backup.sh"}}
	r := newReplayer(jobs, tr, disp)

	results := r.Run(context.Background())

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Job != "backup" {
		t.Errorf("expected job 'backup', got %q", results[0].Job)
	}
	if results[0].Error != nil {
		t.Errorf("unexpected error: %v", results[0].Error)
	}
	if len(disp.called) != 1 || disp.called[0] != "backup" {
		t.Errorf("expected dispatcher called with 'backup', got %v", disp.called)
	}
}

func TestReplay_SkipsRecentlySucceededJob(t *testing.T) {
	tr := tracker.New()
	// Record a very recent success so the job is not considered missed.
	tr.RecordSuccess("backup", time.Now())

	disp := &mockDispatcher{}
	jobs := []config.Job{{Name: "backup", Schedule: "* * * * *", Command: "backup.sh"}}
	r := newReplayer(jobs, tr, disp)

	results := r.Run(context.Background())

	if len(results) != 0 {
		t.Errorf("expected no replays, got %d", len(results))
	}
	if len(disp.called) != 0 {
		t.Errorf("dispatcher should not have been called, got %v", disp.called)
	}
}

func TestReplay_PropagatesDispatchError(t *testing.T) {
	tr := tracker.New()
	disp := &mockDispatcher{err: errors.New("exec failed")}

	jobs := []config.Job{{Name: "cleanup", Schedule: "* * * * *", Command: "clean.sh"}}
	r := newReplayer(jobs, tr, disp)

	results := r.Run(context.Background())

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Error == nil {
		t.Error("expected error, got nil")
	}
}

func TestReplay_SkipsInvalidSchedule(t *testing.T) {
	tr := tracker.New()
	disp := &mockDispatcher{}

	jobs := []config.Job{{Name: "bad", Schedule: "not-a-cron", Command: "cmd"}}
	r := newReplayer(jobs, tr, disp)

	results := r.Run(context.Background())

	if len(results) != 0 {
		t.Errorf("expected 0 results for invalid schedule, got %d", len(results))
	}
}
