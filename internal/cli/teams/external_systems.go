package teams

import (
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/output"
)

func init() { registrations = append(registrations, registerExternalSystems) }

var externalSystemColumns = []output.Column{
	{Header: "Id", Path: "id"}, {Header: "Slug", Path: "slug"}, {Header: "Name", Path: "name"},
	{Header: "Category", Path: "category"}, {Header: "Enabled", Path: "enabled"}, {Header: "TeamState", Path: "teamState"},
	{Header: "Version", Path: "version"},
}

const exSystem = "ext_01JD2Q3R4S5T6V7W8X9Y0Z0001"

func registerExternalSystems(r *cli.Registry) {
	r.Override(cli.Override{Op: "ListTeamExternalSystems", Columns: externalSystemColumns, Examples: []cli.Example{
		{Title: "List the external systems the team may use", Command: "afs teams list-external-systems --team-id " + exTeam},
	}})
	r.Override(cli.Override{Op: "DeclareTeamExternalSystem", Columns: externalSystemColumns, Examples: []cli.Example{
		{Title: "Declare a team-owned external system", Command: `afs teams declare-external-system --team-id ` + exTeam + ` --slug payments-ledger --name "Ledger API" --category other --owner ledger-team@example.com --purpose "Reconcile payouts" --endpoints ledger.example.com:443`},
	}})
	r.Override(cli.Override{Op: "UpdateTeamExternalSystem", Columns: externalSystemColumns, Examples: []cli.Example{
		{Title: "Change an external system's endpoints", Command: "afs teams update-external-system --team-id " + exTeam + " --system-id " + exSystem + " --endpoints ledger.example.com:443 ledger-eu.example.com:443"},
	}})
	r.Override(cli.Override{Op: "SetTeamExternalSystemState", Columns: externalSystemColumns, Examples: []cli.Example{
		{Title: "Enable an external system for the team", Command: "afs teams set-external-system-state --team-id " + exTeam + " --system-id " + exSystem + " --enabled"},
	}})
	r.Override(cli.Override{Op: "DeleteTeamExternalSystem", Examples: []cli.Example{
		{Title: "Delete a team-owned external system", Command: "afs teams delete-external-system --team-id " + exTeam + " --system-id " + exSystem},
	}})
}
