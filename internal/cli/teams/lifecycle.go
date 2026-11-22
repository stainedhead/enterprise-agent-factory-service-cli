package teams

import "github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"

func init() { registrations = append(registrations, registerLifecycle) }

func registerLifecycle(r *cli.Registry) {
	r.Override(cli.Override{Op: "UpdateTeam", Columns: teamColumns, Examples: []cli.Example{
		{Title: "Change a team's description (reads the current version, sends If-Match)", Command: `afs teams update-team --team-id team_01J9ZQ3ABCDEFGHJKMNPQRSTVW --description "Payments platform"`},
		{Title: "Apply a change only if the team is still at version 7", Command: "afs teams update-team --team-id team_01J9ZQ3ABCDEFGHJKMNPQRSTVW --viewer-kubectl --if-match 7"},
	}})
	r.Override(cli.Override{Op: "DeleteTeam", Columns: teamColumns, Examples: []cli.Example{
		{Title: "Delete a team and wait until it is gone", Command: "afs teams delete-team --team-id team_01J9ZQ3ABCDEFGHJKMNPQRSTVW --wait"},
	}})
	r.Override(cli.Override{Op: "SuspendTeam", Columns: teamColumns, Examples: []cli.Example{
		{Title: "Suspend a team (writes are refused until it is resumed)", Command: "afs teams suspend-team --team-id team_01J9ZQ3ABCDEFGHJKMNPQRSTVW"},
	}})
	r.Override(cli.Override{Op: "ResumeTeam", Columns: teamColumns, Examples: []cli.Example{
		{Title: "Resume a suspended team", Command: "afs teams resume-team --team-id team_01J9ZQ3ABCDEFGHJKMNPQRSTVW"},
	}})
}
