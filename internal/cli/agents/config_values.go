package agents

import (
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/output"
)

func init() { registrations = append(registrations, registerConfigValues) }

var configValueColumns = []output.Column{
	{Header: "Name", Path: "name"}, {Header: "Source", Path: "source"}, {Header: "Value", Path: "value"},
	{Header: "Resolved", Path: "resolvedValue"}, {Header: "Version", Path: "version"},
}

func registerConfigValues(r *cli.Registry) {
	r.Override(cli.Override{Op: "ListAgentConfigValues", Columns: configValueColumns, Examples: []cli.Example{
		{Title: "List an agent's configuration values", Command: "afs agents list-config-values --team-id " + exTeam + " --agent-id " + exAgent},
	}})
	r.Override(cli.Override{Op: "PutAgentConfigValue", Columns: configValueColumns, Examples: []cli.Example{
		{Title: "Set a literal value", Command: "afs agents put-config-value --team-id " + exTeam + " --agent-id " + exAgent + " --name LOG_LEVEL --source literal --value debug"},
		{Title: "Take the value from the team's namespace values", Command: "afs agents put-config-value --team-id " + exTeam + " --agent-id " + exAgent + " --name REGION --source namespace"},
	}})
	r.Override(cli.Override{Op: "DeleteAgentConfigValue", Examples: []cli.Example{
		{Title: "Remove a configuration value", Command: "afs agents delete-config-value --team-id " + exTeam + " --agent-id " + exAgent + " --name LOG_LEVEL"},
	}})
}
