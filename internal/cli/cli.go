// Package cli is the afs command line: registry, parsing, global options and dispatch.
// Execution of API commands is wired in by internal/cli/services (generic executor in
// internal/cli/exec); local commands live in their service packages.
package cli

import (
	"context"
	"errors"
	"fmt"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/buildinfo"
)

// Run executes the CLI with args (excluding the program name) and returns the process exit code.
func Run(ctx context.Context, args []string, env *Env, reg *Registry) int {
	if len(args) == 0 {
		printf(env, "%s", rootUsage)
		return ExitUsage
	}
	p, cmd, err := reg.Parse(args)
	if err != nil {
		var ue *UsageError
		if errors.As(err, &ue) {
			printf2(env, "%s\nafs: error: %s\n", ue.Usage, ue.Msg)
			return ExitUsage
		}
		printf2(env, "afs: error: %v\n", err)
		return ExitOther
	}
	if p.G.Version {
		_, _ = fmt.Fprintln(env.Stdout, buildinfo.Line())
		return ExitOK
	}
	inv := &Invocation{Env: env, Cmd: cmd, P: p, Reg: reg}
	if p.Help || cmd == nil {
		if reg.Help == nil {
			_, _ = fmt.Fprint(env.Stdout, rootUsage)
			return ExitOK
		}
		return reg.Help(inv, p.Words)
	}
	if cmd.Local != nil {
		return cmd.Local(ctx, inv)
	}
	if reg.Exec == nil {
		printf2(env, "afs: error: no executor wired\n")
		return ExitOther
	}
	return reg.Exec(ctx, inv)
}

func printf(env *Env, format string, a ...any) { _, _ = fmt.Fprintf(env.Stdout, format, a...) }

func printf2(env *Env, format string, a ...any) { _, _ = fmt.Fprintf(env.Stderr, format, a...) }
