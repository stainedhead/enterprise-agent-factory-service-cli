# Single sign-on

`afs` signs in with the OpenID Connect device authorization flow, so it works over SSH and in containers.

## Set up a profile

```bash
afs configure sso --profile prod
```

This lists the sign-in providers the service offers (no credential is needed for that call) and writes
`sso_provider_id`, `sso_issuer` and `sso_client_id` into the profile.

## Sign in

```bash
afs sso login --profile prod
```

`afs` prints a verification URL and a code. Open the URL in any browser, enter the code and approve. The
session is cached in `~/.afs/sso/cache/` with mode 0600. If the profile has no sign-in settings yet,
`afs sso login` reads them from the service's default provider and saves them.

## Refresh and expiry

When the access token has expired, `afs` refreshes it silently with the refresh token. If there is no
refresh token, or the identity provider rejects it, the command makes no call and exits 255 with
`Run afs sso login --profile <name>`.

## Sign out

```bash
afs sso logout --profile prod
```

This deletes the cached session. The next command that needs it exits 255 with the sign-in hint.

## Who am I

```bash
afs sso get-caller-identity
```
