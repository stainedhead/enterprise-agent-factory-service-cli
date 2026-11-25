# Implementation notes

- The `mcp` package overrides operations in six services (`mcp-servers`, `teams`, `agents`, `pools`, `audit`,
  `reports`); overrides are keyed by operation, so the package boundary follows the feature, not the service.
- `put-team-mcp-server --enabled` is required by the contract; `--no-enabled` switches a server off.
- `start-mcp-oauth-connect` prints the authorization URL and its expiry only; the state is in the URL.
- The approval refusal test checks the AWS-style message; the hint for `approval_not_permitted_by_policy` names
  `list-mcp-reviews`.
