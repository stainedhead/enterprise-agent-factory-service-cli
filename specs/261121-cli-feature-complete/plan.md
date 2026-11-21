# Plan: `afs` CLI, feature complete

| | |
|---|---|
| **Feature** | cli-feature-complete |
| **Date** | 2026-11-21 |
| **Status** | Planning complete; ready for implementation |

## 1. Development approach

- **Contract first, generated second, hand-written last.** The snapshot drives the generated
  descriptors and types; the generic executor makes every operation callable; service overrides add
  only what the PRD asks beyond the generic behaviour.
- **TDD, Red → Green → Refactor** for every task (AGENTS.md). Red is a failing test that names the
  story criterion it proves (`// CLI-US-6: …`). No test reaches a real API or IdP: `httptest.Server`
  fakes only, with the replay validator checking each request against the snapshot.
- **Milestones in order, streams in parallel.** Each milestone: `Mn-0` (contract re-sync, sequential)
  → streams (parallel worktrees, disjoint `Owns`) → `Mn-X` (integration, gate, docs, milestone bump,
  sequential).
- **Quality gate on every commit**: `make gate` (tidy, fmt-check, vet, lint, race tests, build,
  `--version`) plus `make contract-check`, `make generate-check` and `make cover` (≥ 90%).

## 2. Client generation and contract sync

| Step | Command | What it does | Fails when |
|---|---|---|---|
| Pin | `make contract-sync RESTAPI_REF=<branch, tag or commit>` (or `RESTAPI_DIR=../enterprise-agent-factory-service-restapi`) | Copies `api/openapi.yaml` from `-restapi` (`gh api repos/stainedhead/enterprise-agent-factory-service-restapi/contents/api/openapi.yaml?ref=<ref>` or the local clone), writes `api/SOURCE` (`repo`, `ref`, `commit`, `sha256`, `synced_at`), runs `make generate` | the ref does not exist |
| Verify | `make contract-check` (in `gate`, CI) | `sha256sum api/openapi.yaml` equals `api/SOURCE`'s `sha256`; `internal/spec.SnapshotSHA256` equals it too | the snapshot was edited or `SOURCE` not updated |
| Generate | `make generate` | `go run ./tools/afsgen` → `internal/spec/operations_gen.go`; oapi-codegen v2.8.2 → `internal/client/gen/types.gen.go`; `gofmt` | the snapshot lacks `x-afs-cli` on an operation, or a command name is not `<service> <kebab-verb-noun>` |
| Drift | `make generate-check` (CI) | regenerates into a temp dir and `diff -r` against the committed files | generated files were edited or not regenerated |
| Replay | `go test ./internal/cli/...` (contract replay test) | every operation: request validated against the snapshot, example response rendered as json equals the example, `403` → 254 | the CLI builds a request the contract rejects |
| Table | CLI-US-7 test | registry ↔ snapshot both ways, names equal `x-afs-cli` | a command is missing, extra or misnamed |

