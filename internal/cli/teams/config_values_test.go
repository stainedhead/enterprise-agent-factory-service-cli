package teams_test

import (
	"net/http"
	"testing"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/clitest"
)

func TestCommandsDoneConfigValues(t *testing.T) {
	clitest.CheckOperations(t, "config-values", "ListNamespaceConfigValues", "PutNamespaceConfigValue", "DeleteNamespaceConfigValue")
}

// CLI-10: put-config-value reads ListNamespaceConfigValues; a missing value is created without If-Match.
func TestPutNamespaceConfigValueCreateNoIfMatch(t *testing.T) {
	h := func(items string) http.HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) {
			if req.Method == "GET" {
				clitest.Reply(w, 200, items)
				return
			}
			clitest.Reply(w, 200, `{"name":"REGION","value":"eu-west-1","version":2}`)
		}
	}
	args := []string{"teams", "put-config-value", "--team-id", "team_1", "--name", "REGION", "--value", "eu-west-1"}
	r := clitest.Run(t, h(`{"items":[]}`), args...)
	if r.Code != 0 || r.Req.Header.Get("If-Match") != "" || len(r.Requests) != 2 {
		t.Fatalf("create: %+v", r)
	}
	r = clitest.Run(t, h(`{"items":[{"name":"REGION","version":6}]}`), args...)
	if r.Code != 0 || r.Req.Header.Get("If-Match") != `"6"` {
		t.Fatalf("update: %+v", r)
	}
}
