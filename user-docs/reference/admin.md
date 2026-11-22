# afs admin

Generated from `afs admin help`; do not edit by hand.

```text
AFS ADMIN

DESCRIPTION
       Service administration: settings, templates, tags, channels, storage and security operators.

SYNOPSIS
       afs admin <operation> [parameters]

OPERATIONS
       add-security-operator          Add security operator
       approve-price-proposal         Approve model price proposal
       create-api-key                 Create service API key
       create-export                  Create export
       create-external-system         Create catalog system
       create-identity-provider       Create identity provider
       create-import                  Create import
       create-model                   Create model
       create-model-price             Create model price
       create-notification-channel    Create service notification channel
       delete-agent-template          Delete agent template
       delete-api-key                 Delete service API key
       delete-external-system         Delete catalog system
       delete-identity-provider       Delete identity provider
       delete-notification-channel    Delete service notification channel
       delete-service-tag             Delete service tag
       delete-storage-class           Delete storage class
       describe-export                Get export
       describe-external-system       Get catalog system
       describe-identity-provider     Get identity provider
       get-settings                   Get settings
       list-agent-templates           List admin agent templates
       list-api-keys                  List service API keys
       list-collaboration-endpoints   List collaboration endpoints
       list-cost-reconciliations      List cost reconciliations
       list-cost-tag-status           List cost tag status
       list-exports                   List exports
       list-external-systems          List catalog systems
       list-identity-providers        List identity providers
       list-model-prices              List model prices
       list-models                    List admin models
       list-notification-channels     List service notification channels
       list-org-accounts              List org accounts
       list-price-proposals           List model price proposals
       list-security-operators        List security operators
       list-service-tags              List service tags
       list-storage-classes           List storage classes admin
       put-agent-template             Put agent template
       put-service-tag                Put service tag
       put-setting                    Put setting
       put-storage-class              Put storage class
       reject-price-proposal          Reject model price proposal
       remove-security-operator       Remove security operator
       update-external-system         Update catalog system
       update-identity-provider       Update identity provider
       update-model                   Update model
       update-notification-channel    Update service notification channel

To see help text for an operation, run: afs admin <operation> help
```

## add-security-operator

```text
AFS ADMIN ADD-SECURITY-OPERATOR

NAME
       admin add-security-operator

DESCRIPTION
       Add security operator

       Required permission: admin:settings

SYNOPSIS
       afs admin add-security-operator [options]

OPTIONS
       --email (string)
       --group-name (string)
       --identity-provider-id (string)
       --subject-type (string) (required)
          Possible values: user, group
       --user-id (string)
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
       Make a user a security operator

       afs admin add-security-operator --subject-type user --email ana@example.com --identity-provider-id idp_01J9ZQ3ABCDEFGHJKMNPQRSTVW

       Make an identity provider group security operators

       afs admin add-security-operator --subject-type group --group-name security-ops --identity-provider-id idp_01J9ZQ3ABCDEFGHJKMNPQRSTVW

```

## approve-price-proposal

```text
AFS ADMIN APPROVE-PRICE-PROPOSAL

NAME
       admin approve-price-proposal

DESCRIPTION
       Approve model price proposal

       Required permission: admin:pricing

SYNOPSIS
       afs admin approve-price-proposal [options]

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

       afs admin approve-price-proposal --proposal-id mpp_01JD2Q3R4S5T6V7W8X9Y0Z0042 --note "Matches the AWS price list"

```

## create-api-key

```text
AFS ADMIN CREATE-API-KEY

NAME
       admin create-api-key

DESCRIPTION
       Create service API key

       Required permission: admin:settings

SYNOPSIS
       afs admin create-api-key [options]

OPTIONS
       --expires-at (string) (required)
          At most 90 days ahead
       --name (string) (required)
       --role (string) (required)
          Possible values: security_operator
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
       Create a security operator key (the secret is shown once)

       afs admin create-api-key --name cve-pipeline --role security_operator --expires-at 2027-02-19T18:00:00Z

```

## create-export

```text
AFS ADMIN CREATE-EXPORT

NAME
       admin create-export

DESCRIPTION
       Create export

       Required permission: admin:export

SYNOPSIS
       afs admin create-export [options]

OPTIONS
       --team-ids (list)
          Teams to export; omit for every team
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
       Export two teams

       afs admin create-export --team-ids team_01J9ZQ3ABCDEFGHJKMNPQRSTVW team_01J9ZQ3ABCDEFGHJKMNPQRSTVX

       Export every team

       afs admin create-export

```

