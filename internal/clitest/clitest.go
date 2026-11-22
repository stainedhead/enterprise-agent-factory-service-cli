// Package clitest is the shared test harness of the service packages: it runs afs against an
// httptest server, replays the snapshot's documented examples, and checks the "command done"
// criteria (example, columns, goldens, AccessDenied) of every operation (D-C4). Tests only.
package clitest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/api"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli/services"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/golden"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/spec"
)

// Result is one afs run.
type Result struct {
	Code        int
	Out, Err    string
	Req         *http.Request
	ReqBody     []byte
	Requests    []*http.Request
	RequestBody [][]byte
}

// Options tune one run.
type Options struct {
	Stdin string
	Vars  map[string]string
	TTY   bool
}

// Reply writes a JSON response with a request id.
func Reply(w http.ResponseWriter, status int, body string) {
	if status != 204 {
		w.Header().Set("Content-Type", "application/json")
	}
	w.Header().Set("X-Request-Id", "req_test")
	w.WriteHeader(status)
	_, _ = w.Write([]byte(body))
}

// Run starts a server with h, runs afs with args and returns the result with every request seen.
func Run(t *testing.T, h http.HandlerFunc, args ...string) Result {
	return RunWith(t, Options{}, h, args...)
}

// RunWith is Run with stdin, environment and terminal settings.
func RunWith(t *testing.T, o Options, h http.HandlerFunc, args ...string) Result {
	t.Helper()
	var res Result
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		r.Body = io.NopCloser(bytes.NewReader(b))
		res.Requests = append(res.Requests, r.Clone(context.Background()))
		res.RequestBody = append(res.RequestBody, b)
		h(w, r)
	}))
	defer srv.Close()
	reg, err := services.NewRegistry()
	if err != nil {
		t.Fatal(err)
	}
	var out, eb bytes.Buffer
	env := cli.NewEnv(strings.NewReader(o.Stdin), &out, &eb)
	env.Vars = map[string]string{"AFS_ENDPOINT_URL": srv.URL, "AFS_API_KEY": "afsk_unit_test_key"}
	for k, v := range o.Vars {
		env.Vars[k] = v
	}
	env.Home = t.TempDir()
	env.StdinTTY = o.TTY
	env.Sleep = func(context.Context, time.Duration) error { return nil }
	res.Code = cli.Run(context.Background(), args, env, reg)
	res.Out, res.Err = out.String(), eb.String()
	if n := len(res.Requests); n > 0 {
		res.Req, res.ReqBody = res.Requests[n-1], res.RequestBody[n-1]
	}
	return res
}

// Denied answers every request with the API's AccessDenied.
func Denied(permission string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		Reply(w, 403, `{"code":"AccessDenied","message":"not allowed","requestId":"req_d","details":{"requiredPermission":"`+permission+`"}}`)
	}
}

// Fail answers every request with an API error.
func Fail(status int, code, message string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		Reply(w, status, fmt.Sprintf(`{"code":%q,"message":%q,"requestId":"req_e"}`, code, message))
	}
}

// ---- the snapshot's examples ---------------------------------------------------------------

type raw map[string]any

func (r raw) at(path ...string) any {
	var cur any = map[string]any(r)
	for _, p := range path {
		m, ok := cur.(map[string]any)
		if !ok {
			return nil
		}
		cur = m[p]
	}
	return cur
}

func (r raw) resolve(v any) any {
	for i := 0; i < 8; i++ {
		m, ok := v.(map[string]any)
		if !ok {
			return v
		}
		ref, ok := m["$ref"].(string)
		if !ok {
			return v
		}
		v = r.at(strings.Split(strings.TrimPrefix(ref, "#/"), "/")...)
	}
	return v
}

func (r raw) example(media any) (any, bool) {
	m, ok := r.resolve(media).(map[string]any)
	if !ok {
		return nil, false
	}
	if ex, ok := m["example"]; ok {
		return ex, true
	}
	if exs, ok := m["examples"].(map[string]any); ok {
		names := make([]string, 0, len(exs))
		for k := range exs {
			names = append(names, k)
		}
		sort.Strings(names)
		for _, n := range names {
			if e, ok := r.resolve(exs[n]).(map[string]any); ok {
				if v, ok := e["value"]; ok {
					return v, true
				}
			}
		}
	}
	return nil, false
}

func jsonable(v any) any {
	switch x := v.(type) {
	case map[string]any:
		out := map[string]any{}
		for k, e := range x {
			out[k] = jsonable(e)
		}
		return out
	case []any:
		out := make([]any, len(x))
		for i, e := range x {
			out[i] = jsonable(e)
		}
		return out
	case time.Time:
		return x.UTC().Format(time.RFC3339)
	}
	return v
}

var doc raw

func loadDoc(t *testing.T) raw {
	t.Helper()
	if doc == nil {
		var m map[string]any
		if err := yaml.Unmarshal(api.OpenAPI, &m); err != nil {
			t.Fatal(err)
		}
		doc = raw(m)
	}
	return doc
}

