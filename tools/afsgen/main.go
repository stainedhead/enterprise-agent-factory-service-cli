// Command afsgen reads the pinned OpenAPI snapshot and writes internal/spec/operations_gen.go:
// one descriptor per operation, from the x-afs-* extensions, parameters, bodies and responses.
package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"go/format"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

func main() {
	in := flag.String("in", "api/openapi.yaml", "OpenAPI snapshot")
	out := flag.String("out", "internal/spec/operations_gen.go", "generated Go file")
	flag.Parse()
	raw, err := os.ReadFile(*in)
	if err != nil {
		fmt.Fprintln(os.Stderr, "afsgen:", err)
		os.Exit(1)
	}
	src, err := Generate(raw)
	if err != nil {
		fmt.Fprintln(os.Stderr, "afsgen:", err)
		os.Exit(1)
	}
	if err := os.WriteFile(*out, src, 0o644); err != nil { //nolint:gosec // generated source, world-readable
		fmt.Fprintln(os.Stderr, "afsgen:", err)
		os.Exit(1)
	}
}

// doc wraps a decoded OpenAPI document for ref resolution.
type doc struct{ root *yaml.Node }

func (d *doc) mapping(n *yaml.Node) *yaml.Node {
	if n != nil && n.Kind == yaml.DocumentNode && len(n.Content) > 0 {
		return n.Content[0]
	}
	return n
}

// get returns the value node for key in a mapping node, or nil.
func get(n *yaml.Node, key string) *yaml.Node {
	if n == nil || n.Kind != yaml.MappingNode {
		return nil
	}
	for i := 0; i+1 < len(n.Content); i += 2 {
		if n.Content[i].Value == key {
			return n.Content[i+1]
		}
	}
	return nil
}

func keys(n *yaml.Node) []string {
	var out []string
	if n == nil || n.Kind != yaml.MappingNode {
		return out
	}
	for i := 0; i+1 < len(n.Content); i += 2 {
		out = append(out, n.Content[i].Value)
	}
	return out
}

func str(n *yaml.Node) string {
	if n == nil {
		return ""
	}
	return n.Value
}

// resolve follows $ref chains (local refs only) and reports the last component name.
func (d *doc) resolve(n *yaml.Node) (*yaml.Node, string) {
	name := ""
	for n != nil {
		ref := get(n, "$ref")
		if ref == nil {
			return n, name
		}
		parts := strings.Split(strings.TrimPrefix(ref.Value, "#/"), "/")
		cur := d.mapping(d.root)
		for _, p := range parts {
			cur = get(cur, p)
		}
		if parts[len(parts)-2] == "schemas" {
			name = parts[len(parts)-1]
		}
		if cur == nil {
			panic("unresolved ref " + ref.Value)
		}
		// keep sibling keywords such as description of the referencing node
		n = cur
	}
	return n, name
}

type schema struct {
	Type     string
	Nullable bool
	Enum     []string
	Min, Max *string
	MinLen   *string
	MaxLen   *string
	Pattern  string
	Items    *schema
	Ref      string
	WriteOnly, ReadOnly bool
	Example  string
	Desc     string
}

func (d *doc) schemaOf(n *yaml.Node) schema {
	r, name := d.resolve(n)
	s := schema{Ref: name}
	if r == nil {
		s.Type = "any"
		return s
	}
	if v := get(n, "description"); v != nil {
		s.Desc = v.Value
	} else if v := get(r, "description"); v != nil {
		s.Desc = v.Value
	}
	if v := get(r, "writeOnly"); v != nil && v.Value == "true" {
		s.WriteOnly = true
	}
	if v := get(r, "readOnly"); v != nil && v.Value == "true" {
		s.ReadOnly = true
	}
	if t := get(r, "type"); t != nil {
		if t.Kind == yaml.SequenceNode {
			for _, c := range t.Content {
				if c.Value == "null" {
					s.Nullable = true
				} else {
					s.Type = c.Value
				}
			}
		} else {
			s.Type = t.Value
		}
	}
	if s.Type == "" {
		if get(r, "properties") != nil {
			s.Type = "object"
		} else {
			s.Type = "any"
		}
	}
	if e := get(r, "enum"); e != nil {
		for _, c := range e.Content {
			s.Enum = append(s.Enum, c.Value)
		}
	}
	if v := get(r, "minimum"); v != nil {
		s.Min = &v.Value
	}
	if v := get(r, "maximum"); v != nil {
		s.Max = &v.Value
	}
	if v := get(r, "minLength"); v != nil {
		s.MinLen = &v.Value
	}
	if v := get(r, "maxLength"); v != nil {
		s.MaxLen = &v.Value
	}
	s.Pattern = str(get(r, "pattern"))
	if ex := get(r, "example"); ex != nil && ex.Kind == yaml.ScalarNode {
		s.Example = ex.Value
	}
	if it := get(r, "items"); it != nil {
		c := d.schemaOf(it)
		s.Items = &c
	}
	return s
}

