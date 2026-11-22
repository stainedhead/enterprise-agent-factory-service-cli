# Configuration

## Files

| File | Holds |
|---|---|
| `~/.afs/config` | Profiles: `endpoint_url`, `output`, `team_id`, `sso_provider_id`, `sso_issuer`, `sso_client_id` |
| `~/.afs/credentials` | `api_key` per profile. Must be readable by you only (`chmod 600`) |
| `~/.afs/sso/cache/` | Cached single sign-on sessions (mode 0600) |

The format is INI, as in the AWS CLI:

```ini
[default]
endpoint_url = https://afs.example.com
output = json

[profile prod]
endpoint_url = https://afs.prod.example.com
team_id = team_01J...
sso_issuer = https://idp.example.com
sso_client_id = 0oa1cli
```

`AFS_CONFIG_FILE` and `AFS_SHARED_CREDENTIALS_FILE` override the two file paths.

## Precedence

Highest first: command-line option, environment variable, profile in `~/.afs/config`, default.

| Setting | Option | Environment | Profile key |
|---|---|---|---|
| Profile | `--profile` | `AFS_PROFILE` | |
| Endpoint | `--endpoint-url` | `AFS_ENDPOINT_URL` | `endpoint_url` |
| Output format | `--output` | `AFS_DEFAULT_OUTPUT` | `output` |
| API key | | `AFS_API_KEY` | `api_key` (credentials file) |

An unknown profile name is an error (`The config profile (name) could not be found`, exit 255).

## Environment variables

| Variable | Effect |
|---|---|
| `AFS_PROFILE` | Profile to use when `--profile` is absent |
| `AFS_ENDPOINT_URL` | API endpoint (host root, for example `https://afs.example.com`) |
| `AFS_DEFAULT_OUTPUT` | Output format when `--output` is absent |
| `AFS_API_KEY` | API key; first in the credential order |
| `AFS_API_TOKEN_FILE` | Workload token file (default `/var/run/secrets/afs/api/token` when present); last in the credential order |
| `AFS_CONFIG_FILE` | Path of the config file (default `~/.afs/config`) |
| `AFS_SHARED_CREDENTIALS_FILE` | Path of the credentials file (default `~/.afs/credentials`) |
| `AFS_PAGER` | Pager command; empty disables the pager. Falls back to `PAGER`, then `less -R` |
| `HTTPS_PROXY`, `NO_PROXY` | Outbound proxy |
| `SSL_CERT_FILE`, `SSL_CERT_DIR` | Extra trusted certificate authorities (Linux and macOS) |

An invalid value on the command line is a usage error (exit 252); the same invalid value from the environment or a
profile is a configuration error (exit 255) that names where it came from.

## The profile team

When a command needs a team id in its path and you leave `--team-id` out, the profile's `team_id` is used.
It is never used to fill an optional filter, so a list is not narrowed without you asking.

## Global options

`--profile`, `--endpoint-url`, `--output`, `--query`, `--no-paginate`, `--max-items`, `--page-size`,
`--starting-token`, `--cli-input-json`, `--generate-cli-skeleton`, `--debug`, `--no-cli-pager`,
`--color on|off|auto`, `--cli-connect-timeout`, `--cli-read-timeout`, `--version`. They are accepted anywhere
on the command line. Timeouts are in seconds; `0` disables one (defaults: connect 30, read 60).
