# enterprise-agent-factory-service-cli

`afs`, the command line for the Enterprise Agent Factory Service. It follows AWS CLI conventions:
`afs <service> <operation> --options`, named profiles, `--output json|table|text|yaml`, `--query`,
pagination, and AWS CLI exit codes. It talks only to the service's REST API, whose role-based access
decides what each user can do.

Part of [enterprise-agent-factory-service](https://github.com/stainedhead/enterprise-agent-factory-service).

**Status:** milestones M0 (foundation), M1 (teams, members, team audit, settings, identity providers), M2 (agents, secrets, configuration values, external systems, catalogs), M3 (budgets, usage, reports and CSV, notification channels, models and prices) and M4 (API keys, role links, identities, tags, log destinations, restore, report schedules, service administration, exports and imports, `afs completion`, release configuration) done. All 180 operations of the pinned API contract are callable through
the generic executor (`--cli-input-json`, `--generate-cli-skeleton`, pagination, safe updates with `If-Match`,
`--wait`, retries and idempotency keys); `afs configure`, `afs configure sso`, `afs sso login`, `afs sso logout`,
`afs sso get-caller-identity` and `afs sso list-providers` are finished. Per-service polish (examples, table
columns, golden files, `AccessDenied` tests) lands milestone by milestone.

## Examples

```bash
afs configure sso                          # issuer, client id, endpoint for a profile
afs sso login --profile prod
afs teams list-teams --output table
afs agents list-agents --team-id team_01J... --query "items[?status.phase=='Running'].name"
afs budgets put-budget --team-id team_01J... --amount-usd 40000 --period monthly --thresholds 50 80 100
afs secrets create-secret --team-id team_01J... --name github-bot-key --secret-string file://key.pem
```

## Build and test

```bash
make gate     # tidy, contract-check, generate-check, fmt-check, vet, lint, race tests, coverage (90%), build
./bin/afs help
```

Other targets: `make generate` (regenerate descriptors and types from `api/openapi.yaml`),
`make contract-sync RESTAPI_REF=<ref>` (re-pin the API snapshot), `make docs` (regenerate `user-docs/reference/`).

Releases are built by GoReleaser (`.goreleaser.yaml`, `.github/workflows/release.yml`): archives for five platforms, a
cosign-signed `checksums.txt`, and notarization, Authenticode signing and the Homebrew tap as soon as their secrets exist
(see [install](user-docs/install.md)). `afs completion bash|zsh|fish|powershell` prints a completion script.

## Documentation

[user-docs/](user-docs/README.md): install and release verification, getting started, configuration, single sign-on, credentials, output,
pagination, errors and exit codes, troubleshooting, and the generated command reference.

Contributor and agent rules: [AGENTS.md](AGENTS.md). Why this repository exists: [INTENT.md](INTENT.md).
