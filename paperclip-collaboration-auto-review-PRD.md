# Paperclip collaboration: code and design review fixes

**Source:** automated review of branch `claude/project-thread-xqpqky` against `specs/261122-paperclip-collaboration` (release v1.0.1).
**Scope:** `internal/edge`, `internal/cli/collaboration`, `internal/cli/reports`, `internal/cli/runtime`, `internal/cli/parse.go`.

## 1. Executive summary

The change is in good shape. Clean Architecture holds (`edge` imports only `client`, `config`, `redact`; handlers do no business logic),
every new command has examples, columns, goldens and an `AccessDenied` test, the gate is green at 94.7% coverage, and the
credential rules (key from file or stdin only, masked headers, no query string or body in debug, no redirects or cookies)
are proven by sentinel tests. No P0 finding. The review found one behavioural inconsistency in retries and a few test gaps.

## 2. Requirements

### FR-R1 (P1) Retry-After longer than the cap must not be waited for
`internal/edge` caps a `Retry-After` at 20 s (or the read timeout) and then waits that long, while the API client (`client.Do`) refuses
to retry when the server asks for longer than the cap. The two must agree: a pipeline should fail fast, not wait on a hint it cannot honour.
- Acceptance: given a `GET` answered `503` with `Retry-After: 9999`, `invoke-api` makes exactly one request, does not sleep, prints the
  edge error and exits 254. A header at or under the cap is still honoured. *Proven by* `edge.TestDoRetryAfterOverCapNotRetried`.

### FR-R2 (P2) The runtime credential helper has its own tests
`runtime.Credentials` and `runtime.Transport` are exercised only from the `collaboration` package, so `runtime` shows 24% coverage.
- Acceptance: `runtime` tests cover `Credentials` for an API key from the environment, a missing credential, an unknown profile (error), and `Transport` with an injected round tripper. *Proven by* `runtime.TestCredentials`, `runtime.TestTransportUsesInjected`.

### FR-R3 (P2) Unknown-option suggestions are unit tested
`didYouMean` in `internal/cli/parse.go` is proven only through `invoke-api`.
- Acceptance: a parser test shows that an unknown option that is a prefix of a known one names it, one that is not names nothing, and a value after `=` is never echoed. *Proven by* `cli.TestUnknownOptionSuggestsAndNeverEchoesValue`.

### FR-R4 (P2) A key file that is a directory or unreadable gives a clear exit 255
- Acceptance: `--paperclip-key-file <dir>` exits 255 with a message that names the path and never any key text. *Proven by* `edge.TestReadKeyDirectory`.

## 3. Deferred (not fixed in this change)

| Item | Rationale |
|---|---|
| Loopback `http://` origins are exempt from the port, address and domain rules | Matches `endpoint_url` (ADR-012) and is needed for tests; documented in ADR-015 |
| `invoke-api` help lists global options that it refuses | Help is generated from the shared global table; the refusal message explains; changing help generation is out of scope |
| End-to-end proof against a deployed `afs-edge` | Parked (real cluster and AWS) |

## 4. Guidance for the fixes

Use TDD (Red, Green, Refactor) for every fix; review each fix against `AGENTS.md` before the next; fixes are small and independent,
so agent teammates and git worktrees may be used in parallel but this change is small enough to do in sequence. Priorities:
P0 blockers first, then P1, then P2.
