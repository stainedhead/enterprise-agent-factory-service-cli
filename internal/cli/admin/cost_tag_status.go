package admin

import (
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/output"
)

func init() { registrations = append(registrations, registerCostTagStatus) }

func registerCostTagStatus(r *cli.Registry) {
	r.Override(cli.Override{Op: "ListCostTagStatus", Columns: []output.Column{
		{Header: "Key", Path: "key"}, {Header: "AwsStatus", Path: "awsStatus"}, {Header: "SeenInCur", Path: "seenInCur"},
		{Header: "FirstSeenAt", Path: "firstSeenAt"}, {Header: "LastSeenAt", Path: "lastSeenAt"}, {Header: "CheckedAt", Path: "checkedAt"},
	}, Examples: []cli.Example{
		{Title: "Show which cost-allocation tags AWS has activated and seen in the cost report", Command: "afs admin list-cost-tag-status"},
	}})
}