## create-external-system

```text
AFS ADMIN CREATE-EXTERNAL-SYSTEM

NAME
       admin create-external-system

DESCRIPTION
       Create catalog system

       Required permission: admin:external-systems

SYNOPSIS
       afs admin create-external-system [options]

OPTIONS
       --category (string) (required)
          Possible values: dev_tools, collaboration, email, cloud, identity, llm, other
       --description (string)
       --docs-url (string)
       --enabled (boolean)
       --endpoints (list) (required)
       --name (string) (required)
       --slug (string) (required)
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
       Add a system to the catalog

       afs admin create-external-system --slug jira --name Jira --category collaboration --endpoints example.atlassian.net:443 --docs-url https://support.atlassian.com/jira

```

## create-identity-provider

```text
AFS ADMIN CREATE-IDENTITY-PROVIDER

NAME
       admin create-identity-provider

DESCRIPTION
       Create identity provider

       Required permission: admin:identity-providers

SYNOPSIS
       afs admin create-identity-provider [options]

OPTIONS
       --cli-client-id (string) (required)
       --client-secret-arn (string)
       --collab-client-id (string)
          The collaboration (Paperclip edge) client of this provider
       --collab-client-secret-arn (string)
          Reference only; needs collabClientId; never returned by ListSignInProviders
       --console-client-id (string) (required)
       --enabled (boolean)
       --groups-claim (string)
       --is-default (boolean)
       --issuer-url (string) (required)
       --kind (string) (required)
          Possible values: okta, entra_id
       --name (string) (required)
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
       Register an OIDC provider

       afs admin create-identity-provider --name "Corp SSO" --kind okta --issuer-url https://idp.example.com --console-client-id console-app --cli-client-id afs-cli --client-secret-arn arn:aws:secretsmanager:eu-west-1:111122223333:secret:afs/idp-AbCdEf

```

## create-import

```text
AFS ADMIN CREATE-IMPORT

NAME
       admin create-import

DESCRIPTION
       Create import

       Required permission: admin:export

SYNOPSIS
       afs admin create-import [options]

OPTIONS
       --dry-run (boolean)
          Default true: nothing is written
       --account-map (JSON)
          Source AWS account id → target account id
       --bundle (JSON) (required)
          An export bundle: {schemaVersion, sha256, teams[...]}
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
       Dry-run an import (nothing is changed)

       afs admin create-import --bundle file://bundle.json --account-map file://accounts.json

       Apply the import

       afs admin create-import --bundle file://bundle.json --account-map file://accounts.json --no-dry-run

```

## create-model

```text
AFS ADMIN CREATE-MODEL

NAME
       admin create-model

DESCRIPTION
       Create model

       Required permission: admin:models

SYNOPSIS
       afs admin create-model [options]

OPTIONS
       --context-window (integer) (required)
          Range: 1 to unbounded
       --display-name (string) (required)
       --enabled (boolean)
       --id (string) (required)
       --max-output-tokens (integer) (required)
          Required (D28)
          Range: 1 to unbounded
       --provider (string) (required)
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
       Add a model

       afs admin create-model --id anthropic.claude-sonnet-4-5 --provider anthropic --display-name "Claude Sonnet 4.5" --context-window 200000 --max-output-tokens 64000

```

## create-model-price

```text
AFS ADMIN CREATE-MODEL-PRICE

NAME
       admin create-model-price

DESCRIPTION
       Create model price

       Required permission: admin:pricing

SYNOPSIS
       afs admin create-model-price [options]

OPTIONS
       --model-id (string) (required)
       --cache-read-usd-per-mtok (string)
          USD amount as a decimal string with at most 6 fraction digits; never a JSON number
       --cache-write-usd-per-mtok (string)
          USD amount as a decimal string with at most 6 fraction digits; never a JSON number
       --effective-from (string) (required)
          RFC 3339 UTC with Z
       --input-usd-per-mtok (string) (required)
          USD amount as a decimal string with at most 6 fraction digits; never a JSON number
       --output-usd-per-mtok (string) (required)
          USD amount as a decimal string with at most 6 fraction digits; never a JSON number
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
       Record a new price (USD per million tokens, as strings)

       afs admin create-model-price --model-id anthropic.claude-sonnet-4-5 --effective-from 2026-12-01T00:00:00Z --input-usd-per-mtok 3.00 --output-usd-per-mtok 15.00 --cache-read-usd-per-mtok 0.30 --cache-write-usd-per-mtok 3.75

```

