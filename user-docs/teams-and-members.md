# Teams, members, audit, settings and identity providers

Every command here maps to one API call. Run `afs <service> <operation> help` for options and
examples; the [command reference](reference/README.md) lists them all.

## Teams

```bash
afs teams list-teams --state active --query 'items[].[id,displayName]' --output text
afs teams describe-team --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001
afs teams create-team --slug payments --display-name "Payments" --kind persistent \
  --max-budget-usd 500.00 --first-admin '{"email":"ada@example.com"}' \
  --tags '{"project":"checkout"}' --wait
```

- `--generate-cli-skeleton` prints every field of `create-team`; `--cli-input-json file://team.json`
  sends that file's members unchanged.
- `--wait` polls the team until it is `Ready` at the generation your change produced, then prints the
  final team. `Failed` or `Degraded` exits 254 with the conditions on stderr; a timeout exits 255
  (`--wait-timeout 10m`). The change stands either way.
- `update-team`, `delete-team` and the lifecycle actions send `If-Match`. Without `--if-match`, `afs`
  reads the current version, sends it, and after the write succeeds prints
  `Applied to version N; use --if-match to guard against concurrent changes.` on stderr. With `--if-match 7` it sends exactly 7, and a stale version fails with 412 naming the current one.
- `suspend-team` makes every later write in the team fail with `Conflict`; `resume-team` lifts it.
  `delete-team` needs a platform role: a team-admin gets `AccessDenied` (exit 254).
- `put-team-models --model-ids m1 m2` replaces the model list. The resource has no version, so no
  `If-Match` is sent.
- Guardrail fields (`--max-budget-usd`, `--isolation-tier`, ...) need the matching permission; the API
  answers `AccessDenied` naming it. `--ceiling-override-reason` is an ordinary text option.

`afs teams list-revisions` and `describe-revision --revision-id N` show the configuration history.

## Teams and Paperclip

`afs teams create-team` and `update-team` take `--max-collaboration-exposure off|private|public` (replaces the deprecated `--may-expose-collaboration`; see [orchestration](orchestration.md)), `--orchestration`, `--may-expose-collaboration` or `--no-may-expose-collaboration` (service administrators: it lets the
team publish its Paperclip endpoint).

## Members

```bash
afs members add-member --team-id team_... --subject-type user --email ada@example.com --role team_admin
afs members add-member --team-id team_... --subject-type group --group-name payments-devs --role team_viewer
afs members update-member --team-id team_... --binding-id mem_... --role team_viewer
afs members remove-member --team-id team_... --binding-id mem_...
```

The roles are `team_admin`, `team_collaborator` (may use the team's Paperclip, see [Team collaboration](collaboration.md)) and
`team_viewer`.

There is no single-member read, so `update-member` and `remove-member` find the version in the member
list. Pass `--if-match` to skip that read.

## Team audit events

```bash
afs audit list-audit-events --team-id team_... --from 2026-11-01T00:00:00Z --outcome denied --output csv > audit.csv
```

Filters: `--from --to --actor-id --action --outcome --resource-kind --resource-id`. With `--output csv`
every page is fetched and written under one header row.

## Service settings (platform administrators)

```bash
afs admin get-settings --output json
afs admin put-setting --key api.cors_origins --value '["https://console.example.com"]'
afs admin put-setting --key api.idempotency_key_ttl_seconds --value 86400
```

`--value` is sent as JSON when it parses as JSON (a list, number or boolean), otherwise as a string.
An out-of-range value is a `ValidationError` that names the range.

## Identity providers (platform administrators)

```bash
afs admin list-identity-providers
afs admin describe-identity-provider --provider-id idp_...
afs admin create-identity-provider --name "Corp SSO" --kind okta --issuer-url https://idp.example.com \
  --console-client-id console-app --cli-client-id afs-cli --client-secret-arn arn:aws:secretsmanager:...
afs admin delete-identity-provider --provider-id idp_...
```

Output carries the secret's ARN only, never a value. Deleting a provider that teams still use is a
`Conflict` (exit 254).
