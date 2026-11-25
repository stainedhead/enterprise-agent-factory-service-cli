package services_test

import (
	"os"
	"path/filepath"
	"sort"
	"testing"
)

// mcpSpecCandidates are where the 1.3.0 spec lives while it is active and after it is archived.
var mcpSpecCandidates = []string{
	"specs/261125-mcp-gateway/spec.md",
	"specs/archive/261125-mcp-gateway/spec.md",
}

// Release 1.3.0 (governed MCP servers, MCP-CLI-1 to MCP-CLI-4): every functional requirement of the spec names the
// tests that prove it, and each named test exists.
func TestMcpCriteriaTraceToTests(t *testing.T) {
	root := filepath.Join("..", "..", "..")
	var spec []byte
	for _, c := range mcpSpecCandidates {
		if b, err := os.ReadFile(filepath.Join(root, c)); err == nil { //nolint:gosec // fixed repository paths
			spec = b
			break
		}
	}
	if spec == nil {
		t.Fatalf("the 1.3.0 spec is in none of %v", mcpSpecCandidates)
	}
	funcs := testFuncs(t, root)
	rows := frRowRe.FindAllStringSubmatch(string(spec), -1)
	if len(rows) < 14 {
		t.Fatalf("%d requirements found, want at least 14", len(rows))
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
