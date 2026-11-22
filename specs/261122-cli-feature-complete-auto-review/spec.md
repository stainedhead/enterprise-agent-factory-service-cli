# Spec: cli-feature-complete-auto-review

| | |
|---|---|
| **Created** | 2026-11-22 |
| **Source PRD** | [cli-feature-complete-auto-review-PRD.md](cli-feature-complete-auto-review-PRD.md) (moved here by `create-spec`) |
| **Parent** | `specs/archive/261121-cli-feature-complete/` |

## Executive summary
Fix the 20 findings (1 P0, 3 P1, 16 P2) from the automated code and design review of the `afs` CLI
feature-complete branch. The weak spot is data from the service or identity provider trusted too much:
config-file injection (FR-001), clear-text OIDC (FR-002), terminal escapes (FR-003), and one unmet PRD
requirement (CSV `NextToken`, FR-004), plus reliability, validation, convention and test gaps.

## Problem statement
See PRD section 2. Every fix is local to `-cli`; none needs a contract change.

## Goals / non-goals
Goals G1..G5 and non-goals as in PRD section 2. No new commands or options, no contract change, no
change to `--output json` bytes, M6-X2 and M6-X3 stay with the parent spec.

## Functional requirements
FR-001..FR-020 are defined, with acceptance criteria, in the source PRD (section 3). They are not
restated here. Order: FR-001 alone, then FR-002..FR-004, then FR-005..FR-020.

## Non-functional requirements
PRD section 4: security tests of the hostile case, `--output json` byte stable, bounded network waits,
`--debug` covers all outbound HTTP, startup under 100 ms, gate green, coverage >= 90%.

## System architecture
Layers touched: `internal/config`, `internal/sso`, `internal/cli/{configure,sso,exec,wait,safeupdate,runtime}`,
`internal/output`, `internal/client`, `internal/kubeconfig`, `internal/cli/cli.go`. No new layer, no new
dependency direction.

## Scope of changes
Modify the packages above plus `user-docs/` and `docs/`. Create tests. Possibly remove
`internal/client/gen` (FR-017, decision D-R5). Do not edit `api/openapi.yaml` or `api/SOURCE`.

## Breaking changes
Exit codes: bare `afs` and `afs <service>` exit 252 (FR-015); malformed integers exit 252 (FR-009); same
page token twice and undecodable wait status exit 255 (FR-012, FR-013). Listed in
`user-docs/errors-and-exit-codes.md`. Safe-update notice text and timing change (FR-016).

## Success and acceptance criteria
All FR criteria met or deferred with written rationale; `make gate` and `make cover` green; coverage of
`internal/...` >= 90%; `safeupdate` and `exec` >= 90% from their own tests.

## Risks and mitigation
| Risk | Mitigation |
|---|---|
| A strict validator rejects real provider values | Strict sets follow the contract examples; tests with realistic values |
| Golden churn hides regressions | Goldens change only where a FR says so; review the diff |
| Exit-code changes break scripts | Release notes, user-docs |

## Timeline
One session. Order and commits in `plan.md`.

## References
Source PRD (this directory), `AGENTS.md`, parent spec `specs/archive/261121-cli-feature-complete/`.
