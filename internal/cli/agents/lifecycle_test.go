package agents_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/clitest"
)

func TestCommandsDoneLifecycle(t *testing.T) {
	clitest.CheckOperations(t, "lifecycle", "UpdateAgent", "DeleteAgent", "PauseAgent", "ResumeAgent", "RestartAgent")
}

var update = []string{"agents", "update-agent", "--team-id", "team_1", "--agent-id", "agt_1", "--description", "d"}

// CLI-US-6.1: --if-match 7 sends If-Match: "7" and reads nothing first.
func TestUpdateAgentIfMatch7(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) { clitest.Reply(w, 200, `{"id":"agt_1","version":8}`) },
		append(append([]string{}, update...), "--if-match", "7")...)
	if r.Code != 0 || len(r.Requests) != 1 || r.Req.Header.Get("If-Match") != `"7"` || r.Req.Method != "PATCH" {
		t.Fatalf("%+v", r)
	}
}

// CLI-US-6.1: at version 8 the 412 names the current version and the flag to re-run with, from
// the response's ETag or, when it has none, from one more read.
func TestUpdateAgent412Message(t *testing.T) {
	want := "The resource changed (version 8). Re-run with --if-match 8 to apply your change to the latest version."
	withETag := func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("ETag", `"8"`)
		clitest.Reply(w, 412, `{"code":"PreconditionFailed","message":"version mismatch","requestId":"r"}`)
	}
	r := clitest.Run(t, withETag, append(append([]string{}, update...), "--if-match", "7")...)
	if r.Code != 254 || !strings.Contains(r.Err, "(PreconditionFailed)") || !strings.Contains(r.Err, want) {
		t.Fatalf("etag: %+v", r)
	}
	noETag := func(w http.ResponseWriter, req *http.Request) {
		if req.Method == "GET" {
			clitest.Reply(w, 200, `{"id":"agt_1","version":8}`)
			return
		}
		clitest.Reply(w, 412, `{"code":"PreconditionFailed","message":"version mismatch","requestId":"r"}`)
	}
	r = clitest.Run(t, noETag, append(append([]string{}, update...), "--if-match", "7")...)
	if r.Code != 254 || !strings.Contains(r.Err, want) || len(r.Requests) != 2 {
		t.Fatalf("read: %+v", r)
	}
}

// CLI-US-6.2: without --if-match the CLI reads the version, sends it, and says so on stderr.
func TestUpdateAgentReadThenWrite(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, req *http.Request) {
		if req.Method == "GET" {
			clitest.Reply(w, 200, `{"id":"agt_1","version":5}`)
			return
		}
		clitest.Reply(w, 200, `{"id":"agt_1","version":6}`)
	}, update...)
	if r.Code != 0 || r.Req.Header.Get("If-Match") != `"5"` || !strings.Contains(r.Err, "version 5") {
		t.Fatalf("%+v", r)
	}
}

// CLI-US-6.3: delete-agent sends If-Match; a 428 is 254 and says how to fix it.
func TestDeleteAgentIfMatchAnd428(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, req *http.Request) {
		if req.Method == "DELETE" && req.Header.Get("If-Match") == "" {
			clitest.Reply(w, 428, `{"code":"PreconditionRequired","message":"If-Match is required","requestId":"r"}`)
			return
		}
		clitest.Reply(w, 202, `{"id":"agt_1","version":4}`)
	}, "agents", "delete-agent", "--team-id", "team_1", "--agent-id", "agt_1", "--if-match", "4")
	if r.Code != 0 || r.Req.Header.Get("If-Match") != `"4"` {
		t.Fatalf("%+v", r)
	}
	r = clitest.Run(t, clitest.Fail(428, "PreconditionRequired", "If-Match is required"), "agents", "delete-agent", "--team-id", "team_1", "--agent-id", "agt_1", "--if-match", "4")
	if r.Code != 254 || !strings.Contains(r.Err, "(PreconditionRequired)") || !strings.Contains(r.Err, "--if-match") {
		t.Fatalf("428: %+v", r)
	}
}

// CLI-US-23.1: a write to a suspended team's agent is Conflict, 254.
func TestUpdateAgentConflictAfterSuspend(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) {
		clitest.Reply(w, 409, `{"code":"Conflict","message":"team is suspended","requestId":"r","details":{"reason":"team_suspended"}}`)
	}, append(append([]string{}, update...), "--if-match", "1")...)
	if r.Code != 254 || !strings.Contains(r.Err, "(Conflict) when calling the UpdateAgent operation") {
		t.Fatalf("%+v", r)
	}
}
