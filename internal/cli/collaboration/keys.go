package collaboration

import (
	"encoding/json"
	"fmt"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/client"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/output"
)

func init() { registrations = append(registrations, registerKeys) }

// rotatedColumns list what was rotated: names and Secrets Manager ids, never a value (CLI-US-27.1).
var rotatedColumns = []output.Column{
	{Header: "Name", Path: "name"}, {Header: "SecretId", Path: "secretId"}, {Header: "VersionId", Path: "versionId"},
}

func registerKeys(r *cli.Registry) {
	r.Override(cli.Override{Op: "RotateCollaborationKeys", Columns: rotatedColumns, ResultKey: "rotated", OmitEmptyBody: true,
		After: afterRotate, Examples: []cli.Example{
			{Title: "Rotate both of Paperclip's signing secrets (this signs everyone out of Paperclip)", Command: "afs collaboration rotate-collaboration-keys --team-id " + exTeam},
			{Title: "Rotate one secret", Command: "afs collaboration rotate-collaboration-keys --team-id " + exTeam + " --names paperclip-better-auth-secret"},
		}})
}

// afterRotate prints what the rotation does, on stderr (CLI-US-27.2). It never asks first (PC-CLI-D8).
func afterRotate(inv *cli.Invocation, resp *client.Response) int {
	var doc struct {
		SignsOutEveryone    bool `json:"signsOutEveryone"`
		AgentRolloutPending bool `json:"agentRolloutPending"`
	}
	if json.Unmarshal(resp.Body, &doc) != nil {
		return -1
	}
	if doc.SignsOutEveryone {
		_, _ = fmt.Fprintln(inv.Env.Stderr, "Everyone is signed out of Paperclip.")
	}
	if doc.AgentRolloutPending {
		_, _ = fmt.Fprintln(inv.Env.Stderr, "Agents will roll to pick up the new keys.")
	}
	return -1
}
