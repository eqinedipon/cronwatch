package audit_test

import (
	"bufio"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/cronwatch/cronwatch/internal/audit"
)

func newRecorder(t *testing.T) (*audit.Recorder, string) {
	t.Helper()
	path := t.TempDir() + "/audit.log"
	l, err := audit.New(path)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { l.Close() })
	return audit.NewRecorder(l), path
}

func readEntry(t *testing.T, path string) audit.Entry {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	var e audit.Entry
	if err := json.NewDecoder(bufio.NewReader(f)).Decode(&e); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return e
}

func TestRecordSuccess_WritesEntry(t *testing.T) {
	rec, path := newRecorder(t)
	if err := rec.RecordSuccess("nightly", 250.0); err != nil {
		t.Fatal(err)
	}
	e := readEntry(t, path)
	if e.Job != "nightly" || e.Event != audit.EventSuccess || e.Duration != 250.0 {
		t.Errorf("unexpected entry: %+v", e)
	}
}

func TestRecordFailure_WritesEntry(t *testing.T) {
	rec, path := newRecorder(t)
	_ = rec.RecordFailure("sync", "exit status 1", 80.5)
	e := readEntry(t, path)
	if e.Event != audit.EventFailure {
		t.Errorf("event = %q, want failure", e.Event)
	}
	if e.Message != "exit status 1" {
		t.Errorf("message = %q", e.Message)
	}
}

func TestRecordMiss_WritesEntry(t *testing.T) {
	rec, path := newRecorder(t)
	at := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	_ = rec.RecordMiss("report", at)
	e := readEntry(t, path)
	if e.Event != audit.EventMiss {
		t.Errorf("event = %q, want miss", e.Event)
	}
	if !e.Timestamp.Equal(at) {
		t.Errorf("timestamp = %v, want %v", e.Timestamp, at)
	}
}
