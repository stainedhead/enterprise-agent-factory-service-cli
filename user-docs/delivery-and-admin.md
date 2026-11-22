# Keys, identities, tags, logs, schedules and service administration

Operations of milestone M4. Every command is `afs <service> <operation> [--options]`; see the
[command reference](reference/README.md) for every option. The API decides who may do what: a denied call
prints `An error occurred (AccessDenied) ...` and exits 254.

## API keys

```bash
afs api-keys create-api-key --team-id team_01J... --name ci-pipeline --role team_admin --expires-at 2027-02-19T18:00:00Z
afs api-keys list-api-keys --team-id team_01J...
afs api-keys delete-api-key --team-id team_01J... --key-id key_01J...
```

`--role` is `team_admin`, `team_collaborator` or `team_viewer` (a key with `team_collaborator` can call the team's Paperclip, see
[Team collaboration](collaboration.md)).

`create-api-key` prints the key's secret **once**, in the normal output, and says
`Store this key now; it cannot be shown again.` on stderr. `--debug` never prints it. If you repeat the
call with the same `--idempotency-key`, the service replays the response without the secret and the CLI says so:
delete the key and create another if it was lost.

## Agent role links and identities

```bash
afs agents create-role-link --team-id team_01J... --agent-id agt_01J... \
  --role-arn arn:aws:iam::111122223333:role/payments-deployer --profile-name deploy
afs agents describe-role-link --team-id team_01J... --agent-id agt_01J... --link-id arl_01J... \
  --query trustPolicy --output json          # the trust policy alone, to paste into IAM
afs agents verify-role-link --team-id team_01J... --agent-id agt_01J... --link-id arl_01J...
afs teams list-allowed-accounts --team-id team_01J...
```

A new link is `pending_trust` and shows the `externalId` the trust policy must carry. `verify-role-link` ends in
`verified`, or `failed` with the reason in `LastError`. An account outside the team's allowed accounts is a
`ValidationError` (254). `afs agents create-identity` / `list-identities` / `delete-identity` manage the agent's
source-hosting identities.

## Team tags

```bash
afs teams get-team-tags --team-id team_01J...                      # every tag, with its source
afs teams put-team-tags --team-id team_01J... --tags cost-center=CC-2002 project=checkout
```

`--tags` takes `key=value` pairs, JSON or `file://`. The table shows `system`, `service`, `team`,
`service_default` and `missing` tags; send team-owned keys only. A key with the `afs:` or `aws:` prefix, or a
service-owned key, is refused by the API (254). Like every update, `put-team-tags` sends `If-Match`, reading the
version when you do not give `--if-match`.

## Log destinations

```bash
afs teams create-log-destination --team-id team_01J... --name team-cloudwatch --kind cloudwatch \
  --config '{"logGroup":"/afs/payments","region":"us-east-1"}' --wait
afs teams describe-log-destination --team-id team_01J... --destination-id log_01J...
```

`describe-log-destination` shows the configuration and the delivery status. It never shows a credential: the
`SecretId` names the secret.

## Restore a revision

```bash
afs teams restore-revision --team-id team_01J... --revision-id 41 --dry-run   # print the plan, change nothing
afs teams restore-revision --team-id team_01J... --revision-id 41              # apply it
```

## Report schedules

```bash
afs reports create-report-schedule --team-id team_01J... --name weekly-costs --report costs \
  --cadence weekly --channel-id nch_01J... --range 30d --group-by tag:project
afs reports list-report-schedules --team-id team_01J...
```

`--range`, `--from`, `--to`, `--group-by` and `--tag key=value` build the stored query (the same flags as
`afs reports get-team-report`); `--query-json` takes the whole query object (`--query` is the global JMESPath
option). The output shows `NextRunAt`. A channel that belongs to the service or to another team is a
`ValidationError` (254).

## Service administration

```bash
afs audit list-service-audit-events --team-id team_01J... --outcome denied --output csv
afs admin list-notification-channels
afs admin create-notification-channel --name security-email --kind email --config '{"to":["security@example.com"]}' --event-types 'work_order.*'
afs admin put-service-tag --key cost-center --mode required --allowed-values CC-1001 CC-2002
afs admin list-price-proposals --status pending --output table
afs admin approve-price-proposal --proposal-id 1042 --note "Matches the AWS price list"
afs admin list-cost-reconciliations --status over_threshold
afs admin list-org-accounts
```

These need platform roles; a team role gets `AccessDenied` (254). Deleting a notification channel that a setting
names is a `Conflict` (254).

## Exports and imports

```bash
afs admin create-export --team-ids team_01J... team_01J...
afs admin describe-export --export-id exp_01J...          # sha256, sizeBytes and downloadUrl when finished
afs admin create-import --bundle file://bundle.json --account-map file://accounts.json            # dry run
afs admin create-import --bundle file://bundle.json --account-map file://accounts.json --no-dry-run
```

An import is a dry run unless you pass `--no-dry-run`. The request body is `{bundle, accountMap}`. A bundle file over
100 MB is refused (exit 252) before it is read.
