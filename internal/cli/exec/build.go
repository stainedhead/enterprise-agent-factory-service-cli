// Package exec is the generic executor: it turns one parsed command into one API call (plus the
// sanctioned read and polls), paginates, renders and maps errors to exit codes.
package exec

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/secretinput"
)

// Built is a request ready to send.
type Built struct {
	Path        map[string]string
	Query       url.Values
	Body        *cli.Object // nil when the operation has no body
	IfMatch     string      // explicit version ("" when not given)
	IdemKey     string
	Wait        bool
	WaitTimeout time.Duration
}

// Sources are the files and streams Build may read.
type Sources struct {
	Stdin    io.Reader
	ReadFile func(string) ([]byte, error)
	// ReadPassword reads a secret from the terminal without echo; nil when there is no terminal.
	ReadPassword func(prompt string) ([]byte, error)
}

func (s Sources) readFile(p string) ([]byte, error) {
	if s.ReadFile != nil {
		return s.ReadFile(p)
	}
	return os.ReadFile(p) //nolint:gosec // user-selected input file
}

func usage(cmd *cli.Command, format string, a ...any) *cli.UsageError {
	return &cli.UsageError{Msg: fmt.Sprintf(format, a...), Usage: fmt.Sprintf("usage: afs %s [options]\n\nTo see help text, run: afs %s help\n", cmd.Key(), cmd.Key())}
}

// textOrFile returns the text of an option value that may be file://PATH.
func (s Sources) textOrFile(v string) (string, error) {
	if p, ok := strings.CutPrefix(v, "file://"); ok {
		b, err := s.readFile(p)
		if err != nil {
			return "", err
		}
		return string(b), nil
	}
	return v, nil
}

// parseInput reads --cli-input-json (inline JSON or file://PATH).
func parseInput(cmd *cli.Command, p *cli.Parsed, src Sources) (*cli.Object, error) {
	if p.G.CLIInputJSON == "" {
		return cli.NewObject(), nil
	}
	inline := !strings.HasPrefix(p.G.CLIInputJSON, "file://")
	text, err := src.textOrFile(p.G.CLIInputJSON)
	if err != nil {
		return nil, usage(cmd, "cannot read --cli-input-json: %v", err)
	}
	obj, err := cli.ParseObject([]byte(text))
	if err != nil {
		return nil, usage(cmd, "--cli-input-json must be a JSON object: %v", err)
	}
	if inline {
		for _, o := range cmd.Options {
			if o.Sensitive {
				if _, has := obj.Get(o.Wire); has {
					return nil, usage(cmd, "a secret value must not be given inline: put the JSON in a file and use --cli-input-json file://PATH")
				}
			}
		}
	}
	return obj, nil
}

