// Package replay provides the ability to re-trigger missed cron jobs
// by scanning the tracker state and dispatching overdue jobs on demand.
package replay

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/your-org/cronwatch/internal/config"
	"github.com/your-org/cronwatch/internal/schedule"
	"github.com/your-org/cronwatch/internal/tracker"
)

// Dispatcher is the interface used to run a job by name.
type Dispatcher interface {
	Dispatch(ctx context.Context, name, command string) error
}

// Replayer scans for missed jobs and re-dispatches them.
type Replayer struct {
	cfg     *config.Config
	tracker *tracker.Tracker
	disp    Dispatcher
	now     func() time.Time
}

// New creates a new Replayer.
func New(cfg *config.Config, t *tracker.Tracker, d Dispatcher) *Replayer {
	return &Replayer{
		cfg:     cfg,
		tracker: t,
		disp:    d,
		now:     time.Now,
	}
}

// ReplayResult holds the outcome for a single replayed job.
type ReplayResult struct {
	Job   string
	Error error
}

// Run iterates all configured jobs, identifies those that are overdue
// (last success is before the previous scheduled run), and dispatches them.
func (r *Replayer) Run(ctx context.Context) []ReplayResult {
	now := r.now()
	var results []ReplayResult

	for _, job := range r.cfg.Jobs {
		prev, err := schedule.PrevRun(job.Schedule, now)
		if err != nil {
			log.Printf("replay: invalid schedule for job %q: %v", job.Name, err)
			continue
		}

		state, err := r.tracker.Get(job.Name)
		if err != nil {
			// job has never run — eligible for replay
			state = nil
		}

		if state != nil && !state.LastSuccess.IsZero() && state.LastSuccess.After(prev) {
			// job ran successfully after the last scheduled time; skip
			continue
		}

		log.Printf("replay: dispatching missed job %q (prev=%s)", job.Name, prev.Format(time.RFC3339))
		dispErr := r.disp.Dispatch(ctx, job.Name, job.Command)
		if dispErr != nil {
			dispErr = fmt.Errorf("replay dispatch %q: %w", job.Name, dispErr)
		}
		results = append(results, ReplayResult{Job: job.Name, Error: dispErr})
	}

	return results
}
