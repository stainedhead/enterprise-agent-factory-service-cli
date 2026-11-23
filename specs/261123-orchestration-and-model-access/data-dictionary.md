# Data dictionary: orchestration-and-model-access

The CLI holds no data model of its own: every shape below is the API's (OpenAPI schema in parentheses). Output field
names are exactly the API's (`--output json` is the stable machine contract).

| Name | Shape | Where the CLI shows or sends it |
|---|---|---|
| Exposure level | `off` \| `private` \| `public` | `--web-ui-exposure`, `--web-api-exposure`, `--max-collaboration-exposure`, `--exposure` (list filter), shorthand keys |
| Tool | `paperclip` \| `buzz` | `--tool`, `--orchestration Tool=`, `--tool` filter; shown by name (Paperclip, Buzz) in notices |
| `TeamOrchestration` | `{teamId, agentId, agentName, phase, tool, toolVersion, registrationMode, limit, webUi, webApi, connectedAgents, unmanagedToolAgents, version, generation, status}` | `describe-orchestration`, `put-orchestration`, `change-collaboration-tool` |
| `OrchestrationSurface` | `{supported, exposure, effective, reasons[], hostname, url}` | columns `WebUi`, `WebApi`, `EffectiveWebUi`, `EffectiveWebApi`; `supported: false` prints N/A |
| Reasons | `tool_not_supported`, `team_limit`, `service_switch_public`, `service_switch_private`, `base_domain_unset`, `provider_without_collab_client` | printed as returned on stderr for table and text |
| `OrchestrationEndpoint` | one team and surface with a stored exposure other than off | `admin list-orchestration-endpoints`; `migration` progress on stderr |
| `AgentModels` / `AgentModelsInput` | `{access: all\|selected, modelIds[], defaultModelId, effectiveModelIds}` | `--models all \| <id> ...`, `--default-model-id`; columns `DefaultModel`, `Models` |
| `TeamModels` | `{access, modelIds[], items[Model]}` | `list-team-models`, `put-team-models --access --model-ids` |
| `Model` | adds `callable`, `priceState`, `inferenceProfileState`, `origin`, `teamsWithAccess` | columns `Callable`, `Price`, `Profile`, `Origin`, `Teams` |
| `ServiceModelAccess` | `{access, enabledModels, callableModels, version}` | `put-service-model-access` |

Local words: `all` (the one word of `--models` that is not a model id); the shorthand keys of `--orchestration`
(`Tool`, `WebUiExposure`, `WebApiExposure`, `RegistrationMode`, `ModelId`, alias `DefaultModelId`), matched without regard to case.
