package agents

import (
	"encoding/json"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
)

func init() { registrations = append(registrations, registerCollaboration) }

// collaborationOptions are the collaboration-tool options of create-agent and update-agent (CLI-21, CLI-US-28.2).
// They fold into the body's collaboration object: --collaboration-connect and --no-collaboration-connect
// send connect true or false, --paperclip-agent-id sends paperclipAgentId. With none of them the body has
// no collaboration member, so the API default applies.
var collaborationOptions = []cli.Option{
	{Flag: "--collaboration-connect", In: "local", Kind: cli.Bool, Desc: "Connect the agent to the team's collaboration tool (--no-collaboration-connect disconnects it; the default on create is to connect)."},
	{Flag: "--paperclip-agent-id", In: "local", Kind: cli.String, Desc: "The agent's id in Paperclip (Paperclip only, and only while registrationMode is manual)."},
}

func buildCollaboration(p *cli.Parsed, parts *cli.RequestParts) error {
	if !p.Has("collaboration-connect") && !p.Has("paperclip-agent-id") {
		return nil
	}
	if parts.Body == nil {
		return nil
	}
	obj := map[string]json.RawMessage{}
	if raw, ok := parts.Body.Get("collaboration"); ok {
		_ = json.Unmarshal(raw, &obj)
	}
	if p.Has("collaboration-connect") {
		obj["connect"] = json.RawMessage(p.Value("collaboration-connect"))
	}
	if p.Has("paperclip-agent-id") {
		b, _ := json.Marshal(p.Value("paperclip-agent-id"))
		obj["paperclipAgentId"] = b
	}
	b, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	parts.Body.Set("collaboration", b)
	return nil
}

func registerCollaboration(r *cli.Registry) {
	r.Override(cli.Override{Op: "CreateAgent", Extra: collaborationOptions, Build: buildCollaboration, Examples: []cli.Example{
		{Title: "Create a worker and wait until it is ready (it connects to the team's orchestration agent by default)", Command: "afs agents create-agent --team-id " + exTeam + " --name coder --template claude-code --image " + exImage + " --default-model-id anthropic.claude-sonnet-4-5 --wait"},
		{Title: "Create a worker that does not connect to the team's collaboration tool", Command: "afs agents create-agent --team-id " + exTeam + " --name coder --template claude-code --image " + exImage + " --default-model-id anthropic.claude-sonnet-4-5 --no-collaboration-connect"},
	}})
	r.Override(cli.Override{Op: "UpdateAgent", Extra: collaborationOptions, Build: buildCollaboration, Examples: []cli.Example{
		{Title: "Connect an existing worker to the team's collaboration tool (the same effect as afs collaboration connect-agent)", Command: "afs agents update-agent --team-id " + exTeam + " --agent-id " + exAgent + " --collaboration-connect"},
	}})
	r.Override(cli.Override{Op: "ListAgents", Examples: []cli.Example{
		{Title: "Find the team's orchestration agent", Command: "afs agents list-agents --team-id " + exTeam + " --role orchestration"},
	}})
}