**Sync policy.** `-restapi` is mid-implementation on `feat/restapi-feature-complete`; its YAML may still
change. Every milestone starts with `Mn-0`: re-sync to the branch head; if the sha256 changed,
regenerate, run the full suite, record the diff (new/changed operations) in `implementation-notes.md`,
and fix the CLI before streams start. When `-restapi` merges to `main` and tags a release (OQ-S3),
M6-X2 re-pins `api/SOURCE` to the tag (`ref: vX.Y.Z`) in its own commit, as contract §8 asks ("`-cli`
regenerates its client … in its own PR referencing the `-restapi` tag"). If the merge happens earlier,
the next `Mn-0` re-pins to the tag instead of the branch.

## 3. Phase breakdown

| Milestone | Theme | Operations | Streams | Tasks |
|---|---|---|---|---|
| M0 | Foundation: snapshot, generators, registry, transport, config, credentials, SSO, output, executor, help, replay test | 2 (`GetMe`, `ListSignInProviders`) | K (sequential), T, C, S, O, E | 28 |
| M1 | Teams core, members, team audit, settings, identity providers | 24 | A, B, C, D | 10 |
| M2 | Agents, secrets, team config and external systems, models, admin catalogs | 45 | A, B, C, D | 14 |
| M3 | Budgets, on-demand reports, team notification channels, admin models and prices | 18 | A, B, C, D | 10 |
| M4 | API keys, role links, identities, tags, log destinations, restore, schedules, admin service items, exports, completion, releases | 41 | A–G | 17 |
| M5 | Pools, ephemeral teams, kubeconfig, cost tag status | 28 | A, B, C, D | 11 |
| M6 | Work orders, workload token, swarm templates, team directory, service keys, security operators; re-pin; feature complete | 22 | A, B, C, D | 11 |

Totals: 180 operations, 101 tasks (tasks.md).

## 4. Critical path

```
M0-K1 → K2 → K3 → K4 → K5 ─┬─ T1 → T2 → T3 → T4 ─────────┐
                           ├─ C1 → C2 ── S1 → S2 → S3 → S4 ┤
                           ├─ O1 → O2 → O3 → O4 ──────────┤
                           └─ E1 → E2 → E3 → E4 → E5 → E6 → E7 ─┴─ M0-X1 → M0-X2
→ M1 → M2 → M3 → M4 → M5 → M6 (each: Mn-0 → streams → Mn-X)
→ M6-X2 (re-pin to -restapi release tag; waits for OQ-S3)
```

The executor stream (E) is the longest in M0; it builds against the interfaces K4 fixes
(`Transport`, `Renderer`, `CredentialSource`) with fakes, so it does not wait for T, C or O.

## 5. Testing strategy

| Layer | Technique |
|---|---|
| Generators | Golden output for a small fixture OpenAPI; assertions on the real snapshot (182 ops, 180 commands, 52 paginated, 8 CSV, 59 `If-Match`, 44 idempotent) |
| `internal/client` | `httptest.Server`; injected dialer (connect failure vs dropped-after-send), fake clock and jitter source; debug-log golden with redaction |
| `internal/config`, `internal/auth`, `internal/sso` | temp `HOME`; env maps; file modes; `httptest` IdP with discovery, device authorization, token (pending, slow_down, success, refresh, expired) |
| `internal/output` | golden files (`-update`), 80-column check, `NO_COLOR`, pager with a fake command |
| `internal/cli/exec` and services | table-driven per operation: argument parsing, request sent (method, path, query, headers, body), rendering per `--output`, error mapping, exit code |
| Contract replay (all 180) | request validated by `libopenapi-validator` against the snapshot; example 2xx → `--output json` equals the example; `403` → 254 with `operationId` and `requiredPermission` |
| Help | golden help per operation; examples present for every operation at or below `CurrentMilestone` |
| Release | `goreleaser check`; `release --snapshot`; verify job (cosign, osslsigncode, spctl) behind secrets |

Coverage: `internal/...` ≥ 90% (`make cover` fails below; `tools/` and generated files excluded).

## 6. Rollout strategy

- Work on `feat/cli-feature-complete`; stream branches `feat/cli-feature-complete--<stream>` in worktrees,
  merged back by the orchestrator at stream completion (fast-forward or merge commit, never force-push).
- Each milestone exit leaves `feat/cli-feature-complete` green. Per J's instruction, when stable and
  green the orchestrator pushes to `main`.
- First tagged release (`v0.1.0`) after M4 (signed release pipeline exists); `v1.0.0` after M6 and the
  re-pin to the `-restapi` release tag.
- After each release: the root repository's `skills/enterprise-agent-factory-service-cli.md` is updated
  by a root PR (outside this repository; listed in M6-X3 as a hand-off).

## 7. Success metrics

- 180/180 operations with a command, example, golden output and `AccessDenied` test.
- 23/23 stories and every §5.1 criterion with named tests.
- `make gate`, `contract-check`, `generate-check`, `cover` green; govulncheck clean.
- Startup to first request < 100 ms on a warm machine (benchmark in CI, reported).
