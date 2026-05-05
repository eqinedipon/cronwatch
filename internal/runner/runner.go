package runner

import (
	"context"
	"os/exec"
	"time"
)

// Result holds the outcome of a cron job execution.
type Result struct {
	JobName  string
	ExitCode int
	Output   string
	Duration time.Duration
	Err      error
	Started  time.Time
}

// Runner executes shell commands for cron jobs.
type Runner struct {
	timeout time.Duration
}

// New creates a Runner with the given execution timeout.
func New(timeout time.Duration) *Runner {
	if timeout <= 0 {
		timeout = 30 * time.Second
	}
	return &Runner{timeout: timeout}
}

// Run executes the given command string in a shell and returns the result.
func (r *Runner) Run(ctx context.Context, jobName, command string) Result {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	start := time.Now()
	cmd := exec.CommandContext(ctx, "sh", "-c", command)

	out, err := cmd.CombinedOutput()
	duration := time.Since(start)

	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			exitCode = -1
		}
	}

	return Result{
		JobName:  jobName,
		ExitCode: exitCode,
		Output:   string(out),
		Duration: duration,
		Err:      err,
		Started:  start,
	}
}
