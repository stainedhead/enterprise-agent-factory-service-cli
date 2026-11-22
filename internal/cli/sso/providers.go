package sso

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli/runtime"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/client"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/config"
)

// Provider is one enabled sign-in provider (no secret).
type Provider struct {
	ID          string `json:"id"`
	Kind        string `json:"kind"`
	Name        string `json:"name"`
	IssuerURL   string `json:"issuerUrl"`
	CLIClientID string `json:"cliClientId"`
	IsDefault   bool   `json:"isDefault"`
}

// Providers calls ListSignInProviders (public: no credential is sent).
func Providers(ctx context.Context, env *cli.Env, g cli.Globals) ([]Provider, error) {
	sess, err := runtime.NewLenient(env, g, true)
	if err != nil {
		return nil, err
	}
	resp, err := sess.Client.Do(ctx, client.Request{Method: "GET", Path: "/v1/identity-providers", NoAuth: true, Operation: "ListSignInProviders"})
	if err != nil {
		return nil, err
	}
	var doc struct {
		Items []Provider `json:"items"`
	}
	if err := json.Unmarshal(resp.Body, &doc); err != nil {
		return nil, fmt.Errorf("the provider list is not valid JSON: %w", err)
	}
	return doc.Items, nil
}

// Validate checks the fields that are written to ~/.afs/config or sent to the identity provider.
// The provider list is unauthenticated service data, so nothing in it is trusted (FR-001).
func (p Provider) Validate() error {
	bad := func(field string) error {
		return fmt.Errorf("the service returned an invalid sign-in provider (%s)", field)
	}
	if config.ValidateIdentifier("id", p.ID) != nil {
		return bad("id")
	}
	if config.ValidateIssuer("issuerUrl", p.IssuerURL) != nil {
		return bad("issuerUrl")
	}
	if config.ValidateIdentifier("cliClientId", p.CLIClientID) != nil {
		return bad("cliClientId")
	}
	return nil
}

// Pick chooses a provider: the one named by id, else the default, else the only one.
func Pick(ps []Provider, id string) (Provider, error) {
	p, err := pick(ps, id)
	if err != nil {
		return p, err
	}
	return p, p.Validate()
}

func pick(ps []Provider, id string) (Provider, error) {
	if len(ps) == 0 {
		return Provider{}, fmt.Errorf("the service has no sign-in providers enabled")
	}
	for _, p := range ps {
		if id != "" && p.ID == id {
			return p, nil
		}
	}
	if id != "" {
		return Provider{}, fmt.Errorf("provider %s is not enabled", id)
	}
	if len(ps) == 1 {
		return ps[0], nil
	}
	for _, p := range ps {
		if p.IsDefault {
			return p, nil
		}
	}
	return Provider{}, fmt.Errorf("several providers are enabled; run afs configure sso to choose one")
}
