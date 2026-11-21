// Package help renders the three levels of help (afs help, afs <service> help, afs <service>
// <operation> help) from the same definitions the commands run from (CLI-14).
package help

import (
	"fmt"
	"io"
	"strings"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
)

// ServiceSummaries describe each service in one line.
var ServiceSummaries = map[string]string{
	"admin":           "Service administration: settings, templates, tags, channels, storage and security operators.",
	"agents":          "Create and operate agents, their identities, configuration and external systems.",
	"api-keys":        "Create and revoke API keys.",
	"audit":           "Search and export the audit log.",
	"budgets":         "Set budgets and read spend.",
	"completion":      "Print a shell completion script.",
	"configure":       "Create and edit profiles in ~/.afs/config.",
	"members":         "Manage team members and role bindings.",
	"models":          "List and manage models and their prices.",
	"pools":           "Manage worker pools.",
	"reports":         "Run and schedule reports.",
	"secrets":         "Manage secret references and write secret values.",
	"sso":             "Sign in with single sign-on and show who you are.",
	"swarm-templates": "Manage swarm templates.",
	"teams":           "Create and manage teams.",
	"work-orders":     "Manage work orders and their repositories.",
}

var kindNames = map[cli.Kind]string{
	cli.String: "string", cli.Int: "integer", cli.Number: "number", cli.Bool: "boolean",
	cli.List: "list", cli.JSON: "JSON", cli.Any: "string or JSON", cli.Duration: "duration",
}

func rootHelp(w io.Writer, reg *cli.Registry) {
	p := func(f string, a ...any) { _, _ = fmt.Fprintf(w, f, a...) }
	p("AFS(1)\n\nNAME\n       afs - the Enterprise Agent Factory Service command line\n\n")
	p("SYNOPSIS\n       afs [options] <service> <operation> [parameters]\n\n")
	p("SERVICES\n")
	for _, s := range reg.Services() {
		p("       %-16s %s\n", s, ServiceSummaries[s])
	}
	p("\nGLOBAL OPTIONS\n")
	globals(w)
	p("\nTo see help text for a service, run: afs <service> help\n")
}

func globals(w io.Writer) {
	for _, g := range cli.GlobalHelp() {
		_, _ = fmt.Fprintf(w, "       %-22s %s\n", g[0], g[1])
	}
}

func serviceHelp(w io.Writer, reg *cli.Registry, service string) {
	p := func(f string, a ...any) { _, _ = fmt.Fprintf(w, f, a...) }
	p("AFS %s\n\nDESCRIPTION\n       %s\n\nSYNOPSIS\n       afs %s <operation> [parameters]\n\nOPERATIONS\n", strings.ToUpper(service), ServiceSummaries[service], service)
	for _, c := range reg.Commands(service) {
		if c.Name == "" {
			continue
		}
		sum := c.Summary
		if c.Op != nil {
			sum = c.Op.Summary
		}
		p("       %-30s %s\n", c.Name, sum)
	}
	p("\nTo see help text for an operation, run: afs %s <operation> help\n", service)
}

func operationHelp(w io.Writer, c *cli.Command) {
	p := func(f string, a ...any) { _, _ = fmt.Fprintf(w, f, a...) }
	sum := c.Summary
	if c.Op != nil {
		sum = c.Op.Summary
	}
	p("AFS %s\n\nNAME\n       %s\n\nDESCRIPTION\n       %s\n", strings.ToUpper(strings.ReplaceAll(c.Key(), " ", " ")), c.Key(), sum)
	if c.Op != nil {
		if c.Op.Permission != "" && c.Op.Permission != "none" {
			p("\n       Required permission: %s\n", c.Op.Permission)
		} else if len(c.Op.AnyOf) > 0 {
			p("\n       Required permission (any of): %s\n", strings.Join(c.Op.AnyOf, ", "))
		}
	}
	p("\nSYNOPSIS\n       afs %s [options]\n", c.Key())
	if len(c.Options) > 0 {
		p("\nOPTIONS\n")
		for _, o := range c.Options {
			req := ""
			if o.Required {
				req = " (required)"
			}
			p("       %s (%s)%s\n", o.Flag, kindNames[o.Kind], req)
			if o.Desc != "" {
				p("          %s\n", o.Desc)
			}
			if len(o.Enum) > 0 {
				p("          Possible values: %s\n", strings.Join(o.Enum, ", "))
			}
			if o.Min != nil || o.Max != nil {
				p("          Range: %s to %s\n", bound(o.Min), bound(o.Max))
			}
			if o.Sensitive {
				p("          Give file://PATH, fileb://PATH or - (standard input); never the value itself.\n")
			}
		}
	}
	p("\nGLOBAL OPTIONS\n")
	globals(w)
	p("\nEXAMPLES\n")
	if len(c.Examples) == 0 {
		p("       afs %s help\n", c.Key())
	}
	for _, e := range c.Examples {
		p("       %s\n\n       %s\n\n", e.Title, e.Command)
	}
}

func bound(f *float64) string {
	if f == nil {
		return "unbounded"
	}
	return fmt.Sprintf("%v", *f)
}

// Run renders help for the words given: none, a service, or a service and an operation.
func Run(inv *cli.Invocation, words []string) int {
	w := inv.Env.Stdout
	ws := make([]string, 0, len(words))
	for _, x := range words {
		if x != "help" {
			ws = append(ws, x)
		}
	}
	switch {
	case len(ws) == 0:
		rootHelp(w, inv.Reg)
	case len(ws) == 1:
		if len(inv.Reg.Commands(ws[0])) == 0 {
			_, _ = fmt.Fprintf(inv.Env.Stderr, "afs: error: argument service: Invalid choice: %q\n", ws[0])
			return cli.ExitUsage
		}
		serviceHelp(w, inv.Reg, ws[0])
	default:
		c, ok := inv.Reg.Lookup(ws[0], ws[1])
		if !ok {
			_, _ = fmt.Fprintf(inv.Env.Stderr, "afs: error: argument operation: Invalid choice: %q for service %s\n", ws[1], ws[0])
			return cli.ExitUsage
		}
		operationHelp(w, c)
	}
	return cli.ExitOK
}
