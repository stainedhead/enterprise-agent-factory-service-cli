# Enterprise Agent Factory Service — CLI PRD (`-cli`)

| | |
|---|---|
| **Status** | Draft for review, feature-complete scope (2026-11-21); revised after `review-prd` |
| **Author** | J (Enterprise Architecture), drafted with Claude |
| **Parent** | [Root PRD](https://github.com/stainedhead/enterprise-agent-factory-service/blob/main/enterprise-agent-factory-service-initial-PRD.md) (stories US-1..US-18; IF-2; decisions D23–D46); solution design §11 |
| **Consumes** | [API contract](https://github.com/stainedhead/enterprise-agent-factory-service/blob/main/docs/contracts/api-contract.md) via `-restapi/api/openapi.yaml` at a release tag |
| **Provides** | The `afs` command and its machine-readable `--output json`, which pipelines rely on |
| **Precedence** | The API contract is the single source for operation names, CLI names, permissions and enums. Where solution design §11 differs (services `identities`, `usage`, `logs`, `notifications`; `--secret-string`; `admin export-config`), the contract wins and §11 is corrected in the root |

## 1. Purpose

`afs` lets people and automation do everything the console does, with the habits of the AWS CLI:
`afs <service> <operation> [--options]`, named profiles, SSO, `--output json|table|text|yaml|csv`,
`--query`, pagination, `--cli-input-json` and `--generate-cli-skeleton`. One operation is one API call
plus pagination; the only other sanctioned extra calls are the read before an unguarded update (CLI-10)
and status polling under `--wait` (CLI-12). The CLI never decides access; the API does.

## 2. Users

| User | Needs |
|---|---|
| Service admins | Onboard teams, set guardrails, approve prices, manage catalogs and settings from scripts |
| Team-admins | Create agents, secrets, configuration and tags; rotate; pause and restart |
| Viewers, FinOps, SRE | Read status, usage and reports; export CSV; pipe JSON; read-only kubectl |
| CI pipelines | Non-interactive auth with a team API key; stable JSON; exit codes |
| Security operators and security tooling | Create and track CVE work orders, by person (SSO) or with a service-level API key |
| Agents and pool workers | Report work-order progress with their projected ServiceAccount token (D29) |

## 3. Requirements owned

Scope is feature complete (root D18): every operation in the API contract gets a command. Each
command's priority is its operation's priority in the contract.

| Id | Requirement | Product reqs | Pri |
|---|---|---|---|
| CLI-1 | A command for every operation in the API contract §4, named exactly by its `CLI` column and mapping to its `operationId`; priority per operation from the contract. Commands take ids (`--agent-id agt_…`), never names, because names are unique only within a team (D25); the CLI does not resolve names | IF-2, G6 | P0 |
| CLI-2 | Global options on every command: `--profile`, `--endpoint-url`, `--output`, `--query`, `--no-paginate`, `--max-items`, `--page-size`, `--starting-token`, `--cli-input-json`, `--generate-cli-skeleton`, `--debug`, `--no-cli-pager`, `--color on\|off\|auto`, `--cli-connect-timeout`, `--cli-read-timeout`, `--version` | IF-2 | P0 |
| CLI-3 | `afs configure` and `afs configure sso` write `~/.afs/config` (0600). Profile keys: `endpoint_url`, `output`, `sso_provider_id`, `sso_issuer`, `sso_client_id`, `team_id` (default for `--team-id`). Precedence: option → environment (`AFS_PROFILE`, `AFS_ENDPOINT_URL`, `AFS_DEFAULT_OUTPUT`, `AFS_API_KEY`) → profile → defaults (`output json`) | IF-2 | P0 |
| CLI-4 | `afs sso login` (OIDC device grant against the profile's provider, chosen from `sso list-providers`), `sso logout`, `sso get-caller-identity`; token cache `~/.afs/sso/cache` (0600). A session lasts as long as the IdP's refresh token; the CLI refreshes silently when one exists. A token known to be expired locally → exit 255 without a call and the hint `Run afs sso login --profile <name>`; a server `401 Unauthenticated` → AWS-style error, exit 254, same hint (D41) | ACC-1 | P0 |
| CLI-5 | Credential resolution, first match wins: `AFS_API_KEY` (wins even over `--profile`'s SSO token; `--debug` names the source) → the profile's `api_key` in `~/.afs/credentials` → the profile's cached SSO token → the workload token file `AFS_API_TOKEN_FILE` (default `/var/run/secrets/afs/api/token` when it exists; re-read per request, D29). A credentials file or token cache file with a mode wider than 0600 is refused with exit 255 and a `chmod 600` hint | ACC-4 | P1 |
| CLI-6 | Output renderers: `json` exactly the API's fields, `yaml`, `table` and `text` for people, and `csv` (D39). `--query` (JMESPath) applies before rendering. `csv` asks the API for `Accept: text/csv` and writes the body unchanged; it is valid only on operations the contract marks *CSV* (else exit 252); `--query` with `csv` exits 252; auto-pagination concatenates pages keeping one header row | IF-2, RPT-1..9 | P0 |
| CLI-7 | Auto-pagination by default; `--page-size` sets `maxResults` (1–100, else 252); `--max-items N` stops after N items and prints `NextToken`; `--starting-token` resumes; `--no-paginate` returns one page with its `nextToken` | IF-2 | P0 |
| CLI-8 | Errors on stderr: `An error occurred (<code>) when calling the <operationId> operation: <message>`, then `Request ID: <X-Request-Id>`. Exit 0 success, 130 interrupted (Ctrl-C), 252 usage or client-side validation, 254 any API error, 255 other (network, config, local expired token, `--wait` timeout) | IF-2 | P0 |
| CLI-9 | Secret values only from a file (`--secret-value file://…`, `fileb://…`), stdin (`-`) or a no-echo prompt; never a plain argument; never echoed, logged or in `--debug`. `--cli-input-json` for `secrets create-secret` and `secrets put-secret-value` must be `file://…`; inline JSON with a `value` field is refused with exit 252 before any call | SEC-1, SEC-2 | P0 |
| CLI-10 | Safe updates (D40). Every `PUT`, `PATCH` and `DELETE` on an existing resource sends `If-Match`. `--if-match <version>` (or `version` in `--cli-input-json` from a prior `describe-*`) is the safe mode and is sent as given. Without it, the CLI reads the resource, sends its version, and prints `Applied to version N; use --if-match to guard against concurrent changes.` to stderr; if the read returns `404` on a `PUT` that creates, it sends no `If-Match`. `412` prints the current version and exits 254; `428 PreconditionRequired` exits 254 | CFG-1 | P0 |
| CLI-11 | Creates and actions send an `Idempotency-Key` (a ULID per invocation, or `--idempotency-key`), kept for service setting `api.idempotency_key_ttl_seconds` (default 1 hour). Retries: at most 3 attempts in total, exponential backoff from 200 ms doubling with full jitter, `Retry-After` honoured when longer. Retried: network errors and `503` on `GET` and on requests carrying an `Idempotency-Key`; connect failures on any request; `429`; `409` with `idempotency_in_progress` (same key). Never retried: `idempotency_mismatch` (exit 254) and other 4xx. A replay prints the stored response | IF-1, API contract §1 | P0 |
| CLI-12 | `--wait` on mutations of desired-state resources and their actions polls `GET …/status` every 5 s (or `Retry-After`) until `status.phase` is `Ready` (exit 0), `Failed` or `Degraded` (exit 254, conditions printed to stderr), or `--wait-timeout` passes (default 10 min, exit 255). `delete-* --wait` waits for `404` (exit 0). `pools scale-pool --wait` waits until reporting workers equal desired workers (`GetPoolStatus`). Ctrl-C stops waiting, says the mutation stands, exits 130 | TEAM-5, AGT-8 | P1 |
| CLI-13 | Reports: `--range` (values as in the contract: `24h\|7d\|30d\|90d\|mtd`), `--from`, `--to`, `--group-by` (incl. `tag:<key>` and `pool`), `--tag key=value`, `--pool-id` filters; `--output csv` per CLI-6 | RPT-1..9 | P0 |
| CLI-14 | Help: `afs help`, `afs <service> help`, `afs <service> <operation> help`, generated from the same command table as the commands, each operation with at least one example | IF-2 | P0 |
| CLI-15 | Distribution: signed binaries for macOS (arm64, amd64), Linux (arm64, amd64), Windows (amd64); Homebrew tap; `checksums.txt` signed with cosign; `afs --version` prints `afs-cli/<version> contract/<major.minor> <go version> <os>/<arch>` | IF-2 | P1 |
| CLI-16 | `afs teams get-kubeconfig --team-id … [--kubeconfig PATH] [--aws-profile NAME]`: the API returns the kubeconfig content (context `afs-<team-slug>`, exec credential plugin `aws eks get-token` with the role the team's IdP federates to, no token). The CLI merges it as returned into `PATH` (default: first path of `$KUBECONFIG`, else `~/.kube/config`), replacing only entries with the same names, writing atomically. `--aws-profile` adds `AWS_PROFILE=<NAME>` to the exec `env` block. A new file is created 0600; an existing file keeps its mode, with a warning when wider than 0600. `--kubeconfig -` writes the API's YAML unchanged to stdout. A missing `aws` v2 on `PATH` is a warning, not an error | ACC-5 | P2 |
| CLI-17 | Swarm commands: `pools` service (stages, pools, scale, handoffs, pool secret bindings, configuration, external systems, role links) and `swarm-templates` service. `teams create-team --kind ephemeral` requires `--expires-at`, `--teardown-rule` and `--budget-usd` (sent as a lifetime budget with `onExhaust: block`; optional `--budget-thresholds 50 80 100`); the full object is available through `--cli-input-json` | AGT-8, TEAM-7 | P2 |
| CLI-18 | `work-orders` service, including `--repos` shorthand `host/owner/repo=<owners>` where owners are a team id, an email or both comma-separated (`=team_…,ops@example.com`): contains `@` → `owningContact`, starts `team_` → `owningTeamId`, anything else (such as a slug) → exit 252 before any call. No slug lookup (D38); `--output csv` on `list-work-orders` | WO-1..6 | P2 |

### 3.1 Services summary

**The API contract's §4 `CLI` column is the authoritative command list.** This table summarizes it per
service so the shape of each service is visible here; when the two differ, the contract wins and this
table is corrected. CLI-US-7 tests the generated command table against `openapi.yaml`, not this table.

| Service | Operations |
|---|---|
| `teams` | `list-teams`, `create-team`, `describe-team`, `update-team`, `delete-team`, `suspend-team`, `resume-team`, `get-team-status`, `get-kubeconfig`; models `list-team-models`, `put-team-models`; tags `get-team-tags`, `put-team-tags`; `list-allowed-accounts`; config `list-config-values`, `put-config-value`, `delete-config-value`; external systems `list-external-systems`, `set-external-system-state`, `declare-external-system`, `update-external-system`, `delete-external-system`; log destinations `list-log-destinations`, `create-log-destination`, `update-log-destination`, `delete-log-destination`; notification channels `list-notification-channels`, `create-notification-channel`, `update-notification-channel`, `delete-notification-channel`; revisions `list-revisions`, `describe-revision`, `restore-revision` |
| `members`, `api-keys` | `list-members`, `add-member`, `update-member`, `remove-member`; `list-api-keys`, `create-api-key`, `delete-api-key` |
| `agents` | `list-agents`, `create-agent`, `describe-agent`, `update-agent`, `delete-agent`, `pause-agent`, `resume-agent`, `restart-agent`, `get-agent-status`, `list-snapshots`, `create-snapshot`, `list-storage-classes`; role links `list-role-links`, `create-role-link`, `describe-role-link`, `delete-role-link`, `verify-role-link`; identities `list-identities`, `create-identity`, `delete-identity`; config `list-config-values`, `put-config-value`, `delete-config-value`; external systems `list-external-systems`, `set-external-system-state` |
| `secrets` | `list-secrets`, `create-secret`, `describe-secret`, `update-secret`, `put-secret-value`, `rotate-secret`, `delete-secret` |
| `budgets` | `get-budget`, `put-budget`, `get-usage`, `list-costs` |
| `reports` | `get-team-report`, `get-agent-report`, `get-service-report`; schedules `list-report-schedules`, `create-report-schedule`, `update-report-schedule`, `delete-report-schedule` |
| `models`, `audit` | `list-models`; `list-audit-events`, `list-service-audit-events` |
| `admin` | settings `get-settings`, `put-setting`; identity providers `list-identity-providers`, `create-identity-provider`, `update-identity-provider`, `delete-identity-provider`; service tags `list-service-tags`, `put-service-tag`, `delete-service-tag`; models `list-models`, `create-model`, `update-model`; prices `list-model-prices`, `create-model-price`, `list-price-proposals`, `approve-price-proposal`, `reject-price-proposal`; `list-cost-reconciliations`, `list-cost-tag-status`; catalog `list-external-systems`, `create-external-system`, `update-external-system`, `delete-external-system`; storage classes `list-storage-classes`, `put-storage-class`, `delete-storage-class`; `list-org-accounts`; service-level keys `list-api-keys`, `create-api-key`, `delete-api-key`; security operators `list-security-operators`, `add-security-operator`, `remove-security-operator`; `create-export`, `describe-export`, `create-import` |
| `pools` | stages `list-stages`, `create-stage`, `update-stage`, `delete-stage`; `list-pools`, `create-pool`, `describe-pool`, `update-pool`, `delete-pool`, `scale-pool`, `pause-pool`, `resume-pool`, `get-pool-status`; `list-handoffs`, `put-handoffs`; `list-secret-bindings`, `put-secret-bindings`; config `list-config-values`, `put-config-value`, `delete-config-value`; `list-external-systems`, `put-external-system`; role links `list-role-links`, `create-role-link`, `delete-role-link`, `verify-role-link` |
| `swarm-templates` | `list-swarm-templates`, `create-swarm-template`, `describe-swarm-template`, `update-swarm-template`, `delete-swarm-template` (enable and disable: see OQ-C3) |
| `work-orders` | `list-work-orders`, `create-work-order`, `describe-work-order`, `update-work-order`, `cancel-work-order`, `list-repos`, `update-repo`, `list-team-work-orders` |
| `sso`, `configure` | `login`, `logout`, `list-providers`, `get-caller-identity`; `configure`, `configure sso` (local only, no API operation) |

### 3.2 Non-goals

- Client-side workflows that combine several mutations (for example "create team and its agents").
- Name or slug resolution (`--agent-name`, `=team-slug`); commands take ids (D25, D38).
- Any command that reads a secret value: no permission reads one.
- Client-side access decisions: the CLI never skips a call because it predicts a `403`.
- An interactive TUI, plugins or aliases; auto-update (releases go through package managers).
- Local state beyond `~/.afs/config`, `~/.afs/credentials`, the SSO cache and the `GetMe` cache.

## 4. User stories and acceptance criteria

**CLI-US-1 · Sign in once, work all day.** *(US-2; ACC-1; D41)* As a team-admin, I want `afs sso login`
to sign me in through my team's provider and keep me signed in until the IdP's refresh token expires.
- Given a profile with `sso_provider_id`, `sso_issuer` and `sso_client_id`, when I run `afs sso login`,
  then the device code and verification URL are printed, and after I approve, the token is cached with
  mode 0600.
- Given a refresh token, when the access token expires, then the next command refreshes it silently and
  succeeds with no prompt.
- Given a cached token known to be expired with no usable refresh token, then the command makes no call
  and exits 255 with `Run afs sso login --profile <name>`; given a server `401`, then the AWS-style
  `Unauthenticated` error, the same hint and exit 254.
- `afs sso get-caller-identity --output json` prints the API's `GetMe` response unchanged.
- `afs sso logout` deletes the profile's cache files; the next command exits 255 with the hint.

**CLI-US-2 · Script a team onboarding.** *(US-1)* As a service admin, I want to create a team with its
tags and first admin from a JSON file, so that onboarding is repeatable.
- `afs teams create-team --generate-cli-skeleton` prints a JSON skeleton with every field of
  `CreateTeam`; `--cli-input-json file://team.json` sends it unchanged.
- A missing required tag ends with
  `An error occurred (ValidationError) when calling the CreateTeam operation: <message>` and exit 254.
- `--wait` returns when the team's phase is `Ready`, printing the final team; `Failed` exits 254 with
  its conditions on stderr.

**CLI-US-3 · Secrets never in shell history.** *(US-5, US-6)* As a team-admin, I want to create and rotate
secrets without the value appearing anywhere but the request body.
- `afs secrets create-secret --team-id team_… --name gh-token --secret-value file://token.txt` sends the
  file's content; `--secret-value -` reads stdin; with neither, a no-echo prompt asks for it.
- `--secret-value plaintext` (anything not `file://`, `fileb://` or `-`) is refused with exit 252 before
  any call.
- `--cli-input-json '{"name":"gh-token","value":"…"}'` (inline) is refused with exit 252 before any call
  and a message to use `file://`; `--cli-input-json file://secret.json` is accepted.
- With `--debug`, the request log shows `"value": "[redacted]"` and `Authorization: [redacted]`.

**CLI-US-4 · Automation-friendly output.** *(US-11)* As a FinOps analyst, I want reports as JSON or CSV.
- `--output json` fields equal the API's field names and nesting exactly (golden tests per operation).
- `afs reports get-team-report --team-id team_… --kind costs --range 30d --group-by tag:project --output
  csv` sends `Accept: text/csv` and writes the API's CSV byte for byte.
- `afs audit list-audit-events --team-id team_… --output csv` over 3 pages writes one header row and every
  row of the 3 pages.
- `--output csv` on `agents list-agents` (not *CSV* in the contract) exits 252; `--output csv --query …`
  exits 252.
- `--query 'items[?state==\`active\`].slug'` applies before rendering for `json`, `yaml`, `table` and
  `text`.

**CLI-US-5 · Clear errors and exit codes.** *(IF-2)* As a pipeline author, I want consistent exit codes.
- Usage errors (unknown flag, missing required option) → exit 252 with the operation's usage.
- Any API error → stderr AWS-style message naming the `operationId`, then `Request ID: <id>`, exit 254;
  `AccessDenied` includes `details.requiredPermission`.
- Network or config errors → exit 255. Ctrl-C during a request cancels it and exits 130.

**CLI-US-6 · Safe updates.** *(CFG-1; D40)* As a team-admin, I want an update to fail rather than
overwrite someone else's change.
- `afs agents update-agent --if-match 7 …` sends `If-Match: "7"`; when the agent is at version 8, it prints
  `The resource changed (version 8). Re-run with --if-match 8 to apply your change to the latest
  version.` and exits 254.
- Without `--if-match`, the CLI reads, sends the read version and prints
  `Applied to version N; use --if-match to guard against concurrent changes.` to stderr.
- `put-budget`, `put-team-tags` and `delete-agent` send `If-Match` like `update-*`; a server `428
  PreconditionRequired` (for example from a stale CLI) exits 254.

**CLI-US-7 · Every operation has a command, and every command an operation.** *(G6)*
- A test compares the operations in `openapi.yaml` with the generated command table and fails on any
  operation without a command (other than `GetHealth`, `GetReadiness`), on any command without an
  operation (other than `configure`, `configure sso`, `sso login`, `sso logout` and `help`), and on any
  command whose name differs from the contract's `CLI` column for its `operationId`.

**CLI-US-8 · Run a budgeted swarm from a script.** *(US-15, US-17)* As a team-admin of an ephemeral
team, I want to create and scale pools from the command line.
- `afs teams create-team --kind ephemeral` without `--expires-at`, `--teardown-rule` or `--budget-usd`
  exits 252 before any call.
- `afs pools scale-pool --workers 250` exits 252 (range 0–200) before any call; `--workers 20` sends
  `ScalePool` with an `Idempotency-Key` and, with `--wait`, returns 0 when `GetPoolStatus` reports 20
  reporting workers out of 20 desired; on `--wait-timeout` it exits 255.
- `afs pools describe-pool --output table` shows desired, reporting, working, idle and impaired workers.
- `afs pools create-pool` on a persistent team prints the API's `Conflict` and exits 254.

**CLI-US-9 · Security tooling files work orders.** *(US-16; D38)* As a security engineer, I want my
scanner to create a work order with a service-level key.
- With `AFS_API_KEY` set to a security operator key, `afs work-orders create-work-order --cve-id
  CVE-2026-41873 --severity critical --sla-due-at 2026-12-01T00:00:00Z --template-id swt_… --repos
  github.com/payments/ledger-core=team_01JD… github.com/vendor/edi=edi@example.com
  github.com/payments/ops=team_01JE…,ops@example.com` sends one `CreateWorkOrder` with each repository's
  `owningTeamId` and `owningContact` set as classified; a repeat with the same `--idempotency-key`
  replays the stored response.
- `--repos github.com/payments/ledger-core=payments-factory` (a slug) exits 252 before any call, saying
  to use a team id or an email.
- A second open work order for the CVE prints the AWS-style `Conflict` error and exits 254.
- The same key on a team command gets `AccessDenied` (exit 254).

**CLI-US-10 · Workers report progress.** *(US-16; D29)* As a remediation pool worker, I want to report a
repository's PR with the token already mounted in my pod.
- Given `/var/run/secrets/afs/api/token` and no other credential, `afs work-orders update-repo
  --work-order-id wo_… --repo github.com/payments/ledger-core --status pr_opened --pr-url https://…`
  sends the token as `Authorization: Bearer` and succeeds; `sso get-caller-identity` shows
  `type: workload`.
- The same token on `work-orders cancel-work-order` gets `AccessDenied` (exit 254).

**CLI-US-11 · kubectl for viewers.** *(US-18; ACC-5)* As a team viewer, I want read-only kubectl without
handling tokens.
- `afs teams get-kubeconfig --team-id team_…` merges the API's context `afs-<team-slug>` into
  `~/.kube/config` without removing or changing other contexts; a new file is created 0600; the file has
  no token, only the exec block calling `aws eks get-token`.
- `--kubeconfig ./kc.yaml` writes there and leaves `~/.kube/config` untouched; `--aws-profile ops` adds
  `AWS_PROFILE=ops` to the exec `env` block; running it twice leaves one context.
- When the team has viewer kubectl off, the API's `AccessDenied` is printed and exit is 254.

**CLI-US-12 · Stand up an agent.** *(US-4)* As a team-admin, I want to create an agent from a script.
- `afs agents create-agent --team-id team_… --name reviewer --image registry/x@sha256:… --model-id …
  --volume '{"sizeGiB":20,"storageClass":"gp3"}' --wait` returns 0 when the agent's phase is `Ready`.
- `--replicas` is an unknown option (exit 252): the agent has no replica count.
- An image referenced by tag is refused by the API with `ValidationError` (exit 254).
- Two teams can each run `create-agent --name reviewer`; every later command addresses the agent by
  `--agent-id`.

**CLI-US-13 · Link an agent to another account.** *(US-7)* As a team-admin, I want to hand an account
owner a trust policy.
- `afs agents create-role-link --team-id team_… --agent-id agt_… --role-arn arn:aws:iam::…:role/x
  --profile-name deploy` prints the link with `status: pending_trust` and `externalId`.
- `afs agents describe-role-link … --query trustPolicy --output json` prints the trust policy document
  alone, ready to hand over.
- `afs agents verify-role-link …` prints the link with `status` `verified`, or `failed` with the reason;
  an account outside the team's allowed accounts gets `ValidationError` or `AccessDenied` from the API
  (exit 254).

**CLI-US-14 · Set a budget.** *(US-9)* As a team-admin, I want to set my team's budget and thresholds.
- `afs budgets put-budget --team-id team_… --period monthly --amount-usd 100.50 --thresholds 50 80 100`
  sends `amountUsd: "100.50"` as a decimal string and `thresholds: [50, 80, 100]`.
- An amount above the team's ceiling prints the API's `ValidationError` and exits 254.
- `afs budgets get-budget --output table` shows amount, spent, reserved, percent used and forecast.

**CLI-US-15 · Keep prices and settings current.** *(US-10; D45)* As a service admin, I want to review
price proposals and settings from scripts.
- `afs admin list-price-proposals --status pending --output table` lists one row per pending proposal.
- `afs admin approve-price-proposal --proposal-id mpp_… --note "list price change"` prints the decided
  proposal; `reject-price-proposal` likewise; a team-admin gets `AccessDenied` (exit 254).
- `afs admin put-setting --key api.cors_origins --value '["https://console.example.com"]'` sends the JSON
  value; `admin get-settings --output json` prints every setting from the root PRD's settings table.

**CLI-US-16 · Tag a team.** *(US-12)* As a team-admin, I want to set my team's tags.
- `afs teams put-team-tags --team-id team_… --tags project=ledger-modernization business-unit=payments`
  sends team-owned keys only.
- `afs teams get-team-tags --output table` shows every tag with its `source`, including `system` and
  `service` tags, and `missing` for a required key with no value.
- A key with the `afs-` or `aws:` prefix, or a service key, is refused by the API (exit 254).

**CLI-US-17 · Schedule a report.** *(US-11; RPT-8; D23)* As a team-admin, I want a weekly cost report sent
to my team's channel.
- `afs reports create-report-schedule --team-id team_… --name weekly-costs --report costs --cadence
  weekly --channel-id nch_… --range 30d --group-by tag:project` stores the query and prints the schedule
  with `nextRunAt`.
- A service-level or another team's channel gets `ValidationError` (exit 254); a viewer gets
  `AccessDenied` (exit 254).

**CLI-US-18 · Export and import configuration.** *(US-13; CFG-5)* As a service admin, I want to move an
environment by script.
- `afs admin create-export` prints the export with its id; `admin describe-export --export-id exp_…`
  shows its status and, once complete, the bundle reference; the bundle holds no secret value.
- `afs admin create-import --bundle file://bundle.json --account-map file://map.json --dry-run` sends
  `dryRun: true`, prints the plan and changes nothing; without `--dry-run` it applies.

## 5. Non-functional requirements

| Area | Requirement |
|---|---|
| Startup | < 100 ms to first request on a warm machine |
| Reliability | Connect timeout 30 s and read timeout 60 s by default (`--cli-connect-timeout`, `--cli-read-timeout`; 0 = none). Retries per CLI-11. Ctrl-C cancels the in-flight request through its context and exits 130. `--wait` polls every 5 s |
| Observability | `--debug` writes to stderr one line per event, `<RFC 3339 time> DEBUG <component> <message>`, with method, URL, status, duration, `X-Request-Id`, credential source and retry attempt. Redacted: `Authorization`, cookies, and body fields `value`, `secretKey`, `token`, `refreshToken`, `clientSecret`, `password`. Every error prints `Request ID: <X-Request-Id>`. `User-Agent: afs-cli/<version> (<os>/<arch>)` |
| Compatibility | Works against the API contract major version it was built for; reads `X-Afs-Contract-Version` from responses and warns once per day on stderr when the server's minor is newer (OQ-C1) |
| Security | No token, key or secret value in stdout, stderr, logs or shell history; files written 0600. **One exception:** `api-keys create-api-key` and `admin create-api-key` print the response's `secretKey` once to stdout (the only time it exists), with `Store this key now; it cannot be shown again.` on stderr; it is redacted in `--debug` |
| Quality | `internal/...` ≥ 90% coverage; golden files for table, text and help; every operation tested for `AccessDenied` (exit 254) |
| Accessibility | `--color off` and `NO_COLOR` respected; table output readable at 80 columns |

### 5.1 Acceptance for the cross-cutting requirements

- **Global options (CLI-2).** One table-driven test per option: `--profile p2` uses `p2`'s endpoint and
  credential; `--endpoint-url` overrides the profile; `--no-cli-pager` writes without a pager;
  `--color on|off|auto` (auto: off when stdout is not a terminal); `--generate-cli-skeleton` makes no call;
  `--version` prints `afs-cli/1.2.3 contract/1.4 go1.27.4 linux/amd64` and makes no call.
- **Configuration (CLI-3).** One case per precedence step: `--endpoint-url` beats `AFS_ENDPOINT_URL`
  beats the profile's `endpoint_url` beats the default; `AFS_PROFILE` selects the profile when
  `--profile` is absent; `team_id` fills `--team-id` when omitted. `afs configure` writes `~/.afs/config`
  0600 without changing other profiles. `afs configure sso` lists providers from `sso list-providers` and
  writes `sso_provider_id`, `sso_issuer`, `sso_client_id`.
- **Credentials (CLI-5).** `AFS_API_KEY` set and a cached SSO token present → the key is sent, even with
  `--profile`; `--debug` logs `credential source: env AFS_API_KEY` and never the key. `~/.afs/credentials`
  at 0644 → exit 255 with `chmod 600 ~/.afs/credentials` and no call.
- **Pagination (CLI-7).** Against a 3-page fake: default output has one `items` array with every item;
  `--max-items 10` prints 10 items and `NextToken`; `--starting-token <it>` resumes after item 10;
  `--page-size 20` sends `maxResults=20`; `--page-size 101` exits 252; `--no-paginate` returns page 1
  with `nextToken`.
- **Retries and idempotency (CLI-11).** `idempotency_in_progress` twice then success → 3 attempts with the
  same key, exit 0; three times → exit 254. `idempotency_mismatch` → one attempt, exit 254. `429` with
  `Retry-After: 2` waits at least 2 s. A `PATCH` whose connection drops after the request was sent is not
  retried (exit 255); a connect failure on it is.
- **Help (CLI-14).** Golden help per operation; a test fails on an operation with no example; `afs <service>
  help` lists exactly the operations of that service from the command table CLI-US-7 checks.
- **Distribution (CLI-15).** Each release has archives per platform and `checksums.txt` with a cosign
  signature; `cosign verify-blob` against the published identity succeeds; the macOS binary passes
  `spctl --assess` (notarized); the Windows binary has a valid Authenticode signature; `brew install
  afs` from the tap installs a binary whose `afs --version` matches the tag.

## 6. Dependencies

| Dependency | Use | Notes |
|---|---|---|
| `-restapi` `api/openapi.yaml` at a release tag | Generated client and the CLI-US-7 command table | Regenerated in its own PR after each contract change (contract §8) |
| oapi-codegen | Generates `internal/client` types and operations | Pinned in `tools.go`/`go.mod` |
| `github.com/jmespath-community/go-jmespath` | `--query` | |
| IdP device-grant clients (Okta, Entra ID) | `sso login` | One public client per provider with the device authorization grant enabled; its issuer and client id come from `ListSignInProviders` (OQ-C2) |
| AWS CLI v2 on `PATH` | `aws eks get-token` in the kubeconfig exec block | Needed only by kubectl at run time |
| Signing: cosign, Apple Developer ID with notarization, Windows Authenticode certificate | CLI-15 | Certificates held in CI secrets, never in the repository |
| Homebrew tap repository | CLI-15 | Created before the first release; linked only once it exists |
| `X-Afs-Contract-Version` response header | Compatibility NFR, `--version` comparison | Not yet in the API contract (OQ-C1) |

## 7. Delivery

| Milestone | This repository |
|---|---|
| M0 | Generated client from `openapi.yaml`; config, profiles, credentials, output renderers (incl. csv), errors, retries, pagination, `--wait`, help; `sso`, `configure` |
| M1 | `teams` (core, suspend and resume, models, `list-revisions`, `describe-revision`), `members`, `audit list-audit-events`, `admin get-settings`, `put-setting`, identity providers |
| M2 | `agents` (core, snapshots, `list-storage-classes`, config values, external systems), `secrets`, `teams` config values and external systems, `models list-models`, `admin` catalog systems and storage classes |
| M3 | `budgets`, `reports` (on demand), `teams` notification channels, `admin` models and prices |
| M4 | `api-keys`, agent role links and identities, `teams list-allowed-accounts`, `restore-revision`, tags and log destinations, report schedules, `admin` service tags, price proposals, cost reconciliations, org accounts, exports and imports, `audit list-service-audit-events`; signed releases and Homebrew tap |
| M5 | `pools`, ephemeral `create-team` options, `teams get-kubeconfig`, `admin list-cost-tag-status` |
| M6 | `work-orders` (incl. workload token credential), `swarm-templates`, `admin` service-level keys and security operators; feature complete |

## 8. Open questions

| # | Question | Owner | Proposal |
|---|---|---|---|
| OQ-C1 | The API contract has no way to learn the server's contract version. Add response header `X-Afs-Contract-Version: <major.minor>` to contract §1 | root / restapi | Add the header; until then the compatibility warning is off |
| OQ-C2 | Does `ListSignInProviders` return the device-grant client id, or is `clientId` the console's PKCE client? | root / restapi | Add `cliClientId` per provider if the clients differ |
| OQ-C3 | `ListSwarmTemplates … DeleteSwarmTemplate` names `:enable` / `:disable` paths but no operations or CLI names | root | Contract names `EnableSwarmTemplate` / `DisableSwarmTemplate` (`enable-swarm-template` / `disable-swarm-template`) or drops the paths for `update-swarm-template --enabled`; the CLI follows |
| OQ-C4 | `ListTeamDirectory` (D38) is not yet in the contract | root | The CLI exposes it under the contract's CLI name once added |
| OQ-C5 | Shell completion (`afs completion bash\|zsh\|fish\|powershell`) in or out of v1 | J | In, M4, generated from the command table; a local command with no operation |
| OQ-C6 | Token storage: 0600 files (current) or OS keychains (macOS Keychain, Windows Credential Manager, Secret Service) | J | Files in v1; keychain as an opt-in profile key later |
| OQ-9 | Which IAM role the IdP federates to for viewer kubectl (affects the exec block of CLI-16) | J | From root |

## 9. Definition of done (every change)

- [ ] Command follows naming and option rules; help with an example
- [ ] `--output json` equals the API's fields; golden tests for table and text
- [ ] No secret value, token or key can reach output, logs or history (except `secretKey` once on create)
- [ ] AWS-style errors and exit codes; `AccessDenied` test
- [ ] `make gate` green; stories' acceptance criteria covered
