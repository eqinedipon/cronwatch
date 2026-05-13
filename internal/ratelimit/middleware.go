package ratelimit

// Sender is the interface for anything that can deliver an alert message.
type Sender interface {
	Send(msg string) error
}

// GuardedSender wraps a Sender and suppresses duplicate alerts
// for the same message key within the configured cooldown window.
type GuardedSender struct {
	inner   Sender
	limiter *Limiter
}

// NewGuardedSender returns a GuardedSender that gates calls to inner
// using the provided Limiter. The message itself is used as the rate-limit key,
// so distinct messages are tracked independently.
func NewGuardedSender(inner Sender, limiter *Limiter) Sender {
	return &GuardedSender{
		inner:   inner,
		limiter: limiter,
	}
}

// Send delivers the message via the inner Sender only when the rate limiter
// allows it. Suppressed calls return nil without invoking the inner Sender.
func (g *GuardedSender) Send(msg string) error {
	if !g.limiter.Allow(msg) {
		return nil
	}
	return g.inner.Send(msg)
}
