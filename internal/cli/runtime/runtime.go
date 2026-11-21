// Package runtime builds the per-invocation session: resolved settings, credential chain and API client.
package runtime

import (
	"errors"
	"net/http"
	"time"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/auth"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/buildinfo"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/client"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/config"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/spec"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/sso"
)

// ErrNoEndpoint is returned when no endpoint is configured anywhere.
var ErrNoEndpoint = errors.New("no endpoint is configured: use --endpoint-url, set AFS_ENDPOINT_URL or run afs configure")

// Session is what a running API command needs.
type Session struct {
	Settings config.Settings
	Client   *client.Client
	Source   config.Source
}

// Source returns the config source of an invocation.
func Source(env *cli.Env) config.Source { return config.Source{Getenv: env.Getenv, Home: env.Home} }

// Settings resolves the settings of an invocation (exit 255 on error).
func Settings(env *cli.Env, g cli.Globals) (config.Settings, error) {
	return settings(env, g, false)
}

// LenientSettings is Settings for commands that create a profile: an unknown profile is accepted.
func LenientSettings(env *cli.Env, g cli.Globals) (config.Settings, error) {
	return settings(env, g, true)
}

func settings(env *cli.Env, g cli.Globals, lenient bool) (config.Settings, error) {
	return config.Resolve(config.Overrides{Profile: g.Profile, EndpointURL: g.EndpointURL, Output: g.Output, AllowMissing: lenient}, Source(env))
}

// New resolves settings and builds a client. noAuth sends no credential.
func New(env *cli.Env, g cli.Globals, noAuth bool) (*Session, error) {
	return open(env, g, noAuth, false)
}

// NewLenient is New for commands that create a profile (afs configure sso, afs sso login).
func NewLenient(env *cli.Env, g cli.Globals, noAuth bool) (*Session, error) {
	return open(env, g, noAuth, true)
}

func open(env *cli.Env, g cli.Globals, noAuth, lenient bool) (*Session, error) {
	st, err := settings(env, g, lenient)
	if err != nil {
		return nil, err
	}
	if st.EndpointURL == "" {
		return nil, ErrNoEndpoint
	}
	if err := config.ValidateEndpoint(st.EndpointURL); err != nil {
		return nil, err
	}
	src := Source(env)
	rt := env.RoundTripper
	if rt == nil {
		ct := client.DefaultConnectTimeout
		if g.Set["cli-connect-timeout"] {
			ct = time.Duration(g.ConnectTimeout) * time.Second
		}
		rt = client.NewTransport(client.TransportOptions{ConnectTimeout: ct, DialContext: env.Dial})
	}
	rto := client.DefaultReadTimeout
	if g.Set["cli-read-timeout"] {
		rto = time.Duration(g.ReadTimeout) * time.Second
	}
	in := auth.Inputs{Source: src, Settings: st}
	if st.SSOIssuer != "" && st.SSOClientID != "" {
		in.SSO = &sso.Source{Cache: sso.NewCache(env.Home), Issuer: st.SSOIssuer, ClientID: st.SSOClientID, Profile: st.Profile, HTTP: &http.Client{Transport: rt}, Now: env.Clock}
	}
	cfg := client.Config{
		Endpoint: st.EndpointURL, HTTP: rt, ReadTimeout: rto, Stderr: env.Stderr,
		UserAgent: buildinfo.UserAgent(), ContractVersion: spec.ContractVersion,
		Now: env.Now, Sleep: env.Sleep, Rand: env.Rand,
	}
	if !noAuth {
		cfg.Credentials = auth.Chain(in)
	}
	if g.Debug {
		cfg.Debug = env.Stderr
	}
	return &Session{Settings: st, Client: client.New(cfg), Source: src}, nil
}
