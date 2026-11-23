package teams

import (
	"encoding/json"
	"fmt"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli/runtime"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/client"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/output"
)

func init() { registrations = append(registrations, registerModels) }

// modelColumns are the table columns of a team's model (MOD-7): whether the service enabled it, whether the
// gateway can call it (enabled, priced and allowed), its price state and the team's inference profile.
var modelColumns = []output.Column{
	{Header: "Id", Path: "id"}, {Header: "Provider", Path: "provider"}, {Header: "Name", Path: "displayName"},
	{Header: "Enabled", Path: "enabled"}, {Header: "Callable", Path: "callable", Null: "-"}, {Header: "Price", Path: "priceState", Null: "-"},
	{Header: "Profile", Path: "inferenceProfileState", Null: "-"},
	{Header: "Context", Path: "contextWindow"}, {Header: "MaxOutput", Path: "maxOutputTokens"},
}

// afterModels prints, for table and text output, the team's access mode and, for a selection, how many models
// it names, so the rows are read in context (MOD-8). With json nothing extra is printed.
func afterModels(inv *cli.Invocation, resp *client.Response) int {
	st, err := runtime.Settings(inv.Env, inv.P.G)
	if err != nil || (st.Output != "table" && st.Output != "text") {
		return -1
	}
	var doc struct {
		Access   string   `json:"access"`
		ModelIDs []string `json:"modelIds"`
	}
	if json.Unmarshal(resp.Body, &doc) != nil || doc.Access == "" {
		return -1
	}
	switch doc.Access {
	case "selected":
		_, _ = fmt.Fprintf(inv.Env.Stderr, "Team model access: selected (%d in the team's list).\n", len(doc.ModelIDs))
	default:
		_, _ = fmt.Fprintf(inv.Env.Stderr, "Team model access: %s (every model enabled for the service).\n", output.Sanitize(doc.Access))
	}
	return -1
}

func registerModels(r *cli.Registry) {
	r.Override(cli.Override{Op: "ListTeamModels", Columns: modelColumns, After: afterModels, Examples: []cli.Example{
		{Title: "List the models a team may use", Command: "afs teams list-team-models --team-id team_01J9ZQ3ABCDEFGHJKMNPQRSTVW"},
		{Title: "List the ids of the models the gateway can call for the team", Command: "afs teams list-team-models --team-id team_01J9ZQ3ABCDEFGHJKMNPQRSTVW --query 'items[?callable].id' --output text"},
	}})
	r.Override(cli.Override{Op: "PutTeamModels", Columns: modelColumns, After: afterModels,
		Desc: map[string]string{
			"access":   "all (every model enabled for the service, now and later) or selected (the list in --model-ids). A body with only --model-ids is read as selected (deprecated form).",
			"modelIds": "The team's list: 0 to 100 distinct model ids, each enabled for the service. Not accepted with --access all.",
		},
		Notes: []string{"Only service admins change a team's models (teams:set-guardrails). Narrowing never edits an agent's stored selection: restoring a model restores access (MOD-6)."},
		Examples: []cli.Example{
			{Title: "Let the team use every model enabled for the service", Command: "afs teams put-team-models --team-id team_01J9ZQ3ABCDEFGHJKMNPQRSTVW --access all"},
			{Title: "Limit the team to two models", Command: "afs teams put-team-models --team-id team_01J9ZQ3ABCDEFGHJKMNPQRSTVW --access selected --model-ids anthropic.claude-sonnet-4-5 anthropic.claude-haiku-4-5"},
		}})
}
