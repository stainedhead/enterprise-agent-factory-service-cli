# afs teams

Generated from `afs teams help`; do not edit by hand.

```text
AFS TEAMS

DESCRIPTION
       Create and manage teams.

SYNOPSIS
       afs teams <operation> [parameters]

OPERATIONS
       approve-team-mcp-primitive     Approve team MCP primitive
       approve-team-mcp-server        Approve team MCP server
       create-log-destination         Create log destination
       create-notification-channel    Create notification channel
       create-team                    Create team
       create-team-mcp-server         Create team MCP server
       declare-external-system        Declare team external system
       delete-config-value            Delete namespace config value
       delete-external-system         Delete team external system
       delete-log-destination         Delete log destination
       delete-notification-channel    Delete notification channel
       delete-team                    Delete team
       delete-team-mcp-server         Delete team MCP server
       describe-log-destination       Get log destination
       describe-notification-channel  Get notification channel
       describe-revision              Get revision
       describe-team                  Get team
       describe-team-mcp-server       Get team MCP server
       disconnect-mcp-oauth           Disconnect MCP OAuth
       get-kubeconfig                 Get team kubeconfig
       get-team-status                Get team status
       get-team-tags                  Get team tags
       list-allowed-accounts          List allowed accounts
       list-config-values             List namespace config values
       list-external-systems          List team external systems
       list-log-destinations          List log destinations
       list-notification-channels     List notification channels
       list-revisions                 List revisions
       list-team-directory            List team directory
       list-team-mcp-primitives       List team MCP primitives
       list-team-mcp-servers          List team MCP servers
       list-team-models               List team models
       list-teams                     List teams
       put-config-value               Put namespace config value
       put-team-gateway-settings      Put team gateway settings
       put-team-mcp-policy            Put team MCP policy
       put-team-mcp-server            Put team MCP server
       put-team-mcp-server-credential Put team MCP server credential
       put-team-model-endpoints       Put team model endpoints
       put-team-models                Put team models
       put-team-tags                  Put team tags
       refresh-team-mcp-inventory     Refresh team MCP inventory
       reject-team-mcp-primitive      Reject team MCP primitive
       reject-team-mcp-server         Reject team MCP server
       restore-revision               Restore revision
       resume-team                    Resume team
       set-external-system-state      Set team external system state
       start-mcp-oauth-connect        Start MCP OAuth connect
       submit-mcp-server-for-review   Submit MCP server for review
       suspend-team                   Suspend team
       test-mcp-federation            Test MCP federation
       test-notification-channel      Test notification channel
       update-external-system         Update team external system
       update-log-destination         Update log destination
       update-notification-channel    Update notification channel
       update-team                    Update team
       update-team-mcp-server         Update team MCP server

To see help text for an operation, run: afs teams <operation> help
```

## approve-team-mcp-primitive

```text
AFS TEAMS APPROVE-TEAM-MCP-PRIMITIVE

NAME
       teams approve-team-mcp-primitive

DESCRIPTION
       Approve team MCP primitive

       Required permission: team-mcp-servers:approve

SYNOPSIS
       afs teams approve-team-mcp-primitive [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --mcp-server-id (string) (required)
       --primitive-id (string) (required)
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
       Approve a changed tool of the team's server

       afs teams approve-team-mcp-primitive --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --mcp-server-id mcp_01JD2Q3R4S5T6V7W8X9Y0Z0002 --primitive-id mpr_01JD2Q3R4S5T6V7W8X9Y0Z0001 --note "New optional argument"

```

## approve-team-mcp-server

```text
AFS TEAMS APPROVE-TEAM-MCP-SERVER

NAME
       teams approve-team-mcp-server

DESCRIPTION
       Approve team MCP server

       A team admin approves a low-risk team server when mcp.approval_policy allows it; others go to a service admin or a security operator.

       Required permission: team-mcp-servers:approve

SYNOPSIS
       afs teams approve-team-mcp-server [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --mcp-server-id (string) (required)
       --note (string) (required)
       --primitives (JSON)
          Per-primitive decisions as JSON or file://PATH: [{"primitiveId":…, "decision":"approve"|"block"}]; primitives not named keep their state.
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
       Approve the team's low-risk server

       afs teams approve-team-mcp-server --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --mcp-server-id mcp_01JD2Q3R4S5T6V7W8X9Y0Z0002 --note "Read-only lookups"

```

## create-log-destination

```text
AFS TEAMS CREATE-LOG-DESTINATION

NAME
       teams create-log-destination

DESCRIPTION
       Create log destination

       Required permission: logs:write

SYNOPSIS
       afs teams create-log-destination [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --config (JSON) (required)
       --enabled (boolean)
       --kind (string) (required)
          Possible values: cloudwatch, s3, otlp, splunk_hec, datadog
       --name (string) (required)
       --secret-id (string)
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
       Ship logs to CloudWatch and wait until it is applied

       afs teams create-log-destination --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --name team-cloudwatch --kind cloudwatch --config '{"logGroup":"/afs/payments","region":"us-east-1"}' --wait

```

## create-notification-channel

