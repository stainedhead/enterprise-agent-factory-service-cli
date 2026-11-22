# Spec: paperclip-collaboration-auto-review

**Source PRD:** [paperclip-collaboration-auto-review-PRD.md](paperclip-collaboration-auto-review-PRD.md) · **Parent:** `specs/archive/261122-paperclip-collaboration`

## Executive summary
Four small fixes from the automated review of the Paperclip change: make `invoke-api` fail fast on a `Retry-After` above the cap (FR-R1, P1), and close three test and diagnostics gaps (FR-R2 to FR-R4, P2).

## Requirements
| Id | Priority | Requirement |
|---|---|---|
| FR-R1 | P1 | Edge client does not retry when Retry-After exceeds the cap |
| FR-R2 | P2 | `runtime.Credentials` and `Transport` tests |
| FR-R3 | P2 | Parser test for unknown-option suggestion without echoing values |
| FR-R4 | P2 | Key file that is a directory gives exit 255 with a clear message |

## Non-goals
Anything in section 3 of the PRD (deferred items).

## Acceptance
Each requirement's criteria are in the PRD; each names its test.
