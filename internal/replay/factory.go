package replay

import (
	"log"

	"github.com/systemli/cronwatch/internal/config"
	"github.com/systemli/cronwatch/internal/runner"
	"github.com/systemli/cronwatch/internal/tracker"
)

// FromConfig constructs a Replayer from the application config, tracker, and dispatcher.
// Returns nil if replay is disabled or no dispatcher is provided.
func FromConfig(cfg *config.Config, tr *tracker.Tracker, d *runner.Dispatcher) *Replayer {
	if cfg == nil {
		log.Println("[replay] nil config, replay disabled")
		return nil
	}
	if d == nil {
		log.Println("[replay] nil dispatcher, replay disabled")
		return nil
	}
	if tr == nil {
		log.Println("[replay] nil tracker, replay disabled")
		return nil
	}
	if cfg.Replay == nil || !cfg.Replay.Enabled {
		log.Println("[replay] replay not enabled in config")
		return nil
	}

	log.Println("[replay] replay enabled")
	return New(cfg.Jobs, tr, d)
}
