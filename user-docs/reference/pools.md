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
       delete-pool-federated-identity Delete pool federated identity
       delete-role-link               Delete pool role link
       delete-stage                   Delete stage
       describe-pool                  Get pool
       get-effective-pool-mcp         Get effective pool MCP
       get-pool-status                Get pool status
       list-config-values             List pool config values
       list-external-systems          List pool external systems
       list-handoffs                  List pool handoffs
       list-pool-federated-identities List pool federated identities
       list-pools                     List pools
       list-role-links                List pool role links
       list-secret-bindings           List pool secret bindings
       list-stages                    List stages
       pause-pool                     Pause pool
       put-config-value               Put pool config value
       put-external-system            Put pool external system
       put-handoffs                   Put pool handoffs
       put-pool-federated-identity    Put pool federated identity
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
       --mcp (JSON)
          v1.6.0 MCP access of an agent or pool
       --model-id (string)
          Deprecated: use --default-model-id (and --models to choose the set) instead.
       --models (list)
          Models the agent may call: the word all (every model the team may use, now and later) or one or more model ids (a fixed selection). Sent members replace the stored ones.
       --name (string) (required)
          Immutable name (team slug, agent, pool, secret)
       --persona (string)
       --resources (JSON)
       --scratch-size-gib (integer)
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
       --default-model-id (string)
          The model the runtime starts with; it must be one of the agent's models. Without --models on an update, only the default changes.
       --routes (list)
          Model routes the agent may call: the word all (every route of the team, now and later) or one or more route slugs (a fixed selection). Folded into models.routeAccess and models.routes.
       --default-route (string)
          The route of the default model; --default-model-id then names the route's model. Or give --default-model-id <route>/<model> alone.
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
       Create a pool of 10 workers and wait until they report

       afs pools create-pool --team-id team_01J9ZQ3ABCDEFGHJKMNPQRSTVW --name fixer --template claude-code --image registry.example.com/agents/coder@sha256:4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f --models anthropic.claude-sonnet-4-5 anthropic.claude-haiku-4-5 --default-model-id anthropic.claude-sonnet-4-5 --workers 10 --wait

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
       Link a role by name in an account

       afs pools create-role-link --team-id team_01J9ZQ3ABCDEFGHJKMNPQRSTVW --pool-id pool_01J9ZQ3ABCDEFGHJKMNPQRSTVW --account-id 111122223333 --role-name afs-cve-reader --profile-name target

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
          Range: 1 to 50
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
       Add a stage at position 1

       afs pools create-stage --team-id team_01J9ZQ3ABCDEFGHJKMNPQRSTVW --name triage --position 1 --description "Find affected repositories"

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
       Remove a configuration value

       afs pools delete-config-value --team-id team_01J9ZQ3ABCDEFGHJKMNPQRSTVW --pool-id pool_01J9ZQ3ABCDEFGHJKMNPQRSTVW --name LOG_LEVEL

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
       Delete a pool and wait until it is gone

       afs pools delete-pool --team-id team_01J9ZQ3ABCDEFGHJKMNPQRSTVW --pool-id pool_01J9ZQ3ABCDEFGHJKMNPQRSTVW --wait

```

## delete-pool-federated-identity

```text
AFS POOLS DELETE-POOL-FEDERATED-IDENTITY

NAME
       pools delete-pool-federated-identity

DESCRIPTION
       Delete pool federated identity

       Required permission: agent-federated-identities:write

SYNOPSIS
       afs pools delete-pool-federated-identity [options]

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
       Remove a federated identity

       afs pools delete-pool-federated-identity --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --pool-id pool_01JD2Q3R4S5T6V7W8X9Y0Z0001 --name vault

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
       Delete a role link

       afs pools delete-role-link --team-id team_01J9ZQ3ABCDEFGHJKMNPQRSTVW --pool-id pool_01J9ZQ3ABCDEFGHJKMNPQRSTVW --link-id prl_01J9ZQ3ABCDEFGHJKMNPQRSTVW

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
       Delete a stage

       afs pools delete-stage --team-id team_01J9ZQ3ABCDEFGHJKMNPQRSTVW --stage-id stg_01J9ZQ3ABCDEFGHJKMNPQRSTVW

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
       Show a pool (desired workers, state, phase)

       afs pools describe-pool --team-id team_01J9ZQ3ABCDEFGHJKMNPQRSTVW --pool-id pool_01J9ZQ3ABCDEFGHJKMNPQRSTVW

