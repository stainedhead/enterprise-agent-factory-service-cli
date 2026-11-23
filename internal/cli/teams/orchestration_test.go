package teams_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/clitest"
)

func createTeam(t *testing.T, extra ...string) clitest.Result {
	t.Helper()
	args := append([]string{"teams", "create-team", "--slug", "pay", "--display-name", "Pay", "--kind", "persistent", "--max-budget-usd", "100.00",
		"--first-admin", `{"email":"ada@example.com"}`}, extra...)
	return clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) {
		clitest.Reply(w, 201, `{"id":"team_1","slug":"pay","version":1}`)
	}, args...)
}

func compact(b []byte) string { return strings.ReplaceAll(string(b), " ", "") }

// OR-CLI-26.9, ORC-29: --orchestration takes Key=Value pairs and sends the object the API takes.
func TestCreateTeamOrchestrationShorthand(t *testing.T) {
	for name, tc := range map[string]struct {
		arg, want string
	}{
		"tool only":  {"Tool=buzz", `"orchestration":{"tool":"buzz"}`},
		"all keys":   {"Tool=paperclip,WebUiExposure=private,WebApiExposure=off,RegistrationMode=managed,ModelId=m1", `"orchestration":{"tool":"paperclip","webUi":{"exposure":"private"},"webApi":{"exposure":"off"},"registrationMode":"managed","models":{"defaultModelId":"m1"}}`},
		"lower case": {"webuiexposure=public", `"orchestration":{"webUi":{"exposure":"public"}}`},
		"model id":   {"DefaultModelId=m2", `"orchestration":{"models":{"defaultModelId":"m2"}}`},
		"json":       {`{"tool":"buzz","models":{"access":"selected","modelIds":["m1"],"defaultModelId":"m1"}}`, `"orchestration":{"tool":"buzz","models":{"access":"selected","modelIds":["m1"],"defaultModelId":"m1"}}`},
	} {
		r := createTeam(t, "--orchestration", tc.arg)
		if r.Code != 0 || len(r.Requests) != 1 || !strings.Contains(compact(r.ReqBody), tc.want) {
			t.Errorf("%s: exit %d body %s err %s", name, r.Code, r.ReqBody, r.Err)
		}
	}
	if r := createTeam(t); strings.Contains(string(r.ReqBody), "orchestration") {
		t.Errorf("no option, no member: %s", r.ReqBody)
	}
}

// OR-CLI-26.9: a shorthand that cannot be right is a usage error before any call.
func TestCreateTeamOrchestrationShorthandErrors(t *testing.T) {
	for name, tc := range map[string]struct{ arg, want string }{
		"unknown key":  {"Flavour=x", `unknown key "Flavour"`},
		"no equals":    {"Tool", "takes Key=Value pairs or JSON"},
		"bad tool":     {"Tool=slack", "Tool must be one of paperclip, buzz"},
		"bad level":    {"WebUiExposure=internal", "WebUiExposure must be one of off, private, public"},
		"bad mode":     {"RegistrationMode=auto", "RegistrationMode must be one of manual, managed"},
		"empty model":  {"ModelId=", "ModelId cannot be empty"},
		"twice":        {"Tool=buzz,Tool=paperclip", "is given twice"},
		"empty member": {"Tool=buzz,", "takes Key=Value pairs or JSON"},
	} {
		r := createTeam(t, "--orchestration", tc.arg)
		if r.Code != 252 || len(r.Requests) != 0 || !strings.Contains(r.Err, tc.want) {
			t.Errorf("%s: exit %d requests %d %q", name, r.Code, len(r.Requests), r.Err)
		}
	}
}

// ORC-26: the limit is a plain option on create and update; the refusals of the new reasons print hints.
func TestMaxCollaborationExposure(t *testing.T) {
	r := createTeam(t, "--max-collaboration-exposure", "private")
	if r.Code != 0 || !strings.Contains(compact(r.ReqBody), `"maxCollaborationExposure":"private"`) {
		t.Errorf("exit %d body %s", r.Code, r.ReqBody)
	}
	u := clitest.Run(t, collabTeam, "teams", "update-team", "--team-id", "team_1", "--max-collaboration-exposure", "off", "--if-match", "3")
	if u.Code != 0 || string(u.ReqBody) != `{"maxCollaborationExposure":"off"}` {
		t.Errorf("exit %d body %s err %s", u.Code, u.ReqBody, u.Err)
	}
	if bad := createTeam(t, "--max-collaboration-exposure", "world"); bad.Code != 252 {
		t.Errorf("exit %d", bad.Code)
	}
	d := clitest.Run(t, clitest.Denied("teams:set-guardrails"), "teams", "update-team", "--team-id", "team_1", "--max-collaboration-exposure", "public", "--if-match", "3")
	if d.Code != 254 || !strings.Contains(d.Err, "Required permission: teams:set-guardrails") {
		t.Errorf("%d %q", d.Code, d.Err)
	}
}

func TestCreateTeamOrchestrationRefusals(t *testing.T) {
	for _, tc := range []struct {
		status              int
		code, details, want string
	}{
		{400, "ValidationError", `{"fields":[{"field":"orchestration.webUiExposure","reason":"exposure_above_limit"}]}`, "field: orchestration.webUiExposure reason: exposure_above_limit"},
		{400, "ValidationError", `{"fields":[{"field":"orchestration.tool","reason":"tool_not_enabled"}]}`, "That tool is not enabled for the service."},
		{400, "ValidationError", `{"fields":[{"field":"orchestration.models.defaultModelId","reason":"no_model_available"}]}`, "No model can be used as the default."},
		{409, "Conflict", `{"reason":"orchestration_runtime_unset"}`, "orchestration.tools"},
	} {
		r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) {
			clitest.Reply(w, tc.status, `{"code":"`+tc.code+`","message":"refused","requestId":"req_t","details":`+tc.details+`}`)
		}, "teams", "create-team", "--slug", "pay", "--display-name", "Pay", "--kind", "persistent", "--max-budget-usd", "100.00",
			"--first-admin", `{"email":"ada@example.com"}`, "--orchestration", "Tool=buzz")
		if r.Code != 254 || !strings.Contains(r.Err, "when calling the CreateTeam operation") || !strings.Contains(r.Err, tc.want) {
			t.Errorf("%s: exit %d %q", tc.want, r.Code, r.Err)
		}
	}
}

