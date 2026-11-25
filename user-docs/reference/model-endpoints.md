# afs model-endpoints

Generated from `afs model-endpoints help`; do not edit by hand.

```text
AFS MODEL-ENDPOINTS

DESCRIPTION
       

SYNOPSIS
       afs model-endpoints <operation> [parameters]

OPERATIONS
       approve-endpoint-price-proposal Approve endpoint price proposal
       create-endpoint-model-price    Create endpoint model price
       create-model-endpoint          Create model endpoint
       create-team-model-endpoint     Create team model endpoint
       delete-model-endpoint          Delete model endpoint
       delete-team-model-endpoint     Delete team model endpoint
       describe-model-endpoint        Get model endpoint
       import-provider-statement      Import provider statement
       list-endpoint-price-proposals  List endpoint price proposals
       list-endpoint-teams            List endpoint teams
       list-inference-providers       List inference providers
       list-model-endpoints           List model endpoints
       list-team-model-endpoints      List team model endpoints
       put-endpoint-federation        Put endpoint federation
       put-endpoint-models            Put endpoint models
       put-model-endpoint-credential  Put model endpoint credential
       put-team-endpoint-models       Put team endpoint models
       put-team-model-endpoint-credential Put team model endpoint credential
       reject-endpoint-price-proposal Reject endpoint price proposal
       rotate-model-endpoint-credential Rotate model endpoint credential
       update-inference-provider      Update inference provider
       update-model-endpoint          Update model endpoint
       update-team-model-endpoint     Update team model endpoint

To see help text for an operation, run: afs model-endpoints <operation> help
```

## approve-endpoint-price-proposal

```text
AFS MODEL-ENDPOINTS APPROVE-ENDPOINT-PRICE-PROPOSAL

NAME
       model-endpoints approve-endpoint-price-proposal

DESCRIPTION
       Approve endpoint price proposal

       Required permission: admin:pricing

SYNOPSIS
       afs model-endpoints approve-endpoint-price-proposal [options]

OPTIONS
       --proposal-id (string) (required)
       --note (string)
       --idempotency-key (string)
          Idempotency key (default: a new ULID per invocation).

GLOBAL OPTIONS
       --profile                Use a specific profile from ~/.afs/config.
       --endpoint-url           Override the API endpoint URL.
       --output                 Output format: json, yaml, table, text or csv.
       --query                  JMESPath query applied to the result before rendering.
       --no-paginate            Return only the first page of results.
       --max-items              Stop after this many items and print NextToken.
       --page-size              Items per request (maxResults), 1 to 100.
       --starting-token         Resume from a NextToken printed by --max-items.
       --cli-input-json         Read parameters from a JSON document or file://path.
       --generate-cli-skeleton  Print a JSON skeleton of the input and make no call.
       --debug                  Write a redacted request log to stderr.
       --no-cli-pager           Do not send output through a pager.
       --color                  Colored output: on, off or auto.
       --cli-connect-timeout    Connect timeout in seconds (0 for none; default 30).
       --cli-read-timeout       Read timeout in seconds (0 for none; default 60).
       --version                Print the afs version and exit.

EXAMPLES
       Approve a proposal

       afs model-endpoints approve-endpoint-price-proposal --proposal-id epp_01JD2Q3R4S5T6V7W8X9Y0Z0001 --note "Matches the provider's page"

```

## create-endpoint-model-price

```text
AFS MODEL-ENDPOINTS CREATE-ENDPOINT-MODEL-PRICE

NAME
       model-endpoints create-endpoint-model-price

DESCRIPTION
       Create endpoint model price

       Required permission: admin:pricing

SYNOPSIS
       afs model-endpoints create-endpoint-model-price [options]

OPTIONS
       --endpoint-id (string) (required)
       --cached-input-usd-per-mtok (string)
          USD amount as a decimal string with at most 6 fraction digits; never a JSON number
       --effective-from (string) (required)
          RFC 3339 UTC with Z
       --input-usd-per-mtok (string) (required)
          USD amount as a decimal string with at most 6 fraction digits; never a JSON number
       --model-name (string) (required)
       --output-usd-per-mtok (string) (required)
          USD amount as a decimal string with at most 6 fraction digits; never a JSON number
       --request-fee-usd (string)
          USD amount as a decimal string with at most 6 fraction digits; never a JSON number
       --source (string)
          transfer for an internal platform's chargeback rate
          Possible values: manual, transfer
       --idempotency-key (string)
          Idempotency key (default: a new ULID per invocation).

GLOBAL OPTIONS
       --profile                Use a specific profile from ~/.afs/config.
       --endpoint-url           Override the API endpoint URL.
       --output                 Output format: json, yaml, table, text or csv.
       --query                  JMESPath query applied to the result before rendering.
       --no-paginate            Return only the first page of results.
       --max-items              Stop after this many items and print NextToken.
       --page-size              Items per request (maxResults), 1 to 100.
       --starting-token         Resume from a NextToken printed by --max-items.
       --cli-input-json         Read parameters from a JSON document or file://path.
       --generate-cli-skeleton  Print a JSON skeleton of the input and make no call.
       --debug                  Write a redacted request log to stderr.
       --no-cli-pager           Do not send output through a pager.
       --color                  Colored output: on, off or auto.
       --cli-connect-timeout    Connect timeout in seconds (0 for none; default 30).
       --cli-read-timeout       Read timeout in seconds (0 for none; default 60).
       --version                Print the afs version and exit.

EXAMPLES
       Price a model from the first of next month

       afs model-endpoints create-endpoint-model-price --endpoint-id mep_01JD2Q3R4S5T6V7W8X9Y0Z0001 --model-name gpt-4.1 --input-usd-per-mtok 2.00 --output-usd-per-mtok 8.00 --cached-input-usd-per-mtok 0.50 --effective-from 2026-12-01T00:00:00Z

       Record an internal platform's chargeback rate

       afs model-endpoints create-endpoint-model-price --endpoint-id mep_01JD2Q3R4S5T6V7W8X9Y0Z0001 --model-name gpt-4.1 --input-usd-per-mtok 1.80 --output-usd-per-mtok 7.20 --effective-from 2026-12-01T00:00:00Z --source transfer

```

## create-model-endpoint

