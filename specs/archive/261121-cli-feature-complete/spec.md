# Spec: `afs` CLI, feature complete

| | |
|---|---|
| **Spec** | `specs/261121-cli-feature-complete/` (created 2026-11-21) |
| **Source PRD** | [`enterprise-agent-factory-service-cli-PRD.md`](enterprise-agent-factory-service-cli-PRD.md) (moved here by `create-spec`) |
| **Generation source** | `-restapi` `api/openapi.yaml`, OpenAPI 3.1, `info.version` 1.1.0, 182 operations (180 CLI commands + `GetHealth`, `GetReadiness`), read on `-restapi` branch `feat/restapi-feature-complete` at commit `081de14` (sha256 `3456b9e9b0f2c138a3c92ded8edde0422dbe76b3198f93d501fbe74343fb2b0e`) |
| **Contracts** | Root `docs/contracts/api-contract.md` (v1.2.0 on root branch `contract/restapi-m0-deltas`, which aligns the document to `-restapi` deltas CD-1..CD-10); names follow it exactly. Where the YAML and the document differ, the YAML wins (contract header) |
| **Scope** | Feature complete: every milestone M0–M6 of PRD §7, every operation of the contract §4 |
| **Branch** | `feat/cli-feature-complete` |

## 1. Executive summary

Build `afs`, the AWS-CLI-style command line for the Enterprise Agent Factory Service, to feature
completeness: one command per API operation (180), the local commands (`configure`, `configure sso`,
`sso login`, `sso logout`, `completion`, `help`), AWS-style global options, profiles, SSO device grant,
API-key and workload-token credentials, five output formats, JMESPath `--query`, pagination, safe
updates, idempotent retries, `--wait`, generated help, and signed releases. The command table and the
HTTP operation descriptors are **generated** from a pinned, checksummed snapshot of `-restapi`'s
`api/openapi.yaml`; command behaviour is a generic executor plus small per-service overrides.

## 2. Problem statement

People and pipelines need everything the console offers from a shell, with the habits of the AWS CLI
and with a stable machine contract (`--output json`, exit codes). The repository today is a scaffold:
`afs help` and `afs --version` only, and its service list (`identities`, `usage`, `logs`,
`notifications`) is from solution design §11, which the contract supersedes.

## 3. Goals and non-goals

### Goals
- G1. A command for every operation in the contract, named by `x-afs-cli` and mapped to its
  `operationId`; a test proves the table and the OpenAPI agree both ways (CLI-1, CLI-US-7).
- G2. The client is generated from a pinned snapshot of the OpenAPI, with a checksum check and a
  regeneration-drift check in CI, and a documented sync procedure (PRD §6; contract §8).
- G3. AWS CLI ergonomics: global options, profiles, output formats, query, pagination, skeletons,
  exit codes, help (CLI-2..CLI-8, CLI-13, CLI-14).
- G4. Safety: no secret, token or key reaches output, logs or history (except `secretKey` once);
  `If-Match` on every update; idempotent retries (CLI-9..CLI-11).
- G5. Signed, notarized releases with a Homebrew tap (CLI-15).

### Non-goals (PRD §3.2)
- Client-side workflows that combine several mutations.
- Name or slug resolution (`--agent-name`, `=team-slug`); commands take ids.
- Any command that reads a secret value.
- Client-side access decisions; the CLI never skips a call because it predicts `403`.
- TUI, plugins, aliases, auto-update.
- Local state beyond `~/.afs/config`, `~/.afs/credentials`, the SSO cache (and the `GetMe` cache,
  which this spec decides not to build, see D-C12).
- `afs configure set|get|list` (not in the PRD; `configure` and `configure sso` only).
- OS keychain token storage (OQ-C6: files in v1).

## 4. User requirements

### 4.1 Users
Service admins, team-admins, viewers/FinOps/SRE, CI pipelines, security operators and security tooling,
agents and pool workers (PRD §2).

### 4.2 Functional requirements

