# Architectural decision record

Decisions local to this repository. Cross-repo decisions live in the root repository.

## ADR-001 Go for the CLI

- Status: Accepted (J, 2026-11-20)
- Context: the AWS CLI is Python, but users need a single binary for macOS, Linux and agent containers.
- Decision: Go, static binary (`CGO_ENABLED=0`), same toolchain as the API.
- Consequences: AWS CLI conventions are reproduced deliberately (options, output, exit codes), not inherited.

## ADR-002 Client generated from the OpenAPI contract

- Status: Proposed
- Context: the API contract lives in `-restapi/api/openapi.yaml`.
- Decision: generate the HTTP client and types from a pinned copy of the contract; hand-write commands.
- Consequences: a contract change is a regeneration PR here; drift is caught by a CI check against the pinned version.
