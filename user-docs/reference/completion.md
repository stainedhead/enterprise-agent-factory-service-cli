# afs completion

Generated from `afs completion help`; do not edit by hand.

```text
AFS COMPLETION

NAME
       completion

DESCRIPTION
       Print a shell completion script for bash, zsh, fish or PowerShell.

SYNOPSIS
       afs completion [options]

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
       Enable completion in the current bash session

       source <(afs completion bash)

       Install the zsh script

       afs completion zsh > "${fpath[1]}/_afs"

       Install the fish script

       afs completion fish > ~/.config/fish/completions/afs.fish

       Load the PowerShell script in a profile

       afs completion powershell | Out-String | Invoke-Expression

```
