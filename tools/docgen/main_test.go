package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli/services"
)

// The committed reference must equal what docgen renders now.
func TestDocgenMatchesHelp(t *testing.T) {
	reg, err := services.NewRegistry()
	if err != nil {
		t.Fatal(err)
	}
	files, err := render(reg)
	if err != nil {
		t.Fatal(err)
	}
	dir := filepath.Join("..", "..", "user-docs", "reference")
	for name, want := range files {
		got, err := os.ReadFile(filepath.Join(dir, name)) //nolint:gosec // repository file
		if err != nil {
			t.Fatalf("%s: %v (run make docs)", name, err)
		}
		if string(got) != want {
			t.Errorf("%s is stale (run make docs)", name)
		}
	}
	entries, _ := os.ReadDir(dir)
	if len(entries) != len(files) {
		t.Errorf("user-docs/reference has %d files, docgen renders %d", len(entries), len(files))
	}
}
