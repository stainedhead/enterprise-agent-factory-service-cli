package teams

import (
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/output"
)

func init() { registrations = append(registrations, registerConfigValues) }

const exTeam = "team_01JD2Q3R4S5T6V7W8X9Y0Z0001"

var configValueColumns = []output.Column{
	{Header: "Name", Path: "name"}, {Header: "Source", Path: "source"}, {Header: "Value", Path: "value"},
	{Header: "Resolved", Path: "resolvedValue"}, {Header: "Version", Path: "version"},
}

func registerConfigValues(r *cli.Registry) {
	r.Override(cli.Override{Op: "ListNamespaceConfigValues", Columns: configValueColumns, Examples: []cli.Example{
		{Title: "List the team's namespace configuration values", Command: "afs teams list-config-values --team-id " + exTeam},
	}})
	r.Override(cli.Override{Op: "PutNamespaceConfigValue", Columns: configValueColumns, Examples: []cli.Example{
		{Title: "Set a team-wide value (reads the version first; a missing value is created)", Command: "afs teams put-config-value --team-id " + exTeam + " --name REGION --value eu-west-1"},
	}})
	r.Override(cli.Override{Op: "DeleteNamespaceConfigValue", Examples: []cli.Example{
		{Title: "Remove a team-wide value", Command: "afs teams delete-config-value --team-id " + exTeam + " --name REGION"},
	}})
}
