# afs pools

Generated from `afs pools help`; do not edit by hand.

```text
AFS POOLS

DESCRIPTION
       Manage worker pools.

SYNOPSIS
       afs pools <operation> [parameters]

OPERATIONS
       create-pool                    Create pool
       create-role-link               Create pool role link
       create-stage                   Create stage
       delete-config-value            Delete pool config value
       delete-pool                    Delete pool
       delete-role-link               Delete pool role link
       delete-stage                   Delete stage
       describe-pool                  Get pool
       get-pool-status                Get pool status
       list-config-values             List pool config values
       list-external-systems          List pool external systems
       list-handoffs                  List pool handoffs
       list-pools                     List pools
       list-role-links                List pool role links
       list-secret-bindings           List pool secret bindings
       list-stages                    List stages
       pause-pool                     Pause pool
       put-config-value               Put pool config value
       put-external-system            Put pool external system
       put-handoffs                   Put pool handoffs
       put-secret-bindings            Put pool secret bindings
       resume-pool                    Resume pool
       scale-pool                     Scale pool
       update-pool                    Update pool
       update-stage                   Update stage
       verify-role-link               Verify pool role link

To see help text for an operation, run: afs pools <operation> help
```

## create-pool

```text
AFS POOLS CREATE-POOL

NAME
       pools create-pool

DESCRIPTION
       Create pool

       Required permission: pools:write

SYNOPSIS
       afs pools create-pool [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --description (string)
       --image (string) (required)
          Container image pinned by digest
       --model-id (string) (required)
       --name (string) (required)
          Immutable name (team slug, agent, pool, secret)
       --persona (string)
       --resources (JSON)
       --scratch-size-gi-b (integer)
          Range: 1 to 500
       --shared-volume (string)
          Possible values: none, efs
       --stage-id (string)
       --template (string) (required)
       --workers (integer) (required)
          Desired workers
          Range: 0 to 200
       --idempotency-key (string)
          Idempotency key (default: a new ULID per invocation).
       --wait (boolean)
          Wait until the resource is ready (or gone, for a delete).
       --wait-timeout (duration)
          How long --wait waits, for example 90s or 10m (default 10m).

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
       afs pools create-pool help
```

## create-role-link

```text
AFS POOLS CREATE-ROLE-LINK

NAME
       pools create-role-link

DESCRIPTION
       Create pool role link

       Required permission: pools:write

SYNOPSIS
       afs pools create-role-link [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --pool-id (string) (required)
       --account-id (string) (required)
       --max-session-seconds (integer)
          Range: 900 to 43200
       --profile-name (string) (required)
       --role-name (string) (required)
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
       afs pools create-role-link help
```

## create-stage

```text
AFS POOLS CREATE-STAGE

NAME
       pools create-stage

DESCRIPTION
       Create stage

       Required permission: pools:write

SYNOPSIS
       afs pools create-stage [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --description (string)
       --name (string) (required)
          Immutable name (team slug, agent, pool, secret)
       --position (integer) (required)
          Range: 0 to 100
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
       afs pools create-stage help
```

## delete-config-value

```text
AFS POOLS DELETE-CONFIG-VALUE

NAME
       pools delete-config-value

DESCRIPTION
       Delete pool config value

       Required permission: pools:write

SYNOPSIS
       afs pools delete-config-value [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --pool-id (string) (required)
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
       afs pools delete-config-value help
```

## delete-pool

```text
AFS POOLS DELETE-POOL

NAME
       pools delete-pool

DESCRIPTION
       Delete pool

       Required permission: pools:write

SYNOPSIS
       afs pools delete-pool [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --pool-id (string) (required)
       --if-match (string)
          Apply only if the resource is at this version (safe update).
       --wait (boolean)
          Wait until the resource is ready (or gone, for a delete).
       --wait-timeout (duration)
          How long --wait waits, for example 90s or 10m (default 10m).

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
       afs pools delete-pool help
```

## delete-role-link

