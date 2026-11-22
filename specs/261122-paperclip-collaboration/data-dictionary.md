# Data dictionary: paperclip-collaboration

**Date:** 2026-11-22. Types come from `internal/client/gen/types.gen.go` (generated); this file lists what the CLI relies on.

## Entities (wire)

| Name | Members used by the CLI |
|---|---|
| TeamCollaboration | `hostname`, `webUiEnabled`, `restApiEnabled`, `registrationMode`, `effective{webUi,restApi,reasons[]}`, `restApiUrl`, `status{phase,observedGeneration,conditions}`, `generation`, `version` |
| CollaborationLink | `agentId`, `agentName`, `connect`, `paperclipAgentId`, `linkState`, `lastSyncedAt`, `lastError` |
| RotateCollaborationKeysResponse | `rotated[{name,secretId,versionId}]`, `signsOutEveryone`, `agentRolloutPending` |
| CollaborationEndpoint | `slug`, `teamState`, `hostname`, `effective`, `status.phase`, `lastAccessAt` |
| Agent (extended) | `role`, `collaboration{connect,paperclipAgentId,linkState,lastSyncedAt,lastError}` |
| Secret (extended) | `managedBy` |

## Enumerations

`registrationMode`: manual, managed. `linkState`: pending, linked, failed, disabled. `role` (agent): worker, collaboration.
`managedBy`: user, afs, paperclip. Endpoint `state`: Ready, Degraded, Failed, Pending. Member roles: team_admin,
team_collaborator, team_viewer. Report kind: collaboration with group-by surface, reason, person.

## Local value objects (CLI)

| Name | Package | Meaning |
|---|---|---|
| Hint table | `exec` | `details.reason` / `fields[].reason` -> fixed hint text |
| Origin | `collaboration` | validated scheme+host of `restApiUrl` or `--collab-url` |
| PaperclipKey | `collaboration` | trimmed key read from a 0600 file or stdin; never logged |
| Profile key `collab_base_domain` | `config` | optional host pin |

## Request/response types of the local command

`invoke-api` flags: `--team-id --method --path --body --content-type --paperclip-key-file --collab-url --include`.
Output: response body to stdout unchanged; exit 0/254/255/252 per CLI-US-25.
