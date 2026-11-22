# Output

Choose a format with `--output`, `AFS_DEFAULT_OUTPUT` or `output` in the profile.

| Format | Use |
|---|---|
| `json` (default) | The API's response with its field names exactly; the stable format for scripts |
| `yaml` | The same data as YAML |
| `table` | Aligned columns for people; a single object is shown as Field and Value |
| `text` | Tab-separated rows for `grep`, `cut` and `awk` |
| `csv` | The API's CSV export, for operations that offer one (reports, costs, audit events, work orders) |

## --query

`--query` applies a [JMESPath](https://jmespath.org) expression before rendering:

```bash
afs teams list-teams --query "items[].id" --output text
afs teams list-teams --query "items[?kind=='persistent'].slug"
```

An invalid expression is a usage error (exit 252). `--query` cannot be combined with `--output csv`.

## Pager and colour

On a terminal, output goes through a pager: `AFS_PAGER`, else `PAGER`, else `less -R`. An empty `AFS_PAGER`
or `--no-cli-pager` turns it off. `--color on|off|auto` controls colour; `auto` colours only a terminal and
honours `NO_COLOR`.

## Terminal safety

Text from the service can be written by other people (team names, agent names, work-order titles, error
messages). In `table` and `text` output, in error messages and in the sign-in prompt, `afs` replaces control
characters other than newline and tab (ESC, BEL, carriage return, DEL and the C1 range) with visible escapes
such as `\x1b` or `\u009b`, so a value cannot clear the screen, retitle the window or write to the clipboard.
`json` output is the service's own bytes re-indented and is not changed; `yaml` quotes control characters.
Pipe `--output json` through your own tooling when you need the raw value.
