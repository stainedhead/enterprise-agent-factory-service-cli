package audit_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/clitest"
)

func TestCommandsDoneService(t *testing.T) {
	clitest.CheckOperations(t, "service", "ListAdminAuditEvents")
}

// CLI-6: the service audit log streams CSV pages as one table, with the team filter sent.
func TestServiceAuditCSV(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/csv")
		if req.URL.Query().Get("nextToken") == "" {
			w.Header().Set("X-Afs-Next-Token", "p2")
			_, _ = w.Write([]byte("id,action\n1,create\n"))
			return
		}
		_, _ = w.Write([]byte("id,action\n2,delete\n"))
	}, "audit", "list-service-audit-events", "--team-id", "team_1", "--output", "csv")
	if r.Code != 0 || r.Out != "id,action\n1,create\n2,delete\n" || r.Req.URL.Path != "/v1/admin/audit-events" || r.Requests[0].URL.Query().Get("teamId") != "team_1" {
		t.Fatalf("%+v", r)
	}
}

// A team role lacks admin:audit; the API says so and the CLI exits 254.
func TestServiceAuditTeamRoleDenied(t *testing.T) {
	r := clitest.Run(t, clitest.Denied("admin:audit"), "audit", "list-service-audit-events")
	if r.Code != 254 || !strings.Contains(r.Err, "admin:audit") {
		t.Fatalf("%+v", r)
	}
}
