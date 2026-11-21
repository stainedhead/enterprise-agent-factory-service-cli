package cli

import (
	"fmt"
	"strings"
)

// Parsed is the result of parsing a command line.
type Parsed struct {
	G          Globals
	Service    string
	Op         string
	Words      []string            // positional words (for help and local commands)
	Values     map[string][]string // option name (no dashes) -> raw values; booleans hold "true" or "false"
	Positional []string            // extra positional arguments of a local command
	Help       bool                // "help" word or --help
}

// Has reports whether an option was given.
func (p *Parsed) Has(name string) bool { _, ok := p.Values[name]; return ok }

// Value returns the first value of an option, or "".
func (p *Parsed) Value(name string) string {
	if v := p.Values[name]; len(v) > 0 {
		return v[0]
	}
	return ""
}

// UsageError is a command-line error (exit 252).
type UsageError struct {
	Msg   string
	Usage string
}

func (e *UsageError) Error() string { return e.Msg }

const rootUsage = `usage: afs [options] <service> <operation> [parameters]

To see help text, you can run:
  afs help
  afs <service> help
  afs <service> <operation> help
`

func usageErr(format string, a ...any) *UsageError {
	return &UsageError{Msg: fmt.Sprintf(format, a...), Usage: rootUsage}
}

// Parse splits args into globals, a command and its options. Global options may appear anywhere.
// It returns the command (nil for help, version or when only a service is given).
func (r *Registry) Parse(args []string) (*Parsed, *Command, error) {
	p := &Parsed{Values: map[string][]string{}}
	// Phase 1: find the service and operation words; before them only global options may appear.
	var consumed [2]int
	words := 0
	for i := 0; i < len(args) && words < 2; i++ {
		tok := args[i]
		if tok == "-h" {
			p.Help = true
			continue
		}
		if strings.HasPrefix(tok, "--") {
			name, _, hasEq := strings.Cut(tok[2:], "=")
			d, ok := lookupGlobal(name)
			if !ok {
				return p, nil, usageErr("Unknown options: %s (global options only before the operation)", tok)
			}
			if d.kind != gBool && !hasEq {
				i++
			}
			continue
		}
		consumed[words] = i
		words++
		p.Words = append(p.Words, tok)
	}
	// Remaining positional words (for help detection) are found in phase 2.
	var cmd *Command
	skip := map[int]bool{}
	if words >= 1 {
		skip[consumed[0]] = true
		p.Service = p.Words[0]
	}
	if words == 2 {
		skip[consumed[1]] = true
		p.Op = p.Words[1]
	}
	if p.Service == "help" {
		p.Help = true
	}
	if !p.Help && words >= 1 {
		if words == 2 && p.Op == "help" {
			p.Help = true
		} else {
			var ok bool
			if cmd, ok = r.Lookup(p.Service, p.Op); !ok {
				if c2, ok2 := r.Lookup(p.Service, ""); ok2 && c2.Local != nil && words == 2 {
					cmd, p.Positional = c2, []string{p.Op}
					p.Op = ""
				} else if p.Service != "" && len(r.Commands(p.Service)) == 0 {
					return p, nil, usageErr("argument service: Invalid choice: %q", p.Service)
				} else if words == 2 {
					return p, nil, &UsageError{Msg: fmt.Sprintf("argument operation: Invalid choice: %q for service %s", p.Op, p.Service), Usage: fmt.Sprintf("usage: afs %s <operation> [parameters]\n", p.Service)}
				}
			}
		}
	}
	// Phase 2: options.
	var opts []Option
	if cmd != nil {
		opts = cmd.Options
	}
	for i := 0; i < len(args); i++ {
		if skip[i] {
			continue
		}
		tok := args[i]
		if tok == "-h" {
			p.Help = true
			continue
		}
		if !strings.HasPrefix(tok, "--") {
			if tok == "help" {
				p.Help = true
				continue
			}
			if cmd != nil && cmd.Local != nil {
				p.Positional = append(p.Positional, tok)
				continue
			}
			return p, cmd, r.cmdUsage(cmd, "Unexpected argument: %q", tok)
		}
		name, val, hasEq := strings.Cut(tok[2:], "=")
		if d, ok := lookupGlobal(name); ok {
			if d.kind == gBool {
				if err := p.G.setGlobal(name, ""); err != nil {
					return p, cmd, r.cmdUsage(cmd, "%v", err)
				}
				continue
			}
			if !hasEq {
				v, ok := nextValue(args, &i)
				if !ok {
					return p, cmd, r.cmdUsage(cmd, "argument --%s: expected one argument", name)
				}
				val = v
			}
			if err := p.G.setGlobal(name, val); err != nil {
				return p, cmd, r.cmdUsage(cmd, "%v", err)
			}
			continue
		}
		o, negated, ok := findOption(opts, name)
		if !ok {
			return p, cmd, r.cmdUsage(cmd, "Unknown options: %s", tok)
		}
		switch o.Kind {
		case Bool:
			v := "true"
			if negated {
				v = "false"
			}
			if hasEq {
				return p, cmd, r.cmdUsage(cmd, "argument %s: takes no value", o.Flag)
			}
			p.Values[o.Name()] = []string{v}
		case List:
			var vals []string
			if hasEq {
				vals = append(vals, val)
			}
			for i+1 < len(args) && !strings.HasPrefix(args[i+1], "--") {
				i++
				vals = append(vals, args[i])
			}
			if len(vals) == 0 {
				return p, cmd, r.cmdUsage(cmd, "argument %s: expected at least one argument", o.Flag)
			}
			p.Values[o.Name()] = vals
		default:
			if !hasEq {
				v, ok := nextValue(args, &i)
				if !ok || strings.HasPrefix(v, "--") {
					return p, cmd, r.cmdUsage(cmd, "argument %s: expected one argument", o.Flag)
				}
				val = v
			}
			p.Values[o.Name()] = []string{val}
		}
	}
	if p.G.Help {
		p.Help = true
	}
	return p, cmd, nil
}

// nextValue consumes the argument after position *i.
func nextValue(args []string, i *int) (string, bool) {
	if *i+1 >= len(args) {
		return "", false
	}
	*i++
	return args[*i], true
}

// findOption matches --name or --no-name (for Bool options).
func findOption(opts []Option, name string) (Option, bool, bool) {
	for _, o := range opts {
		if o.Name() == name {
			return o, false, true
		}
	}
	if rest, ok := strings.CutPrefix(name, "no-"); ok {
		for _, o := range opts {
			if o.Kind == Bool && o.Name() == rest {
				return o, true, true
			}
		}
	}
	return Option{}, false, false
}

func (r *Registry) cmdUsage(c *Command, format string, a ...any) *UsageError {
	u := rootUsage
	if c != nil {
		u = fmt.Sprintf("usage: afs %s [options]\n\nTo see help text, run: afs %s help\n", c.Key(), c.Key())
	}
	return &UsageError{Msg: fmt.Sprintf(format, a...), Usage: u}
}
