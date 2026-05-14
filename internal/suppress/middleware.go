package suppress

import (
	"context"
	"fmt"

	"github.com/cronwatch/cronwatch/internal/alerter"
)

// GuardedSender wraps an alerter.Sender and skips delivery when the
// (jobName, alertType) pair is within the active suppression window.
type GuardedSender struct {
	inner  alerter.Sender
	window *Window
}

// NewGuardedSender returns a GuardedSender that delegates to inner.
func NewGuardedSender(inner alerter.Sender, w *Window) *GuardedSender {
	return &GuardedSender{inner: inner, window: w}
}

// Send forwards the alert only if the suppression window has not been
// triggered for the same job + alert-type combination.
func (g *GuardedSender) Send(ctx context.Context, a alerter.Alert) error {
	key := fmt.Sprintf("%s:%s", a.JobName, a.Type)
	if g.window.IsSuppressed(key) {
		return nil
	}
	return g.inner.Send(ctx, a)
}
