package alerter

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

// Alert represents a notification to be sent.
type Alert struct {
	JobName   string
	Kind      string // "missed" or "failed"
	Message   string
	Timestamp time.Time
}

// Sender is the interface for sending alerts.
type Sender interface {
	Send(a Alert) error
}

// LogSender writes alerts to the standard logger.
type LogSender struct{}

func (l *LogSender) Send(a Alert) error {
	log.Printf("[ALERT] job=%s kind=%s msg=%s ts=%s",
		a.JobName, a.Kind, a.Message, a.Timestamp.Format(time.RFC3339))
	return nil
}

// WebhookSender posts alerts to an HTTP endpoint.
type WebhookSender struct {
	URL    string
	Client *http.Client
}

func NewWebhookSender(url string) *WebhookSender {
	return &WebhookSender{
		URL:    url,
		Client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (w *WebhookSender) Send(a Alert) error {
	body := fmt.Sprintf(
		`{"job":%q,"kind":%q,"message":%q,"timestamp":%q}`,
		a.JobName, a.Kind, a.Message, a.Timestamp.Format(time.RFC3339),
	)
	resp, err := w.Client.Post(w.URL, "application/json", strings.NewReader(body))
	if err != nil {
		return fmt.Errorf("webhook send: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("webhook returned status %d", resp.StatusCode)
	}
	return nil
}

// MultiSender fans out an alert to multiple senders.
type MultiSender struct {
	Senders []Sender
}

func (m *MultiSender) Send(a Alert) error {
	var errs []string
	for _, s := range m.Senders {
		if err := s.Send(a); err != nil {
			errs = append(errs, err.Error())
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("alerter errors: %s", strings.Join(errs, "; "))
	}
	return nil
}
