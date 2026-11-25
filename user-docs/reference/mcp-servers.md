# afs mcp-servers

Generated from `afs mcp-servers help`; do not edit by hand.

```text
AFS MCP-SERVERS

DESCRIPTION
       Govern the MCP server catalog: register, review, approve, suspend and capture search queries.

SYNOPSIS
       afs mcp-servers <operation> [parameters]

OPERATIONS
       approve-mcp-primitive          Approve MCP primitive
       approve-mcp-server             Approve MCP server
       create-mcp-server              Create MCP server
       delete-mcp-server              Delete MCP server
       describe-mcp-server            Get MCP server
       get-mcp-primitive-change       Get MCP primitive change
       get-mcp-query-capture          Get MCP query capture
       import-mcp-server-manifest     Import MCP server manifest
       list-mcp-primitives            List MCP primitives
       list-mcp-reviews               List MCP reviews
       list-mcp-servers               List MCP servers
       publish-mcp-server             Publish MCP server
       put-mcp-policy                 Put MCP policy
       put-mcp-server-credential      Put MCP server credential
       refresh-mcp-inventory          Refresh MCP inventory
       reject-mcp-primitive           Reject MCP primitive
       reject-mcp-server              Reject MCP server
       request-mcp-server-changes     Request MCP server changes
       resume-mcp-server              Resume MCP server
       retire-mcp-server              Retire MCP server
       start-mcp-query-capture        Start MCP query capture
       stop-mcp-query-capture         Stop MCP query capture
       suspend-mcp-server             Suspend MCP server
       update-mcp-server              Update MCP server

To see help text for an operation, run: afs mcp-servers <operation> help
```

## approve-mcp-primitive

```text
AFS MCP-SERVERS APPROVE-MCP-PRIMITIVE

NAME
       mcp-servers approve-mcp-primitive

DESCRIPTION
       Approve MCP primitive

       Required permission: mcp-servers:review

SYNOPSIS
       afs mcp-servers approve-mcp-primitive [options]

OPTIONS
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
       Approve a changed tool definition

       afs mcp-servers approve-mcp-primitive --mcp-server-id mcp_01JD2Q3R4S5T6V7W8X9Y0Z0001 --primitive-id mpr_01JD2Q3R4S5T6V7W8X9Y0Z0001 --note "Description only"

```

## approve-mcp-server

```text
AFS MCP-SERVERS APPROVE-MCP-SERVER

NAME
       mcp-servers approve-mcp-server

DESCRIPTION
       Approve MCP server

       The approval level follows the server's risk tier and mcp.approval_policy: a team admin, a service admin or a security operator. A caller below the level gets AccessDenied with reason approval_not_permitted_by_policy.

       Required permission: mcp-servers:review

SYNOPSIS
       afs mcp-servers approve-mcp-server [options]

OPTIONS
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
       Approve a server

       afs mcp-servers approve-mcp-server --mcp-server-id mcp_01JD2Q3R4S5T6V7W8X9Y0Z0001 --note "Reviewed with the owning team"

       Approve it and block one tool

       afs mcp-servers approve-mcp-server --mcp-server-id mcp_01JD2Q3R4S5T6V7W8X9Y0Z0001 --note "Approved without delete" --primitives '[{"primitiveId":"mpr_01JD2Q3R4S5T6V7W8X9Y0Z0001","decision":"block"}]'

```

## create-mcp-server