## create-notification-channel

```text
AFS ADMIN CREATE-NOTIFICATION-CHANNEL

NAME
       admin create-notification-channel

DESCRIPTION
       Create service notification channel

       Required permission: admin:settings

SYNOPSIS
       afs admin create-notification-channel [options]

OPTIONS
       --config (JSON) (required)
       --enabled (boolean)
       --event-types (list) (required)
       --kind (string) (required)
          Possible values: slack, email, webhook, jira, eventbridge
       --name (string) (required)
       --secret-id (string)
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
       Create an email channel for security events

       afs admin create-notification-channel --name security-email --kind email --config '{"to":["security@example.com"]}' --event-types 'work_order.*'

```

## delete-agent-template

```text
AFS ADMIN DELETE-AGENT-TEMPLATE

NAME
       admin delete-agent-template

DESCRIPTION
       Delete agent template

       Required permission: admin:agent-templates

SYNOPSIS
       afs admin delete-agent-template [options]

OPTIONS
       --name (string) (required)
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
       Delete a template no agent uses

       afs admin delete-agent-template --name claude-code

```

## delete-api-key

```text
AFS ADMIN DELETE-API-KEY

NAME
       admin delete-api-key

DESCRIPTION
       Delete service API key

       Required permission: admin:settings

SYNOPSIS
       afs admin delete-api-key [options]

OPTIONS
       --key-id (string) (required)
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
       Delete a service-level API key

       afs admin delete-api-key --key-id key_01J9ZQ3ABCDEFGHJKMNPQRSTVW

```

## delete-external-system

```text
AFS ADMIN DELETE-EXTERNAL-SYSTEM

NAME
       admin delete-external-system

DESCRIPTION
       Delete catalog system

       Required permission: admin:external-systems

SYNOPSIS
       afs admin delete-external-system [options]

OPTIONS
       --system-id (string) (required)
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
       Delete a catalog system (built-in ones are refused with 409)

       afs admin delete-external-system --system-id ext_01JD2Q3R4S5T6V7W8X9Y0Z0001

```

## delete-identity-provider

```text
AFS ADMIN DELETE-IDENTITY-PROVIDER

NAME
       admin delete-identity-provider

DESCRIPTION
       Delete identity provider

       Required permission: admin:identity-providers

SYNOPSIS
       afs admin delete-identity-provider [options]

OPTIONS
       --provider-id (string) (required)
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
       Delete a provider that no team uses

       afs admin delete-identity-provider --provider-id idp_01J9ZQ3ABCDEFGHJKMNPQRSTVW

```

## delete-notification-channel

```text
AFS ADMIN DELETE-NOTIFICATION-CHANNEL

NAME
       admin delete-notification-channel

DESCRIPTION
       Delete service notification channel

       Required permission: admin:settings

SYNOPSIS
       afs admin delete-notification-channel [options]

OPTIONS
       --channel-id (string) (required)
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
       Delete a channel (refused while a setting names it)

       afs admin delete-notification-channel --channel-id nch_01JD2Q3R4S5T6V7W8X9Y0Z0001

```

## delete-service-tag

```text
AFS ADMIN DELETE-SERVICE-TAG

NAME
       admin delete-service-tag

DESCRIPTION
       Delete service tag

       Required permission: admin:service-tags

SYNOPSIS
       afs admin delete-service-tag [options]

OPTIONS
       --key (string) (required)
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
       Delete a service tag

       afs admin delete-service-tag --key cost-center

```

## delete-storage-class

```text
AFS ADMIN DELETE-STORAGE-CLASS

NAME
       admin delete-storage-class

DESCRIPTION
       Delete storage class

       Required permission: admin:storage-classes

SYNOPSIS
       afs admin delete-storage-class [options]

OPTIONS
       --name (string) (required)
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
       Delete a storage class no volume uses

       afs admin delete-storage-class --name gp3-standard

```

## describe-export

```text
AFS ADMIN DESCRIBE-EXPORT

NAME
       admin describe-export

DESCRIPTION
       Get export

       Required permission: admin:export

SYNOPSIS
       afs admin describe-export [options]

OPTIONS
       --export-id (string) (required)

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
       Show an export (sha256, size and download URL when it has finished)

       afs admin describe-export --export-id exp_01J9ZQ3ABCDEFGHJKMNPQRSTVW

```

## describe-external-system

