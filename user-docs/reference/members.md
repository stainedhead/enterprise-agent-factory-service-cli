# afs members

Generated from `afs members help`; do not edit by hand.

```text
AFS MEMBERS

DESCRIPTION
       Manage team members and role bindings.

SYNOPSIS
       afs members <operation> [parameters]

OPERATIONS
       add-member                     Add member
       list-members                   List members
       remove-member                  Remove member
       update-member                  Update member

To see help text for an operation, run: afs members <operation> help
```

## add-member

```text
AFS MEMBERS ADD-MEMBER

NAME
       members add-member

DESCRIPTION
       Add member

       Required permission: members:write

SYNOPSIS
       afs members add-member [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --email (string)
       --group-name (string)
       --identity-provider-id (string)
       --role (string) (required)
          Possible values: team_admin, team_collaborator, team_viewer
       --subject-type (string) (required)
          Possible values: user, group
       --user-id (string)
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
       Add a user as team_admin

       afs members add-member --team-id team_01J9ZQ3ABCDEFGHJKMNPQRSTVW --subject-type user --email ada@example.com --role team_admin

       Bind an identity provider group

       afs members add-member --team-id team_01J9ZQ3ABCDEFGHJKMNPQRSTVW --subject-type group --group-name payments-devs --role team_viewer

```

## list-members

```text
AFS MEMBERS LIST-MEMBERS

NAME
       members list-members

DESCRIPTION
       List members

       Required permission: members:read

SYNOPSIS
       afs members list-members [options]

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
       List a team's members and groups

       afs members list-members --team-id team_01J9ZQ3ABCDEFGHJKMNPQRSTVW

```

## remove-member

```text
AFS MEMBERS REMOVE-MEMBER

NAME
       members remove-member

DESCRIPTION
       Remove member

       Required permission: members:write

SYNOPSIS
       afs members remove-member [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --binding-id (string) (required)
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
       Remove a member

       afs members remove-member --team-id team_01J9ZQ3ABCDEFGHJKMNPQRSTVW --binding-id mem_01J9ZQ3ABCDEFGHJKMNPQRSTVW

```

## update-member

```text
AFS MEMBERS UPDATE-MEMBER

NAME
       members update-member

DESCRIPTION
       Update member

       Required permission: members:write

SYNOPSIS
       afs members update-member [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --binding-id (string) (required)
       --role (string) (required)
          Possible values: team_admin, team_collaborator, team_viewer
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
       Change a member's role (reads the version from the member list)

       afs members update-member --team-id team_01J9ZQ3ABCDEFGHJKMNPQRSTVW --binding-id mem_01J9ZQ3ABCDEFGHJKMNPQRSTVW --role team_viewer

```
