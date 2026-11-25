# Changelog

## v1.3.0 - 2026-11-25

Governed MCP servers (API contract v1.6.0). Spec: `specs/261125-mcp-gateway/` (or `specs/archive/261125-mcp-gateway/`
once archived). Guide: `user-docs/mcp-servers.md`.

### Added
- `afs mcp-servers` service, 24 commands: the catalog (`list-`, `create-`, `describe-`, `update-`, `delete-mcp-server`,
  `import-mcp-server-manifest`, `put-mcp-server-credential`, `put-mcp-policy`, `refresh-mcp-inventory`), review
  (`list-mcp-reviews`, `list-mcp-primitives`, `get-mcp-primitive-change`, `approve-`/`reject-mcp-primitive`,
  `approve-`/`reject-mcp-server`, `request-mcp-server-changes`, `publish-`, `suspend-`, `resume-`, `retire-mcp-server`)
  and query capture (`get-`, `start-`, `stop-mcp-query-capture`).
- 18 `afs teams` commands for a team's servers: switches, the team's own servers and their review, credentials,
  policies, OAuth connect and disconnect, `test-mcp-federation`.
- `afs agents` and `afs pools`: `get-effective-*-mcp` and federated identities (`list-`, `put-`, `delete-`); `--mcp` on
  create and update.
- `afs audit list-mcp-events` and `list-team-mcp-events`; `afs reports get-mcp-usage-report` and
  `get-team-mcp-usage-report`.
- `--client-secret` is read like `--secret-value`: from a file, standard input or a no-echo prompt only.
- Hints for the 23 new error reasons; milestones M15 and M16.

### Changed
- Contract snapshot re-pinned to `-restapi` ef92a4c (v1.6.0, 276 operations, 273 commands); `afs --version` prints
  `contract/1.6`.

## v1.2.0 - 2026-11-25

Model endpoint gateway (API contract v1.5.0, section 4.13). Spec: `specs/261125-model-endpoint-gateway/` (or
`specs/archive/261125-model-endpoint-gateway/` once archived). Guide: `user-docs/model-endpoints.md`.

### Added
- `afs model-endpoints` service, 23 commands: the endpoint catalog (`list-`, `create-`, `describe-`, `update-`,
  `delete-model-endpoint`, `put-endpoint-models`), prices and proposals (`create-endpoint-model-price`,
  `list-`, `approve-`, `reject-endpoint-price-proposal`), `list-endpoint-teams`, credentials
  (`put-model-endpoint-credential`, `rotate-model-endpoint-credential`, `put-team-model-endpoint-credential`),
  `put-endpoint-federation`, `import-provider-statement`, provider switches (`list-inference-providers`,
  `update-inference-provider` with `--dry-run`), and team endpoints (`list-`, `create-`, `update-`,
  `delete-team-model-endpoint`, `put-team-endpoint-models`).
- `afs teams put-team-gateway-settings` and `afs teams put-team-model-endpoints`.
- `afs reports get-inference-report` and `get-team-inference-report` (group by route, model, agent, pool, team, hour
  or day; CSV).
- `--routes all|<slug>...` and `--default-route` on `create-agent`, `update-agent`, `create-pool` and `update-pool`;
  `--default-model-id` takes `<route>/<model>`.
- `--route` and `--group-by route` on the usage and cost reports.
- Hints for `slug_taken`, `slug_reserved`, `slug_invalid`, `base_url_not_https`, `provider_disabled`,
  `no_fallback_model`, `endpoint_base_url_placeholder`, `needs_service_admin`, `team_endpoints_not_allowed` and
  `endpoint_enabled`.
- Milestones M13 and M14.

### Changed
- The pinned snapshot is contract v1.5.0 (`-restapi` 7d48ea7, 221 operations, 219 commands; `afs --version` prints
  `contract/1.5`).
- Credential values are read only from `file://`, `fileb://`, `-` (standard input) or a no-echo prompt, and are
  redacted under `--debug`. Credential writes send `If-Match` only when `--if-match` is given (the API returns no
  credential version).
- `import-provider-statement --statement file://PATH` sends a CSV file as `text/csv` with `period` and `source` in the
  query, and a JSON file as the body. The executor can now send a non-JSON request body.
- `update-inference-provider` reads its version from the provider list; `put-team-model-endpoints` reads none.


## v1.1.0 - 2026-11-23

