# Technical details

## Toolchain

Go 1.27.4, golangci-lint v2.15.0, govulncheck v1.8.0.

## Exit codes

0 success · 252 usage/parse error · 254 service error · 255 other (AWS CLI compatible).

## Configuration files

`~/.afs/config` (profiles), `~/.afs/credentials` (API keys), `~/.afs/sso/cache/` (tokens, mode 0600).
