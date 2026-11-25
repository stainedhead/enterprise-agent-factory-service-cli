package modelendpoints_test

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/clitest"
)

const (
	ep   = "mep_01JD2Q3R4S5T6V7W8X9Y0Z0001"
	team = "team_01JD2Q3R4S5T6V7W8X9Y0Z0001"
	key  = "sk-unit-test-0123456789"
)

func reply(status int, body string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) { clitest.Reply(w, status, body) }
}

func writeFile(t *testing.T, name, content string) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	return p
}

func body(t *testing.T, b []byte) map[string]any {
	t.Helper()
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("body %q: %v", b, err)
	}
	return m
}

// FR-002, MEG-CLI-1, MEG-CLI-4: example, columns, goldens and AccessDenied for the service-admin operations.
func TestCommandsDoneModelEndpoints(t *testing.T) {
	clitest.CheckOperations(t, "model-endpoints",
		"ListModelEndpoints", "CreateModelEndpoint", "GetModelEndpoint", "UpdateModelEndpoint", "DeleteModelEndpoint",
		"PutEndpointModels", "CreateEndpointModelPrice", "ListEndpointPriceProposals", "ApproveEndpointPriceProposal",
		"RejectEndpointPriceProposal", "ListEndpointTeams", "PutModelEndpointCredential", "RotateModelEndpointCredential",
		"PutEndpointFederation", "ImportProviderStatement", "ListInferenceProviders", "UpdateInferenceProvider")
}

// FR-002: the team-side operations of the service.
func TestTeamCommandsDoneModelEndpoints(t *testing.T) {
	clitest.CheckOperations(t, "team-endpoints",
		"ListTeamModelEndpoints", "CreateTeamModelEndpoint", "UpdateTeamModelEndpoint", "DeleteTeamModelEndpoint",
		"PutTeamEndpointModels", "PutTeamModelEndpointCredential")
}

// FR-003: list filters reach the query.
func TestListModelEndpointsFilters(t *testing.T) {
	r := clitest.Run(t, reply(200, `{"items":[]}`), "model-endpoints", "list-model-endpoints", "--include-teams",
		"--provider", "openai", "--location", "external", "--pricing-mode", "per_token")
	if r.Code != 0 || r.Req.URL.Path != "/v1/admin/model-endpoints" {
		t.Fatalf("%+v", r)
	}
	q := r.Req.URL.Query()
	for k, v := range map[string]string{"includeTeams": "true", "provider": "openai", "location": "external", "pricingMode": "per_token"} {
		if q.Get(k) != v {
			t.Errorf("%s=%q want %q", k, q.Get(k), v)
		}
	}
	if r := clitest.Run(t, nil, "model-endpoints", "list-model-endpoints", "--location", "moon"); r.Code != 252 {
		t.Errorf("bad enum: %+v", r)
	}
}

// FR-003: create sends the documented body; --protocols is a list, --ca-bundle reads a file, models are JSON.
func TestCreateModelEndpointBody(t *testing.T) {
	ca := writeFile(t, "ca.pem", "-----BEGIN CERTIFICATE-----\nMIIB\n-----END CERTIFICATE-----\n")
	models := writeFile(t, "models.json", `[{"modelName":"gpt-4.1","maxOutputTokens":32768}]`)
	r := clitest.Run(t, reply(201, `{"id":"`+ep+`","slug":"openai-prod","version":1}`), "model-endpoints", "create-model-endpoint",
		"--slug", "openai-prod", "--display-name", "OpenAI production", "--provider", "openai", "--base-url", "https://api.openai.com/v1",
		"--protocols", "openai_chat", "openai_responses", "--location", "external", "--pricing-mode", "per_token",
		"--auth-kind", "bearer", "--credential-mode", "per_team", "--ca-bundle", "file://"+ca, "--models", "file://"+models)
	if r.Code != 0 || r.Req.Method != "POST" || r.Req.URL.Path != "/v1/admin/model-endpoints" || r.Req.Header.Get("Idempotency-Key") == "" {
		t.Fatalf("%+v", r)
	}
	b := body(t, r.ReqBody)
	if p, _ := b["protocols"].([]any); len(p) != 2 || p[1] != "openai_responses" {
		t.Errorf("protocols %v", b["protocols"])
	}
	if !strings.HasPrefix(b["caBundle"].(string), "-----BEGIN CERTIFICATE-----") {
		t.Errorf("caBundle %v", b["caBundle"])
	}
	if m, _ := b["models"].([]any); len(m) != 1 {
		t.Errorf("models %v", b["models"])
	}
	if r := clitest.Run(t, nil, "model-endpoints", "create-model-endpoint", "--slug", "x"); r.Code != 252 || !strings.Contains(r.Err, "required") {
		t.Errorf("missing: %+v", r)
	}
}

