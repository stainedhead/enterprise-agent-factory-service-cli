package sso

import "github.com/stainedhead/enterprise-agent-factory-service-cli/internal/client/gen"

// contractFields references every member of ListSignInProviders that Provider decodes by hand, so a
// rename in the contract fails the build instead of showing up at run time (FR-017, ADR-014).
var _ = func(p gen.SignInProvider) {
	_, _, _, _, _, _ = p.Id, p.Kind, p.Name, p.IssuerUrl, p.CliClientId, p.IsDefault
}
