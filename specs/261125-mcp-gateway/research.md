# Research

- 276 operations; `CompleteMcpOAuthConnect` has no `x-afs-cli` (the browser redirect of an OAuth connect), so 273
  commands. The command table test skips it by name like `GetHealth`.
- oapi-codegen collisions: the inline enum `McpPrimitive.change` is named `McpPrimitiveChange`, which is also a schema;
  the overlay names it `McpPrimitiveChangeKind`. `McpUsageReport.windowDays` is a nullable integer enum (30, 90, null)
  that has no Go enum type; the overlay types it as `int`.
- `McpCredentialRequest` has two write-only secrets (`value`, `clientSecret`) and they are not `SecretValue` refs.
  `value` maps to `--secret-value` as elsewhere; `clientSecret` is sensitive by name. `ceilingOverrideReason` is also
  write-only but not secret, so write-only alone does not make a member sensitive.
- `GetTeamMcpServer` answers a wrapper whose version is `server.version`; the policies, credentials and a team's switch
  have no readable version. Federated identities have no GET of their own; their list is not paginated and items are
  matched by `name`.
