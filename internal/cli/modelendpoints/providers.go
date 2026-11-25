package modelendpoints

import (
	"encoding/json"
	"fmt"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli/runtime"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/client"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/output"
)

func init() { registrations = append(registrations, registerProviders) }

// ProviderColumns are the table columns of an inference provider (screen 68).
var ProviderColumns = []output.Column{
	{Header: "Provider", Path: "provider"}, {Header: "Name", Path: "displayName"}, {Header: "Enabled", Path: "enabled"},
	{Header: "TeamSelection", Path: "teamSelection"}, {Header: "Routes", Path: "routes"}, {Header: "Teams", Path: "teams"},
	{Header: "DirectSystems", Path: "directSystems"}, {Header: "Fallback", Path: "fallbackModel", Null: "-"},
	{Header: "Version", Path: "version"},
}

func registerProviders(r *cli.Registry) {
	r.Override(cli.Override{Op: "ListInferenceProviders", Columns: ProviderColumns, ResultKey: "items", Examples: []cli.Example{
		{Title: "List the providers with their hosts, routes and teams", Command: "afs model-endpoints list-inference-providers"},
		{Title: "Print the hosts of each provider", Command: "afs model-endpoints list-inference-providers --query 'items[].[provider,join(`,`,hosts)]' --output text"},
	}})
	r.Override(cli.Override{Op: "UpdateInferenceProvider", Columns: ProviderColumns, After: afterProvider,
		Notes: []string{
			"Turning a provider off stops every route on its hosts within a minute, service and team endpoints alike; turning it back on restores opt-ins, selections, credentials and prices.",
			"Turning Bedrock off needs --fallback-model <route>/<model> on a provider that is on (no_fallback_model otherwise); --select-for-all-teams selects that route for every team.",
			"--dry-run sends dryRun=true: nothing is written and the response's impact lists the routes, teams and agents affected; with table or text output the impact is also printed on stderr.",
		},
		Examples: []cli.Example{
			{Title: "See what turning Bedrock off would change, without writing", Command: `afs model-endpoints update-inference-provider --provider bedrock --no-enabled --fallback-model azure-openai/gpt-5.1 --select-for-all-teams --reason "No AWS Bedrock in this enterprise" --dry-run`},
			{Title: "Turn Bedrock off and move agents to an Azure route", Command: `afs model-endpoints update-inference-provider --provider bedrock --no-enabled --fallback-model azure-openai/gpt-5.1 --select-for-all-teams --reason "No AWS Bedrock in this enterprise"`},
			{Title: "Turn Azure on, selected per team by service admins", Command: "afs model-endpoints update-inference-provider --provider azure --enabled --team-selection service_admin"},
		}})
}

// afterProvider prints, for table and text output of a dry run, that nothing was written and the impact the
// API computed. JSON output carries the impact in the response and prints nothing extra.
func afterProvider(inv *cli.Invocation, resp *client.Response) int {
	if inv.P.Value("dry-run") != "true" {
		return -1
	}
	st, err := runtime.Settings(inv.Env, inv.P.G)
	if err != nil || (st.Output != "table" && st.Output != "text") {
		return -1
	}
	var doc struct {
		Impact *struct {
			Routes                 int     `json:"routes"`
			Teams                  int     `json:"teams"`
			Agents                 int     `json:"agents"`
			AgentsMovingToFallback int     `json:"agentsMovingToFallback"`
			AgentsDegraded         int     `json:"agentsDegraded"`
			DirectSystems          int     `json:"directSystems"`
			TeamsToSelect          int     `json:"teamsToSelect"`
			FallbackModel          *string `json:"fallbackModel"`
		} `json:"impact"`
	}
	_, _ = fmt.Fprintln(inv.Env.Stderr, "Dry run: nothing was changed.")
	if json.Unmarshal(resp.Body, &doc) != nil || doc.Impact == nil {
		return -1
	}
	im := doc.Impact
	_, _ = fmt.Fprintf(inv.Env.Stderr, "Impact: %d routes, %d teams, %d agents (%d moving to the fallback, %d degraded), %d direct systems, %d teams to select.\n",
		im.Routes, im.Teams, im.Agents, im.AgentsMovingToFallback, im.AgentsDegraded, im.DirectSystems, im.TeamsToSelect)
	if im.FallbackModel != nil {
		_, _ = fmt.Fprintf(inv.Env.Stderr, "Fallback model: %s\n", output.Sanitize(*im.FallbackModel))
	}
	return -1
}
