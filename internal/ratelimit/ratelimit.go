// Package ratelimit provides per-job alert rate limiting to prevent
// alert storms when a job repeatedly fails or is missed.
package ratelimit

import (
	"sync"
	"time"
)

// Limiter tracks the last alert time per job and suppresses duplicate
// alerts that arrive within a configurable cooldown window.
type Limiter struct {
	mu       sync.Mutex
	cooldown time.Duration
	lastSent map[string]time.Time
	now      func() time.Time
}

// New creates a Limiter with the given cooldown duration.
// Alerts for the same job key will be suppressed until cooldown elapses.
func New(cooldown time.Duration) *Limiter {
	return &Limiter{
		cooldown: cooldown,
		lastSent: make(map[string]time.Time),
		now:      time.Now,
	}
}

// Allow returns true if an alert for the given key should be sent.
// It records the current time as the last-sent time when it returns true.
func (l *Limiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.now()
	if last, ok := l.lastSent[key]; ok {
		if now.Sub(last) < l.cooldown {
			return false
		}
	}
	l.lastSent[key] = now
	return true
}

// Reset clears the rate-limit record for a key, allowing the next alert
// to be sent immediately regardless of the cooldown window.
func (l *Limiter) Reset(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.lastSent, key)
}

// Remaining returns how much cooldown time is left for a key.
// Returns 0 if the key is not currently suppressed.
func (l *Limiter) Remaining(key string) time.Duration {
	l.mu.Lock()
	defer l.mu.Unlock()

	last, ok := l.lastSent[key]
	if !ok {
		return 0
	}
	remaining := l.cooldown - l.now().Sub(last)
	if remaining < 0 {
		return 0
	}
	return remaining
}
