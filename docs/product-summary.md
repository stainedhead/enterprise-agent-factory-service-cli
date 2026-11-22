# Product summary

`afs` is the command line for the Enterprise Agent Factory Service (AFS). It gives people, pipelines and
agents everything the web console offers, from a shell, with the habits of the AWS CLI:
`afs <service> <operation> [--options]`, named profiles, `--output json|yaml|table|text|csv`, JMESPath
`--query`, token pagination, `--cli-input-json` and `--generate-cli-skeleton`, and AWS CLI exit codes.

## Who uses it

| User | Typical use |
|---|---|
| Service admins | Teams, settings, identity providers, models and prices, exports and imports, security operators |
| Team admins | Agents, members, secrets, budgets, notification channels, pools, work orders |
| Viewers, FinOps, SRE | Usage and cost reports (CSV), audit search, team status, `kubectl` access |
| CI pipelines | API-key credentials, `--output json` as a stable machine contract, idempotent retries |
| Agents and pool workers | The workload-token credential, work-order reporting, handoffs |

## What it does

- **Every API operation is a command.** All 180 operations of the pinned API contract (`api/openapi.yaml`,
  contract 1.1) are callable, plus the local commands `configure`, `configure sso`, `sso login`, `sso logout`,
  `completion` and `help`.
- **Safe by default.** Updates send `If-Match` (read-then-write when you do not pass `--if-match`); creates send an
  `Idempotency-Key` so retries never duplicate; secret values are accepted only from a file, stdin or a no-echo
  prompt; no token, API key or secret value is printed, logged or written to `--debug` output.
- **The API decides access.** The CLI never skips a call because it predicts a `403`; an `AccessDenied` answer names
  the permission that was missing.
- **Single sign-on and API keys.** OIDC device authorization grant with a 0600 token cache and silent refresh; API keys
  from the environment or `~/.afs/credentials`; the workload token for agents running in the platform.
- **Signed releases** for macOS, Linux and Windows, with cosign-signed checksums and a Homebrew tap.

## Status

Feature complete against contract 1.1 (milestones M0 to M6). The contract snapshot is pinned to the `-restapi`
feature branch until that repository publishes a release tag, at which point it is re-pinned.

See [product-details.md](product-details.md) for the service and operation inventory, and the root repository's PRD
and `docs/solution-design.md` §11 for the product context.