Each FR is the PRD requirement of the same number; the PRD text is normative, this table records how
the spec resolves it.

| FR | PRD | Requirement (summary) | Spec resolution |
|---|---|---|---|
| FR-001 | CLI-1 | A command per contract operation, named by the contract `CLI` column; ids, never names | Generated from `x-afs-cli` in the snapshot; `operations_gen.go` + registry; CLI-US-7 test |
| FR-002 | CLI-2 | Global options on every command (16 listed) | `internal/cli/globals.go`; one table-driven test per option (§5.1) |
| FR-003 | CLI-3 | Local commands; `~/.afs/config` 0600; profile keys; precedence | `internal/config`, `internal/cli/configure`, `internal/cli/completion` |
| FR-004 | CLI-4 | `sso login` (device grant), `logout`, `get-caller-identity`; cache 0600; silent refresh; expired → 255 / server 401 → 254 | `internal/sso`, `internal/cli/sso` |
| FR-005 | CLI-5 | Credential chain `AFS_API_KEY` → credentials file → SSO cache → workload token file; 0600 enforcement | `internal/auth` |
| FR-006 | CLI-6 | `json` (API fields exactly), `yaml`, `table`, `text`, `csv`; `--query` before rendering; csv rules and `X-Afs-Next-Token` pagination | `internal/output`; `json` re-indents the raw body (D-C5) |
| FR-007 | CLI-7 | Auto-pagination; `--page-size` 1–100; `--max-items`; `--starting-token`; `--no-paginate` | `internal/cli/exec/paginate.go`; CLI token format D-C7 |
| FR-008 | CLI-8 | AWS-style errors + `Request ID:`; exit codes 0/130/252/254/255 | `internal/cli/exec/errors.go` |
| FR-009 | CLI-9 | Secret values from file, stdin or no-echo prompt only; inline `--cli-input-json` with `value` refused | `internal/secretinput`, `internal/cli/secrets` |
| FR-010 | CLI-10 | `If-Match` on every update; read-then-write notice; `412`/`428` handling | `internal/cli/safeupdate` with a reader table checked against the snapshot (D-C8) |
| FR-011 | CLI-11 | `Idempotency-Key` ULID; retry policy; replay printing | `internal/client/retry.go` |
| FR-012 | CLI-12 | `--wait` polling, timeouts, delete waits for `404`, `scale-pool --wait` | `internal/cli/wait` with a wait table (D-C9) |
| FR-013 | CLI-13 | Report flags `--range`, `--from`, `--to`, `--group-by`, `--tag`, `--pool-id`; csv | `internal/cli/reports` + `internal/cli/budgets` |
| FR-014 | CLI-14 | Help at three levels, generated from the command table; an example per operation | `internal/cli/help`; examples registered per service |
| FR-015 | CLI-15 | Signed binaries 5 platforms, Homebrew tap, cosign-signed checksums, `--version` format | GoReleaser v2 config + release workflow (D-C14) |
| FR-016 | CLI-16 | `teams get-kubeconfig` merge semantics, modes, `--aws-profile`, `-` | `internal/kubeconfig`, `internal/cli/teams/kubeconfig.go` |
| FR-017 | CLI-17 | `pools` and `swarm-templates` services; ephemeral `create-team` options | `internal/cli/pools`, `internal/cli/swarmtemplates`, `internal/cli/teams/create_ephemeral.go` |
| FR-018 | CLI-18 | `work-orders` service, `--repos` shorthand and owner classification; resolution errors printed | `internal/cli/workorders` |

### 4.3 User stories

CLI-US-1..CLI-US-23 of the PRD §4 are adopted unchanged as acceptance criteria. Their task mapping is
in `tasks.md` §"Traceability"; their wire-level interpretation (where the OpenAPI fixes a detail the
PRD leaves open) is in `implementation-notes.md` §"Decisions".

## 5. Non-functional requirements

