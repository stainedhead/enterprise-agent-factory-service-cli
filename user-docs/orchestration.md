# Team orchestration agent and model access

Contract: API version 1.4. Release: `afs` 1.1.0.

Every team has one **orchestration agent**. It runs one collaboration tool, Paperclip or Buzz, and replaces the
1.0 "collaboration agent". `afs orchestration` reads and changes how it is exposed; model access decides which
models agents may call. The API decides what you may do; the CLI only shapes the request and explains refusals.

```
                  service  --all|selected-->  team  --all|selected-->  agent / pool
 models catalog     (admin put-service-       (teams put-team-          (--models, --default-model-id)
                     model-access)             models)
```

## Orchestration

```bash
afs orchestration describe-orchestration --team-id team_...
afs orchestration put-orchestration --team-id team_... --web-ui-exposure private --web-api-exposure off
afs orchestration put-orchestration --team-id team_... --web-ui-exposure public --web-api-exposure public \
  --registration-mode managed --wait
afs orchestration change-tool --team-id team_... --tool buzz --wait
afs admin list-orchestration-endpoints --output table
```

- Exposure is `off` (inside the team namespace only), `private` (internal ALB, networks routed to the
  environment VPC) or `public` (internet-facing ALB). Both surfaces are sent together: give
  `--web-ui-exposure` and `--web-api-exposure` (or both in `--cli-input-json`).
- A level above the team's limit is refused. A service administrator sets the limit with
  `afs teams update-team --max-collaboration-exposure off|private|public`.
- The tool is `paperclip` or `buzz`. Buzz has no web surfaces: the table shows `N/A`, and only `off` is
  accepted. Changing the tool is its own operation (`change-tool`); `--wait` polls until the change is applied.
- `describe-orchestration` prints a stderr note when a surface is not effective (for example private, not
  public) with the reason from the API.
- `afs admin list-orchestration-endpoints` lists every team's endpoint; while the 1.4 backfill runs it prints
  the progress on stderr.
- `afs teams create-team --orchestration Tool=paperclip,WebUiExposure=off,WebApiExposure=off,ModelId=<id>,VolumeSizeGiB=50,VolumeOnDelete=snapshot_then_delete` sets it at creation (the flat contract fields `webUiExposure`, `webApiExposure`, `modelId` and `volume`).
- `afs collaboration invoke-api` calls the Web API through `GetTeamOrchestration`; if the Web API is off it
  fails with `SurfaceOff` (exit 254).

## Model access

Access is `all` (every model enabled one level up, now and later) or `selected` (a list), at four levels.

```bash
afs admin put-service-model-access --access all
afs teams put-team-models --team-id team_... --access selected --model-ids anthropic.claude-sonnet-4-5 anthropic.claude-haiku-4-5
afs agents create-agent --team-id team_... --name triage --models anthropic.claude-haiku-4-5 \
  --default-model-id anthropic.claude-haiku-4-5
afs agents update-agent --team-id team_... --agent-id agt_... --models all
afs teams list-team-models --team-id team_...
```

- `--models all` or `--models <id> <id>...` on `create-agent`, `update-agent` and the pool operations; `--default-model-id` picks the default.
- Model lists show `Callable` and `Price`: a model can be enabled but not callable (no price, or not enabled one level up).
- Narrowing never edits an agent's stored selection; restoring a model restores access.

## Deprecated, kept until 2.0.0

| Old | Use instead |
|---|---|
| `afs collaboration describe-collaboration`, `put-collaboration` | `afs orchestration describe-orchestration`, `put-orchestration` |
| `afs admin list-collaboration-endpoints` | `afs admin list-orchestration-endpoints` |
| `--may-expose-collaboration` | `--max-collaboration-exposure` |
| `--model-id` on agents and pools | `--default-model-id` (and `--models`) |
| `--role collaboration` | not needed; the orchestration agent is created with the team |
| settings `collab.base_domain`, `collab.public_endpoints` | `orchestration.public_base_domain`, `orchestration.public_endpoints` |

They keep working. Each prints `Warning: ... is deprecated and will be removed in 2.0.0; use ... instead.` on
stderr, and the service's `Deprecation`/`Sunset` headers are reported once when no static warning applied.

## Refusals

Errors name the reason with a one-line fix, for example `exposure_above_limit` (ask an administrator to raise
the team limit), `surface_not_supported_by_tool` (Buzz has no web surfaces), `tool_change_in_progress` (wait and
retry), `model_not_enabled_for_team`, `default_model_not_allowed`, `models_selection_empty` and
`model_ids_not_allowed_for_all`. See [Errors and exit codes](errors-and-exit-codes.md).
