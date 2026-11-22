package collaboration_test

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/clitest"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/golden"
)

const team = "team_01JD2Q3R4S5T6V7W8X9Y0Z0001"

// collab builds a TeamCollaboration document.
func collab(host string, webUI, restAPI bool, reasons string, version, generation, observed int, phase string) string {
	h := "null"
	if host != "" {
		h = `"` + host + `"`
	}
	return fmt.Sprintf(`{"teamId":%q,"agentId":"agt_01JD2Q3R4S5T6V7W8X9Y0Z0001","hostname":%s,"webUiUrl":null,"restApiUrl":null,`+
		`"webUiEnabled":%t,"restApiEnabled":%t,"registrationMode":"manual","effective":{"webUi":false,"restApi":false,"reasons":%s},`+
		`"unmanagedPaperclipAgents":[],"version":%d,"generation":%d,"status":{"phase":%q,"observedGeneration":%d}}`,
		team, h, webUI, restAPI, reasons, version, generation, phase, observed)
}

func TestCommandsDoneCollaboration(t *testing.T) { // CLI-US-24.10, 26.1, 34.2
	clitest.CheckOperations(t, "collaboration", "GetTeamCollaboration", "PutTeamCollaboration",
		"ListCollaborationLinks", "ConnectAgentToCollaboration", "DisconnectAgentFromCollaboration", "RotateCollaborationKeys")
}

func reply(body string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) { clitest.Reply(w, 200, body) }
}

// CLI-US-24.1: one GET, the JSON printed unchanged.
func TestDescribeCollaborationJSONUnchanged(t *testing.T) {
	body := collab("pay.collab.example.com", true, false, `[]`, 4, 4, 4, "Ready")
	r := clitest.Run(t, reply(body), "collaboration", "describe-collaboration", "--team-id", team, "--output", "json")
	if r.Code != 0 || len(r.Requests) != 1 || r.Req.Method != "GET" || r.Req.URL.Path != "/v1/teams/"+team+"/collaboration" {
		t.Fatalf("%+v", r)
	}
	if strings.ReplaceAll(strings.TrimSpace(r.Out), " ", "") == "" || !strings.Contains(strings.ReplaceAll(r.Out, " ", ""), `"hostname":"pay.collab.example.com"`) {
		t.Errorf("out %q", r.Out)
	}
}

// CLI-US-24.1: table and text columns, a null hostname prints -.
func TestDescribeCollaborationTableGolden(t *testing.T) {
	body := collab("", true, false, `[]`, 4, 4, 4, "Ready")
	for _, f := range []string{"table", "text"} {
		r := clitest.Run(t, reply(body), "collaboration", "describe-collaboration", "--team-id", team, "--output", f)
		if r.Code != 0 {
			t.Fatal(r.Err)
		}
		golden.Check(t, "describe-collaboration."+f+".golden", r.Out)
	}
}

// CLI-US-24.2: the not-effective notice goes to stderr for table and text, never for json; exit 0.
func TestDescribeCollaborationNotEffectiveNotice(t *testing.T) {
	body := collab("", true, true, `["base_domain_unset","team_not_allowed"]`, 1, 1, 1, "Ready")
	r := clitest.Run(t, reply(body), "collaboration", "describe-collaboration", "--team-id", team, "--output", "table")
	want := "Not effective: web UI (base_domain_unset, team_not_allowed)\nNot effective: REST API (base_domain_unset, team_not_allowed)\n"
	if r.Code != 0 || !strings.Contains(r.Err, want) {
		t.Fatalf("exit %d stderr %q", r.Code, r.Err)
	}
	if j := clitest.Run(t, reply(body), "collaboration", "describe-collaboration", "--team-id", team, "--output", "json"); j.Err != "" {
		t.Errorf("json printed %q", j.Err)
	}
	none := collab("h.example.com", true, true, `[]`, 1, 1, 1, "Ready")
	if r := clitest.Run(t, reply(none), "collaboration", "describe-collaboration", "--team-id", team, "--output", "text"); r.Err != "" {
		t.Errorf("no reasons must print nothing: %q", r.Err)
	}
}

// CLI-US-24.9: no collaboration agent: AWS-style NotFound plus the create hint, exit 254.
func TestDescribeCollaborationMissingAgentHint(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) {
		clitest.Reply(w, 404, `{"code":"NotFound","message":"no collaboration agent","requestId":"req_1","details":{"reason":"collaboration_agent_missing"}}`)
	}, "collaboration", "describe-collaboration", "--team-id", team)
	if r.Code != 254 || !strings.Contains(r.Err, "(NotFound) when calling the GetTeamCollaboration operation") ||
		!strings.Contains(r.Err, "This team has no collaboration agent. Create one with: afs agents create-agent --role collaboration") {
		t.Errorf("exit %d %q", r.Code, r.Err)
	}
}

