package pools

import (
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/output"
)

func init() { registrations = append(registrations, registerScale) }

// poolStatusColumns show every worker count (CLI-US-8.3).
var poolStatusColumns = []output.Column{
	{Header: "Id", Path: "id"}, {Header: "Phase", Path: "status.phase"}, {Header: "Desired", Path: "workers.desired"},
	{Header: "Reporting", Path: "workers.reporting"}, {Header: "Working", Path: "workers.working"},
	{Header: "Idle", Path: "workers.idle"}, {Header: "Impaired", Path: "workers.impaired"},
	{Header: "Generation", Path: "generation"},
}

func registerScale(r *cli.Registry) {
	r.Override(cli.Override{Op: "ScalePool", Columns: poolColumns, Examples: []cli.Example{
		{Title: "Scale a pool to 20 workers and wait until all 20 report", Command: "afs pools scale-pool --team-id " + exTeam + " --pool-id " + exPool + " --workers 20 --wait"},
	}})
	r.Override(cli.Override{Op: "PausePool", Columns: poolColumns, Examples: []cli.Example{
		{Title: "Pause a pool", Command: "afs pools pause-pool --team-id " + exTeam + " --pool-id " + exPool},
	}})
	r.Override(cli.Override{Op: "ResumePool", Columns: poolColumns, Examples: []cli.Example{
		{Title: "Resume a pool", Command: "afs pools resume-pool --team-id " + exTeam + " --pool-id " + exPool},
	}})
	r.Override(cli.Override{Op: "GetPoolStatus", Columns: poolStatusColumns, Examples: []cli.Example{
		{Title: "Show desired, reporting, working, idle and impaired workers", Command: "afs pools get-pool-status --team-id " + exTeam + " --pool-id " + exPool},
	}})
}
