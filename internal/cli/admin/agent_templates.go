package admin

import (
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/output"
)

func init() { registrations = append(registrations, registerAgentTemplates) }

const exDigest = "sha256:4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f"

var agentTemplateColumns = []output.Column{
	{Header: "Name", Path: "name"}, {Header: "Harness", Path: "harness"}, {Header: "Digest", Path: "currentDigest"},
	{Header: "Enabled", Path: "enabled"}, {Header: "Version", Path: "version"},
}

func registerAgentTemplates(r *cli.Registry) {
	r.Override(cli.Override{Op: "ListAdminAgentTemplates", Columns: agentTemplateColumns, Examples: []cli.Example{
		{Title: "List every agent template", Command: "afs admin list-agent-templates"},
	}})
	r.Override(cli.Override{Op: "PutAgentTemplate", Columns: agentTemplateColumns, Examples: []cli.Example{
		{Title: "Create or replace a template", Command: "afs admin put-agent-template --name claude-code --harness claude-code --current-digest " + exDigest + " --allowed-digests " + exDigest + ` --default-resources '{"cpu":"1","memory":"2Gi"}'`},
	}})
	r.Override(cli.Override{Op: "DeleteAgentTemplate", Examples: []cli.Example{
		{Title: "Delete a template no agent uses", Command: "afs admin delete-agent-template --name claude-code"},
	}})
}
