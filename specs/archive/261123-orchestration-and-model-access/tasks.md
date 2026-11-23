# Tasks: orchestration-and-model-access

## Progress summary

See status.md. Legend: [x] done.

## M11

- [x] T1 Re-pin snapshot to f6d240c, regenerate (spec, client types), fix count and version tests (FR-001)
- [x] T2 Core: Override.Deprecated/DeprecatedOptions/Hidden/Desc/Enum/Notes, Option.Hidden/Deprecated, help output, `exec/deprecation.go` (FR-011)
- [x] T3 `output.Column.NAUnless` and its test
- [x] T4 `orchestration` service: describe, put, change tool, notices, `--wait`, goldens, AccessDenied (FR-002 to FR-004)
- [x] T5 Hints for the 1.4.0 reasons and the changed legacy ones (FR-012)
- [x] T6 `admin list-orchestration-endpoints` with backfill notice (FR-005)
- [x] T7 Deprecated collaboration commands; `invoke-api` on `GetTeamOrchestration` (FR-011, FR-013)

## M12

- [x] T8 `modelaccess` package; agents and pools flags, columns, `--role` hidden (FR-009, FR-010)
- [x] T9 Teams: `--orchestration`, `--max-collaboration-exposure`, `put-team-models --access`, columns, notice (FR-007, FR-008)
- [x] T10 `admin put-service-model-access`, admin and public model columns, settings keys (FR-006, FR-014)
- [x] T11 Swarm template skeleton (FR-015)
- [x] T12 Traceability test over the FR table of spec.md

## Docs, review, release

- [x] T13 `make docs`; README, user-docs, CHANGELOG, AGENTS.md, docs/ (FR-016)
- [x] T14 Root `skills/enterprise-agent-factory-service-cli.md` updated (not committed in the root)
- [x] T15 Code and design review, findings fixed (implementation-notes.md)
- [x] T16 Final gate, PR, CI green, merge
