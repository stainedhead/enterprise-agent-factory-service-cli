package agents_test

import (
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/clitest"
)

func ok201(w http.ResponseWriter, _ *http.Request) { clitest.Reply(w, 201, `{"id":"agt_1"}`) }
func ok200(w http.ResponseWriter, _ *http.Request) {
	clitest.Reply(w, 200, `{"id":"agt_1","version":3}`)
}

func create(t *testing.T, extra ...string) clitest.Result {
	t.Helper()
	args := append([]string{"agents", "create-agent", "--team-id", "team_1", "--name", "coder", "--template", "claude-code", "--image", image}, extra...)
	return clitest.Run(t, ok201, args...)
}

func hasBody(r clitest.Result, want string) bool {
	return strings.Contains(strings.ReplaceAll(string(r.ReqBody), " ", ""), want)
}

// OR-CLI-28.1, MOD-3: --models all or a list of ids, and --default-model-id, become the models member.
func TestCreateAgentModelsBody(t *testing.T) {
	for name, tc := range map[string]struct {
		args []string
		want string
	}{
		"all":          {[]string{"--models", "all"}, `"models":{"access":"all"}`},
		"all+default":  {[]string{"--models", "all", "--default-model-id", "m1"}, `"models":{"access":"all","defaultModelId":"m1"}`},
		"selected":     {[]string{"--models", "m1", "m2", "--default-model-id", "m2"}, `"models":{"access":"selected","modelIds":["m1","m2"],"defaultModelId":"m2"}`},
		"default only": {[]string{"--default-model-id", "m1"}, `"models":{"defaultModelId":"m1"}`},
	} {
		r := create(t, tc.args...)
		if r.Code != 0 || len(r.Requests) != 1 || !hasBody(r, tc.want) {
			t.Errorf("%s: exit %d body %s err %s", name, r.Code, r.ReqBody, r.Err)
		}
		if strings.Contains(string(r.ReqBody), "modelId\"") && !strings.Contains(string(r.ReqBody), "defaultModelId") {
			t.Errorf("%s: unexpected modelId in %s", name, r.ReqBody)
		}
		if r.Err != "" {
			t.Errorf("%s: stderr %q", name, r.Err)
		}
	}
	if r := create(t); strings.Contains(string(r.ReqBody), "models") {
		t.Errorf("no model option, no models member: %s", r.ReqBody)
	}
}

// OR-CLI-28.1: a models word that cannot be right is a usage error before any call.
func TestCreateAgentModelsUsageErrors(t *testing.T) {
	for name, tc := range map[string]struct {
		args []string
		want string
	}{
		"all with ids":          {[]string{"--models", "all", "m1"}, "cannot be combined"},
		"ids with all":          {[]string{"--models", "m1", "all"}, "cannot be combined"},
		"empty id":              {[]string{"--models", ""}, "cannot be empty"},
		"model-id and default":  {[]string{"--model-id", "m1", "--default-model-id", "m2"}, "deprecated name of --default-model-id"},
		"models with no values": {[]string{"--models"}, "expected at least one argument"},
		"default with no value": {[]string{"--default-model-id"}, "expected one argument"},
	} {
		r := create(t, tc.args...)
		if r.Code != 252 || len(r.Requests) != 0 || !strings.Contains(r.Err, tc.want) {
			t.Errorf("%s: exit %d requests %d %q", name, r.Code, len(r.Requests), r.Err)
		}
	}
}

// OR-CLI-28.2: --model-id keeps working and prints one deprecation line on stderr; stdout is unchanged.
func TestModelIDDeprecated(t *testing.T) {
	r := create(t, "--model-id", "anthropic.claude-sonnet-4-5")
	if r.Code != 0 || !hasBody(r, `"modelId":"anthropic.claude-sonnet-4-5"`) {
		t.Fatalf("exit %d body %s", r.Code, r.ReqBody)
	}
	want := "Warning: --model-id is deprecated and will be removed in 2.0.0; use --default-model-id (and --models to choose the set) instead.\n"
	if r.Err != want {
		t.Errorf("stderr %q", r.Err)
	}
	u := clitest.Run(t, ok200, "agents", "update-agent", "--team-id", "team_1", "--agent-id", "agt_1", "--model-id", "m1", "--if-match", "3")
	if u.Code != 0 || !strings.Contains(u.Err, "--model-id is deprecated") {
		t.Errorf("update: exit %d %q", u.Code, u.Err)
	}
	// the member in --cli-input-json counts too
	f := t.TempDir() + "/a.json"
	if err := os.WriteFile(f, []byte(`{"teamId":"team_1","agentId":"agt_1","modelId":"m1","version":3}`), 0o600); err != nil {
		t.Fatal(err)
	}
	u = clitest.Run(t, ok200, "agents", "update-agent", "--cli-input-json", "file://"+f)
	if u.Code != 0 || !strings.Contains(u.Err, "--model-id is deprecated") {
		t.Errorf("input: exit %d %q", u.Code, u.Err)
	}
	// no warning when the option is not used
	if r := create(t, "--default-model-id", "m1"); r.Err != "" {
		t.Errorf("stderr %q", r.Err)
	}
}

