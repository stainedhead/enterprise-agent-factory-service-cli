package mcp

import (
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/output"
)

func init() { registrations = append(registrations, registerAccess) }

var (
	effectiveColumns = []output.Column{
		{Header: "Server", Path: "serverId"}, {Header: "Slug", Path: "slug"}, {Header: "Scope", Path: "scope"},
		{Header: "Mode", Path: "mode"}, {Header: "Url", Path: "url"},
	}
	identityColumns = []output.Column{
		{Header: "Name", Path: "name"}, {Header: "Audience", Path: "audience"}, {Header: "ExpirationSeconds", Path: "expirationSeconds"},
		{Header: "FilePath", Path: "filePath"}, {Header: "EnvVar", Path: "envVar"}, {Header: "Subject", Path: "subject"},
		{Header: "Version", Path: "version"},
	}
	eventColumns = []output.Column{
		{Header: "Id", Path: "id"}, {Header: "At", Path: "at"}, {Header: "Team", Path: "teamId"}, {Header: "Agent", Path: "agentId", Null: "-"},
		{Header: "Server", Path: "serverId", Null: "-"}, {Header: "Method", Path: "method"}, {Header: "Name", Path: "name", Null: "-"},
		{Header: "Outcome", Path: "outcome"}, {Header: "Reason", Path: "reason", Null: "-"}, {Header: "GatewayMs", Path: "gatewayMs", Null: "-"},
	}
	usageColumns = []output.Column{
		{Header: "Key", Path: "key"}, {Header: "Label", Path: "label"}, {Header: "Calls", Path: "calls"}, {Header: "Errors", Path: "errors"},
		{Header: "Refusals", Path: "refusals"}, {Header: "P50Ms", Path: "p50Ms", Null: "-"}, {Header: "P95Ms", Path: "p95Ms", Null: "-"},
		{Header: "Teams", Path: "teams", Null: "-"}, {Header: "Agents", Path: "agents", Null: "-"},
		{Header: "Searches", Path: "searches", Null: "-"}, {Header: "LastUsedAt", Path: "lastUsedAt", Null: "-"},
	}
	identityNote = "A federated identity projects a Kubernetes service account token with the audience into the pod (filePath, envVar), for an external system that trusts the cluster's issuer. The response names the issuer and subject to trust; no secret is stored."
)

func registerAccess(r *cli.Registry) {
	for _, k := range []struct{ svc, kind, flag, id string }{{"agents", "Agent", "--agent-id", exAgent}, {"pools", "Pool", "--pool-id", exPool}} {
		who := " --team-id " + exTeam + " " + k.flag + " " + k.id
		r.Override(cli.Override{Op: "GetEffective" + k.kind + "Mcp", Columns: effectiveColumns, ResultKey: "servers",
			Notes: []string{"The servers the " + k.svc[:len(k.svc)-1] + " reaches through the team gateway after the service, team and " + k.svc[:len(k.svc)-1] + " narrowing, and the AFS_MCP_* variable names it gets."},
			Examples: []cli.Example{
				{Title: "Show the MCP servers it reaches", Command: "afs " + k.svc + " get-effective-" + k.svc[:len(k.svc)-1] + "-mcp" + who},
				{Title: "Show the variable names it gets", Command: "afs " + k.svc + " get-effective-" + k.svc[:len(k.svc)-1] + "-mcp" + who + " --query env --output text"},
			}})
		r.Override(cli.Override{Op: "List" + k.kind + "FederatedIdentities", Columns: identityColumns, Examples: []cli.Example{
			{Title: "List its federated identities", Command: "afs " + k.svc + " list-" + k.svc[:len(k.svc)-1] + "-federated-identities" + who},
		}})
		r.Override(cli.Override{Op: "Put" + k.kind + "FederatedIdentity", Columns: identityColumns, Notes: []string{identityNote}, Examples: []cli.Example{
			{Title: "Project a token for an external system's STS", Command: "afs " + k.svc + " put-" + k.svc[:len(k.svc)-1] + "-federated-identity" + who + " --name vault --audience vault.corp.example --expiration-seconds 3600"},
		}})
		r.Override(cli.Override{Op: "Delete" + k.kind + "FederatedIdentity", Examples: []cli.Example{
			{Title: "Remove a federated identity", Command: "afs " + k.svc + " delete-" + k.svc[:len(k.svc)-1] + "-federated-identity" + who + " --name vault"},
		}})
	}
	r.Override(cli.Override{Op: "ListMcpEvents", Columns: eventColumns, Examples: []cli.Example{
		{Title: "List the refused MCP calls of the last day", Command: "afs audit list-mcp-events --from 2026-11-24T00:00:00Z --to 2026-11-25T00:00:00Z --outcome refused"},
		{Title: "List one agent's calls to one server", Command: "afs audit list-mcp-events --team-id " + exTeam + " --agent-id " + exAgent + " --server-id " + exServer},
	}})
	r.Override(cli.Override{Op: "ListTeamMcpEvents", Columns: eventColumns, Examples: []cli.Example{
		{Title: "List the team's MCP calls that timed out", Command: "afs audit list-team-mcp-events --team-id " + exTeam + " --outcome timeout"},
	}})
	r.Override(cli.Override{Op: "GetMcpUsageReport", Columns: usageColumns, ResultKey: "rows", Examples: []cli.Example{
		{Title: "Calls per server this month", Command: "afs reports get-mcp-usage-report --from 2026-11-01T00:00:00Z --to 2026-12-01T00:00:00Z --group-by server --view usage"},
		{Title: "Servers no team used in 90 days", Command: "afs reports get-mcp-usage-report --from 2026-08-27T00:00:00Z --to 2026-11-25T00:00:00Z --group-by server --view unused --window 90"},
	}})
	r.Override(cli.Override{Op: "GetTeamMcpUsageReport", Columns: usageColumns, ResultKey: "rows", Examples: []cli.Example{
		{Title: "The team's refusals per agent", Command: "afs reports get-team-mcp-usage-report --team-id " + exTeam + " --from 2026-11-01T00:00:00Z --to 2026-12-01T00:00:00Z --group-by agent --view refusals"},
	}})
}
