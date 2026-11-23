package orchestration

import "github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"

// registrations are appended by each file's init and applied in file-name order.
var registrations []func(*cli.Registry)

// Register applies every registration of this service to the registry.
func Register(r *cli.Registry) {
	for _, f := range registrations {
		f(r)
	}
}
