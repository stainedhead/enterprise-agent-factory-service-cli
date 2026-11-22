# Data dictionary: `afs` CLI, feature complete

The CLI owns no server data. This file fixes its local data: the generated descriptors, the command
model, files on disk, and the wire shapes it builds or reads specially. Wire resource shapes are the
snapshot's `components.schemas` (186 schemas) and are not repeated here.

## 1. Generated: `internal/spec` (from `api/openapi.yaml` by `tools/afsgen`)

```go
// Operation describes one OpenAPI operation. One value per operation (182); generated, never edited.
type Operation struct {
    ID          string      // operationId, e.g. "UpdateAgent"
    Method      string      // GET PUT POST PATCH DELETE
    Path        string      // "/v1/teams/{teamId}/agents/{agentId}"
    Service     string      // from x-afs-cli; "" for GetHealth, GetReadiness
    Command     string      // from x-afs-cli, e.g. "update-agent"
    Summary     string      // summary/description, for help
    Permission  string      // x-afs-permission; AnyOf []string for x-afs-any-of
    AnyOf       []string
    Scope       string      // x-afs-scope: public | any | team | service
    Priority    string      // x-afs-priority: P0 | P1 | P2
    CSV         bool        // x-afs-csv
    Params      []Param     // path, query, header (excluding If-Match, Idempotency-Key, X-Request-Id)
    IfMatch     IfMatchMode // None | Required (IfMatch) | Optional (IfMatchPut)
    Idempotent  bool        // takes Idempotency-Key (every POST)
    Body        *Body       // nil when no request body
    Paginated   bool        // takes maxResults/nextToken
    ResultKey   string      // "items" | "rows" | "" (see research RQ-4)
    Success     []int       // 2xx codes, e.g. [202] for DeleteAgent, [200, 202] for TestNotificationChannel
    Produces    []string    // response media types: application/json, text/csv, application/yaml
    Example     json.RawMessage // first 2xx application/json example, for the replay test
    RequestExample json.RawMessage
}

type Param struct {
    Name     string // wire name: teamId, maxResults, dryRun
    In       string // path | query | header
    Required bool
    Schema   Schema
    Example  string
}

type Body struct {
    Required    bool
    MediaType   string   // application/json for every body in the snapshot
    SchemaName  string   // e.g. "UpdateAgentRequest"
    Fields      []Field  // top-level properties, in document order
}

type Field struct {
    Name     string // camelCase wire name
    Required bool
    WriteOnly, ReadOnly bool
    Schema   Schema
}

type Schema struct {
    Type     string   // string | integer | number | boolean | array | object | any ("" in the YAML)
    Nullable bool     // type: [T, "null"] or oneOf with null
    Enum     []string
    Min, Max *float64 // minimum/maximum
    MinLen, MaxLen *int
    Pattern  string
    Items    *Schema  // arrays
    Ref      string   // component name, e.g. "Money", "SecretValue", "Timestamp"
}

// ContractVersion is info.version of the snapshot, e.g. "1.1.0" (D-C13).
const ContractVersion = "…"
// SnapshotSHA256 mirrors api/SOURCE and is checked at test time.
const SnapshotSHA256 = "…"
```

`internal/client/gen/types.gen.go`: oapi-codegen v2.8.2 models (`generate: {models: true}`, package
`gen`). Used only to build typed request bodies in service overrides and in tests; never to decode a
response for output (D-C5).

## 2. Command model: `internal/cli`

```go
// Command is one entry of the command table: a generated operation plus its overrides.
type Command struct {
    Service, Name string          // "agents", "update-agent"
    Op            *spec.Operation // nil for local commands
    Local         LocalFunc       // configure, configure sso, sso login, sso logout, completion, help
    Options       []Option        // derived from Op (exec.DeriveOptions) then adjusted by the override
    Override      *Override
}

// Override is registered by a service package for one operationId.
type Override struct {
    Op        string
    Examples  []Example         // ≥ 1 by the operation's milestone exit (CLI-14)
    Rename    map[string]string // wire param → option, e.g. {"id": "channel-id", "rev": "revision-id"}
    Extra     []Option          // shorthand options, e.g. --repos, --tags, --secret-value, --range
    Build     BuildFunc         // optional: turns parsed options into path/query/body (shorthands)
    Validate  ValidateFunc      // optional client-side checks → exit 252
    Columns   []output.Column   // table columns; text uses the same order
    After     AfterFunc         // optional: post-response behaviour (secretKey notice, kubeconfig merge, test-result exit)
}

type Example struct{ Title, Command string } // Command starts with "afs "

type Option struct {
    Flag      string // "--team-id"
    Wire      string // "teamId"
    In        string // path | query | header | body | local
    Kind      Kind   // String | Int | Number | Bool (--x / --no-x) | List (space-separated) | JSON (object or array of objects; JSON text or file://) | Any (JSON, else string)
    Required  bool
    Enum      []string
    Min, Max  *float64
    Sensitive bool   // secret values: never logged, only file://, fileb://, - or prompt
}
```

