package tracker

import (
	"sync"
	"time"
)

// Status represents the last known state of a monitored job.
type Status int

const (
	StatusUnknown Status = iota
	StatusOK
	StatusMissed
	StatusFailed
)

func (s Status) String() string {
	switch s {
	case StatusOK:
		return "ok"
	case StatusMissed:
		return "missed"
	case StatusFailed:
		return "failed"
	default:
		return "unknown"
	}
}

// JobState holds runtime state for a single cron job.
type JobState struct {
	LastSeen time.Time
	LastStatus Status
	ConsecutiveMisses int
}

// Tracker maintains in-memory state for all monitored jobs.
type Tracker struct {
	mu    sync.RWMutex
	states map[string]*JobState
}

// New creates a new Tracker.
func New() *Tracker {
	return &Tracker{
		states: make(map[string]*JobState),
	}
}

// RecordSuccess marks a job as successfully completed.
func (t *Tracker) RecordSuccess(name string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	s := t.getOrCreate(name)
	s.LastSeen = time.Now()
	s.LastStatus = StatusOK
	s.ConsecutiveMisses = 0
}

// RecordFailure marks a job run as failed.
func (t *Tracker) RecordFailure(name string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	s := t.getOrCreate(name)
	s.LastSeen = time.Now()
	s.LastStatus = StatusFailed
}

// RecordMiss increments the missed-run counter for a job.
func (t *Tracker) RecordMiss(name string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	s := t.getOrCreate(name)
	s.LastStatus = StatusMissed
	s.ConsecutiveMisses++
}

// Get returns a copy of the current state for a job.
func (t *Tracker) Get(name string) (JobState, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	s, ok := t.states[name]
	if !ok {
		return JobState{}, false
	}
	return *s, true
}

func (t *Tracker) getOrCreate(name string) *JobState {
	if _, ok := t.states[name]; !ok {
		t.states[name] = &JobState{}
	}
	return t.states[name]
}