```text
AFS TEAMS CREATE-NOTIFICATION-CHANNEL

NAME
       teams create-notification-channel

DESCRIPTION
       Create notification channel

       Required permission: notifications:write

SYNOPSIS
       afs teams create-notification-channel [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --config (JSON) (required)
       --enabled (boolean)
       --event-types (list) (required)
       --kind (string) (required)
          Possible values: slack, email, webhook, jira, eventbridge
       --name (string) (required)
       --secret-id (string)
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
       Create a Slack channel and wait until it is ready

       afs teams create-notification-channel --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --name payments-slack --kind slack --config '{"channel":"#payments-alerts"}' --event-types 'budget.*' 'agent.status_changed' --secret-id sec_01JD2Q3R4S5T6V7W8X9Y0Z0003 --wait

```

## create-team

```text
AFS TEAMS CREATE-TEAM

NAME
       teams create-team

DESCRIPTION
       Create team

       Required permission: teams:create

SYNOPSIS
       afs teams create-team [options]

OPTIONS
       --allowed-aws-accounts (list)
          Needs teams:set-aws-accounts on update
       --budget (JSON)
          Lifetime budget of an ephemeral team, created with it
       --cost-center (string)
       --description (string)
       --display-name (string) (required)
       --expires-at (string)
          RFC 3339 UTC with Z
       --first-admin (JSON) (required)
          The first team-admin binding: a user (userId or email) or an IdP group
       --identity-provider-id (string)
       --isolation-tier (string)
          Guardrail field: needs teams:set-guardrails
          Possible values: shared, dedicated
       --kind (string) (required)
          Immutable after create
          Possible values: persistent, ephemeral
       --max-budget-usd (string) (required)
          Guardrail field: needs teams:set-guardrails
       --max-collaboration-exposure (string)
          Default off
          Possible values: off, private, public
       --may-add-external-systems (boolean)
          Guardrail field: needs teams:set-guardrails
       --may-expose-collaboration (boolean)
          Deprecated: use --max-collaboration-exposure (off, private or public) instead.
       --may-manage-aws-accounts (boolean)
          Guardrail field: needs teams:set-guardrails
       --model-access (string)
          Possible values: all, selected
       --models (list)
          Model ids enabled for the team
       --orchestration (string or JSON)
          The team's orchestration agent, all optional: Tool=paperclip|buzz,WebUiExposure=off|private|public,WebApiExposure=off|private|public,RegistrationMode=manual|managed,ModelId=<model id>,VolumeSizeGiB=10..1000,VolumeOnDelete=snapshot_then_delete|delete|retain, or JSON (file://PATH). Defaults come from the orchestration.* settings.
       --owner-email (string)
       --quotas (JSON)
       --slug (string) (required)
          Immutable name (team slug, agent, pool, secret)
       --tags (key=value pairs or JSON)
          Team-owned tag values, including required service tags
       --teardown-rule (string)
          Possible values: at_expiry, at_cap_or_expiry, at_work_complete
       --viewer-kubectl (boolean)
          Guardrail field: needs teams:set-guardrails
       --idempotency-key (string)
          Idempotency key (default: a new ULID per invocation).
       --budget-usd (string)
          Lifetime budget of an ephemeral team in USD, for example 500.00 (required with --kind ephemeral).
       --budget-thresholds (list)
          Percent-of-budget alert thresholds of an ephemeral team, for example 50 80 100.
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
       Create an ephemeral team with a lifetime budget

       afs teams create-team --slug cve-sprint --display-name "CVE sprint" --kind ephemeral --max-budget-usd 500.00 --first-admin '{"email":"ada@example.com"}' --expires-at 2026-12-31T00:00:00Z --teardown-rule at_cap_or_expiry --budget-usd 500.00 --budget-thresholds 50 80 100

       Create a team whose orchestration agent runs Buzz

       afs teams create-team --slug payments --display-name "Payments" --kind persistent --max-budget-usd 500.00 --first-admin '{"email":"ada@example.com"}' --orchestration Tool=buzz --wait

       Create a team that may expose its Web UI privately, and serve it that way

       afs teams create-team --slug billing --display-name "Billing" --kind persistent --max-budget-usd 300.00 --first-admin '{"email":"ada@example.com"}' --max-collaboration-exposure private --orchestration Tool=paperclip,WebUiExposure=private,WebApiExposure=off,ModelId=anthropic.claude-sonnet-4-5

       Create a persistent team and wait until it is ready

       afs teams create-team --slug payments --display-name "Payments" --kind persistent --max-budget-usd 500.00 --first-admin '{"email":"ada@example.com"}' --tags project=checkout cost-center=CC-1001 --wait

       Create a team from a file

       afs teams create-team --cli-input-json file://team.json

       Print the input skeleton

       afs teams create-team --generate-cli-skeleton

```

## create-team-mcp-server