// Build turns the parsed options and --cli-input-json into request parts (D-C6). Explicit
// options win over input members; a failure is a usage error (exit 252) before any call.
func Build(cmd *cli.Command, p *cli.Parsed, src Sources, profileTeamID string) (*Built, error) {
	op := cmd.Op
	input, err := parseInput(cmd, p, src)
	if err != nil {
		return nil, err
	}
	if cmd.Precheck != nil {
		if err := cmd.Precheck(p); err != nil {
			return nil, usage(cmd, "%v", err)
		}
	}
	b := &Built{Path: map[string]string{}, Query: url.Values{}}
	if op.Body != nil {
		b.Body = cli.NewObject()
	}
	paramWire := map[string]bool{}
	for _, o := range cmd.Options {
		if o.In == "path" || o.In == "query" {
			paramWire[o.Wire] = true
		}
	}
	// Input members that are not parameters form the body, in the order given.
	for _, k := range input.Keys() {
		v, _ := input.Get(k)
		switch {
		case k == "version":
			b.IfMatch = scalarText(v)
		case k == "Idempotency-Key":
			b.IdemKey = scalarText(v)
		case paramWire[k]:
		case op.Body == nil:
			return nil, usage(cmd, "--cli-input-json: %q is not a parameter of %s", k, cmd.Key())
		default:
			b.Body.Set(k, v)
		}
	}
	for _, o := range cmd.Options {
		switch o.In {
		case "header":
			switch o.Wire {
			case "If-Match":
				if p.Has(o.Name()) {
					b.IfMatch = strings.Trim(p.Value(o.Name()), `"`)
				}
			case "Idempotency-Key":
				if p.Has(o.Name()) {
					b.IdemKey = p.Value(o.Name())
				}
			}
			continue
		case "local":
			continue
		}
		vals, given := p.Values[o.Name()]
		if !given {
			if v, ok := input.Get(o.Wire); ok && (o.In == "path" || o.In == "query") {
				if err := setParam(b, o, scalarText(v)); err != nil {
					return nil, usage(cmd, "%v", err)
				}
			}
			if o.In == "path" && o.Wire == "teamId" && b.Path["teamId"] == "" && profileTeamID != "" {
				b.Path["teamId"] = profileTeamID // D-C10: profile team_id fills a path teamId only
			}
			continue
		}
		raw, text, err := convert(o, vals, src, p)
		if err != nil {
			return nil, usage(cmd, "argument %s: %v", o.Flag, err)
		}
		switch o.In {
		case "path", "query":
			if err := setParam(b, o, text); err != nil {
				return nil, usage(cmd, "%v", err)
			}
		case "body":
			if b.Body == nil {
				return nil, usage(cmd, "%s has no request body", cmd.Key())
			}
			b.Body.Set(o.Wire, raw)
		}
	}
	if p.Has("wait") && p.Value("wait") == "true" {
		b.Wait = true
	}
	if p.Has("wait-timeout") {
		d, err := time.ParseDuration(p.Value("wait-timeout"))
		if err != nil || d <= 0 {
			return nil, usage(cmd, "argument --wait-timeout: expected a duration such as 90s or 10m, got %q", p.Value("wait-timeout"))
		}
		b.WaitTimeout = d
	}
	if cmd.Build != nil {
		parts := &cli.RequestParts{Path: b.Path, Query: b.Query, Body: b.Body}
		if err := cmd.Build(p, parts); err != nil {
			return nil, usage(cmd, "%v", err)
		}
	}
	if cmd.Validate != nil {
		if err := cmd.Validate(p); err != nil {
			return nil, usage(cmd, "%v", err)
		}
	}
	var missing []string
	for _, o := range cmd.Options {
		if !o.Required {
			continue
		}
		switch o.In {
		case "path":
			if b.Path[o.Wire] == "" {
				missing = append(missing, o.Flag)
			}
		case "query":
			if _, ok := b.Query[o.Wire]; !ok {
				missing = append(missing, o.Flag)
			}
		case "body":
			if _, ok := b.Body.Get(o.Wire); !ok {
				if o.Sensitive {
					raw, err := promptSecret(src, o)
					if err != nil {
						return nil, usage(cmd, "argument %s: %v", o.Flag, err)
					}
					b.Body.Set(o.Wire, raw)
					continue
				}
				missing = append(missing, o.Flag)
			}
		}
	}
	if len(missing) > 0 {
		return nil, usage(cmd, "the following arguments are required: %s", strings.Join(missing, ", "))
	}
	return b, nil
}

// promptSecret asks for a missing required secret on the terminal (CLI-9); without a terminal
// it fails with the ways to give the value.
func promptSecret(src Sources, o cli.Option) (json.RawMessage, error) {
	r := secretinput.Reader{}
	if src.ReadPassword != nil {
		r.ReadPassword = func() ([]byte, error) { return src.ReadPassword("Secret value (input is hidden): ") }
	}
	b, err := r.Prompt()
	if err != nil {
		return nil, err
	}
	return json.Marshal(string(b))
}

func setParam(b *Built, o cli.Option, text string) error {
	if err := checkEnum(o, text); err != nil {
		return fmt.Errorf("argument %s: %w", o.Flag, err)
	}
	if o.Kind == cli.Int || o.Kind == cli.Number {
		if err := checkRange(o, text); err != nil {
			return fmt.Errorf("argument %s: %w", o.Flag, err)
		}
	}
	if o.In == "path" {
		b.Path[o.Wire] = text
	} else {
		b.Query.Set(o.Wire, text)
	}
	return nil
}

// scalarText renders a JSON scalar as plain text (strings unquoted).
func scalarText(raw json.RawMessage) string {
	var s string
	if json.Unmarshal(raw, &s) == nil {
		return s
	}
	return strings.TrimSpace(string(raw))
}

