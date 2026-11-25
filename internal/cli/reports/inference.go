package reports

import (
	"fmt"
	"regexp"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
)

func init() { registrations = append(registrations, registerInference) }

// inferenceGroupBy is the groupBy pattern of the inference report (contract v1.5.0): route, model, agent, pool, team,
// hour and day, comma separated. Which dimension suits which report is the API's decision.
var inferenceGroupBy = regexp.MustCompile(`^(route|model|agent|pool|team|hour|day)(,(route|model|agent|pool|team|hour|day))*$`)

// validateInference checks, before any call, that --range is not combined with --from or --to and that
// --group-by is well formed.
func validateInference(p *cli.Parsed) error {
	if p.Has("range") && (p.Has("from") || p.Has("to")) {
		return fmt.Errorf("--range replaces --from and --to: give one or the other")
	}
	if p.Has("group-by") && !inferenceGroupBy.MatchString(p.Value("group-by")) {
		return fmt.Errorf("argument --group-by: %q is not a comma-separated list of route, model, agent, pool, team, hour or day", p.Value("group-by"))
	}
	return nil
}

func registerInference(r *cli.Registry) {
	r.Override(cli.Override{Op: "GetTeamInferenceReport", Dynamic: true, Validate: validateInference,
		Notes: []string{"Latency, time to first token, throughput, errors (upstream errors apart from AFS refusals) and cost, from the gateway's ledger; no prompt or completion text is ever reported."},
		Examples: []cli.Example{
			{Title: "The team's routes over the last day, by route and model", Command: "afs reports get-team-inference-report --team-id " + exTeam + " --range 24h --group-by route,model"},
			{Title: "One route by agent over a week, as CSV", Command: "afs reports get-team-inference-report --team-id " + exTeam + " --range 7d --route openai-prod --group-by agent --output csv"},
		}})
	r.Override(cli.Override{Op: "GetInferenceReport", Dynamic: true, Validate: validateInference,
		Notes: []string{"The service-wide inference report (platform roles); --team-id narrows it to one team."},
		Examples: []cli.Example{
			{Title: "Every route this month, by route", Command: "afs reports get-inference-report --range mtd --group-by route"},
			{Title: "One team's hourly TTFT on a route", Command: "afs reports get-inference-report --team-id " + exTeam + " --route openai-prod --range 24h --group-by hour"},
		}})
}
