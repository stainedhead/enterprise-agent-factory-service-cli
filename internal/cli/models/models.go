package models

import (
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/output"
)

func init() { registrations = append(registrations, registerModels) }

func registerModels(r *cli.Registry) {
	r.Override(cli.Override{Op: "ListModels", Columns: []output.Column{
		{Header: "Id", Path: "id"}, {Header: "Provider", Path: "provider"}, {Header: "Name", Path: "displayName"},
		{Header: "Enabled", Path: "enabled"}, {Header: "Context", Path: "contextWindow"}, {Header: "MaxOutput", Path: "maxOutputTokens"},
	}, Examples: []cli.Example{
		{Title: "List the models the service offers", Command: "afs models list-models"},
		{Title: "Print only the model ids", Command: "afs models list-models --query 'items[].id' --output text"},
	}})
}
