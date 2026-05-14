// Package suppress provides a mechanism to suppress duplicate alerts
// for a job within a configurable time window.
package suppress

import (
	"sync"
	"time"
)

// Window holds suppression state for all tracked keys.
type Window struct {
	mu       sync.Mutex
	records  map[string]time.Time
	duration time.Duration
}

// New creates a Window that suppresses repeated events within d.
func New(d time.Duration) *Window {
	if d <= 0 {
		d = 30 * time.Minute
	}
	return &Window{
		records:  make(map[string]time.Time),
		duration: d,
	}
}

// IsSuppressed reports whether key has been seen within the window.
// If it has not been seen (or the window has expired), it records the
// current time and returns false so the caller may proceed.
func (w *Window) IsSuppressed(key string) bool {
	w.mu.Lock()
	defer w.mu.Unlock()

	now := time.Now()
	if last, ok := w.records[key]; ok && now.Sub(last) < w.duration {
		return true
	}
	w.records[key] = now
	return false
}

// Reset clears suppression state for key, allowing the next event
// through regardless of when it arrives.
func (w *Window) Reset(key string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	delete(w.records, key)
}

// Duration returns the configured suppression window.
func (w *Window) Duration() time.Duration {
	return w.duration
}
