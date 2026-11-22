package apikeys_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/clitest"
)

const secretKey = /* #nosec G101 -- documented example value */ "afsk_3f9aEXAMPLEEXAMPLEEXAMPLEEXAMPLE0000"

const created = `{"id":"key_01JD2Q3R4S5T6V7W8X9Y0Z0001","name":"ci-pipeline","role":"team_admin","keyPrefix":"afsk_3f9a","expiresAt":"2027-02-19T18:00:00Z","lastUsedAt":null,"secretKey":"` + secretKey + `"}`

func TestCommandsDone(t *testing.T) {
	clitest.CheckOperations(t, "keys", "ListApiKeys", "CreateApiKey", "DeleteApiKey")
}

// NFR security exception, D-C25: the secret is printed once, with the notice on stderr.
func TestCreateAPIKeyPrintsSecretOnce(t *testing.T) {
	for _, format := range []string{"json", "table"} {
		r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) { clitest.Reply(w, 201, created) },
			"api-keys", "create-api-key", "--team-id", "team_1", "--name", "ci-pipeline", "--role", "team_admin", "--expires-at", "2027-02-19T18:00:00Z", "--output", format)
		if r.Code != 0 || strings.Count(r.Out, secretKey) != 1 {
			t.Fatalf("%s: %+v", format, r)
		}
		if !strings.Contains(r.Err, "Store this key now; it cannot be shown again.") || strings.Contains(r.Err, secretKey) {
			t.Errorf("%s: stderr %q", format, r.Err)
		}
	}
	// the key is not echoed by any other command
	r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) {
		clitest.Reply(w, 200, `{"items":[{"id":"key_1","name":"ci","role":"team_admin","keyPrefix":"afsk_3f9a"}]}`)
	}, "api-keys", "list-api-keys", "--team-id", "team_1", "--output", "table")
	if r.Code != 0 || strings.Contains(r.Out+r.Err, "afsk_3f9aEXAMPLE") || strings.Contains(r.Err, "Store this key") {
		t.Fatalf("%+v", r)
	}
}

// NFR security: --debug never prints the secret key; the normal output still has it.
func TestCreateAPIKeyDebugRedacted(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) { clitest.Reply(w, 201, created) },
		"api-keys", "create-api-key", "--team-id", "team_1", "--name", "ci-pipeline", "--role", "team_admin", "--expires-at", "2027-02-19T18:00:00Z", "--debug")
	if r.Code != 0 || !strings.Contains(r.Out, secretKey) {
		t.Fatalf("%+v", r)
	}
	if strings.Contains(r.Err, secretKey) || strings.Contains(r.Err, "afsk_unit_test_key") {
		t.Errorf("debug log leaks a key: %s", r.Err)
	}
}

// D52: an idempotent replay has no secret; the notice says to create another key.
func TestCreateAPIKeyReplayNoSecret(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Idempotent-Replayed", "true")
		clitest.Reply(w, 201, `{"id":"key_1","name":"ci-pipeline","role":"team_admin","keyPrefix":"afsk_3f9a"}`)
	}, "api-keys", "create-api-key", "--team-id", "team_1", "--name", "ci-pipeline", "--role", "team_admin", "--expires-at", "2027-02-19T18:00:00Z", "--idempotency-key", "k1")
	if r.Code != 0 || !strings.Contains(r.Err, "This is a replay; the key's secret was shown only to the first request.") || strings.Contains(r.Err, "Store this key") {
		t.Fatalf("%+v", r)
	}
}

// D50: a viewer is refused by the API; the CLI prints it and exits 254.
func TestCreateAPIKeyViewerDenied(t *testing.T) {
	r := clitest.Run(t, clitest.Denied("api-keys:write"), "api-keys", "create-api-key", "--team-id", "team_1", "--name", "x", "--role", "team_viewer", "--expires-at", "2027-02-19T18:00:00Z")
	if r.Code != 254 || !strings.Contains(r.Err, "api-keys:write") {
		t.Fatalf("%+v", r)
	}
}

// CLI-CD-8: ApiKey has no version, so delete sends If-Match only when given.
func TestDeleteAPIKeyIfMatchOnlyWhenGiven(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) { clitest.Reply(w, 204, "") },
		"api-keys", "delete-api-key", "--team-id", "team_1", "--key-id", "key_1")
	if r.Code != 0 || r.Req.Header.Get("If-Match") != "" || r.Req.Method != "DELETE" {
		t.Fatalf("%+v", r)
	}
}

// CLI-US-32.1
func TestCreateApiKeyTeamCollaborator(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) {
		clitest.Reply(w, 201, `{"id":"key_1","name":"pipeline","role":"team_collaborator","secretKey":"afsk_x"}`)
	}, "api-keys", "create-api-key", "--team-id", "team_1", "--name", "pipeline", "--role", "team_collaborator", "--expires-at", "2027-02-19T18:00:00Z")
	if r.Code != 0 || !strings.Contains(string(r.ReqBody), `"role":"team_collaborator"`) {
		t.Errorf("%d %s %s", r.Code, r.ReqBody, r.Err)
	}
	h := clitest.Run(t, nil, "api-keys", "create-api-key", "help")
	if !strings.Contains(h.Out, "team_admin") || !strings.Contains(h.Out, "team_collaborator") || !strings.Contains(h.Out, "team_viewer") {
		t.Errorf("help must list the roles: %q", h.Out)
	}
}
