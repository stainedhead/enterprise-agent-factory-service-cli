# Spec: team orchestration agent and model access in the `afs` CLI

**Feature:** orchestration-and-model-access · **Created:** 2026-11-23 · **Release:** v1.1.0
**Source PRD:** [orchestration-and-model-access-PRD.md](orchestration-and-model-access-PRD.md) (focused copy); root PRD
`team-orchestration-agent-PRD.md` (ORC-*, MOD-*, OR-D*); API contract v1.4.0 section 4.4d.

## 1. Executive Summary

Every team now has exactly one orchestration agent that runs the team's collaboration tool (Paperclip or Buzz) and is
the only part of a team reachable from outside the cluster. Its Web UI and Web API each have an exposure (`off`,
`private`, `public`) within a team limit. Separately, every level (service, team, agent, pool) now chooses between
*all* models and a *selection*. This change gives the `afs` CLI the matching surface and keeps the 1.0.1 spellings
working, with a notice, until 2.0.0.

## 2. Problem Statement

A team-admin cannot script the choice of tool or the exposure of the orchestration agent, a service admin cannot see
which teams expose what, and nobody can choose which models an agent may call from the command line. The 1.0.1
collaboration commands and flags also describe a model that no longer exists (an optional agent with a role).

## 3. Goals / Non-Goals

**Goals:** ORC-29, ORC-30, MOD-8 (CLI part): every new operation has a command with help, example, columns, goldens and
an `AccessDenied` test; the deprecated commands and fields keep working with one stderr line; the new refusal codes print
a readable hint; the pinned contract is re-pinned and everything generated from it regenerated; release v1.1.0.

**Non-goals:** a client-side copy of any server rule (limits, tool support, model sets: the API decides); an
interactive confirmation for *change tool* (the console asks for the slug, the API does not, and a CLI call is one
API call); any removal (2.0.0); the SDK and console (separate repositories, same release).

## 4. Functional requirements