Option derivation rules (generic, `exec.DeriveOptions`):

| Source | Option |
|---|---|
| Path or query parameter `fooBar` | `--foo-bar`; required as the parameter is; `teamId` in the path filled from profile `team_id` when omitted (D-C10) |
| Body field of scalar type | `--foo-bar` of that type; `Money` and other `string` refs stay strings (`--amount-usd 100.50` → `"100.50"`) |
| Body field `boolean` | `--foo-bar` / `--no-foo-bar` |
| Body field array of scalars | `--foo-bar v1 v2 v3` |
| Body field object, array of objects | `--foo-bar '<json>'` or `--foo-bar file://x.json` |
| Body field with no type (`PutSetting.value`) | `--value <json>`; text that is not valid JSON is sent as a JSON string |
| `readOnly` body fields | no option |
| `writeOnly` `SecretValue` fields | `--secret-value file://…|fileb://…|-` (CLI-9) |
| `version` key in `--cli-input-json` | `If-Match` (CLI-10), removed from the body |

Wire names that need a domain name (registered in overrides): `{id}` → `--channel-id`
(notification channels), `--destination-id` (log destinations), `--provider-id` (identity providers);
`{rev}` → `--revision-id`; swarm `{templateId}` → `--swarm-template-id`; report `{kind}` → `--kind`;
`UpdateWorkOrderRepo` `{host}/{owner}/{repo}` → one `--repo host/owner/repo`.

## 3. Safe-update reader table: `internal/cli/safeupdate`

```go
// Reader says how to learn the current version of the resource an If-Match operation changes.
type Reader struct {
    Op       string // the If-Match operation, e.g. "PutNamespaceConfigValue"
    ReadOp   string // "GetTeam", or a list operation, e.g. "ListNamespaceConfigValues"
    Match    string // for a list: the item field that equals a path parameter, e.g. "name"
    Version  string // JSON path of the version in the read result: "version" (default)
    NoVersion bool  // the resource has no version (PutTeamModels): send no If-Match
}
```

A test fails unless every operation with `IfMatch` or `IfMatchPut` in the snapshot has exactly one
`Reader`, every `ReadOp` exists, and every `Match` field exists in the read result schema.

## 4. Wait table: `internal/cli/wait`

```go
type Target struct {
    Ops       []string // mutations that accept --wait
    PollOp    string   // GetTeamStatus | GetAgentStatus | GetPoolStatus | GetLogDestination | GetNotificationChannel
    Ready     func(body) (done bool, failed bool) // phase Ready and observedGeneration ≥ generation of the mutation response; Failed|Degraded → failed
    DeleteGet string   // for delete-*: the Get that must answer 404
}
```

| Resource | Mutations with `--wait` | Poll |
|---|---|---|
| Team | `CreateTeam`, `UpdateTeam`, `SuspendTeam`, `ResumeTeam`, `PutTeamModels`, `DeleteTeam` (→ `GetTeam` 404) | `GetTeamStatus`; final print `GetTeam` |
| Agent | `CreateAgent`, `UpdateAgent`, `PauseAgent`, `ResumeAgent`, `RestartAgent`, `DeleteAgent` (→ `GetAgent` 404) | `GetAgentStatus`; final print `GetAgent` |
| Pool | `CreatePool`, `UpdatePool`, `PausePool`, `ResumePool`, `ScalePool` (reporting = desired), `DeletePool` (→ `GetPool` 404) | `GetPoolStatus`; final print `GetPool` |
| Log destination | `CreateLogDestination`, `UpdateLogDestination`, `DeleteLogDestination` (→ 404) | `GetLogDestination` |
| Notification channel | `CreateNotificationChannel`, `UpdateNotificationChannel`, `DeleteNotificationChannel` (→ 404) | `GetNotificationChannel` |

## 5. Files on disk

