package budgets

import (
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/output"
)

func init() { registrations = append(registrations, registerBudgets) }

const exTeam = "team_01JD2Q3R4S5T6V7W8X9Y0Z0001"

// budgetColumns: amount, spent, reserved, percent used and forecast (CLI-US-14.3).
var budgetColumns = []output.Column{
	{Header: "Period", Path: "period"}, {Header: "Amount", Path: "amountUsd"}, {Header: "Spent", Path: "spentUsd"},
	{Header: "Reserved", Path: "reservedUsd"}, {Header: "PercentUsed", Path: "percentUsed"}, {Header: "Forecast", Path: "forecastUsd"},
	{Header: "OnExhaust", Path: "onExhaust"}, {Header: "Version", Path: "version"},
}

func registerBudgets(r *cli.Registry) {
	r.Override(cli.Override{Op: "GetBudget", Columns: budgetColumns, Examples: []cli.Example{
		{Title: "Show the team's budget and how much is used", Command: "afs budgets get-budget --team-id " + exTeam + " --output table"},
	}})
	r.Override(cli.Override{Op: "PutBudget", Columns: budgetColumns, Examples: []cli.Example{
		{Title: "Set a monthly budget with alert thresholds (percent)", Command: "afs budgets put-budget --team-id " + exTeam + " --period monthly --amount-usd 100.50 --thresholds 50 80 100"},
		{Title: "Replace the budget only if it is still at version 3", Command: "afs budgets put-budget --team-id " + exTeam + " --period monthly --amount-usd 250.00 --on-exhaust notify --if-match 3"},
	}})
}
