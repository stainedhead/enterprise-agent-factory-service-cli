# Research: paperclip-collaboration

**Date:** 2026-11-22 · **Source PRD:** [paperclip-collaboration-PRD.md](paperclip-collaboration-PRD.md)

## Research questions

1. How much of the surface does the generator already give (flags, enums, columns), and what needs overrides?
2. How does the executor print `details.*` and where does the reason hint table plug in?
3. How do the safe-update reader and wait tables handle a singleton resource (`GetTeamCollaboration`)?
4. How do local commands (configure, sso login) get the credential chain, so `invoke-api` can reuse it?
5. How does `redact` and the client debug writer treat headers, and what has to change for `X-Afs-Authorization`?

## Industry standards

AWS CLI conventions (`afs <service> <operation>`, exit codes 252/254/255, `--cli-input-json`), RFC 9110 for HTTP
semantics (no redirect following, `Retry-After`), OWASP guidance on SSRF-style host validation for the origin rules.

## Existing implementations

Findings are recorded in the Findings section as the research phase of step 3 completes.

## API documentation

`api/openapi.yaml` pinned from restapi 99bfcfc (info.version 1.3.1): operations `GetTeamCollaboration`,
`PutTeamCollaboration`, `ListCollaborationLinks`, `ConnectAgentToCollaboration`, `DisconnectAgentFromCollaboration`,
`RotateCollaborationKeys`, `ListCollaborationEndpoints`; edge contract in root `docs/contracts/api-contract.md` §4.4c.

## Best practices

Table-driven tests, golden files, sentinel tests for secrets, `httptest` fakes, generated files never hand-edited.

## Open questions

PC-CLI-OQ5 (both surfaces required on put) is decided: required. All others resolved by contract v1.3.1.

## Findings

- The generic pipeline (enum checks, boolean `--no-` flags, generation-aware wait, If-Match readers, CSV, pagination) already satisfied most of the PRD; the work was overrides, columns, hints, one new local command and tests.
- The pinned snapshot (99bfcfc) lacks the collaboration report kind and groupBy dimensions and the `collab.*` setting keys (RCF-5, RCF-6); open enums keep the CLI free of server lists.
- `invoke-api` needs no new dependency: Go's `http.Transport.RoundTrip` gives no-redirect and no-cookie behaviour by construction.
- A free-form object (`spec`) has no nested properties in the generated descriptors, so skeleton members are supplied per override (D-21).

## References

Root `docs/contracts/api-contract.md`, `paperclip-team-collaboration-PRD.md`, `docs/solution-design.md` §11.
