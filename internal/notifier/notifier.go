package notifier

import (
	"fmt"
	"time"

	"github.com/cronwatch/internal/alerter"
	"github.com/cronwatch/internal/tracker"
)

// EventType represents the kind of notification event.
type EventType string

const (
	EventMiss    EventType = "miss"
	EventFailure EventType = "failure"
	EventRecovery EventType = "recovery"
)

// Event holds the data for a single notification.
type Event struct {
	JobName   string
	Type      EventType
	Message   string
	OccurredAt time.Time
}

// Notifier dispatches alerts based on job state changes.
type Notifier struct {
	sender  alerter.Sender
	tracker *tracker.Tracker
}

// New creates a Notifier wired to the given sender and tracker.
func New(sender alerter.Sender, t *tracker.Tracker) *Notifier {
	return &Notifier{sender: sender, tracker: t}
}

// Notify sends an alert for the given event if conditions are met.
func (n *Notifier) Notify(ev Event) error {
	msg := n.format(ev)
	return n.sender.Send(msg)
}

// NotifyMiss sends a miss alert for the named job.
func (n *Notifier) NotifyMiss(jobName string) error {
	state, ok := n.tracker.Get(jobName)
	if !ok {
		return fmt.Errorf("notifier: unknown job %q", jobName)
	}
	ev := Event{
		JobName:    jobName,
		Type:       EventMiss,
		Message:    fmt.Sprintf("job %q missed its schedule (consecutive misses: %d)", jobName, state.ConsecutiveMisses),
		OccurredAt: time.Now(),
	}
	return n.Notify(ev)
}

// NotifyFailure sends a failure alert for the named job.
func (n *Notifier) NotifyFailure(jobName, output string) error {
	ev := Event{
		JobName:    jobName,
		Type:       EventFailure,
		Message:    fmt.Sprintf("job %q failed: %s", jobName, output),
		OccurredAt: time.Now(),
	}
	return n.Notify(ev)
}

// NotifyRecovery sends a recovery alert when a previously failing job succeeds.
func (n *Notifier) NotifyRecovery(jobName string) error {
	ev := Event{
		JobName:    jobName,
		Type:       EventRecovery,
		Message:    fmt.Sprintf("job %q has recovered and completed successfully", jobName),
		OccurredAt: time.Now(),
	}
	return n.Notify(ev)
}

func (n *Notifier) format(ev Event) string {
	return fmt.Sprintf("[cronwatch][%s] %s at %s", ev.Type, ev.Message, ev.OccurredAt.Format(time.RFC3339))
}