```text
AFS TEAMS CREATE-TEAM-MCP-SERVER

NAME
       teams create-team-mcp-server

DESCRIPTION
       Create team MCP server

       Allowed when the team's gateway settings have mayAddMcpServers (otherwise AccessDenied, reason team_servers_not_allowed). The server starts as a draft; submit it for review.

       Required permission: team-mcp-servers:write

SYNOPSIS
       afs teams create-team-mcp-server [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --auth-header (string)
       --auth-kind (string) (required)
          Possible values: none, static, oauth_client, workload_federation, oauth_user, aws_sigv4
       --ca-bundle (string or JSON)
          PEM bundle of the CA that signs the server's certificate, as text or file://PATH (not a secret).
       --data-classification (string) (required)
          Possible values: public, internal, confidential, restricted
       --display-name (string) (required)
       --doc-url (string)
       --federation (JSON)
          Workload federation as JSON or file://PATH: the token endpoint and audience the gateway exchanges its projected token with.
       --hosted (JSON)
          A hosted server as JSON or file://PATH: {"imageDigest":…, "args":[…], "envNames":[…], "egressHosts":[…], "instanceMode":"per_team", "needsAws":false}.
       --limits (JSON)
          The server's own limits as JSON: {"requestBytes":…, "responseBytes":…, "timeoutSeconds":…, "concurrentPerAgent":…}; unset members take mcp.limits.
       --location (string) (required)
          Possible values: external, corporate, in_namespace
       --owner (string) (required)
       --purpose (string) (required)
       --slug (string) (required)
       --transport (string) (required)
          Possible values: streamable_http, hosted
       --url (string)
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
       Register the team's own server

       afs teams create-team-mcp-server --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --slug payments-tools --display-name "Payments tools" --owner payments-team --purpose "Ledger lookups" --data-classification confidential --transport streamable_http --location corporate --url https://tools.payments.corp.example/mcp --auth-kind workload_federation --federation file://federation.json

```

## declare-external-system

```text
AFS TEAMS DECLARE-EXTERNAL-SYSTEM

NAME
       teams declare-external-system

DESCRIPTION
       Declare team external system

       Required permission: external-systems:write

SYNOPSIS
       afs teams declare-external-system [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --category (string) (required)
          Possible values: dev_tools, collaboration, email, cloud, identity, llm, other, mcp
       --description (string)
       --docs-url (string)
       --endpoints (list) (required)
       --name (string) (required)
       --owner (string) (required)
       --purpose (string) (required)
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
       Declare a team-owned external system

       afs teams declare-external-system --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --slug payments-ledger --name "Ledger API" --category other --owner ledger-team@example.com --purpose "Reconcile payouts" --endpoints ledger.example.com:443

```

## delete-config-value

```text
AFS TEAMS DELETE-CONFIG-VALUE

NAME
       teams delete-config-value

DESCRIPTION
       Delete namespace config value

       Required permission: config:write

SYNOPSIS
       afs teams delete-config-value [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
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
       Remove a team-wide value

       afs teams delete-config-value --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --name REGION

```

## delete-external-system

```text
AFS TEAMS DELETE-EXTERNAL-SYSTEM

NAME
       teams delete-external-system

DESCRIPTION
       Delete team external system

       Required permission: external-systems:write

SYNOPSIS
       afs teams delete-external-system [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
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
       Delete a team-owned external system

       afs teams delete-external-system --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --system-id ext_01JD2Q3R4S5T6V7W8X9Y0Z0001

```

## delete-log-destination

```text
AFS TEAMS DELETE-LOG-DESTINATION

NAME
       teams delete-log-destination

DESCRIPTION
       Delete log destination

       Required permission: logs:write

SYNOPSIS
       afs teams delete-log-destination [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --destination-id (string) (required)
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
       Delete a destination and wait until it is gone

       afs teams delete-log-destination --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --destination-id log_01J9ZQ3ABCDEFGHJKMNPQRSTVW --wait

```

## delete-notification-channel

```text
AFS TEAMS DELETE-NOTIFICATION-CHANNEL

NAME
       teams delete-notification-channel

DESCRIPTION
       Delete notification channel

       Required permission: notifications:write

SYNOPSIS
       afs teams delete-notification-channel [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --channel-id (string) (required)
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
       Delete a channel and wait until it is gone

       afs teams delete-notification-channel --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --channel-id nch_01JD2Q3R4S5T6V7W8X9Y0Z0001 --wait

```

## delete-team

```text
AFS TEAMS DELETE-TEAM

NAME
       teams delete-team

DESCRIPTION
       Delete team

       Required permission: teams:delete

SYNOPSIS
       afs teams delete-team [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
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
       Delete a team and wait until it is gone

       afs teams delete-team --team-id team_01J9ZQ3ABCDEFGHJKMNPQRSTVW --wait

```

## delete-team-mcp-server

```text
AFS TEAMS DELETE-TEAM-MCP-SERVER

NAME
       teams delete-team-mcp-server

DESCRIPTION
       Delete team MCP server

       Required permission: team-mcp-servers:write

SYNOPSIS
       afs teams delete-team-mcp-server [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --mcp-server-id (string) (required)
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
       Delete a draft team server

       afs teams delete-team-mcp-server --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --mcp-server-id mcp_01JD2Q3R4S5T6V7W8X9Y0Z0002

```

## describe-log-destination

```text
AFS TEAMS DESCRIBE-LOG-DESTINATION

NAME
       teams describe-log-destination

DESCRIPTION
       Get log destination

       Required permission: logs:read

SYNOPSIS
       afs teams describe-log-destination [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --destination-id (string) (required)

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
       Show a destination's configuration and delivery status

       afs teams describe-log-destination --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --destination-id log_01J9ZQ3ABCDEFGHJKMNPQRSTVW

```

## describe-notification-channel

```text
AFS TEAMS DESCRIBE-NOTIFICATION-CHANNEL

NAME
       teams describe-notification-channel

DESCRIPTION
       Get notification channel

       Required permission: notifications:read

SYNOPSIS
       afs teams describe-notification-channel [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --channel-id (string) (required)

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
       Show a channel's configuration and delivery status (never a credential)

       afs teams describe-notification-channel --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --channel-id nch_01JD2Q3R4S5T6V7W8X9Y0Z0001

```

