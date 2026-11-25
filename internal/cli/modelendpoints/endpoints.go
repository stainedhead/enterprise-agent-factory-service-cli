package modelendpoints

import (
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/output"
)

func init() { registrations = append(registrations, registerEndpoints) }

// EndpointColumns are the table columns of a model endpoint (service or team-declared).
var EndpointColumns = []output.Column{
	{Header: "Id", Path: "id"}, {Header: "Slug", Path: "slug"}, {Header: "Provider", Path: "provider"},
	{Header: "Location", Path: "location"}, {Header: "Pricing", Path: "pricingMode"}, {Header: "Enabled", Path: "enabled"},
	{Header: "ProviderState", Path: "providerState"}, {Header: "Health", Path: "health.state", Null: "-"},
	{Header: "Credential", Path: "credentialMode"}, {Header: "Teams", Path: "teamsSelected"},
	{Header: "Team", Path: "teamId", Null: "-"}, {Header: "Version", Path: "version"},
}

// ModelColumns are the table columns of an endpoint's models (the items of put-endpoint-models).
var ModelColumns = []output.Column{
	{Header: "Model", Path: "modelName"}, {Header: "Enabled", Path: "enabled"}, {Header: "Callable", Path: "callable"},
	{Header: "Price", Path: "priceState"}, {Header: "InputUsdPerMtok", Path: "price.inputUsdPerMtok", Null: "-"},
	{Header: "OutputUsdPerMtok", Path: "price.outputUsdPerMtok", Null: "-"},
	{Header: "Context", Path: "contextTokens", Null: "-"}, {Header: "MaxOutput", Path: "maxOutputTokens"},
}

var priceColumns = []output.Column{
	{Header: "Model", Path: "modelName"}, {Header: "InputUsdPerMtok", Path: "inputUsdPerMtok"},
	{Header: "OutputUsdPerMtok", Path: "outputUsdPerMtok"}, {Header: "CachedInputUsdPerMtok", Path: "cachedInputUsdPerMtok", Null: "-"},
	{Header: "RequestFeeUsd", Path: "requestFeeUsd", Null: "-"}, {Header: "EffectiveFrom", Path: "effectiveFrom"},
	{Header: "Source", Path: "source"},
}

var proposalColumns = []output.Column{
	{Header: "Id", Path: "id"}, {Header: "Endpoint", Path: "slug"}, {Header: "Model", Path: "modelName"},
	{Header: "InputUsdPerMtok", Path: "inputUsdPerMtok"}, {Header: "OutputUsdPerMtok", Path: "outputUsdPerMtok"},
	{Header: "EffectiveFrom", Path: "effectiveFrom"}, {Header: "Status", Path: "status"}, {Header: "Source", Path: "source"},
	{Header: "DecidedBy", Path: "decidedBy", Null: "-"}, {Header: "Version", Path: "version"},
}

var endpointTeamColumns = []output.Column{
	{Header: "TeamId", Path: "teamId"}, {Header: "Team", Path: "teamSlug", Null: "-"}, {Header: "Selected", Path: "selected"},
	{Header: "SelectedBy", Path: "selectedBy", Null: "-"}, {Header: "Credential", Path: "credentialState"},
	{Header: "Federation", Path: "federationState"}, {Header: "Gateway", Path: "gatewayDeployment"},
	{Header: "GatewayVersion", Path: "gatewayVersion", Null: "-"}, {Header: "Health", Path: "health.state", Null: "-"},
	{Header: "SpendMtdUsd", Path: "spendMtdUsd"},
}

// Desc and Kinds shared by the service and team endpoint writes.
var (
	endpointDesc = map[string]string{
		"caBundle":  "PEM bundle of the CA that signs the endpoint's certificate, as text or file://PATH (not a secret).",
		"protocols": "Protocols the endpoint speaks: openai_chat always; openai_responses and openai_embeddings add those paths (space-separated).",
		"models":    "The endpoint's models as JSON or file://PATH: [{\"modelName\":…, \"maxOutputTokens\":…, \"contextTokens\":…, \"enabled\":…}]; a team endpoint's per_token models also carry \"price\".",
	}
	endpointKinds = map[string]cli.Kind{"caBundle": cli.Any}
)

