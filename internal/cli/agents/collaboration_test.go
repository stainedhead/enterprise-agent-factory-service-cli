package agents_test

import (
	"net/http"
	"os"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/clitest"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/golden"
)

var createCollab = []string{"agents", "create-agent", "--team-id", "team_1", "--name", "paperclip", "--template", "paperclip", "--image", image, "--default-model-id", "m1"}

func with(base []string, more ...string) []string {
	return append(append([]string{}, base...), more...)
}

// CLI-US-28.1
func TestCreateAgentRoleCollaborationBody(t *testing.T) {
	r := clitest.Run(t, createWait, with(createCollab, "--role", "collaboration")...)
	if r.Code != 0 || !strings.Contains(string(r.RequestBody[0]), `"role":"collaboration"`) {
		t.Fatalf("%d %s %s", r.Code, r.RequestBody[0], r.Err)
	}
	w := clitest.Run(t, createWait, with(createCollab, "--role", "worker")...)
	if !strings.Contains(string(w.RequestBody[0]), `"role":"worker"`) {
		t.Errorf("%s", w.RequestBody[0])
	}
}

// ORC-30: --role is hidden and its enum is open, so the API's refusal (role_reserved) is what the person sees.
func TestCreateAgentRoleSentUnchanged(t *testing.T) {
	r := clitest.Run(t, createWait, with(createCollab, "--role", "boss")...)
	if r.Code != 0 || len(r.Requests) != 1 || !strings.Contains(string(r.RequestBody[0]), `"role":"boss"`) {
		t.Errorf("%d %q", r.Code, r.Err)
	}
}

func TestCreateAgentCollaborationWaitReady(t *testing.T) {
	r := clitest.Run(t, createWait, with(createCollab, "--role", "collaboration", "--wait")...)
	if r.Code != 0 || !strings.Contains(r.Out, `"name": "coder"`) {
		t.Errorf("%d %q %q", r.Code, r.Out, r.Err)
	}
}

// CLI-US-28.2
func TestCreateAgentCollaborationConnectBody(t *testing.T) {
	cases := []struct {
		args []string
		want string
	}{
		{[]string{"--collaboration-connect"}, `"collaboration":{"connect":true}`},
		{[]string{"--no-collaboration-connect"}, `"collaboration":{"connect":false}`},
		{[]string{"--collaboration-connect", "--paperclip-agent-id", "ag_7"}, `"collaboration":{"connect":true,"paperclipAgentId":"ag_7"}`},
		{[]string{"--paperclip-agent-id", "ag_7"}, `"collaboration":{"paperclipAgentId":"ag_7"}`},
	}
	for _, c := range cases {
		r := clitest.Run(t, createWait, with(createCollab, c.args...)...)
		if r.Code != 0 || !strings.Contains(string(r.RequestBody[0]), c.want) {
			t.Errorf("%v: %d %s %s", c.args, r.Code, r.RequestBody[0], r.Err)
		}
	}
}

func TestUpdateAgentCollaborationConnectBody(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) { clitest.Reply(w, 200, `{"id":"agt_1","version":3}`) },
		"agents", "update-agent", "--team-id", "team_1", "--agent-id", "agt_1", "--collaboration-connect", "--paperclip-agent-id", "ag_9", "--if-match", "2")
	if r.Code != 0 || string(r.ReqBody) != `{"collaboration":{"connect":true,"paperclipAgentId":"ag_9"}}` {
		t.Errorf("%d %s %s", r.Code, r.ReqBody, r.Err)
	}
}

func TestCreateAgentNoCollaborationMemberByDefault(t *testing.T) {
	r := clitest.Run(t, createWait, createCollab...)
	if strings.Contains(string(r.RequestBody[0]), "collaboration") || strings.Contains(string(r.RequestBody[0]), `"role"`) {
		t.Errorf("%s", r.RequestBody[0])
	}
}

