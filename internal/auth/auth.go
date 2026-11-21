// Package auth resolves the credential of a request: the first match of AFS_API_KEY, the
// profile's api_key in ~/.afs/credentials, the profile's cached SSO token and the workload
// token file (CLI-5).
package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/client"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/config"
)

// LoginRequiredError means the profile's SSO session is missing or expired with no way to
// refresh it: the command makes no call and exits 255 with the hint.
type LoginRequiredError struct {
	Profile string
	Reason  string
}

func (e *LoginRequiredError) Error() string {
	r := e.Reason
	if r == "" {
		r = "You are not signed in"
	}
	return fmt.Sprintf("%s. %s", r, LoginHint(e.Profile))
}

// LoginHint is the instruction printed with every sign-in failure.
func LoginHint(profile string) string { return "Run afs sso login --profile " + profile }

// TokenSource yields a bearer token. ok is false when it has none to offer; err reports a
// failure (a *LoginRequiredError for a missing or expired SSO session).
type TokenSource interface {
	Token(ctx context.Context) (token string, ok bool, err error)
}

// Inputs are what the chain reads.
type Inputs struct {
	Source   config.Source
	Settings config.Settings
	// SSO is the profile's cached SSO token source; nil when the profile has no SSO keys.
	SSO TokenSource
	// Workload is the projected ServiceAccount token source (M6-B); nil until built.
	Workload TokenSource
}

// Chain returns the credential function for a client.
func Chain(in Inputs) client.CredentialFunc {
	return func(ctx context.Context) (client.Credential, error) {
		if key := in.Source.Getenv("AFS_API_KEY"); key != "" {
			return client.Credential{Header: "Bearer " + key, Source: "env AFS_API_KEY"}, nil
		}
		f, err := config.CredentialsFile(in.Source)
		if err != nil {
			return client.Credential{}, err
		}
		if key, ok := f.Get(in.Settings.Profile, "api_key"); ok && key != "" {
			return client.Credential{Header: "Bearer " + key, Source: fmt.Sprintf("credentials file %s profile %s", in.Settings.CredentialsPath, in.Settings.Profile)}, nil
		}
		var pending error
		if in.SSO != nil {
			tok, ok, err := in.SSO.Token(ctx)
			switch {
			case err != nil:
				var lr *LoginRequiredError
				if !errors.As(err, &lr) {
					return client.Credential{}, err
				}
				pending = err
			case ok:
				return client.Credential{Header: "Bearer " + tok, Source: "sso cache profile " + in.Settings.Profile}, nil
			}
		}
		if in.Workload != nil {
			if tok, ok, err := in.Workload.Token(ctx); err != nil {
				return client.Credential{}, err
			} else if ok {
				return client.Credential{Header: "Bearer " + tok, Source: "workload token file"}, nil
			}
		}
		return client.Credential{}, pending
	}
}
