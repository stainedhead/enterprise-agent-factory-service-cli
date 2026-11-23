# afs swarm-templates

Generated from `afs swarm-templates help`; do not edit by hand.

```text
AFS SWARM-TEMPLATES

DESCRIPTION
       Manage swarm templates.

SYNOPSIS
       afs swarm-templates <operation> [parameters]

OPERATIONS
       create-swarm-template          Create swarm template
       delete-swarm-template          Delete swarm template
       describe-swarm-template        Get swarm template
       disable-swarm-template         Disable swarm template
       enable-swarm-template          Enable swarm template
       list-swarm-templates           List swarm templates
       update-swarm-template          Update swarm template

To see help text for an operation, run: afs swarm-templates <operation> help
```

## create-swarm-template

```text
AFS SWARM-TEMPLATES CREATE-SWARM-TEMPLATE

NAME
       swarm-templates create-swarm-template

DESCRIPTION
       Create swarm template

       Required permission: admin:swarm-templates

SYNOPSIS
       afs swarm-templates create-swarm-template [options]

OPTIONS
       --budget-usd (string) (required)
          USD amount as a decimal string with at most 6 fraction digits; never a JSON number
       --description (string)
       --enabled (boolean)
       --kind (string) (required)
          Possible values: cve_remediation, general
       --name (string) (required)
       --spec (JSON) (required)
       --teardown-rule (string) (required)
          Possible values: at_expiry, at_cap_or_expiry, at_work_complete
       --ttl-hours (integer) (required)
          Range: 1 to 8760
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
       Choose the tool and exposure of the remediation team's orchestration agent (spec.orchestration in the spec file)

       afs swarm-templates create-swarm-template --name cve-orchestrated --kind cve_remediation --budget-usd 250.00 --ttl-hours 72 --teardown-rule at_work_complete --spec file://swarm-spec-with-orchestration.json

       Create a CVE remediation template from a spec file

       afs swarm-templates create-swarm-template --name cve-remediation --kind cve_remediation --budget-usd 250.00 --ttl-hours 72 --teardown-rule at_work_complete --spec file://swarm-spec.json

```

## delete-swarm-template

```text
AFS SWARM-TEMPLATES DELETE-SWARM-TEMPLATE

NAME
       swarm-templates delete-swarm-template

DESCRIPTION
       Delete swarm template

       Required permission: admin:swarm-templates

SYNOPSIS
       afs swarm-templates delete-swarm-template [options]

OPTIONS
       --template-id (string) (required)
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
       Delete a template that no open work order uses

       afs swarm-templates delete-swarm-template --template-id swt_01J9ZQ3ABCDEFGHJKMNPQRSTVW

```

## describe-swarm-template

```text
AFS SWARM-TEMPLATES DESCRIBE-SWARM-TEMPLATE

NAME
       swarm-templates describe-swarm-template

DESCRIPTION
       Get swarm template

       Required permission: swarm-templates:read

SYNOPSIS
       afs swarm-templates describe-swarm-template [options]

OPTIONS
       --template-id (string) (required)

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
       Show a swarm template with its spec

       afs swarm-templates describe-swarm-template --template-id swt_01J9ZQ3ABCDEFGHJKMNPQRSTVW --output json

```

## disable-swarm-template

```text
AFS SWARM-TEMPLATES DISABLE-SWARM-TEMPLATE

NAME
       swarm-templates disable-swarm-template

DESCRIPTION
       Disable swarm template

       Required permission: admin:swarm-templates

SYNOPSIS
       afs swarm-templates disable-swarm-template [options]

OPTIONS
       --template-id (string) (required)
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
       Stop new work orders from using a template

       afs swarm-templates disable-swarm-template --template-id swt_01J9ZQ3ABCDEFGHJKMNPQRSTVW

```

## enable-swarm-template

```text
AFS SWARM-TEMPLATES ENABLE-SWARM-TEMPLATE

NAME
       swarm-templates enable-swarm-template

DESCRIPTION
       Enable swarm template

       Required permission: admin:swarm-templates

SYNOPSIS
       afs swarm-templates enable-swarm-template [options]

OPTIONS
       --template-id (string) (required)
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
       Allow new work orders to use a template

       afs swarm-templates enable-swarm-template --template-id swt_01J9ZQ3ABCDEFGHJKMNPQRSTVW

```

## list-swarm-templates

```text
AFS SWARM-TEMPLATES LIST-SWARM-TEMPLATES

NAME
       swarm-templates list-swarm-templates

DESCRIPTION
       List swarm templates

       Required permission: swarm-templates:read

SYNOPSIS
       afs swarm-templates list-swarm-templates [options]

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
       List the swarm templates

       afs swarm-templates list-swarm-templates

```

## update-swarm-template

```text
AFS SWARM-TEMPLATES UPDATE-SWARM-TEMPLATE

NAME
       swarm-templates update-swarm-template

DESCRIPTION
       Update swarm template

       Required permission: admin:swarm-templates

SYNOPSIS
       afs swarm-templates update-swarm-template [options]

OPTIONS
       --template-id (string) (required)
       --budget-usd (string)
          USD amount as a decimal string with at most 6 fraction digits; never a JSON number
       --description (string)
       --spec (JSON)
       --teardown-rule (string)
          Possible values: at_expiry, at_cap_or_expiry, at_work_complete
       --ttl-hours (integer)
          Range: 1 to 8760
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
       Raise a template's budget

       afs swarm-templates update-swarm-template --template-id swt_01J9ZQ3ABCDEFGHJKMNPQRSTVW --budget-usd 500.00

```