```text
AFS ADMIN DESCRIBE-EXTERNAL-SYSTEM

NAME
       admin describe-external-system

DESCRIPTION
       Get catalog system

       Required permission: admin:external-systems

SYNOPSIS
       afs admin describe-external-system [options]

OPTIONS
       --system-id (string) (required)

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
       Show one catalog system

       afs admin describe-external-system --system-id ext_01JD2Q3R4S5T6V7W8X9Y0Z0001

```

## describe-identity-provider

```text
AFS ADMIN DESCRIBE-IDENTITY-PROVIDER

NAME
       admin describe-identity-provider

DESCRIPTION
       Get identity provider

       Required permission: admin:identity-providers

SYNOPSIS
       afs admin describe-identity-provider [options]

OPTIONS
       --provider-id (string) (required)

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
       Show one identity provider (never a secret, only its ARN)

       afs admin describe-identity-provider --provider-id idp_01J9ZQ3ABCDEFGHJKMNPQRSTVW

```

## get-settings

```text
AFS ADMIN GET-SETTINGS

NAME
       admin get-settings

DESCRIPTION
       Get settings

       Required permission: admin:settings

SYNOPSIS
       afs admin get-settings [options]

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
       List every service setting

       afs admin get-settings

       Print the settings as JSON

       afs admin get-settings --output json

```

## list-agent-templates

```text
AFS ADMIN LIST-AGENT-TEMPLATES

NAME
       admin list-agent-templates

DESCRIPTION
       List admin agent templates

       Required permission: admin:agent-templates

SYNOPSIS
       afs admin list-agent-templates [options]

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
       List every agent template

       afs admin list-agent-templates

```

## list-api-keys

```text
AFS ADMIN LIST-API-KEYS

NAME
       admin list-api-keys

DESCRIPTION
       List service API keys

       Required permission: admin:settings

SYNOPSIS
       afs admin list-api-keys [options]

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
       List the service-level API keys

       afs admin list-api-keys

```

## list-collaboration-endpoints

```text
AFS ADMIN LIST-COLLABORATION-ENDPOINTS

NAME
       admin list-collaboration-endpoints

DESCRIPTION
       List collaboration endpoints

       Required permission: admin:settings

SYNOPSIS
       afs admin list-collaboration-endpoints [options]

OPTIONS
       --team-id (string)
       --state (string)
          Possible values: Ready, Degraded, Failed, Pending

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
       afs admin list-collaboration-endpoints help
```

## list-cost-reconciliations

```text
AFS ADMIN LIST-COST-RECONCILIATIONS

NAME
       admin list-cost-reconciliations

DESCRIPTION
       List cost reconciliations

       Required permission: admin:reports

SYNOPSIS
       afs admin list-cost-reconciliations [options]

OPTIONS
       --status (string)
          Possible values: ok, over_threshold, missing_billing

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
       Show the days where billing differs from the gateway

       afs admin list-cost-reconciliations --status over_threshold

```

## list-cost-tag-status

```text
AFS ADMIN LIST-COST-TAG-STATUS

NAME
       admin list-cost-tag-status

DESCRIPTION
       List cost tag status

       Required permission: admin:cost-tag-status

SYNOPSIS
       afs admin list-cost-tag-status [options]

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
       Show which cost-allocation tags AWS has activated and seen in the cost report

       afs admin list-cost-tag-status

```

## list-exports

```text
AFS ADMIN LIST-EXPORTS

NAME
       admin list-exports

DESCRIPTION
       List exports

       Required permission: admin:export

SYNOPSIS
       afs admin list-exports [options]

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
       List exports

       afs admin list-exports

```

## list-external-systems

```text
AFS ADMIN LIST-EXTERNAL-SYSTEMS

NAME
       admin list-external-systems

DESCRIPTION
       List catalog systems

       Required permission: admin:external-systems

SYNOPSIS
       afs admin list-external-systems [options]

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
       List the catalog of external systems

       afs admin list-external-systems

```

## list-identity-providers

```text
AFS ADMIN LIST-IDENTITY-PROVIDERS

NAME
       admin list-identity-providers

DESCRIPTION
       List identity providers

       Required permission: admin:identity-providers

SYNOPSIS
       afs admin list-identity-providers [options]

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
       List identity providers

       afs admin list-identity-providers

```

## list-model-prices

```text
AFS ADMIN LIST-MODEL-PRICES

NAME
       admin list-model-prices

DESCRIPTION
       List model prices

       Required permission: admin:pricing

SYNOPSIS
       afs admin list-model-prices [options]

OPTIONS
       --model-id (string) (required)

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
       Show a model's price history

       afs admin list-model-prices --model-id anthropic.claude-sonnet-4-5

```

