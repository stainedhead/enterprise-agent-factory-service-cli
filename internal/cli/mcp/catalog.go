package mcp

import (
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/output"
)

func init() { registrations = append(registrations, registerCatalog) }

// ServerColumns are the table columns of an MCP server of the catalog (enterprise or team).
var ServerColumns = []output.Column{
	{Header: "Id", Path: "id"}, {Header: "Slug", Path: "slug"}, {Header: "Scope", Path: "scope"},
	{Header: "Team", Path: "teamId", Null: "-"}, {Header: "Location", Path: "location"}, {Header: "Transport", Path: "transport"},
	{Header: "Auth", Path: "authKind"}, {Header: "Review", Path: "reviewStatus"}, {Header: "Risk", Path: "riskTier", Null: "-"},
	{Header: "Published", Path: "published"}, {Header: "Credential", Path: "credentialState"}, {Header: "Version", Path: "version"},
}

// CredentialColumns are the table columns of an MCP credential write: its state, never a value.
var CredentialColumns = []output.Column{
	{Header: "Scope", Path: "scope"}, {Header: "Kind", Path: "kind", Null: "-"}, {Header: "State", Path: "state"},
	{Header: "Subject", Path: "subject", Null: "-"}, {Header: "ExpiresAt", Path: "expiresAt", Null: "-"},
	{Header: "RotatedAt", Path: "rotatedAt", Null: "-"},
}

// PolicyColumns are the table columns of a service or team policy of a server.
var PolicyColumns = []output.Column{
	{Header: "Server", Path: "serverId"}, {Header: "Level", Path: "level"}, {Header: "Team", Path: "teamId", Null: "-"},
	{Header: "DriftPolicy", Path: "driftPolicy", Null: "-"}, {Header: "TrustAnnotations", Path: "trustAnnotations", Null: "-"},
	{Header: "AllowSampling", Path: "allowSampling", Null: "-"}, {Header: "Version", Path: "version"},
}

var (
	serverDesc = map[string]string{
		"hosted":     "A hosted server as JSON or file://PATH: {\"imageDigest\":…, \"args\":[…], \"envNames\":[…], \"egressHosts\":[…], \"instanceMode\":\"per_team\", \"needsAws\":false}.",
		"federation": "Workload federation as JSON or file://PATH: the token endpoint and audience the gateway exchanges its projected token with.",
		"limits":     "The server's own limits as JSON: {\"requestBytes\":…, \"responseBytes\":…, \"timeoutSeconds\":…, \"concurrentPerAgent\":…}; unset members take mcp.limits.",
		"caBundle":   "PEM bundle of the CA that signs the server's certificate, as text or file://PATH (not a secret).",
	}
	serverKinds = map[string]cli.Kind{"caBundle": cli.Any}
	rulesDesc   = map[string]string{
		"rules": "Rules as JSON or file://PATH: [{\"name\":\"delete_*\",\"effect\":\"block\"}, {\"class\":\"destructive\",\"effect\":\"block\"}, {\"name\":\"search\",\"effect\":\"allow\",\"rateLimitPerMinute\":60}]. The most specific rule of a level wins; a block at any level wins over the levels below.",
	}
)

