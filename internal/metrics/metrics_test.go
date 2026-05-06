package metrics

import (
	"testing"
	"time"
)

func TestRecordRun_Success(t *testing.T) {
	c := New()
	c.RecordRun("backup", true, 2*time.Second)

	m, ok := c.Get("backup")
	if !ok {
		t.Fatal("expected metrics for 'backup'")
	}
	if m.TotalRuns != 1 {
		t.Errorf("TotalRuns: got %d, want 1", m.TotalRuns)
	}
	if m.SuccessCount != 1 {
		t.Errorf("SuccessCount: got %d, want 1", m.SuccessCount)
	}
	if m.FailureCount != 0 {
		t.Errorf("FailureCount: got %d, want 0", m.FailureCount)
	}
	if m.LastDuration != 2*time.Second {
		t.Errorf("LastDuration: got %v, want 2s", m.LastDuration)
	}
}

func TestRecordRun_Failure(t *testing.T) {
	c := New()
	c.RecordRun("sync", false, 500*time.Millisecond)

	m, _ := c.Get("sync")
	if m.FailureCount != 1 {
		t.Errorf("FailureCount: got %d, want 1", m.FailureCount)
	}
	if m.SuccessCount != 0 {
		t.Errorf("SuccessCount: got %d, want 0", m.SuccessCount)
	}
}

func TestRecordRun_AvgDuration(t *testing.T) {
	c := New()
	c.RecordRun("job", true, 2*time.Second)
	c.RecordRun("job", true, 4*time.Second)

	m, _ := c.Get("job")
	if m.AvgDuration != 3*time.Second {
		t.Errorf("AvgDuration: got %v, want 3s", m.AvgDuration)
	}
	if m.TotalRuns != 2 {
		t.Errorf("TotalRuns: got %d, want 2", m.TotalRuns)
	}
}

func TestRecordMiss(t *testing.T) {
	c := New()
	c.RecordMiss("nightly")
	c.RecordMiss("nightly")

	m, ok := c.Get("nightly")
	if !ok {
		t.Fatal("expected metrics for 'nightly'")
	}
	if m.MissCount != 2 {
		t.Errorf("MissCount: got %d, want 2", m.MissCount)
	}
}

func TestGet_UnknownJob(t *testing.T) {
	c := New()
	_, ok := c.Get("ghost")
	if ok {
		t.Error("expected false for unknown job")
	}
}

func TestAll_ReturnsAllJobs(t *testing.T) {
	c := New()
	c.RecordRun("alpha", true, time.Second)
	c.RecordRun("beta", false, time.Second)
	c.RecordMiss("gamma")

	all := c.All()
	if len(all) != 3 {
		t.Errorf("All: got %d entries, want 3", len(all))
	}
}
