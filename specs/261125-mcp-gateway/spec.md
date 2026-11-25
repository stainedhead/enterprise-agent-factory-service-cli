# Spec: governed MCP servers in the `afs` CLI

**Feature:** mcp-gateway · **Created:** 2026-11-25 · **Release:** v1.3.0
**Source PRD:** root `mcp-gateway-PRD.md` (decisions answered by J, 2026-11-25). **Contract:** API contract v1.6.0
(55 new operations, 54 of them CLI commands), OpenAPI at `-restapi` `ef92a4c82e028e812d9f485287ebe26d5605375a`.
Cross-repo design: the root PRD and the root contracts; `-restapi` `specs/261125-mcp-gateway/`.

## 1. Executive Summary

Release 1.3.0 lets agents reach MCP servers only through their team gateway, against a governed catalog: service admins
and security operators register, review, approve, suspend and retire servers; team admins switch them on, add their own
when allowed, and narrow them per agent; the gateway holds every credential. This change gives `afs` the matching
surface: the new `mcp-servers` service, the team's MCP commands in `teams`, MCP access and federated identities in
`agents` and `pools`, MCP events in `audit` and MCP usage in `reports`.

## 2. Problem Statement

Without CLI support nobody can script the catalog, its reviews, a team's switches and credentials, or pull the MCP
ledger and usage views into a pipeline.

## 3. Goals / Non-Goals

**Goals:** one command per new CLI operation, each with help, example, table columns, goldens and an `AccessDenied`
test; credential values and OAuth client secrets only from a file, standard input or a no-echo prompt; safe updates that
read the right version (a team's view carries the server's version in `server`); hints for the 23 new error reasons;
`--mcp` on agent and pool create and update; release v1.3.0 docs.

**Non-goals:** any client-side copy of a server rule (risk tiers, approval levels, host deny lists: the API decides);
`CompleteMcpOAuthConnect` (a browser redirect, not a CLI operation); the SDK and console.

## 4. Functional requirements

| FR | Requirement | PRD | Proven by |
|---|---|---|---|
| FR-001 | Re-pin `api/openapi.yaml` and `api/SOURCE` to `-restapi` ef92a4c (v1.6.0); regenerate `internal/spec` and client types; `afs --version` prints `contract/1.6`; 273 commands over 276 operations | MCP-CLI-1 | `buildinfo.TestVersionPrintsContractMinor`, `spec.TestContractVersionAndChecksum`, `services.TestFeatureCompleteStrict` |
| FR-002 | `mcp-servers` service: 24 commands for the catalog, review and query capture, with example, columns, goldens and `AccessDenied` | MCP-CLI-1 | `mcp.TestCommandsDoneMcpServers` |
| FR-003 | The team's 18 MCP commands in `teams` | MCP-CLI-1 | `mcp.TestCommandsDoneTeamMcp` |
| FR-004 | MCP access and federated identities of agents and pools, MCP events, MCP usage reports (12 commands) | MCP-CLI-1 | `mcp.TestCommandsDoneMcpAccess` |
| FR-005 | `list-mcp-servers` filters reach the query; a bad enum is exit 252 | MCP-CLI-1 | `mcp.TestListMcpServersFilters` |
| FR-006 | `create-mcp-server` sends the documented body; `--hosted` reads JSON from a file | MCP-CLI-1 | `mcp.TestCreateMcpServerBody` |
| FR-007 | `--secret-value` and `--client-secret` only from a file, standard input or a prompt; never echoed | MCP-CLI-2 | `mcp.TestMcpCredentialSecretsNeverInline` |
| FR-008 | Updates read the version and send `If-Match`; a team's switch, policies and credentials send none unless given | MCP-CLI-1 | `mcp.TestMcpUpdatesReadVersion`, `safeupdate.TestMcpReaders`, `safeupdate.TestReaderTableCoversAllIfMatchOperations` |
| FR-009 | Policy rules and review decisions are JSON | MCP-CLI-1 | `mcp.TestMcpPolicyAndReviewBodies` |
| FR-010 | Refusals print the AWS-style error, exit 254 | MCP-CLI-3 | `mcp.TestApproveRefusedByPolicy` |
| FR-011 | Effective MCP renders its servers; the usage report renders its rows and sends `view` and `window` | MCP-CLI-1 | `mcp.TestEffectiveAgentMcpAndReport` |
| FR-012 | Every v1.6.0 error reason has a one-line hint naming a command or what to change | MCP-CLI-3 | `exec.TestCollaborationReasonHints` |
| FR-013 | Milestones M15 (catalog) and M16 (team, access, reports); every operation has a milestone and stream | MCP-CLI-1 | `cli.TestEveryOperationHasMilestone`, `services.TestDoneUpToCurrentMilestone` |
| FR-014 | Reference docs regenerated; help goldens of changed commands updated | MCP-CLI-4 | `help.TestHelpGoldens` (and `tools/docgen` TestDocgenMatchesHelp) |

## 5. Acceptance

`make gate` green; every FR row names a test that exists (`services.TestMcpCriteriaTraceToTests`).
