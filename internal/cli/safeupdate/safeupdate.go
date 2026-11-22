// Package safeupdate finds the current version of a resource for operations that take If-Match (CLI-10):
// it reads the resource, or lists its collection and picks the item, and reports the version to send.
package safeupdate

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/client"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/spec"
)

// Doer sends one request.
type Doer interface {
	Do(ctx context.Context, req client.Request) (*client.Response, error)
}

// maxListPages bounds the search through a collection.
const maxListPages = 50

// Reader says how the current version of an operation's resource is read.
type Reader struct {
	Get   *spec.Operation // the GET operation
	Match string          // "" for a direct read; else the path parameter whose value matches an item
	Path  string          // path template to read
	// Composite names the path parameters (host, owner, repo) that together identify an item of the
	// collection as host == first and repo == "second/third" (a work order's repository, D-C17).
	Composite []string
}

// ReadDeniedError is a failure to read the current version because the caller may not read the
// resource. It unwraps to the API error.
type ReadDeniedError struct{ Err *client.APIError }

func (e *ReadDeniedError) Error() string { return e.Err.Error() }

// Unwrap returns the API error.
func (e *ReadDeniedError) Unwrap() error { return e.Err }

// ReadDeniedHint is printed after the API's answer when the version read was denied.
const ReadDeniedHint = "Pass --if-match <version> (the version printed by the previous update-repo)."

// composite lists operations whose resource is one item of a list identified by several path parameters.
var composite = map[string]Reader{
	"UpdateWorkOrderRepo": {Path: "/v1/work-orders/{workOrderId}/repos", Composite: []string{"host", "owner", "repo"}},
}

func segments(p string) []string { return strings.Split(strings.Trim(p, "/"), "/") }

func isParam(s string) bool { return strings.HasPrefix(s, "{") && strings.HasSuffix(s, "}") }

func getOps() map[string]*spec.Operation {
	m := map[string]*spec.Operation{}
	for i := range spec.All() {
		o := &spec.All()[i]
		if o.Method == "GET" {
			m[o.Path] = o
		}
	}
	return m
}

// versionless lists operations whose resource has no version: they send If-Match only when the
// user passes --if-match, and read nothing first (D-C8).
var versionless = map[string]bool{"PutTeamModels": true}

// ReaderFor returns the reader of an operation, or nil when its resource carries no readable
// version (the operation then sends no If-Match).
func ReaderFor(op *spec.Operation) *Reader {
	if versionless[op.ID] {
		return nil
	}
	gets := getOps()
	if c, ok := composite[op.ID]; ok {
		if g, ok := gets[c.Path]; ok {
			c.Get = g
			return &c
		}
		return nil
	}
	segs := segments(op.Path)
	// a trailing literal such as /value or /state belongs to the parent resource
	for len(segs) > 1 && !isParam(segs[len(segs)-1]) {
		if g, ok := gets["/"+strings.Join(segs, "/")]; ok {
			return &Reader{Get: g, Path: g.Path}
		}
		segs = segs[:len(segs)-1]
	}
	full := "/" + strings.Join(segs, "/")
	if g, ok := gets[full]; ok {
		return &Reader{Get: g, Path: full}
	}
	if len(segs) > 1 && isParam(segs[len(segs)-1]) {
		parent := "/" + strings.Join(segs[:len(segs)-1], "/")
		if g, ok := gets[parent]; ok && g.Paginated {
			return &Reader{Get: g, Match: strings.Trim(segs[len(segs)-1], "{}"), Path: parent}
		}
	}
	return nil
}

