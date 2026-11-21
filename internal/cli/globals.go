package cli

import (
	"fmt"
	"strconv"
)

// Exit codes (CLI-8, data-dictionary section 7).
const (
	ExitOK          = 0
	ExitInterrupted = 130
	ExitUsage       = 252
	ExitService     = 254
	ExitOther       = 255
)

// Valid values of the --output and --color options.
var (
	OutputFormats = []string{"json", "yaml", "table", "text", "csv"}
	ColorModes    = []string{"on", "off", "auto"}
)

// Globals are the AWS-style options accepted on every command (CLI-2).
type Globals struct {
	Profile          string
	EndpointURL      string
	Output           string
	Query            string
	NoPaginate       bool
	MaxItems         int
	PageSize         int
	StartingToken    string
	CLIInputJSON     string
	GenerateSkeleton bool
	Debug            bool
	NoCLIPager       bool
	Color            string
	ConnectTimeout   int // seconds
	ReadTimeout      int // seconds
	Version          bool
	Help             bool
	// Set records which options were given explicitly (a zero value can be intentional).
	Set map[string]bool
}

type globalKind int

const (
	gString globalKind = iota
	gInt
	gBool
)

type globalDef struct {
	name string
	kind globalKind
	desc string
}

// globalDefs lists the global options in help order.
var globalDefs = []globalDef{
	{"profile", gString, "Use a specific profile from ~/.afs/config."},
	{"endpoint-url", gString, "Override the API endpoint URL."},
	{"output", gString, "Output format: json, yaml, table, text or csv."},
	{"query", gString, "JMESPath query applied to the result before rendering."},
	{"no-paginate", gBool, "Return only the first page of results."},
	{"max-items", gInt, "Stop after this many items and print NextToken."},
	{"page-size", gInt, "Items per request (maxResults), 1 to 100."},
	{"starting-token", gString, "Resume from a NextToken printed by --max-items."},
	{"cli-input-json", gString, "Read parameters from a JSON document or file://path."},
	{"generate-cli-skeleton", gBool, "Print a JSON skeleton of the input and make no call."},
	{"debug", gBool, "Write a redacted request log to stderr."},
	{"no-cli-pager", gBool, "Do not send output through a pager."},
	{"color", gString, "Colored output: on, off or auto."},
	{"cli-connect-timeout", gInt, "Connect timeout in seconds (0 for none; default 30)."},
	{"cli-read-timeout", gInt, "Read timeout in seconds (0 for none; default 60)."},
	{"version", gBool, "Print the afs version and exit."},
}

// GlobalHelp returns the name and description of every global option, for help text.
func GlobalHelp() [][2]string {
	out := make([][2]string, 0, len(globalDefs))
	for _, d := range globalDefs {
		out = append(out, [2]string{"--" + d.name, d.desc})
	}
	return out
}

func lookupGlobal(name string) (globalDef, bool) {
	for _, d := range globalDefs {
		if d.name == name {
			return d, true
		}
	}
	if name == "help" {
		return globalDef{name: "help", kind: gBool}, true
	}
	return globalDef{}, false
}

func contains(list []string, s string) bool {
	for _, x := range list {
		if x == s {
			return true
		}
	}
	return false
}

// setGlobal stores one parsed global option, validating command-line values (usage errors, 252).
func (g *Globals) setGlobal(name, val string) error {
	if g.Set == nil {
		g.Set = map[string]bool{}
	}
	g.Set[name] = true
	asInt := func() (int, error) {
		n, err := strconv.Atoi(val)
		if err != nil || n < 0 {
			return 0, fmt.Errorf("argument --%s: expected a non-negative integer, got %q", name, val)
		}
		return n, nil
	}
	var err error
	switch name {
	case "profile":
		g.Profile = val
	case "endpoint-url":
		g.EndpointURL = val
	case "output":
		if !contains(OutputFormats, val) {
			return fmt.Errorf("argument --output: invalid choice %q (choose from json, yaml, table, text, csv)", val)
		}
		g.Output = val
	case "query":
		g.Query = val
	case "no-paginate":
		g.NoPaginate = true
	case "max-items":
		g.MaxItems, err = asInt()
	case "page-size":
		g.PageSize, err = asInt()
	case "starting-token":
		g.StartingToken = val
	case "cli-input-json":
		g.CLIInputJSON = val
	case "generate-cli-skeleton":
		g.GenerateSkeleton = true
	case "debug":
		g.Debug = true
	case "no-cli-pager":
		g.NoCLIPager = true
	case "color":
		if !contains(ColorModes, val) {
			return fmt.Errorf("argument --color: invalid choice %q (choose from on, off, auto)", val)
		}
		g.Color = val
	case "cli-connect-timeout":
		g.ConnectTimeout, err = asInt()
	case "cli-read-timeout":
		g.ReadTimeout, err = asInt()
	case "version":
		g.Version = true
	case "help":
		g.Help = true
	}
	return err
}
