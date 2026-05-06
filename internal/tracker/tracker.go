package tracker

import (
	"sync"
	"time"
)

// State holds the last known status of a cron job.
type State struct {
	Name        string    `json:"name"`
	LastSuccess time.Time `json:"last_success"`
	LastFailure time.Time `json:"last_failure"`
	MissCount   int       `json:"miss_count"`
	FailCount   int       `json:"fail_count"`
}

// Tracker stores job states in memory.
type Tracker struct {
	mu     sync.RWMutex
	states map[string]*State
}

// New creates an empty Tracker.
func New() *Tracker {
	return &Tracker{states: make(map[string]*State)}
}

func (t *Tracker) getOrCreate(name string) *State {
	if s, ok := t.states[name]; ok {
		return s
	}
	s := &State{Name: name}
	t.states[name] = s
	return s
}

// RecordSuccess marks a successful run at the given time.
func (t *Tracker) RecordSuccess(name string, at time.Time) {
	t.mu.Lock()
	defer t.mu.Unlock()
	s := t.getOrCreate(name)
	s.LastSuccess = at
	s.MissCount = 0
}

// RecordFailure marks a failed run.
func (t *Tracker) RecordFailure(name string, at time.Time) {
	t.mu.Lock()
	defer t.mu.Unlock()
	s := t.getOrCreate(name)
	s.LastFailure = at
	s.FailCount++
}

// RecordMiss increments the miss counter for a job.
func (t *Tracker) RecordMiss(name string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	s := t.getOrCreate(name)
	s.MissCount++
}

// Get returns the state for a job, or false if unknown.
func (t *Tracker) Get(name string) (State, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	s, ok := t.states[name]
	if !ok {
		return State{}, false
	}
	return *s, true
}

// All returns a snapshot of all job states.
func (t *Tracker) All() map[string]State {
	t.mu.RLock()
	defer t.mu.RUnlock()
	out := make(map[string]State, len(t.states))
	for k, v := range t.states {
		out[k] = *v
	}
	return out
}
