package circuitbreaker

import (
	"testing"
	"time"
)

func TestNew_DefaultsApplied(t *testing.T) {
	b := New(0, 0)
	if b.threshold != 3 {
		t.Fatalf("expected default threshold 3, got %d", b.threshold)
	}
	if b.cooldown != 30*time.Second {
		t.Fatalf("expected default cooldown 30s, got %v", b.cooldown)
	}
}

func TestAllow_ClosedByDefault(t *testing.T) {
	b := New(3, time.Second)
	if !b.Allow() {
		t.Fatal("expected Allow() true for closed breaker")
	}
}

func TestRecordFailure_OpensAtThreshold(t *testing.T) {
	b := New(3, time.Minute)
	b.RecordFailure()
	b.RecordFailure()
	if b.State() != StateClosed {
		t.Fatal("expected still closed after 2 failures")
	}
	b.RecordFailure()
	if b.State() != StateOpen {
		t.Fatalf("expected open after threshold, got %s", b.State())
	}
}

func TestAllow_ReturnsFalseWhenOpen(t *testing.T) {
	b := New(1, time.Minute)
	b.RecordFailure()
	if b.Allow() {
		t.Fatal("expected Allow() false when open")
	}
}

func TestAllow_HalfOpenAfterCooldown(t *testing.T) {
	b := New(1, 10*time.Millisecond)
	b.RecordFailure()
	time.Sleep(20 * time.Millisecond)
	if !b.Allow() {
		t.Fatal("expected Allow() true after cooldown (half-open)")
	}
	if b.State() != StateHalfOpen {
		t.Fatalf("expected half-open, got %s", b.State())
	}
}

func TestRecordSuccess_ClosesBreakerFromHalfOpen(t *testing.T) {
	b := New(1, 10*time.Millisecond)
	b.RecordFailure()
	time.Sleep(20 * time.Millisecond)
	b.Allow() // transitions to half-open
	b.RecordSuccess()
	if b.State() != StateClosed {
		t.Fatalf("expected closed after success, got %s", b.State())
	}
}

func TestRecordSuccess_ResetFailureCount(t *testing.T) {
	b := New(3, time.Minute)
	b.RecordFailure()
	b.RecordFailure()
	b.RecordSuccess()
	// After success, two more failures should not open (need 3 again)
	b.RecordFailure()
	b.RecordFailure()
	if b.State() != StateClosed {
		t.Fatal("expected still closed after reset + 2 failures")
	}
}

func TestState_String(t *testing.T) {
	for _, tc := range []struct {
		s    State
		want string
	}{
		{StateClosed, "closed"},
		{StateOpen, "open"},
		{StateHalfOpen, "half-open"},
		{State(99), "unknown"},
	} {
		if got := tc.s.String(); got != tc.want {
			t.Errorf("State(%d).String() = %q, want %q", tc.s, got, tc.want)
		}
	}
}
