# Changelog

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
