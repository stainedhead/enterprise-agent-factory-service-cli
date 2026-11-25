# Model endpoints and the team gateway

Since 1.2.0 every team has a model gateway. Agents call OpenAI-compatible **routes**: service endpoints the team opts
in to, and the team's own endpoints. Calls are metered and capped against the team's budget, like Bedrock, and the
upstream credential stays with AFS. Agents never hold it. A service admin also switches whole **providers** on and off,
Bedrock included.

The commands are in the `model-endpoints` service, plus two in `teams` and two in `reports`. The full option lists are
in the [command reference](reference/model-endpoints.md).

## Service admins: the endpoint catalog

```bash
# Add an endpoint. It is created switched off.
afs model-endpoints create-model-endpoint --slug openai-prod --display-name "OpenAI production" \
  --provider openai --base-url https://api.openai.com/v1 --protocols openai_chat openai_responses \
  --location external --pricing-mode per_token --auth-kind bearer --credential-mode per_team \
  --models file://openai-models.json

# Price a model, review proposed prices, then switch the endpoint on.
afs model-endpoints create-endpoint-model-price --endpoint-id mep_… --model-name gpt-4.1 \
  --input-usd-per-mtok 2.00 --output-usd-per-mtok 8.00 --effective-from 2026-12-01T00:00:00Z
afs model-endpoints list-endpoint-price-proposals --status pending
afs model-endpoints approve-endpoint-price-proposal --proposal-id epp_… --note "Matches the provider's page"
afs model-endpoints update-model-endpoint --endpoint-id mep_… --enabled

# See the teams that use it: credential and federation state, gateway version, spend this month.
afs model-endpoints list-endpoint-teams --endpoint-id mep_…
```

`--models` takes JSON or `file://PATH`:
`[{"modelName":"gpt-4.1","maxOutputTokens":32768,"contextTokens":1047576}]`. `--ca-bundle` takes PEM text or
`file://PATH`. A CA bundle is not a secret.

A delete is refused while the endpoint is enabled (`endpoint_enabled`), so switch it off first.

## Credentials

A credential value is read only from a file, from standard input or from a prompt that does not echo. It is never
taken as a plain argument, so it stays out of your shell history. The API never returns a value: you see its state
(`set` or `missing`) and when it was last rotated.

```bash
# A team's key on a per-team endpoint, from a file.
afs model-endpoints put-model-endpoint-credential --endpoint-id mep_… --scope team_… --secret-value file://team.key

# The shared key from standard input.
your-secret-store get openai/shared | afs model-endpoints put-model-endpoint-credential --endpoint-id mep_… --scope _shared --secret-value -

# Rotate. When --secret-value is left out on a terminal, you are asked for the value.
afs model-endpoints rotate-model-endpoint-credential --endpoint-id mep_… --scope team_…

# A team-admin sets the team's key of a team endpoint, or of a per-team service endpoint.
afs model-endpoints put-team-model-endpoint-credential --team-id team_… --endpoint-id mep_… --secret-value file://team.key
```

The API does not return a credential's version, so these commands send no `If-Match` unless you give `--if-match`. To
replace a service credential that is already set, use `rotate-model-endpoint-credential`. Under `--debug` both the
value and the `Authorization` header are shown as `[redacted]`.

Azure and Google endpoints can use workload identity federation instead of a key. In that case no secret is stored:

```bash
afs model-endpoints put-endpoint-federation --endpoint-id mep_… --team-id team_… \
  --tenant-id 11111111-2222-3333-4444-555555555555 --client-id 66666666-7777-8888-9999-000000000000
```

The response names the issuer, subject and audience to trust on the provider's side.

## Providers

```bash
afs model-endpoints list-inference-providers

# Turn Azure on. Its endpoints are then selected per team by service admins.
afs model-endpoints update-inference-provider --provider azure --enabled --team-selection service_admin

# Turning Bedrock off needs a fallback model. Look at the impact first; nothing is written.
afs model-endpoints update-inference-provider --provider bedrock --no-enabled \
  --fallback-model azure-openai/gpt-5.1 --select-for-all-teams --reason "No AWS Bedrock here" --dry-run
# Then send the same command without --dry-run.
```

With `--dry-run` and table or text output, the impact is also printed on stderr: routes, teams, agents moving to the
fallback, degraded agents, and teams to select.

## Provider statements

Import a provider's monthly statement so spend can be reconciled per team. The statement is always a file. A CSV file
(columns `teamId,model,amountUsd,tokens`) is sent as is, and needs `--period` and `--source`:

```bash
afs model-endpoints import-provider-statement --endpoint-id mep_… --period 2026-11 --source opencode_csv \
  --statement file://opencode-2026-11.csv
```

A JSON file can hold the whole body (`{"period", "source", "rows"}`) or only the rows. Options you give on the command
line take precedence over the file.

## Team-admins: routes for the team

```bash
# The service routes the team uses, and the ones it may opt in to.
afs model-endpoints list-team-model-endpoints --team-id team_…
afs model-endpoints list-team-model-endpoints --team-id team_… --relation available

# Opt in. An endpoint left out of the list is deselected.
afs teams put-team-model-endpoints --team-id team_… --endpoint-ids mep_A mep_B

# Declare the team's own endpoint. This needs mayAddModelEndpoints, which a service admin sets.
afs model-endpoints create-team-model-endpoint --team-id team_… --slug coder-vllm --display-name "Coder vLLM" \
  --base-url http://vllm.afs-t-payments.svc:8000/v1 --protocols openai_chat --location in_namespace \
  --pricing-mode capacity --auth-kind none --owner payments-platform --purpose "Code completion" \
  --models file://coder-models.json
```

A service admin chooses the team's gateway deployment and whether the team may declare endpoints:

```bash
afs teams put-team-gateway-settings --team-id team_… --deployment per_namespace --may-add-model-endpoints
```

## Agents and pools

The model options of `create-agent`, `update-agent`, `create-pool` and `update-pool` also take routes:

```bash
afs agents update-agent --team-id team_… --agent-id agt_… --routes openai-prod coder-vllm \
  --default-model-id openai-prod/gpt-4.1
afs pools update-pool --team-id team_… --pool-id pool_… --routes all
```

`--routes all` follows the team's routes now and later. With a list of slugs, the agent uses only those routes. The
default model can be `<route>/<model>`, or the model's name together with `--default-route`.

## Inference report

```bash
afs reports get-team-inference-report --team-id team_… --range 24h --group-by route,model
afs reports get-team-inference-report --team-id team_… --range 7d --route openai-prod --group-by agent --output csv
afs reports get-inference-report --range mtd --group-by route        # service-wide, platform roles
```

The report shows latency, time to first token, throughput, errors (upstream errors and AFS refusals are counted
separately) and cost. It never includes prompt or completion text. The usage and cost reports also accept `--route` and
`--group-by route`.

## Common refusals

| Reason | What to do |
|---|---|
| `provider_disabled` | The provider is off; a service admin turns it on with `update-inference-provider --enabled` |
| `no_fallback_model` | Add `--fallback-model <route>/<model>` when you turn Bedrock off |
| `needs_service_admin` | The endpoint needs extra terms, or its provider is selected by service admins |
| `team_endpoints_not_allowed` | A service admin allows team endpoints with `put-team-gateway-settings --may-add-model-endpoints` |
| `slug_taken`, `slug_reserved`, `slug_invalid` | Choose another `--slug` (2 to 31 characters: lowercase letters, digits or hyphens, starting with a letter) |
| `base_url_not_https` | Only an `in_namespace` team endpoint may use `http://` |
| `endpoint_enabled` | Switch the endpoint off before you delete it |
