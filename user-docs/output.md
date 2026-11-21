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
