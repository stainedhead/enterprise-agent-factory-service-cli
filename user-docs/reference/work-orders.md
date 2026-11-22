# afs work-orders

Generated from `afs work-orders help`; do not edit by hand.

```text
AFS WORK-ORDERS

DESCRIPTION
       Manage work orders and their repositories.

SYNOPSIS
       afs work-orders <operation> [parameters]

OPERATIONS
       cancel-work-order              Cancel work order
       create-work-order              Create work order
       describe-work-order            Get work order
       list-repos                     List work order repos
       list-team-work-orders          List team work orders
       list-work-orders               List work orders
       update-repo                    Update work order repo
       update-work-order              Update work order

To see help text for an operation, run: afs work-orders <operation> help
```

## cancel-work-order

```text
AFS WORK-ORDERS CANCEL-WORK-ORDER

NAME
       work-orders cancel-work-order

DESCRIPTION
       Cancel work order

       Required permission: work-orders:write

SYNOPSIS
       afs work-orders cancel-work-order [options]

OPTIONS
       --work-order-id (string) (required)
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
       Cancel a work order

       afs work-orders cancel-work-order --work-order-id wo_01J9ZQ3ABCDEFGHJKMNPQRSTVW --reason "CVE withdrawn upstream"

```

## create-work-order

```text
AFS WORK-ORDERS CREATE-WORK-ORDER

NAME
       work-orders create-work-order

DESCRIPTION
       Create work order

       Required permission: work-orders:write

SYNOPSIS
       afs work-orders create-work-order [options]

OPTIONS
       --budget-usd (string)
          USD amount as a decimal string with at most 6 fraction digits; never a JSON number
       --cve-id (string)
       --kind (string)
          Possible values: cve_remediation, general
       --repos (list) (required)
       --severity (string) (required)
          Possible values: critical, high, medium, low
       --sla-due-at (string) (required)
          RFC 3339 UTC with Z
       --template-id (string) (required)
       --title (string) (required)
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
       Open a CVE work order; the team configuration decides who owns the first repository

       afs work-orders create-work-order --template-id swt_01J9ZQ3ABCDEFGHJKMNPQRSTVW --title "Patch CVE-2026-1234" --kind cve_remediation --cve-id CVE-2026-1234 --severity critical --sla-due-at 2026-12-01T00:00:00Z --repos github.com/acme/ledger github.com/acme/billing=team_01J9ZQ3ABCDEFGHJKMNPQRSTVX

       Name an owning contact by email; retry safely with an idempotency key

       afs work-orders create-work-order --template-id swt_01J9ZQ3ABCDEFGHJKMNPQRSTVW --title "Upgrade the runtime" --severity medium --sla-due-at 2026-12-15T00:00:00Z --repos github.com/acme/ledger=ledger-team@example.com --idempotency-key upgrade-2026-12

```

## describe-work-order

```text
AFS WORK-ORDERS DESCRIBE-WORK-ORDER

NAME
       work-orders describe-work-order

DESCRIPTION
       Get work order

       Required permission: work-orders:read

SYNOPSIS
       afs work-orders describe-work-order [options]

OPTIONS
       --work-order-id (string) (required)

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
       Show a work order and its repositories

       afs work-orders describe-work-order --work-order-id wo_01J9ZQ3ABCDEFGHJKMNPQRSTVW --output json

```

## list-repos

```text
AFS WORK-ORDERS LIST-REPOS

NAME
       work-orders list-repos

DESCRIPTION
       List work order repos

       Required permission: work-orders:read

SYNOPSIS
       afs work-orders list-repos [options]

OPTIONS
       --work-order-id (string) (required)

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
       List the repositories of a work order

       afs work-orders list-repos --work-order-id wo_01J9ZQ3ABCDEFGHJKMNPQRSTVW

```

## list-team-work-orders

```text
AFS WORK-ORDERS LIST-TEAM-WORK-ORDERS

NAME
       work-orders list-team-work-orders

DESCRIPTION
       List team work orders

       Required permission: work-orders:read

SYNOPSIS
       afs work-orders list-team-work-orders [options]

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
       List a team's work orders

       afs work-orders list-team-work-orders --team-id team_01J9ZQ3ABCDEFGHJKMNPQRSTVW --output json

```

## list-work-orders

```text
AFS WORK-ORDERS LIST-WORK-ORDERS

NAME
       work-orders list-work-orders

DESCRIPTION
       List work orders

       Required permission: work-orders:read

SYNOPSIS
       afs work-orders list-work-orders [options]

OPTIONS
       --status (string)
          Possible values: open, assigned, in_progress, pr_opened, done, failed, cancelled
       --severity (string)
          Possible values: critical, high, medium, low
       --sla-state (string)
          Possible values: on_track, due_soon, breached, met, missed
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
       List open critical work orders

       afs work-orders list-work-orders --status open --severity critical

       Export every work order as CSV

       afs work-orders list-work-orders --output csv

```

## update-repo

```text
AFS WORK-ORDERS UPDATE-REPO

NAME
       work-orders update-repo

DESCRIPTION
       Update work order repo

       Required permission: work-orders:report

SYNOPSIS
       afs work-orders update-repo [options]

OPTIONS
       --work-order-id (string) (required)
       --host (string)
       --owner (string)
       --repo (string) (required)
       --detail (string)
       --pr-number (integer)
          Range: 1 to unbounded
       --pr-url (string)
       --status (string) (required)
          Possible values: pending, in_progress, pr_opened, merged, not_affected, failed, skipped
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
       Report a pull request for one repository

       afs work-orders update-repo --work-order-id wo_01J9ZQ3ABCDEFGHJKMNPQRSTVW --repo github.com/acme/ledger --status pr_opened --pr-url https://github.com/acme/ledger/pull/42 --pr-number 42

       Report with the version printed by the previous update (a workload token cannot read it)

       afs work-orders update-repo --work-order-id wo_01J9ZQ3ABCDEFGHJKMNPQRSTVW --repo github.com/acme/ledger --status merged --if-match 3

```

## update-work-order

```text
AFS WORK-ORDERS UPDATE-WORK-ORDER

NAME
       work-orders update-work-order

DESCRIPTION
       Update work order

       Required permission: work-orders:write

SYNOPSIS
       afs work-orders update-work-order [options]

OPTIONS
       --work-order-id (string) (required)
       --add-repos (list)
       --budget-usd (string)
          USD amount as a decimal string with at most 6 fraction digits; never a JSON number
       --sla-due-at (string)
          RFC 3339 UTC with Z
       --title (string)
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
       Move the SLA date and add a repository

       afs work-orders update-work-order --work-order-id wo_01J9ZQ3ABCDEFGHJKMNPQRSTVW --sla-due-at 2026-12-20T00:00:00Z --add-repos github.com/acme/gateway=team_01J9ZQ3ABCDEFGHJKMNPQRSTVX

```
