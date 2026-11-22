package teams_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/clitest"
)

const restorePlan = `{"dryRun":true,"revision":41,"changes":[{"action":"update","resourceKind":"agents","resourceId":"agt_1","fields":[{"field":"modelId","from":"a","to":"b"}]}],"skipped":[]}`

func TestCommandsDoneRestore(t *testing.T) {
	clitest.CheckOperations(t, "restore", "ListAllowedAccounts", "RestoreRevision")
}

// CLI-US-20.1: --dry-run sends dryRun: true and prints the plan.
func TestRestoreRevisionDryRun(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) { clitest.Reply(w, 200, restorePlan) },
		"teams", "restore-revision", "--team-id", "team_1", "--revision-id", "41", "--dry-run")
	if r.Code != 0 || string(r.ReqBody) != `{"dryRun":true}` || r.Req.URL.Path != "/v1/teams/team_1/revisions/41:restore" || !strings.Contains(r.Out, `"dryRun": true`) {
		t.Fatalf("%+v body=%s", r, r.ReqBody)
	}
}

// CLI-US-20.1: without --dry-run the restore is applied (no dryRun true sent).
func TestRestoreRevisionApply(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) {
		clitest.Reply(w, 200, strings.Replace(restorePlan, `"dryRun":true`, `"dryRun":false`, 1))
	}, "teams", "restore-revision", "--team-id", "team_1", "--revision-id", "41")
	if r.Code != 0 || strings.Contains(string(r.ReqBody), `"dryRun":true`) || !strings.Contains(r.Out, `"dryRun": false`) {
		t.Fatalf("%+v body=%s", r, r.ReqBody)
	}
}

// CLI-US-20.2: a viewer is refused by the API (254).
func TestRestoreRevisionViewerDenied(t *testing.T) {
	r := clitest.Run(t, clitest.Denied("revisions:restore"), "teams", "restore-revision", "--team-id", "team_1", "--revision-id", "41")
	if r.Code != 254 || !strings.Contains(r.Err, "(AccessDenied) when calling the RestoreRevision operation") || !strings.Contains(r.Err, "revisions:restore") {
		t.Fatalf("%+v", r)
	}
}
