package teams_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/clitest"
)

func collabTeam(w http.ResponseWriter, _ *http.Request) {
	clitest.Reply(w, 200, `{"id":"team_1","slug":"pay","mayExposeCollaboration":true,"version":4}`)
}

// CLI-US-31.1
func TestUpdateTeamMayExposeCollaborationBody(t *testing.T) {
	r := clitest.Run(t, collabTeam, "teams", "update-team", "--team-id", "team_1", "--may-expose-collaboration", "--if-match", "3")
	if r.Code != 0 || string(r.ReqBody) != `{"mayExposeCollaboration":true}` {
		t.Errorf("%d %s %s", r.Code, r.ReqBody, r.Err)
	}
	r = clitest.Run(t, collabTeam, "teams", "update-team", "--team-id", "team_1", "--no-may-expose-collaboration", "--if-match", "3")
	if string(r.ReqBody) != `{"mayExposeCollaboration":false}` {
		t.Errorf("%s", r.ReqBody)
	}
	r = clitest.Run(t, collabTeam, "teams", "update-team", "--team-id", "team_1", "--display-name", "Pay", "--if-match", "3")
	if strings.Contains(string(r.ReqBody), "mayExpose") {
		t.Errorf("omitted when neither is given: %s", r.ReqBody)
	}
}

func TestCreateTeamMayExposeCollaborationBody(t *testing.T) {
	r := clitest.Run(t, collabTeam, "teams", "create-team", "--slug", "pay", "--display-name", "Pay", "--kind", "persistent", "--max-budget-usd", "100.00",
		"--first-admin", `{"email":"ada@example.com"}`, "--may-expose-collaboration")
	if r.Code != 0 || !strings.Contains(string(r.ReqBody), `"mayExposeCollaboration":true`) {
		t.Errorf("%d %s %s", r.Code, r.ReqBody, r.Err)
	}
}

func TestUpdateTeamMayExposeCollaborationAccessDenied254(t *testing.T) {
	r := clitest.Run(t, clitest.Denied("teams:set-guardrails"), "teams", "update-team", "--team-id", "team_1", "--may-expose-collaboration", "--if-match", "3")
	if r.Code != 254 || !strings.Contains(r.Err, "Required permission: teams:set-guardrails") {
		t.Errorf("%d %q", r.Code, r.Err)
	}
}

func TestCreateTeamHostnameHeld409(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) {
		clitest.Reply(w, 409, `{"code":"Conflict","message":"slug held","requestId":"req_h","details":{"reason":"hostname_held"}}`)
	}, "teams", "create-team", "--slug", "pay", "--display-name", "Pay", "--kind", "persistent", "--max-budget-usd", "100.00", "--first-admin", `{"email":"ada@example.com"}`)
	if r.Code != 254 || !strings.Contains(r.Err, "That slug's collaboration hostname is held for 30 days after the deletion of a team that had one; choose another slug.") {
		t.Errorf("%d %q", r.Code, r.Err)
	}
}

func TestDescribeTeamMayExposeJSONUnchanged(t *testing.T) {
	r := clitest.Run(t, collabTeam, "teams", "describe-team", "--team-id", "team_1", "--output", "json")
	if !strings.Contains(r.Out, `"mayExposeCollaboration": true`) {
		t.Errorf("%q", r.Out)
	}
}