| Area | Requirement (PRD §5) | Verification |
|---|---|---|
| Startup | < 100 ms to first request on a warm machine | `BenchmarkStartup` + a CI timing check on `afs --version` and a first request to a local fake (M0-X1) |
| Reliability | Connect 30 s, read 60 s, `0` = none; retries per CLI-11; Ctrl-C → 130; `--wait` every 5 s | `internal/client` tests with injected dialer and clock |
| Observability | `--debug` line format, redaction list, `Request ID:`, `User-Agent: afs-cli/<v> (<os>/<arch>)` | `internal/client/debug_test.go` golden |
| Compatibility | Built for one contract major; reads `X-Afs-Contract-Version`; one notice per session when the server's minor is newer | `internal/client/version_test.go` |
| Security | No token, key or secret value in stdout/stderr/logs/history; files 0600; `secretKey` printed once on create | Redaction tests, file-mode tests, `secretKey` tests (M4-A, M6-D) |
| Quality | `internal/...` ≥ 90% coverage; golden table/text/help; `AccessDenied` test per operation | `make cover` gate (M0-K1 makes it fail under 90%); contract replay test (M0-E7) |
| Accessibility | `--color off`, `NO_COLOR`; table readable at 80 columns | `internal/output` tests |

## 6. System architecture

See `architecture.md`. Layers: `cmd/afs` → `internal/cli` (registry, executor, per-service overrides)
→ `internal/client` (generated descriptors and types, transport), `internal/config`, `internal/auth`,
`internal/sso`, `internal/output`, `internal/kubeconfig`, `internal/secretinput`. Generators live in
`tools/`; the snapshot lives in `api/`.

## 7. Scope of changes

### 7.1 Create
- `api/openapi.yaml` (snapshot), `api/SOURCE`, `api/embed.go`, `api/snapshot_test.go`
- `tools/afsgen/` (descriptor and command-table generator), `tools/docgen/` (command reference)
- `internal/spec/operations_gen.go` (generated descriptors)
- `internal/client/gen/types.gen.go` (oapi-codegen types), `internal/client/gen/oapi-codegen.yaml`
- `internal/client/` (transport, retry, debug, version notice, errors)
- `internal/config/`, `internal/auth/`, `internal/sso/`, `internal/output/`, `internal/secretinput/`,
  `internal/kubeconfig/`, `internal/redact/`
- `internal/cli/{registry,globals,parse}.go`, `internal/cli/exec/`, `internal/cli/safeupdate/`,
  `internal/cli/wait/`, `internal/cli/help/`, `internal/cli/services/`, `internal/cli/milestones.go`
- One package per service: `internal/cli/{teams,members,apikeys,agents,secrets,budgets,reports,models,audit,admin,pools,swarmtemplates,workorders,sso,configure,completion}/`
- `.goreleaser.yaml`, `.github/workflows/release.yml`, contract checks in `.github/workflows/ci.yml`
- `user-docs/` getting started, configuration, SSO, credentials, output, troubleshooting, generated
  command reference

### 7.2 Modify
- `internal/cli/cli.go` (replace the §11 service list with the registry), `cmd/afs/main.go` (signal
  context, exit code), `internal/buildinfo/buildinfo.go` (contract version in `--version`)
- `Makefile` (`generate`, `generate-check`, `contract-check`, `contract-sync`, `cover` threshold,
  `docs`, `release-snapshot`), `go.mod`/`go.sum`
- `README.md`, `docs/technical-details.md`, `docs/architectural-decision-record.md` (ADR-002 accepted,
  ADR-003..ADR-006), `AGENTS.md` layout table

