package apikeys

import (
	"encoding/json"
	"fmt"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/client"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/output"
)

func init() { registrations = append(registrations, registerKeys) }

// KeyColumns are the table columns of an ApiKey.
var KeyColumns = []output.Column{
	{Header: "Id", Path: "id"}, {Header: "Name", Path: "name"}, {Header: "Role", Path: "role"},
	{Header: "KeyPrefix", Path: "keyPrefix"}, {Header: "ExpiresAt", Path: "expiresAt"}, {Header: "LastUsedAt", Path: "lastUsedAt"},
}

// CreatedKeyColumns add the one-time secret (D-C25).
var CreatedKeyColumns = append(append([]output.Column{}, KeyColumns...), output.Column{Header: "SecretKey", Path: "secretKey"})

// ShowSecretNotice is the After hook of the key-creating operations (D-C25): the secret is in the
// normal output, once; a replay carries none and says so.
func ShowSecretNotice(inv *cli.Invocation, resp *client.Response) int {
	var body struct {
		SecretKey string `json:"secretKey"`
	}
	_ = json.Unmarshal(resp.Body, &body)
	switch {
	case body.SecretKey != "":
		_, _ = fmt.Fprintln(inv.Env.Stderr, "Store this key now; it cannot be shown again.")
	case resp.Replayed:
		_, _ = fmt.Fprintln(inv.Env.Stderr, "This is a replay; the key's secret was shown only to the first request. Delete the key and create another if it was lost.")
	}
	return -1
}

func registerKeys(r *cli.Registry) {
	r.Override(cli.Override{Op: "ListApiKeys", Columns: KeyColumns, Examples: []cli.Example{
		{Title: "List a team's API keys", Command: "afs api-keys list-api-keys --team-id team_01J9ZQ3ABCDEFGHJKMNPQRSTVW"},
	}})
	r.Override(cli.Override{Op: "CreateApiKey", Columns: CreatedKeyColumns, After: ShowSecretNotice, Examples: []cli.Example{
		{Title: "Create a team API key (the secret is shown once)", Command: "afs api-keys create-api-key --team-id team_01J9ZQ3ABCDEFGHJKMNPQRSTVW --name ci-pipeline --role team_admin --expires-at 2027-02-19T18:00:00Z"},
	}})
	r.Override(cli.Override{Op: "DeleteApiKey", Examples: []cli.Example{
		{Title: "Delete an API key", Command: "afs api-keys delete-api-key --team-id team_01J9ZQ3ABCDEFGHJKMNPQRSTVW --key-id key_01J9ZQ3ABCDEFGHJKMNPQRSTVW"},
	}})
}
