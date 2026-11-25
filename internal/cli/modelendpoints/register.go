package modelendpoints

import "github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"

// registrations are appended by each file's init and applied in file-name order.
var registrations []func(*cli.Registry)

// Register applies every registration of this service to the registry.
func Register(r *cli.Registry) {
	for _, f := range registrations {
		f(r)
	}
}

// Example ids, shaped like the contract's.
const (
	exEndpoint     = "mep_01JD2Q3R4S5T6V7W8X9Y0Z0001"
	exTeamEndpoint = "mep_01JD2Q3R4S5T6V7W8X9Y0Z0002"
	exTeam         = "team_01JD2Q3R4S5T6V7W8X9Y0Z0001"
	exProposal     = "epp_01JD2Q3R4S5T6V7W8X9Y0Z0001"
)
