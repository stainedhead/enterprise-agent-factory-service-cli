package budgets

import (
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/output"
)

func init() { registrations = append(registrations, registerUsage) }

func registerUsage(r *cli.Registry) {
	// Usage rows are shaped by --group-by: the columns follow the response (key.agent, key.model, ...).
	r.Override(cli.Override{Op: "GetUsage", Dynamic: true, Examples: []cli.Example{
		{Title: "Usage by agent since the first of the month", Command: "afs budgets get-usage --team-id " + exTeam + " --from 2026-11-01T00:00:00Z --group-by agent"},
		{Title: "Daily usage as CSV", Command: "afs budgets get-usage --team-id " + exTeam + " --group-by day --output csv > usage.csv"},
	}})
	r.Override(cli.Override{Op: "ListCosts", Columns: []output.Column{
		{Header: "Day", Path: "day"}, {Header: "Source", Path: "source"}, {Header: "Agent", Path: "agentId"},
		{Header: "Pool", Path: "poolId"}, {Header: "Model", Path: "modelId"}, {Header: "CostUsd", Path: "costUsd"},
		{Header: "Calls", Path: "calls"},
	}, Examples: []cli.Example{
		{Title: "List daily costs (gateway spend and the AWS cost and usage report)", Command: "afs budgets list-costs --team-id " + exTeam + " --from 2026-11-01T00:00:00Z --to 2026-11-21T00:00:00Z"},
		{Title: "Export costs as CSV", Command: "afs budgets list-costs --team-id " + exTeam + " --output csv > costs.csv"},
	}})
}
