package modelendpoints

import (
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/output"
)

func init() { registrations = append(registrations, registerTeam) }

// TeamEndpointColumns are the table columns of the team's view of the endpoints it may use (screen 61).
var TeamEndpointColumns = []output.Column{
	{Header: "Id", Path: "id"}, {Header: "Slug", Path: "slug"}, {Header: "Relation", Path: "relation"},
	{Header: "Provider", Path: "provider"}, {Header: "Pricing", Path: "pricingMode"}, {Header: "Enabled", Path: "enabled"},
	{Header: "Credential", Path: "credentialState"}, {Header: "Health", Path: "health.state", Null: "-"},
	{Header: "SelectedBy", Path: "selectedBy", Null: "-"}, {Header: "OffReason", Path: "offReason", Null: "-"},
	{Header: "Version", Path: "version"},
}

func registerTeam(r *cli.Registry) {
	r.Override(cli.Override{Op: "ListTeamModelEndpoints", Columns: TeamEndpointColumns, Examples: []cli.Example{
		{Title: "List the routes the team uses and may use", Command: "afs model-endpoints list-team-model-endpoints --team-id " + exTeam},
		{Title: "List the service endpoints the team may opt in to", Command: "afs model-endpoints list-team-model-endpoints --team-id " + exTeam + " --relation available"},
	}})
	r.Override(cli.Override{Op: "CreateTeamModelEndpoint", Columns: EndpointColumns, Desc: endpointDesc, Kinds: endpointKinds,
		Notes: []string{"Needs the team's mayAddModelEndpoints (team_endpoints_not_allowed otherwise; a service admin sets it with afs teams put-team-gateway-settings). The endpoint is created switched off."},
		Examples: []cli.Example{
			{Title: "Declare a vLLM server in the team's namespace", Command: `afs model-endpoints create-team-model-endpoint --team-id ` + exTeam + ` --slug coder-vllm --display-name "Coder vLLM" --base-url http://vllm.afs-t-payments.svc:8000/v1 --protocols openai_chat --location in_namespace --pricing-mode capacity --auth-kind none --owner payments-platform --purpose "Code completion for the payments agents" --models file://coder-models.json`},
		}})
	r.Override(cli.Override{Op: "UpdateTeamModelEndpoint", Columns: EndpointColumns, Desc: endpointDesc, Kinds: endpointKinds, Examples: []cli.Example{
		{Title: "Switch a team endpoint on", Command: "afs model-endpoints update-team-model-endpoint --team-id " + exTeam + " --endpoint-id " + exTeamEndpoint + " --enabled"},
		{Title: "Scrape the server's metrics port", Command: "afs model-endpoints update-team-model-endpoint --team-id " + exTeam + " --endpoint-id " + exTeamEndpoint + " --metrics-port 8000"},
	}})
	r.Override(cli.Override{Op: "DeleteTeamModelEndpoint", Examples: []cli.Example{
		{Title: "Delete a team endpoint", Command: "afs model-endpoints delete-team-model-endpoint --team-id " + exTeam + " --endpoint-id " + exTeamEndpoint},
	}})
	r.Override(cli.Override{Op: "PutTeamEndpointModels", Columns: ModelColumns, ResultKey: "items", Desc: endpointDesc,
		Notes: []string{"Replaces the team endpoint's models. On a per_token endpoint each model carries its price inline; capacity and free endpoints take none."},
		Examples: []cli.Example{
			{Title: "Replace a team endpoint's models from a file", Command: "afs model-endpoints put-team-endpoint-models --team-id " + exTeam + " --endpoint-id " + exTeamEndpoint + " --models file://coder-models.json"},
		}})
	r.Override(cli.Override{Op: "PutTeamModelEndpointCredential", Columns: CredentialColumns,
		Notes: []string{"Sets the team's credential of a team endpoint, or of a per_team service endpoint the team opted in to.", credentialNote},
		Examples: []cli.Example{
			{Title: "Set the team's key from a file", Command: "afs model-endpoints put-team-model-endpoint-credential --team-id " + exTeam + " --endpoint-id " + exEndpoint + " --secret-value file://openai-payments.key"},
			{Title: "Set it from standard input (piped from a secret store)", Command: "afs model-endpoints put-team-model-endpoint-credential --team-id " + exTeam + " --endpoint-id " + exEndpoint + " --secret-value -"},
		}})
}