```text
AFS MODEL-ENDPOINTS CREATE-MODEL-ENDPOINT

NAME
       model-endpoints create-model-endpoint

DESCRIPTION
       Create model endpoint

       The endpoint is created switched off; turn it on with update-model-endpoint --enabled once its models are priced and its credential is set. The base URL must be https and its host must not belong to a provider that is off.

       Required permission: admin:model-endpoints

SYNOPSIS
       afs model-endpoints create-model-endpoint [options]

OPTIONS
       --auth-header (string)
          Header name, required when authKind is header
       --auth-kind (string) (required)
          Possible values: none, bearer, header, azure_federated, google_federated
       --base-url (string) (required)
          https for external and corporate; in_namespace may be http://<service>.<namespace>.svc:<port>
       --ca-bundle (string or JSON)
          PEM bundle of the CA that signs the endpoint's certificate, as text or file://PATH (not a secret).
       --credential-mode (string)
          Possible values: shared, per_team
       --description (string)
       --display-name (string) (required)
       --external-system-id (string)
       --location (string) (required)
          Possible values: external, corporate
       --max-body-mib (integer)
          Range: 1 to 100
       --models (JSON)
          The endpoint's models as JSON or file://PATH: [{"modelName":…, "maxOutputTokens":…, "contextTokens":…, "enabled":…}]; a team endpoint's per_token models also carry "price".
       --pricing-mode (string) (required)
          Possible values: per_token, capacity, free
       --protocols (list) (required)
          Protocols the endpoint speaks: openai_chat always; openai_responses and openai_embeddings add those paths (space-separated).
       --provider (string) (required)
          Provider key of inference_providers (bedrock, openai, openrouter, nous, opencode, azure, google, other) or internal for corporate and in-...
       --requires-terms (boolean)
       --slug (string) (required)
       --timeout-seconds (integer)
          Range: 5 to 3600
       --use-provider-cost (boolean)
       --idempotency-key (string)
          Idempotency key (default: a new ULID per invocation).

GLOBAL OPTIONS
       --profile                Use a specific profile from ~/.afs/config.
       --endpoint-url           Override the API endpoint URL.
       --output                 Output format: json, yaml, table, text or csv.
       --query                  JMESPath query applied to the result before rendering.
       --no-paginate            Return only the first page of results.
       --max-items              Stop after this many items and print NextToken.
       --page-size              Items per request (maxResults), 1 to 100.
       --starting-token         Resume from a NextToken printed by --max-items.
       --cli-input-json         Read parameters from a JSON document or file://path.
       --generate-cli-skeleton  Print a JSON skeleton of the input and make no call.
       --debug                  Write a redacted request log to stderr.
       --no-cli-pager           Do not send output through a pager.
       --color                  Colored output: on, off or auto.
       --cli-connect-timeout    Connect timeout in seconds (0 for none; default 30).
       --cli-read-timeout       Read timeout in seconds (0 for none; default 60).
       --version                Print the afs version and exit.

EXAMPLES
       Add an OpenAI route with a credential per team

       afs model-endpoints create-model-endpoint --slug openai-prod --display-name "OpenAI production" --provider openai --base-url https://api.openai.com/v1 --protocols openai_chat openai_responses --location external --pricing-mode per_token --auth-kind bearer --credential-mode per_team --models file://openai-models.json

       Add a corporate capacity endpoint with its own CA

       afs model-endpoints create-model-endpoint --slug corp-llm --display-name "Corporate LLM" --provider internal --base-url https://llm.corp.example.com/v1 --protocols openai_chat --location corporate --pricing-mode capacity --auth-kind header --auth-header X-Api-Key --ca-bundle file://corp-ca.pem

```

## create-team-model-endpoint

```text
AFS MODEL-ENDPOINTS CREATE-TEAM-MODEL-ENDPOINT

NAME
       model-endpoints create-team-model-endpoint

DESCRIPTION
       Create team model endpoint

       Needs the team's mayAddModelEndpoints (team_endpoints_not_allowed otherwise; a service admin sets it with afs teams put-team-gateway-settings). The endpoint is created switched off.

       Required permission: model-endpoints:write

SYNOPSIS
       afs model-endpoints create-team-model-endpoint [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --auth-header (string)
          Header name, required when authKind is header
       --auth-kind (string) (required)
          Possible values: none, bearer, header, azure_federated, google_federated
       --base-url (string) (required)
          https for external and corporate; in_namespace may be http://<service>.<namespace>.svc:<port>
       --ca-bundle (string or JSON)
          PEM bundle of the CA that signs the endpoint's certificate, as text or file://PATH (not a secret).
       --description (string)
       --display-name (string) (required)
       --location (string) (required)
          Possible values: external, corporate, in_namespace
       --max-body-mib (integer)
          Range: 1 to 100
       --metrics-port (integer)
          Range: 1 to 65535
       --models (JSON)
          The endpoint's models as JSON or file://PATH: [{"modelName":…, "maxOutputTokens":…, "contextTokens":…, "enabled":…}]; a team endpoint's per_token models also carry "price".
       --owner (string) (required)
       --pricing-mode (string) (required)
          Possible values: per_token, capacity, free
       --protocols (list) (required)
          Protocols the endpoint speaks: openai_chat always; openai_responses and openai_embeddings add those paths (space-separated).
       --purpose (string) (required)
       --slug (string) (required)
          Unique in the team and not a service slug
       --timeout-seconds (integer)
          Range: 5 to 3600
       --idempotency-key (string)
          Idempotency key (default: a new ULID per invocation).

GLOBAL OPTIONS
       --profile                Use a specific profile from ~/.afs/config.
       --endpoint-url           Override the API endpoint URL.
       --output                 Output format: json, yaml, table, text or csv.
       --query                  JMESPath query applied to the result before rendering.
       --no-paginate            Return only the first page of results.
       --max-items              Stop after this many items and print NextToken.
       --page-size              Items per request (maxResults), 1 to 100.
       --starting-token         Resume from a NextToken printed by --max-items.
       --cli-input-json         Read parameters from a JSON document or file://path.
       --generate-cli-skeleton  Print a JSON skeleton of the input and make no call.
       --debug                  Write a redacted request log to stderr.
       --no-cli-pager           Do not send output through a pager.
       --color                  Colored output: on, off or auto.
       --cli-connect-timeout    Connect timeout in seconds (0 for none; default 30).
       --cli-read-timeout       Read timeout in seconds (0 for none; default 60).
       --version                Print the afs version and exit.

EXAMPLES
       Declare a vLLM server in the team's namespace

       afs model-endpoints create-team-model-endpoint --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --slug coder-vllm --display-name "Coder vLLM" --base-url http://vllm.afs-t-payments.svc:8000/v1 --protocols openai_chat --location in_namespace --pricing-mode capacity --auth-kind none --owner payments-platform --purpose "Code completion for the payments agents" --models file://coder-models.json

```

## delete-model-endpoint

