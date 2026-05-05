package watcher

import (
	"testing"
	"time"

	"github.com/cronwatch/internal/alerter"
	"github.com/cronwatch/internal/config"
	"github.com/cronwatch/internal/tracker"
)

func makeConfig(name, sched string, grace int) *config.Config {
	return &config.Config{
		Jobs: []config.Job{
			{Name: name, Schedule: sched, GraceSeconds: grace},
		},
	}
}

func TestCheck_NoMissWhenRecentSuccess(t *testing.T) {
	cfg := makeConfig("backup", "*/5 * * * *", 60)
	tr := tracker.New()
	log := &alerter.LogSender{}
	w := New(cfg, tr, log)

	now := time.Now()
	// Record a success just now so it is after the previous scheduled run.
	tr.RecordSuccess("backup")

	w.check(now)

	state, err := tr.Get("backup")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if state.Misses != 0 {
		t.Errorf("expected 0 misses, got %d", state.Misses)
	}
}

func TestCheck_MissRecordedWhenOverdue(t *testing.T) {
	cfg := makeConfig("report", "0 9 * * *", 300)
	tr := tracker.New()
	log := &alerter.LogSender{}
	w := New(cfg, tr, log)

	// Seed an old success so the job is known but stale.
	tr.RecordSuccess("report")
	tr.ForceLastSuccess("report", time.Now().Add(-48*time.Hour))

	// Simulate checking well after the 09:00 window plus grace.
	now := time.Now()
	w.check(now)

	state, err := tr.Get("report")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if state.Misses == 0 {
		t.Error("expected at least one miss to be recorded")
	}
}

func TestCheck_UnknownJobSkipped(t *testing.T) {
	cfg := makeConfig("newjob", "*/10 * * * *", 120)
	tr := tracker.New()
	log := &alerter.LogSender{}
	w := New(cfg, tr, log)

	// No RecordSuccess call — job is unknown to tracker.
	w.check(time.Now())

	_, err := tr.Get("newjob")
	if err == nil {
		t.Error("expected error for unknown job, got nil")
	}
}

func TestStartStop(t *testing.T) {
	cfg := makeConfig("ping", "* * * * *", 30)
	tr := tracker.New()
	log := &alerter.LogSender{}
	w := New(cfg, tr, log)

	w.Start(50 * time.Millisecond)
	time.Sleep(120 * time.Millisecond)
	w.Stop()
	// If we reach here without panic the loop started and stopped cleanly.
}