func versioned(get string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			clitest.Reply(w, 200, get)
		case "DELETE":
			w.WriteHeader(204)
		default:
			clitest.Reply(w, 200, `{"id":"`+ep+`","version":5}`)
		}
	}
}

// FR-004: update and delete read the endpoint's version and send it as If-Match.
func TestUpdateModelEndpointReadsVersion(t *testing.T) {
	r := clitest.Run(t, versioned(`{"id":"`+ep+`","version":4}`), "model-endpoints", "update-model-endpoint", "--endpoint-id", ep, "--enabled")
	if r.Code != 0 || len(r.Requests) != 2 || r.Requests[0].URL.Path != "/v1/admin/model-endpoints/"+ep || r.Req.Method != "PATCH" || r.Req.Header.Get("If-Match") != `"4"` {
		t.Fatalf("%+v", r)
	}
	if b := body(t, r.ReqBody); b["enabled"] != true {
		t.Errorf("body %v", b)
	}
	r = clitest.Run(t, versioned(`{"id":"`+ep+`","version":4}`), "model-endpoints", "delete-model-endpoint", "--endpoint-id", ep)
	if r.Code != 0 || r.Req.Method != "DELETE" || r.Req.Header.Get("If-Match") != `"4"` {
		t.Fatalf("%+v", r)
	}
}

// FR-004: team endpoint writes read the version from the team's endpoint list.
func TestTeamEndpointWritesReadVersionFromList(t *testing.T) {
	list := `{"items":[{"id":"mep_01JD2Q3R4S5T6V7W8X9Y0Z0009","version":1},{"id":"` + ep + `","version":7}]}`
	for _, args := range [][]string{
		{"update-team-model-endpoint", "--team-id", team, "--endpoint-id", ep, "--no-enabled"},
		{"delete-team-model-endpoint", "--team-id", team, "--endpoint-id", ep},
		{"put-team-endpoint-models", "--team-id", team, "--endpoint-id", ep, "--models", `[{"modelName":"m","maxOutputTokens":10}]`},
	} {
		r := clitest.Run(t, versioned(list), append([]string{"model-endpoints"}, args...)...)
		if r.Code != 0 || r.Requests[0].URL.Path != "/v1/teams/"+team+"/model-endpoints" || r.Req.Header.Get("If-Match") != `"7"` {
			t.Errorf("%s: %+v", args[0], r)
		}
	}
}

var credentialCommands = [][]string{
	{"put-model-endpoint-credential", "--endpoint-id", ep, "--scope", team},
	{"rotate-model-endpoint-credential", "--endpoint-id", ep, "--scope", "_shared"},
	{"put-team-model-endpoint-credential", "--team-id", team, "--endpoint-id", ep},
}

const credentialAnswer = `{"endpointId":"` + ep + `","scope":"_shared","state":"set","rotatedAt":"2026-11-25T18:00:00Z"}`

