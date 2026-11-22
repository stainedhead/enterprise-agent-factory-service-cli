package agents

import (
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/output"
)

func init() { registrations = append(registrations, registerAgents) }

const (
	exTeam  = "team_01JD2Q3R4S5T6V7W8X9Y0Z0001"
	exAgent = "agt_01JD2Q3R4S5T6V7W8X9Y0Z0001"
	exImage = "registry.example.com/agents/coder@sha256:4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f"
)

var agentColumns = []output.Column{
	{Header: "Id", Path: "id"}, {Header: "Name", Path: "name"}, {Header: "Template", Path: "template"},
	{Header: "Model", Path: "modelId"}, {Header: "Paused", Path: "paused"}, {Header: "Phase", Path: "status.phase"},
	{Header: "Role", Path: "role"}, {Header: "Paperclip", Path: "collaboration.linkState", Null: "-"},
	{Header: "Version", Path: "version"},
}

var statusColumns = []output.Column{
	{Header: "Id", Path: "id"}, {Header: "Kind", Path: "kind"}, {Header: "Generation", Path: "generation"},
	{Header: "ObservedGeneration", Path: "status.observedGeneration"}, {Header: "Phase", Path: "status.phase"},
}

func registerAgents(r *cli.Registry) {
	r.Override(cli.Override{Op: "ListAgents", Columns: agentColumns, Examples: []cli.Example{
		{Title: "List a team's agents", Command: "afs agents list-agents --team-id " + exTeam},
		{Title: "Print agent ids and phases", Command: "afs agents list-agents --team-id " + exTeam + " --query 'items[].[id,status.phase]' --output text"},
	}})
	r.Override(cli.Override{Op: "GetAgent", Columns: agentColumns, Examples: []cli.Example{
		{Title: "Show one agent", Command: "afs agents describe-agent --team-id " + exTeam + " --agent-id " + exAgent},
	}})
	r.Override(cli.Override{Op: "GetAgentStatus", Columns: statusColumns, Examples: []cli.Example{
		{Title: "Show reconcile status", Command: "afs agents get-agent-status --team-id " + exTeam + " --agent-id " + exAgent},
	}})
	r.Override(cli.Override{Op: "CreateAgent", Columns: agentColumns, Examples: []cli.Example{
		{Title: "Create an agent with a volume and wait until it is ready", Command: "afs agents create-agent --team-id " + exTeam + " --name coder --template claude-code --image " + exImage +
			` --model-id anthropic.claude-sonnet-4-5 --volume '{"sizeGiB":20,"storageClass":"gp3-standard","onDelete":"snapshot_then_delete"}' --wait`},
		{Title: "Create an agent from a file", Command: "afs agents create-agent --team-id " + exTeam + " --cli-input-json file://agent.json"},
	}})
}
