package alerter

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var testAlert = Alert{
	JobName:   "backup",
	Kind:      "missed",
	Message:   "job did not run within grace period",
	Timestamp: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
}

func TestLogSender_Send(t *testing.T) {
	s := &LogSender{}
	if err := s.Send(testAlert); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestWebhookSender_Send_Success(t *testing.T) {
	var received string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		received = r.Header.Get("Content-Type")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	s := NewWebhookSender(server.URL)
	if err := s.Send(testAlert); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received != "application/json" {
		t.Errorf("expected application/json, got %s", received)
	}
}

func TestWebhookSender_Send_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	s := NewWebhookSender(server.URL)
	if err := s.Send(testAlert); err == nil {
		t.Fatal("expected error for non-2xx status")
	}
}

func TestMultiSender_AllSucceed(t *testing.T) {
	m := &MultiSender{Senders: []Sender{&LogSender{}, &LogSender{}}}
	if err := m.Send(testAlert); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMultiSender_PartialFailure(t *testing.T) {
	badSender := NewWebhookSender("http://127.0.0.1:1") // unreachable
	m := &MultiSender{Senders: []Sender{&LogSender{}, badSender}}
	if err := m.Send(testAlert); err == nil {
		t.Fatal("expected error when one sender fails")
	}
}