type opInfo struct {
	id, method, path string
	node             *yaml.Node
}

var (
	skipHeaders  = map[string]bool{"If-Match": true, "Idempotency-Key": true, "X-Request-Id": true}
	methodsOrder = []string{"get", "put", "post", "delete", "patch"}
	cliName      = regexp.MustCompile(`^([a-z]+(?:-[a-z]+)*) ([a-z]+(?:-[a-z]+)*)$`)
)

// Generate returns the gofmt'd source of operations_gen.go for the snapshot.
func Generate(raw []byte) ([]byte, error) {
	var root yaml.Node
	if err := yaml.Unmarshal(raw, &root); err != nil {
		return nil, err
	}
	d := &doc{root: &root}
	top := d.mapping(&root)
	version := str(get(get(top, "info"), "version"))
	sum := sha256.Sum256(raw)

	var ops []opInfo
	paths := get(top, "paths")
	for _, p := range keys(paths) {
		item := get(paths, p)
		for _, m := range methodsOrder {
			if o := get(item, m); o != nil {
				ops = append(ops, opInfo{id: str(get(o, "operationId")), method: strings.ToUpper(m), path: p, node: o})
			}
		}
	}
	sort.Slice(ops, func(i, j int) bool { return ops[i].id < ops[j].id })

	var b bytes.Buffer
	b.WriteString("// Code generated by tools/afsgen from api/openapi.yaml; DO NOT EDIT.\n\npackage spec\n\n")
	fmt.Fprintf(&b, "// ContractVersion is info.version of the snapshot.\nconst ContractVersion = %q\n\n", version)
	fmt.Fprintf(&b, "// SnapshotSHA256 is the checksum of the snapshot these descriptors were generated from.\nconst SnapshotSHA256 = %q\n\n", hex.EncodeToString(sum[:]))
	b.WriteString("var operations = []Operation{\n")
	for _, o := range ops {
		if err := d.emitOp(&b, o); err != nil {
			return nil, fmt.Errorf("%s: %w", o.id, err)
		}
	}
	b.WriteString("}\n")
	return format.Source(b.Bytes())
}

