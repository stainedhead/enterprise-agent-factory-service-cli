package admin_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/clitest"
)

func TestCommandsDoneCostTagStatus(t *testing.T) {
	clitest.CheckOperations(t, "cost-tag-status", "ListCostTagStatus")
}

// TAG-7: one row per tag key with AWS's activation status and whether the cost report carries it.
func TestListCostTagStatusTable(t *testing.T) {
	body := `{"items":[{"key":"afs:team","awsStatus":"active","seenInCur":true,"firstSeenAt":"2026-10-02T00:00:00Z","lastSeenAt":"2026-11-20T00:00:00Z","checkedAt":"2026-11-21T18:00:00Z"},{"key":"cost-center","awsStatus":"inactive","seenInCur":false,"firstSeenAt":null,"lastSeenAt":null,"checkedAt":"2026-11-21T18:00:00Z"}]}`
	r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) { clitest.Reply(w, 200, body) }, "admin", "list-cost-tag-status", "--output", "table")
	if r.Code != 0 || !strings.Contains(r.Out, "inactive") || len(strings.Split(strings.TrimSpace(r.Out), "\n")) != 3 {
		t.Fatalf("%+v", r)
	}
}
