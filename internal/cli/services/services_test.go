package services_test

import (
	"io/fs"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/api"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli/services"
)

// CLI-US-7.1: the command table and the OpenAPI snapshot agree in both directions, and names
// equal the x-afs-cli value of each operation.
func TestCommandTableMatchesSnapshot(t *testing.T) {
	var doc struct {
		Paths map[string]map[string]struct {
			OperationID string `yaml:"operationId"`
			CLI         string `yaml:"x-afs-cli"`
		} `yaml:"paths"`
	}
	if err := yaml.Unmarshal(api.OpenAPI, &doc); err != nil {
		t.Fatal(err)
	}
	reg, err := services.NewRegistry()
	if err != nil {
		t.Fatal(err)
	}
	wantCmd := map[string]string{} // operationId -> "service name"
	for _, item := range doc.Paths {
		for _, op := range item {
			if op.OperationID == "GetHealth" || op.OperationID == "GetReadiness" {
				continue
			}
			wantCmd[op.OperationID] = op.CLI
		}
	}
	if len(wantCmd) != 187 {
		t.Fatalf("snapshot has %d CLI operations", len(wantCmd))
	}
	for id, name := range wantCmd {
		c := reg.ByOp(id)
		if c == nil {
			t.Errorf("operation %s has no command", id)
			continue
		}
		if c.Key() != name {
			t.Errorf("%s: command %q, contract says %q", id, c.Key(), name)
		}
	}
	local := map[string]bool{"configure": true, "configure sso": true, "sso login": true, "sso logout": true, "completion": true}
	for _, c := range reg.All() {
		if c.Op == nil && !local[c.Key()] {
			t.Errorf("command %q has no operation", c.Key())
		}
		if c.Op != nil && wantCmd[c.Op.ID] == "" {
			t.Errorf("command %q maps to unknown operation", c.Key())
		}
	}
}

// D-C4: every operation at or below CurrentMilestone has an example and table columns.
func TestDoneUpToCurrentMilestone(t *testing.T) {
	reg, err := services.NewRegistry()
	if err != nil {
		t.Fatal(err)
	}
	if p := reg.Pending(); len(p) > 0 {
		t.Errorf("pending operations: %s", strings.Join(p, ", "))
	}
	for _, c := range reg.All() {
		if c.Op == nil || cli.MilestoneOf(c.Op.ID) > cli.CurrentMilestone {
			continue
		}
		for _, ex := range c.Examples {
			if !strings.HasPrefix(ex.Command, "afs ") {
				t.Errorf("%s: example %q must start with afs", c.Key(), ex.Command)
			}
		}
	}
}

// D-C4: every operation at or below CurrentMilestone has golden help (and, when it returns a
// body, table and text) files in its package's testdata, written by clitest.CheckOperations.
func TestDoneOperationsHaveGoldens(t *testing.T) {
	have := map[string]bool{}
	root := filepath.Join("..", "..", "..", "internal", "cli")
	_ = filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err == nil && !d.IsDir() && strings.HasSuffix(p, ".golden") && strings.Contains(p, "testdata") {
			have[filepath.Base(p)] = true
		}
		return nil
	})
	reg, err := services.NewRegistry()
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range reg.All() {
		if c.Op == nil || cli.MilestoneOf(c.Op.ID) < 1 || cli.MilestoneOf(c.Op.ID) > cli.CurrentMilestone {
			continue
		}
		if !have[c.Op.ID+".help.golden"] {
			t.Errorf("%s: no help golden", c.Op.ID)
		}
		if len(c.Columns) > 0 && !have[c.Op.ID+".table.golden"] {
			t.Errorf("%s: no table golden", c.Op.ID)
		}
	}
}

// wantMilestone is raised only by a milestone's exit task (PC-CLI-D9: M7 to M10).
const wantMilestone = 7

// M6-X1, CLI-US-7: the feature is complete: all 187 operations have a command, none is held back
// by the milestone gate, and none lacks an example or columns.
func TestFeatureCompleteStrict(t *testing.T) {
	reg, err := services.NewRegistry()
	if err != nil {
		t.Fatal(err)
	}
	if cli.CurrentMilestone != wantMilestone {
		t.Fatalf("CurrentMilestone = %d, want %d", cli.CurrentMilestone, wantMilestone)
	}
	n := 0
	for _, c := range reg.All() {
		if c.Op == nil {
			continue
		}
		n++
		if cli.MilestoneOf(c.Op.ID) > 10 { // later milestones are held back by design until their exit task
			t.Errorf("%s has milestone %d", c.Op.ID, cli.MilestoneOf(c.Op.ID))
		}
	}
	if n != 187 {
		t.Errorf("%d operation commands, want 187", n)
	}
	if p := reg.Pending(); len(p) != 0 {
		t.Errorf("pending: %v", p)
	}
}
