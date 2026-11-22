package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/spec"
)

func testEnv() (*Env, *bytes.Buffer, *bytes.Buffer) {
	var out, errb bytes.Buffer
	return &Env{Stdin: strings.NewReader(""), Stdout: &out, Stderr: &errb, Vars: map[string]string{}}, &out, &errb
}

func testRegistry(t *testing.T) *Registry {
	t.Helper()
	r := NewRegistry(spec.All())
	if err := r.Finalize(); err != nil {
		t.Fatal(err)
	}
	return r
}

type captured struct{ inv *Invocation }

func withExec(r *Registry, c *captured) {
	r.Exec = func(ctx context.Context, inv *Invocation) int { c.inv = inv; return ExitOK }
}

// CLI-2: global options are accepted before and after the operation, with = or a separate value.
func TestGlobalOptionsAnywhere(t *testing.T) {
	cases := [][]string{
		{"--profile", "p2", "--output", "table", "teams", "list-teams"},
		{"teams", "list-teams", "--profile", "p2", "--output", "table"},
		{"--profile=p2", "teams", "--output=table", "list-teams"},
	}
	for _, args := range cases {
		r := testRegistry(t)
		var c captured
		withExec(r, &c)
		env, _, errb := testEnv()
		if code := Run(context.Background(), args, env, r); code != ExitOK {
			t.Fatalf("%v: exit %d: %s", args, code, errb)
		}
		g := c.inv.P.G
		if g.Profile != "p2" || g.Output != "table" || c.inv.Cmd.Key() != "teams list-teams" {
			t.Errorf("%v: parsed %+v cmd %s", args, g, c.inv.Cmd.Key())
		}
	}
}

func TestAllGlobalOptionsParse(t *testing.T) {
	r := testRegistry(t)
	var c captured
	withExec(r, &c)
	env, _, errb := testEnv()
	args := []string{"teams", "list-teams", "--no-paginate", "--max-items", "10", "--page-size", "20", "--starting-token", "tok",
		"--cli-input-json", "{}", "--generate-cli-skeleton", "--debug", "--no-cli-pager", "--color", "off",
		"--cli-connect-timeout", "5", "--cli-read-timeout", "7", "--query", "items", "--endpoint-url", "http://x"}
	if code := Run(context.Background(), args, env, r); code != ExitOK {
		t.Fatalf("exit %d: %s", code, errb)
	}
	g := c.inv.P.G
	if !g.NoPaginate || g.MaxItems != 10 || g.PageSize != 20 || g.StartingToken != "tok" || g.CLIInputJSON != "{}" ||
		!g.GenerateSkeleton || !g.Debug || !g.NoCLIPager || g.Color != "off" || g.ConnectTimeout != 5 || g.ReadTimeout != 7 ||
		g.Query != "items" || g.EndpointURL != "http://x" || !g.Set["page-size"] {
		t.Errorf("globals = %+v", g)
	}
}

// CLI-US-5.1: unknown services, operations and options exit 252.
func TestUnknownServiceAndOptionExit252(t *testing.T) {
	for _, args := range [][]string{
		{"nope", "list-things"},
		{"teams", "nope"},
		{"teams", "list-teams", "--replicas", "3"},
		{"--unknown", "teams", "list-teams"},
		{"teams", "list-teams", "extra"},
		{"teams", "list-teams", "--output", "xml"},
		{"teams", "list-teams", "--color", "blue"},
		{"teams", "list-teams", "--max-items", "abc"},
		{"teams", "list-teams", "--profile"},
		{"teams", "describe-team", "--team-id"},
		{"teams", "describe-team", "--team-id", "--debug"},
		{},
	} {
		r := testRegistry(t)
		var c captured
		withExec(r, &c)
		env, _, errb := testEnv()
		if code := Run(context.Background(), args, env, r); code != ExitUsage {
			t.Errorf("%v: exit %d, want 252 (stderr %q)", args, code, errb)
		}
		if c.inv != nil {
			t.Errorf("%v: executed", args)
		}
	}
}

