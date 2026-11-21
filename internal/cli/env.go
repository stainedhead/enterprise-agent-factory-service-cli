package cli

import (
	"context"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

// Env is everything a command touches outside its arguments, injected so tests run without a
// network, a real home directory or a terminal.
type Env struct {
	Stdin          io.Reader
	Stdout, Stderr io.Writer
	// Vars is the process environment. Never nil.
	Vars map[string]string
	// Home is the user's home directory ("" when unknown).
	Home string
	// StdinTTY and StdoutTTY say whether the streams are terminals.
	StdinTTY, StdoutTTY bool
	// Now, Sleep and Rand are the clock, the cancellable sleeper and the jitter source.
	Now   func() time.Time
	Sleep func(ctx context.Context, d time.Duration) error
	Rand  func() float64
	// RoundTripper replaces the default HTTP transport (tests).
	RoundTripper http.RoundTripper
	// Dial replaces net.Dialer (tests).
	Dial func(ctx context.Context, network, addr string) (net.Conn, error)
	// LookPath finds a program on PATH.
	LookPath func(string) (string, error)
	// StartPager starts the pager command writing to out; nil uses os/exec.
	StartPager func(command string, out io.Writer) (io.WriteCloser, func() error, error)
}

// Getenv returns the value of an environment variable.
func (e *Env) Getenv(key string) string { return e.Vars[key] }

// LookupEnv returns a variable and whether it is set.
func (e *Env) LookupEnv(key string) (string, bool) {
	v, ok := e.Vars[key]
	return v, ok
}

// SleepCtx sleeps with the injected sleeper, or really.
func (e *Env) SleepCtx(ctx context.Context, d time.Duration) error {
	if e.Sleep != nil {
		return e.Sleep(ctx, d)
	}
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
		return nil
	}
}

// Clock returns the current time from the injected clock.
func (e *Env) Clock() time.Time {
	if e.Now != nil {
		return e.Now()
	}
	return time.Now()
}

// NewEnv builds the real environment for the process; terminal flags are filled by the caller.
func NewEnv(stdin io.Reader, stdout, stderr io.Writer) *Env {
	vars := map[string]string{}
	for _, kv := range os.Environ() {
		if k, v, ok := strings.Cut(kv, "="); ok {
			vars[k] = v
		}
	}
	home, _ := os.UserHomeDir()
	return &Env{Stdin: stdin, Stdout: stdout, Stderr: stderr, Vars: vars, Home: home, LookPath: exec.LookPath}
}
