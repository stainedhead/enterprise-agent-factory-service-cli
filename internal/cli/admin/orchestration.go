package admin

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli/runtime"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/client"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/output"
)

func init() { registrations = append(registrations, registerOrchestration) }

// orchestrationEndpointColumns are the columns of one team surface whose stored exposure is not off (ORC-25,
// OR-CLI-26). Effective can differ from Exposure; the reasons are in the JSON and CSV.
var orchestrationEndpointColumns = []output.Column{
	{Header: "Team", Path: "slug"}, {Header: "TeamState", Path: "teamState"}, {Header: "Tool", Path: "tool"},
	{Header: "Surface", Path: "surface"}, {Header: "Exposure", Path: "exposure"}, {Header: "Effective", Path: "effective"},
	{Header: "Limit", Path: "limit"}, {Header: "Hostname", Path: "hostname", Null: "-"},
	{Header: "Phase", Path: "phase"}, {Header: "LastAccessAt", Path: "lastAccessAt", Null: "-"},
}

// serviceModelAccessColumns are the columns of the service-wide model access (MOD-1).
var serviceModelAccessColumns = []output.Column{
	{Header: "Access", Path: "access"}, {Header: "EnabledModels", Path: "enabledModels"},
	{Header: "CallableModels", Path: "callableModels"}, {Header: "Version", Path: "version"},
}

// afterEndpoints prints, for table and text output, the progress of the 1.1.0 backfill when it is not complete
// (S12): how many teams have an orchestration agent, and the teams that are blocked with the reason. The values
// are printed as returned, sanitized. With json and csv nothing extra is printed.
func afterEndpoints(inv *cli.Invocation, resp *client.Response) int {
	st, err := runtime.Settings(inv.Env, inv.P.G)
	if err != nil || (st.Output != "table" && st.Output != "text") {
		return -1
	}
	var doc struct {
		Migration struct {
			State                  string `json:"state"`
			TeamsTotal             int    `json:"teamsTotal"`
			TeamsWithOrchestration int    `json:"teamsWithOrchestration"`
			TeamsPending           int    `json:"teamsPending"`
			Blocked                []struct {
				TeamID string `json:"teamId"`
				Reason string `json:"reason"`
			} `json:"blocked"`
		} `json:"migration"`
	}
	if json.Unmarshal(resp.Body, &doc) != nil || doc.Migration.State == "" || doc.Migration.State == "complete" {
		return -1
	}
	m := doc.Migration
	_, _ = fmt.Fprintf(inv.Env.Stderr, "Orchestration agents: %s (%d of %d teams have one, %d pending).\n", output.Sanitize(m.State), m.TeamsWithOrchestration, m.TeamsTotal, m.TeamsPending)
	for _, b := range m.Blocked {
		_, _ = fmt.Fprintf(inv.Env.Stderr, "Blocked: %s (%s)\n", output.Sanitize(b.TeamID), output.Sanitize(strings.TrimSpace(b.Reason)))
	}
	return -1
}

func registerOrchestration(r *cli.Registry) {
	r.Override(cli.Override{Op: "ListOrchestrationEndpoints", Columns: orchestrationEndpointColumns, After: afterEndpoints, Examples: []cli.Example{
		{Title: "List every team surface that is exposed, privately or publicly", Command: "afs admin list-orchestration-endpoints"},
		{Title: "List the Buzz teams whose orchestration agent is not ready", Command: "afs admin list-orchestration-endpoints --tool buzz --state Failed"},
		{Title: "List the effectively public surfaces as CSV", Command: "afs admin list-orchestration-endpoints --exposure public --output csv"},
	}})
	r.Override(cli.Override{Op: "PutServiceModelAccess", Columns: serviceModelAccessColumns,
		Desc: map[string]string{"access": "all (every catalog model is enabled for the service, including models added later; each model's enabled flag is ignored) or selected (each model's enabled flag decides)."},
		Examples: []cli.Example{
			{Title: "Enable every catalog model for the service", Command: "afs admin put-service-model-access --access all"},
			{Title: "Go back to choosing models one by one", Command: "afs admin put-service-model-access --access selected"},
		}})
	r.Override(cli.Override{Op: "ListCollaborationEndpoints", Deprecated: "afs admin list-orchestration-endpoints"})
	r.Override(cli.Override{Op: "PutSetting", OpenEnum: []string{"key"},
		Notes: []string{"collab.base_domain and collab.public_endpoints are read-only until 2.0.0; use orchestration.public_base_domain and orchestration.public_endpoints."},
		Examples: []cli.Example{
			{Title: "Set the base domain of the public orchestration hostnames", Command: `afs admin put-setting --key orchestration.public_base_domain --value '"collab.example.com"'`},
			{Title: "Set the base domain of the private (internal ALB) hostnames", Command: `afs admin put-setting --key orchestration.private_base_domain --value '"collab.corp.example.com"'`},
			{Title: "Lower the default limit of new teams to private", Command: `afs admin put-setting --key orchestration.default_limit --value '"private"'`},
			{Title: "Choose the default model of new agents", Command: `afs admin put-setting --key models.default_model_id --value '"anthropic.claude-sonnet-4-5"'`},
		}})
}
