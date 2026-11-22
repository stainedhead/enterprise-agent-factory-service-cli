package budgets_test

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/clitest"
)

func TestCommandsDoneUsage(t *testing.T) {
	clitest.CheckOperations(t, "usage", "GetUsage", "ListCosts")
}

// CLI-6, CLI-13: get-usage --output csv sends Accept: text/csv with the window and grouping, and
// writes the CSV byte for byte.
func TestGetUsageCSV(t *testing.T) {
	csv := "agent,model,calls,costUsd\r\nagt_1,m1,12,2.50\r\n"
	r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/csv")
		_, _ = w.Write([]byte(csv))
	}, "budgets", "get-usage", "--team-id", "team_1", "--from", "2026-11-01T00:00:00Z", "--to", "2026-11-21T00:00:00Z", "--group-by", "agent", "--output", "csv")
	if r.Code != 0 || r.Out != csv || r.Req.Header.Get("Accept") != "text/csv" {
		t.Fatalf("%+v", r)
	}
	q := r.Req.URL.Query()
	if q.Get("from") != "2026-11-01T00:00:00Z" || q.Get("to") != "2026-11-21T00:00:00Z" || q.Get("groupBy") != "agent" {
		t.Errorf("query %v", q)
	}
}

// CLI-13: --group-by is one of agent, model, day.
func TestGetUsageGroupByEnum252(t *testing.T) {
	if r := clitest.Run(t, nil, "budgets", "get-usage", "--team-id", "team_1", "--group-by", "planet"); r.Code != 252 {
		t.Fatalf("%+v", r)
	}
}

// Usage rows keep their group key: the table spreads it into key.<member> columns.
func TestGetUsageTableShowsKey(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) {
		clitest.Reply(w, 200, `{"groupBy":["agent"],"rows":[{"calls":3,"costUsd":"1.00","key":{"agent":"agt_1"}}],"totals":{"calls":3}}`)
	}, "budgets", "get-usage", "--team-id", "team_1", "--output", "table")
	if r.Code != 0 || !strings.Contains(r.Out, "key.agent") || !strings.Contains(r.Out, "agt_1") {
		t.Fatalf("%+v", r)
	}
}

// CLI-7: list-costs follows nextToken and merges the pages.
func TestListCostsPaginates(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, req *http.Request) {
		switch req.URL.Query().Get("nextToken") {
		case "":
			clitest.Reply(w, 200, `{"items":[{"day":"2026-11-19","costUsd":"1.00"},{"day":"2026-11-20","costUsd":"2.00"}],"nextToken":"p2"}`)
		default:
			clitest.Reply(w, 200, `{"items":[{"day":"2026-11-21","costUsd":"3.00"}]}`)
		}
	}, "budgets", "list-costs", "--team-id", "team_1")
	var doc struct{ Items []map[string]any }
	if err := json.Unmarshal([]byte(r.Out), &doc); err != nil || r.Code != 0 || len(doc.Items) != 3 || len(r.Requests) != 2 {
		t.Fatalf("%+v err=%v", r, err)
	}
}
