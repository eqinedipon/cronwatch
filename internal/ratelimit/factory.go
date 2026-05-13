package ratelimit

import (
	"time"

	"github.com/cronwatch/cronwatch/internal/config"
)

const defaultCooldown = 15 * time.Minute

// FromConfig constructs a Limiter from the alerting section of the config.
// If no cooldown is specified, defaultCooldown is used.
func FromConfig(cfg *config.Config) *Limiter {
	if cfg == nil || cfg.Alerting == nil {
		return New(defaultCooldown)
	}
	if cfg.Alerting.CooldownMinutes <= 0 {
		return New(defaultCooldown)
	}
	cooldown := time.Duration(cfg.Alerting.CooldownMinutes) * time.Minute
	return New(cooldown)
}
