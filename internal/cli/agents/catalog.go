package agents

import (
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/output"
)

func init() { registrations = append(registrations, registerCatalog) }

var storageClassColumns = []output.Column{
	{Header: "Name", Path: "name"}, {Header: "Type", Path: "volumeType"}, {Header: "Enabled", Path: "enabled"},
	{Header: "Description", Path: "description"},
}

var templateColumns = []output.Column{
	{Header: "Name", Path: "name"}, {Header: "Harness", Path: "harness"}, {Header: "Digest", Path: "currentDigest"},
	{Header: "Capabilities", Path: "capabilities", Null: "-"}, {Header: "Enabled", Path: "enabled"}, {Header: "Version", Path: "version"},
}

func registerCatalog(r *cli.Registry) {
	r.Override(cli.Override{Op: "ListStorageClasses", Columns: storageClassColumns, Examples: []cli.Example{
		{Title: "List the storage classes agents may use", Command: "afs agents list-storage-classes"},
	}})
	r.Override(cli.Override{Op: "ListAgentTemplates", Columns: templateColumns, Examples: []cli.Example{
		{Title: "List the agent templates", Command: "afs agents list-agent-templates"},
		{Title: "List the templates that can run Paperclip", Command: "afs agents list-agent-templates --capability paperclip"},
	}})
}
