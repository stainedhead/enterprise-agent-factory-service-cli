package services_test

import (
	"os"
	"path/filepath"
	"sort"
	"testing"
)

// megSpecCandidates are where the 1.2.0 spec lives while it is active and after it is archived.
var megSpecCandidates = []string{
	"specs/261125-model-endpoint-gateway/spec.md",
	"specs/archive/261125-model-endpoint-gateway/spec.md",
}

// Release 1.2.0 (model endpoint gateway, MEG-CLI-1 to MEG-CLI-4): every functional requirement of the spec names the
// tests that prove it, and each named test exists.
func TestModelEndpointCriteriaTraceToTests(t *testing.T) {
	root := filepath.Join("..", "..", "..")
	var spec []byte
	for _, c := range megSpecCandidates {
		if b, err := os.ReadFile(filepath.Join(root, c)); err == nil { //nolint:gosec // fixed repository paths
			spec = b
			break
		}
	}
	if spec == nil {
		t.Fatalf("the 1.2.0 spec is in none of %v", megSpecCandidates)
	}
	funcs := testFuncs(t, root)
	rows := frRowRe.FindAllStringSubmatch(string(spec), -1)
	if len(rows) < 15 {
		t.Fatalf("%d requirements found, want at least 15", len(rows))
	}
	var missing []string
	for _, row := range rows {
		names := provenRe.FindAllStringSubmatch(row[0], -1)
		if len(names) == 0 {
			t.Errorf("%s names no test", row[1])
		}
		for _, m := range names {
			if !funcs[m[1]][m[2]] {
				missing = append(missing, row[1]+": "+m[1]+"."+m[2])
			}
		}
	}
	sort.Strings(missing)
	for _, m := range missing {
		t.Errorf("no such test: %s", m)
	}
}
