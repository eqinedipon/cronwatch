package alerter

import (
	"fmt"

	"cronwatch/internal/config"
)

// FromConfig builds a Sender based on the application configuration.
// If no alerting config is provided, a LogSender is returned as the default.
func FromConfig(cfg *config.Config) (Sender, error) {
	if cfg.Alerting == nil {
		return &LogSender{}, nil
	}

	var senders []Sender

	if cfg.Alerting.Log {
		senders = append(senders, &LogSender{})
	}

	if cfg.Alerting.WebhookURL != "" {
		senders = append(senders, NewWebhookSender(cfg.Alerting.WebhookURL))
	}

	switch len(senders) {
	case 0:
		return nil, fmt.Errorf("alerting config present but no senders configured")
	case 1:
		return senders[0], nil
	default:
		return &MultiSender{Senders: senders}, nil
	}
}
