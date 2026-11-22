# Pagination

List operations page their results. By default `afs` follows every page and prints one merged result.

| Option | Effect |
|---|---|
| `--no-paginate` | Print the first page only, with `NextToken` if there are more |
| `--max-items N` | Stop after N items and print `NextToken` |
| `--starting-token T` | Resume from a `NextToken` printed earlier |
| `--page-size N` | Items per request (1 to 100) |

```bash
afs teams list-teams --max-items 10
# ... "NextToken": "eyJ0Ijoi..."
afs teams list-teams --max-items 10 --starting-token eyJ0Ijoi...
```

`NextToken` is opaque; pass it back unchanged. A token that `afs` did not print is a usage error.

## CSV

`--output csv` streams every page, writing the header once. It cannot be combined with `--query` or
`--max-items`. With `--no-paginate` only the first page is printed, and when more rows exist `afs` writes
`NextToken: <token>` to stderr (stdout holds only the CSV page, so redirecting it stays clean).

```bash
afs budgets list-costs --team-id team_1 --output csv --no-paginate > costs.csv
# stderr: NextToken: eyJ...
afs budgets list-costs --team-id team_1 --output csv --no-paginate --starting-token eyJ... >> costs.csv
```

A resumed CSV run (`--starting-token`) leaves out the header row, so appending the pages yields one valid
file. CSV tokens come from the service's `X-Afs-Next-Token` header and are not the same as the JSON
`NextToken`; a token from one format cannot be used with the other.
