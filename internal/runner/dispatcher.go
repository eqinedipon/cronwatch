package runner

import (
	"context"
	"log"
	"time"

	"cronwatch/internal/alerter"
	"cronwatch/internal/config"
	"cronwatch/internal/tracker"
)

// Dispatcher runs jobs and records results via the tracker and alerter.
type Dispatcher struct {
	runner  *Runner
	tracker *tracker.Tracker
	alerter alerter.Sender
}

// NewDispatcher creates a Dispatcher with the provided dependencies.
func NewDispatcher(r *Runner, tr *tracker.Tracker, al alerter.Sender) *Dispatcher {
	return &Dispatcher{runner: r, tracker: tr, alerter: al}
}

// Dispatch executes a job and records success or failure.
func (d *Dispatcher) Dispatch(ctx context.Context, job config.Job) {
	if job.Command == "" {
		log.Printf("[runner] job %q has no command, skipping", job.Name)
		return
	}

	log.Printf("[runner] starting job %q: %s", job.Name, job.Command)
	res := d.runner.Run(ctx, job.Name, job.Command)

	if res.Err != nil || res.ExitCode != 0 {
		log.Printf("[runner] job %q failed (exit %d): %v", job.Name, res.ExitCode, res.Err)
		d.tracker.RecordFailure(job.Name)

		msg := alerter.Message{
			JobName: job.Name,
			Kind:    alerter.KindFailure,
			Detail:  res.Output,
			At:      res.Started,
		}
		if err := d.alerter.Send(ctx, msg); err != nil {
			log.Printf("[runner] alert send error for job %q: %v", job.Name, err)
		}
		return
	}

	log.Printf("[runner] job %q succeeded in %s", job.Name, res.Duration.Round(time.Millisecond))
	d.tracker.RecordSuccess(job.Name)
}
