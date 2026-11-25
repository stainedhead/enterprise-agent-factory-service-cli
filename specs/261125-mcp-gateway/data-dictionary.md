# Data dictionary

| Wire type | Columns |
|---|---|
| `McpServer` | Id, Slug, Scope, Team, Location, Transport, Auth, Review, Risk, Published, Credential, Version |
| `TeamMcpServer` | Id, Slug (from `server`), Managed, State, Mode, Review, Credential, Federation |
| `McpCredential` | Scope, Kind, State, Subject, ExpiresAt, RotatedAt (never a value) |
| `McpPolicy` | Server, Level, Team, DriftPolicy, TrustAnnotations, AllowSampling, Version |
| `McpPrimitive` | Id, Kind, Name, Status, Change, LastSeen |
| `McpReviewItem` | Kind, Server, Slug, Team, Primitive, Risk, Approver, Since |
| `EffectiveMcp.servers` | Server, Slug, Scope, Mode, Url |
| `FederatedIdentity` | Name, Audience, ExpirationSeconds, FilePath, EnvVar, Subject, Version |
| `McpEvent` | Id, At, Team, Agent, Server, Method, Name, Outcome, Reason, GatewayMs |
| `McpUsageReport.rows` | Key, Label, Calls, Errors, Refusals, P50Ms, P95Ms, Teams, Agents, Searches, LastUsedAt |
