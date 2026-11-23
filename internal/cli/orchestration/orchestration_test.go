package orchestration_test

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/clitest"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/golden"
)

const team = "team_01JD2Q3R4S5T6V7W8X9Y0Z0001"

// orch builds a TeamOrchestration document of a Paperclip team: webUi stored and effective levels, webApi off.
func orch(tool, uiStored, uiEff, reasons string, version, generation, observed int, phase string) string {
	host := "null"
	if uiEff != "off" {
		host = `"pay.collab.corp.example.com"`
	}
	supported := tool == "paperclip"
	return fmt.Sprintf(`{"teamId":%q,"agentId":"agt_01JD2Q3R4S5T6V7W8X9Y0Z0001","agentName":"orchestrator","phase":%q,`+
		`"tool":%q,"toolVersion":null,"registrationMode":"managed","limit":"private",`+
		`"webUi":{"supported":%t,"exposure":%q,"effective":%q,"reasons":%s,"hostname":%s,"url":null},`+
		`"webApi":{"supported":%t,"exposure":"off","effective":"off","reasons":[],"hostname":null,"url":null},`+
		`"connectedAgents":4,"unmanagedToolAgents":[],"version":%d,"generation":%d,`+
		`"status":{"phase":%q,"observedGeneration":%d}}`,
		team, phase, tool, supported, uiStored, uiEff, reasons, host, supported, version, generation, phase, observed)
}

func reply(body string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) { clitest.Reply(w, 200, body) }
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}

func TestCommandsDoneOrchestration(t *testing.T) { // OR-CLI-26.1: example, columns, goldens, AccessDenied
	clitest.CheckOperations(t, "orchestration", "GetTeamOrchestration", "PutTeamOrchestration", "ChangeCollaborationTool")
}

// OR-CLI-26.1: one GET, the JSON printed unchanged.
func TestDescribeOrchestrationJSONUnchanged(t *testing.T) {
	body := orch("paperclip", "private", "private", `[]`, 7, 7, 7, "Ready")
	r := clitest.Run(t, reply(body), "orchestration", "describe-orchestration", "--team-id", team, "--output", "json")
	if r.Code != 0 || len(r.Requests) != 1 || r.Req.Method != "GET" || r.Req.URL.Path != "/v1/teams/"+team+"/orchestration" {
		t.Fatalf("%+v", r)
	}
	if !strings.Contains(strings.ReplaceAll(r.Out, " ", ""), `"hostname":"pay.collab.corp.example.com"`) || r.Err != "" {
		t.Errorf("out %q err %q", r.Out, r.Err)
	}
}

// OR-CLI-26.1: table and text columns for a Paperclip team (limit lowers the level) and for a Buzz team (N/A).
func TestDescribeOrchestrationGoldens(t *testing.T) {
	for name, body := range map[string]string{
		"paperclip": orch("paperclip", "public", "private", `["team_limit"]`, 7, 7, 7, "Ready"),
		"buzz":      orch("buzz", "off", "off", `["tool_not_supported"]`, 3, 3, 3, "Ready"),
	} {
		for _, f := range []string{"table", "text"} {
			r := clitest.Run(t, reply(body), "orchestration", "describe-orchestration", "--team-id", team, "--output", f)
			if r.Code != 0 {
				t.Fatal(r.Err)
			}
			golden.Check(t, "describe-"+name+"."+f+".golden", r.Out)
		}
	}
}

