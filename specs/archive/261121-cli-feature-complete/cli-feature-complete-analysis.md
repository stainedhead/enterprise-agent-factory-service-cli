# Dev-Flow Process Analysis

**Feature:** cli-feature-complete
**Spec directory:** specs/archive/261121-cli-feature-complete
**Report generated:** 2026-11-22

---

## 1. Executive Summary

Built `afs`, the AWS-CLI-style command line for the Enterprise Agent Factory Service, to feature completeness: one command per API operation (180), the local commands (configure, configure sso, sso login/logout, completion, help), AWS-style global options, profiles, SSO device grant, API-key and workload-token credentials, five output formats, JMESPath `--query`, pagination, safe updates, idempotent retries, `--wait` and generated help. The command table and HTTP descriptors are generated from a pinned, checksummed snapshot of `-restapi`'s `api/openapi.yaml`. A code and design review then produced 20 findings (1 P0, 3 P1, 16 P2), all handled in the same run.

**Total runtime:** 3 h 24 min by git (first branch commit 2026-11-21T22:11:44Z to a47c81f at 2026-11-22T01:35:27Z, before this report); 3 h 34 min from the dashboard Process Start (2026-11-21T22:01:20Z).
**Overall assessment:** Fast and complete. All 20 review findings were handled (19 fixed in full, FR-010 criterion 2 deferred with rationale D-R9). `make gate` and `make cover` are green (94.7%, minimum 90%). `make vuln` could not run in the sandbox. Open items for the owner remain (re-pin, release and signing hand-offs).

---

## 2. Step-by-Step Timing

Git timestamps are authoritative. Steps 1-2 and 4-6 use the dashboard where git agrees; steps 3 and 7-11 are bounded by commits.

| Step | Name | Start | End | Runtime (min) | Key Outputs |
|---|---|---|---|---|---|
| 1 | Create Spec from PRD | 22:01:20Z | 22:11:44Z | 10 | spec directory, 8 spec files |
| 2 | Review Spec | 22:11:44Z | 22:13:33Z | 2 | spec review fixes |
| 3 | Implement Product (M0-M6) | 22:13:33Z (first code commit 23:25:42Z) | 00:44:21Z | 150 | 180 operations, generator, executor, SSO, output, docs |
| 4 | Documentation and User Docs | 00:44:21Z | 00:52:25Z | 8 | product, technical, user docs, ADR-007..011 |
| 5 | Code and Design Review | 00:47:14Z | 00:54:52Z | 8 | review PRD with 20 findings |
| 6 | Prepare Review PRD | 00:54:52Z | 00:55:54Z | 1 | PRD finalized |
| 7 | Archive Original Spec | 00:56:03Z | 00:57:02Z | 1 | specs/archive/261121-cli-feature-complete |
| 8 | Spec Review Fixes (create-spec on review PRD) | 00:57:02Z | 00:57:30Z | 0.5 | review spec |
| 9 | Implement Review Fixes | 00:57:30Z | 01:34:00Z | 37 | 8 fix commits, ADR-012..014 |
| 10 | Archive Fixes Spec | 01:34:00Z | 01:34:24Z | 0.5 | nested archive |
| 11 | Final Quality Pass | 01:34:24Z | 01:35:27Z | 1 | gate and cover green |
| 12 | Process Analysis Report | 01:35:27Z | see commit list | n/a | this report |
| 13 | Archive Spec (verify layout) | | | | layout verified |
| 14 | Open Pull Request | | | | PR into main |

**Notable observations:**
- Step 3 dominates (about 70 percent of the run). Git shows no commits between 22:13Z and 23:25Z: that hour was design and planning, not committed work.
- M0 (foundation) took 32 min of commits; M1-M6 took 39 min together, because the generic executor and generated descriptors made each service mostly overrides plus tests.
- The dashboard step 3 start (22:30Z) and step 10/11 times were estimates; they are corrected here to git boundaries. Step 11 originally showed an end earlier than its start.

---

## 3. Commit and Push Summary

**Total commits at the time of the report:** 41 (plus this report)

