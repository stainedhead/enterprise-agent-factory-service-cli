package agents

import (
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli/modelaccess"
)

func init() { registrations = append(registrations, registerModels) }

// registerModels gives create-agent and update-agent the model options of MOD-8 and takes --role out of
// help: an agent made with create-agent is always a worker, and the API answers any other value with its own
// ValidationError (role_reserved), so the CLI sends it unchanged (ORC-30).
func registerModels(r *cli.Registry) {
	create := modelaccess.Override("CreateAgent")
	create.Hidden = []string{"role"}
	create.OpenEnum = []string{"role"}
	create.Examples = []cli.Example{
		{Title: "Create a worker that may call any model the team may use, with a default model", Command: "afs agents create-agent --team-id " + exTeam + " --name coder --template claude-code --image " + exImage + " --models all --default-model-id anthropic.claude-sonnet-4-5"},
		{Title: "Create a worker limited to two models", Command: "afs agents create-agent --team-id " + exTeam + " --name reviewer --template claude-code --image " + exImage + " --models anthropic.claude-sonnet-4-5 anthropic.claude-haiku-4-5 --default-model-id anthropic.claude-haiku-4-5"},
	}
	r.Override(create)
	update := modelaccess.Override("UpdateAgent")
	update.Examples = []cli.Example{
		{Title: "Change the agent's default model only", Command: "afs agents update-agent --team-id " + exTeam + " --agent-id " + exAgent + " --default-model-id anthropic.claude-haiku-4-5"},
		{Title: "Limit the agent to one model", Command: "afs agents update-agent --team-id " + exTeam + " --agent-id " + exAgent + " --models anthropic.claude-sonnet-4-5 --default-model-id anthropic.claude-sonnet-4-5"},
		{Title: "Let the agent call every model its team may use", Command: "afs agents update-agent --team-id " + exTeam + " --agent-id " + exAgent + " --models all"},
	}
	r.Override(update)
	r.Override(cli.Override{Op: "ListAgents",
		Desc: map[string]string{"role": "worker or orchestration (collaboration is read as orchestration until 2.0.0)."}})
}
