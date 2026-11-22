# Team collaboration (Paperclip)

Every team may run one **collaboration agent**, Paperclip, and connect its other agents to it. The
`afs collaboration` commands read and change the team's collaboration settings, list and change the link of
each worker agent, and (for a service administrator) list every team's endpoint. They call the AFS API; the
API decides what you may do.

Contract: API version 1.3. A CLI built for 1.3 against an older service answers `404` on the new paths; a
service newer than the CLI prints a one-time notice and keeps working.

## Create the collaboration agent

The collaboration agent is an ordinary agent with `--role collaboration`, made from a template whose
capabilities include `paperclip`:

```bash
afs agents list-agent-templates --capability paperclip
afs agents create-agent --team-id team_... --name paperclip --role collaboration --template paperclip \
  --image registry.example.com/agents/paperclip@sha256:... --model-id anthropic.claude-sonnet-4-5 --wait
afs agents list-agents --team-id team_... --role collaboration
```

A team has at most one. A second `create-agent --role collaboration` is the API's `409 Conflict`. The role
cannot be changed afterwards, so `update-agent` has no `--role`.

## See and change the settings

```bash
afs collaboration describe-collaboration --team-id team_...
afs collaboration put-collaboration --team-id team_... --web-ui-enabled --no-rest-api-enabled
afs collaboration put-collaboration --team-id team_... --web-ui-enabled --rest-api-enabled \
  --registration-mode managed --wait
```

- `put-collaboration` replaces both switches, so it needs `--web-ui-enabled` or `--no-web-ui-enabled`, and
  `--rest-api-enabled` or `--no-rest-api-enabled` (or both in `--cli-input-json`). `--registration-mode` is
  `manual` or `managed`; leave it out to keep the current value.
- Like other updates it sends `If-Match`: pass `--if-match 4` to guard against a concurrent change, or let
  `afs` read the version and say so on stderr.
- `--wait` polls until the endpoint is `Ready` at the generation of your change (a stale `Ready` does not
  count). `Failed` or `Degraded` exits 254 and prints the conditions; a timeout exits 255.
- `describe-collaboration` with `--output table` or `text` prints one stderr line per surface that is not
  effective, for example `Not effective: web UI (base_domain_unset)`. A surface is effective only when your
  choice, the team's guardrail, the service setting and (for the web UI) a collaboration sign-in client all hold.
- When the API refuses a change, each failing field is printed as `field: <field> reason: <reason>` with a
  one-line hint, for example `Ask a service admin to set collab.base_domain (afs admin put-setting).`

## Connect agents

```bash
afs collaboration list-collaboration-links --team-id team_... [--link-state failed] [--connect | --no-connect]
afs collaboration connect-agent --team-id team_... --agent-id agt_... [--paperclip-agent-id ag_...]
afs collaboration disconnect-agent --team-id team_... --agent-id agt_...
```

`--paperclip-agent-id` is accepted only while the registration mode is `manual` (in `managed` mode the
reconciler records the id and the API answers with the hint `Registration is managed; the reconciler records the id.`).
`afs agents update-agent --collaboration-connect` has the same effect as `connect-agent`.
`afs agents create-agent` takes `--collaboration-connect` or `--no-collaboration-connect` for workers; with
neither, the API default applies (connected when the team has a collaboration agent). The collaboration agent itself and
pool workers cannot be connected.

## Secrets in Paperclip

`afs secrets list-secrets --managed-by user|afs|paperclip` filters by who manages a secret, and the table has a
`ManagedBy` column. Secrets managed by AFS (the Paperclip boot secrets) or by Paperclip cannot be changed,
deleted or given a new value with `update-secret`, `delete-secret` or `put-secret-value`: the API answers
`409 Conflict` and `afs` adds the hint. `rotate-secret` on an AFS-managed secret sends no value (AFS generates it).
The master key `paperclip-secrets-master-key` cannot be rotated. The names `paperclip` and `paperclip-...` are
reserved, except `paperclip-agent-api-key` on a worker agent (manual registration).

## Call the Paperclip API from a pipeline

`afs collaboration invoke-api` sends one request to your team's Paperclip REST API. Your AFS credential
(an API key in `AFS_API_KEY`, or your `afs sso login` session) proves who you are to AFS; a Paperclip key
proves what you may do inside Paperclip. The REST API must be switched on for the team
(`afs collaboration put-collaboration --rest-api-enabled`).

