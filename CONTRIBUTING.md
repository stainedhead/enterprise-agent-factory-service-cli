# Contributing to enterprise-agent-factory-service-cli

Start with the product-wide
[contributing guide](https://github.com/stainedhead/enterprise-agent-factory-service/blob/main/CONTRIBUTING.md);
this file adds what is specific to `afs`. The full rules are in [AGENTS.md](AGENTS.md).

## Before you start

- Read this repository's [PRD](enterprise-agent-factory-service-cli-PRD.md), the
  [API contract](https://github.com/stainedhead/enterprise-agent-factory-service/blob/main/docs/contracts/api-contract.md)
  and [docs/technical-details.md](docs/technical-details.md).
- Create or update a spec in `specs/YYMMDD-<feature>/` and keep its `status.md` current.

## Where a change goes

| Change | Where |
|---|---|
| A new or changed API operation | `make contract-sync RESTAPI_REF=<ref>` (it also runs `make generate`); the command appears from the contract |
| Examples, table columns, shorthands for an operation | `internal/cli/<service>/` overrides |
| Global options, parsing, profiles | `internal/cli`, `internal/config` |
| Credentials, SSO | `internal/auth`, `internal/sso` |
| Output formats, `--query` | `internal/output` |

Command rules (AWS CLI conventions): plural service nouns, kebab-case verb-noun operations that map to an
OpenAPI `operationId`, `--kebab-case` parameters, **one operation, one API call**. The CLI never skips a call
because it predicts a 403; the API decides.

## Make the change

1. **Red, Green, Refactor**, table-driven, against an `httptest` server; no test calls a real API.
2. Golden files for `table`, `text` and `help` output (run the package tests with `-update` to refresh, then review the diff).
3. Every operation keeps an `AccessDenied` test (exit 254, AWS-style message).
4. Run `make docs` if help text changed, so `user-docs/reference/` stays in step.

Review checklist:

- [ ] `--output json` matches the API's field names exactly
- [ ] No token, API key or secret value can reach stdout, stderr, logs or shell history
- [ ] Errors use the AWS-style message and the right exit code
- [ ] Context passed through; requests cancel on Ctrl-C

## Run the gate

```bash
make gate     # tidy, contract-check, generate-check, fmt-check, vet, lint, race tests, coverage (90%), build
./bin/afs --version
```

## Open a pull request

Branch from `main`, reference the spec directory in your commits, and include before and after output for any
change a user will see. CI must be green.