## describe-revision

```text
AFS TEAMS DESCRIBE-REVISION

NAME
       teams describe-revision

DESCRIPTION
       Get revision

       Required permission: revisions:read

SYNOPSIS
       afs teams describe-revision [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --revision-id (integer) (required)
          Range: 1 to unbounded

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
       Show one revision with its snapshot

       afs teams describe-revision --team-id team_01J9ZQ3ABCDEFGHJKMNPQRSTVW --revision-id 12 --output json

```

## describe-team

```text
AFS TEAMS DESCRIBE-TEAM

NAME
       teams describe-team

DESCRIPTION
       Get team

       Required permission: teams:read

SYNOPSIS
       afs teams describe-team [options]

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
       Show one team

       afs teams describe-team --team-id team_01J9ZQ3ABCDEFGHJKMNPQRSTVW

```

## describe-team-mcp-server

```text
AFS TEAMS DESCRIBE-TEAM-MCP-SERVER

NAME
       teams describe-team-mcp-server

DESCRIPTION
       Get team MCP server

       Required permission: team-mcp-servers:read

SYNOPSIS
       afs teams describe-team-mcp-server [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --mcp-server-id (string) (required)

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
       Show a server as the team sees it

       afs teams describe-team-mcp-server --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --mcp-server-id mcp_01JD2Q3R4S5T6V7W8X9Y0Z0001

```

## disconnect-mcp-oauth

```text
AFS TEAMS DISCONNECT-MCP-OAUTH

NAME
       teams disconnect-mcp-oauth

DESCRIPTION
       Disconnect MCP OAuth

       Required permission: team-mcp-servers:write

SYNOPSIS
       afs teams disconnect-mcp-oauth [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --mcp-server-id (string) (required)
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
       Drop the team's OAuth connection

       afs teams disconnect-mcp-oauth --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --mcp-server-id mcp_01JD2Q3R4S5T6V7W8X9Y0Z0001

```

## get-kubeconfig

```text
AFS TEAMS GET-KUBECONFIG

NAME
       teams get-kubeconfig

DESCRIPTION
       Get team kubeconfig

       Required permission: teams:get-kubeconfig

SYNOPSIS
       afs teams get-kubeconfig [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --kubeconfig (string)
          Kubeconfig file to merge into (default: the first path of $KUBECONFIG, else ~/.kube/config); - prints the API's YAML to stdout.
       --aws-profile (string)
          Set AWS_PROFILE to this profile in the exec credential plugin.

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
       Merge the team's context into ~/.kube/config

       afs teams get-kubeconfig --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001

       Write a separate file and use an AWS profile for the token

       afs teams get-kubeconfig --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --kubeconfig ./kc.yaml --aws-profile ops

       Print the YAML only

       afs teams get-kubeconfig --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --kubeconfig -

```

## get-team-status

```text
AFS TEAMS GET-TEAM-STATUS

NAME
       teams get-team-status

DESCRIPTION
       Get team status

       Required permission: teams:read

SYNOPSIS
       afs teams get-team-status [options]

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
       Show reconcile status

       afs teams get-team-status --team-id team_01J9ZQ3ABCDEFGHJKMNPQRSTVW

```

## get-team-tags

```text
AFS TEAMS GET-TEAM-TAGS

NAME
       teams get-team-tags

DESCRIPTION
       Get team tags

       Required permission: tags:read

SYNOPSIS
       afs teams get-team-tags [options]

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
       Show every tag of a team and where it comes from

       afs teams get-team-tags --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001

       Print the tags as JSON

       afs teams get-team-tags --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --output json

```

## list-allowed-accounts

```text
AFS TEAMS LIST-ALLOWED-ACCOUNTS

NAME
       teams list-allowed-accounts

DESCRIPTION
       List allowed accounts

       Required permission: identities:read

SYNOPSIS
       afs teams list-allowed-accounts [options]

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
       List the AWS accounts a team may link roles in

       afs teams list-allowed-accounts --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001

```

## list-config-values

```text
AFS TEAMS LIST-CONFIG-VALUES

NAME
       teams list-config-values

DESCRIPTION
       List namespace config values

       Required permission: config:read

SYNOPSIS
       afs teams list-config-values [options]

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
       List the team's namespace configuration values

       afs teams list-config-values --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001

```

## list-external-systems

```text
AFS TEAMS LIST-EXTERNAL-SYSTEMS

NAME
       teams list-external-systems

DESCRIPTION
       List team external systems

       Required permission: external-systems:read

SYNOPSIS
       afs teams list-external-systems [options]

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
       List the external systems the team may use

       afs teams list-external-systems --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001

```

## list-log-destinations

```text
AFS TEAMS LIST-LOG-DESTINATIONS

NAME
       teams list-log-destinations

DESCRIPTION
       List log destinations

       Required permission: logs:read

SYNOPSIS
       afs teams list-log-destinations [options]

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
       List a team's log destinations

       afs teams list-log-destinations --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001

```

## list-notification-channels

```text
AFS TEAMS LIST-NOTIFICATION-CHANNELS

NAME
       teams list-notification-channels

DESCRIPTION
       List notification channels

       Required permission: notifications:read

SYNOPSIS
       afs teams list-notification-channels [options]

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
       List the team's notification channels

       afs teams list-notification-channels --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001

```