// CLI-US-28.3: contract 1.4 documents role on UpdateAgent (immutable); the flag is sent and the API judges it.
func TestUpdateAgentRoleSentAndJudgedByAPI(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) {
		clitest.Reply(w, 400, `{"code":"ValidationError","message":"role cannot change","requestId":"req_r","details":{"fields":[{"field":"role","reason":"immutable"}]}}`)
	}, "agents", "update-agent", "--team-id", "team_1", "--agent-id", "agt_1", "--role", "collaboration", "--if-match", "1")
	if r.Code != 254 || !strings.Contains(string(r.ReqBody), `"role":"collaboration"`) {
		t.Errorf("%d %q %s", r.Code, r.Err, r.ReqBody)
	}
}

func TestUpdateAgentRoleChangeRefused254(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) {
		clitest.Reply(w, 400, `{"code":"ValidationError","message":"role cannot change","requestId":"req_r","details":{"fields":[{"field":"role","reason":"immutable"}]}}`)
	}, "agents", "update-agent", "--team-id", "team_1", "--agent-id", "agt_1", "--cli-input-json", writeInput(t, `{"role":"collaboration","version":1}`))
	if r.Code != 254 || !strings.Contains(r.Err, "(ValidationError) when calling the UpdateAgent operation") {
		t.Errorf("%d %q", r.Code, r.Err)
	}
}

// CLI-US-28.4
func TestListAgentsRoleFilter(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) { clitest.Reply(w, 200, `{"items":[]}`) },
		"agents", "list-agents", "--team-id", "team_1", "--role", "collaboration")
	if r.Code != 0 || r.Req.URL.Query().Get("role") != "collaboration" {
		t.Errorf("%d %v", r.Code, r.Req.URL.Query())
	}
}

// The role filter is the contract's closed enum (worker, orchestration, collaboration): a value outside it is a usage error.
func TestListAgentsRoleClosedEnum(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) { clitest.Reply(w, 200, `{"items":[]}`) },
		"agents", "list-agents", "--team-id", "team_1", "--role", "boss")
	if r.Code != 252 || !strings.Contains(r.Err, "invalid choice") {
		t.Errorf("%d %q", r.Code, r.Err)
	}
}

func TestAgentTableRoleAndLinkColumnsGolden(t *testing.T) {
	body := `{"items":[` +
		`{"id":"agt_1","name":"orchestrator","template":"orchestration","models":{"access":"all","defaultModelId":"m1"},"paused":false,"role":"orchestration","status":{"phase":"Ready"},"version":2},` +
		`{"id":"agt_2","name":"coder","template":"claude-code","models":{"access":"selected","defaultModelId":"m1"},"paused":false,"role":"worker","collaboration":{"connect":true,"linkState":"linked"},"status":{"phase":"Ready"},"version":5},` +
		`{"id":"agt_3","name":"fixer","template":"claude-code","models":{"access":"selected","defaultModelId":"m1"},"paused":false,"role":"worker","status":{"phase":"Ready"},"version":1}]}`
	for _, f := range []string{"table", "text"} {
		r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) { clitest.Reply(w, 200, body) },
			"agents", "list-agents", "--team-id", "team_1", "--output", f)
		golden.Check(t, "agent-role-link."+f+".golden", r.Out)
		for _, line := range strings.Split(r.Out, "\n") {
			if f == "table" && len([]rune(line)) > 80 {
				t.Errorf("line over 80 columns: %q", line)
			}
		}
	}
}

// CLI-US-28.5
func TestDescribeAgentCollaborationJSONUnchanged(t *testing.T) {
	body := `{"id":"agt_2","role":"worker","collaboration":{"connect":true,"paperclipAgentId":"ag_1","linkState":"linked","lastSyncedAt":"2026-11-22T10:00:00Z","lastError":null}}`
	r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) { clitest.Reply(w, 200, body) },
		"agents", "describe-agent", "--team-id", "team_1", "--agent-id", "agt_2", "--output", "json")
	got := strings.Join(strings.Fields(r.Out), "")
	if !strings.Contains(got, `"role":"worker"`) || !strings.Contains(got, `"collaboration":{"connect":true,"paperclipAgentId":"ag_1","linkState":"linked","lastSyncedAt":"2026-11-22T10:00:00Z","lastError":null}`) {
		t.Errorf("%s", r.Out)
	}
}

