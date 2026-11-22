package exec_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli/exec"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli/safeupdate"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli/services"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/client"
)

// runEnv is run with control over the environment: variables, terminal, password prompt, sleeping, context.
func runEnv(t *testing.T, h http.HandlerFunc, mut func(*cli.Env), args ...string) result {
	t.Helper()
	srv := newServer(t, h)
	reg, err := services.NewRegistry()
	if err != nil {
		t.Fatal(err)
	}
	var out, eb bytes.Buffer
	env := cli.NewEnv(strings.NewReader(""), &out, &eb)
	env.Vars = map[string]string{"AFS_ENDPOINT_URL": srv, "AFS_API_KEY": "afsk_unit_test_key"}
	env.Home = t.TempDir()
	env.Sleep = func(context.Context, time.Duration) error { return nil }
	if mut != nil {
		mut(env)
	}
	code := cli.Run(context.Background(), args, env, reg)
	return result{code, out.String(), eb.String()}
}

func echo(w http.ResponseWriter, _ *http.Request) { jsonReply(w, 200, `{"id":"x","version":3}`) }

// FR-019: numeric limits, input JSON members and header members.
func TestOptionRangesAndInputJSONMembers(t *testing.T) {
	r := run(t, echo, "pools", "scale-pool", "--team-id", "t", "--pool-id", "p", "--workers", "201")
	if r.code != 252 || !strings.Contains(r.errOut, "above the maximum 200") {
		t.Errorf("max: %+v", r)
	}
	r = run(t, echo, "pools", "create-pool", "--team-id", "t", "--scratch-size-gib", "0")
	if r.code != 252 || !strings.Contains(r.errOut, "below the minimum 1") {
		t.Errorf("min: %+v", r)
	}
	var hdr http.Header
	var path string
	h := func(w http.ResponseWriter, req *http.Request) {
		hdr, path = req.Header.Clone(), req.URL.Path
		jsonReply(w, 200, `{"id":"x"}`)
	}
	r = run(t, h, "pools", "scale-pool", "--cli-input-json", `{"teamId":"team_1","poolId":"pool_1","workers":3,"version":4,"Idempotency-Key":"k-1"}`)
	if r.code != 0 || hdr.Get("Idempotency-Key") != "k-1" || path != "/v1/teams/team_1/pools/pool_1:scale" {
		t.Errorf("input json: %+v %v %s", r, hdr, path)
	}
	// an input-JSON parameter is checked like an option
	r = run(t, echo, "budgets", "list-costs", "--cli-input-json", `{"teamId":"team_1","from":"x"}`, "--output", "csv", "--no-paginate")
	if r.code == 252 {
		t.Errorf("valid query parameter from input json refused: %+v", r)
	}
	r = run(t, echo, "teams", "update-team", "--team-id", "team_1", "--wait", "--wait-timeout", "soon")
	if r.code != 252 || !strings.Contains(r.errOut, "--wait-timeout") {
		t.Errorf("wait-timeout: %+v", r)
	}
}

func TestTagPairsAndSecretPrompt(t *testing.T) {
	var body string
	h := func(w http.ResponseWriter, req *http.Request) {
		if req.Method == "GET" {
			jsonReply(w, 200, `{"id":"team_1","version":3}`)
			return
		}
		b := new(bytes.Buffer)
		_, _ = b.ReadFrom(req.Body)
		body = b.String()
		jsonReply(w, 200, `{"id":"team_1","version":4}`)
	}
	r := run(t, h, "teams", "put-team-tags", "--team-id", "team_1", "--tags", "project=x", "owner=y")
	if r.code != 0 || !strings.Contains(body, `"key":"project"`) || !strings.Contains(body, `"value":"y"`) {
		t.Errorf("pairs: %+v %s", r, body)
	}
	r = run(t, h, "teams", "put-team-tags", "--team-id", "team_1", "--tags", `[{"key":"a","value":"b"}]`)
	if r.code != 0 || !strings.Contains(body, `"key":"a"`) {
		t.Errorf("json: %+v %s", r, body)
	}
	if r := run(t, h, "teams", "put-team-tags", "--team-id", "team_1", "--tags", "novalue"); r.code != 252 {
		t.Errorf("bad pair: %+v", r)
	}

	// a missing required secret is prompted for on a terminal, refused without one
	args := []string{"secrets", "create-secret", "--team-id", "team_1", "--name", "n", "--scope", "team"}
	r = runEnv(t, h, func(e *cli.Env) {
		e.StdinTTY = true
		e.ReadPassword = func(string) ([]byte, error) { return []byte("hunter2"), nil }
	}, args...)
	if r.code != 0 || !strings.Contains(body, `"value":"hunter2"`) || strings.Contains(r.out+r.errOut, "hunter2") {
		t.Errorf("prompt: %+v %s", r, body)
	}
	r = runEnv(t, h, nil, args...)
	if r.code != 252 || strings.Contains(r.errOut, "hunter2") {
		t.Errorf("no terminal: %+v", r)
	}
	r = runEnv(t, h, func(e *cli.Env) {
		e.StdinTTY = true
		e.ReadPassword = func(string) ([]byte, error) { return nil, errors.New("no tty") }
	}, args...)
	if r.code != 252 {
		t.Errorf("prompt failure: %+v", r)
	}
}