| Path | Format | Mode | Content |
|---|---|---|---|
| `~/.afs/config` | INI, AWS style: `[default]`, `[profile <name>]` | 0600 (written); read refuses nothing | `endpoint_url`, `output`, `sso_provider_id`, `sso_issuer`, `sso_client_id`, `team_id` |
| `~/.afs/credentials` | INI: `[<name>]` | must be 0600 or narrower, else exit 255 | `api_key = afsk_…` |
| `~/.afs/sso/cache/<sha1(issuer + "\n" + client_id)>.json` | JSON | must be 0600, else exit 255 | `{issuer, clientId, accessToken, refreshToken?, expiresAt (RFC 3339), tokenType}` |
| `$AFS_API_TOKEN_FILE` or `/var/run/secrets/afs/api/token` | raw JWT | not checked (mounted by Kubernetes) | projected ServiceAccount token, re-read per request |
| kubeconfig (`--kubeconfig`, first `$KUBECONFIG` entry, `~/.kube/config`) | YAML | new file 0600; existing mode kept, warning when wider | merged API kubeconfig |

`~` is `os.UserHomeDir()`; tests set `HOME` (and `USERPROFILE`) to a temp dir. `AFS_CONFIG_FILE`,
`AFS_SHARED_CREDENTIALS_FILE` override the two paths (AWS parity; used by tests and CI).

## 6. Pagination token (`--max-items` / `--starting-token`)

`NextToken` printed by the CLI is `base64url(JSON{"t": "<server nextToken or empty>", "s": <items to skip in that page>})`.
It is opaque to users, valid as long as the server token (24 h, same caller and operation).

## 7. Exit codes

| Code | Meaning |
|---|---|
| 0 | Success (including `--wait` reaching `Ready`, or `404` after a delete) |
| 130 | Interrupted (SIGINT); the in-flight request is cancelled; under `--wait` the mutation stands |
| 252 | Usage or client-side validation (unknown option, missing required, enum, range, `csv` misuse, inline secret, `--repos` slug, bundle > 100 MB) |
| 254 | Any API error response (4xx/5xx after retries), `412`, `428`, `--wait` reaching `Failed`/`Degraded`, a failed notification test |
| 255 | Network, config, file-mode, locally expired token, `--wait` timeout |

## 8. Redaction list (`internal/redact`)

Headers: `Authorization`, `Cookie`, `Set-Cookie`. JSON body fields at any depth: `value`, `secretKey`,
`token`, `refreshToken`, `accessToken`, `clientSecret`, `password`, `idToken`, `deviceCode`. Replaced by
`"[redacted]"`. Query parameters: none carry secrets (checked against the snapshot in a test).
`accessToken`, `idToken`, `deviceCode` are added to the PRD's list because the SSO flow handles them.

## 9. Special wire shapes the CLI builds

| Operation | CLI input | Request |
|---|---|---|
| `CreateWorkOrder` | `--repos github.com/payments/ledger-core=team_01…,ops@example.com …` | `repos[] {host: "github.com", repo: "payments/ledger-core", owningTeamId?, owningContact?}` |
| `UpdateWorkOrderRepo` | `--repo github.com/payments/ledger-core` | path `…/repos/github.com/payments/ledger-core` (segments escaped) |
| `PutTeamTags` | `--tags project=ledger business-unit=payments` | `{tags: [{key, value}]}` |
| `PutBudget` | `--amount-usd 100.50 --thresholds 50 80 100` | `{amountUsd: "100.50", thresholds: [50, 80, 100]}` |
| `CreateTeam --kind ephemeral` | `--expires-at --teardown-rule --budget-usd [--budget-thresholds]` | `{kind, expiresAt, teardownRule, budget: {period: "lifetime", amountUsd, thresholds?}}` (`TeamBudgetInput` has no `onExhaust`; an ephemeral budget always blocks, D-C19) |
| `CreateImport` | `--bundle file://b.json --account-map file://m.json [--no-dry-run]` | `?dryRun=true|false`, body `{bundle: <b.json>, accountMap: <m.json>}` |
| `CreateSecret`, `PutSecretValue` | `--secret-value file://…|fileb://…|-`, or prompt | `value` (string, UTF-8, ≤ 65536) |
| `RestoreRevision` | `--dry-run` | body `{dryRun: true}` |
| reports | `--range --from --to --group-by --tag k=v --pool-id` | query `range`, `from`, `to`, `groupBy`, `tag.k=v`, `poolId` |
