package admin_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/clitest"
)

func TestCommandsDoneIdentityProviders(t *testing.T) {
	clitest.CheckOperations(t, "identity-providers", "ListIdentityProviders", "GetIdentityProvider", "CreateIdentityProvider", "UpdateIdentityProvider", "DeleteIdentityProvider")
}

// CLI-1: the path parameter is --provider-id.
func TestProviderIDRename(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) {
		clitest.Reply(w, 200, `{"id":"idp_1","name":"Corp","clientSecretArn":"arn:aws:secretsmanager:eu-west-1:111122223333:secret:x"}`)
	}, "admin", "describe-identity-provider", "--provider-id", "idp_1")
	if r.Code != 0 || r.Req.URL.Path != "/v1/admin/identity-providers/idp_1" || !strings.Contains(r.Out, "clientSecretArn") {
		t.Fatalf("%+v", r)
	}
}

// CLI-1: deleting a provider that teams use is a 409 that is printed (254).
func TestDeleteIdentityProviderInUse409(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, req *http.Request) {
		if req.Method == "GET" {
			clitest.Reply(w, 200, `{"id":"idp_1","version":3}`)
			return
		}
		clitest.Reply(w, 409, `{"code":"Conflict","message":"identity provider is used by 2 teams","requestId":"r"}`)
	}, "admin", "delete-identity-provider", "--provider-id", "idp_1")
	if r.Code != 254 || !strings.Contains(r.Err, "(Conflict) when calling the DeleteIdentityProvider operation: identity provider is used by 2 teams") {
		t.Fatalf("%+v", r)
	}
}

// NFR security: the output carries the secret's ARN and never a value.
func TestCreateIdentityProviderNoSecretValue(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) {
		clitest.Reply(w, 201, `{"id":"idp_1","name":"Corp","clientSecretArn":"arn:aws:secretsmanager:eu-west-1:111122223333:secret:x"}`)
	}, "admin", "create-identity-provider", "--name", "Corp", "--kind", "okta", "--issuer-url", "https://idp.example.com", "--console-client-id", "c", "--cli-client-id", "l", "--client-secret-arn", "arn:aws:secretsmanager:eu-west-1:111122223333:secret:x")
	if r.Code != 0 || strings.Contains(strings.ToLower(r.Out), "clientsecret\"") {
		t.Fatalf("%+v", r)
	}
	if !strings.Contains(string(r.ReqBody), `"clientSecretArn"`) {
		t.Errorf("body %s", r.ReqBody)
	}
}
