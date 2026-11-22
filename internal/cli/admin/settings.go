package admin

import (
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/output"
)

func init() { registrations = append(registrations, registerSettings) }

var settingColumns = []output.Column{
	{Header: "Key", Path: "key"}, {Header: "Value", Path: "value"}, {Header: "Type", Path: "type"},
	{Header: "Default", Path: "default"}, {Header: "Version", Path: "version"},
}

func registerSettings(r *cli.Registry) {
	r.Override(cli.Override{Op: "GetSettings", Columns: settingColumns, Examples: []cli.Example{
		{Title: "List every service setting", Command: "afs admin get-settings"},
		{Title: "Print the settings as JSON", Command: "afs admin get-settings --output json"},
	}})
	r.Override(cli.Override{Op: "PutSetting", Columns: settingColumns, Examples: []cli.Example{
		{Title: "Set a list value (JSON)", Command: `afs admin put-setting --key api.cors_origins --value '["https://console.example.com"]'`},
		{Title: "Set a number", Command: "afs admin put-setting --key api.idempotency_key_ttl_seconds --value 86400"},
	}})
}
