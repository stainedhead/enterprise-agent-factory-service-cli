package mcp_test

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
	srv   = "mcp_01JD2Q3R4S5T6V7W8X9Y0Z0001"
	team  = "team_01JD2Q3R4S5T6V7W8X9Y0Z0001"
	agent = "agt_01JD2Q3R4S5T6V7W8X9Y0Z0001"
	token = "ghp_unit_test_0123456789" //nolint:gosec // G101: a test fixture, not a credential
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

// FR-002: the mcp-servers service: example, columns, goldens and AccessDenied for each operation.
func TestCommandsDoneMcpServers(t *testing.T) {
	clitest.CheckOperations(t, "mcp-servers",
		"ListMcpServers", "CreateMcpServer", "ImportMcpServerManifest", "GetMcpServer", "UpdateMcpServer", "DeleteMcpServer",
		"PutMcpServerCredential", "PutMcpPolicy", "RefreshMcpInventory", "ListMcpReviews", "ListMcpPrimitives",
		"GetMcpPrimitiveChange", "ApproveMcpPrimitive", "RejectMcpPrimitive", "ApproveMcpServer", "RejectMcpServer",
		"RequestMcpServerChanges", "PublishMcpServer", "SuspendMcpServer", "ResumeMcpServer", "RetireMcpServer",
		"GetMcpQueryCapture", "StartMcpQueryCapture", "StopMcpQueryCapture")
}

// FR-003: the team's MCP commands in the teams service.
func TestCommandsDoneTeamMcp(t *testing.T) {
	clitest.CheckOperations(t, "team-mcp",
		"ListTeamMcpServers", "GetTeamMcpServer", "PutTeamMcpServer", "CreateTeamMcpServer", "UpdateTeamMcpServer",
		"DeleteTeamMcpServer", "SubmitMcpServerForReview", "PutTeamMcpServerCredential", "PutTeamMcpPolicy",
		"RefreshTeamMcpInventory", "ListTeamMcpPrimitives", "ApproveTeamMcpPrimitive", "RejectTeamMcpPrimitive",
		"ApproveTeamMcpServer", "RejectTeamMcpServer", "StartMcpOAuthConnect", "DisconnectMcpOAuth", "TestMcpFederation")
}

// FR-004: MCP access and federated identities of agents and pools, MCP events and usage reports.
func TestCommandsDoneMcpAccess(t *testing.T) {
	clitest.CheckOperations(t, "mcp-access",
		"GetEffectiveAgentMcp", "ListAgentFederatedIdentities", "PutAgentFederatedIdentity", "DeleteAgentFederatedIdentity",
		"GetEffectivePoolMcp", "ListPoolFederatedIdentities", "PutPoolFederatedIdentity", "DeletePoolFederatedIdentity",
		"ListMcpEvents", "ListTeamMcpEvents", "GetMcpUsageReport", "GetTeamMcpUsageReport")
}

// FR-005: list filters reach the query; a bad enum is a usage error before any call.
func TestListMcpServersFilters(t *testing.T) {
	r := clitest.Run(t, reply(200, `{"items":[]}`), "mcp-servers", "list-mcp-servers", "--scope", "team",
		"--review-status", "pending_review", "--team-id", team, "--location", "corporate")
	if r.Code != 0 || r.Req.URL.Path != "/v1/mcp-servers" {
		t.Fatalf("%+v", r)
	}
	q := r.Req.URL.Query()
	for k, v := range map[string]string{"scope": "team", "reviewStatus": "pending_review", "teamId": team, "location": "corporate"} {
		if q.Get(k) != v {
			t.Errorf("%s=%q want %q", k, q.Get(k), v)
		}
	}
	if r := clitest.Run(t, nil, "mcp-servers", "list-mcp-servers", "--scope", "galaxy"); r.Code != 252 {
		t.Errorf("bad enum: %+v", r)
	}
}

