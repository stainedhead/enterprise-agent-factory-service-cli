package admin_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/clitest"
)

func TestCommandsDoneAgentTemplates(t *testing.T) {
	clitest.CheckOperations(t, "agent-templates", "ListAdminAgentTemplates", "PutAgentTemplate", "DeleteAgentTemplate")
}

const digest = "sha256:4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f"

// CLI-1: --allowed-digests is a list, --default-resources is JSON.
func TestPutAgentTemplateListAndJSON(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, req *http.Request) {
		if req.Method == "GET" {
			clitest.Reply(w, 404, `{"code":"NotFound","message":"none","requestId":"r"}`)
			return
		}
		clitest.Reply(w, 200, `{"name":"claude-code","version":1}`)
	}, "admin", "put-agent-template", "--name", "claude-code", "--harness", "claude-code", "--current-digest", digest,
		"--allowed-digests", digest, digest, "--default-resources", `{"cpu":"1"}`)
	if r.Code != 0 || !strings.Contains(string(r.ReqBody), `"allowedDigests":["`+digest+`","`+digest+`"]`) || !strings.Contains(string(r.ReqBody), `"defaultResources":{"cpu":"1"}`) {
		t.Fatalf("%+v body=%s", r, r.ReqBody)
	}
}

// CLI-1: deleting a template that agents use is a 409 (254).
func TestDeleteAgentTemplateInUse409(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, req *http.Request) {
		if req.Method == "GET" {
			clitest.Reply(w, 200, `{"items":[{"name":"claude-code","version":3}]}`)
			return
		}
		clitest.Reply(w, 409, `{"code":"Conflict","message":"used by agents","requestId":"r","details":{"dependents":["agt_1"]}}`)
	}, "admin", "delete-agent-template", "--name", "claude-code")
	if r.Code != 254 || !strings.Contains(r.Err, "Dependents: agt_1") {
		t.Fatalf("%+v", r)
	}
}
