package audit

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

// EventType classifies audit log entries.
type EventType string

const (
	EventSuccess EventType = "success"
	EventFailure EventType = "failure"
	EventMiss    EventType = "miss"
)

// Entry represents a single audit log record.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Job       string    `json:"job"`
	Event     EventType `json:"event"`
	Message   string    `json:"message,omitempty"`
	Duration  float64   `json:"duration_ms,omitempty"`
}

// Logger writes structured audit entries to a file.
type Logger struct {
	mu  sync.Mutex
	enc *json.Encoder
	f   *os.File
}

// New opens (or creates) the audit log file at path and returns a Logger.
func New(path string) (*Logger, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, err
	}
	return &Logger{enc: json.NewEncoder(f), f: f}, nil
}

// Log writes an entry to the audit log.
func (l *Logger) Log(e Entry) error {
	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now().UTC()
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.enc.Encode(e)
}

// Close flushes and closes the underlying file.
func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.f.Close()
}
