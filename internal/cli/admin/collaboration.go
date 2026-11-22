package admin

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/output"
)

func init() { registrations = append(registrations, registerCollaboration) }

// endpointColumns are the columns of a team's collaboration endpoint (CLI-US-30.2). The surfaces are the
// effective ones; a team that never signed in has no last access, shown as "-".
var endpointColumns = []output.Column{
	{Header: "Team", Path: "slug"}, {Header: "TeamState", Path: "teamState"}, {Header: "Hostname", Path: "hostname", Null: "-"},
	{Header: "WebUi", Path: "effective.webUi"}, {Header: "RestApi", Path: "effective.restApi"},
	{Header: "Phase", Path: "phase"}, {Header: "LastAccessAt", Path: "lastAccessAt", Null: "-"},
}

// secretsManagerARN is the shape of a Secrets Manager secret ARN (any AWS partition). The check only keeps
// a secret value out of an option that names a reference (CLI-US-31.2); the API validates the rest.
var secretsManagerARN = regexp.MustCompile(`^arn:aws[a-z-]*:secretsmanager:[a-z0-9-]+:[0-9]{12}:secret:[A-Za-z0-9/_+=.@-]+$`)

// buildCollabSecretARN refuses a --collab-client-secret-arn (or the member of --cli-input-json) that is not
// an ARN before any call; the message never repeats the value.
func buildCollabSecretARN(_ *cli.Parsed, parts *cli.RequestParts) error {
	if parts.Body == nil {
		return nil
	}
	raw, ok := parts.Body.Get("collabClientSecretArn")
	if !ok {
		return nil
	}
	var s string
	if json.Unmarshal(raw, &s) != nil || !secretsManagerARN.MatchString(s) {
		return fmt.Errorf("--collab-client-secret-arn takes the ARN of a Secrets Manager secret (arn:aws:secretsmanager:<region>:<account>:secret:<name>), never the secret value")
	}
	return nil
}

func registerCollaboration(r *cli.Registry) {
	r.Override(cli.Override{Op: "ListCollaborationEndpoints", Columns: endpointColumns, Examples: []cli.Example{
		{Title: "List every team's Paperclip endpoint", Command: "afs admin list-collaboration-endpoints"},
		{Title: "List the endpoints that are not ready, as CSV", Command: "afs admin list-collaboration-endpoints --state Failed --output csv"},
		{Title: "Show one team's endpoint", Command: "afs admin list-collaboration-endpoints --team-id team_01J9ZQ3ABCDEFGHJKMNPQRSTVW"},
	}})
	collab := []cli.Example{
		{Title: "Register the collaboration sign-in client (a reference to the secret, never its value)", Command: "afs admin update-identity-provider --provider-id " + idp + " --collab-client-id paperclip-web --collab-client-secret-arn arn:aws:secretsmanager:eu-west-1:111122223333:secret:afs/collab-AbCdEf"},
	}
	r.Override(cli.Override{Op: "CreateIdentityProvider", Build: buildCollabSecretARN})
	r.Override(cli.Override{Op: "UpdateIdentityProvider", Build: buildCollabSecretARN, Examples: collab})
	r.Override(cli.Override{Op: "PutSetting", Examples: []cli.Example{
		{Title: "Set the base domain of the Paperclip endpoints", Command: `afs admin put-setting --key collab.base_domain --value '"collab.example.com"'`},
		{Title: "Switch the public Paperclip endpoints on", Command: "afs admin put-setting --key collab.public_endpoints --value true"},
	}})
}
