package teams_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/clitest"
)

func TestCommandsDoneLifecycle(t *testing.T) {
	clitest.CheckOperations(t, "lifecycle", "UpdateTeam", "DeleteTeam", "SuspendTeam", "ResumeTeam")
}

// CLI-US-23.1: after suspend, a write is refused with Conflict (254); resume lifts the refusal.
func TestSuspendThenUpdateAgentConflict254(t *testing.T) {
	suspended := false
	h := func(w http.ResponseWriter, req *http.Request) {
		switch {
		case strings.HasSuffix(req.URL.Path, ":suspend"):
			suspended = true
			clitest.Reply(w, 200, `{"id":"team_1","state":"suspended"}`)
		case strings.HasSuffix(req.URL.Path, ":resume"):
			suspended = false
			clitest.Reply(w, 200, `{"id":"team_1","state":"active"}`)
		case suspended:
			clitest.Reply(w, 409, `{"code":"Conflict","message":"team is suspended","requestId":"r","details":{"reason":"team_suspended"}}`)
		default:
			clitest.Reply(w, 200, `{"id":"agt_1","version":4}`)
		}
	}
	if r := clitest.Run(t, h, "teams", "suspend-team", "--team-id", "team_1"); r.Code != 0 {
		t.Fatalf("%+v", r)
	}
	r := clitest.Run(t, h, "agents", "update-agent", "--team-id", "team_1", "--agent-id", "agt_1", "--description", "x", "--if-match", "4")
	if r.Code != 254 || !strings.Contains(r.Err, "(Conflict) when calling the UpdateAgent operation") {
		t.Fatalf("%+v", r)
	}
	if r := clitest.Run(t, h, "teams", "resume-team", "--team-id", "team_1"); r.Code != 0 {
		t.Fatalf("%+v", r)
	}
	if r := clitest.Run(t, h, "agents", "update-agent", "--team-id", "team_1", "--agent-id", "agt_1", "--description", "x", "--if-match", "4"); r.Code != 0 {
		t.Fatalf("%+v", r)
	}
}

// CLI-US-23.2: a team-admin's delete-team is the API's AccessDenied, 254.
func TestDeleteTeamAccessDenied(t *testing.T) {
	r := clitest.Run(t, clitest.Denied("teams:delete"), "teams", "delete-team", "--team-id", "team_1", "--if-match", "3")
	if r.Code != 254 || !strings.Contains(r.Err, "(AccessDenied) when calling the DeleteTeam operation") || !strings.Contains(r.Err, "teams:delete") {
		t.Fatalf("%+v", r)
	}
}

// CLI-12: delete-team --wait polls until the team answers 404.
func TestDeleteTeamWaitFor404(t *testing.T) {
	gets := 0
	r := clitest.Run(t, func(w http.ResponseWriter, req *http.Request) {
		if req.Method == "DELETE" {
			clitest.Reply(w, 202, `{"id":"team_1","state":"deleting"}`)
			return
		}
		gets++
		if gets < 3 {
			clitest.Reply(w, 200, `{"id":"team_1"}`)
			return
		}
		clitest.Fail(404, "NotFound", "gone")(w, req)
	}, "teams", "delete-team", "--team-id", "team_1", "--if-match", "3", "--wait")
	if r.Code != 0 || gets != 3 {
		t.Fatalf("%+v gets=%d", r, gets)
	}
}

// CLI-10: update-team without --if-match reads the version and sends it; with --if-match 7 sends 7.
func TestUpdateTeamIfMatch(t *testing.T) {
	var got []string
	h := func(w http.ResponseWriter, req *http.Request) {
		if req.Method == "GET" {
			clitest.Reply(w, 200, `{"id":"team_1","version":5}`)
			return
		}
		got = append(got, req.Header.Get("If-Match"))
		clitest.Reply(w, 200, `{"id":"team_1","version":6}`)
	}
	if r := clitest.Run(t, h, "teams", "update-team", "--team-id", "team_1", "--description", "d"); r.Code != 0 {
		t.Fatalf("%+v", r)
	}
	if r := clitest.Run(t, h, "teams", "update-team", "--team-id", "team_1", "--description", "d", "--if-match", "7"); r.Code != 0 {
		t.Fatalf("%+v", r)
	}
	if strings.Join(got, ",") != `"5","7"` {
		t.Fatalf("If-Match %v", got)
	}
}

// The writeOnly ceilingOverrideReason is an ordinary option, not a secret value.
func TestCeilingOverrideReasonIsOrdinaryOption(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) {
		clitest.Reply(w, 200, `{"id":"team_1","version":2}`)
	}, "teams", "update-team", "--team-id", "team_1", "--ceiling-override-usd", "9000.00", "--ceiling-override-reason", "launch", "--if-match", "1")
	if r.Code != 0 || !strings.Contains(string(r.ReqBody), `"ceilingOverrideReason":"launch"`) {
		t.Fatalf("%+v body=%s", r, r.ReqBody)
	}
}