// CLI-US-5.1 / 5.1-global-options: --version prints the contract line and makes no call.
func TestVersionString(t *testing.T) {
	r := testRegistry(t)
	var c captured
	withExec(r, &c)
	env, out, _ := testEnv()
	if code := Run(context.Background(), []string{"--version"}, env, r); code != ExitOK {
		t.Fatal(code)
	}
	if !strings.HasPrefix(out.String(), "afs-cli/") || !strings.Contains(out.String(), " contract/1.3 go") || c.inv != nil {
		t.Errorf("version = %q", out)
	}
}

func TestParseOptionKinds(t *testing.T) {
	r := testRegistry(t)
	cases := []struct {
		args []string
		name string
		want []string
	}{
		{[]string{"budgets", "put-budget", "--thresholds", "50", "80", "100", "--period", "monthly"}, "thresholds", []string{"50", "80", "100"}},
		{[]string{"budgets", "put-budget", "--period=monthly"}, "period", []string{"monthly"}},
		{[]string{"teams", "update-team", "--no-viewer-kubectl"}, "viewer-kubectl", []string{"false"}},
		{[]string{"teams", "update-team", "--viewer-kubectl"}, "viewer-kubectl", []string{"true"}},
		{[]string{"agents", "create-agent", "--volume", `{"sizeGiB":20}`}, "volume", []string{`{"sizeGiB":20}`}},
	}
	for _, c := range cases {
		p, cmd, err := r.Parse(c.args)
		if err != nil || cmd == nil {
			t.Fatalf("%v: %v", c.args, err)
		}
		got := p.Values[c.name]
		if strings.Join(got, "|") != strings.Join(c.want, "|") {
			t.Errorf("%v: %s = %v, want %v", c.args, c.name, got, c.want)
		}
	}
	if _, _, err := r.Parse([]string{"teams", "update-team", "--viewer-kubectl=true"}); err == nil {
		t.Error("bool with value accepted")
	}
	if _, _, err := r.Parse([]string{"budgets", "put-budget", "--thresholds"}); err == nil {
		t.Error("empty list accepted")
	}
}

func TestHelpWords(t *testing.T) {
	r := testRegistry(t)
	for _, args := range [][]string{{"help"}, {"teams", "help"}, {"teams", "list-teams", "help"}, {"teams", "list-teams", "--help"}, {"-h"}} {
		p, _, err := r.Parse(args)
		if err != nil || !p.Help {
			t.Errorf("%v: help=%v err=%v", args, p != nil && p.Help, err)
		}
	}
	env, out, _ := testEnv()
	if code := Run(context.Background(), []string{"help"}, env, r); code != ExitOK || !strings.Contains(out.String(), "usage: afs") {
		t.Errorf("fallback help: %d %q", code, out)
	}
}

func TestLocalCommandPositional(t *testing.T) {
	r := NewRegistry(spec.All())
	var got []string
	r.AddLocal(&Command{Service: "completion", Local: func(ctx context.Context, inv *Invocation) int { got = inv.P.Positional; return 7 }})
	r.AddLocal(&Command{Service: "completion", Local: nil})
	if err := r.Finalize(); err == nil {
		t.Error("duplicate local command not reported")
	}
	r = NewRegistry(spec.All())
	r.AddLocal(&Command{Service: "completion", Local: func(ctx context.Context, inv *Invocation) int { got = inv.P.Positional; return 7 }})
	r.AddLocal(&Command{Service: "configure", Name: "sso", Local: func(ctx context.Context, inv *Invocation) int { return 8 }})
	if err := r.Finalize(); err != nil {
		t.Fatal(err)
	}
	env, _, _ := testEnv()
	if code := Run(context.Background(), []string{"completion", "bash"}, env, r); code != 7 || len(got) != 1 || got[0] != "bash" {
		t.Errorf("completion: code %d positional %v", code, got)
	}
	if code := Run(context.Background(), []string{"configure", "sso"}, env, r); code != 8 {
		t.Errorf("configure sso: %d", code)
	}
}

