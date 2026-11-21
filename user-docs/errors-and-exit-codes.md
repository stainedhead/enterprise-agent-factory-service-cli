# Errors and exit codes

## Format

An error from the service is printed on stderr in the AWS style:

```text
An error occurred (AccessDenied) when calling the GetTeam operation: not allowed
Required permission: teams:read
Request ID: req_01J...
```

Quote the request id when you ask for help.

## Exit codes

| Code | Meaning |
|---|---|
| 0 | Success |
| 130 | Interrupted (Ctrl-C). With `--wait`, the change that was already sent still stands |
| 252 | Usage error: unknown option, missing or invalid value, bad JSON. No request was made |
| 254 | The service answered with an error (including `AccessDenied` and a resource that reached `Failed`) |
| 255 | Anything else: network failure, timeout, configuration or credential problem, `--wait-timeout` |

## Retries

`afs` makes up to 3 attempts for a throttled (429), unavailable (503) or dropped request, with backoff and
`Retry-After`. Requests that create something carry an `Idempotency-Key`, so a retry cannot create it twice.
A request that may have been applied is not retried unless it is safe to.

## Safe updates

Updates and deletes send `If-Match` with the resource's version. If you do not pass `--if-match <version>`,
`afs` reads the current version first and tells you which one it used. A `412` means someone changed the
resource; read it again and retry. `428` means a version is required.

## --debug

`--debug` writes the requests and responses to stderr with credentials and secret values redacted.