```text
AFS POOLS DELETE-ROLE-LINK

NAME
       pools delete-role-link

DESCRIPTION
       Delete pool role link

       Required permission: pools:write

SYNOPSIS
       afs pools delete-role-link [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --pool-id (string) (required)
       --link-id (string) (required)
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
       afs pools delete-role-link help
```

## delete-stage

```text
AFS POOLS DELETE-STAGE

NAME
       pools delete-stage

DESCRIPTION
       Delete stage

       Required permission: pools:write

SYNOPSIS
       afs pools delete-stage [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --stage-id (string) (required)
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
       afs pools delete-stage help
```

## describe-pool

```text
AFS POOLS DESCRIBE-POOL

NAME
       pools describe-pool

DESCRIPTION
       Get pool

       Required permission: pools:read

SYNOPSIS
       afs pools describe-pool [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --pool-id (string) (required)

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
       afs pools describe-pool help
```

## get-pool-status

```text
AFS POOLS GET-POOL-STATUS

NAME
       pools get-pool-status

DESCRIPTION
       Get pool status

       Required permission: pools:read

SYNOPSIS
       afs pools get-pool-status [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --pool-id (string) (required)

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
       afs pools get-pool-status help
```

## list-config-values

```text
AFS POOLS LIST-CONFIG-VALUES

NAME
       pools list-config-values

DESCRIPTION
       List pool config values

       Required permission: pools:read

SYNOPSIS
       afs pools list-config-values [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --pool-id (string) (required)

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
       afs pools list-config-values help
```

## list-external-systems

```text
AFS POOLS LIST-EXTERNAL-SYSTEMS

NAME
       pools list-external-systems

DESCRIPTION
       List pool external systems

       Required permission: pools:read

SYNOPSIS
       afs pools list-external-systems [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --pool-id (string) (required)

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
       afs pools list-external-systems help
```

## list-handoffs

```text
AFS POOLS LIST-HANDOFFS

NAME
       pools list-handoffs

DESCRIPTION
       List pool handoffs

       Required permission: pools:read

SYNOPSIS
       afs pools list-handoffs [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --pool-id (string) (required)

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
       afs pools list-handoffs help
```

## list-pools

```text
AFS POOLS LIST-POOLS

NAME
       pools list-pools

DESCRIPTION
       List pools

       Required permission: pools:read

SYNOPSIS
       afs pools list-pools [options]

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
       afs pools list-pools help
```

## list-role-links

```text
AFS POOLS LIST-ROLE-LINKS

NAME
       pools list-role-links

DESCRIPTION
       List pool role links

       Required permission: pools:read

SYNOPSIS
       afs pools list-role-links [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --pool-id (string) (required)

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
       afs pools list-role-links help
```

## list-secret-bindings

```text
AFS POOLS LIST-SECRET-BINDINGS

NAME
       pools list-secret-bindings

DESCRIPTION
       List pool secret bindings

       Required permission: pools:read

SYNOPSIS
       afs pools list-secret-bindings [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --pool-id (string) (required)

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
       afs pools list-secret-bindings help
```

## list-stages

```text
AFS POOLS LIST-STAGES

NAME
       pools list-stages

DESCRIPTION
       List stages

       Required permission: pools:read

SYNOPSIS
       afs pools list-stages [options]

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
       afs pools list-stages help
```

## pause-pool

```text
AFS POOLS PAUSE-POOL

NAME
       pools pause-pool

DESCRIPTION
       Pause pool

       Required permission: pools:write

SYNOPSIS
       afs pools pause-pool [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --pool-id (string) (required)
       --idempotency-key (string)
          Idempotency key (default: a new ULID per invocation).
       --wait (boolean)
          Wait until the resource is ready (or gone, for a delete).
       --wait-timeout (duration)
          How long --wait waits, for example 90s or 10m (default 10m).

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
       afs pools pause-pool help
```

## put-config-value

```text
AFS POOLS PUT-CONFIG-VALUE

NAME
       pools put-config-value

DESCRIPTION
       Put pool config value

       Required permission: pools:write

SYNOPSIS
       afs pools put-config-value [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --pool-id (string) (required)
       --name (string) (required)
       --source (string) (required)
          Possible values: literal, namespace
       --value (string)
          Required when source is literal
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
       afs pools put-config-value help
```

