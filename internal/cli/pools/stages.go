package pools

import (
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/output"
)

func init() { registrations = append(registrations, registerStages) }

const (
	exTeam = "team_01J9ZQ3ABCDEFGHJKMNPQRSTVW"
	exPool = "pool_01J9ZQ3ABCDEFGHJKMNPQRSTVW"
	exLink = "prl_01J9ZQ3ABCDEFGHJKMNPQRSTVW"
)

var stageColumns = []output.Column{
	{Header: "Id", Path: "id"}, {Header: "Name", Path: "name"}, {Header: "Position", Path: "position"},
	{Header: "Description", Path: "description"}, {Header: "Version", Path: "version"},
}

func registerStages(r *cli.Registry) {
	r.Override(cli.Override{Op: "ListStages", Columns: stageColumns, Examples: []cli.Example{
		{Title: "List a swarm team's stages", Command: "afs pools list-stages --team-id " + exTeam},
	}})
	r.Override(cli.Override{Op: "CreateStage", Columns: stageColumns, Examples: []cli.Example{
		{Title: "Add a stage at position 1", Command: `afs pools create-stage --team-id ` + exTeam + ` --name triage --position 1 --description "Find affected repositories"`},
	}})
	r.Override(cli.Override{Op: "UpdateStage", Columns: stageColumns, Examples: []cli.Example{
		{Title: "Move a stage", Command: "afs pools update-stage --team-id " + exTeam + " --stage-id stg_01J9ZQ3ABCDEFGHJKMNPQRSTVW --position 2"},
	}})
	r.Override(cli.Override{Op: "DeleteStage", Examples: []cli.Example{
		{Title: "Delete a stage", Command: "afs pools delete-stage --team-id " + exTeam + " --stage-id stg_01J9ZQ3ABCDEFGHJKMNPQRSTVW"},
	}})
}