func (d *doc) emitOp(b *bytes.Buffer, o opInfo) error {
	n := o.node
	cli := str(get(n, "x-afs-cli"))
	service, command := "", ""
	if !strings.HasPrefix(cli, "none") {
		m := cliName.FindStringSubmatch(cli)
		if m == nil {
			return fmt.Errorf("x-afs-cli %q is not '<service> <kebab-verb-noun>'", cli)
		}
		service, command = m[1], m[2]
	}
	if get(n, "x-afs-cli") == nil {
		return fmt.Errorf("no x-afs-cli")
	}
	fmt.Fprintf(b, "{ID: %q, Method: %q, Path: %q, Service: %q, Command: %q, Summary: %q,\n", o.id, o.method, o.path, service, command, str(get(n, "summary")))
	fmt.Fprintf(b, "Permission: %q, Scope: %q, Priority: %q, CSV: %v,\n", str(get(n, "x-afs-permission")), str(get(n, "x-afs-scope")), str(get(n, "x-afs-priority")), str(get(n, "x-afs-csv")) == "true")
	if any := get(n, "x-afs-any-of"); any != nil {
		var q []string
		for _, c := range any.Content {
			q = append(q, strconv.Quote(c.Value))
		}
		b.WriteString("AnyOf: []string{" + strings.Join(q, ", ") + "},\n")
	}
	// parameters
	var params []string
	ifMatch, idem, hasMax, hasNext := "None", false, false, false
	if get(n, "parameters") == nil {
		return fmt.Errorf("no parameters")
	}
	for _, pn := range get(n, "parameters").Content {
		r, _ := d.resolve(pn)
		name, in := str(get(r, "name")), str(get(r, "in"))
		switch {
		case in == "header" && name == "If-Match":
			if str(get(r, "required")) == "true" {
				ifMatch = "Required"
			} else {
				ifMatch = "Optional"
			}
			continue
		case in == "header" && name == "Idempotency-Key":
			idem = true
			continue
		case in == "header" && skipHeaders[name]:
			continue
		case in == "header":
			return fmt.Errorf("unsupported header parameter %s", name)
		case in == "query" && name == "maxResults":
			hasMax = true
			continue
		case in == "query" && name == "nextToken":
			hasNext = true
			continue
		}
		s := d.schemaOf(get(r, "schema"))
		desc := str(get(r, "description"))
		params = append(params, fmt.Sprintf("{Name: %q, In: %q, Required: %v, Desc: %q, Example: %q, Schema: %s},\n",
			name, in, str(get(r, "required")) == "true", short(desc), s.Example, goSchema(s)))
	}
	if len(params) > 0 {
		b.WriteString("Params: []Param{\n" + strings.Join(params, "") + "},\n")
	}
	fmt.Fprintf(b, "IfMatch: IfMatch%s, Idempotent: %v, Paginated: %v,\n", ifMatch, idem, hasMax && hasNext)
	// body
	if rb := get(n, "requestBody"); rb != nil {
		content := get(rb, "content")
		mt := keys(content)[0]
		sn, name := d.resolve(get(get(content, mt), "schema"))
		fmt.Fprintf(b, "Body: &Body{Required: %v, MediaType: %q, SchemaName: %q, Fields: []Field{\n", str(get(rb, "required")) == "true", mt, name)
		req := map[string]bool{}
		if r := get(sn, "required"); r != nil {
			for _, c := range r.Content {
				req[c.Value] = true
			}
		}
		props := get(sn, "properties")
		for _, k := range keys(props) {
			s := d.schemaOf(get(props, k))
			fmt.Fprintf(b, "{Name: %q, Required: %v, WriteOnly: %v, ReadOnly: %v, Desc: %q, Schema: %s},\n", k, req[k], s.WriteOnly, s.ReadOnly, short(s.Desc), goSchema(s))
		}
		b.WriteString("}},\n")
	}
	// responses
	var success []int
	var produces []string
	resultKey := ""
	resp := get(n, "responses")
	for _, code := range keys(resp) {
		if !strings.HasPrefix(code, "2") {
			continue
		}
		c, _ := strconv.Atoi(code)
		success = append(success, c)
		r, _ := d.resolve(get(resp, code))
		content := get(r, "content")
		for _, mt := range keys(content) {
			if !contains(produces, mt) {
				produces = append(produces, mt)
			}
			if mt == "application/json" && resultKey == "" {
				sn, _ := d.resolve(get(get(content, mt), "schema"))
				props := get(sn, "properties")
				switch {
				case get(props, "rows") != nil:
					resultKey = "rows"
				case get(props, "items") != nil:
					resultKey = "items"
				}
			}
		}
	}
	sort.Ints(success)
	if hasMax && hasNext {
		fmt.Fprintf(b, "ResultKey: %q,\n", resultKey)
	}
	var sc []string
	for _, c := range success {
		sc = append(sc, strconv.Itoa(c))
	}
	b.WriteString("Success: []int{" + strings.Join(sc, ", ") + "},\n")
	if len(produces) > 0 {
		var q []string
		for _, p := range produces {
			q = append(q, strconv.Quote(p))
		}
		b.WriteString("Produces: []string{" + strings.Join(q, ", ") + "},\n")
	}
	b.WriteString("},\n")
	return nil
}

func contains(l []string, s string) bool {
	for _, x := range l {
		if x == s {
			return true
		}
	}
	return false
}

// short keeps help text to its first sentence, at most 140 bytes.
func short(s string) string {
	s = strings.Join(strings.Fields(s), " ")
	if i := strings.Index(s, ". "); i > 0 {
		s = s[:i+1]
	}
	if len(s) > 140 {
		s = s[:137] + "..."
	}
	return s
}

func goSchema(s schema) string {
	var parts []string
	parts = append(parts, fmt.Sprintf("Type: %q", s.Type))
	if s.Nullable {
		parts = append(parts, "Nullable: true")
	}
	if len(s.Enum) > 0 {
		q := make([]string, len(s.Enum))
		for i, e := range s.Enum {
			q[i] = strconv.Quote(e)
		}
		parts = append(parts, "Enum: []string{"+strings.Join(q, ", ")+"}")
	}
	if s.Min != nil {
		parts = append(parts, "Min: f64("+*s.Min+")")
	}
	if s.Max != nil {
		parts = append(parts, "Max: f64("+*s.Max+")")
	}
	if s.MinLen != nil {
		parts = append(parts, "MinLen: i64("+*s.MinLen+")")
	}
	if s.MaxLen != nil {
		parts = append(parts, "MaxLen: i64("+*s.MaxLen+")")
	}
	if s.Pattern != "" {
		parts = append(parts, fmt.Sprintf("Pattern: %q", s.Pattern))
	}
	if s.Items != nil {
		parts = append(parts, "Items: &"+goSchema(*s.Items))
	}
	if s.Ref != "" {
		parts = append(parts, fmt.Sprintf("Ref: %q", s.Ref))
	}
	return "Schema{" + strings.Join(parts, ", ") + "}"
}