func checkEnum(o cli.Option, v string) error {
	if len(o.Enum) == 0 {
		return nil
	}
	for _, e := range o.Enum {
		if e == v {
			return nil
		}
	}
	return fmt.Errorf("invalid choice %q (choose from %s)", v, strings.Join(o.Enum, ", "))
}

func checkRange(o cli.Option, v string) error {
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return fmt.Errorf("expected a number, got %q", v)
	}
	if o.Min != nil && f < *o.Min {
		return fmt.Errorf("%s is below the minimum %v", v, *o.Min)
	}
	if o.Max != nil && f > *o.Max {
		return fmt.Errorf("%s is above the maximum %v", v, *o.Max)
	}
	return nil
}

// convert turns the raw values of an option into JSON (for bodies) and text (for parameters).
func convert(o cli.Option, vals []string, src Sources, p *cli.Parsed) (json.RawMessage, string, error) {
	first := vals[0]
	if o.Sensitive {
		b, err := secretinput.Reader{Stdin: src.Stdin, ReadFile: src.ReadFile}.Read(first)
		if err != nil {
			return nil, "", err
		}
		raw, _ := json.Marshal(string(b))
		return raw, string(b), nil
	}
	switch o.Kind {
	case cli.String, cli.Duration:
		if err := checkEnum(o, first); err != nil {
			return nil, "", err
		}
		raw, _ := json.Marshal(first)
		return raw, first, nil
	case cli.Int:
		if _, err := strconv.ParseInt(first, 10, 64); err != nil {
			return nil, "", fmt.Errorf("expected an integer, got %q", first)
		}
		if err := checkRange(o, first); err != nil {
			return nil, "", err
		}
		return json.RawMessage(first), first, nil
	case cli.Number:
		if err := checkRange(o, first); err != nil {
			return nil, "", err
		}
		return json.RawMessage(first), first, nil
	case cli.Bool:
		return json.RawMessage(first), first, nil
	case cli.List:
		items := make([]json.RawMessage, len(vals))
		for i, v := range vals {
			switch o.Elem {
			case "integer":
				if _, err := strconv.ParseInt(v, 10, 64); err != nil {
					return nil, "", fmt.Errorf("expected integers, got %q", v)
				}
				items[i] = json.RawMessage(v)
			case "number":
				if _, err := strconv.ParseFloat(v, 64); err != nil {
					return nil, "", fmt.Errorf("expected numbers, got %q", v)
				}
				items[i] = json.RawMessage(v)
			case "boolean":
				if v != "true" && v != "false" {
					return nil, "", fmt.Errorf("expected true or false, got %q", v)
				}
				items[i] = json.RawMessage(v)
			default:
				if err := checkEnum(o, v); err != nil {
					return nil, "", err
				}
				items[i], _ = json.Marshal(v)
			}
		}
		raw, _ := json.Marshal(items)
		return raw, strings.Join(vals, ","), nil
	case cli.TagPairs:
		if !strings.HasPrefix(first, "{") && !strings.HasPrefix(first, "[") && !strings.HasPrefix(first, "file://") {
			pairs := make([]map[string]string, 0, len(vals))
			for _, kv := range vals {
				k, v, ok := strings.Cut(kv, "=")
				if !ok || k == "" {
					return nil, "", fmt.Errorf("expected key=value pairs, JSON or file://PATH, got %q", kv)
				}
				pairs = append(pairs, map[string]string{"key": k, "value": v})
			}
			raw, _ := json.Marshal(pairs)
			return raw, strings.Join(vals, ","), nil
		}
		fallthrough
	case cli.JSON, cli.Any:
		text, err := src.textOrFile(first)
		if err != nil {
			return nil, "", err
		}
		var buf bytes.Buffer
		if json.Valid([]byte(text)) {
			if err := json.Compact(&buf, []byte(text)); err != nil {
				return nil, "", err
			}
			return buf.Bytes(), text, nil
		}
		if o.Kind == cli.JSON || o.Kind == cli.TagPairs {
			return nil, "", fmt.Errorf("expected JSON or file://PATH")
		}
		raw, _ := json.Marshal(text)
		return raw, text, nil
	}
	_ = p
	return nil, "", fmt.Errorf("unsupported option kind")
}
