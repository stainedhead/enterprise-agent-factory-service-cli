# Tasks: `afs` CLI, feature complete

| | |
|---|---|
| **Feature** | cli-feature-complete |
| **Date** | 2026-11-21 |
| **Status** | Feature complete (M0-M6 implemented); M6-X2 and M6-X3 open |

## Progress summary

99 / 101 tasks complete. M0 28/28 · M1 10/10 · M2 14/14 · M3 10/10 · M4 17/17 · M5 11/11 · M6 9/11 (M6-X2 blocked on OQ-S3; M6-X3 hand-offs listed in implementation-notes for J).

## Parallel work rules

1. **Milestones are sequential.** Inside a milestone: `Mn-0` (re-sync, sequential) → streams in
   parallel → `Mn-X` (exit, sequential). A stream starts when its `Depends` are merged into
   `feat/cli-feature-complete`.
2. **One stream, one worktree, one branch:** `git worktree add ../afs-cli-<stream> -b
   feat/cli-feature-complete--<stream>` from the integration branch. Tasks inside a stream run in order.
3. **Owns is exclusive.** A stream edits only the paths on its `Owns` line. Owning `x.go` includes
   `x_test.go` and golden files under `testdata/x/` of the same package, so two streams in one package
   never share a test or golden file. Anything else is a request to the orchestrator.
4. **Never edited by streams:** `specs/**` (status.md, tasks.md, implementation-notes.md, the dashboard),
   `Makefile`, `.github/workflows/**`, `api/**`, `internal/spec/**`, `internal/client/gen/**`,
   `internal/cli/milestones.go`, `internal/cli/services/services.go`, `README.md`, `docs/**`. These
   belong to `-K` and `-X` tasks and the orchestrator. `go.mod`/`go.sum`: only modules from spec.md
   §7.3; conflicts are resolved by the orchestrator (D-C28).
5. **Report, don't write spec files.** When a task is green, the worker reports: task id, tests added
   (names), decisions, edge cases, deviations. The orchestrator updates status.md, tasks.md and
   implementation-notes.md on the integration branch after every task (AGENTS.md: mandatory).
6. **TDD per task:** Red (the listed tests fail for the right reason) → Green → Refactor; `make gate`
   green in the worktree before reporting. Every test that proves a criterion carries a comment naming
   it, e.g. `// CLI-US-6.1`.
7. **No real network:** `httptest.Server` for `afs-api` and IdPs; injected clock, dialer, stdin, TTY,
   `HOME`. Story criteria are cited as `CLI-US-n.k` = story n, bullet k (PRD §4), `5.1-<area>` = PRD §5.1.
8. **Service registration** never edits a shared file: each owned file adds its registrations with
   `func init() { registrations = append(registrations, registerX) }`, and the package's `register.go`
   (written once by M0-K5) applies them in file-name order. Several files may register parts of one
   operation (D-C32).

Estimates: S ≤ 2 h, M ≤ 4 h, L ≤ 8 h of focused work.

---

## M0 — Foundation (2 operations; 28 tasks)

| Stream | Owns | Tasks |
|---|---|---|
| M0-K kernel (sequential, first) | `api/`, `tools/afsgen/`, `internal/spec/`, `internal/client/gen/`, `internal/cli/{cli,registry,globals,parse,env,milestones}.go`, `internal/cli/services/`, `internal/cli/<service>/doc.go` + empty `register.go` for every service package, `internal/buildinfo/`, `cmd/afs/`, `Makefile`, `.github/workflows/ci.yml`, `go.mod` | K1–K5 |
| M0-T transport | `internal/client/` (except `gen/`), `internal/redact/` | T1–T4 |
| M0-C config and credentials | `internal/config/`, `internal/auth/` (except `workload.go`) | C1–C2 |
| M0-S SSO and configure | `internal/sso/`, `internal/cli/sso/`, `internal/cli/configure/` | S1–S4 |
| M0-O output | `internal/output/` | O1–O4 |
| M0-E executor | `internal/cli/exec/`, `internal/cli/safeupdate/`, `internal/cli/wait/`, `internal/cli/help/`, `internal/cli/replay_test.go`, `internal/cli/testdata/replay/` | E1–E7 |
| M0-X exit (sequential, last) | integration; `README.md`, `docs/`, `user-docs/`, `tools/docgen/`, `Makefile`, CI | X1–X2 |

T, C, O and E start after K5 and run in parallel; S starts after C1. E builds against the K4
interfaces with fakes and does not wait for T, C or O.

### M0-K1 — Snapshot, provenance and checksum gate
- **Owns:** `api/openapi.yaml`, `api/SOURCE`, `api/embed.go`, `api/snapshot_test.go`, `Makefile` (`contract-check`, `contract-sync`, `cover` threshold, `gate` adds `contract-check`), `.github/workflows/ci.yml`
- **Depends:** —
- **Covers:** PRD §6 (pinned contract), D-C1, ADR-002
- **Red:** `TestSnapshotChecksumMatchesSource` (sha256 of `openapi.yaml` = `SOURCE.sha256`); `TestSnapshotIsOpenAPI31` (`openapi: 3.1.0`, `info.version` parses as semver); a tampered copy in `testdata/` fails the check
- **Done when:** snapshot copied verbatim from `-restapi` `feat/restapi-feature-complete` head (record commit; `081de14` at spec time); `make contract-check` in `gate` and CI; `make contract-sync RESTAPI_REF=… | RESTAPI_DIR=…` works against a local clone; `make cover` fails below 90% for `internal/...` excluding generated files
- **Est:** M

### M0-K2 — `afsgen`: operation descriptors
- **Owns:** `tools/afsgen/`, `internal/spec/`, `Makefile` (`generate`, `generate-check`), CI step
- **Depends:** K1
- **Covers:** CLI-1, D-C2, data-dictionary §1, research RQ-3..RQ-8
- **Red:** golden test of `afsgen` on `tools/afsgen/testdata/mini.yaml` (path/query/header params, `$ref` params, `IfMatch`/`IfMatchPut`, body fields, nullable unions, `rows` vs `items`, CSV, 202/204); `internal/spec` tests on the real snapshot: 182 operations, 180 with `Service`, 52 `Paginated`, 8 `CSV`, 59 `IfMatch != None` (required on every `PATCH`/`DELETE`, optional on every `PUT`), 44 `Idempotent`; every `Command` matches `^[a-z]+(-[a-z]+)*$`; `ContractVersion == "1.1.0"`; `SnapshotSHA256` equals `api/SOURCE`
- **Done when:** `make generate` writes deterministic, gofmt'd output sorted by `operationId`; `make generate-check` fails on a hand edit
- **Est:** L

### M0-K3 — Generated types (oapi-codegen)
- **Owns:** `internal/client/gen/`, `Makefile` (`generate` runs oapi-codegen `v2.8.2`), `go.mod` (`oapi-codegen/runtime v1.9.1`)
- **Depends:** K2
- **Covers:** PRD §6 oapi-codegen, D-C3
- **Red:** `TestCreateTeamRequestRoundTrip`, `TestPutBudgetMoneyIsString` (marshal → `"amountUsd":"100.50"`)
- **Done when:** `generate-check` covers `types.gen.go`; lint excludes generated files only
- **Est:** S

### M0-K4 — Registry, parser, global options, exit codes, interfaces
- **Owns:** `internal/cli/{cli,registry,globals,parse,env}.go` (+ tests), `cmd/afs/main.go`, `internal/buildinfo/`
- **Depends:** K2
- **Covers:** CLI-1, CLI-2 (parsing of all 16 options), CLI-8 (codes), CLI-US-7, CLI-US-5.1, 5.1-global-options (`--version`), D-C13, D-C23
- **Red:** `TestCommandTableMatchesSnapshot` (CLI-US-7: every operation except `GetHealth`/`GetReadiness` has a command; every command has an operation except `configure`, `configure sso`, `sso login`, `sso logout`, `completion`, `help`; names equal `x-afs-cli`); `TestGlobalOptionsAnywhere`; `TestUnknownServiceAndOptionExit252`; `TestVersionString` (`afs-cli/1.2.3 contract/1.1 go1.27.4 linux/amd64`, no call); `TestSigintExit130` (`cmd/afs` with a fake signal)
- **Done when:** the §11 `Services` list is gone; `Env` (stdin, stdout, stderr, env map, home, TTY flags, clock, `http.RoundTripper`) is injected everywhere; interfaces `Transport`, `Renderer`, `CredentialSource` defined for the parallel streams
- **Est:** L

### M0-K5 — Service packages, overrides and milestone map
- **Owns:** `internal/cli/services/`, `internal/cli/{teams,members,apikeys,agents,secrets,budgets,reports,models,audit,admin,pools,swarmtemplates,workorders,sso,configure,completion}/{doc,register}.go`, `internal/cli/milestones.go`, `internal/cli/override.go`
- **Depends:** K4
- **Covers:** D-C4, D-C32, CLI-14 (example rule scaffold)
- **Red:** `TestEveryOperationHasMilestone` (Appendix A, all 180); `TestDoneUpToCurrentMilestone` (every operation at or below `CurrentMilestone` has ≥ 1 example and table columns; `CurrentMilestone = M0`, so `GetMe` and `ListSignInProviders` are listed as pending until S4 — the test is red until S4 lands, by design, and K5 ships it skipped with `t.Skip("enabled by M0-X1")`); `TestOverrideConflicts` (two `Columns` for one operation → error)
- **Done when:** every service package exists and is imported once from `services.go`; each package's `register.go` applies its `registrations` slice (rule 8); registration merge rules implemented
- **Est:** M

### M0-T1 — Request building, response capture, API errors
- **Owns:** `internal/client/{client,request,errors,version}.go` + tests
- **Depends:** K5
- **Covers:** CLI-8 (decode), NFR observability (`User-Agent`), NFR compatibility, D-C11, D-C13
- **Red:** `TestPathEscapingPerSegment` (`host/owner/repo` values), `TestUserAgent` (`afs-cli/<v> (<os>/<arch>)`), `TestAcceptHeaders` (json, `text/csv`, `application/yaml`), `TestAPIErrorDecode` (envelope → `Code`, `Message`, `RequestID`, `Details`; non-JSON 5xx → `InternalError` with status), `TestContractVersionNoticeOncePerProcess` (newer minor → one stderr notice; same/older → none; different major → notice); `TestUnexpectedMediaType255`
- **Est:** M

### M0-T2 — Retries and idempotency keys
- **Owns:** `internal/client/retry.go`, `internal/client/idempotency.go` + tests
- **Depends:** T1
- **Covers:** CLI-11, 5.1-retries (all four bullets), D-C24
- **Red:** `TestIdempotencyInProgressTwiceThenSuccess` (3 attempts, same key, exit 0); `…ThreeTimes` (254); `TestIdempotencyMismatchNotRetried`; `Test429RetryAfter2WaitsAtLeast2s` (fake clock); `TestPatchDroppedAfterSendNotRetried` (255) vs `TestPatchConnectFailureRetried`; `Test503RetriedOnGetAndIdempotentPostOnly`; `TestULIDPerInvocationOrFlag` (`--idempotency-key` honoured); `TestReplayHeaderSurfaced` (`Idempotent-Replayed: true` exposed to the executor)
- **Est:** M

### M0-T3 — Timeouts and cancellation
- **Owns:** `internal/client/transport.go` + tests
- **Depends:** T1
- **Covers:** NFR reliability, CLI-US-5.3 (Ctrl-C), D-C21
- **Red:** `TestConnectTimeoutDefault30s` / `--cli-connect-timeout 0` = none (injected dialer); `TestReadTimeoutHeadersAndBody`; `TestContextCancelAbortsInFlight` (→ `context.Canceled`, mapped to 130 by the executor)
- **Est:** S

