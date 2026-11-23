# afs orchestration

Generated from `afs orchestration help`; do not edit by hand.

```text
AFS ORCHESTRATION

DESCRIPTION
       Manage a team's orchestration agent: its collaboration tool and who can reach it.

SYNOPSIS
       afs orchestration <operation> [parameters]

OPERATIONS
       change-collaboration-tool      Change collaboration tool
       describe-orchestration         Get team orchestration
       put-orchestration              Put team orchestration

To see help text for an operation, run: afs orchestration <operation> help
```

## change-collaboration-tool

```text
AFS ORCHESTRATION CHANGE-COLLABORATION-TOOL

NAME
       orchestration change-collaboration-tool

DESCRIPTION
       Change collaboration tool

       The orchestration agent is snapshotted (kept 90 days), then started with the new tool and an empty data directory: no content moves between tools. Surfaces the new tool lacks are set to off, and connected workers are registered again. Buzz has no Web UI and no Web API.

       Required permission: collaboration:manage

SYNOPSIS
       afs orchestration change-collaboration-tool [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --tool (string) (required)
          Possible values: paperclip, buzz
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
       Switch the team's collaboration tool to Buzz and wait until it is ready

       afs orchestration change-collaboration-tool --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --tool buzz --wait

       Switch back to Paperclip

       afs orchestration change-collaboration-tool --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --tool paperclip

```

## describe-orchestration

```text
AFS ORCHESTRATION DESCRIBE-ORCHESTRATION

NAME
       orchestration describe-orchestration

DESCRIPTION
       Get team orchestration

       Required permission: collaboration:read

SYNOPSIS
       afs orchestration describe-orchestration [options]

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
       Show a team's orchestration agent, its tool and what is effective

       afs orchestration describe-orchestration --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001

       Print the Web UI address

       afs orchestration describe-orchestration --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --query webUi.url --output text

```

## put-orchestration

```text
AFS ORCHESTRATION PUT-ORCHESTRATION

NAME
       orchestration put-orchestration

DESCRIPTION
       Put team orchestration

       Sets the exposure of both surfaces in one call: off (inside the namespace only), private (through the internal ALB from networks routed to the environment VPC) or public (through the internet-facing ALB). A level above the team's limit is refused; the limit is set by a service admin (afs teams update-team --max-collaboration-exposure).

       Required permission: collaboration:manage

SYNOPSIS
       afs orchestration put-orchestration [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --registration-mode (string)
          Omitted keeps the current value
          Possible values: manual, managed
       --if-match (string)
          Apply only if the resource is at this version (safe update).
       --web-ui-exposure (string)
          Exposure of the Web UI: off, private (internal ALB) or public (internet-facing ALB). Both surfaces are sent; give both unless --cli-input-json has them.
          Possible values: off, private, public
       --web-api-exposure (string)
          Exposure of the Web API: off, private or public. Buzz has no Web API: only off is accepted for it.
          Possible values: off, private, public
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
       Serve the Web UI privately and keep the Web API off

       afs orchestration put-orchestration --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --web-ui-exposure private --web-api-exposure off

       Publish both surfaces, use managed registration and wait until it is ready

       afs orchestration put-orchestration --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --web-ui-exposure public --web-api-exposure public --registration-mode managed --wait

```
