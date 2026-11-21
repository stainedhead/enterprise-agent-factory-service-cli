package cli

import (
	"context"
	"fmt"
	"sort"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/spec"
)

// LocalFunc runs a local command (configure, sso login, completion ...) that calls no API operation.
type LocalFunc func(ctx context.Context, inv *Invocation) int

// ExecFunc runs an API command: the generic executor.
type ExecFunc func(ctx context.Context, inv *Invocation) int

// HelpFunc renders help for the given positional words (none, a service, or a service and operation).
type HelpFunc func(inv *Invocation, words []string) int

// Command is one entry of the command table: a generated operation plus its overrides, or a local command.
type Command struct {
	Service string
	Name    string // "" for a service-level local command such as "configure"
	Op      *spec.Operation
	Local   LocalFunc
	Summary string // local commands only; operations use Op.Summary
	Options []Option
	Override
}

// Key is "service name".
func (c *Command) Key() string {
	if c.Name == "" {
		return c.Service
	}
	return c.Service + " " + c.Name
}

// Invocation is what a running command receives.
type Invocation struct {
	Env *Env
	Cmd *Command
	P   *Parsed
	Reg *Registry
}

// Registry is the command table: generated operations plus service overrides and local commands.
type Registry struct {
	// Exec and Help are wired by the composition root (internal/cli/services); cli cannot
	// import them without a cycle.
	Exec ExecFunc
	Help HelpFunc

	cmds      map[string]*Command
	byOp      map[string]*Command
	overrides map[string]*Override
	extra     map[string][]Option
	errs      []error
	final     bool
}

// NewRegistry builds a command for every operation that has a CLI name.
func NewRegistry(ops []spec.Operation) *Registry {
	r := &Registry{cmds: map[string]*Command{}, byOp: map[string]*Command{}, overrides: map[string]*Override{}, extra: map[string][]Option{}}
	for i := range ops {
		op := &ops[i]
		if op.Service == "" {
			continue
		}
		c := &Command{Service: op.Service, Name: op.Command, Op: op}
		r.cmds[c.Key()] = c
		r.byOp[op.ID] = c
	}
	return r
}

// Override registers (merges) an override for one operationId.
func (r *Registry) Override(ov Override) {
	if _, ok := r.byOp[ov.Op]; !ok {
		r.errs = append(r.errs, fmt.Errorf("override for unknown operation %q", ov.Op))
		return
	}
	cur, ok := r.overrides[ov.Op]
	if !ok {
		cur = &Override{Op: ov.Op}
		r.overrides[ov.Op] = cur
	}
	if err := cur.merge(ov); err != nil {
		r.errs = append(r.errs, err)
	}
}

// AddOptions appends options (such as --wait) to an operation's command.
func (r *Registry) AddOptions(opID string, opts ...Option) {
	r.extra[opID] = append(r.extra[opID], opts...)
}

// AddLocal registers a local command.
func (r *Registry) AddLocal(c *Command) {
	if _, dup := r.cmds[c.Key()]; dup {
		r.errs = append(r.errs, fmt.Errorf("command %q registered twice", c.Key()))
		return
	}
	r.cmds[c.Key()] = c
}

// Finalize applies overrides, derives options and reports registration errors.
func (r *Registry) Finalize() error {
	if len(r.errs) > 0 {
		return r.errs[0]
	}
	for id, c := range r.byOp {
		if ov := r.overrides[id]; ov != nil {
			c.Override = *ov
		}
		c.Options = DeriveOptions(c.Op, &c.Override)
		c.Options = append(c.Options, r.extra[id]...)
		seen := map[string]bool{}
		for _, o := range c.Options {
			if seen[o.Flag] {
				return fmt.Errorf("%s: option %s defined twice", c.Key(), o.Flag)
			}
			if _, isGlobal := lookupGlobal(o.Name()); isGlobal {
				return fmt.Errorf("%s: option %s collides with a global option", c.Key(), o.Flag)
			}
			seen[o.Flag] = true
		}
	}
	r.final = true
	return nil
}

// Lookup finds a command; ok is false when there is none.
func (r *Registry) Lookup(service, name string) (*Command, bool) {
	c, ok := r.cmds[service+" "+name]
	if name == "" {
		c, ok = r.cmds[service]
	}
	return c, ok
}

// ByOp returns the command of an operation.
func (r *Registry) ByOp(id string) *Command { return r.byOp[id] }

// Services lists service names, sorted.
func (r *Registry) Services() []string {
	seen := map[string]bool{}
	for _, c := range r.cmds {
		seen[c.Service] = true
	}
	out := make([]string, 0, len(seen))
	for s := range seen {
		out = append(out, s)
	}
	sort.Strings(out)
	return out
}

// Commands lists a service's commands sorted by name.
func (r *Registry) Commands(service string) []*Command {
	var out []*Command
	for _, c := range r.cmds {
		if c.Service == service {
			out = append(out, c)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

// All lists every command.
func (r *Registry) All() []*Command {
	out := make([]*Command, 0, len(r.cmds))
	for _, c := range r.cmds {
		out = append(out, c)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Key() < out[j].Key() })
	return out
}