// OR-CLI-26.2: the notices go to stderr for table and text, never for json; the exit code stays 0.
func TestDescribeOrchestrationNotices(t *testing.T) {
	lowered := orch("paperclip", "public", "private", `["team_limit"]`, 7, 7, 7, "Ready")
	r := clitest.Run(t, reply(lowered), "orchestration", "describe-orchestration", "--team-id", team, "--output", "table")
	if want := "Not effective: Web UI is private, not public (team_limit)\n"; r.Code != 0 || r.Err != want {
		t.Fatalf("exit %d stderr %q", r.Code, r.Err)
	}
	buzz := orch("buzz", "off", "off", `["tool_not_supported"]`, 3, 3, 3, "Ready")
	r = clitest.Run(t, reply(buzz), "orchestration", "describe-orchestration", "--team-id", team, "--output", "text")
	if want := "Web UI: N/A for Buzz\nWeb API: N/A for Buzz\n"; r.Code != 0 || r.Err != want {
		t.Fatalf("exit %d stderr %q", r.Code, r.Err)
	}
	if j := clitest.Run(t, reply(lowered), "orchestration", "describe-orchestration", "--team-id", team, "--output", "json"); j.Err != "" {
		t.Errorf("json printed %q", j.Err)
	}
	fine := orch("paperclip", "private", "private", `[]`, 7, 7, 7, "Ready")
	if r := clitest.Run(t, reply(fine), "orchestration", "describe-orchestration", "--team-id", team); r.Err != "" {
		t.Errorf("an effective team must print nothing: %q", r.Err)
	}
	if r := clitest.Run(t, reply(`{"tool":`), "orchestration", "describe-orchestration", "--team-id", team, "--output", "table"); r.Code == 0 && r.Err != "" {
		t.Errorf("a body that is not JSON prints no notice: %q", r.Err)
	}
}

// OR-CLI-26.7: the orchestration agent is not there yet: the AWS-style NotFound and the reason hint.
func TestDescribeOrchestrationPendingHint(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) {
		clitest.Reply(w, 404, `{"code":"NotFound","message":"pending","requestId":"req_1","details":{"reason":"orchestration_agent_pending"}}`)
	}, "orchestration", "describe-orchestration", "--team-id", team)
	if r.Code != 254 || !strings.Contains(r.Err, "(NotFound) when calling the GetTeamOrchestration operation") ||
		!strings.Contains(r.Err, "still being created by the 1.1.0 upgrade") {
		t.Errorf("exit %d %q", r.Code, r.Err)
	}
}

func putHandler(version int, got *[]string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		*got = append(*got, r.Method)
		if r.Method == "GET" {
			w.Header().Set("ETag", fmt.Sprintf(`"%d"`, version))
			clitest.Reply(w, 200, orch("paperclip", "private", "private", `[]`, version, version, version, "Ready"))
			return
		}
		clitest.Reply(w, 200, orch("paperclip", "public", "private", `["team_limit"]`, version+1, version+1, version, "Pending"))
	}
}

// OR-CLI-26.3: both surfaces are sent; registrationMode only when given.
func TestPutOrchestrationBody(t *testing.T) {
	var m []string
	r := clitest.Run(t, putHandler(4, &m), "orchestration", "put-orchestration", "--team-id", team,
		"--web-ui-exposure", "private", "--web-api-exposure", "off", "--registration-mode", "managed", "--if-match", "4")
	if r.Code != 0 || r.Req.Method != "PUT" || r.Req.URL.Path != "/v1/teams/"+team+"/orchestration" || len(r.Requests) != 1 {
		t.Fatalf("%+v", r)
	}
	if got := strings.ReplaceAll(string(r.ReqBody), " ", ""); got != `{"webUi":{"exposure":"private"},"webApi":{"exposure":"off"},"registrationMode":"managed"}` &&
		got != `{"registrationMode":"managed","webApi":{"exposure":"off"},"webUi":{"exposure":"private"}}` &&
		got != `{"registrationMode":"managed","webUi":{"exposure":"private"},"webApi":{"exposure":"off"}}` {
		t.Errorf("body %s", r.ReqBody)
	}
	r = clitest.Run(t, putHandler(4, &m), "orchestration", "put-orchestration", "--team-id", team,
		"--web-ui-exposure", "off", "--web-api-exposure", "public", "--if-match", "4")
	if strings.Contains(string(r.ReqBody), "registrationMode") || !strings.Contains(string(r.ReqBody), `"webApi":{"exposure":"public"}`) {
		t.Errorf("body %s", r.ReqBody)
	}
}

