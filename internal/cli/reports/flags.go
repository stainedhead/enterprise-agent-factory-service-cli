package reports

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
)

// groupByPattern is the groupBy pattern: team, agent, model, system, day, hour, pool, tag:<key> and, for
// the collaboration report (API contract 4.8, v1.3), surface, reason and person, and route (v1.5.0) on usage and
// costs, comma separated (CLI-13,
// CLI-US-32.2). Whether a dimension suits a report kind is the API's decision: the CLI checks the shape only.
var groupByPattern = regexp.MustCompile(`^(team|agent|model|system|day|hour|pool|surface|reason|person|route|tag:[A-Za-z0-9_.:/=+@-]+)(,(team|agent|model|system|day|hour|pool|surface|reason|person|route|tag:[A-Za-z0-9_.:/=+@-]+))*$`)

// groupByDesc replaces the snapshot's description of groupBy, which the generator shortens.
const groupByDesc = "Comma-separated: team, agent, model, system, day, hour, tag:<key>, pool (ephemeral teams), surface, reason and person (collaboration), and route (usage and costs, v1.5.0)."

// tagOption is --tag key=value, repeated by listing several pairs; each becomes the filter tag.<key>=<value>.
var tagOption = cli.Option{Flag: "--tag", In: "local", Kind: cli.List, Elem: "string",
	Desc: "Filter by a tag, key=value (several allowed); sent as tag.<key>=<value>."}

// flags is the part of every report operation that is the same: --tag, and the checks that fail
// with exit 252 before any call.
func flags(op string) cli.Override {
	return cli.Override{
		Op:    op,
		Extra: []cli.Option{tagOption},
		Desc:  map[string]string{"groupBy": groupByDesc},
		Validate: func(p *cli.Parsed) error {
			if p.Has("range") && (p.Has("from") || p.Has("to")) {
				return fmt.Errorf("--range replaces --from and --to: give one or the other")
			}
			if p.Has("group-by") && !groupByPattern.MatchString(p.Value("group-by")) {
				return fmt.Errorf("argument --group-by: %q is not a comma-separated list of team, agent, model, system, day, hour, pool, surface, reason, person, route or tag:<key>", p.Value("group-by"))
			}
			for _, kv := range p.Values["tag"] {
				if k, _, ok := strings.Cut(kv, "="); !ok || k == "" {
					return fmt.Errorf("argument --tag: expected key=value, got %q", kv)
				}
			}
			return nil
		},
		Build: func(p *cli.Parsed, parts *cli.RequestParts) error {
			for _, kv := range p.Values["tag"] {
				k, v, _ := strings.Cut(kv, "=")
				parts.Query.Add("tag."+k, v)
			}
			return nil
		},
	}
}
