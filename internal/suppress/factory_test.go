package suppress_test

import (
	"testing"
	"time"

	"github.com/cronwatch/cronwatch/internal/config"
	"github.com/cronwatch/cronwatch/internal/suppress"
)

func TestFromConfig_NilConfig(t *testing.T) {
	w := suppress.FromConfig(nil)
	if w.Duration() != 30*time.Minute {
		t.Fatalf("expected 30m default, got %v", w.Duration())
	}
}

func TestFromConfig_NilAlerting(t *testing.T) {
	w := suppress.FromConfig(&config.Config{})
	if w.Duration() != 30*time.Minute {
		t.Fatalf("expected 30m default, got %v", w.Duration())
	}
}

func TestFromConfig_ZeroMinutes(t *testing.T) {
	cfg := &config.Config{Alerting: &config.Alerting{SuppressWindowMinutes: 0}}
	w := suppress.FromConfig(cfg)
	if w.Duration() != 30*time.Minute {
		t.Fatalf("expected 30m default for zero, got %v", w.Duration())
	}
}

func TestFromConfig_CustomMinutes(t *testing.T) {
	cfg := &config.Config{Alerting: &config.Alerting{SuppressWindowMinutes: 60}}
	w := suppress.FromConfig(cfg)
	if w.Duration() != 60*time.Minute {
		t.Fatalf("expected 60m, got %v", w.Duration())
	}
}

func TestFromConfig_NegativeMinutes(t *testing.T) {
	cfg := &config.Config{Alerting: &config.Alerting{SuppressWindowMinutes: -10}}
	w := suppress.FromConfig(cfg)
	if w.Duration() != 30*time.Minute {
		t.Fatalf("expected 30m default for negative, got %v", w.Duration())
	}
}