```text
AFS MCP-SERVERS CREATE-MCP-SERVER

NAME
       mcp-servers create-mcp-server

DESCRIPTION
       Create MCP server

       The server is created as a draft. Set its credential, refresh its inventory, then approve it; teams see it once it is approved and published.

       Required permission: mcp-servers:write

SYNOPSIS
       afs mcp-servers create-mcp-server [options]

OPTIONS
       --allow-sampling (boolean)
       --auth-header (string)
       --auth-kind (string) (required)
          Possible values: none, static, oauth_client, workload_federation, oauth_user, aws_sigv4
       --ca-bundle (string or JSON)
          PEM bundle of the CA that signs the server's certificate, as text or file://PATH (not a secret).
       --credential-mode (string)
          Possible values: shared, per_team
       --data-classification (string) (required)
          Possible values: public, internal, confidential, restricted
       --display-name (string) (required)
       --doc-url (string)
       --drift-policy (string)
          Possible values: hold, auto_approve_compatible
       --federation (JSON)
          Workload federation as JSON or file://PATH: the token endpoint and audience the gateway exchanges its projected token with.
       --hosted (JSON)
          A hosted server as JSON or file://PATH: {"imageDigest":…, "args":[…], "envNames":[…], "egressHosts":[…], "instanceMode":"per_team", "needsAws":false}.
       --limits (JSON)
          The server's own limits as JSON: {"requestBytes":…, "responseBytes":…, "timeoutSeconds":…, "concurrentPerAgent":…}; unset members take mcp.limits.
       --location (string) (required)
          Possible values: external, corporate, in_namespace
       --owner (string)
       --purpose (string)
       --slug (string) (required)
       --transport (string) (required)
          Possible values: streamable_http, hosted
       --trust-annotations (boolean)
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
       Register a remote server with a static credential

       afs mcp-servers create-mcp-server --slug github --display-name "GitHub" --owner platform-team --purpose "Issues and pull requests" --data-classification internal --transport streamable_http --location external --url https://api.githubcopilot.com/mcp/ --auth-kind static --credential-mode per_team

       Register a hosted stdio server

       afs mcp-servers create-mcp-server --slug jira --display-name "Jira" --data-classification confidential --transport hosted --location in_namespace --auth-kind static --hosted file://jira-hosted.json

```

## delete-mcp-server

```text
AFS MCP-SERVERS DELETE-MCP-SERVER

NAME
       mcp-servers delete-mcp-server

DESCRIPTION
       Delete MCP server

       Only a draft or rejected server can be deleted; retire an approved one instead.

       Required permission: mcp-servers:write

SYNOPSIS
       afs mcp-servers delete-mcp-server [options]

OPTIONS
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
       Delete a draft

       afs mcp-servers delete-mcp-server --mcp-server-id mcp_01JD2Q3R4S5T6V7W8X9Y0Z0001

```

## describe-mcp-server

```text
AFS MCP-SERVERS DESCRIBE-MCP-SERVER

NAME
       mcp-servers describe-mcp-server

DESCRIPTION
       Get MCP server

       Required permission: mcp-servers:read

SYNOPSIS
       afs mcp-servers describe-mcp-server [options]

OPTIONS
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
       Show a server

       afs mcp-servers describe-mcp-server --mcp-server-id mcp_01JD2Q3R4S5T6V7W8X9Y0Z0001

       Show why it is high risk

       afs mcp-servers describe-mcp-server --mcp-server-id mcp_01JD2Q3R4S5T6V7W8X9Y0Z0001 --query riskReasons --output text

```

## get-mcp-primitive-change

```text
AFS MCP-SERVERS GET-MCP-PRIMITIVE-CHANGE

NAME
       mcp-servers get-mcp-primitive-change

DESCRIPTION
       Get MCP primitive change

       Required permission: mcp-servers:read

SYNOPSIS
       afs mcp-servers get-mcp-primitive-change [options]

OPTIONS
       --mcp-server-id (string) (required)
       --primitive-id (string) (required)

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
       Show what changed in a held definition

       afs mcp-servers get-mcp-primitive-change --mcp-server-id mcp_01JD2Q3R4S5T6V7W8X9Y0Z0001 --primitive-id mpr_01JD2Q3R4S5T6V7W8X9Y0Z0001

       Show only the changed fields

       afs mcp-servers get-mcp-primitive-change --mcp-server-id mcp_01JD2Q3R4S5T6V7W8X9Y0Z0001 --primitive-id mpr_01JD2Q3R4S5T6V7W8X9Y0Z0001 --query diff

```

## get-mcp-query-capture

```text
AFS MCP-SERVERS GET-MCP-QUERY-CAPTURE

NAME
       mcp-servers get-mcp-query-capture

DESCRIPTION
       Get MCP query capture

       Required permission: mcp-events:read

SYNOPSIS
       afs mcp-servers get-mcp-query-capture [options]

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
       Show whether search queries are being captured

       afs mcp-servers get-mcp-query-capture

```

## import-mcp-server-manifest

