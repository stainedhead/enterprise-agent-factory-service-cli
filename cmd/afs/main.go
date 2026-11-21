// Command afs is the Agent Factory CLI. It follows AWS CLI conventions:
// "afs <service> <operation> [--options]".
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli/services"
)

func main() { os.Exit(run(os.Args[1:], cli.NewEnv(os.Stdin, os.Stdout, os.Stderr))) }

// run executes one invocation with a context cancelled by SIGINT (exit 130, CLI-8).
func run(args []string, env *cli.Env) int {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	reg, err := services.NewRegistry()
	if err != nil {
		_, _ = fmt.Fprintln(env.Stderr, "afs: internal error:", err)
		return cli.ExitOther
	}
	return runWith(ctx, args, env, reg)
}

// runWith maps an interrupted run to exit 130.
func runWith(ctx context.Context, args []string, env *cli.Env, reg *cli.Registry) int {
	code := cli.Run(ctx, args, env, reg)
	if ctx.Err() != nil && code != cli.ExitOK {
		return cli.ExitInterrupted
	}
	return code
}
