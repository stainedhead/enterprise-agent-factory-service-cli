# Changelog

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
