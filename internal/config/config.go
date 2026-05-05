package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Job defines a single cron job to monitor.
type Job struct {
	Name     string        `yaml:"name"`
	Schedule string        `yaml:"schedule"`
	Timeout  time.Duration `yaml:"timeout"`
	Alert    AlertConfig   `yaml:"alert"`
}

// AlertConfig holds alerting settings for a job.
type AlertConfig struct {
	OnMiss   bool `yaml:"on_miss"`
	OnFail   bool `yaml:"on_fail"`
}

// NotifierConfig holds global notification settings.
type NotifierConfig struct {
	Email   string `yaml:"email"`
	Webhook string `yaml:"webhook"`
}

// Config is the top-level cronwatch configuration.
type Config struct {
	LogLevel string         `yaml:"log_level"`
	Notifier NotifierConfig `yaml:"notifier"`
	Jobs     []Job          `yaml:"jobs"`
}

// Load reads and parses a YAML config file from the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file %q: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config file %q: %w", path, err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &cfg, nil
}

// validate checks that required fields are present and valid.
func (c *Config) validate() error {
	if len(c.Jobs) == 0 {
		return fmt.Errorf("at least one job must be defined")
	}
	for i, job := range c.Jobs {
		if job.Name == "" {
			return fmt.Errorf("job[%d]: name is required", i)
		}
		if job.Schedule == "" {
			return fmt.Errorf("job %q: schedule is required", job.Name)
		}
	}
	return nil
}
