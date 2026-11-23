package exec

import (
	"fmt"
	"io"
	"sort"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/client"
)

// hints is the one table of one-line hints for the reasons of API contract section 4.4b, 4.4c and 4.6
// (PC-CLI-D6, CLI-20). A hint names commands and permissions only: it contains no text supplied by the
// server, it adds words and never a decision, and a reason that is not in the table prints no hint.
// The reasons are the values of details.reason and of details.fields[].reason.
var hints = map[string]string{ //nolint:gosec // G101: the words "key" and "secret" in the hint texts are not credentials
	// section 4.4b and 4.4
	"collaboration_agent_exists":  "This team already has an orchestration agent; see afs orchestration describe-orchestration.",
	"collaboration_agent_missing": "This team has no orchestration agent yet; one is created with every team, and the 1.1.0 upgrade creates it for older teams (afs orchestration describe-orchestration).",
	"team_suspended":              "The team is suspended. A service admin or team-admin can resume it: afs teams resume-team --team-id <id>.",
	"managed_registration":        "Registration is managed; the reconciler records the id.",
	"hostname_held":               "That slug's collaboration hostname is held for 30 days after the deletion of a team that had one; choose another slug.",
	// validation fields of PutTeamCollaboration
	"collab_base_domain_unset":           "Ask a service admin to set orchestration.public_base_domain (afs admin put-setting).",
	"collab_public_endpoints_off":        "Ask a service admin to switch orchestration.public_endpoints on (afs admin put-setting).",
	"collaboration_not_allowed_for_team": "Ask a service admin to allow it for this team (afs teams update-team --max-collaboration-exposure).",
	"provider_without_collab_client":     "Ask a service admin to register a collaboration client on the team's identity provider (afs admin update-identity-provider --collab-client-id).",
	"collaboration_agent_not_ready":      "Wait until the collaboration agent is Ready (afs agents get-agent-status).",
	// section 4.6 secrets
	"managed_secret":           "This secret is managed by AFS or Paperclip and cannot be changed here.",
	"master_key_not_rotatable": "The master key cannot be rotated: Paperclip's stored identities would become undecryptable.",
	"reserved_name":            "That secret name is reserved for Paperclip; choose another name.",
	// section 4.4c the edge
	"surface_disabled":      "The Web API is off for this team. A team-admin can turn it on: afs orchestration put-orchestration --team-id <id> --web-api-exposure private|public --web-ui-exposure ...",
	"no_team_role":          "Your credential has no role on this team. A team-admin can add one: afs members add-member --role team_collaborator.",
	"permission_required":   "Ask a team-admin for a role that holds the permission named above.",
	"ambiguous_credentials": "Send one AFS credential only.",
	"malformed_credential":  "The AFS credential is malformed; check AFS_API_KEY or run afs sso login.",
	"unsupported_encoding":  "The request body encoding is not supported; send the body uncompressed.",
	"body_too_large":        "The request body is larger than the service limit (collab.max_body_mb); send a smaller body.",
	"token_invalid":         "Run afs sso login to sign in again.",
	"token_expired":         "Run afs sso login to sign in again.",
	"key_expired":           "Create a new API key (afs api-keys create-api-key) and set AFS_API_KEY.",
	"key_revoked":           "Create a new API key (afs api-keys create-api-key) and set AFS_API_KEY.",
	"session_expired":       "The Paperclip web session has expired; sign in again in the browser.",
	"paused":                "The team's Paperclip is paused. Resume it with afs agents resume-agent.",
	"starting":              "Paperclip is starting; try again in a minute.",
	"upstream_unavailable":  "Paperclip did not answer; see afs orchestration describe-orchestration.",
	// section 4.4d the orchestration agent (v1.4.0)
	"orchestration_agent_pending":     "The team's orchestration agent is still being created by the 1.1.0 upgrade; try again in a few minutes (afs orchestration describe-orchestration).",
	"orchestration_agent_exists":      "This team already has an orchestration agent; see afs orchestration describe-orchestration.",
	"orchestration_agent_required":    "The orchestration agent is deleted only with its team. Pause or restart it (afs agents pause-agent, afs agents restart-agent) or change its tool (afs orchestration change-collaboration-tool).",
	"orchestration_field_fixed":       "The orchestration agent's name, role, template, image, IAM identity and collaboration settings are fixed; change its tool with afs orchestration change-collaboration-tool.",
	"orchestration_runtime_unset":     "A service admin must set the runtime of that tool in the orchestration.tools setting (afs admin put-setting --key orchestration.tools).",
	"tool_not_enabled":                "That tool is not enabled for the service. A service admin enables it in the orchestration.tools setting (afs admin put-setting --key orchestration.tools).",
	"tool_unchanged":                  "The team already uses that tool.",
	"tool_change_in_progress":         "A tool change is already running; follow it with afs orchestration describe-orchestration.",
	"exposure_above_limit":            "That level is above the team's limit. A service admin can raise the limit: afs teams update-team --max-collaboration-exposure off|private|public.",
	"surface_not_supported_by_tool":   "This collaboration tool has no such surface (Buzz has neither a Web UI nor a Web API); use off.",
	"registration_mode_not_supported": "This collaboration tool registers workers itself; use --registration-mode managed or omit it.",
	"role_reserved":                   "Agents you create are workers; the team's orchestration agent is created with the team.",
	"name_reserved":                   "The name orchestrator is reserved for the team's orchestration agent; choose another name.",
	// section 4.4 and 4.2 model access (v1.4.0)
	"no_model_available":            "No model can be used as the default. Ask a service admin to enable and price one (afs models list-models).",
	"model_not_enabled_for_service": "That model is not enabled for the service. A service admin can enable it (afs admin update-model) or switch service access to all (afs admin put-service-model-access --access all).",
	"model_not_enabled_for_team":    "That model is not allowed for the team. A service admin can change the team's models (afs teams put-team-models --access all|selected).",
	"models_selection_empty":        "Name at least one model with --models, or use --models all.",
	"default_model_not_allowed":     "The default model must be one the agent may use; add it to --models or choose another --default-model-id.",
	"model_ids_not_allowed_for_all": "Model ids are not accepted with --access all; use --access selected, or leave out --model-ids.",
	"conflicting_fields":            "--model-id and --default-model-id name different models; give one of them (--model-id is deprecated).",
	"setting_renamed":               "That setting was renamed in 1.1.0: collab.base_domain is now orchestration.public_base_domain and collab.public_endpoints is orchestration.public_endpoints (afs admin get-settings lists every key).",
	"service_model_access_all":      "Service model access is all, so a model's enabled flag cannot change; switch it first with afs admin put-service-model-access --access selected.",
}

// Hint returns the hint of a reason, or "" for one that has none.
func Hint(reason string) string { return hints[reason] }

// HintReasons lists the reasons that have a hint, sorted (tests and docs).
func HintReasons() []string {
	out := make([]string, 0, len(hints))
	for k := range hints {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// printDetails writes what an API error's details add after the Request ID: one
// "field: <field> reason: <reason>" line per validation failure (CLI-US-24.7), then the hint of
// details.reason and of each distinct field reason. Nothing here decides anything.
func printDetails(w io.Writer, ae *client.APIError) {
	seen := map[string]bool{}
	fields := ae.Fields()
	for _, f := range fields {
		_, _ = fmt.Fprintf(w, "field: %s reason: %s\n", f.Field, f.Reason)
	}
	reasons := []string{ae.Reason()}
	for _, f := range fields {
		reasons = append(reasons, f.Reason)
	}
	for _, r := range reasons {
		if h := hints[r]; h != "" && !seen[h] {
			seen[h] = true
			_, _ = fmt.Fprintln(w, h)
		}
	}
}
