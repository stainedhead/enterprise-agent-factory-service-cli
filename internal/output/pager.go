package output

import (
	"context"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
)

// PagerOptions decide whether output goes through a pager (D-C22).
type PagerOptions struct {
	Out     io.Writer
	TTY     bool
	NoPager bool
	Getenv  func(string) string
	// Start starts the pager command; nil runs it with os/exec, writing to Out.
	Start func(ctx context.Context, command string, out io.Writer) (io.WriteCloser, func() error, error)
}

// PagerCommand returns the pager command: AFS_PAGER, else PAGER, else "less -R". An AFS_PAGER that
// is set but empty disables paging.
func PagerCommand(getenv func(string) string, afsPagerSet bool) (string, bool) {
	if v := getenv("AFS_PAGER"); v != "" {
		return v, true
	}
	if afsPagerSet {
		return "", false
	}
	if v := getenv("PAGER"); v != "" {
		return v, true
	}
	return "less -R", true
}

type nopCloser struct{ io.Writer }

func (nopCloser) Close() error { return nil }

// OpenPager returns the writer to use and a function to call when done (it closes the pager and
// waits for it). Without a terminal, with --no-cli-pager or with an empty AFS_PAGER the output is
// written directly.
func OpenPager(ctx context.Context, o PagerOptions, afsPagerSet bool) (io.WriteCloser, func() error) {
	direct := func() (io.WriteCloser, func() error) { return nopCloser{o.Out}, func() error { return nil } }
	if !o.TTY || o.NoPager {
		return direct()
	}
	cmd, ok := PagerCommand(o.Getenv, afsPagerSet)
	if !ok {
		return direct()
	}
	start := o.Start
	if start == nil {
		start = startProcess
	}
	w, wait, err := start(ctx, cmd, o.Out)
	if err != nil {
		return direct()
	}
	return w, func() error {
		_ = w.Close()
		return wait()
	}
}

func startProcess(ctx context.Context, command string, out io.Writer) (io.WriteCloser, func() error, error) {
	fields := strings.Fields(command)
	if len(fields) == 0 {
		return nil, nil, io.ErrClosedPipe
	}
	c := exec.CommandContext(ctx, fields[0], fields[1:]...) //nolint:gosec // the user's own pager command
	// on cancellation ask the pager to stop (it restores the terminal) rather than kill it
	c.Cancel = func() error { return c.Process.Signal(os.Interrupt) }
	c.WaitDelay = 2 * time.Second
	c.Stdout = out
	in, err := c.StdinPipe()
	if err != nil {
		return nil, nil, err
	}
	if err := c.Start(); err != nil {
		return nil, nil, err
	}
	return in, c.Wait, nil
}
