package modelendpoints

import (
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/output"
)

func init() { registrations = append(registrations, registerCredentials) }

// CredentialColumns are the table columns of a credential write: its state and last rotation, never a value.
var CredentialColumns = []output.Column{
	{Header: "Endpoint", Path: "endpointId"}, {Header: "Scope", Path: "scope"}, {Header: "State", Path: "state"},
	{Header: "RotatedAt", Path: "rotatedAt", Null: "-"},
}

var federationColumns = []output.Column{
	{Header: "Team", Path: "teamId"}, {Header: "Kind", Path: "kind"}, {Header: "Status", Path: "status"},
	{Header: "ExpectedIssuer", Path: "expectedIssuer", Null: "-"}, {Header: "ExpectedSubject", Path: "expectedSubject", Null: "-"},
	{Header: "ExpectedAudience", Path: "expectedAudience", Null: "-"}, {Header: "ExpectedRoleArn", Path: "expectedRoleArn", Null: "-"},
	{Header: "Version", Path: "version", Null: "-"},
}

// credentialNote is said on every credential write (MEG-CLI-3).
const credentialNote = //nolint:gosec // G101: help text about credentials, not a credential
"The value is read from file://PATH, fileb://PATH or - (standard input), or asked for without echo when it is omitted on a terminal; it is never a plain argument, and the API never returns it (only its state and last rotation). The credential's version is not returned, so no If-Match is sent unless --if-match is given."

func registerCredentials(r *cli.Registry) {
	r.Override(cli.Override{Op: "PutModelEndpointCredential", Columns: CredentialColumns,
		Notes: []string{"--scope is a team id for a per_team endpoint, or _shared for a shared one.", credentialNote,
			"To replace a value that is already set without knowing its version, use rotate-model-endpoint-credential."},
		Examples: []cli.Example{
			{Title: "Set a team's key on a per-team endpoint from a file", Command: "afs model-endpoints put-model-endpoint-credential --endpoint-id " + exEndpoint + " --scope " + exTeam + " --secret-value file://openai-payments.key"},
			{Title: "Set the shared key from standard input", Command: "afs model-endpoints put-model-endpoint-credential --endpoint-id " + exEndpoint + " --scope _shared --secret-value -"},
		}})
	r.Override(cli.Override{Op: "RotateModelEndpointCredential", Columns: CredentialColumns,
		Notes: []string{"Replaces the value that is set (NotFound when none is); the previous value is no longer used.", "The value is read from file://PATH, fileb://PATH or - (standard input), or asked for without echo when it is omitted on a terminal; it is never a plain argument."},
		Examples: []cli.Example{
			{Title: "Rotate a team's key, asked for without echo", Command: "afs model-endpoints rotate-model-endpoint-credential --endpoint-id " + exEndpoint + " --scope " + exTeam},
			{Title: "Rotate the shared key from a file", Command: "afs model-endpoints rotate-model-endpoint-credential --endpoint-id " + exEndpoint + " --scope _shared --secret-value file://new.key"},
		}})
	r.Override(cli.Override{Op: "PutEndpointFederation", Columns: federationColumns,
		Notes: []string{"Workload identity federation for an azure_federated or google_federated endpoint: metadata only, no secret is stored. The response names the issuer, subject and audience (or role) to trust on the provider's side."},
		Examples: []cli.Example{
			{Title: "Federate a team's gateway with an Azure app registration", Command: "afs model-endpoints put-endpoint-federation --endpoint-id " + exEndpoint + " --team-id " + exTeam + " --tenant-id 11111111-2222-3333-4444-555555555555 --client-id 66666666-7777-8888-9999-000000000000"},
			{Title: "Federate with a Google workload identity pool", Command: "afs model-endpoints put-endpoint-federation --endpoint-id " + exEndpoint + " --team-id " + exTeam + " --project-number 123456789012 --pool-id afs-gateways --provider-id eks --project-id payments-ai --location us-central1"},
		}})
}
