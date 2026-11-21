# Research: `afs` CLI, feature complete

| | |
|---|---|
| **Feature** | cli-feature-complete |
| **Date** | 2026-11-21 |
| **Source PRD** | [enterprise-agent-factory-service-cli-PRD.md](enterprise-agent-factory-service-cli-PRD.md) |

## Research questions

| # | Question (from PRD §6, §8 and the contract) | Answer | Status |
|---|---|---|---|
| RQ-1 | Can oapi-codegen generate Go types from the OpenAPI 3.1 snapshot (55 `type: [T, "null"]` unions, 24 `oneOf`)? | Yes. Spike on 2026-11-21: oapi-codegen v2.8.2, `generate: {models: true}`, produced 7,437 lines; `go vet` clean with `oapi-codegen/runtime` v1.9.1. Nullable unions become `*T` with `omitempty`, so an explicit `null` is lost on re-encode (see RQ-2) | Answered |
| RQ-2 | Can `--output json` print "exactly the API's fields" if responses are decoded into generated types? | No: typed decoding drops unknown fields from a newer minor and explicit `null`s (`mergedAt: null` in the `UpdateWorkOrderRepo` example), and reorders keys. Decision D-C5: `json` re-indents the raw body; other renderers decode into an ordered generic tree | Answered |
| RQ-3 | Where does the command table come from? | Each of the 182 operations carries `x-afs-cli` (`"<service> <operation>"`, or `none (not a CLI operation)` for `GetHealth`, `GetReadiness`), `x-afs-csv`, `x-afs-priority`, `x-afs-permission`, `x-afs-scope`, `x-afs-sdk`, `x-afs-screen`, `x-afs-events`. The 180 `x-afs-cli` values equal the PRD §3.1 table exactly (checked by script 2026-11-21: no differences beyond the local commands) | Answered |
| RQ-4 | Which operations paginate and what is the item key? | 52 operations take `maxResults`/`nextToken`. 48 return `items`; `GetTeamReport`, `GetAgentReport`, `GetServiceReport`, `GetUsage` return `rows` (with `range`, `groupBy`, `totals`, `notes`). Generator rule: `rows` when present, else `items` | Answered |
| RQ-5 | Which operations need `If-Match`, and is it required? | 59 operations: parameter `IfMatch` (required: every `PATCH` and `DELETE`) or `IfMatchPut` (optional: every `PUT`, CD-2). `428` when required and absent. `PutTeamModels`' resource (`TeamModels`) has no `version`, so no `If-Match` can be sent for it | Answered |
| RQ-6 | Does a `412` carry the current version, so the CLI can print "The resource changed (version 8)"? | No: the `PreconditionFailed` response has no `ETag` and `ErrorDetails` has no version field. The CLI reads the resource after a `412` (D-C8) and requests an `ETag` on `412` (CLI-CD-3) | Answered, delta requested |
| RQ-7 | Which `CREATE`/action operations accept `Idempotency-Key`? | Every `POST` (44 operations, including `CreateImport` and the `:verb` actions). None of the `PUT`/`PATCH`/`DELETE` operations | Answered |
| RQ-8 | Which operations answer CSV? | 8: `GetUsage`, `ListCosts`, `GetTeamReport`, `GetAgentReport`, `GetServiceReport`, `ListAuditEvents`, `ListAdminAuditEvents`, `ListWorkOrders` (`x-afs-csv: true` and a `text/csv` response). Matches contract §4 *CSV* marks | Answered |
| RQ-9 | What is the `CreateImport` request? | `POST /v1/imports?dryRun=true|false` (default `true`), body `ImportRequest {bundle (object, required), accountMap (object)}`, up to 100 MB. The PRD's "the bundle file is the request body" is read as "the bundle file is the `bundle` member of the body" (D-C18) | Answered |
| RQ-10 | How does the CLI get the device-grant client? | `ListSignInProviders` (public) returns `id`, `kind`, `name`, `issuerUrl`, `consoleClientId`, `cliClientId`, `isDefault` (D48). The device authorization and token endpoints come from the issuer's OIDC discovery document (`/.well-known/openid-configuration`). Note CD-9: `-restapi` cannot yet serve distinct client ids until `-db` adds a column (root commit `5c579ae` maps them to `client_id`/`cli_client_id`, `-db` 000017) | Answered; depends on `-db`/`-restapi` |
| RQ-11 | What contract version does the CLI report? | The snapshot's `info.version` is `1.1.0` while the root contract document is v1.2.0; the `X-Afs-Contract-Version` header example is `1.1.0`. The CLI takes `info.version` from the snapshot (OpenAPI wins) and requests the bump (CLI-CD-2) | Answered, delta requested |
| RQ-12 | How are report tag filters sent? | Contract §4.9: query parameters `tag.<key>=<value>`. The OpenAPI does not declare them (only the `GroupBy` description mentions them), so request validation in the replay test must allow `tag.*` (CLI-CD-4) | Answered, delta requested |
| RQ-13 | Can a workload principal read a work-order repository to learn its version before `UpdateWorkOrderRepo` (a `PATCH`, `If-Match` required)? | No: the workload principal holds only `work-orders:report`; `ListWorkOrderRepos` needs `work-orders:read`, and there is no single-repository `GET`. CLI-US-10 cannot be met with read-then-write as specified (CLI-CD-1; D-C17) | Open on `-restapi`/root (OQ-S2) |
| RQ-14 | Which signing toolchain meets CLI-15 without proprietary tooling? | GoReleaser v2 (OSS): `signs` (cosign `sign-blob` keyless on `checksums.txt`), `notarize.macos` (quill, Apple Developer ID), Windows Authenticode through a `signs` step with `osslsigncode`, Homebrew via the tap repository. Verified by a release-verify job (cosign `verify-blob`, `osslsigncode verify`, `spctl --assess` on a macOS runner) | Answered; secrets pending (OQ-S1) |