## list-models

```text
AFS ADMIN LIST-MODELS

NAME
       admin list-models

DESCRIPTION
       List admin models

       Required permission: admin:models

SYNOPSIS
       afs admin list-models [options]

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
       List every model, enabled or not

       afs admin list-models

```

## list-notification-channels

```text
AFS ADMIN LIST-NOTIFICATION-CHANNELS

NAME
       admin list-notification-channels

DESCRIPTION
       List service notification channels

       Required permission: admin:settings

SYNOPSIS
       afs admin list-notification-channels [options]

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
       List the service-level notification channels

       afs admin list-notification-channels

```

## list-org-accounts

```text
AFS ADMIN LIST-ORG-ACCOUNTS

NAME
       admin list-org-accounts

DESCRIPTION
       List org accounts

       Required permission: admin:accounts

SYNOPSIS
       afs admin list-org-accounts [options]

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
       List the AWS Organization accounts the service knows

       afs admin list-org-accounts

```

## list-price-proposals

```text
AFS ADMIN LIST-PRICE-PROPOSALS

NAME
       admin list-price-proposals

DESCRIPTION
       List model price proposals

       Required permission: admin:pricing

SYNOPSIS
       afs admin list-price-proposals [options]

OPTIONS
       --status (string)
          Possible values: pending, approved, rejected, superseded

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
       List the pending price proposals

       afs admin list-price-proposals --status pending --output table

```

## list-security-operators

```text
AFS ADMIN LIST-SECURITY-OPERATORS

NAME
       admin list-security-operators

DESCRIPTION
       List security operators

       Required permission: admin:settings

SYNOPSIS
       afs admin list-security-operators [options]

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
       List the security operators

       afs admin list-security-operators

```

## list-service-tags

```text
AFS ADMIN LIST-SERVICE-TAGS

NAME
       admin list-service-tags

DESCRIPTION
       List service tags

       Required permission: admin:service-tags

SYNOPSIS
       afs admin list-service-tags [options]

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
       List the service tags

       afs admin list-service-tags

```

## list-storage-classes

```text
AFS ADMIN LIST-STORAGE-CLASSES

NAME
       admin list-storage-classes

DESCRIPTION
       List storage classes admin

       Required permission: admin:storage-classes

SYNOPSIS
       afs admin list-storage-classes [options]

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
       List every storage class

       afs admin list-storage-classes

```

## put-agent-template

```text
AFS ADMIN PUT-AGENT-TEMPLATE

NAME
       admin put-agent-template

DESCRIPTION
       Put agent template

       Required permission: admin:agent-templates

SYNOPSIS
       afs admin put-agent-template [options]

OPTIONS
       --name (string) (required)
       --allowed-digests (list) (required)
       --capabilities (list)
       --current-digest (string) (required)
       --default-resources (JSON)
       --description (string)
       --enabled (boolean)
       --harness (string) (required)
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
       Create or replace a template

       afs admin put-agent-template --name claude-code --harness claude-code --current-digest sha256:4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f --allowed-digests sha256:4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f --default-resources '{"cpu":"1","memory":"2Gi"}'

```

## put-service-tag

```text
AFS ADMIN PUT-SERVICE-TAG

NAME
       admin put-service-tag

DESCRIPTION
       Put service tag

       Required permission: admin:service-tags

SYNOPSIS
       afs admin put-service-tag [options]

OPTIONS
       --key (string) (required)
       --allowed-values (list)
       --default-value (string)
       --description (string)
       --mode (string) (required)
          Possible values: fixed, required
       --value (string)
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
       Require a cost-center tag with allowed values

       afs admin put-service-tag --key cost-center --mode required --allowed-values CC-1001 CC-2002 --description "Finance cost center"

       Apply a fixed tag to every team

       afs admin put-service-tag --key environment --mode fixed --value production

```

## put-setting

```text
AFS ADMIN PUT-SETTING

NAME
       admin put-setting

DESCRIPTION
       Put setting

       Required permission: admin:settings

SYNOPSIS
       afs admin put-setting [options]

OPTIONS
       --key (string) (required)
          Possible values: secrets.path_prefix, pricing.price_list_sync, billing.variance_alert_pct, api.idempotency_key_ttl_seconds, teams.ephemeral_max_days, work_orders.sla_warning_hours, billing.tag_activation_acknowledged, work_orders.notification_channel_id, reports.idle_cpu_millicores, egress.denied_spike_threshold, gateway.reservation_ttl_seconds, api.cors_origins
       --value (string or JSON) (required)
          JSON value of the setting's type
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
       Set a list value (JSON)

       afs admin put-setting --key api.cors_origins --value '["https://console.example.com"]'

       Set a number

       afs admin put-setting --key api.idempotency_key_ttl_seconds --value 86400

```

