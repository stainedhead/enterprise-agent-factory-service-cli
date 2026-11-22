package admin_test

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/clitest"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/golden"
)

const endpoint = `{"teamId":"team_1","slug":"payments","teamState":"active","hostname":"payments.collab.example.com","agentId":"agt_1","webUiEnabled":true,"restApiEnabled":false,"effective":{"webUi":true,"restApi":false},"mayExposeCollaboration":true,"phase":"Ready","lastAccessAt":null}`

func TestCommandsDoneCollaborationEndpoints(t *testing.T) {
	clitest.CheckOperations(t, "collaboration-endpoints", "ListCollaborationEndpoints")
}

// CLI-US-30.1
func TestListCollaborationEndpointsFilters(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) { clitest.Reply(w, 200, `{"items":[`+endpoint+`]}`) },
		"admin", "list-collaboration-endpoints", "--team-id", "team_1", "--state", "Ready")
	q := r.Req.URL.Query()
	if r.Code != 0 || r.Req.URL.Path != "/v1/admin/collaboration-endpoints" || q.Get("teamId") != "team_1" || q.Get("state") != "Ready" {
		t.Errorf("%d %s %v", r.Code, r.Req.URL.Path, q)
	}
}

func TestListCollaborationEndpointsStateEnum252(t *testing.T) {
	r := clitest.Run(t, nil, "admin", "list-collaboration-endpoints", "--state", "Bad")
	if r.Code != 252 || len(r.Requests) != 0 {
		t.Errorf("%d", r.Code)
	}
}

func TestListCollaborationEndpointsTeamIdNotDefaulted(t *testing.T) {
	home := t.TempDir()
	writeProfile(t, home, "team_01J9ZQ3ABCDEFGHJKMNPQRSTVW")
	r := clitest.RunWith(t, clitest.Options{Home: home}, func(w http.ResponseWriter, _ *http.Request) { clitest.Reply(w, 200, `{"items":[]}`) },
		"admin", "list-collaboration-endpoints")
	if r.Code != 0 || r.Req.URL.Query().Has("teamId") {
		t.Errorf("the profile team_id must not fill the query filter: %d %v", r.Code, r.Req.URL.Query())
	}
}

// CLI-US-30.2
func TestCollaborationEndpointsTableGolden(t *testing.T) {
	for _, f := range []string{"table", "text"} {
		r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) { clitest.Reply(w, 200, `{"items":[`+endpoint+`]}`) },
			"admin", "list-collaboration-endpoints", "--output", f)
		golden.Check(t, "collaboration-endpoints."+f+".golden", r.Out)
	}
}

func TestListCollaborationEndpointsCSVThreePages(t *testing.T) {
	pages := []string{"team,hostname\nt1,h1\n", "team,hostname\nt2,h2\n", "team,hostname\nt3,h3\n"}
	n := 0
	r := clitest.Run(t, func(w http.ResponseWriter, req *http.Request) {
		if req.Header.Get("Accept") != "text/csv" {
			t.Errorf("Accept = %q", req.Header.Get("Accept"))
		}
		w.Header().Set("Content-Type", "text/csv")
		if n < 2 {
			w.Header().Set("X-Afs-Next-Token", fmt.Sprintf("n%d", n+1))
		}
		w.Header().Set("X-Request-Id", "req_c")
		_, _ = w.Write([]byte(pages[n]))
		n++
	}, "admin", "list-collaboration-endpoints", "--output", "csv")
	if r.Code != 0 || len(r.Requests) != 3 || r.Out != "team,hostname\nt1,h1\nt2,h2\nt3,h3\n" {
		t.Errorf("%d %d requests %q %s", r.Code, len(r.Requests), r.Out, r.Err)
	}
}

// CLI-US-30.3
func TestListCollaborationEndpointsAccessDenied254(t *testing.T) {
	r := clitest.Run(t, clitest.Denied("admin:settings"), "admin", "list-collaboration-endpoints")
	if r.Code != 254 || !strings.Contains(r.Err, "Required permission: admin:settings") {
		t.Errorf("%d %q", r.Code, r.Err)
	}
}

const arn = "arn:aws:secretsmanager:eu-west-1:111122223333:secret:afs/collab-AbCdEf"

func idpUpdate(extra ...string) []string {
	return append([]string{"admin", "update-identity-provider", "--provider-id", "idp_1", "--if-match", "2"}, extra...)
}

