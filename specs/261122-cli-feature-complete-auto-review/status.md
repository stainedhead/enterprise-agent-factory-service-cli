# Status: cli-feature-complete-auto-review

| | |
|---|---|
| **Feature** | Review fixes for the `afs` CLI feature-complete branch |
| **Created** | 2026-11-22 |
| **Branch** | `feat/cli-feature-complete` |
| **Dashboard** | [cli-feature-complete-status.md](../archive/261121-cli-feature-complete/cli-feature-complete-status.md) |

**MANDATORY:** update after every task.

## Overall progress

| Phase | Name | Status | Notes |
|---|---|---|---|
| 0 | Initial research (review PRD) | Complete | PRD reviewed in step 6 |
| 1 | Specification | Complete | spec.md |
| 2 | Research and data modelling | Complete | research.md, data-dictionary.md |
| 3 | Architecture and planning | Complete | architecture.md, plan.md |
| 4 | Task breakdown | Complete | tasks.md, one task per FR |
| 5 | Implementation | In Progress | 16 / 20 FRs |
| 6 | Completion and archival | Not Started | |

## Task log
Updated after every fix: FR, commit, tests added, review note.

| FR | Priority | State | Commit subject (git log --grep FR-nnn) | Notes |
|---|---|---|---|---|
| FR-001 | P0 | Done | reject INI injection and invalid SSO provider fields | Set returns an error on control chars and structural chars; ValidateIssuer/ValidateIdentifier; provider validated in Pick and configure sso. Review: second pass tried tab, NUL, DEL, C1, bracket keys, http issuer; all refused. Gate green, 93.5%. |
| FR-002 | P1 | Done | require https and a matching issuer for OIDC discovery and endpoints | Discover validates issuer, issuer equality, endpoint schemes; sso_issuer validated on read (non-lenient) and write. Review: no secret form sent in any rejected case. Gate green, 93.5%. |
| FR-003 | P1 | Done | sanitize terminal control characters in table, text, errors and prompts | output.Sanitize + sanitizing writer in ReportError; device prompt and provider list sanitized; json/yaml unchanged (D-R7). Review: bypass attempts with C1, DEL, invalid UTF-8 covered. 93.6%. |
| FR-004 | P1 | Done | print NextToken on stderr for CSV with --no-paginate; resume without header | fetchCSV reports X-Afs-Next-Token; CSVWriter.Resume. Review: stdout stays clean CSV. 93.6%. |
| FR-006,FR-007,FR-008 | P2 | Done | IdP timeouts, access-token validation, --debug for OIDC traffic | FR-006 IdPClient (transport + read timeout); FR-007 token()/ErrNoAccessToken, Bearer only; FR-008 sso.WithDebug logs method, redacted URL, status, duration, never bodies. Review: no secret reaches debug output. 93.7%. |
| FR-005,FR-009,FR-010,FR-011 | P2 | Done (FR-010 criterion 2 deferred, D-R9) | endpoint userinfo refused and redacted, strict number options, dot path segments refused, Retry-After cap | FR-005 ValidateEndpoint + redactURL (user name masked too); FR-009 strictInt/strictNumber (252); FR-010 '.'/'..' refused, id-pattern check deferred (about 60 tests use abbreviated ids); FR-011 20s or read-timeout cap, debug log reason. Review: all hostile cases tested. 94.0%. |
| FR-012,FR-013,FR-016 | P2 | Done | pagination loop guard, --wait Retry-After and undecodable status, notice after the write | FR-012 seen-token set (json+csv) and non-list result key error; FR-013 pollInterval (capped Retry-After) + ErrUndecodable; FR-016 PRD text after success (ADR-013, D-R4). Review: null list still empty. 94.0%. |
| FR-018,FR-020 | P2 | Done | sensitive option help lists safe forms and is not required; pager runs under the command context | FR-018 help.go + goldens + reference regenerated; FR-020 OpenPager(ctx), pager gets SIGINT on cancel (not killed, so the terminal is restored). 94.1%. |
