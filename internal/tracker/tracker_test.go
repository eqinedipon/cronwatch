package tracker_test

import (
	"testing"
	"time"

	"github.com/cronwatch/cronwatch/internal/tracker"
)

func TestRecordSuccess(t *testing.T) {
	tr := tracker.New()
	before := time.Now()
	tr.RecordSuccess("backup")

	state, ok := tr.Get("backup")
	if !ok {
		t.Fatal("expected state to exist after RecordSuccess")
	}
	if state.LastStatus != tracker.StatusOK {
		t.Errorf("expected StatusOK, got %v", state.LastStatus)
	}
	if state.ConsecutiveMisses != 0 {
		t.Errorf("expected 0 consecutive misses, got %d", state.ConsecutiveMisses)
	}
	if state.LastSeen.Before(before) {
		t.Error("LastSeen should be set to approximately now")
	}
}

func TestRecordFailure(t *testing.T) {
	tr := tracker.New()
	tr.RecordFailure("cleanup")

	state, ok := tr.Get("cleanup")
	if !ok {
		t.Fatal("expected state to exist after RecordFailure")
	}
	if state.LastStatus != tracker.StatusFailed {
		t.Errorf("expected StatusFailed, got %v", state.LastStatus)
	}
}

func TestRecordMiss_Increments(t *testing.T) {
	tr := tracker.New()
	tr.RecordMiss("report")
	tr.RecordMiss("report")

	state, _ := tr.Get("report")
	if state.ConsecutiveMisses != 2 {
		t.Errorf("expected 2 consecutive misses, got %d", state.ConsecutiveMisses)
	}
	if state.LastStatus != tracker.StatusMissed {
		t.Errorf("expected StatusMissed, got %v", state.LastStatus)
	}
}

func TestRecordSuccess_ResetsMisses(t *testing.T) {
	tr := tracker.New()
	tr.RecordMiss("sync")
	tr.RecordMiss("sync")
	tr.RecordSuccess("sync")

	state, _ := tr.Get("sync")
	if state.ConsecutiveMisses != 0 {
		t.Errorf("expected misses reset to 0, got %d", state.ConsecutiveMisses)
	}
}

func TestGet_UnknownJob(t *testing.T) {
	tr := tracker.New()
	_, ok := tr.Get("nonexistent")
	if ok {
		t.Error("expected false for unknown job")
	}
}

func TestStatusString(t *testing.T) {
	cases := []struct {
		status tracker.Status
		want   string
	}{
		{tracker.StatusOK, "ok"},
		{tracker.StatusMissed, "missed"},
		{tracker.StatusFailed, "failed"},
		{tracker.StatusUnknown, "unknown"},
	}
	for _, c := range cases {
		if got := c.status.String(); got != c.want {
			t.Errorf("Status(%d).String() = %q, want %q", c.status, got, c.want)
		}
	}
}
