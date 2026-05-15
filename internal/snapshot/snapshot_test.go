package snapshot_test

import (
	"testing"
	"time"

	"github.com/cronwatch/internal/snapshot"
	"github.com/cronwatch/internal/tracker"
)

// fakeTracker implements snapshot.Tracker for testing.
type fakeTracker struct {
	states map[string]tracker.State
}

func (f *fakeTracker) All() map[string]tracker.State { return f.states }

func newTracker(states map[string]tracker.State) *fakeTracker {
	return &fakeTracker{states: states}
}

func TestCapture_EmptyTracker(t *testing.T) {
	s := snapshot.New(newTracker(nil))
	snap := s.Capture()
	if len(snap.Jobs) != 0 {
		t.Fatalf("expected 0 jobs, got %d", len(snap.Jobs))
	}
	if snap.CapturedAt.IsZero() {
		t.Fatal("expected non-zero CapturedAt")
	}
}

func TestCapture_OKJob(t *testing.T) {
	now := time.Now().UTC()
	states := map[string]tracker.State{
		"backup": {LastSuccess: now, ConsecutiveMisses: 0, LastExitCode: 0},
	}
	s := snapshot.New(newTracker(states))
	snap := s.Capture()

	if j, ok := snap.Jobs["backup"]; !ok {
		t.Fatal("expected job 'backup' in snapshot")
	} else if j.LastStatus != "ok" {
		t.Errorf("expected status ok, got %s", j.LastStatus)
	}
}

func TestCapture_MissingJob(t *testing.T) {
	states := map[string]tracker.State{
		"sync": {ConsecutiveMisses: 3, LastExitCode: 0},
	}
	s := snapshot.New(newTracker(states))
	snap := s.Capture()

	if j := snap.Jobs["sync"]; j.LastStatus != "missing" {
		t.Errorf("expected status missing, got %s", j.LastStatus)
	}
	if snap.Jobs["sync"].MissCount != 3 {
		t.Errorf("expected MissCount 3, got %d", snap.Jobs["sync"].MissCount)
	}
}

func TestCapture_FailedJob(t *testing.T) {
	states := map[string]tracker.State{
		"report": {ConsecutiveMisses: 0, LastExitCode: 1},
	}
	s := snapshot.New(newTracker(states))
	snap := s.Capture()

	if j := snap.Jobs["report"]; j.LastStatus != "failed" {
		t.Errorf("expected status failed, got %s", j.LastStatus)
	}
}

func TestLatest_NilBeforeCapture(t *testing.T) {
	s := snapshot.New(newTracker(nil))
	if s.Latest() != nil {
		t.Fatal("expected nil before first capture")
	}
}

func TestLatest_ReturnsMostRecent(t *testing.T) {
	states := map[string]tracker.State{
		"job1": {LastExitCode: 0},
	}
	s := snapshot.New(newTracker(states))
	s.Capture()
	l := s.Latest()
	if l == nil {
		t.Fatal("expected non-nil latest after capture")
	}
	if _, ok := l.Jobs["job1"]; !ok {
		t.Error("expected job1 in latest snapshot")
	}
}
