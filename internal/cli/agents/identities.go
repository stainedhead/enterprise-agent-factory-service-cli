package agents

import (
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/output"
)

func init() { registrations = append(registrations, registerIdentities) }

var identityColumns = []output.Column{
	{Header: "Id", Path: "id"}, {Header: "Kind", Path: "kind"}, {Header: "ExternalRef", Path: "externalRef"},
	{Header: "SecretId", Path: "secretId"}, {Header: "Status", Path: "status"}, {Header: "Version", Path: "version"},
}

func registerIdentities(r *cli.Registry) {
	r.Override(cli.Override{Op: "ListAgentIdentities", Columns: identityColumns, Examples: []cli.Example{
		{Title: "List an agent's identities", Command: "afs agents list-identities --team-id " + exTeam + " --agent-id " + exAgent},
	}})
	r.Override(cli.Override{Op: "CreateAgentIdentity", Columns: identityColumns, Examples: []cli.Example{
		{Title: "Register a GitHub App identity", Command: "afs agents create-identity --team-id " + exTeam + " --agent-id " + exAgent + " --kind github_app --external-ref github.com/example-org/apps/afs-coder --secret-id sec_01J9ZQ3ABCDEFGHJKMNPQRSTVW"},
	}})
	r.Override(cli.Override{Op: "DeleteAgentIdentity", Examples: []cli.Example{
		{Title: "Delete an identity", Command: "afs agents delete-identity --team-id " + exTeam + " --agent-id " + exAgent + " --identity-id aid_01J9ZQ3ABCDEFGHJKMNPQRSTVW"},
	}})
}
