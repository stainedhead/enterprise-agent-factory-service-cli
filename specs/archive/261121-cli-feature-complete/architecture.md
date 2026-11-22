# Architecture: `afs` CLI, feature complete

| | |
|---|---|
| **Feature** | cli-feature-complete |
| **Date** | 2026-11-21 |
| **Status** | Draft for implementation |

## 1. Overview

```
                     api/openapi.yaml  (snapshot of -restapi, sha256 in api/SOURCE)
                            │  make generate
              ┌─────────────┴──────────────┐
     tools/afsgen                    oapi-codegen v2.8.2
              │                            │
 internal/spec/operations_gen.go   internal/client/gen/types.gen.go
              │                            │ (request bodies in overrides/tests only)
              ▼                            │
cmd/afs ──► internal/cli  (registry = generated operations + service overrides + local commands)
              │  parse globals → resolve profile → pick command → derive options → parse
              │
              ├─► internal/cli/exec       build request · paginate · render · map errors → exit code
              │      ├─► internal/cli/safeupdate   If-Match reader table (CLI-10)
              │      └─► internal/cli/wait         status polling (CLI-12)
              ├─► internal/cli/<service>  overrides: examples, renames, shorthands, columns, after-hooks
              ├─► internal/output         json · yaml · table · text · csv · --query · pager · color
              ├─► internal/config         ~/.afs/config, ~/.afs/credentials, env precedence
              ├─► internal/auth           credential chain (env key → credentials → SSO → workload token)
              ├─► internal/sso            device grant, cache, refresh
              ├─► internal/secretinput    file:// fileb:// - prompt
              ├─► internal/kubeconfig     merge for teams get-kubeconfig
              └─► internal/client         HTTP transport: auth header, UA, retries, idempotency,
                                          timeouts, debug log with redaction, contract-version notice
```

Dependency rule (AGENTS.md): `cmd/afs` → `internal/cli` → (`internal/client`, `internal/config`,
`internal/output`, `internal/auth`, `internal/sso`, …). Leaf packages never import `internal/cli`.
Command handlers parse, call one client method (plus the sanctioned read and polls), render — no
business logic. `internal/cli` takes `io.Writer`s, an injected `Env` (environment, home dir, stdin,
TTY flags, clock, HTTP transport) so every test runs without a network or a real home.

## 2. Components

| Package | Responsibility | Key types |
|---|---|---|
| `api` | Snapshot, `SOURCE` (repo, ref, commit, sha256, synced-at), `embed.go` (`//go:embed openapi.yaml`, imported by tests only) | `api.OpenAPI []byte` |
| `tools/afsgen` | Reads the snapshot (`gopkg.in/yaml.v3`, local `$ref` resolution), writes `internal/spec/operations_gen.go` sorted by `operationId`; refuses an operation without `x-afs-cli` | — |
| `tools/docgen` | Writes `user-docs/reference/<service>.md` from the command table and help text | — |
| `internal/spec` | Generated `[]Operation`, `ContractVersion`, `SnapshotSHA256`; lookup by id | `Operation`, `Param`, `Body`, `Schema` |
| `internal/client/gen` | Generated request/response types | — |
| `internal/client` | `Client.Do(ctx, Request) (*Response, error)`: URL build (path escaping per segment), headers (`Authorization`, `User-Agent`, `Accept`, `Content-Type`, `If-Match`, `Idempotency-Key`), retries, timeouts, debug log, `X-Afs-Contract-Version` notice, error decode | `Request`, `Response{Status, Header, Body []byte}`, `APIError{Status, Code, Message, RequestID, Details json.RawMessage}` |
| `internal/redact` | Header and JSON-field redaction, shared by `client` and `cli` | `Headers`, `JSON` |
| `internal/config` | INI read/write (0600, other profiles untouched), `Resolve(opts, env) Settings` with precedence option → env → profile → default | `Profile`, `Settings` |
| `internal/auth` | `Resolve(settings, env) (Credential, Source)`; each credential is a `func(ctx) (header string, err)` so the workload token is re-read per request and SSO can refresh | `Credential`, `Source` (for `--debug`) |
| `internal/sso` | `Login` (discovery → device authorization → poll token endpoint), `Cache` (load/save/delete, 0600 checks), `Token(ctx)` with refresh | `Cache`, `Token` |
| `internal/output` | Renderers; `Query(jmespath)`; table at ≤ 80 columns with truncation marks; text (AWS tab-separated); csv writer passthrough with header de-duplication; pager; colour | `Renderer`, `Column` |
| `internal/secretinput` | Read a secret from `file://`, `fileb://`, `-`, or a no-echo prompt; refuses anything else | `Read(arg, env) ([]byte, error)` |
| `internal/kubeconfig` | Parse, merge by name (clusters, users, contexts), optional exec `env` injection, atomic write with mode rules | `Merge`, `WriteFile` |
| `internal/cli` | Registry, global options, dispatch, local commands, exit codes | `Registry`, `Command`, `Override`, `Option` |
| `internal/cli/exec` | Generic executor: options → request, `--cli-input-json`, `--generate-cli-skeleton`, pagination, CSV paging, render, error → exit | `Run(ctx, cmd, args) int` |
| `internal/cli/safeupdate` | Reader table and read-then-write logic, `412`/`428` messages | `Reader` |
| `internal/cli/wait` | Wait table and polling loop (5 s or `Retry-After`, `--wait-timeout`, Ctrl-C) | `Target` |
| `internal/cli/help` | Help rendering for three levels, golden files, examples check | — |
| `internal/cli/services` | Imports every service package and calls its `Register(*cli.Registry)` | — |
| `internal/cli/<service>` | One file per resource group; each registers `Override`s | — |

