package alerter

import (
	"testing"

	"cronwatch/internal/config"
)

func TestFromConfig_NilAlerting(t *testing.T) {
	cfg := &config.Config{}
	s, err := FromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := s.(*LogSender); !ok {
		t.Errorf("expected *LogSender, got %T", s)
	}
}

func TestFromConfig_LogOnly(t *testing.T) {
	cfg := &config.Config{Alerting: &config.AlertingConfig{Log: true}}
	s, err := FromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := s.(*LogSender); !ok {
		t.Errorf("expected *LogSender, got %T", s)
	}
}

func TestFromConfig_WebhookOnly(t *testing.T) {
	cfg := &config.Config{Alerting: &config.AlertingConfig{WebhookURL: "http://example.com/hook"}}
	s, err := FromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := s.(*WebhookSender); !ok {
		t.Errorf("expected *WebhookSender, got %T", s)
	}
}

func TestFromConfig_MultiSender(t *testing.T) {
	cfg := &config.Config{Alerting: &config.AlertingConfig{
		Log:        true,
		WebhookURL: "http://example.com/hook",
	}}
	s, err := FromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := s.(*MultiSender); !ok {
		t.Errorf("expected *MultiSender, got %T", s)
	}
}

func TestFromConfig_NoSendersConfigured(t *testing.T) {
	cfg := &config.Config{Alerting: &config.AlertingConfig{}}
	_, err := FromConfig(cfg)
	if err == nil {
		t.Fatal("expected error when alerting block has no senders")
	}
}