## Existing implementations

- This repository: scaffold only (`internal/cli/cli.go` with the §11 service list, exit-code constants,
  `buildinfo`). Nothing to preserve beyond the exit codes and the Makefile/CI shape.
- `-restapi` bundles `api/src` into `api/openapi.yaml` with `tools/openapibundle` and parses it with
  `pb33f/libopenapi` v0.44.0; its `api/contract_test.go` validates requests and responses with
  `libopenapi-validator` v0.16.5. The CLI's replay test uses the same validator version so both sides
  read the contract identically.
- `-restapi` `api/contract/operations.json` lists `{operationId, method, path, permission, scope, csv}`
  — the same facts the CLI generator reads from the YAML extensions; the CLI does not copy it (one
  source: the YAML).

## API documentation notes (from the snapshot)

- Paths include the `/v1` prefix (`/v1/teams/{teamId}`); `GetHealth`/`GetReadiness` are `/healthz`,
  `/readyz`. The profile's `endpoint_url` is therefore the host root (`https://afs.example.com`).
- Security schemes: `idpBearer`, `apiKeyBearer` (`afsk_…`), `workloadBearer`; all `Authorization: Bearer`.
- Error envelope `{code, message, requestId, details}`; `ErrorCode` enum has the 10 `afs-api` codes;
  `ErrorDetails` fields: `ambiguousRepos`, `dependents`, `fields`, `missingTags`, `reason`,
  `requiredPermission`, `unownedRepos`.
- `Money` is a decimal string. `Setting.value` and `PutSettingRequest.value` have no type (any JSON).
- `GetTeamKubeconfig` answers `application/yaml`.
- Delete status codes (CD-1): `202` with a body for `DeleteTeam`, `DeleteAgent`, `DeletePool`,
  `DeleteLogDestination`, `DeleteNotificationChannel`; `204` otherwise.
- `TestNotificationChannel`: `200` `{testId, delivered, error}` or `202` `{testId, pending: true}` (CD-3).
- `UpdateWorkOrderRepo` path is `/v1/work-orders/{workOrderId}/repos/{host}/{owner}/{repo}` (three path
  segments, each escaped).

## Best practices adopted

- AWS CLI v2 conventions for `--cli-input-json` (parameters and body members in one object, explicit
  options override it), `--generate-cli-skeleton` (prints and exits, no call), `NextToken` for
  `--max-items`, `text` output (tab-separated), pager (`AFS_PAGER`/`PAGER`, default `less -R` on a TTY).
- `golden` files with an `-update` flag; `httptest` servers for every request; fake clocks for retries
  and `--wait`.

## Open questions

| # | Question | Owner | Proposal |
|---|---|---|---|
| OQ-S1 | Signing material: Apple Developer ID (notarization), Windows Authenticode certificate, cosign keyless identity (GitHub OIDC on the release workflow), Homebrew tap repository `stainedhead/homebrew-afs` | J | Build the pipeline with `--snapshot` dry runs now; enable signing jobs when CI secrets exist; link the tap only once it exists |
| OQ-S2 | CLI-CD-1: how does a workload principal satisfy `If-Match` on `UpdateWorkOrderRepo`? | `-restapi` / root | Grant workload principals `ListWorkOrderRepos` on their assigned work order (preferred), or return the repository `version` to them another way; CLI fallback D-C17 |
| OQ-S3 | When will `-restapi` merge `feat/restapi-feature-complete` and tag a release to re-pin to? | `-restapi` | M6-X2 re-pins to the tag; until then each milestone re-syncs to the branch head (Mn-0) |
| OQ-C6 | Keychain storage | J | Files in v1 (PRD) |
| OQ-9 | IAM role for viewer kubectl | J (root) | No CLI impact: the CLI merges the kubeconfig exactly as the API returns it |

## References

- `-restapi/api/openapi.yaml` @ `081de14`; `-restapi/specs/261121-restapi-feature-complete/implementation-notes.md`
- Root `docs/contracts/api-contract.md` v1.2.0 (§1 conventions, §2 auth, §4 operations, §6 errors, §8 change process)
- RFC 8628 (OAuth 2.0 Device Authorization Grant); OpenID Connect Discovery 1.0
- AWS CLI v2 user guide: global options, pagination, `--cli-input-json`, output formats