```text
AFS MODEL-ENDPOINTS DELETE-MODEL-ENDPOINT

NAME
       model-endpoints delete-model-endpoint

DESCRIPTION
       Delete model endpoint

       Refused while the endpoint is enabled (endpoint_enabled): switch it off first.

       Required permission: admin:model-endpoints

SYNOPSIS
       afs model-endpoints delete-model-endpoint [options]

OPTIONS
       --endpoint-id (string) (required)
       --if-match (string)
          Apply only if the resource is at this version (safe update).

GLOBAL OPTIONS
       --profile                Use a specific profile from ~/.afs/config.
       --endpoint-url           Override the API endpoint URL.
       --output                 Output format: json, yaml, table, text or csv.
       --query                  JMESPath query applied to the result before rendering.
       --no-paginate            Return only the first page of results.
       --max-items              Stop after this many items and print NextToken.
       --page-size              Items per request (maxResults), 1 to 100.
       --starting-token         Resume from a NextToken printed by --max-items.
       --cli-input-json         Read parameters from a JSON document or file://path.
       --generate-cli-skeleton  Print a JSON skeleton of the input and make no call.
       --debug                  Write a redacted request log to stderr.
       --no-cli-pager           Do not send output through a pager.
       --color                  Colored output: on, off or auto.
       --cli-connect-timeout    Connect timeout in seconds (0 for none; default 30).
       --cli-read-timeout       Read timeout in seconds (0 for none; default 60).
       --version                Print the afs version and exit.

EXAMPLES
       Switch an endpoint off, then delete it

       afs model-endpoints delete-model-endpoint --endpoint-id mep_01JD2Q3R4S5T6V7W8X9Y0Z0001

```

## delete-team-model-endpoint

```text
AFS MODEL-ENDPOINTS DELETE-TEAM-MODEL-ENDPOINT

NAME
       model-endpoints delete-team-model-endpoint

DESCRIPTION
       Delete team model endpoint

       Required permission: model-endpoints:write

SYNOPSIS
       afs model-endpoints delete-team-model-endpoint [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --endpoint-id (string) (required)
       --if-match (string)
          Apply only if the resource is at this version (safe update).

GLOBAL OPTIONS
       --profile                Use a specific profile from ~/.afs/config.
       --endpoint-url           Override the API endpoint URL.
       --output                 Output format: json, yaml, table, text or csv.
       --query                  JMESPath query applied to the result before rendering.
       --no-paginate            Return only the first page of results.
       --max-items              Stop after this many items and print NextToken.
       --page-size              Items per request (maxResults), 1 to 100.
       --starting-token         Resume from a NextToken printed by --max-items.
       --cli-input-json         Read parameters from a JSON document or file://path.
       --generate-cli-skeleton  Print a JSON skeleton of the input and make no call.
       --debug                  Write a redacted request log to stderr.
       --no-cli-pager           Do not send output through a pager.
       --color                  Colored output: on, off or auto.
       --cli-connect-timeout    Connect timeout in seconds (0 for none; default 30).
       --cli-read-timeout       Read timeout in seconds (0 for none; default 60).
       --version                Print the afs version and exit.

EXAMPLES
       Delete a team endpoint

       afs model-endpoints delete-team-model-endpoint --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --endpoint-id mep_01JD2Q3R4S5T6V7W8X9Y0Z0002

```

## describe-model-endpoint

```text
AFS MODEL-ENDPOINTS DESCRIBE-MODEL-ENDPOINT

NAME
       model-endpoints describe-model-endpoint

DESCRIPTION
       Get model endpoint

       Required permission: admin:model-endpoints

SYNOPSIS
       afs model-endpoints describe-model-endpoint [options]

OPTIONS
       --endpoint-id (string) (required)

GLOBAL OPTIONS
       --profile                Use a specific profile from ~/.afs/config.
       --endpoint-url           Override the API endpoint URL.
       --output                 Output format: json, yaml, table, text or csv.
       --query                  JMESPath query applied to the result before rendering.
       --no-paginate            Return only the first page of results.
       --max-items              Stop after this many items and print NextToken.
       --page-size              Items per request (maxResults), 1 to 100.
       --starting-token         Resume from a NextToken printed by --max-items.
       --cli-input-json         Read parameters from a JSON document or file://path.
       --generate-cli-skeleton  Print a JSON skeleton of the input and make no call.
       --debug                  Write a redacted request log to stderr.
       --no-cli-pager           Do not send output through a pager.
       --color                  Colored output: on, off or auto.
       --cli-connect-timeout    Connect timeout in seconds (0 for none; default 30).
       --cli-read-timeout       Read timeout in seconds (0 for none; default 60).
       --version                Print the afs version and exit.

EXAMPLES
       Show an endpoint

       afs model-endpoints describe-model-endpoint --endpoint-id mep_01JD2Q3R4S5T6V7W8X9Y0Z0001

       Show its models and prices

       afs model-endpoints describe-model-endpoint --endpoint-id mep_01JD2Q3R4S5T6V7W8X9Y0Z0001 --query models --output table

```

## import-provider-statement

```text
AFS MODEL-ENDPOINTS IMPORT-PROVIDER-STATEMENT

NAME
       model-endpoints import-provider-statement

DESCRIPTION
       Import provider statement

       Imports a provider's monthly statement so spend can be reconciled per team (variance over 5% is alerted). A route shows Not reconciled until its first import.

       Required permission: admin:pricing

SYNOPSIS
       afs model-endpoints import-provider-statement [options]

OPTIONS
       --endpoint-id (string) (required)
       --period (string)
          The statement's month, YYYY-MM.
       --source (string)
          The statement's format: nous_statement or opencode_csv.
          Possible values: nous_statement, opencode_csv
       --idempotency-key (string)
          Idempotency key (default: a new ULID per invocation).
       --statement (string)
          The provider's statement, file://PATH only: a CSV file (columns teamId,model,amountUsd,tokens) is sent as text/csv; a JSON file is either the whole body {period, source, rows} or the array of rows.

GLOBAL OPTIONS
       --profile                Use a specific profile from ~/.afs/config.
       --endpoint-url           Override the API endpoint URL.
       --output                 Output format: json, yaml, table, text or csv.
       --query                  JMESPath query applied to the result before rendering.
       --no-paginate            Return only the first page of results.
       --max-items              Stop after this many items and print NextToken.
       --page-size              Items per request (maxResults), 1 to 100.
       --starting-token         Resume from a NextToken printed by --max-items.
       --cli-input-json         Read parameters from a JSON document or file://path.
       --generate-cli-skeleton  Print a JSON skeleton of the input and make no call.
       --debug                  Write a redacted request log to stderr.
       --no-cli-pager           Do not send output through a pager.
       --color                  Colored output: on, off or auto.
       --cli-connect-timeout    Connect timeout in seconds (0 for none; default 30).
       --cli-read-timeout       Read timeout in seconds (0 for none; default 60).
       --version                Print the afs version and exit.

EXAMPLES
       Import OpenCode's usage CSV for November

       afs model-endpoints import-provider-statement --endpoint-id mep_01JD2Q3R4S5T6V7W8X9Y0Z0001 --period 2026-11 --source opencode_csv --statement file://opencode-2026-11.csv

       Import a statement prepared as JSON

       afs model-endpoints import-provider-statement --endpoint-id mep_01JD2Q3R4S5T6V7W8X9Y0Z0001 --statement file://statement-2026-11.json

```

