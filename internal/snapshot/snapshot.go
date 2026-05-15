// Package snapshot provides point-in-time state capture for all monitored jobs.
package snapshot

import (
	"sync"
	"time"

	"github.com/cronwatch/internal/tracker"
)

// JobSnapshot holds a frozen view of a single job's state.
type JobSnapshot struct {
	Name      string    `json:"name"`
	LastRun   time.Time `json:"last_run,omitempty"`
	LastStatus string   `json:"last_status"`
	MissCount  int      `json:"miss_count"`
	CapturedAt time.Time `json:"captured_at"`
}

// Snapshot holds a frozen view of all jobs at a point in time.
type Snapshot struct {
	CapturedAt time.Time              `json:"captured_at"`
	Jobs       map[string]JobSnapshot `json:"jobs"`
}

// Tracker is the subset of tracker.Tracker used by the snapshotter.
type Tracker interface {
	All() map[string]tracker.State
}

// Snapshotter captures and stores periodic snapshots of job state.
type Snapshotter struct {
	mu      sync.RWMutex
	tracker Tracker
	latest  *Snapshot
}

// New creates a new Snapshotter backed by the given tracker.
func New(t Tracker) *Snapshotter {
	return &Snapshotter{tracker: t}
}

// Capture takes a new snapshot of all job states and stores it.
func (s *Snapshotter) Capture() Snapshot {
	now := time.Now().UTC()
	all := s.tracker.All()

	jobs := make(map[string]JobSnapshot, len(all))
	for name, state := range all {
		status := "ok"
		if state.ConsecutiveMisses > 0 {
			status = "missing"
		}
		if state.LastExitCode != 0 {
			status = "failed"
		}
		jobs[name] = JobSnapshot{
			Name:       name,
			LastRun:    state.LastSuccess,
			LastStatus: status,
			MissCount:  state.ConsecutiveMisses,
			CapturedAt: now,
		}
	}

	snap := Snapshot{CapturedAt: now, Jobs: jobs}

	s.mu.Lock()
	s.latest = &snap
	s.mu.Unlock()

	return snap
}

// Latest returns the most recent snapshot, or nil if none has been taken.
func (s *Snapshotter) Latest() *Snapshot {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.latest
}
