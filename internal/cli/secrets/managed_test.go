package secrets_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/clitest"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/golden"
)

// CLI-US-29.1
func TestListSecretsManagedByFilter(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) {
		clitest.Reply(w, 200, `{"items":[{"id":"sec_1","name":"paperclip-better-auth-secret","scope":"team","managedBy":"afs","version":1}]}`)
	}, "secrets", "list-secrets", "--team-id", "team_1", "--managed-by", "afs", "--output", "json")
	if r.Code != 0 || r.Req.URL.Query().Get("managedBy") != "afs" || !strings.Contains(r.Out, `"managedBy": "afs"`) {
		t.Errorf("%d %v %q", r.Code, r.Req.URL.Query(), r.Out)
	}
}

func TestListSecretsManagedByEnum252(t *testing.T) {
	r := clitest.Run(t, nil, "secrets", "list-secrets", "--team-id", "team_1", "--managed-by", "robot")
	if r.Code != 252 || len(r.Requests) != 0 {
		t.Errorf("%d", r.Code)
	}
}

func TestSecretTableManagedByColumnGolden(t *testing.T) {
	body := `{"items":[{"id":"sec_1","name":"github-pat","scope":"team","managedBy":"user","version":1},{"id":"sec_2","name":"paperclip-secrets-master-key","scope":"team","managedBy":"afs","version":1}]}`
	for _, f := range []string{"table", "text"} {
		r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) { clitest.Reply(w, 200, body) }, "secrets", "list-secrets", "--team-id", "team_1", "--output", f)
		golden.Check(t, "managed-by."+f+".golden", r.Out)
	}
}

func managedConflict(w http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		clitest.Reply(w, 200, `{"id":"sec_1","version":3,"managedBy":"afs"}`)
		return
	}
	clitest.Reply(w, 409, `{"code":"Conflict","message":"managed by AFS","requestId":"req_c","details":{"reason":"managed_secret"}}`)
}

// CLI-US-29.2: the CLI does not read managedBy to decide; the API's 409 prints with the hint.
func TestManagedSecretMutations409(t *testing.T) {
	for _, args := range [][]string{
		{"secrets", "update-secret", "--team-id", "team_1", "--secret-id", "sec_1", "--description", "x", "--if-match", "3"},
		{"secrets", "delete-secret", "--team-id", "team_1", "--secret-id", "sec_1", "--if-match", "3"},
	} {
		r := clitest.Run(t, managedConflict, args...)
		if r.Code != 254 || !strings.Contains(r.Err, "(Conflict) when calling the") || !strings.Contains(r.Err, "This secret is managed by AFS or Paperclip and cannot be changed here.") {
			t.Errorf("%v: %d %q", args, r.Code, r.Err)
		}
		if len(r.Requests) != 1 {
			t.Errorf("%v: %d requests; managedBy must not be read first", args, len(r.Requests))
		}
	}
	r := clitest.RunWith(t, clitest.Options{Stdin: secret}, managedConflict, "secrets", "put-secret-value", "--team-id", "team_1", "--secret-id", "sec_1", "--secret-value", "-", "--if-match", "3")
	if r.Code != 254 || !strings.Contains(r.Err, "cannot be changed here") {
		t.Errorf("put-secret-value: %d %q", r.Code, r.Err)
	}
}

func TestManagedSecretRefusalNeverEchoesValue(t *testing.T) {
	r := clitest.RunWith(t, clitest.Options{Stdin: secret}, managedConflict, "secrets", "put-secret-value", "--team-id", "team_1", "--secret-id", "sec_1", "--secret-value", "-", "--debug", "--if-match", "3")
	if r.Code != 254 || strings.Contains(r.Out, "s3cr3t") || strings.Contains(r.Err, "s3cr3t") {
		t.Errorf("%d out %q err %q", r.Code, r.Out, r.Err)
	}
}

// CLI-US-29.3
func TestRotateManagedSecret(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) {
		clitest.Reply(w, 200, `{"id":"sec_1","name":"paperclip-better-auth-secret","scope":"team","managedBy":"afs","version":2}`)
	}, "secrets", "rotate-secret", "--team-id", "team_1", "--secret-id", "sec_1")
	if r.Code != 0 || len(r.ReqBody) != 0 || r.Req.Header.Get("Idempotency-Key") == "" || len(r.Requests) != 1 {
		t.Errorf("%d body %q", r.Code, r.ReqBody)
	}
	r = clitest.Run(t, managedConflict, "secrets", "rotate-secret", "--team-id", "team_1", "--secret-id", "sec_1")
	if r.Code != 254 || !strings.Contains(r.Err, "cannot be changed here") {
		t.Errorf("%d %q", r.Code, r.Err)
	}
}

func TestRotateMasterKey254(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) {
		clitest.Reply(w, 400, `{"code":"ValidationError","message":"not rotatable","requestId":"req_m","details":{"fields":[{"field":"secretId","reason":"master_key_not_rotatable"}]}}`)
	}, "secrets", "rotate-secret", "--team-id", "team_1", "--secret-id", "sec_1")
	if r.Code != 254 || !strings.Contains(r.Err, "The master key cannot be rotated: Paperclip's stored identities would become undecryptable.") {
		t.Errorf("%d %q", r.Code, r.Err)
	}
}

// CLI-US-29.4
func TestCreateSecretReservedName254(t *testing.T) {
	r := clitest.RunWith(t, clitest.Options{Stdin: secret}, func(w http.ResponseWriter, _ *http.Request) {
		clitest.Reply(w, 400, `{"code":"ValidationError","message":"reserved","requestId":"req_n","details":{"fields":[{"field":"name","reason":"reserved_name"}]}}`)
	}, "secrets", "create-secret", "--team-id", "team_1", "--name", "paperclip-x", "--scope", "team", "--secret-value", "-")
	if r.Code != 254 || !strings.Contains(r.Err, "That secret name is reserved for Paperclip") {
		t.Errorf("%d %q", r.Code, r.Err)
	}
}

func TestCreateAgentKeySecretManualRegistration(t *testing.T) {
	r := clitest.RunWith(t, clitest.Options{Stdin: secret}, func(w http.ResponseWriter, _ *http.Request) {
		clitest.Reply(w, 201, `{"id":"sec_9","name":"paperclip-agent-api-key","scope":"agent","managedBy":"user","version":1}`)
	}, "secrets", "create-secret", "--team-id", "team_1", "--name", "paperclip-agent-api-key", "--scope", "agent", "--agent-id", "agt_1", "--secret-value", "-")
	if r.Code != 0 || !strings.Contains(string(r.ReqBody), `"name":"paperclip-agent-api-key"`) || strings.Contains(r.Out, "s3cr3t") {
		t.Errorf("%d %s", r.Code, r.ReqBody)
	}
}

// CLI-US-29.5
func TestCreateSecretAccessDeniedCollaborator254(t *testing.T) {
	r := clitest.RunWith(t, clitest.Options{Stdin: secret}, clitest.Denied("secrets:write"),
		"secrets", "create-secret", "--team-id", "team_1", "--name", "x", "--scope", "team", "--secret-value", "-")
	if r.Code != 254 || !strings.Contains(r.Err, "Required permission: secrets:write") {
		t.Errorf("%d %q", r.Code, r.Err)
	}
}
