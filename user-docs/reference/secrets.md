# afs secrets

Generated from `afs secrets help`; do not edit by hand.

```text
AFS SECRETS

DESCRIPTION
       Manage secret references and write secret values.

SYNOPSIS
       afs secrets <operation> [parameters]

OPERATIONS
       create-secret                  Create secret
       delete-secret                  Delete secret
       describe-secret                Get secret
       list-secrets                   List secrets
       put-secret-value               Put secret value
       rotate-secret                  Rotate secret
       update-secret                  Update secret

To see help text for an operation, run: afs secrets <operation> help
```

## create-secret

```text
AFS SECRETS CREATE-SECRET

NAME
       secrets create-secret

DESCRIPTION
       Create secret

       Required permission: secrets:write

SYNOPSIS
       afs secrets create-secret [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --agent-id (string)
       --description (string)
       --external-system-id (string)
       --name (string) (required)
       --pool-id (string)
       --rotation (JSON)
       --scope (string) (required)
          Possible values: team, agent, pool
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
       Create a secret from a file

       afs secrets create-secret --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --name github-pat --scope team --secret-value file://token.txt

       Create a secret from standard input (printf '%s' "$TOKEN" | ...)

       afs secrets create-secret --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --name github-pat --scope team --secret-value -

       Create a secret and type the value at a hidden prompt

       afs secrets create-secret --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --name github-pat --scope team

```

## delete-secret

```text
AFS SECRETS DELETE-SECRET

NAME
       secrets delete-secret

DESCRIPTION
       Delete secret

       Required permission: secrets:write

SYNOPSIS
       afs secrets delete-secret [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --secret-id (string) (required)
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
       Delete a secret that nothing uses

       afs secrets delete-secret --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --secret-id sec_01JD2Q3R4S5T6V7W8X9Y0Z0001

```

## describe-secret

```text
AFS SECRETS DESCRIBE-SECRET

NAME
       secrets describe-secret

DESCRIPTION
       Get secret

       Required permission: secrets:read-metadata

SYNOPSIS
       afs secrets describe-secret [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --secret-id (string) (required)

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
       Show one secret's metadata

       afs secrets describe-secret --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --secret-id sec_01JD2Q3R4S5T6V7W8X9Y0Z0001

```

## list-secrets

```text
AFS SECRETS LIST-SECRETS

NAME
       secrets list-secrets

DESCRIPTION
       List secrets

       Required permission: secrets:read-metadata

SYNOPSIS
       afs secrets list-secrets [options]

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
       List a team's secrets (names and metadata, never values)

       afs secrets list-secrets --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001

```

## put-secret-value

```text
AFS SECRETS PUT-SECRET-VALUE

NAME
       secrets put-secret-value

DESCRIPTION
       Put secret value

       Required permission: secrets:write

SYNOPSIS
       afs secrets put-secret-value [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --secret-id (string) (required)
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
       Replace a secret's value from a file (prints nothing)

       afs secrets put-secret-value --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --secret-id sec_01JD2Q3R4S5T6V7W8X9Y0Z0001 --secret-value file://token.txt

```

## rotate-secret

```text
AFS SECRETS ROTATE-SECRET

NAME
       secrets rotate-secret

DESCRIPTION
       Rotate secret

       Required permission: secrets:write

SYNOPSIS
       afs secrets rotate-secret [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --secret-id (string) (required)
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
       Rotate a secret to a new value

       afs secrets rotate-secret --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --secret-id sec_01JD2Q3R4S5T6V7W8X9Y0Z0001 --secret-value file://new-token.txt

```

## update-secret

```text
AFS SECRETS UPDATE-SECRET

NAME
       secrets update-secret

DESCRIPTION
       Update secret

       Required permission: secrets:write

SYNOPSIS
       afs secrets update-secret [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --secret-id (string) (required)
       --description (string)
       --external-system-id (string)
       --rotation (JSON)
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
       Change a secret's description

       afs secrets update-secret --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --secret-id sec_01JD2Q3R4S5T6V7W8X9Y0Z0001 --description "Token for the build bot"

```