| Commit | Timestamp | Message |
|---|---|---|
| 8e2ff07 | 2026-11-21T22:11:44 | spec: create specs/261121-cli-feature-complete from the CLI PRD (dev-flow step 1) |
| 2366670 | 2026-11-21T22:13:33 | spec: review-spec fixes for 261121-cli-feature-complete (dev-flow step 2) |
| b371c98 | 2026-11-21T23:25:42 | M0-K1: pin OpenAPI snapshot with provenance and checksum gate |
| 8f034f0 | 2026-11-21T23:28:45 | M0-K2, K3: afsgen operation descriptors and oapi-codegen types from the pinned snapshot |
| 17f3ea3 | 2026-11-21T23:32:40 | M0-K4, K5: registry, parser, global options, services skeleton, milestone map; status |
| 0f2004b | 2026-11-21T23:35:00 | M0-T1..T4: HTTP transport with retries, idempotency keys, timeouts, debug log and redaction |
| fe63ef6 | 2026-11-21T23:36:59 | M0-C1, C2: config and credentials files, precedence, credential chain |
| 003e752 | 2026-11-21T23:38:14 | M0-S1..S3: SSO token cache, device authorization grant, silent refresh |
| 9e536eb | 2026-11-21T23:40:05 | M0-O1..O4: output renderers (json, yaml, table, text, csv), JMESPath query, pager, color |
| 905b5fa | 2026-11-21T23:43:14 | M0-E1 (part): request builder, secretinput |
| 78176a7 | 2026-11-21T23:46:53 | M0-E1..E6: executor, pagination, safe update, wait, errors, help |
| 965568b | 2026-11-21T23:48:29 | M0-S4: sso login/logout, get-caller-identity, list-providers, configure, configure sso |
| 4477155 | 2026-11-21T23:50:27 | M0-E7: contract replay test for all 180 operations |
| 47bfe24 | 2026-11-21T23:55:55 | M0-X1: globals tests, startup test, help goldens, coverage gate, CI contract checks |
| 61f1d9a | 2026-11-21T23:57:33 | M0-X2: docs, generated command reference, user-docs, ADRs, status |
| 3ab776c | 2026-11-22T00:05:37 | M1: teams, members, team audit, settings, identity providers; shared clitest harness (specs/261121-cli-feature-complete) |
| c8994ca | 2026-11-22T00:11:23 | M2: agents, secrets, config values, external systems, models, admin catalogs (specs/261121-cli-feature-complete) |
| 6d4e106 | 2026-11-22T00:15:52 | M3: budgets, usage, reports and CSV, notification channels, admin models and prices (specs/261121-cli-feature-complete) |
| 2e1b7b8 | 2026-11-22T00:18:29 | Table output never truncates id columns; goldens regenerated (specs/261121-cli-feature-complete) |
| 82b6764 | 2026-11-22T00:28:53 | M4: API keys, role links, identities, tags, log destinations, restore, schedules, service admin, exports, completion, release config (specs/261121-cli-feature-complete) |
| 557a184 | 2026-11-22T00:35:05 | M5: pools, stages, ephemeral teams, kubeconfig, cost tag status (specs/261121-cli-feature-complete) |
| a99706d | 2026-11-22T00:44:21 | M6: work orders, workload token, swarm templates, directory, service keys, security operators; feature complete (specs/261121-cli-feature-complete) |
| 5393e0f | 2026-11-22T00:47:14 | Step 4: product, technical and user docs match what was built (specs/261121-cli-feature-complete) |
| 45dfe46 | 2026-11-22T00:52:25 | Step 4: correct the secret example and the safe-update notice wording (specs/261121-cli-feature-complete) |
| 2afe66b | 2026-11-22T00:54:52 | Step 5: code and design review findings as a PRD (specs/261121-cli-feature-complete) |
| 6684168 | 2026-11-22T00:55:54 | Step 6: review PRD checked with review-prd and finalized (specs/261121-cli-feature-complete) |
| 171db5c | 2026-11-22T00:56:03 | Dashboard: correct step 6 runtime (specs/261121-cli-feature-complete) |
| 05cb5ac | 2026-11-22T00:57:02 | Step 7: archive the original spec (specs/archive/261121-cli-feature-complete) |
| 7d75b19 | 2026-11-22T00:57:30 | Step 8: create the review-fix spec (specs/261122-cli-feature-complete-auto-review) |
| 495b126 | 2026-11-22T01:00:08 | FR-001: reject INI injection and invalid SSO provider fields |
| ad8572c | 2026-11-22T01:02:38 | FR-002: require https and a matching issuer for OIDC discovery and endpoints |
| 348e828 | 2026-11-22T01:06:16 | FR-003: sanitize terminal control characters in table, text, errors and prompts |
| 7d055a5 | 2026-11-22T01:08:09 | FR-004: print NextToken on stderr for CSV with --no-paginate; resume without header |
| cfdc805 | 2026-11-22T01:12:59 | FR-006,FR-007,FR-008: IdP timeouts, access-token validation, --debug for OIDC traffic |
| ed8d854 | 2026-11-22T01:17:21 | FR-005,FR-009,FR-010,FR-011: endpoint userinfo refused and redacted, strict number options, dot path segments refused, Retry-After cap |
| 74e9079 | 2026-11-22T01:20:39 | FR-012,FR-013,FR-016: pagination loop guard, --wait Retry-After and undecodable status, notice after the write |
| 5945d85 | 2026-11-22T01:23:26 | FR-018,FR-020: sensitive option help lists safe forms and is not required; pager runs under the command context |
| e5a6740 | 2026-11-22T01:28:15 | FR-014,FR-015,FR-017: kubeconfig merge validation, exit 252 for a missing operation, compile-time contract references |
| 5c8a139 | 2026-11-22T01:34:00 | FR-019: close test gaps; ambiguous list match is an error; temp files removed after a failed save |
| 174b65e | 2026-11-22T01:34:24 | Step 10: archive the review-fix spec inside the original spec's archive |
| a47c81f | 2026-11-22T01:35:27 | Step 11: final quality pass, stable (gate and cover green, coverage 94.7%) |

