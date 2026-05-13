package circuitbreaker

import (
	"fmt"

	"github.com/yourorg/cronwatch/internal/alerter"
)

// GuardedSender wraps an alerter.Sender and gates sends through a Breaker.
// When the circuit is open, sends are skipped to avoid hammering a broken
// downstream webhook endpoint.
type GuardedSender struct {
	inner   alerter.Sender
	breaker *Breaker
}

// NewGuardedSender wraps inner with circuit-breaker protection using b.
func NewGuardedSender(inner alerter.Sender, b *Breaker) *GuardedSender {
	return &GuardedSender{inner: inner, breaker: b}
}

// Send forwards the alert through the inner sender only when the circuit
// allows it. A successful send resets the breaker; any error trips it.
func (g *GuardedSender) Send(jobName, event, detail string) error {
	if !g.breaker.Allow() {
		return fmt.Errorf("circuit open: alert suppressed for job %q", jobName)
	}
	err := g.inner.Send(jobName, event, detail)
	if err != nil {
		g.breaker.RecordFailure()
		return err
	}
	g.breaker.RecordSuccess()
	return nil
}
