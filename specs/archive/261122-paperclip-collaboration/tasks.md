# Tasks: paperclip-collaboration

**Date:** 2026-11-22 · **Status:** Planning (see status.md for progress)

## Progress summary

35/38 tasks complete (M7 to M10 done; X1 to X3 continue).

## M7 (C1)

- [x] P7.1 Re-pin snapshot, regenerate, fix count/version tests, milestone and stream entries for 7 operations (FR-001)
- [x] P7.2 Redact edge headers in `internal/redact` and the client debug writer (FR-008)
- [x] P7.3 Reason hint table and `details.fields[]` printing in `exec` (FR-009)
- [x] P7.4 `collaboration describe-collaboration` + notice + missing-agent hint (FR-002)
- [x] P7.5 `collaboration put-collaboration` body, required surfaces, If-Match reader, `--wait` (FR-003)
- [x] P7.6 `list-collaboration-links`, `connect-agent`, `disconnect-agent` (FR-004)
- [x] P7.7 `agents` flags, validation, columns, skeleton, `list-agent-templates --capability`, `admin put-agent-template --capabilities` (FR-005)
- [x] P7.8 `secrets` managed-by, column, hints (FR-006)
- [x] P7.9 `teams` guardrail flag, `hostname_held` hint, `members`/`api-keys` role (FR-007)
- [x] P7.10 safeupdate and wait table checks (US-34.3)
- [x] P7.11 Docs first cut; raise `CurrentMilestone` to 7; gate; commit; push

## M8 (C2)

- [x] P8.1 `admin list-collaboration-endpoints` columns, CSV (FR-010)
- [x] P8.2 Identity provider `--collab-client-id`, `--collab-client-secret-arn` validation (FR-011)
- [x] P8.3 `collab.*` settings pass-through tests; redaction of ARN vs secret (FR-011)
- [x] P8.4 Config profile key `collab_base_domain`
- [x] P8.5 Docs; `CurrentMilestone` 8; gate; commit; push

## M9 (C3)

- [x] P9.1 Origin validation (FR-012)
- [x] P9.2 Paperclip key reader (file mode, stdin, shape)
- [x] P9.3 Credential chain and workload refusal
- [x] P9.4 Request transport: no redirects, no cookies, retries, streaming, timeouts
- [x] P9.5 Edge error mapping, HttpError, hints
- [x] P9.6 Debug logging with redaction, query dropped
- [x] P9.7 Global option refusals; help, examples, goldens
- [x] P9.8 `rotate-collaboration-keys` (FR-013)
- [x] P9.9 Sentinel tests; docs; `CurrentMilestone` 9; gate; commit; push

## M10 (C4)

- [x] P10.1 swarm-templates passthrough tests and skeleton (FR-014)
- [x] P10.2 pools no collaboration option test
- [x] P10.3 Reports kind collaboration, group-by (FR-015)
- [x] P10.4 Traceability test over CLI-US-24..34
- [x] P10.5 Docs complete, CHANGELOG, `user-docs`
- [x] P10.6 `CurrentMilestone` 10; gate; commit; push

## Cross-cutting

- [ ] X1 `make gate` green and coverage >= 90% at each milestone
- [ ] X2 status.md and dashboard updated after each task
- [ ] X3 Follow-up: re-pin SOURCE to the restapi merge commit
