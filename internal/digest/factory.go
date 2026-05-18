package digest

import (
	"time"

	"github.com/user/cronwatch/internal/alerter"
	"github.com/user/cronwatch/internal/config"
	"github.com/user/cronwatch/internal/metrics"
	"github.com/user/cronwatch/internal/tracker"
)

// Sender is the subset of alerter.Sender used by the digest scheduler.
type Sender interface {
	Send(subject, body string) error
}

// Scheduler drives periodic digest delivery.
type Scheduler struct {
	reporter *Reporter
	sender   Sender
	period   time.Duration
	stop     chan struct{}
}

// FromConfig constructs a Scheduler from application config.
// Returns nil if digest reporting is not configured.
func FromConfig(cfg *config.Config, m *metrics.Metrics, t *tracker.Tracker, s alerter.Sender) *Scheduler {
	if cfg == nil || cfg.Digest == nil || cfg.Digest.IntervalMinutes <= 0 {
		return nil
	}
	names := make([]string, 0, len(cfg.Jobs))
	for _, j := range cfg.Jobs {
		names = append(names, j.Name)
	}
	return &Scheduler{
		reporter: New(m, t, names),
		sender:   s,
		period:   time.Duration(cfg.Digest.IntervalMinutes) * time.Minute,
		stop:     make(chan struct{}),
	}
}

// Start begins the periodic digest loop in a goroutine.
func (sc *Scheduler) Start() {
	go func() {
		ticker := time.NewTicker(sc.period)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				sc.send()
			case <-sc.stop:
				return
			}
		}
	}()
}

// Stop halts the digest scheduler.
func (sc *Scheduler) Stop() { close(sc.stop) }

func (sc *Scheduler) send() {
	s := sc.reporter.Build(sc.period)
	_ = sc.sender.Send("cronwatch digest", Format(s))
}
