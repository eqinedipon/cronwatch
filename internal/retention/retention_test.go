package retention

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeFile(t *testing.T, dir, name string, age time.Duration) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte("log"), 0o644); err != nil {
		t.Fatalf("writeFile: %v", err)
	}
	mod := time.Now().Add(-age)
	if err := os.Chtimes(path, mod, mod); err != nil {
		t.Fatalf("chtimes: %v", err)
	}
	return path
}

func TestClean_RemovesExpiredFiles(t *testing.T) {
	dir := t.TempDir()
	old := writeFile(t, dir, "audit-old.log", 10*24*time.Hour)
	writeFile(t, dir, "audit-new.log", 1*time.Hour)

	c := New(Policy{Dir: dir, MaxAge: 7 * 24 * time.Hour, Pattern: "audit-*.log"})
	removed, err := c.Clean()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(removed) != 1 || removed[0] != old {
		t.Errorf("expected %s removed, got %v", old, removed)
	}
	if _, err := os.Stat(old); !os.IsNotExist(err) {
		t.Errorf("expected old file to be deleted")
	}
}

func TestClean_RespectsMaxFiles(t *testing.T) {
	dir := t.TempDir()
	for i := 0; i < 5; i++ {
		age := time.Duration(i) * time.Hour
		writeFile(t, dir, filepath.Base(t.TempDir())+".log", age)
	}
	// Create named files we can reason about.
	writeFile(t, dir, "audit-a.log", 1*time.Hour)
	writeFile(t, dir, "audit-b.log", 2*time.Hour)
	writeFile(t, dir, "audit-c.log", 3*time.Hour)

	c := New(Policy{Dir: dir, MaxFiles: 2, Pattern: "audit-*.log"})
	removed, err := c.Clean()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(removed) == 0 {
		t.Errorf("expected files to be removed")
	}
}

func TestClean_NoFiles_ReturnsEmpty(t *testing.T) {
	dir := t.TempDir()
	c := New(Policy{Dir: dir, MaxAge: time.Hour, Pattern: "audit-*.log"})
	removed, err := c.Clean()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(removed) != 0 {
		t.Errorf("expected no removals, got %v", removed)
	}
}

func TestClean_DefaultPattern(t *testing.T) {
	dir := t.TempDir()
	c := New(Policy{Dir: dir, MaxAge: time.Hour})
	if c.policy.Pattern != "audit-*.log" {
		t.Errorf("expected default pattern, got %q", c.policy.Pattern)
	}
}

func TestClean_NoMaxAge_KeepsFiles(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "audit-keep.log", 30*24*time.Hour)

	c := New(Policy{Dir: dir, MaxAge: 0, MaxFiles: 0, Pattern: "audit-*.log"})
	removed, err := c.Clean()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(removed) != 0 {
		t.Errorf("expected no removals without constraints, got %v", removed)
	}
}
