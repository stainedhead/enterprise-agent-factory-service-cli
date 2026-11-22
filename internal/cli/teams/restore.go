package teams

import (
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/output"
)

func init() { registrations = append(registrations, registerRestore) }

var allowedAccountColumns = []output.Column{
	{Header: "AccountId", Path: "accountId"}, {Header: "Name", Path: "name"}, {Header: "Status", Path: "status"},
}

// restoreColumns show the outcome of a restore (or its plan under --dry-run).
var restoreColumns = []output.Column{
	{Header: "Revision", Path: "revision"}, {Header: "DryRun", Path: "dryRun"}, {Header: "Changes", Path: "changes"}, {Header: "Skipped", Path: "skipped"},
}

func registerRestore(r *cli.Registry) {
	r.Override(cli.Override{Op: "ListAllowedAccounts", Columns: allowedAccountColumns, Examples: []cli.Example{
		{Title: "List the AWS accounts a team may link roles in", Command: "afs teams list-allowed-accounts --team-id " + exTeam},
	}})
	r.Override(cli.Override{Op: "RestoreRevision", Columns: restoreColumns, Rename: map[string]string{"rev": "revision-id"}, Examples: []cli.Example{
		{Title: "Show what restoring revision 41 would change", Command: "afs teams restore-revision --team-id " + exTeam + " --revision-id 41 --dry-run"},
		{Title: "Restore revision 41", Command: "afs teams restore-revision --team-id " + exTeam + " --revision-id 41"},
	}})
}
