package admin

import (
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/output"
)

func init() { registrations = append(registrations, registerIdentityProviders) }

var idpColumns = []output.Column{
	{Header: "Id", Path: "id"}, {Header: "Name", Path: "name"}, {Header: "Kind", Path: "kind"},
	{Header: "Issuer", Path: "issuerUrl"}, {Header: "Enabled", Path: "enabled"}, {Header: "Default", Path: "isDefault"},
	{Header: "Version", Path: "version"},
}

const idp = "idp_01J9ZQ3ABCDEFGHJKMNPQRSTVW"

func registerIdentityProviders(r *cli.Registry) {
	rename := map[string]string{"id": "provider-id"}
	r.Override(cli.Override{Op: "ListIdentityProviders", Columns: idpColumns, Examples: []cli.Example{
		{Title: "List identity providers", Command: "afs admin list-identity-providers"},
	}})
	r.Override(cli.Override{Op: "GetIdentityProvider", Rename: rename, Columns: idpColumns, Examples: []cli.Example{
		{Title: "Show one identity provider (never a secret, only its ARN)", Command: "afs admin describe-identity-provider --provider-id " + idp},
	}})
	r.Override(cli.Override{Op: "CreateIdentityProvider", Columns: idpColumns, Examples: []cli.Example{
		{Title: "Register an OIDC provider", Command: `afs admin create-identity-provider --name "Corp SSO" --kind okta --issuer-url https://idp.example.com --console-client-id console-app --cli-client-id afs-cli --client-secret-arn arn:aws:secretsmanager:eu-west-1:111122223333:secret:afs/idp-AbCdEf`},
	}})
	r.Override(cli.Override{Op: "UpdateIdentityProvider", Rename: rename, Columns: idpColumns, Examples: []cli.Example{
		{Title: "Disable a provider", Command: "afs admin update-identity-provider --provider-id " + idp + " --no-enabled"},
	}})
	r.Override(cli.Override{Op: "DeleteIdentityProvider", Rename: rename, Examples: []cli.Example{
		{Title: "Delete a provider that no team uses", Command: "afs admin delete-identity-provider --provider-id " + idp},
	}})
}
