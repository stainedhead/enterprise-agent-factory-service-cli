package agents

import "github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"

func init() { registrations = append(registrations, registerLifecycle) }

func registerLifecycle(r *cli.Registry) {
	r.Override(cli.Override{Op: "UpdateAgent", Columns: agentColumns, Examples: []cli.Example{
		{Title: "Change the default model (reads the current version, sends If-Match)", Command: "afs agents update-agent --team-id " + exTeam + " --agent-id " + exAgent + " --default-model-id anthropic.claude-haiku-4-5"},
		{Title: "Apply a change only if the agent is still at version 7", Command: "afs agents update-agent --team-id " + exTeam + " --agent-id " + exAgent + " --description 'Reviews pull requests' --if-match 7"},
	}})
	r.Override(cli.Override{Op: "DeleteAgent", Columns: agentColumns, Examples: []cli.Example{
		{Title: "Delete an agent and wait until it is gone", Command: "afs agents delete-agent --team-id " + exTeam + " --agent-id " + exAgent + " --wait"},
	}})
	r.Override(cli.Override{Op: "PauseAgent", Columns: agentColumns, Examples: []cli.Example{
		{Title: "Pause an agent", Command: "afs agents pause-agent --team-id " + exTeam + " --agent-id " + exAgent},
	}})
	r.Override(cli.Override{Op: "ResumeAgent", Columns: agentColumns, Examples: []cli.Example{
		{Title: "Resume a paused agent", Command: "afs agents resume-agent --team-id " + exTeam + " --agent-id " + exAgent},
	}})
	r.Override(cli.Override{Op: "RestartAgent", Columns: agentColumns, Examples: []cli.Example{
		{Title: "Restart an agent's pods", Command: "afs agents restart-agent --team-id " + exTeam + " --agent-id " + exAgent + " --wait"},
	}})
}
