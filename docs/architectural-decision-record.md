# Architectural decision record

Decisions local to this repository. Cross-repo decisions live in the root repository.

## ADR-001 Go for the CLI

- Status: Accepted (J, 2026-11-20)
- Context: the AWS CLI is Python, but users need a single binary for macOS, Linux and agent containers.
- Decision: Go, static binary (`CGO_ENABLED=0`), same toolchain as the API.
- Consequences: AWS CLI conventions are reproduced deliberately (options, output, exit codes), not inherited.

## ADR-002 Client generated from the OpenAPI contract

- Status: Accepted (2026-11-21, M0)
- Context: the API contract lives in `-restapi/api/openapi.yaml`.
- Decision: generate the HTTP client and types from a pinned copy of the contract; hand-write commands.
- Consequences: a contract change is a regeneration PR here; drift is caught by a CI check against the pinned version.

## ADR-003 One generic executor, per-operation overrides

- Status: Accepted (2026-11-21, M0)
- Context: the contract has 180 CLI operations; hand-written handlers would drift and collide.
- Decision: every operation is runnable from a generated descriptor by one executor; a service package adds
  examples, table columns and shorthands as an override merged per operation. A test requires every operation
  at or below the current milestone to be finished.
- Consequences: new operations from a contract sync work at once; polish is data, not code.

## ADR-004 Hand-written argument parser

- Status: Accepted (2026-11-21, M0)
- Context: AWS semantics (global options anywhere, `--opt=value`, list values up to the next option) and a
  generated option table; startup must stay under 100 ms.
- Decision: a small parser over the option table instead of cobra or pflag.
- Consequences: help, completion and docs come from the same table; no framework to track.

## ADR-005 Raw JSON body for `--output json`

- Status: Accepted (2026-11-21, M0)
- Context: scripts need the API's field names, order, nulls and unknown fields of a newer minor version.
- Decision: print the response body re-indented; decode into an ordered tree only for yaml, table, text and
  `--query`.
- Consequences: `json` output is the API's, byte-stable apart from indentation.

## ADR-006 Work-order repositories are classified on the client; the version read is composite

- Status: Accepted (2026-11-22, M6)
- Context: `--repos` takes `host/owner/repo[=owners]`; the API wants `WorkOrderRepoInput` objects. `UpdateWorkOrderRepo` needs
  `If-Match`, but a repository has no single `GET`, and a workload principal cannot list a work order's repositories.
- Decision: the CLI splits and classifies owners (`@` is an owning contact, `team_` an owning team id, anything else is a
  usage error, exit 252) and sends no owner when there is no `=`, leaving the team configuration to the service.
  Without `--if-match`, `update-repo` reads `ListWorkOrderRepos` and matches `host` and `owner/repo`; when that read is
  denied it prints the API's answer and `Pass --if-match <version> ...`, exit 254. The workload token is the last
  credential in the chain and is re-read for every request.
- Consequences: no client-side ownership decisions beyond the syntax; workers must pass `--if-match` until the API grants
  the read (contract delta CLI-CD-1).


## ADR-007 Own descriptor generator; oapi-codegen for types only

- Status: Accepted (2026-11-21, M0)
- Context: operation descriptors carry `x-afs-*` facts that oapi-codegen does not model, and tool dependencies would
  widen `govulncheck`'s scope.
- Decision: a dedicated generator (`tools/afsgen`) writes the descriptors; oapi-codegen v2.8.2 is used for types only and
  runs from the Makefile, not from `go.mod`. `make generate-check` fails on drift.
- Consequences: the runtime module graph holds only what the binary uses; a contract change is `make contract-sync` and a
  regeneration diff.

## ADR-008 Safe updates and `--wait` are tables, checked against the contract

- Status: Accepted (2026-11-21, M0)
- Context: 59 operations require `If-Match`, and many mutations are asynchronous.
- Decision: a reader table maps each `If-Match` operation to the read that returns its version; a wait table maps each
  waitable mutation to its status read. Tests check both tables against the snapshot. A wait ends at `Ready` with
  `observedGeneration` at or above the mutation's `generation`, so a stale `Ready` does not end it.
- Consequences: a new update or async operation in the contract fails a test until it has a table entry.

## ADR-009 Credentials in files, not the OS keychain, in v1

- Status: Accepted (2026-11-21, M0; OQ-C6)
- Context: `afs` runs on laptops, in CI and in minimal agent containers with no keychain.
- Decision: API keys in `~/.afs/credentials` and SSO tokens in `~/.afs/sso/cache/`, both required to be 0600; the
  workload token is read from its mounted file for every request. No `GetMe` cache (D-C12): `AccessDenied` already names the
  required permission.
- Consequences: one code path on every platform; keychain storage can be added later as an extra credential source.

## ADR-010 CLI pagination token wraps the server token

- Status: Accepted (2026-11-21, M0)
- Context: `--max-items` can stop mid-page, which a server `nextToken` cannot express.
- Decision: the printed `NextToken` is `base64url({"t": <server token>, "s": <items to skip>})`. `--max-items` with
  `--output csv` is a usage error because CSV pages cannot be cut without re-parsing.
- Consequences: `--starting-token` resumes exactly; the token is opaque to users, as in the AWS CLI.

## ADR-011 Signed releases with GoReleaser and keyless cosign

- Status: Accepted (2026-11-21, M4)
- Context: CLI-15 requires signed binaries on five platforms and a Homebrew tap without long-lived signing keys in CI.
- Decision: GoReleaser v2; `checksums.txt` signed keyless with cosign from the release workflow's OIDC identity; quill for
  macOS notarization; `osslsigncode` for Authenticode; jobs gated on their secrets.