// Version reads the current version. ok is false when the resource does not exist (404) or
// carries no version; err is any other failure.
func Version(ctx context.Context, d Doer, op *spec.Operation, path map[string]string) (version string, ok bool, err error) {
	r := ReaderFor(op)
	if r == nil {
		return "", false, nil
	}
	pp := map[string]string{}
	for k, v := range path {
		pp[k] = v
	}
	var want string
	if r.Composite != nil {
		return versionInList(ctx, d, r, pp, func(m map[string]json.RawMessage) bool {
			var host, repo string
			_ = json.Unmarshal(m["host"], &host)
			_ = json.Unmarshal(m["repo"], &repo)
			return host == pp[r.Composite[0]] && repo == pp[r.Composite[1]]+"/"+pp[r.Composite[2]]
		})
	}
	if r.Match != "" {
		want = pp[r.Match]
		delete(pp, r.Match)
	}
	if r.Match == "" {
		resp, err := d.Do(ctx, client.Request{Method: "GET", Path: client.BuildPath(r.Path, pp), Operation: r.Get.ID})
		if err != nil {
			var ae *client.APIError
			if errors.As(err, &ae) && ae.Status == 404 {
				return "", false, nil
			}
			return "", false, err
		}
		return versionOf(resp.Body)
	}
	token := ""
	for page := 0; page < maxListPages; page++ {
		q := url.Values{}
		if token != "" {
			q.Set("nextToken", token)
		}
		resp, err := d.Do(ctx, client.Request{Method: "GET", Path: client.BuildPath(r.Path, pp), Query: q, Operation: r.Get.ID})
		if err != nil {
			var ae *client.APIError
			if errors.As(err, &ae) && ae.Status == 404 {
				return "", false, nil
			}
			return "", false, err
		}
		var doc map[string]json.RawMessage
		if json.Unmarshal(resp.Body, &doc) != nil {
			return "", false, nil
		}
		var items []json.RawMessage
		key := r.Get.ResultKey
		if key == "" {
			key = "items"
		}
		_ = json.Unmarshal(doc[key], &items)
		it, err := pickItem(items, want, r.Match)
		if err != nil {
			return "", false, err
		}
		if it != nil {
			return versionOf(it)
		}
		token = ""
		_ = json.Unmarshal(doc["nextToken"], &token)
		if token == "" {
			break
		}
	}
	return "", false, nil
}

// ErrAmbiguous is returned when more than one item of a list matches the resource.
var ErrAmbiguous = errors.New("more than one item matches; pass --if-match <version> to name the version")

// pickItem finds the item for want. Fields are tried in the order id, name, key, then the path parameter's
// name, over the whole page, so an id match beats a name match on another item; two items that match on the
// same field are ambiguous and no version is guessed.
func pickItem(items []json.RawMessage, want, match string) (json.RawMessage, error) {
	keys := []string{"id", "name", "key"}
	if match != "" && match != "id" && match != "name" && match != "key" {
		keys = append(keys, match)
	}
	for _, k := range keys {
		var found json.RawMessage
		n := 0
		for _, it := range items {
			var m map[string]json.RawMessage
			if json.Unmarshal(it, &m) != nil {
				continue
			}
			var s string
			if json.Unmarshal(m[k], &s) == nil && s == want {
				found = it
				n++
			}
		}
		switch {
		case n == 1:
			return found, nil
		case n > 1:
			return nil, ErrAmbiguous
		}
	}
	return nil, nil
}

// versionInList pages through the collection of r and returns the version of the item that match accepts.
func versionInList(ctx context.Context, d Doer, r *Reader, pp map[string]string, match func(map[string]json.RawMessage) bool) (string, bool, error) {
	listPath := map[string]string{}
	for k, v := range pp {
		if strings.Contains(r.Path, "{"+k+"}") {
			listPath[k] = v
		}
	}
	token := ""
	for page := 0; page < maxListPages; page++ {
		q := url.Values{}
		if token != "" {
			q.Set("nextToken", token)
		}
		resp, err := d.Do(ctx, client.Request{Method: "GET", Path: client.BuildPath(r.Path, listPath), Query: q, Operation: r.Get.ID})
		if err != nil {
			var ae *client.APIError
			if errors.As(err, &ae) {
				if ae.Status == 404 {
					return "", false, nil
				}
				if ae.Status == 403 {
					return "", false, &ReadDeniedError{Err: ae}
				}
			}
			return "", false, err
		}
		var doc struct {
			Items     []json.RawMessage `json:"items"`
			NextToken string            `json:"nextToken"`
		}
		if json.Unmarshal(resp.Body, &doc) != nil {
			return "", false, nil
		}
		for _, it := range doc.Items {
			var m map[string]json.RawMessage
			if json.Unmarshal(it, &m) == nil && match(m) {
				return versionOf(it)
			}
		}
		if token = doc.NextToken; token == "" {
			break
		}
	}
	return "", false, nil
}

func versionOf(body []byte) (string, bool, error) {
	var m map[string]json.RawMessage
	if err := json.Unmarshal(body, &m); err != nil {
		return "", false, nil
	}
	v, ok := m["version"]
	if !ok {
		return "", false, nil
	}
	var s string
	if json.Unmarshal(v, &s) == nil {
		return s, s != "", nil
	}
	return strings.TrimSpace(string(v)), true, nil
}

// Notice is the line printed to stderr after a successful write whose version the CLI read itself
// (PRD CLI-10, D-R4).
func Notice(version string) string {
	return fmt.Sprintf("Applied to version %s; use --if-match to guard against concurrent changes.", version)
}
