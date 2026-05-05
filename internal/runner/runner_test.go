package runner

import (
	"context"
	"strings"
	"testing"
	"time"
)

func TestRun_Success(t *testing.T) {
	r := New(5 * time.Second)
	res := r.Run(context.Background(), "echo-job", "echo hello")

	if res.Err != nil {
		t.Fatalf("expected no error, got %v", res.Err)
	}
	if res.ExitCode != 0 {
		t.Errorf("expected exit code 0, got %d", res.ExitCode)
	}
	if !strings.Contains(res.Output, "hello") {
		t.Errorf("expected output to contain 'hello', got %q", res.Output)
	}
	if res.JobName != "echo-job" {
		t.Errorf("expected job name 'echo-job', got %q", res.JobName)
	}
	if res.Duration <= 0 {
		t.Error("expected positive duration")
	}
}

func TestRun_Failure(t *testing.T) {
	r := New(5 * time.Second)
	res := r.Run(context.Background(), "fail-job", "exit 2")

	if res.ExitCode != 2 {
		t.Errorf("expected exit code 2, got %d", res.ExitCode)
	}
	if res.Err == nil {
		t.Error("expected non-nil error for failing command")
	}
}

func TestRun_Timeout(t *testing.T) {
	r := New(100 * time.Millisecond)
	res := r.Run(context.Background(), "slow-job", "sleep 10")

	if res.ExitCode == 0 {
		t.Error("expected non-zero exit code on timeout")
	}
	if res.Err == nil {
		t.Error("expected error on timeout")
	}
}

func TestRun_DefaultTimeout(t *testing.T) {
	r := New(0)
	if r.timeout != 30*time.Second {
		t.Errorf("expected default timeout 30s, got %v", r.timeout)
	}
}

func TestRun_CapturesOutput(t *testing.T) {
	r := New(5 * time.Second)
	res := r.Run(context.Background(), "output-job", "echo stdout; echo stderr >&2")

	if !strings.Contains(res.Output, "stdout") {
		t.Errorf("expected stdout in output, got %q", res.Output)
	}
	if !strings.Contains(res.Output, "stderr") {
		t.Errorf("expected stderr in output, got %q", res.Output)
	}
}
