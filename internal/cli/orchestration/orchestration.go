package orchestration

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli/runtime"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/client"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/output"
)

func init() { registrations = append(registrations, registerOrchestration) }

const exTeam = "team_01JD2Q3R4S5T6V7W8X9Y0Z0001"

// Levels are the exposure levels of a surface (ORC-12), lowest first.
var Levels = []string{"off", "private", "public"}

// Columns are the table columns of the team's orchestration resource (OR-CLI-26). A surface the tool does not
// have (Buzz has neither) prints N/A; a hostname that is not served prints "-".
var Columns = []output.Column{
	{Header: "Tool", Path: "tool"},
	{Header: "WebUi", Path: "webUi.exposure", NAUnless: "webUi.supported"},
	{Header: "WebApi", Path: "webApi.exposure", NAUnless: "webApi.supported"},
	{Header: "Limit", Path: "limit"},
	{Header: "EffectiveWebUi", Path: "webUi.effective", NAUnless: "webUi.supported"},
	{Header: "EffectiveWebApi", Path: "webApi.effective", NAUnless: "webApi.supported"},
	{Header: "Hostname", Path: "webUi.hostname", Null: "-"},
	{Header: "Mode", Path: "registrationMode"},
	{Header: "Phase", Path: "phase"},
	{Header: "Version", Path: "version"},
}

// exposureOptions are the options of put-orchestration that stand for the webUi and webApi members.
var exposureOptions = []cli.Option{
	{Flag: "--web-ui-exposure", In: "local", Kind: cli.String, Enum: Levels,
		Desc: "Exposure of the Web UI: off, private (internal ALB) or public (internet-facing ALB). Both surfaces are sent; give both unless --cli-input-json has them."},
	{Flag: "--web-api-exposure", In: "local", Kind: cli.String, Enum: Levels,
		Desc: "Exposure of the Web API: off, private or public. Buzz has no Web API: only off is accepted for it."},
}

// buildExposure folds the two exposure options into the webUi and webApi members and refuses a body that
// lacks either: the API takes both surfaces on every write, and one call changes both (ORC-25).
func buildExposure(p *cli.Parsed, parts *cli.RequestParts) error {
	for _, s := range []struct{ flag, member string }{{"web-ui-exposure", "webUi"}, {"web-api-exposure", "webApi"}} {
		if p.Has(s.flag) {
			v := p.Value(s.flag)
			if !isLevel(v) {
				return fmt.Errorf("argument --%s: invalid choice %q (choose from %s)", s.flag, v, strings.Join(Levels, ", "))
			}
			b, _ := json.Marshal(map[string]string{"exposure": v})
			parts.Body.Set(s.member, b)
		}
		if _, ok := parts.Body.Get(s.member); !ok {
			return fmt.Errorf("the following arguments are required: --web-ui-exposure, --web-api-exposure (the API takes both surfaces in one write; --cli-input-json may give %s)", s.member)
		}
	}
	return nil
}

func isLevel(v string) bool {
	for _, l := range Levels {
		if l == v {
			return true
		}
	}
	return false
}

func registerOrchestration(r *cli.Registry) {
	r.Override(cli.Override{Op: "GetTeamOrchestration", Columns: Columns, After: AfterDescribe, Examples: []cli.Example{
		{Title: "Show a team's orchestration agent, its tool and what is effective", Command: "afs orchestration describe-orchestration --team-id " + exTeam},
		{Title: "Print the Web UI address", Command: "afs orchestration describe-orchestration --team-id " + exTeam + " --query webUi.url --output text"},
	}})
	r.Override(cli.Override{Op: "PutTeamOrchestration", Columns: Columns, After: AfterDescribe,
		Extra: exposureOptions, Build: buildExposure, Optional: []string{"webUi", "webApi"}, Hidden: []string{"webUi", "webApi"},
		Notes: []string{"Sets the exposure of both surfaces in one call: off (inside the namespace only), private (through the internal ALB from networks routed to the environment VPC) or public (through the internet-facing ALB). A level above the team's limit is refused; the limit is set by a service admin (afs teams update-team --max-collaboration-exposure)."},
		Examples: []cli.Example{
			{Title: "Serve the Web UI privately and keep the Web API off", Command: "afs orchestration put-orchestration --team-id " + exTeam + " --web-ui-exposure private --web-api-exposure off"},
			{Title: "Publish both surfaces, use managed registration and wait until it is ready", Command: "afs orchestration put-orchestration --team-id " + exTeam + " --web-ui-exposure public --web-api-exposure public --registration-mode managed --wait"},
		}})
	r.Override(cli.Override{Op: "ChangeCollaborationTool", Columns: Columns, After: AfterDescribe,
		Notes: []string{"The orchestration agent is snapshotted (kept 90 days), then started with the new tool and an empty data directory: no content moves between tools. Surfaces the new tool lacks are set to off, and connected workers are registered again. Buzz has no Web UI and no Web API."},
		Examples: []cli.Example{
			{Title: "Switch the team's collaboration tool to Buzz and wait until it is ready", Command: "afs orchestration change-collaboration-tool --team-id " + exTeam + " --tool buzz --wait"},
			{Title: "Switch back to Paperclip", Command: "afs orchestration change-collaboration-tool --team-id " + exTeam + " --tool paperclip"},
		}})
}

// AfterDescribe prints, for table and text output, one stderr line for each surface the tool does not have
// and one for each surface that is not effective at the level chosen, with the API's reasons as returned
// (OR-CLI-26). With json nothing extra is printed, and the exit code is never changed.
func AfterDescribe(inv *cli.Invocation, resp *client.Response) int {
	st, err := runtime.Settings(inv.Env, inv.P.G)
	if err != nil || (st.Output != "table" && st.Output != "text") {
		return -1
	}
	var doc struct {
		Tool   string  `json:"tool"`
		WebUI  surface `json:"webUi"`
		WebAPI surface `json:"webApi"`
	}
	if json.Unmarshal(resp.Body, &doc) != nil {
		return -1
	}
	for _, s := range []struct {
		name string
		s    surface
	}{{"Web UI", doc.WebUI}, {"Web API", doc.WebAPI}} {
		switch {
		case s.s.Supported != nil && !*s.s.Supported:
			_, _ = fmt.Fprintf(inv.Env.Stderr, "%s: N/A for %s\n", s.name, output.Sanitize(ToolName(doc.Tool)))
		case s.s.Effective != s.s.Exposure && len(s.s.Reasons) > 0:
			_, _ = fmt.Fprintf(inv.Env.Stderr, "Not effective: %s is %s, not %s (%s)\n", s.name, output.Sanitize(s.s.Effective), output.Sanitize(s.s.Exposure), output.Sanitize(strings.Join(s.s.Reasons, ", ")))
		}
	}
	return -1
}

type surface struct {
	Supported *bool    `json:"supported"`
	Exposure  string   `json:"exposure"`
	Effective string   `json:"effective"`
	Reasons   []string `json:"reasons"`
}

// ToolName is the display name of a collaboration tool: Paperclip or Buzz; an unknown tool is shown as given.
func ToolName(tool string) string {
	switch tool {
	case "paperclip":
		return "Paperclip"
	case "buzz":
		return "Buzz"
	}
	return tool
}