// CLI-US-24.8
func TestDescribeCollaborationAccessDenied254(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) {
		clitest.Reply(w, 404, `{"code":"NotFound","message":"not found","requestId":"req_1"}`)
	}, "collaboration", "describe-collaboration", "--team-id", team)
	if r.Code != 254 || !strings.Contains(r.Err, "(NotFound)") {
		t.Errorf("no binding: %d %q", r.Code, r.Err)
	}
	d := clitest.Run(t, clitest.Denied("collaboration:read"), "collaboration", "describe-collaboration", "--team-id", team)
	if d.Code != 254 || !strings.Contains(d.Err, "Required permission: collaboration:read") {
		t.Errorf("denied: %d %q", d.Code, d.Err)
	}
}

func putHandler(version int, got *[]string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		*got = append(*got, r.Method)
		if r.Method == "GET" {
			w.Header().Set("ETag", fmt.Sprintf(`"%d"`, version))
			clitest.Reply(w, 200, collab("h.example.com", true, false, `[]`, version, version, version, "Ready"))
			return
		}
		clitest.Reply(w, 200, collab("h.example.com", true, false, `[]`, version+1, version+1, version, "Pending"))
	}
}

// CLI-US-24.3
func TestPutCollaborationBody(t *testing.T) {
	var m []string
	r := clitest.Run(t, putHandler(4, &m), "collaboration", "put-collaboration", "--team-id", team,
		"--web-ui-enabled", "--no-rest-api-enabled", "--registration-mode", "managed", "--if-match", "4")
	if r.Code != 0 || r.Req.Method != "PUT" || r.Req.URL.Path != "/v1/teams/"+team+"/collaboration" {
		t.Fatalf("%+v", r)
	}
	if got := strings.ReplaceAll(string(r.ReqBody), " ", ""); got != `{"registrationMode":"managed","restApiEnabled":false,"webUiEnabled":true}` &&
		got != `{"webUiEnabled":true,"restApiEnabled":false,"registrationMode":"managed"}` {
		t.Errorf("body %s", r.ReqBody)
	}
	r = clitest.Run(t, putHandler(4, &m), "collaboration", "put-collaboration", "--team-id", team,
		"--web-ui-enabled", "--no-rest-api-enabled", "--if-match", "4")
	if strings.Contains(string(r.ReqBody), "registrationMode") {
		t.Errorf("registrationMode must be omitted: %s", r.ReqBody)
	}
}

func TestPutCollaborationRegistrationModeEnum252(t *testing.T) {
	r := clitest.Run(t, reply("{}"), "collaboration", "put-collaboration", "--team-id", team,
		"--web-ui-enabled", "--rest-api-enabled", "--registration-mode", "auto")
	if r.Code != 252 || len(r.Requests) != 0 {
		t.Errorf("exit %d requests %d %q", r.Code, len(r.Requests), r.Err)
	}
}

// CLI-US-24.4
func TestPutCollaborationBothSurfacesRequired252(t *testing.T) {
	for name, args := range map[string][]string{
		"rest-api-enabled": {"--web-ui-enabled"},
		"web-ui-enabled":   {"--no-rest-api-enabled"},
	} {
		r := clitest.Run(t, reply("{}"), append([]string{"collaboration", "put-collaboration", "--team-id", team}, args...)...)
		if r.Code != 252 || len(r.Requests) != 0 || !strings.Contains(r.Err, name) {
			t.Errorf("%s: exit %d requests %d %q", name, r.Code, len(r.Requests), r.Err)
		}
	}
}

func TestPutCollaborationInputJSON(t *testing.T) {
	dir := t.TempDir()
	f := dir + "/collab.json"
	writeFile(t, f, `{"teamId":"`+team+`","webUiEnabled":false,"restApiEnabled":true,"version":4}`)
	var m []string
	r := clitest.Run(t, putHandler(4, &m), "collaboration", "put-collaboration", "--cli-input-json", "file://"+f)
	if r.Code != 0 || !strings.Contains(string(r.ReqBody), `"restApiEnabled":true`) || !strings.Contains(string(r.ReqBody), `"webUiEnabled":false`) {
		t.Errorf("exit %d body %s err %s", r.Code, r.ReqBody, r.Err)
	}
}

