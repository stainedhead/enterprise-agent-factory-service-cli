// Package buildinfo carries version metadata stamped in at link time by the Makefile.
package buildinfo

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/spec"
)

// Set with -ldflags "-X github.com/stainedhead/enterprise-agent-factory-service-cli/internal/buildinfo.Version=...".
var (
	Version = "dev"
	Commit  = "unknown"
	Date    = "unknown"
)

// ContractMajorMinor is the API contract version this build was generated for, as major.minor.
func ContractMajorMinor() string {
	parts := strings.Split(spec.ContractVersion, ".")
	if len(parts) < 2 {
		return spec.ContractVersion
	}
	return parts[0] + "." + parts[1]
}

// Format renders the --version line (CLI-15): afs-cli/<version> contract/<major.minor> <go> <os>/<arch>.
func Format(version, contract, goVersion, goos, goarch string) string {
	return fmt.Sprintf("afs-cli/%s contract/%s %s %s/%s", version, contract, goVersion, goos, goarch)
}

// Line returns the --version line of this build.
func Line() string {
	return Format(Version, ContractMajorMinor(), runtime.Version(), runtime.GOOS, runtime.GOARCH)
}

// UserAgent is the User-Agent header value: afs-cli/<version> (<os>/<arch>).
func UserAgent() string {
	return fmt.Sprintf("afs-cli/%s (%s/%s)", Version, runtime.GOOS, runtime.GOARCH)
}