## list-endpoint-price-proposals

```text
AFS MODEL-ENDPOINTS LIST-ENDPOINT-PRICE-PROPOSALS

NAME
       model-endpoints list-endpoint-price-proposals

DESCRIPTION
       List endpoint price proposals

       Required permission: admin:pricing

SYNOPSIS
       afs model-endpoints list-endpoint-price-proposals [options]

OPTIONS
       --status (string)
          Possible values: pending, approved, rejected, superseded
       --endpoint-id (string)

GLOBAL OPTIONS
       --profile                Use a specific profile from ~/.afs/config.
       --endpoint-url           Override the API endpoint URL.
       --output                 Output format: json, yaml, table, text or csv.
       --query                  JMESPath query applied to the result before rendering.
       --no-paginate            Return only the first page of results.
       --max-items              Stop after this many items and print NextToken.
       --page-size              Items per request (maxResults), 1 to 100.
       --starting-token         Resume from a NextToken printed by --max-items.
       --cli-input-json         Read parameters from a JSON document or file://path.
       --generate-cli-skeleton  Print a JSON skeleton of the input and make no call.
       --debug                  Write a redacted request log to stderr.
       --no-cli-pager           Do not send output through a pager.
       --color                  Colored output: on, off or auto.
       --cli-connect-timeout    Connect timeout in seconds (0 for none; default 30).
       --cli-read-timeout       Read timeout in seconds (0 for none; default 60).
       --version                Print the afs version and exit.

EXAMPLES
       List pending endpoint price proposals

       afs model-endpoints list-endpoint-price-proposals --status pending

       List one endpoint's proposals

       afs model-endpoints list-endpoint-price-proposals --endpoint-id mep_01JD2Q3R4S5T6V7W8X9Y0Z0001

```

## list-endpoint-teams

```text
AFS MODEL-ENDPOINTS LIST-ENDPOINT-TEAMS

NAME
       model-endpoints list-endpoint-teams

DESCRIPTION
       List endpoint teams

       Required permission: admin:model-endpoints

SYNOPSIS
       afs model-endpoints list-endpoint-teams [options]

OPTIONS
       --endpoint-id (string) (required)

GLOBAL OPTIONS
       --profile                Use a specific profile from ~/.afs/config.
       --endpoint-url           Override the API endpoint URL.
       --output                 Output format: json, yaml, table, text or csv.
       --query                  JMESPath query applied to the result before rendering.
       --no-paginate            Return only the first page of results.
       --max-items              Stop after this many items and print NextToken.
       --page-size              Items per request (maxResults), 1 to 100.
       --starting-token         Resume from a NextToken printed by --max-items.
       --cli-input-json         Read parameters from a JSON document or file://path.
       --generate-cli-skeleton  Print a JSON skeleton of the input and make no call.
       --debug                  Write a redacted request log to stderr.
       --no-cli-pager           Do not send output through a pager.
       --color                  Colored output: on, off or auto.
       --cli-connect-timeout    Connect timeout in seconds (0 for none; default 30).
       --cli-read-timeout       Read timeout in seconds (0 for none; default 60).
       --version                Print the afs version and exit.

EXAMPLES
       Show which teams use an endpoint, with credential and federation state and spend

       afs model-endpoints list-endpoint-teams --endpoint-id mep_01JD2Q3R4S5T6V7W8X9Y0Z0001

       List the teams whose credential is missing

       afs model-endpoints list-endpoint-teams --endpoint-id mep_01JD2Q3R4S5T6V7W8X9Y0Z0001 --query "items[?credentialState=='missing'].teamId" --output text

```

## list-inference-providers

```text
AFS MODEL-ENDPOINTS LIST-INFERENCE-PROVIDERS

NAME
       model-endpoints list-inference-providers

DESCRIPTION
       List inference providers

       Required permission: admin:model-endpoints

SYNOPSIS
       afs model-endpoints list-inference-providers [options]

GLOBAL OPTIONS
       --profile                Use a specific profile from ~/.afs/config.
       --endpoint-url           Override the API endpoint URL.
       --output                 Output format: json, yaml, table, text or csv.
       --query                  JMESPath query applied to the result before rendering.
       --no-paginate            Return only the first page of results.
       --max-items              Stop after this many items and print NextToken.
       --page-size              Items per request (maxResults), 1 to 100.
       --starting-token         Resume from a NextToken printed by --max-items.
       --cli-input-json         Read parameters from a JSON document or file://path.
       --generate-cli-skeleton  Print a JSON skeleton of the input and make no call.
       --debug                  Write a redacted request log to stderr.
       --no-cli-pager           Do not send output through a pager.
       --color                  Colored output: on, off or auto.
       --cli-connect-timeout    Connect timeout in seconds (0 for none; default 30).
       --cli-read-timeout       Read timeout in seconds (0 for none; default 60).
       --version                Print the afs version and exit.

EXAMPLES
       List the providers with their hosts, routes and teams

       afs model-endpoints list-inference-providers

       Print the hosts of each provider

       afs model-endpoints list-inference-providers --query 'items[].[provider,join(`,`,hosts)]' --output text

```

## list-model-endpoints

```text
AFS MODEL-ENDPOINTS LIST-MODEL-ENDPOINTS

NAME
       model-endpoints list-model-endpoints

DESCRIPTION
       List model endpoints

       Required permission: admin:model-endpoints

SYNOPSIS
       afs model-endpoints list-model-endpoints [options]

OPTIONS
       --include-teams (boolean)
          Also list team-declared endpoints
       --provider (string)
       --location (string)
          Possible values: external, corporate, in_namespace
       --pricing-mode (string)
          Possible values: per_token, capacity, free

GLOBAL OPTIONS
       --profile                Use a specific profile from ~/.afs/config.
       --endpoint-url           Override the API endpoint URL.
       --output                 Output format: json, yaml, table, text or csv.
       --query                  JMESPath query applied to the result before rendering.
       --no-paginate            Return only the first page of results.
       --max-items              Stop after this many items and print NextToken.
       --page-size              Items per request (maxResults), 1 to 100.
       --starting-token         Resume from a NextToken printed by --max-items.
       --cli-input-json         Read parameters from a JSON document or file://path.
       --generate-cli-skeleton  Print a JSON skeleton of the input and make no call.
       --debug                  Write a redacted request log to stderr.
       --no-cli-pager           Do not send output through a pager.
       --color                  Colored output: on, off or auto.
       --cli-connect-timeout    Connect timeout in seconds (0 for none; default 30).
       --cli-read-timeout       Read timeout in seconds (0 for none; default 60).
       --version                Print the afs version and exit.

EXAMPLES
       List the service's model endpoints

       afs model-endpoints list-model-endpoints

       Include team-declared endpoints and keep only per-token external routes

       afs model-endpoints list-model-endpoints --include-teams --location external --pricing-mode per_token

```

