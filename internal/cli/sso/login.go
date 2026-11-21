package sso

import (
	"context"
	"fmt"
	"net/http"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli/exec"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli/runtime"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/config"
	issso "github.com/stainedhead/enterprise-agent-factory-service-cli/internal/sso"
)

func init() {
	registrations = append(registrations, func(r *cli.Registry) {
		r.AddLocal(&cli.Command{Service: "sso", Name: "login", Summary: "Sign in with single sign-on (device authorization) and cache the session.", Local: login,
			Override: cli.Override{Examples: []cli.Example{{Title: "Sign in with the default profile", Command: "afs sso login"}, {Title: "Sign in with a named profile", Command: "afs sso login --profile prod"}}}})
		r.AddLocal(&cli.Command{Service: "sso", Name: "logout", Summary: "Delete the cached single sign-on session of the profile.", Local: logout,
			Override: cli.Override{Examples: []cli.Example{{Title: "Sign out of the default profile", Command: "afs sso logout"}}}})
	})
}

func fail(inv *cli.Invocation, profile string, err error) int {
	return exec.ReportError(inv.Env.Stderr, "sso", profile, err)
}

func login(ctx context.Context, inv *cli.Invocation) int {
	env := inv.Env
	st, err := runtime.LenientSettings(env, inv.P.G)
	if err != nil {
		return fail(inv, "default", err)
	}
	issuer, clientID := st.SSOIssuer, st.SSOClientID
	if issuer == "" || clientID == "" {
		ps, err := Providers(ctx, env, inv.P.G)
		if err != nil {
			return fail(inv, st.Profile, err)
		}
		p, err := Pick(ps, st.SSOProviderID)
		if err != nil {
			return fail(inv, st.Profile, err)
		}
		issuer, clientID = p.IssuerURL, p.CLIClientID
		if _, err := config.WriteProfile(runtime.Source(env), st.Profile, map[string]string{"sso_provider_id": p.ID, "sso_issuer": issuer, "sso_client_id": clientID}, []string{"sso_provider_id", "sso_issuer", "sso_client_id"}); err != nil {
			return fail(inv, st.Profile, err)
		}
	}
	rt := env.RoundTripper
	if rt == nil {
		rt = http.DefaultTransport
	}
	tok, err := issso.Login(ctx, issso.LoginOptions{HTTP: &http.Client{Transport: rt}, Out: env.Stderr, Now: env.Clock, Sleep: env.SleepCtx}, issuer, clientID)
	if err != nil {
		return fail(inv, st.Profile, err)
	}
	if err := issso.NewCache(env.Home).Save(tok); err != nil {
		return fail(inv, st.Profile, err)
	}
	_, _ = fmt.Fprintf(env.Stdout, "Signed in. Session cached for profile %s.\n", st.Profile)
	return cli.ExitOK
}

func logout(_ context.Context, inv *cli.Invocation) int {
	env := inv.Env
	st, err := runtime.LenientSettings(env, inv.P.G)
	if err != nil {
		return fail(inv, "default", err)
	}
	if st.SSOIssuer == "" || st.SSOClientID == "" {
		_, _ = fmt.Fprintf(env.Stdout, "Profile %s has no single sign-on session.\n", st.Profile)
		return cli.ExitOK
	}
	if err := issso.NewCache(env.Home).Delete(st.SSOIssuer, st.SSOClientID); err != nil {
		return fail(inv, st.Profile, err)
	}
	_, _ = fmt.Fprintf(env.Stdout, "Signed out of profile %s.\n", st.Profile)
	return cli.ExitOK
}
