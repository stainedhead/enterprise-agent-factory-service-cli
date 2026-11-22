package audit

import (
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/output"
)

func init() { registrations = append(registrations, registerTeamAudit) }

func registerTeamAudit(r *cli.Registry) {
	r.Override(cli.Override{Op: "ListAuditEvents", Columns: []output.Column{
		{Header: "Time", Path: "occurredAt"}, {Header: "Actor", Path: "actorId"}, {Header: "Action", Path: "action"},
		{Header: "Outcome", Path: "outcome"}, {Header: "Kind", Path: "resourceKind"}, {Header: "Resource", Path: "resourceId"},
		{Header: "Id", Path: "id"},
	}, Examples: []cli.Example{
		{Title: "List a team's audit events as CSV", Command: "afs audit list-audit-events --team-id team_01J9ZQ3ABCDEFGHJKMNPQRSTVW --from 2026-11-01T00:00:00Z --outcome denied --output csv"},
		{Title: "Events of one resource", Command: "afs audit list-audit-events --team-id team_01J9ZQ3ABCDEFGHJKMNPQRSTVW --resource-kind agent --resource-id agt_01J9ZQ3ABCDEFGHJKMNPQRSTVW"},
	}})
}
