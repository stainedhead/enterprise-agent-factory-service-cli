package main

import (
	"bytes"
	"context"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/spec"
)

func TestVersionRuns(t *testing.T) {
	var out, errb bytes.Buffer
	env := cli.NewEnv(strings.NewReader(""), &out, &errb)
	if code := run([]string{"--version"}, env); code != 0 || !strings.HasPrefix(out.String(), "afs-cli/") {
		t.Fatalf("code %d out %q err %q", code, out.String(), errb.String())
	}
}

// CLI-US-5.3: SIGINT cancels the context and the process exits 130.
func TestSigintExit130(t *testing.T) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	reg := cli.NewRegistry(spec.All())
	reg.AddLocal(&cli.Command{Service: "wait", Local: func(ctx context.Context, inv *cli.Invocation) int {
		<-ctx.Done()
		return cli.ExitOther
	}})
	if err := reg.Finalize(); err != nil {
		t.Fatal(err)
	}
	go func() {
		time.Sleep(50 * time.Millisecond)
		_ = syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	}()
	var out, errb bytes.Buffer
	env := cli.NewEnv(strings.NewReader(""), &out, &errb)
	if code := runWith(ctx, []string{"wait"}, env, reg); code != cli.ExitInterrupted {
		t.Fatalf("exit %d, want 130", code)
	}
}
