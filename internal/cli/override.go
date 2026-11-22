package cli

import (
	"fmt"
	"net/url"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/client"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/output"
)

// Example is one usage example shown in help; Command starts with "afs ".
type Example struct{ Title, Command string }

// RequestParts is what a Build hook may adjust before the request is sent.
type RequestParts struct {
	Path  map[string]string // path parameter values by wire name
	Query url.Values
	Body  *Object // nil when the operation has no body
}

// BuildFunc turns parsed options into request parts (shorthands such as --repos).
type BuildFunc func(p *Parsed, parts *RequestParts) error

// ValidateFunc is a client-side check that fails with exit 252 before any call.
type ValidateFunc func(p *Parsed) error

// AfterFunc runs after a successful response; it may print extras and return an exit code
// (use -1 to keep the executor's).
type AfterFunc func(inv *Invocation, resp *client.Response) int

// Override carries what one service adds to an operation beyond the generic behaviour.
// Several registrations for one operationId are merged (D-C32).
type Override struct {
	Op       string
	Examples []Example
	Rename   map[string]string // wire name -> option name without dashes
	// Kinds overrides how a body member is read from the command line (wire name -> kind).
	Kinds map[string]Kind
	// ResultKey names the member of the response that holds the rows of a table, when it is not
	// the operation's list member (for example "tags").
	ResultKey string
	Extra     []Option
	Build     BuildFunc
	Validate  ValidateFunc
	// Precheck runs before any option value or file is read, so it can refuse input by its size or
	// name without reading it.
	Precheck ValidateFunc
	Columns  []output.Column
	// Dynamic says the table columns follow the response (report rows are shaped by --group-by),
	// so the operation is done without a fixed column list.
	Dynamic bool
	After   AfterFunc
	// Optional lists wire names whose option is not required on the command line because a Build
	// hook fills it from another option (update-repo's --repo host/owner/repo).
	Optional []string
	// Accept is the media type to ask for when the operation answers in something other than JSON
	// (application/yaml for a kubeconfig); the response is then given to Handle.
	Accept string
	// Handle takes over a successful response instead of the generic rendering (kubeconfig).
	Handle AfterFunc
}

// merge folds b into a; two registrations setting Columns, After or the same Rename key conflict.
func (a *Override) merge(b Override) error {
	a.Examples = append(a.Examples, b.Examples...)
	a.Dynamic = a.Dynamic || b.Dynamic
	a.Extra = append(a.Extra, b.Extra...)
	a.Optional = append(a.Optional, b.Optional...)
	if b.ResultKey != "" {
		a.ResultKey = b.ResultKey
	}
	for k, v := range b.Kinds {
		if a.Kinds == nil {
			a.Kinds = map[string]Kind{}
		}
		a.Kinds[k] = v
	}
	for k, v := range b.Rename {
		if _, dup := a.Rename[k]; dup {
			return fmt.Errorf("%s: rename of %q registered twice", b.Op, k)
		}
		if a.Rename == nil {
			a.Rename = map[string]string{}
		}
		a.Rename[k] = v
	}
	if b.Columns != nil {
		if a.Columns != nil {
			return fmt.Errorf("%s: Columns registered twice", b.Op)
		}
		a.Columns = b.Columns
	}
	if b.Accept != "" {
		a.Accept = b.Accept
	}
	if b.Handle != nil {
		if a.Handle != nil {
			return fmt.Errorf("%s: Handle registered twice", b.Op)
		}
		a.Handle = b.Handle
	}
	if b.After != nil {
		if a.After != nil {
			return fmt.Errorf("%s: After registered twice", b.Op)
		}
		a.After = b.After
	}
	if b.Build != nil {
		prev, next := a.Build, b.Build
		if prev == nil {
			a.Build = next
		} else {
			a.Build = func(p *Parsed, parts *RequestParts) error {
				if err := prev(p, parts); err != nil {
					return err
				}
				return next(p, parts)
			}
		}
	}
	if b.Precheck != nil {
		prev, next := a.Precheck, b.Precheck
		if prev == nil {
			a.Precheck = next
		} else {
			a.Precheck = func(p *Parsed) error {
				if err := prev(p); err != nil {
					return err
				}
				return next(p)
			}
		}
	}
	if b.Validate != nil {
		prev, next := a.Validate, b.Validate
		if prev == nil {
			a.Validate = next
		} else {
			a.Validate = func(p *Parsed) error {
				if err := prev(p); err != nil {
					return err
				}
				return next(p)
			}
		}
	}
	return nil
}