// CLI-US-24.5
func TestPutCollaborationIfMatchGuarded(t *testing.T) {
	var m []string
	r := clitest.Run(t, putHandler(4, &m), "collaboration", "put-collaboration", "--team-id", team, "--web-ui-enabled", "--rest-api-enabled", "--if-match", "4")
	if r.Code != 0 || len(r.Requests) != 1 || r.Req.Header.Get("If-Match") != `"4"` || strings.Contains(r.Err, "Applied to version") {
		t.Errorf("%+v", r)
	}
}

func TestPutCollaborationReadsVersion(t *testing.T) {
	var m []string
	r := clitest.Run(t, putHandler(7, &m), "collaboration", "put-collaboration", "--team-id", team, "--web-ui-enabled", "--rest-api-enabled")
	if r.Code != 0 || strings.Join(m, ",") != "GET,PUT" || r.Req.Header.Get("If-Match") != `"7"` {
		t.Fatalf("%v %+v", m, r)
	}
	if !strings.Contains(r.Err, "Applied to version 7; use --if-match to guard against concurrent changes.") {
		t.Errorf("stderr %q", r.Err)
	}
}

func TestPutCollaboration412Message(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("ETag", `"5"`)
		clitest.Reply(w, 412, `{"code":"PreconditionFailed","message":"stale","requestId":"req_4"}`)
	}, "collaboration", "put-collaboration", "--team-id", team, "--web-ui-enabled", "--rest-api-enabled", "--if-match", "4")
	if r.Code != 254 || !strings.Contains(r.Err, "The resource changed (version 5). Re-run with --if-match 5 to apply your change to the latest version.") {
		t.Errorf("exit %d %q", r.Code, r.Err)
	}
}

// CLI-US-24.5 and PC-CLI-D2: a 404 on the version read sends no If-Match.
func TestPutCollaborationMissingAgent409(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, req *http.Request) {
		if req.Method == "GET" {
			clitest.Reply(w, 404, `{"code":"NotFound","message":"none","requestId":"req_g","details":{"reason":"collaboration_agent_missing"}}`)
			return
		}
		clitest.Reply(w, 409, `{"code":"Conflict","message":"no collaboration agent","requestId":"req_p","details":{"reason":"collaboration_agent_missing"}}`)
	}, "collaboration", "put-collaboration", "--team-id", team, "--web-ui-enabled", "--rest-api-enabled")
	if r.Code != 254 || r.Req.Header.Get("If-Match") != "" || !strings.Contains(r.Err, "(Conflict) when calling the PutTeamCollaboration operation") ||
		!strings.Contains(r.Err, "This team has no collaboration agent.") {
		t.Errorf("exit %d if-match %q %q", r.Code, r.Req.Header.Get("If-Match"), r.Err)
	}
}

// CLI-US-24.7
func TestPutCollaborationValidationFieldsPrinted(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, req *http.Request) {
		clitest.Reply(w, 400, `{"code":"ValidationError","message":"invalid","requestId":"req_v","details":{"fields":[{"field":"webUiEnabled","reason":"collab_base_domain_unset"},{"field":"restApiEnabled","reason":"collab_public_endpoints_off"}]}}`)
	}, "collaboration", "put-collaboration", "--team-id", team, "--web-ui-enabled", "--rest-api-enabled", "--if-match", "1")
	if r.Code != 254 {
		t.Fatal(r.Code)
	}
	for _, want := range []string{"Request ID: req_v\nfield: webUiEnabled reason: collab_base_domain_unset\nfield: restApiEnabled reason: collab_public_endpoints_off\n",
		"Ask a service admin to set collab.base_domain (afs admin put-setting).", "Ask a service admin to switch collab.public_endpoints on"} {
		if !strings.Contains(r.Err, want) {
			t.Errorf("missing %q in %q", want, r.Err)
		}
	}
}

func TestPutCollaborationAccessDenied254(t *testing.T) {
	r := clitest.Run(t, clitest.Denied("collaboration:manage"), "collaboration", "put-collaboration", "--team-id", team,
		"--web-ui-enabled", "--rest-api-enabled", "--if-match", "1")
	if r.Code != 254 || !strings.Contains(r.Err, "(AccessDenied) when calling the PutTeamCollaboration operation") || !strings.Contains(r.Err, "Required permission: collaboration:manage") {
		t.Errorf("%d %q", r.Code, r.Err)
	}
}

