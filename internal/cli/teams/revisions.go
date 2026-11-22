package teams

import (
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/output"
)

func init() { registrations = append(registrations, registerRevisions) }

var revisionColumns = []output.Column{
	{Header: "Revision", Path: "revision"}, {Header: "Kind", Path: "resourceKind"}, {Header: "Resource", Path: "resourceId"},
	{Header: "Operation", Path: "operation"}, {Header: "ChangedBy", Path: "changedBy"}, {Header: "ChangedAt", Path: "changedAt"},
}

func registerRevisions(r *cli.Registry) {
	r.Override(cli.Override{Op: "ListRevisions", Columns: revisionColumns, Examples: []cli.Example{
		{Title: "List the configuration revisions of a team", Command: "afs teams list-revisions --team-id team_01J9ZQ3ABCDEFGHJKMNPQRSTVW"},
	}})
	r.Override(cli.Override{Op: "GetRevision", Columns: revisionColumns, Examples: []cli.Example{
		{Title: "Show one revision with its snapshot", Command: "afs teams describe-revision --team-id team_01J9ZQ3ABCDEFGHJKMNPQRSTVW --rev 12 --output json"},
	}})
}
