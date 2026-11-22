package audit

import (
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/output"
)

func init() { registrations = append(registrations, registerServiceAudit) }

func registerServiceAudit(r *cli.Registry) {
	r.Override(cli.Override{Op: "ListAdminAuditEvents", Columns: []output.Column{
		{Header: "Time", Path: "occurredAt"}, {Header: "Team", Path: "teamId"}, {Header: "Actor", Path: "actorId"},
		{Header: "Action", Path: "action"}, {Header: "Outcome", Path: "outcome"}, {Header: "Kind", Path: "resourceKind"},
		{Header: "Resource", Path: "resourceId"}, {Header: "Id", Path: "id"},
	}, Examples: []cli.Example{
		{Title: "Service-wide audit events as CSV (platform roles)", Command: "afs audit list-service-audit-events --from 2026-11-01T00:00:00Z --output csv"},
		{Title: "Denied events of one team", Command: "afs audit list-service-audit-events --team-id team_01J9ZQ3ABCDEFGHJKMNPQRSTVW --outcome denied"},
	}})
}
