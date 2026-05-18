// Package digest generates periodic summary reports of cron job health.
package digest

import (
	"fmt"
	"strings"
	"time"

	"github.com/user/cronwatch/internal/metrics"
	"github.com/user/cronwatch/internal/tracker"
)

// Reporter builds and sends a digest summary via an alerter-compatible sender.
type Reporter struct {
	metrics  *metrics.Metrics
	tracker  *tracker.Tracker
	jobNames []string
}

// Summary holds the digest data for a reporting period.
type Summary struct {
	GeneratedAt time.Time
	Period      time.Duration
	Jobs        []JobSummary
}

// JobSummary holds per-job stats for the digest.
type JobSummary struct {
	Name        string
	TotalRuns   int
	Failures    int
	Misses      int
	AvgDuration time.Duration
	Status      string
}

// New creates a Reporter for the given job names.
func New(m *metrics.Metrics, t *tracker.Tracker, jobNames []string) *Reporter {
	return &Reporter{metrics: m, tracker: t, jobNames: jobNames}
}

// Build assembles a Summary snapshot for the given period.
func (r *Reporter) Build(period time.Duration) Summary {
	s := Summary{
		GeneratedAt: time.Now().UTC(),
		Period:      period,
	}
	for _, name := range r.jobNames {
		m, _ := r.metrics.Get(name)
		st, _ := r.tracker.Get(name)
		status := "ok"
		if st.ConsecutiveMisses > 0 {
			status = "missing"
		} else if m.Failures > 0 {
			status = "degraded"
		}
		s.Jobs = append(s.Jobs, JobSummary{
			Name:        name,
			TotalRuns:   m.Successes + m.Failures,
			Failures:    m.Failures,
			Misses:      st.ConsecutiveMisses,
			AvgDuration: m.AvgDuration,
			Status:      status,
		})
	}
	return s
}

// Format renders a Summary as a human-readable text block.
func Format(s Summary) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "=== cronwatch digest (%s) — %s ===\n",
		s.Period.String(), s.GeneratedAt.Format(time.RFC3339))
	for _, j := range s.Jobs {
		fmt.Fprintf(&sb, "  [%s] %s  runs=%d failures=%d misses=%d avg=%s\n",
			j.Status, j.Name, j.TotalRuns, j.Failures, j.Misses, j.AvgDuration.Round(time.Millisecond))
	}
	return sb.String()
}
