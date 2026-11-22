# Product details

The command inventory as built. The command reference generated from the same registry as `afs help` is in
[`user-docs/reference/`](../user-docs/reference/README.md); the task guides are in [`user-docs/`](../user-docs/README.md).

## Services

Operation names are the kebab-case form of the contract's `x-afs-cli` names and map one-to-one to OpenAPI
`operationId`s (`list-teams` ↔ `ListTeams`). A test proves the command table and the contract agree in both
directions.

| Service | Operations | Covers |
|---|---|---|
| `teams` | 37 | Teams (persistent and ephemeral), status, suspend and resume, tags, models, configuration values, external systems, log destinations, notification channels, revisions and restore, the team directory, `get-kubeconfig` |
| `agents` | 26 | Agents and their lifecycle (pause, resume, restart), identities, role links, snapshots, configuration values, external-system state, agent templates and storage classes |
| `pools` | 26 | Worker pools, stages, scaling (`scale-pool --wait`), pause and resume, handoffs, secret bindings, role links, configuration values, external systems |
| `admin` | 46 | Service settings, identity providers, agent templates, service tags, notification channels, storage classes, external systems, models and prices, price proposals, cost reconciliations and cost-tag status, organization accounts, exports and imports, service API keys, security operators |
| `work-orders` | 8 | Work orders and their repositories (`--repos host/owner/repo[=owners]`); workers report progress with `update-repo` |
| `secrets` | 7 | Secret references and write-only secret values (create, rotate, delete) |
| `reports` | 7 | Agent, team and service reports (`--range`, `--from`, `--to`, `--group-by`, `--tag`, `--pool-id`), CSV, report schedules |
| `swarm-templates` | 7 | Swarm templates |
| `budgets` | 4 | Budgets and thresholds, usage, costs |
| `members` | 4 | Team members and their roles |
| `api-keys` | 3 | Team API keys (the secret key is shown once) |
| `audit` | 2 | Team and service audit event search |
| `models` | 1 | Models available to the caller |
| `sso` | 4 | `login` and `logout` (local), `get-caller-identity`, `list-providers` |
| `configure` | local | `afs configure`, `afs configure sso` |
| `completion` | local | `afs completion bash|zsh|fish|powershell` |

## Behaviour common to every operation

| Feature | Behaviour |
|---|---|
| Global options | `--profile`, `--endpoint-url`, `--output`, `--query`, `--no-paginate`, `--max-items`, `--page-size`, `--starting-token`, `--cli-input-json`, `--generate-cli-skeleton`, `--debug`, `--no-cli-pager`, `--color`, `--cli-connect-timeout`, `--cli-read-timeout`, `--version`, accepted anywhere on the line |
| Output | `json` is the API's body re-indented (field names, order, `null`s and unknown fields preserved); `yaml`, `table`, `text`, `csv` for people and spreadsheets; `--query` (JMESPath) applies before rendering |
| Pagination | Automatic; `--max-items` prints a `NextToken` that `--starting-token` resumes; `--page-size` 1 to 100 |
| Safe updates | `If-Match` on every update; without `--if-match` the CLI reads the current version and prints `Read the current version N; sending If-Match: "N"` on stderr; `412` reports the current version |
| Idempotency and retries | Creates send an `Idempotency-Key` (ULID); at most 3 attempts with full-jitter backoff or `Retry-After`; a request is re-sent after a dropped connection only when it is a `GET` or carries an idempotency key |
| `--wait` | Polls until `Ready` at the mutation's generation; `Failed` or `Degraded` exits 254; `--wait-timeout` (default 10m) |
| Errors | `An error occurred (<Code>) when calling the <Operation> operation: <message>` and `Request ID:` on stderr |
| Exit codes | 0 success, 130 interrupted, 252 usage, 254 service error, 255 other (network, configuration) |
| Help | `afs help`, `afs <service> help`, `afs <service> <operation> help`; every operation has an example |

## Credentials

Checked in order: `AFS_API_KEY`, `api_key` in `~/.afs/credentials` (must be 0600), the SSO token cache for the profile,
then the workload token file (`AFS_API_TOKEN_FILE`, default `/var/run/secrets/afs/api/token`, re-read for every request).

## Known contract gaps

These are behaviours the CLI handles today and that a contract change would simplify. Each is recorded as a contract
delta for `-restapi`:

- `work-orders update-repo` needs `--if-match` from a workload principal, because that principal cannot list a work
  order's repositories to read the version (CLI-CD-1).
- `pools describe-pool` shows the desired worker count; the reporting, working, idle and impaired counts are in
  `pools get-pool-status` (CLI-CD-6).
- `412 PreconditionFailed` carries no version; the CLI uses `ETag` when present, else one more read (CLI-CD-3).