func TestPagingOptionErrorsAndEdges(t *testing.T) {
	var mu sync.Mutex
	var q []string
	for _, args := range [][]string{
		{"teams", "list-teams", "--max-items", "0"},
		{"teams", "list-teams", "--page-size", "0"},
		{"teams", "describe-team", "--team-id", "team_1", "--max-items", "1"},
		{"teams", "list-teams", "--output", "csv"},
		{"teams", "list-teams", "--starting-token", "!!!notatoken"},
	} {
		if r := run(t, teamsPages(&mu, &q), args...); r.code != 252 {
			t.Errorf("%v: %+v", args, r)
		}
	}
	// a page that is not a JSON object
	if r := run(t, func(w http.ResponseWriter, _ *http.Request) { jsonReply(w, 200, `[1]`) }, "teams", "list-teams"); r.code != 255 {
		t.Errorf("array page: %+v", r)
	}
	// a cut mid-page prints a token that resumes after the last printed item
	r := run(t, teamsPages(&mu, &q), "teams", "list-teams", "--max-items", "3")
	if r.code != 0 || count(r.out) != 3 || !strings.Contains(r.out, "NextToken") {
		t.Fatalf("%+v", r)
	}
	// csv with an explicit page size
	var query string
	r = run(t, func(w http.ResponseWriter, req *http.Request) {
		query = req.URL.RawQuery
		w.Header().Set("Content-Type", "text/csv")
		_, _ = w.Write([]byte("a\n1\n"))
	}, "budgets", "list-costs", "--team-id", "team_1", "--output", "csv", "--page-size", "5")
	if r.code != 0 || !strings.Contains(query, "maxResults=5") {
		t.Errorf("csv page size: %+v %s", r, query)
	}
	// a failing page ends the CSV listing with the API error
	r = run(t, func(w http.ResponseWriter, _ *http.Request) {
		jsonReply(w, 403, `{"code":"AccessDenied","message":"no","requestId":"r"}`)
	}, "budgets", "list-costs", "--team-id", "team_1", "--output", "csv")
	if r.code != 254 {
		t.Errorf("csv error: %+v", r)
	}
}

func TestRunErrorPaths(t *testing.T) {
	// an invalid configured output format is a configuration error (255)
	r := runEnv(t, echo, func(e *cli.Env) { e.Vars["AFS_DEFAULT_OUTPUT"] = "bogus" }, "teams", "describe-team", "--team-id", "team_1")
	if r.code != 255 || !strings.Contains(r.errOut, "output") {
		t.Errorf("output: %+v", r)
	}
	// the version read of a safe update fails
	r = run(t, func(w http.ResponseWriter, req *http.Request) {
		jsonReply(w, 500, `{"code":"InternalError","message":"down","requestId":"r"}`)
	}, "teams", "put-team-tags", "--team-id", "team_1", "--tags", "a=b")
	if r.code != 254 {
		t.Errorf("version read: %+v", r)
	}
	// a handled command's own error
	r = run(t, func(w http.ResponseWriter, _ *http.Request) {
		jsonReply(w, 500, `{"code":"InternalError","message":"down","requestId":"r"}`)
	}, "teams", "get-kubeconfig", "--team-id", "team_1", "--kubeconfig", "-")
	if r.code != 254 {
		t.Errorf("handled: %+v", r)
	}
	// Ctrl-C while --wait polls: the change stands
	r = runEnv(t, func(w http.ResponseWriter, req *http.Request) {
		if req.Method == "POST" {
			jsonReply(w, 202, `{"id":"team_1","generation":3}`)
			return
		}
		jsonReply(w, 200, `{"status":{"phase":"Pending","observedGeneration":3}}`)
	}, func(e *cli.Env) {
		e.Sleep = func(context.Context, time.Duration) error { return context.Canceled }
	}, "teams", "resume-team", "--team-id", "team_1", "--wait")
	if r.code != 130 || !strings.Contains(r.errOut, "the request was sent and the change stands") {
		t.Errorf("interrupted wait: %+v", r)
	}
	// --generate-cli-skeleton makes no request
	r = runEnv(t, func(http.ResponseWriter, *http.Request) { t.Error("request sent") }, nil, "teams", "create-team", "--generate-cli-skeleton")
	if r.code != 0 || !strings.Contains(r.out, "slug") {
		t.Errorf("skeleton: %+v", r)
	}
	// a read from a file that cannot be read is a usage error
	dir := t.TempDir()
	p := filepath.Join(dir, "secret.txt")
	if err := os.WriteFile(p, []byte("v"), 0o600); err != nil {
		t.Fatal(err)
	}
	r = run(t, echo, "secrets", "create-secret", "--team-id", "team_1", "--name", "n", "--scope", "team", "--secret-value", "file://"+filepath.Join(dir, "missing"))
	if r.code != 252 {
		t.Errorf("missing secret file: %+v", r)
	}
}

