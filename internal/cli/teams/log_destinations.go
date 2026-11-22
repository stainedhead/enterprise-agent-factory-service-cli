package teams

import (
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/output"
)

func init() { registrations = append(registrations, registerLogDestinations) }

// logDestinationColumns never include a credential: config is non-secret settings and secretId
// names the secret (CLI-US-19.2).
var logDestinationColumns = []output.Column{
	{Header: "Id", Path: "id"}, {Header: "Name", Path: "name"}, {Header: "Kind", Path: "kind"},
	{Header: "Enabled", Path: "enabled"}, {Header: "Phase", Path: "status.phase"}, {Header: "Config", Path: "config"},
	{Header: "SecretId", Path: "secretId"}, {Header: "Version", Path: "version"},
}

const exLog = "log_01J9ZQ3ABCDEFGHJKMNPQRSTVW"

func registerLogDestinations(r *cli.Registry) {
	r.Override(cli.Override{Op: "ListLogDestinations", Columns: logDestinationColumns, Rename: map[string]string{"id": "destination-id"}, Examples: []cli.Example{
		{Title: "List a team's log destinations", Command: "afs teams list-log-destinations --team-id " + exTeam},
	}})
	r.Override(cli.Override{Op: "CreateLogDestination", Columns: logDestinationColumns, Examples: []cli.Example{
		{Title: "Ship logs to CloudWatch and wait until it is applied", Command: `afs teams create-log-destination --team-id ` + exTeam + ` --name team-cloudwatch --kind cloudwatch --config '{"logGroup":"/afs/payments","region":"us-east-1"}' --wait`},
	}})
	r.Override(cli.Override{Op: "GetLogDestination", Columns: logDestinationColumns, Rename: map[string]string{"id": "destination-id"}, Examples: []cli.Example{
		{Title: "Show a destination's configuration and delivery status", Command: "afs teams describe-log-destination --team-id " + exTeam + " --destination-id " + exLog},
	}})
	r.Override(cli.Override{Op: "UpdateLogDestination", Columns: logDestinationColumns, Rename: map[string]string{"id": "destination-id"}, Examples: []cli.Example{
		{Title: "Disable a destination", Command: "afs teams update-log-destination --team-id " + exTeam + " --destination-id " + exLog + " --no-enabled"},
	}})
	r.Override(cli.Override{Op: "DeleteLogDestination", Columns: logDestinationColumns, Rename: map[string]string{"id": "destination-id"}, Examples: []cli.Example{
		{Title: "Delete a destination and wait until it is gone", Command: "afs teams delete-log-destination --team-id " + exTeam + " --destination-id " + exLog + " --wait"},
	}})
}
