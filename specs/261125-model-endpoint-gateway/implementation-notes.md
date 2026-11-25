# Implementation notes: model-endpoint-gateway

## Decisions

- **Service names from the snapshot.** `x-afs-cli` puts 23 operations in `model-endpoints`, two in `teams` and two in
  `reports`; the CLI follows it unchanged (the contract is the single source of names).
- **One option for a name that is both query and body** (`ImportProviderStatement`'s `period`, `source`): the body
  member is not derived; the override's Build moves the values into a JSON body or leaves them in the query for CSV.
  No other operation in the snapshot has such a pair (the option-collision test would fail).
- **Non-JSON request body.** `RequestParts.Raw` and `ContentType` (carried by `exec.Built`, sent by
  `client.Request.ContentType`). Only `import-provider-statement` uses it. Required body members are not checked when
  a raw body is sent; the override checks what CSV needs (`--period`, `--source`).
- **`--statement` is a file only.** Inline text and `-` are refused by Precheck before anything is read; files over
  16 MiB are refused by size. The first non-space byte chooses JSON object, JSON array or CSV.
- **Safe updates.** `UpdateInferenceProvider` reads its version from the non-paginated `ListInferenceProviders`
  (matched on `provider`). `PutTeamModelEndpoints` and `PutTeamModelEndpointCredential` are versionless: the generic
  rule would have read the team's or the endpoint's version, which is not the version of what they write. Five
  operations without a reader are listed by name in `TestReaderTableCoversAllIfMatchOperations`; each takes an optional
  `If-Match`.
- **Credentials** use the existing `SecretValue` handling: `--secret-value` reads `file://`, `fileb://` or `-`, or
  prompts without echo; `--value` is an unknown option (exit 252, value not echoed).
- **Dry run.** `update-inference-provider --dry-run` is the derived `dryRun` query option; an After hook prints
  "Dry run: nothing was changed." and the impact on stderr for table and text output.
- **Report group-by text.** The snapshot's description of `groupBy` on the usage reports is shortened by the generator
  ("rep..."); the reports override gives the full text.
- **`--ca-bundle`** is read as `Any` so it takes PEM text or `file://PATH`.

## Gotchas

- The container restarted during T3. The uncommitted work was intact on disk; from then on each step was committed and
  pushed.
- `make vuln` locally builds govulncheck with an older toolchain; with `GOTOOLCHAIN=go1.27.4` it then needs
  `vuln.go.dev`, which the sandbox proxy blocks. CI runs it.
- gosec G101 flags the help text constant about credentials; it carries a `nolint` with the reason.

## Review (AGENTS.md checklist)

- Naming and options follow the AWS CLI rules; every command has help and examples.
- `--output json` is the API's body unchanged (no new rendering of JSON).
- No value can reach stdout, stderr, logs or history: tests cover file, stdin, prompt, plain argument, `--value` and
  `--debug`.
- Errors are AWS-style with exit 254 (`AccessDenied` per command); usage errors exit 252 before any call.
- Context passes through; no new network calls beyond the sanctioned version read.
- No hardcoded endpoints, tenant ids or credentials (examples use contract-shaped placeholder ids).
