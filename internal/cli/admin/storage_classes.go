package admin

import (
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/output"
)

func init() { registrations = append(registrations, registerStorageClasses) }

var storageClassColumns = []output.Column{
	{Header: "Name", Path: "name"}, {Header: "Type", Path: "volumeType"}, {Header: "Enabled", Path: "enabled"},
	{Header: "Description", Path: "description"}, {Header: "Version", Path: "version"},
}

func registerStorageClasses(r *cli.Registry) {
	r.Override(cli.Override{Op: "ListStorageClassesAdmin", Columns: storageClassColumns, Examples: []cli.Example{
		{Title: "List every storage class", Command: "afs admin list-storage-classes"},
	}})
	r.Override(cli.Override{Op: "PutStorageClass", Columns: storageClassColumns, Examples: []cli.Example{
		{Title: "Create or replace a storage class", Command: `afs admin put-storage-class --name gp3-standard --volume-type gp3 --description "General purpose SSD"`},
	}})
	r.Override(cli.Override{Op: "DeleteStorageClass", Examples: []cli.Example{
		{Title: "Delete a storage class no volume uses", Command: "afs admin delete-storage-class --name gp3-standard"},
	}})
}
