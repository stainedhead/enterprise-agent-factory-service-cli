package collaboration_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/clitest"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/golden"
)

const rotated = `{"rotated":[{"name":"paperclip-better-auth-secret","secretId":"sec_01JD2Q3R4S5T6V7W8X9Y0Z0001","versionId":"3b2c9f1e-0000-4000-8000-000000000001"},{"name":"paperclip-tool-action-signing-secret","secretId":"sec_01JD2Q3R4S5T6V7W8X9Y0Z0002","versionId":"3b2c9f1e-0000-4000-8000-000000000002"}],"signsOutEveryone":true,"agentRolloutPending":true}`

// CLI-US-27.1
func TestRotateCollaborationKeysBody(t *testing.T) {
	r := clitest.Run(t, reply(rotated), "collaboration", "rotate-collaboration-keys", "--team-id", team)
	if r.Code != 0 || r.Req.Method != "POST" || r.Req.URL.Path != "/v1/teams/"+team+"/collaboration:rotate-keys" || len(r.ReqBody) != 0 ||
		r.Req.Header.Get("Idempotency-Key") == "" || len(r.Requests) != 1 || r.Req.Header.Get("If-Match") != "" {
		t.Fatalf("%+v body %q", r, r.ReqBody)
	}
	r = clitest.Run(t, reply(rotated), "collaboration", "rotate-collaboration-keys", "--team-id", team, "--names", "paperclip-better-auth-secret", "paperclip-tool-action-signing-secret")
	if string(r.ReqBody) != `{"names":["paperclip-better-auth-secret","paperclip-tool-action-signing-secret"]}` {
		t.Errorf("body %s", r.ReqBody)
	}
}

func TestRotateCollaborationKeysTable(t *testing.T) {
	for _, f := range []string{"table", "text"} {
		r := clitest.Run(t, reply(rotated), "collaboration", "rotate-collaboration-keys", "--team-id", team, "--output", f)
		golden.Check(t, "rotate-keys."+f+".golden", r.Out)
	}
}

// CLI-US-27.2
func TestRotateCollaborationKeysNotices(t *testing.T) {
	r := clitest.Run(t, reply(rotated), "collaboration", "rotate-collaboration-keys", "--team-id", team)
	if r.Code != 0 || !strings.Contains(r.Err, "Everyone is signed out of Paperclip.\n") || !strings.Contains(r.Err, "Agents will roll to pick up the new keys.\n") {
		t.Errorf("%d %q", r.Code, r.Err)
	}
	none := strings.Replace(strings.Replace(rotated, `"signsOutEveryone":true`, `"signsOutEveryone":false`, 1), `"agentRolloutPending":true`, `"agentRolloutPending":false`, 1)
	r = clitest.Run(t, reply(none), "collaboration", "rotate-collaboration-keys", "--team-id", team)
	if r.Err != "" {
		t.Errorf("no notice expected: %q", r.Err)
	}
	var prompted bool
	pw := ""
	r = clitest.RunWith(t, clitest.Options{Password: &pw, Prompts: &[]string{}, TTY: true}, func(w http.ResponseWriter, _ *http.Request) { prompted = true; clitest.Reply(w, 200, rotated) },
		"collaboration", "rotate-collaboration-keys", "--team-id", team)
	if !prompted || r.Code != 0 {
		t.Error("the command must run without any confirmation")
	}
}

func TestRotateCollaborationKeysNoValue(t *testing.T) {
	r := clitest.Run(t, reply(rotated), "collaboration", "rotate-collaboration-keys", "--team-id", team, "--debug")
	for _, bad := range []string{"value", "plaintext"} {
		if strings.Contains(strings.ToLower(r.Out), `"`+bad+`"`) {
			t.Errorf("output has a %s member: %s", bad, r.Out)
		}
	}
}

// CLI-US-27.3
func TestRotateCollaborationKeysMasterKey254(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) {
		clitest.Reply(w, 400, `{"code":"ValidationError","message":"not rotatable","requestId":"req_k","details":{"fields":[{"field":"names","reason":"master_key_not_rotatable"}]}}`)
	}, "collaboration", "rotate-collaboration-keys", "--team-id", team, "--names", "paperclip-secrets-master-key")
	if r.Code != 254 || string(r.ReqBody) != `{"names":["paperclip-secrets-master-key"]}` ||
		!strings.Contains(r.Err, "The master key cannot be rotated: Paperclip's stored identities would become undecryptable.") {
		t.Errorf("%d body %s %q", r.Code, r.ReqBody, r.Err)
	}
}

// CLI-US-27.4
func TestRotateCollaborationKeysAccessDenied254(t *testing.T) {
	r := clitest.Run(t, clitest.Denied("collaboration:manage"), "collaboration", "rotate-collaboration-keys", "--team-id", team)
	if r.Code != 254 || !strings.Contains(r.Err, "Required permission: collaboration:manage") {
		t.Errorf("%d %q", r.Code, r.Err)
	}
}