// FR-005, MEG-CLI-3: the value comes from a file, standard input or a no-echo prompt, and is sent as value.
func TestCredentialValueOnlyFromFileStdinOrPrompt(t *testing.T) {
	f := writeFile(t, "key", key)
	for _, args := range credentialCommands {
		base := append([]string{"model-endpoints"}, args...)
		for name, run := range map[string]func() clitest.Result{
			"file": func() clitest.Result {
				return clitest.Run(t, reply(200, credentialAnswer), append(base, "--secret-value", "file://"+f)...)
			},
			"stdin": func() clitest.Result {
				return clitest.RunWith(t, clitest.Options{Stdin: key}, reply(200, credentialAnswer), append(base, "--secret-value", "-")...)
			},
			"prompt": func() clitest.Result {
				pw := key
				return clitest.RunWith(t, clitest.Options{Password: &pw}, reply(200, credentialAnswer), base...)
			},
		} {
			r := run()
			if r.Code != 0 || body(t, r.ReqBody)["value"] != key {
				t.Errorf("%s %s: %+v", args[0], name, r)
			}
			if strings.Contains(r.Out, key) || strings.Contains(r.Err, key) {
				t.Errorf("%s %s: value printed", args[0], name)
			}
		}
		if r := clitest.Run(t, nil, base...); r.Code != 252 || !strings.Contains(r.Err, "no terminal") {
			t.Errorf("%s without a value and a terminal: %+v", args[0], r)
		}
	}
}

// FR-005, TC-CLI-02: a plain argument and --value are refused before any call, and the value is never echoed.
func TestCredentialPlainArgumentRefused(t *testing.T) {
	noCall := func(w http.ResponseWriter, _ *http.Request) { t.Error("a call was made") }
	for _, args := range credentialCommands {
		base := append([]string{"model-endpoints"}, args...)
		for _, extra := range [][]string{{"--secret-value", key}, {"--value", key}} {
			r := clitest.Run(t, noCall, append(base, extra...)...)
			if r.Code != 252 || strings.Contains(r.Err, key) || strings.Contains(r.Out, key) {
				t.Errorf("%s %s: %+v", args[0], extra[0], r)
			}
		}
	}
}

// FR-006, TC-CLI-04: under --debug neither the value nor the API key is logged.
func TestCredentialDebugRedacted(t *testing.T) {
	for _, args := range credentialCommands {
		r := clitest.RunWith(t, clitest.Options{Stdin: key}, reply(200, credentialAnswer), append(append([]string{"model-endpoints"}, args...), "--secret-value", "-", "--debug")...)
		if r.Code != 0 || strings.Contains(r.Err, key) || strings.Contains(r.Err, "afsk_unit_test_key") {
			t.Fatalf("%s leak: %+v", args[0], r)
		}
		if !strings.Contains(r.Err, "[redacted]") {
			t.Errorf("%s: nothing shown as redacted", args[0])
		}
	}
}

// FR-007: credential writes read no version and send If-Match only when --if-match is given.
func TestCredentialSendsNoVersionRead(t *testing.T) {
	for _, args := range [][]string{credentialCommands[0], credentialCommands[2]} {
		base := append([]string{"model-endpoints"}, args...)
		r := clitest.RunWith(t, clitest.Options{Stdin: key}, reply(200, credentialAnswer), append(base, "--secret-value", "-")...)
		if r.Code != 0 || len(r.Requests) != 1 || r.Req.Method != "PUT" || r.Req.Header.Get("If-Match") != "" {
			t.Errorf("%s: %+v", args[0], r)
		}
		r = clitest.RunWith(t, clitest.Options{Stdin: key}, reply(200, credentialAnswer), append(base, "--secret-value", "-", "--if-match", "3")...)
		if r.Code != 0 || r.Req.Header.Get("If-Match") != `"3"` {
			t.Errorf("%s with --if-match: %+v", args[0], r)
		}
	}
}

func providers(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		clitest.Reply(w, 200, `{"items":[{"provider":"bedrock","version":3},{"provider":"azure","version":2}]}`)
		return
	}
	clitest.Reply(w, 200, `{"provider":"bedrock","displayName":"Amazon Bedrock","enabled":false,"teamSelection":"team_admin","routes":1,"teams":12,"directSystems":0,"fallbackModel":"azure-openai/gpt-5.1","hosts":[],"version":4,`+
		`"impact":{"routes":1,"teams":12,"agents":40,"agentsMovingToFallback":31,"agentsDegraded":2,"directSystems":0,"teamsToSelect":9,"fallbackModel":"azure-openai/gpt-5.1"}}`)
}

