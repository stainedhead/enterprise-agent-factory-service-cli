package reports_test

import (
	"net/http"
	"testing"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/clitest"
)

// FR-011: the two inference reports.
func TestCommandsDoneInference(t *testing.T) {
	clitest.CheckOperations(t, "inference", "GetInferenceReport", "GetTeamInferenceReport")
}

func inference(w http.ResponseWriter, _ *http.Request) {
	clitest.Reply(w, 200, `{"groupBy":["route"],"rows":[{"route":"openai-prod","calls":3,"ttftP95Ms":1310}],"totals":{"calls":3}}`)
}

// FR-011: every filter reaches the query of the right path.
func TestInferenceReportQuery(t *testing.T) {
	r := clitest.Run(t, inference, "reports", "get-team-inference-report", "--team-id", "team_1", "--range", "7d",
		"--group-by", "route,agent", "--route", "openai-prod", "--model-id", "gpt-4.1", "--agent-id", "agt_1")
	if r.Code != 0 || r.Req.URL.Path != "/v1/teams/team_1/reports/inference" {
		t.Fatalf("%+v", r)
	}
	q := r.Req.URL.Query()
	for k, v := range map[string]string{"range": "7d", "groupBy": "route,agent", "route": "openai-prod", "modelId": "gpt-4.1", "agentId": "agt_1"} {
		if q.Get(k) != v {
			t.Errorf("%s=%q want %q", k, q.Get(k), v)
		}
	}
	r = clitest.Run(t, inference, "reports", "get-inference-report", "--team-id", "team_2", "--from", "2026-11-01T00:00:00Z", "--to", "2026-11-02T00:00:00Z")
	if q := r.Req.URL.Query(); r.Code != 0 || r.Req.URL.Path != "/v1/reports/inference" || q.Get("teamId") != "team_2" || q.Get("from") == "" {
		t.Fatalf("%+v", r)
	}
	c := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/csv")
		_, _ = w.Write([]byte("route,calls\nopenai-prod,3\n"))
	}, "reports", "get-inference-report", "--output", "csv")
	if c.Code != 0 || c.Req.Header.Get("Accept") != "text/csv" {
		t.Errorf("csv: %+v", c)
	}
}

// FR-011: --range is exclusive with --from and --to and --group-by is checked, before any call.
func TestInferenceReportUsageErrors(t *testing.T) {
	noCall := func(http.ResponseWriter, *http.Request) { t.Error("a call was made") }
	for name, args := range map[string][]string{
		"range and from": {"--range", "7d", "--from", "2026-11-01T00:00:00Z"},
		"bad group-by":   {"--group-by", "planet"},
		"tag group-by":   {"--group-by", "tag:project"},
		"bad range":      {"--range", "1y"},
	} {
		if r := clitest.Run(t, noCall, append([]string{"reports", "get-inference-report"}, args...)...); r.Code != 252 {
			t.Errorf("%s: %+v", name, r)
		}
	}
}

// FR-013: usage and cost reports take --route and --group-by route (v1.5.0).
func TestRouteFilterAndGroupBy(t *testing.T) {
	r := clitest.Run(t, ok, "reports", "get-team-report", "--team-id", "team_1", "--kind", "costs", "--group-by", "route,model", "--route", "openai-prod")
	if q := r.Req.URL.Query(); r.Code != 0 || q.Get("groupBy") != "route,model" || q.Get("route") != "openai-prod" {
		t.Fatalf("%+v", r)
	}
}
