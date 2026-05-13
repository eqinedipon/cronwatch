package circuitbreaker

import (
	"time"

	"github.com/cronwatch/internal/alerter"
	"github.com/cronwatch/internal/config"
)

const (
	defaultFailureThreshold = 3
	defaultCooldown         = 30 * time.Second
)

// FromConfig constructs a GuardedSender wrapping the provided Sender using
// circuit-breaker settings from cfg. If no circuit-breaker config is present
// the sender is returned unwrapped.
func FromConfig(cfg *config.Config, sender alerter.Sender) alerter.Sender {
	if cfg == nil || cfg.Alerting == nil || cfg.Alerting.CircuitBreaker == nil {
		return sender
	}

	cb := cfg.Alerting.CircuitBreaker

	threshold := cb.FailureThreshold
	if threshold <= 0 {
		threshold = defaultFailureThreshold
	}

	cooldown := time.Duration(cb.CooldownSeconds) * time.Second
	if cooldown <= 0 {
		cooldown = defaultCooldown
	}

	breaker := New(Options{
		FailureThreshold: threshold,
		Cooldown:         cooldown,
	})

	return NewGuardedSender(sender, breaker)
}