## list-revisions

```text
AFS TEAMS LIST-REVISIONS

NAME
       teams list-revisions

DESCRIPTION
       List revisions

       Required permission: revisions:read

SYNOPSIS
       afs teams list-revisions [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
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
       List the configuration revisions of a team

       afs teams list-revisions --team-id team_01J9ZQ3ABCDEFGHJKMNPQRSTVW

```

## list-team-directory

```text
AFS TEAMS LIST-TEAM-DIRECTORY

NAME
       teams list-team-directory

DESCRIPTION
       List team directory

       Required permission: teams:read-directory

SYNOPSIS
       afs teams list-team-directory [options]

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
       Look up team ids to name as repository owners

       afs teams list-team-directory

```

## list-team-mcp-primitives

```text
AFS TEAMS LIST-TEAM-MCP-PRIMITIVES

NAME
       teams list-team-mcp-primitives

DESCRIPTION
       List team MCP primitives

       Required permission: team-mcp-servers:read

SYNOPSIS
       afs teams list-team-mcp-primitives [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --mcp-server-id (string) (required)
       --kind (string)
          Possible values: tool, resource, resource_template, prompt
       --status (string)
          Possible values: approved, pending_review, blocked

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
       List the tools of a server for the team

       afs teams list-team-mcp-primitives --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --mcp-server-id mcp_01JD2Q3R4S5T6V7W8X9Y0Z0001 --kind tool

```

## list-team-mcp-servers

```text
AFS TEAMS LIST-TEAM-MCP-SERVERS

NAME
       teams list-team-mcp-servers

DESCRIPTION
       List team MCP servers

       Required permission: team-mcp-servers:read

SYNOPSIS
       afs teams list-team-mcp-servers [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --state (string)
          Possible values: on, off, pending

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
       List the servers a team can use

       afs teams list-team-mcp-servers --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001

       List only the servers that are on

       afs teams list-team-mcp-servers --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --state on

```

## list-team-models

```text
AFS TEAMS LIST-TEAM-MODELS

NAME
       teams list-team-models

DESCRIPTION
       List team models

       Required permission: teams:read

SYNOPSIS
       afs teams list-team-models [options]

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
       List the models a team may use

       afs teams list-team-models --team-id team_01J9ZQ3ABCDEFGHJKMNPQRSTVW

       List the ids of the models the gateway can call for the team

       afs teams list-team-models --team-id team_01J9ZQ3ABCDEFGHJKMNPQRSTVW --query 'items[?callable].id' --output text

```

## list-teams

```text
AFS TEAMS LIST-TEAMS

NAME
       teams list-teams

DESCRIPTION
       List teams

       Required permission: teams:read

SYNOPSIS
       afs teams list-teams [options]

OPTIONS
       --state (string)
          Possible values: active, suspended, deleting, deleted
       --kind (string)
          Possible values: persistent, ephemeral

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
       List the teams you can see

       afs teams list-teams

       List suspended teams as ids and names only

       afs teams list-teams --state suspended --query 'items[].[id,displayName]' --output text

```

## put-config-value

```text
AFS TEAMS PUT-CONFIG-VALUE

NAME
       teams put-config-value

DESCRIPTION
       Put namespace config value

       Required permission: config:write

SYNOPSIS
       afs teams put-config-value [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --name (string) (required)
       --value (string) (required)
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
       Set a team-wide value (reads the version first; a missing value is created)

       afs teams put-config-value --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --name REGION --value eu-west-1

```

## put-team-gateway-settings

```text
AFS TEAMS PUT-TEAM-GATEWAY-SETTINGS

NAME
       teams put-team-gateway-settings

DESCRIPTION
       Put team gateway settings

       shared sends the team's agents through the service's shared gateway; per_namespace runs the team's own gateway in its namespace (agents roll to it once it is ready). mayAddModelEndpoints lets the team's admins declare their own endpoints.

       Required permission: teams:set-guardrails

SYNOPSIS
       afs teams put-team-gateway-settings [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --deployment (string)
          Possible values: shared, per_namespace
       --may-add-mcp-servers (boolean)
          v1.6.0
       --may-add-model-endpoints (boolean)
       --mcp-placement (string)
          v1.6.0
          Possible values: combined, separate
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
       Give a team its own gateway

       afs teams put-team-gateway-settings --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --deployment per_namespace

       Let a team declare its own model endpoints

       afs teams put-team-gateway-settings --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --may-add-model-endpoints

```

## put-team-mcp-policy

```text
AFS TEAMS PUT-TEAM-MCP-POLICY

NAME
       teams put-team-mcp-policy

DESCRIPTION
       Put team MCP policy

       Required permission: team-mcp-servers:write

SYNOPSIS
       afs teams put-team-mcp-policy [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --mcp-server-id (string) (required)
       --rules (JSON) (required)
          Rules as JSON or file://PATH: [{"name":"delete_*","effect":"block"}, {"class":"destructive","effect":"block"}, {"name":"search","effect":"allow","rateLimitPerMinute":60}]. The most specific rule of a level wins; a block at any level wins over the levels below.
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
       Block the delete tools for the team

       afs teams put-team-mcp-policy --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --mcp-server-id mcp_01JD2Q3R4S5T6V7W8X9Y0Z0001 --rules '[{"name":"delete_*","effect":"block"}]'

```

