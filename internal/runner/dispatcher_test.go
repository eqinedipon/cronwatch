package runner

import (
	"context"
	"testing"
	"time"

	"cronwatch/internal/alerter"
	"cronwatch/internal/config"
	"cronwatch/internal/tracker"
)

type mockSender struct {
	called  int
	lastMsg alerter.Message
}

func (m *mockSender) Send(_ context.Context, msg alerter.Message) error {
	m.called++
	m.lastMsg = msg
	return nil
}

func newTestDispatcher(t *testing.T, sender alerter.Sender) (*Dispatcher, *tracker.Tracker) {
	t.Helper()
	tr := tracker.New()
	r := New(5 * time.Second)
	return NewDispatcher(r, tr, sender), tr
}

func TestDispatch_Success(t *testing.T) {
	sender := &mockSender{}
	d, tr := newTestDispatcher(t, sender)

	job := config.Job{Name: "ok-job", Command: "exit 0", Schedule: "* * * * *"}
	d.Dispatch(context.Background(), job)

	st, _ := tr.Get("ok-job")
	if st.Successes != 1 {
		t.Errorf("expected 1 success, got %d", st.Successes)
	}
	if sender.called != 0 {
		t.Errorf("expected no alert on success, got %d", sender.called)
	}
}

func TestDispatch_Failure(t *testing.T) {
	sender := &mockSender{}
	d, tr := newTestDispatcher(t, sender)

	job := config.Job{Name: "bad-job", Command: "exit 1", Schedule: "* * * * *"}
	d.Dispatch(context.Background(), job)

	st, _ := tr.Get("bad-job")
	if st.Failures != 1 {
		t.Errorf("expected 1 failure, got %d", st.Failures)
	}
	if sender.called != 1 {
		t.Errorf("expected 1 alert, got %d", sender.called)
	}
	if sender.lastMsg.Kind != alerter.KindFailure {
		t.Errorf("expected KindFailure, got %v", sender.lastMsg.Kind)
	}
}

func TestDispatch_NoCommand(t *testing.T) {
	sender := &mockSender{}
	d, _ := newTestDispatcher(t, sender)

	job := config.Job{Name: "empty-job", Command: "", Schedule: "* * * * *"}
	d.Dispatch(context.Background(), job) // should not panic

	if sender.called != 0 {
		t.Errorf("expected no alert for empty command, got %d", sender.called)
	}
}
