package metrics

import (
	"sync"
	"time"
)

// JobMetrics holds runtime statistics for a single cron job.
type JobMetrics struct {
	JobName      string
	TotalRuns    int64
	SuccessCount int64
	FailureCount int64
	MissCount    int64
	LastRunAt    time.Time
	LastDuration time.Duration
	AvgDuration  time.Duration
}

// Collector accumulates metrics across all tracked jobs.
type Collector struct {
	mu   sync.RWMutex
	jobs map[string]*JobMetrics
}

// New returns an initialised Collector.
func New() *Collector {
	return &Collector{
		jobs: make(map[string]*JobMetrics),
	}
}

// RecordRun records the outcome of a single job execution.
func (c *Collector) RecordRun(name string, success bool, duration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	m := c.getOrCreate(name)
	m.TotalRuns++
	m.LastRunAt = time.Now()
	m.LastDuration = duration

	// rolling average
	m.AvgDuration = time.Duration((int64(m.AvgDuration)*(m.TotalRuns-1) + int64(duration)) / m.TotalRuns)

	if success {
		m.SuccessCount++
	} else {
		m.FailureCount++
	}
}

// RecordMiss increments the miss counter for a job.
func (c *Collector) RecordMiss(name string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.getOrCreate(name).MissCount++
}

// Get returns a copy of the metrics for the given job.
// The second return value is false when the job is unknown.
func (c *Collector) Get(name string) (JobMetrics, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	m, ok := c.jobs[name]
	if !ok {
		return JobMetrics{}, false
	}
	return *m, true
}

// All returns a snapshot of metrics for every known job.
func (c *Collector) All() []JobMetrics {
	c.mu.RLock()
	defer c.mu.RUnlock()
	out := make([]JobMetrics, 0, len(c.jobs))
	for _, m := range c.jobs {
		out = append(out, *m)
	}
	return out
}

// getOrCreate must be called with c.mu held for writing.
func (c *Collector) getOrCreate(name string) *JobMetrics {
	if m, ok := c.jobs[name]; ok {
		return m
	}
	m := &JobMetrics{JobName: name}
	c.jobs[name] = m
	return m
}