### M0-T4 — Debug log and redaction
- **Owns:** `internal/redact/`, `internal/client/debug.go` + tests
- **Depends:** T1
- **Covers:** NFR observability, CLI-US-3.4, 5.1-credentials (`--debug` names the source, never the key)
- **Red:** golden `debug.golden` (`<RFC 3339> DEBUG <component> <message>`: method, URL, status, duration, `X-Request-Id`, credential source, attempt); `TestRedactsHeaders` (`Authorization`, cookies); `TestRedactsBodyFields` (data-dictionary §8, nested and in arrays); `TestNoQueryParamCarriesSecret` (snapshot scan)
- **Est:** M

### M0-C1 — Config and credentials files, precedence
- **Owns:** `internal/config/` + tests
- **Depends:** K5
- **Covers:** CLI-3 (files, keys, precedence), 5.1-configuration bullets 1–3, D-C10, D-C15, D-C33, D-C35
- **Red:** `TestPrecedenceEndpointURL` (option > `AFS_ENDPOINT_URL` > profile > default); `TestAFSProfileSelectsProfile`; `TestTeamIDDefaultsPathOnly`; `TestDefaultOutputJSON`; `TestWriteKeepsOtherProfilesAnd0600`; `TestConfigAndCredentialsFileEnvOverrides`; `TestHTTPEndpointOnlyLoopback`; `TestUnknownProfile255`; `TestInvalidOutputFromEnvOrProfile255`
- **Est:** M

### M0-C2 — Credential chain and file modes
- **Owns:** `internal/auth/` (except `workload.go`) + tests
- **Depends:** C1, K4 (`CredentialSource`)
- **Covers:** CLI-5 (env key, credentials file, SSO slot, workload slot reserved), 5.1-credentials, D-C29
- **Red:** `TestEnvKeyWinsOverProfileSSO` (CLI-5; key sent even with `--profile`); `TestCredentialsFile0644Refused` (255, `chmod 600 ~/.afs/credentials`, no call); `TestSourceNamedForDebug`; `TestNoCredentialForPublicOperation` (`ListSignInProviders` sends none)
- **Est:** M

### M0-S1 — SSO token cache and expiry rules
- **Owns:** `internal/sso/cache.go` + tests
- **Depends:** C1
- **Covers:** CLI-4 (cache 0600, expiry), CLI-US-1.3, D41
- **Red:** `TestCacheWrites0600`; `TestCacheWiderThan0600Refused` (255); `TestExpiredNoRefreshExit255NoCall` (`Run afs sso login --profile <name>`); `TestServer401Exit254WithHint` (executor-level test with a fake transport)
- **Est:** M

### M0-S2 — Device authorization grant (`sso login`)
- **Owns:** `internal/sso/login.go`, `internal/cli/sso/login.go` + tests
- **Depends:** S1, K5
- **Covers:** CLI-4, CLI-US-1.1, D48
- **Red:** fake IdP (discovery, device authorization, token): `TestLoginPrintsCodeAndURLThenCaches0600`; `TestLoginHonoursSlowDownAndPending`; `TestLoginUsesCliClientIdFromListSignInProviders`; `TestLoginExpiredDeviceCode` (255)
- **Est:** L

### M0-S3 — Silent refresh
- **Owns:** `internal/sso/refresh.go` + tests
- **Depends:** S2
- **Covers:** CLI-4, CLI-US-1.2
- **Red:** `TestExpiredAccessTokenRefreshedSilently` (no prompt, command succeeds, cache updated 0600); `TestRefreshRejectedFallsBackToHint` (255)
- **Est:** S

### M0-S4 — `sso logout`, `get-caller-identity`, `list-providers`, `configure`, `configure sso`
- **Owns:** `internal/cli/sso/{logout,register}.go`, `internal/cli/configure/` + tests
- **Depends:** S3, C2
- **Covers:** CLI-3 (`configure`, `configure sso`), CLI-4, CLI-US-1.4, CLI-US-1.5, 5.1-configuration bullets 4–5; operations `GetMe`, `ListSignInProviders`
- **Red:** `TestGetCallerIdentityJSONUnchanged`; `TestLogoutDeletesCacheThenNextCommandExit255`; `TestConfigureWritesProfile0600` (injected stdin); `TestConfigureSSOListsProvidersAndWritesKeys` (`sso_provider_id`, `sso_issuer`, `sso_client_id`); overrides with examples and columns for both operations
- **Est:** M

### M0-O1 — `json` and `yaml` renderers
- **Owns:** `internal/output/{json,yaml,tree}.go` + tests
- **Depends:** K5
- **Covers:** CLI-6, CLI-US-4.1, D-C5
- **Red:** `TestJSONPreservesOrderNullsUnknownFields` (raw body re-indented byte-stably); `TestYAMLFromOrderedTree`; `TestNumbersKeepPrecision` (`UseNumber`)
- **Est:** S

### M0-O2 — `--query` (JMESPath)
- **Owns:** `internal/output/query.go` + tests
- **Depends:** O1
- **Covers:** CLI-6, CLI-US-4.6
- **Red:** `TestQueryBeforeRenderJSONYAMLTableText` (``items[?state==`active`].slug``); `TestInvalidQueryExit252`
- **Est:** S

### M0-O3 — `table` and `text` renderers, colour
- **Owns:** `internal/output/{table,text,color,golden}.go` + tests, `internal/output/testdata/`
- **Depends:** O1
- **Covers:** CLI-6, NFR accessibility, 5.1-global-options (`--color on|off|auto`)
- **Red:** golden table/text for a fixture; `TestTableFitsIn80Columns` (truncation marker); `TestColorAutoOffWhenNotTTY`; `TestNoColorEnv`; `TestGenericColumnsFallback` (scalar top-level fields when an operation has no `Columns`); `TestEmptyListRendering` (table header only, text nothing); golden helper with `-update`
- **Est:** M

### M0-O4 — `csv` passthrough and pager
- **Owns:** `internal/output/{csv,pager}.go` + tests
- **Depends:** O1
- **Covers:** CLI-6 (csv body unchanged, one header row), 5.1-global-options (`--no-cli-pager`), D-C22
- **Red:** `TestCSVPagesKeepOneHeader`; `TestCSVBytesUnchanged`; `TestCSVStreamsPageByPage`; `TestPagerOnlyOnTTY`; `TestNoCLIPager`; `TestEmptyAFSPagerDisables`
- **Est:** S

### M0-E1 — Option derivation, input JSON, skeletons
- **Owns:** `internal/cli/exec/{options,input,skeleton,build}.go` + tests
- **Depends:** K5
- **Covers:** CLI-1 (ids not names), CLI-2 (`--cli-input-json`, `--generate-cli-skeleton`), CLI-US-2.1, CLI-US-5.1, CLI-US-12.2 (`--replicas` unknown → 252), data-dictionary §2 rules, D-C6, D-C10
- **Red:** `TestDerivedOptionsPerKind` (string, int, Money string, bool `--x/--no-x`, list, JSON object, `file://`, any-type `--value`); `TestEnumAndRangeExit252BeforeCall`; `TestMissingRequiredExit252WithUsage`; `TestPathTeamIDFromProfile`; `TestCLIInputJSONBodyUnchangedAndOptionsWin`; `TestInputVersionBecomesIfMatch`; `TestSkeletonEveryWritableFieldNoCall` (all 180 operations)
- **Est:** L

### M0-E2 — Pagination (JSON and CSV)
- **Owns:** `internal/cli/exec/paginate.go` + tests
- **Depends:** E1, K4 (`Transport`)
- **Covers:** CLI-6 (csv rules), CLI-7, 5.1-pagination (all bullets), CLI-US-4.3, CLI-US-4.4, CLI-US-4.5, D-C7
- **Red:** 3-page fake: `TestAutoPaginateMergesItems`; `TestMaxItems10PrintsNextToken`; `TestStartingTokenResumesAfterItem10`; `TestPageSize20SendsMaxResults20`; `TestPageSize101Exit252`; `TestNoPaginateReturnsPage1WithNextToken`; `TestRowsKeyForReports`; `TestCSVFollowsXAfsNextTokenUntilAbsent`; `TestCSVOnNonCSVOperationExit252` (`agents list-agents`); `TestCSVWithQueryExit252`; `TestCSVWithMaxItemsExit252`
- **Est:** L

### M0-E3 — Safe updates (`If-Match`)
- **Owns:** `internal/cli/safeupdate/` + tests
- **Depends:** E1
- **Covers:** CLI-10, CLI-US-6 (generic forms of all three bullets), D-C8
- **Red:** `TestReaderTableCoversAllIfMatchOperations` (59; read operations exist; match fields exist); `TestIfMatchFlagSentAsGiven` (`--if-match 7` → `If-Match: "7"`); `TestReadThenWriteNotice`; `TestPut404SendsNoIfMatch`; `TestNoVersionResourceSendsNone` (`PutTeamModels`); `Test412PrintsCurrentVersion` (from `ETag`, else a read) exit 254; `Test428Exit254NotRetried`; `TestListReaderPaginatesToFindItem`
- **Est:** L

### M0-E4 — `--wait`
- **Owns:** `internal/cli/wait/` + tests
- **Depends:** E1
- **Covers:** CLI-12, CLI-US-2.3 (generic), CLI-US-8.2 (`--wait-timeout` → 255; pool rule wired in M5-A3), D-C9
- **Red:** `TestWaitReadyExit0PrintsFinalResource`; `TestWaitFailedOrDegradedExit254WithConditions`; `TestWaitIgnoresStaleReadyByGeneration`; `TestWaitTimeoutExit255`; `TestDeleteWaitsFor404`; `TestWaitHonoursRetryAfter`; `TestCtrlCDuringWaitExit130MutationStands`; `TestWaitTableOpsExist` (snapshot)
- **Est:** M

### M0-E5 — Error rendering and exit codes
- **Owns:** `internal/cli/exec/{errors,exit,result}.go` + tests
- **Depends:** E1
- **Covers:** CLI-8, CLI-US-5.2, CLI-US-5.3, D-C26
- **Red:** `TestAWSStyleErrorAndRequestID` (`An error occurred (<code>) when calling the <operationId> operation: <message>` then `Request ID: <id>`); `TestAccessDeniedIncludesRequiredPermission`; `TestNetworkErrorExit255`; `TestUnauthenticatedHint`; `TestDelete202PrintsBody204PrintsNothing`; `TestReplayPrintsStoredResponse`
- **Est:** M

### M0-E6 — Help
- **Owns:** `internal/cli/help/` + tests, `internal/cli/help/testdata/`
- **Depends:** E1
- **Covers:** CLI-14, 5.1-help (all bullets)
- **Red:** golden `afs help`, `afs sso help`, `afs sso get-caller-identity help`; `TestServiceHelpListsExactlyItsOperations` (against the CLI-US-7 table); `TestOperationHelpHasExample` (enforced through `TestDoneUpToCurrentMilestone`); help shows options, types, required marks, enum values, the permission (`x-afs-permission`) and the example
- **Est:** M

