package agents

import (
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/output"
)

func init() { registrations = append(registrations, registerSnapshots) }

var snapshotColumns = []output.Column{
	{Header: "Id", Path: "id"}, {Header: "Agent", Path: "agentId"}, {Header: "State", Path: "state"},
	{Header: "SizeGiB", Path: "sizeGiB"}, {Header: "RequestedAt", Path: "requestedAt"}, {Header: "CompletedAt", Path: "completedAt"},
}

func registerSnapshots(r *cli.Registry) {
	r.Override(cli.Override{Op: "ListAgentSnapshots", Columns: snapshotColumns, Examples: []cli.Example{
		{Title: "List an agent's volume snapshots", Command: "afs agents list-snapshots --team-id " + exTeam + " --agent-id " + exAgent},
	}})
	r.Override(cli.Override{Op: "CreateAgentSnapshot", Columns: snapshotColumns, Examples: []cli.Example{
		{Title: "Take a snapshot of the agent's volume", Command: `afs agents create-snapshot --team-id ` + exTeam + ` --agent-id ` + exAgent + ` --description "before upgrade"`},
	}})
}
