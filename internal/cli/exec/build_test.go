package exec_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli/services"
)

// capture runs a command and returns what the server received.
func capture(t *testing.T, stdin string, args ...string) (result, string, http.Header) {
	t.Helper()
	var body string
	var hdr http.Header
	srv := func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		body, hdr = string(b), r.Header.Clone()
		jsonReply(w, 200, `{"id":"x"}`)
	}
	reg, err := services.NewRegistry()
	if err != nil {
		t.Fatal(err)
	}
	_ = reg
	res := runStdin(t, srv, stdin, args...)
	return res, body, hdr
}

func runStdin(t *testing.T, h http.HandlerFunc, stdin string, args ...string) result {
	t.Helper()
	reg, _ := services.NewRegistry()
	s := newServer(t, h)
	var out, eb bytes.Buffer
	env := cli.NewEnv(strings.NewReader(stdin), &out, &eb)
	env.Vars = map[string]string{"AFS_ENDPOINT_URL": s, "AFS_API_KEY": "afsk_unit_test_key"}
	env.Home = t.TempDir()
	code := cli.Run(context.Background(), args, env, reg)
	return result{code, out.String(), eb.String()}
}

func TestOptionKindsReachTheBody(t *testing.T) {
	dir := t.TempDir()
	adm := filepath.Join(dir, "admin.json")
	if err := os.WriteFile(adm, []byte(`{ "email": "a@b.c" }`), 0o600); err != nil {
		t.Fatal(err)
	}
	r, body, _ := capture(t, "", "teams", "create-team", "--slug", "pay", "--display-name", "Pay", "--kind", "persistent",
		"--max-budget-usd", "10.50", "--models", "m1", "m2", "--first-admin", "file://"+adm, "--may-add-external-systems", "--no-may-manage-aws-accounts")
	if r.code != 0 {
		t.Fatalf("%+v", r)
	}
	for _, want := range []string{`"models":["m1","m2"]`, `"firstAdmin":{"email":"a@b.c"}`, `"mayAddExternalSystems":true`, `"mayManageAwsAccounts":false`, `"maxBudgetUsd":"10.50"`} {
		if !strings.Contains(body, want) {
			t.Errorf("body lacks %s: %s", want, body)
		}
	}
	if r, _, _ := capture(t, "", "teams", "create-team", "--slug", "p", "--display-name", "P", "--kind", "persistent", "--max-budget-usd", "1", "--first-admin", "not json"); r.code != 252 {
		t.Errorf("bad JSON: %+v", r)
	}
	if r, _, _ := capture(t, "", "teams", "create-team", "--slug", "p", "--display-name", "P", "--kind", "persistent", "--max-budget-usd", "1", "--first-admin", "file:///nonexistent"); r.code != 252 {
		t.Errorf("missing file: %+v", r)
	}
}

func TestAnyOptionJSONOrString(t *testing.T) {
	_, body, _ := capture(t, "", "admin", "put-setting", "--key", "teams.ephemeral_max_days", "--value", "30", "--if-match", "1")
	if !strings.Contains(body, `"value":30`) {
		t.Errorf("%s", body)
	}
	_, body, _ = capture(t, "", "admin", "put-setting", "--key", "api.cors_origins", "--value", "plain text", "--if-match", "1")
	if !strings.Contains(body, `"value":"plain text"`) {
		t.Errorf("%s", body)
	}
}

func TestSecretValueFromStdinFileAndNeverInline(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "v.txt")
	if err := os.WriteFile(f, []byte("s3cret"), 0o600); err != nil {
		t.Fatal(err)
	}
	base := []string{"secrets", "put-secret-value", "--team-id", "team_1", "--secret-id", "sec_1", "--if-match", "1"}
	r, body, _ := capture(t, "", append(base, "--secret-value", "file://"+f)...)
	if r.code != 0 || !strings.Contains(body, `"value":"s3cret"`) && !strings.Contains(body, "s3cret") {
		t.Fatalf("%+v %s", r, body)
	}
	r, body, _ = capture(t, "from-stdin", append(base, "--secret-value", "-")...)
	if r.code != 0 || !strings.Contains(body, "from-stdin") {
		t.Fatalf("%+v %s", r, body)
	}
	r, _, _ = capture(t, "", append(base, "--secret-value", "plain-on-the-command-line")...)
	if r.code != 252 || strings.Contains(r.errOut+r.out, "plain-on-the-command-line") {
		t.Fatalf("%+v", r)
	}
	r, _, _ = capture(t, "", append(base[:len(base):len(base)], "--cli-input-json", `{"value":"inline-secret"}`)...)
	if r.code != 252 || strings.Contains(r.errOut, "inline-secret") {
		t.Fatalf("%+v", r)
	}
}

