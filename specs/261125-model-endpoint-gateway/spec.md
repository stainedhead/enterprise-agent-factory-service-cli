# Spec: model endpoint gateway in the `afs` CLI

**Feature:** model-endpoint-gateway · **Created:** 2026-11-25 · **Release:** v1.2.0
**Source PRD:** root `model-endpoint-gateway-PRD.md` (approved by J, 2026-11-25), Appendix D.4 (MEG-CLI-1 to
MEG-CLI-4), stories US-70 to US-81, test cases TC-CLI-01 to TC-CLI-05. **Contract:** API contract v1.5.0 section 4.13,
OpenAPI at `-restapi` `7d48ea70d2f9cb77ef4308b0225ca9d4168dea60`. The CLI's own stories are section 4.2 of the
[CLI PRD](../../enterprise-agent-factory-service-cli-PRD.md).

## 1. Executive Summary

Release 1.2.0 gives every team a model gateway: agents call OpenAI-compatible routes (service endpoints a team opts in
to, and the team's own endpoints) through AFS, metered and capped like Bedrock, with the upstream credential held by
AFS. Service admins switch whole providers on and off, Bedrock included. This change gives the `afs` CLI the matching
surface: the new `model-endpoints` service, two `teams` commands, two inference reports, and route options on the agent
and pool model commands.

## 2. Problem Statement

Without CLI support a service admin cannot script the endpoint catalog, prices, credentials or provider switches, a
team-admin cannot script route selection or team endpoints, and nobody can pull the inference report into a pipeline.

## 3. Goals / Non-Goals

**Goals:** MEG-CLI-1 to MEG-CLI-4: one command per new operation (27 operations in the v1.5.0 snapshot), each with help,
example, table columns, goldens and an `AccessDenied` test; credential values only from a file, standard input or a
no-echo prompt, redacted under `--debug`; the provider statement import takes a file (CSV or JSON); `--dry-run` on
`update-inference-provider`; route options on agents and pools; readable hints for the new error reasons; release
v1.2.0 docs.

**Non-goals:** any client-side copy of a server rule (provider hosts, slug reservations, the Bedrock fallback guard:
the API decides); the MCP operations of contract v1.6.0 (release 1.3.0, another change); the SDK and console.

## 4. Functional requirements

| FR | Requirement | PRD | Proven by |
|---|---|---|---|
| FR-001 | Re-pin `api/openapi.yaml` and `api/SOURCE` to `-restapi` 7d48ea7; regenerate `internal/spec` and client types; `afs --version` prints `contract/1.5`; 219 commands over 221 operations | MEG-CLI-1 | `buildinfo.TestVersionPrintsContractMinor`, `spec.TestContractVersionAndChecksum`, `services.TestFeatureCompleteStrict` |
| FR-002 | `model-endpoints` service: a command for each of the 25 operations of section 4.13 it owns, with example, columns, goldens and `AccessDenied` | MEG-CLI-1, MEG-CLI-4 | `modelendpoints.TestCommandsDoneModelEndpoints`, `modelendpoints.TestTeamCommandsDoneModelEndpoints` |
| FR-003 | `list-model-endpoints` filters reach the query; `create-model-endpoint` sends the documented body; `--protocols` is a list | MEG-CLI-1 | `modelendpoints.TestListModelEndpointsFilters`, `modelendpoints.TestCreateModelEndpointBody` |
| FR-004 | `update-model-endpoint`, `delete-model-endpoint`, `update-team-model-endpoint`, `delete-team-model-endpoint` read the version and send `If-Match` | MEG-CLI-1 | `modelendpoints.TestUpdateModelEndpointReadsVersion`, `modelendpoints.TestTeamEndpointWritesReadVersionFromList` |
| FR-005 | `put-model-endpoint-credential`, `rotate-model-endpoint-credential`, `put-team-model-endpoint-credential` take the value only from `file://`, `fileb://`, `-` or a no-echo prompt; a plain argument and `--value` are refused with exit 252 before any call | MEG-CLI-3, TC-CLI-02 | `modelendpoints.TestCredentialValueOnlyFromFileStdinOrPrompt`, `modelendpoints.TestCredentialPlainArgumentRefused` |
| FR-006 | Under `--debug` a credential call logs neither the value nor `Authorization` | MEG-CLI-3, TC-CLI-04 | `modelendpoints.TestCredentialDebugRedacted` |
| FR-007 | Credential writes send no `If-Match` unless `--if-match` is given (the credential's version is not readable) | MEG-CLI-3 | `modelendpoints.TestCredentialSendsNoVersionRead` |
| FR-008 | `list-inference-providers`; `update-inference-provider --provider bedrock --no-enabled --fallback-model azure-openai/gpt-5.1 --select-for-all-teams --reason …` sends the documented body with the version read from the list; `--dry-run` sends `dryRun=true` and prints the impact on stderr | TC-CLI-05 | `modelendpoints.TestUpdateInferenceProviderBedrockOff`, `modelendpoints.TestUpdateInferenceProviderDryRun`, `safeupdate.TestInferenceProviderVersionFromList` |
| FR-009 | `import-provider-statement --statement file://PATH`: a CSV file is sent as `text/csv` with `period` and `source` in the query; a JSON file is the body; anything but a file is refused | MEG-CLI-1, US-79 | `modelendpoints.TestImportStatementCSV`, `modelendpoints.TestImportStatementJSON`, `modelendpoints.TestImportStatementUsageErrors` |
| FR-010 | `teams put-team-gateway-settings` and `teams put-team-model-endpoints` with example, columns, goldens, `AccessDenied` | MEG-CLI-2 | `teams.TestCommandsDoneGateway`, `teams.TestPutTeamModelEndpointsBody` |
| FR-011 | `reports get-inference-report` and `get-team-inference-report`: query filters, `--range` exclusive with `--from`/`--to`, group-by checked, CSV | MEG-CLI-1, US-73 | `reports.TestCommandsDoneInference`, `reports.TestInferenceReportQuery`, `reports.TestInferenceReportUsageErrors` |
| FR-012 | Agents and pools: `--routes all\|<slug>…` and `--default-route` fold into `models`; `--default-model-id <route>/<model>` is sent unchanged | MEG-CLI-2 | `modelaccess.TestRoutesFoldIntoModels`, `agents.TestCreateAgentRoutes`, `pools.TestPoolRoutes` |
| FR-013 | Usage and cost reports accept `--route` and `--group-by route` | MEG-CLI-2 | `reports.TestRouteFilterAndGroupBy` |
| FR-014 | Hints for the v1.5.0 error reasons | MEG-CLI-1 | `exec.TestModelEndpointHints` |
| FR-015 | Reference docs regenerated; usage examples in `user-docs/`; CHANGELOG 1.2.0; README, AGENTS.md and the CLI PRD updated | MEG-CLI-4 | `docs.TestReferenceUpToDate` |

## 5. Acceptance

Every row above names tests that exist (`services.TestModelEndpointCriteriaTraceToTests`); `make gate` is green and
`internal/...` coverage stays at or above 90%.
