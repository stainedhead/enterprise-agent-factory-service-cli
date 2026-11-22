package teams_test

import (
	"strings"
	"testing"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/clitest"
)

func TestCommandsDoneExternalSystems(t *testing.T) {
	clitest.CheckOperations(t, "external-systems", "ListTeamExternalSystems", "DeclareTeamExternalSystem", "UpdateTeamExternalSystem", "SetTeamExternalSystemState", "DeleteTeamExternalSystem")
}

// CLI-1: declaring a system when mayAddExternalSystems is off is the API's AccessDenied (254).
func TestDeclareExternalSystemRefused(t *testing.T) {
	r := clitest.Run(t, clitest.Denied("teams:declare-external-system"),
		"teams", "declare-external-system", "--team-id", "team_1", "--slug", "ledger", "--name", "Ledger", "--category", "other",
		"--owner", "o@example.com", "--purpose", "p", "--endpoints", "ledger.example.com:443")
	if r.Code != 254 || !strings.Contains(r.Err, "(AccessDenied) when calling the DeclareTeamExternalSystem operation") {
		t.Fatalf("%+v", r)
	}
	if !strings.Contains(string(r.ReqBody), `"endpoints":["ledger.example.com:443"]`) {
		t.Errorf("body %s", r.ReqBody)
	}
}
