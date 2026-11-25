package mcp

import (
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/output"
)

func init() { registrations = append(registrations, registerTeam) }

// TeamServerColumns are the table columns of a server as one team sees it.
var TeamServerColumns = []output.Column{
	{Header: "Id", Path: "server.id"}, {Header: "Slug", Path: "server.slug"}, {Header: "Managed", Path: "managed"},
	{Header: "State", Path: "state"}, {Header: "Mode", Path: "mode"}, {Header: "Review", Path: "server.reviewStatus"},
	{Header: "Credential", Path: "server.credentialState"}, {Header: "Federation", Path: "federationState", Null: "-"},
}

var oauthColumns = []output.Column{{Header: "AuthorizationUrl", Path: "authorizationUrl"}, {Header: "ExpiresAt", Path: "expiresAt"}}

func registerTeam(r *cli.Registry) {
	r.Override(cli.Override{Op: "ListTeamMcpServers", Columns: TeamServerColumns, Examples: []cli.Example{
		{Title: "List the servers a team can use", Command: "afs teams list-team-mcp-servers --team-id " + exTeam},
		{Title: "List only the servers that are on", Command: "afs teams list-team-mcp-servers --team-id " + exTeam + " --state on"},
	}})
	r.Override(cli.Override{Op: "GetTeamMcpServer", Columns: TeamServerColumns, Examples: []cli.Example{
		{Title: "Show a server as the team sees it", Command: "afs teams describe-team-mcp-server --team-id " + exTeam + " --mcp-server-id " + exServer},
	}})
	r.Override(cli.Override{Op: "PutTeamMcpServer", Columns: TeamServerColumns,
		Notes: []string{"Only a team on its own gateway (per_namespace) can switch a server on. read_only mode admits only tools the review marked read-only."},
		Examples: []cli.Example{
			{Title: "Switch a server on, read-only", Command: "afs teams put-team-mcp-server --team-id " + exTeam + " --mcp-server-id " + exServer + " --enabled --mode read_only"},
			{Title: "Switch it off", Command: "afs teams put-team-mcp-server --team-id " + exTeam + " --mcp-server-id " + exServer + " --no-enabled"},
		}})
	r.Override(cli.Override{Op: "CreateTeamMcpServer", Columns: ServerColumns, Desc: serverDesc, Kinds: serverKinds,
		Notes: []string{"Allowed when the team's gateway settings have mayAddMcpServers (otherwise AccessDenied, reason team_servers_not_allowed). The server starts as a draft; submit it for review."},
		Examples: []cli.Example{
			{Title: "Register the team's own server", Command: `afs teams create-team-mcp-server --team-id ` + exTeam + ` --slug payments-tools --display-name "Payments tools" --owner payments-team --purpose "Ledger lookups" --data-classification confidential --transport streamable_http --location corporate --url https://tools.payments.corp.example/mcp --auth-kind workload_federation --federation file://federation.json`},
		}})
	r.Override(cli.Override{Op: "UpdateTeamMcpServer", Columns: ServerColumns, Desc: serverDesc, Kinds: serverKinds, Examples: []cli.Example{
		{Title: "Point the team's server at a new URL", Command: "afs teams update-team-mcp-server --team-id " + exTeam + " --mcp-server-id " + exTeamSrv + " --url https://tools-v2.payments.corp.example/mcp"},
	}})
	r.Override(cli.Override{Op: "DeleteTeamMcpServer", Examples: []cli.Example{
		{Title: "Delete a draft team server", Command: "afs teams delete-team-mcp-server --team-id " + exTeam + " --mcp-server-id " + exTeamSrv},
	}})
	r.Override(cli.Override{Op: "SubmitMcpServerForReview", Columns: ServerColumns, Examples: []cli.Example{
		{Title: "Submit the team's server for review", Command: "afs teams submit-mcp-server-for-review --team-id " + exTeam + " --mcp-server-id " + exTeamSrv},
	}})
	r.Override(cli.Override{Op: "PutTeamMcpServerCredential", Columns: CredentialColumns, Notes: []string{secretNote}, Examples: []cli.Example{
		{Title: "Set the team's token for a per_team server, asked for without echo", Command: "afs teams put-team-mcp-server-credential --team-id " + exTeam + " --mcp-server-id " + exServer + " --secret-value -"},
	}})
	r.Override(cli.Override{Op: "PutTeamMcpPolicy", Columns: PolicyColumns, Desc: rulesDesc, Examples: []cli.Example{
		{Title: "Block the delete tools for the team", Command: `afs teams put-team-mcp-policy --team-id ` + exTeam + ` --mcp-server-id ` + exServer + ` --rules '[{"name":"delete_*","effect":"block"}]'`},
	}})
	r.Override(cli.Override{Op: "RefreshTeamMcpInventory", Columns: ServerColumns, Examples: []cli.Example{
		{Title: "Refresh the team server's inventory", Command: "afs teams refresh-team-mcp-inventory --team-id " + exTeam + " --mcp-server-id " + exTeamSrv},
	}})
	r.Override(cli.Override{Op: "ListTeamMcpPrimitives", Columns: PrimitiveColumns, Examples: []cli.Example{
		{Title: "List the tools of a server for the team", Command: "afs teams list-team-mcp-primitives --team-id " + exTeam + " --mcp-server-id " + exServer + " --kind tool"},
	}})
	r.Override(cli.Override{Op: "ApproveTeamMcpPrimitive", Columns: PrimitiveColumns, Examples: []cli.Example{
		{Title: "Approve a changed tool of the team's server", Command: "afs teams approve-team-mcp-primitive --team-id " + exTeam + " --mcp-server-id " + exTeamSrv + " --primitive-id " + exPrimitive + ` --note "New optional argument"`},
	}})
	r.Override(cli.Override{Op: "RejectTeamMcpPrimitive", Columns: PrimitiveColumns, Examples: []cli.Example{
		{Title: "Block a new tool of the team's server", Command: "afs teams reject-team-mcp-primitive --team-id " + exTeam + " --mcp-server-id " + exTeamSrv + " --primitive-id " + exPrimitive + ` --note "Not needed"`},
	}})
	r.Override(cli.Override{Op: "ApproveTeamMcpServer", Columns: ServerColumns, Desc: reviewDesc,
		Notes: []string{"A team admin approves a low-risk team server when mcp.approval_policy allows it; others go to a service admin or a security operator."},
		Examples: []cli.Example{
			{Title: "Approve the team's low-risk server", Command: "afs teams approve-team-mcp-server --team-id " + exTeam + " --mcp-server-id " + exTeamSrv + ` --note "Read-only lookups"`},
		}})
	r.Override(cli.Override{Op: "RejectTeamMcpServer", Columns: ServerColumns, Desc: reviewDesc, Examples: []cli.Example{
		{Title: "Reject the team's server", Command: "afs teams reject-team-mcp-server --team-id " + exTeam + " --mcp-server-id " + exTeamSrv + ` --note "Use the enterprise server"`},
	}})
	r.Override(cli.Override{Op: "StartMcpOAuthConnect", Columns: oauthColumns,
		Notes: []string{"Open the authorization URL in a browser before it expires; the service stores the resulting token and the gateway refreshes it. The state is single-use."},
		Examples: []cli.Example{
			{Title: "Connect the team's account to an oauth_user server", Command: "afs teams start-mcp-oauth-connect --team-id " + exTeam + " --mcp-server-id " + exServer},
			{Title: "Print only the URL to open", Command: "afs teams start-mcp-oauth-connect --team-id " + exTeam + " --mcp-server-id " + exServer + " --query authorizationUrl --output text"},
		}})
	r.Override(cli.Override{Op: "DisconnectMcpOAuth", Columns: CredentialColumns, Examples: []cli.Example{
		{Title: "Drop the team's OAuth connection", Command: "afs teams disconnect-mcp-oauth --team-id " + exTeam + " --mcp-server-id " + exServer},
	}})
	r.Override(cli.Override{Op: "TestMcpFederation", Columns: TeamServerColumns,
		Notes: []string{"Accepted (202): the team gateway exchanges its token once; read federationState afterwards."},
		Examples: []cli.Example{
			{Title: "Test the team gateway's token exchange", Command: "afs teams test-mcp-federation --team-id " + exTeam + " --mcp-server-id " + exTeamSrv},
		}})
}
