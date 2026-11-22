# Architecture: cli-feature-complete-auto-review

Status: Draft. Created 2026-11-22.

## Overview
No structural change. Command handlers still call one client method and render. New shared helpers:
- `internal/output.Sanitize` (FR-003), used by table, text, errors and the device prompt.
- `config.ValidateEndpoint` is reused for `sso_issuer` and provider `issuerUrl` (FR-001, FR-002).
- `config.File.Set` validates its arguments (FR-001); the single choke point for INI writes.
- `client` gets a `Retry-After` cap and URL redaction (FR-005, FR-011).

## Layer responsibilities
Unchanged: `cmd/afs` -> `internal/cli` -> `client`, `config`, `output`, `sso`, `kubeconfig`.

## Data flow, integration points
Hostile data enters at: provider list (config writes), OIDC discovery and token responses, API
responses (printing), kubeconfig response, pagination headers. Each now has a validation edge.

## Architectural decisions
In implementation-notes.md (D-R1..).
