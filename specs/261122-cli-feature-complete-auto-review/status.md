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
| 5 | Implementation | In Progress | 3 / 20 FRs |
| 6 | Completion and archival | Not Started | |

## Task log
Updated after every fix: FR, commit, tests added, review note.

| FR | Priority | State | Commit subject (git log --grep FR-nnn) | Notes |
|---|---|---|---|---|
| FR-001 | P0 | Done | reject INI injection and invalid SSO provider fields | Set returns an error on control chars and structural chars; ValidateIssuer/ValidateIdentifier; provider validated in Pick and configure sso. Review: second pass tried tab, NUL, DEL, C1, bracket keys, http issuer; all refused. Gate green, 93.5%. |
| FR-002 | P1 | Done | require https and a matching issuer for OIDC discovery and endpoints | Discover validates issuer, issuer equality, endpoint schemes; sso_issuer validated on read (non-lenient) and write. Review: no secret form sent in any rejected case. Gate green, 93.5%. |
| FR-003 | P1 | Done | sanitize terminal control characters in table, text, errors and prompts | output.Sanitize + sanitizing writer in ReportError; device prompt and provider list sanitized; json/yaml unchanged (D-R7). Review: bypass attempts with C1, DEL, invalid UTF-8 covered. 93.6%. |
