package replay

import (
	"testing"

	"github.com/systemli/cronwatch/internal/config"
	"github.com/systemli/cronwatch/internal/runner"
	"github.com/systemli/cronwatch/internal/tracker"
)

func makeDispatcher() *runner.Dispatcher {
	return runner.NewDispatcher(nil)
}

func makeTracker() *tracker.Tracker {
	return tracker.New()
}

func TestFromConfig_NilConfig(t *testing.T) {
	r := FromConfig(nil, makeTracker(), makeDispatcher())
	if r != nil {
		t.Fatal("expected nil replayer for nil config")
	}
}

func TestFromConfig_NilDispatcher(t *testing.T) {
	cfg := &config.Config{Replay: &config.ReplayConfig{Enabled: true}}
	r := FromConfig(cfg, makeTracker(), nil)
	if r != nil {
		t.Fatal("expected nil replayer for nil dispatcher")
	}
}

func TestFromConfig_NilTracker(t *testing.T) {
	cfg := &config.Config{Replay: &config.ReplayConfig{Enabled: true}}
	r := FromConfig(cfg, nil, makeDispatcher())
	if r != nil {
		t.Fatal("expected nil replayer for nil tracker")
	}
}

func TestFromConfig_ReplayDisabled(t *testing.T) {
	cfg := &config.Config{Replay: &config.ReplayConfig{Enabled: false}}
	r := FromConfig(cfg, makeTracker(), makeDispatcher())
	if r != nil {
		t.Fatal("expected nil replayer when replay disabled")
	}
}

func TestFromConfig_NilReplaySection(t *testing.T) {
	cfg := &config.Config{}
	r := FromConfig(cfg, makeTracker(), makeDispatcher())
	if r != nil {
		t.Fatal("expected nil replayer when replay config section is nil")
	}
}

func TestFromConfig_ReplayEnabled(t *testing.T) {
	cfg := &config.Config{
		Jobs:   []config.Job{{Name: "job1", Schedule: "@hourly"}},
		Replay: &config.ReplayConfig{Enabled: true},
	}
	r := FromConfig(cfg, makeTracker(), makeDispatcher())
	if r == nil {
		t.Fatal("expected non-nil replayer when replay is enabled")
	}
}
