# afs reports

Generated from `afs reports help`; do not edit by hand.

```text
AFS REPORTS

DESCRIPTION
       Run and schedule reports.

SYNOPSIS
       afs reports <operation> [parameters]

OPERATIONS
       create-report-schedule         Create report schedule
       delete-report-schedule         Delete report schedule
       get-agent-report               Get agent report
       get-inference-report           Get inference report
       get-mcp-usage-report           Get MCP usage report
       get-service-report             Get service report
       get-team-inference-report      Get team inference report
       get-team-mcp-usage-report      Get team MCP usage report
       get-team-report                Get team report
       list-report-schedules          List report schedules
       update-report-schedule         Update report schedule

To see help text for an operation, run: afs reports <operation> help
```

## create-report-schedule

```text
AFS REPORTS CREATE-REPORT-SCHEDULE

NAME
       reports create-report-schedule

DESCRIPTION
       Create report schedule

       Required permission: reports:schedule

SYNOPSIS
       afs reports create-report-schedule [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --cadence (string) (required)
          Possible values: daily, weekly, monthly
       --channel-id (string) (required)
       --enabled (boolean)
       --format (string)
          Possible values: csv
       --name (string) (required)
       --query-json (JSON)
          Stored report query parameters
       --report (string) (required)
          Known values: usage, costs, activity, traffic
       --idempotency-key (string)
          Idempotency key (default: a new ULID per invocation).
       --range (string)
          Report range: 24h, 7d, 30d, 90d or mtd.
          Possible values: 24h, 7d, 30d, 90d, mtd
       --from (string)
          Start of the report window (RFC 3339); use with --to instead of --range.
       --to (string)
          End of the report window (RFC 3339).
       --group-by (string)
          Comma-separated dimensions: team, agent, model, system, day, hour, pool, tag:<key>; surface, reason and person for the collaboration report.
       --tag (list)
          Filter by a tag, key=value (several allowed); sent as tag.<key>=<value>.

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
       Email the cost-by-project report every week

       afs reports create-report-schedule --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --name weekly-costs --report costs --cadence weekly --channel-id nch_01JD2Q3R4S5T6V7W8X9Y0Z0001 --range 30d --group-by tag:project

```

## delete-report-schedule

```text
AFS REPORTS DELETE-REPORT-SCHEDULE

NAME
       reports delete-report-schedule

DESCRIPTION
       Delete report schedule

       Required permission: reports:schedule

SYNOPSIS
       afs reports delete-report-schedule [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --schedule-id (string) (required)
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
       Delete a schedule

       afs reports delete-report-schedule --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --schedule-id rsc_01J9ZQ3ABCDEFGHJKMNPQRSTVW

```

## get-agent-report

```text
AFS REPORTS GET-AGENT-REPORT

NAME
       reports get-agent-report

DESCRIPTION
       Get agent report

       Required permission: reports:read

SYNOPSIS
       afs reports get-agent-report [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --agent-id (string) (required)
       --kind (string) (required)
          Possible values: usage, costs, activity, traffic
       --from (string)
       --to (string)
       --range (string)
          Instead of from and to
          Possible values: 24h, 7d, 30d, 90d, mtd
       --group-by (string)
          Comma-separated: team, agent, model, system, day, hour, tag:<key>, pool (ephemeral teams), surface, reason and person (collaboration), and route (usage and costs, v1.5.0).
       --model-id (string)
       --route (string)
          v1.5.0, usage and costs reports; a route slug, bedrock for Bedrock calls
       --system-id (string)
       --pool-id (string)
          Ephemeral teams only
       --tag (list)
          Filter by a tag, key=value (several allowed); sent as tag.<key>=<value>.

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
       One agent's activity over 24 hours

       afs reports get-agent-report --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --agent-id agt_01JD2Q3R4S5T6V7W8X9Y0Z0001 --kind activity --range 24h

```

## get-inference-report

