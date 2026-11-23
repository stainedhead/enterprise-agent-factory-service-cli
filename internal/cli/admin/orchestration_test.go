package admin_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/clitest"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/golden"
)

const orchEndpoint = `{"teamId":"team_1","slug":"payments","teamState":"active","agentId":"agt_1","tool":"paperclip","surface":"webUi","exposure":"public","effective":"private","reasons":["team_limit"],"limit":"private","hostname":"payments.collab.corp.example.com","phase":"Ready","lastAccessAt":null}`

func migration(state string) string {
	return `"migration":{"state":"` + state + `","teamsTotal":12,"teamsWithOrchestration":9,"teamsPending":3,"blocked":[{"teamId":"team_7","reason":"orchestration_runtime_unset"}]}`
}

func TestCommandsDoneOrchestrationAdmin(t *testing.T) { // OR-CLI-26, OR-CLI-28
	clitest.CheckOperations(t, "orchestration", "ListOrchestrationEndpoints", "PutServiceModelAccess")
}

// ORC-25: the filters are query parameters; the profile team_id never fills teamId; the enums are checked first.
func TestListOrchestrationEndpointsFilters(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) {
		clitest.Reply(w, 200, `{"items":[`+orchEndpoint+`],`+migration("complete")+`}`)
	}, "admin", "list-orchestration-endpoints", "--exposure", "public", "--tool", "buzz", "--team-id", "team_1", "--state", "Ready")
	q := r.Req.URL.Query()
	if r.Code != 0 || r.Req.URL.Path != "/v1/admin/orchestration-endpoints" || q.Get("exposure") != "public" || q.Get("tool") != "buzz" || q.Get("teamId") != "team_1" || q.Get("state") != "Ready" {
		t.Errorf("%d %s %v", r.Code, r.Req.URL.Path, q)
	}
	for _, bad := range [][]string{{"--exposure", "off"}, {"--tool", "slack"}, {"--state", "Bad"}} {
		if r := clitest.Run(t, nil, append([]string{"admin", "list-orchestration-endpoints"}, bad...)...); r.Code != 252 || len(r.Requests) != 0 {
			t.Errorf("%v: %d", bad, r.Code)
		}
	}
	home := t.TempDir()
	writeProfile(t, home, "team_01J9ZQ3ABCDEFGHJKMNPQRSTVW")
	n := clitest.RunWith(t, clitest.Options{Home: home}, func(w http.ResponseWriter, _ *http.Request) {
		clitest.Reply(w, 200, `{"items":[],`+migration("complete")+`}`)
	},
		"admin", "list-orchestration-endpoints")
	if n.Code != 0 || n.Req.URL.Query().Has("teamId") {
		t.Errorf("the profile team_id must not fill the query filter: %d %v", n.Code, n.Req.URL.Query())
	}
}

func TestOrchestrationEndpointsTableGolden(t *testing.T) {
	for _, f := range []string{"table", "text"} {
		r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) {
			clitest.Reply(w, 200, `{"items":[`+orchEndpoint+`],`+migration("complete")+`}`)
		}, "admin", "list-orchestration-endpoints", "--output", f)
		golden.Check(t, "orchestration-endpoints."+f+".golden", r.Out)
	}
}

// ORC-25: the backfill progress goes to stderr for table and text while it is not complete; never for json or csv.
func TestOrchestrationEndpointsMigrationNotice(t *testing.T) {
	h := func(state string) http.HandlerFunc {
		return func(w http.ResponseWriter, _ *http.Request) {
			clitest.Reply(w, 200, `{"items":[`+orchEndpoint+`],`+migration(state)+`}`)
		}
	}
	r := clitest.Run(t, h("blocked"), "admin", "list-orchestration-endpoints", "--output", "table")
	want := "Orchestration agents: blocked (9 of 12 teams have one, 3 pending).\nBlocked: team_7 (orchestration_runtime_unset)\n"
	if r.Code != 0 || r.Err != want {
		t.Errorf("exit %d err %q", r.Code, r.Err)
	}
	if c := clitest.Run(t, h("complete"), "admin", "list-orchestration-endpoints", "--output", "text"); c.Err != "" {
		t.Errorf("complete prints nothing: %q", c.Err)
	}
	if j := clitest.Run(t, h("running"), "admin", "list-orchestration-endpoints", "--output", "json"); j.Err != "" {
		t.Errorf("json printed %q", j.Err)
	}
	if n := clitest.Run(t, ok, "admin", "list-orchestration-endpoints", "--output", "table"); n.Err != "" {
		t.Errorf("no migration member prints nothing: %q", n.Err)
	}
}

func TestListOrchestrationEndpointsCSV(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, req *http.Request) {
		if req.Header.Get("Accept") != "text/csv" {
			t.Errorf("Accept %q", req.Header.Get("Accept"))
		}
		w.Header().Set("Content-Type", "text/csv")
		_, _ = w.Write([]byte("teamId,slug\r\nteam_1,payments\r\n"))
	}, "admin", "list-orchestration-endpoints", "--output", "csv")
	if r.Code != 0 || !strings.Contains(r.Out, "team_1,payments") {
		t.Errorf("%d %q %q", r.Code, r.Out, r.Err)
	}
}

// ORC-30: the 1.0.1 endpoint list keeps working and says what replaces it.
func TestListCollaborationEndpointsDeprecated(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) { clitest.Reply(w, 200, `{"items":[]}`) }, "admin", "list-collaboration-endpoints")
	want := "Warning: afs admin list-collaboration-endpoints is deprecated and will be removed in 2.0.0; use afs admin list-orchestration-endpoints instead.\n"
	if r.Code != 0 || r.Err != want {
		t.Errorf("exit %d err %q", r.Code, r.Err)
	}
	h := clitest.Run(t, nil, "admin", "list-collaboration-endpoints", "help")
	if !strings.Contains(h.Out, "Deprecated: kept until 2.0.0; use afs admin list-orchestration-endpoints instead.") {
		t.Errorf("%q", h.Out)
	}
}