```

## get-effective-pool-mcp

```text
AFS POOLS GET-EFFECTIVE-POOL-MCP

NAME
       pools get-effective-pool-mcp

DESCRIPTION
       Get effective pool MCP

       The servers the pool reaches through the team gateway after the service, team and pool narrowing, and the AFS_MCP_* variable names it gets.

       Required permission: team-mcp-servers:read

SYNOPSIS
       afs pools get-effective-pool-mcp [options]

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
       Show the MCP servers it reaches

       afs pools get-effective-pool-mcp --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --pool-id pool_01JD2Q3R4S5T6V7W8X9Y0Z0001

       Show the variable names it gets

       afs pools get-effective-pool-mcp --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --pool-id pool_01JD2Q3R4S5T6V7W8X9Y0Z0001 --query env --output text

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
       Show desired, reporting, working, idle and impaired workers

       afs pools get-pool-status --team-id team_01J9ZQ3ABCDEFGHJKMNPQRSTVW --pool-id pool_01J9ZQ3ABCDEFGHJKMNPQRSTVW

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
       List a pool's configuration values

       afs pools list-config-values --team-id team_01J9ZQ3ABCDEFGHJKMNPQRSTVW --pool-id pool_01J9ZQ3ABCDEFGHJKMNPQRSTVW

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
       List the external systems a pool may reach

       afs pools list-external-systems --team-id team_01J9ZQ3ABCDEFGHJKMNPQRSTVW --pool-id pool_01J9ZQ3ABCDEFGHJKMNPQRSTVW

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
       Show where a pool hands work to

       afs pools list-handoffs --team-id team_01J9ZQ3ABCDEFGHJKMNPQRSTVW --pool-id pool_01J9ZQ3ABCDEFGHJKMNPQRSTVW

```

## list-pool-federated-identities

```text
AFS POOLS LIST-POOL-FEDERATED-IDENTITIES

NAME
       pools list-pool-federated-identities

DESCRIPTION
       List pool federated identities

       Required permission: pools:read

SYNOPSIS
       afs pools list-pool-federated-identities [options]

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
       List its federated identities

       afs pools list-pool-federated-identities --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --pool-id pool_01JD2Q3R4S5T6V7W8X9Y0Z0001

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
       List a team's pools

       afs pools list-pools --team-id team_01J9ZQ3ABCDEFGHJKMNPQRSTVW

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
       List a pool's role links

       afs pools list-role-links --team-id team_01J9ZQ3ABCDEFGHJKMNPQRSTVW --pool-id pool_01J9ZQ3ABCDEFGHJKMNPQRSTVW

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
       List a pool's secret bindings

       afs pools list-secret-bindings --team-id team_01J9ZQ3ABCDEFGHJKMNPQRSTVW --pool-id pool_01J9ZQ3ABCDEFGHJKMNPQRSTVW

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
       List a swarm team's stages

       afs pools list-stages --team-id team_01J9ZQ3ABCDEFGHJKMNPQRSTVW

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
       Pause a pool

       afs pools pause-pool --team-id team_01J9ZQ3ABCDEFGHJKMNPQRSTVW --pool-id pool_01J9ZQ3ABCDEFGHJKMNPQRSTVW

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
       Set a literal value

       afs pools put-config-value --team-id team_01J9ZQ3ABCDEFGHJKMNPQRSTVW --pool-id pool_01J9ZQ3ABCDEFGHJKMNPQRSTVW --name LOG_LEVEL --source literal --value debug

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
       Disable an external system for a pool

       afs pools put-external-system --team-id team_01J9ZQ3ABCDEFGHJKMNPQRSTVW --pool-id pool_01J9ZQ3ABCDEFGHJKMNPQRSTVW --system-id ext_01J9ZQ3ABCDEFGHJKMNPQRSTVW --no-enabled

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
       Hand off to two pools

       afs pools put-handoffs --team-id team_01J9ZQ3ABCDEFGHJKMNPQRSTVW --pool-id pool_01J9ZQ3ABCDEFGHJKMNPQRSTVW --to-pool-ids pool_01J9ZQ3ABCDEFGHJKMNPQRSTVX pool_01J9ZQ3ABCDEFGHJKMNPQRSTVY

