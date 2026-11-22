package teams_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/clitest"
)

func TestCommandsDoneModelsAndRevisions(t *testing.T) {
	clitest.CheckOperations(t, "models", "ListTeamModels", "PutTeamModels")
	clitest.CheckOperations(t, "revisions", "ListRevisions", "GetRevision")
}

// D-C8: put-team-models sends no If-Match and prints no version notice; --model-ids is a list.
func TestPutTeamModelsNoIfMatch(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) {
		clitest.Reply(w, 200, `{"items":[{"id":"m1"},{"id":"m2"}]}`)
	}, "teams", "put-team-models", "--team-id", "team_1", "--model-ids", "m1", "m2")
	if r.Code != 0 || len(r.Requests) != 1 || r.Req.Method != "PUT" {
		t.Fatalf("%+v", r)
	}
	if r.Req.Header.Get("If-Match") != "" || strings.Contains(r.Err, "version") {
		t.Errorf("If-Match %q stderr %q", r.Req.Header.Get("If-Match"), r.Err)
	}
	if string(r.ReqBody) != `{"modelIds":["m1","m2"]}` {
		t.Errorf("body %s", r.ReqBody)
	}
}

// CLI-1: describe-revision takes the revision number as --rev.
func TestDescribeRevisionPath(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) {
		clitest.Reply(w, 200, `{"revision":12,"operation":"update"}`)
	}, "teams", "describe-revision", "--team-id", "team_1", "--rev", "12")
	if r.Code != 0 || r.Req.URL.Path != "/v1/teams/team_1/revisions/12" {
		t.Fatalf("%+v", r)
	}
}
