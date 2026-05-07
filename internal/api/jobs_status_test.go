package api

import (
	"testing"
	"time"
)

type fakeState struct {
	missCount int
	healthy   bool
}

func (f fakeState) GetMissCount() int { return f.missCount }
func (f fakeState) IsHealthy() bool   { return f.healthy }

func TestDeriveStatus_OK(t *testing.T) {
	s := fakeState{healthy: true, missCount: 0}
	if got := deriveStatus(s); got != "ok" {
		t.Errorf("expected ok, got %s", got)
	}
}

func TestDeriveStatus_Missing(t *testing.T) {
	s := fakeState{healthy: false, missCount: 2}
	if got := deriveStatus(s); got != "missing" {
		t.Errorf("expected missing, got %s", got)
	}
}

func TestDeriveStatus_Failed(t *testing.T) {
	s := fakeState{healthy: false, missCount: 0}
	if got := deriveStatus(s); got != "failed" {
		t.Errorf("expected failed, got %s", got)
	}
}

func TestNilIfZero_Zero(t *testing.T) {
	var zero time.Time
	if nilIfZero(zero) != nil {
		t.Error("expected nil for zero time")
	}
}

func TestNilIfZero_NonZero(t *testing.T) {
	now := time.Now()
	if nilIfZero(now) == nil {
		t.Error("expected non-nil for non-zero time")
	}
}
