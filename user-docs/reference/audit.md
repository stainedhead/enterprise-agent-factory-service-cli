# afs audit

Generated from `afs audit help`; do not edit by hand.

```text
AFS AUDIT

DESCRIPTION
       Search and export the audit log.

SYNOPSIS
       afs audit <operation> [parameters]

OPERATIONS
       list-audit-events              List audit events
       list-mcp-events                List MCP events
       list-service-audit-events      List admin audit events
       list-team-mcp-events           List team MCP events

To see help text for an operation, run: afs audit <operation> help
```

## list-audit-events

```text
AFS AUDIT LIST-AUDIT-EVENTS

NAME
       audit list-audit-events

DESCRIPTION
       List audit events

       Required permission: audit:read

SYNOPSIS
       afs audit list-audit-events [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --from (string)
       --to (string)
       --actor-id (string)
       --action (string)
       --outcome (string)
          Possible values: allowed, denied, failed
       --resource-kind (string)
       --resource-id (string)

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
       List a team's audit events as CSV

       afs audit list-audit-events --team-id team_01J9ZQ3ABCDEFGHJKMNPQRSTVW --from 2026-11-01T00:00:00Z --outcome denied --output csv

       Events of one resource

       afs audit list-audit-events --team-id team_01J9ZQ3ABCDEFGHJKMNPQRSTVW --resource-kind agent --resource-id agt_01J9ZQ3ABCDEFGHJKMNPQRSTVW

```

## list-mcp-events

```text
AFS AUDIT LIST-MCP-EVENTS

NAME
       audit list-mcp-events

DESCRIPTION
       List MCP events

       Required permission: mcp-events:read

SYNOPSIS
       afs audit list-mcp-events [options]

OPTIONS
       --from (string)
       --to (string)
       --team-id (string)
       --agent-id (string)
       --pool-id (string)
       --server-id (string)
       --method (string)
       --name (string)
       --outcome (string)
          Possible values: ok, tool_error, protocol_error, refused, upstream_unavailable, timeout, cancelled
       --reason (string)

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
       List the refused MCP calls of the last day

       afs audit list-mcp-events --from 2026-11-24T00:00:00Z --to 2026-11-25T00:00:00Z --outcome refused

       List one agent's calls to one server

       afs audit list-mcp-events --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --agent-id agt_01JD2Q3R4S5T6V7W8X9Y0Z0001 --server-id mcp_01JD2Q3R4S5T6V7W8X9Y0Z0001

```

## list-service-audit-events

```text
AFS AUDIT LIST-SERVICE-AUDIT-EVENTS

NAME
       audit list-service-audit-events

DESCRIPTION
       List admin audit events

       Required permission: admin:audit

SYNOPSIS
       afs audit list-service-audit-events [options]

OPTIONS
       --from (string)
       --to (string)
       --actor-id (string)
       --action (string)
       --outcome (string)
          Possible values: allowed, denied, failed
       --resource-kind (string)
       --resource-id (string)
       --team-id (string)

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
       Service-wide audit events as CSV (platform roles)

       afs audit list-service-audit-events --from 2026-11-01T00:00:00Z --output csv

       Denied events of one team

       afs audit list-service-audit-events --team-id team_01J9ZQ3ABCDEFGHJKMNPQRSTVW --outcome denied

```

## list-team-mcp-events

```text
AFS AUDIT LIST-TEAM-MCP-EVENTS

NAME
       audit list-team-mcp-events

DESCRIPTION
       List team MCP events

       Required permission: mcp-events:read

SYNOPSIS
       afs audit list-team-mcp-events [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --from (string)
       --to (string)
       --agent-id (string)
       --pool-id (string)
       --server-id (string)
       --method (string)
       --name (string)
       --outcome (string)
          Possible values: ok, tool_error, protocol_error, refused, upstream_unavailable, timeout, cancelled
       --reason (string)

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
       List the team's MCP calls that timed out

       afs audit list-team-mcp-events --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --outcome timeout

```
