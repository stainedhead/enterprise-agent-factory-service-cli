## Summary

Adds Paperclip team collaboration to the `afs` CLI (API contract v1.3.1, release **v1.0.1**): the `collaboration` service (describe, put with If-Match and `--wait`, links, connect, disconnect, rotate keys), the local command `afs collaboration invoke-api` for pipelines (AFS credential in `X-Afs-Authorization`, Paperclip key from a 0600 file or stdin, origin validation, no redirects or cookies, narrow retries, edge errors mapped to AWS-style messages), the collaboration agent role and managed secrets, the collaborator role, the collaboration report, swarm template `spec.collaboration`, and `admin list-collaboration-endpoints`. Credentials never reach stdout, stderr, debug output or process arguments (sentinel tests).

The contract snapshot is pinned to `-restapi` commit 3588074 (`api/SOURCE`); `afs --version` prints `contract/1.3`.

## Process

Implemented with the dev-flow `implm-frm-prd` workflow, 14 steps. Analysis: [paperclip-collaboration-analysis.md](paperclip-collaboration-analysis.md). Spec: `specs/archive/261122-paperclip-collaboration/` with the review-fixes spec nested inside it (4 findings, all fixed).

## Tests

`make gate` green; `internal/...` coverage 94.8% (floor 90%). `services.TestPaperclipCriteriaTraceToTests` checks that every criterion of CLI-US-24 to CLI-US-35 names an existing test.

## Parked (needs a real cluster or AWS)

| Item | Needs |
|---|---|
| End-to-end proofs against a deployed `afs-edge` (collaboration agent, connect agents, pipeline creates an issue) | Dev environment with the edge and `collab.base_domain` |
| Release signing, Homebrew tap | CI secrets, tap repository |

## Root contract follow-ups

See `implementation-notes.md`: RCF-1 (`admin list-agent-templates` has no capability filter), RCF-2 (`managedBy` in `managed_secret` details), RCF-3 and RCF-7 (root `skills/` entry), RCF-5 (`collab.*` keys missing from the `PutSetting` key enum).

## Release steps for J

1. After the `-restapi` PR merges, re-pin `api/SOURCE` to the merge commit (`api/openapi.yaml` is expected to be unchanged) in a follow-up.
2. Merge this PR, then push the tag: **`git tag v1.0.1 && git push origin v1.0.1`**. The Makefile and GoReleaser stamp the version from the tag.

`CHANGELOG.md` carries the v1.0.1 entry.