## list-team-model-endpoints

```text
AFS MODEL-ENDPOINTS LIST-TEAM-MODEL-ENDPOINTS

NAME
       model-endpoints list-team-model-endpoints

DESCRIPTION
       List team model endpoints

       Required permission: model-endpoints:read

SYNOPSIS
       afs model-endpoints list-team-model-endpoints [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --relation (string)
          Possible values: selected, own, available

GLOBAL OPTIONS
       --profile                Use a specific profile from ~/.afs/config.
       --endpoint-url           Override the API endpoint URL.
       --output                 Output format: json, yaml, table, text or csv.
       --query                  JMESPath query applied to the result before rendering.
       --no-paginate            Return only the first page of results.
       --max-items              Stop after this many items and print NextToken.
       --page-size              Items per request (maxResults), 1 to 100.
       --starting-token         Resume from a NextToken printed by --max-items.
       --cli-input-json         Read parameters from a JSON document or file://path.
       --generate-cli-skeleton  Print a JSON skeleton of the input and make no call.
       --debug                  Write a redacted request log to stderr.
       --no-cli-pager           Do not send output through a pager.
       --color                  Colored output: on, off or auto.
       --cli-connect-timeout    Connect timeout in seconds (0 for none; default 30).
       --cli-read-timeout       Read timeout in seconds (0 for none; default 60).
       --version                Print the afs version and exit.

EXAMPLES
       List the routes the team uses and may use

       afs model-endpoints list-team-model-endpoints --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001

       List the service endpoints the team may opt in to

       afs model-endpoints list-team-model-endpoints --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --relation available

```

## put-endpoint-federation

```text
AFS MODEL-ENDPOINTS PUT-ENDPOINT-FEDERATION

NAME
       model-endpoints put-endpoint-federation

DESCRIPTION
       Put endpoint federation

       Workload identity federation for an azure_federated or google_federated endpoint: metadata only, no secret is stored. The response names the issuer, subject and audience (or role) to trust on the provider's side.

       Required permission: admin:model-endpoints

SYNOPSIS
       afs model-endpoints put-endpoint-federation [options]

OPTIONS
       --endpoint-id (string) (required)
       --team-id (string) (required)
          The team whose gateway is federated
       --client-id (string)
          Azure: the managed identity or app registration
       --location (string)
          Google, for example us-central1 or global
       --pool-id (string)
          Google workload identity pool
       --project-id (string)
          Google
       --project-number (string)
          Google
       --provider-id (string)
          Google pool provider
       --tenant-id (string)
          Azure
       --if-match (string)
          Apply only if the resource is at this version (safe update).

GLOBAL OPTIONS
       --profile                Use a specific profile from ~/.afs/config.
       --endpoint-url           Override the API endpoint URL.
       --output                 Output format: json, yaml, table, text or csv.
       --query                  JMESPath query applied to the result before rendering.
       --no-paginate            Return only the first page of results.
       --max-items              Stop after this many items and print NextToken.
       --page-size              Items per request (maxResults), 1 to 100.
       --starting-token         Resume from a NextToken printed by --max-items.
       --cli-input-json         Read parameters from a JSON document or file://path.
       --generate-cli-skeleton  Print a JSON skeleton of the input and make no call.
       --debug                  Write a redacted request log to stderr.
       --no-cli-pager           Do not send output through a pager.
       --color                  Colored output: on, off or auto.
       --cli-connect-timeout    Connect timeout in seconds (0 for none; default 30).
       --cli-read-timeout       Read timeout in seconds (0 for none; default 60).
       --version                Print the afs version and exit.

EXAMPLES
       Federate a team's gateway with an Azure app registration

       afs model-endpoints put-endpoint-federation --endpoint-id mep_01JD2Q3R4S5T6V7W8X9Y0Z0001 --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --tenant-id 11111111-2222-3333-4444-555555555555 --client-id 66666666-7777-8888-9999-000000000000

       Federate with a Google workload identity pool

       afs model-endpoints put-endpoint-federation --endpoint-id mep_01JD2Q3R4S5T6V7W8X9Y0Z0001 --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --project-number 123456789012 --pool-id afs-gateways --provider-id eks --project-id payments-ai --location us-central1

```

## put-endpoint-models

```text
AFS MODEL-ENDPOINTS PUT-ENDPOINT-MODELS

NAME
       model-endpoints put-endpoint-models

DESCRIPTION
       Put endpoint models

       Replaces the endpoint's model list. A model is callable once it is enabled and, on a per_token endpoint, priced.

       Required permission: admin:model-endpoints

SYNOPSIS
       afs model-endpoints put-endpoint-models [options]

OPTIONS
       --endpoint-id (string) (required)
       --models (JSON) (required)
          The endpoint's models as JSON or file://PATH: [{"modelName":…, "maxOutputTokens":…, "contextTokens":…, "enabled":…}]; a team endpoint's per_token models also carry "price".
       --if-match (string)
          Apply only if the resource is at this version (safe update).

GLOBAL OPTIONS
       --profile                Use a specific profile from ~/.afs/config.
       --endpoint-url           Override the API endpoint URL.
       --output                 Output format: json, yaml, table, text or csv.
       --query                  JMESPath query applied to the result before rendering.
       --no-paginate            Return only the first page of results.
       --max-items              Stop after this many items and print NextToken.
       --page-size              Items per request (maxResults), 1 to 100.
       --starting-token         Resume from a NextToken printed by --max-items.
       --cli-input-json         Read parameters from a JSON document or file://path.
       --generate-cli-skeleton  Print a JSON skeleton of the input and make no call.
       --debug                  Write a redacted request log to stderr.
       --no-cli-pager           Do not send output through a pager.
       --color                  Colored output: on, off or auto.
       --cli-connect-timeout    Connect timeout in seconds (0 for none; default 30).
       --cli-read-timeout       Read timeout in seconds (0 for none; default 60).
       --version                Print the afs version and exit.

EXAMPLES
       Replace the endpoint's models from a file

       afs model-endpoints put-endpoint-models --endpoint-id mep_01JD2Q3R4S5T6V7W8X9Y0Z0001 --models file://openai-models.json

```

