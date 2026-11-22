package teams

import (
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/output"
)

func init() { registrations = append(registrations, registerTags) }

// tagColumns show every tag with where it comes from (system, service, team, service_default, missing).
var tagColumns = []output.Column{
	{Header: "Key", Path: "key"}, {Header: "Value", Path: "value"}, {Header: "Source", Path: "source"},
}

func registerTags(r *cli.Registry) {
	r.Override(cli.Override{Op: "GetTeamTags", Columns: tagColumns, ResultKey: "tags", Examples: []cli.Example{
		{Title: "Show every tag of a team and where it comes from", Command: "afs teams get-team-tags --team-id " + exTeam},
		{Title: "Print the tags as JSON", Command: "afs teams get-team-tags --team-id " + exTeam + " --output json"},
	}})
	r.Override(cli.Override{Op: "PutTeamTags", Columns: tagColumns, ResultKey: "tags", Kinds: map[string]cli.Kind{"tags": cli.TagPairs}, Examples: []cli.Example{
		{Title: "Replace the team-owned tags", Command: "afs teams put-team-tags --team-id " + exTeam + " --tags cost-center=CC-2002 project=checkout"},
	}})
}
