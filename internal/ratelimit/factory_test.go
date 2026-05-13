package ratelimit

import (
	"testing"
	"time"

	"github.com/cronwatch/cronwatch/internal/config"
)

func TestFromConfig_NilConfig(t *testing.T) {
	l := FromConfig(nil)
	if l == nil {
		t.Fatal("expected non-nil limiter")
	}
	if l.cooldown != defaultCooldown {
		t.Fatalf("expected default cooldown %v, got %v", defaultCooldown, l.cooldown)
	}
}

func TestFromConfig_NilAlerting(t *testing.T) {
	l := FromConfig(&config.Config{})
	if l.cooldown != defaultCooldown {
		t.Fatalf("expected default cooldown, got %v", l.cooldown)
	}
}

func TestFromConfig_ZeroCooldown(t *testing.T) {
	cfg := &config.Config{
		Alerting: &config.Alerting{CooldownMinutes: 0},
	}
	l := FromConfig(cfg)
	if l.cooldown != defaultCooldown {
		t.Fatalf("expected default cooldown for zero value, got %v", l.cooldown)
	}
}

func TestFromConfig_CustomCooldown(t *testing.T) {
	cfg := &config.Config{
		Alerting: &config.Alerting{CooldownMinutes: 30},
	}
	l := FromConfig(cfg)
	want := 30 * time.Minute
	if l.cooldown != want {
		t.Fatalf("expected %v, got %v", want, l.cooldown)
	}
}

func TestFromConfig_NegativeCooldown(t *testing.T) {
	cfg := &config.Config{
		Alerting: &config.Alerting{CooldownMinutes: -5},
	}
	l := FromConfig(cfg)
	if l.cooldown != defaultCooldown {
		t.Fatalf("expected default cooldown for negative value, got %v", l.cooldown)
	}
}
