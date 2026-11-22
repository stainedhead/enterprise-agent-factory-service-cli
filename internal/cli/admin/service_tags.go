package admin

import (
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/output"
)

func init() { registrations = append(registrations, registerServiceTags) }

var serviceTagColumns = []output.Column{
	{Header: "Key", Path: "key"}, {Header: "Mode", Path: "mode"}, {Header: "Value", Path: "value"},
	{Header: "DefaultValue", Path: "defaultValue"}, {Header: "AllowedValues", Path: "allowedValues"},
	{Header: "Description", Path: "description"}, {Header: "Version", Path: "version"},
}

func registerServiceTags(r *cli.Registry) {
	r.Override(cli.Override{Op: "ListServiceTags", Columns: serviceTagColumns, Examples: []cli.Example{
		{Title: "List the service tags", Command: "afs admin list-service-tags"},
	}})
	r.Override(cli.Override{Op: "PutServiceTag", Columns: serviceTagColumns, Examples: []cli.Example{
		{Title: "Require a cost-center tag with allowed values", Command: `afs admin put-service-tag --key cost-center --mode required --allowed-values CC-1001 CC-2002 --description "Finance cost center"`},
		{Title: "Apply a fixed tag to every team", Command: "afs admin put-service-tag --key environment --mode fixed --value production"},
	}})
	r.Override(cli.Override{Op: "DeleteServiceTag", Examples: []cli.Example{
		{Title: "Delete a service tag", Command: "afs admin delete-service-tag --key cost-center"},
	}})
}
