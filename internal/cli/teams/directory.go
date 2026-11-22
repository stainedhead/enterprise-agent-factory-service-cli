package teams

import (
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/output"
)

func init() { registrations = append(registrations, registerDirectory) }

func registerDirectory(r *cli.Registry) {
	r.Override(cli.Override{Op: "ListTeamDirectory", Columns: []output.Column{
		{Header: "Id", Path: "id"}, {Header: "Slug", Path: "slug"}, {Header: "DisplayName", Path: "displayName"},
	}, Examples: []cli.Example{
		{Title: "Look up team ids to name as repository owners", Command: "afs teams list-team-directory"},
	}})
}
