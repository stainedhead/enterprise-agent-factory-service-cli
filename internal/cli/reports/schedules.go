package reports

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/output"
)

func init() { registrations = append(registrations, registerSchedules) }

var scheduleColumns = []output.Column{
	{Header: "Id", Path: "id"}, {Header: "Name", Path: "name"}, {Header: "Report", Path: "report"}, {Header: "Cadence", Path: "cadence"},
	{Header: "Enabled", Path: "enabled"}, {Header: "NextRunAt", Path: "nextRunAt"}, {Header: "LastStatus", Path: "lastStatus"},
	{Header: "ChannelId", Path: "channelId"},
}

const exChannel = "nch_01JD2Q3R4S5T6V7W8X9Y0Z0001"

// queryOptions build the stored report query from the same flags the reports use (CLI-US-17.1);
// --query-json gives the whole object instead (the name --query is the global JMESPath option).
var queryOptions = []cli.Option{
	{Flag: "--range", In: "local", Kind: cli.String, Enum: []string{"24h", "7d", "30d", "90d", "mtd"}, Desc: "Report range: 24h, 7d, 30d, 90d or mtd."},
	{Flag: "--from", In: "local", Kind: cli.String, Desc: "Start of the report window (RFC 3339); use with --to instead of --range."},
	{Flag: "--to", In: "local", Kind: cli.String, Desc: "End of the report window (RFC 3339)."},
	{Flag: "--group-by", In: "local", Kind: cli.String, Desc: "Comma-separated dimensions: team, agent, model, system, day, hour, pool, tag:<key>; surface, reason and person for the collaboration report."},
	tagOption,
}

var queryFlagNames = []string{"range", "from", "to", "group-by", "tag"}

func anyQueryFlag(p *cli.Parsed) bool {
	for _, n := range queryFlagNames {
		if p.Has(n) {
			return true
		}
	}
	return false
}

// scheduleQuery folds the report flags into the body's query member, over --query-json when given.
func scheduleQuery(p *cli.Parsed, parts *cli.RequestParts) error {
	if !anyQueryFlag(p) || parts.Body == nil {
		return nil
	}
	q := map[string]any{}
	if raw, ok := parts.Body.Get("query"); ok {
		if err := json.Unmarshal(raw, &q); err != nil {
			return fmt.Errorf("--query-json must be a JSON object: %w", err)
		}
	}
	if p.Has("range") {
		q["range"] = p.Value("range")
	}
	if p.Has("from") {
		q["from"] = p.Value("from")
	}
	if p.Has("to") {
		q["to"] = p.Value("to")
	}
	if p.Has("group-by") {
		q["groupBy"] = strings.Split(p.Value("group-by"), ",")
	}
	if tags := p.Values["tag"]; len(tags) > 0 {
		filters, _ := q["filters"].(map[string]any)
		if filters == nil {
			filters = map[string]any{}
		}
		for _, kv := range tags {
			k, v, _ := strings.Cut(kv, "=")
			filters["tag."+k] = v
		}
		q["filters"] = filters
	}
	raw, err := json.Marshal(q)
	if err != nil {
		return err
	}
	parts.Body.Set("query", raw)
	return nil
}

func scheduleFlags(op string) cli.Override {
	v := flags(op) // the shared checks; its Build is replaced below
	return cli.Override{Op: op, Extra: queryOptions, Validate: v.Validate, Build: scheduleQuery}
}

func registerSchedules(r *cli.Registry) {
	r.Override(cli.Override{Op: "ListReportSchedules", Columns: scheduleColumns, Examples: []cli.Example{
		{Title: "List a team's scheduled reports", Command: "afs reports list-report-schedules --team-id " + exTeam},
	}})
	r.Override(scheduleFlags("CreateReportSchedule"))
	r.Override(cli.Override{Op: "CreateReportSchedule", Columns: scheduleColumns, Examples: []cli.Example{
		{Title: "Email the cost-by-project report every week", Command: "afs reports create-report-schedule --team-id " + exTeam + " --name weekly-costs --report costs --cadence weekly --channel-id " + exChannel + " --range 30d --group-by tag:project"},
	}})
	r.Override(scheduleFlags("UpdateReportSchedule"))
	r.Override(cli.Override{Op: "UpdateReportSchedule", Columns: scheduleColumns, Examples: []cli.Example{
		{Title: "Pause a schedule", Command: "afs reports update-report-schedule --team-id " + exTeam + " --schedule-id rsc_01J9ZQ3ABCDEFGHJKMNPQRSTVW --no-enabled"},
	}})
	r.Override(cli.Override{Op: "DeleteReportSchedule", Examples: []cli.Example{
		{Title: "Delete a schedule", Command: "afs reports delete-report-schedule --team-id " + exTeam + " --schedule-id rsc_01J9ZQ3ABCDEFGHJKMNPQRSTVW"},
	}})
}
