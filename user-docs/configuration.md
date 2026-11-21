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

## The profile team

When a command needs a team id in its path and you leave `--team-id` out, the profile's `team_id` is used.
It is never used to fill an optional filter, so a list is not narrowed without you asking.

## Global options

`--profile`, `--endpoint-url`, `--output`, `--query`, `--no-paginate`, `--max-items`, `--page-size`,
`--starting-token`, `--cli-input-json`, `--generate-cli-skeleton`, `--debug`, `--no-cli-pager`,
`--color on|off|auto`, `--cli-connect-timeout`, `--cli-read-timeout`, `--version`. They are accepted anywhere
on the command line. Timeouts are in seconds; `0` disables one (defaults: connect 30, read 60).
