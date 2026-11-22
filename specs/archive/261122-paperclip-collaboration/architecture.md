# Architecture: paperclip-collaboration

**Date:** 2026-11-22 · **Status:** Draft (confirmed in step 3)

## Architecture overview

The CLI is table-driven: the operation descriptors generated from the pinned OpenAPI give every operation a command,
options and a request builder. Service packages add `cli.Override` registrations (examples, columns, kinds, build and
validate hooks) and local commands. The executor (`internal/cli/exec`) sends, paginates, waits and maps errors.

## Component architecture

```
cmd/afs -> internal/cli (registry, parse, override) -> services (composition root)
   collaboration/   overrides for 6 operations + local invoke-api (origin, key, transport, errors)
   agents|secrets|teams|members|apikeys|admin|reports|swarmtemplates  small overrides
   exec/hints.go    reason hint table + details printing
   safeupdate|wait  reader/wait table entries
   redact           edge headers
   client           debug writer uses redact
```

## Layer responsibilities

Handlers parse flags, call one client method, render. Hints add words, never decisions. No access prediction.

## Data flow

`describe/put/...`: options -> Build -> HTTP (If-Match from reader) -> JSON -> columns. `invoke-api`: credential chain ->
`GetTeamCollaboration` (unless `--collab-url`) -> origin validation -> one request to the edge with two credentials ->
body streamed to stdout.

## Integration points

AFS API (`/v1/...`); the team edge (`https://<host>/...`) for `invoke-api` only.

## Architectural decisions

See implementation-notes.md (PC-CLI-D1..D12 plus decisions made while implementing).
