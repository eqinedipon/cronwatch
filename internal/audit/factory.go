package audit

import (
	"fmt"
	"io"
	"os"
)

// Config holds audit logging configuration sourced from the main config file.
type Config struct {
	Enabled bool   `yaml:"enabled"`
	Path    string `yaml:"path"`
}

// FromConfig constructs a Logger from the given Config.
// If audit logging is disabled or no path is set, a no-op logger is returned.
func FromConfig(cfg Config) (*Logger, error) {
	if !cfg.Enabled || cfg.Path == "" {
		return newNopLogger(), nil
	}
	l, err := New(cfg.Path)
	if err != nil {
		return nil, fmt.Errorf("audit: open log file: %w", err)
	}
	return l, nil
}

// newNopLogger returns a Logger that discards all output.
func newNopLogger() *Logger {
	return &Logger{
		enc: newJSONEncoder(io.Discard),
		f:   os.NewFile(0, "/dev/null"),
	}
}
