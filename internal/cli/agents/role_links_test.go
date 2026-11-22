package agents_test

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/clitest"
)

const roleLink = `{"id":"arl_01JD2Q3R4S5T6V7W8X9Y0Z0001","accountId":"444455556666","externalId":"afs-4N8P2Q","profileName":"deploy","roleArn":"arn:aws:iam::444455556666:role/payments-deployer","status":"pending_trust","lastError":null,"lastVerifiedAt":null,"maxSessionSeconds":3600,"version":1,"trustPolicy":{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":["sts:AssumeRole"],"Principal":{"AWS":"arn:aws:iam::111122223333:role/afs-payments-coder"},"Condition":{"StringEquals":{"sts:ExternalId":"afs-4N8P2Q"}}}]}}`

func TestCommandsDoneRoleLinks(t *testing.T) {
	clitest.CheckOperations(t, "role-links", "ListAgentRoleLinks", "CreateAgentRoleLink", "GetAgentRoleLink", "VerifyAgentRoleLink", "DeleteAgentRoleLink")
}

func TestCommandsDoneIdentities(t *testing.T) {
	clitest.CheckOperations(t, "identities", "ListAgentIdentities", "CreateAgentIdentity", "DeleteAgentIdentity")
}

// CLI-US-13.1: create-role-link prints pending_trust and the external id.
func TestCreateRoleLinkPendingTrust(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) { clitest.Reply(w, 201, roleLink) },
		"agents", "create-role-link", "--team-id", "team_1", "--agent-id", "agt_1",
		"--role-arn", "arn:aws:iam::444455556666:role/payments-deployer", "--profile-name", "deploy", "--output", "table")
	if r.Code != 0 || !strings.Contains(r.Out, "pending_trust") || !strings.Contains(r.Out, "afs-4N8P2Q") {
		t.Fatalf("%+v", r)
	}
	var body map[string]any
	_ = json.Unmarshal(r.ReqBody, &body)
	if r.Req.Method != "POST" || r.Req.URL.Path != "/v1/teams/team_1/agents/agt_1/role-links" || body["roleArn"] != "arn:aws:iam::444455556666:role/payments-deployer" || body["profileName"] != "deploy" {
		t.Errorf("%s %s %s", r.Req.Method, r.Req.URL.Path, r.ReqBody)
	}
	if r.Req.Header.Get("Idempotency-Key") == "" {
		t.Error("no Idempotency-Key")
	}
}

// CLI-US-13.2: --query trustPolicy prints the document alone, ready to paste into IAM.
func TestDescribeRoleLinkTrustPolicyQuery(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) { clitest.Reply(w, 200, roleLink) },
		"agents", "describe-role-link", "--team-id", "team_1", "--agent-id", "agt_1", "--link-id", "arl_1", "--query", "trustPolicy", "--output", "json")
	if r.Code != 0 {
		t.Fatalf("%+v", r)
	}
	var doc map[string]any
	if err := json.Unmarshal([]byte(r.Out), &doc); err != nil {
		t.Fatalf("not JSON: %v\n%s", err, r.Out)
	}
	if doc["Version"] != "2012-10-17" || doc["id"] != nil || len(doc) != 2 {
		t.Errorf("not the trust policy alone: %s", r.Out)
	}
}

// CLI-US-13.3: verify-role-link shows verified or failed with the reason.
func TestVerifyRoleLinkFailedReason(t *testing.T) {
	failed := strings.Replace(strings.Replace(roleLink, `"status":"pending_trust"`, `"status":"failed"`, 1), `"lastError":null`, `"lastError":"AccessDenied: the role does not trust the agent role"`, 1)
	r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) { clitest.Reply(w, 202, failed) },
		"agents", "verify-role-link", "--team-id", "team_1", "--agent-id", "agt_1", "--link-id", "arl_1", "--output", "table")
	if r.Code != 0 || !strings.Contains(r.Out, "failed") || !strings.Contains(r.Out, "AccessDenied: the role does not trust the agent role") {
		t.Fatalf("%+v", r)
	}
	if r.Req.URL.Path != "/v1/teams/team_1/agents/agt_1/role-links/arl_1:verify" {
		t.Errorf("path %s", r.Req.URL.Path)
	}
	ok := strings.Replace(roleLink, `"status":"pending_trust"`, `"status":"verified"`, 1)
	r = clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) { clitest.Reply(w, 202, ok) },
		"agents", "verify-role-link", "--team-id", "team_1", "--agent-id", "agt_1", "--link-id", "arl_1", "--output", "table")
	if r.Code != 0 || !strings.Contains(r.Out, "verified") {
		t.Fatalf("%+v", r)
	}
}

// CLI-US-13.3: an account outside the team's allowed accounts is the API's ValidationError, 254.
func TestRoleLinkAccountNotAllowed254(t *testing.T) {
	r := clitest.Run(t, clitest.Fail(400, "ValidationError", "account 999988887777 is not in the team's allowed accounts"),
		"agents", "create-role-link", "--team-id", "team_1", "--agent-id", "agt_1", "--role-arn", "arn:aws:iam::999988887777:role/x", "--profile-name", "deploy")
	if r.Code != 254 || !strings.Contains(r.Err, "(ValidationError) when calling the CreateAgentRoleLink operation") || !strings.Contains(r.Err, "allowed accounts") {
		t.Fatalf("%+v", r)
	}
}
