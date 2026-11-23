package teams

import (
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/output"
)

func init() { registrations = append(registrations, registerTeams) }

// teamColumns are the table columns of every operation that returns a Team.
var teamColumns = []output.Column{
	{Header: "Id", Path: "id"}, {Header: "Slug", Path: "slug"}, {Header: "Name", Path: "displayName"},
	{Header: "Kind", Path: "kind"}, {Header: "State", Path: "state"}, {Header: "Phase", Path: "status.phase"},
	{Header: "Limit", Path: "maxCollaborationExposure", Null: "-"}, {Header: "Models", Path: "modelAccess", Null: "-"},
	{Header: "Version", Path: "version"},
}

// statusColumns are the columns of the *Status operations (id, kind, generation, phase).
var statusColumns = []output.Column{
	{Header: "Id", Path: "id"}, {Header: "Kind", Path: "kind"}, {Header: "Generation", Path: "generation"},
	{Header: "ObservedGeneration", Path: "status.observedGeneration"}, {Header: "Phase", Path: "status.phase"},
}

func registerTeams(r *cli.Registry) {
	r.Override(cli.Override{Op: "ListTeams", Columns: teamColumns, Examples: []cli.Example{
		{Title: "List the teams you can see", Command: "afs teams list-teams"},
		{Title: "List suspended teams as ids and names only", Command: "afs teams list-teams --state suspended --query 'items[].[id,displayName]' --output text"},
	}})
	r.Override(cli.Override{Op: "GetTeam", Columns: teamColumns, Examples: []cli.Example{
		{Title: "Show one team", Command: "afs teams describe-team --team-id team_01J9ZQ3ABCDEFGHJKMNPQRSTVW"},
	}})
	r.Override(cli.Override{Op: "GetTeamStatus", Columns: statusColumns, Examples: []cli.Example{
		{Title: "Show reconcile status", Command: "afs teams get-team-status --team-id team_01J9ZQ3ABCDEFGHJKMNPQRSTVW"},
	}})
	r.Override(cli.Override{Op: "CreateTeam", Columns: teamColumns, Kinds: map[string]cli.Kind{"tags": cli.TagPairs}, Examples: []cli.Example{
		{Title: "Create a persistent team and wait until it is ready", Command: `afs teams create-team --slug payments --display-name "Payments" --kind persistent --max-budget-usd 500.00 --first-admin '{"email":"ada@example.com"}' --tags project=checkout cost-center=CC-1001 --wait`},
		{Title: "Create a team from a file", Command: "afs teams create-team --cli-input-json file://team.json"},
		{Title: "Print the input skeleton", Command: "afs teams create-team --generate-cli-skeleton"},
	}})
}
