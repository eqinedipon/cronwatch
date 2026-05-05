package schedule_test

import (
	"testing"
	"time"

	"github.com/cronwatch/cronwatch/internal/schedule"
)

func mustParse(s string) time.Time {
	t, err := time.Parse("2006-01-02 15:04", s)
	if err != nil {
		panic(err)
	}
	return t
}

func TestNextRun(t *testing.T) {
	from := mustParse("2024-01-15 09:00")
	next, err := schedule.NextRun("30 10 * * *", from)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := mustParse("2024-01-15 10:30")
	if !next.Equal(expected) {
		t.Errorf("NextRun: got %v, want %v", next, expected)
	}
}

func TestNextRun_InvalidExpr(t *testing.T) {
	_, err := schedule.NextRun("not-a-cron", time.Now())
	if err == nil {
		t.Error("expected error for invalid cron expression")
	}
}

func TestPrevRun(t *testing.T) {
	// Every hour at :00 — previous run before 09:45 should be 09:00
	from := mustParse("2024-01-15 09:45")
	prev, err := schedule.PrevRun("0 * * * *", from)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := mustParse("2024-01-15 09:00")
	if !prev.Equal(expected) {
		t.Errorf("PrevRun: got %v, want %v", prev, expected)
	}
}

func TestIsDue_WithinGrace(t *testing.T) {
	// Job runs every hour; check 2 minutes after the hour — within 5 min grace
	at := mustParse("2024-01-15 10:02")
	due, err := schedule.IsDue("0 * * * *", at, 5*time.Minute)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !due {
		t.Error("expected job to be due within grace period")
	}
}

func TestIsDue_OutsideGrace(t *testing.T) {
	// Job runs every hour; check 10 minutes after the hour — outside 5 min grace
	at := mustParse("2024-01-15 10:10")
	due, err := schedule.IsDue("0 * * * *", at, 5*time.Minute)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if due {
		t.Error("expected job to be outside grace period")
	}
}