func TestRunWithoutExecutor(t *testing.T) {
	r := testRegistry(t)
	env, _, errb := testEnv()
	if code := Run(context.Background(), []string{"teams", "list-teams"}, env, r); code != ExitOther || !strings.Contains(errb.String(), "no executor") {
		t.Errorf("code %d %q", code, errb)
	}
}

func TestEnvHelpers(t *testing.T) {
	env, _, _ := testEnv()
	env.Vars["A"] = "1"
	if env.Getenv("A") != "1" {
		t.Error("Getenv")
	}
	if _, ok := env.LookupEnv("B"); ok {
		t.Error("LookupEnv")
	}
	if env.Clock().IsZero() {
		t.Error("Clock")
	}
	if err := env.SleepCtx(context.Background(), time.Millisecond); err != nil {
		t.Error(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := env.SleepCtx(ctx, time.Hour); err == nil {
		t.Error("cancelled sleep returned nil")
	}
	slept := time.Duration(0)
	env.Sleep = func(ctx context.Context, d time.Duration) error { slept = d; return nil }
	_ = env.SleepCtx(context.Background(), 3*time.Second)
	if slept != 3*time.Second {
		t.Error("injected sleeper unused")
	}
	env.Now = func() time.Time { return time.Unix(5, 0) }
	if env.Clock().Unix() != 5 {
		t.Error("injected clock unused")
	}
	if e := NewEnv(strings.NewReader(""), &bytes.Buffer{}, &bytes.Buffer{}); e.Vars == nil || e.LookPath == nil {
		t.Error("NewEnv")
	}
}

func TestObject(t *testing.T) {
	o, err := ParseObject([]byte(`{"b":1,"a":{"x":null},"c":[1,2]}`))
	if err != nil {
		t.Fatal(err)
	}
	o.Set("d", json.RawMessage(`"z"`))
	o.Set("b", json.RawMessage(`2`))
	o.Delete("c")
	o.Delete("missing")
	out, _ := json.Marshal(o)
	if string(out) != `{"b":2,"a":{"x":null},"d":"z"}` {
		t.Errorf("marshal = %s", out)
	}
	if o.Len() != 3 || strings.Join(o.Keys(), ",") != "b,a,d" {
		t.Error("keys")
	}
	if _, ok := o.Get("a"); !ok {
		t.Error("Get")
	}
	for _, bad := range []string{`[1]`, `{"a":`, ``, `{"a":1`} {
		if _, err := ParseObject([]byte(bad)); err == nil {
			t.Errorf("%q parsed", bad)
		}
	}
}

// FR-015: no operation is a usage error (252) like the AWS CLI; explicit help still exits 0.
func TestMissingCommandOrOperationExits252(t *testing.T) {
	reg := testRegistry(t)
	cases := []struct {
		args []string
		want string
	}{
		{nil, "the following arguments are required: command"},
		{[]string{"--profile", "p"}, "the following arguments are required: command"},
		{[]string{"teams"}, "the following arguments are required: operation"},
	}
	for _, c := range cases {
		env, out, eb := testEnv()
		if code := Run(context.Background(), c.args, env, reg); code != ExitUsage || !strings.Contains(eb.String(), "afs: error: "+c.want) || out.Len() != 0 {
			t.Errorf("%v: code %d out %q err %q", c.args, code, out.String(), eb.String())
		}
	}
	for _, args := range [][]string{{"help"}, {"teams", "help"}, {"--help"}, {"help", "teams"}, {"-h"}} {
		env, _, eb := testEnv()
		if code := Run(context.Background(), args, env, reg); code != ExitOK {
			t.Errorf("%v: code %d err %q", args, code, eb.String())
		}
	}
}
