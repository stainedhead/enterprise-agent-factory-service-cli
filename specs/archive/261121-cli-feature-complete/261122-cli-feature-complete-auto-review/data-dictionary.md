# Data dictionary: cli-feature-complete-auto-review

No new persisted data. Values validated or changed by the fixes:

| Value | Where | Rule |
|---|---|---|
| INI section, key, value | `~/.afs/config`, `~/.afs/credentials` | no control chars; no `=`, `[`, `]` in keys or sections |
| Provider `issuerUrl`, `cliClientId`, `id` | `ListSignInProviders` | https absolute URL (loopback http), strict charset |
| `sso_issuer`, discovered endpoints | profile, discovery doc | https (loopback http), issuer equality |
| Endpoint URL | option, env, profile | no userinfo, query, fragment |
| Path parameter | `BuildPath` | non-empty, not `.`/`..`, matches snapshot pattern |
| Integer/number option | `exec/build.go` | strict JSON number, re-encoded |
| Page token | `nextToken`, `X-Afs-Next-Token` | must differ from the previous one |
| Control characters | table/text cells, errors, device prompt | replaced with `\xNN` escapes |