// FR-006: create sends the documented body; --hosted reads JSON from a file.
func TestCreateMcpServerBody(t *testing.T) {
	hosted := writeFile(t, "hosted.json", `{"imageDigest":"registry.example.com/jira@sha256:0123","args":["--stdio"],"envNames":["JIRA_TOKEN"]}`)
	r := clitest.Run(t, reply(201, `{"id":"`+srv+`","slug":"jira","version":1}`), "mcp-servers", "create-mcp-server",
		"--slug", "jira", "--display-name", "Jira", "--data-classification", "confidential", "--transport", "hosted",
		"--location", "in_namespace", "--auth-kind", "static", "--hosted", "file://"+hosted)
	if r.Code != 0 || r.Req.Method != "POST" || r.Req.URL.Path != "/v1/mcp-servers" {
		t.Fatalf("%+v", r)
	}
	b := body(t, r.ReqBody)
	h, _ := b["hosted"].(map[string]any)
	if b["slug"] != "jira" || b["transport"] != "hosted" || h == nil || h["imageDigest"] == nil {
		t.Fatalf("body %v", b)
	}
	if r.Req.Header.Get("Idempotency-Key") == "" {
		t.Error("no idempotency key")
	}
}

// FR-007: a credential's value and client secret come from a file or standard input, never inline, and are never echoed.
func TestMcpCredentialSecretsNeverInline(t *testing.T) {
	f := writeFile(t, "tok", token)
	r := clitest.Run(t, reply(200, `{"scope":"`+team+`","state":"set","scopes":[]}`), "mcp-servers", "put-mcp-server-credential",
		"--mcp-server-id", srv, "--scope", team, "--secret-value", "file://"+f)
	if r.Code != 0 || body(t, r.ReqBody)["value"] != token {
		t.Fatalf("%+v %s", r, r.ReqBody)
	}
	if strings.Contains(r.Out+r.Err, token) {
		t.Error("the value was echoed")
	}
	if r := clitest.Run(t, nil, "mcp-servers", "put-mcp-server-credential", "--mcp-server-id", srv, "--scope", team, "--secret-value", token); r.Code != 252 {
		t.Errorf("inline value accepted: %+v", r)
	}
	if r := clitest.Run(t, nil, "mcp-servers", "put-mcp-server-credential", "--mcp-server-id", srv, "--scope", "_shared", "--client-id", "afs", "--client-secret", token); r.Code != 252 {
		t.Errorf("inline client secret accepted: %+v", r)
	}
	r = clitest.RunWith(t, clitest.Options{Stdin: token}, reply(200, `{"scope":"_shared","state":"set","scopes":[]}`), "mcp-servers",
		"put-mcp-server-credential", "--mcp-server-id", srv, "--scope", "_shared", "--client-id", "afs", "--client-secret", "-")
	if r.Code != 0 || body(t, r.ReqBody)["clientSecret"] != token || body(t, r.ReqBody)["clientId"] != "afs" {
		t.Fatalf("stdin client secret: %+v %s", r, r.ReqBody)
	}
}

// FR-008: update reads the server's version first and sends If-Match; a team server's version is read from "server".
func TestMcpUpdatesReadVersion(t *testing.T) {
	h := func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "GET" && r.URL.Path == "/v1/mcp-servers/"+srv:
			clitest.Reply(w, 200, `{"id":"`+srv+`","version":4}`)
		case r.Method == "GET":
			clitest.Reply(w, 200, `{"server":{"id":"`+srv+`","version":6},"enabled":true}`)
		default:
			clitest.Reply(w, 200, `{"id":"`+srv+`","version":7}`)
		}
	}
	r := clitest.Run(t, h, "mcp-servers", "update-mcp-server", "--mcp-server-id", srv, "--drift-policy", "auto_approve_compatible")
	if r.Code != 0 || len(r.Requests) != 2 || r.Requests[1].Header.Get("If-Match") != `"4"` {
		t.Fatalf("%+v", r)
	}
	r = clitest.Run(t, h, "teams", "update-team-mcp-server", "--team-id", team, "--mcp-server-id", srv, "--purpose", "Lookups")
	if r.Code != 0 || len(r.Requests) != 2 || r.Requests[1].Header.Get("If-Match") != `"6"` {
		t.Fatalf("team: %+v", r)
	}
	// a team's switch has no version of its own: no read, no If-Match
	r = clitest.Run(t, reply(200, `{"server":{"id":"`+srv+`"},"enabled":true,"mode":"full","state":"on","managed":"enterprise"}`),
		"teams", "put-team-mcp-server", "--team-id", team, "--mcp-server-id", srv, "--enabled", "--mode", "full")
	if r.Code != 0 || len(r.Requests) != 1 || r.Req.Header.Get("If-Match") != "" || body(t, r.ReqBody)["enabled"] != true {
		t.Fatalf("switch: %+v %s", r, r.ReqBody)
	}
}

