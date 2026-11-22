# Getting started

`afs` is the command line for the Enterprise Agent Factory Service. Commands have the shape
`afs <service> <operation> [--options]`, like the AWS CLI.

## Install

Build from source with Go (1.27 or later):

```bash
make build            # writes bin/afs
./bin/afs --version   # afs-cli/<version> contract/<major.minor> <go> <os>/<arch>
```

Put `bin/afs` on your `PATH`. Signed release archives and Homebrew are described in [Install](install.md).

## Point afs at the service

```bash
afs configure
# AFS endpoint URL [None]: https://afs.example.com
# Default output format (json, yaml, table, text, csv) [json]:
# Default team id (optional) [None]:
```

This writes the `default` profile in `~/.afs/config`. The endpoint must use `https` (plain `http` is
accepted only for `localhost` and other loopback addresses).

## Sign in

```bash
afs configure sso        # choose the sign-in provider the service offers
afs sso login            # prints a URL and a code; approve it in your browser
afs sso get-caller-identity
```

For automation, use an API key instead; see [Credentials](credentials.md).

## Run commands

```bash
afs teams list-teams --output table
afs teams describe-team --team-id team_01J...
afs teams create-team --generate-cli-skeleton > team.json   # a JSON skeleton; edit it, then:
afs teams create-team --cli-input-json file://team.json
```

Every command has help: `afs help`, `afs teams help`, `afs teams create-team help`.