## put-model-endpoint-credential

```text
AFS MODEL-ENDPOINTS PUT-MODEL-ENDPOINT-CREDENTIAL

NAME
       model-endpoints put-model-endpoint-credential

DESCRIPTION
       Put model endpoint credential

       --scope is a team id for a per_team endpoint, or _shared for a shared one.

       The value is read from file://PATH, fileb://PATH or - (standard input), or asked for without echo when it is omitted on a terminal; it is never a plain argument, and the API never returns it (only its state and last rotation). The credential's version is not returned, so no If-Match is sent unless --if-match is given.

       To replace a value that is already set without knowing its version, use rotate-model-endpoint-credential.

       Required permission: admin:model-endpoints

SYNOPSIS
       afs model-endpoints put-model-endpoint-credential [options]

OPTIONS
       --endpoint-id (string) (required)
       --scope (string) (required)
          A team id (per_team credentials) or _shared
       --secret-value (file://PATH | fileb://PATH | -; prompted without echo when omitted on a terminal)
          Sent once, passed to Secrets Manager and dropped; never returned, logged or audited
          Never pass the value itself on the command line.
       --if-match (string)
          Apply only if the resource is at this version (safe update).

GLOBAL OPTIONS
       --profile                Use a specific profile from ~/.afs/config.
       --endpoint-url           Override the API endpoint URL.
       --output                 Output format: json, yaml, table, text or csv.
       --query                  JMESPath query applied to the result before rendering.
       --no-paginate            Return only the first page of results.
       --max-items              Stop after this many items and print NextToken.
       --page-size              Items per request (maxResults), 1 to 100.
       --starting-token         Resume from a NextToken printed by --max-items.
       --cli-input-json         Read parameters from a JSON document or file://path.
       --generate-cli-skeleton  Print a JSON skeleton of the input and make no call.
       --debug                  Write a redacted request log to stderr.
       --no-cli-pager           Do not send output through a pager.
       --color                  Colored output: on, off or auto.
       --cli-connect-timeout    Connect timeout in seconds (0 for none; default 30).
       --cli-read-timeout       Read timeout in seconds (0 for none; default 60).
       --version                Print the afs version and exit.

EXAMPLES
       Set a team's key on a per-team endpoint from a file

       afs model-endpoints put-model-endpoint-credential --endpoint-id mep_01JD2Q3R4S5T6V7W8X9Y0Z0001 --scope team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --secret-value file://openai-payments.key

       Set the shared key from standard input

       afs model-endpoints put-model-endpoint-credential --endpoint-id mep_01JD2Q3R4S5T6V7W8X9Y0Z0001 --scope _shared --secret-value -

```

## put-team-endpoint-models

```text
AFS MODEL-ENDPOINTS PUT-TEAM-ENDPOINT-MODELS

NAME
       model-endpoints put-team-endpoint-models

DESCRIPTION
       Put team endpoint models

       Replaces the team endpoint's models. On a per_token endpoint each model carries its price inline; capacity and free endpoints take none.

       Required permission: model-endpoints:write

SYNOPSIS
       afs model-endpoints put-team-endpoint-models [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --endpoint-id (string) (required)
       --models (JSON) (required)
          The endpoint's models as JSON or file://PATH: [{"modelName":…, "maxOutputTokens":…, "contextTokens":…, "enabled":…}]; a team endpoint's per_token models also carry "price".
       --if-match (string)
          Apply only if the resource is at this version (safe update).

GLOBAL OPTIONS
       --profile                Use a specific profile from ~/.afs/config.
       --endpoint-url           Override the API endpoint URL.
       --output                 Output format: json, yaml, table, text or csv.
       --query                  JMESPath query applied to the result before rendering.
       --no-paginate            Return only the first page of results.
       --max-items              Stop after this many items and print NextToken.
       --page-size              Items per request (maxResults), 1 to 100.
       --starting-token         Resume from a NextToken printed by --max-items.
       --cli-input-json         Read parameters from a JSON document or file://path.
       --generate-cli-skeleton  Print a JSON skeleton of the input and make no call.
       --debug                  Write a redacted request log to stderr.
       --no-cli-pager           Do not send output through a pager.
       --color                  Colored output: on, off or auto.
       --cli-connect-timeout    Connect timeout in seconds (0 for none; default 30).
       --cli-read-timeout       Read timeout in seconds (0 for none; default 60).
       --version                Print the afs version and exit.

EXAMPLES
       Replace a team endpoint's models from a file

       afs model-endpoints put-team-endpoint-models --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --endpoint-id mep_01JD2Q3R4S5T6V7W8X9Y0Z0002 --models file://coder-models.json

```

## put-team-model-endpoint-credential

```text
AFS MODEL-ENDPOINTS PUT-TEAM-MODEL-ENDPOINT-CREDENTIAL

NAME
       model-endpoints put-team-model-endpoint-credential

DESCRIPTION
       Put team model endpoint credential

       Sets the team's credential of a team endpoint, or of a per_team service endpoint the team opted in to.

       The value is read from file://PATH, fileb://PATH or - (standard input), or asked for without echo when it is omitted on a terminal; it is never a plain argument, and the API never returns it (only its state and last rotation). The credential's version is not returned, so no If-Match is sent unless --if-match is given.

       Required permission: model-endpoints:write

SYNOPSIS
       afs model-endpoints put-team-model-endpoint-credential [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --endpoint-id (string) (required)
       --secret-value (file://PATH | fileb://PATH | -; prompted without echo when omitted on a terminal)
          Sent once, passed to Secrets Manager and dropped; never returned, logged or audited
          Never pass the value itself on the command line.
       --if-match (string)
          Apply only if the resource is at this version (safe update).

GLOBAL OPTIONS
       --profile                Use a specific profile from ~/.afs/config.
       --endpoint-url           Override the API endpoint URL.
       --output                 Output format: json, yaml, table, text or csv.
       --query                  JMESPath query applied to the result before rendering.
       --no-paginate            Return only the first page of results.
       --max-items              Stop after this many items and print NextToken.
       --page-size              Items per request (maxResults), 1 to 100.
       --starting-token         Resume from a NextToken printed by --max-items.
       --cli-input-json         Read parameters from a JSON document or file://path.
       --generate-cli-skeleton  Print a JSON skeleton of the input and make no call.
       --debug                  Write a redacted request log to stderr.
       --no-cli-pager           Do not send output through a pager.
       --color                  Colored output: on, off or auto.
       --cli-connect-timeout    Connect timeout in seconds (0 for none; default 30).
       --cli-read-timeout       Read timeout in seconds (0 for none; default 60).
       --version                Print the afs version and exit.

EXAMPLES
       Set the team's key from a file

       afs model-endpoints put-team-model-endpoint-credential --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --endpoint-id mep_01JD2Q3R4S5T6V7W8X9Y0Z0001 --secret-value file://openai-payments.key

       Set it from standard input (piped from a secret store)

       afs model-endpoints put-team-model-endpoint-credential --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --endpoint-id mep_01JD2Q3R4S5T6V7W8X9Y0Z0001 --secret-value -

```