### 7.3 Dependencies (pinned in M0-K1; versions at least two weeks old on 2026-11-21)
| Module | Version | Use |
|---|---|---|
| `github.com/oapi-codegen/oapi-codegen/v2` (run via `go run …@v2.8.2` in the Makefile) | v2.8.2 | Types from the snapshot; spiked on 2026-11-21 against the 3.1 snapshot: 7,437 lines, `go vet` clean |
| `github.com/oapi-codegen/runtime` | v1.9.1 | Types used by generated code |
| `github.com/jmespath-community/go-jmespath` | v1.1.1 | `--query` |
| `gopkg.in/yaml.v3` | v3.0.1 | `--output yaml`, kubeconfig merge, generator input (same as `-restapi`) |
| `github.com/oklog/ulid/v2` | v2.1.2 | `Idempotency-Key` |
| `golang.org/x/term` | v0.47.0 | No-echo prompt, TTY detection |
| `golang.org/x/oauth2` | v0.37.0 | Device authorization grant and refresh |
| `github.com/pb33f/libopenapi-validator` (tests only) | v0.16.5 (with `libopenapi` v0.44.0, as `-restapi`) | Validates every request the CLI sends against the snapshot in the contract replay test |

## 8. Breaking changes

- The service list changes from solution design §11 to the contract: `identities`, `usage`, `logs` and
  `notifications` services disappear (their operations live under `agents`, `budgets`, `teams`);
  `reports`, `pools`, `swarm-templates`, `completion` appear. The scaffold had no implemented
  operations, so nothing released breaks.
- No API, config or schema change in other repositories. Contract deltas the CLI needs from `-restapi`
  are requests, recorded in `implementation-notes.md` §"Contract deltas requested" (CLI-CD-1..CLI-CD-6).

## 9. Success and acceptance criteria

- Every PRD acceptance criterion (CLI-US-1..CLI-US-23, §5.1) has a named test (task list).
- Quality gate per AGENTS.md (`make gate`) plus `make contract-check`, `make generate-check` and
  `make cover` (≥ 90% for `internal/...`) green on every commit.
- CLI-US-7 completeness test strict (no pending commands) at M6 exit.
- Contract replay test: for each of the 180 operations, the request the CLI builds validates against
  the snapshot, `--output json` of the documented example response equals the example, and a `403`
  response exits 254 naming the `operationId` and `requiredPermission`.
- Release dry run (`goreleaser release --snapshot`) green in CI; signed release verified as PRD §5.1
  "Distribution" when the signing secrets exist (OQ-S1).

## 10. Risks and mitigation

| Risk | Mitigation |
|---|---|
| `-restapi` changes `openapi.yaml` while its own implementation runs | Snapshot + checksum; a re-sync task opens every milestone (Mn-0); final re-pin to the `-restapi` release tag after it merges to `main` (M6-X2) |
| oapi-codegen's OpenAPI 3.1 support is partial | Spiked green; types are used only to build request bodies; output never round-trips through types (D-C5); fallback recorded (D-C3) |
| 180 commands overwhelm review | Generic executor; per-service overrides are small; contract replay test covers all operations at once |
| Parallel streams collide | Explicit `Owns` per stream; spec files and `go.mod` owned by the orchestrator; per-file registration (tasks.md §"Parallel work rules") |
| Workload principals cannot read work-order repos to learn the version for `If-Match` | CLI-CD-1 requested; fallback behaviour defined (D-C17) |
| Signing certificates and the Homebrew tap do not exist yet | Release pipeline built and dry-run tested; signing steps gated on secrets (OQ-S1) |

## 11. Timeline and milestones

PRD §7 milestones M0–M6, each a sequence: `Mn-0` contract re-sync → parallel streams → `Mn-X` exit
(gate, docs, milestone counter). Operation-to-milestone map: `tasks.md` Appendix A.

## 12. References

- PRD: `specs/261121-cli-feature-complete/enterprise-agent-factory-service-cli-PRD.md`
- Root PRD: `enterprise-agent-factory-service/enterprise-agent-factory-service-initial-PRD.md` (US-1..US-26, IF-2, D23–D57)
- API contract: root `docs/contracts/api-contract.md` v1.2.0
- OpenAPI: `-restapi/api/openapi.yaml` (branch `feat/restapi-feature-complete`, commit `081de14`)
- `-restapi` contract deltas: `-restapi/specs/261121-restapi-feature-complete/implementation-notes.md` (CD-1..CD-10)
- Solution design §11 (superseded by the contract where they differ)
