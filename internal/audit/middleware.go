package audit

import "time"

// Recorder wraps a Logger and exposes helper methods used by other packages
// to record job lifecycle events without constructing Entry structs manually.
type Recorder struct {
	logger *Logger
}

// NewRecorder wraps an existing Logger.
func NewRecorder(l *Logger) *Recorder {
	return &Recorder{logger: l}
}

// RecordSuccess logs a successful job execution.
func (r *Recorder) RecordSuccess(job string, durationMs float64) error {
	return r.logger.Log(Entry{
		Job:      job,
		Event:    EventSuccess,
		Duration: durationMs,
	})
}

// RecordFailure logs a failed job execution with an optional message.
func (r *Recorder) RecordFailure(job, message string, durationMs float64) error {
	return r.logger.Log(Entry{
		Job:      job,
		Event:    EventFailure,
		Message:  message,
		Duration: durationMs,
	})
}

// RecordMiss logs a missed job execution.
func (r *Recorder) RecordMiss(job string, at time.Time) error {
	return r.logger.Log(Entry{
		Timestamp: at,
		Job:       job,
		Event:     EventMiss,
	})
}