// ReportError: every detail the service can attach is printed, and nothing else.
func TestReportErrorDetailLines(t *testing.T) {
	api := func(status int, details string) *client.APIError {
		e := &client.APIError{Status: status, Code: "ValidationError", Message: "bad", RequestID: "req_1", Header: http.Header{}}
		e.Details = []byte(details)
		return e
	}
	var b bytes.Buffer
	code := exec.ReportError(&b, "CreateWorkOrder", "p", api(400, `{"ambiguousRepos":[{"repo":"github.com/a/b","teamIds":["team_1","team_2"]}],"unownedRepos":["github.com/a/c"],"dependents":["agt_1"],"missingTags":["cost"]}`))
	out := b.String()
	for _, want := range []string{"Ambiguous repository: github.com/a/b is owned by more than one team (team_1, team_2)", "Unowned repositories: github.com/a/c", "Name the owner of each in --repos", "Dependents: agt_1", "Missing tags: cost", "Request ID: req_1"} {
		if code != 254 || !strings.Contains(out, want) {
			t.Errorf("missing %q in\n%s", want, out)
		}
	}
	b.Reset()
	exec.ReportError(&b, "UpdateWorkOrderRepo", "p", &safeupdate.ReadDeniedError{Err: api(403, `{"requiredPermission":"work-orders:read"}`)})
	if !strings.Contains(b.String(), safeupdate.ReadDeniedHint) || !strings.Contains(b.String(), "Required permission: work-orders:read") {
		t.Errorf("read denied:\n%s", b.String())
	}
	b.Reset()
	if code := exec.ReportError(&b, "X", "p", fmt.Errorf("wrapped: %w", context.Canceled)); code != 130 || b.Len() != 0 {
		t.Errorf("canceled: %d %q", code, b.String())
	}
	b.Reset()
	if code := exec.ReportError(&b, "X", "prod", errors.New("plain")); code != 255 || !strings.Contains(b.String(), "afs: error: plain") {
		t.Errorf("plain: %d %q", code, b.String())
	}
}

// A profile with SSO settings and no cached session asks for a sign-in (255).
func TestLoginRequiredIsReported(t *testing.T) {
	var out, eb bytes.Buffer
	reg, _ := services.NewRegistry()
	env := cli.NewEnv(strings.NewReader(""), &out, &eb)
	srv := newServer(t, echo)
	env.Vars = map[string]string{"AFS_ENDPOINT_URL": srv}
	env.Home = t.TempDir()
	if err := os.MkdirAll(filepath.Join(env.Home, ".afs"), 0o700); err != nil {
		t.Fatal(err)
	}
	cfg := "[default]\nsso_issuer = https://idp.example.com\nsso_client_id = c\n"
	if err := os.WriteFile(filepath.Join(env.Home, ".afs", "config"), []byte(cfg), 0o600); err != nil {
		t.Fatal(err)
	}
	code := cli.Run(context.Background(), []string{"teams", "describe-team", "--team-id", "team_1"}, env, reg)
	if code != 255 || !strings.Contains(eb.String(), "afs sso login") {
		t.Errorf("%d %q", code, eb.String())
	}
}