func registerEndpoints(r *cli.Registry) {
	r.Override(cli.Override{Op: "ListModelEndpoints", Columns: EndpointColumns, Examples: []cli.Example{
		{Title: "List the service's model endpoints", Command: "afs model-endpoints list-model-endpoints"},
		{Title: "Include team-declared endpoints and keep only per-token external routes", Command: "afs model-endpoints list-model-endpoints --include-teams --location external --pricing-mode per_token"},
	}})
	r.Override(cli.Override{Op: "CreateModelEndpoint", Columns: EndpointColumns, Desc: endpointDesc, Kinds: endpointKinds,
		Notes: []string{"The endpoint is created switched off; turn it on with update-model-endpoint --enabled once its models are priced and its credential is set. The base URL must be https and its host must not belong to a provider that is off."},
		Examples: []cli.Example{
			{Title: "Add an OpenAI route with a credential per team", Command: `afs model-endpoints create-model-endpoint --slug openai-prod --display-name "OpenAI production" --provider openai --base-url https://api.openai.com/v1 --protocols openai_chat openai_responses --location external --pricing-mode per_token --auth-kind bearer --credential-mode per_team --models file://openai-models.json`},
			{Title: "Add a corporate capacity endpoint with its own CA", Command: `afs model-endpoints create-model-endpoint --slug corp-llm --display-name "Corporate LLM" --provider internal --base-url https://llm.corp.example.com/v1 --protocols openai_chat --location corporate --pricing-mode capacity --auth-kind header --auth-header X-Api-Key --ca-bundle file://corp-ca.pem`},
		}})
	r.Override(cli.Override{Op: "GetModelEndpoint", Columns: EndpointColumns, Examples: []cli.Example{
		{Title: "Show an endpoint", Command: "afs model-endpoints describe-model-endpoint --endpoint-id " + exEndpoint},
		{Title: "Show its models and prices", Command: "afs model-endpoints describe-model-endpoint --endpoint-id " + exEndpoint + " --query models --output table"},
	}})
	r.Override(cli.Override{Op: "UpdateModelEndpoint", Columns: EndpointColumns, Desc: endpointDesc, Kinds: endpointKinds, Examples: []cli.Example{
		{Title: "Switch an endpoint on", Command: "afs model-endpoints update-model-endpoint --endpoint-id " + exEndpoint + " --enabled"},
		{Title: "Raise the timeout and require extra terms", Command: "afs model-endpoints update-model-endpoint --endpoint-id " + exEndpoint + " --timeout-seconds 600 --requires-terms"},
	}})
	r.Override(cli.Override{Op: "DeleteModelEndpoint",
		Notes: []string{"Refused while the endpoint is enabled (endpoint_enabled): switch it off first."},
		Examples: []cli.Example{
			{Title: "Switch an endpoint off, then delete it", Command: "afs model-endpoints delete-model-endpoint --endpoint-id " + exEndpoint},
		}})
	r.Override(cli.Override{Op: "PutEndpointModels", Columns: ModelColumns, ResultKey: "items", Desc: endpointDesc,
		Notes: []string{"Replaces the endpoint's model list. A model is callable once it is enabled and, on a per_token endpoint, priced."},
		Examples: []cli.Example{
			{Title: "Replace the endpoint's models from a file", Command: "afs model-endpoints put-endpoint-models --endpoint-id " + exEndpoint + " --models file://openai-models.json"},
		}})
	r.Override(cli.Override{Op: "CreateEndpointModelPrice", Columns: priceColumns, Examples: []cli.Example{
		{Title: "Price a model from the first of next month", Command: "afs model-endpoints create-endpoint-model-price --endpoint-id " + exEndpoint + " --model-name gpt-4.1 --input-usd-per-mtok 2.00 --output-usd-per-mtok 8.00 --cached-input-usd-per-mtok 0.50 --effective-from 2026-12-01T00:00:00Z"},
		{Title: "Record an internal platform's chargeback rate", Command: "afs model-endpoints create-endpoint-model-price --endpoint-id " + exEndpoint + " --model-name gpt-4.1 --input-usd-per-mtok 1.80 --output-usd-per-mtok 7.20 --effective-from 2026-12-01T00:00:00Z --source transfer"},
	}})
	r.Override(cli.Override{Op: "ListEndpointPriceProposals", Columns: proposalColumns, Examples: []cli.Example{
		{Title: "List pending endpoint price proposals", Command: "afs model-endpoints list-endpoint-price-proposals --status pending"},
		{Title: "List one endpoint's proposals", Command: "afs model-endpoints list-endpoint-price-proposals --endpoint-id " + exEndpoint},
	}})
	r.Override(cli.Override{Op: "ApproveEndpointPriceProposal", Columns: proposalColumns, Examples: []cli.Example{
		{Title: "Approve a proposal", Command: "afs model-endpoints approve-endpoint-price-proposal --proposal-id " + exProposal + ` --note "Matches the provider's page"`},
	}})
	r.Override(cli.Override{Op: "RejectEndpointPriceProposal", Columns: proposalColumns, Examples: []cli.Example{
		{Title: "Reject a proposal", Command: "afs model-endpoints reject-endpoint-price-proposal --proposal-id " + exProposal + ` --note "Waiting for the contract rate"`},
	}})
	r.Override(cli.Override{Op: "ListEndpointTeams", Columns: endpointTeamColumns, Examples: []cli.Example{
		{Title: "Show which teams use an endpoint, with credential and federation state and spend", Command: "afs model-endpoints list-endpoint-teams --endpoint-id " + exEndpoint},
		{Title: "List the teams whose credential is missing", Command: "afs model-endpoints list-endpoint-teams --endpoint-id " + exEndpoint + ` --query "items[?credentialState=='missing'].teamId" --output text`},
	}})
}
