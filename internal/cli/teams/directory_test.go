package teams_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/clitest"
)

func TestCommandsDoneDirectory(t *testing.T) {
	clitest.CheckOperations(t, "directory", "ListTeamDirectory")
}

func TestListTeamDirectory(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) {
		clitest.Reply(w, 200, `{"items":[{"id":"team_01J9ZQ3ABCDEFGHJKMNPQRSTVW","slug":"payments","displayName":"Payments"}]}`)
	}, "teams", "list-team-directory", "--output", "table")
	if r.Code != 0 || r.Req.URL.Path != "/v1/team-directory" || !strings.Contains(r.Out, "team_01J9ZQ3ABCDEFGHJKMNPQRSTVW") || !strings.Contains(r.Out, "payments") {
		t.Fatalf("%+v", r)
	}
}

// D-56: a team-admin is refused (254).
func TestListTeamDirectoryTeamAdminDenied(t *testing.T) {
	r := clitest.Run(t, clitest.Denied("teams:read-directory"), "teams", "list-team-directory")
	if r.Code != 254 || !strings.Contains(r.Err, "teams:read-directory") {
		t.Fatalf("%+v", r)
	}
}
