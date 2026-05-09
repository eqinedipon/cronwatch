package audit_test

import (
	"bufio"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/cronwatch/cronwatch/internal/audit"
)

func TestNew_CreatesFile(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "audit-*.log")
	f.Close()
	if err != nil {
		t.Fatal(err)
	}
	l, err := audit.New(f.Name())
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer l.Close()
}

func TestLog_WritesEntry(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/audit.log"

	l, err := audit.New(path)
	if err != nil {
		t.Fatal(err)
	}

	entry := audit.Entry{
		Timestamp: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		Job:       "backup",
		Event:     audit.EventSuccess,
		Duration:  123.45,
	}
	if err := l.Log(entry); err != nil {
		t.Fatalf("Log: %v", err)
	}
	l.Close()

	f, _ := os.Open(path)
	defer f.Close()
	var got audit.Entry
	if err := json.NewDecoder(bufio.NewReader(f)).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.Job != "backup" {
		t.Errorf("job = %q, want backup", got.Job)
	}
	if got.Event != audit.EventSuccess {
		t.Errorf("event = %q, want success", got.Event)
	}
	if got.Duration != 123.45 {
		t.Errorf("duration = %f, want 123.45", got.Duration)
	}
}

func TestLog_SetsTimestampWhenZero(t *testing.T) {
	dir := t.TempDir()
	l, _ := audit.New(dir + "/a.log")
	defer l.Close()

	before := time.Now().UTC()
	_ = l.Log(audit.Entry{Job: "j", Event: audit.EventMiss})
	after := time.Now().UTC()

	f, _ := os.Open(dir + "/a.log")
	defer f.Close()
	var got audit.Entry
	_ = json.NewDecoder(f).Decode(&got)

	if got.Timestamp.Before(before) || got.Timestamp.After(after) {
		t.Errorf("timestamp %v not in [%v, %v]", got.Timestamp, before, after)
	}
}

func TestNew_InvalidPath(t *testing.T) {
	_, err := audit.New("/no/such/dir/audit.log")
	if err == nil {
		t.Error("expected error for invalid path, got nil")
	}
}

func TestLog_MultipleEntries(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/multi.log"

	l, err := audit.New(path)
	if err != nil {
		t.Fatal(err)
	}

	events := []audit.EventType{audit.EventSuccess, audit.EventMiss, audit.EventFailure}
	for _, ev := range events {
		if err := l.Log(audit.Entry{Job: "myjob", Event: ev}); err != nil {
			t.Fatalf("Log(%s): %v", ev, err)
		}
	}
	l.Close()

	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var entries []audit.Entry
	for scanner.Scan() {
		var e audit.Entry
		if err := json.Unmarshal(scanner.Bytes(), &e); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}
		entries = append(entries, e)
	}

	if len(entries) != len(events) {
		t.Fatalf("got %d entries, want %d", len(entries), len(events))
	}
	for i, e := range entries {
		if e.Event != events[i] {
			t.Errorf("entry[%d].Event = %q, want %q", i, e.Event, events[i])
		}
	}
}
