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
