// Package ratelimit — middleware wraps an alerter.Sender and gates sends
// through the Limiter so repeated alerts for the same job are suppressed.
package ratelimit

import (
	"fmt"

	"github.com/cronwatch/cronwatch/internal/alerter"
)

// GuardedSender wraps an alerter.Sender and suppresses duplicate alerts
// for the same job within the configured cooldown window.
type GuardedSender struct {
	inner   alerter.Sender
	limiter *Limiter
}

// NewGuardedSender returns a GuardedSender that delegates to inner only when
// the Limiter permits the alert for the given job key.
func NewGuardedSender(inner alerter.Sender, limiter *Limiter) *GuardedSender {
	return &GuardedSender{inner: inner, limiter: limiter}
}

// Send forwards the alert to the underlying sender if the rate limiter allows
// it. When suppressed, Send returns a descriptive error instead of nil so
// callers can distinguish a suppressed send from a delivery failure.
func (g *GuardedSender) Send(jobName, event, detail string) error {
	key := jobName + ":" + event
	if !g.limiter.Allow(key) {
		remaining := g.limiter.Remaining(key)
		return fmt.Errorf("ratelimit: alert for %q suppressed, retry in %s", key, remaining)
	}
	return g.inner.Send(jobName, event, detail)
}