// OR-CLI-26.3: a missing surface, or a level that is not off, private or public, is a usage error (252) before any call.
func TestPutOrchestrationUsageErrors(t *testing.T) {
	for name, tc := range map[string]struct {
		args []string
		want string
	}{
		"missing web api": {[]string{"--web-ui-exposure", "off"}, "--web-api-exposure"},
		"missing web ui":  {[]string{"--web-api-exposure", "off"}, "--web-ui-exposure"},
		"bad level":       {[]string{"--web-ui-exposure", "internal", "--web-api-exposure", "off"}, `invalid choice "internal"`},
		"bad mode":        {[]string{"--web-ui-exposure", "off", "--web-api-exposure", "off", "--registration-mode", "auto"}, "registration-mode"},
	} {
		r := clitest.Run(t, reply("{}"), append([]string{"orchestration", "put-orchestration", "--team-id", team}, tc.args...)...)
		if r.Code != 252 || len(r.Requests) != 0 || !strings.Contains(r.Err, tc.want) {
			t.Errorf("%s: exit %d requests %d %q", name, r.Code, len(r.Requests), r.Err)
		}
	}
}

func TestPutOrchestrationInputJSON(t *testing.T) {
	f := t.TempDir() + "/orch.json"
	writeFile(t, f, `{"teamId":"`+team+`","webUi":{"exposure":"public"},"webApi":{"exposure":"private"},"version":4}`)
	var m []string
	r := clitest.Run(t, putHandler(4, &m), "orchestration", "put-orchestration", "--cli-input-json", "file://"+f)
	if r.Code != 0 || !strings.Contains(string(r.ReqBody), `"webUi":{"exposure":"public"}`) || !strings.Contains(string(r.ReqBody), `"webApi":{"exposure":"private"}`) {
		t.Errorf("exit %d body %s err %s", r.Code, r.ReqBody, r.Err)
	}
	// an option beats the input member
	r = clitest.Run(t, putHandler(4, &m), "orchestration", "put-orchestration", "--cli-input-json", "file://"+f, "--web-ui-exposure", "off")
	if !strings.Contains(string(r.ReqBody), `"webUi":{"exposure":"off"}`) {
		t.Errorf("body %s", r.ReqBody)
	}
}

// OR-CLI-26.4: the version is read and sent when --if-match is not given, and said so.
func TestPutOrchestrationReadsVersion(t *testing.T) {
	var m []string
	r := clitest.Run(t, putHandler(7, &m), "orchestration", "put-orchestration", "--team-id", team, "--web-ui-exposure", "private", "--web-api-exposure", "off")
	if r.Code != 0 || strings.Join(m, ",") != "GET,PUT" || r.Req.Header.Get("If-Match") != `"7"` {
		t.Fatalf("%v %+v", m, r)
	}
	if !strings.Contains(r.Err, "Applied to version 7; use --if-match to guard against concurrent changes.") {
		t.Errorf("stderr %q", r.Err)
	}
}

func TestPutOrchestration412Message(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("ETag", `"5"`)
		clitest.Reply(w, 412, `{"code":"PreconditionFailed","message":"stale","requestId":"req_4"}`)
	}, "orchestration", "put-orchestration", "--team-id", team, "--web-ui-exposure", "off", "--web-api-exposure", "off", "--if-match", "4")
	if r.Code != 254 || !strings.Contains(r.Err, "Re-run with --if-match 5") {
		t.Errorf("exit %d %q", r.Code, r.Err)
	}
}

// OR-CLI-26.5: the validation fields of the new refusal reasons are printed with one hint each.
func TestPutOrchestrationRefusalReasons(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) {
		clitest.Reply(w, 400, `{"code":"ValidationError","message":"invalid","requestId":"req_v","details":{"fields":[`+
			`{"field":"webUi.exposure","reason":"exposure_above_limit"},{"field":"webApi.exposure","reason":"exposure_above_limit"},`+
			`{"field":"registrationMode","reason":"registration_mode_not_supported"}]}}`)
	}, "orchestration", "put-orchestration", "--team-id", team, "--web-ui-exposure", "public", "--web-api-exposure", "public", "--if-match", "1")
	if r.Code != 254 {
		t.Fatal(r.Code)
	}
	for _, want := range []string{"field: webUi.exposure reason: exposure_above_limit\nfield: webApi.exposure reason: exposure_above_limit\nfield: registrationMode reason: registration_mode_not_supported\n",
		"That level is above the team's limit. A service admin can raise the limit: afs teams update-team --max-collaboration-exposure",
		"registers workers itself"} {
		if !strings.Contains(r.Err, want) {
			t.Errorf("missing %q in %q", want, r.Err)
		}
	}
	if strings.Count(r.Err, "above the team's limit") != 1 {
		t.Error("the same hint must print once")
	}
}