```text
AFS REPORTS GET-INFERENCE-REPORT

NAME
       reports get-inference-report

DESCRIPTION
       Get inference report

       The service-wide inference report (platform roles); --team-id narrows it to one team.

       Required permission: admin:reports

SYNOPSIS
       afs reports get-inference-report [options]

OPTIONS
       --team-id (string)
          Filter (service reports only)
       --from (string)
       --to (string)
       --range (string)
          Instead of from and to
          Possible values: 24h, 7d, 30d, 90d, mtd
       --group-by (string)
          Comma-separated: route, model, agent, hour, day (and team on the service report, pool in ephemeral teams)
       --route (string)
          Filter by route (bedrock for Bedrock)
       --model-id (string)
       --agent-id (string)
       --pool-id (string)
          Ephemeral teams only

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
       Every route this month, by route

       afs reports get-inference-report --range mtd --group-by route

       One team's hourly TTFT on a route

       afs reports get-inference-report --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --route openai-prod --range 24h --group-by hour

```

## get-mcp-usage-report

```text
AFS REPORTS GET-MCP-USAGE-REPORT

NAME
       reports get-mcp-usage-report

DESCRIPTION
       Get MCP usage report

       Required permission: mcp-events:read

SYNOPSIS
       afs reports get-mcp-usage-report [options]

OPTIONS
       --from (string)
       --to (string)
       --group-by (string)
          Possible values: server, primitive, team, agent, day
       --view (string)
          Possible values: usage, unused, refusals, discovery, adoption
       --window (integer)
          Days, for view unused
          Possible values: 30, 90

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
       Calls per server this month

       afs reports get-mcp-usage-report --from 2026-11-01T00:00:00Z --to 2026-12-01T00:00:00Z --group-by server --view usage

       Servers no team used in 90 days

       afs reports get-mcp-usage-report --from 2026-08-27T00:00:00Z --to 2026-11-25T00:00:00Z --group-by server --view unused --window 90

```

## get-service-report

```text
AFS REPORTS GET-SERVICE-REPORT

NAME
       reports get-service-report

DESCRIPTION
       Get service report

       Required permission: admin:reports

SYNOPSIS
       afs reports get-service-report [options]

OPTIONS
       --kind (string) (required)
          Possible values: usage, costs, activity, traffic, collaboration
       --from (string)
       --to (string)
       --range (string)
          Instead of from and to
          Possible values: 24h, 7d, 30d, 90d, mtd
       --group-by (string)
          Comma-separated: team, agent, model, system, day, hour, tag:<key>, pool (ephemeral teams), surface, reason and person (collaboration), and route (usage and costs, v1.5.0).
       --team-id (string)
          Filter (service reports only)
       --agent-id (string)
       --model-id (string)
       --route (string)
          v1.5.0, usage and costs reports; a route slug, bedrock for Bedrock calls
       --system-id (string)
       --pool-id (string)
          Ephemeral teams only
       --tag (list)
          Filter by a tag, key=value (several allowed); sent as tag.<key>=<value>.

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
       Service-wide costs by team this month (platform roles)

       afs reports get-service-report --kind costs --range mtd --group-by team

       Paperclip use by team (platform roles)

       afs reports get-service-report --kind collaboration --group-by team

```

## get-team-inference-report

```text
AFS REPORTS GET-TEAM-INFERENCE-REPORT

NAME
       reports get-team-inference-report

DESCRIPTION
       Get team inference report

       Latency, time to first token, throughput, errors (upstream errors apart from AFS refusals) and cost, from the gateway's ledger; no prompt or completion text is ever reported.

       Required permission: reports:read

SYNOPSIS
       afs reports get-team-inference-report [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --from (string)
       --to (string)
       --range (string)
          Instead of from and to
          Possible values: 24h, 7d, 30d, 90d, mtd
       --group-by (string)
          Comma-separated: route, model, agent, hour, day (and team on the service report, pool in ephemeral teams)
       --route (string)
          Filter by route (bedrock for Bedrock)
       --model-id (string)
       --agent-id (string)
       --pool-id (string)
          Ephemeral teams only

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
       The team's routes over the last day, by route and model

       afs reports get-team-inference-report --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --range 24h --group-by route,model

       One route by agent over a week, as CSV

       afs reports get-team-inference-report --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --range 7d --route openai-prod --group-by agent --output csv

```

