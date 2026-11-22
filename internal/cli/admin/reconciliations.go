package admin

import (
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/output"
)

func init() { registrations = append(registrations, registerReconciliations) }

func registerReconciliations(r *cli.Registry) {
	r.Override(cli.Override{Op: "ListCostReconciliations", Columns: []output.Column{
		{Header: "Day", Path: "day"}, {Header: "TeamId", Path: "teamId"}, {Header: "GatewayUsd", Path: "gatewayUsd"},
		{Header: "BilledUsd", Path: "billedUsd"}, {Header: "VariancePct", Path: "variancePct"}, {Header: "Status", Path: "status"},
	}, Examples: []cli.Example{
		{Title: "Show the days where billing differs from the gateway", Command: "afs admin list-cost-reconciliations --status over_threshold"},
	}})
	r.Override(cli.Override{Op: "ListOrgAccounts", Columns: []output.Column{
		{Header: "AccountId", Path: "accountId"}, {Header: "Name", Path: "name"}, {Header: "Status", Path: "status"}, {Header: "SyncedAt", Path: "syncedAt"},
	}, Examples: []cli.Example{
		{Title: "List the AWS Organization accounts the service knows", Command: "afs admin list-org-accounts"},
	}})
}