func TestPutOrchestrationSuspended409(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) {
		clitest.Reply(w, 409, `{"code":"Conflict","message":"suspended","requestId":"req_s","details":{"reason":"team_suspended"}}`)
	}, "orchestration", "put-orchestration", "--team-id", team, "--web-ui-exposure", "off", "--web-api-exposure", "off", "--if-match", "1")
	if r.Code != 254 || !strings.Contains(r.Err, "The team is suspended.") {
		t.Errorf("%d %q", r.Code, r.Err)
	}
}

// OR-CLI-26.6: --wait polls GetTeamOrchestration until Ready at the mutation's generation.
func TestPutOrchestrationWait(t *testing.T) {
	polls := 0
	r := clitest.Run(t, func(w http.ResponseWriter, req *http.Request) {
		switch {
		case req.Method == "PUT":
			clitest.Reply(w, 200, orch("paperclip", "private", "private", `[]`, 5, 5, 4, "Pending"))
		case req.Header.Get("If-Match") == "" && polls == 0 && req.URL.Query().Get("poll") == "":
			polls++
			w.Header().Set("ETag", `"4"`)
			clitest.Reply(w, 200, orch("paperclip", "private", "private", `[]`, 4, 4, 4, "Ready")) // the version read
		case polls == 1:
			polls++
			clitest.Reply(w, 200, orch("paperclip", "private", "private", `[]`, 5, 5, 4, "Ready")) // stale
		default:
			polls++
			clitest.Reply(w, 200, orch("paperclip", "private", "private", `[]`, 5, 5, 5, "Ready"))
		}
	}, "orchestration", "put-orchestration", "--team-id", team, "--web-ui-exposure", "private", "--web-api-exposure", "off", "--wait", "--output", "json")
	if r.Code != 0 || polls < 3 || !strings.Contains(strings.ReplaceAll(r.Out, " ", ""), `"phase":"Ready"`) {
		t.Errorf("exit %d polls %d out %q err %q", r.Code, polls, r.Out, r.Err)
	}
}

func TestPutOrchestrationAccessDenied(t *testing.T) {
	r := clitest.Run(t, clitest.Denied("collaboration:manage"), "orchestration", "put-orchestration", "--team-id", team,
		"--web-ui-exposure", "off", "--web-api-exposure", "off", "--if-match", "1")
	if r.Code != 254 || !strings.Contains(r.Err, "(AccessDenied) when calling the PutTeamOrchestration operation") || !strings.Contains(r.Err, "Required permission: collaboration:manage") {
		t.Errorf("%d %q", r.Code, r.Err)
	}
}

func changeHandler(body *string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		*body = r.Header.Get("Idempotency-Key")
		clitest.Reply(w, 202, orch("buzz", "off", "off", `["tool_not_supported"]`, 8, 8, 7, "Pending"))
	}
}

// OR-CLI-26.8: one POST with an idempotency key; the body is the tool only.
func TestChangeCollaborationTool(t *testing.T) {
	var key string
	r := clitest.Run(t, changeHandler(&key), "orchestration", "change-collaboration-tool", "--team-id", team, "--tool", "buzz", "--output", "table")
	if r.Code != 0 || len(r.Requests) != 1 || r.Req.Method != "POST" || r.Req.URL.Path != "/v1/teams/"+team+"/orchestration:change-tool" {
		t.Fatalf("%+v", r)
	}
	if string(r.ReqBody) != `{"tool":"buzz"}` || key == "" {
		t.Errorf("body %s key %q", r.ReqBody, key)
	}
	if !strings.Contains(r.Err, "Web UI: N/A for Buzz") {
		t.Errorf("stderr %q", r.Err)
	}
	r = clitest.Run(t, changeHandler(&key), "orchestration", "change-collaboration-tool", "--team-id", team, "--tool", "buzz", "--idempotency-key", "k-1")
	if key != "k-1" {
		t.Errorf("key %q", key)
	}
}

