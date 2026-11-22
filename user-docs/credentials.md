# Credentials

`afs` looks for a credential in this order and uses the first it finds:

1. `AFS_API_KEY` in the environment.
2. `api_key` for the profile in `~/.afs/credentials`.
3. The profile's cached single sign-on session ([Single sign-on](sso.md)).
4. The workload token file: `AFS_API_TOKEN_FILE`, else `/var/run/secrets/afs/api/token` when it exists. Agents and pool
   workers running in the platform get this file projected into their pod, so they need no setup. The file is re-read
   for every request, so a rotated token is picked up at once. See [Work orders](work-orders.md).

```ini
# ~/.afs/credentials  (chmod 600)
[default]
api_key = afsk_...
```

A credentials file that others can read is refused (exit 255).

## Keeping secrets out of history

- Never put a token or API key in a command line. Use the environment or the credentials file.
- Commands that take a secret value (for example `afs secrets put-secret-value --secret-value`) accept only
  `file://PATH`, `fileb://PATH` or `-` (standard input). A value typed on the command line is refused
  (exit 252). The value is sent once and is never printed or logged.
- `--cli-input-json` with a secret member must come from a file (`file://PATH`), not inline.
- `--debug` redacts the `Authorization` header and secret fields.
- An API key's secret (`secretKey`) is shown once when the key is created; store it then.
