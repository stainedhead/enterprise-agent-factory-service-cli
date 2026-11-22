# Swarms: ephemeral teams, stages and pools

Operations of milestone M5. A swarm team is an ephemeral team that runs pools of workers through stages. The API
decides who may do what; a refusal prints `An error occurred (AccessDenied) ...` and exits 254.

## Create an ephemeral team

```bash
afs teams create-team --slug cve-sprint --display-name "CVE sprint" --kind ephemeral \
  --max-budget-usd 500.00 --first-admin '{"email":"ada@example.com"}' \
  --expires-at 2026-12-31T00:00:00Z --teardown-rule at_cap_or_expiry \
  --budget-usd 500.00 --budget-thresholds 50 80 100 --wait
```

`--kind ephemeral` needs `--expires-at`, `--teardown-rule` (`at_expiry`, `at_cap_or_expiry`, `at_work_complete`)
and `--budget-usd`; leaving one out is a usage error (exit 252) before any call. `--budget-usd` and
`--budget-thresholds` are sent as a lifetime budget `{period: lifetime, amountUsd, thresholds}`; an ephemeral
team's budget always blocks at the cap. The whole object can come from `--cli-input-json file://team.json`
(`afs teams create-team --generate-cli-skeleton` prints the shape) and is checked the same way.

## Stages and pools

```bash
afs pools create-stage --team-id team_01J... --name triage --position 1
afs pools list-stages --team-id team_01J...
afs pools create-pool --team-id team_01J... --name fixer --template claude-code \
  --image registry.example.com/agents/coder@sha256:... --model-id anthropic.claude-sonnet-4-5 --workers 10 --wait
afs pools describe-pool --team-id team_01J... --pool-id pool_01J... --output table
afs pools get-pool-status --team-id team_01J... --pool-id pool_01J... --output table
```

Pools exist only on ephemeral teams; `create-pool` on a persistent team is the API's `Conflict` (254).
`describe-pool` shows the desired worker count, state and phase; `get-pool-status` shows desired, reporting,
working, idle and impaired workers.

## Scale, pause and resume

```bash
afs pools scale-pool --team-id team_01J... --pool-id pool_01J... --workers 20 --wait
afs pools pause-pool --team-id team_01J... --pool-id pool_01J...
afs pools resume-pool --team-id team_01J... --pool-id pool_01J...
```

`--workers` is 0 to 200; a value outside the range exits 252 before any call. The request carries an
`Idempotency-Key` (yours with `--idempotency-key`, else a new one), so a retry does not scale twice. With `--wait`,
`scale-pool` returns 0 when every desired worker reports, 254 when the pool is `Failed` or `Degraded`, and 255 when
`--wait-timeout` (default 10 minutes) passes first; the change stands either way.

## Handoffs, secrets, configuration, external systems, role links

```bash
afs pools put-handoffs --team-id team_01J... --pool-id pool_01J... --to-pool-ids pool_01J... pool_01J...
afs pools put-secret-bindings --team-id team_01J... --pool-id pool_01J... --items file://bindings.json
afs pools put-config-value --team-id team_01J... --pool-id pool_01J... --name LOG_LEVEL --source literal --value debug
afs pools put-external-system --team-id team_01J... --pool-id pool_01J... --system-id ext_01J... --no-enabled
afs pools create-role-link --team-id team_01J... --pool-id pool_01J... --account-id 111122223333 --role-name afs-cve-reader --profile-name target
afs pools verify-role-link --team-id team_01J... --pool-id pool_01J... --link-id prl_01J...
```

The `put-*` commands replace the whole list and send `If-Match` (reading the version when you do not give
`--if-match`). Secret bindings carry secret ids and targets, never values.

## Cost-allocation tag status

```bash
afs admin list-cost-tag-status
```

Shows, per tag key, whether AWS has activated it as a cost-allocation tag and whether it has appeared in the cost report.
