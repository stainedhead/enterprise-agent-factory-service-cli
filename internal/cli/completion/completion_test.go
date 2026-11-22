package completion_test

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli/completion"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli/services"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/clitest"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/golden"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/spec"
)

// fixture is a three-operation command table, small enough to read its golden scripts.
func fixture(t *testing.T) *cli.Registry {
	t.Helper()
	var ops []spec.Operation
	for _, id := range []string{"ListTeams", "SuspendTeam", "ListAuditEvents"} {
		ops = append(ops, *spec.ByID(id))
	}
	r := cli.NewRegistry(ops)
	completion.Register(r)
	if err := r.Finalize(); err != nil {
		t.Fatal(err)
	}
	return r
}

// CLI-3, D-C27: a golden script per shell.
func TestCompletionGoldenScripts(t *testing.T) {
	r := fixture(t)
	for _, shell := range completion.Shells {
		var b bytes.Buffer
		if err := completion.Generate(&b, shell, r); err != nil {
			t.Fatal(err)
		}
		golden.Check(t, shell+".golden", b.String())
	}
}

// CLI-3: every command and option of the real table is in every script.
func TestCompletionListsEveryCommand(t *testing.T) {
	r, err := services.NewRegistry()
	if err != nil {
		t.Fatal(err)
	}
	for _, shell := range completion.Shells {
		var b bytes.Buffer
		if err := completion.Generate(&b, shell, r); err != nil {
			t.Fatal(err)
		}
		script := b.String()
		for _, c := range r.All() {
			if c.Name != "" && !strings.Contains(script, c.Name) {
				t.Errorf("%s: no %s", shell, c.Key())
			}
			if c.Name == "" && !strings.Contains(script, c.Service) {
				t.Errorf("%s: no service %s", shell, c.Service)
			}
			for _, o := range c.Options {
				if !strings.Contains(script, strings.TrimPrefix(o.Flag, "--")) {
					t.Errorf("%s: %s has no %s", shell, c.Key(), o.Flag)
				}
			}
		}
	}
}

func TestCompletionUnknownShell252(t *testing.T) {
	r := clitest.Run(t, nil, "completion", "tcsh")
	if r.Code != 252 || !strings.Contains(r.Err, "unknown shell") || r.Out != "" {
		t.Fatalf("%+v", r)
	}
	if r := clitest.Run(t, nil, "completion"); r.Code != 252 {
		t.Fatalf("no shell: %+v", r)
	}
}

// The scripts must be valid in the shells that are installed (CLI-3).
func TestCompletionScriptsParse(t *testing.T) {
	for shell, args := range map[string][]string{"bash": {"-n"}, "zsh": {"-n"}, "fish": {"--no-execute"}} {
		path, err := exec.LookPath(shell)
		if err != nil {
			continue
		}
		r := clitest.Run(t, nil, "completion", shell)
		if r.Code != 0 {
			t.Fatalf("%s: %+v", shell, r)
		}
		file := filepath.Join(t.TempDir(), "afs."+shell)
		if err := os.WriteFile(file, []byte(r.Out), 0o600); err != nil {
			t.Fatal(err)
		}
		if out, err := exec.CommandContext(context.Background(), path, append(args, file)...).CombinedOutput(); err != nil { //nolint:gosec // fixed shell names
			t.Errorf("%s -n: %v\n%s", shell, err, out)
		}
	}
}

// The bash script completes services, operations, options and enumerated values.
func TestCompletionBashBehaviour(t *testing.T) {
	bash, err := exec.LookPath("bash")
	if err != nil {
		t.Skip("no bash")
	}
	r := clitest.Run(t, nil, "completion", "bash")
	file := filepath.Join(t.TempDir(), "afs.bash")
	if err := os.WriteFile(file, []byte(r.Out), 0o600); err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct{ words, cword, want string }{
		{`afs te`, "1", "teams"},
		{`afs teams create-t`, "2", "create-team"},
		{`afs teams create-team --ki`, "3", "--kind"},
		{`afs teams create-team --kind ""`, "4", "persistent ephemeral"},
		{`afs --profile x tea`, "3", "teams"},
		{`afs teams list-teams --output ""`, "4", "json yaml table text csv"},
	} {
		script := `source ` + file + `; COMP_WORDS=(` + c.words + `); COMP_CWORD=` + c.cword + `; _afs; echo "${COMPREPLY[@]}"`
		out, err := exec.CommandContext(context.Background(), bash, "-c", script).CombinedOutput() //nolint:gosec // test script
		if err != nil || !strings.Contains(string(out), c.want) {
			t.Errorf("%s: %q (%v), want %q", c.words, out, err, c.want)
		}
	}
}