## put-storage-class

```text
AFS ADMIN PUT-STORAGE-CLASS

NAME
       admin put-storage-class

DESCRIPTION
       Put storage class

       Required permission: admin:storage-classes

SYNOPSIS
       afs admin put-storage-class [options]

OPTIONS
       --name (string) (required)
       --description (string)
       --enabled (boolean)
       --volume-type (string) (required)
          Possible values: gp3, io2
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
       Create or replace a storage class

       afs admin put-storage-class --name gp3-standard --volume-type gp3 --description "General purpose SSD"

```

## reject-price-proposal

```text
AFS ADMIN REJECT-PRICE-PROPOSAL

NAME
       admin reject-price-proposal

DESCRIPTION
       Reject model price proposal

       Required permission: admin:pricing

SYNOPSIS
       afs admin reject-price-proposal [options]

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

       afs admin reject-price-proposal --proposal-id mpp_01JD2Q3R4S5T6V7W8X9Y0Z0042 --note "Waiting for the contract rate"

```

## remove-security-operator

```text
AFS ADMIN REMOVE-SECURITY-OPERATOR

NAME
       admin remove-security-operator

DESCRIPTION
       Remove security operator

       Required permission: admin:settings

SYNOPSIS
       afs admin remove-security-operator [options]

OPTIONS
       --binding-id (string) (required)
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
       Remove a security operator

       afs admin remove-security-operator --binding-id rb_01J9ZQ3ABCDEFGHJKMNPQRSTVW

```

## update-external-system

```text
AFS ADMIN UPDATE-EXTERNAL-SYSTEM

NAME
       admin update-external-system

DESCRIPTION
       Update catalog system

       Required permission: admin:external-systems

SYNOPSIS
       afs admin update-external-system [options]

OPTIONS
       --system-id (string) (required)
       --description (string)
       --docs-url (string)
       --enabled (boolean)
          Service-level switch
       --endpoints (list)
       --name (string)
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
       Disable a catalog system

       afs admin update-external-system --system-id ext_01JD2Q3R4S5T6V7W8X9Y0Z0001 --no-enabled

```

## update-identity-provider

```text
AFS ADMIN UPDATE-IDENTITY-PROVIDER

NAME
       admin update-identity-provider

DESCRIPTION
       Update identity provider

       Required permission: admin:identity-providers

SYNOPSIS
       afs admin update-identity-provider [options]

OPTIONS
       --provider-id (string) (required)
       --cli-client-id (string)
       --client-secret-arn (string)
       --collab-client-id (string)
          The collaboration (Paperclip edge) client of this provider
       --collab-client-secret-arn (string)
          Reference only; needs collabClientId; never returned by ListSignInProviders
       --console-client-id (string)
       --enabled (boolean)
       --groups-claim (string)
       --is-default (boolean)
       --name (string)
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
       Disable a provider

       afs admin update-identity-provider --provider-id idp_01J9ZQ3ABCDEFGHJKMNPQRSTVW --no-enabled

```

## update-model

```text
AFS ADMIN UPDATE-MODEL

NAME
       admin update-model

DESCRIPTION
       Update model

       Required permission: admin:models

SYNOPSIS
       afs admin update-model [options]

OPTIONS
       --model-id (string) (required)
       --context-window (integer)
          Range: 1 to unbounded
       --display-name (string)
       --enabled (boolean)
       --max-output-tokens (integer)
          Range: 1 to unbounded
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
       Disable a model

       afs admin update-model --model-id anthropic.claude-sonnet-4-5 --no-enabled

```

## update-notification-channel

```text
AFS ADMIN UPDATE-NOTIFICATION-CHANNEL

NAME
       admin update-notification-channel

DESCRIPTION
       Update service notification channel

       Required permission: admin:settings

SYNOPSIS
       afs admin update-notification-channel [options]

OPTIONS
       --channel-id (string) (required)
       --config (JSON)
       --enabled (boolean)
       --event-types (list)
       --name (string)
       --secret-id (string)
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
       Disable a channel

       afs admin update-notification-channel --channel-id nch_01JD2Q3R4S5T6V7W8X9Y0Z0001 --no-enabled

```
