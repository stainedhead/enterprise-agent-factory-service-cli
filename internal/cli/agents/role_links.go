package agents

import (
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/output"
)

func init() { registrations = append(registrations, registerRoleLinks) }

// roleLinkColumns are the table columns of an agent role link. The trust policy is shown by
// describe-role-link --query trustPolicy (CLI-US-13.2), not in the table.
var roleLinkColumns = []output.Column{
	{Header: "Id", Path: "id"}, {Header: "Profile", Path: "profileName"}, {Header: "AccountId", Path: "accountId"},
	{Header: "RoleArn", Path: "roleArn"}, {Header: "Status", Path: "status"}, {Header: "ExternalId", Path: "externalId"},
	{Header: "LastVerifiedAt", Path: "lastVerifiedAt"}, {Header: "LastError", Path: "lastError"},
}

func registerRoleLinks(r *cli.Registry) {
	r.Override(cli.Override{Op: "ListAgentRoleLinks", Columns: roleLinkColumns, Examples: []cli.Example{
		{Title: "List an agent's role links", Command: "afs agents list-role-links --team-id " + exTeam + " --agent-id " + exAgent},
	}})
	r.Override(cli.Override{Op: "CreateAgentRoleLink", Columns: roleLinkColumns, Examples: []cli.Example{
		{Title: "Link a role in another account (prints the trust policy's external id)", Command: "afs agents create-role-link --team-id " + exTeam + " --agent-id " + exAgent + " --role-arn arn:aws:iam::111122223333:role/payments-deployer --profile-name deploy"},
	}})
	r.Override(cli.Override{Op: "GetAgentRoleLink", Columns: roleLinkColumns, Examples: []cli.Example{
		{Title: "Show a role link", Command: "afs agents describe-role-link --team-id " + exTeam + " --agent-id " + exAgent + " --link-id arl_01J9ZQ3ABCDEFGHJKMNPQRSTVW"},
		{Title: "Print only the trust policy to paste into IAM", Command: "afs agents describe-role-link --team-id " + exTeam + " --agent-id " + exAgent + " --link-id arl_01J9ZQ3ABCDEFGHJKMNPQRSTVW --query trustPolicy --output json"},
	}})
	r.Override(cli.Override{Op: "VerifyAgentRoleLink", Columns: roleLinkColumns, Examples: []cli.Example{
		{Title: "Check that the trust policy is in place", Command: "afs agents verify-role-link --team-id " + exTeam + " --agent-id " + exAgent + " --link-id arl_01J9ZQ3ABCDEFGHJKMNPQRSTVW"},
	}})
	r.Override(cli.Override{Op: "DeleteAgentRoleLink", Examples: []cli.Example{
		{Title: "Delete a role link", Command: "afs agents delete-role-link --team-id " + exTeam + " --agent-id " + exAgent + " --link-id arl_01J9ZQ3ABCDEFGHJKMNPQRSTVW"},
	}})
}
