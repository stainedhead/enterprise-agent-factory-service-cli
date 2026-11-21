# Enterprise Agent Factory Service — CLI PRD (`-cli`)

| | |
|---|---|
| **Status** | Draft for review (2026-11-21) |
| **Author** | J (Enterprise Architecture), drafted with Claude |
| **Parent** | [Root PRD](https://github.com/stainedhead/enterprise-agent-factory-service/blob/main/enterprise-agent-factory-service-initial-PRD.md) (stories US-1..US-16; IF-2); solution design §11 |
| **Consumes** | [API contract](https://github.com/stainedhead/enterprise-agent-factory-service/blob/main/docs/contracts/api-contract.md) via `-restapi/api/openapi.yaml` at a release tag |
| **Provides** | The `afs` command and its machine-readable `--output json`, which pipelines rely on |

## 1. Purpose

`afs` lets people and automation do everything the console does, with the habits of the AWS CLI:
`afs <service> <operation> [--options]`, named profiles, SSO, `--output json|table|text|yaml`,
`--query`, pagination, `--cli-input-json` and `--generate-cli-skeleton`. One operation is one API call
(plus pagination). The CLI never decides access; the API does.

## 2. Users

| User | Needs |
|---|---|
| Service admins | Onboard teams, set guardrails, approve prices, manage catalogs from scripts |
| Team-admins | Create agents, secrets, configuration and tags; rotate; pause and restart |
| Viewers, FinOps, SRE | Read status, usage and reports; export CSV; pipe JSON |
| CI pipelines | Non-interactive auth with a team API key; stable JSON; exit codes |

## 3. Requirements owned

| Id | Requirement | Product reqs | Pri |
|---|---|---|---|
| CLI-1 | A command for every operation in the API contract §4, named by its `CLI` column, mapping to its `operationId` | IF-2, G6 | P0 |
| CLI-2 | Global options: `--profile`, `--endpoint-url`, `--output`, `--query`, `--no-paginate`, `--max-items`, `--page-size`, `--starting-token`, `--cli-input-json`, `--generate-cli-skeleton`, `--debug`, `--no-cli-pager`, `--color`, `--version` | IF-2 | P0 |
| CLI-3 | `afs configure` and `afs configure sso` write `~/.afs/config`; precedence: option → environment (`AFS_PROFILE`, `AFS_ENDPOINT_URL`, `AFS_DEFAULT_OUTPUT`, `AFS_API_KEY`) → profile → defaults | IF-2 | P0 |
| CLI-4 | `afs sso login` (OIDC device grant against the profile's provider, chosen from `sso list-providers`), `sso logout`, `sso get-caller-identity`; token cache `~/.afs/sso/cache` with 0600 | ACC-1 | P0 |
| CLI-5 | API key auth from `AFS_API_KEY` or `~/.afs/credentials` (0600) | ACC-4 | P1 |
| CLI-6 | Output renderers: `json` exactly the API's fields, `yaml`, `table` and `text` for people; `--query` JMESPath before rendering | IF-2 | P0 |
| CLI-7 | Auto-pagination; with `--max-items`, print `NextToken` for `--starting-token` | IF-2 | P0 |
| CLI-8 | Errors: `An error occurred (<code>) when calling the <Operation> operation: <message>` on stderr; exit 0 success, 252 usage, 254 service error, 255 other | IF-2 | P0 |
| CLI-9 | Secret values only from a file (`--secret-value file://…`, `fileb://…`), stdin (`-`) or a no-echo prompt; never a plain argument; never echoed, logged or in `--debug` | SEC-1, SEC-2 | P0 |
| CLI-10 | Updates send `If-Match` from a preceding read (or `--if-match`); a `412` prints the current version and exits 254 | CFG-1 | P0 |
| CLI-11 | Creates send an `Idempotency-Key` (generated, or `--idempotency-key`) and retry safely on network errors | API contract §1 | P0 |
| CLI-12 | `--wait` on mutations polls status until `Ready`, `Degraded` or `Failed` (timeout `--wait-timeout`, default 10 min) | TEAM-5 | P1 |
| CLI-13 | Reports: `--range`, `--from`, `--to`, `--group-by` (incl. `tag:<key>`), `--tag key=value` filters, `--output csv` for CSV from the API | RPT-1..9 | P0 |
| CLI-14 | Help for every service and operation, with at least one example, generated from the same definitions as the commands | IF-2 | P0 |
| CLI-15 | Distribution: signed binaries for macOS (arm64, amd64), Linux (arm64, amd64), Windows (amd64); Homebrew tap; checksums; `afs --version` shows version and API contract version | IF-2 | P1 |

### 3.1 Services and representative operations

| Service | Operations (full list from the API contract) |
|---|---|
| `teams` | `list-teams`, `create-team`, `describe-team`, `update-team`, `delete-team`, `suspend-team`, `resume-team`, `get-team-status`, `list-team-models`, `put-team-models`, `get-team-tags`, `put-team-tags`, `list-allowed-accounts`, `list-config-values`, `put-config-value`, `delete-config-value`, `list-external-systems`, `set-external-system-state`, `declare-external-system`, `list-log-destinations` …, `list-revisions`, `restore-revision` |
| `members`, `api-keys` | `list-members`, `add-member`, `update-member`, `remove-member`; `list-api-keys`, `create-api-key`, `delete-api-key` |
| `agents` | `list-agents`, `create-agent`, `describe-agent`, `update-agent`, `delete-agent`, `pause-agent`, `resume-agent`, `restart-agent`, `get-agent-status`, `list-snapshots`, `create-snapshot`, `list-role-links`, `create-role-link`, `verify-role-link`, `list-config-values`, `put-config-value`, `list-external-systems`, `set-external-system-state` |
| `secrets` | `list-secrets`, `create-secret`, `describe-secret`, `update-secret`, `put-secret-value`, `rotate-secret`, `delete-secret` |
| `budgets`, `reports` | `get-budget`, `put-budget`, `get-usage`, `list-costs`; `get-team-report`, `get-agent-report`, `get-service-report`, `list-report-schedules` … |
| `models`, `audit` | `list-models`; `list-audit-events`, `list-service-audit-events` |
| `admin` | settings, identity providers, service tags, models, prices, price proposals, catalog systems, storage classes, org accounts, cost reconciliations, exports and imports |
| `sso`, `configure` | `login`, `logout`, `list-providers`, `get-caller-identity`; `configure`, `configure sso` |
| `work-orders` | Phase 3 |

## 4. User stories and acceptance criteria

**CLI-US-1 · Sign in once, work all day.** *(US-2)* As a team-admin, I want `afs sso login` to sign me in
through my provider and keep me signed in for the session.
- Given a profile with an SSO start URL and provider, when I run `afs sso login`, then the device code and
  verification URL are printed, and after I approve, the token is cached with mode 0600.
- When the token expires, then the next command fails with exit 255 and the hint
  `Run afs sso login --profile <name>`.
- `afs sso get-caller-identity --output json` prints the API's `GetMe` response unchanged.

**CLI-US-2 · Script a team onboarding.** *(US-1)* As a service admin, I want to create a team with its
tags and first admin from a JSON file, so that onboarding is repeatable.
- `afs teams create-team --generate-cli-skeleton` prints a JSON skeleton with every field of
  `CreateTeam`; `--cli-input-json file://team.json` sends it unchanged.
- A missing required tag ends with
  `An error occurred (ValidationError) when calling the CreateTeam operation: …` and exit 254.
- `--wait` returns when the team is `Ready`, printing the final team.

**CLI-US-3 · Secrets never in shell history.** *(US-6)* As a team-admin, I want to create and rotate
secrets without the value appearing anywhere but the request body.
- `afs secrets create-secret --team-id … --name gh-token --secret-value file://token.txt` sends the file's
  content; `--secret-value -` reads stdin; with neither, a no-echo prompt asks for it.
- `--secret-value plaintext` (anything not `file://`, `fileb://` or `-`) is refused with exit 252 before
  any call.
- With `--debug`, the request log shows `"value": "[redacted]"` and `Authorization: [redacted]`.

**CLI-US-4 · Automation-friendly output.** *(US-11)* As a FinOps analyst, I want reports as JSON or CSV.
- `--output json` fields equal the API's field names and nesting exactly (golden tests per operation).
- `afs reports get-team-report --team-id … --kind costs --range 30d --group-by tag:project --output csv`
  writes the API's CSV byte for byte.
- `--query 'items[?state==\`active\`].slug'` applies before rendering for every output format.

**CLI-US-5 · Clear errors and exit codes.** *(US-2)* As a pipeline author, I want consistent exit codes.
- Usage errors (unknown flag, missing required option) → exit 252 with the operation's usage.
- Any API error → stderr AWS-style message and exit 254; `AccessDenied` includes the required
  permission from `details`.
- Network or config errors → exit 255.

**CLI-US-6 · Safe updates.** *(US-13)* As a team-admin, I want an update to fail rather than overwrite
someone else's change.
- `update-*` reads the resource, sends `If-Match` with its version, and on `412` prints
  `The resource changed (version N). Re-run to apply your change to the latest version.` with exit 254.

**CLI-US-7 · Every console action has a command.** *(G6)*
- A test compares the operations in `openapi.yaml` with the command table and fails on any operation
  without a command (other than health).

## 5. Non-functional requirements

| Area | Requirement |
|---|---|
| Startup | < 100 ms to first request on a warm machine |
| Compatibility | Works against the API contract major version it was built for; warns when the server's contract minor is newer |
| Security | No token, key or secret in stdout, stderr, logs or shell history; files written 0600 |
| Quality | `internal/...` ≥ 90% coverage; golden files for table, text and help; every operation tested for `AccessDenied` (exit 254) |
| Accessibility | `--color off` and `NO_COLOR` respected; table output readable at 80 columns |

## 6. Delivery

| Milestone | This repository |
|---|---|
| M0 | Generated client from `openapi.yaml`; config, profiles, output, errors; `sso` |
| M1 | `teams`, `members`, `api-keys` |
| M2 | `agents`, `secrets`, config values, external systems |
| M3 | `budgets`, `reports`, `models` |
| M4 | `admin` (incl. price proposals, service tags), export and import, signed releases |

## 7. Definition of done (every change)

- [ ] Command follows naming and option rules; help with an example
- [ ] `--output json` equals the API's fields; golden tests for table and text
- [ ] No secret value, token or key can reach output, logs or history
- [ ] AWS-style errors and exit codes; `AccessDenied` test
- [ ] `make gate` green; stories' acceptance criteria covered
