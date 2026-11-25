package exec

import (
	"bytes"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/client"
)

// contractReasons are the details.reason and details.fields[].reason values of API contract
// section 4.4b, 4.4c and 4.6 (CLI-US-24.7, 5.2 hints).
var contractReasons = []string{
	"collaboration_agent_exists", "collaboration_agent_missing", "team_suspended", "managed_secret",
	"master_key_not_rotatable", "reserved_name", "managed_registration", "surface_disabled", "no_team_role",
	"permission_required", "ambiguous_credentials", "malformed_credential", "body_too_large",
	"unsupported_encoding", "token_invalid", "token_expired", "key_expired", "key_revoked", "session_expired",
	"paused", "starting", "upstream_unavailable", "hostname_held",
	"collab_base_domain_unset", "collab_public_endpoints_off", "collaboration_not_allowed_for_team",
	"provider_without_collab_client", "collaboration_agent_not_ready",
	// API contract section 4.4d and section 4.4 (v1.4.0): ORC-29, MOD-8, CLI-US-26 to CLI-US-28
	"orchestration_agent_pending", "orchestration_agent_exists", "orchestration_agent_required", "orchestration_field_fixed",
	"orchestration_runtime_unset", "tool_not_enabled", "tool_unchanged", "tool_change_in_progress", "exposure_above_limit",
	"surface_not_supported_by_tool", "registration_mode_not_supported", "role_reserved", "name_reserved",
	"no_model_available", "model_not_enabled_for_service", "model_not_enabled_for_team", "models_selection_empty",
	"default_model_not_allowed", "model_ids_not_allowed_for_all", "conflicting_fields", "setting_renamed", "service_model_access_all",
}

// modelEndpointReasons are the new error reasons of API contract section 4.13 (v1.5.0).
var modelEndpointReasons = []string{"slug_taken", "slug_reserved", "slug_invalid", "base_url_not_https", "provider_disabled",
	"no_fallback_model", "endpoint_base_url_placeholder", "needs_service_admin", "team_endpoints_not_allowed", "endpoint_enabled"}

// mcpReasons are the new error reasons of API contract v1.6.0 (governed MCP servers).
var mcpReasons = []string{"url_not_https", "host_denied", "image_not_pinned", "image_unsigned", "registry_not_allowed", "auth_kind_unsupported", "service_account_required", "audience_not_allowed", "audience_reserved", "capture_window_too_long", "mcp_server_not_enabled_for_team", "invalid_policy_rule", "team_servers_not_allowed", "approval_not_permitted_by_policy", "not_team_server", "server_in_use", "server_not_published", "server_not_approved", "team_gateway_required", "invalid_review_transition", "federated_identity_limit", "governed_by_mcp_server", "capture_active"}

// CLI-US-24.7, PC-CLI-D6, 5.2 hints: every reason of the contract has a hint, the table has no
// reason the contract does not name, a hint is one line, and an unknown reason prints none.
func TestCollaborationReasonHints(t *testing.T) {
	want := map[string]bool{}
	for _, r := range append(append(append([]string{}, contractReasons...), modelEndpointReasons...), mcpReasons...) {
		want[r] = true
		h := Hint(r)
		if h == "" {
			t.Errorf("reason %q has no hint", r)
		}
		if strings.ContainsAny(h, "\n\r\t\x1b") || strings.TrimSpace(h) != h {
			t.Errorf("hint of %q is not one clean line: %q", r, h)
		}
	}
	for _, r := range HintReasons() {
		if !want[r] {
			t.Errorf("hint for %q, a reason the contract does not name", r)
		}
	}
	if Hint("something_new") != "" || Hint("") != "" {
		t.Error("an unknown reason must print no hint")
	}
}

func apiErr(status int, code, details string) *client.APIError {
	return &client.APIError{Status: status, Code: code, Message: "the message", RequestID: "req_9",
		Header: http.Header{}, Details: []byte(details)}
}

// CLI-US-24.7: the field lines follow the Request ID, then the hints; exit 254.
func TestReportErrorPrintsFieldsAndHints(t *testing.T) {
	var b bytes.Buffer
	e := apiErr(400, "ValidationError", `{"fields":[{"field":"webUiEnabled","reason":"collab_base_domain_unset"},{"field":"restApiEnabled","reason":"collab_base_domain_unset"},{"field":"x","reason":"unknown_reason"}]}`)
	if code := ReportError(&b, "PutTeamCollaboration", "default", e); code != 254 {
		t.Fatalf("exit %d", code)
	}
	want := "Request ID: req_9\nfield: webUiEnabled reason: collab_base_domain_unset\nfield: restApiEnabled reason: collab_base_domain_unset\nfield: x reason: unknown_reason\nAsk a service admin to set orchestration.public_base_domain (afs admin put-setting).\n"
	if !strings.HasSuffix(b.String(), want) {
		t.Errorf("got %q", b.String())
	}
	if strings.Count(b.String(), "Ask a service admin") != 1 {
		t.Error("the same hint must print once")
	}
}

func TestReportErrorReasonHintAndUnknownReason(t *testing.T) {
	var b bytes.Buffer
	ReportError(&b, "PutTeamCollaboration", "default", apiErr(409, "Conflict", `{"reason":"team_suspended"}`))
	if !strings.Contains(b.String(), Hint("team_suspended")) {
		t.Errorf("%q", b.String())
	}
	b.Reset()
	if code := ReportError(&b, "X", "default", apiErr(409, "Conflict", `{"reason":"brand_new"}`)); code != 254 || strings.Contains(b.String(), "brand_new") {
		t.Errorf("exit %d %q", code, b.String())
	}
	b.Reset()
	ReportError(&b, "X", "default", errors.New("plain"))
	if strings.Contains(b.String(), "field:") {
		t.Error("no details for a plain error")
	}
}

// FR-014 (v1.5.0): every new error reason of section 4.13 has a hint that names a command or what to change.
func TestModelEndpointHints(t *testing.T) {
	for _, reason := range modelEndpointReasons {
		if Hint(reason) == "" {
			t.Errorf("%s has no hint", reason)
		}
	}
}
