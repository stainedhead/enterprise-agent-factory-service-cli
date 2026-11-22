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
	"collaboration_agent_exists":  "This team already has a collaboration agent; see afs agents list-agents --role collaboration.",
	"collaboration_agent_missing": "This team has no collaboration agent. Create one with: afs agents create-agent --role collaboration ...",
	"team_suspended":              "The team is suspended. A service admin or team-admin can resume it: afs teams resume-team --team-id <id>.",
	"managed_registration":        "Registration is managed; the reconciler records the id.",
	"hostname_held":               "That slug's collaboration hostname is held for 30 days after the deletion of a team that had one; choose another slug.",
	// validation fields of PutTeamCollaboration
	"collab_base_domain_unset":           "Ask a service admin to set collab.base_domain (afs admin put-setting).",
	"collab_public_endpoints_off":        "Ask a service admin to switch collab.public_endpoints on (afs admin put-setting).",
	"collaboration_not_allowed_for_team": "Ask a service admin to allow it for this team (afs teams update-team --may-expose-collaboration).",
	"provider_without_collab_client":     "Ask a service admin to register a collaboration client on the team's identity provider (afs admin update-identity-provider --collab-client-id).",
	"collaboration_agent_not_ready":      "Wait until the collaboration agent is Ready (afs agents get-agent-status).",
	// section 4.6 secrets
	"managed_secret":           "This secret is managed by AFS or Paperclip and cannot be changed here.",
	"master_key_not_rotatable": "The master key cannot be rotated: Paperclip's stored identities would become undecryptable.",
	"reserved_name":            "That secret name is reserved for Paperclip; choose another name.",
	// section 4.4c the edge
	"surface_disabled":      "The REST API is switched off for this team. A team-admin can turn it on: afs collaboration put-collaboration --team-id <id> --rest-api-enabled ...",
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
	"upstream_unavailable":  "Paperclip did not answer; see afs collaboration describe-collaboration.",
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