// MOD-3: update sends only the members given, so omitted ones are kept by the API.
func TestUpdateAgentModels(t *testing.T) {
	u := clitest.Run(t, ok200, "agents", "update-agent", "--team-id", "team_1", "--agent-id", "agt_1", "--default-model-id", "m2", "--if-match", "3")
	if u.Code != 0 || string(u.ReqBody) != `{"models":{"defaultModelId":"m2"}}` {
		t.Errorf("exit %d body %s", u.Code, u.ReqBody)
	}
	u = clitest.Run(t, ok200, "agents", "update-agent", "--team-id", "team_1", "--agent-id", "agt_1", "--models", "all", "--if-match", "3")
	if string(u.ReqBody) != `{"models":{"access":"all"}}` {
		t.Errorf("body %s", u.ReqBody)
	}
	// an object from --cli-input-json is kept and --default-model-id is folded into it
	f := t.TempDir() + "/a.json"
	if err := os.WriteFile(f, []byte(`{"teamId":"team_1","agentId":"agt_1","models":{"access":"selected","modelIds":["m1","m2"]},"version":3}`), 0o600); err != nil {
		t.Fatal(err)
	}
	u = clitest.Run(t, ok200, "agents", "update-agent", "--cli-input-json", "file://"+f, "--default-model-id", "m1")
	if !hasBody(u, `"models":{"access":"selected","modelIds":["m1","m2"],"defaultModelId":"m1"}`) {
		t.Errorf("body %s err %s", u.ReqBody, u.Err)
	}
	f2 := t.TempDir() + "/b.json"
	if err := os.WriteFile(f2, []byte(`{"teamId":"team_1","agentId":"agt_1","models":5,"version":3}`), 0o600); err != nil {
		t.Fatal(err)
	}
	if u = clitest.Run(t, ok200, "agents", "update-agent", "--cli-input-json", "file://"+f2); u.Code != 252 || !strings.Contains(u.Err, "--models takes the word all or model ids") {
		t.Errorf("exit %d %q", u.Code, u.Err)
	}
}

// ORC-30: --role is out of help but still sent; the API's own refusal is printed with its hint (exit 254).
func TestCreateAgentRoleHiddenAndRefusedByAPI(t *testing.T) {
	h := clitest.Run(t, nil, "agents", "create-agent", "help")
	if strings.Contains(h.Out, "--role") || !strings.Contains(h.Out, "--default-model-id") || !strings.Contains(h.Out, "--models (list)") {
		t.Errorf("help: %q", h.Out)
	}
	r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) {
		clitest.Reply(w, 400, `{"code":"ValidationError","message":"invalid","requestId":"req_r","details":{"fields":[{"field":"role","reason":"role_reserved"}]}}`)
	}, "agents", "create-agent", "--team-id", "team_1", "--name", "p", "--template", "t", "--image", image, "--role", "collaboration")
	if r.Code != 254 || !hasBody(r, `"role":"collaboration"`) || !strings.Contains(r.Err, "field: role reason: role_reserved") ||
		!strings.Contains(r.Err, "Agents you create are workers") {
		t.Errorf("exit %d body %s err %q", r.Code, r.ReqBody, r.Err)
	}
}

// ORC-30: the orchestration agent's refusals print their hints.
func TestAgentOrchestrationRefusals(t *testing.T) {
	for _, tc := range []struct {
		op, reason, want string
		args             []string
	}{
		{"DeleteAgent", "orchestration_agent_required", "deleted only with its team", []string{"agents", "delete-agent", "--team-id", "team_1", "--agent-id", "agt_1", "--if-match", "3"}},
		{"UpdateAgent", "orchestration_field_fixed", "are fixed", []string{"agents", "update-agent", "--team-id", "team_1", "--agent-id", "agt_1", "--description", "x", "--if-match", "3"}},
		{"CreateAgent", "name_reserved", "reserved for the team's orchestration agent", []string{"agents", "create-agent", "--team-id", "team_1", "--name", "orchestrator", "--template", "t", "--image", image}},
		{"CreateAgent", "default_model_not_allowed", "must be one the agent may use", []string{"agents", "create-agent", "--team-id", "team_1", "--name", "c", "--template", "t", "--image", image, "--default-model-id", "m9"}},
	} {
		r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) {
			clitest.Reply(w, 409, `{"code":"Conflict","message":"refused","requestId":"req_o","details":{"reason":"`+tc.reason+`"}}`)
		}, tc.args...)
		if r.Code != 254 || !strings.Contains(r.Err, "when calling the "+tc.op+" operation") || !strings.Contains(r.Err, tc.want) {
			t.Errorf("%s: exit %d %q", tc.reason, r.Code, r.Err)
		}
	}
}

func TestListAgentsRoleOpenEnum(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) { clitest.Reply(w, 200, `{"items":[]}`) },
		"agents", "list-agents", "--team-id", "team_1", "--role", "orchestration")
	if r.Code != 0 || r.Req.URL.Query().Get("role") != "orchestration" {
		t.Errorf("exit %d %s %q", r.Code, r.Req.URL.RawQuery, r.Err)
	}
	h := clitest.Run(t, nil, "agents", "list-agents", "help")
	if !strings.Contains(h.Out, "Known values: worker, orchestration") {
		t.Errorf("%q", h.Out)
	}
}