var bedrockOff = []string{"model-endpoints", "update-inference-provider", "--provider", "bedrock", "--no-enabled",
	"--fallback-model", "azure-openai/gpt-5.1", "--select-for-all-teams", "--reason", "No AWS Bedrock in this enterprise"}

// FR-008, TC-CLI-05: Bedrock off with a fallback: the documented body, the version read from the provider list.
func TestUpdateInferenceProviderBedrockOff(t *testing.T) {
	r := clitest.Run(t, providers, bedrockOff...)
	if r.Code != 0 || len(r.Requests) != 2 || r.Requests[0].URL.Path != "/v1/admin/inference-providers" {
		t.Fatalf("%+v", r)
	}
	if r.Req.Method != "PATCH" || r.Req.URL.Path != "/v1/admin/inference-providers/bedrock" || r.Req.Header.Get("If-Match") != `"3"` || r.Req.URL.Query().Has("dryRun") {
		t.Fatalf("%+v", r.Req)
	}
	b := body(t, r.ReqBody)
	want := map[string]any{"enabled": false, "fallbackModel": "azure-openai/gpt-5.1", "selectForAllTeams": true, "reason": "No AWS Bedrock in this enterprise"}
	if len(b) != len(want) {
		t.Errorf("body %v", b)
	}
	for k, v := range want {
		if b[k] != v {
			t.Errorf("%s = %v, want %v", k, b[k], v)
		}
	}
	if strings.Contains(r.Err, "Dry run") {
		t.Errorf("a real write said dry run: %q", r.Err)
	}
	d := clitest.Run(t, clitest.Denied("admin:inference-providers"), bedrockOff...)
	if d.Code != 254 || !strings.Contains(d.Err, "(AccessDenied) when calling the UpdateInferenceProvider operation") {
		t.Errorf("AccessDenied: %+v", d)
	}
}

// FR-008: --dry-run sends dryRun=true and prints the impact on stderr for table and text output only.
func TestUpdateInferenceProviderDryRun(t *testing.T) {
	r := clitest.Run(t, providers, append(append([]string{}, bedrockOff...), "--dry-run", "--output", "table")...)
	if r.Code != 0 || r.Req.URL.Query().Get("dryRun") != "true" {
		t.Fatalf("%+v", r)
	}
	for _, s := range []string{"Dry run: nothing was changed.", "Impact: 1 routes, 12 teams, 40 agents (31 moving to the fallback, 2 degraded), 0 direct systems, 9 teams to select.", "Fallback model: azure-openai/gpt-5.1"} {
		if !strings.Contains(r.Err, s) {
			t.Errorf("stderr lacks %q:\n%s", s, r.Err)
		}
	}
	j := clitest.Run(t, providers, append(append([]string{}, bedrockOff...), "--dry-run", "--output", "json")...)
	if j.Code != 0 || strings.Contains(j.Err, "Dry run") || !strings.Contains(j.Out, `"agentsMovingToFallback"`) {
		t.Errorf("json: %+v", j)
	}
	n := clitest.Run(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			providers(w, r)
			return
		}
		clitest.Reply(w, 200, `{"provider":"bedrock","version":3}`)
	}, append(append([]string{}, bedrockOff...), "--dry-run", "--output", "text")...)
	if n.Code != 0 || !strings.Contains(n.Err, "Dry run") || strings.Contains(n.Err, "Impact") {
		t.Errorf("no impact: %+v", n)
	}
}

const statementAnswer = `{"id":"pst_01JD2Q3R4S5T6V7W8X9Y0Z0001","endpointId":"` + ep + `","period":"2026-11","source":"opencode_csv","rowCount":1,"totalUsd":"212.40","importedAt":"2026-11-25T18:00:00Z"}`

