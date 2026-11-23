package pools

import (
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli/modelaccess"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/output"
)

func init() { registrations = append(registrations, registerPools) }

// poolColumns show the desired worker count, state and phase (D-C31); the reporting, working, idle
// and impaired counts are in get-pool-status.
var poolColumns = []output.Column{
	{Header: "Id", Path: "id"}, {Header: "Name", Path: "name"}, {Header: "Stage", Path: "stageId"},
	{Header: "Template", Path: "template"}, {Header: "DefaultModel", Path: "models.defaultModelId", Null: "-"},
	{Header: "Workers", Path: "workers"}, {Header: "State", Path: "state"},
	{Header: "Phase", Path: "status.phase"}, {Header: "Version", Path: "version"},
}

var scratchRename = map[string]string{"scratchSizeGiB": "scratch-size-gib"}

func registerPools(r *cli.Registry) {
	r.Override(cli.Override{Op: "ListPools", Columns: poolColumns, Examples: []cli.Example{
		{Title: "List a team's pools", Command: "afs pools list-pools --team-id " + exTeam},
	}})
	r.Override(cli.Override{Op: "CreatePool", Columns: poolColumns, Rename: scratchRename, Examples: []cli.Example{
		{Title: "Create a pool of 10 workers and wait until they report", Command: `afs pools create-pool --team-id ` + exTeam + ` --name fixer --template claude-code --image registry.example.com/agents/coder@sha256:4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f --models anthropic.claude-sonnet-4-5 anthropic.claude-haiku-4-5 --default-model-id anthropic.claude-sonnet-4-5 --workers 10 --wait`},
	}})
	r.Override(modelaccess.Override("CreatePool"))
	r.Override(modelaccess.Override("UpdatePool"))
	r.Override(cli.Override{Op: "UpdatePool", Examples: []cli.Example{
		{Title: "Let the pool's workers call every model the team may use", Command: "afs pools update-pool --team-id " + exTeam + " --pool-id " + exPool + " --models all"},
	}})
	r.Override(cli.Override{Op: "GetPool", Columns: poolColumns, Examples: []cli.Example{
		{Title: "Show a pool (desired workers, state, phase)", Command: "afs pools describe-pool --team-id " + exTeam + " --pool-id " + exPool},
	}})
	r.Override(cli.Override{Op: "UpdatePool", Columns: poolColumns, Rename: scratchRename, Examples: []cli.Example{
		{Title: "Change a pool's persona", Command: `afs pools update-pool --team-id ` + exTeam + ` --pool-id ` + exPool + ` --persona "Minimal-change dependency upgrader"`},
	}})
	r.Override(cli.Override{Op: "DeletePool", Columns: poolColumns, Examples: []cli.Example{
		{Title: "Delete a pool and wait until it is gone", Command: "afs pools delete-pool --team-id " + exTeam + " --pool-id " + exPool + " --wait"},
	}})
}
