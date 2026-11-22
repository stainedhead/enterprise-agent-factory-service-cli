package teams

import (
	"encoding/json"
	"fmt"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/client"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/output"
)

func init() { registrations = append(registrations, registerNotificationChannels) }

const exChannel = "nch_01JD2Q3R4S5T6V7W8X9Y0Z0001"

var channelColumns = []output.Column{
	{Header: "Id", Path: "id"}, {Header: "Name", Path: "name"}, {Header: "Kind", Path: "kind"},
	{Header: "Enabled", Path: "enabled"}, {Header: "Phase", Path: "status.phase"}, {Header: "Healthy", Path: "delivery.healthy"},
	{Header: "LastDelivered", Path: "delivery.lastDeliveredAt"}, {Header: "LastError", Path: "delivery.lastError"},
	{Header: "Version", Path: "version"},
}

func registerNotificationChannels(r *cli.Registry) {
	rename := map[string]string{"id": "channel-id"}
	r.Override(cli.Override{Op: "ListNotificationChannels", Columns: channelColumns, Examples: []cli.Example{
		{Title: "List the team's notification channels", Command: "afs teams list-notification-channels --team-id " + exTeam},
	}})
	r.Override(cli.Override{Op: "GetNotificationChannel", Rename: rename, Columns: channelColumns, Examples: []cli.Example{
		{Title: "Show a channel's configuration and delivery status (never a credential)", Command: "afs teams describe-notification-channel --team-id " + exTeam + " --channel-id " + exChannel},
	}})
	r.Override(cli.Override{Op: "CreateNotificationChannel", Columns: channelColumns, Examples: []cli.Example{
		{Title: "Create a Slack channel and wait until it is ready", Command: `afs teams create-notification-channel --team-id ` + exTeam + ` --name payments-slack --kind slack --config '{"channel":"#payments-alerts"}' --event-types 'budget.*' 'agent.status_changed' --secret-id sec_01JD2Q3R4S5T6V7W8X9Y0Z0003 --wait`},
	}})
	r.Override(cli.Override{Op: "UpdateNotificationChannel", Rename: rename, Columns: channelColumns, Examples: []cli.Example{
		{Title: "Disable a channel", Command: "afs teams update-notification-channel --team-id " + exTeam + " --channel-id " + exChannel + " --no-enabled"},
	}})
	r.Override(cli.Override{Op: "DeleteNotificationChannel", Rename: rename, Columns: channelColumns, Examples: []cli.Example{
		{Title: "Delete a channel and wait until it is gone", Command: "afs teams delete-notification-channel --team-id " + exTeam + " --channel-id " + exChannel + " --wait"},
	}})
	r.Override(cli.Override{Op: "TestNotificationChannel", Rename: rename, Dynamic: true, After: afterTest, Examples: []cli.Example{
		{Title: "Send a test notification (exits 254 when delivery fails)", Command: "afs teams test-notification-channel --team-id " + exTeam + " --channel-id " + exChannel},
	}})
}

// afterTest implements D-C20: a delivered test is exit 0; a failed delivery prints the AWS-style
// DeliveryFailed error and exits 254; a test still running (202) says how to look again.
func afterTest(inv *cli.Invocation, resp *client.Response) int {
	var doc struct {
		Delivered bool    `json:"delivered"`
		Pending   bool    `json:"pending"`
		Error     *string `json:"error"`
	}
	if json.Unmarshal(resp.Body, &doc) != nil {
		return -1
	}
	w := inv.Env.Stderr
	switch {
	case doc.Pending || resp.Status == 202:
		_, _ = fmt.Fprintln(w, "The test is still running; run afs teams describe-notification-channel to see the result.")
		return cli.ExitOK
	case !doc.Delivered:
		msg := "the test notification was not delivered"
		if doc.Error != nil && *doc.Error != "" {
			msg = *doc.Error
		}
		_, _ = fmt.Fprintf(w, "\nAn error occurred (DeliveryFailed) when calling the TestNotificationChannel operation: %s\n", msg)
		if id := resp.RequestID(); id != "" {
			_, _ = fmt.Fprintf(w, "Request ID: %s\n", id)
		}
		return cli.ExitService
	}
	return -1
}