// FR-009: a CSV statement is sent unchanged as text/csv with period and source in the query.
func TestImportStatementCSV(t *testing.T) {
	csv := "teamId,model,amountUsd,tokens\n" + team + ",Hermes-4-405B,212.40,48100000\n"
	f := writeFile(t, "statement.csv", csv)
	r := clitest.Run(t, reply(201, statementAnswer), "model-endpoints", "import-provider-statement", "--endpoint-id", ep,
		"--period", "2026-11", "--source", "opencode_csv", "--statement", "file://"+f)
	if r.Code != 0 || r.Req.Method != "POST" || r.Req.URL.Path != "/v1/admin/model-endpoints/"+ep+"/statements" {
		t.Fatalf("%+v", r)
	}
	if r.Req.Header.Get("Content-Type") != "text/csv" || string(r.ReqBody) != csv {
		t.Errorf("content type %q body %q", r.Req.Header.Get("Content-Type"), r.ReqBody)
	}
	if q := r.Req.URL.Query(); q.Get("period") != "2026-11" || q.Get("source") != "opencode_csv" {
		t.Errorf("query %v", q)
	}
}

// FR-009: a JSON statement is the body; the command line wins over the file; a rows array is the rows.
func TestImportStatementJSON(t *testing.T) {
	whole := writeFile(t, "s.json", `{"period":"2026-10","source":"nous_statement","rows":[{"teamId":"`+team+`","model":"m","amountUsd":"1.00"}]}`)
	r := clitest.Run(t, reply(201, statementAnswer), "model-endpoints", "import-provider-statement", "--endpoint-id", ep,
		"--statement", "file://"+whole, "--period", "2026-11")
	if r.Code != 0 || r.Req.Header.Get("Content-Type") != "application/json" || r.Req.URL.RawQuery != "" {
		t.Fatalf("%+v", r)
	}
	b := body(t, r.ReqBody)
	if b["period"] != "2026-11" || b["source"] != "nous_statement" || len(b["rows"].([]any)) != 1 {
		t.Errorf("body %v", b)
	}
	rows := writeFile(t, "rows.json", ` [{"teamId":"`+team+`","model":"m","amountUsd":"1.00"}]`)
	r = clitest.Run(t, reply(201, statementAnswer), "model-endpoints", "import-provider-statement", "--endpoint-id", ep,
		"--statement", "file://"+rows, "--period", "2026-11", "--source", "opencode_csv")
	if b := body(t, r.ReqBody); r.Code != 0 || b["source"] != "opencode_csv" || len(b["rows"].([]any)) != 1 {
		t.Errorf("rows array: %+v %v", r, b)
	}
}

// FR-009: anything but a file, a CSV without period and source, or a body without rows is refused with 252.
func TestImportStatementUsageErrors(t *testing.T) {
	noCall := func(w http.ResponseWriter, _ *http.Request) { t.Error("a call was made") }
	csv := writeFile(t, "s.csv", "teamId,model,amountUsd\n")
	empty := writeFile(t, "e.csv", "  \n")
	badObj := writeFile(t, "b.json", `{"period":`)
	badArr := writeFile(t, "a.json", `[{`)
	badPeriod := writeFile(t, "p.json", `{"period":202611,"source":"opencode_csv","rows":[]}`)
	base := []string{"model-endpoints", "import-provider-statement", "--endpoint-id", ep}
	for name, extra := range map[string][]string{
		"inline":          {"--statement", "teamId,model"},
		"stdin":           {"--statement", "-"},
		"missing file":    {"--statement", "file:///nonexistent/statement.csv"},
		"csv no period":   {"--statement", "file://" + csv, "--source", "opencode_csv"},
		"empty file":      {"--statement", "file://" + empty},
		"bad object":      {"--statement", "file://" + badObj},
		"bad array":       {"--statement", "file://" + badArr, "--period", "2026-11", "--source", "opencode_csv"},
		"period number":   {"--statement", "file://" + badPeriod},
		"no statement":    {"--period", "2026-11", "--source", "opencode_csv"},
		"bad source enum": {"--statement", "file://" + csv, "--period", "2026-11", "--source", "pdf"},
	} {
		if r := clitest.Run(t, noCall, append(append([]string{}, base...), extra...)...); r.Code != 252 {
			t.Errorf("%s: %+v", name, r)
		}
	}
}
