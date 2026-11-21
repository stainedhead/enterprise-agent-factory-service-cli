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

## Exit codes

0 success · 130 interrupted · 252 usage/parse error · 254 service error · 255 other (AWS CLI compatible).

## Configuration files

`~/.afs/config` (profiles), `~/.afs/credentials` (API keys, mode 0600), `~/.afs/sso/cache/` (tokens, mode 0600).

## Quality gate

`make gate`: `go mod tidy`, `contract-check` (snapshot matches `api/SOURCE`), `generate-check` (generated files
match the snapshot), `fmt-check`, `go vet`, golangci-lint, race tests, `cover` (`internal/...` at least 90%,
generated code and the golden helper excluded), build. The contract replay test sends every operation's
documented example to a fake server and validates the request against the snapshot.
