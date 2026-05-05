# cronwatch

Lightweight daemon that monitors cron job execution and sends alerts on missed or failed runs.

---

## Installation

```bash
go install github.com/yourname/cronwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourname/cronwatch.git && cd cronwatch && go build -o cronwatch .
```

---

## Usage

Define your monitored jobs in a `cronwatch.yaml` config file:

```yaml
jobs:
  - name: daily-backup
    schedule: "0 2 * * *"
    timeout: 30m
    alert:
      email: ops@example.com

  - name: hourly-sync
    schedule: "0 * * * *"
    timeout: 5m
    alert:
      slack_webhook: https://hooks.slack.com/services/your/webhook/url
```

Start the daemon:

```bash
cronwatch --config cronwatch.yaml
```

Wrap an existing cron command to report its exit status:

```bash
# In your crontab
0 2 * * * cronwatch run --job daily-backup -- /usr/local/bin/backup.sh
```

cronwatch will send an alert if the job exceeds its timeout, exits with a non-zero code, or fails to run within the expected window.

---

## Configuration

| Field | Description |
|---|---|
| `schedule` | Standard cron expression for the expected run time |
| `timeout` | Maximum allowed runtime before the job is flagged |
| `grace_period` | Extra time allowed after the scheduled window before a missed-run alert fires (default: `5m`) |
| `alert` | Notification target (email, Slack, or webhook) |

### Alert targets

| Key | Description |
|---|---|
| `email` | Send alert to the specified email address |
| `slack_webhook` | Post alert to a Slack incoming webhook URL |
| `webhook` | POST alert payload to an arbitrary HTTP endpoint |

---

## License

MIT © [yourname](https://github.com/yourname)