## 3. Layer responsibilities

- **Generated layer** (`internal/spec`, `internal/client/gen`): facts from the contract only. Never
  hand-edited; `make generate-check` fails CI on drift.
- **Transport** (`internal/client`): knows HTTP, retries and redaction; knows nothing about commands.
- **Executor** (`internal/cli/exec`): knows the conventions every operation shares (options,
  pagination, safe updates, wait, errors). Operation-specific behaviour is not allowed here; it goes in
  an `Override`.
- **Service overrides**: shorthands, domain option names, examples, table columns, after-hooks. Each is
  registered per `operationId` from a file owned by exactly one stream (tasks.md).

## 4. Data flow: one command

1. `cmd/afs/main.go` builds `Env` with `signal.NotifyContext(SIGINT)` and calls `cli.Run(ctx, args, env)`.
2. `cli.Run` parses global options anywhere on the line (AWS style), resolves settings (`internal/config`),
   looks up `<service> <operation>` in the registry (unknown → 252 with usage).
3. `exec.Run` derives options, applies `--cli-input-json`, then explicit options (they win), then the
   override's `Build`/`Validate`; any failure → 252 before a call. `--generate-cli-skeleton` prints and
   returns 0 here.
4. Safe update (if `IfMatch != None`): `--if-match` or input `version` → header; else the `Reader`'s
   read, notice on stderr; `404` on an optional-`If-Match` `PUT` → no header.
5. `client.Do` sends with credentials from `internal/auth`; retries per CLI-11 reuse the same
   `Idempotency-Key`.
6. Pagination loop (auto, `--max-items`, `--starting-token`, `--no-paginate`); JSON pages are merged
   into one `items`/`rows` array, CSV pages concatenated under one header.
7. Render: `--query` then renderer; `json` without `--query` re-indents the raw body bytes.
8. `--wait`: poll the `Target` until done/failed/timeout; print the final resource.
9. Errors: `An error occurred (<code>) when calling the <operationId> operation: <message>` +
   `Request ID: <id>` on stderr; exit code per data-dictionary §7.

## 5. Sequences

### 5.1 `afs agents update-agent --team-id t --agent-id a --description x` (no `--if-match`)
```
cli → client: GET /v1/teams/t/agents/a                       (Reader: GetAgent)
client ← 200 {version: 7}
cli → stderr: Applied to version 7; use --if-match to guard against concurrent changes.   (printed after success)
cli → client: PATCH /v1/teams/t/agents/a  If-Match: "7"  {description: x}
client ← 200 Agent → stdout (json)
       ← 412 → cli reads again (GetAgent) → "The resource changed (version 8). Re-run with --if-match 8 …", exit 254
```

### 5.2 `afs sso login`
```
cli → afs-api: GET /v1/identity-providers (no credential)  → pick sso_provider_id (or prompt in configure sso)
cli → IdP: GET {issuer}/.well-known/openid-configuration
cli → IdP: POST device_authorization_endpoint (client_id = cliClientId, scope "openid offline_access")
cli → stderr: verification URL + user code
loop: POST token_endpoint (device_code) every interval, honour slow_down / authorization_pending
cache ← {accessToken, refreshToken, expiresAt} (0600)
```

### 5.3 CSV auto-pagination
```
GET …?…  Accept: text/csv  → 200 body1, X-Afs-Next-Token: n1 → write body1 whole
GET …&nextToken=n1          → 200 body2, X-Afs-Next-Token: n2 → write body2 minus its header row
GET …&nextToken=n2          → 200 body3 (no header)            → write body3 minus header; stop
```

## 6. Integration points

| System | How | Tested with |
|---|---|---|
| `afs-api` | HTTPS, OpenAPI snapshot | `httptest.Server` fakes; replay test validating each request with `libopenapi-validator` against the snapshot |
| IdP (Okta, Entra ID) | OIDC discovery, RFC 8628 device grant, refresh grant | `httptest` IdP fake |
| Kubernetes config | file merge only; `aws` looked up on `PATH` for a warning | temp dirs, fake `PATH` |
| Release | GoReleaser v2, cosign, quill, osslsigncode, Homebrew tap | `goreleaser check`, `release --snapshot --clean --skip=sign,publish,notarize` in CI |

## 7. Architectural decisions

The full decision list (D-C1..D-C35) is in `implementation-notes.md`. Repository ADRs to write in
`docs/architectural-decision-record.md` during implementation:

- ADR-002 (accept): client generated from a pinned snapshot; snapshot + checksum + regeneration check.
- ADR-003: command table and operation descriptors generated by `tools/afsgen` from `x-afs-*`; one
  generic executor; per-operation overrides.
- ADR-004: `--output json` prints the raw response body (re-indented); generated types never decode
  output.
- ADR-005: AWS-style INI config and credentials files; SSO cache keyed by issuer and client id.
- ADR-006: GoReleaser v2 OSS for signed releases (cosign, quill notarization, osslsigncode).
