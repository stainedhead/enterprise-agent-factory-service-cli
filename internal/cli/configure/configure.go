package configure

import (
	"bufio"
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli/exec"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli/runtime"
	clisso "github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli/sso"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/config"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/output"
)

func init() {
	registrations = append(registrations, func(r *cli.Registry) {
		r.AddLocal(&cli.Command{Service: "configure", Summary: "Write a profile in ~/.afs/config. Use \"afs configure sso\" to choose a sign-in provider.", Local: run,
			Override: cli.Override{Examples: []cli.Example{
				{Title: "Set up the default profile", Command: "afs configure"},
				{Title: "Set up a named profile", Command: "afs configure --profile prod"},
				{Title: "Choose the sign-in provider of a profile", Command: "afs configure sso --profile prod"},
			}}})
	})
}

func run(ctx context.Context, inv *cli.Invocation) int {
	if len(inv.P.Positional) > 0 {
		if inv.P.Positional[0] == "sso" && len(inv.P.Positional) == 1 {
			return configureSSO(ctx, inv)
		}
		_, _ = fmt.Fprintf(inv.Env.Stderr, "afs: error: unknown argument %q: use afs configure or afs configure sso\n", inv.P.Positional[0])
		return cli.ExitUsage
	}
	env := inv.Env
	st, err := runtime.LenientSettings(env, inv.P.G)
	if err != nil {
		// a missing profile is what configure creates
		st = config.Settings{Profile: profileName(inv)}
	}
	in := bufio.NewReader(env.Stdin)
	ask := func(label, cur string) string {
		if cur != "" {
			_, _ = fmt.Fprintf(env.Stderr, "%s [%s]: ", label, cur)
		} else {
			_, _ = fmt.Fprintf(env.Stderr, "%s [None]: ", label)
		}
		line, _ := in.ReadString('\n')
		if v := strings.TrimSpace(line); v != "" {
			return v
		}
		return cur
	}
	endpoint := ask("AFS endpoint URL", st.EndpointURL)
	if endpoint != "" {
		if err := config.ValidateEndpoint(endpoint); err != nil {
			return exec.ReportError(env.Stderr, "configure", st.Profile, err)
		}
	}
	outFmt := ask("Default output format (json, yaml, table, text, csv)", orDefault(st.Output, "json"))
	if !contains(cli.OutputFormats, outFmt) {
		_, _ = fmt.Fprintf(env.Stderr, "afs: error: invalid output format %q (choose from %s)\n", outFmt, strings.Join(cli.OutputFormats, ", "))
		return cli.ExitUsage
	}
	team := ask("Default team id (optional)", st.TeamID)
	path, err := config.WriteProfile(runtime.Source(env), st.Profile, map[string]string{"endpoint_url": endpoint, "output": outFmt, "team_id": team}, []string{"endpoint_url", "output", "team_id"})
	if err != nil {
		return exec.ReportError(env.Stderr, "configure", st.Profile, err)
	}
	_, _ = fmt.Fprintf(env.Stdout, "Profile %s written to %s\n", st.Profile, path)
	return cli.ExitOK
}

func profileName(inv *cli.Invocation) string {
	n, _ := config.ProfileName(config.Overrides{Profile: inv.P.G.Profile}, runtime.Source(inv.Env))
	return n
}

func orDefault(s, d string) string {
	if s == "" {
		return d
	}
	return s
}

func contains(l []string, s string) bool {
	for _, x := range l {
		if x == s {
			return true
		}
	}
	return false
}

func configureSSO(ctx context.Context, inv *cli.Invocation) int {
	env := inv.Env
	name := profileName(inv)
	st, err := runtime.LenientSettings(env, inv.P.G)
	if err != nil {
		st = config.Settings{Profile: name}
	}
	ps, err := clisso.Providers(ctx, env, inv.P.G)
	if err != nil {
		return exec.ReportError(env.Stderr, "ListSignInProviders", st.Profile, err)
	}
	if len(ps) == 0 {
		_, _ = fmt.Fprintln(env.Stderr, "afs: error: the service has no sign-in providers enabled")
		return cli.ExitOther
	}
	chosen := ps[0]
	if len(ps) > 1 {
		for i, p := range ps {
			def := ""
			if p.IsDefault {
				def = " (default)"
			}
			_, _ = fmt.Fprintf(env.Stderr, "%d) %s [%s]%s\n", i+1, output.Sanitize(p.Name), output.Sanitize(p.Kind), def)
		}
		_, _ = fmt.Fprint(env.Stderr, "Choose a provider by number: ")
		line, _ := bufio.NewReader(env.Stdin).ReadString('\n')
		n, convErr := strconv.Atoi(strings.TrimSpace(line))
		switch {
		case convErr == nil && n >= 1 && n <= len(ps):
			chosen = ps[n-1]
		case strings.TrimSpace(line) == "":
			p, err := clisso.Pick(ps, "")
			if err != nil {
				_, _ = fmt.Fprintf(env.Stderr, "afs: error: %v\n", err)
				return cli.ExitUsage
			}
			chosen = p
		default:
			_, _ = fmt.Fprintf(env.Stderr, "afs: error: %q is not a number between 1 and %d\n", strings.TrimSpace(line), len(ps))
			return cli.ExitUsage
		}
	}
	if err := chosen.Validate(); err != nil {
		return exec.ReportError(env.Stderr, "configure", st.Profile, err)
	}
	kv := map[string]string{"sso_provider_id": chosen.ID, "sso_issuer": chosen.IssuerURL, "sso_client_id": chosen.CLIClientID}
	if st.EndpointURL != "" && st.Sources["endpoint_url"] != "profile" {
		kv["endpoint_url"] = st.EndpointURL
	}
	path, err := config.WriteProfile(runtime.Source(env), st.Profile, kv, []string{"endpoint_url", "sso_provider_id", "sso_issuer", "sso_client_id"})
	if err != nil {
		return exec.ReportError(env.Stderr, "configure", st.Profile, err)
	}
	_, _ = fmt.Fprintf(env.Stdout, "Profile %s now signs in with %s. Run afs sso login --profile %s. (%s)\n", st.Profile, chosen.Name, st.Profile, path)
	return cli.ExitOK
}
