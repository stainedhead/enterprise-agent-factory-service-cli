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
