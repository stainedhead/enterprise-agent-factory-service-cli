# Plan: paperclip-collaboration

**Date:** 2026-11-22 · **Status:** Planning

## Development approach

Red-Green-Refactor per task, table-driven tests, goldens via `-update`, httptest fakes through `clitest`. Generated files are
only produced by `make contract-sync`, `make generate`, `make docs`.

## Phase breakdown

1. M7: re-pin; milestone tables; redaction; hints table; `collaboration` (describe, put, links, connect, disconnect);
   safeupdate/wait tables; agents; secrets; teams/members/api-keys; docs first cut.
2. M8: `admin list-collaboration-endpoints`; identity provider collab client flags; settings pass-through tests; notice help.
3. M9: `invoke-api`; `rotate-collaboration-keys`; managed-mode hints.
4. M10: swarm-templates; reports; traceability test; docs; CHANGELOG; release notes.

## Critical path

Re-pin -> milestone tables -> hints/redact -> collaboration service -> invoke-api (largest, security heavy).

## Testing strategy

Per-operation `CheckOperations`; criterion-named tests (the traceability test greps for `CLI-US-n.k` test names);
sentinel tests for credentials; coverage gate `make cover` >= 90% on `internal/...`.

## Rollout strategy

Four commits (one per milestone) plus review-fix commits on branch `claude/project-thread-xqpqky`; one PR; tag `v1.0.1` pushed by J.

## Success metrics

All CLI-US-24..35 criteria map to passing tests; gate green; zero hand-edited generated files.