Snapshot refreshed to `-restapi` f8b4404 (still 194 operations): the real enums replace the open-enum workarounds
(`put-setting --key` lists the 14 `orchestration.*` and `models.*` keys; `list-agents --role` is worker, orchestration
or collaboration), `update-agent --role` is documented, `teams create-team --orchestration` sends the flat
`webUiExposure`, `webApiExposure`, `modelId` and `volume` (`VolumeSizeGiB`, `VolumeOnDelete`) fields. `put-setting --key`
stays open so the API answers a renamed `collab.*` key with `setting_renamed` (hint included). The generated
types use an overlay to name two colliding inline enums.

Team orchestration agent and model access (API contract v1.4.0). Spec: `specs/archive/261123-orchestration-and-model-access/`
(or `specs/archive/` once archived). Guide: `user-docs/orchestration.md`.

### Added
- `afs orchestration` service: `describe-orchestration`, `put-orchestration` (`--web-ui-exposure`, `--web-api-exposure`
  off|private|public, `--registration-mode`, If-Match, `--wait`), `change-tool` (`--tool paperclip|buzz`, `--wait`);
  `N/A` for Buzz web surfaces and stderr notes for surfaces that are not effective.
- `afs admin list-orchestration-endpoints` (with backfill progress on stderr) and `afs admin put-service-model-access`.
- Model access: `--models all|ids` and `--default-model-id` on agents and pools; `teams put-team-models --access`;
  `Callable` and `Price` columns on model lists; `teams create-team/update-team --orchestration` and `--max-collaboration-exposure`.
- Hints for the new refusal codes (`exposure_above_limit`, `surface_not_supported_by_tool`, `tool_change_in_progress`,
  `model_not_enabled_for_team`, `default_model_not_allowed` and others); `SurfaceOff` for `invoke-api`.
- Milestones M11 and M12; settings keys `orchestration.*` and `models.*` for `put-setting`.

### Changed
- The pinned snapshot is contract v1.4.0 (`-restapi` f6d240c, 194 operations, 192 commands; `afs --version` prints `contract/1.4`).
- `collaboration invoke-api` reads `GetTeamOrchestration`. Agent table: `Template` column dropped, `Paperclip` renamed `Link`, `DefaultModel`, `Models` and `Role` added.
- `--role` on `create-agent` is hidden and an open enum; the swarm template skeleton shows `spec.orchestration`.

### Deprecated (kept until 2.0.0, a warning on stderr)
- `collaboration describe-collaboration`, `put-collaboration`; `admin list-collaboration-endpoints`;
  `--may-expose-collaboration`; `--model-id` on agents and pools; settings `collab.base_domain`, `collab.public_endpoints`.

## v1.0.1 - 2026-11-22

Paperclip team collaboration (API contract v1.3.1). Spec: `specs/archive/261122-paperclip-collaboration/`
(or `specs/261122-paperclip-collaboration/` until it is archived).

### Added
- `afs collaboration` service: `describe-collaboration`, `put-collaboration` (If-Match, `--wait`),
  `list-collaboration-links`, `connect-agent`, `disconnect-agent`, `rotate-collaboration-keys`.
- `afs collaboration invoke-api`: one request to a team's Paperclip REST API with the AFS credential in
  `X-Afs-Authorization` and a Paperclip key read from a 0600 file or stdin (never an argument); origin
  validation, no redirects or cookies, narrow retries, edge errors mapped to AWS-style messages.
- `afs admin list-collaboration-endpoints` (table and CSV); identity provider `--collab-client-id` and
  `--collab-client-secret-arn`; the seven `collab.*` settings through `admin put-setting`.
- `agents create-agent --role collaboration`, `--collaboration-connect`, `--paperclip-agent-id`;
  `Role`, `Paperclip` and `Capabilities` columns; `--capability` filter on `list-agent-templates`;
  `admin put-agent-template --capabilities`.
- `secrets list-secrets --managed-by`; a `ManagedBy` column.
- `teams --may-expose-collaboration`; role `team_collaborator` on `members` and `api-keys`.
- Report kind `collaboration` with `--group-by surface|reason|person`; swarm template `spec.collaboration`.
- Profile key `collab_base_domain`; reason hints for the collaboration errors; `Patterns:` for `SecretDetected`.

### Changed
- The pinned API snapshot is re-pinned to contract v1.3.1 (`-restapi` 763e939) (`afs --version` prints `contract/1.3`);
  generated types, descriptors and the reference docs are regenerated.
- `X-Afs-Authorization` joins the headers masked in `--debug` output; unknown options echo only their name and
  suggest the options that start with it.
- The report schedule `--report` is an open enum (`put-setting --key` is closed again at 763e939, which lists the `collab.*` keys): the API judges values the snapshot does not list.

### Release
- Tag `v1.0.1`. The version is stamped by the Makefile from the tag (`make build VERSION=1.0.1`).
