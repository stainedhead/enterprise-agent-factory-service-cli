package admin

import (
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/output"
)

func init() { registrations = append(registrations, registerModels) }

var adminModelColumns = []output.Column{
	{Header: "Id", Path: "id"}, {Header: "Provider", Path: "provider"}, {Header: "Name", Path: "displayName"},
	{Header: "Enabled", Path: "enabled"}, {Header: "Callable", Path: "callable", Null: "-"}, {Header: "Price", Path: "priceState", Null: "-"},
	{Header: "Origin", Path: "origin", Null: "-"}, {Header: "Teams", Path: "teamsWithAccess", Null: "-"},
	{Header: "Context", Path: "contextWindow"}, {Header: "Version", Path: "version"},
}

func registerModels(r *cli.Registry) {
	r.Override(cli.Override{Op: "ListAdminModels", Columns: adminModelColumns, Examples: []cli.Example{
		{Title: "List every model, enabled or not", Command: "afs admin list-models"},
	}})
	r.Override(cli.Override{Op: "CreateModel", Columns: adminModelColumns, Examples: []cli.Example{
		{Title: "Add a model", Command: `afs admin create-model --id anthropic.claude-sonnet-4-5 --provider anthropic --display-name "Claude Sonnet 4.5" --context-window 200000 --max-output-tokens 64000`},
	}})
	r.Override(cli.Override{Op: "UpdateModel", Columns: adminModelColumns, Examples: []cli.Example{
		{Title: "Disable a model (refused while service model access is all)", Command: "afs admin update-model --model-id anthropic.claude-sonnet-4-5 --no-enabled"},
	}})
}
