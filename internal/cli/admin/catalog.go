package admin

import (
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/output"
)

func init() { registrations = append(registrations, registerCatalog) }

const exSystem = "ext_01JD2Q3R4S5T6V7W8X9Y0Z0001"

var catalogColumns = []output.Column{
	{Header: "Id", Path: "id"}, {Header: "Slug", Path: "slug"}, {Header: "Name", Path: "name"},
	{Header: "Category", Path: "category"}, {Header: "Enabled", Path: "enabled"}, {Header: "Builtin", Path: "builtin"},
	{Header: "Version", Path: "version"},
}

func registerCatalog(r *cli.Registry) {
	r.Override(cli.Override{Op: "ListCatalogSystems", Columns: catalogColumns, Examples: []cli.Example{
		{Title: "List the catalog of external systems", Command: "afs admin list-external-systems"},
	}})
	r.Override(cli.Override{Op: "GetCatalogSystem", Columns: catalogColumns, Examples: []cli.Example{
		{Title: "Show one catalog system", Command: "afs admin describe-external-system --system-id " + exSystem},
	}})
	r.Override(cli.Override{Op: "CreateCatalogSystem", Columns: catalogColumns, Examples: []cli.Example{
		{Title: "Add a system to the catalog", Command: "afs admin create-external-system --slug jira --name Jira --category collaboration --endpoints example.atlassian.net:443 --docs-url https://support.atlassian.com/jira"},
	}})
	r.Override(cli.Override{Op: "UpdateCatalogSystem", Columns: catalogColumns, Examples: []cli.Example{
		{Title: "Disable a catalog system", Command: "afs admin update-external-system --system-id " + exSystem + " --no-enabled"},
	}})
	r.Override(cli.Override{Op: "DeleteCatalogSystem", Examples: []cli.Example{
		{Title: "Delete a catalog system (built-in ones are refused with 409)", Command: "afs admin delete-external-system --system-id " + exSystem},
	}})
}