```text
AFS MCP-SERVERS IMPORT-MCP-SERVER-MANIFEST

NAME
       mcp-servers import-mcp-server-manifest

DESCRIPTION
       Import MCP server manifest

       Required permission: mcp-servers:write

SYNOPSIS
       afs mcp-servers import-mcp-server-manifest [options]

OPTIONS
       --manifest (JSON) (required)
          An MCP Registry server.json (preview format).
       --slug (string)
          Overrides the slug derived from the manifest name
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
       Register a server from its server.json manifest

       afs mcp-servers import-mcp-server-manifest --manifest file://server.json --slug github

```

## list-mcp-primitives

```text
AFS MCP-SERVERS LIST-MCP-PRIMITIVES

NAME
       mcp-servers list-mcp-primitives

DESCRIPTION
       List MCP primitives

       Required permission: mcp-servers:read

SYNOPSIS
       afs mcp-servers list-mcp-primitives [options]

OPTIONS
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
       List a server's tools

       afs mcp-servers list-mcp-primitives --mcp-server-id mcp_01JD2Q3R4S5T6V7W8X9Y0Z0001 --kind tool

       List the definitions held for review

       afs mcp-servers list-mcp-primitives --mcp-server-id mcp_01JD2Q3R4S5T6V7W8X9Y0Z0001 --status pending_review

```

## list-mcp-reviews

```text
AFS MCP-SERVERS LIST-MCP-REVIEWS

NAME
       mcp-servers list-mcp-reviews

DESCRIPTION
       List MCP reviews

       Required permission: mcp-servers:review

SYNOPSIS
       afs mcp-servers list-mcp-reviews [options]

OPTIONS
       --kind (string)
          Possible values: server, primitive
       --risk-tier (string)
          Possible values: low, elevated, high
       --team-id (string)
       --approver (string)
          Possible values: security, team_admin

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
       List everything that waits for a security approval

       afs mcp-servers list-mcp-reviews --approver security

       List held definition changes of high-risk servers

       afs mcp-servers list-mcp-reviews --kind primitive --risk-tier high

```

## list-mcp-servers

```text
AFS MCP-SERVERS LIST-MCP-SERVERS

NAME
       mcp-servers list-mcp-servers

DESCRIPTION
       List MCP servers

       Required permission: mcp-servers:read

SYNOPSIS
       afs mcp-servers list-mcp-servers [options]

OPTIONS
       --scope (string)
          Possible values: enterprise, team, all
       --review-status (string)
          Possible values: draft, pending_review, approved, rejected, suspended
       --risk-tier (string)
          Possible values: low, elevated, high
       --team-id (string)
       --location (string)
          Possible values: external, corporate, in_namespace

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
       List the catalog

       afs mcp-servers list-mcp-servers

       List the enterprise servers that wait for review

       afs mcp-servers list-mcp-servers --scope enterprise --review-status pending_review

       List one team's own servers

       afs mcp-servers list-mcp-servers --scope team --team-id team_01JD2Q3R4S5T6V7W8X9Y0Z0001

```

## publish-mcp-server

```text
AFS MCP-SERVERS PUBLISH-MCP-SERVER

NAME
       mcp-servers publish-mcp-server

DESCRIPTION
       Publish MCP server

       Required permission: mcp-servers:write

SYNOPSIS
       afs mcp-servers publish-mcp-server [options]

OPTIONS
       --mcp-server-id (string) (required)
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
       Publish an approved server to every team's catalog

       afs mcp-servers publish-mcp-server --mcp-server-id mcp_01JD2Q3R4S5T6V7W8X9Y0Z0001 --note "Available to all teams"

```

## put-mcp-policy

```text
AFS MCP-SERVERS PUT-MCP-POLICY

NAME
       mcp-servers put-mcp-policy

DESCRIPTION
       Put MCP policy

       Required permission: mcp-servers:write

SYNOPSIS
       afs mcp-servers put-mcp-policy [options]

OPTIONS
       --mcp-server-id (string) (required)
       --allow-sampling (boolean)
       --drift-policy (string)
          Possible values: hold, auto_approve_compatible
       --rules (JSON) (required)
          Rules as JSON or file://PATH: [{"name":"delete_*","effect":"block"}, {"class":"destructive","effect":"block"}, {"name":"search","effect":"allow","rateLimitPerMinute":60}]. The most specific rule of a level wins; a block at any level wins over the levels below.
       --trust-annotations (boolean)
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
       Block every destructive tool of a server

       afs mcp-servers put-mcp-policy --mcp-server-id mcp_01JD2Q3R4S5T6V7W8X9Y0Z0001 --rules '[{"class":"destructive","effect":"block"}]'

       Replace the service policy from a file and trust the server's annotations

       afs mcp-servers put-mcp-policy --mcp-server-id mcp_01JD2Q3R4S5T6V7W8X9Y0Z0001 --rules file://github-policy.json --trust-annotations

```