// CLI-US-28.6
func TestCreateAgentCollaborationErrors254(t *testing.T) {
	cases := []struct {
		name, status, code, details, hint string
	}{
		{"exists", "409", "Conflict", `{"reason":"collaboration_agent_exists"}`, "This team already has an orchestration agent; see afs orchestration describe-orchestration."},
		{"template", "400", "ValidationError", `{"fields":[{"field":"template","reason":"not_paperclip_capable"}]}`, "field: template reason: not_paperclip_capable"},
		{"identity", "400", "ValidationError", `{"fields":[{"field":"iamIdentity","reason":"team_shared"}]}`, "field: iamIdentity reason: team_shared"},
		{"snapshot", "400", "ValidationError", `{"fields":[{"field":"volume.snapshotSchedule","reason":"off"}]}`, "field: volume.snapshotSchedule reason: off"},
		{"worker", "400", "ValidationError", `{"fields":[{"field":"template","reason":"paperclip_on_worker"}]}`, "field: template reason: paperclip_on_worker"},
		{"connect", "400", "ValidationError", `{"fields":[{"field":"collaboration.connect","reason":"no_collaboration_agent"}]}`, "field: collaboration.connect reason: no_collaboration_agent"},
	}
	for _, c := range cases {
		status := 409
		if c.status == "400" {
			status = 400
		}
		r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) {
			clitest.Reply(w, status, `{"code":"`+c.code+`","message":"refused","requestId":"req_x","details":`+c.details+`}`)
		}, createCollab...)
		if r.Code != 254 || !strings.Contains(r.Err, "("+c.code+") when calling the CreateAgent operation") || !strings.Contains(r.Err, c.hint) {
			t.Errorf("%s: %d %q", c.name, r.Code, r.Err)
		}
	}
}

func TestCreateAgentCollaboration503Retried(t *testing.T) {
	var n atomic.Int32
	r := clitest.Run(t, func(w http.ResponseWriter, req *http.Request) {
		if n.Add(1) == 1 {
			clitest.Reply(w, 503, `{"code":"Unavailable","message":"secrets manager unavailable","requestId":"req_u"}`)
			return
		}
		clitest.Reply(w, 201, `{"id":"agt_1"}`)
	}, with(createCollab, "--role", "collaboration")...)
	if r.Code != 0 || len(r.Requests) != 2 || r.Requests[0].Header.Get("Idempotency-Key") == "" ||
		r.Requests[0].Header.Get("Idempotency-Key") != r.Requests[1].Header.Get("Idempotency-Key") {
		t.Errorf("%d %d requests", r.Code, len(r.Requests))
	}
}

// CLI-US-28.7
func TestListAgentTemplatesCapabilityFilter(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) {
		clitest.Reply(w, 200, `{"items":[{"name":"paperclip","harness":"paperclip","currentDigest":"sha256:aa","capabilities":["paperclip"],"enabled":true,"version":1}]}`)
	}, "agents", "list-agent-templates", "--capability", "paperclip", "--output", "table")
	if r.Code != 0 || r.Req.URL.Query().Get("capability") != "paperclip" || !strings.Contains(r.Out, "Capabilities") || !strings.Contains(r.Out, "paperclip") {
		t.Errorf("%d %v %q", r.Code, r.Req.URL.Query(), r.Out)
	}
}

// CLI-US-28.8
func TestCreateAgentCollaborationAccessDenied254(t *testing.T) {
	r := clitest.Run(t, clitest.Denied("agents:write"), with(createCollab, "--role", "collaboration")...)
	if r.Code != 254 || !strings.Contains(r.Err, "Required permission: agents:write") {
		t.Errorf("%d %q", r.Code, r.Err)
	}
}

// CLI-US-28.9
func TestCreateAgentSkeletonHasRoleAndCollaboration(t *testing.T) {
	r := clitest.Run(t, nil, "agents", "create-agent", "--generate-cli-skeleton")
	if r.Code != 0 || len(r.Requests) != 0 || !strings.Contains(r.Out, `"role"`) || !strings.Contains(r.Out, `"collaboration"`) {
		t.Errorf("%d %q", r.Code, r.Out)
	}
}

func writeInput(t *testing.T, content string) string {
	t.Helper()
	f := t.TempDir() + "/in.json"
	if err := os.WriteFile(f, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	return "file://" + f
}
