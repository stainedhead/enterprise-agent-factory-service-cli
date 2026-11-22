package members_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/clitest"
)

func TestCommandsDone(t *testing.T) {
	clitest.CheckOperations(t, "members", "ListMembers", "AddMember", "UpdateMember", "RemoveMember")
}

// CLI-10: remove-member and update-member read the version from ListMembers (no single GET).
func TestRemoveMemberReadsVersionFromList(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, req *http.Request) {
		if req.Method == "GET" {
			clitest.Reply(w, 200, `{"items":[{"id":"mem_a","version":2},{"id":"mem_b","version":9}]}`)
			return
		}
		w.WriteHeader(204)
	}, "members", "remove-member", "--team-id", "team_1", "--binding-id", "mem_b")
	if r.Code != 0 || r.Req.Method != "DELETE" || r.Req.Header.Get("If-Match") != `"9"` || r.Out != "" {
		t.Fatalf("%+v if-match=%q", r, r.Req.Header.Get("If-Match"))
	}
	if !strings.HasSuffix(r.Req.URL.Path, "/members/mem_b") || !strings.HasSuffix(r.Requests[0].URL.Path, "/members") {
		t.Errorf("paths %s %s", r.Requests[0].URL.Path, r.Req.URL.Path)
	}
}

func TestUpdateMemberReadsVersionFromList(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, req *http.Request) {
		if req.Method == "GET" {
			clitest.Reply(w, 200, `{"items":[{"id":"mem_b","version":4}]}`)
			return
		}
		clitest.Reply(w, 200, `{"id":"mem_b","role":"team_viewer","version":5}`)
	}, "members", "update-member", "--team-id", "team_1", "--binding-id", "mem_b", "--role", "team_viewer")
	if r.Code != 0 || r.Req.Header.Get("If-Match") != `"4"` || string(r.ReqBody) != `{"role":"team_viewer"}` {
		t.Fatalf("%+v body=%s", r, r.ReqBody)
	}
}
