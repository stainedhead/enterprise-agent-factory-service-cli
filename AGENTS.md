# AGENTS.md

Rules for AI agents and human contributors working in this repository.

## Project summary

`enterprise-agent-factory-service-cli` is `afs`, the command line interface to the Enterprise Agent
Factory Service. People and automation use it to manage teams, agents, members, secrets, budgets and the
rest of the service through the REST API. It follows the **AWS CLI design**: `afs <service> <operation>
[--options]`, named profiles, `--output json|table|text|yaml`, `--query` (JMESPath), token pagination,
`--cli-input-json` and `--generate-cli-skeleton`.

It is one child of the [enterprise-agent-factory-service](https://github.com/stainedhead/enterprise-agent-factory-service)
root (PRD: `enterprise-agent-factory-service-initial-PRD.md`; CLI design: `docs/solution-design.md` §11).
It talks only to the API in
[enterprise-agent-factory-service-restapi](https://github.com/stainedhead/enterprise-agent-factory-service-restapi),
whose `api/openapi.yaml` is the contract.

Status: milestones M0-M5 done (generic executor for all 180 operations; per-operation polish through M5). See `specs/` for the plan.

## Toolchain

| Tool | Version | Notes |
|---|---|---|
| Go | 1.27.4 | Pinned by `go.mod`; older Go downloads it automatically |
| golangci-lint | v2.15.0 | v2 config; `make lint` builds it with the module's toolchain |
| govulncheck | v1.8.0 | `make vuln` |

**golangci-lint v2 config rules:** formatters under `formatters:`; no `gosimple`; options under
`linters.settings:`.

## Layout

```
cmd/afs/                  entry point; calls cli.Run and exits with its code (130 on Ctrl-C)
api/                      pinned OpenAPI snapshot (openapi.yaml, SOURCE) and its embed
tools/afsgen/             generates internal/spec/operations_gen.go from the snapshot
tools/docgen/             generates user-docs/reference/ from the command table
internal/spec/            operation descriptors (generated) and lookups
internal/cli/             registry, parsing, global options, overrides, Run
internal/cli/<service>/   one package per service: overrides, examples, columns, local commands
internal/cli/services/    composition root: registers every service, wires exec and help
internal/cli/exec/        the generic executor: build, skeleton, pagination, errors, run
internal/cli/safeupdate/  If-Match version reads
internal/cli/wait/        --wait polling
internal/cli/help/        three-level help
internal/cli/runtime/     per-invocation settings, credentials and client
internal/client/          HTTP client: retries, idempotency, redaction, debug; gen/ holds generated types
internal/auth/            credential chain (API key, credentials file, SSO, workload token)
internal/config/          ~/.afs/config and ~/.afs/credentials, profiles, precedence
internal/output/          json, yaml, table, text, csv renderers, JMESPath, pager
internal/sso/             OIDC device flow, token cache, refresh
internal/secretinput/     secret values from file, stdin or prompt only
internal/redact/          header and JSON redaction
internal/buildinfo/       version stamped by the Makefile
docs/  user-docs/  specs/
```

## Command design rules (AWS CLI conventions)

- **Shape.** `afs <service> <operation>`. Services are plural nouns (`teams`, `agents`, `secrets`,
  `budgets`, `members`, `api-keys`, `audit`, `work-orders`, `models`, `admin`) plus `configure` and `sso`.
  Operations are kebab-case verb-noun: `list-teams`, `describe-team`, `create-agent`, `update-agent`,
  `delete-agent`, `put-budget`, and action verbs such as `suspend-team`, `restart-agent`, `rotate-secret`.
- **One operation, one API call** (plus pagination). Operation names map to OpenAPI `operationId`s
  (`list-teams` ↔ `ListTeams`). Do not add client-side workflows that combine several mutations.
- **Parameters** are `--kebab-case` mirrors of API fields (`--team-id`, `--amount-usd`). Lists take
  space-separated values (`--thresholds 50 80 100`). Structured input takes JSON or `file://` / `fileb://`.
- **Global options** on every command: `--profile`, `--endpoint-url`, `--output`, `--query`,
  `--no-paginate`, `--max-items`, `--page-size`, `--starting-token`, `--cli-input-json`,
  `--generate-cli-skeleton`, `--debug`, `--no-cli-pager`, `--color on|off|auto`, `--version`.
- **Configuration precedence:** command-line option → environment (`AFS_PROFILE`, `AFS_ENDPOINT_URL`,
  `AFS_DEFAULT_OUTPUT`, `AFS_API_KEY`) → profile in `~/.afs/config` → defaults.
- **Pagination:** auto-paginate by default; with `--max-items`, print `NextToken` for `--starting-token`.
- **Output:** `json` is the stable machine contract (field names exactly as the API returns them);
  `table` and `text` are for humans; `--query` applies JMESPath before rendering.
- **Errors** go to stderr as `An error occurred (<Code>) when calling the <Operation> operation: <message>`.
  Exit codes: 0 success, 252 usage/parse error, 254 service error, 255 other (network, config).
- **Help:** `afs help`, `afs <service> help`, `afs <service> <operation> help`, generated from the same
  definitions as the commands. Every operation has an example.

## Security and RBAC rules

- **The API decides access.** The CLI never skips a call because it predicts a 403. It may read
  `GET /v1/me` (cached 5 minutes) to improve messages, for example naming the role required.
- **Credentials:** SSO tokens are cached under `~/.afs/sso/cache` with `0600` permissions; API keys come
  from `AFS_API_KEY` or `~/.afs/credentials`. Never print a token or API key, including with `--debug`
  (redact `Authorization`).
- **Secret values** are accepted only from a file, stdin or an interactive no-echo prompt, never as a
  plain argument that lands in shell history. They are sent once and never echoed or logged.

## Architecture rules

- Clean Architecture: `cmd/afs` → `internal/cli` → `internal/client`, `internal/config`, `internal/output`.
  Command handlers parse flags, call one client method, and render; no business logic.
- `internal/cli` takes `io.Writer`s and an injected client so commands are testable without a network.

## TDD: Red → Green → Refactor

All three phases are mandatory; refactor is not optional.

- Table-driven tests per operation: argument parsing, the request sent (against an `httptest` server),
  the rendering for each `--output`, error mapping and exit codes.
- Golden files for `table` and `text` output and for `help` text (`-update` flag to refresh).
- `internal/...` holds ≥ 90% coverage (`make cover`). No test calls a real API.
- Every operation has a test for the API's `AccessDenied` response (exit 254, AWS-style message).

## Quality gate

All of these must pass before every commit (`make gate` runs them in order):

```bash
go mod tidy
make fmt-check
go vet ./...
make lint
go test -race ./...
make build           # bin/afs, CGO_ENABLED=0
./bin/afs --version
```

## Code design review (before every commit)

- [ ] Command follows the naming and option rules above; help and example present
- [ ] `--output json` matches the API's field names exactly
- [ ] No token, API key or secret value can reach stdout, stderr, logs or shell history
- [ ] Errors use the AWS-style message and correct exit code
- [ ] Context passed through; requests cancel on Ctrl-C
- [ ] No hardcoded endpoints, tenant ids or credentials

## Agent workflow

1. **Understand:** this repository's PRD (`enterprise-agent-factory-service-cli-PRD.md`), the shared contracts in the root's `docs/contracts/`, the root PRD, solution design §11, `api/openapi.yaml` in `-restapi`, and the active spec.
2. **Plan:** the operations, flags and outputs affected.
3. **Red → Green → Refactor.**
4. **Quality gate:** `make gate`.
5. **Design review:** the checklist above.
6. **Document:** `user-docs/` command reference and examples; the root repo's `skills/` entry once released.
7. **Update status:** `specs/<feature>/status.md` after every task. Mandatory.
8. **Final verify:** rerun the quality gate.

## user-docs/ rule

`user-docs/` holds only files that help someone install, configure and use `afs`: installation,
`afs configure` and profiles, SSO login, command reference, examples, troubleshooting. It is not for
design, requirements, spec or process material, and must not link into `specs/`.

## Git

Clear commit messages that reference the spec directory when there is one. Do not commit `bin/`, `dist/`
or any `~/.afs` content. Work on a branch and open a pull request into `main`.

## Dev-Flow Skills

The following skills are available for the developer workflow. Invoke them with the listed slash command.

| Skill | Invocation | Purpose |
|---|---|---|
| `init-dev-flow` | `/init-dev-flow` | One-time repo setup; updates AGENTS.md |
| `create-prd` | `/create-prd [title]` | Interactive PRD authoring; writes to repo root |
| `review-prd` | `/review-prd [file]` | PRD quality review |
| `create-spec` | `/create-spec [prd-file]` | Spec creation from PRD; YYMMDD-prefixed directory |
| `review-spec` | `/review-spec [spec-dir]` | Spec quality review |
| `archive-spec` | `/archive-spec [spec-dir]` | Move completed spec to specs/archive/ |
| `implm-frm-prd` | `/implm-frm-prd [prd-file]` | Implement from PRD (11 steps) |
| `implm-frm-change-dtls` | `/implm-frm-change-dtls [ticket-or-desc]` | Implement from ticket or description (12 steps) |
| `implm-from-spec` | `/implm-from-spec [spec-dir]` | Full 11-step orchestrated implementation |
| `create-review` | `/create-review` | Code and design review artifact |
| `review-code` | `/review-code` | Code review of current branch |
| `write-flow-analys` | `/write-flow-analys [spec-dir]` | Process analysis report (final step of implm-from-spec) |

## Feature Specification Workflow

### Specs Directory Structure

All feature development uses the `specs/` directory for planning and tracking. Each feature gets its own subdirectory named `YYMMDD-<feature-name>/`.

**Directory Structure:**
```
specs/
└── YYMMDD-<feature-name>/
    ├── spec.md                  # Feature specification and requirements
    ├── status.md                # **CRITICAL**: Phase progress tracking (update after each task)
    ├── plan.md                  # Implementation plan and architecture decisions
    ├── tasks.md                 # Task breakdown and progress tracking
    ├── research.md              # Research findings, API docs, examples
    ├── data-dictionary.md       # Data structures, types, schemas
    ├── architecture.md          # System architecture and component design
    └── implementation-notes.md  # Implementation details, gotchas, decisions
```

### Progressive Documentation Build

Documents are created progressively as the feature develops:

**Phase 0: Initial Research (PRD/Feature Research)**
- Input: Product Requirement Document, RFC, or feature research
- Purpose: Understand the problem, gather requirements, identify constraints
- **Update status.md**: Mark Phase 0 as "In Progress"

**Phase 1: Specification (spec.md)**
- Define what the feature does, user requirements, acceptance criteria, goals and non-goals
- **Update status.md**: Mark Phase 0 complete, Phase 1 in progress

**Phase 2: Research & Data Modeling (research.md, data-dictionary.md)**
- Gather documentation, explore existing code, define domain entities and data structures
- **Update status.md**: Mark Phase 1 complete, Phase 2 in progress

**Phase 3: Architecture & Planning (architecture.md, plan.md)**
- Design implementation approach, identify affected layers, document component architecture
- **Update status.md**: Mark Phase 2 complete, Phase 3 in progress

**Phase 4: Task Breakdown (tasks.md)**
- Break down work into concrete, testable tasks with dependencies and estimates
- **Update status.md**: Mark Phase 3 complete, Phase 4 in progress

**Phase 5: Implementation (code + implementation-notes.md)**
- Follow TDD (Red-Green-Refactor), record decisions, document edge cases
- **Update status.md**: After EACH task completion — MANDATORY

**Phase 6: Completion & Archival**
- Update product documentation, move spec to `specs/archive/`, capture lessons learned
- **Verify status.md**: Must show 100% completion before archiving

**MANDATORY**: Update `status.md` after completing each task or phase. This file is the single source of truth for progress tracking.

### Specs Workflow Rules

- **Create the spec directory** before starting any new feature work
- **Update progressively** — specs are living documents, not written once
- **Update status.md ALWAYS** after completing each task, phase, or milestone — this is MANDATORY
- **Reference from commits** — link to the spec directory in commit messages
- **Archive completed specs** — move to `specs/archive/` when 100% complete in status.md
- **Version control** — commit specs alongside code for team visibility
