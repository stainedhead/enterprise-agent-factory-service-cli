# PRD (focused copy): team orchestration agent and model access in the `afs` CLI

**Source:** the root repository's `team-orchestration-agent-PRD.md` (status: approved by J, 2026-11-23, release 1.1.0 in every
repository). This file keeps only what the CLI must do; the root document wins on any difference. Contract:
`docs/contracts/api-contract.md` v1.4.0, section 4.4d; OpenAPI `api/openapi.yaml` of `-restapi` at f6d240c (194 operations).

## Requirements that reach the CLI

| ID | Requirement | Pri |
|---|---|---|
| ORC-29 | New service `orchestration`: `describe-orchestration`, `put-orchestration --web-ui-exposure off\|private\|public --web-api-exposure ... [--registration-mode]`, `change-collaboration-tool --tool paperclip\|buzz`; `admin list-orchestration-endpoints`. `teams create-team --orchestration Tool=...,WebUiExposure=...,WebApiExposure=...,ModelId=...` and `--max-collaboration-exposure`. Help, examples, columns, goldens and `AccessDenied` tests for each | P0 |
| ORC-30 | `collaboration describe-collaboration`, `put-collaboration` and `admin list-collaboration-endpoints` keep working, print one deprecation line to stderr naming the replacement, and are removed in 2.0.0. `agents create-agent --role` is removed from help; passing it returns the API's `ValidationError` (exit 254). `collaboration invoke-api` picks the hostname of the Web API's effective level and fails with exit 254 and "Web API is off (N/A for Buzz)" when the tool has none | P0 |
| MOD-8 (CLI part) | `agents create-agent`/`update-agent --models all \| <id> ... --default-model-id <id>`, the same on pools; `teams put-team-models --access all \| selected [--model-ids ...]`; `admin put-service-model-access --access all \| selected`; table columns `DefaultModel`, `Models` | P0 |

## CLI stories (root PRD section 7.6; numbered `OR-CLI-26` to `OR-CLI-28` here, because the 1.0.1 Paperclip stories already use `CLI-US-26` to `CLI-US-28`)

**OR-CLI-26 · Orchestration commands.** *(refines US-34, US-37, US-38)* The commands of ORC-29 with help, examples,
table columns (`Tool`, `WebUi`, `WebApi`, `Limit`, `Effective`), goldens and `AccessDenied` tests (exit 254).

**OR-CLI-27 · Deprecations.** *(refines US-40)* The deprecated commands of ORC-30 print one stderr line and still
succeed; `invoke-api` follows the effective level and exits 254 for N/A surfaces with the AWS-style message.

**OR-CLI-28 · Model access options.** *(refines US-41)* The flags of MOD-8 with help, examples, goldens and
`AccessDenied` tests; `--model-id` keeps working with a deprecation line on stderr.
