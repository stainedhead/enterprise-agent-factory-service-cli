package cli

import (
	"strings"
	"unicode"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/spec"
)

// Kind is how an option's command-line value is read and sent.
type Kind int

// Option kinds (data-dictionary section 2).
const (
	String Kind = iota
	Int
	Number
	Bool // --x / --no-x
	List // space-separated values up to the next option
	JSON // object or array of objects: JSON text or file://
	Any  // JSON when valid, else a string
	Duration
	// TagPairs is an array of {key, value} objects given as `key=value ...` pairs (a list option) or as
	// JSON or file://PATH (tags).
	TagPairs
)

// Option is one command-specific option derived from the operation.
type Option struct {
	Flag      string // "--team-id"
	Wire      string // "teamId", "If-Match"
	In        string // path | query | header | body | local
	Kind      Kind
	Elem      string // List element type: string | integer | number | boolean
	Required  bool
	Enum      []string
	Min, Max  *float64
	Sensitive bool // secret values: only file://, fileb:// or -
	// OpenEnum says Enum lists the known values only: the API decides what else is valid, so another
	// value is sent as given (put-setting --key: the collab.* settings, CLI-US-31.3).
	OpenEnum bool
	Desc     string
	// Hidden options are accepted but not listed in help.
	Hidden bool
	// Deprecated, when set, names what replaces the option; help shows it instead of Desc and giving the
	// option prints one stderr line.
	Deprecated string
}

// Name is the flag without dashes.
func (o Option) Name() string { return strings.TrimPrefix(o.Flag, "--") }

// Kebab converts camelCase to kebab-case: teamId -> team-id.
func Kebab(s string) string {
	var b strings.Builder
	rs := []rune(s)
	for i, r := range rs {
		if unicode.IsUpper(r) {
			if i > 0 && (unicode.IsLower(rs[i-1]) || unicode.IsDigit(rs[i-1])) {
				b.WriteByte('-')
			}
			b.WriteRune(unicode.ToLower(r))
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}

func kindOf(s spec.Schema) (Kind, string) {
	switch s.Type {
	case "integer":
		return Int, ""
	case "number":
		return Number, ""
	case "boolean":
		return Bool, ""
	case "array":
		if s.Items != nil {
			switch s.Items.Type {
			case "string", "integer", "number", "boolean":
				return List, s.Items.Type
			}
		}
		return JSON, ""
	case "object":
		return JSON, ""
	case "any":
		return Any, ""
	}
	return String, ""
}

// DeriveOptions builds the option list of an operation from its descriptor (CLI-1, CLI-2):
// path and query parameters, writable body members, If-Match and Idempotency-Key, then the
// override's renames and extra options.
func DeriveOptions(op *spec.Operation, ov *Override) []Option {
	rename := map[string]string{}
	if ov != nil {
		rename = ov.Rename
	}
	flag := func(wire string) string {
		if r, ok := rename[wire]; ok {
			return "--" + r
		}
		return "--" + Kebab(wire)
	}
	var out []Option
	for _, p := range op.Params {
		k, elem := kindOf(p.Schema)
		out = append(out, Option{Flag: flag(p.Name), Wire: p.Name, In: p.In, Kind: k, Elem: elem, Required: p.Required,
			Enum: p.Schema.Enum, Min: p.Schema.Min, Max: p.Schema.Max, Desc: p.Desc})
	}
	if op.Body != nil {
		for _, f := range op.Body.Fields {
			if f.ReadOnly {
				continue
			}
			k, elem := kindOf(f.Schema)
			if ov != nil {
				if kk, ok := ov.Kinds[f.Name]; ok {
					k = kk
				}
			}
			if k == List && elem == "" {
				elem = "string" // a Kinds override that turns an object member into a list of words
			}
			o := Option{Flag: flag(f.Name), Wire: f.Name, In: "body", Kind: k, Elem: elem, Required: f.Required,
				Enum: f.Schema.Enum, Min: f.Schema.Min, Max: f.Schema.Max, Desc: f.Desc}
			if _, isGlobal := lookupGlobal(o.Name()); isGlobal {
				// A body member named like a global option (report schedule "query") gets a -json suffix.
				o.Flag += "-json"
			}
			if f.Schema.Ref == "SecretValue" {
				o.Flag, o.Sensitive, o.Kind = "--secret-value", true, String
			}
			out = append(out, o)
		}
	}
	if op.IfMatch != spec.IfMatchNone {
		out = append(out, Option{Flag: "--if-match", Wire: "If-Match", In: "header", Kind: String,
			Desc: "Apply only if the resource is at this version (safe update)."})
	}
	if op.Idempotent {
		out = append(out, Option{Flag: "--idempotency-key", Wire: "Idempotency-Key", In: "header", Kind: String,
			Desc: "Idempotency key (default: a new ULID per invocation)."})
	}
	if ov != nil {
		out = append(out, ov.Extra...)
		for i := range out {
			for _, w := range ov.Optional {
				if out[i].Wire == w {
					out[i].Required = false
				}
			}
			for _, w := range ov.OpenEnum {
				if out[i].Wire == w {
					out[i].OpenEnum = true
				}
			}
			for _, w := range ov.Hidden {
				if out[i].Wire == w {
					out[i].Hidden = true
				}
			}
			if d, ok := ov.Desc[out[i].Wire]; ok && out[i].In != "local" {
				out[i].Desc = d
			}
			if e, ok := ov.Enum[out[i].Wire]; ok && out[i].In != "local" {
				out[i].Enum = e
			}
			if d, ok := ov.DeprecatedOptions[out[i].Wire]; ok && out[i].In != "local" {
				out[i].Deprecated = d
			}
		}
	}
	return out
}