// ExampleResponse returns the documented success status and example body of an operation.
func ExampleResponse(t *testing.T, op *spec.Operation) (int, []byte) {
	t.Helper()
	d := loadDoc(t)
	om, _ := d.at("paths", op.Path, strings.ToLower(op.Method)).(map[string]any)
	resp, _ := om["responses"].(map[string]any)
	for _, code := range op.Success {
		r, _ := d.resolve(resp[strconv.Itoa(code)]).(map[string]any)
		content, _ := r["content"].(map[string]any)
		if ex, ok := d.example(content["application/json"]); ok {
			b, _ := json.Marshal(jsonable(ex))
			return code, b
		}
	}
	code := 200
	if len(op.Success) > 0 {
		code = op.Success[0]
	}
	return code, nil
}

// InputJSON builds a --cli-input-json document of an operation from the snapshot's examples
// and writes it to a file under dir; it returns the file:// argument.
func InputJSON(t *testing.T, op *spec.Operation, dir string) string {
	t.Helper()
	d := loadDoc(t)
	in := map[string]any{}
	for _, p := range op.Params {
		switch {
		case p.Example != "":
			in[p.Name] = p.Example
		case len(p.Schema.Enum) > 0:
			in[p.Name] = p.Schema.Enum[0]
		case p.In == "path":
			in[p.Name] = "x1"
		}
	}
	if om, _ := d.at("paths", op.Path, strings.ToLower(op.Method)).(map[string]any); om != nil {
		if rb, ok := d.resolve(om["requestBody"]).(map[string]any); ok {
			if content, ok := rb["content"].(map[string]any); ok {
				if ex, ok := d.example(content["application/json"]); ok {
					if m, ok := jsonable(ex).(map[string]any); ok {
						for k, v := range m {
							in[k] = v
						}
					}
				}
			}
		}
	}
	if op.IfMatch != spec.IfMatchNone {
		in["version"] = 1
	}
	b, _ := json.Marshal(in)
	f := filepath.Join(dir, op.ID+".json")
	if err := os.WriteFile(f, b, 0o600); err != nil {
		t.Fatal(err)
	}
	return "file://" + f
}

// Split breaks a shell-like line into words, honouring single and double quotes.
func Split(line string) []string {
	var out []string
	var cur strings.Builder
	var q rune
	in := false
	for _, r := range line {
		switch {
		case q != 0:
			if r == q {
				q = 0
			} else {
				cur.WriteRune(r)
			}
		case r == '\'' || r == '"':
			q, in = r, true
		case r == ' ':
			if in || cur.Len() > 0 {
				out = append(out, cur.String())
				cur.Reset()
				in = false
			}
		default:
			cur.WriteRune(r)
			in = true
		}
	}
	if in || cur.Len() > 0 {
		out = append(out, cur.String())
	}
	return out
}

// CheckOperations proves "command done" (tasks.md M1) for each operation: at least one example
// that parses, table columns, golden table/text/help files under testdata/<group>/, and the
// AccessDenied answer (exit 254, AWS-style message).
func CheckOperations(t *testing.T, group string, ids ...string) {
	t.Helper()
	reg, err := services.NewRegistry()
	if err != nil {
		t.Fatal(err)
	}
	for _, id := range ids {
		cmd := reg.ByOp(id)
		if cmd == nil {
			t.Errorf("%s: no command", id)
			continue
		}
		op := cmd.Op
		t.Run(id, func(t *testing.T) {
			if len(cmd.Examples) == 0 { // CLI-14
				t.Fatalf("%s has no example", id)
			}
			for _, ex := range cmd.Examples {
				words := Split(ex.Command)
				if len(words) < 3 || words[0] != "afs" || words[1] != cmd.Service || words[2] != cmd.Name {
					t.Errorf("example %q does not call %s", ex.Command, cmd.Key())
					continue
				}
				if _, _, err := reg.Parse(words[1:]); err != nil {
					t.Errorf("example %q does not parse: %v", ex.Command, err)
				}
			}
			status, body := ExampleResponse(t, op)
			hasBody := len(body) > 0 && status != 204
			if hasBody && len(cmd.Columns) == 0 {
				t.Fatalf("%s has no table columns", id)
			}
			in := InputJSON(t, op, t.TempDir())
			args := []string{op.Service, op.Command, "--cli-input-json", in}
			if op.Paginated {
				args = append(args, "--no-paginate")
			}
			h := func(w http.ResponseWriter, _ *http.Request) { Reply(w, status, string(body)) }
			if hasBody {
				for _, format := range []string{"table", "text"} {
					r := Run(t, h, append(append([]string{}, args...), "--output", format)...)
					if r.Code != 0 {
						t.Fatalf("%s: exit %d: %s", format, r.Code, r.Err)
					}
					golden.Check(t, group+"/"+id+"."+format+".golden", r.Out)
				}
			}
			r := Run(t, nil, op.Service, op.Command, "help")
			if r.Code != 0 {
				t.Fatalf("help: exit %d", r.Code)
			}
			golden.Check(t, group+"/"+id+".help.golden", r.Out)
			d := Run(t, Denied(op.Permission), args...) // AccessDenied
			if d.Code != 254 || !strings.Contains(d.Err, "(AccessDenied) when calling the "+id+" operation") {
				t.Errorf("AccessDenied: exit %d stderr %q", d.Code, d.Err)
			}
		})
	}
}
