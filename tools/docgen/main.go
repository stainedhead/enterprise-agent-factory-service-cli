// Command docgen writes the command reference in user-docs/reference from the command table and
// the help text, so the documentation cannot drift from the commands (D-C30).
package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli/help"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli/services"
)

func helpText(reg *cli.Registry, words ...string) (string, error) {
	var out, eb bytes.Buffer
	env := cli.NewEnv(strings.NewReader(""), &out, &eb)
	env.Vars = map[string]string{}
	if code := cli.Run(context.Background(), append(append([]string{}, words...), "help"), env, reg); code != 0 {
		return "", fmt.Errorf("help %v: exit %d: %s", words, code, eb.String())
	}
	return out.String(), nil
}

// render returns the reference files by name.
func render(reg *cli.Registry) (map[string]string, error) {
	files := map[string]string{}
	var index strings.Builder
	index.WriteString("# Command reference\n\nGenerated from `afs help` by `make docs`; do not edit by hand.\n\n| Service | What it does |\n|---|---|\n")
	for _, svc := range reg.Services() {
		fmt.Fprintf(&index, "| [%s](%s.md) | %s |\n", svc, svc, help.ServiceSummaries[svc])
		var b strings.Builder
		fmt.Fprintf(&b, "# afs %s\n\nGenerated from `afs %s help`; do not edit by hand.\n\n", svc, svc)
		top, err := helpText(reg, svc)
		if err != nil {
			return nil, err
		}
		fmt.Fprintf(&b, "```text\n%s```\n", top)
		for _, c := range reg.Commands(svc) {
			if c.Name == "" {
				continue
			}
			text, err := helpText(reg, svc, c.Name)
			if err != nil {
				return nil, err
			}
			fmt.Fprintf(&b, "\n## %s\n\n```text\n%s```\n", c.Name, text)
		}
		files[svc+".md"] = b.String()
	}
	files["README.md"] = index.String()
	return files, nil
}

func main() {
	dir := flag.String("out", "user-docs/reference", "output directory")
	flag.Parse()
	reg, err := services.NewRegistry()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	files, err := render(reg)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := os.MkdirAll(*dir, 0o755); err != nil { //nolint:gosec // documentation
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	for name, body := range files {
		if err := os.WriteFile(filepath.Join(*dir, name), []byte(body), 0o644); err != nil { //nolint:gosec // documentation
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}
