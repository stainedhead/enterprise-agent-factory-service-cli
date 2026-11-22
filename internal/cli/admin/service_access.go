package admin

import (
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli/apikeys"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/output"
)

func init() { registrations = append(registrations, registerServiceAccess) }

var operatorColumns = []output.Column{
	{Header: "Id", Path: "id"}, {Header: "SubjectType", Path: "subjectType"}, {Header: "Email", Path: "email"},
	{Header: "GroupName", Path: "groupName"}, {Header: "UserId", Path: "userId"}, {Header: "IdentityProviderId", Path: "identityProviderId"},
}

func registerServiceAccess(r *cli.Registry) {
	r.Override(cli.Override{Op: "ListServiceApiKeys", Columns: apikeys.KeyColumns, Examples: []cli.Example{
		{Title: "List the service-level API keys", Command: "afs admin list-api-keys"},
	}})
	r.Override(cli.Override{Op: "CreateServiceApiKey", Columns: apikeys.CreatedKeyColumns, After: apikeys.ShowSecretNotice, Examples: []cli.Example{
		{Title: "Create a security operator key (the secret is shown once)", Command: "afs admin create-api-key --name cve-pipeline --role security_operator --expires-at 2027-02-19T18:00:00Z"},
	}})
	r.Override(cli.Override{Op: "DeleteServiceApiKey", Examples: []cli.Example{
		{Title: "Delete a service-level API key", Command: "afs admin delete-api-key --key-id key_01J9ZQ3ABCDEFGHJKMNPQRSTVW"},
	}})
	r.Override(cli.Override{Op: "ListSecurityOperators", Columns: operatorColumns, Examples: []cli.Example{
		{Title: "List the security operators", Command: "afs admin list-security-operators"},
	}})
	r.Override(cli.Override{Op: "AddSecurityOperator", Columns: operatorColumns, Examples: []cli.Example{
		{Title: "Make a user a security operator", Command: "afs admin add-security-operator --subject-type user --email ana@example.com --identity-provider-id idp_01J9ZQ3ABCDEFGHJKMNPQRSTVW"},
		{Title: "Make an identity provider group security operators", Command: "afs admin add-security-operator --subject-type group --group-name security-ops --identity-provider-id idp_01J9ZQ3ABCDEFGHJKMNPQRSTVW"},
	}})
	r.Override(cli.Override{Op: "RemoveSecurityOperator", Examples: []cli.Example{
		{Title: "Remove a security operator", Command: "afs admin remove-security-operator --binding-id rb_01J9ZQ3ABCDEFGHJKMNPQRSTVW"},
	}})
}
