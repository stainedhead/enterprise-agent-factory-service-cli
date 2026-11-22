# afs api-keys

Generated from `afs api-keys help`; do not edit by hand.

```text
AFS API-KEYS

DESCRIPTION
       Create and revoke API keys.

SYNOPSIS
       afs api-keys <operation> [parameters]

OPERATIONS
       create-api-key                 Create API key
       delete-api-key                 Delete API key
       list-api-keys                  List API keys

To see help text for an operation, run: afs api-keys <operation> help
```

## create-api-key

```text
AFS API-KEYS CREATE-API-KEY

NAME
       api-keys create-api-key

DESCRIPTION
       Create API key

       Required permission: api-keys:write

SYNOPSIS
       afs api-keys create-api-key [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --expires-at (string) (required)
          At most 90 days ahead
       --name (string) (required)
       --role (string) (required)
          Possible values: team_admin, team_viewer
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
       Create a team API key (the secret is shown once)

       afs api-keys create-api-key --team-id team_01J9ZQ3ABCDEFGHJKMNPQRSTVW --name ci-pipeline --role team_admin --expires-at 2027-02-19T18:00:00Z

```

## delete-api-key

```text
AFS API-KEYS DELETE-API-KEY

NAME
       api-keys delete-api-key

DESCRIPTION
       Delete API key

       Required permission: api-keys:write

SYNOPSIS
       afs api-keys delete-api-key [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
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
       Delete an API key

       afs api-keys delete-api-key --team-id team_01J9ZQ3ABCDEFGHJKMNPQRSTVW --key-id key_01J9ZQ3ABCDEFGHJKMNPQRSTVW

```

## list-api-keys

```text
AFS API-KEYS LIST-API-KEYS

NAME
       api-keys list-api-keys

DESCRIPTION
       List API keys

       Required permission: api-keys:read

SYNOPSIS
       afs api-keys list-api-keys [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only

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
       List a team's API keys

       afs api-keys list-api-keys --team-id team_01J9ZQ3ABCDEFGHJKMNPQRSTVW

```
