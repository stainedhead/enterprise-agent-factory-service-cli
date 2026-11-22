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
		for _, it := range items {
			var m map[string]json.RawMessage
			if json.Unmarshal(it, &m) != nil {
				continue
			}
			for _, k := range []string{"id", "name", "key", r.Match} {
				var s string
				if json.Unmarshal(m[k], &s) == nil && s == want {
					return versionOf(it)
				}
			}
		}
		token = ""
		_ = json.Unmarshal(doc["nextToken"], &token)
		if token == "" {
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

// Notice is the line printed to stderr when the CLI read the version itself.
func Notice(version string) string {
	return fmt.Sprintf("Read the current version %s; sending If-Match: \"%s\". Pass --if-match to choose a version.", version, version)
}