## put-mcp-server-credential

```text
AFS MCP-SERVERS PUT-MCP-SERVER-CREDENTIAL

NAME
       mcp-servers put-mcp-server-credential

DESCRIPTION
       Put MCP server credential

       --scope is _shared for a shared credential, or a team id for a per_team server.

       --secret-value (a static credential) and --client-secret (an OAuth client secret) are read from file://PATH, fileb://PATH or - (standard input); they are never plain arguments, and the API never returns them (only the credential's state, subject and expiry).

       Required permission: mcp-servers:write

SYNOPSIS
       afs mcp-servers put-mcp-server-credential [options]

OPTIONS
       --mcp-server-id (string) (required)
       --scope (string) (required)
          _shared or a team id
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
       Set a team's token from a file

       afs mcp-servers put-mcp-server-credential --mcp-server-id mcp_01JD2Q3R4S5T6V7W8X9Y0Z0001 --scope team_01JD2Q3R4S5T6V7W8X9Y0Z0001 --secret-value file://github-payments.token

       Set the shared OAuth client from standard input

       afs mcp-servers put-mcp-server-credential --mcp-server-id mcp_01JD2Q3R4S5T6V7W8X9Y0Z0001 --scope _shared --client-id afs-gateway --client-secret -

```

## refresh-mcp-inventory

```text
AFS MCP-SERVERS REFRESH-MCP-INVENTORY

NAME
       mcp-servers refresh-mcp-inventory

DESCRIPTION
       Refresh MCP inventory

       Accepted (202): the gateway lists the server's tools, resources and prompts; new and changed definitions are held for review.

       Required permission: mcp-servers:write

SYNOPSIS
       afs mcp-servers refresh-mcp-inventory [options]

OPTIONS
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
       Refresh a server's inventory now

       afs mcp-servers refresh-mcp-inventory --mcp-server-id mcp_01JD2Q3R4S5T6V7W8X9Y0Z0001

```

## reject-mcp-primitive

```text
AFS MCP-SERVERS REJECT-MCP-PRIMITIVE

NAME
       mcp-servers reject-mcp-primitive

DESCRIPTION
       Reject MCP primitive

       Required permission: mcp-servers:review

SYNOPSIS
       afs mcp-servers reject-mcp-primitive [options]

OPTIONS
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
       Block a new tool

       afs mcp-servers reject-mcp-primitive --mcp-server-id mcp_01JD2Q3R4S5T6V7W8X9Y0Z0001 --primitive-id mpr_01JD2Q3R4S5T6V7W8X9Y0Z0001 --note "Writes to production"

```

## reject-mcp-server

```text
AFS MCP-SERVERS REJECT-MCP-SERVER

NAME
       mcp-servers reject-mcp-server

DESCRIPTION
       Reject MCP server

       The approval level follows the server's risk tier and mcp.approval_policy: a team admin, a service admin or a security operator. A caller below the level gets AccessDenied with reason approval_not_permitted_by_policy.

       Required permission: mcp-servers:review

SYNOPSIS
       afs mcp-servers reject-mcp-server [options]

OPTIONS
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
       Reject a server

       afs mcp-servers reject-mcp-server --mcp-server-id mcp_01JD2Q3R4S5T6V7W8X9Y0Z0001 --note "Sends data to an unapproved region"

```

## request-mcp-server-changes

```text
AFS MCP-SERVERS REQUEST-MCP-SERVER-CHANGES

NAME
       mcp-servers request-mcp-server-changes

DESCRIPTION
       Request MCP server changes

       Required permission: mcp-servers:review

SYNOPSIS
       afs mcp-servers request-mcp-server-changes [options]

OPTIONS
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
       Send a server back to its owner

       afs mcp-servers request-mcp-server-changes --mcp-server-id mcp_01JD2Q3R4S5T6V7W8X9Y0Z0001 --note "Pin the image by digest"

```

## resume-mcp-server

