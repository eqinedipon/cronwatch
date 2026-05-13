package ratelimit_test

import (
	"errors"
	"testing"
	"time"

	"github.com/cronwatch/cronwatch/internal/ratelimit"
)

type mockSender struct {
	calls []string
	err   error
}

func (m *mockSender) Send(msg string) error {
	m.calls = append(m.calls, msg)
	return m.err
}

func newGuardedSender(cooldown time.Duration) (*mockSender, ratelimit.Sender) {
	inner := &mockSender{}
	limiter := ratelimit.New(cooldown)
	guarded := ratelimit.NewGuardedSender(inner, limiter)
	return inner, guarded
}

func TestGuardedSender_FirstCallDelivered(t *testing.T) {
	inner, guarded := newGuardedSender(5 * time.Minute)

	if err := guarded.Send("job.fail"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(inner.calls) != 1 {
		t.Errorf("expected 1 call, got %d", len(inner.calls))
	}
}

func TestGuardedSender_SecondCallSuppressed(t *testing.T) {
	inner, guarded := newGuardedSender(5 * time.Minute)

	_ = guarded.Send("job.fail")
	_ = guarded.Send("job.fail")

	if len(inner.calls) != 1 {
		t.Errorf("expected 1 call (suppressed), got %d", len(inner.calls))
	}
}

func TestGuardedSender_DifferentKeysIndependent(t *testing.T) {
	inner, guarded := newGuardedSender(5 * time.Minute)

	_ = guarded.Send("job.a")
	_ = guarded.Send("job.b")

	if len(inner.calls) != 2 {
		t.Errorf("expected 2 calls for distinct keys, got %d", len(inner.calls))
	}
}

func TestGuardedSender_PropagatesError(t *testing.T) {
	inner := &mockSender{err: errors.New("send failed")}
	limiter := ratelimit.New(5 * time.Minute)
	guarded := ratelimit.NewGuardedSender(inner, limiter)

	err := guarded.Send("job.fail")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGuardedSender_AllowedAfterCooldown(t *testing.T) {
	inner, guarded := newGuardedSender(10 * time.Millisecond)

	_ = guarded.Send("job.fail")
	time.Sleep(20 * time.Millisecond)
	_ = guarded.Send("job.fail")

	if len(inner.calls) != 2 {
		t.Errorf("expected 2 calls after cooldown, got %d", len(inner.calls))
	}
}
