# Architecture: model-endpoint-gateway

No new layer. The generic executor still turns one command into one API call.

- `internal/cli/modelendpoints/` (new service package): overrides (examples, columns, notes, Build hooks) for the 23
  operations whose `x-afs-cli` service is `model-endpoints`. Registered in `internal/cli/services`.
- `internal/cli/teams/gateway.go`: `put-team-gateway-settings`, `put-team-model-endpoints`.
- `internal/cli/reports/inference.go`: the two inference reports (dynamic columns, range and group-by checks).
- `internal/cli/modelaccess`: `--routes` and `--default-route` fold into the `models` member.
- Core changes, each small and tested:
  - `cli.DeriveOptions`: a body member that has the name of a query parameter of the same operation is not a second
    option (the query option stands; the override places the value).
  - `cli.RequestParts.Raw` and `ContentType`, carried by `exec.Built` and sent through the new
    `client.Request.ContentType`: a Build hook may send a non-JSON body (the provider statement CSV).
  - `safeupdate`: `UpdateInferenceProvider` reads its version from the non-paginated provider list;
    `PutTeamModelEndpoints` and `PutTeamModelEndpointCredential` are versionless.
  - `exec/hints.go`: the v1.5.0 reasons.
- Secrets: the credential bodies use the `SecretValue` schema, so the generic `--secret-value` option applies
  (`internal/secretinput`: file, stdin or prompt; never a plain argument); `redact` masks `value` in debug logs.