// MOD-1, MOD-8: --access is required and is all or selected; the body is the access only; no version is read.
func TestPutServiceModelAccess(t *testing.T) {
	h := func(w http.ResponseWriter, _ *http.Request) {
		clitest.Reply(w, 200, `{"access":"selected","enabledModels":3,"callableModels":2,"version":2}`)
	}
	r := clitest.Run(t, h, "admin", "put-service-model-access", "--access", "selected", "--output", "table")
	if r.Code != 0 || len(r.Requests) != 1 || r.Req.Method != "PUT" || r.Req.URL.Path != "/v1/admin/model-access" || string(r.ReqBody) != `{"access":"selected"}` {
		t.Fatalf("%+v", r)
	}
	if !strings.Contains(r.Out, "CallableModels") || !strings.Contains(r.Out, "selected") {
		t.Errorf("%q", r.Out)
	}
	if g := clitest.Run(t, h, "admin", "put-service-model-access", "--access", "all", "--if-match", "2"); g.Req.Header.Get("If-Match") != `"2"` {
		t.Errorf("If-Match %q", g.Req.Header.Get("If-Match"))
	}
	for name, args := range map[string][]string{"missing": {}, "bad": {"--access", "some"}} {
		if e := clitest.Run(t, h, append([]string{"admin", "put-service-model-access"}, args...)...); e.Code != 252 || len(e.Requests) != 0 {
			t.Errorf("%s: %d", name, e.Code)
		}
	}
	d := clitest.Run(t, clitest.Denied("admin:models"), "admin", "put-service-model-access", "--access", "all")
	if d.Code != 254 || !strings.Contains(d.Err, "Required permission: admin:models") {
		t.Errorf("%d %q", d.Code, d.Err)
	}
}

// MOD-1: while service access is all, a model's enabled flag cannot change; the hint says how to switch.
func TestUpdateModelRefusedWhileServiceAccessAll(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) {
		clitest.Reply(w, 409, `{"code":"Conflict","message":"service access is all","requestId":"req_a","details":{"reason":"service_model_access_all"}}`)
	}, "admin", "update-model", "--model-id", "m1", "--no-enabled", "--if-match", "1")
	if r.Code != 254 || !strings.Contains(r.Err, "afs admin put-service-model-access --access selected") {
		t.Errorf("%d %q", r.Code, r.Err)
	}
}

// ORC-27: the new settings are known to help, are sent as given, and a renamed 1.0.1 key is the API's refusal.
func TestPutSettingOrchestrationKeys(t *testing.T) {
	for key, val := range map[string]string{
		"orchestration.public_base_domain":      `"collab.example.com"`,
		"orchestration.private_base_domain":     `"collab.corp.example.com"`,
		"orchestration.default_limit":           `"private"`,
		"orchestration.default_web_ui_exposure": `"private"`,
		"orchestration.default_volume_gib":      `20`,
		"orchestration.tools":                   `{"paperclip":{"enabled":true,"default":true}}`,
		"models.service_access":                 `"all"`,
		"models.default_model_id":               `"anthropic.claude-sonnet-4-5"`,
		"models.default_team_models":            `["m1"]`,
	} {
		r := putSetting(t, ok, "--key", key, "--value", val, "--if-match", "1")
		if r.Code != 0 || string(r.ReqBody) != `{"value":`+val+`}` || r.Req.URL.Path != "/v1/admin/settings/"+key {
			t.Errorf("%s: %d %s %s", key, r.Code, r.ReqBody, r.Err)
		}
	}
	h := clitest.Run(t, nil, "admin", "put-setting", "help")
	if !strings.Contains(h.Out, "Known values:") || !strings.Contains(h.Out, "orchestration.default_limit") || strings.Contains(h.Out, "collab.base_domain,") || strings.Contains(h.Out, "collab.public_endpoints,") {
		t.Errorf("help must name the new keys and not the renamed ones: %q", h.Out)
	}
	r := putSetting(t, clitest.Fail(400, "ValidationError", "renamed"), "--key", "collab.base_domain", "--value", `"x.example.com"`, "--if-match", "1")
	if r.Code != 254 || !strings.Contains(r.Err, "(ValidationError) when calling the PutSetting operation") {
		t.Errorf("%d %q", r.Code, r.Err)
	}
	rn := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) {
		clitest.Reply(w, 400, `{"code":"ValidationError","message":"renamed","requestId":"req_s","details":{"reason":"setting_renamed"}}`)
	}, "admin", "put-setting", "--key", "collab.base_domain", "--value", `"x.example.com"`, "--if-match", "1")
	if !strings.Contains(rn.Err, "collab.base_domain is now orchestration.public_base_domain") {
		t.Errorf("%q", rn.Err)
	}
}

// MOD-7: the admin model list shows callable, price, origin and the number of teams with access.
func TestAdminModelsColumns(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) {
		clitest.Reply(w, 200, `{"items":[{"id":"m1","provider":"anthropic","displayName":"M1","enabled":true,"callable":false,"priceState":"pending","origin":"bedrock_sync","teamsWithAccess":4,"contextWindow":200000,"version":1}]}`)
	}, "admin", "list-models", "--output", "table")
	for _, want := range []string{"Callable", "Price", "Origin", "Teams", "bedrock", "pending"} {
		if !strings.Contains(r.Out, want) {
			t.Errorf("missing %q in %q", want, r.Out)
		}
	}
}
