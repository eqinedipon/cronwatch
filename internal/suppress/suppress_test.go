package suppress_test

import (
	"testing"
	"time"

	"github.com/cronwatch/cronwatch/internal/suppress"
)

func newWindow(d time.Duration) *suppress.Window {
	return suppress.New(d)
}

func TestIsSuppressed_FirstCallAllowed(t *testing.T) {
	w := newWindow(time.Minute)
	if w.IsSuppressed("job-a") {
		t.Fatal("expected first call to be allowed")
	}
}

func TestIsSuppressed_SecondCallSuppressed(t *testing.T) {
	w := newWindow(time.Minute)
	w.IsSuppressed("job-a")
	if !w.IsSuppressed("job-a") {
		t.Fatal("expected second call within window to be suppressed")
	}
}

func TestIsSuppressed_AllowedAfterWindowExpires(t *testing.T) {
	w := newWindow(10 * time.Millisecond)
	w.IsSuppressed("job-a")
	time.Sleep(20 * time.Millisecond)
	if w.IsSuppressed("job-a") {
		t.Fatal("expected call after window expiry to be allowed")
	}
}

func TestIsSuppressed_IndependentKeys(t *testing.T) {
	w := newWindow(time.Minute)
	w.IsSuppressed("job-a")
	if w.IsSuppressed("job-b") {
		t.Fatal("job-b should not be suppressed by job-a")
	}
}

func TestReset_ClearsState(t *testing.T) {
	w := newWindow(time.Minute)
	w.IsSuppressed("job-a")
	w.Reset("job-a")
	if w.IsSuppressed("job-a") {
		t.Fatal("expected reset to clear suppression")
	}
}

func TestNew_DefaultDuration(t *testing.T) {
	w := suppress.New(0)
	if w.Duration() != 30*time.Minute {
		t.Fatalf("expected default 30m, got %v", w.Duration())
	}
}

func TestNew_NegativeDuration(t *testing.T) {
	w := suppress.New(-5 * time.Second)
	if w.Duration() != 30*time.Minute {
		t.Fatalf("expected default 30m for negative input, got %v", w.Duration())
	}
}
