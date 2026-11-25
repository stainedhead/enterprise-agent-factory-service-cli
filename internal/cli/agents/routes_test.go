package agents_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/clitest"
)

// FR-012 (v1.5.0): create-agent takes --routes and --default-route; a <route>/<model> default is sent unchanged.
func TestCreateAgentRoutes(t *testing.T) {
	ok := func(w http.ResponseWriter, _ *http.Request) { clitest.Reply(w, 201, `{"id":"agt_1","version":1}`) }
	base := []string{"agents", "create-agent", "--team-id", "team_1", "--name", "coder", "--template", "claude-code",
		"--image", "registry.example.com/a@sha256:" + strings.Repeat("a", 64)}
	r := clitest.Run(t, ok, append(append([]string{}, base...), "--models", "all", "--routes", "openai-prod", "--default-route", "openai-prod", "--default-model-id", "gpt-4.1")...)
	if r.Code != 0 || !strings.Contains(string(r.ReqBody), `"models":{"access":"all","defaultModelId":"gpt-4.1","routeAccess":"selected","routes":["openai-prod"],"defaultRoute":"openai-prod"}`) {
		t.Errorf("exit %d body %s err %q", r.Code, r.ReqBody, r.Err)
	}
	r = clitest.Run(t, ok, append(append([]string{}, base...), "--default-model-id", "azure-openai/gpt-5.1")...)
	if r.Code != 0 || !strings.Contains(string(r.ReqBody), `"models":{"defaultModelId":"azure-openai/gpt-5.1"}`) {
		t.Errorf("exit %d body %s", r.Code, r.ReqBody)
	}
}
