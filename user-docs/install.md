# Install `afs`

## Download a release

Each release (`v*` tag) publishes archives for darwin/arm64, darwin/amd64, linux/arm64, linux/amd64
and windows/amd64, with `checksums.txt`.

```bash
# Linux, amd64 (adjust the version and platform)
curl -LO https://github.com/stainedhead/enterprise-agent-factory-service-cli/releases/download/v0.1.0/afs_0.1.0_linux_amd64.tar.gz
curl -LO https://github.com/stainedhead/enterprise-agent-factory-service-cli/releases/download/v0.1.0/checksums.txt
sha256sum --ignore-missing -c checksums.txt
tar -xzf afs_0.1.0_linux_amd64.tar.gz afs && sudo mv afs /usr/local/bin/
afs --version        # afs-cli/0.1.0 contract/1.1 go1.27.4 linux/amd64
```

`afs --version` prints the CLI version, the API contract version it was built for (`contract/<major.minor>`),
the Go version and the platform.

## Verify a release

`checksums.txt` is signed keyless with cosign from the release workflow's identity. Check it before you
trust the archives:

```bash
cosign verify-blob \
  --certificate checksums.txt.pem --signature checksums.txt.sig \
  --certificate-identity-regexp '^https://github.com/stainedhead/enterprise-agent-factory-service-cli/.github/workflows/release.yml@' \
  --certificate-oidc-issuer https://token.actions.githubusercontent.com \
  checksums.txt
```

macOS binaries are notarized and the Windows binary is Authenticode-signed once the release pipeline has
its signing certificates. Until then those two checks do not apply, and a draft release says so in its notes.
On macOS, `spctl --assess --type execute $(which afs)` shows whether a binary is notarized.

## Homebrew

```bash
brew install stainedhead/tap/afs
```

The formula is published when the tap repository exists. Until then use the archives.

## Build from source

```bash
git clone https://github.com/stainedhead/enterprise-agent-factory-service-cli.git
cd enterprise-agent-factory-service-cli
make build           # bin/afs, version "dev"
```

## Shell completion

```bash
source <(afs completion bash)                                    # bash, current session
afs completion zsh > "${fpath[1]}/_afs"                          # zsh
afs completion fish > ~/.config/fish/completions/afs.fish        # fish
afs completion powershell | Out-String | Invoke-Expression       # PowerShell
```

The scripts complete services, operations, options and the values of enumerated options (`--output`,
`--kind`, `--role` ...). They are generated from the command table, so they match the installed version.
