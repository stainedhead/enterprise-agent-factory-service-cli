# User docs

How to install, configure and use the `afs` CLI.

| Document | For |
|---|---|
| [Getting started](getting-started.md) | Install, point `afs` at a service, sign in, run a first command |
| [Configuration](configuration.md) | Profiles, `~/.afs/config`, environment variables, precedence |
| [Single sign-on](sso.md) | `afs configure sso`, `afs sso login`, refresh and `afs sso logout` |
| [Credentials](credentials.md) | API keys, the credential order, keeping secrets out of history |
| [Teams and members](teams-and-members.md) | Teams, members, team audit events, service settings, identity providers |
| [Agents and secrets](agents-and-secrets.md) | Agents, secret input rules, configuration values, external systems, platform catalogs |
| [Budgets and reports](budgets-and-reports.md) | Budgets, usage, costs, reports and CSV, notification channels, models and prices |
| [Keys, tags, logs, schedules and service administration](delivery-and-admin.md) | API keys, role links, identities, team tags, log destinations, restore, report schedules, service audit, price proposals, exports and imports |
| [Swarms](swarms.md) | Ephemeral teams, stages, pools, scaling, handoffs, pool secrets and role links, cost-tag status |
| [Work orders](work-orders.md) | Work orders, `--repos` owners, reporting from workers with the workload token, swarm templates, the team directory, service keys, security operators |
| [kubectl access](kubectl.md) | `afs teams get-kubeconfig`: merging, `--kubeconfig`, `--aws-profile` |
| [Install](install.md) | Releases, verifying them, Homebrew, shell completion |
| [Output](output.md) | `--output json|yaml|table|text|csv`, `--query`, the pager, colour |
| [Pagination](pagination.md) | Auto-pagination, `--max-items`, `--starting-token`, `--page-size`, CSV |
| [Errors and exit codes](errors-and-exit-codes.md) | Error format, exit codes, retries, `--debug` |
| [Troubleshooting](troubleshooting.md) | Common problems and what to do |
| [Command reference](reference/README.md) | Every service and operation (generated from `afs help`) |

## The rest of the product

These guides cover this repository only. For the service as a whole, start with the
[service overview](https://github.com/stainedhead/enterprise-agent-factory-service#readme) and the
[adoption guide](https://github.com/stainedhead/enterprise-agent-factory-service/blob/main/user-docs/adoption-guide.md).
