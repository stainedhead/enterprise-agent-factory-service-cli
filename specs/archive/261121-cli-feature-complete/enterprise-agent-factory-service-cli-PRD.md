# Enterprise Agent Factory Service — CLI PRD (`-cli`)

| | |
|---|---|
| **Status** | Draft for review, feature-complete scope (2026-11-21); revised after `review-prd`. **Extended 2026-11-22 for the Paperclip team collaboration change** ([change PRD](https://github.com/stainedhead/enterprise-agent-factory-service/blob/main/paperclip-team-collaboration-PRD.md), contract v1.3.0): requirements CLI-19 to CLI-27, stories CLI-US-24 to CLI-US-35, milestones M7 to M10, decisions PC-CLI-D1 to D11, open questions PC-CLI-OQ1 to OQ5; revised after a second `review-prd` pass (section 10). Not started |
| **Author** | J (Enterprise Architecture), drafted with Claude |
| **Parent** | [Root PRD](https://github.com/stainedhead/enterprise-agent-factory-service/blob/main/enterprise-agent-factory-service-initial-PRD.md) (stories US-1..US-26; IF-2; decisions D23–D57); solution design §11. Change PRD: [Paperclip team collaboration](https://github.com/stainedhead/enterprise-agent-factory-service/blob/main/paperclip-team-collaboration-PRD.md) (stories US-27..US-33, COL-*, SEC-7, decisions PC-D1..PC-D19; its §7.6 assigns CLI-US-24 and CLI-US-25 to this repository) |
| **Consumes** | [API contract](https://github.com/stainedhead/enterprise-agent-factory-service/blob/main/docs/contracts/api-contract.md) (v1.3.0 for the Paperclip change: §2.1, §3, §4.4, §4.4b, §4.4c, §4.6, §7) via `-restapi/api/openapi.yaml` at a release tag; for the one local command `invoke-api`, also the public edge of §4.4c |
| **Provides** | The `afs` command and its machine-readable `--output json`, which pipelines rely on |
| **Precedence** | The API contract is the single source for operation names, CLI names, permissions and enums. Where solution design §11 differs (services `identities`, `usage`, `logs`, `notifications`; `--secret-string`; `admin export-config`), the contract wins and §11 is corrected in the root |

## 1. Purpose

`afs` lets people and automation do everything the console does, with the habits of the AWS CLI:
`afs <service> <operation> [--options]`, named profiles, SSO, `--output json|table|text|yaml|csv`,
`--query`, pagination, `--cli-input-json` and `--generate-cli-skeleton`. One operation is one API call
plus pagination; the only other sanctioned extra calls are the read before an unguarded update (CLI-10)
and status polling under `--wait` (CLI-12). The CLI never decides access; the API does.

### 1.1 Change: Paperclip team collaboration

Every team may now run one collaboration agent (Paperclip), connect its other agents to it, and publish its web
UI and REST API through a governed public endpoint (change PRD §1). Without CLI support a team-admin could not
script any of it, and an automation author had no supported way to call the team's Paperclip API with an AFS
credential beside a Paperclip key. This repository delivers: the `collaboration` service (six operations),
`admin list-collaboration-endpoints`, the flag and column changes on `agents`, `secrets`, `teams`, `members`,
`api-keys`, `admin` and `reports`, and one local command, `afs collaboration invoke-api`. The measure: every
operation of API contract §4.4b has a command with help, example, columns, goldens and an `AccessDenied` test, and a
pipeline can create a Paperclip issue with one command and a team key, with no credential visible in any output
(CLI-US-24 to CLI-US-34). The CLI stays a thin client: access is decided by the API and by the edge, never here.

## 2. Users

| User | Needs |
|---|---|
| Service admins | Onboard teams, set guardrails, approve prices, manage catalogs and settings from scripts |
| Team-admins | Create agents, secrets, configuration and tags; rotate; pause and restart |
| Viewers, FinOps, SRE | Read status, usage and reports; export CSV; pipe JSON; read-only kubectl |
| CI pipelines | Non-interactive auth with a team API key; stable JSON; exit codes |
| Security operators and security tooling | Create and track CVE work orders, by person (SSO) or with a service-level API key |
| Agents and pool workers | Report work-order progress with their projected ServiceAccount token (D29) |
| Team collaborators and automation authors (change PRD) | Hold the `team_collaborator` role or a team API key with it; call the team's Paperclip REST API from a pipeline with `invoke-api`; cannot change the team's AFS configuration |

## 3. Requirements owned

Scope is feature complete (root D18): every operation in the API contract gets a command. Each
command's priority is its operation's priority in the contract.

| Id | Requirement | Product reqs | Pri |
|---|---|---|---|
| CLI-1 | A command for every operation in the API contract §4, named exactly by its `CLI` column and mapping to its `operationId`; priority per operation from the contract. Commands take ids (`--agent-id agt_…`), never names, because names are unique only within a team (D25); the CLI does not resolve names | IF-2, G6 | P0 |
| CLI-2 | Global options on every command: `--profile`, `--endpoint-url`, `--output`, `--query`, `--no-paginate`, `--max-items`, `--page-size`, `--starting-token`, `--cli-input-json`, `--generate-cli-skeleton`, `--debug`, `--no-cli-pager`, `--color on\|off\|auto`, `--cli-connect-timeout`, `--cli-read-timeout`, `--version` | IF-2 | P0 |
| CLI-3 | Local commands that call no API operation: `afs configure`, `afs sso login`, `afs sso logout` and `afs completion bash\|zsh\|fish\|powershell` (generated from the command table). `afs configure` and `afs configure sso` write `~/.afs/config` (0600). Profile keys: `endpoint_url`, `output`, `sso_provider_id`, `sso_issuer`, `sso_client_id`, `team_id` (default for `--team-id`). Precedence: option → environment (`AFS_PROFILE`, `AFS_ENDPOINT_URL`, `AFS_DEFAULT_OUTPUT`, `AFS_API_KEY`) → profile → defaults (`output json`) | IF-2 | P0 |
| CLI-4 | `afs sso login` (OIDC device grant against the profile's provider, chosen from `sso list-providers`), `sso logout`, `sso get-caller-identity`; token cache `~/.afs/sso/cache` (0600). A session lasts as long as the IdP's refresh token; the CLI refreshes silently when one exists. A token known to be expired locally → exit 255 without a call and the hint `Run afs sso login --profile <name>`; a server `401 Unauthenticated` → AWS-style error, exit 254, same hint (D41) | ACC-1 | P0 |
| CLI-5 | Credential resolution, first match wins: `AFS_API_KEY` (wins even over `--profile`'s SSO token; `--debug` names the source) → the profile's `api_key` in `~/.afs/credentials` → the profile's cached SSO token → the workload token file `AFS_API_TOKEN_FILE` (default `/var/run/secrets/afs/api/token` when it exists; re-read per request, D29). A credentials file or token cache file with a mode wider than 0600 is refused with exit 255 and a `chmod 600` hint | ACC-4 | P1 |
| CLI-6 | Output renderers: `json` exactly the API's fields, `yaml`, `table` and `text` for people, and `csv` (D39). `--query` (JMESPath) applies before rendering. `csv` asks the API for `Accept: text/csv` and writes the body unchanged; it is valid only on operations the contract marks *CSV* (else exit 252); `--query` with `csv` exits 252; auto-pagination reads each next page's token from the response header `X-Afs-Next-Token` (a CSV body has no `nextToken` field; the header is absent on the last page) and concatenates the pages keeping one header row | IF-2, RPT-1..9 | P0 |
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
| CLI-18 | `work-orders` service, including `--repos` shorthand `host/owner/repo[=<owners>]`. With `=`, owners are a team id, an email or both comma-separated (`=team_…,ops@example.com`): contains `@` → `owningContact`, starts `team_` → `owningTeamId`, anything else (such as a slug) → exit 252 before any call. Without `=` (`host/owner/repo` alone) the repository is sent with no owner and `afs-api` resolves it from the teams' `AFS_OWNED_REPOS` configuration (D38); an explicit owner always wins. An ambiguous or unowned repository is the API's `ValidationError` (exit 254) with `details.ambiguousRepos` or `details.unownedRepos` printed. No slug lookup; `--output csv` on `list-work-orders` | WO-1..6 | P2 |
| CLI-19 | The `collaboration` service: `describe-collaboration` (`GetTeamCollaboration`), `put-collaboration` (`PutTeamCollaboration`), `list-collaboration-links` (`ListCollaborationLinks`), `connect-agent` (`ConnectAgentToCollaboration`), `disconnect-agent` (`DisconnectAgentFromCollaboration`), `rotate-collaboration-keys` (`RotateCollaborationKeys`), and `admin list-collaboration-endpoints` (`ListCollaborationEndpoints`), named exactly by the contract's `CLI` column (CLI-1). Priorities per contract v1.3.1 (PC-D28): P0 in C1 describe, put, list-links, connect, disconnect; P1 rotate; P0 in C2 endpoints. The two connect actions and the rotate action send an `Idempotency-Key` (CLI-11) and have no `--wait` | COL-34, COL-10, COL-7 | P0 |
| CLI-20 | `put-collaboration` takes `--web-ui-enabled \| --no-web-ui-enabled` and `--rest-api-enabled \| --no-rest-api-enabled` (both required unless `--cli-input-json` supplies them: the PUT replaces both) and `--registration-mode manual\|managed`. It follows CLI-10 (If-Match from a read of `GetTeamCollaboration`; a `404` on that read sends no `If-Match`, PC-CLI-D2) and `--wait` (CLI-12) polls `GetTeamCollaboration` until `status.phase` is `Ready` at or above the PUT's `generation`. Every `details.fields[].reason` and every `details.reason` of §4.4b prints with a one-line hint from a small table of hints (PC-CLI-D6); the hints add words, never decisions | COL-20, COL-34 | P0 |
| CLI-21 | `agents create-agent` and `update-agent`: `--role worker\|collaboration` (create only), `--collaboration-connect \| --no-collaboration-connect`, `--paperclip-agent-id`; `--role collaboration` with any of the other three exits 252; `list-agents --role`; `list-agent-templates --capability`; `admin put-agent-template --capabilities`. Agent tables gain `Role` and `Paperclip` columns; templates a `Capabilities` column | COL-1, COL-2, COL-10, AGT-11 | P0 |
| CLI-22 | `secrets list-secrets --managed-by user\|afs\|paperclip`; a `ManagedBy` column; the API's `409 managed_secret`, `master_key_not_rotatable` and `reserved_name` errors print with hints. The CLI never reads `managedBy` to decide anything and never prints a secret value (CLI-9 unchanged) | SEC-7, COL-7 | P0 |
| CLI-23 | `admin list-collaboration-endpoints --team-id --state`, table columns, `--output csv` (the operation is *CSV*, CLI-6) | COL-23, COL-36 | P0 (C2) |
| CLI-24 | Local command `afs collaboration invoke-api` (CLI-US-25): one request to the team's edge with the profile's AFS credential in `X-Afs-Authorization` and the Paperclip key (file or stdin, never an argument) in `Authorization`; origin from `GetTeamCollaboration.restApiUrl` or `--collab-url`, validated; no redirects, no cookies; edge errors map to AWS-style messages and exit 254; Paperclip's own errors exit 254 with the body on stdout; workload tokens refused locally | COL-19, COL-27 | P1 |
| CLI-25 | Redaction for the edge: `X-Afs-Authorization` joins the masked headers everywhere; `invoke-api --debug` logs no body and no query string; the AFS token, API key and Paperclip key never reach stdout, stderr, logs, the process arguments or shell history. `SecretDetected` errors print pattern types only | SEC-1, COL-26 | P0 |
| CLI-26 | The generated surface that the contract changes without new operations: `--may-expose-collaboration \| --no-may-expose-collaboration` on `teams create-team` and `update-team`; `--collab-client-id` and `--collab-client-secret-arn` (an ARN, validated as one, never a value) on `admin create-identity-provider` and `update-identity-provider`; role `team_collaborator` on `members` and `api-keys`; report kind `collaboration` with `--group-by surface\|reason\|person`; the seven `collab.*` settings through `admin put-setting` (no client-side keys or ranges) | ACC-2, COL-17, COL-23, COL-36 | P0 |
| CLI-27 | Contract compatibility: the pinned snapshot is re-pinned to the `-restapi` release that carries contract v1.3.0, then types, descriptors, the command table, reference docs and goldens are regenerated (never hand-edited); the reader and wait tables and the milestone table are extended and checked against the snapshot (ADR-008); `afs --version` prints `contract/1.3`. The profile gains one optional key, `collab_base_domain` (PC-CLI-D5) | API contract §8, §9 | P0 |

The rows CLI-19 to CLI-27 were added for the Paperclip change (contract v1.3.0). Product requirements `COL-*`
and `SEC-7` are in the [change PRD](https://github.com/stainedhead/enterprise-agent-factory-service/blob/main/paperclip-team-collaboration-PRD.md) §6; this repository owns only the client side of them.

### 3.1 Services summary

**The API contract's §4 `CLI` column is the authoritative command list.** This table summarizes it per
service so the shape of each service is visible here; when the two differ, the contract wins and this
table is corrected. CLI-US-7 tests the generated command table against `openapi.yaml`, not this table.

| Service | Operations |
|---|---|
| `teams` | `list-teams`, `create-team`, `describe-team`, `update-team`, `delete-team`, `suspend-team`, `resume-team`, `get-team-status`, `get-kubeconfig`, `list-team-directory`; models `list-team-models`, `put-team-models`; tags `get-team-tags`, `put-team-tags`; `list-allowed-accounts`; config `list-config-values`, `put-config-value`, `delete-config-value`; external systems `list-external-systems`, `set-external-system-state`, `declare-external-system`, `update-external-system`, `delete-external-system`; log destinations `list-log-destinations`, `create-log-destination`, `describe-log-destination`, `update-log-destination`, `delete-log-destination`; notification channels `list-notification-channels`, `create-notification-channel`, `describe-notification-channel`, `update-notification-channel`, `delete-notification-channel`, `test-notification-channel`; revisions `list-revisions`, `describe-revision`, `restore-revision` |
| `members`, `api-keys` | `list-members`, `add-member`, `update-member`, `remove-member`; `list-api-keys`, `create-api-key`, `delete-api-key` |
| `agents` | `list-agents`, `create-agent`, `describe-agent`, `update-agent`, `delete-agent`, `pause-agent`, `resume-agent`, `restart-agent`, `get-agent-status`, `list-snapshots`, `create-snapshot`, `list-storage-classes`, `list-agent-templates`; role links `list-role-links`, `create-role-link`, `describe-role-link`, `delete-role-link`, `verify-role-link`; identities `list-identities`, `create-identity`, `delete-identity`; config `list-config-values`, `put-config-value`, `delete-config-value`; external systems `list-external-systems`, `set-external-system-state` |
| `secrets` | `list-secrets`, `create-secret`, `describe-secret`, `update-secret`, `put-secret-value`, `rotate-secret`, `delete-secret` |
| `budgets` | `get-budget`, `put-budget`, `get-usage`, `list-costs` |
| `reports` | `get-team-report`, `get-agent-report`, `get-service-report`; schedules `list-report-schedules`, `create-report-schedule`, `update-report-schedule`, `delete-report-schedule` |
| `models`, `audit` | `list-models`; `list-audit-events`, `list-service-audit-events` |
| `admin` | settings `get-settings`, `put-setting`; identity providers `list-identity-providers`, `create-identity-provider`, `describe-identity-provider`, `update-identity-provider`, `delete-identity-provider`; service notification channels `list-notification-channels`, `create-notification-channel`, `update-notification-channel`, `delete-notification-channel`; agent templates `list-agent-templates`, `put-agent-template`, `delete-agent-template`; service tags `list-service-tags`, `put-service-tag`, `delete-service-tag`; models `list-models`, `create-model`, `update-model`; prices `list-model-prices`, `create-model-price`, `list-price-proposals`, `approve-price-proposal`, `reject-price-proposal`; `list-cost-reconciliations`, `list-cost-tag-status`; catalog `list-external-systems`, `create-external-system`, `describe-external-system`, `update-external-system`, `delete-external-system`; storage classes `list-storage-classes`, `put-storage-class`, `delete-storage-class`; `list-org-accounts`; service-level keys `list-api-keys`, `create-api-key`, `delete-api-key`; security operators `list-security-operators`, `add-security-operator`, `remove-security-operator`; `create-export`, `describe-export`, `list-exports`, `create-import` |
| `pools` | stages `list-stages`, `create-stage`, `update-stage`, `delete-stage`; `list-pools`, `create-pool`, `describe-pool`, `update-pool`, `delete-pool`, `scale-pool`, `pause-pool`, `resume-pool`, `get-pool-status`; `list-handoffs`, `put-handoffs`; `list-secret-bindings`, `put-secret-bindings`; config `list-config-values`, `put-config-value`, `delete-config-value`; `list-external-systems`, `put-external-system`; role links `list-role-links`, `create-role-link`, `delete-role-link`, `verify-role-link` |
| `swarm-templates` | `list-swarm-templates`, `create-swarm-template`, `describe-swarm-template`, `update-swarm-template`, `delete-swarm-template`, `enable-swarm-template`, `disable-swarm-template` |
| `work-orders` | `list-work-orders`, `create-work-order`, `describe-work-order`, `update-work-order`, `cancel-work-order`, `list-repos`, `update-repo`, `list-team-work-orders` |
| `collaboration` (v1.3.0) | `describe-collaboration`, `put-collaboration`, `list-collaboration-links`, `connect-agent`, `disconnect-agent`, `rotate-collaboration-keys`; local: `invoke-api` (not an operation). `admin` also gains `list-collaboration-endpoints`; `agents` gains `--role`, `--collaboration-connect`, `--paperclip-agent-id`, `--capability` on templates; `secrets` gains `--managed-by` |
| `sso`, `configure`, `completion` | `login`, `logout`, `list-providers`, `get-caller-identity`; `configure`, `configure sso`, `completion` (`login`, `logout`, `configure` and `completion` are local only, no API operation) |

### 3.2 Non-goals

- Client-side workflows that combine several mutations (for example "create team and its agents").
- Name or slug resolution (`--agent-name`, `=team-slug`); commands take ids (D25, D38).
- Any command that reads a secret value: no permission reads one.
- Client-side access decisions: the CLI never skips a call because it predicts a `403`.
- An interactive TUI, plugins or aliases; auto-update (releases go through package managers).
- Local state beyond `~/.afs/config`, `~/.afs/credentials`, the SSO cache and the `GetMe` cache.
- For the Paperclip change: no command that opens a browser or signs in to the web UI; no WebSocket or
  server-sent-event client in `invoke-api`; no commands for Paperclip's own objects (issues, agents, keys: `invoke-api`
  is the escape hatch, not a wrapper); no creation, listing or revocation of Paperclip keys by the CLI (they are secrets
  managed by AFS or by the team-admin, SEC-7, COL-12); no client-side secret scanning (the edge scans, COL-24); no
  client-side copy of the reserved secret names, the collaboration settings' ranges or the permission matrix; no
  `collaboration create-collaboration` (the collaboration agent is created by `create-agent --role collaboration`);
  no automatic connection of existing agents (use `connect-agent`, COL-10).

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
- Given a CSV response whose `X-Afs-Next-Token` header is set, the CLI sends the next request with that
  token and stops when the header is absent [fake server with 3 pages].
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
  operation (other than `configure`, `configure sso`, `sso login`, `sso logout`, `completion` and `help`), and on any
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

**CLI-US-9 · Security tooling files work orders.** *(US-16; WO-1; D38)* As a security engineer, I want my
scanner to create a work order with a service-level key.
- With `AFS_API_KEY` set to a security operator key, `afs work-orders create-work-order --cve-id
  CVE-2026-41873 --severity critical --sla-due-at 2026-12-01T00:00:00Z --template-id swt_… --repos
  github.com/payments/ledger-core=team_01JD… github.com/vendor/edi=edi@example.com
  github.com/payments/ops=team_01JE…,ops@example.com` sends one `CreateWorkOrder` with each repository's
  `owningTeamId` and `owningContact` set as classified; a repeat with the same `--idempotency-key`
  replays the stored response.
- `--repos github.com/payments/ledger-core=payments-factory` (a slug) exits 252 before any call, saying
  to use a team id or an email.
- Given `--repos github.com/payments/ledger-core` with no `=`, then the request carries that repository
  with no `owningTeamId` and no `owningContact`, and the response's `ownerSource: team_config` is printed
  (`--output json` unchanged); an explicit `=team_…` in the same command sends the owner and wins.
- Given the API answers `ValidationError` with `details.ambiguousRepos` or `details.unownedRepos`, then the
  AWS-style error names the repositories and teams, says to give `=team_…` or `=email`, and exits 254.
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
  an account outside the team's allowed accounts gets `ValidationError` from the API (exit 254).

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

**CLI-US-18 · Export and import configuration.** *(US-13; CFG-5; D49)* As a service admin, I want to move an
environment by script.
- `afs admin create-export` prints the export with its id; `admin describe-export --export-id exp_…`
  prints `sha256`, `sizeBytes` and, once complete, `downloadUrl` (valid 7 days); `admin list-exports`
  lists exports; the bundle holds no secret value.
- `afs admin create-import --bundle file://bundle.json --account-map file://map.json` sends
  `?dryRun=true` unless `--no-dry-run` is given; the bundle file is the request body; with `--no-dry-run`
  it sends `?dryRun=false` and applies.
- A bundle file over 100 MB exits 252 before any call.

**CLI-US-19 · Test a notification channel.** *(US-19; OPS-3)* As a team-admin, I want to know a channel
delivers before I rely on it.
- `afs teams test-notification-channel --team-id team_… --channel-id nch_…` prints the API's result; when
  delivery fails, it prints the failure (channel status and error) and exits 254 with the AWS-style message
  of the `TestNotificationChannel` operation.
- `afs teams describe-notification-channel` and `describe-log-destination` print the stored configuration
  and last delivery status, never a credential value.

**CLI-US-20 · Restore a revision.** *(US-20; CFG-6)* As a team-admin, I want to see what a restore would do
and then apply it.
- `afs teams restore-revision --team-id team_… --revision-id … --dry-run` sends `dryRun: true`, prints the
  plan and changes nothing; without `--dry-run` it applies.
- A viewer running it gets `AccessDenied` (exit 254).

**CLI-US-21 · Settings are validated.** *(US-21; §6.13)* As a service admin, I want a bad setting refused.
- `afs admin put-setting --key billing.variance_alert_pct --value 150` prints
  `An error occurred (ValidationError) when calling the PutSetting operation: <message naming the range>`
  and exits 254.
- A team-admin running `put-setting` gets `AccessDenied` (exit 254).

**CLI-US-22 · Disable a swarm template.** *(US-22; WO-2)* As a service admin, I want to retire a template.
- `afs swarm-templates disable-swarm-template --swarm-template-id swt_…` sends `DisableSwarmTemplate`;
  `enable-swarm-template` reverses it; a security operator gets `AccessDenied` (exit 254).
- `afs work-orders create-work-order … --template-id swt_…` naming the disabled template prints the API's
  `ValidationError` and exits 254.

**CLI-US-23 · Suspend and delete a team.** *(US-23; TEAM-9)* As a team-admin, I want suspend to be
predictable and delete to be a service admin's act.
- `afs teams suspend-team --team-id team_…` succeeds; a following write such as `agents update-agent` prints
  `Conflict` (the team is suspended) and exits 254; `resume-team` lifts it.
- `afs teams delete-team` by a team-admin prints `AccessDenied` and exits 254.


### 4.1 Paperclip team collaboration (contract v1.3.0)

Added by the change PRD [paperclip-team-collaboration-PRD.md](https://github.com/stainedhead/enterprise-agent-factory-service/blob/main/paperclip-team-collaboration-PRD.md)
§7.6 (CLI-US-24, CLI-US-25), split finer here (CLI-US-26 to CLI-US-34). Names, flags and shapes follow API
contract v1.3.0 §4.4b and §4.4c; the contract wins on any difference. Each criterion is numbered
`CLI-US-n.k` so the traceability test can find the test that names it, and carries its `Proven by` test.
Test names are Go test functions in the package named; fakes are `httptest` servers through `clitest`
(the pinned snapshot's own examples are used by the generic tests, CLI-US-34).

**CLI-US-24 · See and change a team's collaboration settings.** *(refines US-27, US-30, US-33; COL-20, COL-34; CLI-19, CLI-20)*
As a team-admin, I want to read the state of my team's Paperclip endpoint and switch its web UI and REST
API from a script, so that I publish only what we use.
- 24.1 Given `afs collaboration describe-collaboration --team-id team_…`, then one `GetTeamCollaboration`
  (`GET /v1/teams/{teamId}/collaboration`) is sent and `--output json` prints the response unchanged. With
  `--output table` the columns are `Hostname`, `WebUi` (`webUiEnabled`), `RestApi` (`restApiEnabled`),
  `EffectiveWebUi` (`effective.webUi`), `EffectiveRestApi`, `Mode` (`registrationMode`), `Phase`
  (`status.phase`) and `Version`; a `null` hostname prints `-`. *Proven by*
  `collaboration.TestDescribeCollaborationJSONUnchanged`, `collaboration.TestDescribeCollaborationTableGolden`
  (`testdata/describe-collaboration.table.golden`, `…text.golden`).
- 24.2 Given a response whose `effective.reasons` is non-empty, then with `--output table` or `text` one
  stderr line per surface says `Not effective: <web UI|REST API> (<reason>[, <reason>])` using the contract's
  reason values (`team_not_allowed`, `service_endpoints_off`, `base_domain_unset`,
  `provider_without_collab_client`), and the exit code stays 0; with `--output json` nothing extra is
  printed. *Proven by* `collaboration.TestDescribeCollaborationNotEffectiveNotice`.
- 24.3 Given `put-collaboration --team-id team_… --web-ui-enabled --no-rest-api-enabled --registration-mode
  managed`, then one `PUT /v1/teams/{teamId}/collaboration` carries `{"webUiEnabled":true,
  "restApiEnabled":false,"registrationMode":"managed"}`; with `--registration-mode` omitted the member is
  omitted (the API keeps the current value). `--registration-mode` accepts `manual` or `managed`, anything
  else exits 252 before any call. *Proven by* `collaboration.TestPutCollaborationBody`,
  `collaboration.TestPutCollaborationRegistrationModeEnum252`.
- 24.4 Given `put-collaboration` without `--web-ui-enabled | --no-web-ui-enabled` or without
  `--rest-api-enabled | --no-rest-api-enabled` (the PUT replaces both), then it exits 252 naming the missing
  option, before any call; `--cli-input-json file://collab.json` supplies both. *Proven by*
  `collaboration.TestPutCollaborationBothSurfacesRequired252`, `collaboration.TestPutCollaborationInputJSON`.
- 24.5 Safe updates (CLI-10). With `--if-match 4` the request carries `If-Match: "4"` and no read is made;
  without it the CLI reads `GetTeamCollaboration`, sends its `version`, and prints
  `Applied to version N; use --if-match to guard against concurrent changes.` to stderr. A `412` prints
  `The resource changed (version 5). Re-run with --if-match 5 to apply your change to the latest version.`
  and exits 254. If the version read answers `404`, no `If-Match` is sent and the API's own answer
  (`409 Conflict`, `details.reason = collaboration_agent_missing`) is what the user sees (PC-CLI-D2).
  *Proven by* `collaboration.TestPutCollaborationIfMatchGuarded`, `collaboration.TestPutCollaborationReadsVersion`,
  `collaboration.TestPutCollaboration412Message`, `collaboration.TestPutCollaborationMissingAgent409`.
- 24.6 `--wait` polls `GetTeamCollaboration` every 5 s (or `Retry-After`) until `status.phase` is `Ready` and
  `status.observedGeneration` is at least the `generation` of the PUT response (exit 0, prints the final
  resource), or `Failed` or `Degraded` (exit 254, conditions including `EndpointReady` and `PaperclipReady`
  on stderr), or `--wait-timeout` passes (default 10 min, exit 255). A stale `Ready` with a lower
  `observedGeneration` does not end the wait. *Proven by* `collaboration.TestPutCollaborationWaitReady`,
  `collaboration.TestPutCollaborationWaitIgnoresStaleReady`, `collaboration.TestPutCollaborationWaitDegraded254`,
  `collaboration.TestPutCollaborationWaitTimeout255`, `wait.TestCollaborationWaitTableMatchesSnapshot`.
- 24.7 Validation answers from the API print AWS-style and exit 254, with every `details.fields[]` entry as
  one line `field: <field> reason: <reason>` on stderr after `Request ID`: turning a surface on with
  `collab_base_domain_unset`, `collab_public_endpoints_off`, `collaboration_not_allowed_for_team`,
  `provider_without_collab_client`, and `registrationMode managed` with `collaboration_agent_not_ready`.
  Each reason also prints the hint of PC-CLI-D6 (for example `Ask a service admin to set collab.base_domain
  (afs admin put-setting).`). Switching a surface off is never refused client-side. *Proven by*
  `collaboration.TestPutCollaborationValidationFieldsPrinted`, `collaboration.TestCollaborationReasonHints`.
- 24.8 `AccessDenied`: a viewer or collaborator running `put-collaboration` gets the API's `403`, printed as
  `An error occurred (AccessDenied) when calling the PutTeamCollaboration operation: …` with
  `Required permission: collaboration:manage` and exit 254; a viewer running `describe-collaboration` is
  allowed, and a caller with no binding gets the API's `404` (exit 254). A suspended team's `409` with
  `details.reason = team_suspended` prints its hint and exits 254. *Proven by*
  `collaboration.TestPutCollaborationAccessDenied254`, `collaboration.TestDescribeCollaborationAccessDenied254`,
  `collaboration.TestPutCollaborationSuspended409`.
- 24.9 `describe-collaboration` on a team with no collaboration agent prints
  `An error occurred (NotFound) when calling the GetTeamCollaboration operation: …`, the hint `This team has
  no collaboration agent. Create one with: afs agents create-agent --role collaboration …`, and exits 254.
  *Proven by* `collaboration.TestDescribeCollaborationMissingAgentHint`.
- 24.10 Every operation of this story has help with at least one example and a golden help page
  (`afs collaboration help`, `afs collaboration put-collaboration help`); the `put-collaboration` example
  is the equivalent shown on the console's T16 switch (`afs collaboration put-collaboration --team-id … --web-ui-enabled
  --no-rest-api-enabled`). *Proven by* `collaboration.TestCommandsDoneCollaboration`,
  `help.TestCollaborationHelpGolden`.

**CLI-US-25 · Call the Paperclip API from a pipeline.** *(refines US-31; COL-19, COL-27; CLI-24, CLI-25)*
As an automation author, I want one command that calls my team's Paperclip REST API with my AFS credential
and a Paperclip key, so that a pipeline can create and update work without a hand-built `curl`.

`afs collaboration invoke-api --team-id <id> --method <verb> --path <path> [--body file://… | fileb://…]
--paperclip-key-file <file | ->` is a **local command** (API contract §4.4b): it is not an operation and has
no `operationId` of its own; it sends one request to the team's edge (§4.4c) and prints the response.
Options: `--team-id` (or the profile's `team_id`), `--method` (`GET`, `HEAD`, `POST`, `PUT`, `PATCH`,
`DELETE`, case-insensitive; default `GET`), `--path` (required), `--body`, `--content-type` (default
`application/json` when a body is present), `--paperclip-key-file` (required), `--collab-url` (optional
override of the host, PC-CLI-D1), `--include` (print the status line and response headers before the body).
- 25.1 Given the REST API is on, a team API key with the collaborator role in `AFS_API_KEY`, and
  `--paperclip-key-file ./pc.key`, when I run `afs collaboration invoke-api --team-id team_… --method POST
  --path /api/issues --body file://issue.json --paperclip-key-file ./pc.key`, then exactly one
  `GetTeamCollaboration` and one request to the origin of the resource's `restApiUrl` are made; the second
  carries `X-Afs-Authorization: Bearer <the AFS credential>`, `Authorization: Bearer <contents of pc.key>`,
  the body, `Content-Type: application/json`, `Accept: application/json`, `User-Agent: afs-cli/<version>
  (<os>/<arch>)` and **no** `Cookie`; the response body is written to stdout unchanged and the exit is 0.
  *Proven by* `collaboration.TestInvokeAPIRequestShape`, `collaboration.TestInvokeAPIResponseBodyUnchanged`.
- 25.2 Given an SSO profile, then the SSO access token is sent in `X-Afs-Authorization` (silent refresh as in
  CLI-4); given `AFS_API_KEY`, the key; the credential chain of CLI-5 is used and `--debug` names the source
  without printing it. A workload token (`AFS_API_TOKEN_FILE`) is refused with exit 255 and `The collaboration
  endpoint does not admit workload credentials; use an API key or afs sso login.` before any call to the
  edge (PC-CLI-D3). A token known to be expired locally exits 255 with the login hint and makes no call.
  *Proven by* `collaboration.TestInvokeAPICredentialChain`, `collaboration.TestInvokeAPIWorkloadTokenRefused255`,
  `collaboration.TestInvokeAPIExpiredToken255`.
- 25.3 Given `--paperclip-key-file -`, then the key is read from stdin; a file whose mode is wider than 0600
  is refused with exit 255 and `chmod 600 <file>` (as CLI-5); trailing whitespace is trimmed; an empty
  key, a key longer than 4096 bytes or one containing whitespace or control characters exits 252 before any
  call. There is **no** option that takes the key as an argument: `--paperclip-key <value>` is an unknown
  option (exit 252) and the usage names `--paperclip-key-file`. `--body -` exits 252 (stdin is for the key).
  *Proven by* `collaboration.TestInvokeAPIKeyFromStdin`, `collaboration.TestInvokeAPIKeyFileMode255`,
  `collaboration.TestInvokeAPIKeyShapeRefused252`, `collaboration.TestInvokeAPIKeyAsArgumentUnknown252`.
- 25.4 Client-side validation, all exit 252 before any call: `--method` not in the list above; `--path` not
  starting with `/`, containing a scheme or host, `..` segments, control characters or a space; `--path`
  beginning `/_afs/` (the edge's reserved paths are never forwarded and are not Paperclip's); `--body` with
  `GET` or `HEAD`; a `--body` file over 100 MB (the service setting's maximum, `collab.max_body_mb`); a
  `--body` not `file://` or `fileb://`. The edge's own limit (`collab.max_body_mb`, default 25) is the API's
  `400 ValidationError` with `details.reason = body_too_large` and exits 254. *Proven by*
  `collaboration.TestInvokeAPIValidation252`, `collaboration.TestInvokeAPIReservedPathRefused252`,
  `collaboration.TestInvokeAPIBodyTooLargeFromEdge254`.
- 25.5 Where the request goes. The origin comes from `restApiUrl` of the `GetTeamCollaboration` answer, or from
  `--collab-url` (then no `GetTeamCollaboration` is made). `restApiUrl` is `https://<hostname>/api`, the base of Paperclip's REST API, so a full `--path /api/…` equals it plus a suffix (PC-D31). It must be `https://` (an `http://` loopback host is
  accepted for tests, as for `endpoint_url`), carry no userinfo, query, fragment or port other than 443, and a
  DNS name (not an IP literal; the apex `<collab.base_domain>` serves only the OIDC callback and is not a team host); if the profile sets `collab_base_domain` the host must be `<label>.<that
  domain>` (PC-CLI-D5). A `null` `restApiUrl` (service base domain unset) exits 255 with `Public
  collaboration is not configured for this service (collab.base_domain is unset).` Only the URL's scheme and
  host are used; `--path` is appended to the origin and the URL's own path (`/api`) is not (PC-CLI-D4).
  Redirects are **not** followed: a `3xx` is printed like any response (with `--include`, its `Location`). No
  cookie jar: a `Set-Cookie` in a response is never stored or sent. *Proven by*
  `collaboration.TestInvokeAPIOriginFromRestApiURL`, `collaboration.TestInvokeAPICollabURLSkipsRead`,
  `collaboration.TestInvokeAPIOriginRefused255` (table: http, userinfo, port, IP literal, foreign domain),
  `collaboration.TestInvokeAPINullRestApiURL255`, `collaboration.TestInvokeAPINoRedirectFollowed`,
  `collaboration.TestInvokeAPINoCookieJar`.
- 25.6 Errors from the edge (the §6 envelope, API contract §4.4c table) print AWS-style with the operation name
  `InvokeApi` and exit 254, then `Request ID: <requestId>`:
  `403 AccessDenied` `surface_disabled` → `An error occurred (AccessDenied) when calling the InvokeApi
  operation: …` plus `The REST API is switched off for this team. A team-admin can turn it on: afs
  collaboration put-collaboration --team-id <id> --rest-api-enabled …`; `403 permission_required` prints
  `Required permission: <details.requiredPermission>`; `403 no_team_role`; `401 Unauthenticated` with
  `token_expired`, `token_invalid`, `key_expired`, `key_revoked` prints the login or key hint (D41);
  `400` with `ambiguous_credentials`, `malformed_credential`, `unsupported_encoding` (`host_mismatch` no longer exists, PC-D23);
  `404 NotFound` (no detail, ever; only for an unknown host or both surfaces not effective, PC-D31) prints `The team's Paperclip is not reachable through the public endpoint; run afs collaboration describe-collaboration to see why.`; a caller with a binding gets `403 surface_disabled` instead when exactly one surface is off;
  `422 SecretDetected` prints the message and `Patterns: <details.patterns joined>` (pattern types only) and
  the hint `Store the credential as an AFS secret instead (afs secrets create-secret).`; `429 Throttled`
  and `503 Unavailable` (`paused`, `team_suspended`, `starting`, `upstream_unavailable`) print the reason.
  *Proven by* `collaboration.TestInvokeAPIEdgeErrors254` (table over every status and reason above, each
  with golden stderr), `collaboration.TestInvokeAPISurfaceDisabledHint`, `collaboration.TestInvokeAPISecretDetectedPatternsOnly`.
- 25.7 A response that is not the §6 envelope (Paperclip's own `4xx` or `5xx`, or a WAF block's default body)
  exits 254 with `An error occurred (HttpError) when calling the InvokeApi operation: HTTP <status>
  <reason phrase>` on stderr; the response body is still written to stdout so a pipeline can read
  Paperclip's error; the body is never copied into the stderr message (it is Paperclip content, ADR-012).
  A `2xx` exits 0 (`204` and `HEAD` print nothing); a `3xx` exits 0. *Proven by*
  `collaboration.TestInvokeAPIPaperclipErrorBodyOnStdout254`, `collaboration.TestInvokeAPIWAFBlock254`,
  `collaboration.TestInvokeAPI204AndHead`.
- 25.8 Retries (CLI-11, narrowed): `GET` and `HEAD` retry `429` and `503` (honouring `Retry-After`) and connect
  failures, at most 3 attempts; any other method retries only connect failures, never after the request was
  sent; no `Idempotency-Key` is sent (the edge does not define one). *Proven by*
  `collaboration.TestInvokeAPIRetriesGetOn503`, `collaboration.TestInvokeAPIPostNotRetriedAfterSend`,
  `collaboration.TestInvokeAPINoIdempotencyKey`.
- 25.9 Global options. `--profile`, `--endpoint-url` (the AFS API used for the one read), `--debug`,
  `--cli-connect-timeout`, `--cli-read-timeout`, `--color`, `--no-cli-pager` and `--version` behave as usual.
  `--query`, `--output`, `--no-paginate`, `--max-items`, `--page-size`, `--starting-token`,
  `--cli-input-json`, `--generate-cli-skeleton` and `--wait` exit 252 with `invoke-api prints the response
  unchanged; pipe it to a JSON tool instead.` (PC-CLI-D7). Streaming responses (server-sent events) and
  WebSocket are not supported: the read timeout applies and a `101` is an error (exit 255). *Proven by*
  `collaboration.TestInvokeAPIGlobalOptionsRefused252`, `collaboration.TestInvokeAPIReadTimeout255`.
- 25.10 Redaction (CLI-25). With `--debug`, stderr shows one line per event with method, host, path **without
  its query string**, status, duration, `X-Request-Id`, credential source and attempt; `Authorization` and
  `X-Afs-Authorization` print `[redacted]`; request and response **bodies are never logged**; `Set-Cookie`
  and `Cookie` print `[redacted]`. A test drives a run with a recognisable AFS token, API key, Paperclip key and
  a `?token=…` query and asserts that none of those strings, nor the body, appears anywhere in stdout (other
  than the response body itself), stderr, `--debug` output or the process arguments. *Proven by*
  `collaboration.TestInvokeAPIDebugRedactsEverything`, `collaboration.TestInvokeAPISecretsNeverInArgsOrStderr`,
  `redact.TestRedactsEdgeHeaders`.
- 25.11 The command is listed by `afs collaboration help` and in the generated command reference with a
  golden help page and two examples (one `GET`, one `POST` with `--body`); it is excluded from CLI-US-7's
  operation comparison by name (as `configure`, `sso login`). *Proven by* `help.TestInvokeAPIHelpGolden`,
  `services.TestOperationTableExemptsLocalCommands`.

**CLI-US-26 · Connect and disconnect agents, and see the links.** *(refines US-28; COL-10, COL-12, COL-14; CLI-19)*
As a team-admin, I want to connect existing agents to Paperclip in bulk from a script and see their link
state, so that agents that predate the collaboration agent can pull work.
- 26.1 `afs collaboration connect-agent --team-id team_… --agent-id agt_… [--paperclip-agent-id ag_…]` sends
  one `POST …/collaboration/links/{agentId}:connect` with body `{"paperclipAgentId":"ag_…"}` (or `{}`), an
  `Idempotency-Key`, and prints the `CollaborationLink`; `--output table` columns `AgentId`, `Name`
  (`agentName`), `Connect`, `PaperclipAgentId`, `LinkState`, `LastSyncedAt`. `disconnect-agent` sends
  `…:disconnect` with no body and the same columns. *Proven by*
  `collaboration.TestConnectAgentRequest`, `collaboration.TestDisconnectAgentRequest`,
  `collaboration.TestLinkColumnsGolden`.
- 26.2 A repeat with the same `--idempotency-key` replays the stored response (`Idempotent-Replayed`) and
  prints it; the CLI-11 retry rules apply. Neither operation takes `If-Match`, has `--wait` or reads a version
  (`--wait` is an unknown option, exit 252). *Proven by* `collaboration.TestConnectAgentReplay`,
  `collaboration.TestConnectAgentNoWaitNoIfMatch`.
- 26.3 API refusals print AWS-style and exit 254: `ValidationError` for the collaboration agent itself
  (`field: agentId`), for a pool, and for `--paperclip-agent-id` while `registrationMode` is `managed`
  (`field: paperclipAgentId`, `reason: managed_registration`, with the hint `Registration is managed; the
  reconciler records the id.`); `NotFound` for an unknown agent; `409` `collaboration_agent_missing` or
  `team_suspended`. *Proven by* `collaboration.TestConnectAgentValidation254`,
  `collaboration.TestConnectAgentManagedRegistrationHint`.
- 26.4 `afs collaboration list-collaboration-links --team-id team_… [--link-state pending|linked|failed|disabled]
  [--connect | --no-connect]` sends `GET …/collaboration/links?linkState=…&connect=true|false`, paginates by
  default (CLI-7), and `--connect` omitted sends no `connect` query member. `--link-state` outside the four
  values exits 252. `lastError` is printed as returned (it never holds a key or content, contract §4.4b);
  with `--output table` the `LastError` column is truncated with the table's mark. *Proven by*
  `collaboration.TestListCollaborationLinksFilters`, `collaboration.TestListCollaborationLinksPagination`,
  `collaboration.TestListCollaborationLinksLinkStateEnum252`, `collaboration.TestLinksTableGolden`.
- 26.5 `AccessDenied` (exit 254, `Required permission: collaboration:manage`) for `connect-agent` and
  `disconnect-agent` by a collaborator or viewer; `collaboration:read` for `list-collaboration-links`.
  *Proven by* `collaboration.TestConnectAgentAccessDenied254`, `collaboration.TestDisconnectAgentAccessDenied254`,
  `collaboration.TestListCollaborationLinksAccessDenied254`.

**CLI-US-27 · Rotate Paperclip's keys.** *(refines US-27; COL-7; CLI-19)* As a team-admin, I want to rotate
Paperclip's signing secrets, knowing it signs everyone out.
- 27.1 `afs collaboration rotate-collaboration-keys --team-id team_… [--names paperclip-better-auth-secret
  paperclip-tool-action-signing-secret]` sends one `POST /v1/teams/{teamId}/collaboration:rotate-keys` with an
  `Idempotency-Key` and body `{"names":[…]}` (no body when `--names` is omitted; the API rotates both). It
  prints the response; `--output table` lists `Name`, `SecretId`, `VersionId` from `rotated`. *Proven by*
  `collaboration.TestRotateCollaborationKeysBody`, `collaboration.TestRotateCollaborationKeysTable`.
- 27.2 After a success the CLI prints to stderr `Everyone is signed out of Paperclip.` when
  `signsOutEveryone` is true and `Agents will roll to pick up the new keys.` when `agentRolloutPending` is
  true. It never prompts for confirmation (PC-CLI-D8). No secret value exists in the response or any output.
  *Proven by* `collaboration.TestRotateCollaborationKeysNotices`, `collaboration.TestRotateCollaborationKeysNoValue`.
- 27.3 Naming `paperclip-secrets-master-key` in `--names` is **not** refused by the CLI (the API decides); the
  API's `400 ValidationError` (`field: names`, `reason: master_key_not_rotatable`) prints with the hint `The
  master key cannot be rotated: Paperclip's stored identities would become undecryptable.` and exits 254.
  `afs secrets rotate-secret` on the master key's secret id gets the same error (CLI-US-29). *Proven by*
  `collaboration.TestRotateCollaborationKeysMasterKey254`.
- 27.4 `AccessDenied` for collaborator, viewer and security operator (`Required permission:
  collaboration:manage`, exit 254). *Proven by* `collaboration.TestRotateCollaborationKeysAccessDenied254`.

**CLI-US-28 · Create the collaboration agent and connect workers.** *(refines US-27, US-28; COL-1, COL-2, COL-10; CLI-21)*
As a team-admin, I want `create-agent` to create the team's Paperclip and to connect workers, so that one
`--cli-input-json` file stands up a team.
- 28.1 `afs agents create-agent --team-id team_… --name paperclip --role collaboration --template paperclip
  --image …@sha256:… --model-id … [--volume …] --wait` sends `"role":"collaboration"` and returns 0 when the
  agent's phase is `Ready` (with the default 10 min `--wait-timeout`, since the product target is 5 min at
  p95). `--role` accepts `worker` (the default; omitted when not given) or `collaboration`, anything else exits
  252. *Proven by* `agents.TestCreateAgentRoleCollaborationBody`, `agents.TestCreateAgentRoleEnum252`,
  `agents.TestCreateAgentCollaborationWaitReady`.
- 28.2 `--collaboration-connect` and `--no-collaboration-connect` send `{"collaboration":{"connect":true|false}}`;
  `--paperclip-agent-id ag_…` adds `paperclipAgentId`. All three are valid on `create-agent` and `update-agent`
  for workers (`update-agent` is the same effect as `connect-agent`, contract §4.4). With none of them the body
  has no `collaboration` member, so the API default applies (`true` when the team has a collaboration agent). Given
  `--role collaboration` together with any of the three options, the CLI exits 252 before any call (the API
  would refuse `collaboration` on that role). *Proven by* `agents.TestCreateAgentCollaborationConnectBody`,
  `agents.TestUpdateAgentCollaborationConnectBody`, `agents.TestCreateAgentCollaborationRoleWithConnect252`,
  `agents.TestCreateAgentNoCollaborationMemberByDefault`.
- 28.3 `update-agent` has no `--role` option: it exits 252 as an unknown option (`--replicas` precedent,
  CLI-US-12); a change of role through `--cli-input-json` is the API's `ValidationError` (exit 254). *Proven by*
  `agents.TestUpdateAgentRoleUnknown252`, `agents.TestUpdateAgentRoleChangeRefused254`.
- 28.4 `list-agents --role worker|collaboration` sends `?role=`; any other value exits 252. The agent table gains
  `Role` and `Paperclip` columns after `Phase` (`role`; `collaboration.linkState`, `-` on the collaboration agent
  and on agents with no `collaboration` object); the table stays within 80 columns through the existing
  truncation (`output.TestTableFitsIn80Columns` stays green). *Proven by* `agents.TestListAgentsRoleFilter`, `agents.TestListAgentsRoleEnum252`,
  `agents.TestAgentTableRoleAndPaperclipColumnsGolden`.
- 28.5 `describe-agent --output json` prints `role` and `collaboration` `{connect, paperclipAgentId, linkState,
  lastSyncedAt, lastError}` unchanged. *Proven by* `agents.TestDescribeAgentCollaborationJSONUnchanged`.
- 28.6 The API's refusals print AWS-style and exit 254: `409 Conflict` with `details.reason =
  collaboration_agent_exists` (hint `This team already has a collaboration agent; see afs agents list-agents --role collaboration.`), `ValidationError` naming the template (`field: template`), `iamIdentity team_shared`,
  `volume.snapshotSchedule off`, a `paperclip` template on a worker, `collaboration.connect true` with no
  collaboration agent (`field: collaboration.connect`), and `503 Unavailable` when Secrets Manager cannot be
  written (retried per CLI-11 because `create-agent` carries an `Idempotency-Key`). *Proven by*
  `agents.TestCreateAgentCollaborationErrors254` (table), `agents.TestCreateAgentCollaboration503Retried`.
- 28.7 `afs agents list-agent-templates [--capability paperclip]` and `afs admin list-agent-templates
  [--capability paperclip]` send `?capability=`; templates show `capabilities` in `--output json` and a
  `Capabilities` column in tables; `afs admin put-agent-template --capabilities paperclip` sends
  `capabilities: ["paperclip"]` (space-separated list, CLI-2 style; `--capabilities` with no value sends `[]`).
  *Proven by* `agents.TestListAgentTemplatesCapabilityFilter`, `admin.TestPutAgentTemplateCapabilities`,
  `admin.TestAgentTemplatesCapabilitiesColumnGolden`.
- 28.8 `AccessDenied` (exit 254, `Required permission: agents:write`) for a viewer or collaborator creating a
  collaboration agent. *Proven by* `agents.TestCreateAgentCollaborationAccessDenied254`.
- 28.9 Skeleton: `afs agents create-agent --generate-cli-skeleton` includes `role` and `collaboration` with their
  contract members and makes no call. *Proven by* `agents.TestCreateAgentSkeletonHasRoleAndCollaboration`.

**CLI-US-29 · See who manages each secret.** *(refines US-27, US-32; SEC-7, COL-7; CLI-22)* As a team-admin, I want
to tell my own secrets from the ones AFS generated or Paperclip created, so that I do not try to change what I
cannot.
- 29.1 `afs secrets list-secrets --team-id team_… [--managed-by user|afs|paperclip]` sends `?managedBy=`; any other
  value exits 252. `--output json` prints `managedBy` unchanged; the table gains a `ManagedBy` column. *Proven by*
  `secrets.TestListSecretsManagedByFilter`, `secrets.TestListSecretsManagedByEnum252`,
  `secrets.TestSecretTableManagedByColumnGolden`.
- 29.2 `update-secret`, `delete-secret` and `put-secret-value` on an `afs` or `paperclip` secret print the API's
  `409 Conflict` (`details.reason = managed_secret`) with the hint `This secret is managed by <AFS|Paperclip>
  and cannot be changed here.` and exit 254. The CLI does not read `managedBy` first (the API decides), and a
  `put-secret-value` still reads the value only from a file, stdin or the prompt (CLI-9) and never echoes it
  even when the call is refused. *Proven by* `secrets.TestManagedSecretMutations409`,
  `secrets.TestManagedSecretRefusalNeverEchoesValue`.
- 29.3 `rotate-secret` on an `afs` secret sends no body (AFS generates the value, PC-D30) and prints the rotation; on
  `paperclip-secrets-master-key` it is the API's `ValidationError` (`reason: master_key_not_rotatable`), and on a
  `paperclip`-managed secret the API's `409 Conflict` (`managed_secret`); both exit 254. *Proven by* `secrets.TestRotateManagedSecret`, `secrets.TestRotateMasterKey254`.
- 29.4 `create-secret --name paperclip` or any `paperclip-…` name other than `paperclip-agent-api-key` on an
  agent-scope secret prints the API's `ValidationError` (`reason: reserved_name`) and exits 254 (the CLI keeps
  no copy of the reserved list); `create-secret --name paperclip-agent-api-key --scope agent --agent-id agt_…
  --secret-value file://…` sends the request normally (manual registration, COL-12). Binding a `paperclip`
  secret is the API's `ValidationError`. *Proven by* `secrets.TestCreateSecretReservedName254`,
  `secrets.TestCreateAgentKeySecretManualRegistration`.
- 29.5 `AccessDenied` (exit 254) for a collaborator running `create-secret` (`secrets:write`). *Proven by*
  `secrets.TestCreateSecretAccessDeniedCollaborator254`.

**CLI-US-30 · Service admin: see every team's endpoint.** *(refines US-33; COL-23, COL-36; CLI-23)* As a service
admin, I want a list of every team's Paperclip endpoint and its state from a script.
- 30.1 `afs admin list-collaboration-endpoints [--team-id team_…] [--state Ready|Degraded|Failed|Pending]` sends
  `GET /v1/admin/collaboration-endpoints?teamId=&state=` and paginates by default; `--state` outside the four values
  exits 252 (`--team-id` is the query filter here, not a path parameter, and falls back to the profile's
  `team_id` only for path parameters, so it is not defaulted). *Proven by*
  `admin.TestListCollaborationEndpointsFilters`, `admin.TestListCollaborationEndpointsStateEnum252`,
  `admin.TestListCollaborationEndpointsTeamIdNotDefaulted`.
- 30.2 `--output table` columns: `Team` (`slug`), `TeamState`, `Hostname`, `WebUi` (`effective.webUi`), `RestApi`
  (`effective.restApi`), `Phase`, `LastAccessAt`; a `null` `lastAccessAt` prints `-`. `--output csv` is valid
  (the operation is *CSV* in the contract): it sends `Accept: text/csv`, follows `X-Afs-Next-Token` over 3 pages
  and writes one header row (CLI-6). *Proven by* `admin.TestCollaborationEndpointsTableGolden`,
  `admin.TestListCollaborationEndpointsCSVThreePages`.
- 30.3 `AccessDenied` (exit 254, `Required permission: admin:settings`) for a team-admin, a collaborator and a
  security operator. *Proven by* `admin.TestListCollaborationEndpointsAccessDenied254`.

**CLI-US-31 · Service admin: guardrail, provider client and settings.** *(refines US-33, US-29; COL-17, COL-23, ACC-1; CLI-26)*
As a service admin, I want to switch collaboration per team, register the collaboration sign-in client, and set the
`collab.*` settings from scripts.
- 31.1 `afs teams create-team` and `afs teams update-team` accept `--may-expose-collaboration` and
  `--no-may-expose-collaboration` (sent as `mayExposeCollaboration`; omitted when neither is given); a team-admin
  using them gets the API's `403 AccessDenied` with `Required permission: teams:set-guardrails` (exit 254).
  `describe-team --output json` prints `mayExposeCollaboration` unchanged. `create-team` with a slug whose hostname is
  held prints the API's `409 Conflict` (`details.reason = hostname_held`) with the hint `That slug's collaboration
  hostname is held for 30 days after the deletion of a team that had one; choose another slug.` (exit 254), proven by
  `teams.TestCreateTeamHostnameHeld409`. *Proven by* `teams.TestUpdateTeamMayExposeCollaborationBody`,
  `teams.TestCreateTeamMayExposeCollaborationBody`, `teams.TestUpdateTeamMayExposeCollaborationAccessDenied254`.
- 31.2 `afs admin create-identity-provider` and `update-identity-provider` accept `--collab-client-id <id>` and
  `--collab-client-secret-arn <arn>` (an ARN of a Secrets Manager secret, never the value; the option name says
  so and the help explains it). A value that is not an `arn:aws:secretsmanager:…` string exits 252 before any call.
  `describe-identity-provider` prints both members unchanged; `--debug` does **not** redact the ARN (it is a
  reference, not a secret) but still redacts `clientSecret` and any `…Secret` member that is not an ARN
  reference. *Proven by* `admin.TestIdentityProviderCollabClientBody`, `admin.TestIdentityProviderCollabSecretArnShape252`,
  `redact.TestCollabClientSecretArnNotRedactedButSecretIs`.
- 31.3 `afs admin put-setting --key collab.base_domain --value '"collab.example.com"'` and the other six `collab.*`
  settings send the JSON value as for every setting (CLI-US-15); the CLI carries no list of keys or ranges: an
  out-of-range value (`collab.session_ttl_minutes 5`, `collab.max_body_mb 500`, `collab.content_scan "maybe"`,
  `collab.allowed_source_cidrs ["10.0.0.0/33"]`) prints the API's `ValidationError` naming the range and exits 254.
  `admin get-settings --output json` prints the settings unchanged. *Proven by*
  `admin.TestPutSettingCollabKeysPassThrough`, `admin.TestPutSettingCollabOutOfRange254`.
- 31.4 `AccessDenied` for a team-admin running `put-setting` or `create-identity-provider` is already covered by
  CLI-US-15 and CLI-US-21; this story adds the cases for the new members. *Proven by*
  `admin.TestIdentityProviderCollabClientAccessDenied254`.

**CLI-US-32 · Collaborator role and the collaboration report.** *(refines US-29, US-31; ACC-2, COL-36; CLI-26)* As a
team-admin, I want to give a person or a pipeline the collaborator role, and to read who uses Paperclip.
- 32.1 `afs members add-member`, `members update-member` and `api-keys create-api-key` accept `--role
  team_collaborator` (the API enum is `team_admin`, `team_collaborator`, `team_viewer`; for `create-api-key` also
  `security_operator`, which the API still refuses on a team key); a service-level key or a service-level binding
  with `team_collaborator` is the API's `ValidationError` (exit 254). The help lists the three roles. *Proven by*
  `members.TestAddMemberTeamCollaborator`, `apikeys.TestCreateApiKeyTeamCollaborator`,
  `members.TestMemberRoleEnumHelpGolden`.
- 32.2 `afs reports get-team-report --kind collaboration --range 7d --group-by day,person` and `afs reports
  get-service-report --kind collaboration --group-by team` are accepted: the CLI's client-side `--group-by` check
  gains `surface`, `reason` and `person` (valid only with the collaboration kind, as the contract says;
  other kinds with them are the API's `ValidationError`, the CLI no longer rejects them client-side). `--output csv`
  works as for other reports; the table columns follow the response (`Dynamic`). `create-report-schedule --report
  collaboration` is the API's `ValidationError` (the report is not schedulable, PC-D13), exit 254. *Proven by*
  `reports.TestCollaborationReportGroupByAccepted`, `reports.TestGroupByPatternMessageListsNewDimensions`,
  `reports.TestCollaborationReportScheduleRefused254`, `reports.TestCollaborationReportCSV`.
- 32.3 `AccessDenied` (exit 254) for a security operator running `get-team-report --kind collaboration` (team
  scoped routes are refused to them). *Proven by* `reports.TestCollaborationReportAccessDenied254`.

**CLI-US-33 · Edge credentials are never shown.** *(refines US-31, US-32; SEC-1, COL-26; CLI-25)* As security, I want
the new credentials treated like every other.
- 33.1 `internal/redact` masks the headers `Authorization`, `Cookie`, `Set-Cookie`, `Proxy-Authorization` and the new
  `X-Afs-Authorization`, and `--debug` output of every command (not only `invoke-api`) uses it. The value of the
  edge cookie `__Host-afs-collab` never appears in any output. *Proven by* `redact.TestRedactsEdgeHeaders`,
  `client.TestDebugRedactsXAfsAuthorization`.
- 33.2 The query member `code` (the edge's one-time handoff code) and any `token`, `key` or `secret` member are
  dropped from URLs in `--debug` lines for `invoke-api`; for the AFS API's own calls the existing rule
  (CLI-US-3, `TestNoQueryParamCarriesSecret`) stays. *Proven by* `collaboration.TestInvokeAPIDebugDropsQuery`.
- 33.3 `rotate-collaboration-keys`, `connect-agent` and every secret path never print or log a secret value; a
  test greps the stdout, stderr and `--debug` output of all seven new operations for a planted value and finds it
  nowhere. *Proven by* `collaboration.TestNoSecretValueInAnyCollaborationOutput`.

**CLI-US-34 · The command table, contract pin and documentation follow the contract.** *(extends CLI-US-7; refines
US-27 to US-33; CLI-1, CLI-14, CLI-27)* As a maintainer, I want the new operations to exist because the pinned
contract has them, with help, examples, columns, goldens and docs.
- 34.1 After the re-pin the command-table test (CLI-US-7) passes with seven more operations and fails on any of them
  without a command or with a name that differs from the contract's `CLI` column:
  `collaboration describe-collaboration`, `put-collaboration`, `list-collaboration-links`, `connect-agent`,
  `disconnect-agent`, `rotate-collaboration-keys`, and `admin list-collaboration-endpoints`. `invoke-api` is
  listed as a local command (not an operation). *Proven by* `services.TestOperationTableMatchesSnapshot`,
  `services.TestOperationTableExemptsLocalCommands`.
- 34.2 Every new operation has at least one example, table columns (or `Dynamic`), a golden help page, a golden
  table and text rendering, and an `AccessDenied` test; `CurrentMilestone` is raised only by the exit task of the
  milestone (PC-CLI-D9: M7 to M10), and `Pending()` is empty at each. *Proven by*
  `cli.TestPendingEmptyAtCurrentMilestone`, `services.TestEveryOperationHasAccessDeniedTest`,
  `help.TestEveryOperationHasExample`.
- 34.3 The `If-Match` reader table and the `--wait` table are checked against the snapshot (ADR-008): `PutTeamCollaboration`
  has a reader (`GetTeamCollaboration`, same path) and a wait entry (the resource's own `status`); the
  two new POST actions have neither. A test fails on any divergence. *Proven by*
  `safeupdate.TestReaderTableMatchesSnapshot`, `wait.TestCollaborationWaitTableMatchesSnapshot`.
- 34.4 `api/openapi.yaml` and `api/SOURCE` are re-pinned with `make contract-sync` from the `-restapi` ref that carries
  contract v1.3.0 (PC-CLI-D10), then `make generate`, `make docs`; `afs --version` prints `contract/1.3`; the
  snapshot test (`api/snapshot_test.go`) passes on the new sha256. Between the re-pin and the last milestone the
  operations of later milestones exist as generated commands but are not "done" (`milestoneOf`). *Proven by*
  `api.TestSnapshotMatchesSource`, `buildinfo.TestVersionPrintsContractMinor`.
- 34.5 Documentation: `user-docs/collaboration.md` (new: connect to Paperclip, the public endpoint, `invoke-api`
  with a team key, secrets in Paperclip, troubleshooting by edge error); `user-docs/agents-and-secrets.md`,
  `teams-and-members.md` and `delivery-and-admin.md` updated; `user-docs/reference/` regenerated
  (`collaboration.md`, and the changed `agents.md`, `secrets.md`, `admin.md`, `members.md`, `api-keys.md`,
  `reports.md`, `teams.md`); none of it links into `specs/`. The root repository's
  `skills/enterprise-agent-factory-service-cli.md` gains the new commands in the same cycle (CLI-14 workflow
  step 6). *Proven by* `docs.TestReferenceUpToDate` (`make docs` leaves no diff).

**CLI-US-35 · Swarm templates may carry a collaboration agent.** *(refines US-33 and PC-D29; `SwarmTemplate.spec.collaboration`)*
The swarm template commands take the spec as a file (`--spec file://…`), so no new option is needed.
- 35.1 `afs swarm-templates create-swarm-template … --spec file://spec.json` with `"collaboration":{"templateId":"…"}` in the
  spec sends the spec unchanged, and `describe-swarm-template --output json` prints `spec.collaboration` unchanged;
  `--generate-cli-skeleton` shows the optional member. *Proven by* `swarmtemplates.TestCreateSwarmTemplateSpecCollaborationUnchanged`,
  `swarmtemplates.TestSwarmTemplateSkeletonHasCollaboration`.
- 35.2 A template without the `paperclip` capability is the API's `ValidationError` (`field: spec.collaboration.templateId`),
  printed AWS-style with exit 254; a security operator gets `AccessDenied` as for every swarm template write. *Proven by*
  `swarmtemplates.TestSwarmTemplateCollaborationTemplateRefused254`, `swarmtemplates.TestSwarmTemplateCollaborationAccessDenied254`.
- 35.3 Pools never connect to Paperclip: `pools` commands gain no collaboration option. Priority P1, M10 (C4). *Proven by*
  `pools.TestNoCollaborationOptionOnPoolCommands`.

## 5. Non-functional requirements

| Area | Requirement |
|---|---|
| Startup | < 100 ms to first request on a warm machine |
| Reliability | Connect timeout 30 s and read timeout 60 s by default (`--cli-connect-timeout`, `--cli-read-timeout`; 0 = none). Retries per CLI-11. Ctrl-C cancels the in-flight request through its context and exits 130. `--wait` polls every 5 s |
| Observability | `--debug` writes to stderr one line per event, `<RFC 3339 time> DEBUG <component> <message>`, with method, URL, status, duration, `X-Request-Id`, credential source and retry attempt. Redacted: `Authorization`, cookies, and body fields `value`, `secretKey`, `token`, `refreshToken`, `clientSecret`, `password`. Every error prints `Request ID: <X-Request-Id>`. `User-Agent: afs-cli/<version> (<os>/<arch>)` |
| Compatibility | Works against the API contract major version it was built for; reads `X-Afs-Contract-Version` (full semver) from responses and prints one notice per session on stderr when the server's minor is newer (D47) |
| Security | No token, key or secret value in stdout, stderr, logs or shell history; files written 0600. **One exception:** `api-keys create-api-key` and `admin create-api-key` print the response's `secretKey` once to stdout (the only time it exists), with `Store this key now; it cannot be shown again.` on stderr; it is redacted in `--debug` |
| Quality | `internal/...` ≥ 90% coverage; golden files for table, text and help; every operation tested for `AccessDenied` (exit 254) |
| Accessibility | `--color off` and `NO_COLOR` respected; table output readable at 80 columns |
| Paperclip: performance | `invoke-api` adds at most one `GetTeamCollaboration` round trip (none with `--collab-url`) and streams the response body to stdout without buffering it; startup stays under 100 ms with the new service registered (`startup_bench_test.go`) |
| Paperclip: security | The AFS credential goes only to the AFS API and to the validated edge origin, over TLS 1.2 or later (the transport rules of `endpoint_url`, ADR-012); never to a redirect target; the Paperclip key only in the `Authorization` header of that one request; neither credential is in the process arguments, `--debug`, stdout (except Paperclip's own response body), stderr or any file the CLI writes; the Paperclip key file must be 0600 or stdin; a workload token is never sent to the edge. Data from the API (`restApiUrl`, `lastError`, error messages) is untrusted and passes through `output.Sanitize` (ADR-012) |
| Paperclip: reliability | Retries per CLI-US-25.8; timeouts per the existing flags; Ctrl-C cancels the in-flight request and exits 130; a `503` page or WAF body never panics the error mapper (non-JSON bodies are handled) |
| Paperclip: compatibility | A CLI built for contract 1.2 against a 1.3 server keeps working and prints the one-per-session minor-version notice (D47); `describe-agent` of such a CLI simply shows the new members. A CLI built for 1.3 against a 1.2 server gets `404` on the new paths, printed AWS-style with exit 254 |
| Paperclip: observability | `--debug` on `invoke-api` writes one line per event with component `collab` (origin read, request sent, response, retry, redaction applied) in the CLI-wide format, naming the credential source and never a credential, body or query string; every error prints `Request ID: <X-Request-Id>` from the edge's response; `User-Agent` is the CLI's. No telemetry is sent anywhere |
| Paperclip: quality | The new packages (`internal/cli/collaboration`, the changed parts of `agents`, `secrets`, `admin`, `redact`) hold at least 90% coverage; golden files for every new table, text and help page; every operation tested for `AccessDenied`; `internal/cli/services` traceability test extended to CLI-US-24 to CLI-US-34 |

### 5.1 Acceptance for the cross-cutting requirements

The Paperclip rows are in section 5.2 below the original rows.

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

### 5.2 Acceptance for the Paperclip cross-cutting requirements

- **Hints (PC-CLI-D6).** One table-driven test, `collaboration.TestCollaborationReasonHints`, walks every
  `details.reason` value of API contract §4.4b, §4.4c, §4.6 (`collaboration_agent_exists`,
  `collaboration_agent_missing`, `team_suspended`, `managed_secret`, `master_key_not_rotatable`, `reserved_name`,
  `managed_registration`, `surface_disabled`, `no_team_role`, `permission_required`, `ambiguous_credentials`,
  `malformed_credential`, `body_too_large`, `unsupported_encoding`, `token_invalid`,
  `token_expired`, `key_expired`, `key_revoked`, `session_expired`, `paused`, `starting`, `upstream_unavailable`) and
  `details.fields[].reason` values, and fails when a value has no hint or a hint contains text not in the table (a
  hint never includes server-supplied text). An unknown reason prints no hint and still exits 254.
- **Origin safety (CLI-US-25.5).** `collaboration.TestInvokeAPIOriginRefused255` is table-driven over: `http://host`,
  `https://user:pw@host`, `https://host:8443`, `https://10.0.0.1`, `https://[::1]`, `https://host/?q=1`,
  `https://host/#f`, a host outside `collab_base_domain` and an empty host; each exits 255 with no request to the
  edge and a message that names the rule, not the value of any credential.
- **No credential anywhere.** `collaboration.TestInvokeAPIDebugRedactsEverything` and
  `collaboration.TestNoSecretValueInAnyCollaborationOutput` plant distinct sentinel strings for the AFS token, API
  key, Paperclip key and a `?token=` query, run every success and failure path of CLI-US-24 to CLI-US-27 and CLI-US-25,
  and assert the sentinels are absent from stdout (except as part of a response body the fake server itself returned),
  stderr, `--debug` output and the recorded request line and arguments of the process.
- **Exit codes.** A table in `collaboration` tests pins exit 0, 130, 252, 254 and 255 for each path above: usage and
  client-side validation 252, any API or edge answer 254, local credential, file, origin and network problems 255.
- **Compatibility (CLI-27).** `api.TestSnapshotMatchesSource` and `make contract-check`, `make generate-check` pass on
  the re-pinned snapshot; a snapshot from before the re-pin (contract 1.2) still builds in a test fixture so the
  1.2 to 1.3 notice (D47) is proven; `services.TestOperationTableMatchesSnapshot` fails when a v1.3.0 operation is
  removed from the snapshot.

## 6. Dependencies

| Dependency | Use | Notes |
|---|---|---|
| `-restapi` `api/openapi.yaml` at a release tag | Generated client and the CLI-US-7 command table | Regenerated in its own PR after each contract change (contract §8) |
| oapi-codegen | Generates `internal/client` types and operations | Pinned in `tools.go`/`go.mod` |
| `github.com/jmespath-community/go-jmespath` | `--query` | |
| IdP device-grant clients (Okta, Entra ID) | `sso login` | One public client per provider with the device authorization grant enabled; its issuer and `cliClientId` come from `ListSignInProviders` (D48) |
| AWS CLI v2 on `PATH` | `aws eks get-token` in the kubeconfig exec block | Needed only by kubectl at run time |
| Signing: cosign, Apple Developer ID with notarization, Windows Authenticode certificate | CLI-15 | Certificates held in CI secrets, never in the repository |
| Homebrew tap repository | CLI-15 | Created before the first release; linked only once it exists |
| `-restapi` release whose `api/openapi.yaml` carries contract v1.3.0 (the seven §4.4b operations, `role` and `collaboration` on agents, `managedBy`, `capabilities`, `collabClientId`, `mayExposeCollaboration`, `team_collaborator`, report kind `collaboration`) | CLI-19 to CLI-27 | Not available until `-restapi` merges its C1 work; this repository re-pins then (PC-CLI-D10). The pinned snapshot, `api/SOURCE`, `internal/client/gen`, `internal/spec/operations_gen.go`, `user-docs/reference/` and every golden file that lists operations must be resynced and regenerated at that point; **do not hand-edit them now** |
| `-restapi` `afs-edge` deployed in a dev environment, with `collab.base_domain` set, a team provider that has a collaboration client, and a team with a collaboration agent | End-to-end proof of `invoke-api` (CLI-US-25) in C3; unit tests use a fake edge | Change PRD C2 and C3 exit criteria |
| Contract v1.3.1 (PC-D20 to PC-D31), which answers PC-CLI-OQ1 to OQ4 | `invoke-api` origin, path and error rules | Root; merged (1f60ecd) |
| Root `skills/enterprise-agent-factory-service-cli.md` | New commands for agents | Updated in the same cycle as the release (workflow step 6) |
| Console's T16 shows the equivalent `afs collaboration put-collaboration` command (change PRD §5) | Parity of examples | `-website`; the example text in CLI-US-24.10 is the source |

## 7. Delivery

| Milestone | This repository |
|---|---|
| M0 | Generated client from `openapi.yaml`; config, profiles, credentials, output renderers (incl. csv), errors, retries, pagination, `--wait`, help; `sso`, `configure` |
| M1 | `teams` (core, suspend and resume, models, `list-revisions`, `describe-revision`), `members`, `audit list-audit-events`, `admin get-settings`, `put-setting`, identity providers (incl. `describe-identity-provider`) |
| M2 | `agents` (core, snapshots, `list-storage-classes`, `list-agent-templates`, config values, external systems), `secrets`, `teams` config values and external systems, `models list-models`, `admin` catalog systems (incl. `describe-external-system`), storage classes and agent templates (`list-agent-templates`, `put-agent-template`, `delete-agent-template`) |
| M3 | `budgets`, `reports` (on demand), `teams` notification channels (`list-`, `create-`, `describe-`, `update-`, `delete-` and `test-notification-channel`), `admin` models and prices |
| M4 | `api-keys`, agent role links and identities, `teams list-allowed-accounts`, `restore-revision`, tags and log destinations (incl. `describe-log-destination`), report schedules, `admin` service notification channels, `admin` service tags, price proposals, cost reconciliations, org accounts, exports and imports (incl. `list-exports`), `audit list-service-audit-events`; `completion` (OQ-C5); signed releases and Homebrew tap |
| M5 | `pools`, ephemeral `create-team` options, `teams get-kubeconfig`, `admin list-cost-tag-status` |
| M6 | `work-orders` (incl. workload token credential), `swarm-templates` (incl. `enable-swarm-template`, `disable-swarm-template`), `teams list-team-directory`, `admin` service-level keys and security operators; feature complete |

Paperclip change (change PRD section 10, milestones C1 to C4; CLI milestone numbers continue M6, PC-CLI-D9). Each
milestone is a pull request into `main` that raises `CurrentMilestone` only in its exit task, and each starts only after
the `-restapi` release that provides its operations (the re-pin is the first task, PC-CLI-D10).

| Milestone | Change PRD | This repository | Exit criteria |
|---|---|---|---|
| M7 | C1 Collaboration agent | Re-pin to contract v1.3.0 (CLI-27, CLI-US-34); `collaboration describe-collaboration`, `put-collaboration` (If-Match, `--wait`), `list-collaboration-links`, `connect-agent`, `disconnect-agent` (CLI-US-24, CLI-US-26: all P0 in C1, PC-D28); `agents` flags, columns, `--role`, `--capability`, `admin put-agent-template --capabilities` (CLI-US-28); `secrets --managed-by` and hints (CLI-US-29); `teams --may-expose-collaboration`, `members` and `api-keys` `team_collaborator` (CLI-US-31.1, CLI-US-32.1); redaction of `X-Afs-Authorization` (CLI-US-33.1); reference docs and `user-docs/collaboration.md` first cut | `make gate` green; every P0 operation of §4.4b has example, columns, goldens and `AccessDenied` test; the end-to-end script of the change PRD (create the collaboration agent, connect two agents) passes against dev with `afs` only; `Pending()` empty at milestone 7 |
| M8 | C2 Public web UI | `admin list-collaboration-endpoints` (CLI-US-30, P0 in C2, PC-D28); `admin create-identity-provider` and `update-identity-provider` `--collab-client-id`, `--collab-client-secret-arn`; `collab.*` settings through `put-setting` verified (CLI-US-31.2 to 31.4); help text for the effective-surface notice (CLI-US-24.2); docs for the public endpoint | Gate green; a service admin sets the base domain, the collaboration client and a team guardrail from the CLI in dev; `Pending()` empty at milestone 8 |
| M9 | C3 REST API and managed mode | `collaboration invoke-api` (CLI-US-25, CLI-US-33.2); `rotate-collaboration-keys` (CLI-US-27); managed-mode hints (`managed_registration`) | Gate green; a pipeline creates a Paperclip issue with a collaborator team key against the dev edge; all redaction tests green; `Pending()` empty at milestone 9 |
| M10 | C4 People and reports | swarm-template `spec.collaboration` goldens (CLI-US-35); report kind `collaboration` and the `--group-by` additions (CLI-US-32.2, 32.3); final docs, root `skills/` entry, traceability test over CLI-US-24 to CLI-US-34; release | Gate green; `CurrentMilestone` is 10; every acceptance criterion of CLI-US-24 to CLI-US-34 names a passing test |

## 8. Open questions

| # | Question | Owner | Status or proposal |
|---|---|---|---|
| OQ-C1 | The API contract has no way to learn the server's contract version. Add response header `X-Afs-Contract-Version: <major.minor>` to contract §1 | root / restapi | Resolved: D47. `X-Afs-Contract-Version` is in the API contract §1 |
| OQ-C2 | Does `ListSignInProviders` return the device-grant client id, or is `clientId` the console's PKCE client? | root / restapi | Resolved: D48. `ListSignInProviders` returns `cliClientId` |
| OQ-C3 | `ListSwarmTemplates … DeleteSwarmTemplate` names `:enable` / `:disable` paths but no operations or CLI names | root | Resolved: the contract names `EnableSwarmTemplate` / `DisableSwarmTemplate` (`enable-swarm-template`, `disable-swarm-template`) |
| OQ-C4 | `ListTeamDirectory` (D38) is not yet in the contract | root | Resolved: `ListTeamDirectory` is in the contract as `teams list-team-directory` (D56) |
| OQ-C5 | Shell completion (`afs completion bash\|zsh\|fish\|powershell`) in or out of v1 | J | Resolved: accepted. `completion` is in, M4, a local command |
| OQ-C6 | Token storage: 0600 files (current) or OS keychains (macOS Keychain, Windows Credential Manager, Secret Service) | J | Files in v1; keychain as an opt-in profile key later |
| OQ-9 | Which IAM role the IdP federates to for viewer kubectl (affects the exec block of CLI-16) | J | From root |
| PC-CLI-OQ1 | Contract said `invoke-api` calls no `/v1` operation but needs the hostname | root | **Resolved, PC-D31 (contract v1.3.1):** one `GetTeamCollaboration` read unless `--collab-url`, then only the edge (PC-CLI-D1) |
| PC-CLI-OQ2 | Is `/api` Paperclip's REST prefix, and does `--path` include it? | root | **Resolved, PC-D31:** `restApiUrl` is `https://<hostname>/api`; `--path` is the full path; the CLI uses only the origin (PC-CLI-D4) |
| PC-CLI-OQ3 | Can `invoke-api` tell why a request failed when the edge answers `404`? | root | **Resolved, PC-D31:** `404` only for an unknown host or both surfaces off; with one surface off a bound caller gets `403 surface_disabled`. The `404` message points to `describe-collaboration` |
| PC-CLI-OQ4 | Do once-only Paperclip keys need special handling in `invoke-api`? | Agentic Teammates | **Resolved, PC-D31:** the edge passes them through unchanged and never logs bodies; the CLI prints the response as returned (no special case), docs warn |
| PC-CLI-OQ5 | Should `put-collaboration` fill an omitted surface from the current value instead of requiring both flags? | J | Built: both required (explicit, one API call, no hidden read). Revisit after the console and CLI are used together |

### 8.1 Decisions for the Paperclip change (aligned to contract v1.3.1) (Claude, for J's review)

| # | Question | Decision |
|---|---|---|
| PC-CLI-D1 | How does `invoke-api` learn the host? | One `GetTeamCollaboration` read (a sanctioned extra call beside CLI-10 and CLI-12), taking only the origin of `restApiUrl`; `--collab-url https://<host>` skips the read for a pipeline that already knows it. The read needs `collaboration:read`, which every role that may use the REST API holds |
| PC-CLI-D2 | What does `put-collaboration` do when the version read answers `404`? | Treat it like the existing "PUT that creates" case: send no `If-Match` and show the API's `409 collaboration_agent_missing`, rather than a `404` for an operation the user did not call |
| PC-CLI-D3 | May `invoke-api` send a workload token to the edge? | No. The edge never admits workload principals (§2.1), and a pod's projected token must not travel to a public host. This is a credential-handling rule, not a prediction of `403`, so it does not conflict with "the CLI never skips a call because it predicts a 403" |
| PC-CLI-D4 | `--path` and the `/api` in `restApiUrl` | `--path` is the full path Paperclip serves; the CLI uses only the scheme and host of `restApiUrl` (PC-CLI-OQ2) |
| PC-CLI-D5 | How is a hostile `restApiUrl` contained? | Same transport rules as `endpoint_url` plus no port, userinfo, query or IP literal, no redirects, and an optional profile key `collab_base_domain` that, when set, pins the host to `<label>.<that domain>`. The key is optional so nothing breaks, and recommended in the user docs. The service that holds the credential is the one that supplies the URL, so the residual risk equals trusting `afs-api` (ADR-012) |
| PC-CLI-D6 | Hints for reasons | The CLI prints a short fixed hint under the AWS-style error for each known `details.reason`, from one table in `internal/cli/exec`; hints name commands and permissions only, contain no server text, and an unknown reason prints none. It is help, not a decision: the API's answer and exit 254 are unchanged |
| PC-CLI-D7 | Global options on `invoke-api` | `--query`, `--output`, pagination options, `--cli-input-json`, `--generate-cli-skeleton` and `--wait` exit 252: the command prints the response unchanged, and a half-supported `--query` on non-JSON bodies would hide errors. `--include` exists for headers |
| PC-CLI-D8 | Confirmation before `rotate-collaboration-keys`? | None. The CLI has no interactive confirmations (scripts first); the effects are printed on stderr afterwards and the API audits the call |
| PC-CLI-D9 | Milestone numbering | M7 to M10 map to C1 to C4 so `milestoneOf`, `CurrentMilestone` and `Pending()` keep working per milestone; operations of later milestones exist as generated commands from the re-pin but are not "done" until their milestone |
| PC-CLI-D10 | When is the snapshot re-pinned? | Once, at the start of M7, from the first `-restapi` ref whose `api/openapi.yaml` carries all seven §4.4b operations (the provider lands them together); later `-restapi` releases are re-pinned in their own PRs as usual (contract §8). Until `-restapi` merges, this repository changes no code |
| PC-CLI-D12 | Alignment to contract v1.3.1 | PC-CLI-OQ1 to OQ4 closed by PC-D31; `host_mismatch` removed (PC-D23); priorities of links, connect, disconnect and endpoints per PC-D28; `rotate-secret` per PC-D30; `hostname_held` and swarm `spec.collaboration` added (CLI-US-31.1, CLI-US-35); the apex serves only the callback, so it is never a valid `invoke-api` origin |
| PC-CLI-D11 | Where does this PRD live? | The repository's root `enterprise-agent-factory-service-cli-PRD.md` is a stub that points to the archived spec, so the stub stays and this PRD is extended where it is, in `specs/archive/261121-cli-feature-complete/`. The implementation spec for the change is a new `specs/YYMMDD-paperclip-collaboration/` that links here, and when it starts the PRD moves with it (the stub's own rule) |

## 9. Definition of done (every change)

- [ ] Command follows naming and option rules; help with an example
- [ ] `--output json` equals the API's fields; golden tests for table and text
- [ ] No secret value, token or key can reach output, logs or history (except `secretKey` once on create)
- [ ] AWS-style errors and exit codes; `AccessDenied` test
- [ ] `make gate` green; stories' acceptance criteria covered
- [ ] For a collaboration command: no AFS token, API key, Paperclip key or secret value in stdout, stderr, `--debug`, arguments or history (sentinel tests green); `X-Afs-Authorization` redacted
- [ ] For a collaboration command: reason hints table updated for any new `details.reason`; no client-side copy of server rules (permissions, ranges, reserved names, enums beyond what the contract lists as an option)
- [ ] The pinned snapshot, generated files, reference docs and goldens regenerated from the contract, never hand-edited; `make contract-check` and `make generate-check` green; `afs --version` shows the new contract minor
- [ ] `Pending()` empty at the new `CurrentMilestone`; every new operation has an `AccessDenied` test and a golden help page
- [ ] `user-docs/` updated and the root `skills/` entry in step; no link into `specs/`

## 10. Review record

`review-prd` pass for the Paperclip change, 2026-11-22 (dev-flow `review-prd`, run by hand; the reviewer acted as the
user for step 5). Scope of the review: the additions of this change (header, section 1.1, CLI-19 to CLI-27, section 3.2
additions, CLI-US-24 to CLI-US-34, NFR rows, section 5.2, dependencies, M7 to M10, PC-CLI-OQ1 to OQ5, PC-CLI-D1 to D11,
definition of done) read against the change PRD (sections 6, 7, 10, 12) and API contract v1.3.0. Verdict after fixes:
**Ready for spec**.

| Dimension | Status | Note |
|---|---|---|
| Problem clarity | Pass | Section 1.1 names the gap (no scripting of the collaboration agent, no supported way to call Paperclip with an AFS credential) and the affected users, with a measure |
| Scope definition | Pass | Requirements CLI-19 to CLI-27 are the goals; section 3.2 gains eight explicit non-goals for the change |
| Functional completeness | Pass | Every CLI row of the change PRD (§7.6) and every operation of contract §4.4b maps to a requirement and a story; generated-only changes (roles, guardrail, capabilities, provider client, report kind, settings) are owned by CLI-26 |
| NFR coverage | Pass | Performance, security, reliability, compatibility, observability and quality rows for the change; section 5.2 gives the cross-cutting acceptance |
| Acceptance criteria | Pass | Each story has numbered Given/When/Then criteria with exact option names, exit codes, error codes, request shapes and a named test; no vague criteria remain |
| Dependency identification | Pass | The `-restapi` release (and which files must be resynced), a dev edge, contract clarifications, root skills entry, console parity |
| Open questions | Pass | PC-CLI-OQ1 to OQ5, each with owner and the default built; decisions PC-CLI-D1 to D11 |

Findings from the first read of the draft, all fixed in this pass:

- **Warning, acceptance criteria:** 31.1 claimed a `Collaboration` column on `describe-team` that no contract text supports;
  removed, the member is proven in `--output json`. 28.4 pointed at CLI-US-34 for the 80-column rule; it now names the
  existing `TestTableFitsIn80Columns`. 28.6 invented a server-supplied agent id in a hint; hints now contain no server text (PC-CLI-D6).
- **Warning, acceptance criteria:** the reports' client-side `--group-by` pattern (`internal/cli/reports/flags.go`) rejects
  `surface`, `reason` and `person`; CLI-US-32.2 now requires the check to accept them.
- **Warning, dependencies:** no statement that the pinned snapshot and everything generated from it must be resynced;
  added to section 6 with the do-not-hand-edit rule, plus PC-CLI-D10 (when to re-pin) and PC-CLI-D9 (milestones).
- **Warning, open questions:** three contract gaps affect `invoke-api` (hostname source, `/api` prefix, undetailed `404`);
  recorded as PC-CLI-OQ1 to OQ3 with the defaults built (all four since resolved by contract v1.3.1, see below).
- **Warning, NFR coverage:** the credential-handling risks of a command that sends two credentials to a public host
  were not stated; added the security NFR, the workload-token refusal (PC-CLI-D3), origin validation and no redirects
  (PC-CLI-D5, CLI-US-25.5) and the sentinel redaction tests (section 5.2).
- **Info:** the stub at the repository root is kept and this PRD is extended where it lives (PC-CLI-D11).

**Aligned to v1.3.1 (2026-11-22).** The root contracts were corrected (API contract v1.3.1, PC-D20 to PC-D31). This PRD
was updated to match: open questions PC-CLI-OQ1 to OQ4 closed (PC-D31), `host_mismatch` dropped (PC-D23), 404 versus
403 `surface_disabled` stated, priorities of links, connect, disconnect and endpoints moved (PC-D28), `rotate-secret`
and `hostname_held` criteria added (PC-D30), CLI-US-35 added for swarm template collaboration (PC-D29), decision
PC-CLI-D12. Remaining open: PC-CLI-OQ5 only. Dimensions unchanged: all Pass.
