# Technical details

## Toolchain

Go 1.27.4, golangci-lint v2.15.0, govulncheck v1.8.0, oapi-codegen v2.8.2 (types, run from the Makefile).

## How the commands exist

`api/openapi.yaml` is a pinned copy of the API contract (`api/SOURCE` records repo, ref, commit, sha256).
`tools/afsgen` turns it into `internal/spec/operations_gen.go` (one descriptor per operation: path,
parameters, body fields, `x-afs-*` facts). `internal/cli.DeriveOptions` turns a descriptor into options.
`internal/cli/<service>` packages add per-operation overrides (examples, columns, shorthands) and local
commands. `internal/cli/services` assembles the registry and wires the executor and help.

## Request path

`cli.Run` parses (`internal/cli/parse.go`) and calls the executor (`internal/cli/exec`): `Build` merges
options and `--cli-input-json`; `safeupdate` reads the current version for `If-Match`; `internal/client`
sends the request (credential chain in `internal/auth`, retries, idempotency keys, redacted debug log);
pagination merges pages; `wait` polls for `--wait`; `internal/output` renders; errors map to exit codes.

## Package map

| Package | Responsibility |
|---|---|
| `cmd/afs` | Entry point: calls `cli.Run` with the process streams and exits with its code |
| `internal/cli` | Parser, global options, option derivation, registry, override merging, milestone table |
| `internal/cli/<service>` | Per-service overrides (examples, columns, shorthands, validators, body builders) and local commands |
| `internal/cli/services` | Composition: assembles every service into the registry and wires executor and help |
| `internal/cli/exec` | Generic executor: build request, paginate, render, map errors to exit codes, skeletons |
| `internal/cli/safeupdate` | Reader table for the 59 `If-Match` operations; read-then-write and `412` handling |
| `internal/cli/wait` | Wait table and polling for `--wait` |
| `internal/cli/help`, `completion` | Help pages and shell completion rendered from the registry |
| `internal/spec` | Generated operation descriptors (`operations_gen.go`) |
| `internal/client` | HTTP transport, timeouts, retries, idempotency keys, contract-version notice, redacted debug log; `gen/` holds generated types |
| `internal/auth` | Credential chain: `AFS_API_KEY`, credentials file, SSO cache, workload token |
| `internal/sso` | OIDC device authorization grant, token cache (0600), silent refresh |
| `internal/config` | INI parsing of `~/.afs/config` and `~/.afs/credentials`, precedence |
| `internal/output` | json, yaml, table, text and csv renderers, JMESPath `--query`, pager, colour |
| `internal/redact` | Header and JSON-field redaction for anything logged |
| `internal/secretinput` | Secret values from `file://`, `fileb://`, stdin or a no-echo prompt |
| `internal/kubeconfig` | Node-preserving kubeconfig merge and atomic write for `teams get-kubeconfig` |
| `tools/afsgen`, `tools/docgen` | Descriptor generator and command-reference generator |

`internal/cli` takes `io.Writer`s, an environment map and an injected HTTP client, so every command is tested against an
`httptest` server without a network.

## Security controls

- **Transport:** `endpoint_url` must be `https://`; `http://` is accepted only for loopback hosts; credentials, a query
  or a fragment in the URL are refused. Proxies come from
  `HTTPS_PROXY`/`NO_PROXY`; custom roots from `SSL_CERT_FILE`/`SSL_CERT_DIR`.
- **Redaction:** `--debug` logs pass through `internal/redact`, which masks `Authorization`, `Cookie`,
  `Set-Cookie` and `Proxy-Authorization` headers and secret-bearing JSON fields at any depth.
- **Secret input:** secret values are never accepted as plain arguments; an inline `--cli-input-json` that carries a
  value is refused (252).
- **File modes:** the credentials file and SSO cache must be 0600 or they are refused; files the CLI creates are 0600
  in 0700 directories; the kubeconfig is written atomically (temp file, fsync, rename).
- **Untrusted service data (ADR-012):** `config.File.Set` refuses control and structural characters, so no value can add
  lines or sections to `~/.afs/config`; SSO provider fields are validated before they are stored; OIDC discovery must be
  https (loopback http only), name the configured issuer, and return a Bearer access token; `table`, `text`, errors and
  prompts pass through `output.Sanitize`; the service's kubeconfig is merged only when it runs `aws eks get-token` and its
  entries start with `afs`.
- **One-time secrets:** a created API key's `secretKey` is printed once with a stderr warning; a replay never prints it.

## Timeouts, retries and waits

`--cli-connect-timeout` (default 30 s) bounds dial and TLS handshake; `--cli-read-timeout` (default 60 s) bounds the wait
for headers and each body read; `0` disables either. Retries: at most 3 attempts, full-jitter backoff from 200 ms, or
`Retry-After` when longer, capped at 20 s or the read timeout when smaller (over the cap the error is returned). Identity
provider calls share the connect and read timeouts, and `--debug` logs them. `--wait` polls every 5 s (or `Retry-After`,
capped) until `--wait-timeout` (default `10m`). Pagination stops with exit 255 on a repeated page token.

## Exit codes

0 success · 130 interrupted · 252 usage/parse error (including a missing operation) · 254 service error · 255 other
(AWS CLI compatible).

## Configuration files

`~/.afs/config` (profiles), `~/.afs/credentials` (API keys, mode 0600), `~/.afs/sso/cache/` (tokens, mode 0600).

## Releases

GoReleaser v2 (`.goreleaser.yaml`, `.github/workflows/release.yml`) builds darwin/arm64, darwin/amd64, linux/arm64,
linux/amd64 and windows/amd64 archives with `CGO_ENABLED=0`, signs `checksums.txt` keyless with cosign, notarizes macOS
binaries, signs Windows binaries with Authenticode and updates the Homebrew tap. Signing jobs run only when their secrets
exist. CI (`.github/workflows/ci.yml`) runs the gate, govulncheck, a cross-build matrix and a GoReleaser snapshot.

## Quality gate

`make gate`: `go mod tidy`, `contract-check` (snapshot matches `api/SOURCE`), `generate-check` (generated files
match the snapshot), `fmt-check`, `go vet`, golangci-lint, race tests, `cover` (`internal/...` at least 90%,
generated code and the golden helper excluded), build. The contract replay test sends every operation's
documented example to a fake server and validates the request against the snapshot.
