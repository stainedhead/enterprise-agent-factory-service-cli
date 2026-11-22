package collaboration_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/clitest"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/golden"
)

const agent = "agt_01JD2Q3R4S5T6V7W8X9Y0Z0001"

const link = `{"agentId":"agt_01JD2Q3R4S5T6V7W8X9Y0Z0001","agentName":"reviewer-bot","connect":true,"paperclipAgentId":"ag_8f2","linkState":"linked","lastSyncedAt":"2026-11-22T10:00:00Z","lastError":null,"version":2}`

// CLI-US-26.1
func TestConnectAgentRequest(t *testing.T) {
	r := clitest.Run(t, reply(link), "collaboration", "connect-agent", "--team-id", team, "--agent-id", agent, "--paperclip-agent-id", "ag_8f2")
	if r.Code != 0 || r.Req.Method != "POST" || r.Req.URL.Path != "/v1/teams/"+team+"/collaboration/links/"+agent+":connect" ||
		string(r.ReqBody) != `{"paperclipAgentId":"ag_8f2"}` || r.Req.Header.Get("Idempotency-Key") == "" {
		t.Fatalf("%+v body %s", r, r.ReqBody)
	}
	r = clitest.Run(t, reply(link), "collaboration", "connect-agent", "--team-id", team, "--agent-id", agent)
	if r.Code != 0 || strings.TrimSpace(string(r.ReqBody)) != `{}` {
		t.Errorf("body %q", r.ReqBody)
	}
}

func TestDisconnectAgentRequest(t *testing.T) {
	r := clitest.Run(t, reply(strings.Replace(link, `"connect":true`, `"connect":false`, 1)), "collaboration", "disconnect-agent", "--team-id", team, "--agent-id", agent)
	if r.Code != 0 || r.Req.Method != "POST" || !strings.HasSuffix(r.Req.URL.Path, "/collaboration/links/"+agent+":disconnect") ||
		len(r.ReqBody) != 0 || r.Req.Header.Get("Idempotency-Key") == "" {
		t.Fatalf("%+v body %q", r, r.ReqBody)
	}
}

func TestLinkColumnsGolden(t *testing.T) {
	for _, f := range []string{"table", "text"} {
		r := clitest.Run(t, reply(link), "collaboration", "connect-agent", "--team-id", team, "--agent-id", agent, "--output", f)
		golden.Check(t, "link."+f+".golden", r.Out)
	}
}

// CLI-US-26.2
func TestConnectAgentReplay(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Idempotent-Replayed", "true")
		clitest.Reply(w, 200, link)
	}, "collaboration", "connect-agent", "--team-id", team, "--agent-id", agent, "--idempotency-key", "k-1")
	if r.Code != 0 || r.Req.Header.Get("Idempotency-Key") != "k-1" || !strings.Contains(r.Err, "stored response") || !strings.Contains(r.Out, "reviewer-bot") {
		t.Errorf("%+v", r)
	}
}

func TestConnectAgentNoWaitNoIfMatch(t *testing.T) {
	for _, op := range []string{"connect-agent", "disconnect-agent"} {
		r := clitest.Run(t, reply(link), "collaboration", op, "--team-id", team, "--agent-id", agent, "--wait")
		if r.Code != 252 || len(r.Requests) != 0 {
			t.Errorf("%s --wait: %d", op, r.Code)
		}
		r = clitest.Run(t, reply(link), "collaboration", op, "--team-id", team, "--agent-id", agent)
		if len(r.Requests) != 1 || r.Req.Header.Get("If-Match") != "" {
			t.Errorf("%s: %d requests, If-Match %q", op, len(r.Requests), r.Req.Header.Get("If-Match"))
		}
	}
}

// CLI-US-26.3
func TestConnectAgentValidation254(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) {
		clitest.Reply(w, 400, `{"code":"ValidationError","message":"the collaboration agent cannot be linked","requestId":"req_c","details":{"fields":[{"field":"agentId","reason":"invalid"}]}}`)
	}, "collaboration", "connect-agent", "--team-id", team, "--agent-id", agent)
	if r.Code != 254 || !strings.Contains(r.Err, "field: agentId reason: invalid") {
		t.Errorf("%d %q", r.Code, r.Err)
	}
}