```text
AFS MCP-SERVERS RESUME-MCP-SERVER

NAME
       mcp-servers resume-mcp-server

DESCRIPTION
       Resume MCP server

       Required permission: mcp-servers:suspend

SYNOPSIS
       afs mcp-servers resume-mcp-server [options]

OPTIONS
       --mcp-server-id (string) (required)
       --reason (string) (required)
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
       Resume a suspended server

       afs mcp-servers resume-mcp-server --mcp-server-id mcp_01JD2Q3R4S5T6V7W8X9Y0Z0001 --reason "Vendor incident closed"

```

## retire-mcp-server

```text
AFS MCP-SERVERS RETIRE-MCP-SERVER

NAME
       mcp-servers retire-mcp-server

DESCRIPTION
       Retire MCP server

       Required permission: mcp-servers:write

SYNOPSIS
       afs mcp-servers retire-mcp-server [options]

OPTIONS
       --mcp-server-id (string) (required)
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
       Retire a server teams no longer use

       afs mcp-servers retire-mcp-server --mcp-server-id mcp_01JD2Q3R4S5T6V7W8X9Y0Z0001 --note "Replaced by github-v2"

```

## start-mcp-query-capture

```text
AFS MCP-SERVERS START-MCP-QUERY-CAPTURE

NAME
       mcp-servers start-mcp-query-capture

DESCRIPTION
       Start MCP query capture

       While a capture runs the ledger keeps the text of search queries; it stops by itself at the end of the window.

       Required permission: mcp-settings:write

SYNOPSIS
       afs mcp-servers start-mcp-query-capture [options]

OPTIONS
       --days (integer) (required)
          At most mcp.query_capture_max_days (400 capture_window_too_long)
          Range: 1 to 14
       --reason (string) (required)
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
       Capture search queries for a week

       afs mcp-servers start-mcp-query-capture --days 7 --reason "Tune the search ranking"

```

## stop-mcp-query-capture

```text
AFS MCP-SERVERS STOP-MCP-QUERY-CAPTURE

NAME
       mcp-servers stop-mcp-query-capture

DESCRIPTION
       Stop MCP query capture

       Required permission: mcp-settings:write

SYNOPSIS
       afs mcp-servers stop-mcp-query-capture [options]

OPTIONS
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
       Stop the capture now

       afs mcp-servers stop-mcp-query-capture

```

## suspend-mcp-server

```text
AFS MCP-SERVERS SUSPEND-MCP-SERVER

NAME
       mcp-servers suspend-mcp-server

DESCRIPTION
       Suspend MCP server

       Every team gateway refuses calls to a suspended server within a minute; agents keep their configuration.

       Required permission: mcp-servers:suspend

SYNOPSIS
       afs mcp-servers suspend-mcp-server [options]

OPTIONS
       --mcp-server-id (string) (required)
       --reason (string) (required)
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
       Suspend a server during an incident

       afs mcp-servers suspend-mcp-server --mcp-server-id mcp_01JD2Q3R4S5T6V7W8X9Y0Z0001 --reason "Vendor incident INC-1234"

```

## update-mcp-server

```text
AFS MCP-SERVERS UPDATE-MCP-SERVER

NAME
       mcp-servers update-mcp-server

DESCRIPTION
       Update MCP server

       Required permission: mcp-servers:write

SYNOPSIS
       afs mcp-servers update-mcp-server [options]

OPTIONS
       --mcp-server-id (string) (required)
       --allow-sampling (boolean)
       --auth-header (string)
       --auth-kind (string)
          Possible values: none, static, oauth_client, workload_federation, oauth_user, aws_sigv4
       --ca-bundle (string or JSON)
          PEM bundle of the CA that signs the server's certificate, as text or file://PATH (not a secret).
       --credential-mode (string)
          Possible values: shared, per_team
       --data-classification (string)
          Possible values: public, internal, confidential, restricted
       --default-for-new-teams (boolean)
       --display-name (string)
       --doc-url (string)
       --drift-policy (string)
          Possible values: hold, auto_approve_compatible
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
       --trust-annotations (boolean)
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
       Approve compatible definition changes automatically

       afs mcp-servers update-mcp-server --mcp-server-id mcp_01JD2Q3R4S5T6V7W8X9Y0Z0001 --drift-policy auto_approve_compatible

       Raise the timeout

       afs mcp-servers update-mcp-server --mcp-server-id mcp_01JD2Q3R4S5T6V7W8X9Y0Z0001 --limits '{"timeoutSeconds":120}'

```
