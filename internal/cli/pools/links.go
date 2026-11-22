package pools

import (
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/output"
)

func init() { registrations = append(registrations, registerLinks) }

var (
	poolExternalSystemColumns = []output.Column{
		{Header: "Id", Path: "id"}, {Header: "Slug", Path: "slug"}, {Header: "Name", Path: "name"},
		{Header: "Category", Path: "category"}, {Header: "Enabled", Path: "enabled"}, {Header: "TeamState", Path: "teamState"},
		{Header: "Version", Path: "version"},
	}

	poolRoleLinkColumns = []output.Column{
		{Header: "Id", Path: "id"}, {Header: "Profile", Path: "profileName"}, {Header: "AccountId", Path: "accountId"},
		{Header: "RoleName", Path: "roleName"}, {Header: "Status", Path: "status"}, {Header: "ExternalId", Path: "externalId"},
		{Header: "LastVerifiedAt", Path: "lastVerifiedAt"}, {Header: "LastError", Path: "lastError"},
	}
)

func registerLinks(r *cli.Registry) {
	r.Override(cli.Override{Op: "ListPoolExternalSystems", Columns: poolExternalSystemColumns, Examples: []cli.Example{
		{Title: "List the external systems a pool may reach", Command: "afs pools list-external-systems --team-id " + exTeam + " --pool-id " + exPool},
	}})
	r.Override(cli.Override{Op: "PutPoolExternalSystem", Columns: poolExternalSystemColumns, Examples: []cli.Example{
		{Title: "Disable an external system for a pool", Command: "afs pools put-external-system --team-id " + exTeam + " --pool-id " + exPool + " --system-id ext_01J9ZQ3ABCDEFGHJKMNPQRSTVW --no-enabled"},
	}})
	r.Override(cli.Override{Op: "ListPoolRoleLinks", Columns: poolRoleLinkColumns, Examples: []cli.Example{
		{Title: "List a pool's role links", Command: "afs pools list-role-links --team-id " + exTeam + " --pool-id " + exPool},
	}})
	r.Override(cli.Override{Op: "CreatePoolRoleLink", Columns: poolRoleLinkColumns, Examples: []cli.Example{
		{Title: "Link a role by name in an account", Command: "afs pools create-role-link --team-id " + exTeam + " --pool-id " + exPool + " --account-id 111122223333 --role-name afs-cve-reader --profile-name target"},
	}})
	r.Override(cli.Override{Op: "VerifyPoolRoleLink", Columns: poolRoleLinkColumns, Examples: []cli.Example{
		{Title: "Check that the trust policy is in place", Command: "afs pools verify-role-link --team-id " + exTeam + " --pool-id " + exPool + " --link-id " + exLink},
	}})
	r.Override(cli.Override{Op: "DeletePoolRoleLink", Examples: []cli.Example{
		{Title: "Delete a role link", Command: "afs pools delete-role-link --team-id " + exTeam + " --pool-id " + exPool + " --link-id " + exLink},
	}})
}