// CLI-US-31.2
func TestIdentityProviderCollabClientBody(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) {
		clitest.Reply(w, 200, `{"id":"idp_1","name":"Corp","collabClientId":"paperclip-web","collabClientSecretArn":"`+arn+`","version":3}`)
	}, idpUpdate("--collab-client-id", "paperclip-web", "--collab-client-secret-arn", arn)...)
	if r.Code != 0 || string(r.ReqBody) != `{"collabClientId":"paperclip-web","collabClientSecretArn":"`+arn+`"}` {
		t.Errorf("%d %s %s", r.Code, r.ReqBody, r.Err)
	}
	d := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) {
		clitest.Reply(w, 200, `{"id":"idp_1","collabClientId":"paperclip-web","collabClientSecretArn":"`+arn+`"}`)
	}, "admin", "describe-identity-provider", "--provider-id", "idp_1", "--output", "json", "--debug")
	if !strings.Contains(d.Out, arn) || !strings.Contains(d.Err, arn) {
		t.Errorf("the ARN is a reference and is printed, even by --debug: out %q err %q", d.Out, d.Err)
	}
}

func TestIdentityProviderCollabSecretArnShape252(t *testing.T) {
	for _, bad := range []string{"hunter2-the-secret-value", "arn:aws:s3:::bucket", "arn:aws:secretsmanager:eu-west-1:1234:secret:x", ""} {
		r := clitest.Run(t, nil, idpUpdate("--collab-client-secret-arn", bad)...)
		if r.Code != 252 || len(r.Requests) != 0 || (bad != "" && strings.Contains(r.Err, bad)) {
			t.Errorf("%q: %d requests %d err %q", bad, r.Code, len(r.Requests), r.Err)
		}
	}
	in := writeInput(t, `{"collabClientSecretArn":"not-an-arn","version":2}`)
	if r := clitest.Run(t, nil, "admin", "update-identity-provider", "--provider-id", "idp_1", "--cli-input-json", in); r.Code != 252 {
		t.Errorf("input json: %d", r.Code)
	}
	c := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) { clitest.Reply(w, 201, `{"id":"idp_1"}`) },
		"admin", "create-identity-provider", "--name", "Corp", "--kind", "okta", "--issuer-url", "https://idp.example.com", "--console-client-id", "c", "--cli-client-id", "l",
		"--client-secret-arn", arn, "--collab-client-id", "p", "--collab-client-secret-arn", "nope")
	if c.Code != 252 {
		t.Errorf("create: %d", c.Code)
	}
}

// CLI-US-31.3
func TestPutSettingCollabKeysPassThrough(t *testing.T) {
	cases := map[string]string{
		"collab.base_domain":          `"collab.example.com"`,
		"collab.public_endpoints":     `true`,
		"collab.session_ttl_minutes":  `480`,
		"collab.max_body_mb":          `25`,
		"collab.content_scan":         `"enforce"`,
		"collab.allowed_source_cidrs": `["10.0.0.0/8"]`,
	}
	for key, val := range cases {
		r := putSetting(t, ok, "--key", key, "--value", val, "--if-match", "1")
		if r.Code != 0 || string(r.ReqBody) != `{"value":`+val+`}` || r.Req.URL.Path != "/v1/admin/settings/"+key {
			t.Errorf("%s: %d %s", key, r.Code, r.ReqBody)
		}
	}
}

func TestPutSettingCollabOutOfRange254(t *testing.T) {
	r := putSetting(t, clitest.Fail(422, "ValidationError", "collab.max_body_mb must be between 1 and 100"), "--key", "collab.max_body_mb", "--value", "500", "--if-match", "1")
	if r.Code != 254 || !strings.Contains(r.Err, "collab.max_body_mb must be between 1 and 100") {
		t.Errorf("%d %q", r.Code, r.Err)
	}
}

// CLI-US-31.4
func TestIdentityProviderCollabClientAccessDenied254(t *testing.T) {
	r := clitest.Run(t, clitest.Denied("admin:identity-providers"), idpUpdate("--collab-client-id", "paperclip-web")...)
	if r.Code != 254 || !strings.Contains(r.Err, "(AccessDenied) when calling the UpdateIdentityProvider operation") {
		t.Errorf("%d %q", r.Code, r.Err)
	}
}

func writeInput(t *testing.T, content string) string {
	t.Helper()
	f := t.TempDir() + "/in.json"
	if err := os.WriteFile(f, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	return "file://" + f
}

func writeProfile(t *testing.T, home, teamID string) {
	t.Helper()
	if err := os.MkdirAll(home+"/.afs", 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(home+"/.afs/config", []byte("[default]\nteam_id = "+teamID+"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
}
