package suppress_test

import (
	"errors"
	"testing"
	"time"

	"github.com/cronwatch/internal/suppress"
)

type mockSender struct {
	calls []string
	err   error
}

func (m *mockSender) Send(key, message string) error {
	m.calls = append(m.calls, key)
	return m.err
}

func newGuardedSender(windowMinutes int) (*suppress.GuardedSender, *mockSender) {
	inner := &mockSender{}
	w := suppress.New(time.Duration(windowMinutes) * time.Minute)
	gs := suppress.NewGuardedSender(w, inner)
	return gs, inner
}

func TestGuardedSender_FirstCallDelivered(t *testing.T) {
	gs, inner := newGuardedSender(5)
	err := gs.Send("job-a", "first alert")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(inner.calls) != 1 {
		t.Fatalf("expected 1 call, got %d", len(inner.calls))
	}
}

func TestGuardedSender_SecondCallSuppressed(t *testing.T) {
	gs, inner := newGuardedSender(5)
	_ = gs.Send("job-a", "first")
	err := gs.Send("job-a", "second")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(inner.calls) != 1 {
		t.Fatalf("expected 1 call (suppressed second), got %d", len(inner.calls))
	}
}

func TestGuardedSender_DifferentKeysIndependent(t *testing.T) {
	gs, inner := newGuardedSender(5)
	_ = gs.Send("job-a", "msg")
	_ = gs.Send("job-b", "msg")
	if len(inner.calls) != 2 {
		t.Fatalf("expected 2 calls for different keys, got %d", len(inner.calls))
	}
}

func TestGuardedSender_PropagatesInnerError(t *testing.T) {
	inner := &mockSender{err: errors.New("send failed")}
	w := suppress.New(5 * time.Minute)
	gs := suppress.NewGuardedSender(w, inner)
	err := gs.Send("job-a", "msg")
	if err == nil {
		t.Fatal("expected error from inner sender")
	}
}

func TestGuardedSender_AllowedAfterWindowExpires(t *testing.T) {
	inner := &mockSender{}
	w := suppress.New(10 * time.Millisecond)
	gs := suppress.NewGuardedSender(w, inner)
	_ = gs.Send("job-a", "first")
	time.Sleep(20 * time.Millisecond)
	_ = gs.Send("job-a", "second after window")
	if len(inner.calls) != 2 {
		t.Fatalf("expected 2 calls after window expired, got %d", len(inner.calls))
	}
}