### M0-E7 — Contract replay test (all operations)
- **Owns:** `internal/cli/replay_test.go`, `internal/cli/testdata/replay/`, `go.mod` (test-only `libopenapi-validator v0.16.5`)
- **Depends:** E2, E3, E5
- **Covers:** CLI-US-4.1 (`--output json` equals the API's fields, per operation), NFR quality (`AccessDenied` test per operation), contract §8.4 (replay the OpenAPI examples)
- **Red:** for each of the 180 operations, args built from the snapshot's parameter and request examples (`--cli-input-json`): (a) the request validates against the snapshot (`libopenapi-validator`; `tag.*` query keys allow-listed, CLI-CD-4); (b) a fake answering the documented 2xx example gives `--output json` semantically equal to the example; (c) a fake answering `403 AccessDenied` gives exit 254 and stderr naming the `operationId` and `requiredPermission`
- **Est:** L

### M0-X1 — M0 integration and gate
- **Owns:** integration of K/T/C/S/O/E branches; `internal/cli/milestones.go` (`CurrentMilestone = M0`, un-skip `TestDoneUpToCurrentMilestone`); `internal/cli/globals_test.go` §5.1 table (one case per global option: `--profile p2`, `--endpoint-url`, `--no-cli-pager`, `--color`, `--generate-cli-skeleton`, `--version`); `internal/cli/startup_bench_test.go`; CI workflow
- **Depends:** all M0 streams
- **Covers:** 5.1-global-options, NFR startup (< 100 ms to first request against a local fake, reported in CI), NFR quality (≥ 90%)
- **Done when:** `make gate`, `contract-check`, `generate-check`, `cover` green; streams merged without conflicts
- **Est:** M

### M0-X2 — Docs and generated reference
- **Owns:** `tools/docgen/`, `Makefile` (`docs`), `user-docs/{README,getting-started,configuration,sso,credentials,output,pagination,errors-and-exit-codes,troubleshooting}.md`, `user-docs/reference/`, `README.md`, `docs/technical-details.md`, `docs/architectural-decision-record.md` (ADR-002 accepted, ADR-003..ADR-005), `AGENTS.md` layout table
- **Depends:** X1
- **Covers:** CLAUDE.md docs rules (README current, user-docs per feature), D-C30
- **Red:** `TestDocgenMatchesHelp` (`user-docs/reference/*.md` regenerated equals committed)
- **Est:** M

---

## M1 — Teams, members, team audit, settings, identity providers (24 operations; 10 tasks)

| Stream | Owns | Tasks |
|---|---|---|
| M1-A teams core | `internal/cli/teams/{teams,lifecycle,models,revisions}.go` + tests + `testdata/` | A1–A4 |
| M1-B members | `internal/cli/members/` | B1 |
| M1-C team audit | `internal/cli/audit/team.go` + tests | C1 |
| M1-D settings and identity providers | `internal/cli/admin/{settings,identity_providers}.go` + tests | D1–D2 |

Every service task below also delivers, for each of its operations: ≥ 1 example, table columns,
golden `table`/`text`/`help`, and an `AccessDenied` (254) case in the package's table test. This is
called **"command done"** and is not repeated per task.

### M1-0 — Contract re-sync
- **Owns:** `api/`, `internal/spec/`, `internal/client/gen/` (orchestrator)
- **Depends:** M0-X2
- **Covers:** D-C1 sync policy
- **Done when:** snapshot re-synced to the `-restapi` branch head (or release tag if merged); if the sha256 changed: regenerate, full suite green, operation diff recorded in implementation-notes.md. An operation added by `-restapi` is mapped in `internal/cli/milestones.go` and Appendix A to the current or a later milestone (and given a stream) before streams start; a renamed or removed one updates its override in the owning stream's first task
- **Est:** S

### M1-A1 — `list-teams`, `describe-team`, `get-team-status`
- **Depends:** M1-0 · **Covers:** CLI-1, CLI-US-4.6 (`--query` on `list-teams`), command done · **Red:** request/rendering table tests; `--state`/`--kind` filters · **Est:** S

### M1-A2 — `create-team`
- **Depends:** A1 · **Covers:** CLI-US-2.1 (skeleton with every `CreateTeam` field; `file://team.json` sent unchanged), CLI-US-2.2 (missing tag → `ValidationError`, 254), CLI-US-2.3 (`--wait` → `Ready` prints final team; `Failed` → 254 with conditions) · **Red:** `TestCreateTeamSkeleton`, `TestCreateTeamInputFileUnchanged`, `TestCreateTeamMissingTag254`, `TestCreateTeamWaitReady`, `TestCreateTeamWaitFailed` · **Est:** M

### M1-A3 — `update-team`, `delete-team`, `suspend-team`, `resume-team`
- **Depends:** A1 · **Covers:** CLI-US-23.1 (suspend; following write → `Conflict` 254; resume lifts), CLI-US-23.2 (team-admin `delete-team` → `AccessDenied` 254), CLI-10 on `update-team`/`delete-team`, CLI-12 `--wait` · **Red:** `TestSuspendThenUpdateAgentConflict254` (fake returns `team_suspended`), `TestDeleteTeamAccessDenied`, `TestDeleteTeamWaitFor404` · **Est:** M

### M1-A4 — `list-team-models`, `put-team-models`, `list-revisions`, `describe-revision`
- **Depends:** A1 · **Covers:** CLI-1, D-C8 (`put-team-models` sends no `If-Match`) · **Red:** `TestPutTeamModelsNoIfMatch`, `--model-ids` list option · **Est:** S

### M1-B1 — `members`: `list-members`, `add-member`, `update-member`, `remove-member`
- **Depends:** M1-0 · **Covers:** CLI-1, CLI-10 (update/remove read via `ListMembers`), command done · **Red:** `TestRemoveMemberReadsVersionFromList` · **Est:** M

### M1-C1 — `audit list-audit-events`
- **Depends:** M1-0 · **Covers:** CLI-US-4.4 (3 CSV pages → one header, every row), CLI-6, filters `--from --to --actor-id --action --outcome --resource-kind --resource-id` · **Red:** `TestAuditCSV3Pages` · **Est:** S

### M1-D1 — `admin get-settings`, `admin put-setting`
- **Depends:** M1-0 · **Covers:** CLI-US-15.3 (`--value '["https://console.example.com"]'` sent as JSON; `get-settings --output json` prints every setting), CLI-US-21.1 (out-of-range → `ValidationError` naming the range, 254), CLI-US-21.2 (team-admin → `AccessDenied` 254) · **Red:** `TestPutSettingJSONValue`, `TestPutSettingPlainStringValue`, `TestPutSettingValidationError`, `TestPutSettingAccessDenied` · **Est:** S

### M1-D2 — identity providers (5 operations)
- **Depends:** M1-0 · **Covers:** CLI-1, `--provider-id` rename, `update`/`delete` `409` in use printed (254), no secret in output (`clientSecretArn` only) · **Red:** `TestDeleteIdentityProviderInUse409` · **Est:** M

### M1-X1 — M1 exit
- **Owns:** integration; `CurrentMilestone = M1`; `make docs`; `user-docs/` examples for teams and members; README
- **Depends:** all M1 streams · **Covers:** CLI-14 for M1 operations (`TestDoneUpToCurrentMilestone`), gate · **Est:** S

---

## M2 — Agents, secrets, team config and external systems, models, admin catalogs (45 operations; 14 tasks)

| Stream | Owns | Tasks |
|---|---|---|
| M2-A agents | `internal/cli/agents/{agents,lifecycle,snapshots,catalog,config_values,external_systems}.go` + tests | A1–A5 |
| M2-B secrets | `internal/secretinput/`, `internal/cli/secrets/` | B1–B3 |
| M2-C team config and external systems | `internal/cli/teams/{config_values,external_systems}.go` + tests | C1–C2 |
| M2-D admin catalogs and models | `internal/cli/admin/{catalog,storage_classes,agent_templates}.go`, `internal/cli/models/` | D1–D2 |

### M2-0 — Contract re-sync
- As M1-0. **Depends:** M1-X1 · **Est:** S

### M2-A1 — `list-agents`, `describe-agent`, `get-agent-status`
- **Depends:** M2-0 · **Covers:** CLI-1, CLI-US-4.5 (`--output csv` on `list-agents` → 252, product test), command done · **Est:** S

### M2-A2 — `create-agent`
- **Depends:** A1 · **Covers:** CLI-US-12.1 (`--volume '{…}' --wait` → 0 on `Ready`), CLI-US-12.2 (`--replicas` → 252), CLI-US-12.3 (image by tag → `ValidationError` 254), CLI-US-12.4 (two teams same `--name`; later commands by `--agent-id`) · **Red:** `TestCreateAgentVolumeJSONAndWait`, `TestCreateAgentReplicasUnknown252`, `TestCreateAgentImageTag254`, `TestCreateAgentSameNameTwoTeams` · **Est:** M

### M2-A3 — `update-agent`, `delete-agent`, `pause-agent`, `resume-agent`, `restart-agent`
- **Depends:** A1 · **Covers:** CLI-US-6.1 (`--if-match 7` → `If-Match: "7"`; at version 8 → `The resource changed (version 8). Re-run with --if-match 8 to apply your change to the latest version.`, 254), CLI-US-6.2 (read, send, notice), CLI-US-6.3 (`delete-agent` sends `If-Match`; `428` → 254), CLI-US-23.1 (update after suspend → `Conflict`) · **Red:** `TestUpdateAgentIfMatch7`, `TestUpdateAgent412Message`, `TestUpdateAgentReadThenWrite`, `TestDeleteAgentIfMatchAnd428` · **Est:** M

### M2-A4 — `list-snapshots`, `create-snapshot`, `list-storage-classes`, `list-agent-templates`
- **Depends:** A1 · **Covers:** CLI-1, command done · **Est:** S

### M2-A5 — agent config values (3) and external systems (2)
- **Depends:** A1 · **Covers:** CLI-1, `put-config-value --source literal|namespace --value`, CLI-10 via list readers · **Est:** M

### M2-B1 — `internal/secretinput`
- **Depends:** M2-0 · **Covers:** CLI-9 (file, fileb, stdin, no-echo prompt; plain argument refused), D-C16 · **Red:** `TestFileAndFilebVerbatim`, `TestStdinDash`, `TestPromptNoEcho` (fake terminal), `TestPlaintextRefused252`, `TestNoTTYNoValue252`, `TestNonUTF8OrTooLong252` · **Est:** M

### M2-B2 — `create-secret`, `put-secret-value`
- **Depends:** B1 · **Covers:** CLI-US-3.1, CLI-US-3.2, CLI-US-3.3 (inline `--cli-input-json` with `value` → 252 before any call; `file://` accepted), CLI-US-3.4 (`--debug` shows `"value": "[redacted]"` and `Authorization: [redacted]`), PRD NFR security (value never in output) · **Red:** `TestCreateSecretFromFile`, `TestCreateSecretStdin`, `TestCreateSecretPrompt`, `TestSecretValuePlaintext252`, `TestInlineInputJSONWithValue252`, `TestInputJSONFileAccepted`, `TestDebugRedactsValue`, `TestPutSecretValue204PrintsNothing` · **Est:** M

### M2-B3 — `list-secrets`, `describe-secret`, `update-secret`, `rotate-secret`, `delete-secret`
- **Depends:** B1 · **Covers:** CLI-1, `delete-secret` `409` with `details.dependents` printed (254), command done · **Est:** S

### M2-C1 — team config values (3)
- **Depends:** M2-0 · **Covers:** CLI-1, CLI-10 (`put-config-value` reads `ListNamespaceConfigValues`, `404` → create without `If-Match`) · **Red:** `TestPutNamespaceConfigValueCreateNoIfMatch` · **Est:** S

### M2-C2 — team external systems (5)
- **Depends:** M2-0 · **Covers:** CLI-1, `declare-external-system` refused when `mayAddExternalSystems` is off (API's `AccessDenied`, 254) · **Est:** M

### M2-D1 — `models list-models`, admin catalog systems (5)
- **Depends:** M2-0 · **Covers:** CLI-1, `delete-external-system` built-in → `409` (254) · **Est:** M

### M2-D2 — admin storage classes (3), admin agent templates (3)
- **Depends:** M2-0 · **Covers:** CLI-1, `--allowed-digests` list, `delete-*` in use → `409` · **Est:** M

### M2-X1 — M2 exit
- As M1-X1 for M2; `user-docs/` secrets guide (secret input rules). **Depends:** all M2 streams · **Est:** S

---

## M3 — Budgets, reports, team notification channels, admin models and prices (18 operations; 10 tasks)

| Stream | Owns | Tasks |
|---|---|---|
| M3-A budgets | `internal/cli/budgets/` | A1–A2 |
| M3-B reports on demand | `internal/cli/reports/{reports,flags}.go` + tests | B1–B2 |
| M3-C notification channels | `internal/cli/teams/notification_channels.go` + tests | C1–C2 |
| M3-D admin models and prices | `internal/cli/admin/{models,prices}.go` + tests | D1–D2 |

### M3-0 — Contract re-sync
- As M1-0. **Depends:** M2-X1 · **Est:** S

### M3-A1 — `get-budget`, `put-budget`
- **Depends:** M3-0 · **Covers:** CLI-US-14.1 (`amountUsd: "100.50"`, `thresholds: [50, 80, 100]`), CLI-US-14.2 (above ceiling → `ValidationError` 254), CLI-US-14.3 (table: amount, spent, reserved, percent used, forecast), CLI-US-6.3 (`put-budget` sends `If-Match`) · **Red:** `TestPutBudgetBody`, `TestPutBudgetAboveCeiling254`, `TestGetBudgetTableGolden`, `TestPutBudgetIfMatch` · **Est:** M

### M3-A2 — `get-usage`, `list-costs`
- **Depends:** M3-0 · **Covers:** CLI-6 csv, CLI-13 (`--from`, `--to`, `--group-by agent|model|day`), RPT-1..9 via CSV · **Red:** `TestGetUsageCSV`, `TestListCostsPaginates` · **Est:** S

### M3-B1 — report flags
- **Depends:** M3-0 · **Covers:** CLI-13 (`--range 24h|7d|30d|90d|mtd`, `--from`/`--to`, `--group-by` incl. `tag:<key>` and `pool`, `--tag key=value` → `tag.key=value`, `--pool-id`, `--team-id`/`--agent-id`/`--model-id`/`--system-id` filters), validation → 252 · **Red:** `TestReportFlagsToQuery`, `TestRangeEnum252`, `TestTagFilterShorthand` · **Est:** S

### M3-B2 — `get-team-report`, `get-agent-report`, `get-service-report`
- **Depends:** B1 · **Covers:** CLI-US-4.2 (`--kind costs --range 30d --group-by tag:project --output csv` sends `Accept: text/csv`, writes CSV byte for byte), CLI-US-4.3 (3-page CSV via `X-Afs-Next-Token`), JSON `rows` pagination · **Red:** `TestTeamReportCSVByteForByte`, `TestReportCSV3Pages`, `TestServiceReportAccessDeniedForTeamAdmin` · **Est:** M

### M3-C1 — notification channels CRUD (5)
- **Depends:** M3-0 · **Covers:** CLI-US-19.2 (`describe-notification-channel` prints configuration and delivery status, never a credential), `--channel-id` rename, `--wait` on create/update/delete, CLI-10 · **Est:** M

### M3-C2 — `test-notification-channel`
- **Depends:** C1 · **Covers:** CLI-US-19.1 (failure → channel status and error, AWS-style message of `TestNotificationChannel`, 254), D-C20 · **Red:** `TestNotificationTestDelivered0`, `TestNotificationTestFailed254`, `TestNotificationTestPending202Exit0` · **Est:** S

### M3-D1 — admin models: `list-models`, `create-model`, `update-model`
- **Depends:** M3-0 · **Covers:** CLI-1, `create-model` without `--max-output-tokens` → 252 (required, D28) · **Est:** S

### M3-D2 — admin prices: `list-model-prices`, `create-model-price`
- **Depends:** M3-0 · **Covers:** CLI-1, Money options as strings · **Est:** S

### M3-X1 — M3 exit
- As M1-X1 for M3; `user-docs/` reports and CSV guide. **Depends:** all M3 streams · **Est:** S

---

## M4 — API keys, role links, identities, tags, log destinations, restore, schedules, admin service items, exports, completion, releases (41 operations; 17 tasks)

| Stream | Owns | Tasks |
|---|---|---|
| M4-A api-keys | `internal/cli/apikeys/` | A1 |
| M4-B agent identity | `internal/cli/agents/{role_links,identities}.go` + tests | B1–B2 |
| M4-C teams delivery and history | `internal/cli/teams/{tags,log_destinations,allowed_accounts,restore}.go` + tests | C1–C3 |
| M4-D schedules and service audit | `internal/cli/reports/schedules.go`, `internal/cli/audit/service.go` + tests | D1–D2 |
| M4-E admin service items | `internal/cli/admin/{service_channels,service_tags,price_proposals,reconciliations,org_accounts,exports}.go` + tests | E1–E4 |
| M4-F completion | `internal/cli/completion/` | F1 |
| M4-G releases | `.goreleaser.yaml`, `.github/workflows/release.yml`, `user-docs/install.md` | G1–G2 |

### M4-0 — Contract re-sync
- As M1-0. **Depends:** M3-X1 · **Est:** S

### M4-A1 — `list-api-keys`, `create-api-key`, `delete-api-key`
- **Depends:** M4-0 · **Covers:** NFR security exception (`secretKey` once to stdout, `Store this key now; it cannot be shown again.` on stderr, redacted in `--debug`), D52 replay notice, D-C25, D50 (viewer → `AccessDenied`) · **Red:** `TestCreateAPIKeyPrintsSecretOnce`, `TestCreateAPIKeyDebugRedacted`, `TestCreateAPIKeyReplayNoSecret` · **Est:** M

### M4-B1 — agent role links (5)
- **Depends:** M4-0 · **Covers:** CLI-US-13.1 (`create-role-link --role-arn --profile-name` prints `pending_trust` and `externalId`), CLI-US-13.2 (`describe-role-link --query trustPolicy --output json` prints the document alone), CLI-US-13.3 (`verify-role-link` `verified`/`failed` with reason; account outside allowed accounts → `ValidationError` 254) · **Red:** `TestCreateRoleLinkPendingTrust`, `TestDescribeRoleLinkTrustPolicyQuery`, `TestVerifyRoleLinkFailedReason`, `TestRoleLinkAccountNotAllowed254` · **Est:** M

### M4-B2 — agent identities (3)
- **Depends:** M4-0 · **Covers:** CLI-1 · **Est:** S

### M4-C1 — `get-team-tags`, `put-team-tags`
- **Depends:** M4-0 · **Covers:** CLI-US-16.1 (`--tags k=v …` sends team-owned keys only), CLI-US-16.2 (table shows every tag with `source`, incl. `system`, `service`, `missing`), CLI-US-16.3 (`afs-`/`aws:` prefix or service key → API refusal 254), CLI-US-6.3 (`put-team-tags` `If-Match`) · **Red:** `TestPutTeamTagsShorthand`, `TestGetTeamTagsTableGolden`, `TestPutTeamTagsReservedPrefix254`, `TestPutTeamTagsIfMatch` · **Est:** M

### M4-C2 — log destinations (5)
- **Depends:** M4-0 · **Covers:** CLI-US-19.2 (`describe-log-destination` prints configuration and last delivery status, no credential), `--destination-id`, `--wait` · **Est:** M

### M4-C3 — `list-allowed-accounts`, `restore-revision`
- **Depends:** M4-0 · **Covers:** CLI-US-20.1 (`--dry-run` sends `dryRun: true`, prints the plan; without it applies), CLI-US-20.2 (viewer → `AccessDenied` 254), `--revision-id` · **Red:** `TestRestoreRevisionDryRun`, `TestRestoreRevisionApply`, `TestRestoreRevisionViewerDenied` · **Est:** S

### M4-D1 — report schedules (4)
- **Depends:** M4-0 · **Covers:** CLI-US-17.1 (`create-report-schedule --name --report costs --cadence weekly --channel-id --range 30d --group-by tag:project` stores the query, prints `nextRunAt`), CLI-US-17.2 (service-level or other team's channel → `ValidationError` 254; viewer → `AccessDenied` 254) · **Red:** `TestCreateReportScheduleQueryFromFlags`, `TestReportScheduleForeignChannel254`, `TestReportScheduleViewerDenied` · **Est:** M

### M4-D2 — `audit list-service-audit-events`
- **Depends:** M4-0 · **Covers:** CLI-6 csv, `admin:audit` denial for team roles · **Est:** S

### M4-E1 — admin service notification channels (4)
- **Depends:** M4-0 · **Covers:** CLI-1, delete named in a setting → `409` · **Est:** S

### M4-E2 — admin service tags (3)
- **Depends:** M4-0 · **Covers:** CLI-1, `--mode fixed|required`, `--allowed-values` list · **Est:** S

### M4-E3 — price proposals (3), cost reconciliations, org accounts
- **Depends:** M4-0 · **Covers:** CLI-US-15.1 (`list-price-proposals --status pending --output table` one row per proposal), CLI-US-15.2 (`approve-price-proposal --proposal-id --note` prints the decided proposal; `reject-price-proposal` likewise; team-admin → `AccessDenied` 254) · **Red:** `TestListPriceProposalsTableGolden`, `TestApprovePriceProposal`, `TestRejectPriceProposal`, `TestApproveAccessDeniedTeamAdmin` · **Est:** M

### M4-E4 — exports and imports (4)
- **Depends:** M4-0 · **Covers:** CLI-US-18.1 (`create-export` prints id; `describe-export` prints `sha256`, `sizeBytes`, `downloadUrl` when complete; `list-exports`), CLI-US-18.2 (`create-import --bundle file:// --account-map file://` sends `?dryRun=true` unless `--no-dry-run`; body `{bundle, accountMap}`, D-C18), CLI-US-18.3 (bundle > 100 MB → 252 before any call) · **Red:** `TestCreateImportDryRunDefault`, `TestCreateImportNoDryRun`, `TestCreateImportBodyWrapsFiles`, `TestCreateImportOver100MB252` (sparse file) · **Est:** M

### M4-F1 — `afs completion bash|zsh|fish|powershell`
- **Depends:** M4-0 · **Covers:** CLI-3 (`completion`, generated from the command table), OQ-C5, D-C27 · **Red:** golden scripts per shell; `TestCompletionListsEveryCommand`; `TestCompletionUnknownShell252`; bash script syntax check (`bash -n`) when bash is present · **Est:** M

### M4-G1 — GoReleaser configuration and snapshot build
- **Depends:** M4-0 · **Covers:** CLI-15 (5 platform archives, `checksums.txt`, `--version` ldflags with the tag), D-C14 · **Red:** CI job `goreleaser check` + `goreleaser release --snapshot --clean --skip=sign,publish,notarize`; a script test that the built binary's `--version` matches `afs-cli/<snapshot version> contract/1.1` · **Est:** M

### M4-G2 — Signing, notarization, Homebrew tap, verification
- **Depends:** H1 · **Covers:** CLI-15, 5.1-distribution (cosign `verify-blob` against the published identity; `spctl --assess` on macOS; Authenticode valid; `brew install afs` from the tap matches the tag) · **Done when:** release workflow on `v*` tags: cosign keyless `sign-blob checksums.txt`, quill notarization, `osslsigncode` signing, Homebrew formula push, then a verify job running the four checks; signing steps gated on secrets (OQ-S1), documented in `user-docs/install.md` · **Est:** L

### M4-X1 — M4 exit
- As M1-X1 for M4; first release `v0.1.0` cut only when OQ-S1 is resolved (else documented as pending). **Depends:** all M4 streams · **Est:** S

---

## M5 — Pools, ephemeral teams, kubeconfig, cost tag status (28 operations; 11 tasks)

| Stream | Owns | Tasks |
|---|---|---|
| M5-A pools | `internal/cli/pools/` | A1–A5 |
| M5-B ephemeral teams | `internal/cli/teams/create_ephemeral.go` + tests | B1 |
| M5-C kubeconfig | `internal/kubeconfig/`, `internal/cli/teams/kubeconfig.go` + tests | C1–C2 |
| M5-D cost tag status | `internal/cli/admin/cost_tag_status.go` + tests | D1 |

### M5-0 — Contract re-sync
- As M1-0. **Depends:** M4-X1 · **Est:** S

### M5-A1 — stages (4)
- **Depends:** M5-0 · **Covers:** CLI-17 (stages), `--position` · **Est:** S

### M5-A2 — `list-pools`, `create-pool`, `describe-pool`, `update-pool`, `delete-pool`
- **Depends:** M5-0 · **Covers:** CLI-US-8.3 (via D-C31: `describe-pool` table shows desired workers, state, phase; the five counts in `get-pool-status`, A3), CLI-US-8.4 (`create-pool` on a persistent team → `Conflict` 254), `--wait` · **Red:** `TestCreatePoolPersistentConflict254`, `TestDescribePoolTableGolden` · **Est:** M

### M5-A3 — `scale-pool`, `pause-pool`, `resume-pool`, `get-pool-status`
- **Depends:** A2 · **Covers:** CLI-US-8.2 (`--workers 250` → 252; `--workers 20` sends `ScalePool` with `Idempotency-Key`; `--wait` returns 0 when 20 of 20 report; `--wait-timeout` → 255), CLI-US-8.3 (`get-pool-status --output table` shows desired, reporting, working, idle, impaired), CLI-12 pool rule · **Red:** `TestScalePoolRange252`, `TestScalePoolIdempotencyKey`, `TestScalePoolWaitReporting`, `TestScalePoolWaitTimeout255`, `TestGetPoolStatusTableGolden` · **Est:** M

### M5-A4 — handoffs (2), secret bindings (2), pool config values (3)
- **Depends:** A2 · **Covers:** CLI-17, `--to-pool-ids` list, CLI-10 on `PUT`s · **Est:** M

### M5-A5 — pool external systems (2), pool role links (4)
- **Depends:** A2 · **Covers:** CLI-17 · **Est:** S

### M5-B1 — ephemeral `create-team` options
- **Depends:** M5-0 · **Covers:** CLI-17 (`--kind ephemeral` requires `--expires-at`, `--teardown-rule`, `--budget-usd`; `--budget-thresholds`; full object through `--cli-input-json`), CLI-US-8.1 (missing one → 252 before any call), D-C19, D-C32 (extends the M1-A2 override from its own file) · **Red:** `TestEphemeralRequiresAllThree252`, `TestEphemeralBudgetBody`, `TestEphemeralViaInputJSON` · **Est:** S

### M5-C1 — `internal/kubeconfig` merge library
- **Depends:** M5-0 · **Covers:** CLI-16 (merge replacing same-named entries only; atomic write; new file 0600; existing mode kept with warning when wider; `AWS_PROFILE` env injection) · **Red:** `TestMergeKeepsOtherContexts`, `TestMergeTwiceOneContext`, `TestNewFile0600`, `TestExistingModeKeptWarnWide`, `TestAtomicWrite`, `TestExecEnvAWSProfile` · **Est:** M

### M5-C2 — `teams get-kubeconfig`
- **Depends:** C1 · **Covers:** CLI-US-11.1 (merges `afs-<team-slug>` into `~/.kube/config`, no token, exec `aws eks get-token`), CLI-US-11.2 (`--kubeconfig ./kc.yaml`; `--aws-profile ops`; twice → one context), CLI-US-11.3 (viewer kubectl off → `AccessDenied` 254), `--kubeconfig -` to stdout, default path from `$KUBECONFIG`, missing `aws` → warning · **Red:** `TestGetKubeconfigMergeDefaultPath`, `TestGetKubeconfigExplicitPath`, `TestGetKubeconfigStdout`, `TestGetKubeconfigAccessDenied`, `TestMissingAWSWarns` · **Est:** M

### M5-D1 — `admin list-cost-tag-status`
- **Depends:** M5-0 · **Covers:** CLI-1, TAG-7 table · **Est:** S

### M5-X1 — M5 exit
- As M1-X1 for M5; `user-docs/` swarms and kubectl guides. **Depends:** all M5 streams · **Est:** S

---

## M6 — Work orders, workload token, swarm templates, directory, service keys, security operators; feature complete (22 operations; 11 tasks)

| Stream | Owns | Tasks |
|---|---|---|
| M6-A work orders | `internal/cli/workorders/` | A1–A3 |
| M6-B workload token | `internal/auth/workload.go` + tests | B1 |
| M6-C swarm templates | `internal/cli/swarmtemplates/` | C1 |
| M6-D directory and service-level access | `internal/cli/teams/directory.go`, `internal/cli/admin/{service_api_keys,security_operators}.go` + tests | D1–D2 |

### M6-0 — Contract re-sync
- As M1-0; if `-restapi` has merged and tagged, re-pin to the tag here. **Depends:** M5-X1 · **Est:** S

### M6-A1 — `--repos` shorthand and owner classification
- **Depends:** M6-0 · **Covers:** CLI-18 (`host/owner/repo[=<owners>]`; `@` → `owningContact`, `team_` → `owningTeamId`, other → 252; no `=` → no owner), CLI-US-9.2 (slug → 252 saying to use a team id or an email) · **Red:** `TestReposClassification` (table: team, email, both, none, slug, malformed) · **Est:** S

### M6-A2 — `create-work-order`, `list-work-orders`, `describe-work-order`
- **Depends:** A1 · **Covers:** CLI-US-9.1 (security-operator key via `AFS_API_KEY`; one `CreateWorkOrder` with owners as classified; same `--idempotency-key` replays), CLI-US-9.3 (no `=` → no owner fields; `ownerSource: team_config` printed; explicit `=team_…` wins), CLI-US-9.4 (`ambiguousRepos`/`unownedRepos` → error naming repositories and teams, says to give `=team_…` or `=email`, 254), CLI-US-9.5 (second open work order → `Conflict` 254), CLI-US-9.6 (same key on a team command → `AccessDenied` 254), CLI-US-22.2 (disabled template → `ValidationError` 254), `list-work-orders --output csv` · **Red:** `TestCreateWorkOrderReposBody`, `TestCreateWorkOrderReplay`, `TestCreateWorkOrderTeamConfigOwner`, `TestCreateWorkOrderAmbiguous254`, `TestCreateWorkOrderUnowned254`, `TestCreateWorkOrderDuplicateCVE409`, `TestServiceKeyOnTeamCommandDenied`, `TestDisabledTemplate254`, `TestListWorkOrdersCSV` · **Est:** L

### M6-A3 — `update-work-order`, `cancel-work-order`, `list-repos`, `update-repo`, `list-team-work-orders`
- **Depends:** A1 · **Covers:** CLI-US-10.1 (`update-repo --work-order-id --repo host/owner/repo --status pr_opened --pr-url` with the workload token → `Authorization: Bearer`, succeeds — tested with `--if-match` and with a read-granting fake, D-C17), CLI-US-10.2 (workload token on `cancel-work-order` → `AccessDenied` 254), D-C17 hint when the read is denied, security operator on `list-team-work-orders` → 403 (API) · **Red:** `TestUpdateRepoPathSegments`, `TestUpdateRepoWithIfMatch`, `TestUpdateRepoReadDeniedHint254`, `TestCancelWorkOrderWorkloadDenied` · **Est:** M

### M6-B1 — workload token credential
- **Depends:** M6-0 · **Covers:** CLI-5 (`AFS_API_TOKEN_FILE`, default `/var/run/secrets/afs/api/token` when it exists, re-read per request, last in the chain), CLI-US-10.1 (`sso get-caller-identity` shows `type: workload`), D-C29 · **Red:** `TestWorkloadTokenUsedWhenNoOtherCredential`, `TestWorkloadTokenReReadPerRequest`, `TestWorkloadTokenLastInChain`, `TestGetCallerIdentityWorkload` · **Est:** S

### M6-C1 — swarm templates (7)
- **Depends:** M6-0 · **Covers:** CLI-17, CLI-US-22.1 (`disable-swarm-template --swarm-template-id` sends `DisableSwarmTemplate`; `enable-swarm-template` reverses; security operator → `AccessDenied` 254), `delete-swarm-template` with an open work order → `409` · **Red:** `TestDisableEnableSwarmTemplate`, `TestSwarmTemplateSecurityOperatorDenied` · **Est:** M

### M6-D1 — `teams list-team-directory`
- **Depends:** M6-0 · **Covers:** CLI-1, D56 (team-admin → `AccessDenied`) · **Est:** S

### M6-D2 — admin service-level API keys (3), security operators (3)
- **Depends:** M6-0 · **Covers:** NFR security (`admin create-api-key` prints `secretKey` once, D-C25), CLI-1 · **Red:** `TestCreateServiceAPIKeySecretOnce`, `TestAddSecurityOperator` · **Est:** M

### M6-X1 — Feature complete
- **Owns:** integration; `CurrentMilestone = M6`; `TestDoneUpToCurrentMilestone` strict (180/180); `TestEveryStoryCriterionHasATest` (scans test comments for every `CLI-US-n.k` and `5.1-*` id listed in tasks.md Traceability); docs (`make docs`, README, `user-docs/work-orders.md`, ADR-006); AGENTS.md status line
- **Depends:** all M6 streams · **Covers:** G1–G4, success criteria · **Est:** M

### M6-X2 — Re-pin to the `-restapi` release tag
- **Owns:** `api/SOURCE`, `api/openapi.yaml`, generated files (own commit/PR referencing the tag, contract §8.3)
- **Depends:** M6-X1, OQ-S3 (`-restapi` merged to `main` and tagged) · **Covers:** PRD §6 ("at a release tag"), D-C1 · **Done when:** `SOURCE.ref` is the tag; full suite green; diff recorded · **Est:** S

### M6-X3 — Hand-offs outside this repository
- **Owns:** nothing in this repository except implementation-notes.md (orchestrator)
- **Depends:** M6-X1 · **Covers:** AGENTS.md workflow step 6 (root `skills/` entry once released), contract deltas CLI-CD-1..CLI-CD-6 · **Done when:** `-restapi` and the root have been told (message or issue) of CLI-CD-1..6 and the PRD corrections; the root `skills/enterprise-agent-factory-service-cli.md` update is requested with the release tag · **Est:** S

---

## Traceability

### Requirements → tasks

| PRD | Tasks |
|---|---|
| CLI-1 | M0-K2, K4, K5, E1, E7, S4; M1-A1..A4, B1, C1, D1, D2; M2-A1..A5, B2, B3, C1, C2, D1, D2; M3-A1, A2, B2, C1, C2, D1, D2; M4-A1, B1, B2, C1..C3, D1, D2, E1..E4; M5-A1..A5, C2, D1; M6-A2, A3, C1, D1, D2 (Appendix A lists each operation's stream) |
| CLI-2 | M0-K4, E1, O3, O4, T3, X1 |
| CLI-3 | M0-C1, S4, M4-F1 |
| CLI-4 | M0-S1, S2, S3, S4 |
| CLI-5 | M0-C2, M6-B1 |
| CLI-6 | M0-O1..O4, E2; M1-C1, M3-A2, M3-B2, M4-D2, M6-A2 |
| CLI-7 | M0-E2 |
| CLI-8 | M0-K4, T1, E5 |
| CLI-9 | M2-B1, M2-B2 |
| CLI-10 | M0-E3; M1-A3, M1-A4, M1-B1, M2-A3, M2-A5, M2-C1, M3-A1, M4-C1, M5-A4, M6-A3 |
| CLI-11 | M0-T2 |
| CLI-12 | M0-E4; M1-A2, M2-A2, M5-A3 |
| CLI-13 | M3-B1, M3-B2, M3-A2 |
| CLI-14 | M0-E6, K5 (+ "command done" in every service task, enforced at each `-X1`) |
| CLI-15 | M0-K4 (`--version`), M4-G1, M4-G2 |
| CLI-16 | M5-C1, M5-C2 |
| CLI-17 | M5-A1..A5, M5-B1, M6-C1 |
| CLI-18 | M6-A1, M6-A2, M6-A3 |
| NFR startup | M0-X1 |
| NFR reliability | M0-T2, T3, E4 |
| NFR observability | M0-T1, T4 |
| NFR compatibility | M0-T1 |
| NFR security | M0-T4, M2-B1, B2, M4-A1, M6-D2 |
| NFR quality | M0-K1 (`cover`), E7, O3 |
| NFR accessibility | M0-O3 |
| PRD §6 dependencies | M0-K1 (snapshot), K3 (oapi-codegen), O2 (jmespath), S2 (IdP), M5-C2 (aws on PATH), M4-G2 (signing, tap), M6-X2 (release tag) |
| PRD §9 definition of done | "command done" rule + `-X1` gates |

### Stories → tasks

| Story | Criteria → task |
|---|---|
| CLI-US-1 | .1 S2 · .2 S3 · .3 S1 · .4 S4 · .5 S4 |
| CLI-US-2 | .1 M1-A2 (generic: M0-E1) · .2 M1-A2 · .3 M1-A2 (generic: M0-E4) |
| CLI-US-3 | .1 M2-B2 · .2 M2-B1, B2 · .3 M2-B2 · .4 M0-T4, M2-B2 |
| CLI-US-4 | .1 M0-O1, E7 · .2 M3-B2 · .3 M0-E2, M3-B2 · .4 M0-E2, M1-C1 · .5 M0-E2, M2-A1 · .6 M0-O2, M1-A1 |
| CLI-US-5 | .1 M0-K4, E1 · .2 M0-E5 · .3 M0-E5, T3, K4 |
| CLI-US-6 | .1 M0-E3, M2-A3 · .2 M0-E3, M2-A3 · .3 M0-E3, M2-A3, M3-A1, M4-C1 |
| CLI-US-7 | .1 M0-K4 (strict at M6-X1) |
| CLI-US-8 | .1 M5-B1 · .2 M5-A3 (generic: M0-E4) · .3 M5-A2, A3 (D-C31) · .4 M5-A2 |
| CLI-US-9 | .1 M6-A2 · .2 M6-A1 · .3 M6-A2 · .4 M6-A2 · .5 M6-A2 · .6 M6-A2 |
| CLI-US-10 | .1 M6-A3, M6-B1 (D-C17) · .2 M6-A3 |
| CLI-US-11 | .1 M5-C2 · .2 M5-C1, C2 · .3 M5-C2 |
| CLI-US-12 | .1 M2-A2 · .2 M2-A2 (generic: M0-E1) · .3 M2-A2 · .4 M2-A2 |
| CLI-US-13 | .1 M4-B1 · .2 M4-B1 · .3 M4-B1 |
| CLI-US-14 | .1 M3-A1 · .2 M3-A1 · .3 M3-A1 |
| CLI-US-15 | .1 M4-E3 · .2 M4-E3 · .3 M1-D1 |
| CLI-US-16 | .1 M4-C1 · .2 M4-C1 · .3 M4-C1 |
| CLI-US-17 | .1 M4-D1 · .2 M4-D1 |
| CLI-US-18 | .1 M4-E4 · .2 M4-E4 · .3 M4-E4 |
| CLI-US-19 | .1 M3-C2 · .2 M3-C1, M4-C2 |
| CLI-US-20 | .1 M4-C3 · .2 M4-C3 |
| CLI-US-21 | .1 M1-D1 · .2 M1-D1 |
| CLI-US-22 | .1 M6-C1 · .2 M6-A2 |
| CLI-US-23 | .1 M1-A3, M2-A3 · .2 M1-A3 |
| 5.1 global options | M0-K4, X1, O3, O4, E1 |
| 5.1 configuration | M0-C1, S4 |
| 5.1 credentials | M0-C2, T4 |
| 5.1 pagination | M0-E2 |
| 5.1 retries | M0-T2 |
| 5.1 help | M0-E6 |
| 5.1 distribution | M4-G1, M4-G2 |

---

## Appendix A — Operation → milestone → stream (180 operations)

Source for `internal/cli/milestones.go` (M0-K5). Generated from the snapshot's `x-afs-cli` and PRD §7.
"If-Match" is the snapshot's parameter: *required* (`PATCH`, `DELETE`) or *optional* (`PUT`, CD-2).

| # | operationId | Command | Method and path | CSV | If-Match | Milestone | Stream |
|---|---|---|---|:-:|---|---|---|
| 1 | `GetMe` | `sso get-caller-identity` | `GET /v1/me` |  |  | M0 | M0-S |
| 2 | `ListSignInProviders` | `sso list-providers` | `GET /v1/identity-providers` |  |  | M0 | M0-S |
| 3 | `CreateTeam` | `teams create-team` | `POST /v1/teams` |  |  | M1 | M1-A |
| 4 | `DeleteTeam` | `teams delete-team` | `DELETE /v1/teams/{teamId}` |  | required | M1 | M1-A |
| 5 | `GetRevision` | `teams describe-revision` | `GET /v1/teams/{teamId}/revisions/{rev}` |  |  | M1 | M1-A |
| 6 | `GetTeam` | `teams describe-team` | `GET /v1/teams/{teamId}` |  |  | M1 | M1-A |
| 7 | `GetTeamStatus` | `teams get-team-status` | `GET /v1/teams/{teamId}/status` |  |  | M1 | M1-A |
| 8 | `ListRevisions` | `teams list-revisions` | `GET /v1/teams/{teamId}/revisions` |  |  | M1 | M1-A |
| 9 | `ListTeamModels` | `teams list-team-models` | `GET /v1/teams/{teamId}/models` |  |  | M1 | M1-A |
| 10 | `ListTeams` | `teams list-teams` | `GET /v1/teams` |  |  | M1 | M1-A |
| 11 | `PutTeamModels` | `teams put-team-models` | `PUT /v1/teams/{teamId}/models` |  | optional | M1 | M1-A |
| 12 | `ResumeTeam` | `teams resume-team` | `POST /v1/teams/{teamId}:resume` |  |  | M1 | M1-A |
| 13 | `SuspendTeam` | `teams suspend-team` | `POST /v1/teams/{teamId}:suspend` |  |  | M1 | M1-A |
| 14 | `UpdateTeam` | `teams update-team` | `PATCH /v1/teams/{teamId}` |  | required | M1 | M1-A |
| 15 | `AddMember` | `members add-member` | `POST /v1/teams/{teamId}/members` |  |  | M1 | M1-B |
| 16 | `ListMembers` | `members list-members` | `GET /v1/teams/{teamId}/members` |  |  | M1 | M1-B |
| 17 | `RemoveMember` | `members remove-member` | `DELETE /v1/teams/{teamId}/members/{bindingId}` |  | required | M1 | M1-B |
| 18 | `UpdateMember` | `members update-member` | `PATCH /v1/teams/{teamId}/members/{bindingId}` |  | required | M1 | M1-B |
| 19 | `ListAuditEvents` | `audit list-audit-events` | `GET /v1/teams/{teamId}/audit-events` | CSV |  | M1 | M1-C |
| 20 | `CreateIdentityProvider` | `admin create-identity-provider` | `POST /v1/admin/identity-providers` |  |  | M1 | M1-D |
| 21 | `DeleteIdentityProvider` | `admin delete-identity-provider` | `DELETE /v1/admin/identity-providers/{id}` |  | required | M1 | M1-D |
| 22 | `GetIdentityProvider` | `admin describe-identity-provider` | `GET /v1/admin/identity-providers/{id}` |  |  | M1 | M1-D |
| 23 | `GetSettings` | `admin get-settings` | `GET /v1/admin/settings` |  |  | M1 | M1-D |
| 24 | `ListIdentityProviders` | `admin list-identity-providers` | `GET /v1/admin/identity-providers` |  |  | M1 | M1-D |
| 25 | `PutSetting` | `admin put-setting` | `PUT /v1/admin/settings/{key}` |  | optional | M1 | M1-D |
| 26 | `UpdateIdentityProvider` | `admin update-identity-provider` | `PATCH /v1/admin/identity-providers/{id}` |  | required | M1 | M1-D |
| 27 | `CreateAgent` | `agents create-agent` | `POST /v1/teams/{teamId}/agents` |  |  | M2 | M2-A |
| 28 | `CreateAgentSnapshot` | `agents create-snapshot` | `POST /v1/teams/{teamId}/agents/{agentId}/snapshots` |  |  | M2 | M2-A |
| 29 | `DeleteAgent` | `agents delete-agent` | `DELETE /v1/teams/{teamId}/agents/{agentId}` |  | required | M2 | M2-A |
| 30 | `DeleteAgentConfigValue` | `agents delete-config-value` | `DELETE /v1/teams/{teamId}/agents/{agentId}/config-values/{name}` |  | required | M2 | M2-A |
| 31 | `GetAgent` | `agents describe-agent` | `GET /v1/teams/{teamId}/agents/{agentId}` |  |  | M2 | M2-A |
| 32 | `GetAgentStatus` | `agents get-agent-status` | `GET /v1/teams/{teamId}/agents/{agentId}/status` |  |  | M2 | M2-A |
| 33 | `ListAgentTemplates` | `agents list-agent-templates` | `GET /v1/agent-templates` |  |  | M2 | M2-A |
| 34 | `ListAgents` | `agents list-agents` | `GET /v1/teams/{teamId}/agents` |  |  | M2 | M2-A |
| 35 | `ListAgentConfigValues` | `agents list-config-values` | `GET /v1/teams/{teamId}/agents/{agentId}/config-values` |  |  | M2 | M2-A |
| 36 | `ListAgentExternalSystems` | `agents list-external-systems` | `GET /v1/teams/{teamId}/agents/{agentId}/external-systems` |  |  | M2 | M2-A |
| 37 | `ListAgentSnapshots` | `agents list-snapshots` | `GET /v1/teams/{teamId}/agents/{agentId}/snapshots` |  |  | M2 | M2-A |
| 38 | `ListStorageClasses` | `agents list-storage-classes` | `GET /v1/storage-classes` |  |  | M2 | M2-A |
| 39 | `PauseAgent` | `agents pause-agent` | `POST /v1/teams/{teamId}/agents/{agentId}:pause` |  |  | M2 | M2-A |
| 40 | `PutAgentConfigValue` | `agents put-config-value` | `PUT /v1/teams/{teamId}/agents/{agentId}/config-values/{name}` |  | optional | M2 | M2-A |
| 41 | `RestartAgent` | `agents restart-agent` | `POST /v1/teams/{teamId}/agents/{agentId}:restart` |  |  | M2 | M2-A |
| 42 | `ResumeAgent` | `agents resume-agent` | `POST /v1/teams/{teamId}/agents/{agentId}:resume` |  |  | M2 | M2-A |
| 43 | `SetAgentExternalSystemState` | `agents set-external-system-state` | `PUT /v1/teams/{teamId}/agents/{agentId}/external-systems/{systemId}/state` |  | optional | M2 | M2-A |
| 44 | `UpdateAgent` | `agents update-agent` | `PATCH /v1/teams/{teamId}/agents/{agentId}` |  | required | M2 | M2-A |
| 45 | `CreateSecret` | `secrets create-secret` | `POST /v1/teams/{teamId}/secrets` |  |  | M2 | M2-B |
| 46 | `DeleteSecret` | `secrets delete-secret` | `DELETE /v1/teams/{teamId}/secrets/{secretId}` |  | required | M2 | M2-B |
| 47 | `GetSecret` | `secrets describe-secret` | `GET /v1/teams/{teamId}/secrets/{secretId}` |  |  | M2 | M2-B |
| 48 | `ListSecrets` | `secrets list-secrets` | `GET /v1/teams/{teamId}/secrets` |  |  | M2 | M2-B |
| 49 | `PutSecretValue` | `secrets put-secret-value` | `PUT /v1/teams/{teamId}/secrets/{secretId}/value` |  | optional | M2 | M2-B |
| 50 | `RotateSecret` | `secrets rotate-secret` | `POST /v1/teams/{teamId}/secrets/{secretId}:rotate` |  |  | M2 | M2-B |
| 51 | `UpdateSecret` | `secrets update-secret` | `PATCH /v1/teams/{teamId}/secrets/{secretId}` |  | required | M2 | M2-B |
| 52 | `DeclareTeamExternalSystem` | `teams declare-external-system` | `POST /v1/teams/{teamId}/external-systems` |  |  | M2 | M2-C |
| 53 | `DeleteNamespaceConfigValue` | `teams delete-config-value` | `DELETE /v1/teams/{teamId}/config-values/{name}` |  | required | M2 | M2-C |
| 54 | `DeleteTeamExternalSystem` | `teams delete-external-system` | `DELETE /v1/teams/{teamId}/external-systems/{systemId}` |  | required | M2 | M2-C |
| 55 | `ListNamespaceConfigValues` | `teams list-config-values` | `GET /v1/teams/{teamId}/config-values` |  |  | M2 | M2-C |
| 56 | `ListTeamExternalSystems` | `teams list-external-systems` | `GET /v1/teams/{teamId}/external-systems` |  |  | M2 | M2-C |
| 57 | `PutNamespaceConfigValue` | `teams put-config-value` | `PUT /v1/teams/{teamId}/config-values/{name}` |  | optional | M2 | M2-C |
| 58 | `SetTeamExternalSystemState` | `teams set-external-system-state` | `PUT /v1/teams/{teamId}/external-systems/{systemId}/state` |  | optional | M2 | M2-C |
| 59 | `UpdateTeamExternalSystem` | `teams update-external-system` | `PATCH /v1/teams/{teamId}/external-systems/{systemId}` |  | required | M2 | M2-C |
| 60 | `CreateCatalogSystem` | `admin create-external-system` | `POST /v1/admin/external-systems` |  |  | M2 | M2-D |
| 61 | `DeleteAgentTemplate` | `admin delete-agent-template` | `DELETE /v1/admin/agent-templates/{name}` |  | required | M2 | M2-D |
| 62 | `DeleteCatalogSystem` | `admin delete-external-system` | `DELETE /v1/admin/external-systems/{systemId}` |  | required | M2 | M2-D |
| 63 | `DeleteStorageClass` | `admin delete-storage-class` | `DELETE /v1/admin/storage-classes/{name}` |  | required | M2 | M2-D |
| 64 | `GetCatalogSystem` | `admin describe-external-system` | `GET /v1/admin/external-systems/{systemId}` |  |  | M2 | M2-D |
| 65 | `ListAdminAgentTemplates` | `admin list-agent-templates` | `GET /v1/admin/agent-templates` |  |  | M2 | M2-D |
| 66 | `ListCatalogSystems` | `admin list-external-systems` | `GET /v1/admin/external-systems` |  |  | M2 | M2-D |
| 67 | `ListStorageClassesAdmin` | `admin list-storage-classes` | `GET /v1/admin/storage-classes` |  |  | M2 | M2-D |
| 68 | `PutAgentTemplate` | `admin put-agent-template` | `PUT /v1/admin/agent-templates/{name}` |  | optional | M2 | M2-D |
| 69 | `PutStorageClass` | `admin put-storage-class` | `PUT /v1/admin/storage-classes/{name}` |  | optional | M2 | M2-D |
| 70 | `UpdateCatalogSystem` | `admin update-external-system` | `PATCH /v1/admin/external-systems/{systemId}` |  | required | M2 | M2-D |
| 71 | `ListModels` | `models list-models` | `GET /v1/models` |  |  | M2 | M2-D |
| 72 | `GetBudget` | `budgets get-budget` | `GET /v1/teams/{teamId}/budget` |  |  | M3 | M3-A |
| 73 | `GetUsage` | `budgets get-usage` | `GET /v1/teams/{teamId}/usage` | CSV |  | M3 | M3-A |
| 74 | `ListCosts` | `budgets list-costs` | `GET /v1/teams/{teamId}/costs` | CSV |  | M3 | M3-A |
| 75 | `PutBudget` | `budgets put-budget` | `PUT /v1/teams/{teamId}/budget` |  | optional | M3 | M3-A |
| 76 | `GetAgentReport` | `reports get-agent-report` | `GET /v1/teams/{teamId}/agents/{agentId}/reports/{kind}` | CSV |  | M3 | M3-B |
| 77 | `GetServiceReport` | `reports get-service-report` | `GET /v1/reports/{kind}` | CSV |  | M3 | M3-B |
| 78 | `GetTeamReport` | `reports get-team-report` | `GET /v1/teams/{teamId}/reports/{kind}` | CSV |  | M3 | M3-B |
| 79 | `CreateNotificationChannel` | `teams create-notification-channel` | `POST /v1/teams/{teamId}/notification-channels` |  |  | M3 | M3-C |
| 80 | `DeleteNotificationChannel` | `teams delete-notification-channel` | `DELETE /v1/teams/{teamId}/notification-channels/{id}` |  | required | M3 | M3-C |
| 81 | `GetNotificationChannel` | `teams describe-notification-channel` | `GET /v1/teams/{teamId}/notification-channels/{id}` |  |  | M3 | M3-C |
| 82 | `ListNotificationChannels` | `teams list-notification-channels` | `GET /v1/teams/{teamId}/notification-channels` |  |  | M3 | M3-C |
| 83 | `TestNotificationChannel` | `teams test-notification-channel` | `POST /v1/teams/{teamId}/notification-channels/{id}:test` |  |  | M3 | M3-C |
| 84 | `UpdateNotificationChannel` | `teams update-notification-channel` | `PATCH /v1/teams/{teamId}/notification-channels/{id}` |  | required | M3 | M3-C |
| 85 | `CreateModel` | `admin create-model` | `POST /v1/admin/models` |  |  | M3 | M3-D |
| 86 | `CreateModelPrice` | `admin create-model-price` | `POST /v1/admin/models/{modelId}/prices` |  |  | M3 | M3-D |
| 87 | `ListModelPrices` | `admin list-model-prices` | `GET /v1/admin/models/{modelId}/prices` |  |  | M3 | M3-D |
| 88 | `ListAdminModels` | `admin list-models` | `GET /v1/admin/models` |  |  | M3 | M3-D |
| 89 | `UpdateModel` | `admin update-model` | `PATCH /v1/admin/models/{modelId}` |  | required | M3 | M3-D |
| 90 | `CreateApiKey` | `api-keys create-api-key` | `POST /v1/teams/{teamId}/api-keys` |  |  | M4 | M4-A |
| 91 | `DeleteApiKey` | `api-keys delete-api-key` | `DELETE /v1/teams/{teamId}/api-keys/{keyId}` |  | required | M4 | M4-A |
| 92 | `ListApiKeys` | `api-keys list-api-keys` | `GET /v1/teams/{teamId}/api-keys` |  |  | M4 | M4-A |
| 93 | `CreateAgentIdentity` | `agents create-identity` | `POST /v1/teams/{teamId}/agents/{agentId}/identities` |  |  | M4 | M4-B |
| 94 | `CreateAgentRoleLink` | `agents create-role-link` | `POST /v1/teams/{teamId}/agents/{agentId}/role-links` |  |  | M4 | M4-B |
| 95 | `DeleteAgentIdentity` | `agents delete-identity` | `DELETE /v1/teams/{teamId}/agents/{agentId}/identities/{identityId}` |  | required | M4 | M4-B |
| 96 | `DeleteAgentRoleLink` | `agents delete-role-link` | `DELETE /v1/teams/{teamId}/agents/{agentId}/role-links/{linkId}` |  | required | M4 | M4-B |
| 97 | `GetAgentRoleLink` | `agents describe-role-link` | `GET /v1/teams/{teamId}/agents/{agentId}/role-links/{linkId}` |  |  | M4 | M4-B |
| 98 | `ListAgentIdentities` | `agents list-identities` | `GET /v1/teams/{teamId}/agents/{agentId}/identities` |  |  | M4 | M4-B |
| 99 | `ListAgentRoleLinks` | `agents list-role-links` | `GET /v1/teams/{teamId}/agents/{agentId}/role-links` |  |  | M4 | M4-B |
| 100 | `VerifyAgentRoleLink` | `agents verify-role-link` | `POST /v1/teams/{teamId}/agents/{agentId}/role-links/{linkId}:verify` |  |  | M4 | M4-B |
| 101 | `CreateLogDestination` | `teams create-log-destination` | `POST /v1/teams/{teamId}/log-destinations` |  |  | M4 | M4-C |
| 102 | `DeleteLogDestination` | `teams delete-log-destination` | `DELETE /v1/teams/{teamId}/log-destinations/{id}` |  | required | M4 | M4-C |
| 103 | `GetLogDestination` | `teams describe-log-destination` | `GET /v1/teams/{teamId}/log-destinations/{id}` |  |  | M4 | M4-C |
| 104 | `GetTeamTags` | `teams get-team-tags` | `GET /v1/teams/{teamId}/tags` |  |  | M4 | M4-C |
| 105 | `ListAllowedAccounts` | `teams list-allowed-accounts` | `GET /v1/teams/{teamId}/allowed-accounts` |  |  | M4 | M4-C |
| 106 | `ListLogDestinations` | `teams list-log-destinations` | `GET /v1/teams/{teamId}/log-destinations` |  |  | M4 | M4-C |
| 107 | `PutTeamTags` | `teams put-team-tags` | `PUT /v1/teams/{teamId}/tags` |  | optional | M4 | M4-C |
| 108 | `RestoreRevision` | `teams restore-revision` | `POST /v1/teams/{teamId}/revisions/{rev}:restore` |  |  | M4 | M4-C |
| 109 | `UpdateLogDestination` | `teams update-log-destination` | `PATCH /v1/teams/{teamId}/log-destinations/{id}` |  | required | M4 | M4-C |
| 110 | `ListAdminAuditEvents` | `audit list-service-audit-events` | `GET /v1/admin/audit-events` | CSV |  | M4 | M4-D |
| 111 | `CreateReportSchedule` | `reports create-report-schedule` | `POST /v1/teams/{teamId}/report-schedules` |  |  | M4 | M4-D |
| 112 | `DeleteReportSchedule` | `reports delete-report-schedule` | `DELETE /v1/teams/{teamId}/report-schedules/{scheduleId}` |  | required | M4 | M4-D |
| 113 | `ListReportSchedules` | `reports list-report-schedules` | `GET /v1/teams/{teamId}/report-schedules` |  |  | M4 | M4-D |
| 114 | `UpdateReportSchedule` | `reports update-report-schedule` | `PATCH /v1/teams/{teamId}/report-schedules/{scheduleId}` |  | required | M4 | M4-D |
| 115 | `ApproveModelPriceProposal` | `admin approve-price-proposal` | `POST /v1/admin/model-price-proposals/{proposalId}:approve` |  |  | M4 | M4-E |
| 116 | `CreateExport` | `admin create-export` | `POST /v1/exports` |  |  | M4 | M4-E |
| 117 | `CreateImport` | `admin create-import` | `POST /v1/imports` |  |  | M4 | M4-E |
| 118 | `CreateServiceNotificationChannel` | `admin create-notification-channel` | `POST /v1/admin/notification-channels` |  |  | M4 | M4-E |
| 119 | `DeleteServiceNotificationChannel` | `admin delete-notification-channel` | `DELETE /v1/admin/notification-channels/{id}` |  | required | M4 | M4-E |
| 120 | `DeleteServiceTag` | `admin delete-service-tag` | `DELETE /v1/admin/service-tags/{key}` |  | required | M4 | M4-E |
| 121 | `GetExport` | `admin describe-export` | `GET /v1/exports/{exportId}` |  |  | M4 | M4-E |
| 122 | `ListCostReconciliations` | `admin list-cost-reconciliations` | `GET /v1/admin/cost-reconciliations` |  |  | M4 | M4-E |
| 123 | `ListExports` | `admin list-exports` | `GET /v1/exports` |  |  | M4 | M4-E |
| 124 | `ListServiceNotificationChannels` | `admin list-notification-channels` | `GET /v1/admin/notification-channels` |  |  | M4 | M4-E |
| 125 | `ListOrgAccounts` | `admin list-org-accounts` | `GET /v1/admin/org-accounts` |  |  | M4 | M4-E |
| 126 | `ListModelPriceProposals` | `admin list-price-proposals` | `GET /v1/admin/model-price-proposals` |  |  | M4 | M4-E |
| 127 | `ListServiceTags` | `admin list-service-tags` | `GET /v1/admin/service-tags` |  |  | M4 | M4-E |
| 128 | `PutServiceTag` | `admin put-service-tag` | `PUT /v1/admin/service-tags/{key}` |  | optional | M4 | M4-E |
| 129 | `RejectModelPriceProposal` | `admin reject-price-proposal` | `POST /v1/admin/model-price-proposals/{proposalId}:reject` |  |  | M4 | M4-E |
| 130 | `UpdateServiceNotificationChannel` | `admin update-notification-channel` | `PATCH /v1/admin/notification-channels/{id}` |  | required | M4 | M4-E |
| 131 | `CreatePool` | `pools create-pool` | `POST /v1/teams/{teamId}/pools` |  |  | M5 | M5-A |
| 132 | `CreatePoolRoleLink` | `pools create-role-link` | `POST /v1/teams/{teamId}/pools/{poolId}/role-links` |  |  | M5 | M5-A |
| 133 | `CreateStage` | `pools create-stage` | `POST /v1/teams/{teamId}/stages` |  |  | M5 | M5-A |
| 134 | `DeletePoolConfigValue` | `pools delete-config-value` | `DELETE /v1/teams/{teamId}/pools/{poolId}/config-values/{name}` |  | required | M5 | M5-A |
| 135 | `DeletePool` | `pools delete-pool` | `DELETE /v1/teams/{teamId}/pools/{poolId}` |  | required | M5 | M5-A |
| 136 | `DeletePoolRoleLink` | `pools delete-role-link` | `DELETE /v1/teams/{teamId}/pools/{poolId}/role-links/{linkId}` |  | required | M5 | M5-A |
| 137 | `DeleteStage` | `pools delete-stage` | `DELETE /v1/teams/{teamId}/stages/{stageId}` |  | required | M5 | M5-A |
| 138 | `GetPool` | `pools describe-pool` | `GET /v1/teams/{teamId}/pools/{poolId}` |  |  | M5 | M5-A |
| 139 | `GetPoolStatus` | `pools get-pool-status` | `GET /v1/teams/{teamId}/pools/{poolId}/status` |  |  | M5 | M5-A |
| 140 | `ListPoolConfigValues` | `pools list-config-values` | `GET /v1/teams/{teamId}/pools/{poolId}/config-values` |  |  | M5 | M5-A |
| 141 | `ListPoolExternalSystems` | `pools list-external-systems` | `GET /v1/teams/{teamId}/pools/{poolId}/external-systems` |  |  | M5 | M5-A |
| 142 | `ListPoolHandoffs` | `pools list-handoffs` | `GET /v1/teams/{teamId}/pools/{poolId}/handoffs` |  |  | M5 | M5-A |
| 143 | `ListPools` | `pools list-pools` | `GET /v1/teams/{teamId}/pools` |  |  | M5 | M5-A |
| 144 | `ListPoolRoleLinks` | `pools list-role-links` | `GET /v1/teams/{teamId}/pools/{poolId}/role-links` |  |  | M5 | M5-A |
| 145 | `ListPoolSecretBindings` | `pools list-secret-bindings` | `GET /v1/teams/{teamId}/pools/{poolId}/secret-bindings` |  |  | M5 | M5-A |
| 146 | `ListStages` | `pools list-stages` | `GET /v1/teams/{teamId}/stages` |  |  | M5 | M5-A |
| 147 | `PausePool` | `pools pause-pool` | `POST /v1/teams/{teamId}/pools/{poolId}:pause` |  |  | M5 | M5-A |
| 148 | `PutPoolConfigValue` | `pools put-config-value` | `PUT /v1/teams/{teamId}/pools/{poolId}/config-values/{name}` |  | optional | M5 | M5-A |
| 149 | `PutPoolExternalSystem` | `pools put-external-system` | `PUT /v1/teams/{teamId}/pools/{poolId}/external-systems/{systemId}` |  | optional | M5 | M5-A |
| 150 | `PutPoolHandoffs` | `pools put-handoffs` | `PUT /v1/teams/{teamId}/pools/{poolId}/handoffs` |  | optional | M5 | M5-A |
| 151 | `PutPoolSecretBindings` | `pools put-secret-bindings` | `PUT /v1/teams/{teamId}/pools/{poolId}/secret-bindings` |  | optional | M5 | M5-A |
| 152 | `ResumePool` | `pools resume-pool` | `POST /v1/teams/{teamId}/pools/{poolId}:resume` |  |  | M5 | M5-A |
| 153 | `ScalePool` | `pools scale-pool` | `POST /v1/teams/{teamId}/pools/{poolId}:scale` |  |  | M5 | M5-A |
| 154 | `UpdatePool` | `pools update-pool` | `PATCH /v1/teams/{teamId}/pools/{poolId}` |  | required | M5 | M5-A |
| 155 | `UpdateStage` | `pools update-stage` | `PATCH /v1/teams/{teamId}/stages/{stageId}` |  | required | M5 | M5-A |
| 156 | `VerifyPoolRoleLink` | `pools verify-role-link` | `POST /v1/teams/{teamId}/pools/{poolId}/role-links/{linkId}:verify` |  |  | M5 | M5-A |
| 157 | `GetTeamKubeconfig` | `teams get-kubeconfig` | `GET /v1/teams/{teamId}/kubeconfig` |  |  | M5 | M5-C |
| 158 | `ListCostTagStatus` | `admin list-cost-tag-status` | `GET /v1/admin/cost-tag-status` |  |  | M5 | M5-D |
| 159 | `CancelWorkOrder` | `work-orders cancel-work-order` | `POST /v1/work-orders/{workOrderId}:cancel` |  |  | M6 | M6-A |
| 160 | `CreateWorkOrder` | `work-orders create-work-order` | `POST /v1/work-orders` |  |  | M6 | M6-A |
| 161 | `GetWorkOrder` | `work-orders describe-work-order` | `GET /v1/work-orders/{workOrderId}` |  |  | M6 | M6-A |
| 162 | `ListWorkOrderRepos` | `work-orders list-repos` | `GET /v1/work-orders/{workOrderId}/repos` |  |  | M6 | M6-A |
| 163 | `ListTeamWorkOrders` | `work-orders list-team-work-orders` | `GET /v1/teams/{teamId}/work-orders` |  |  | M6 | M6-A |
| 164 | `ListWorkOrders` | `work-orders list-work-orders` | `GET /v1/work-orders` | CSV |  | M6 | M6-A |
| 165 | `UpdateWorkOrderRepo` | `work-orders update-repo` | `PATCH /v1/work-orders/{workOrderId}/repos/{host}/{owner}/{repo}` |  | required | M6 | M6-A |
| 166 | `UpdateWorkOrder` | `work-orders update-work-order` | `PATCH /v1/work-orders/{workOrderId}` |  | required | M6 | M6-A |
| 167 | `CreateSwarmTemplate` | `swarm-templates create-swarm-template` | `POST /v1/admin/swarm-templates` |  |  | M6 | M6-C |
| 168 | `DeleteSwarmTemplate` | `swarm-templates delete-swarm-template` | `DELETE /v1/admin/swarm-templates/{templateId}` |  | required | M6 | M6-C |
| 169 | `GetSwarmTemplate` | `swarm-templates describe-swarm-template` | `GET /v1/admin/swarm-templates/{templateId}` |  |  | M6 | M6-C |
| 170 | `DisableSwarmTemplate` | `swarm-templates disable-swarm-template` | `POST /v1/admin/swarm-templates/{templateId}:disable` |  |  | M6 | M6-C |
| 171 | `EnableSwarmTemplate` | `swarm-templates enable-swarm-template` | `POST /v1/admin/swarm-templates/{templateId}:enable` |  |  | M6 | M6-C |
| 172 | `ListSwarmTemplates` | `swarm-templates list-swarm-templates` | `GET /v1/admin/swarm-templates` |  |  | M6 | M6-C |
| 173 | `UpdateSwarmTemplate` | `swarm-templates update-swarm-template` | `PATCH /v1/admin/swarm-templates/{templateId}` |  | required | M6 | M6-C |
| 174 | `AddSecurityOperator` | `admin add-security-operator` | `POST /v1/admin/security-operators` |  |  | M6 | M6-D |
| 175 | `CreateServiceApiKey` | `admin create-api-key` | `POST /v1/admin/api-keys` |  |  | M6 | M6-D |
| 176 | `DeleteServiceApiKey` | `admin delete-api-key` | `DELETE /v1/admin/api-keys/{keyId}` |  | required | M6 | M6-D |
| 177 | `ListServiceApiKeys` | `admin list-api-keys` | `GET /v1/admin/api-keys` |  |  | M6 | M6-D |
| 178 | `ListSecurityOperators` | `admin list-security-operators` | `GET /v1/admin/security-operators` |  |  | M6 | M6-D |
| 179 | `RemoveSecurityOperator` | `admin remove-security-operator` | `DELETE /v1/admin/security-operators/{bindingId}` |  | required | M6 | M6-D |
| 180 | `ListTeamDirectory` | `teams list-team-directory` | `GET /v1/team-directory` |  |  | M6 | M6-D |
