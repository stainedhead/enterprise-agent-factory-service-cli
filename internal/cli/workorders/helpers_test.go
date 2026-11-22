package workorders_test

import (
	"os"
	"path/filepath"
	"testing"
)

func tokenFile(t *testing.T) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "token")
	if err := os.WriteFile(p, []byte("workload-token\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	return p
}
