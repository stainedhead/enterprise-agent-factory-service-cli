package exec

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli/runtime"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli/safeupdate"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli/wait"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/client"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/output"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/spec"
)

// WaitOptions are the --wait and --wait-timeout options added to every operation in wait.Table.
func WaitOptions() []cli.Option {
	return []cli.Option{
		{Flag: "--wait", In: "local", Kind: cli.Bool, Desc: "Wait until the resource is ready (or gone, for a delete)."},
		{Flag: "--wait-timeout", In: "local", Kind: cli.Duration, Desc: "How long --wait waits, for example 90s or 10m (default 10m)."},
	}
}

func usageExit(env *cli.Env, err error) int {
	var ue *cli.UsageError
	if errors.As(err, &ue) {
		_, _ = fmt.Fprintf(env.Stderr, "%s\nafs: error: %s\n", ue.Usage, ue.Msg)
	} else {
		_, _ = fmt.Fprintf(env.Stderr, "\nafs: error: %v\n", err)
	}
	return cli.ExitUsage
}

// Run executes one API command: build, send (paginating, reading versions, waiting), render.
func Run(ctx context.Context, inv *cli.Invocation) int {
	env, cmd, g := inv.Env, inv.Cmd, inv.P.G
	op := cmd.Op
	if g.GenerateSkeleton {
		raw, _ := json.Marshal(Skeleton(op))
		if err := output.WriteJSON(env.Stdout, raw); err != nil {
			return cli.ExitOther
		}
		return cli.ExitOK
	}
	st, err := runtime.Settings(env, g)
	if err != nil {
		return ReportError(env.Stderr, op.ID, "default", err)
	}
	b, err := Build(cmd, inv.P, Sources{Stdin: env.Stdin, ReadPassword: env.ReadPassword}, st.TeamID)
	if err != nil {
		return usageExit(env, err)
	}
	format := st.Output
	if err := checkPaging(cmd, g, format); err != nil {
		return usageExit(env, err)
	}
	sess, err := runtime.New(env, g, op.Scope == "public")
	if err != nil {
		return ReportError(env.Stderr, op.ID, st.Profile, err)
	}
	c := sess.Client
	fail := func(err error) int {
		return ReportError(env.Stderr, op.ID, st.Profile, currentVersionOn412(ctx, c, op, b.Path, err))
	}

	if op.IfMatch != spec.IfMatchNone && b.IfMatch == "" {
		v, ok, err := safeupdate.Version(ctx, c, op, b.Path)
		if err != nil {
			return fail(err)
		}
		if ok {
			b.IfMatch = v
			_, _ = fmt.Fprintln(env.Stderr, safeupdate.Notice(v))
		}
	}

	var body []byte
	if b.Body != nil {
		if body, err = json.Marshal(b.Body); err != nil {
			return fail(err)
		}
	}
	accept := ""
	if format == "csv" && op.CSV {
		accept = "text/csv"
	}
	if cmd.Handle != nil && cmd.Accept != "" {
		accept = cmd.Accept
	}
	send := func(ctx context.Context, q url.Values) (*client.Response, error) {
		return c.Do(ctx, client.Request{
			Method: op.Method, Path: client.BuildPath(op.Path, b.Path), Query: q, Body: body, Accept: accept,
			IfMatch: b.IfMatch, Idempotent: op.Idempotent, IdempotencyKey: b.IdemKey, NoAuth: op.Scope == "public", Operation: op.ID,
		})
	}

	// the output sink
	_, pagerSet := env.LookupEnv("AFS_PAGER")
	sink, closeSink := output.OpenPager(output.PagerOptions{Out: env.Stdout, TTY: env.StdoutTTY, NoPager: g.NoCLIPager, Getenv: env.Getenv, Start: env.StartPager}, pagerSet)
	defer func() { _ = closeSink() }()

	var resp *client.Response
	var data []byte
	switch {
	case accept == "text/csv":
		pg := &pages{do: send, op: op, g: g, base: b.Query, errw: env.Stderr}
		resp, err = pg.fetchCSV(ctx, &output.CSVWriter{W: sink, Resume: g.StartingToken != ""})
		if err != nil {
			return fail(err)
		}
		return finish(inv, resp)
	case op.Paginated:
		pg := &pages{do: send, op: op, g: g, base: b.Query}
		data, resp, err = pg.fetchJSON(ctx)
		var uf *usageFail
		if errors.As(err, &uf) {
			return usageExit(env, uf.error)
		}
	case cmd.Handle != nil:
		resp, err = send(ctx, b.Query)
		if err != nil {
			return fail(err)
		}
		return cmd.Handle(inv, resp)
	default:
		resp, err = send(ctx, b.Query)
		if resp != nil {
			data = resp.Body
		}
	}
	if err != nil {
		return fail(err)
	}
	if resp.Replayed {
		_, _ = fmt.Fprintln(env.Stderr, "The service answered from the stored response of an earlier request with this Idempotency-Key.")
	}

	if t, ok := wait.Table[op.ID]; ok && b.Wait {
		poll, get := wait.Paths(t, b.Path, resp.Body)
		final, err := wait.Until(ctx, c, t, resp.Body, poll, get, wait.Options{Timeout: b.WaitTimeout, Now: env.Now, Sleep: env.Sleep, Interval: 0})
		if err != nil {
			if errors.Is(err, context.Canceled) {
				_, _ = fmt.Fprintln(env.Stderr, "Interrupted while waiting: the request was sent and the change stands.")
			}
			if errors.Is(err, wait.ErrTimeout) {
				return fail(fmt.Errorf("%w; the request was sent and the change stands", err))
			}
			return fail(err)
		}
		if final != nil {
			data = final
		} else {
			data = nil
		}
	}

	if len(bytes.TrimSpace(data)) > 0 && resp.Status != 204 {
		itemsKey := op.ResultKey
		if cmd.ResultKey != "" {
			itemsKey = cmd.ResultKey
		}
		ro := output.RenderOptions{Format: format, Query: g.Query, Options: output.Options{
			Columns: cmd.Columns, ItemsKey: itemsKey,
			Color: output.ColorEnabled(g.Color, env.StdoutTTY, env.Getenv),
		}}
		if err := output.Render(sink, data, ro); err != nil {
			return usageExit(env, err)
		}
	}
	return finish(inv, resp)
}

// currentVersionOn412 makes a 412 carry the resource's current version in its ETag header: the
// response's own when present, else one more read (D-C8). Any failure leaves err unchanged.
func currentVersionOn412(ctx context.Context, c safeupdate.Doer, op *spec.Operation, path map[string]string, err error) error {
	var ae *client.APIError
	if !errors.As(err, &ae) || ae.Status != 412 || ae.Header.Get("ETag") != "" {
		return err
	}
	if v, ok, rerr := safeupdate.Version(ctx, c, op, path); rerr == nil && ok {
		if ae.Header == nil {
			ae.Header = http.Header{}
		}
		ae.Header.Set("ETag", `"`+v+`"`)
	}
	return err
}

// finish runs the command's After hook.
func finish(inv *cli.Invocation, resp *client.Response) int {
	if inv.Cmd.After != nil {
		if code := inv.Cmd.After(inv, resp); code >= 0 {
			return code
		}
	}
	return cli.ExitOK
}
