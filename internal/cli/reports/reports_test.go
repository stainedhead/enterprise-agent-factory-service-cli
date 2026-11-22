package reports_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/clitest"
)

func TestCommandsDone(t *testing.T) {
	clitest.CheckOperations(t, "reports", "GetTeamReport", "GetAgentReport", "GetServiceReport")
}

func ok(w http.ResponseWriter, _ *http.Request) {
	clitest.Reply(w, 200, `{"kind":"costs","groupBy":["agent"],"rows":[{"agent":"agt_1","calls":2,"costUsd":"1.50"}],"totals":{"calls":2}}`)
}

// CLI-13: every report flag reaches the query.
func TestReportFlagsToQuery(t *testing.T) {
	r := clitest.Run(t, ok, "reports", "get-team-report", "--team-id", "team_1", "--kind", "costs", "--range", "30d",
		"--group-by", "tag:project,agent", "--agent-id", "agt_1", "--model-id", "m1", "--system-id", "ext_1", "--pool-id", "pool_1")
	if r.Code != 0 || r.Req.URL.Path != "/v1/teams/team_1/reports/costs" {
		t.Fatalf("%+v", r)
	}
	q := r.Req.URL.Query()
	for k, v := range map[string]string{"range": "30d", "groupBy": "tag:project,agent", "agentId": "agt_1", "modelId": "m1", "systemId": "ext_1", "poolId": "pool_1"} {
		if q.Get(k) != v {
			t.Errorf("%s=%q want %q", k, q.Get(k), v)
		}
	}
	r = clitest.Run(t, ok, "reports", "get-service-report", "--kind", "usage", "--from", "2026-11-01T00:00:00Z", "--to", "2026-11-08T00:00:00Z", "--team-id", "team_2")
	q = r.Req.URL.Query()
	if r.Code != 0 || q.Get("from") != "2026-11-01T00:00:00Z" || q.Get("to") != "2026-11-08T00:00:00Z" || q.Get("teamId") != "team_2" {
		t.Fatalf("%+v q=%v", r, q)
	}
}

// CLI-13: --range is one of 24h, 7d, 30d, 90d, mtd; it cannot be combined with --from/--to; --group-by is checked.
func TestRangeEnum252(t *testing.T) {
	base := []string{"reports", "get-team-report", "--team-id", "team_1", "--kind", "costs"}
	for name, extra := range map[string][]string{
		"bad range":      {"--range", "1y"},
		"range and from": {"--range", "7d", "--from", "2026-11-01T00:00:00Z"},
		"bad group-by":   {"--group-by", "planet"},
		"trailing comma": {"--range", "7d", "--group-by", "agent,"},
	} {
		args := append([]string{}, base...)
		args = append(args, extra...)
		if r := clitest.Run(t, func(http.ResponseWriter, *http.Request) { t.Errorf("%s: a call was made", name) }, args...); r.Code != 252 {
			t.Errorf("%s: %+v", name, r)
		}
	}
	if r := clitest.Run(t, nil, "reports", "get-team-report", "--team-id", "team_1", "--kind", "mood"); r.Code != 252 {
		t.Errorf("kind: %+v", r)
	}
}

// CLI-13: --tag key=value becomes the filter tag.<key>=<value>; several pairs are allowed.
func TestTagFilterShorthand(t *testing.T) {
	r := clitest.Run(t, ok, "reports", "get-team-report", "--team-id", "team_1", "--kind", "costs", "--tag", "project=checkout", "env=prod")
	q := r.Req.URL.Query()
	if r.Code != 0 || q.Get("tag.project") != "checkout" || q.Get("tag.env") != "prod" {
		t.Fatalf("%+v q=%v", r, q)
	}
	for _, bad := range []string{"project", "=x"} {
		if r := clitest.Run(t, ok, "reports", "get-team-report", "--team-id", "team_1", "--kind", "costs", "--tag", bad); r.Code != 252 {
			t.Errorf("--tag %s: %+v", bad, r)
		}
	}
}

// CLI-US-4.2: --output csv sends Accept: text/csv and writes the API's CSV byte for byte.
func TestTeamReportCSVByteForByte(t *testing.T) {
	csv := "tag:project,costUsd\r\ncheckout,25.35\r\n"
	r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/csv")
		_, _ = w.Write([]byte(csv))
	}, "reports", "get-team-report", "--team-id", "team_1", "--kind", "costs", "--range", "30d", "--group-by", "tag:project", "--output", "csv")
	if r.Code != 0 || r.Out != csv || r.Req.Header.Get("Accept") != "text/csv" {
		t.Fatalf("%+v", r)
	}
}

// CLI-US-4.3: a CSV report over 3 pages follows X-Afs-Next-Token and writes one header.
func TestReportCSV3Pages(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/csv")
		switch req.URL.Query().Get("nextToken") {
		case "":
			w.Header().Set("X-Afs-Next-Token", "t2")
			_, _ = w.Write([]byte("agent,costUsd\r\nagt_1,1.00\r\n"))
		case "t2":
			w.Header().Set("X-Afs-Next-Token", "t3")
			_, _ = w.Write([]byte("agent,costUsd\r\nagt_2,2.00\r\n"))
		default:
			_, _ = w.Write([]byte("agent,costUsd\r\nagt_3,3.00\r\n"))
		}
	}, "reports", "get-team-report", "--team-id", "team_1", "--kind", "costs", "--output", "csv")
	want := "agent,costUsd\r\nagt_1,1.00\r\nagt_2,2.00\r\nagt_3,3.00\r\n"
	if r.Code != 0 || r.Out != want || len(r.Requests) != 3 {
		t.Fatalf("%+v out=%q", r, r.Out)
	}
}

// CLI-7: the JSON report's rows are paginated and merged.
func TestReportJSONRowsPaginate(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Query().Get("nextToken") == "" {
			clitest.Reply(w, 200, `{"kind":"usage","rows":[{"agent":"agt_1"},{"agent":"agt_2"}],"nextToken":"p2"}`)
			return
		}
		clitest.Reply(w, 200, `{"kind":"usage","rows":[{"agent":"agt_3"}]}`)
	}, "reports", "get-agent-report", "--team-id", "team_1", "--agent-id", "agt_1", "--kind", "usage")
	if r.Code != 0 || strings.Count(r.Out, `"agent"`) != 3 || len(r.Requests) != 2 {
		t.Fatalf("%+v", r)
	}
}

// CLI-US-4.2: the table follows --group-by: one column per group key and metric.
func TestReportTableColumnsFollowGroupBy(t *testing.T) {
	r := clitest.Run(t, ok, "reports", "get-team-report", "--team-id", "team_1", "--kind", "costs", "--output", "table")
	head := strings.SplitN(r.Out, "\n", 2)[0]
	if r.Code != 0 || !strings.Contains(head, "agent") || !strings.Contains(head, "costUsd") {
		t.Fatalf("%+v", r)
	}
}

// RPT-1..9, CLI-US-21-style: the service report is for platform roles; a team-admin gets AccessDenied.
func TestServiceReportAccessDeniedForTeamAdmin(t *testing.T) {
	r := clitest.Run(t, clitest.Denied("reports:service"), "reports", "get-service-report", "--kind", "costs")
	if r.Code != 254 || !strings.Contains(r.Err, "(AccessDenied) when calling the GetServiceReport operation") || !strings.Contains(r.Err, "reports:service") {
		t.Fatalf("%+v", r)
	}
}
