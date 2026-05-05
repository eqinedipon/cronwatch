package config_test

import (
	"os"
	"testing"

	"github.com/yourorg/cronwatch/internal/config"
)

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "cronwatch-*.yaml")
	if err != nil {
		t.Fatalf("creating temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("writing temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestLoad_ValidConfig(t *testing.T) {
	yaml := `
log_level: info
notifier:
  email: ops@example.com
jobs:
  - name: backup
    schedule: "0 2 * * *"
    timeout: 30m
    alert:
      on_miss: true
      on_fail: true
`
	path := writeTempConfig(t, yaml)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(cfg.Jobs) != 1 {
		t.Errorf("expected 1 job, got %d", len(cfg.Jobs))
	}
	if cfg.Jobs[0].Name != "backup" {
		t.Errorf("expected job name 'backup', got %q", cfg.Jobs[0].Name)
	}
	if cfg.Notifier.Email != "ops@example.com" {
		t.Errorf("expected email 'ops@example.com', got %q", cfg.Notifier.Email)
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := config.Load("/nonexistent/path/config.yaml")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoad_NoJobs(t *testing.T) {
	yaml := `log_level: debug\njobs: []\n`
	path := writeTempConfig(t, yaml)
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected validation error for empty jobs, got nil")
	}
}

func TestLoad_JobMissingSchedule(t *testing.T) {
	yaml := `
jobs:
  - name: orphan
`
	path := writeTempConfig(t, yaml)
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected validation error for missing schedule, got nil")
	}
}
