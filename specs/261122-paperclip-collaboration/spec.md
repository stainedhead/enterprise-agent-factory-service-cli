# Spec: Paperclip team collaboration in the `afs` CLI

**Feature:** paperclip-collaboration · **Created:** 2026-11-22 · **Release:** v1.0.1
**Source PRD:** [specs/261122-paperclip-collaboration/paperclip-collaboration-PRD.md](paperclip-collaboration-PRD.md) (focused copy);
full PRD: [specs/archive/261121-cli-feature-complete/enterprise-agent-factory-service-cli-PRD.md](../archive/261121-cli-feature-complete/enterprise-agent-factory-service-cli-PRD.md)

## 1. Executive Summary

Every AFS team may run one collaboration agent (Paperclip), connect its workers to it, and publish its web UI and
REST API through a governed public endpoint (API contract v1.3.1). This change gives the `afs` CLI the matching
surface: a `collaboration` service (six operations), `admin list-collaboration-endpoints`, the flag, column and hint
changes on `agents`, `secrets`, `teams`, `members`, `api-keys`, `admin`, `reports` and `swarm-templates`, and one local
command, `afs collaboration invoke-api`, that calls a team's Paperclip REST API with an AFS credential in
`X-Afs-Authorization` and a Paperclip key in `Authorization`. The CLI remains a thin client: the API and the edge
decide access.

## 2. Problem Statement

A team-admin cannot script the collaboration settings, and an automation author has no supported way to call the
team's Paperclip API with an AFS credential beside a Paperclip key. Without CLI support, pipelines hand-build `curl`
calls and risk leaking credentials.

## 3. Goals / Non-Goals

**Goals:** CLI-19 to CLI-27 (see PRD section 2): every §4.4b operation has a command with help, example, columns, goldens
and an `AccessDenied` test; `invoke-api` creates a Paperclip issue from a pipeline with no credential visible in any
output; the pinned contract is re-pinned and everything generated from it regenerated; release v1.0.1.

**Non-goals:** PRD section 2.1 (no browser or web UI sign-in, no WebSocket/SSE client, no wrappers for Paperclip's own
objects, no Paperclip key management, no client-side secret scanning, no client-side copies of server rules, no
`create-collaboration`, no automatic connection of existing agents).

## 4. Functional requirements

| FR | Requirement | PRD | Milestone |
|---|---|---|---|
| FR-001 | Re-pin `api/openapi.yaml`/`api/SOURCE` to restapi 99bfcfc; regenerate spec, client types, docs; `afs --version` prints `contract/1.3` | CLI-27, US-34.1, 34.4 | M7 |
| FR-002 | `collaboration describe-collaboration`, JSON unchanged, table/text columns, not-effective notice, missing-agent hint | CLI-19, US-24.1, 24.2, 24.9 | M7 |
| FR-003 | `collaboration put-collaboration`: body, both surfaces required, If-Match, `--wait` (generation-aware), validation fields and hints | CLI-20, US-24.3-24.8, 24.10 | M7 |
| FR-004 | `collaboration list-collaboration-links`, `connect-agent`, `disconnect-agent` | CLI-19, US-26 | M7 |
| FR-005 | `agents`: `--role`, `--collaboration-connect`, `--paperclip-agent-id`, role filter, template capability, columns, skeleton, hints | CLI-21, US-28 | M7 |
| FR-006 | `secrets`: `--managed-by`, `ManagedBy` column, hints for `managed_secret`, `master_key_not_rotatable`, `reserved_name` | CLI-22, US-29 | M7 |
| FR-007 | `teams --may-expose-collaboration`, `hostname_held` hint, `members`/`api-keys` role `team_collaborator` | CLI-26, US-31.1, US-32.1 | M7 |
| FR-008 | Redaction: `X-Afs-Authorization`, `Cookie`, `Set-Cookie`, `Proxy-Authorization`; `__Host-afs-collab` never printed | CLI-25, US-33.1 | M7 |
| FR-009 | Reason hints table (one table in `internal/cli/exec`, no server text) | PC-CLI-D6, section 5.2 | M7 |
| FR-010 | `admin list-collaboration-endpoints` with filters, columns, CSV | CLI-23, US-30 | M8 |
| FR-011 | `admin create/update-identity-provider --collab-client-id --collab-client-secret-arn` (ARN shape validated); `collab.*` settings pass through | CLI-26, US-31.2-31.4 | M8 |
| FR-012 | `collaboration invoke-api` local command with origin safety, credential chain, key handling, edge error mapping, redaction | CLI-24, CLI-25, US-25, US-33.2 | M9 |
| FR-013 | `collaboration rotate-collaboration-keys` with notices and no secret values | CLI-19, US-27 | M9 |
| FR-014 | `swarm-templates` `spec.collaboration` passthrough and goldens; pools gain no collaboration option | US-35 | M10 |
| FR-015 | Reports: kind `collaboration`, `--group-by` additions | CLI-26, US-32.2, 32.3 | M10 |
| FR-016 | Docs (`user-docs/collaboration.md`, reference regenerated), traceability test over CLI-US-24..34, CHANGELOG, release notes | US-34.5 | M10 |

