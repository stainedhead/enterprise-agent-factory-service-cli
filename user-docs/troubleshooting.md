# Troubleshooting

| Symptom | Cause and fix |
|---|---|
| `no endpoint is configured` (255) | Run `afs configure`, set `AFS_ENDPOINT_URL`, or pass `--endpoint-url` |
| `endpoint_url must use https` (255) | Use an `https://` endpoint; `http://` works only for loopback hosts |
| `the service returned an invalid sign-in provider (<field>)` (255) | The service's provider settings are malformed or not https; nothing was written. Ask the service administrator to correct the identity provider |
| `sso_issuer must use https` or `names a different issuer` (255) | The profile's `sso_issuer`, or what the identity provider reports, is not a safe https issuer. Fix the profile with `afs configure sso` or the provider's settings |
| `Run afs sso login --profile <name>` (255) | The sign-in session is missing or expired. Run the command it names |
| `An error occurred (Unauthenticated)` (254) | The credential was rejected. Check `AFS_API_KEY`, the credentials file and `afs sso get-caller-identity` |
| `An error occurred (AccessDenied)` (254) | Your role lacks the permission shown on the `Required permission` line |
| `The config profile (x) could not be found` (255) | Check `--profile` and `AFS_PROFILE`; list profiles in `~/.afs/config` |
| Credentials file refused (255) | `chmod 600 ~/.afs/credentials` |
| `the following arguments are required` (252) | Add the option named in the message, or use `--cli-input-json` |
| A secret value is refused (252) | Pass it as `file://PATH` or `-`, never on the command line |
| Certificate errors behind a corporate proxy | Set `SSL_CERT_FILE` or `SSL_CERT_DIR` to a bundle that includes the proxy's CA; `HTTPS_PROXY` and `NO_PROXY` are honoured |
| A command waits too long | `--cli-read-timeout` (seconds) bounds the wait for a reply; `--wait-timeout` bounds `--wait` |
| A newer service version notice | `afs` is older than the service's API contract. Upgrade `afs` when convenient |

Run with `--debug` to see what is sent. `afs <service> <operation> help` shows the options of any command.
