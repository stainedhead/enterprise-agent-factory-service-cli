package mcp

import "github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"

// registrations are appended by each file's init and applied in file-name order.
var registrations []func(*cli.Registry)

// Register applies every registration of this package to the registry.
func Register(r *cli.Registry) {
	for _, f := range registrations {
		f(r)
	}
}

// Example ids, shaped like the contract's.
const (
	exServer    = "mcp_01JD2Q3R4S5T6V7W8X9Y0Z0001"
	exTeamSrv   = "mcp_01JD2Q3R4S5T6V7W8X9Y0Z0002"
	exPrimitive = "mpr_01JD2Q3R4S5T6V7W8X9Y0Z0001"
	exTeam      = "team_01JD2Q3R4S5T6V7W8X9Y0Z0001"
	exAgent     = "agt_01JD2Q3R4S5T6V7W8X9Y0Z0001"
	exPool      = "pool_01JD2Q3R4S5T6V7W8X9Y0Z0001"
)

// secretNote is said on every MCP credential write.
const secretNote = //nolint:gosec // G101: help text about credentials, not a credential
"--secret-value (a static credential) and --client-secret (an OAuth client secret) are read from file://PATH, fileb://PATH or - (standard input); they are never plain arguments, and the API never returns them (only the credential's state, subject and expiry)."