## put-team-mcp-server

```text
AFS TEAMS PUT-TEAM-MCP-SERVER

NAME
       teams put-team-mcp-server

DESCRIPTION
       Put team MCP server

       Only a team on its own gateway (per_namespace) can switch a server on. read_only mode admits only tools the review marked read-only.

       Required permission: team-mcp-servers:write

SYNOPSIS
       afs teams put-team-mcp-server [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --mcp-server-id (string) (required)
       --enabled (boolean) (required)
       --mode (string)
          Possible values: full, read_only
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
       Switch a server on, read-only

       afs teams put-team-mcp-server --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --mcp-server-id mcp_01JD2Q3R4S5T6V7W8X9Y0Z0001 --enabled --mode read_only

       Switch it off

       afs teams put-team-mcp-server --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --mcp-server-id mcp_01JD2Q3R4S5T6V7W8X9Y0Z0001 --no-enabled

```

## put-team-mcp-server-credential

```text
AFS TEAMS PUT-TEAM-MCP-SERVER-CREDENTIAL

NAME
       teams put-team-mcp-server-credential

DESCRIPTION
       Put team MCP server credential

       --secret-value (a static credential) and --client-secret (an OAuth client secret) are read from file://PATH, fileb://PATH or - (standard input); they are never plain arguments, and the API never returns them (only the credential's state, subject and expiry).

       Required permission: team-mcp-servers:write

SYNOPSIS
       afs teams put-team-mcp-server-credential [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --mcp-server-id (string) (required)
       --client-id (string)
       --client-secret (file://PATH | fileb://PATH | -; prompted without echo when omitted on a terminal)
          An OAuth client secret; never returned
          Never pass the value itself on the command line.
       --secret-value (file://PATH | fileb://PATH | -; prompted without echo when omitted on a terminal)
          A static credential; never returned
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
       Set the team's token for a per_team server, asked for without echo

       afs teams put-team-mcp-server-credential --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --mcp-server-id mcp_01JD2Q3R4S5T6V7W8X9Y0Z0001 --secret-value -

```

## put-team-model-endpoints

```text
AFS TEAMS PUT-TEAM-MODEL-ENDPOINTS

NAME
       teams put-team-model-endpoints

DESCRIPTION
       Put team model endpoints

       Replaces the team's opted-in service endpoints: an endpoint left out is deselected. Endpoints that need extra terms, or whose provider is selected by service admins, need a service admin (needs_service_admin). No If-Match is sent unless --if-match is given.

       Required permission: model-endpoints:write

SYNOPSIS
       afs teams put-team-model-endpoints [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --endpoint-ids (list) (required)
          Service endpoints to opt in to; one left out is deselected
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
       Opt the team in to two service routes

       afs teams put-team-model-endpoints --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --endpoint-ids mep_01JD2Q3R4S5T6V7W8X9Y0Z0001 mep_01JD2Q3R4S5T6V7W8X9Y0Z0003

```

## put-team-models

```text
AFS TEAMS PUT-TEAM-MODELS

NAME
       teams put-team-models

DESCRIPTION
       Put team models

       Only service admins change a team's models (teams:set-guardrails). Narrowing never edits an agent's stored selection: restoring a model restores access (MOD-6).

       Required permission: teams:set-guardrails

SYNOPSIS
       afs teams put-team-models [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --access (string)
          all (every model enabled for the service, now and later) or selected (the list in --model-ids). A body with only --model-ids is read as selected (deprecated form).
          Possible values: all, selected
       --model-ids (list)
          The team's list: 0 to 100 distinct model ids, each enabled for the service. Not accepted with --access all.
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
       Let the team use every model enabled for the service

       afs teams put-team-models --team-id team_01J9ZQ3ABCDEFGHJKMNPQRSTVW --access all

       Limit the team to two models

       afs teams put-team-models --team-id team_01J9ZQ3ABCDEFGHJKMNPQRSTVW --access selected --model-ids anthropic.claude-sonnet-4-5 anthropic.claude-haiku-4-5

```

## put-team-tags

```text
AFS TEAMS PUT-TEAM-TAGS

NAME
       teams put-team-tags

DESCRIPTION
       Put team tags

       Required permission: tags:write

SYNOPSIS
       afs teams put-team-tags [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --tags (key=value pairs or JSON) (required)
          Team-owned keys only
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
       Replace the team-owned tags

       afs teams put-team-tags --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --tags cost-center=CC-2002 project=checkout

```

## refresh-team-mcp-inventory

```text
AFS TEAMS REFRESH-TEAM-MCP-INVENTORY

NAME
       teams refresh-team-mcp-inventory

DESCRIPTION
       Refresh team MCP inventory

       Required permission: team-mcp-servers:write

SYNOPSIS
       afs teams refresh-team-mcp-inventory [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --mcp-server-id (string) (required)
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
       Refresh the team server's inventory

       afs teams refresh-team-mcp-inventory --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --mcp-server-id mcp_01JD2Q3R4S5T6V7W8X9Y0Z0002

```

## reject-team-mcp-primitive

