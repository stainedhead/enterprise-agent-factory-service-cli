package swarmtemplates

import (
	"encoding/json"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/output"
)

func init() { registrations = append(registrations, registerTemplates) }

const exTemplate = "swt_01J9ZQ3ABCDEFGHJKMNPQRSTVW"

var templateColumns = []output.Column{
	{Header: "Id", Path: "id"}, {Header: "Name", Path: "name"}, {Header: "Kind", Path: "kind"}, {Header: "Enabled", Path: "enabled"},
	{Header: "BudgetUsd", Path: "budgetUsd"}, {Header: "TtlHours", Path: "ttlHours"}, {Header: "TeardownRule", Path: "teardownRule"},
	{Header: "Version", Path: "version"},
}

// specSkeleton shows the optional spec.orchestration member (ORC-7) in --generate-cli-skeleton; the spec is
// sent unchanged, so this is documentation only (CLI-US-35.1). spec.collaboration is deprecated until 2.0.0
// and is no longer shown.
var specSkeleton = map[string]json.RawMessage{"spec": json.RawMessage(`{"orchestration":{"tool":"","webUiExposure":"","webApiExposure":"","defaultModelId":""}}`)}

func registerTemplates(r *cli.Registry) {
	r.Override(cli.Override{Op: "CreateSwarmTemplate", SkeletonMembers: specSkeleton, Examples: []cli.Example{
		{Title: "Choose the tool and exposure of the remediation team's orchestration agent (spec.orchestration in the spec file)", Command: "afs swarm-templates create-swarm-template --name cve-orchestrated --kind cve_remediation --budget-usd 250.00 --ttl-hours 72 --teardown-rule at_work_complete --spec file://swarm-spec-with-orchestration.json"},
	}})
	r.Override(cli.Override{Op: "UpdateSwarmTemplate", SkeletonMembers: specSkeleton})
	r.Override(cli.Override{Op: "ListSwarmTemplates", Columns: templateColumns, Examples: []cli.Example{
		{Title: "List the swarm templates", Command: "afs swarm-templates list-swarm-templates"},
	}})
	r.Override(cli.Override{Op: "CreateSwarmTemplate", Columns: templateColumns, Examples: []cli.Example{
		{Title: "Create a CVE remediation template from a spec file", Command: "afs swarm-templates create-swarm-template --name cve-remediation --kind cve_remediation --budget-usd 250.00 --ttl-hours 72 --teardown-rule at_work_complete --spec file://swarm-spec.json"},
	}})
	r.Override(cli.Override{Op: "GetSwarmTemplate", Columns: templateColumns, Examples: []cli.Example{
		{Title: "Show a swarm template with its spec", Command: "afs swarm-templates describe-swarm-template --template-id " + exTemplate + " --output json"},
	}})
	r.Override(cli.Override{Op: "UpdateSwarmTemplate", Columns: templateColumns, Examples: []cli.Example{
		{Title: "Raise a template's budget", Command: "afs swarm-templates update-swarm-template --template-id " + exTemplate + " --budget-usd 500.00"},
	}})
	r.Override(cli.Override{Op: "DeleteSwarmTemplate", Examples: []cli.Example{
		{Title: "Delete a template that no open work order uses", Command: "afs swarm-templates delete-swarm-template --template-id " + exTemplate},
	}})
	r.Override(cli.Override{Op: "EnableSwarmTemplate", Columns: templateColumns, Examples: []cli.Example{
		{Title: "Allow new work orders to use a template", Command: "afs swarm-templates enable-swarm-template --template-id " + exTemplate},
	}})
	r.Override(cli.Override{Op: "DisableSwarmTemplate", Columns: templateColumns, Examples: []cli.Example{
		{Title: "Stop new work orders from using a template", Command: "afs swarm-templates disable-swarm-template --template-id " + exTemplate},
	}})
}
