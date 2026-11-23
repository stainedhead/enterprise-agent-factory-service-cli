# afs collaboration

Generated from `afs collaboration help`; do not edit by hand.

```text
AFS COLLABORATION

DESCRIPTION
       

SYNOPSIS
       afs collaboration <operation> [parameters]

OPERATIONS
       connect-agent                  Connect agent to collaboration
       describe-collaboration         Get team collaboration
       disconnect-agent               Disconnect agent from collaboration
       invoke-api                     Call a team's Paperclip REST API (the orchestration agent's Web API) with your AFS credential and a Paperclip key.
       list-collaboration-links       List collaboration links
       put-collaboration              Put team collaboration
       rotate-collaboration-keys      Rotate collaboration keys

To see help text for an operation, run: afs collaboration <operation> help
```

## connect-agent

```text
AFS COLLABORATION CONNECT-AGENT

NAME
       collaboration connect-agent

DESCRIPTION
       Connect agent to collaboration

       Required permission: collaboration:manage

SYNOPSIS
       afs collaboration connect-agent [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --agent-id (string) (required)
       --paperclip-agent-id (string)
          Accepted only while registrationMode is manual
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
       Connect an agent to the team's collaboration tool

       afs collaboration connect-agent --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --agent-id agt_01JD2Q3R4S5T6V7W8X9Y0Z0001

       Connect an agent in manual registration mode, naming its Paperclip agent

       afs collaboration connect-agent --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --agent-id agt_01JD2Q3R4S5T6V7W8X9Y0Z0001 --paperclip-agent-id ag_8f2c1d

```

## describe-collaboration

```text
AFS COLLABORATION DESCRIBE-COLLABORATION

NAME
       collaboration describe-collaboration

DESCRIPTION
       Get team collaboration

       Deprecated: kept until 2.0.0; use afs orchestration describe-orchestration instead.

       Required permission: collaboration:read

SYNOPSIS
       afs collaboration describe-collaboration [options]

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
       Show a team's Paperclip endpoint and what is effective

       afs collaboration describe-collaboration --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001

       Print the public hostname only

       afs collaboration describe-collaboration --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --query hostname --output text

```

## disconnect-agent

```text
AFS COLLABORATION DISCONNECT-AGENT

NAME
       collaboration disconnect-agent

DESCRIPTION
       Disconnect agent from collaboration

       Required permission: collaboration:manage

SYNOPSIS
       afs collaboration disconnect-agent [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --agent-id (string) (required)
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
       Disconnect an agent from the team's collaboration tool

       afs collaboration disconnect-agent --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --agent-id agt_01JD2Q3R4S5T6V7W8X9Y0Z0001

```

## invoke-api

```text
AFS COLLABORATION INVOKE-API

NAME
       collaboration invoke-api

DESCRIPTION
       Call a team's Paperclip REST API (the orchestration agent's Web API) with your AFS credential and a Paperclip key.

SYNOPSIS
       afs collaboration invoke-api [options]

OPTIONS
       --team-id (string)
          The team (or the profile's team_id). Not needed with --collab-url.
       --method (string)
          GET, HEAD, POST, PUT, PATCH or DELETE (default GET).
       --path (string) (required)
          The Paperclip path, starting with / (for example /api/issues).
       --body (string)
          The request body: file://PATH or fileb://PATH.
       --content-type (string)
          The body's media type (default application/json).
       --paperclip-key-file (string) (required)
          A file (mode 0600) holding the Paperclip key, or - for standard input. The key is never an argument.
       --collab-url (string)
          Use this https origin instead of the Web API address of the team's orchestration agent.
       --include (boolean)
          Print the status line and response headers before the body.
       --wait (boolean)
          Not supported: invoke-api prints the response unchanged.

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
       Read the team's companies through the Web API at its effective level (private or public)

       afs collaboration invoke-api --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --path /api/companies --paperclip-key-file ./pc.key

       Create an issue from a pipeline

       afs collaboration invoke-api --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --method POST --path /api/issues --body file://issue.json --paperclip-key-file ./pc.key

```

## list-collaboration-links

```text
AFS COLLABORATION LIST-COLLABORATION-LINKS

NAME
       collaboration list-collaboration-links

DESCRIPTION
       List collaboration links

       Required permission: collaboration:read

SYNOPSIS
       afs collaboration list-collaboration-links [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --link-state (string)
          Possible values: pending, linked, failed, disabled
       --connect (boolean)

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
       List the link state of every worker agent

       afs collaboration list-collaboration-links --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001

       List the links that failed

       afs collaboration list-collaboration-links --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --link-state failed

```

## put-collaboration

```text
AFS COLLABORATION PUT-COLLABORATION

NAME
       collaboration put-collaboration

DESCRIPTION
       Put team collaboration

       Deprecated: kept until 2.0.0; use afs orchestration put-orchestration instead.

       The booleans stand for exposure levels: on read, public is true and anything else false; on write, true sets the highest level within the team's limit and false sets off. Private exposure cannot be expressed with booleans.

       Required permission: collaboration:manage

SYNOPSIS
       afs collaboration put-collaboration [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --registration-mode (string)
          Omitted keeps the current value
          Possible values: manual, managed
       --rest-api-enabled (boolean) (required)
       --web-ui-enabled (boolean) (required)
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
       Publish the web UI and keep the REST API off

       afs collaboration put-collaboration --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --web-ui-enabled --no-rest-api-enabled

       Turn both surfaces on, use managed registration and wait until it is ready

       afs collaboration put-collaboration --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --web-ui-enabled --rest-api-enabled --registration-mode managed --wait

```

## rotate-collaboration-keys

```text
AFS COLLABORATION ROTATE-COLLABORATION-KEYS

NAME
       collaboration rotate-collaboration-keys

DESCRIPTION
       Rotate collaboration keys

       Required permission: collaboration:manage

SYNOPSIS
       afs collaboration rotate-collaboration-keys [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --names (list)
          Default is both.
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
       Rotate both of Paperclip's signing secrets (this signs everyone out of Paperclip)

       afs collaboration rotate-collaboration-keys --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001

       Rotate one secret

       afs collaboration rotate-collaboration-keys --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --names paperclip-better-auth-secret

```
