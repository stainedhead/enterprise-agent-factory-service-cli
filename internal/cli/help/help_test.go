package help_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli/services"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/golden"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/spec"
)

func afs(t *testing.T, args ...string) (int, string, string) {
	t.Helper()
	reg, err := services.NewRegistry()
	if err != nil {
		t.Fatal(err)
	}
	var out, eb bytes.Buffer
	env := cli.NewEnv(strings.NewReader(""), &out, &eb)
	env.Vars = map[string]string{}
	code := cli.Run(context.Background(), args, env, reg)
	return code, out.String(), eb.String()
}

// 5.1-help: help text for every form is held in golden files.
func TestHelpGoldens(t *testing.T) {
	for name, args := range map[string][]string{
		"root.golden":                    {"help"},
		"sso.golden":                     {"sso", "help"},
		"sso-get-caller-identity.golden": {"sso", "get-caller-identity", "help"},
		"configure.golden":               {"configure", "help"},
		"teams-create-team.golden":       {"teams", "create-team", "help"},
	} {
		code, out, _ := afs(t, args...)
		if code != 0 {
			t.Fatalf("%v: exit %d", args, code)
		}
		golden.Check(t, name, out)
	}
}

func TestHelpForms(t *testing.T) {
	_, a, _ := afs(t, "teams", "create-team", "help")
	_, b, _ := afs(t, "teams", "create-team", "--help")
	_, c, _ := afs(t, "help", "teams", "create-team")
	if a != b || a != c {
		t.Error("the help forms differ")
	}
}

func TestServiceHelpListsExactlyItsOperations(t *testing.T) {
	want := map[string][]string{}
	for _, op := range spec.Commands() {
		want[op.Service] = append(want[op.Service], op.Command)
	}
	for svc, cmds := range want {
		_, out, _ := afs(t, svc, "help")
		for _, c := range cmds {
			if !strings.Contains(out, "\n       "+c+" ") {
				t.Errorf("%s help lacks %s", svc, c)
			}
		}
	}
}

func TestOperationHelpShowsOptionsPermissionAndExample(t *testing.T) {
	_, out, _ := afs(t, "sso", "get-caller-identity", "help")
	for _, want := range []string{"EXAMPLES", "afs sso get-caller-identity", "GLOBAL OPTIONS"} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q", want)
		}
	}
	_, out, _ = afs(t, "teams", "create-team", "help")
	for _, want := range []string{"(required)", "Possible values: persistent, ephemeral", "Required permission: teams:create"} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q", want)
		}
	}
	_, out, _ = afs(t, "secrets", "put-secret-value", "help")
	if !strings.Contains(out, "Never pass the value itself") {
		t.Error("secret option help does not say how to pass the value")
	}
}

func TestHelpInvalidChoices252(t *testing.T) {
	for _, args := range [][]string{{"nope", "help"}, {"teams", "nope", "help"}} {
		if code, _, errOut := afs(t, args...); code != 252 || !strings.Contains(errOut, "Invalid choice") {
			t.Errorf("%v: %d %q", args, code, errOut)
		}
	}
}

// FR-018: a sensitive option is not "(required)" (a hidden prompt is offered) and its help lists the safe forms.
func TestSensitiveOptionHelpListsSafeFormsAndIsNotRequired(t *testing.T) {
	for _, op := range spec.All() {
		if op.Body == nil {
			continue
		}
		hasSecret := false
		for _, f := range op.Body.Fields {
			if f.Schema.Ref == "SecretValue" {
				hasSecret = true
			}
		}
		if !hasSecret {
			continue
		}
		code, out, _ := afs(t, op.Service, op.Command, "help")
		if code != 0 {
			t.Fatalf("%s: exit %d", op.ID, code)
		}
		line := ""
		for _, l := range strings.Split(out, "\n") {
			if strings.Contains(l, "--secret-value") && strings.Contains(l, "(") {
				line = strings.TrimSpace(l)
				break
			}
		}
		if line != "--secret-value (file://PATH | fileb://PATH | -; prompted without echo when omitted on a terminal)" {
			t.Errorf("%s: option line = %q", op.ID, line)
		}
	}
}