func TestCLIInputJSONBodyUnchangedAndOptionsWin(t *testing.T) {
	_, body, _ := capture(t, "", "teams", "create-team", "--cli-input-json",
		`{"slug":"a","displayName":"First","kind":"persistent","firstAdmin":{"email":"x@y.z"},"maxBudgetUsd":"5","extra":{"keep":[1,2]}}`, "--display-name", "Second")
	if !strings.Contains(body, `"displayName":"Second"`) || !strings.Contains(body, `"extra":{"keep":[1,2]}`) {
		t.Errorf("%s", body)
	}
	if r, _, _ := capture(t, "", "teams", "create-team", "--cli-input-json", "[1]"); r.code != 252 {
		t.Errorf("array: %+v", r)
	}
	if r, _, _ := capture(t, "", "teams", "create-team", "--cli-input-json", "file:///nonexistent.json"); r.code != 252 {
		t.Errorf("file: %+v", r)
	}
	if r, _, _ := capture(t, "", "teams", "describe-team", "--cli-input-json", `{"teamId":"team_1","bogus":1}`); r.code != 252 {
		t.Errorf("unknown member: %+v", r)
	}
}

func TestRangeAndEnumErrors252(t *testing.T) {
	for _, args := range [][]string{
		{"teams", "list-teams", "--page-size", "x"},
		{"teams", "create-team", "--models"},
	} {
		if r, _, _ := capture(t, "", args...); r.code != 252 {
			t.Errorf("%v: %+v", args, r)
		}
	}
}

func TestPathTeamIDFromProfile(t *testing.T) {
	var path string
	reg, _ := services.NewRegistry()
	s := newServer(t, func(w http.ResponseWriter, r *http.Request) { path = r.URL.Path; jsonReply(w, 200, `{}`) })
	home := t.TempDir()
	if err := os.MkdirAll(filepath.Join(home, ".afs"), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(home, ".afs", "config"), []byte("[default]\nendpoint_url = "+s+"\nteam_id = team_prof\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	var out, eb bytes.Buffer
	env := cli.NewEnv(strings.NewReader(""), &out, &eb)
	env.Vars = map[string]string{"AFS_API_KEY": "k"}
	env.Home = home
	if code := cli.Run(context.Background(), []string{"teams", "describe-team"}, env, reg); code != 0 || path != "/v1/teams/team_prof" {
		t.Fatalf("code=%d path=%s %s", code, path, eb.String())
	}
}

func TestSkeletonShape(t *testing.T) {
	r, _, _ := capture(t, "", "pools", "delete-stage", "--generate-cli-skeleton")
	for _, want := range []string{`"teamId"`, `"stageId"`, `"version"`} {
		if !strings.Contains(r.out, want) {
			t.Errorf("skeleton lacks %s: %s", want, r.out)
		}
	}
	r, _, _ = capture(t, "", "teams", "create-team", "--generate-cli-skeleton")
	for _, want := range []string{`"models": []`, `"firstAdmin": {}`, `"Idempotency-Key"`, `"mayAddExternalSystems": false`} {
		if !strings.Contains(r.out, want) {
			t.Errorf("skeleton lacks %s: %s", want, r.out)
		}
	}
}

func TestUsageErrorText(t *testing.T) {
	e := &cli.UsageError{Msg: "m"}
	if e.Error() != "m" {
		t.Fatal(e.Error())
	}
	p := &cli.Parsed{Values: map[string][]string{"a": {"1", "2"}}}
	if p.Value("a") != "1" || p.Value("zz") != "" {
		t.Fatal("Value")
	}
}
