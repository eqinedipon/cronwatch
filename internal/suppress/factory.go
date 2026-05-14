package suppress

import (
	"time"

	"github.com/cronwatch/cronwatch/internal/config"
)

// FromConfig constructs a Window from the application configuration.
// If no suppression window is configured a sensible default is applied.
func FromConfig(cfg *config.Config) *Window {
	if cfg == nil || cfg.Alerting == nil {
		return New(0)
	}

	d := time.Duration(cfg.Alerting.SuppressWindowMinutes) * time.Minute
	return New(d)
}
