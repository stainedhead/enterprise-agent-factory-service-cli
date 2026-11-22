# Budgets, usage, reports and CSV, notification channels, models and prices

## Budgets

```bash
afs budgets get-budget --team-id team_... --output table
afs budgets put-budget --team-id team_... --period monthly --amount-usd 100.50 --thresholds 50 80 100
```

`--amount-usd` is sent as a decimal string (`"100.50"`), `--thresholds` as numbers (percent of the
budget). An amount above the team's ceiling is the API's `ValidationError` (exit 254). The table shows
amount, spent, reserved, percent used and forecast. `put-budget` sends `If-Match`: pass `--if-match N`, or
let `afs` read the current version (it says so on stderr).

## Usage and costs

```bash
afs budgets get-usage --team-id team_... --from 2026-11-01T00:00:00Z --group-by agent     # agent | model | day
afs budgets list-costs --team-id team_... --output csv > costs.csv
```

Table output of `get-usage` has one column per group key (`key.agent`, `key.model`) plus the metrics.

## Reports

```bash
afs reports get-team-report --team-id team_... --kind costs --range 30d --group-by tag:project --output csv
afs reports get-team-report --team-id team_... --kind usage --range 7d --group-by agent,model --tag project=checkout
afs reports get-agent-report --team-id team_... --agent-id agt_... --kind activity --range 24h
afs reports get-service-report --kind costs --range mtd --group-by team          # platform roles
```

| Option | Meaning |
|---|---|
| `--kind` | `usage`, `costs`, `activity` or `traffic` (required) |
| `--range` | `24h`, `7d`, `30d`, `90d` or `mtd`; replaces `--from`/`--to` (giving both is exit 252) |
| `--from`, `--to` | RFC 3339 UTC timestamps |
| `--group-by` | comma-separated `team`, `agent`, `model`, `system`, `day`, `hour`, `pool` (ephemeral teams), `tag:<key>` |
| `--tag key=value` | filter by a tag; several pairs allowed; sent as `tag.<key>=<value>` |
| `--agent-id --model-id --system-id --pool-id --team-id` | filters (`--team-id` on the service report only) |

Invalid values exit 252 before any call. Table columns follow `--group-by`.

### CSV

With `--output csv`, `afs` sends `Accept: text/csv` and writes the API's CSV exactly as returned (line
endings included). A large report arrives in pages: `afs` follows the `X-Afs-Next-Token` header and writes
one header row followed by every row of every page, streaming page by page. CSV works for reports,
`get-usage`, `list-costs`, audit events and the other operations marked CSV in `afs <service> <operation> help`;
elsewhere, and with `--query`, `--output csv` exits 252.

## Notification channels

```bash
afs teams create-notification-channel --team-id team_... --name payments-slack --kind slack \
  --config '{"channel":"#payments-alerts"}' --event-types 'budget.*' 'agent.status_changed' --secret-id sec_... --wait
afs teams describe-notification-channel --team-id team_... --channel-id nch_...
afs teams test-notification-channel --team-id team_... --channel-id nch_...
```

`describe` prints the stored configuration and the last delivery status; credentials are referenced by
secret id and never shown. `test-notification-channel`:

| Result | Output | Exit |
|---|---|---|
| delivered | the result | 0 |
| not delivered | the result, then `An error occurred (DeliveryFailed) when calling the TestNotificationChannel operation: <error>` | 254 |
| still running (HTTP 202) | the result, and a note to run `describe-notification-channel` | 0 |

## Models and prices (platform administrators)

```bash
afs admin list-models
afs admin create-model --id anthropic.claude-sonnet-4-5 --provider anthropic --display-name "Claude Sonnet 4.5" \
  --context-window 200000 --max-output-tokens 64000
afs admin create-model-price --model-id anthropic.claude-sonnet-4-5 --effective-from 2026-12-01T00:00:00Z \
  --input-usd-per-mtok 3.00 --output-usd-per-mtok 15.00
afs admin list-model-prices --model-id anthropic.claude-sonnet-4-5
```

`--max-output-tokens` is required. Prices are USD per million tokens, given and sent as strings.
