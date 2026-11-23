# Research: orchestration-and-model-access

## Inputs read

- Root `team-orchestration-agent-PRD.md` (ORC-1 to ORC-32, MOD-1 to MOD-8, OR-D1 to OR-D17), sections 6.6 (CLI), 6.9 (model access), 7.6 (CLI stories).
- Root `docs/contracts/api-contract.md` v1.4.0, section 4.4d, the changes to 4.2 (teams) and 4.4 (agents, pools, models), the error table of section 5.
- `-restapi` `api/openapi.yaml` at f6d240c (194 operations, five new: `GetTeamOrchestration`, `PutTeamOrchestration`,
  `ChangeCollaborationTool`, `ListOrchestrationEndpoints`, `PutServiceModelAccess`) and its `CHANGELOG.md` 1.1.0.
- Existing code: `internal/cli/collaboration`, `admin/collaboration.go`, `agents/collaboration.go`, the generic executor.

## What the snapshot changes for the CLI

| Area | Change in the snapshot | CLI consequence |
|---|---|---|
| Operations | 5 new, `x-afs-cli` names given | `orchestration` service (3), `admin list-orchestration-endpoints`, `admin put-service-model-access` |
| `Agent`, `Pool` | `models {access, modelIds, defaultModelId, effectiveModelIds}`; `modelId` deprecated and no longer required; `Agent.role` is `worker` or `orchestration`, read-only | `--models`, `--default-model-id`; `--model-id` warns; `--role` hidden |
| `CreateAgentRequest.role` | enum `worker` only, description `role_reserved` | hidden, open enum, API judges |
| `CreateTeamRequest` | `orchestration {tool, webUi, webApi, registrationMode, models, resources}`, `maxCollaborationExposure`, `modelAccess` | `--orchestration` shorthand, flags derived |
| `Team` | `orchestrationAgentId`, `maxCollaborationExposure`, `modelAccess`; `mayExposeCollaboration` deprecated | columns `Limit`, `Models`; deprecated flag warns |
| `TeamModels` / `Model` | `access`, `modelIds`; `callable`, `priceState`, `inferenceProfileState`, `origin`, `teamsWithAccess` | columns and an access notice |
| Deprecated operations | `GetTeamCollaboration`, `PutTeamCollaboration`, `ListCollaborationEndpoints` with `Deprecation`, `Sunset`, `Link` headers | help note, one stderr line, header-driven line for unknown cases |
| `PutSetting` key enum | lists only the 1.0.1 keys (and still `collab.base_domain`, `collab.public_endpoints`); the service (`internal/domain/settings_orchestration.go`) accepts the 14 new keys and refuses the two renamed ones with `setting_renamed` | the key is an open enum with the corrected list (D-O6) |
| `ListAgents ?role=` | enum is `worker`, `collaboration`; the contract text says `worker`, `orchestration` | open enum with the contract's values (D-O5) |

## Findings to report to `-restapi` (not blocking)

1. `api/openapi.yaml` `PutSetting` `key` enum lacks the 14 new keys and still lists the two renamed ones.
2. `api/openapi.yaml` `ListAgents` `role` enum lacks `orchestration`.
3. `CreateTeamRequest.orchestration` has no `modelId` shorthand and no `volume` (contract 4.2 lists both); the CLI sends `models.defaultModelId` for `ModelId`.
