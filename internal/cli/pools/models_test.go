package pools_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/clitest"
)

const poolImage = "registry.example.com/a@sha256:"

func okPool(w http.ResponseWriter, _ *http.Request) {
	clitest.Reply(w, 201, `{"id":"pool_1","version":2,"models":{"access":"all","defaultModelId":"m1"}}`)
}

func createPool(t *testing.T, extra ...string) clitest.Result {
	t.Helper()
	args := append([]string{"pools", "create-pool", "--team-id", "team_1", "--name", "fixer", "--template", "claude-code", "--image", poolImage + strings.Repeat("a", 64), "--workers", "10"}, extra...)
	return clitest.Run(t, okPool, args...)
}

// MOD-3, MOD-8, OR-CLI-28: the same model options as agents, on pools.
func TestPoolModelsBody(t *testing.T) {
	r := createPool(t, "--models", "m1", "m2", "--default-model-id", "m1")
	if r.Code != 0 || !strings.Contains(strings.ReplaceAll(string(r.ReqBody), " ", ""), `"models":{"access":"selected","modelIds":["m1","m2"],"defaultModelId":"m1"}`) || r.Err != "" {
		t.Errorf("exit %d body %s err %q", r.Code, r.ReqBody, r.Err)
	}
	u := clitest.Run(t, okPool, "pools", "update-pool", "--team-id", "team_1", "--pool-id", "pool_1", "--models", "all", "--if-match", "2")
	if u.Code != 0 || string(u.ReqBody) != `{"models":{"access":"all"}}` {
		t.Errorf("exit %d body %s", u.Code, u.ReqBody)
	}
	if e := createPool(t, "--models", "all", "m1"); e.Code != 252 || len(e.Requests) != 0 {
		t.Errorf("exit %d", e.Code)
	}
}

// OR-CLI-28: --model-id keeps working on pools with one deprecation line.
func TestPoolModelIDDeprecated(t *testing.T) {
	r := createPool(t, "--model-id", "m1")
	want := "Warning: --model-id is deprecated and will be removed in 2.0.0; use --default-model-id (and --models to choose the set) instead.\n"
	if r.Code != 0 || r.Err != want || !strings.Contains(string(r.ReqBody), `"modelId":"m1"`) {
		t.Errorf("exit %d err %q body %s", r.Code, r.Err, r.ReqBody)
	}
	u := clitest.Run(t, okPool, "pools", "update-pool", "--team-id", "team_1", "--pool-id", "pool_1", "--model-id", "m1", "--if-match", "2")
	if u.Err != want {
		t.Errorf("update: %q", u.Err)
	}
}

// ORC-30: the reserved name applies to pools; the model refusals print their hints.
func TestPoolRefusals(t *testing.T) {
	for _, tc := range []struct{ reason, field, want string }{
		{"name_reserved", "name", "reserved for the team's orchestration agent"},
		{"model_not_enabled_for_team", "models.modelIds[0]", "afs teams put-team-models"},
		{"models_selection_empty", "models", "--models all"},
	} {
		r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) {
			clitest.Reply(w, 400, `{"code":"ValidationError","message":"invalid","requestId":"req_p","details":{"fields":[{"field":"`+tc.field+`","reason":"`+tc.reason+`"}]}}`)
		}, "pools", "create-pool", "--team-id", "team_1", "--name", "orchestrator", "--template", "t", "--image", poolImage+strings.Repeat("a", 64), "--workers", "1")
		if r.Code != 254 || !strings.Contains(r.Err, "field: "+tc.field+" reason: "+tc.reason) || !strings.Contains(r.Err, tc.want) {
			t.Errorf("%s: exit %d %q", tc.reason, r.Code, r.Err)
		}
	}
}