## 5. Non-functional requirements

See PRD section 4: performance (at most one extra read; streaming body), security (credential scope, TLS, no redirects,
no cookies, 0600 key file), reliability (retries, timeouts, Ctrl-C, non-JSON bodies), compatibility (1.2 CLI against 1.3
server and the reverse), observability (`--debug` events with component `collab`), quality (>= 90% coverage, goldens).

## 6. System architecture

Layers follow `AGENTS.md`: `cmd/afs` -> `internal/cli` -> `internal/client`, `internal/config`, `internal/output`.
New: `internal/cli/collaboration` (service package, one registration file per concern), `internal/cli/exec/hints.go`
(reason hint table), `internal/redact` (edge headers), `internal/collab` is **not** created: the invoke-api origin and
request logic lives in `internal/cli/collaboration` because it is a command handler plus a small transport helper.
Changed: `agents`, `secrets`, `teams`, `members`, `apikeys`, `admin`, `reports`, `swarmtemplates`, `safeupdate`, `wait`,
`services`, `milestones*`, `client` (debug), `config` (profile key `collab_base_domain`).

## 7. Scope of changes

Create: `internal/cli/collaboration/*`, `internal/cli/exec/hints.go`, `user-docs/collaboration.md`, `CHANGELOG.md`,
goldens under each package's `testdata/`. Modify: the packages above, `api/*`, generated files, `user-docs/*`.
Do not edit the root repository (its `skills/` entry is a recorded follow-up).

## 8. Breaking changes

None to existing commands. Behaviour changes: `reports --group-by` no longer rejects `surface|reason|person`
client-side; `afs --version` prints `contract/1.3`; the agent table gains two columns, templates and secrets one each.

## 9. Success and acceptance criteria

Every criterion of CLI-US-24 to CLI-US-35 in the PRD names a passing test (traceability test, M10). Quality gate green:
`make gate` (tidy, contract-check, generate-check, fmt-check, vet, lint, race tests, cover >= 90%, build, `--version`).

## 10. Risks and mitigation

| Risk | Mitigation |
|---|---|
| Contract pin is an unmerged restapi commit | Pin recorded in `api/SOURCE`; follow-up to re-pin SOURCE to the merge commit |
| Credential leak through `invoke-api` | Sentinel tests over stdout, stderr, debug, args; origin validation; no redirects/cookies |
| Hostile `restApiUrl` | PC-CLI-D5 rules plus optional `collab_base_domain` |
| Generated surface differs from PRD wording | Contract wins; record in implementation-notes.md |

## 11. Timeline and milestones

M7 (C1), M8 (C2), M9 (C3), M10 (C4); each ends with a commit, a push and `CurrentMilestone` raised in its exit task.

## 12. References

PRD (this directory); API contract v1.3.1; change PRD (root); `docs/solution-design.md` §11 (root).
