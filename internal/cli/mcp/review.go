package mcp

import (
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/output"
)

func init() { registrations = append(registrations, registerReview) }

// PrimitiveColumns are the table columns of a server's tools, resources and prompts.
var PrimitiveColumns = []output.Column{
	{Header: "Id", Path: "id"}, {Header: "Kind", Path: "kind"}, {Header: "Name", Path: "name"},
	{Header: "Status", Path: "status"}, {Header: "Change", Path: "change", Null: "-"}, {Header: "LastSeen", Path: "lastSeen"},
}

var (
	reviewColumns = []output.Column{
		{Header: "Kind", Path: "kind"}, {Header: "Server", Path: "serverId"}, {Header: "Slug", Path: "slug"},
		{Header: "Team", Path: "teamId", Null: "-"}, {Header: "Primitive", Path: "name", Null: "-"},
		{Header: "Risk", Path: "riskTier", Null: "-"}, {Header: "Approver", Path: "approver"}, {Header: "Since", Path: "since"},
	}
	changeColumns = []output.Column{
		{Header: "Name", Path: "primitive.name"}, {Header: "Kind", Path: "primitive.kind"}, {Header: "Status", Path: "primitive.status"},
		{Header: "Change", Path: "primitive.change", Null: "-"}, {Header: "CrossesRiskLine", Path: "crossesRiskLine"},
	}
	captureColumns = []output.Column{
		{Header: "Id", Path: "id"}, {Header: "StartedBy", Path: "startedBy"}, {Header: "StartedAt", Path: "startedAt"},
		{Header: "Until", Path: "until"}, {Header: "Reason", Path: "reason"}, {Header: "StoppedAt", Path: "stoppedAt", Null: "-"},
	}
	captureStateColumns = []output.Column{
		{Header: "Id", Path: "capture.id", Null: "-"}, {Header: "StartedBy", Path: "capture.startedBy", Null: "-"},
		{Header: "Until", Path: "capture.until", Null: "-"}, {Header: "Reason", Path: "capture.reason", Null: "-"},
	}
	reviewDesc = map[string]string{
		"primitives": "Per-primitive decisions as JSON or file://PATH: [{\"primitiveId\":…, \"decision\":\"approve\"|\"block\"}]; primitives not named keep their state.",
	}
)