func TestPutCollaborationSuspended409(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) {
		clitest.Reply(w, 409, `{"code":"Conflict","message":"suspended","requestId":"req_s","details":{"reason":"team_suspended"}}`)
	}, "collaboration", "put-collaboration", "--team-id", team, "--web-ui-enabled", "--rest-api-enabled", "--if-match", "1")
	if r.Code != 254 || !strings.Contains(r.Err, "The team is suspended.") {
		t.Errorf("%d %q", r.Code, r.Err)
	}
}

// waitServer answers the PUT with generation 5 and each poll with the next phase of script.
func waitServer(script []string, polls *int) http.HandlerFunc {
	var mu sync.Mutex
	return func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		if r.Method == "PUT" {
			clitest.Reply(w, 200, collab("h.example.com", true, false, `[]`, 5, 5, 4, "Pending"))
			return
		}
		if r.Header.Get("If-Match") == "" && *polls == 0 && len(script) > 0 && script[0] == "version" {
			script = script[1:]
			w.Header().Set("ETag", `"4"`)
			clitest.Reply(w, 200, collab("h.example.com", true, false, `[]`, 4, 4, 4, "Ready"))
			return
		}
		i := *polls
		if i >= len(script) {
			i = len(script) - 1
		}
		*polls++
		clitest.Reply(w, 200, script[i])
	}
}

// CLI-US-24.6
func TestPutCollaborationWaitReady(t *testing.T) {
	var polls int
	script := []string{"version",
		collab("h.example.com", true, false, `[]`, 5, 5, 5, "Pending"),
		collab("h.example.com", true, false, `[]`, 5, 5, 5, "Ready")}
	r := clitest.Run(t, waitServer(script, &polls), "collaboration", "put-collaboration", "--team-id", team,
		"--web-ui-enabled", "--no-rest-api-enabled", "--wait", "--output", "json")
	if r.Code != 0 || polls != 2 || !strings.Contains(r.Out, `"phase": "Ready"`) && !strings.Contains(r.Out, `"phase":"Ready"`) {
		t.Errorf("exit %d polls %d out %q err %q", r.Code, polls, r.Out, r.Err)
	}
}

func TestPutCollaborationWaitIgnoresStaleReady(t *testing.T) {
	var polls int
	script := []string{"version",
		collab("h.example.com", true, false, `[]`, 5, 5, 4, "Ready"), // stale: observedGeneration 4 < generation 5
		collab("h.example.com", true, false, `[]`, 5, 5, 5, "Ready")}
	r := clitest.Run(t, waitServer(script, &polls), "collaboration", "put-collaboration", "--team-id", team,
		"--web-ui-enabled", "--no-rest-api-enabled", "--wait")
	if r.Code != 0 || polls != 2 {
		t.Errorf("exit %d polls %d err %q", r.Code, polls, r.Err)
	}
}

func TestPutCollaborationWaitDegraded254(t *testing.T) {
	var polls int
	deg := strings.Replace(collab("h.example.com", true, false, `[]`, 5, 5, 5, "Degraded"), `"observedGeneration":5}`,
		`"observedGeneration":5,"conditions":[{"type":"EndpointReady","status":"False"},{"type":"PaperclipReady","status":"True"}]}`, 1)
	r := clitest.Run(t, waitServer([]string{"version", deg}, &polls), "collaboration", "put-collaboration", "--team-id", team,
		"--web-ui-enabled", "--no-rest-api-enabled", "--wait")
	if r.Code != 254 || !strings.Contains(r.Err, "EndpointReady") || !strings.Contains(r.Err, "PaperclipReady") {
		t.Errorf("exit %d err %q", r.Code, r.Err)
	}
}

func TestPutCollaborationWaitTimeout255(t *testing.T) {
	var polls int
	clock := time.Date(2026, 11, 22, 0, 0, 0, 0, time.UTC)
	r := clitest.RunWith(t, clitest.Options{
		Now:   func() time.Time { return clock },
		Sleep: func(_ context.Context, d time.Duration) error { clock = clock.Add(d); return nil },
	}, waitServer([]string{"version", collab("h.example.com", true, false, `[]`, 5, 5, 4, "Pending")}, &polls),
		"collaboration", "put-collaboration", "--team-id", team, "--web-ui-enabled", "--no-rest-api-enabled", "--wait", "--wait-timeout", "30s")
	if r.Code != 255 || !strings.Contains(r.Err, "the change stands") {
		t.Errorf("exit %d err %q", r.Code, r.Err)
	}
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}