## reject-endpoint-price-proposal

```text
AFS MODEL-ENDPOINTS REJECT-ENDPOINT-PRICE-PROPOSAL

NAME
       model-endpoints reject-endpoint-price-proposal

DESCRIPTION
       Reject endpoint price proposal

       Required permission: admin:pricing

SYNOPSIS
       afs model-endpoints reject-endpoint-price-proposal [options]

OPTIONS
       --proposal-id (string) (required)
       --note (string)
       --idempotency-key (string)
          Idempotency key (default: a new ULID per invocation).

GLOBAL OPTIONS
       --profile                Use a specific profile from ~/.afs/config.
       --endpoint-url           Override the API endpoint URL.
       --output                 Output format: json, yaml, table, text or csv.
       --query                  JMESPath query applied to the result before rendering.
       --no-paginate            Return only the first page of results.
       --max-items              Stop after this many items and print NextToken.
       --page-size              Items per request (maxResults), 1 to 100.
       --starting-token         Resume from a NextToken printed by --max-items.
       --cli-input-json         Read parameters from a JSON document or file://path.
       --generate-cli-skeleton  Print a JSON skeleton of the input and make no call.
       --debug                  Write a redacted request log to stderr.
       --no-cli-pager           Do not send output through a pager.
       --color                  Colored output: on, off or auto.
       --cli-connect-timeout    Connect timeout in seconds (0 for none; default 30).
       --cli-read-timeout       Read timeout in seconds (0 for none; default 60).
       --version                Print the afs version and exit.

EXAMPLES
       Reject a proposal

       afs model-endpoints reject-endpoint-price-proposal --proposal-id epp_01JD2Q3R4S5T6V7W8X9Y0Z0001 --note "Waiting for the contract rate"

```

## rotate-model-endpoint-credential

```text
AFS MODEL-ENDPOINTS ROTATE-MODEL-ENDPOINT-CREDENTIAL

NAME
       model-endpoints rotate-model-endpoint-credential

DESCRIPTION
       Rotate model endpoint credential

       Replaces the value that is set (NotFound when none is); the previous value is no longer used.

       The value is read from file://PATH, fileb://PATH or - (standard input), or asked for without echo when it is omitted on a terminal; it is never a plain argument.

       Required permission: admin:model-endpoints

SYNOPSIS
       afs model-endpoints rotate-model-endpoint-credential [options]

OPTIONS
       --endpoint-id (string) (required)
       --scope (string) (required)
          A team id (per_team credentials) or _shared
       --secret-value (file://PATH | fileb://PATH | -; prompted without echo when omitted on a terminal)
          Sent once, passed to Secrets Manager and dropped; never returned, logged or audited
          Never pass the value itself on the command line.
       --idempotency-key (string)
          Idempotency key (default: a new ULID per invocation).

GLOBAL OPTIONS
       --profile                Use a specific profile from ~/.afs/config.
       --endpoint-url           Override the API endpoint URL.
       --output                 Output format: json, yaml, table, text or csv.
       --query                  JMESPath query applied to the result before rendering.
       --no-paginate            Return only the first page of results.
       --max-items              Stop after this many items and print NextToken.
       --page-size              Items per request (maxResults), 1 to 100.
       --starting-token         Resume from a NextToken printed by --max-items.
       --cli-input-json         Read parameters from a JSON document or file://path.
       --generate-cli-skeleton  Print a JSON skeleton of the input and make no call.
       --debug                  Write a redacted request log to stderr.
       --no-cli-pager           Do not send output through a pager.
       --color                  Colored output: on, off or auto.
       --cli-connect-timeout    Connect timeout in seconds (0 for none; default 30).
       --cli-read-timeout       Read timeout in seconds (0 for none; default 60).
       --version                Print the afs version and exit.

EXAMPLES
       Rotate a team's key, asked for without echo

       afs model-endpoints rotate-model-endpoint-credential --endpoint-id mep_01JD2Q3R4S5T6V7W8X9Y0Z0001 --scope team_01JD2Q3R4S5T6V7W8X9Y0Z0001

       Rotate the shared key from a file

       afs model-endpoints rotate-model-endpoint-credential --endpoint-id mep_01JD2Q3R4S5T6V7W8X9Y0Z0001 --scope _shared --secret-value file://new.key

```

## update-inference-provider

```text
AFS MODEL-ENDPOINTS UPDATE-INFERENCE-PROVIDER

NAME
       model-endpoints update-inference-provider

DESCRIPTION
       Update inference provider

       Turning a provider off stops every route on its hosts within a minute, service and team endpoints alike; turning it back on restores opt-ins, selections, credentials and prices.

       Turning Bedrock off needs --fallback-model <route>/<model> on a provider that is on (no_fallback_model otherwise); --select-for-all-teams selects that route for every team.

       --dry-run sends dryRun=true: nothing is written and the response's impact lists the routes, teams and agents affected; with table or text output the impact is also printed on stderr.

       Required permission: admin:inference-providers

SYNOPSIS
       afs model-endpoints update-inference-provider [options]

OPTIONS
       --provider (string) (required)
       --dry-run (boolean)
          Return the impact without writing
       --enabled (boolean)
       --fallback-model (string)
          Bedrock off only: becomes models.default_model_id; its route must be enabled, on a provider that is on
       --reason (string)
          Recorded as disabledReason and in the audit record
       --select-for-all-teams (boolean)
          Bedrock off only: select the fallback route for every team that lacks it (selectedBy service_admin) in the same transaction
       --team-selection (string)
          Possible values: team_admin, service_admin
       --if-match (string)
          Apply only if the resource is at this version (safe update).

GLOBAL OPTIONS
       --profile                Use a specific profile from ~/.afs/config.
       --endpoint-url           Override the API endpoint URL.
       --output                 Output format: json, yaml, table, text or csv.
       --query                  JMESPath query applied to the result before rendering.
       --no-paginate            Return only the first page of results.
       --max-items              Stop after this many items and print NextToken.
       --page-size              Items per request (maxResults), 1 to 100.
       --starting-token         Resume from a NextToken printed by --max-items.
       --cli-input-json         Read parameters from a JSON document or file://path.
       --generate-cli-skeleton  Print a JSON skeleton of the input and make no call.
       --debug                  Write a redacted request log to stderr.
       --no-cli-pager           Do not send output through a pager.
       --color                  Colored output: on, off or auto.
       --cli-connect-timeout    Connect timeout in seconds (0 for none; default 30).
       --cli-read-timeout       Read timeout in seconds (0 for none; default 60).
       --version                Print the afs version and exit.

EXAMPLES
       See what turning Bedrock off would change, without writing

       afs model-endpoints update-inference-provider --provider bedrock --no-enabled --fallback-model azure-openai/gpt-5.1 --select-for-all-teams --reason "No AWS Bedrock in this enterprise" --dry-run

       Turn Bedrock off and move agents to an Azure route

       afs model-endpoints update-inference-provider --provider bedrock --no-enabled --fallback-model azure-openai/gpt-5.1 --select-for-all-teams --reason "No AWS Bedrock in this enterprise"

       Turn Azure on, selected per team by service admins

       afs model-endpoints update-inference-provider --provider azure --enabled --team-selection service_admin

```

