package circuitbreaker_test

import (
	"errors"
	"testing"
	"time"

	"github.com/cronwatch/internal/alerter"
	"github.com/cronwatch/internal/circuitbreaker"
	"github.com/cronwatch/internal/config"
)

type stubSender struct{ called int }

func (s *stubSender) Send(_ alerter.Alert) error {
	s.called++
	return nil
}

func TestFromConfig_NilConfig(t *testing.T) {
	sender := &stubSender{}
	result := circuitbreaker.FromConfig(nil, sender)
	if result != sender {
		t.Fatal("expected original sender when config is nil")
	}
}

func TestFromConfig_NilAlerting(t *testing.T) {
	cfg := &config.Config{}
	sender := &stubSender{}
	result := circuitbreaker.FromConfig(cfg, sender)
	if result != sender {
		t.Fatal("expected original sender when alerting config is nil")
	}
}

func TestFromConfig_NilCircuitBreaker(t *testing.T) {
	cfg := &config.Config{Alerting: &config.AlertingConfig{}}
	sender := &stubSender{}
	result := circuitbreaker.FromConfig(cfg, sender)
	if result != sender {
		t.Fatal("expected original sender when circuit breaker config is nil")
	}
}

func TestFromConfig_DefaultThresholdAndCooldown(t *testing.T) {
	cfg := &config.Config{
		Alerting: &config.AlertingConfig{
			CircuitBreaker: &config.CircuitBreakerConfig{},
		},
	}
	sender := &stubSender{}
	guarded := circuitbreaker.FromConfig(cfg, sender)
	if guarded == sender {
		t.Fatal("expected a wrapped GuardedSender")
	}
	// Verify defaults: 3 failures should open the circuit.
	for i := 0; i < 3; i++ {
		_ = guarded.Send(alerter.Alert{JobName: "job", Kind: "failure", Err: errors.New("boom")})
	}
	// After threshold failures the circuit should be open and block sends.
	err := guarded.Send(alerter.Alert{JobName: "job", Kind: "failure"})
	if err == nil {
		t.Fatal("expected circuit-open error after default threshold")
	}
}

func TestFromConfig_CustomSettings(t *testing.T) {
	cfg := &config.Config{
		Alerting: &config.AlertingConfig{
			CircuitBreaker: &config.CircuitBreakerConfig{
				FailureThreshold: 1,
				CooldownSeconds:  int(time.Hour.Seconds()),
			},
		},
	}
	sender := &stubSender{}
	guarded := circuitbreaker.FromConfig(cfg, sender)
	// One failure should open the circuit immediately.
	_ = guarded.Send(alerter.Alert{JobName: "j", Kind: "failure", Err: errors.New("e")})
	err := guarded.Send(alerter.Alert{JobName: "j", Kind: "failure"})
	if err == nil {
		t.Fatal("expected circuit-open error after threshold of 1")
	}
}
