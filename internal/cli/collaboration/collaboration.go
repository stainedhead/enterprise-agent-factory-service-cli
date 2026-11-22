package collaboration

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli/runtime"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/client"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/output"
)

func init() { registrations = append(registrations, registerCollaboration) }

const (
	exTeam  = "team_01JD2Q3R4S5T6V7W8X9Y0Z0001"
	exAgent = "agt_01JD2Q3R4S5T6V7W8X9Y0Z0001"
)

// collaborationColumns are the table columns of the team's collaboration resource (CLI-US-24.1); a
// hostname that is not set yet (collab.base_domain unset) prints "-".
var collaborationColumns = []output.Column{
	{Header: "Hostname", Path: "hostname", Null: "-"}, {Header: "WebUi", Path: "webUiEnabled"}, {Header: "RestApi", Path: "restApiEnabled"},
	{Header: "EffectiveWebUi", Path: "effective.webUi"}, {Header: "EffectiveRestApi", Path: "effective.restApi"},
	{Header: "Mode", Path: "registrationMode"}, {Header: "Phase", Path: "status.phase"}, {Header: "Version", Path: "version"},
}

// linkColumns are the columns of one agent's link (CLI-US-26.1).
var linkColumns = []output.Column{
	{Header: "AgentId", Path: "agentId"}, {Header: "Name", Path: "agentName"}, {Header: "Connect", Path: "connect"},
	{Header: "PaperclipAgentId", Path: "paperclipAgentId", Null: "-"}, {Header: "LinkState", Path: "linkState"},
	{Header: "LastSyncedAt", Path: "lastSyncedAt", Null: "-"},
}

// listLinkColumns add the reconciler's last error, truncated by the table (CLI-US-26.4).
var listLinkColumns = append(append([]output.Column{}, linkColumns...), output.Column{Header: "LastError", Path: "lastError", Null: "-"})

func registerCollaboration(r *cli.Registry) {
	r.Override(cli.Override{Op: "GetTeamCollaboration", Columns: collaborationColumns, After: afterDescribe, Examples: []cli.Example{
		{Title: "Show a team's Paperclip endpoint and what is effective", Command: "afs collaboration describe-collaboration --team-id " + exTeam},
		{Title: "Print the public hostname only", Command: "afs collaboration describe-collaboration --team-id " + exTeam + " --query hostname --output text"},
	}})
	r.Override(cli.Override{Op: "PutTeamCollaboration", Columns: collaborationColumns, Examples: []cli.Example{
		{Title: "Publish the web UI and keep the REST API off", Command: "afs collaboration put-collaboration --team-id " + exTeam + " --web-ui-enabled --no-rest-api-enabled"},
		{Title: "Turn both surfaces on, use managed registration and wait until it is ready", Command: "afs collaboration put-collaboration --team-id " + exTeam + " --web-ui-enabled --rest-api-enabled --registration-mode managed --wait"},
	}})
	r.Override(cli.Override{Op: "ListCollaborationLinks", Columns: listLinkColumns, Examples: []cli.Example{
		{Title: "List the link state of every worker agent", Command: "afs collaboration list-collaboration-links --team-id " + exTeam},
		{Title: "List the links that failed", Command: "afs collaboration list-collaboration-links --team-id " + exTeam + " --link-state failed"},
	}})
	r.Override(cli.Override{Op: "ConnectAgentToCollaboration", Columns: linkColumns, Examples: []cli.Example{
		{Title: "Connect an agent to the team's Paperclip", Command: "afs collaboration connect-agent --team-id " + exTeam + " --agent-id " + exAgent},
		{Title: "Connect an agent in manual registration mode, naming its Paperclip agent", Command: "afs collaboration connect-agent --team-id " + exTeam + " --agent-id " + exAgent + " --paperclip-agent-id ag_8f2c1d"},
	}})
	r.Override(cli.Override{Op: "DisconnectAgentFromCollaboration", Columns: linkColumns, Examples: []cli.Example{
		{Title: "Disconnect an agent from the team's Paperclip", Command: "afs collaboration disconnect-agent --team-id " + exTeam + " --agent-id " + exAgent},
	}})
}

// afterDescribe prints, for table and text output, one stderr line per surface that is not effective
// while the response lists reasons (CLI-US-24.2). With json nothing extra is printed; the exit code
// is never changed. The reason values are printed as returned (the contract's four values).
func afterDescribe(inv *cli.Invocation, resp *client.Response) int {
	st, err := runtime.Settings(inv.Env, inv.P.G)
	if err != nil || (st.Output != "table" && st.Output != "text") {
		return -1
	}
	var doc struct {
		Effective struct {
			WebUI   bool     `json:"webUi"`
			RestAPI bool     `json:"restApi"`
			Reasons []string `json:"reasons"`
		} `json:"effective"`
	}
	if json.Unmarshal(resp.Body, &doc) != nil || len(doc.Effective.Reasons) == 0 {
		return -1
	}
	reasons := output.Sanitize(strings.Join(doc.Effective.Reasons, ", "))
	if !doc.Effective.WebUI {
		_, _ = fmt.Fprintf(inv.Env.Stderr, "Not effective: web UI (%s)\n", reasons)
	}
	if !doc.Effective.RestAPI {
		_, _ = fmt.Fprintf(inv.Env.Stderr, "Not effective: REST API (%s)\n", reasons)
	}
	return -1
}
