package digest_test

import (
	"strings"
	"testing"
	"time"

	"github.com/user/cronwatch/internal/digest"
	"github.com/user/cronwatch/internal/metrics"
	"github.com/user/cronwatch/internal/tracker"
)

func newReporter(names []string) *digest.Reporter {
	m := metrics.New()
	t := tracker.New()
	for _, n := range names {
		t.RecordSuccess(n)
	}
	return digest.New(m, t, names)
}

func TestBuild_EmptyJobs(t *testing.T) {
	r := newReporter(nil)
	s := r.Build(time.Hour)
	if len(s.Jobs) != 0 {
		t.Fatalf("expected 0 jobs, got %d", len(s.Jobs))
	}
	if s.Period != time.Hour {
		t.Errorf("expected period %s, got %s", time.Hour, s.Period)
	}
}

func TestBuild_OKJob(t *testing.T) {
	m := metrics.New()
	tr := tracker.New()
	tr.RecordSuccess("backup")
	m.RecordRun("backup", true, 200*time.Millisecond)

	r := digest.New(m, tr, []string{"backup"})
	s := r.Build(24 * time.Hour)

	if len(s.Jobs) != 1 {
		t.Fatalf("expected 1 job")
	}
	j := s.Jobs[0]
	if j.Name != "backup" {
		t.Errorf("unexpected name %q", j.Name)
	}
	if j.Status != "ok" {
		t.Errorf("expected status ok, got %q", j.Status)
	}
	if j.TotalRuns != 1 {
		t.Errorf("expected 1 run, got %d", j.TotalRuns)
	}
}

func TestBuild_MissingJob(t *testing.T) {
	m := metrics.New()
	tr := tracker.New()
	tr.RecordSuccess("sync")
	tr.RecordMiss("sync")

	r := digest.New(m, tr, []string{"sync"})
	s := r.Build(time.Hour)

	if s.Jobs[0].Status != "missing" {
		t.Errorf("expected missing, got %q", s.Jobs[0].Status)
	}
}

func TestBuild_DegradedJob(t *testing.T) {
	m := metrics.New()
	tr := tracker.New()
	tr.RecordSuccess("report")
	m.RecordRun("report", false, 50*time.Millisecond)

	r := digest.New(m, tr, []string{"report"})
	s := r.Build(time.Hour)

	if s.Jobs[0].Status != "degraded" {
		t.Errorf("expected degraded, got %q", s.Jobs[0].Status)
	}
}

func TestFormat_ContainsJobName(t *testing.T) {
	m := metrics.New()
	tr := tracker.New()
	tr.RecordSuccess("cleanup")

	r := digest.New(m, tr, []string{"cleanup"})
	s := r.Build(time.Hour)
	out := digest.Format(s)

	if !strings.Contains(out, "cleanup") {
		t.Errorf("expected output to contain job name, got:\n%s", out)
	}
	if !strings.Contains(out, "cronwatch digest") {
		t.Errorf("expected header in output")
	}
}
