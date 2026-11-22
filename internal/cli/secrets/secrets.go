package secrets

import (
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/output"
)

func init() { registrations = append(registrations, registerSecrets) }

const (
	exTeam = "team_01JD2Q3R4S5T6V7W8X9Y0Z0001"
	exID   = "sec_01JD2Q3R4S5T6V7W8X9Y0Z0001"
)

// secretColumns carry names and metadata only; the value is never returned by the API.
var secretColumns = []output.Column{
	{Header: "Id", Path: "id"}, {Header: "Name", Path: "name"}, {Header: "Scope", Path: "scope"},
	{Header: "Agent", Path: "agentId"}, {Header: "Pool", Path: "poolId"}, {Header: "ManagedBy", Path: "managedBy"},
	{Header: "LastChanged", Path: "lastChangedAt"}, {Header: "Version", Path: "version"},
}

func registerSecrets(r *cli.Registry) {
	r.Override(cli.Override{Op: "ListSecrets", Columns: secretColumns, Examples: []cli.Example{
		{Title: "List a team's secrets (names and metadata, never values)", Command: "afs secrets list-secrets --team-id " + exTeam},
		{Title: "List the secrets that AFS generated", Command: "afs secrets list-secrets --team-id " + exTeam + " --managed-by afs"},
	}})
	r.Override(cli.Override{Op: "GetSecret", Columns: secretColumns, Examples: []cli.Example{
		{Title: "Show one secret's metadata", Command: "afs secrets describe-secret --team-id " + exTeam + " --secret-id " + exID},
	}})
	r.Override(cli.Override{Op: "CreateSecret", Columns: secretColumns, Examples: []cli.Example{
		{Title: "Create a secret from a file", Command: "afs secrets create-secret --team-id " + exTeam + " --name github-pat --scope team --secret-value file://token.txt"},
		{Title: "Create a secret from standard input (printf '%s' \"$TOKEN\" | ...)", Command: "afs secrets create-secret --team-id " + exTeam + " --name github-pat --scope team --secret-value -"},
		{Title: "Create a secret and type the value at a hidden prompt", Command: "afs secrets create-secret --team-id " + exTeam + " --name github-pat --scope team"},
	}})
	r.Override(cli.Override{Op: "PutSecretValue", Examples: []cli.Example{
		{Title: "Replace a secret's value from a file (prints nothing)", Command: "afs secrets put-secret-value --team-id " + exTeam + " --secret-id " + exID + " --secret-value file://token.txt"},
	}})
	r.Override(cli.Override{Op: "UpdateSecret", Columns: secretColumns, Examples: []cli.Example{
		{Title: "Change a secret's description", Command: `afs secrets update-secret --team-id ` + exTeam + ` --secret-id ` + exID + ` --description "Token for the build bot"`},
	}})
	r.Override(cli.Override{Op: "RotateSecret", Columns: secretColumns, OmitEmptyBody: true, Examples: []cli.Example{
		{Title: "Rotate a secret to a new value", Command: "afs secrets rotate-secret --team-id " + exTeam + " --secret-id " + exID + " --secret-value file://new-token.txt"},
		{Title: "Rotate a secret that AFS generates (no value, so none is sent)", Command: "afs secrets rotate-secret --team-id " + exTeam + " --secret-id " + exID + ""},
	}})
	r.Override(cli.Override{Op: "DeleteSecret", Examples: []cli.Example{
		{Title: "Delete a secret that nothing uses", Command: "afs secrets delete-secret --team-id " + exTeam + " --secret-id " + exID},
	}})
}
