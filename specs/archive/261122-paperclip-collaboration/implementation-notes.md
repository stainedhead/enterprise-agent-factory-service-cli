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
| D-11 | `cli.Override.OpenEnum` and `Option.OpenEnum`: the help says "Known values" and a value outside the pinned enum is sent for the API to judge; now used only by the report schedule `--report`; `put-setting --key` returned to a closed enum at 763e939 (RCF-5 resolved) | The pinned snapshot's `key` enum lacks the seven `collab.*` settings the contract lists (RCF-5); CLI-US-31.3 says the CLI carries no list of keys |
| D-12 | `--collab-client-secret-arn` is checked against a Secrets Manager ARN shape (any partition, 12-digit account) in a Build hook, covering `--cli-input-json` too; the error never repeats the value | CLI-US-31.2 |
| D-13 | Settings gained `CollabBaseDomain` (profile key `collab_base_domain`), read per profile; `afs configure` does not prompt for it | Optional key (PC-CLI-D5) |

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
| RCF-5 | Resolved: `-restapi` 763e939 lists the seven `collab.*` keys in the `PutSetting` enum; the CLI was re-pinned to it and `put-setting --key` is an ordinary (closed) enum again. Originally: `-restapi` `api/openapi.yaml` (99bfcfc): the `PutSetting`/`GetSettings` key enum lacks `collab.base_domain`, `collab.public_endpoints`, `collab.session_ttl_minutes`, `collab.max_body_mb`, `collab.content_scan`, `collab.allowed_source_cidrs`, `collab.denied_spike_threshold` although the root contract §4.13 lists them; the CLI treats `--key` as an open enum until the snapshot lists them |
| RCF-6 | Resolved: `-restapi` 3588074 lists the collaboration report kind and the `surface`, `reason` and `person` dimensions; the CLI was re-pinned to it and `--kind` is an ordinary (closed) enum again. No action left |
| RCF-7 | The root `skills/enterprise-agent-factory-service-cli.md` also needs `collaboration invoke-api` and the `--kind collaboration` report (extends RCF-3) |
| RCF-4 | Re-pin `api/SOURCE` to the `-restapi` merge commit when PR #7 lands (api/openapi.yaml is expected to be unchanged, so only SOURCE changes) |

## Parked

| Item | Why parked | Needs |
|---|---|---|
| End-to-end proofs against a dev environment (create collaboration agent, connect two agents; pipeline creates a Paperclip issue; service admin sets base domain) | No real cluster, AWS account or deployed `afs-edge` | A dev environment with `afs-edge`, `collab.base_domain`, a team provider with a collaboration client |
| Release signing, Homebrew tap | Credentials and tap repo not in this repository (as in v1.0.0) | CI secrets, tap repository |

## Lessons learned

- Open enums keep the CLI free of server lists but need a test that proves the value reaches the API.
- A generic checker (traceability test over PRD criteria) caught seven PRD test names that did not exist.
- Do not overwrite an existing test file: check `git status` before `cat >`.

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

## M9 notes

- D-14 Edge package. `invoke-api` lives in `internal/edge` (origin, key, request) plus `internal/cli/collaboration/invoke.go`; the edge client uses `RoundTrip` directly, so no redirect is followed and no cookie jar exists. It sends `Accept-Encoding: identity` and streams the body with an idle read timeout.
- D-15 Envelope test. A response is the edge's section 6 envelope only when it is JSON with `code`, `message` and a request id (body or `X-Request-Id`); Paperclip's own `{"code": ...}` bodies stay Paperclip's (HttpError, body on stdout). Error bodies are held up to 1 MiB for that test.
- D-16 Loopback exemption. An `http://` loopback origin is exempt from the port, address and `collab_base_domain` rules (as `endpoint_url` is exempt from https), so tests and local stubs work. Anything else is strict.
- D-17 Local problems. An unreadable `--body` file, a key file wider than 0600, a missing credential, a refused origin and a null `restApiUrl` exit 255; bad shapes exit 252.
- D-18 `--wait` is declared as an option only so it can be refused with the PC-CLI-D7 sentence instead of a generic unknown-option error.
- D-19 Unknown options echo only the option name (never `=value`) and suggest options that start with it, so `--paperclip-key=<value>` cannot reach stderr and points at `--paperclip-key-file`.
- D-20 Reports. After the re-pin to 3588074 the snapshot carries the collaboration kind and dimensions; `--kind` is a normal enum and the client-side `--group-by` shape check allows `surface`, `reason` and `person` (CLI-US-32.2). `create-report-schedule --report` is an open enum so `--report collaboration` reaches the API, which refuses it (PC-D13).
- D-21 `Override.SkeletonMembers` replaces a body member's placeholder in `--generate-cli-skeleton`; swarm template create and update show `spec.collaboration.templateId` this way (the generic skeleton prints `"spec": {}` because the snapshot loader keeps no nested properties).
- D-22 PRD test names. Seven names in the PRD's `Proven by` lines pointed at tests that did not exist under that name; the PRD now names the existing equivalents (`exec.TestCollaborationReasonHints`, `services.TestCommandTableMatchesSnapshot`, `services.TestDoneUpToCurrentMilestone`, `help.TestOperationHelpShowsOptionsPermissionAndExample`, `api.TestSnapshotChecksumMatchesSource`, `docs.TestDocgenMatchesHelp`) and `services.TestPaperclipCriteriaTraceToTests` fails when a named test is missing.
