# Implementation notes: orchestration-and-model-access

## Decisions

| ID | Decision | Why |
|---|---|---|
| D-O1 | New milestones M11 (orchestration, endpoints) and M12 (model access); `CurrentMilestone` 12 | the existing milestone gate (`Pending`) then proves every new operation has an example and columns |
| D-O2 | The 1.1.0 stories are called `OR-CLI-26` to `OR-CLI-28` in tests and docs | the 1.0.1 Paperclip stories already own `CLI-US-26` to `CLI-US-28` and the archived traceability test reads them |
| D-O3 | `put-orchestration` takes `--web-ui-exposure` and `--web-api-exposure` as local options folded into `webUi` and `webApi`; the derived `--web-ui` and `--web-api` JSON options are hidden | the AWS-style flag is a level, not a JSON object; the API requires both surfaces on every write, so the CLI requires both (or input JSON) rather than reading one first (one operation, one call) |
| D-O4 | `--models` is a list: the word `all`, or model ids; `--default-model-id` is separate; `--model-id` stays with a warning | PRD MOD-8; the combination `all` plus ids is a syntax error, not a server rule |
| D-O5 | `agents list-agents --role` is an open enum with `worker`, `orchestration` | the snapshot's enum is `worker`, `collaboration` while contract 4.4d says `?role=worker\|orchestration`; the API decides |
| D-O6 | `admin put-setting --key` is an open enum: the snapshot's keys without `collab.base_domain` and `collab.public_endpoints`, plus the 14 new keys | the service refuses the renamed pair with `setting_renamed` and accepts the new keys, which the snapshot does not list |
| D-O7 | `--role` on `create-agent` is hidden with an open enum | PRD ORC-30: out of help, the API's `ValidationError` (`role_reserved`) is what a caller who passes it sees (exit 254) |
| D-O8 | A deprecated operation prints its line before the request; an option's line is printed when it was given by flag or in `--cli-input-json`; the `Deprecation` header prints a line only when the CLI printed none | the stderr line must not depend on the call succeeding; no duplicates |
| D-O9 | Agent tables drop the `Template` column and the Paperclip column becomes `Link` | with `DefaultModel` and `Models` added, ten columns no longer fit 80 columns, and the orchestration agent's template is hidden from non-admins (CD-32) |
| D-O10 | `invoke-api` refuses a Web API that is off with the code `SurfaceOff` (exit 254) before sending any key | PRD ORC-30; the message is the one the PRD names |
| D-O11 | No confirmation prompt for `change-collaboration-tool`; its help says what happens | one operation, one API call; the destructive step is recoverable (the API keeps a 90-day snapshot) |
| D-O12 | `ModelId` in the `--orchestration` shorthand is sent as `models.defaultModelId` | the snapshot's `CreateTeamRequest.orchestration` has no `modelId` (contract 4.2 calls it shorthand) |
| D-O13 | `put-team-models` stays versionless (no If-Match read) | the contract says the ETag of the models list is the team's version, but the operation's `If-Match` is optional on PUT and 1.0.1 behaviour stands (D-C8) |

## Gotchas

- `wait.Table` polls `GetTeamOrchestration` for `put-orchestration` and `change-collaboration-tool`; the resource has `generation` and `status.observedGeneration`, so a stale `Ready` is ignored.
- The private hostnames of an orchestration agent are under `orchestration.private_base_domain`; the profile key `collab_base_domain` (an optional guard on `invoke-api`) must name the domain of the level in use, or be unset.
- `--query 'items[?callable].id'` in two examples relies on the JMESPath filter, which the renderer supports.

## Review (step 5, against the AGENTS.md checklist)

See the section "Code and design review" at the end of this file after the review step.

## Code and design review

Checked against the AGENTS.md list: naming and options follow the AWS CLI rules, help and an example exist for every new operation; `--output json` is untouched (columns and notices only affect table, text and stderr); no secret, token or key reaches any output (no new secret input); errors use the AWS-style message, `SurfaceOff` exits 254; context is passed through and no endpoint or credential is hardcoded; every new operation has an AccessDenied test via `clitest.CheckOperations`. Findings: none critical or important. Minor, fixed: tests named CLI-US-26..28 collided with 1.0.1 numbering and were renamed OR-CLI-26..28; the agent table exceeded 80 columns (Template dropped). Open items for `-restapi` are in research.md (setting key enum, ListAgents role enum, `modelId` and `volume` on CreateTeamRequest.orchestration).
