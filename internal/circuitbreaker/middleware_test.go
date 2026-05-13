package circuitbreaker

import (
	"errors"
	"testing"
	"time"
)

// stubSender is a minimal alerter.Sender for testing.
type stubSender struct {
	calls int
	err   error
}

func (s *stubSender) Send(_, _, _ string) error {
	s.calls++
	return s.err
}

func newGuardedSender(threshold int, cooldown time.Duration, err error) (*GuardedSender, *stubSender, *Breaker) {
	stub := &stubSender{err: err}
	b := New(threshold, cooldown)
	gs := NewGuardedSender(stub, b)
	return gs, stub, b
}

func TestGuardedSender_SuccessForwardedAndBreakerReset(t *testing.T) {
	gs, stub, b := newGuardedSender(3, time.Minute, nil)
	if err := gs.Send("job1", "miss", ""); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stub.calls != 1 {
		t.Fatalf("expected 1 inner call, got %d", stub.calls)
	}
	if b.State() != StateClosed {
		t.Fatalf("expected closed after success, got %s", b.State())
	}
}

func TestGuardedSender_FailureTripsBreaker(t *testing.T) {
	gs, stub, b := newGuardedSender(2, time.Minute, errors.New("timeout"))
	gs.Send("job1", "miss", "") //nolint:errcheck
	gs.Send("job1", "miss", "") //nolint:errcheck
	if b.State() != StateOpen {
		t.Fatalf("expected open after threshold failures, got %s", b.State())
	}
	if stub.calls != 2 {
		t.Fatalf("expected 2 inner calls, got %d", stub.calls)
	}
}

func TestGuardedSender_OpenCircuitSuppressesInnerCall(t *testing.T) {
	gs, stub, _ := newGuardedSender(1, time.Minute, errors.New("fail"))
	gs.Send("job1", "miss", "") //nolint:errcheck — opens the circuit
	prev := stub.calls
	err := gs.Send("job1", "miss", "")
	if err == nil {
		t.Fatal("expected error when circuit open")
	}
	if stub.calls != prev {
		t.Fatal("inner sender should not be called when circuit is open")
	}
}

func TestGuardedSender_RecoverAfterCooldown(t *testing.T) {
	gs, stub, b := newGuardedSender(1, 10*time.Millisecond, nil)
	// Force open by recording a failure directly.
	b.RecordFailure()
	time.Sleep(20 * time.Millisecond)
	if err := gs.Send("job1", "miss", ""); err != nil {
		t.Fatalf("unexpected error after cooldown: %v", err)
	}
	if stub.calls != 1 {
		t.Fatalf("expected inner call after recovery, got %d", stub.calls)
	}
	if b.State() != StateClosed {
		t.Fatalf("expected closed after successful probe, got %s", b.State())
	}
}