// ORC-27, OR-CLI-27: --may-expose-collaboration keeps working and prints one deprecation line.
func TestMayExposeCollaborationDeprecated(t *testing.T) {
	want := "Warning: --may-expose-collaboration is deprecated and will be removed in 2.0.0; use --max-collaboration-exposure (off, private or public) instead.\n"
	r := clitest.Run(t, collabTeam, "teams", "update-team", "--team-id", "team_1", "--may-expose-collaboration", "--if-match", "3")
	if r.Code != 0 || r.Err != want {
		t.Errorf("update: exit %d %q", r.Code, r.Err)
	}
	if c := createTeam(t, "--no-may-expose-collaboration"); c.Code != 0 || c.Err != want {
		t.Errorf("create: exit %d %q", c.Code, c.Err)
	}
	h := clitest.Run(t, nil, "teams", "update-team", "help")
	if !strings.Contains(h.Out, "Deprecated: use --max-collaboration-exposure") {
		t.Errorf("help: %q", h.Out)
	}
}

// OR-CLI-28.3, MOD-7: list-team-models shows the access mode on stderr for table and text, never for json.
func TestListTeamModelsAccessNotice(t *testing.T) {
	body := `{"access":"selected","modelIds":["m1","m2"],"items":[{"id":"m1","provider":"anthropic","displayName":"M1","enabled":true,"callable":true,"priceState":"approved","inferenceProfileState":"active","contextWindow":200000,"maxOutputTokens":64000}]}`
	h := func(w http.ResponseWriter, _ *http.Request) { clitest.Reply(w, 200, body) }
	r := clitest.Run(t, h, "teams", "list-team-models", "--team-id", "team_1", "--output", "table")
	if r.Code != 0 || r.Err != "Team model access: selected (2 in the team's list).\n" || !strings.Contains(r.Out, "Callable") || !strings.Contains(r.Out, "active") {
		t.Errorf("exit %d err %q out %q", r.Code, r.Err, r.Out)
	}
	all := `{"access":"all","modelIds":[],"items":[]}`
	r = clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) { clitest.Reply(w, 200, all) }, "teams", "list-team-models", "--team-id", "team_1", "--output", "text")
	if r.Err != "Team model access: all (every model enabled for the service).\n" {
		t.Errorf("err %q", r.Err)
	}
	if j := clitest.Run(t, h, "teams", "list-team-models", "--team-id", "team_1", "--output", "json"); j.Err != "" {
		t.Errorf("json printed %q", j.Err)
	}
	if r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) { clitest.Reply(w, 200, `{"items":[]}`) }, "teams", "list-team-models", "--team-id", "team_1", "--output", "table"); r.Err != "" {
		t.Errorf("a body without access prints no notice: %q", r.Err)
	}
}

// MOD-8: put-team-models --access all|selected [--model-ids ...]; the API's refusals print their hints.
func TestPutTeamModelsAccess(t *testing.T) {
	ok := func(w http.ResponseWriter, _ *http.Request) {
		clitest.Reply(w, 200, `{"access":"all","modelIds":[],"items":[]}`)
	}
	r := clitest.Run(t, ok, "teams", "put-team-models", "--team-id", "team_1", "--access", "all")
	if r.Code != 0 || string(r.ReqBody) != `{"access":"all"}` || r.Req.Header.Get("If-Match") != "" {
		t.Errorf("exit %d body %s", r.Code, r.ReqBody)
	}
	r = clitest.Run(t, ok, "teams", "put-team-models", "--team-id", "team_1", "--access", "selected", "--model-ids", "m1", "m2")
	if string(r.ReqBody) != `{"access":"selected","modelIds":["m1","m2"]}` {
		t.Errorf("body %s", r.ReqBody)
	}
	if bad := clitest.Run(t, ok, "teams", "put-team-models", "--team-id", "team_1", "--access", "some"); bad.Code != 252 {
		t.Errorf("exit %d", bad.Code)
	}
	for _, reason := range []string{"model_ids_not_allowed_for_all", "model_not_enabled_for_service"} {
		e := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) {
			clitest.Reply(w, 400, `{"code":"ValidationError","message":"invalid","requestId":"req_m","details":{"fields":[{"field":"modelIds","reason":"`+reason+`"}]}}`)
		}, "teams", "put-team-models", "--team-id", "team_1", "--access", "all", "--model-ids", "m1")
		if e.Code != 254 || !strings.Contains(e.Err, "field: modelIds reason: "+reason) || strings.Count(e.Err, "\n") < 4 {
			t.Errorf("%s: exit %d %q", reason, e.Code, e.Err)
		}
	}
	d := clitest.Run(t, clitest.Denied("teams:set-guardrails"), "teams", "put-team-models", "--team-id", "team_1", "--access", "all")
	if d.Code != 254 || !strings.Contains(d.Err, "Required permission: teams:set-guardrails") {
		t.Errorf("%d %q", d.Code, d.Err)
	}
}
