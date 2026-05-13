package retention

import (
	"time"

	"github.com/cronwatch/internal/config"
)

const (
	defaultMaxAge   = 7 * 24 * time.Hour
	defaultMaxFiles = 30
	defaultPattern  = "audit-*.log"
)

// FromConfig constructs a Cleaner from the application config.
// If auditing is not configured, it returns nil.
func FromConfig(cfg *config.Config) *Cleaner {
	if cfg == nil || cfg.Audit == nil {
		return nil
	}

	dir := cfg.Audit.Dir
	if dir == "" {
		return nil
	}

	maxAge := defaultMaxAge
	if cfg.Audit.RetentionDays > 0 {
		maxAge = time.Duration(cfg.Audit.RetentionDays) * 24 * time.Hour
	}

	maxFiles := defaultMaxFiles
	if cfg.Audit.MaxFiles > 0 {
		maxFiles = cfg.Audit.MaxFiles
	}

	return New(Policy{
		Dir:      dir,
		MaxAge:   maxAge,
		MaxFiles: maxFiles,
		Pattern:  defaultPattern,
	})
}