Branch `feat/cli-feature-complete` was pushed after each step group. A pull request into `main` is opened in Step 14 (not merged); see the PR for its number.

---

## 4. Spec vs. Implementation Comparison

| Phase | Planned (spec) | Actual (git log) | Difference | Notes |
|---|---|---|---|---|
| Foundation (M0) | milestone M0, 28 tasks | 23:25Z-23:57Z (32 min) | n/a | spec gives sizes, not durations |
| Service milestones (M1-M6) | six milestones | 00:05Z-00:44Z (39 min) | n/a | all 180 operations, contract replay test |
| Documentation | step 4 | 8 min | n/a | two commits |
| Review | step 5-6 | 9 min | n/a | 20 findings |
| Review fixes | FR-001 first, P1s, then P2s | 37 min, 12 commits | n/a | order followed |
| Final pass | step 11 | 1 min | n/a | gate and cover green |

**Phases skipped:** none. Open for the owner, not implementable here: M6-X2 (re-pin to a `-restapi` release tag), M6-X3 hand-offs (CLI-CD-1..6, signing material, Homebrew tap `stainedhead/homebrew-afs`).
**Phases added:** the review-fix spec (steps 5-10) and ADR-012..014. Deferred: FR-010 criterion 2 (path id pattern check) because the pinned patterns reject the abbreviated ids used across about 60 tests (D-R9). FR-017 took the compile-time-reference option (ADR-014) because the sandbox refused `git rm -r` on the generated package.

---

## 5. Token / Message Usage

Exact token counts unavailable. Estimate by step complexity: the orchestrator ran as one long session that was compacted once (during Step 11-12) and resumed; implementation, review and fixes ran in that session and in sub-agents. Step 3 was the heaviest consumer, step 9 the second. No figure is claimed beyond that.

---

## 6. Process Observations

### What worked well
- Generating descriptors and the command table from a pinned, checksummed contract made 180 operations tractable and gave a replay test covering all of them.
- The review step found a reproduced P0 (config-file injection) and the TDD fix order (P0 alone, then P1, then P2) kept each commit small and gate-green.
- The new tests exposed two real bugs: ambiguous list matching in safe update and temp-file leaks on failed saves.

### What caused delays or rework
- FR-010 criterion 2 was attempted, broke about 60 tests, and was reverted and deferred.
- A shared scratchpad caused one output file to be overwritten by sibling agents; golangci-lint also reported a parallel run several times. Both were worked around with a private output directory and retries.
- A red-phase stash made silent-IdP tests hang; stopped and restored without loss.
- Long `make gate` runs exceed the 120 s tool limit and had to run in the background and be polled.
- The dashboard times for steps 3 and 7-11 were first guessed, then corrected to git.
- `make vuln` is unavailable in the sandbox and must be checked in CI.

### Recommendations for future runs
- Record step times from git at each step end, not from memory.
- Give each agent a private output directory.
- Make the spec's test fixtures use contract-valid ids from the start, so pattern checks can be enforced.
- Run `make vuln` in CI on this branch before merge.

---

## 7. Manual vs. Automated Comparison

**Estimated manual duration:** about 4-6 weeks for one senior developer (180 operations with generated plumbing, SSO device grant, five renderers, pagination, safe update, docs, a structured review and 20 fixes), assuming no meetings and excluding code review turnaround.
**Actual automated runtime:** 3 h 24 min from first to last commit before this report.
**Efficiency gain:** roughly 10-15x on elapsed time. This is a rough estimate; the human review and the open owner items are not included in the automated figure.
