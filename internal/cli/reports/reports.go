package reports

import "github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"

func init() { registrations = append(registrations, registerReports) }

const (
	exTeam  = "team_01JD2Q3R4S5T6V7W8X9Y0Z0001"
	exAgent = "agt_01JD2Q3R4S5T6V7W8X9Y0Z0001"
)

func registerReports(r *cli.Registry) {
	for _, op := range []string{"GetTeamReport", "GetAgentReport", "GetServiceReport"} {
		r.Override(flags(op))
		r.Override(cli.Override{Op: op, Dynamic: true}) // columns follow --group-by
	}
	r.Override(cli.Override{Op: "GetTeamReport", Examples: []cli.Example{
		{Title: "Cost by project tag over 30 days, as CSV", Command: "afs reports get-team-report --team-id " + exTeam + " --kind costs --range 30d --group-by tag:project --output csv"},
		{Title: "Usage by agent and model, filtered by tag", Command: "afs reports get-team-report --team-id " + exTeam + " --kind usage --range 7d --group-by agent,model --tag project=checkout"},
	}})
	r.Override(cli.Override{Op: "GetAgentReport", Examples: []cli.Example{
		{Title: "One agent's activity over 24 hours", Command: "afs reports get-agent-report --team-id " + exTeam + " --agent-id " + exAgent + " --kind activity --range 24h"},
	}})
	r.Override(cli.Override{Op: "GetServiceReport", Examples: []cli.Example{
		{Title: "Service-wide costs by team this month (platform roles)", Command: "afs reports get-service-report --kind costs --range mtd --group-by team"},
	}})
}
