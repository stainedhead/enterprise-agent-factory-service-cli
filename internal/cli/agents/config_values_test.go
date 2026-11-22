package agents_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/clitest"
)

func TestCommandsDoneConfigAndExternal(t *testing.T) {
	clitest.CheckOperations(t, "config-values", "ListAgentConfigValues", "PutAgentConfigValue", "DeleteAgentConfigValue")
	clitest.CheckOperations(t, "external-systems", "ListAgentExternalSystems", "SetAgentExternalSystemState")
}

func putValue(items string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method == "GET" {
			clitest.Reply(w, 200, items)
			return
		}
		clitest.Reply(w, 200, `{"name":"LOG_LEVEL","source":"literal","value":"debug","version":4}`)
	}
}

var put = []string{"agents", "put-config-value", "--team-id", "team_1", "--agent-id", "agt_1", "--name", "LOG_LEVEL", "--source", "literal", "--value", "debug"}

// CLI-10: put-config-value reads the version from the list (matched by name).
func TestPutAgentConfigValueReadsVersionFromList(t *testing.T) {
	r := clitest.Run(t, putValue(`{"items":[{"name":"OTHER","version":9},{"name":"LOG_LEVEL","version":3}]}`), put...)
	if r.Code != 0 || r.Req.Header.Get("If-Match") != `"3"` || string(r.ReqBody) != `{"source":"literal","value":"debug"}` {
		t.Fatalf("%+v body=%s", r, r.ReqBody)
	}
}

// CLI-10: a value that does not exist yet is created without If-Match.
func TestPutAgentConfigValueCreateNoIfMatch(t *testing.T) {
	r := clitest.Run(t, putValue(`{"items":[]}`), put...)
	if r.Code != 0 || r.Req.Header.Get("If-Match") != "" || strings.Contains(r.Err, "version") {
		t.Fatalf("%+v", r)
	}
}

// CLI-1: --source is an enumeration.
func TestPutAgentConfigValueSourceEnum252(t *testing.T) {
	args := append([]string{}, put...)
	args[len(args)-3] = "elsewhere"
	if r := clitest.Run(t, nil, args...); r.Code != 252 {
		t.Fatalf("%+v", r)
	}
}
