package watcher

import (
	"log"
	"time"

	"github.com/cronwatch/internal/alerter"
	"github.com/cronwatch/internal/config"
	"github.com/cronwatch/internal/schedule"
	"github.com/cronwatch/internal/tracker"
)

// Watcher periodically checks all configured cron jobs and alerts on misses.
type Watcher struct {
	cfg     *config.Config
	tracker *tracker.Tracker
	alerter alerter.Sender
	ticker  *time.Ticker
	stop    chan struct{}
}

// New creates a new Watcher.
func New(cfg *config.Config, tr *tracker.Tracker, al alerter.Sender) *Watcher {
	return &Watcher{
		cfg:     cfg,
		tracker: tr,
		alerter: al,
		stop:    make(chan struct{}),
	}
}

// Start begins the watch loop, checking jobs every interval.
func (w *Watcher) Start(interval time.Duration) {
	w.ticker = time.NewTicker(interval)
	go func() {
		for {
			select {
			case t := <-w.ticker.C:
				w.check(t)
			case <-w.stop:
				return
			}
		}
	}()
}

// Stop halts the watch loop.
func (w *Watcher) Stop() {
	w.ticker.Stop()
	close(w.stop)
}

// check inspects each job at time now and records a miss if overdue.
func (w *Watcher) check(now time.Time) {
	for _, job := range w.cfg.Jobs {
		state, err := w.tracker.Get(job.Name)
		if err != nil {
			// Never seen before — skip until first heartbeat.
			continue
		}

		prev, err := schedule.PrevRun(job.Schedule, now)
		if err != nil {
			log.Printf("watcher: invalid schedule for job %q: %v", job.Name, err)
			continue
		}

		grace := time.Duration(job.GraceSeconds) * time.Second
		deadline := prev.Add(grace)

		if now.After(deadline) && state.LastSuccess.Before(prev) {
			w.tracker.RecordMiss(job.Name)
			msg := "job " + job.Name + " missed its scheduled run at " + prev.Format(time.RFC3339)
			if alertErr := w.alerter.Send(job.Name, msg); alertErr != nil {
				log.Printf("watcher: alert failed for job %q: %v", job.Name, alertErr)
			}
		}
	}
}
