package admin

import (
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/output"
)

func init() { registrations = append(registrations, registerServiceChannels) }

const exServiceChannel = "nch_01JD2Q3R4S5T6V7W8X9Y0Z0001"

var serviceChannelColumns = []output.Column{
	{Header: "Id", Path: "id"}, {Header: "Name", Path: "name"}, {Header: "Kind", Path: "kind"},
	{Header: "Enabled", Path: "enabled"}, {Header: "EventTypes", Path: "eventTypes"}, {Header: "Phase", Path: "status.phase"},
	{Header: "Healthy", Path: "delivery.healthy"}, {Header: "LastError", Path: "delivery.lastError"}, {Header: "Version", Path: "version"},
}

func registerServiceChannels(r *cli.Registry) {
	rename := map[string]string{"id": "channel-id"}
	r.Override(cli.Override{Op: "ListServiceNotificationChannels", Columns: serviceChannelColumns, Examples: []cli.Example{
		{Title: "List the service-level notification channels", Command: "afs admin list-notification-channels"},
	}})
	r.Override(cli.Override{Op: "CreateServiceNotificationChannel", Columns: serviceChannelColumns, Examples: []cli.Example{
		{Title: "Create an email channel for security events", Command: `afs admin create-notification-channel --name security-email --kind email --config '{"to":["security@example.com"]}' --event-types 'work_order.*'`},
	}})
	r.Override(cli.Override{Op: "UpdateServiceNotificationChannel", Rename: rename, Columns: serviceChannelColumns, Examples: []cli.Example{
		{Title: "Disable a channel", Command: "afs admin update-notification-channel --channel-id " + exServiceChannel + " --no-enabled"},
	}})
	r.Override(cli.Override{Op: "DeleteServiceNotificationChannel", Rename: rename, Examples: []cli.Example{
		{Title: "Delete a channel (refused while a setting names it)", Command: "afs admin delete-notification-channel --channel-id " + exServiceChannel},
	}})
}