## put-external-system

```text
AFS POOLS PUT-EXTERNAL-SYSTEM

NAME
       pools put-external-system

DESCRIPTION
       Put pool external system

       Required permission: pools:write

SYNOPSIS
       afs pools put-external-system [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --pool-id (string) (required)
       --system-id (string) (required)
       --enabled (boolean) (required)
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
       afs pools put-external-system help
```

## put-handoffs

```text
AFS POOLS PUT-HANDOFFS

NAME
       pools put-handoffs

DESCRIPTION
       Put pool handoffs

       Required permission: pools:write

SYNOPSIS
       afs pools put-handoffs [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --pool-id (string) (required)
       --to-pool-ids (list) (required)
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
       afs pools put-handoffs help
```

## put-secret-bindings

```text
AFS POOLS PUT-SECRET-BINDINGS

NAME
       pools put-secret-bindings

DESCRIPTION
       Put pool secret bindings

       Required permission: pools:write

SYNOPSIS
       afs pools put-secret-bindings [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --pool-id (string) (required)
       --items (JSON) (required)
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
       afs pools put-secret-bindings help
```

## resume-pool

```text
AFS POOLS RESUME-POOL

NAME
       pools resume-pool

DESCRIPTION
       Resume pool

       Required permission: pools:write

SYNOPSIS
       afs pools resume-pool [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --pool-id (string) (required)
       --idempotency-key (string)
          Idempotency key (default: a new ULID per invocation).
       --wait (boolean)
          Wait until the resource is ready (or gone, for a delete).
       --wait-timeout (duration)
          How long --wait waits, for example 90s or 10m (default 10m).

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
       afs pools resume-pool help
```

## scale-pool

```text
AFS POOLS SCALE-POOL

NAME
       pools scale-pool

DESCRIPTION
       Scale pool

       Required permission: pools:write

SYNOPSIS
       afs pools scale-pool [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --pool-id (string) (required)
       --workers (integer) (required)
          Range: 0 to 200
       --idempotency-key (string)
          Idempotency key (default: a new ULID per invocation).
       --wait (boolean)
          Wait until the resource is ready (or gone, for a delete).
       --wait-timeout (duration)
          How long --wait waits, for example 90s or 10m (default 10m).

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
       afs pools scale-pool help
```

## update-pool

```text
AFS POOLS UPDATE-POOL

NAME
       pools update-pool

DESCRIPTION
       Update pool

       Required permission: pools:write

SYNOPSIS
       afs pools update-pool [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --pool-id (string) (required)
       --description (string)
       --image (string)
          Container image pinned by digest
       --model-id (string)
       --persona (string)
       --resources (JSON)
       --scratch-size-gi-b (integer)
          Range: 1 to 500
       --shared-volume (string)
          Possible values: none, efs
       --stage-id (string)
       --template (string)
       --workers (integer)
          Desired workers
          Range: 0 to 200
       --if-match (string)
          Apply only if the resource is at this version (safe update).
       --wait (boolean)
          Wait until the resource is ready (or gone, for a delete).
       --wait-timeout (duration)
          How long --wait waits, for example 90s or 10m (default 10m).

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
       afs pools update-pool help
```

## update-stage

```text
AFS POOLS UPDATE-STAGE

NAME
       pools update-stage

DESCRIPTION
       Update stage

       Required permission: pools:write

SYNOPSIS
       afs pools update-stage [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --stage-id (string) (required)
       --description (string)
       --position (integer)
          Range: 0 to 100
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
       afs pools update-stage help
```

## verify-role-link

```text
AFS POOLS VERIFY-ROLE-LINK

NAME
       pools verify-role-link

DESCRIPTION
       Verify pool role link

       Required permission: pools:write

SYNOPSIS
       afs pools verify-role-link [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --pool-id (string) (required)
       --link-id (string) (required)
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
       afs pools verify-role-link help
```
