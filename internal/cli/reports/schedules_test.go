package reports_test

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/clitest"
)

const schedule = `{"id":"rsc_1","name":"weekly-costs","report":"costs","cadence":"weekly","channelId":"nch_1","enabled":true,"format":"csv","nextRunAt":"2026-11-28T06:00:00Z","lastStatus":null,"version":1,"query":{"range":"30d","groupBy":["tag:project"]}}`

func TestCommandsDoneSchedules(t *testing.T) {
	clitest.CheckOperations(t, "schedules", "ListReportSchedules", "CreateReportSchedule", "UpdateReportSchedule", "DeleteReportSchedule")
}

// CLI-US-17.1: the report flags become the stored query, and nextRunAt is printed.
func TestCreateReportScheduleQueryFromFlags(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) { clitest.Reply(w, 201, schedule) },
		"reports", "create-report-schedule", "--team-id", "team_1", "--name", "weekly-costs", "--report", "costs", "--cadence", "weekly",
		"--channel-id", "nch_1", "--range", "30d", "--group-by", "tag:project", "--tag", "env=prod", "--output", "table")
	if r.Code != 0 || !strings.Contains(r.Out, "2026-11-28T06:00:00Z") {
		t.Fatalf("%+v", r)
	}
	var body struct {
		Name    string `json:"name"`
		Channel string `json:"channelId"`
		Query   struct {
			Range   string            `json:"range"`
			GroupBy []string          `json:"groupBy"`
			Filters map[string]string `json:"filters"`
		} `json:"query"`
	}
	if err := json.Unmarshal(r.ReqBody, &body); err != nil {
		t.Fatal(err)
	}
	if body.Query.Range != "30d" || len(body.Query.GroupBy) != 1 || body.Query.GroupBy[0] != "tag:project" || body.Query.Filters["tag.env"] != "prod" || body.Channel != "nch_1" {
		t.Errorf("body %s", r.ReqBody)
	}
	if r.Req.Header.Get("Idempotency-Key") == "" || r.Req.URL.Path != "/v1/teams/team_1/report-schedules" {
		t.Errorf("%s %v", r.Req.URL.Path, r.Req.Header)
	}
	// --query-json is the whole object; the flags are folded over it
	r = clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) { clitest.Reply(w, 201, schedule) },
		"reports", "create-report-schedule", "--team-id", "team_1", "--name", "n", "--report", "usage", "--cadence", "daily", "--channel-id", "nch_1",
		"--query-json", `{"groupBy":["agent"],"range":"7d"}`, "--range", "24h")
	if r.Code != 0 || !strings.Contains(string(r.ReqBody), `"groupBy":["agent"]`) || !strings.Contains(string(r.ReqBody), `"range":"24h"`) {
		t.Fatalf("%+v body=%s", r, r.ReqBody)
	}
	// the same checks as the reports: --range excludes --from
	r = clitest.Run(t, nil, "reports", "create-report-schedule", "--team-id", "team_1", "--name", "n", "--report", "usage", "--cadence", "daily", "--channel-id", "nch_1", "--range", "7d", "--from", "2026-01-01T00:00:00Z")
	if r.Code != 252 || len(r.Requests) != 0 {
		t.Fatalf("%+v", r)
	}
}

// CLI-US-17.2: a service-level or another team's channel is the API's ValidationError (254).
func TestReportScheduleForeignChannel254(t *testing.T) {
	r := clitest.Run(t, clitest.Fail(400, "ValidationError", "channel nch_other does not belong to this team"),
		"reports", "create-report-schedule", "--team-id", "team_1", "--name", "n", "--report", "usage", "--cadence", "daily", "--channel-id", "nch_other")
	if r.Code != 254 || !strings.Contains(r.Err, "(ValidationError) when calling the CreateReportSchedule operation") {
		t.Fatalf("%+v", r)
	}
}

// CLI-US-17.2: a viewer is refused by the API (254).
func TestReportScheduleViewerDenied(t *testing.T) {
	r := clitest.Run(t, clitest.Denied("reports:schedule"),
		"reports", "create-report-schedule", "--team-id", "team_1", "--name", "n", "--report", "usage", "--cadence", "daily", "--channel-id", "nch_1")
	if r.Code != 254 || !strings.Contains(r.Err, "reports:schedule") {
		t.Fatalf("%+v", r)
	}
}
