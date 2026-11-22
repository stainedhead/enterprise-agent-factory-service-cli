package workorders

import (
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/output"
)

func init() { registrations = append(registrations, registerWorkOrders) }

const (
	exTeam      = "team_01J9ZQ3ABCDEFGHJKMNPQRSTVW"
	exWorkOrder = "wo_01J9ZQ3ABCDEFGHJKMNPQRSTVW"
	exTemplate  = "swt_01J9ZQ3ABCDEFGHJKMNPQRSTVW"
)

var workOrderColumns = []output.Column{
	{Header: "Id", Path: "id"}, {Header: "Title", Path: "title"}, {Header: "Severity", Path: "severity"},
	{Header: "Status", Path: "status"}, {Header: "SlaState", Path: "slaState"}, {Header: "SlaDueAt", Path: "slaDueAt"},
	{Header: "TeamId", Path: "teamId"}, {Header: "Version", Path: "version"},
}

var repoColumns = []output.Column{
	{Header: "Host", Path: "host"}, {Header: "Repo", Path: "repo"}, {Header: "Status", Path: "status"},
	{Header: "PrUrl", Path: "prUrl"}, {Header: "OwningTeamId", Path: "owningTeamId"},
	{Header: "OwningContact", Path: "owningContact"}, {Header: "OwnerSource", Path: "ownerSource"}, {Header: "Version", Path: "version"},
}

func registerWorkOrders(r *cli.Registry) {
	r.Override(cli.Override{Op: "CreateWorkOrder", Columns: workOrderColumns, Kinds: map[string]cli.Kind{"repos": cli.List},
		Build: expandRepos("repos"), Examples: []cli.Example{
			{Title: "Open a CVE work order; the team configuration decides who owns the first repository", Command: "afs work-orders create-work-order --template-id " + exTemplate + " --title \"Patch CVE-2026-1234\" --kind cve_remediation --cve-id CVE-2026-1234 --severity critical --sla-due-at 2026-12-01T00:00:00Z --repos github.com/acme/ledger github.com/acme/billing=team_01J9ZQ3ABCDEFGHJKMNPQRSTVX"},
			{Title: "Name an owning contact by email; retry safely with an idempotency key", Command: "afs work-orders create-work-order --template-id " + exTemplate + " --title \"Upgrade the runtime\" --severity medium --sla-due-at 2026-12-15T00:00:00Z --repos github.com/acme/ledger=ledger-team@example.com --idempotency-key upgrade-2026-12"},
		}})
	r.Override(cli.Override{Op: "ListWorkOrders", Columns: workOrderColumns, Examples: []cli.Example{
		{Title: "List open critical work orders", Command: "afs work-orders list-work-orders --status open --severity critical"},
		{Title: "Export every work order as CSV", Command: "afs work-orders list-work-orders --output csv"},
	}})
	r.Override(cli.Override{Op: "GetWorkOrder", Columns: workOrderColumns, Examples: []cli.Example{
		{Title: "Show a work order and its repositories", Command: "afs work-orders describe-work-order --work-order-id " + exWorkOrder + " --output json"},
	}})
	r.Override(cli.Override{Op: "UpdateWorkOrder", Columns: workOrderColumns, Kinds: map[string]cli.Kind{"addRepos": cli.List},
		Build: expandRepos("addRepos"), Examples: []cli.Example{
			{Title: "Move the SLA date and add a repository", Command: "afs work-orders update-work-order --work-order-id " + exWorkOrder + " --sla-due-at 2026-12-20T00:00:00Z --add-repos github.com/acme/gateway=team_01J9ZQ3ABCDEFGHJKMNPQRSTVX"},
		}})
	r.Override(cli.Override{Op: "CancelWorkOrder", Columns: workOrderColumns, Examples: []cli.Example{
		{Title: "Cancel a work order", Command: "afs work-orders cancel-work-order --work-order-id " + exWorkOrder + " --reason \"CVE withdrawn upstream\""},
	}})
	r.Override(cli.Override{Op: "ListWorkOrderRepos", Columns: repoColumns, Examples: []cli.Example{
		{Title: "List the repositories of a work order", Command: "afs work-orders list-repos --work-order-id " + exWorkOrder},
	}})
	r.Override(cli.Override{Op: "UpdateWorkOrderRepo", Columns: repoColumns, Optional: []string{"host", "owner"}, Build: splitRepoPath,
		Examples: []cli.Example{
			{Title: "Report a pull request for one repository", Command: "afs work-orders update-repo --work-order-id " + exWorkOrder + " --repo github.com/acme/ledger --status pr_opened --pr-url https://github.com/acme/ledger/pull/42 --pr-number 42"},
			{Title: "Report with the version printed by the previous update (a workload token cannot read it)", Command: "afs work-orders update-repo --work-order-id " + exWorkOrder + " --repo github.com/acme/ledger --status merged --if-match 3"},
		}})
	r.Override(cli.Override{Op: "ListTeamWorkOrders", Columns: workOrderColumns, Examples: []cli.Example{
		{Title: "List a team's work orders", Command: "afs work-orders list-team-work-orders --team-id " + exTeam + " --output json"},
	}})
}