func TestConnectAgentManagedRegistrationHint(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) {
		clitest.Reply(w, 400, `{"code":"ValidationError","message":"registration is managed","requestId":"req_m","details":{"fields":[{"field":"paperclipAgentId","reason":"managed_registration"}]}}`)
	}, "collaboration", "connect-agent", "--team-id", team, "--agent-id", agent, "--paperclip-agent-id", "ag_1")
	if r.Code != 254 || !strings.Contains(r.Err, "field: paperclipAgentId reason: managed_registration") || !strings.Contains(r.Err, "Registration is managed; the reconciler records the id.") {
		t.Errorf("%d %q", r.Code, r.Err)
	}
}

// CLI-US-26.4
func TestListCollaborationLinksFilters(t *testing.T) {
	r := clitest.Run(t, reply(`{"items":[`+link+`]}`), "collaboration", "list-collaboration-links", "--team-id", team, "--link-state", "failed", "--no-connect")
	q := r.Req.URL.Query()
	if r.Code != 0 || q.Get("linkState") != "failed" || q.Get("connect") != "false" {
		t.Fatalf("%v %s", q, r.Err)
	}
	r = clitest.Run(t, reply(`{"items":[`+link+`]}`), "collaboration", "list-collaboration-links", "--team-id", team, "--connect")
	if r.Req.URL.Query().Get("connect") != "true" || r.Req.URL.Query().Has("linkState") {
		t.Errorf("%v", r.Req.URL.Query())
	}
	r = clitest.Run(t, reply(`{"items":[]}`), "collaboration", "list-collaboration-links", "--team-id", team)
	if r.Req.URL.Query().Has("connect") {
		t.Errorf("connect must be omitted: %v", r.Req.URL.Query())
	}
}

func TestListCollaborationLinksPagination(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Query().Get("nextToken") == "" {
			clitest.Reply(w, 200, `{"items":[`+link+`],"nextToken":"n2"}`)
			return
		}
		clitest.Reply(w, 200, `{"items":[`+strings.Replace(link, "reviewer-bot", "second-bot", 1)+`]}`)
	}, "collaboration", "list-collaboration-links", "--team-id", team, "--output", "text")
	if r.Code != 0 || len(r.Requests) != 2 || !strings.Contains(r.Out, "reviewer-bot") || !strings.Contains(r.Out, "second-bot") {
		t.Errorf("%d requests, out %q", len(r.Requests), r.Out)
	}
}

func TestListCollaborationLinksLinkStateEnum252(t *testing.T) {
	r := clitest.Run(t, reply(`{}`), "collaboration", "list-collaboration-links", "--team-id", team, "--link-state", "broken")
	if r.Code != 252 || len(r.Requests) != 0 {
		t.Errorf("%d", r.Code)
	}
}

func TestLinksTableGolden(t *testing.T) {
	body := `{"items":[` + link + `,` + strings.Replace(strings.Replace(link, `"lastError":null`, `"lastError":"agent key not yet written"`, 1), "linked", "failed", 1) + `]}`
	r := clitest.Run(t, reply(body), "collaboration", "list-collaboration-links", "--team-id", team, "--output", "table")
	golden.Check(t, "links.table.golden", r.Out)
}

// CLI-US-26.5
func TestConnectAgentAccessDenied254(t *testing.T) {
	r := clitest.Run(t, clitest.Denied("collaboration:manage"), "collaboration", "connect-agent", "--team-id", team, "--agent-id", agent)
	if r.Code != 254 || !strings.Contains(r.Err, "Required permission: collaboration:manage") {
		t.Errorf("%d %q", r.Code, r.Err)
	}
}

func TestDisconnectAgentAccessDenied254(t *testing.T) {
	r := clitest.Run(t, clitest.Denied("collaboration:manage"), "collaboration", "disconnect-agent", "--team-id", team, "--agent-id", agent)
	if r.Code != 254 || !strings.Contains(r.Err, "(AccessDenied) when calling the DisconnectAgentFromCollaboration operation") {
		t.Errorf("%d %q", r.Code, r.Err)
	}
}

func TestListCollaborationLinksAccessDenied254(t *testing.T) {
	r := clitest.Run(t, clitest.Denied("collaboration:read"), "collaboration", "list-collaboration-links", "--team-id", team)
	if r.Code != 254 || !strings.Contains(r.Err, "Required permission: collaboration:read") {
		t.Errorf("%d %q", r.Code, r.Err)
	}
}
