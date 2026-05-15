package suppress_test

import (
	"testing"

	"github.com/cronwatch/internal/config"
	"github.com/cronwatch/internal/suppress"
)

func TestFromConfig_NilConfig(t *testing.T) {
	gs := suppress.FromConfig(nil, nil)
	if gs != nil {
		t.Fatal("expected nil GuardedSender for nil config")
	}
}

func TestFromConfig_NilAlerting(t *testing.T) {
	cfg := &config.Config{}
	gs := suppress.FromConfig(cfg, nil)
	if gs != nil {
		t.Fatal("expected nil GuardedSender when alerting config is nil")
	}
}

func TestFromConfig_ZeroMinutes(t *testing.T) {
	cfg := &config.Config{
		Alerting: &config.AlertingConfig{
			SuppressWindowMinutes: 0,
		},
	}
	gs := suppress.FromConfig(cfg, &mockSender{})
	if gs != nil {
		t.Fatal("expected nil GuardedSender for zero suppress window")
	}
}

func TestFromConfig_CustomMinutes(t *testing.T) {
	cfg := &config.Config{
		Alerting: &config.AlertingConfig{
			SuppressWindowMinutes: 10,
		},
	}
	gs := suppress.FromConfig(cfg, &mockSender{})
	if gs == nil {
		t.Fatal("expected non-nil GuardedSender")
	}
}

func TestFromConfig_NegativeMinutes(t *testing.T) {
	cfg := &config.Config{
		Alerting: &config.AlertingConfig{
			SuppressWindowMinutes: -5,
		},
	}
	gs := suppress.FromConfig(cfg, &mockSender{})
	if gs != nil {
		t.Fatal("expected nil GuardedSender for negative suppress window")
	}
}
