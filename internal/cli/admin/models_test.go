package admin_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/clitest"
)

func TestCommandsDoneModels(t *testing.T) {
	clitest.CheckOperations(t, "models", "ListAdminModels", "CreateModel", "UpdateModel")
}

// CLI-1, D28: --max-output-tokens is required; a missing one is a usage error before any call.
func TestCreateModelRequiresMaxOutputTokens252(t *testing.T) {
	r := clitest.Run(t, func(http.ResponseWriter, *http.Request) { t.Error("a call was made") },
		"admin", "create-model", "--id", "m1", "--provider", "anthropic", "--display-name", "M", "--context-window", "200000")
	if r.Code != 252 || !strings.Contains(r.Err, "--max-output-tokens") {
		t.Fatalf("%+v", r)
	}
}

func TestCreateModelBody(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) { clitest.Reply(w, 201, `{"id":"m1","version":1}`) },
		"admin", "create-model", "--id", "m1", "--provider", "anthropic", "--display-name", "M", "--context-window", "200000", "--max-output-tokens", "64000")
	if r.Code != 0 || !strings.Contains(string(r.ReqBody), `"maxOutputTokens":64000`) || !strings.Contains(string(r.ReqBody), `"contextWindow":200000`) {
		t.Fatalf("%+v body=%s", r, r.ReqBody)
	}
}

// CLI-10: update-model reads the version from the model list.
func TestUpdateModelReadsVersionFromList(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, req *http.Request) {
		if req.Method == "GET" {
			clitest.Reply(w, 200, `{"items":[{"id":"m1","version":3}]}`)
			return
		}
		clitest.Reply(w, 200, `{"id":"m1","version":4}`)
	}, "admin", "update-model", "--model-id", "m1", "--no-enabled")
	if r.Code != 0 || r.Req.Header.Get("If-Match") != `"3"` {
		t.Fatalf("%+v", r)
	}
}
