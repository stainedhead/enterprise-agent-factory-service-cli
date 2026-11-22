package members

import (
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/output"
)

func init() { registrations = append(registrations, registerMembers) }

var memberColumns = []output.Column{
	{Header: "Id", Path: "id"}, {Header: "Subject", Path: "subjectType"}, {Header: "User", Path: "userId"},
	{Header: "Email", Path: "email"}, {Header: "Group", Path: "groupName"}, {Header: "Role", Path: "role"},
	{Header: "Version", Path: "version"},
}

const team = "team_01J9ZQ3ABCDEFGHJKMNPQRSTVW"

func registerMembers(r *cli.Registry) {
	r.Override(cli.Override{Op: "ListMembers", Columns: memberColumns, Examples: []cli.Example{
		{Title: "List a team's members and groups", Command: "afs members list-members --team-id " + team},
	}})
	r.Override(cli.Override{Op: "AddMember", Columns: memberColumns, Examples: []cli.Example{
		{Title: "Add a user as team_admin", Command: "afs members add-member --team-id " + team + " --subject-type user --email ada@example.com --role team_admin"},
		{Title: "Bind an identity provider group", Command: "afs members add-member --team-id " + team + " --subject-type group --group-name payments-devs --role team_viewer"},
	}})
	r.Override(cli.Override{Op: "UpdateMember", Columns: memberColumns, Examples: []cli.Example{
		{Title: "Change a member's role (reads the version from the member list)", Command: "afs members update-member --team-id " + team + " --binding-id mem_01J9ZQ3ABCDEFGHJKMNPQRSTVW --role team_viewer"},
	}})
	r.Override(cli.Override{Op: "RemoveMember", Examples: []cli.Example{
		{Title: "Remove a member", Command: "afs members remove-member --team-id " + team + " --binding-id mem_01J9ZQ3ABCDEFGHJKMNPQRSTVW"},
	}})
}
