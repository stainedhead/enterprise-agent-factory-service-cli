package buildinfo

import (
	"strings"
	"testing"
)

// CLI-US-34.4: the version line carries the contract minor of the pinned snapshot (contract v1.4.x).
func TestVersionPrintsContractMinor(t *testing.T) {
	if !strings.Contains(Line(), " contract/1.5 ") {
		t.Errorf("Line() = %q", Line())
	}
	if ua := UserAgent(); !strings.HasPrefix(ua, "afs-cli/") || !strings.Contains(ua, "/") {
		t.Errorf("UserAgent() = %q", ua)
	}
}
