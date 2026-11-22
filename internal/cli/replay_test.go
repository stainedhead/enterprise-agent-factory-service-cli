package cli_test

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
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/pb33f/libopenapi"
	validator "github.com/pb33f/libopenapi-validator"
	verrors "github.com/pb33f/libopenapi-validator/errors"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"gopkg.in/yaml.v3"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/api"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli/services"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/spec"
)

// raw is the snapshot as generic YAML, used to read request and response examples.
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

// resolve follows a local $ref.
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

// example reads the example of a media type object (example, or the first of examples).
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

func loadRaw(t *testing.T) raw {
	t.Helper()
	var m map[string]any
	if err := yaml.Unmarshal(api.OpenAPI, &m); err != nil {
		t.Fatal(err)
	}
	return raw(m)
}

// jsonable turns yaml-decoded values into JSON-marshalable ones.
func jsonable(v any) any {
	switch x := v.(type) {
	case map[string]any:
		out := map[string]any{}
		for k, e := range x {
			out[k] = jsonable(e)
		}
		return out
	case map[any]any:
		out := map[string]any{}
		for k, e := range x {
			out[fmt.Sprint(k)] = jsonable(e)
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

func method(m string) string { return strings.ToLower(m) }

type replay struct {
	t       *testing.T
	doc     raw
	val     validator.Validator
	model   *libopenapi.DocumentModel[v3.Document]
	reg     *cli.Registry
	dir     string
	mu      sync.Mutex
	invalid []string
}

func newReplay(t *testing.T) *replay {
	t.Helper()
	d, err := libopenapi.NewDocument(api.OpenAPI)
	if err != nil {
		t.Fatal(err)
	}
	model, merrs := d.BuildV3Model()
	if merrs != nil {
		t.Fatal(merrs)
	}
	v, errs := validator.NewValidator(d)
	if len(errs) > 0 {
		t.Fatal(errs)
	}
	reg, err := services.NewRegistry()
	if err != nil {
		t.Fatal(err)
	}
	return &replay{t: t, doc: loadRaw(t), val: v, model: model, reg: reg, dir: t.TempDir()}
}

// input builds the --cli-input-json document of an operation from the snapshot's examples.
func (rp *replay) input(op *spec.Operation) map[string]any {
	in := map[string]any{}
	for _, p := range op.Params {
		if p.Example != "" {
			in[p.Name] = p.Example
		} else if len(p.Schema.Enum) > 0 {
			in[p.Name] = p.Schema.Enum[0]
		} else if p.In == "path" {
			in[p.Name] = "x1"
		}
	}
	opNode := rp.doc.at("paths", op.Path, method(op.Method))
	om, _ := opNode.(map[string]any)
	if om != nil {
		if rb, ok := rp.doc.resolve(om["requestBody"]).(map[string]any); ok {
			if content, ok := rb["content"].(map[string]any); ok {
				if ex, ok := rp.doc.example(content["application/json"]); ok {
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
	return in
}

// success returns the documented 2xx status and example body.
func (rp *replay) success(op *spec.Operation) (int, any, bool) {
	om, _ := rp.doc.at("paths", op.Path, method(op.Method)).(map[string]any)
	resp, _ := om["responses"].(map[string]any)
	for _, code := range op.Success {
		r, _ := rp.doc.resolve(resp[strconv.Itoa(code)]).(map[string]any)
		content, _ := r["content"].(map[string]any)
		if ex, ok := rp.doc.example(content["application/json"]); ok {
			return code, jsonable(ex), true
		}
	}
	code := 200
	if len(op.Success) > 0 {
		code = op.Success[0]
	}
	return code, nil, false
}

func normalize(v any) any {
	switch x := v.(type) {
	case map[string]any:
		out := map[string]any{}
		for k, e := range x {
			if k == "nextToken" || k == "NextToken" {
				continue
			}
			out[k] = normalize(e)
		}
		return out
	case []any:
		out := make([]any, len(x))
		for i, e := range x {
			out[i] = normalize(e)
		}
		return out
	}
	return v
}

func (rp *replay) run(op *spec.Operation, status int, body []byte, extra ...string) (int, string, string, *http.Request) {
	t := rp.t
	in := rp.input(op)
	b, _ := json.Marshal(in)
	f := filepath.Join(rp.dir, op.ID+".json")
	if err := os.WriteFile(f, b, 0o600); err != nil {
		t.Fatal(err)
	}
	var seen *http.Request
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bodyBytes, _ := io.ReadAll(r.Body)
		r.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		clone := r.Clone(context.Background())
		clone.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		clone.URL.Scheme, clone.URL.Host = "https", "afs.example.com"
		clone.RequestURI = ""
		seen = clone
		if status < 300 {
			// the validator cannot match "{id}:action" templates by path, so hand it the path item
			item, found := rp.model.Model.Paths.PathItems.Get(op.Path)
			var ok bool
			var errs []*verrors.ValidationError
			if found {
				ok, errs = rp.val.ValidateHttpRequestWithPathItem(clone, item, op.Path)
			} else {
				ok, errs = rp.val.ValidateHttpRequest(clone)
			}
			if !ok {
				for _, e := range errs {
					msg := e.Message + " " + e.Reason
					if strings.Contains(msg, "tag.") {
						continue // CLI-CD-4: tag.<key> query filters are open-ended
					}
					rp.mu.Lock()
					rp.invalid = append(rp.invalid, op.ID+": "+msg)
					rp.mu.Unlock()
				}
			}
		}
		if status != 204 {
			w.Header().Set("Content-Type", "application/json")
		}
		w.Header().Set("X-Request-Id", "req_replay")
		w.WriteHeader(status)
		_, _ = w.Write(body)
	}))
	defer srv.Close()
	var out, eb bytes.Buffer
	env := cli.NewEnv(strings.NewReader(""), &out, &eb)
	env.Vars = map[string]string{"AFS_ENDPOINT_URL": srv.URL, "AFS_API_KEY": "afsk_replay"}
	env.Home = rp.dir
	env.Sleep = func(context.Context, time.Duration) error { return nil }
	args := []string{op.Service, op.Command, "--cli-input-json", "file://" + f}
	if op.Paginated {
		args = append(args, "--no-paginate")
	}
	args = append(args, extra...)
	code := cli.Run(context.Background(), args, env, rp.reg)
	return code, out.String(), eb.String(), seen
}

func TestContractReplayAllOperations(t *testing.T) {
	rp := newReplay(t)
	n := 0
	for _, op := range spec.Commands() {
		op := op
		if cmd := rp.reg.ByOp(op.ID); cmd == nil {
			t.Errorf("%s has no command", op.ID)
			continue
		}
		n++
		t.Run(op.ID, func(t *testing.T) {
			rp.t = t
			status, ex, hasEx := rp.success(op)
			var body []byte
			if hasEx {
				body, _ = json.Marshal(ex)
			}
			code, out, errOut, _ := rp.run(op, status, body)
			if code != 0 {
				t.Fatalf("exit %d: %s", code, errOut)
			}
			if hasEx {
				var got any
				if err := json.Unmarshal([]byte(out), &got); err != nil {
					t.Fatalf("output is not JSON: %v\n%s", err, out)
				}
				var want any
				_ = json.Unmarshal(body, &want)
				if !reflect.DeepEqual(normalize(got), normalize(want)) {
					t.Errorf("output differs from the documented example\n got: %s\nwant: %s", out, body)
				}
			}
			// (c) AccessDenied
			denied := []byte(`{"code":"AccessDenied","message":"denied","requestId":"req_d","details":{"requiredPermission":"` + op.Permission + `"}}`)
			code, _, errOut, _ = rp.run(op, 403, denied)
			if code != 254 || !strings.Contains(errOut, "(AccessDenied) when calling the "+op.ID+" operation") || !strings.Contains(errOut, op.Permission) {
				t.Errorf("AccessDenied: exit %d stderr %q", code, errOut)
			}
		})
	}
	if n != 187 {
		t.Errorf("replayed %d operations, want 187", n)
	}
	rp.mu.Lock()
	defer rp.mu.Unlock()
	if len(rp.invalid) > 0 {
		sort.Strings(rp.invalid)
		t.Errorf("%d requests do not validate against the snapshot:\n%s", len(rp.invalid), strings.Join(rp.invalid, "\n"))
	}
}

// The validator must reject a request that breaks the snapshot, or the replay proves nothing.
func TestReplayValidatorRejectsBadRequest(t *testing.T) {
	rp := newReplay(t)
	item, _ := rp.model.Model.Paths.PathItems.Get("/v1/teams")
	req, _ := http.NewRequestWithContext(context.Background(), "POST", "https://afs.example.com/v1/teams", strings.NewReader(`{"slug":5}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer x")
	if ok, _ := rp.val.ValidateHttpRequestWithPathItem(req, item, "/v1/teams"); ok {
		t.Fatal("a body that breaks the schema validated")
	}
}