// FR-009: policy rules are JSON; review decisions carry a note and per-primitive decisions.
func TestMcpPolicyAndReviewBodies(t *testing.T) {
	r := clitest.Run(t, reply(200, `{"serverId":"`+srv+`","level":"service","teamId":null,"rules":[],"version":2}`), "mcp-servers", "put-mcp-policy",
		"--mcp-server-id", srv, "--rules", `[{"class":"destructive","effect":"block"}]`, "--trust-annotations")
	b := body(t, r.ReqBody)
	rules, _ := b["rules"].([]any)
	if r.Code != 0 || len(rules) != 1 || b["trustAnnotations"] != true || r.Req.Header.Get("If-Match") != "" {
		t.Fatalf("%+v %v", r, b)
	}
	r = clitest.Run(t, reply(200, `{"id":"`+srv+`","version":3}`), "mcp-servers", "approve-mcp-server", "--mcp-server-id", srv,
		"--note", "ok", "--primitives", `[{"primitiveId":"mpr_1","decision":"block"}]`)
	b = body(t, r.ReqBody)
	if r.Code != 0 || r.Req.URL.Path != "/v1/mcp-servers/"+srv+":approve" || b["note"] != "ok" || b["primitives"] == nil {
		t.Fatalf("%+v %v", r, b)
	}
}

// FR-010: approval refused by policy is a service error with the AWS-style message and a hint.
func TestApproveRefusedByPolicy(t *testing.T) {
	r := clitest.Run(t, reply(403, `{"code":"AccessDenied","message":"approval needs a security operator","reason":"approval_not_permitted_by_policy","requestId":"req_1"}`),
		"teams", "approve-team-mcp-server", "--team-id", team, "--mcp-server-id", srv, "--note", "ok")
	if r.Code != 254 || !strings.Contains(r.Err, "(AccessDenied) when calling the ApproveTeamMcpServer operation") {
		t.Fatalf("%+v", r)
	}
	if !strings.Contains(r.Err, "security operator") {
		t.Errorf("no hint: %q", r.Err)
	}
}

// FR-011: the effective MCP of an agent renders its servers; the report renders its rows.
func TestEffectiveAgentMcpAndReport(t *testing.T) {
	r := clitest.Run(t, reply(200, `{"subjectKind":"agent","subjectId":"`+agent+`","access":"all","mode":"search","env":["AFS_MCP_URL"],"servers":[{"serverId":"`+srv+`","slug":"github","scope":"enterprise","mode":"full","url":"https://afs-gateway.afs-t-a.svc:8443/mcp/github","primitives":3}]}`),
		"agents", "get-effective-agent-mcp", "--team-id", team, "--agent-id", agent, "--output", "text")
	if r.Code != 0 || !strings.Contains(r.Out, "github") || r.Req.URL.Path != "/v1/teams/"+team+"/agents/"+agent+"/mcp" {
		t.Fatalf("%+v", r)
	}
	r = clitest.Run(t, reply(200, `{"groupBy":"server","view":"unused","from":"2026-08-27T00:00:00Z","to":"2026-11-25T00:00:00Z","windowDays":90,"rows":[{"key":"`+srv+`","label":"github","calls":0,"errors":0,"refusals":0}]}`),
		"reports", "get-mcp-usage-report", "--from", "2026-08-27T00:00:00Z", "--to", "2026-11-25T00:00:00Z", "--group-by", "server", "--view", "unused", "--window", "90", "--output", "table")
	q := r.Req.URL.Query()
	if r.Code != 0 || q.Get("view") != "unused" || q.Get("window") != "90" || !strings.Contains(r.Out, "github") {
		t.Fatalf("%+v", r)
	}
}