```text
AFS TEAMS REJECT-TEAM-MCP-PRIMITIVE

NAME
       teams reject-team-mcp-primitive

DESCRIPTION
       Reject team MCP primitive

       Required permission: team-mcp-servers:approve

SYNOPSIS
       afs teams reject-team-mcp-primitive [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --mcp-server-id (string) (required)
       --primitive-id (string) (required)
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
       Block a new tool of the team's server

       afs teams reject-team-mcp-primitive --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --mcp-server-id mcp_01JD2Q3R4S5T6V7W8X9Y0Z0002 --primitive-id mpr_01JD2Q3R4S5T6V7W8X9Y0Z0001 --note "Not needed"

```

## reject-team-mcp-server

```text
AFS TEAMS REJECT-TEAM-MCP-SERVER

NAME
       teams reject-team-mcp-server

DESCRIPTION
       Reject team MCP server

       Required permission: team-mcp-servers:approve

SYNOPSIS
       afs teams reject-team-mcp-server [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --mcp-server-id (string) (required)
       --note (string) (required)
       --primitives (JSON)
          Per-primitive decisions as JSON or file://PATH: [{"primitiveId":…, "decision":"approve"|"block"}]; primitives not named keep their state.
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
       Reject the team's server

       afs teams reject-team-mcp-server --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --mcp-server-id mcp_01JD2Q3R4S5T6V7W8X9Y0Z0002 --note "Use the enterprise server"

```

## restore-revision

```text
AFS TEAMS RESTORE-REVISION

NAME
       teams restore-revision

DESCRIPTION
       Restore revision

       Required permission: revisions:restore

SYNOPSIS
       afs teams restore-revision [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --revision-id (integer) (required)
          Range: 1 to unbounded
       --dry-run (boolean)
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
       Show what restoring revision 41 would change

       afs teams restore-revision --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --revision-id 41 --dry-run

       Restore revision 41

       afs teams restore-revision --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --revision-id 41

```

## resume-team

```text
AFS TEAMS RESUME-TEAM

NAME
       teams resume-team

DESCRIPTION
       Resume team

       Required permission: teams:suspend

SYNOPSIS
       afs teams resume-team [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
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
       Resume a suspended team

       afs teams resume-team --team-id team_01J9ZQ3ABCDEFGHJKMNPQRSTVW

```

## set-external-system-state

```text
AFS TEAMS SET-EXTERNAL-SYSTEM-STATE

NAME
       teams set-external-system-state

DESCRIPTION
       Set team external system state

       Required permission: external-systems:write

SYNOPSIS
       afs teams set-external-system-state [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
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
       Enable an external system for the team

       afs teams set-external-system-state --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --system-id ext_01JD2Q3R4S5T6V7W8X9Y0Z0001 --enabled

```

## start-mcp-oauth-connect

```text
AFS TEAMS START-MCP-OAUTH-CONNECT

NAME
       teams start-mcp-oauth-connect

DESCRIPTION
       Start MCP OAuth connect

       Open the authorization URL in a browser before it expires; the service stores the resulting token and the gateway refreshes it. The state is single-use.

       Required permission: team-mcp-servers:write

SYNOPSIS
       afs teams start-mcp-oauth-connect [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --mcp-server-id (string) (required)
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
       Connect the team's account to an oauth_user server

       afs teams start-mcp-oauth-connect --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --mcp-server-id mcp_01JD2Q3R4S5T6V7W8X9Y0Z0001

       Print only the URL to open

       afs teams start-mcp-oauth-connect --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --mcp-server-id mcp_01JD2Q3R4S5T6V7W8X9Y0Z0001 --query authorizationUrl --output text

```

## submit-mcp-server-for-review

```text
AFS TEAMS SUBMIT-MCP-SERVER-FOR-REVIEW

NAME
       teams submit-mcp-server-for-review

DESCRIPTION
       Submit MCP server for review

       Required permission: team-mcp-servers:write

SYNOPSIS
       afs teams submit-mcp-server-for-review [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --mcp-server-id (string) (required)
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
       Submit the team's server for review

       afs teams submit-mcp-server-for-review --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --mcp-server-id mcp_01JD2Q3R4S5T6V7W8X9Y0Z0002

```

## suspend-team

```text
AFS TEAMS SUSPEND-TEAM

NAME
       teams suspend-team

DESCRIPTION
       Suspend team

       Required permission: teams:suspend

SYNOPSIS
       afs teams suspend-team [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
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
       Suspend a team (writes are refused until it is resumed)

       afs teams suspend-team --team-id team_01J9ZQ3ABCDEFGHJKMNPQRSTVW

```

## test-mcp-federation

```text
AFS TEAMS TEST-MCP-FEDERATION

NAME
       teams test-mcp-federation

DESCRIPTION
       Test MCP federation

       Accepted (202): the team gateway exchanges its token once; read federationState afterwards.

       Required permission: team-mcp-servers:write

SYNOPSIS
       afs teams test-mcp-federation [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --mcp-server-id (string) (required)
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
       Test the team gateway's token exchange

       afs teams test-mcp-federation --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --mcp-server-id mcp_01JD2Q3R4S5T6V7W8X9Y0Z0002

```

## test-notification-channel

```text
AFS TEAMS TEST-NOTIFICATION-CHANNEL

NAME
       teams test-notification-channel

DESCRIPTION
       Test notification channel

       Required permission: notifications:write

SYNOPSIS
       afs teams test-notification-channel [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --channel-id (string) (required)
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
       Send a test notification (exits 254 when delivery fails)

       afs teams test-notification-channel --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --channel-id nch_01JD2Q3R4S5T6V7W8X9Y0Z0001

```