## get-team-mcp-usage-report

```text
AFS REPORTS GET-TEAM-MCP-USAGE-REPORT

NAME
       reports get-team-mcp-usage-report

DESCRIPTION
       Get team MCP usage report

       Required permission: mcp-events:read

SYNOPSIS
       afs reports get-team-mcp-usage-report [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --from (string)
       --to (string)
       --group-by (string)
          Possible values: server, primitive, agent, day
       --view (string)
          Possible values: usage, unused, refusals, discovery
       --window (integer)
          Days, for view unused
          Possible values: 30, 90

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
       The team's refusals per agent

       afs reports get-team-mcp-usage-report --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --from 2026-11-01T00:00:00Z --to 2026-12-01T00:00:00Z --group-by agent --view refusals

```

## get-team-report

```text
AFS REPORTS GET-TEAM-REPORT

NAME
       reports get-team-report

DESCRIPTION
       Get team report

       Required permission: reports:read

SYNOPSIS
       afs reports get-team-report [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --kind (string) (required)
          Possible values: usage, costs, activity, traffic, collaboration
       --from (string)
       --to (string)
       --range (string)
          Instead of from and to
          Possible values: 24h, 7d, 30d, 90d, mtd
       --group-by (string)
          Comma-separated: team, agent, model, system, day, hour, tag:<key>, pool (ephemeral teams), surface, reason and person (collaboration), and route (usage and costs, v1.5.0).
       --agent-id (string)
       --model-id (string)
       --route (string)
          v1.5.0, usage and costs reports; a route slug, bedrock for Bedrock calls
       --system-id (string)
       --pool-id (string)
          Ephemeral teams only
       --tag (list)
          Filter by a tag, key=value (several allowed); sent as tag.<key>=<value>.

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
       Cost by project tag over 30 days, as CSV

       afs reports get-team-report --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --kind costs --range 30d --group-by tag:project --output csv

       Usage by agent and model, filtered by tag

       afs reports get-team-report --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --kind usage --range 7d --group-by agent,model --tag project=checkout

       Who used Paperclip in the last week, by day and person

       afs reports get-team-report --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --kind collaboration --range 7d --group-by day,person

```

## list-report-schedules

```text
AFS REPORTS LIST-REPORT-SCHEDULES

NAME
       reports list-report-schedules

DESCRIPTION
       List report schedules

       Required permission: reports:read

SYNOPSIS
       afs reports list-report-schedules [options]

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
       List a team's scheduled reports

       afs reports list-report-schedules --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001

```

## update-report-schedule

```text
AFS REPORTS UPDATE-REPORT-SCHEDULE

NAME
       reports update-report-schedule

DESCRIPTION
       Update report schedule

       Required permission: reports:schedule

SYNOPSIS
       afs reports update-report-schedule [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --schedule-id (string) (required)
       --cadence (string)
          Possible values: daily, weekly, monthly
       --channel-id (string)
       --enabled (boolean)
       --name (string)
       --query-json (JSON)
          Stored report query parameters
       --if-match (string)
          Apply only if the resource is at this version (safe update).
       --range (string)
          Report range: 24h, 7d, 30d, 90d or mtd.
          Possible values: 24h, 7d, 30d, 90d, mtd
       --from (string)
          Start of the report window (RFC 3339); use with --to instead of --range.
       --to (string)
          End of the report window (RFC 3339).
       --group-by (string)
          Comma-separated dimensions: team, agent, model, system, day, hour, pool, tag:<key>; surface, reason and person for the collaboration report.
       --tag (list)
          Filter by a tag, key=value (several allowed); sent as tag.<key>=<value>.

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
       Pause a schedule

       afs reports update-report-schedule --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --schedule-id rsc_01J9ZQ3ABCDEFGHJKMNPQRSTVW --no-enabled

```
