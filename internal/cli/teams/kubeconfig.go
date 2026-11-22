package teams

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/client"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/kubeconfig"
)

func init() { registrations = append(registrations, registerKubeconfig) }

func registerKubeconfig(r *cli.Registry) {
	r.Override(cli.Override{
		Op: "GetTeamKubeconfig", Dynamic: true, // the answer is a kubeconfig, not a table
		Accept: "application/yaml", Handle: handleKubeconfig,
		Extra: []cli.Option{
			{Flag: "--kubeconfig", In: "local", Kind: cli.String, Desc: "Kubeconfig file to merge into (default: the first path of $KUBECONFIG, else ~/.kube/config); - prints the API's YAML to stdout."},
			{Flag: "--aws-profile", In: "local", Kind: cli.String, Desc: "Set AWS_PROFILE to this profile in the exec credential plugin."},
		},
		Examples: []cli.Example{
			{Title: "Merge the team's context into ~/.kube/config", Command: "afs teams get-kubeconfig --team-id " + exTeam},
			{Title: "Write a separate file and use an AWS profile for the token", Command: "afs teams get-kubeconfig --team-id " + exTeam + " --kubeconfig ./kc.yaml --aws-profile ops"},
			{Title: "Print the YAML only", Command: "afs teams get-kubeconfig --team-id " + exTeam + " --kubeconfig -"},
		},
	})
}

// kubeconfigPath is the --kubeconfig value, else the first path of $KUBECONFIG, else ~/.kube/config.
func kubeconfigPath(inv *cli.Invocation) (string, error) {
	if p := inv.P.Value("kubeconfig"); p != "" {
		return p, nil
	}
	if kc := inv.Env.Getenv("KUBECONFIG"); kc != "" {
		if first := strings.Split(kc, string(os.PathListSeparator))[0]; first != "" {
			return first, nil
		}
	}
	if inv.Env.Home == "" {
		return "", fmt.Errorf("cannot find the home directory; pass --kubeconfig PATH")
	}
	return filepath.Join(inv.Env.Home, ".kube", "config"), nil
}

// handleKubeconfig merges the returned kubeconfig into the user's file (CLI-16), or prints it.
func handleKubeconfig(inv *cli.Invocation, resp *client.Response) int {
	env := inv.Env
	if inv.P.Value("kubeconfig") == "-" {
		_, _ = env.Stdout.Write(resp.Body)
		return cli.ExitOK
	}
	path, err := kubeconfigPath(inv)
	if err != nil {
		_, _ = fmt.Fprintf(env.Stderr, "\nafs: error: %v\n", err)
		return cli.ExitUsage
	}
	res, err := kubeconfig.Apply(path, resp.Body, kubeconfig.Options{AWSProfile: inv.P.Value("aws-profile")})
	if err != nil {
		_, _ = fmt.Fprintf(env.Stderr, "\nafs: error: %v\n", err)
		return cli.ExitOther
	}
	if res.Wide {
		_, _ = fmt.Fprintf(env.Stderr, "warning: %s has mode %04o, wider than 0600; the mode was kept. Run chmod 600 %s to restrict it.\n", path, res.Mode, path)
	}
	if env.LookPath != nil {
		if _, err := env.LookPath("aws"); err != nil {
			_, _ = fmt.Fprintln(env.Stderr, "warning: the aws CLI (v2) was not found on PATH; kubectl runs `aws eks get-token` to authenticate.")
		}
	}
	verb := "Updated"
	if res.Created {
		verb = "Created"
	}
	_, _ = fmt.Fprintf(env.Stdout, "%s %s with context %s (current-context set).\n", verb, path, res.Context)
	return cli.ExitOK
}