## update-external-system

```text
AFS TEAMS UPDATE-EXTERNAL-SYSTEM

NAME
       teams update-external-system

DESCRIPTION
       Update team external system

       Required permission: external-systems:write

SYNOPSIS
       afs teams update-external-system [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --system-id (string) (required)
       --description (string)
       --docs-url (string)
       --endpoints (list)
       --name (string)
       --owner (string)
       --purpose (string)
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
       Change an external system's endpoints

       afs teams update-external-system --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --system-id ext_01JD2Q3R4S5T6V7W8X9Y0Z0001 --endpoints ledger.example.com:443 ledger-eu.example.com:443

```

## update-log-destination

```text
AFS TEAMS UPDATE-LOG-DESTINATION

NAME
       teams update-log-destination

DESCRIPTION
       Update log destination

       Required permission: logs:write

SYNOPSIS
       afs teams update-log-destination [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --destination-id (string) (required)
       --config (JSON)
       --enabled (boolean)
       --name (string)
       --secret-id (string)
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
       Disable a destination

       afs teams update-log-destination --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --destination-id log_01J9ZQ3ABCDEFGHJKMNPQRSTVW --no-enabled

```

## update-notification-channel

```text
AFS TEAMS UPDATE-NOTIFICATION-CHANNEL

NAME
       teams update-notification-channel

DESCRIPTION
       Update notification channel

       Required permission: notifications:write

SYNOPSIS
       afs teams update-notification-channel [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --channel-id (string) (required)
       --config (JSON)
       --enabled (boolean)
       --event-types (list)
       --name (string)
       --secret-id (string)
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
       Disable a channel

       afs teams update-notification-channel --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --channel-id nch_01JD2Q3R4S5T6V7W8X9Y0Z0001 --no-enabled

```

## update-team

```text
AFS TEAMS UPDATE-TEAM

NAME
       teams update-team

DESCRIPTION
       Update team

       Required permission: teams:write

SYNOPSIS
       afs teams update-team [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --allowed-aws-accounts (list)
       --ceiling-override-expires-at (string or JSON)
       --ceiling-override-reason (string)
          Required with an override; stored in the audit record only
       --ceiling-override-usd (string or JSON)
       --cost-center (string)
       --description (string)
       --display-name (string)
       --identity-provider-id (string)
       --isolation-tier (string)
          Possible values: shared, dedicated
       --kind (string)
          Immutable; a different value is a ValidationError
          Possible values: persistent, ephemeral
       --max-budget-usd (string)
          USD amount as a decimal string with at most 6 fraction digits; never a JSON number
       --max-collaboration-exposure (string)
          Possible values: off, private, public
       --may-add-external-systems (boolean)
       --may-expose-collaboration (boolean)
          Deprecated: use --max-collaboration-exposure (off, private or public) instead.
       --may-manage-aws-accounts (boolean)
       --owner-email (string)
       --quotas (JSON)
       --viewer-kubectl (boolean)
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
       Change a team's description (reads the current version, sends If-Match)

       afs teams update-team --team-id team_01J9ZQ3ABCDEFGHJKMNPQRSTVW --description "Payments platform"

       Apply a change only if the team is still at version 7

       afs teams update-team --team-id team_01J9ZQ3ABCDEFGHJKMNPQRSTVW --viewer-kubectl --if-match 7

       Cap how far the team's orchestration agent may be exposed (service admins)

       afs teams update-team --team-id team_01J9ZQ3ABCDEFGHJKMNPQRSTVW --max-collaboration-exposure private

```

## update-team-mcp-server

```text
AFS TEAMS UPDATE-TEAM-MCP-SERVER

NAME
       teams update-team-mcp-server

DESCRIPTION
       Update team MCP server

       Required permission: team-mcp-servers:write

SYNOPSIS
       afs teams update-team-mcp-server [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --mcp-server-id (string) (required)
       --auth-header (string)
       --auth-kind (string)
          Possible values: none, static, oauth_client, workload_federation, oauth_user, aws_sigv4
       --ca-bundle (string or JSON)
          PEM bundle of the CA that signs the server's certificate, as text or file://PATH (not a secret).
       --data-classification (string)
          Possible values: public, internal, confidential, restricted
       --display-name (string)
       --doc-url (string)
       --federation (JSON)
          Workload federation as JSON or file://PATH: the token endpoint and audience the gateway exchanges its projected token with.
       --hosted (JSON)
          A hosted server as JSON or file://PATH: {"imageDigest":…, "args":[…], "envNames":[…], "egressHosts":[…], "instanceMode":"per_team", "needsAws":false}.
       --limits (JSON)
          The server's own limits as JSON: {"requestBytes":…, "responseBytes":…, "timeoutSeconds":…, "concurrentPerAgent":…}; unset members take mcp.limits.
       --location (string)
          Possible values: external, corporate, in_namespace
       --owner (string)
       --purpose (string)
       --transport (string)
          Possible values: streamable_http, hosted
       --url (string)
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
       Point the team's server at a new URL

       afs teams update-team-mcp-server --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --mcp-server-id mcp_01JD2Q3R4S5T6V7W8X9Y0Z0002 --url https://tools-v2.payments.corp.example/mcp

```