```

## put-pool-federated-identity

```text
AFS POOLS PUT-POOL-FEDERATED-IDENTITY

NAME
       pools put-pool-federated-identity

DESCRIPTION
       Put pool federated identity

       A federated identity projects a Kubernetes service account token with the audience into the pod (filePath, envVar), for an external system that trusts the cluster's issuer. The response names the issuer and subject to trust; no secret is stored.

       Required permission: agent-federated-identities:write

SYNOPSIS
       afs pools put-pool-federated-identity [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --pool-id (string) (required)
       --name (string) (required)
       --audience (string) (required)
          Must match federation.allowed_audiences (400 audience_not_allowed); sts.amazonaws.com and afs-gateway are reserved
       --env-var (string)
       --expiration-seconds (integer)
          Range: 600 to 86400
       --file-path (string)
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
       Project a token for an external system's STS

       afs pools put-pool-federated-identity --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --pool-id pool_01JD2Q3R4S5T6V7W8X9Y0Z0001 --name vault --audience vault.corp.example --expiration-seconds 3600

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
       Bind a secret as a file (the value is never shown)

       afs pools put-secret-bindings --team-id team_01J9ZQ3ABCDEFGHJKMNPQRSTVW --pool-id pool_01J9ZQ3ABCDEFGHJKMNPQRSTVW --items '[{"secretId":"sec_01J9ZQ3ABCDEFGHJKMNPQRSTVW","delivery":"file","target":"/var/run/afs/secrets/github-app.pem"}]'

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
       Resume a pool

       afs pools resume-pool --team-id team_01J9ZQ3ABCDEFGHJKMNPQRSTVW --pool-id pool_01J9ZQ3ABCDEFGHJKMNPQRSTVW

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
       Scale a pool to 20 workers and wait until all 20 report

       afs pools scale-pool --team-id team_01J9ZQ3ABCDEFGHJKMNPQRSTVW --pool-id pool_01J9ZQ3ABCDEFGHJKMNPQRSTVW --workers 20 --wait

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
       --mcp (JSON)
          v1.6.0 MCP access of an agent or pool
       --model-id (string)
          Deprecated: use --default-model-id (and --models to choose the set) instead.
       --models (list)
          Models the agent may call: the word all (every model the team may use, now and later) or one or more model ids (a fixed selection). Sent members replace the stored ones.
       --persona (string)
       --resources (JSON)
       --scratch-size-gib (integer)
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
       --default-model-id (string)
          The model the runtime starts with; it must be one of the agent's models. Without --models on an update, only the default changes.
       --routes (list)
          Model routes the agent may call: the word all (every route of the team, now and later) or one or more route slugs (a fixed selection). Folded into models.routeAccess and models.routes.
       --default-route (string)
          The route of the default model; --default-model-id then names the route's model. Or give --default-model-id <route>/<model> alone.
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
       Let the pool's workers call every model the team may use

       afs pools update-pool --team-id team_01J9ZQ3ABCDEFGHJKMNPQRSTVW --pool-id pool_01J9ZQ3ABCDEFGHJKMNPQRSTVW --models all

       Change a pool's persona

       afs pools update-pool --team-id team_01J9ZQ3ABCDEFGHJKMNPQRSTVW --pool-id pool_01J9ZQ3ABCDEFGHJKMNPQRSTVW --persona "Minimal-change dependency upgrader"

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
          Range: 1 to 50
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
       Move a stage

       afs pools update-stage --team-id team_01J9ZQ3ABCDEFGHJKMNPQRSTVW --stage-id stg_01J9ZQ3ABCDEFGHJKMNPQRSTVW --position 2

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
       Check that the trust policy is in place

       afs pools verify-role-link --team-id team_01J9ZQ3ABCDEFGHJKMNPQRSTVW --pool-id pool_01J9ZQ3ABCDEFGHJKMNPQRSTVW --link-id prl_01J9ZQ3ABCDEFGHJKMNPQRSTVW

```
