# afs budgets

Generated from `afs budgets help`; do not edit by hand.

```text
AFS BUDGETS

DESCRIPTION
       Set budgets and read spend.

SYNOPSIS
       afs budgets <operation> [parameters]

OPERATIONS
       get-budget                     Get budget
       get-usage                      Get usage
       list-costs                     List costs
       put-budget                     Put budget

To see help text for an operation, run: afs budgets <operation> help
```

## get-budget

```text
AFS BUDGETS GET-BUDGET

NAME
       budgets get-budget

DESCRIPTION
       Get budget

       Required permission: budgets:read

SYNOPSIS
       afs budgets get-budget [options]

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
       afs budgets get-budget help
```

## get-usage

```text
AFS BUDGETS GET-USAGE

NAME
       budgets get-usage

DESCRIPTION
       Get usage

       Required permission: usage:read

SYNOPSIS
       afs budgets get-usage [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --from (string)
       --to (string)
       --group-by (string)
          Possible values: agent, model, day

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
       afs budgets get-usage help
```

## list-costs

```text
AFS BUDGETS LIST-COSTS

NAME
       budgets list-costs

DESCRIPTION
       List costs

       Required permission: usage:read

SYNOPSIS
       afs budgets list-costs [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --from (string)
       --to (string)

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
       afs budgets list-costs help
```

## put-budget

```text
AFS BUDGETS PUT-BUDGET

NAME
       budgets put-budget

DESCRIPTION
       Put budget

       Required permission: budgets:write

SYNOPSIS
       afs budgets put-budget [options]

OPTIONS
       --team-id (string) (required)
          Team id; taken from the path only
       --amount-usd (string) (required)
          USD amount as a decimal string with at most 6 fraction digits; never a JSON number
       --expires-at (string or JSON)
       --on-exhaust (string)
          notify needs teams:set-guardrails and a persistent team (D33)
          Possible values: block, notify
       --period (string) (required)
          Possible values: monthly, lifetime
       --thresholds (list)
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
       afs budgets put-budget help
```