func registerCatalog(r *cli.Registry) {
	r.Override(cli.Override{Op: "ListMcpServers", Columns: ServerColumns, Examples: []cli.Example{
		{Title: "List the catalog", Command: "afs mcp-servers list-mcp-servers"},
		{Title: "List the enterprise servers that wait for review", Command: "afs mcp-servers list-mcp-servers --scope enterprise --review-status pending_review"},
		{Title: "List one team's own servers", Command: "afs mcp-servers list-mcp-servers --scope team --team-id " + exTeam},
	}})
	r.Override(cli.Override{Op: "CreateMcpServer", Columns: ServerColumns, Desc: serverDesc, Kinds: serverKinds,
		Notes: []string{"The server is created as a draft. Set its credential, refresh its inventory, then approve it; teams see it once it is approved and published."},
		Examples: []cli.Example{
			{Title: "Register a remote server with a static credential", Command: `afs mcp-servers create-mcp-server --slug github --display-name "GitHub" --owner platform-team --purpose "Issues and pull requests" --data-classification internal --transport streamable_http --location external --url https://api.githubcopilot.com/mcp/ --auth-kind static --credential-mode per_team`},
			{Title: "Register a hosted stdio server", Command: `afs mcp-servers create-mcp-server --slug jira --display-name "Jira" --data-classification confidential --transport hosted --location in_namespace --auth-kind static --hosted file://jira-hosted.json`},
		}})
	r.Override(cli.Override{Op: "ImportMcpServerManifest", Columns: ServerColumns, Examples: []cli.Example{
		{Title: "Register a server from its server.json manifest", Command: "afs mcp-servers import-mcp-server-manifest --manifest file://server.json --slug github"},
	}})
	r.Override(cli.Override{Op: "GetMcpServer", Columns: ServerColumns, Examples: []cli.Example{
		{Title: "Show a server", Command: "afs mcp-servers describe-mcp-server --mcp-server-id " + exServer},
		{Title: "Show why it is high risk", Command: "afs mcp-servers describe-mcp-server --mcp-server-id " + exServer + " --query riskReasons --output text"},
	}})
	r.Override(cli.Override{Op: "UpdateMcpServer", Columns: ServerColumns, Desc: serverDesc, Kinds: serverKinds, Examples: []cli.Example{
		{Title: "Approve compatible definition changes automatically", Command: "afs mcp-servers update-mcp-server --mcp-server-id " + exServer + " --drift-policy auto_approve_compatible"},
		{Title: "Raise the timeout", Command: `afs mcp-servers update-mcp-server --mcp-server-id ` + exServer + ` --limits '{"timeoutSeconds":120}'`},
	}})
	r.Override(cli.Override{Op: "DeleteMcpServer",
		Notes: []string{"Only a draft or rejected server can be deleted; retire an approved one instead."},
		Examples: []cli.Example{
			{Title: "Delete a draft", Command: "afs mcp-servers delete-mcp-server --mcp-server-id " + exServer},
		}})
	r.Override(cli.Override{Op: "PutMcpServerCredential", Columns: CredentialColumns,
		Notes: []string{"--scope is _shared for a shared credential, or a team id for a per_team server.", secretNote},
		Examples: []cli.Example{
			{Title: "Set a team's token from a file", Command: "afs mcp-servers put-mcp-server-credential --mcp-server-id " + exServer + " --scope " + exTeam + " --secret-value file://github-payments.token"},
			{Title: "Set the shared OAuth client from standard input", Command: "afs mcp-servers put-mcp-server-credential --mcp-server-id " + exServer + " --scope _shared --client-id afs-gateway --client-secret -"},
		}})
	r.Override(cli.Override{Op: "PutMcpPolicy", Columns: PolicyColumns, Desc: rulesDesc, Examples: []cli.Example{
		{Title: "Block every destructive tool of a server", Command: `afs mcp-servers put-mcp-policy --mcp-server-id ` + exServer + ` --rules '[{"class":"destructive","effect":"block"}]'`},
		{Title: "Replace the service policy from a file and trust the server's annotations", Command: "afs mcp-servers put-mcp-policy --mcp-server-id " + exServer + " --rules file://github-policy.json --trust-annotations"},
	}})
	r.Override(cli.Override{Op: "RefreshMcpInventory", Columns: ServerColumns,
		Notes: []string{"Accepted (202): the gateway lists the server's tools, resources and prompts; new and changed definitions are held for review."},
		Examples: []cli.Example{
			{Title: "Refresh a server's inventory now", Command: "afs mcp-servers refresh-mcp-inventory --mcp-server-id " + exServer},
		}})
}