func registerReview(r *cli.Registry) {
	r.Override(cli.Override{Op: "ListMcpReviews", Columns: reviewColumns, Examples: []cli.Example{
		{Title: "List everything that waits for a security approval", Command: "afs mcp-servers list-mcp-reviews --approver security"},
		{Title: "List held definition changes of high-risk servers", Command: "afs mcp-servers list-mcp-reviews --kind primitive --risk-tier high"},
	}})
	r.Override(cli.Override{Op: "ListMcpPrimitives", Columns: PrimitiveColumns, Examples: []cli.Example{
		{Title: "List a server's tools", Command: "afs mcp-servers list-mcp-primitives --mcp-server-id " + exServer + " --kind tool"},
		{Title: "List the definitions held for review", Command: "afs mcp-servers list-mcp-primitives --mcp-server-id " + exServer + " --status pending_review"},
	}})
	r.Override(cli.Override{Op: "GetMcpPrimitiveChange", Columns: changeColumns, Examples: []cli.Example{
		{Title: "Show what changed in a held definition", Command: "afs mcp-servers get-mcp-primitive-change --mcp-server-id " + exServer + " --primitive-id " + exPrimitive},
		{Title: "Show only the changed fields", Command: "afs mcp-servers get-mcp-primitive-change --mcp-server-id " + exServer + " --primitive-id " + exPrimitive + " --query diff"},
	}})
	r.Override(cli.Override{Op: "ApproveMcpPrimitive", Columns: PrimitiveColumns, Examples: []cli.Example{
		{Title: "Approve a changed tool definition", Command: "afs mcp-servers approve-mcp-primitive --mcp-server-id " + exServer + " --primitive-id " + exPrimitive + ` --note "Description only"`},
	}})
	r.Override(cli.Override{Op: "RejectMcpPrimitive", Columns: PrimitiveColumns, Examples: []cli.Example{
		{Title: "Block a new tool", Command: "afs mcp-servers reject-mcp-primitive --mcp-server-id " + exServer + " --primitive-id " + exPrimitive + ` --note "Writes to production"`},
	}})
	approval := []string{"The approval level follows the server's risk tier and mcp.approval_policy: a team admin, a service admin or a security operator. A caller below the level gets AccessDenied with reason approval_not_permitted_by_policy."}
	r.Override(cli.Override{Op: "ApproveMcpServer", Columns: ServerColumns, Desc: reviewDesc, Notes: approval, Examples: []cli.Example{
		{Title: "Approve a server", Command: "afs mcp-servers approve-mcp-server --mcp-server-id " + exServer + ` --note "Reviewed with the owning team"`},
		{Title: "Approve it and block one tool", Command: `afs mcp-servers approve-mcp-server --mcp-server-id ` + exServer + ` --note "Approved without delete" --primitives '[{"primitiveId":"` + exPrimitive + `","decision":"block"}]'`},
	}})
	r.Override(cli.Override{Op: "RejectMcpServer", Columns: ServerColumns, Desc: reviewDesc, Notes: approval, Examples: []cli.Example{
		{Title: "Reject a server", Command: "afs mcp-servers reject-mcp-server --mcp-server-id " + exServer + ` --note "Sends data to an unapproved region"`},
	}})
	r.Override(cli.Override{Op: "RequestMcpServerChanges", Columns: ServerColumns, Desc: reviewDesc, Examples: []cli.Example{
		{Title: "Send a server back to its owner", Command: "afs mcp-servers request-mcp-server-changes --mcp-server-id " + exServer + ` --note "Pin the image by digest"`},
	}})
	r.Override(cli.Override{Op: "PublishMcpServer", Columns: ServerColumns, Examples: []cli.Example{
		{Title: "Publish an approved server to every team's catalog", Command: "afs mcp-servers publish-mcp-server --mcp-server-id " + exServer + ` --note "Available to all teams"`},
	}})
	r.Override(cli.Override{Op: "SuspendMcpServer", Columns: ServerColumns,
		Notes: []string{"Every team gateway refuses calls to a suspended server within a minute; agents keep their configuration."},
		Examples: []cli.Example{
			{Title: "Suspend a server during an incident", Command: "afs mcp-servers suspend-mcp-server --mcp-server-id " + exServer + ` --reason "Vendor incident INC-1234"`},
		}})
	r.Override(cli.Override{Op: "ResumeMcpServer", Columns: ServerColumns, Examples: []cli.Example{
		{Title: "Resume a suspended server", Command: "afs mcp-servers resume-mcp-server --mcp-server-id " + exServer + ` --reason "Vendor incident closed"`},
	}})
	r.Override(cli.Override{Op: "RetireMcpServer", Columns: ServerColumns, Examples: []cli.Example{
		{Title: "Retire a server teams no longer use", Command: "afs mcp-servers retire-mcp-server --mcp-server-id " + exServer + ` --note "Replaced by github-v2"`},
	}})
	r.Override(cli.Override{Op: "GetMcpQueryCapture", Columns: captureStateColumns, Examples: []cli.Example{
		{Title: "Show whether search queries are being captured", Command: "afs mcp-servers get-mcp-query-capture"},
	}})
	r.Override(cli.Override{Op: "StartMcpQueryCapture", Columns: captureColumns,
		Notes: []string{"While a capture runs the ledger keeps the text of search queries; it stops by itself at the end of the window."},
		Examples: []cli.Example{
			{Title: "Capture search queries for a week", Command: `afs mcp-servers start-mcp-query-capture --days 7 --reason "Tune the search ranking"`},
		}})
	r.Override(cli.Override{Op: "StopMcpQueryCapture", Columns: captureColumns, Examples: []cli.Example{
		{Title: "Stop the capture now", Command: "afs mcp-servers stop-mcp-query-capture"},
	}})
}
