# Data dictionary: model-endpoint-gateway

Wire shapes are the snapshot's; the CLI prints JSON unchanged. Table columns:

| Resource | Columns |
|---|---|
| `ModelEndpoint` | Id, Slug, Provider, Location, Pricing (`pricingMode`), Enabled, Provider (`providerState`), Health (`health.state`), Credential (`credentialMode`), Teams (`teamsSelected`), Team (`teamId`, `-`), Version |
| `TeamModelEndpoint` | Id, Slug, Relation, Provider, Pricing, Enabled, Credential (`credentialState`), Health, SelectedBy, OffReason, Version |
| `EndpointModel` | Model, Enabled, Callable, Price (`priceState`), Input, Output, Context |
| `EndpointModelPrice` | Model, Input, Output, CachedInput, RequestFee, EffectiveFrom, Source |
| `EndpointPriceProposal` | Id, Endpoint (`slug`), Model, Input, Output, EffectiveFrom, Status, Source, Version |
| `EndpointTeam` | TeamId, Slug, Selected, SelectedBy, Credential, Federation, Gateway, GatewayVersion, Health, SpendMtd |
| `ModelEndpointCredential` | Endpoint, Scope, State, RotatedAt (never a value) |
| `EndpointFederation` | Team, Kind, Status, ExpectedIssuer, ExpectedSubject, ExpectedAudience, ExpectedRoleArn, Version |
| `ProviderStatement` | Id, Endpoint, Period, Source, Rows, TotalUsd, ImportedAt |
| `InferenceProvider` | Provider, Name, Enabled, TeamSelection, Routes, Teams, DirectSystems, Fallback, Version |
| `TeamGatewaySettings` | Team, Deployment, MayAddEndpoints, State, Url, Version |
| `TeamModelEndpointSelections` | EndpointId, Slug, SelectedBy (rows of `selections`) |
| `InferenceReport` | dynamic (rows follow `--group-by`) |

Local options:

| Command | Option | Meaning |
|---|---|---|
| `import-provider-statement` | `--statement file://PATH` | CSV (sent as `text/csv`) or JSON (`{period, source, rows}` or a rows array) |
| agents and pools model commands | `--routes all\|<slug>…`, `--default-route <slug>` | `models.routeAccess`, `models.routes`, `models.defaultRoute` |
