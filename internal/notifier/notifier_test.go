package notifier_test

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/cronwatch/internal/notifier"
	"github.com/cronwatch/internal/tracker"
)

// fakeSender captures messages sent to it.
type fakeSender struct {
	Messages []string
	Err      error
}

func (f *fakeSender) Send(msg string) error {
	if f.Err != nil {
		return f.Err
	}
	f.Messages = append(f.Messages, msg)
	return nil
}

func newNotifier(t *testing.T) (*notifier.Notifier, *fakeSender, *tracker.Tracker) {
	t.Helper()
	tr := tracker.New()
	tr.RecordSuccess("myjob", time.Now())
	sender := &fakeSender{}
	n := notifier.New(sender, tr)
	return n, sender, tr
}

func TestNotifyMiss_SendsAlert(t *testing.T) {
	n, sender, tr := newNotifier(t)
	tr.RecordMiss("myjob")

	if err := n.NotifyMiss("myjob"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(sender.Messages) != 1 {
		t.Fatalf("expected 1 message, got %d", len(sender.Messages))
	}
	if !strings.Contains(sender.Messages[0], "miss") {
		t.Errorf("expected message to contain 'miss', got: %s", sender.Messages[0])
	}
}

func TestNotifyMiss_UnknownJob(t *testing.T) {
	tr := tracker.New()
	n := notifier.New(&fakeSender{}, tr)

	err := n.NotifyMiss("ghost")
	if err == nil {
		t.Fatal("expected error for unknown job, got nil")
	}
}

func TestNotifyFailure_SendsAlert(t *testing.T) {
	n, sender, _ := newNotifier(t)

	if err := n.NotifyFailure("myjob", "exit status 1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(sender.Messages[0], "failure") {
		t.Errorf("expected 'failure' in message, got: %s", sender.Messages[0])
	}
	if !strings.Contains(sender.Messages[0], "exit status 1") {
		t.Errorf("expected output in message, got: %s", sender.Messages[0])
	}
}

func TestNotifyRecovery_SendsAlert(t *testing.T) {
	n, sender, _ := newNotifier(t)

	if err := n.NotifyRecovery("myjob"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(sender.Messages[0], "recovery") {
		t.Errorf("expected 'recovery' in message, got: %s", sender.Messages[0])
	}
}

func TestNotify_SenderError_Propagated(t *testing.T) {
	tr := tracker.New()
	tr.RecordSuccess("myjob", time.Now())
	sender := &fakeSender{Err: errors.New("webhook unreachable")}
	n := notifier.New(sender, tr)

	err := n.NotifyRecovery("myjob")
	if err == nil {
		t.Fatal("expected error from sender, got nil")
	}
}