```bash
# read
afs collaboration invoke-api --team-id team_... --path /api/companies --paperclip-key-file ./pc.key

# create an issue
afs collaboration invoke-api --team-id team_... --method POST --path /api/issues \
  --body file://issue.json --paperclip-key-file ./pc.key
```

- **The Paperclip key never goes on the command line.** Put it in a file with mode `0600` (`chmod 600 pc.key`)
  or pass `--paperclip-key-file -` and pipe it in. `--paperclip-key <value>` does not exist.
- **The response is printed unchanged**, so pipe it to `jq`. `--include` adds the status line and headers.
  `--query`, `--output`, `--wait` and the pagination options do not apply and exit 252.
- **Exit codes.** 0 for a 2xx or 3xx answer; 254 for any error answer (AFS's own errors print
  `An error occurred (<Code>) when calling the InvokeApi operation`, Paperclip's own errors print
  `HTTP <status>` and leave Paperclip's body on stdout); 252 for a bad option; 255 for a local problem
  (key file mode, origin, network, no credential).
- **Where it goes.** The host comes from the team's `restApiUrl`, one read of `describe-collaboration`. Use
  `--collab-url https://<team>.<base domain>` to skip the read. Only `https` on port 443 and a DNS name are
  accepted; set `collab_base_domain` in your profile to pin the host to your service's domain. Redirects are
  not followed and cookies are never kept.
- **Workload tokens are refused.** A pipeline inside the cluster must use an API key.
- **Retries.** `GET` and `HEAD` retry `429`/`503` up to three tries; other methods retry only a failed connection.
- **`--debug`** logs the method, host, path (no query string), status and timing. Credentials, bodies and
  cookies are never logged.

## Roles

Give a person or a pipeline the **collaborator** role with `--role team_collaborator` on
`afs members add-member`, `members update-member` or `afs api-keys create-api-key`. A collaborator may use the team's
Paperclip but cannot change the team's AFS configuration. Service-level keys and bindings cannot hold it.

## Service administrators

```bash
afs teams update-team --team-id team_... --may-expose-collaboration      # the per-team guardrail
afs admin put-agent-template --name paperclip --harness paperclip --current-digest sha256:... \
  --allowed-digests sha256:... --capabilities paperclip
```

A team slug whose collaboration hostname is held (30 days after the deletion of a team that had one) is refused with
`409 Conflict` and a hint to choose another slug.

## Troubleshooting

| You see | Do |
|---|---|
| `This team has no collaboration agent.` | `afs agents create-agent --role collaboration ...` |
| `The team is suspended.` | Resume the team: `afs teams resume-team --team-id team_...` |
| `Not effective: ... (team_not_allowed)` | A service admin turns the guardrail on: `afs teams update-team --may-expose-collaboration` |
| `Not effective: ... (base_domain_unset)` | A service admin sets `collab.base_domain` with `afs admin put-setting` |
| `This secret is managed by AFS or Paperclip and cannot be changed here.` | Change it in Paperclip, or rotate with `rotate-secret` if AFS generates it |

## Service administrator setup

A service administrator prepares the service once, then switches teams on:

```bash
afs admin put-setting --key collab.base_domain --value '"collab.example.com"'
afs admin put-setting --key collab.public_endpoints --value true
afs admin update-identity-provider --provider-id idp_... \
  --collab-client-id paperclip-web --collab-client-secret-arn arn:aws:secretsmanager:eu-west-1:111122223333:secret:afs/collab-AbCdEf
afs teams update-team --team-id team_... --may-expose-collaboration
afs admin list-collaboration-endpoints [--team-id team_...] [--state Failed] [--output csv]
```

- `--collab-client-secret-arn` takes the **ARN** of a Secrets Manager secret, never the secret's value; a value that is not an ARN
  exits 252 before any call, and `--debug` shows the ARN (a reference) but never a client secret.
- The seven `collab.*` settings (`base_domain`, `public_endpoints`, `session_ttl_minutes`, `max_body_mb`, `content_scan`,
  `allowed_source_cidrs` and `denied_spike_threshold`) go through `put-setting`; the value is JSON, and the service names the valid range when it refuses one.
- `list-collaboration-endpoints` shows each team's hostname, effective surfaces, phase and last access; with `--output csv` it
  pages through the whole list.