func TestChangeCollaborationToolUsageErrors(t *testing.T) {
	for name, args := range map[string][]string{"missing": {}, "unknown": {"--tool", "slack"}} {
		r := clitest.Run(t, reply("{}"), append([]string{"orchestration", "change-collaboration-tool", "--team-id", team}, args...)...)
		if r.Code != 252 || len(r.Requests) != 0 {
			t.Errorf("%s: exit %d requests %d %q", name, r.Code, len(r.Requests), r.Err)
		}
	}
}

// OR-CLI-26.8: the refusal reasons of a tool change print their hints.
func TestChangeCollaborationToolRefusals(t *testing.T) {
	for _, tc := range []struct {
		status       int
		code, reason string
		want         string
	}{
		{409, "Conflict", "tool_change_in_progress", "A tool change is already running"},
		{409, "Conflict", "orchestration_runtime_unset", "orchestration.tools"},
		{409, "Conflict", "team_suspended", "The team is suspended."},
		{400, "ValidationError", "tool_unchanged", "The team already uses that tool."},
		{400, "ValidationError", "tool_not_enabled", "That tool is not enabled for the service."},
	} {
		r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) {
			clitest.Reply(w, tc.status, fmt.Sprintf(`{"code":%q,"message":"no","requestId":"req_c","details":{"reason":%q}}`, tc.code, tc.reason))
		}, "orchestration", "change-collaboration-tool", "--team-id", team, "--tool", "buzz")
		if r.Code != 254 || !strings.Contains(r.Err, "("+tc.code+") when calling the ChangeCollaborationTool operation") || !strings.Contains(r.Err, tc.want) {
			t.Errorf("%s: exit %d %q", tc.reason, r.Code, r.Err)
		}
	}
}

func TestChangeCollaborationToolAccessDenied(t *testing.T) {
	r := clitest.Run(t, clitest.Denied("collaboration:manage"), "orchestration", "change-collaboration-tool", "--team-id", team, "--tool", "buzz")
	if r.Code != 254 || !strings.Contains(r.Err, "(AccessDenied) when calling the ChangeCollaborationTool operation") {
		t.Errorf("%d %q", r.Code, r.Err)
	}
}

func TestChangeCollaborationToolWait(t *testing.T) {
	polls := 0
	r := clitest.Run(t, func(w http.ResponseWriter, req *http.Request) {
		if req.Method == "POST" {
			clitest.Reply(w, 202, orch("buzz", "off", "off", `[]`, 8, 8, 7, "Pending"))
			return
		}
		polls++
		if polls < 2 {
			clitest.Reply(w, 200, orch("buzz", "off", "off", `[]`, 8, 8, 8, "Pending"))
			return
		}
		clitest.Reply(w, 200, orch("buzz", "off", "off", `[]`, 8, 8, 8, "Ready"))
	}, "orchestration", "change-collaboration-tool", "--team-id", team, "--tool", "buzz", "--wait", "--output", "json")
	if r.Code != 0 || polls != 2 {
		t.Errorf("exit %d polls %d err %q", r.Code, polls, r.Err)
	}
}

func TestHelpAndSkeleton(t *testing.T) {
	r := clitest.Run(t, nil, "orchestration", "help")
	if r.Code != 0 || !strings.Contains(r.Out, "describe-orchestration") || !strings.Contains(r.Out, "put-orchestration") || !strings.Contains(r.Out, "change-collaboration-tool") {
		t.Errorf("%q", r.Out)
	}
	r = clitest.Run(t, nil, "orchestration", "put-orchestration", "help")
	if strings.Contains(r.Out, "--web-ui (JSON)") || !strings.Contains(r.Out, "--web-ui-exposure (string)") {
		t.Errorf("hidden members must stay out of help: %q", r.Out)
	}
	sk := clitest.Run(t, nil, "orchestration", "put-orchestration", "--generate-cli-skeleton")
	if sk.Code != 0 || !strings.Contains(sk.Out, "webUi") || !strings.Contains(sk.Out, "webApi") {
		t.Errorf("%q", sk.Out)
	}
}
