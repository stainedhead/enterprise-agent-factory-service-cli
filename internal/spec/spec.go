// Package spec holds the operation descriptors generated from the pinned OpenAPI snapshot.
// The data lives in operations_gen.go (tools/afsgen); this file has the types and lookups.
package spec

import "sort"

// IfMatchMode says whether an operation takes If-Match.
type IfMatchMode int

// If-Match modes: none, required (PATCH, DELETE) or optional (PUT, CD-2).
const (
	IfMatchNone IfMatchMode = iota
	IfMatchRequired
	IfMatchOptional
)

// Operation describes one OpenAPI operation. One value per operation; generated, never edited.
type Operation struct {
	ID         string // operationId, e.g. "UpdateAgent"
	Method     string // GET PUT POST PATCH DELETE
	Path       string // "/v1/teams/{teamId}/agents/{agentId}"
	Service    string // from x-afs-cli; "" for GetHealth and GetReadiness
	Command    string // from x-afs-cli, e.g. "update-agent"
	Summary    string
	Permission string   // x-afs-permission
	AnyOf      []string // x-afs-any-of
	Scope      string   // public | any | team | service
	Priority   string   // P0 | P1 | P2
	CSV        bool     // x-afs-csv
	Params     []Param  // path and query (maxResults/nextToken excluded; If-Match, Idempotency-Key, X-Request-Id excluded)
	IfMatch    IfMatchMode
	Idempotent bool     // takes Idempotency-Key (every POST)
	Body       *Body    // nil when no request body
	Paginated  bool     // takes maxResults and nextToken
	ResultKey  string   // "items" | "rows" | "" for paginated operations
	Success    []int    // 2xx codes
	Produces   []string // media types of the success responses
}

// Param is a path or query parameter.
type Param struct {
	Name     string // wire name: teamId, dryRun
	In       string // path | query
	Required bool
	Desc     string
	Example  string
	Schema   Schema
}

// Body is a JSON request body.
type Body struct {
	Required   bool
	MediaType  string
	SchemaName string
	Fields     []Field
}

// Field is a top-level property of a body.
type Field struct {
	Name      string
	Required  bool
	WriteOnly bool
	ReadOnly  bool
	Desc      string
	Schema    Schema
}

// Schema is the resolved type of a parameter or field.
type Schema struct {
	Type     string // string | integer | number | boolean | array | object | any
	Nullable bool
	Enum     []string
	Min, Max *float64
	MinLen   *int
	MaxLen   *int
	Pattern  string
	Items    *Schema
	Ref      string // component schema name, e.g. "Money", "SecretValue"
}

// All returns every operation, sorted by operationId.
func All() []Operation { return operations }

var byID map[string]*Operation

// ByID finds an operation by operationId, or nil.
func ByID(id string) *Operation {
	if byID == nil {
		m := make(map[string]*Operation, len(operations))
		for i := range operations {
			m[operations[i].ID] = &operations[i]
		}
		byID = m
	}
	return byID[id]
}

// Commands returns the operations that have a CLI command, ordered by service then command.
func Commands() []*Operation {
	var out []*Operation
	for i := range operations {
		if operations[i].Service != "" {
			out = append(out, &operations[i])
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Service != out[j].Service {
			return out[i].Service < out[j].Service
		}
		return out[i].Command < out[j].Command
	})
	return out
}

func f64(v float64) *float64 { return &v }
func i64(v int) *int         { return &v }
