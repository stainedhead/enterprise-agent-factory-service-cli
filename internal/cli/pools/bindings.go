package pools

import (
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/output"
)

func init() { registrations = append(registrations, registerBindings) }

var (
	handoffColumns = []output.Column{{Header: "ToPoolIds", Path: "toPoolIds"}, {Header: "Version", Path: "version"}}

	secretBindingColumns = []output.Column{
		{Header: "SecretId", Path: "secretId"}, {Header: "Delivery", Path: "delivery"}, {Header: "Target", Path: "target"},
	}

	poolConfigColumns = []output.Column{
		{Header: "Name", Path: "name"}, {Header: "Source", Path: "source"}, {Header: "Value", Path: "value"},
		{Header: "Resolved", Path: "resolvedValue"}, {Header: "UsedBy", Path: "usedBy"}, {Header: "Version", Path: "version"},
	}
)

func registerBindings(r *cli.Registry) {
	r.Override(cli.Override{Op: "ListPoolHandoffs", Columns: handoffColumns, Examples: []cli.Example{
		{Title: "Show where a pool hands work to", Command: "afs pools list-handoffs --team-id " + exTeam + " --pool-id " + exPool},
	}})
	r.Override(cli.Override{Op: "PutPoolHandoffs", Columns: handoffColumns, Examples: []cli.Example{
		{Title: "Hand off to two pools", Command: "afs pools put-handoffs --team-id " + exTeam + " --pool-id " + exPool + " --to-pool-ids pool_01J9ZQ3ABCDEFGHJKMNPQRSTVX pool_01J9ZQ3ABCDEFGHJKMNPQRSTVY"},
	}})
	r.Override(cli.Override{Op: "ListPoolSecretBindings", Columns: secretBindingColumns, Examples: []cli.Example{
		{Title: "List a pool's secret bindings", Command: "afs pools list-secret-bindings --team-id " + exTeam + " --pool-id " + exPool},
	}})
	r.Override(cli.Override{Op: "PutPoolSecretBindings", Columns: secretBindingColumns, Examples: []cli.Example{
		{Title: "Bind a secret as a file (the value is never shown)", Command: `afs pools put-secret-bindings --team-id ` + exTeam + ` --pool-id ` + exPool + ` --items '[{"secretId":"sec_01J9ZQ3ABCDEFGHJKMNPQRSTVW","delivery":"file","target":"/var/run/afs/secrets/github-app.pem"}]'`},
	}})
	r.Override(cli.Override{Op: "ListPoolConfigValues", Columns: poolConfigColumns, Examples: []cli.Example{
		{Title: "List a pool's configuration values", Command: "afs pools list-config-values --team-id " + exTeam + " --pool-id " + exPool},
	}})
	r.Override(cli.Override{Op: "PutPoolConfigValue", Columns: poolConfigColumns, Examples: []cli.Example{
		{Title: "Set a literal value", Command: "afs pools put-config-value --team-id " + exTeam + " --pool-id " + exPool + " --name LOG_LEVEL --source literal --value debug"},
	}})
	r.Override(cli.Override{Op: "DeletePoolConfigValue", Examples: []cli.Example{
		{Title: "Remove a configuration value", Command: "afs pools delete-config-value --team-id " + exTeam + " --pool-id " + exPool + " --name LOG_LEVEL"},
	}})
}