- Consequences: releases are verifiable with `cosign verify-blob` against the workflow identity; notarization, Authenticode
  and the tap switch on when their secrets are provisioned.

## ADR-012 Service and identity-provider data is untrusted

- Status: Accepted (2026-11-22, review fixes)
- Context: the review of the feature-complete branch found that provider fields written to `~/.afs/config` could inject
  lines into the file, that OIDC endpoints were not required to use https, and that service text reached the terminal raw.
- Decision: every value that comes from the service or the identity provider is validated at the edge before it is stored,
  sent or printed. `config.File.Set` refuses control characters and structural characters; provider `issuerUrl` must be an
  https URL (loopback http allowed) and ids use a strict character set; OIDC discovery must name the configured issuer and
  use https endpoints; `table`, `text`, error output and prompts pass through `output.Sanitize`; endpoint URLs may not carry
  credentials. `--output json` keeps the service's bytes (ADR-005).
- Consequences: a hostile or misconfigured service fails a command with exit 255 instead of changing local state. DEL and C1
  characters can still pass through `--output json`.

## ADR-013 Safe-update notice after the write

- Status: Accepted (2026-11-22, review fixes; supersedes the wording of D-C8 as built)
- Context: PRD CLI-10 specifies `Applied to version N; use --if-match to guard against concurrent changes.` after a
  successful write. The build printed a different sentence before the write.
- Decision: the CLI prints the PRD text on stderr after the write succeeds, only when it read the version itself. No notice
  appears when the write fails or when `--if-match` is given.
- Consequences: matches the root PRD, so no deviation is reported to the root PRD owner.

## ADR-014 Generated types are compile-time references for the fields the CLI reads

- Status: Accepted (2026-11-22, review fixes FR-017; amends ADR-002 and ADR-007)
- Context: `internal/client/gen` (oapi-codegen types) was imported by nothing but its own test, so ADR-002's "the client is
  generated from the contract" was true only of the descriptors. Commands decode into ordered trees or `json.RawMessage`
  because `--output json` must keep the API's bytes (ADR-005) and unknown fields must survive.
- Decision: keep the types, and make them useful where the CLI reads response members itself: the SSO provider list, the
  `--wait` status evaluator and the safe-update version readers reference each member they read (`contract.go` in
  `internal/cli/sso`, `wait` and `safeupdate`). A contract rename of such a member fails the build. Decoding for output
  still goes through trees. Removing the package was considered (OQ-5) and not done in this change.
- Consequences: `make generate-check` still guards both generated outputs; a rename that the references do not cover is
  caught by the contract replay and example tests.

## ADR-015 invoke-api is a local command with its own narrow HTTP client

- Status: Accepted (2026-11-22, release v1.0.1, spec `261122-paperclip-collaboration`)
- Context: a team's Paperclip REST API is served by the AFS edge (API contract 4.4c), not by `afs-api`, with two
  credentials in two headers and different retry, redirect and error rules than the operations of the contract.
- Decision: `afs collaboration invoke-api` is a local command (no `operationId`) built on `internal/edge`. It uses the
  transport's `RoundTrip` directly (no redirect, no cookie jar), validates the origin (https, port 443, a DNS name, optionally one
  label under the profile's `collab_base_domain`), reads the Paperclip key only from a 0600 file or stdin, refuses workload
  tokens, retries only `GET`/`HEAD` on 429/503 and any method on a failed connection, maps the edge's error envelope to
  AWS-style errors and prints Paperclip's own errors as `HttpError` with the body on stdout. Debug lines carry no query string and no body.
- Consequences: the generic executor stays free of edge rules; the credential never leaves the AFS API and the validated edge origin.

## ADR-016 Open enums for values the pinned snapshot does not list

- Status: Accepted (2026-11-22, release v1.0.1)
- Context: the pinned snapshot lacks the `collab.*` setting keys that contract 1.3 defines (root contract follow-up RCF-5), and the collaboration report is not schedulable.
- Decision: `put-setting --key` and `create/update-report-schedule --report` are open enums (`Override.OpenEnum`): help says "Known values" and any other value is sent for the API to judge. The CLI holds no list of server keys.

## ADR-017 Deprecations are declared on overrides and warned on stderr

- Status: Accepted (2026-11-23, release v1.1.0, spec `261123-orchestration-and-model-access`)
- Context: contract 1.4 replaces the collaboration operations and `--model-id`, and keeps them until 2.0.0.
- Decision: `Override.Deprecated`, `DeprecatedOptions`, `Hidden`, `Desc`, `Enum` and `Notes` describe a deprecation once; help prints it and `exec/deprecation.go` prints one stderr line per call (by option too, also from `--cli-input-json`). When no static line applies, the service's `Deprecation`/`Sunset` headers produce one. Output on stdout is unchanged.
- Consequences: scripts keep working and parse the same stdout; removal in 2.0.0 is a table edit.

## ADR-018 Orchestration exposure is two local options folded into the body; model access is one builder

- Status: Accepted (2026-11-23, release v1.1.0)
- Context: the API takes `webUi`/`webApi` objects and `{access, modelIds}`; one operation maps to one call, with no client-side workflow.
- Decision: `--web-ui-exposure` and `--web-api-exposure` (both required unless `--cli-input-json` supplies them) and `--models all|ids` plus `--default-model-id` are built into the request body by hooks (`orchestration.buildExposure`, `modelaccess.Build`); the JSON-shaped options stay available but hidden. Snapshot gaps (setting keys, agent role enum) are handled with open enums. `output.Column.NAUnless` renders `N/A` for Buzz surfaces.
- Consequences: no policy is predicted client-side; the API's refusal reasons are explained by `exec/hints.go`.
