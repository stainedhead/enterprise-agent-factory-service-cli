package agents_test

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/clitest"
)

func TestCommandsDoneA(t *testing.T) {
	clitest.CheckOperations(t, "agents", "ListAgents", "GetAgent", "GetAgentStatus", "CreateAgent")
}

// CLI-US-4.5: an operation without CSV refuses --output csv before any call (252).
func TestListAgentsCSVRefused252(t *testing.T) {
	r := clitest.Run(t, func(http.ResponseWriter, *http.Request) { t.Error("a call was made") },
		"agents", "list-agents", "--team-id", "team_1", "--output", "csv")
	if r.Code != 252 || len(r.Requests) != 0 {
		t.Fatalf("%+v", r)
	}
}

const image = "registry.example.com/agents/coder@sha256:4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f"

func createWait(w http.ResponseWriter, req *http.Request) {
	switch {
	case req.Method == "POST":
		clitest.Reply(w, 201, `{"id":"agt_1","generation":1}`)
	case strings.HasSuffix(req.URL.Path, "/status"):
		clitest.Reply(w, 200, `{"id":"agt_1","generation":1,"status":{"phase":"Ready","observedGeneration":1}}`)
	default:
		clitest.Reply(w, 200, `{"id":"agt_1","name":"coder","status":{"phase":"Ready"}}`)
	}
}

// CLI-US-12.1: --volume takes JSON and --wait exits 0 once the agent is Ready.
func TestCreateAgentVolumeJSONAndWait(t *testing.T) {
	vol := `{"sizeGiB":20,"storageClass":"gp3-standard","onDelete":"snapshot_then_delete"}`
	r := clitest.Run(t, createWait, "agents", "create-agent", "--team-id", "team_1", "--name", "coder", "--template", "claude-code",
		"--image", image, "--model-id", "m1", "--volume", vol, "--wait")
	if r.Code != 0 || !strings.Contains(string(r.Requests[0].URL.Path), "/v1/teams/team_1/agents") || !strings.Contains(string(r.RequestBody[0]), `"volume":`+vol) {
		t.Fatalf("%+v body=%s", r, r.RequestBody[0])
	}
	if !strings.Contains(r.Out, `"name": "coder"`) {
		t.Errorf("final agent not printed: %s", r.Out)
	}
}

// CLI-US-12.2: --replicas is not an option: 252, nothing sent.
func TestCreateAgentReplicasUnknown252(t *testing.T) {
	r := clitest.Run(t, nil, "agents", "create-agent", "--team-id", "team_1", "--name", "c", "--template", "t", "--image", image, "--model-id", "m", "--replicas", "3")
	if r.Code != 252 || !strings.Contains(r.Err, "--replicas") || len(r.Requests) != 0 {
		t.Fatalf("%+v", r)
	}
}

// CLI-US-12.3: an image by tag is the API's ValidationError (254).
func TestCreateAgentImageTag254(t *testing.T) {
	r := clitest.Run(t, clitest.Fail(422, "ValidationError", "image must be pinned by digest"),
		"agents", "create-agent", "--team-id", "team_1", "--name", "c", "--template", "t", "--image", "registry.example.com/coder:latest", "--model-id", "m")
	if r.Code != 254 || !strings.Contains(r.Err, "(ValidationError) when calling the CreateAgent operation: image must be pinned by digest") {
		t.Fatalf("%+v", r)
	}
}

// CLI-US-12.4: the same name in two teams is fine; later commands address the agent by id.
func TestCreateAgentSameNameTwoTeams(t *testing.T) {
	h := func(w http.ResponseWriter, req *http.Request) {
		if req.Method == "POST" {
			team := strings.Split(req.URL.Path, "/")[3]
			clitest.Reply(w, 201, fmt.Sprintf(`{"id":"agt_%s","name":"coder"}`, team))
			return
		}
		clitest.Reply(w, 200, `{"id":"agt_a","name":"coder"}`)
	}
	for _, team := range []string{"a", "b"} {
		r := clitest.Run(t, h, "agents", "create-agent", "--team-id", team, "--name", "coder", "--template", "t", "--image", image, "--model-id", "m")
		if r.Code != 0 || !strings.Contains(r.Out, "agt_"+team) {
			t.Fatalf("%s: %+v", team, r)
		}
	}
	r := clitest.Run(t, h, "agents", "describe-agent", "--team-id", "a", "--agent-id", "agt_a")
	if r.Code != 0 || r.Req.URL.Path != "/v1/teams/a/agents/agt_a" {
		t.Fatalf("%+v", r)
	}
}
