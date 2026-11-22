# Agents, secrets, configuration values and external systems

## Agents

```bash
afs agents list-agents --team-id team_...
afs agents create-agent --team-id team_... --name coder --template claude-code \
  --image registry.example.com/agents/coder@sha256:... --model-id anthropic.claude-sonnet-4-5 \
  --volume '{"sizeGiB":20,"storageClass":"gp3-standard","onDelete":"snapshot_then_delete"}' --wait
afs agents describe-agent --team-id team_... --agent-id agt_...
afs agents pause-agent|resume-agent|restart-agent --team-id team_... --agent-id agt_...
```

- Agents are addressed by `--agent-id`; two teams may use the same `--name`.
- Images must be pinned by digest; a tag is a `ValidationError` from the API (exit 254).
- `--wait` returns 0 when the agent is `Ready` at the generation your change produced and prints the
  final agent; `Failed` is 254, a timeout is 255.
- `update-agent` and `delete-agent` send `If-Match`. Pass `--if-match 7` to choose the version; without
  it `afs` reads the current one and tells you on stderr. When the agent moved on, the error says
  `The resource changed (version 8). Re-run with --if-match 8 to apply your change to the latest version.`
- Snapshots: `afs agents create-snapshot` and `list-snapshots`. Catalogs: `list-storage-classes`,
  `list-agent-templates`, and `afs models list-models`.
- Agent configuration values: `list-config-values`, `put-config-value --name N --source literal --value V`
  (or `--source namespace` to take the team's value), `delete-config-value`.
- Paperclip: `create-agent --role collaboration` creates the team's collaboration agent; `--collaboration-connect`
  or `--no-collaboration-connect` and `--paperclip-agent-id` control whether a worker is connected; `list-agents --role`
  and `list-agent-templates --capability paperclip` filter. The agent table has `Role` and `Paperclip` columns. See
  [Team collaboration](collaboration.md).
- Per-agent external systems: `list-external-systems`, `set-external-system-state --system-id ... --enabled` or `--no-enabled`.

## Secrets

`afs` never takes a secret value as an ordinary argument, so it never lands in shell history, process
listings or logs. Give it one of:

| How | Example |
|---|---|
| A file (read byte for byte, no newline trimming) | `--secret-value file://token.txt` (`fileb://` is the same) |
| Standard input | `printf '%s' "$TOKEN" \| afs secrets create-secret ... --secret-value -` |
| A hidden prompt (terminal only) | leave `--secret-value` out |

`afs secrets list-secrets --managed-by user|afs|paperclip` filters by who manages a secret (a `ManagedBy` column shows it); AFS and
Paperclip secrets are read-only here and answer `409 Conflict` (see [Team collaboration](collaboration.md)).

A plain argument (`--secret-value hunter2`) exits 252 before any call. With no terminal and no value the
command exits 252 and lists the options above. The value must be UTF-8 text of 1 to 65,536 bytes.

```bash
afs secrets create-secret --team-id team_... --name github-pat --scope team --secret-value file://token.txt
afs secrets put-secret-value --team-id team_... --secret-id sec_... --secret-value -     # prints nothing (204)
afs secrets rotate-secret --team-id team_... --secret-id sec_... --secret-value file://new.txt
afs secrets list-secrets --team-id team_...           # names, scope and metadata, never values
afs secrets delete-secret --team-id team_... --secret-id sec_...
```

- `--cli-input-json` may carry the value only from a file (`--cli-input-json file://secret.json`); inline
  JSON with a `value` is refused (252).
- The value is never printed. With `--debug` the request log shows `"value": "[redacted]"` and
  `Authorization: [redacted]`.
- `delete-secret` fails with `Conflict` when agents use the secret, and lists them: `Dependents: agt_1, agt_2`.

## Team configuration values and external systems

```bash
afs teams list-config-values --team-id team_...
afs teams put-config-value --team-id team_... --name REGION --value eu-west-1   # creates or replaces
afs teams list-external-systems --team-id team_...
afs teams declare-external-system --team-id team_... --slug ledger --name "Ledger API" --category other \
  --owner ledger-team@example.com --purpose "Reconcile payouts" --endpoints ledger.example.com:443
afs teams set-external-system-state --team-id team_... --system-id ext_... --enabled
```

`put-config-value` reads the list to find the current version; a value that does not exist yet is
created without `If-Match`. Declaring a system when the team may not add external systems is the API's
`AccessDenied`.

## Platform catalogs (administrators)

```bash
afs admin list-external-systems            # create-/update-/delete-external-system, describe-external-system
afs admin list-storage-classes             # put-storage-class --name N --volume-type gp3|io2, delete-storage-class
afs admin list-agent-templates             # put-agent-template --name N --harness H --current-digest D --allowed-digests D1 D2
```

Deleting a built-in system, a storage class in use or a template agents use is a `Conflict` (254).
