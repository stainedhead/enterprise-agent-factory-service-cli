package admin_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/clitest"
)

func TestCommandsDoneServiceAccess(t *testing.T) {
	clitest.CheckOperations(t, "serviceaccess", "ListServiceApiKeys", "CreateServiceApiKey", "DeleteServiceApiKey",
		"ListSecurityOperators", "AddSecurityOperator", "RemoveSecurityOperator")
}

const createdKey = `{"id":"key_1","name":"cve","role":"security_operator","keyPrefix":"afsk_ab","expiresAt":"2027-02-19T18:00:00Z","secretKey":"afsk_SECRET_ONCE"}`

// NFR security, D-C25: create-api-key prints the secret once and says to store it.
func TestCreateServiceAPIKeySecretOnce(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) { clitest.Reply(w, 201, createdKey) },
		"admin", "create-api-key", "--name", "cve", "--role", "security_operator", "--expires-at", "2027-02-19T18:00:00Z")
	if r.Code != 0 || strings.Count(r.Out, "afsk_SECRET_ONCE") != 1 || !strings.Contains(r.Err, "cannot be shown again") || strings.Contains(r.Err, "afsk_SECRET_ONCE") {
		t.Fatalf("%+v", r)
	}
	l := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) {
		clitest.Reply(w, 200, `{"items":[{"id":"key_1","name":"cve","role":"security_operator","keyPrefix":"afsk_ab"}]}`)
	}, "admin", "list-api-keys", "--output", "table")
	if l.Code != 0 || strings.Contains(l.Out, "SECRET") || !strings.Contains(l.Out, "afsk_ab") {
		t.Fatalf("%+v", l)
	}
}

// NFR security: a replayed create says the secret is gone and never prints one.
func TestCreateServiceAPIKeyReplayNoSecret(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Idempotent-Replayed", "true")
		clitest.Reply(w, 201, strings.Replace(createdKey, `,"secretKey":"afsk_SECRET_ONCE"`, "", 1))
	}, "admin", "create-api-key", "--name", "cve", "--role", "security_operator", "--expires-at", "2027-02-19T18:00:00Z", "--idempotency-key", "k")
	if r.Code != 0 || !strings.Contains(r.Err, "replay") {
		t.Fatalf("%+v", r)
	}
}

// CLI-1: add-security-operator sends the subject as given.
func TestAddSecurityOperator(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) {
		clitest.Reply(w, 201, `{"id":"rb_1","subjectType":"user","email":"ana@example.com","version":1}`)
	}, "admin", "add-security-operator", "--subject-type", "user", "--email", "ana@example.com", "--identity-provider-id", "idp_1")
	if r.Code != 0 || r.Req.URL.Path != "/v1/admin/security-operators" || !strings.Contains(string(r.ReqBody), `"subjectType":"user"`) || !strings.Contains(string(r.ReqBody), `"email":"ana@example.com"`) {
		t.Fatalf("%+v body=%s", r, r.ReqBody)
	}
}

func TestServiceAccessTeamAdminDenied(t *testing.T) {
	for _, args := range [][]string{{"admin", "list-security-operators"}, {"admin", "list-api-keys"}, {"admin", "remove-security-operator", "--binding-id", "rb_1"}} {
		r := clitest.Run(t, clitest.Denied("admin:settings"), args...)
		if r.Code != 254 || !strings.Contains(r.Err, "admin:settings") {
			t.Errorf("%v: %+v", args, r)
		}
	}
}
