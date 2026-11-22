# Implementation notes: paperclip-collaboration

**Date:** 2026-11-22. Updated as work proceeds; records decisions, edge cases and deviations.

## Technical decisions

Inherited decisions PC-CLI-D1..D12 are in the PRD section 7.1. Decisions taken while implementing are listed below.

| # | Decision | Why |
|---|---|---|
| D-1 | Pinned restapi branch commit 99bfcfc (info.version 1.3.1); `api/SOURCE` records `ref: claude/project-thread-xqpqky`; follow-up to re-pin to the merge commit | PR #7 unmerged; `Further restapi commits leave api/openapi.yaml unchanged` |
| D-2 | `output.Column` gained `Null` (text shown for a null or absent value) | PRD wants `-` for a null hostname, a missing Paperclip link state and empty capabilities; `Scalar(nil)` prints an empty cell |
| D-3 | `cli.Override.OmitEmptyBody` sends no body for an optional body with no member; used by `rotate-secret` (and `rotate-collaboration-keys`); `connect-agent` still sends `{}` | PRD: rotate sends no body, connect sends `{}` |
| D-4 | The `managed_secret` hint reads "managed by AFS or Paperclip" | The PRD text `<AFS\|Paperclip>` would need the value of `managedBy`, which the CLI never reads and a hint never contains server text |
| D-5 | "Not effective" notice: one stderr line per surface whose `effective` flag is false, when `effective.reasons` is non-empty, listing all reasons | Contract `effective.reasons[]` is not per surface; PRD says one line per surface |
| D-6 | `ListCollaborationEndpoints` table uses the top-level `phase` member | The contract's `CollaborationEndpoint` has `phase`, not `status.phase` (PRD wrote `Phase` only) |
| D-7 | `afs admin list-agent-templates` has no `--capability` option | Contract: only `ListAgentTemplates` accepts `?capability=` (PRD 28.7 over-asked); see root follow-ups |
| D-8 | Hints table sits in `internal/cli/exec/hints.go`; the field lines and hints print after `Request ID` for every API error, not only collaboration | One code path; unknown reasons print nothing |
| D-9 | `clitest.Options` gained `Now` and `Sleep` so a `--wait` timeout is proven with a virtual clock | A timeout test must not wait 10 minutes |
| D-10 | Redaction: members named `...Secret` whose value is a string not starting `arn:` are masked; `...Arn` members never | CLI-US-31.2: `collabClientSecretArn` is logged, `clientSecret` is not |

## Edge cases and solutions

(none yet)

## Deviations from plan

(none yet)

## Root contract follow-ups

| # | Follow-up |
|---|---|
| RCF-1 | PRD/contract: `admin list-agent-templates` has no `capability` filter in the contract (§4.4 names it for `ListAgentTemplates` only). The PRD story CLI-US-28.7 and the living PRD should drop `afs admin list-agent-templates --capability` or the contract should add it |
| RCF-2 | PRD/contract: the `managed_secret` hint cannot name the manager without reading `managedBy`; the contract could add `details.managedBy` so the hint can say AFS or Paperclip |
| RCF-3 | The root `skills/enterprise-agent-factory-service-cli.md` must gain the new commands (root repository is not edited here); the command list is in `user-docs/collaboration.md` and `user-docs/reference/collaboration.md` |
| RCF-4 | Re-pin `api/SOURCE` to the `-restapi` merge commit when PR #7 lands (api/openapi.yaml is expected to be unchanged, so only SOURCE changes) |

## Parked

| Item | Why parked | Needs |
|---|---|---|
| End-to-end proofs against a dev environment (create collaboration agent, connect two agents; pipeline creates a Paperclip issue; service admin sets base domain) | No real cluster, AWS account or deployed `afs-edge` | A dev environment with `afs-edge`, `collab.base_domain`, a team provider with a collaboration client |
| Release signing, Homebrew tap | Credentials and tap repo not in this repository (as in v1.0.0) | CI secrets, tap repository |

## Lessons learned

(at completion)

## Review of the spec (step 2)

| Dimension | Status | Note |
|---|---|---|
| Acceptance criteria | Pass | CLI-US-24..35 numbered, each with a named test |
| Technical approach | Pass | table-driven overrides plus one local command |
| Component coverage | Pass | spec.md section 6 and 7 name packages and files |
| Edge case handling | Pass (after fix) | spec.md section 13 added |
| Open questions | Pass (after fix) | spec.md section 14 added |
| Out-of-scope clarity | Pass | PRD section 2.1 |
| Status initialization | Pass | Phase 0 in progress |

Verdict: implementation-ready.
