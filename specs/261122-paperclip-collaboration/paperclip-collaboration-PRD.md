# Paperclip team collaboration: CLI change PRD (`-cli`)

| | |
|---|---|
| **Status** | Ready for implementation (2026-11-22). Release **v1.0.1** |
| **Author** | J (Enterprise Architecture), drafted with Claude |
| **Source** | Extracted from the living CLI PRD, [`specs/archive/261121-cli-feature-complete/enterprise-agent-factory-service-cli-PRD.md`](specs/archive/261121-cli-feature-complete/enterprise-agent-factory-service-cli-PRD.md) (sections 1.1, 3, 3.2, 4.1, 5, 5.2, 6, 7, 8, 8.1, 9, merged as PR #7). That file stays where it is and remains the complete PRD; this focused copy drives the implementation spec |
| **Parent** | [Change PRD: Paperclip team collaboration](https://github.com/stainedhead/enterprise-agent-factory-service/blob/main/paperclip-team-collaboration-PRD.md); [API contract v1.3.1](https://github.com/stainedhead/enterprise-agent-factory-service/blob/main/docs/contracts/api-contract.md) (§2.1, §3, §4.4, §4.4b, §4.4c, §4.6, §7) |
| **Contract source** | `-restapi` branch `claude/project-thread-xqpqky`, commit `99bfcfc67dace785ba991f2025fa34def040f047` (`api/openapi.yaml`, info.version 1.3.1), pinned until the `-restapi` PR merges; then re-pin SOURCE to the merge commit |
| **Release** | Semver **v1.0.1**, tag `v1.0.1`; versions come from the tag through the Makefile, so no version constant changes |

## 1. Purpose


Every team may now run one collaboration agent (Paperclip), connect its other agents to it, and publish its web
UI and REST API through a governed public endpoint (change PRD §1). Without CLI support a team-admin could not
script any of it, and an automation author had no supported way to call the team's Paperclip API with an AFS
credential beside a Paperclip key. This repository delivers: the `collaboration` service (six operations),
`admin list-collaboration-endpoints`, the flag and column changes on `agents`, `secrets`, `teams`, `members`,
`api-keys`, `admin` and `reports`, and one local command, `afs collaboration invoke-api`. The measure: every
operation of API contract §4.4b has a command with help, example, columns, goldens and an `AccessDenied` test, and a
pipeline can create a Paperclip issue with one command and a team key, with no credential visible in any output
(CLI-US-24 to CLI-US-34). The CLI stays a thin client: access is decided by the API and by the edge, never here.


## 2. Requirements owned (CLI-19 to CLI-27)

Scope is feature complete (root D18): every operation in the API contract gets a command. Each
command's priority is its operation's priority in the contract.

| Id | Requirement | Product reqs | Pri |
|---|---|---|---|
| CLI-19 | The `collaboration` service: `describe-collaboration` (`GetTeamCollaboration`), `put-collaboration` (`PutTeamCollaboration`), `list-collaboration-links` (`ListCollaborationLinks`), `connect-agent` (`ConnectAgentToCollaboration`), `disconnect-agent` (`DisconnectAgentFromCollaboration`), `rotate-collaboration-keys` (`RotateCollaborationKeys`), and `admin list-collaboration-endpoints` (`ListCollaborationEndpoints`), named exactly by the contract's `CLI` column (CLI-1). Priorities per contract v1.3.1 (PC-D28): P0 in C1 describe, put, list-links, connect, disconnect; P1 rotate; P0 in C2 endpoints. The two connect actions and the rotate action send an `Idempotency-Key` (CLI-11) and have no `--wait` | COL-34, COL-10, COL-7 | P0 |
| CLI-20 | `put-collaboration` takes `--web-ui-enabled \| --no-web-ui-enabled` and `--rest-api-enabled \| --no-rest-api-enabled` (both required unless `--cli-input-json` supplies them: the PUT replaces both) and `--registration-mode manual\|managed`. It follows CLI-10 (If-Match from a read of `GetTeamCollaboration`; a `404` on that read sends no `If-Match`, PC-CLI-D2) and `--wait` (CLI-12) polls `GetTeamCollaboration` until `status.phase` is `Ready` at or above the PUT's `generation`. Every `details.fields[].reason` and every `details.reason` of §4.4b prints with a one-line hint from a small table of hints (PC-CLI-D6); the hints add words, never decisions | COL-20, COL-34 | P0 |
| CLI-21 | `agents create-agent` and `update-agent`: `--role worker\|collaboration` (create only), `--collaboration-connect \| --no-collaboration-connect`, `--paperclip-agent-id`; `--role collaboration` with any of the other three exits 252; `list-agents --role`; `list-agent-templates --capability`; `admin put-agent-template --capabilities`. Agent tables gain `Role` and `Paperclip` columns; templates a `Capabilities` column | COL-1, COL-2, COL-10, AGT-11 | P0 |
| CLI-22 | `secrets list-secrets --managed-by user\|afs\|paperclip`; a `ManagedBy` column; the API's `409 managed_secret`, `master_key_not_rotatable` and `reserved_name` errors print with hints. The CLI never reads `managedBy` to decide anything and never prints a secret value (CLI-9 unchanged) | SEC-7, COL-7 | P0 |
| CLI-23 | `admin list-collaboration-endpoints --team-id --state`, table columns, `--output csv` (the operation is *CSV*, CLI-6) | COL-23, COL-36 | P0 (C2) |
| CLI-24 | Local command `afs collaboration invoke-api` (CLI-US-25): one request to the team's edge with the profile's AFS credential in `X-Afs-Authorization` and the Paperclip key (file or stdin, never an argument) in `Authorization`; origin from `GetTeamCollaboration.restApiUrl` or `--collab-url`, validated; no redirects, no cookies; edge errors map to AWS-style messages and exit 254; Paperclip's own errors exit 254 with the body on stdout; workload tokens refused locally | COL-19, COL-27 | P1 |
| CLI-25 | Redaction for the edge: `X-Afs-Authorization` joins the masked headers everywhere; `invoke-api --debug` logs no body and no query string; the AFS token, API key and Paperclip key never reach stdout, stderr, logs, the process arguments or shell history. `SecretDetected` errors print pattern types only | SEC-1, COL-26 | P0 |
| CLI-26 | The generated surface that the contract changes without new operations: `--may-expose-collaboration \| --no-may-expose-collaboration` on `teams create-team` and `update-team`; `--collab-client-id` and `--collab-client-secret-arn` (an ARN, validated as one, never a value) on `admin create-identity-provider` and `update-identity-provider`; role `team_collaborator` on `members` and `api-keys`; report kind `collaboration` with `--group-by surface\|reason\|person`; the seven `collab.*` settings through `admin put-setting` (no client-side keys or ranges) | ACC-2, COL-17, COL-23, COL-36 | P0 |
| CLI-27 | Contract compatibility: the pinned snapshot is re-pinned to the `-restapi` release that carries contract v1.3.0, then types, descriptors, the command table, reference docs and goldens are regenerated (never hand-edited); the reader and wait tables and the milestone table are extended and checked against the snapshot (ADR-008); `afs --version` prints `contract/1.3`. The profile gains one optional key, `collab_base_domain` (PC-CLI-D5) | API contract §8, §9 | P0 |

Product requirements `COL-*` and `SEC-7` are in the change PRD section 6; this repository owns only the client side of them.

### 2.1 Non-goals

- For the Paperclip change: no command that opens a browser or signs in to the web UI; no WebSocket or
  server-sent-event client in `invoke-api`; no commands for Paperclip's own objects (issues, agents, keys: `invoke-api`
  is the escape hatch, not a wrapper); no creation, listing or revocation of Paperclip keys by the CLI (they are secrets
  managed by AFS or by the team-admin, SEC-7, COL-12); no client-side secret scanning (the edge scans, COL-24); no
  client-side copy of the reserved secret names, the collaboration settings' ranges or the permission matrix; no
  `collaboration create-collaboration` (the collaboration agent is created by `create-agent --role collaboration`);
  no automatic connection of existing agents (use `connect-agent`, COL-10).

## 3. User stories and acceptance criteria

Added by the change PRD [paperclip-team-collaboration-PRD.md](https://github.com/stainedhead/enterprise-agent-factory-service/blob/main/paperclip-team-collaboration-PRD.md)
§7.6 (CLI-US-24, CLI-US-25), split finer here (CLI-US-26 to CLI-US-34). Names, flags and shapes follow API
contract v1.3.0 §4.4b and §4.4c; the contract wins on any difference. Each criterion is numbered
`CLI-US-n.k` so the traceability test can find the test that names it, and carries its `Proven by` test.
Test names are Go test functions in the package named; fakes are `httptest` servers through `clitest`
(the pinned snapshot's own examples are used by the generic tests, CLI-US-34).

**CLI-US-24 · See and change a team's collaboration settings.** *(refines US-27, US-30, US-33; COL-20, COL-34; CLI-19, CLI-20)*
As a team-admin, I want to read the state of my team's Paperclip endpoint and switch its web UI and REST
API from a script, so that I publish only what we use.
- 24.1 Given `afs collaboration describe-collaboration --team-id team_…`, then one `GetTeamCollaboration`
  (`GET /v1/teams/{teamId}/collaboration`) is sent and `--output json` prints the response unchanged. With
  `--output table` the columns are `Hostname`, `WebUi` (`webUiEnabled`), `RestApi` (`restApiEnabled`),
  `EffectiveWebUi` (`effective.webUi`), `EffectiveRestApi`, `Mode` (`registrationMode`), `Phase`
  (`status.phase`) and `Version`; a `null` hostname prints `-`. *Proven by*
  `collaboration.TestDescribeCollaborationJSONUnchanged`, `collaboration.TestDescribeCollaborationTableGolden`
  (`testdata/describe-collaboration.table.golden`, `…text.golden`).
- 24.2 Given a response whose `effective.reasons` is non-empty, then with `--output table` or `text` one
  stderr line per surface says `Not effective: <web UI|REST API> (<reason>[, <reason>])` using the contract's
  reason values (`team_not_allowed`, `service_endpoints_off`, `base_domain_unset`,
  `provider_without_collab_client`), and the exit code stays 0; with `--output json` nothing extra is
  printed. *Proven by* `collaboration.TestDescribeCollaborationNotEffectiveNotice`.
- 24.3 Given `put-collaboration --team-id team_… --web-ui-enabled --no-rest-api-enabled --registration-mode
  managed`, then one `PUT /v1/teams/{teamId}/collaboration` carries `{"webUiEnabled":true,
  "restApiEnabled":false,"registrationMode":"managed"}`; with `--registration-mode` omitted the member is
  omitted (the API keeps the current value). `--registration-mode` accepts `manual` or `managed`, anything
  else exits 252 before any call. *Proven by* `collaboration.TestPutCollaborationBody`,
  `collaboration.TestPutCollaborationRegistrationModeEnum252`.
- 24.4 Given `put-collaboration` without `--web-ui-enabled | --no-web-ui-enabled` or without
  `--rest-api-enabled | --no-rest-api-enabled` (the PUT replaces both), then it exits 252 naming the missing
  option, before any call; `--cli-input-json file://collab.json` supplies both. *Proven by*
  `collaboration.TestPutCollaborationBothSurfacesRequired252`, `collaboration.TestPutCollaborationInputJSON`.
- 24.5 Safe updates (CLI-10). With `--if-match 4` the request carries `If-Match: "4"` and no read is made;
  without it the CLI reads `GetTeamCollaboration`, sends its `version`, and prints
  `Applied to version N; use --if-match to guard against concurrent changes.` to stderr. A `412` prints
  `The resource changed (version 5). Re-run with --if-match 5 to apply your change to the latest version.`
  and exits 254. If the version read answers `404`, no `If-Match` is sent and the API's own answer
  (`409 Conflict`, `details.reason = collaboration_agent_missing`) is what the user sees (PC-CLI-D2).
  *Proven by* `collaboration.TestPutCollaborationIfMatchGuarded`, `collaboration.TestPutCollaborationReadsVersion`,
  `collaboration.TestPutCollaboration412Message`, `collaboration.TestPutCollaborationMissingAgent409`.
- 24.6 `--wait` polls `GetTeamCollaboration` every 5 s (or `Retry-After`) until `status.phase` is `Ready` and
  `status.observedGeneration` is at least the `generation` of the PUT response (exit 0, prints the final
  resource), or `Failed` or `Degraded` (exit 254, conditions including `EndpointReady` and `PaperclipReady`
  on stderr), or `--wait-timeout` passes (default 10 min, exit 255). A stale `Ready` with a lower
  `observedGeneration` does not end the wait. *Proven by* `collaboration.TestPutCollaborationWaitReady`,
  `collaboration.TestPutCollaborationWaitIgnoresStaleReady`, `collaboration.TestPutCollaborationWaitDegraded254`,
  `collaboration.TestPutCollaborationWaitTimeout255`, `wait.TestCollaborationWaitTableMatchesSnapshot`.
- 24.7 Validation answers from the API print AWS-style and exit 254, with every `details.fields[]` entry as
  one line `field: <field> reason: <reason>` on stderr after `Request ID`: turning a surface on with
  `collab_base_domain_unset`, `collab_public_endpoints_off`, `collaboration_not_allowed_for_team`,
  `provider_without_collab_client`, and `registrationMode managed` with `collaboration_agent_not_ready`.
  Each reason also prints the hint of PC-CLI-D6 (for example `Ask a service admin to set collab.base_domain
  (afs admin put-setting).`). Switching a surface off is never refused client-side. *Proven by*
  `collaboration.TestPutCollaborationValidationFieldsPrinted`, `exec.TestCollaborationReasonHints`.
- 24.8 `AccessDenied`: a viewer or collaborator running `put-collaboration` gets the API's `403`, printed as
  `An error occurred (AccessDenied) when calling the PutTeamCollaboration operation: …` with
  `Required permission: collaboration:manage` and exit 254; a viewer running `describe-collaboration` is
  allowed, and a caller with no binding gets the API's `404` (exit 254). A suspended team's `409` with
  `details.reason = team_suspended` prints its hint and exits 254. *Proven by*
  `collaboration.TestPutCollaborationAccessDenied254`, `collaboration.TestDescribeCollaborationAccessDenied254`,
  `collaboration.TestPutCollaborationSuspended409`.
- 24.9 `describe-collaboration` on a team with no collaboration agent prints
  `An error occurred (NotFound) when calling the GetTeamCollaboration operation: …`, the hint `This team has
  no collaboration agent. Create one with: afs agents create-agent --role collaboration …`, and exits 254.
  *Proven by* `collaboration.TestDescribeCollaborationMissingAgentHint`.
- 24.10 Every operation of this story has help with at least one example and a golden help page
  (`afs collaboration help`, `afs collaboration put-collaboration help`); the `put-collaboration` example
  is the equivalent shown on the console's T16 switch (`afs collaboration put-collaboration --team-id … --web-ui-enabled
  --no-rest-api-enabled`). *Proven by* `collaboration.TestCommandsDoneCollaboration`,
  `help.TestCollaborationHelpGolden`.

**CLI-US-25 · Call the Paperclip API from a pipeline.** *(refines US-31; COL-19, COL-27; CLI-24, CLI-25)*
As an automation author, I want one command that calls my team's Paperclip REST API with my AFS credential
and a Paperclip key, so that a pipeline can create and update work without a hand-built `curl`.

`afs collaboration invoke-api --team-id <id> --method <verb> --path <path> [--body file://… | fileb://…]
--paperclip-key-file <file | ->` is a **local command** (API contract §4.4b): it is not an operation and has
no `operationId` of its own; it sends one request to the team's edge (§4.4c) and prints the response.
Options: `--team-id` (or the profile's `team_id`), `--method` (`GET`, `HEAD`, `POST`, `PUT`, `PATCH`,
`DELETE`, case-insensitive; default `GET`), `--path` (required), `--body`, `--content-type` (default
`application/json` when a body is present), `--paperclip-key-file` (required), `--collab-url` (optional
override of the host, PC-CLI-D1), `--include` (print the status line and response headers before the body).
- 25.1 Given the REST API is on, a team API key with the collaborator role in `AFS_API_KEY`, and
  `--paperclip-key-file ./pc.key`, when I run `afs collaboration invoke-api --team-id team_… --method POST
  --path /api/issues --body file://issue.json --paperclip-key-file ./pc.key`, then exactly one
  `GetTeamCollaboration` and one request to the origin of the resource's `restApiUrl` are made; the second
  carries `X-Afs-Authorization: Bearer <the AFS credential>`, `Authorization: Bearer <contents of pc.key>`,
  the body, `Content-Type: application/json`, `Accept: application/json`, `User-Agent: afs-cli/<version>
  (<os>/<arch>)` and **no** `Cookie`; the response body is written to stdout unchanged and the exit is 0.
  *Proven by* `collaboration.TestInvokeAPIRequestShape`, `collaboration.TestInvokeAPIResponseBodyUnchanged`.
- 25.2 Given an SSO profile, then the SSO access token is sent in `X-Afs-Authorization` (silent refresh as in
  CLI-4); given `AFS_API_KEY`, the key; the credential chain of CLI-5 is used and `--debug` names the source
  without printing it. A workload token (`AFS_API_TOKEN_FILE`) is refused with exit 255 and `The collaboration
  endpoint does not admit workload credentials; use an API key or afs sso login.` before any call to the
  edge (PC-CLI-D3). A token known to be expired locally exits 255 with the login hint and makes no call.
  *Proven by* `collaboration.TestInvokeAPICredentialChain`, `collaboration.TestInvokeAPIWorkloadTokenRefused255`,
  `collaboration.TestInvokeAPIExpiredToken255`.
- 25.3 Given `--paperclip-key-file -`, then the key is read from stdin; a file whose mode is wider than 0600
  is refused with exit 255 and `chmod 600 <file>` (as CLI-5); trailing whitespace is trimmed; an empty
  key, a key longer than 4096 bytes or one containing whitespace or control characters exits 252 before any
  call. There is **no** option that takes the key as an argument: `--paperclip-key <value>` is an unknown
  option (exit 252) and the usage names `--paperclip-key-file`. `--body -` exits 252 (stdin is for the key).
  *Proven by* `collaboration.TestInvokeAPIKeyFromStdin`, `collaboration.TestInvokeAPIKeyFileMode255`,
  `collaboration.TestInvokeAPIKeyShapeRefused252`, `collaboration.TestInvokeAPIKeyAsArgumentUnknown252`.
- 25.4 Client-side validation, all exit 252 before any call: `--method` not in the list above; `--path` not
  starting with `/`, containing a scheme or host, `..` segments, control characters or a space; `--path`
  beginning `/_afs/` (the edge's reserved paths are never forwarded and are not Paperclip's); `--body` with
  `GET` or `HEAD`; a `--body` file over 100 MB (the service setting's maximum, `collab.max_body_mb`); a
  `--body` not `file://` or `fileb://`. The edge's own limit (`collab.max_body_mb`, default 25) is the API's
  `400 ValidationError` with `details.reason = body_too_large` and exits 254. *Proven by*
  `collaboration.TestInvokeAPIValidation252`, `collaboration.TestInvokeAPIReservedPathRefused252`,
  `collaboration.TestInvokeAPIBodyTooLargeFromEdge254`.
- 25.5 Where the request goes. The origin comes from `restApiUrl` of the `GetTeamCollaboration` answer, or from
  `--collab-url` (then no `GetTeamCollaboration` is made). `restApiUrl` is `https://<hostname>/api`, the base of Paperclip's REST API, so a full `--path /api/…` equals it plus a suffix (PC-D31). It must be `https://` (an `http://` loopback host is
  accepted for tests, as for `endpoint_url`), carry no userinfo, query, fragment or port other than 443, and a
  DNS name (not an IP literal; the apex `<collab.base_domain>` serves only the OIDC callback and is not a team host); if the profile sets `collab_base_domain` the host must be `<label>.<that
  domain>` (PC-CLI-D5). A `null` `restApiUrl` (service base domain unset) exits 255 with `Public
  collaboration is not configured for this service (collab.base_domain is unset).` Only the URL's scheme and
  host are used; `--path` is appended to the origin and the URL's own path (`/api`) is not (PC-CLI-D4).
  Redirects are **not** followed: a `3xx` is printed like any response (with `--include`, its `Location`). No
  cookie jar: a `Set-Cookie` in a response is never stored or sent. *Proven by*
  `collaboration.TestInvokeAPIOriginFromRestApiURL`, `collaboration.TestInvokeAPICollabURLSkipsRead`,
  `collaboration.TestInvokeAPIOriginRefused255` (table: http, userinfo, port, IP literal, foreign domain),
  `collaboration.TestInvokeAPINullRestApiURL255`, `collaboration.TestInvokeAPINoRedirectFollowed`,
  `collaboration.TestInvokeAPINoCookieJar`.
- 25.6 Errors from the edge (the §6 envelope, API contract §4.4c table) print AWS-style with the operation name
  `InvokeApi` and exit 254, then `Request ID: <requestId>`:
  `403 AccessDenied` `surface_disabled` → `An error occurred (AccessDenied) when calling the InvokeApi
  operation: …` plus `The REST API is switched off for this team. A team-admin can turn it on: afs
  collaboration put-collaboration --team-id <id> --rest-api-enabled …`; `403 permission_required` prints
  `Required permission: <details.requiredPermission>`; `403 no_team_role`; `401 Unauthenticated` with
  `token_expired`, `token_invalid`, `key_expired`, `key_revoked` prints the login or key hint (D41);
  `400` with `ambiguous_credentials`, `malformed_credential`, `unsupported_encoding` (`host_mismatch` no longer exists, PC-D23);
  `404 NotFound` (no detail, ever; only for an unknown host or both surfaces not effective, PC-D31) prints `The team's Paperclip is not reachable through the public endpoint; run afs collaboration describe-collaboration to see why.`; a caller with a binding gets `403 surface_disabled` instead when exactly one surface is off;
  `422 SecretDetected` prints the message and `Patterns: <details.patterns joined>` (pattern types only) and
  the hint `Store the credential as an AFS secret instead (afs secrets create-secret).`; `429 Throttled`
  and `503 Unavailable` (`paused`, `team_suspended`, `starting`, `upstream_unavailable`) print the reason.
  *Proven by* `collaboration.TestInvokeAPIEdgeErrors254` (table over every status and reason above, each
  with golden stderr), `collaboration.TestInvokeAPISurfaceDisabledHint`, `collaboration.TestInvokeAPISecretDetectedPatternsOnly`.
- 25.7 A response that is not the §6 envelope (Paperclip's own `4xx` or `5xx`, or a WAF block's default body)
  exits 254 with `An error occurred (HttpError) when calling the InvokeApi operation: HTTP <status>
  <reason phrase>` on stderr; the response body is still written to stdout so a pipeline can read
  Paperclip's error; the body is never copied into the stderr message (it is Paperclip content, ADR-012).
  A `2xx` exits 0 (`204` and `HEAD` print nothing); a `3xx` exits 0. *Proven by*
  `collaboration.TestInvokeAPIPaperclipErrorBodyOnStdout254`, `collaboration.TestInvokeAPIWAFBlock254`,
  `collaboration.TestInvokeAPI204AndHead`.
- 25.8 Retries (CLI-11, narrowed): `GET` and `HEAD` retry `429` and `503` (honouring `Retry-After`) and connect
  failures, at most 3 attempts; any other method retries only connect failures, never after the request was
  sent; no `Idempotency-Key` is sent (the edge does not define one). *Proven by*
  `collaboration.TestInvokeAPIRetriesGetOn503`, `collaboration.TestInvokeAPIPostNotRetriedAfterSend`,
  `collaboration.TestInvokeAPINoIdempotencyKey`.
- 25.9 Global options. `--profile`, `--endpoint-url` (the AFS API used for the one read), `--debug`,
  `--cli-connect-timeout`, `--cli-read-timeout`, `--color`, `--no-cli-pager` and `--version` behave as usual.
  `--query`, `--output`, `--no-paginate`, `--max-items`, `--page-size`, `--starting-token`,
  `--cli-input-json`, `--generate-cli-skeleton` and `--wait` exit 252 with `invoke-api prints the response
  unchanged; pipe it to a JSON tool instead.` (PC-CLI-D7). Streaming responses (server-sent events) and
  WebSocket are not supported: the read timeout applies and a `101` is an error (exit 255). *Proven by*
  `collaboration.TestInvokeAPIGlobalOptionsRefused252`, `collaboration.TestInvokeAPIReadTimeout255`.
- 25.10 Redaction (CLI-25). With `--debug`, stderr shows one line per event with method, host, path **without
  its query string**, status, duration, `X-Request-Id`, credential source and attempt; `Authorization` and
  `X-Afs-Authorization` print `[redacted]`; request and response **bodies are never logged**; `Set-Cookie`
  and `Cookie` print `[redacted]`. A test drives a run with a recognisable AFS token, API key, Paperclip key and
  a `?token=…` query and asserts that none of those strings, nor the body, appears anywhere in stdout (other
  than the response body itself), stderr, `--debug` output or the process arguments. *Proven by*
  `collaboration.TestInvokeAPIDebugRedactsEverything`, `collaboration.TestInvokeAPISecretsNeverInArgsOrStderr`,
  `redact.TestRedactsEdgeHeaders`.
- 25.11 The command is listed by `afs collaboration help` and in the generated command reference with a
  golden help page and two examples (one `GET`, one `POST` with `--body`); it is excluded from CLI-US-7's
  operation comparison by name (as `configure`, `sso login`). *Proven by* `help.TestInvokeAPIHelpGolden`,
  `services.TestOperationTableExemptsLocalCommands`.

**CLI-US-26 · Connect and disconnect agents, and see the links.** *(refines US-28; COL-10, COL-12, COL-14; CLI-19)*
As a team-admin, I want to connect existing agents to Paperclip in bulk from a script and see their link
state, so that agents that predate the collaboration agent can pull work.
- 26.1 `afs collaboration connect-agent --team-id team_… --agent-id agt_… [--paperclip-agent-id ag_…]` sends
  one `POST …/collaboration/links/{agentId}:connect` with body `{"paperclipAgentId":"ag_…"}` (or `{}`), an
  `Idempotency-Key`, and prints the `CollaborationLink`; `--output table` columns `AgentId`, `Name`
  (`agentName`), `Connect`, `PaperclipAgentId`, `LinkState`, `LastSyncedAt`. `disconnect-agent` sends
  `…:disconnect` with no body and the same columns. *Proven by*
  `collaboration.TestConnectAgentRequest`, `collaboration.TestDisconnectAgentRequest`,
  `collaboration.TestLinkColumnsGolden`.
- 26.2 A repeat with the same `--idempotency-key` replays the stored response (`Idempotent-Replayed`) and
  prints it; the CLI-11 retry rules apply. Neither operation takes `If-Match`, has `--wait` or reads a version
  (`--wait` is an unknown option, exit 252). *Proven by* `collaboration.TestConnectAgentReplay`,
  `collaboration.TestConnectAgentNoWaitNoIfMatch`.
- 26.3 API refusals print AWS-style and exit 254: `ValidationError` for the collaboration agent itself
  (`field: agentId`), for a pool, and for `--paperclip-agent-id` while `registrationMode` is `managed`
  (`field: paperclipAgentId`, `reason: managed_registration`, with the hint `Registration is managed; the
  reconciler records the id.`); `NotFound` for an unknown agent; `409` `collaboration_agent_missing` or
  `team_suspended`. *Proven by* `collaboration.TestConnectAgentValidation254`,
  `collaboration.TestConnectAgentManagedRegistrationHint`.
- 26.4 `afs collaboration list-collaboration-links --team-id team_… [--link-state pending|linked|failed|disabled]
  [--connect | --no-connect]` sends `GET …/collaboration/links?linkState=…&connect=true|false`, paginates by
  default (CLI-7), and `--connect` omitted sends no `connect` query member. `--link-state` outside the four
  values exits 252. `lastError` is printed as returned (it never holds a key or content, contract §4.4b);
  with `--output table` the `LastError` column is truncated with the table's mark. *Proven by*
  `collaboration.TestListCollaborationLinksFilters`, `collaboration.TestListCollaborationLinksPagination`,
  `collaboration.TestListCollaborationLinksLinkStateEnum252`, `collaboration.TestLinksTableGolden`.
- 26.5 `AccessDenied` (exit 254, `Required permission: collaboration:manage`) for `connect-agent` and
  `disconnect-agent` by a collaborator or viewer; `collaboration:read` for `list-collaboration-links`.
  *Proven by* `collaboration.TestConnectAgentAccessDenied254`, `collaboration.TestDisconnectAgentAccessDenied254`,
  `collaboration.TestListCollaborationLinksAccessDenied254`.

**CLI-US-27 · Rotate Paperclip's keys.** *(refines US-27; COL-7; CLI-19)* As a team-admin, I want to rotate
Paperclip's signing secrets, knowing it signs everyone out.
- 27.1 `afs collaboration rotate-collaboration-keys --team-id team_… [--names paperclip-better-auth-secret
  paperclip-tool-action-signing-secret]` sends one `POST /v1/teams/{teamId}/collaboration:rotate-keys` with an
  `Idempotency-Key` and body `{"names":[…]}` (no body when `--names` is omitted; the API rotates both). It
  prints the response; `--output table` lists `Name`, `SecretId`, `VersionId` from `rotated`. *Proven by*
  `collaboration.TestRotateCollaborationKeysBody`, `collaboration.TestRotateCollaborationKeysTable`.
- 27.2 After a success the CLI prints to stderr `Everyone is signed out of Paperclip.` when
  `signsOutEveryone` is true and `Agents will roll to pick up the new keys.` when `agentRolloutPending` is
  true. It never prompts for confirmation (PC-CLI-D8). No secret value exists in the response or any output.
  *Proven by* `collaboration.TestRotateCollaborationKeysNotices`, `collaboration.TestRotateCollaborationKeysNoValue`.
- 27.3 Naming `paperclip-secrets-master-key` in `--names` is **not** refused by the CLI (the API decides); the
  API's `400 ValidationError` (`field: names`, `reason: master_key_not_rotatable`) prints with the hint `The
  master key cannot be rotated: Paperclip's stored identities would become undecryptable.` and exits 254.
  `afs secrets rotate-secret` on the master key's secret id gets the same error (CLI-US-29). *Proven by*
  `collaboration.TestRotateCollaborationKeysMasterKey254`.
- 27.4 `AccessDenied` for collaborator, viewer and security operator (`Required permission:
  collaboration:manage`, exit 254). *Proven by* `collaboration.TestRotateCollaborationKeysAccessDenied254`.

**CLI-US-28 · Create the collaboration agent and connect workers.** *(refines US-27, US-28; COL-1, COL-2, COL-10; CLI-21)*
As a team-admin, I want `create-agent` to create the team's Paperclip and to connect workers, so that one
`--cli-input-json` file stands up a team.
- 28.1 `afs agents create-agent --team-id team_… --name paperclip --role collaboration --template paperclip
  --image …@sha256:… --model-id … [--volume …] --wait` sends `"role":"collaboration"` and returns 0 when the
  agent's phase is `Ready` (with the default 10 min `--wait-timeout`, since the product target is 5 min at
  p95). `--role` accepts `worker` (the default; omitted when not given) or `collaboration`, anything else exits
  252. *Proven by* `agents.TestCreateAgentRoleCollaborationBody`, `agents.TestCreateAgentRoleEnum252`,
  `agents.TestCreateAgentCollaborationWaitReady`.
- 28.2 `--collaboration-connect` and `--no-collaboration-connect` send `{"collaboration":{"connect":true|false}}`;
  `--paperclip-agent-id ag_…` adds `paperclipAgentId`. All three are valid on `create-agent` and `update-agent`
  for workers (`update-agent` is the same effect as `connect-agent`, contract §4.4). With none of them the body
  has no `collaboration` member, so the API default applies (`true` when the team has a collaboration agent). Given
  `--role collaboration` together with any of the three options, the CLI exits 252 before any call (the API
  would refuse `collaboration` on that role). *Proven by* `agents.TestCreateAgentCollaborationConnectBody`,
  `agents.TestUpdateAgentCollaborationConnectBody`, `agents.TestCreateAgentCollaborationRoleWithConnect252`,
  `agents.TestCreateAgentNoCollaborationMemberByDefault`.
- 28.3 `update-agent` has no `--role` option: it exits 252 as an unknown option (`--replicas` precedent,
  CLI-US-12); a change of role through `--cli-input-json` is the API's `ValidationError` (exit 254). *Proven by*
  `agents.TestUpdateAgentRoleUnknown252`, `agents.TestUpdateAgentRoleChangeRefused254`.
- 28.4 `list-agents --role worker|collaboration` sends `?role=`; any other value exits 252. The agent table gains
  `Role` and `Paperclip` columns after `Phase` (`role`; `collaboration.linkState`, `-` on the collaboration agent
  and on agents with no `collaboration` object); the table stays within 80 columns through the existing
  truncation (`output.TestTableFitsIn80Columns` stays green). *Proven by* `agents.TestListAgentsRoleFilter`, `agents.TestListAgentsRoleEnum252`,
  `agents.TestAgentTableRoleAndPaperclipColumnsGolden`.
- 28.5 `describe-agent --output json` prints `role` and `collaboration` `{connect, paperclipAgentId, linkState,
  lastSyncedAt, lastError}` unchanged. *Proven by* `agents.TestDescribeAgentCollaborationJSONUnchanged`.
- 28.6 The API's refusals print AWS-style and exit 254: `409 Conflict` with `details.reason =
  collaboration_agent_exists` (hint `This team already has a collaboration agent; see afs agents list-agents --role collaboration.`), `ValidationError` naming the template (`field: template`), `iamIdentity team_shared`,
  `volume.snapshotSchedule off`, a `paperclip` template on a worker, `collaboration.connect true` with no
  collaboration agent (`field: collaboration.connect`), and `503 Unavailable` when Secrets Manager cannot be
  written (retried per CLI-11 because `create-agent` carries an `Idempotency-Key`). *Proven by*
  `agents.TestCreateAgentCollaborationErrors254` (table), `agents.TestCreateAgentCollaboration503Retried`.
- 28.7 `afs agents list-agent-templates [--capability paperclip]` and `afs admin list-agent-templates
  [--capability paperclip]` send `?capability=`; templates show `capabilities` in `--output json` and a
  `Capabilities` column in tables; `afs admin put-agent-template --capabilities paperclip` sends
  `capabilities: ["paperclip"]` (space-separated list, CLI-2 style; `--capabilities` with no value sends `[]`).
  *Proven by* `agents.TestListAgentTemplatesCapabilityFilter`, `admin.TestPutAgentTemplateCapabilities`,
  `admin.TestAgentTemplatesCapabilitiesColumnGolden`.
- 28.8 `AccessDenied` (exit 254, `Required permission: agents:write`) for a viewer or collaborator creating a
  collaboration agent. *Proven by* `agents.TestCreateAgentCollaborationAccessDenied254`.
- 28.9 Skeleton: `afs agents create-agent --generate-cli-skeleton` includes `role` and `collaboration` with their
  contract members and makes no call. *Proven by* `agents.TestCreateAgentSkeletonHasRoleAndCollaboration`.

**CLI-US-29 · See who manages each secret.** *(refines US-27, US-32; SEC-7, COL-7; CLI-22)* As a team-admin, I want
to tell my own secrets from the ones AFS generated or Paperclip created, so that I do not try to change what I
cannot.
- 29.1 `afs secrets list-secrets --team-id team_… [--managed-by user|afs|paperclip]` sends `?managedBy=`; any other
  value exits 252. `--output json` prints `managedBy` unchanged; the table gains a `ManagedBy` column. *Proven by*
  `secrets.TestListSecretsManagedByFilter`, `secrets.TestListSecretsManagedByEnum252`,
  `secrets.TestSecretTableManagedByColumnGolden`.
- 29.2 `update-secret`, `delete-secret` and `put-secret-value` on an `afs` or `paperclip` secret print the API's
  `409 Conflict` (`details.reason = managed_secret`) with the hint `This secret is managed by <AFS|Paperclip>
  and cannot be changed here.` and exit 254. The CLI does not read `managedBy` first (the API decides), and a
  `put-secret-value` still reads the value only from a file, stdin or the prompt (CLI-9) and never echoes it
  even when the call is refused. *Proven by* `secrets.TestManagedSecretMutations409`,
  `secrets.TestManagedSecretRefusalNeverEchoesValue`.
- 29.3 `rotate-secret` on an `afs` secret sends no body (AFS generates the value, PC-D30) and prints the rotation; on
  `paperclip-secrets-master-key` it is the API's `ValidationError` (`reason: master_key_not_rotatable`), and on a
  `paperclip`-managed secret the API's `409 Conflict` (`managed_secret`); both exit 254. *Proven by* `secrets.TestRotateManagedSecret`, `secrets.TestRotateMasterKey254`.
- 29.4 `create-secret --name paperclip` or any `paperclip-…` name other than `paperclip-agent-api-key` on an
  agent-scope secret prints the API's `ValidationError` (`reason: reserved_name`) and exits 254 (the CLI keeps
  no copy of the reserved list); `create-secret --name paperclip-agent-api-key --scope agent --agent-id agt_…
  --secret-value file://…` sends the request normally (manual registration, COL-12). Binding a `paperclip`
  secret is the API's `ValidationError`. *Proven by* `secrets.TestCreateSecretReservedName254`,
  `secrets.TestCreateAgentKeySecretManualRegistration`.
- 29.5 `AccessDenied` (exit 254) for a collaborator running `create-secret` (`secrets:write`). *Proven by*
  `secrets.TestCreateSecretAccessDeniedCollaborator254`.

**CLI-US-30 · Service admin: see every team's endpoint.** *(refines US-33; COL-23, COL-36; CLI-23)* As a service
admin, I want a list of every team's Paperclip endpoint and its state from a script.
- 30.1 `afs admin list-collaboration-endpoints [--team-id team_…] [--state Ready|Degraded|Failed|Pending]` sends
  `GET /v1/admin/collaboration-endpoints?teamId=&state=` and paginates by default; `--state` outside the four values
  exits 252 (`--team-id` is the query filter here, not a path parameter, and falls back to the profile's
  `team_id` only for path parameters, so it is not defaulted). *Proven by*
  `admin.TestListCollaborationEndpointsFilters`, `admin.TestListCollaborationEndpointsStateEnum252`,
  `admin.TestListCollaborationEndpointsTeamIdNotDefaulted`.
- 30.2 `--output table` columns: `Team` (`slug`), `TeamState`, `Hostname`, `WebUi` (`effective.webUi`), `RestApi`
  (`effective.restApi`), `Phase`, `LastAccessAt`; a `null` `lastAccessAt` prints `-`. `--output csv` is valid
  (the operation is *CSV* in the contract): it sends `Accept: text/csv`, follows `X-Afs-Next-Token` over 3 pages
  and writes one header row (CLI-6). *Proven by* `admin.TestCollaborationEndpointsTableGolden`,
  `admin.TestListCollaborationEndpointsCSVThreePages`.
- 30.3 `AccessDenied` (exit 254, `Required permission: admin:settings`) for a team-admin, a collaborator and a
  security operator. *Proven by* `admin.TestListCollaborationEndpointsAccessDenied254`.

**CLI-US-31 · Service admin: guardrail, provider client and settings.** *(refines US-33, US-29; COL-17, COL-23, ACC-1; CLI-26)*
As a service admin, I want to switch collaboration per team, register the collaboration sign-in client, and set the
`collab.*` settings from scripts.
- 31.1 `afs teams create-team` and `afs teams update-team` accept `--may-expose-collaboration` and
  `--no-may-expose-collaboration` (sent as `mayExposeCollaboration`; omitted when neither is given); a team-admin
  using them gets the API's `403 AccessDenied` with `Required permission: teams:set-guardrails` (exit 254).
  `describe-team --output json` prints `mayExposeCollaboration` unchanged. `create-team` with a slug whose hostname is
  held prints the API's `409 Conflict` (`details.reason = hostname_held`) with the hint `That slug's collaboration
  hostname is held for 30 days after the deletion of a team that had one; choose another slug.` (exit 254), proven by
  `teams.TestCreateTeamHostnameHeld409`. *Proven by* `teams.TestUpdateTeamMayExposeCollaborationBody`,
  `teams.TestCreateTeamMayExposeCollaborationBody`, `teams.TestUpdateTeamMayExposeCollaborationAccessDenied254`.
- 31.2 `afs admin create-identity-provider` and `update-identity-provider` accept `--collab-client-id <id>` and
  `--collab-client-secret-arn <arn>` (an ARN of a Secrets Manager secret, never the value; the option name says
  so and the help explains it). A value that is not an `arn:aws:secretsmanager:…` string exits 252 before any call.
  `describe-identity-provider` prints both members unchanged; `--debug` does **not** redact the ARN (it is a
  reference, not a secret) but still redacts `clientSecret` and any `…Secret` member that is not an ARN
  reference. *Proven by* `admin.TestIdentityProviderCollabClientBody`, `admin.TestIdentityProviderCollabSecretArnShape252`,
  `redact.TestCollabClientSecretArnNotRedactedButSecretIs`.
- 31.3 `afs admin put-setting --key collab.base_domain --value '"collab.example.com"'` and the other six `collab.*`
  settings send the JSON value as for every setting (CLI-US-15); the CLI carries no list of keys or ranges: an
  out-of-range value (`collab.session_ttl_minutes 5`, `collab.max_body_mb 500`, `collab.content_scan "maybe"`,
  `collab.allowed_source_cidrs ["10.0.0.0/33"]`) prints the API's `ValidationError` naming the range and exits 254.
  `admin get-settings --output json` prints the settings unchanged. *Proven by*
  `admin.TestPutSettingCollabKeysPassThrough`, `admin.TestPutSettingCollabOutOfRange254`.
- 31.4 `AccessDenied` for a team-admin running `put-setting` or `create-identity-provider` is already covered by
  CLI-US-15 and CLI-US-21; this story adds the cases for the new members. *Proven by*
  `admin.TestIdentityProviderCollabClientAccessDenied254`.

**CLI-US-32 · Collaborator role and the collaboration report.** *(refines US-29, US-31; ACC-2, COL-36; CLI-26)* As a
team-admin, I want to give a person or a pipeline the collaborator role, and to read who uses Paperclip.
- 32.1 `afs members add-member`, `members update-member` and `api-keys create-api-key` accept `--role
  team_collaborator` (the API enum is `team_admin`, `team_collaborator`, `team_viewer`; for `create-api-key` also
  `security_operator`, which the API still refuses on a team key); a service-level key or a service-level binding
  with `team_collaborator` is the API's `ValidationError` (exit 254). The help lists the three roles. *Proven by*
  `members.TestAddMemberTeamCollaborator`, `apikeys.TestCreateApiKeyTeamCollaborator`,
  `members.TestMemberRoleEnumHelpGolden`.
- 32.2 `afs reports get-team-report --kind collaboration --range 7d --group-by day,person` and `afs reports
  get-service-report --kind collaboration --group-by team` are accepted: the CLI's client-side `--group-by` check
  gains `surface`, `reason` and `person` (valid only with the collaboration kind, as the contract says;
  other kinds with them are the API's `ValidationError`, the CLI no longer rejects them client-side). `--output csv`
  works as for other reports; the table columns follow the response (`Dynamic`). `create-report-schedule --report
  collaboration` is the API's `ValidationError` (the report is not schedulable, PC-D13), exit 254. *Proven by*
  `reports.TestCollaborationReportGroupByAccepted`, `reports.TestGroupByPatternMessageListsNewDimensions`,
  `reports.TestCollaborationReportScheduleRefused254`, `reports.TestCollaborationReportCSV`.
- 32.3 `AccessDenied` (exit 254) for a security operator running `get-team-report --kind collaboration` (team
  scoped routes are refused to them). *Proven by* `reports.TestCollaborationReportAccessDenied254`.

**CLI-US-33 · Edge credentials are never shown.** *(refines US-31, US-32; SEC-1, COL-26; CLI-25)* As security, I want
the new credentials treated like every other.
- 33.1 `internal/redact` masks the headers `Authorization`, `Cookie`, `Set-Cookie`, `Proxy-Authorization` and the new
  `X-Afs-Authorization`, and `--debug` output of every command (not only `invoke-api`) uses it. The value of the
  edge cookie `__Host-afs-collab` never appears in any output. *Proven by* `redact.TestRedactsEdgeHeaders`,
  `client.TestDebugRedactsXAfsAuthorization`.
- 33.2 The query member `code` (the edge's one-time handoff code) and any `token`, `key` or `secret` member are
  dropped from URLs in `--debug` lines for `invoke-api`; for the AFS API's own calls the existing rule
  (CLI-US-3, `TestNoQueryParamCarriesSecret`) stays. *Proven by* `collaboration.TestInvokeAPIDebugDropsQuery`.
- 33.3 `rotate-collaboration-keys`, `connect-agent` and every secret path never print or log a secret value; a
  test greps the stdout, stderr and `--debug` output of all seven new operations for a planted value and finds it
  nowhere. *Proven by* `collaboration.TestNoSecretValueInAnyCollaborationOutput`.

**CLI-US-34 · The command table, contract pin and documentation follow the contract.** *(extends CLI-US-7; refines
US-27 to US-33; CLI-1, CLI-14, CLI-27)* As a maintainer, I want the new operations to exist because the pinned
contract has them, with help, examples, columns, goldens and docs.
- 34.1 After the re-pin the command-table test (CLI-US-7) passes with seven more operations and fails on any of them
  without a command or with a name that differs from the contract's `CLI` column:
  `collaboration describe-collaboration`, `put-collaboration`, `list-collaboration-links`, `connect-agent`,
  `disconnect-agent`, `rotate-collaboration-keys`, and `admin list-collaboration-endpoints`. `invoke-api` is
  listed as a local command (not an operation). *Proven by* `services.TestCommandTableMatchesSnapshot`,
  `services.TestOperationTableExemptsLocalCommands`.
- 34.2 Every new operation has at least one example, table columns (or `Dynamic`), a golden help page, a golden
  table and text rendering, and an `AccessDenied` test; `CurrentMilestone` is raised only by the exit task of the
  milestone (PC-CLI-D9: M7 to M10), and `Pending()` is empty at each. *Proven by*
  `services.TestDoneUpToCurrentMilestone`, `services.TestDoneOperationsHaveGoldens`,
  `help.TestOperationHelpShowsOptionsPermissionAndExample`.
- 34.3 The `If-Match` reader table and the `--wait` table are checked against the snapshot (ADR-008): `PutTeamCollaboration`
  has a reader (`GetTeamCollaboration`, same path) and a wait entry (the resource's own `status`); the
  two new POST actions have neither. A test fails on any divergence. *Proven by*
  `safeupdate.TestReaderTableMatchesSnapshot`, `wait.TestCollaborationWaitTableMatchesSnapshot`.
- 34.4 `api/openapi.yaml` and `api/SOURCE` are re-pinned with `make contract-sync` from the `-restapi` ref that carries
  contract v1.3.0 (PC-CLI-D10), then `make generate`, `make docs`; `afs --version` prints `contract/1.3`; the
  snapshot test (`api/snapshot_test.go`) passes on the new sha256. Between the re-pin and the last milestone the
  operations of later milestones exist as generated commands but are not "done" (`milestoneOf`). *Proven by*
  `api.TestSnapshotChecksumMatchesSource`, `buildinfo.TestVersionPrintsContractMinor`.
- 34.5 Documentation: `user-docs/collaboration.md` (new: connect to Paperclip, the public endpoint, `invoke-api`
  with a team key, secrets in Paperclip, troubleshooting by edge error); `user-docs/agents-and-secrets.md`,
  `teams-and-members.md` and `delivery-and-admin.md` updated; `user-docs/reference/` regenerated
  (`collaboration.md`, and the changed `agents.md`, `secrets.md`, `admin.md`, `members.md`, `api-keys.md`,
  `reports.md`, `teams.md`); none of it links into `specs/`. The root repository's
  `skills/enterprise-agent-factory-service-cli.md` gains the new commands in the same cycle (CLI-14 workflow
  step 6). *Proven by* `docs.TestDocgenMatchesHelp` (`make docs` leaves no diff).

**CLI-US-35 · Swarm templates may carry a collaboration agent.** *(refines US-33 and PC-D29; `SwarmTemplate.spec.collaboration`)*
The swarm template commands take the spec as a file (`--spec file://…`), so no new option is needed.
- 35.1 `afs swarm-templates create-swarm-template … --spec file://spec.json` with `"collaboration":{"templateId":"…"}` in the
  spec sends the spec unchanged, and `describe-swarm-template --output json` prints `spec.collaboration` unchanged;
  `--generate-cli-skeleton` shows the optional member. *Proven by* `swarmtemplates.TestCreateSwarmTemplateSpecCollaborationUnchanged`,
  `swarmtemplates.TestSwarmTemplateSkeletonHasCollaboration`.
- 35.2 A template without the `paperclip` capability is the API's `ValidationError` (`field: spec.collaboration.templateId`),
  printed AWS-style with exit 254; a security operator gets `AccessDenied` as for every swarm template write. *Proven by*
  `swarmtemplates.TestSwarmTemplateCollaborationTemplateRefused254`, `swarmtemplates.TestSwarmTemplateCollaborationAccessDenied254`.
- 35.3 Pools never connect to Paperclip: `pools` commands gain no collaboration option. Priority P1, M10 (C4). *Proven by*
  `pools.TestNoCollaborationOptionOnPoolCommands`.

## 4. Non-functional requirements

| Area | Requirement |
|---|---|
| Paperclip: performance | `invoke-api` adds at most one `GetTeamCollaboration` round trip (none with `--collab-url`) and streams the response body to stdout without buffering it; startup stays under 100 ms with the new service registered (`startup_bench_test.go`) |
| Paperclip: security | The AFS credential goes only to the AFS API and to the validated edge origin, over TLS 1.2 or later (the transport rules of `endpoint_url`, ADR-012); never to a redirect target; the Paperclip key only in the `Authorization` header of that one request; neither credential is in the process arguments, `--debug`, stdout (except Paperclip's own response body), stderr or any file the CLI writes; the Paperclip key file must be 0600 or stdin; a workload token is never sent to the edge. Data from the API (`restApiUrl`, `lastError`, error messages) is untrusted and passes through `output.Sanitize` (ADR-012) |
| Paperclip: reliability | Retries per CLI-US-25.8; timeouts per the existing flags; Ctrl-C cancels the in-flight request and exits 130; a `503` page or WAF body never panics the error mapper (non-JSON bodies are handled) |
| Paperclip: compatibility | A CLI built for contract 1.2 against a 1.3 server keeps working and prints the one-per-session minor-version notice (D47); `describe-agent` of such a CLI simply shows the new members. A CLI built for 1.3 against a 1.2 server gets `404` on the new paths, printed AWS-style with exit 254 |
| Paperclip: observability | `--debug` on `invoke-api` writes one line per event with component `collab` (origin read, request sent, response, retry, redaction applied) in the CLI-wide format, naming the credential source and never a credential, body or query string; every error prints `Request ID: <X-Request-Id>` from the edge's response; `User-Agent` is the CLI's. No telemetry is sent anywhere |
| Paperclip: quality | The new packages (`internal/cli/collaboration`, the changed parts of `agents`, `secrets`, `admin`, `redact`) hold at least 90% coverage; golden files for every new table, text and help page; every operation tested for `AccessDenied`; `internal/cli/services` traceability test extended to CLI-US-24 to CLI-US-34 |

### 4.1 Acceptance for the Paperclip cross-cutting requirements

- **Hints (PC-CLI-D6).** One table-driven test, `exec.TestCollaborationReasonHints`, walks every
  `details.reason` value of API contract §4.4b, §4.4c, §4.6 (`collaboration_agent_exists`,
  `collaboration_agent_missing`, `team_suspended`, `managed_secret`, `master_key_not_rotatable`, `reserved_name`,
  `managed_registration`, `surface_disabled`, `no_team_role`, `permission_required`, `ambiguous_credentials`,
  `malformed_credential`, `body_too_large`, `unsupported_encoding`, `token_invalid`,
  `token_expired`, `key_expired`, `key_revoked`, `session_expired`, `paused`, `starting`, `upstream_unavailable`) and
  `details.fields[].reason` values, and fails when a value has no hint or a hint contains text not in the table (a
  hint never includes server-supplied text). An unknown reason prints no hint and still exits 254.
- **Origin safety (CLI-US-25.5).** `collaboration.TestInvokeAPIOriginRefused255` is table-driven over: `http://host`,
  `https://user:pw@host`, `https://host:8443`, `https://10.0.0.1`, `https://[::1]`, `https://host/?q=1`,
  `https://host/#f`, a host outside `collab_base_domain` and an empty host; each exits 255 with no request to the
  edge and a message that names the rule, not the value of any credential.
- **No credential anywhere.** `collaboration.TestInvokeAPIDebugRedactsEverything` and
  `collaboration.TestNoSecretValueInAnyCollaborationOutput` plant distinct sentinel strings for the AFS token, API
  key, Paperclip key and a `?token=` query, run every success and failure path of CLI-US-24 to CLI-US-27 and CLI-US-25,
  and assert the sentinels are absent from stdout (except as part of a response body the fake server itself returned),
  stderr, `--debug` output and the recorded request line and arguments of the process.
- **Exit codes.** A table in `collaboration` tests pins exit 0, 130, 252, 254 and 255 for each path above: usage and
  client-side validation 252, any API or edge answer 254, local credential, file, origin and network problems 255.
- **Compatibility (CLI-27).** `api.TestSnapshotChecksumMatchesSource` and `make contract-check`, `make generate-check` pass on
  the re-pinned snapshot; a snapshot from before the re-pin (contract 1.2) still builds in a test fixture so the
  1.2 to 1.3 notice (D47) is proven; `services.TestCommandTableMatchesSnapshot` fails when a v1.3.0 operation is
  removed from the snapshot.

## 5. Dependencies

| Dependency | Use | Notes |
|---|---|---|
| `-restapi` release whose `api/openapi.yaml` carries contract v1.3.0 (the seven §4.4b operations, `role` and `collaboration` on agents, `managedBy`, `capabilities`, `collabClientId`, `mayExposeCollaboration`, `team_collaborator`, report kind `collaboration`) | CLI-19 to CLI-27 | Not available until `-restapi` merges its C1 work; this repository re-pins then (PC-CLI-D10). The pinned snapshot, `api/SOURCE`, `internal/client/gen`, `internal/spec/operations_gen.go`, `user-docs/reference/` and every golden file that lists operations must be resynced and regenerated at that point; **do not hand-edit them now** |
| `-restapi` `afs-edge` deployed in a dev environment, with `collab.base_domain` set, a team provider that has a collaboration client, and a team with a collaboration agent | End-to-end proof of `invoke-api` (CLI-US-25) in C3; unit tests use a fake edge | Change PRD C2 and C3 exit criteria |
| Contract v1.3.1 (PC-D20 to PC-D31), which answers PC-CLI-OQ1 to OQ4 | `invoke-api` origin, path and error rules | Root; merged (1f60ecd) |
| Root `skills/enterprise-agent-factory-service-cli.md` | New commands for agents | Updated in the same cycle as the release (workflow step 6) |
| Console's T16 shows the equivalent `afs collaboration put-collaboration` command (change PRD §5) | Parity of examples | `-website`; the example text in CLI-US-24.10 is the source |

## 6. Delivery

Paperclip change (change PRD section 10, milestones C1 to C4; CLI milestone numbers continue M6, PC-CLI-D9). Each
milestone is a pull request into `main` that raises `CurrentMilestone` only in its exit task, and each starts only after
the `-restapi` release that provides its operations (the re-pin is the first task, PC-CLI-D10).

| Milestone | Change PRD | This repository | Exit criteria |
|---|---|---|---|
| M7 | C1 Collaboration agent | Re-pin to contract v1.3.0 (CLI-27, CLI-US-34); `collaboration describe-collaboration`, `put-collaboration` (If-Match, `--wait`), `list-collaboration-links`, `connect-agent`, `disconnect-agent` (CLI-US-24, CLI-US-26: all P0 in C1, PC-D28); `agents` flags, columns, `--role`, `--capability`, `admin put-agent-template --capabilities` (CLI-US-28); `secrets --managed-by` and hints (CLI-US-29); `teams --may-expose-collaboration`, `members` and `api-keys` `team_collaborator` (CLI-US-31.1, CLI-US-32.1); redaction of `X-Afs-Authorization` (CLI-US-33.1); reference docs and `user-docs/collaboration.md` first cut | `make gate` green; every P0 operation of §4.4b has example, columns, goldens and `AccessDenied` test; the end-to-end script of the change PRD (create the collaboration agent, connect two agents) passes against dev with `afs` only; `Pending()` empty at milestone 7 |
| M8 | C2 Public web UI | `admin list-collaboration-endpoints` (CLI-US-30, P0 in C2, PC-D28); `admin create-identity-provider` and `update-identity-provider` `--collab-client-id`, `--collab-client-secret-arn`; `collab.*` settings through `put-setting` verified (CLI-US-31.2 to 31.4); help text for the effective-surface notice (CLI-US-24.2); docs for the public endpoint | Gate green; a service admin sets the base domain, the collaboration client and a team guardrail from the CLI in dev; `Pending()` empty at milestone 8 |
| M9 | C3 REST API and managed mode | `collaboration invoke-api` (CLI-US-25, CLI-US-33.2); `rotate-collaboration-keys` (CLI-US-27); managed-mode hints (`managed_registration`) | Gate green; a pipeline creates a Paperclip issue with a collaborator team key against the dev edge; all redaction tests green; `Pending()` empty at milestone 9 |
| M10 | C4 People and reports | swarm-template `spec.collaboration` goldens (CLI-US-35); report kind `collaboration` and the `--group-by` additions (CLI-US-32.2, 32.3); final docs, root `skills/` entry, traceability test over CLI-US-24 to CLI-US-34; release | Gate green; `CurrentMilestone` is 10; every acceptance criterion of CLI-US-24 to CLI-US-34 names a passing test |

## 7. Open questions

| # | Question | Owner | Status or proposal |
|---|---|---|---|
| PC-CLI-OQ1 | Contract said `invoke-api` calls no `/v1` operation but needs the hostname | root | **Resolved, PC-D31 (contract v1.3.1):** one `GetTeamCollaboration` read unless `--collab-url`, then only the edge (PC-CLI-D1) |
| PC-CLI-OQ2 | Is `/api` Paperclip's REST prefix, and does `--path` include it? | root | **Resolved, PC-D31:** `restApiUrl` is `https://<hostname>/api`; `--path` is the full path; the CLI uses only the origin (PC-CLI-D4) |
| PC-CLI-OQ3 | Can `invoke-api` tell why a request failed when the edge answers `404`? | root | **Resolved, PC-D31:** `404` only for an unknown host or both surfaces off; with one surface off a bound caller gets `403 surface_disabled`. The `404` message points to `describe-collaboration` |
| PC-CLI-OQ4 | Do once-only Paperclip keys need special handling in `invoke-api`? | Agentic Teammates | **Resolved, PC-D31:** the edge passes them through unchanged and never logs bodies; the CLI prints the response as returned (no special case), docs warn |
| PC-CLI-OQ5 | Should `put-collaboration` fill an omitted surface from the current value instead of requiring both flags? | J | Built: both required (explicit, one API call, no hidden read). Revisit after the console and CLI are used together |

### 7.1 Decisions

| # | Question | Decision |
|---|---|---|
| PC-CLI-D1 | How does `invoke-api` learn the host? | One `GetTeamCollaboration` read (a sanctioned extra call beside CLI-10 and CLI-12), taking only the origin of `restApiUrl`; `--collab-url https://<host>` skips the read for a pipeline that already knows it. The read needs `collaboration:read`, which every role that may use the REST API holds |
| PC-CLI-D2 | What does `put-collaboration` do when the version read answers `404`? | Treat it like the existing "PUT that creates" case: send no `If-Match` and show the API's `409 collaboration_agent_missing`, rather than a `404` for an operation the user did not call |
| PC-CLI-D3 | May `invoke-api` send a workload token to the edge? | No. The edge never admits workload principals (§2.1), and a pod's projected token must not travel to a public host. This is a credential-handling rule, not a prediction of `403`, so it does not conflict with "the CLI never skips a call because it predicts a 403" |
| PC-CLI-D4 | `--path` and the `/api` in `restApiUrl` | `--path` is the full path Paperclip serves; the CLI uses only the scheme and host of `restApiUrl` (PC-CLI-OQ2) |
| PC-CLI-D5 | How is a hostile `restApiUrl` contained? | Same transport rules as `endpoint_url` plus no port, userinfo, query or IP literal, no redirects, and an optional profile key `collab_base_domain` that, when set, pins the host to `<label>.<that domain>`. The key is optional so nothing breaks, and recommended in the user docs. The service that holds the credential is the one that supplies the URL, so the residual risk equals trusting `afs-api` (ADR-012) |
| PC-CLI-D6 | Hints for reasons | The CLI prints a short fixed hint under the AWS-style error for each known `details.reason`, from one table in `internal/cli/exec`; hints name commands and permissions only, contain no server text, and an unknown reason prints none. It is help, not a decision: the API's answer and exit 254 are unchanged |
| PC-CLI-D7 | Global options on `invoke-api` | `--query`, `--output`, pagination options, `--cli-input-json`, `--generate-cli-skeleton` and `--wait` exit 252: the command prints the response unchanged, and a half-supported `--query` on non-JSON bodies would hide errors. `--include` exists for headers |
| PC-CLI-D8 | Confirmation before `rotate-collaboration-keys`? | None. The CLI has no interactive confirmations (scripts first); the effects are printed on stderr afterwards and the API audits the call |
| PC-CLI-D9 | Milestone numbering | M7 to M10 map to C1 to C4 so `milestoneOf`, `CurrentMilestone` and `Pending()` keep working per milestone; operations of later milestones exist as generated commands from the re-pin but are not "done" until their milestone |
| PC-CLI-D10 | When is the snapshot re-pinned? | Once, at the start of M7, from the first `-restapi` ref whose `api/openapi.yaml` carries all seven §4.4b operations (the provider lands them together); later `-restapi` releases are re-pinned in their own PRs as usual (contract §8). Until `-restapi` merges, this repository changes no code |
| PC-CLI-D12 | Alignment to contract v1.3.1 | PC-CLI-OQ1 to OQ4 closed by PC-D31; `host_mismatch` removed (PC-D23); priorities of links, connect, disconnect and endpoints per PC-D28; `rotate-secret` per PC-D30; `hostname_held` and swarm `spec.collaboration` added (CLI-US-31.1, CLI-US-35); the apex serves only the callback, so it is never a valid `invoke-api` origin |
| PC-CLI-D11 | Where does this PRD live? | The repository's root `enterprise-agent-factory-service-cli-PRD.md` is a stub that points to the archived spec, so the stub stays and this PRD is extended where it is, in `specs/archive/261121-cli-feature-complete/`. The implementation spec for the change is a new `specs/YYMMDD-paperclip-collaboration/` that links here, and when it starts the PRD moves with it (the stub's own rule) |

### 7.2 Decisions taken while implementing (v1.0.1 run)

See `implementation-notes.md` in the spec directory; each decision is also added here when it changes behaviour.

## 8. Definition of done

- [ ] Command follows naming and option rules; help with an example
- [ ] `--output json` equals the API's fields; golden tests for table and text
- [ ] No secret value, token or key can reach output, logs or history (except `secretKey` once on create)
- [ ] AWS-style errors and exit codes; `AccessDenied` test
- [ ] `make gate` green; stories' acceptance criteria covered
- [ ] For a collaboration command: no AFS token, API key, Paperclip key or secret value in stdout, stderr, `--debug`, arguments or history (sentinel tests green); `X-Afs-Authorization` redacted
- [ ] For a collaboration command: reason hints table updated for any new `details.reason`; no client-side copy of server rules (permissions, ranges, reserved names, enums beyond what the contract lists as an option)
- [ ] The pinned snapshot, generated files, reference docs and goldens regenerated from the contract, never hand-edited; `make contract-check` and `make generate-check` green; `afs --version` shows the new contract minor
- [ ] `Pending()` empty at the new `CurrentMilestone`; every new operation has an `AccessDenied` test and a golden help page
- [ ] `user-docs/` updated and the root `skills/` entry in step; no link into `specs/`
