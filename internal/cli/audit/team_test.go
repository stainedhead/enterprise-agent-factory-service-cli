package audit_test

import (
	"net/http"
	"testing"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/clitest"
)

func TestCommandsDone(t *testing.T) {
	clitest.CheckOperations(t, "team", "ListAuditEvents")
}

// CLI-US-4.4, CLI-6: three CSV pages become one header and every row, and the filters are sent.
func TestAuditCSV3Pages(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, req *http.Request) {
		if req.Header.Get("Accept") != "text/csv" {
			t.Errorf("accept %q", req.Header.Get("Accept"))
		}
		w.Header().Set("Content-Type", "text/csv")
		switch req.URL.Query().Get("nextToken") {
		case "":
			w.Header().Set("X-Afs-Next-Token", "p2")
			_, _ = w.Write([]byte("id,action\nevt_1,create\nevt_2,update\n"))
		case "p2":
			w.Header().Set("X-Afs-Next-Token", "p3")
			_, _ = w.Write([]byte("id,action\nevt_3,delete\n"))
		default:
			_, _ = w.Write([]byte("id,action\nevt_4,create\n"))
		}
	}, "audit", "list-audit-events", "--team-id", "team_1", "--from", "2026-11-01T00:00:00Z", "--to", "2026-11-30T00:00:00Z",
		"--actor-id", "usr_1", "--action", "create", "--outcome", "denied", "--resource-kind", "agent", "--resource-id", "agt_1", "--output", "csv")
	want := "id,action\nevt_1,create\nevt_2,update\nevt_3,delete\nevt_4,create\n"
	if r.Code != 0 || r.Out != want || len(r.Requests) != 3 {
		t.Fatalf("%+v", r)
	}
	q := r.Requests[0].URL.Query()
	for k, v := range map[string]string{"from": "2026-11-01T00:00:00Z", "to": "2026-11-30T00:00:00Z", "actorId": "usr_1", "action": "create", "outcome": "denied", "resourceKind": "agent", "resourceId": "agt_1"} {
		if q.Get(k) != v {
			t.Errorf("%s=%q", k, q.Get(k))
		}
	}
}
