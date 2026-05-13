package ratelimit

import (
	"testing"
	"time"
)

func newLimiter(cooldown time.Duration) *Limiter {
	l := New(cooldown)
	return l
}

func TestAllow_FirstCallAlwaysAllowed(t *testing.T) {
	l := newLimiter(time.Minute)
	if !l.Allow("job1") {
		t.Fatal("expected first call to be allowed")
	}
}

func TestAllow_SecondCallSuppressed(t *testing.T) {
	l := newLimiter(time.Minute)
	l.Allow("job1")
	if l.Allow("job1") {
		t.Fatal("expected second call within cooldown to be suppressed")
	}
}

func TestAllow_AllowedAfterCooldown(t *testing.T) {
	now := time.Now()
	l := newLimiter(time.Minute)
	l.now = func() time.Time { return now }
	l.Allow("job1")

	// advance past cooldown
	l.now = func() time.Time { return now.Add(2 * time.Minute) }
	if !l.Allow("job1") {
		t.Fatal("expected call after cooldown to be allowed")
	}
}

func TestAllow_IndependentKeys(t *testing.T) {
	l := newLimiter(time.Minute)
	l.Allow("job1")
	if !l.Allow("job2") {
		t.Fatal("expected different key to be allowed independently")
	}
}

func TestReset_ClearsLimit(t *testing.T) {
	l := newLimiter(time.Minute)
	l.Allow("job1")
	l.Reset("job1")
	if !l.Allow("job1") {
		t.Fatal("expected allow after reset")
	}
}

func TestRemaining_ZeroWhenNotSeen(t *testing.T) {
	l := newLimiter(time.Minute)
	if r := l.Remaining("job1"); r != 0 {
		t.Fatalf("expected 0 remaining for unseen key, got %v", r)
	}
}

func TestRemaining_PositiveAfterAllow(t *testing.T) {
	now := time.Now()
	l := newLimiter(time.Minute)
	l.now = func() time.Time { return now }
	l.Allow("job1")

	l.now = func() time.Time { return now.Add(10 * time.Second) }
	if r := l.Remaining("job1"); r != 50*time.Second {
		t.Fatalf("expected 50s remaining, got %v", r)
	}
}

func TestRemaining_ZeroAfterCooldownExpires(t *testing.T) {
	now := time.Now()
	l := newLimiter(time.Minute)
	l.now = func() time.Time { return now }
	l.Allow("job1")

	l.now = func() time.Time { return now.Add(2 * time.Minute) }
	if r := l.Remaining("job1"); r != 0 {
		t.Fatalf("expected 0 remaining after cooldown, got %v", r)
	}
}
