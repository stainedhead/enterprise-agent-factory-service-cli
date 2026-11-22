package agents

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

func registerExternalSystems(r *cli.Registry) {
	r.Override(cli.Override{Op: "ListAgentExternalSystems", Columns: externalSystemColumns, Examples: []cli.Example{
		{Title: "List the external systems an agent may reach", Command: "afs agents list-external-systems --team-id " + exTeam + " --agent-id " + exAgent},
	}})
	r.Override(cli.Override{Op: "SetAgentExternalSystemState", Columns: externalSystemColumns, Examples: []cli.Example{
		{Title: "Disable an external system for one agent", Command: "afs agents set-external-system-state --team-id " + exTeam + " --agent-id " + exAgent + " --system-id ext_01JD2Q3R4S5T6V7W8X9Y0Z0001 --no-enabled"},
	}})
}
