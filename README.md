# enterprise-agent-factory-service-cli

`afs`, the command line for the Enterprise Agent Factory Service. It follows AWS CLI conventions:
`afs <service> <operation> --options`, named profiles, `--output json|table|text|yaml`, `--query`,
pagination, and AWS CLI exit codes. It talks only to the service's REST API, whose role-based access
decides what each user can do.

Part of [enterprise-agent-factory-service](https://github.com/stainedhead/enterprise-agent-factory-service).

**Status:** scaffolded. `afs help` and `afs --version` work; service operations come next.

## Examples (target design)

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
make gate     # tidy, fmt-check, vet, golangci-lint, race tests, build
./bin/afs help
```

Contributor and agent rules: [AGENTS.md](AGENTS.md). Why this repository exists: [INTENT.md](INTENT.md).
