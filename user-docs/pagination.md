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
`--max-items`. With `--no-paginate` only the first page is printed.