## update-model-endpoint

```text
AFS MODEL-ENDPOINTS UPDATE-MODEL-ENDPOINT

NAME
       model-endpoints update-model-endpoint

DESCRIPTION
       Update model endpoint

       Required permission: admin:model-endpoints

SYNOPSIS
       afs model-endpoints update-model-endpoint [options]

OPTIONS
       --endpoint-id (string) (required)
       --auth-header (string)
          Header name, required when authKind is header
       --auth-kind (string)
          Possible values: none, bearer, header, azure_federated, google_federated
       --base-url (string)
          https for external and corporate; in_namespace may be http://<service>.<namespace>.svc:<port>
       --ca-bundle (string or JSON)
          PEM bundle of the CA that signs the endpoint's certificate, as text or file://PATH (not a secret).
       --credential-mode (string)
          Possible values: shared, per_team
       --description (string)
       --display-name (string)
       --enabled (boolean)
          Refused while the provider is off (provider_disabled) or the base URL holds a placeholder (endpoint_base_url_placeholder)
       --external-system-id (string)
       --location (string)
          Possible values: external, corporate
       --max-body-mib (integer)
          Range: 1 to 100
       --pricing-mode (string)
          Possible values: per_token, capacity, free
       --protocols (list)
          Protocols the endpoint speaks: openai_chat always; openai_responses and openai_embeddings add those paths (space-separated).
       --provider (string)
       --requires-terms (boolean)
       --timeout-seconds (integer)
          Range: 5 to 3600
       --use-provider-cost (boolean)
       --if-match (string)
          Apply only if the resource is at this version (safe update).

GLOBAL OPTIONS
       --profile                Use a specific profile from ~/.afs/config.
       --endpoint-url           Override the API endpoint URL.
       --output                 Output format: json, yaml, table, text or csv.
       --query                  JMESPath query applied to the result before rendering.
       --no-paginate            Return only the first page of results.
       --max-items              Stop after this many items and print NextToken.
       --page-size              Items per request (maxResults), 1 to 100.
       --starting-token         Resume from a NextToken printed by --max-items.
       --cli-input-json         Read parameters from a JSON document or file://path.
       --generate-cli-skeleton  Print a JSON skeleton of the input and make no call.
       --debug                  Write a redacted request log to stderr.
       --no-cli-pager           Do not send output through a pager.
       --color                  Colored output: on, off or auto.
       --cli-connect-timeout    Connect timeout in seconds (0 for none; default 30).
       --cli-read-timeout       Read timeout in seconds (0 for none; default 60).
       --version                Print the afs version and exit.

EXAMPLES
       Switch an endpoint on

       afs model-endpoints update-model-endpoint --endpoint-id mep_01JD2Q3R4S5T6V7W8X9Y0Z0001 --enabled

       Raise the timeout and require extra terms

       afs model-endpoints update-model-endpoint --endpoint-id mep_01JD2Q3R4S5T6V7W8X9Y0Z0001 --timeout-seconds 600 --requires-terms

```

## update-team-model-endpoint

```text
AFS MODEL-ENDPOINTS UPDATE-TEAM-MODEL-ENDPOINT

NAME
       model-endpoints update-team-model-endpoint

DESCRIPTION
       Update team model endpoint

       Required permission: model-endpoints:write

SYNOPSIS
       afs model-endpoints update-team-model-endpoint [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --endpoint-id (string) (required)
       --auth-header (string)
          Header name, required when authKind is header
       --auth-kind (string)
          Possible values: none, bearer, header, azure_federated, google_federated
       --base-url (string)
          https for external and corporate; in_namespace may be http://<service>.<namespace>.svc:<port>
       --ca-bundle (string or JSON)
          PEM bundle of the CA that signs the endpoint's certificate, as text or file://PATH (not a secret).
       --description (string)
       --display-name (string)
       --enabled (boolean)
       --location (string)
          Possible values: external, corporate, in_namespace
       --max-body-mib (integer)
          Range: 1 to 100
       --metrics-port (integer)
          Range: 1 to 65535
       --owner (string)
       --pricing-mode (string)
          Possible values: per_token, capacity, free
       --protocols (list)
          Protocols the endpoint speaks: openai_chat always; openai_responses and openai_embeddings add those paths (space-separated).
       --purpose (string)
       --timeout-seconds (integer)
          Range: 5 to 3600
       --if-match (string)
          Apply only if the resource is at this version (safe update).

GLOBAL OPTIONS
       --profile                Use a specific profile from ~/.afs/config.
       --endpoint-url           Override the API endpoint URL.
       --output                 Output format: json, yaml, table, text or csv.
       --query                  JMESPath query applied to the result before rendering.
       --no-paginate            Return only the first page of results.
       --max-items              Stop after this many items and print NextToken.
       --page-size              Items per request (maxResults), 1 to 100.
       --starting-token         Resume from a NextToken printed by --max-items.
       --cli-input-json         Read parameters from a JSON document or file://path.
       --generate-cli-skeleton  Print a JSON skeleton of the input and make no call.
       --debug                  Write a redacted request log to stderr.
       --no-cli-pager           Do not send output through a pager.
       --color                  Colored output: on, off or auto.
       --cli-connect-timeout    Connect timeout in seconds (0 for none; default 30).
       --cli-read-timeout       Read timeout in seconds (0 for none; default 60).
       --version                Print the afs version and exit.

EXAMPLES
       Switch a team endpoint on

       afs model-endpoints update-team-model-endpoint --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --endpoint-id mep_01JD2Q3R4S5T6V7W8X9Y0Z0002 --enabled

       Scrape the server's metrics port

       afs model-endpoints update-team-model-endpoint --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --endpoint-id mep_01JD2Q3R4S5T6V7W8X9Y0Z0002 --metrics-port 8000

```