| FR | Requirement | PRD | Proven by |
|---|---|---|---|
| FR-001 | Re-pin `api/openapi.yaml` and `api/SOURCE` to `-restapi` f6d240c; regenerate `internal/spec`, client types and the user-docs reference; `afs --version` prints `contract/1.4`; 192 commands over 194 operations | ORC-29 | `buildinfo.TestVersionPrintsContractMinor`, `spec.TestContractVersionAndChecksum`, `spec.TestSnapshotCounts`, `services.TestFeatureCompleteStrict` |
| FR-002 | `orchestration describe-orchestration`: JSON unchanged, columns with N/A for a surface the tool lacks, stderr notices for lowered and N/A surfaces, the pending-agent hint | ORC-29, OR-CLI-26 | `orchestration.TestDescribeOrchestrationJSONUnchanged`, `orchestration.TestDescribeOrchestrationGoldens`, `orchestration.TestDescribeOrchestrationNotices`, `orchestration.TestDescribeOrchestrationPendingHint` |
| FR-003 | `orchestration put-orchestration`: `--web-ui-exposure`, `--web-api-exposure` (both required unless input JSON has them), `--registration-mode`, If-Match read, `--wait`, refusal reasons with hints | ORC-29, OR-CLI-26 | `orchestration.TestPutOrchestrationBody`, `orchestration.TestPutOrchestrationUsageErrors`, `orchestration.TestPutOrchestrationInputJSON`, `orchestration.TestPutOrchestrationReadsVersion`, `orchestration.TestPutOrchestrationWait`, `orchestration.TestPutOrchestrationRefusalReasons` |
| FR-004 | `orchestration change-collaboration-tool --tool paperclip\|buzz`: one POST with an idempotency key, `--wait`, refusal hints | ORC-29, OR-CLI-26 | `orchestration.TestChangeCollaborationTool`, `orchestration.TestChangeCollaborationToolUsageErrors`, `orchestration.TestChangeCollaborationToolRefusals`, `orchestration.TestChangeCollaborationToolWait` |
| FR-005 | `admin list-orchestration-endpoints`: filters, columns, CSV, backfill progress on stderr | ORC-29, OR-CLI-26 | `admin.TestListOrchestrationEndpointsFilters`, `admin.TestOrchestrationEndpointsTableGolden`, `admin.TestOrchestrationEndpointsMigrationNotice`, `admin.TestListOrchestrationEndpointsCSV` |
| FR-006 | `admin put-service-model-access --access all\|selected`; `update-model` refusal hint; `admin list-models` columns | MOD-8, OR-CLI-28 | `admin.TestPutServiceModelAccess`, `admin.TestUpdateModelRefusedWhileServiceAccessAll`, `admin.TestAdminModelsColumns` |
| FR-007 | `teams create-team --orchestration` (shorthand and JSON) and `--max-collaboration-exposure`; `update-team --max-collaboration-exposure`; team table columns | ORC-26, ORC-29 | `teams.TestCreateTeamOrchestrationShorthand`, `teams.TestCreateTeamOrchestrationShorthandErrors`, `teams.TestMaxCollaborationExposure`, `teams.TestCreateTeamOrchestrationRefusals` |
| FR-008 | `teams put-team-models --access`; `list-team-models` columns and access notice; `models list-models` columns | MOD-7, MOD-8 | `teams.TestPutTeamModelsAccess`, `teams.TestListTeamModelsAccessNotice` |
| FR-009 | `agents create-agent`/`update-agent --models all\|ids --default-model-id`; `--role` out of help; role filter open | MOD-8, ORC-30, OR-CLI-28 | `agents.TestCreateAgentModelsBody`, `agents.TestCreateAgentModelsUsageErrors`, `agents.TestUpdateAgentModels`, `agents.TestCreateAgentRoleHiddenAndRefusedByAPI`, `agents.TestListAgentsRoleOpenEnum` |
| FR-010 | The same model options on pools | MOD-8, OR-CLI-28 | `pools.TestPoolModelsBody`, `pools.TestPoolRefusals` |
| FR-011 | Deprecated operations (`describe-collaboration`, `put-collaboration`, `admin list-collaboration-endpoints`) and options (`--model-id`, `--may-expose-collaboration`) keep working with one stderr line and a help note; the service's `Deprecation` header prints one line when the CLI had no static notice | ORC-30, OR-CLI-27, OR-CLI-28 | `agents.TestModelIDDeprecated`, `pools.TestPoolModelIDDeprecated`, `teams.TestMayExposeCollaborationDeprecated`, `admin.TestListCollaborationEndpointsDeprecated`, `services.TestServiceDeprecationHeaderLine` |
| FR-012 | Hints for every new refusal reason (and the changed legacy ones), one clean line each, no server text | ORC-29 | `exec.TestCollaborationReasonHints`, `exec.TestReportErrorPrintsFieldsAndHints`, `agents.TestAgentOrchestrationRefusals` |
| FR-013 | `collaboration invoke-api` reads the Web API address of the effective level with one `GetTeamOrchestration`; a Web API that is off, or a tool without one, is exit 254 `Web API is off` | ORC-30, OR-CLI-27 | `collaboration.TestInvokeAPIWebAPIOff254`, `collaboration.TestInvokeAPIBuzzHasNoWebAPI254`, `collaboration.TestInvokeAPIOrchestrationReadFailures`, `collaboration.TestInvokeAPIRequestShape` |
| FR-014 | `admin put-setting` knows the `orchestration.*` and `models.*` keys; the renamed 1.0.1 keys are the API's refusal | ORC-27 | `admin.TestPutSettingOrchestrationKeys` |
| FR-015 | Swarm template skeleton shows `spec.orchestration` | ORC-7 | `swarmtemplates.TestSwarmTemplateSkeletonHasOrchestration` |
| FR-016 | Docs: README, user-docs, CHANGELOG 1.1.0, AGENTS.md status, the root skills entry | all | `docs.TestDocgenMatchesHelp` |

## 5. Acceptance criteria

Every FR names its tests above; `internal/cli/services` checks that each named test exists
(`TestOrchestrationCriteriaTraceToTests`). `make gate` is green with `internal/...` coverage at or above 90%.

## 6. Decisions

See [implementation-notes.md](implementation-notes.md) (D-O1 onward).
