package sso

import (
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/output"
)

func init() {
	registrations = append(registrations, func(r *cli.Registry) {
		r.Override(cli.Override{
			Op: "GetMe",
			Examples: []cli.Example{
				{Title: "Show who the credential belongs to", Command: "afs sso get-caller-identity"},
				{Title: "Show only the email", Command: "afs sso get-caller-identity --query principal.email --output text"},
			},
			Columns: []output.Column{{Header: "Id", Path: "principal.id"}, {Header: "Type", Path: "principal.type"}, {Header: "Email", Path: "principal.email"}, {Header: "ServiceAdmin", Path: "serviceAdmin"}},
		})
		r.Override(cli.Override{
			Op: "ListSignInProviders",
			Examples: []cli.Example{
				{Title: "List the sign-in providers", Command: "afs sso list-providers --output table"},
			},
			Columns: []output.Column{{Header: "Id", Path: "id"}, {Header: "Name", Path: "name"}, {Header: "Kind", Path: "kind"}, {Header: "Issuer", Path: "issuerUrl"}, {Header: "Default", Path: "isDefault"}},
		})
	})
}
