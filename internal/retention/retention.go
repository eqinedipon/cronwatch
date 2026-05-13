package retention

import (
	"os"
	"path/filepath"
	"sort"
	"time"
)

// Policy defines how long audit log files are retained.
type Policy struct {
	// MaxAge is the maximum age of a log file before it is deleted.
	MaxAge time.Duration
	// MaxFiles is the maximum number of log files to keep (0 = unlimited).
	MaxFiles int
	// Dir is the directory containing audit log files.
	Dir string
	// Pattern is the glob pattern used to match log files.
	Pattern string
}

// Cleaner removes old audit log files according to a retention policy.
type Cleaner struct {
	policy Policy
	now    func() time.Time
}

// New creates a new Cleaner with the given policy.
func New(p Policy) *Cleaner {
	if p.Pattern == "" {
		p.Pattern = "audit-*.log"
	}
	return &Cleaner{policy: p, now: time.Now}
}

// Clean removes files that exceed the policy constraints.
// It returns the list of files removed and any error encountered.
func (c *Cleaner) Clean() ([]string, error) {
	matches, err := filepath.Glob(filepath.Join(c.policy.Dir, c.policy.Pattern))
	if err != nil {
		return nil, err
	}

	type entry struct {
		path    string
		modTime time.Time
	}

	var entries []entry
	for _, m := range matches {
		info, err := os.Stat(m)
		if err != nil {
			continue
		}
		entries = append(entries, entry{path: m, modTime: info.ModTime()})
	}

	// Sort newest first.
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].modTime.After(entries[j].modTime)
	})

	var removed []string
	for i, e := range entries {
		expired := c.policy.MaxAge > 0 && c.now().Sub(e.modTime) > c.policy.MaxAge
		tooMany := c.policy.MaxFiles > 0 && i >= c.policy.MaxFiles
		if expired || tooMany {
			if err := os.Remove(e.path); err == nil {
				removed = append(removed, e.path)
			}
		}
	}
	return removed, nil
}
