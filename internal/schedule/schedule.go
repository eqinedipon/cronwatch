package schedule

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

// NextRun returns the next scheduled run time after 'from' for the given cron expression.
func NextRun(expr string, from time.Time) (time.Time, error) {
	p := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	sched, err := p.Parse(expr)
	if err != nil {
		return time.Time{}, fmt.Errorf("parse cron expression %q: %w", expr, err)
	}
	return sched.Next(from), nil
}

// PrevRun returns the most recent scheduled run time at or before 'from'.
// It works by stepping back in time until it finds the previous tick.
func PrevRun(expr string, from time.Time) (time.Time, error) {
	p := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	sched, err := p.Parse(expr)
	if err != nil {
		return time.Time{}, fmt.Errorf("parse cron expression %q: %w", expr, err)
	}

	// Step back one minute at a time to find the previous run.
	// We search up to 366 days back to handle edge cases.
	candidate := from.Add(-time.Minute)
	for i := 0; i < 366*24*60; i++ {
		next := sched.Next(candidate)
		if !next.After(from) {
			return next, nil
		}
		candidate = candidate.Add(-time.Minute)
	}
	return time.Time{}, fmt.Errorf("could not determine previous run for %q", expr)
}

// IsDue reports whether a cron expression was due within the window ending at 'at'
// with the given grace period.
func IsDue(expr string, at time.Time, grace time.Duration) (bool, error) {
	prev, err := PrevRun(expr, at)
	if err != nil {
		return false, err
	}
	return at.Sub(prev) <= grace, nil
}
