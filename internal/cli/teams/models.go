package teams

import (
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/output"
)

func init() { registrations = append(registrations, registerModels) }

// modelColumns are the table columns of a Model.
var modelColumns = []output.Column{
	{Header: "Id", Path: "id"}, {Header: "Provider", Path: "provider"}, {Header: "Name", Path: "displayName"},
	{Header: "Enabled", Path: "enabled"}, {Header: "Context", Path: "contextWindow"}, {Header: "MaxOutput", Path: "maxOutputTokens"},
}

func registerModels(r *cli.Registry) {
	r.Override(cli.Override{Op: "ListTeamModels", Columns: modelColumns, Examples: []cli.Example{
		{Title: "List the models a team may use", Command: "afs teams list-team-models --team-id team_01J9ZQ3ABCDEFGHJKMNPQRSTVW"},
	}})
	r.Override(cli.Override{Op: "PutTeamModels", Columns: modelColumns, Examples: []cli.Example{
		{Title: "Replace the team's model list", Command: "afs teams put-team-models --team-id team_01J9ZQ3ABCDEFGHJKMNPQRSTVW --model-ids claude-sonnet claude-haiku"},
	}})
}
