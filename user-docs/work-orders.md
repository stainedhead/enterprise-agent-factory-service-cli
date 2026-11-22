# Work orders, swarm templates and service access

Operations of milestone M6. The API decides who may do what; a refusal prints
`An error occurred (AccessDenied) ...` and exits 254.

## Create a work order

```bash
afs work-orders create-work-order --template-id swt_01J... --title "Patch CVE-2026-1234" \
  --kind cve_remediation --cve-id CVE-2026-1234 --severity critical --sla-due-at 2026-12-01T00:00:00Z \
  --repos github.com/acme/ledger github.com/acme/billing=team_01J... github.com/acme/gateway=ops@example.com
```

Each `--repos` item is `host/owner/repo` with an optional `=owners`:

| Item | Sent |
|---|---|
| `github.com/acme/ledger` | no owner; the team configuration decides (the answer shows `ownerSource: team_config`) |
| `...=team_01J...` | `owningTeamId` |
| `...=ops@example.com` | `owningContact` |
| `...=team_01J...,ops@example.com` | both |
| `...=payments` | refused before any call (exit 252): use a team id or an email address |

`afs teams list-team-directory` lists team ids (service-level roles only). When the service cannot choose an owner, the
error names the repositories (`Ambiguous repository ...` with the teams, or `Unowned repositories ...`) and says to give
`=team_...` or `=email`; exit 254. A second open work order for the same CVE and repository is `Conflict`; a disabled
template is a `ValidationError`. The creating key is a security-operator key in `AFS_API_KEY`; use `--idempotency-key`
to retry safely.

```bash
afs work-orders list-work-orders --status open --severity critical --output csv
afs work-orders describe-work-order --work-order-id wo_01J...
afs work-orders list-repos --work-order-id wo_01J... --output table
afs work-orders update-work-order --work-order-id wo_01J... --sla-due-at 2026-12-20T00:00:00Z --add-repos github.com/acme/x=team_01J...
afs work-orders cancel-work-order --work-order-id wo_01J... --reason "CVE withdrawn"
afs work-orders list-team-work-orders --team-id team_01J...
```

## Report progress from a worker

A worker with only the projected service-account token needs no setup: the token file named by `AFS_API_TOKEN_FILE`
(default `/var/run/secrets/afs/api/token` when it exists) is the last credential tried and is re-read for every request.
`afs sso get-caller-identity` shows `type: workload`.

```bash
afs work-orders update-repo --work-order-id wo_01J... --repo github.com/acme/ledger \
  --status pr_opened --pr-url https://github.com/acme/ledger/pull/42 --pr-number 42 --if-match 3
```

`--repo host/owner/repo` fills the three path segments (`--host`, `--owner` and `--repo ledger` also work). Without
`--if-match` the CLI reads the repository's version first; a workload principal cannot make that read today, so the
call ends with the API's `AccessDenied` and `Pass --if-match <version> (the version printed by the previous
update-repo).` (exit 254). Pass the `version` the previous `update-repo` printed.
A workload token cannot cancel a work order (`AccessDenied`).

## Swarm templates (service admin)

```bash
afs swarm-templates create-swarm-template --name cve-remediation --kind cve_remediation --budget-usd 250.00 \
  --ttl-hours 72 --teardown-rule at_work_complete --spec file://swarm-spec.json
afs swarm-templates disable-swarm-template --template-id swt_01J...
afs swarm-templates enable-swarm-template --template-id swt_01J...
afs swarm-templates delete-swarm-template --template-id swt_01J...
```

Deleting a template that an open work order uses is a `Conflict` listing the work orders under `Dependents`.

## Service-level keys and security operators

```bash
afs admin create-api-key --name cve-pipeline --role security_operator --expires-at 2027-02-19T18:00:00Z
afs admin list-api-keys
afs admin add-security-operator --subject-type user --email ana@example.com --identity-provider-id idp_01J...
afs admin list-security-operators
```

`create-api-key` prints the `secretKey` once; store it then. A replayed create carries no secret.
The full option list of every command is in [reference/work-orders.md](reference/work-orders.md),
[reference/swarm-templates.md](reference/swarm-templates.md) and [reference/admin.md](reference/admin.md).
