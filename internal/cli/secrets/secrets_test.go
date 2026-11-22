package secrets_test

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/clitest"
)

func TestCommandsDone(t *testing.T) {
	clitest.CheckOperations(t, "secrets", "ListSecrets", "GetSecret", "CreateSecret", "PutSecretValue", "UpdateSecret", "RotateSecret", "DeleteSecret")
}

const secret = "s3cr3t-Value\n"

func created(w http.ResponseWriter, _ *http.Request) {
	clitest.Reply(w, 201, `{"id":"sec_1","name":"github-pat","scope":"team","version":1}`)
}

func createArgs(extra ...string) []string {
	return append([]string{"secrets", "create-secret", "--team-id", "team_1", "--name", "github-pat", "--scope", "team"}, extra...)
}

func assertValueSentButNeverPrinted(t *testing.T, r clitest.Result) {
	t.Helper()
	if r.Code != 0 {
		t.Fatalf("%+v", r)
	}
	if !strings.Contains(string(r.ReqBody), `"value":"s3cr3t-Value\n"`) {
		t.Errorf("value not sent verbatim: %s", r.ReqBody)
	}
	if strings.Contains(r.Out, "s3cr3t") || strings.Contains(r.Err, "s3cr3t") {
		t.Errorf("value reached the output: %q %q", r.Out, r.Err)
	}
}

// CLI-US-3.1: the value comes from a file, byte for byte.
func TestCreateSecretFromFile(t *testing.T) {
	f := filepath.Join(t.TempDir(), "token.txt")
	if err := os.WriteFile(f, []byte(secret), 0o600); err != nil {
		t.Fatal(err)
	}
	assertValueSentButNeverPrinted(t, clitest.Run(t, created, createArgs("--secret-value", "file://"+f)...))
	assertValueSentButNeverPrinted(t, clitest.Run(t, created, createArgs("--secret-value", "fileb://"+f)...))
}

// CLI-US-3.2: the value comes from standard input.
func TestCreateSecretStdin(t *testing.T) {
	assertValueSentButNeverPrinted(t, clitest.RunWith(t, clitest.Options{Stdin: secret}, created, createArgs("--secret-value", "-")...))
}

// CLI-US-3.3 (prompt): with a terminal and no value, afs asks without echo.
func TestCreateSecretPrompt(t *testing.T) {
	pw := secret
	var prompts []string
	r := clitest.RunWith(t, clitest.Options{Password: &pw, Prompts: &prompts}, created, createArgs()...)
	assertValueSentButNeverPrinted(t, r)
	if len(prompts) != 1 || strings.Contains(prompts[0], "s3cr3t") {
		t.Errorf("prompts %v", prompts)
	}
}

// CLI-US-3.3: no terminal and no value is a usage error naming the ways to give it.
func TestNoTTYNoValue252(t *testing.T) {
	r := clitest.Run(t, created, createArgs()...)
	if r.Code != 252 || len(r.Requests) != 0 || !strings.Contains(r.Err, "file://") {
		t.Fatalf("%+v", r)
	}
}

// CLI-US-3.3: a plain argument is refused before any call and never echoed.
func TestSecretValuePlaintext252(t *testing.T) {
	r := clitest.Run(t, created, createArgs("--secret-value", "hunter2")...)
	if r.Code != 252 || len(r.Requests) != 0 || strings.Contains(r.Err, "hunter2") || strings.Contains(r.Out, "hunter2") {
		t.Fatalf("%+v", r)
	}
}

// CLI-9: a value that is not UTF-8, empty or over 64 KiB is refused (252).
func TestNonUTF8OrTooLong252(t *testing.T) {
	dir := t.TempDir()
	for name, content := range map[string][]byte{"bin": {0xff, 0xfe, 0x00}, "empty": {}, "big": []byte(strings.Repeat("x", 65537))} {
		f := filepath.Join(dir, name)
		if err := os.WriteFile(f, content, 0o600); err != nil {
			t.Fatal(err)
		}
		if r := clitest.Run(t, created, createArgs("--secret-value", "file://"+f)...); r.Code != 252 || len(r.Requests) != 0 {
			t.Errorf("%s: %+v", name, r)
		}
	}
}

// CLI-US-3.3: --cli-input-json with a value inline is refused (252) before any call.
func TestInlineInputJSONWithValue252(t *testing.T) {
	r := clitest.Run(t, created, "secrets", "create-secret", "--team-id", "team_1", "--cli-input-json",
		`{"name":"n","scope":"team","value":"hunter2"}`)
	if r.Code != 252 || len(r.Requests) != 0 || strings.Contains(r.Err, "hunter2") {
		t.Fatalf("%+v", r)
	}
}

// CLI-US-3.3: the same document in a file is accepted.
func TestInputJSONFileAccepted(t *testing.T) {
	f := filepath.Join(t.TempDir(), "secret.json")
	if err := os.WriteFile(f, []byte(`{"name":"n","scope":"team","value":"hunter2"}`), 0o600); err != nil {
		t.Fatal(err)
	}
	r := clitest.Run(t, created, "secrets", "create-secret", "--team-id", "team_1", "--cli-input-json", "file://"+f)
	if r.Code != 0 || !strings.Contains(string(r.ReqBody), `"value":"hunter2"`) || strings.Contains(r.Out+r.Err, "hunter2") {
		t.Fatalf("%+v", r)
	}
}

// CLI-US-3.4, NFR security: --debug shows the value and the credential redacted.
func TestDebugRedactsValue(t *testing.T) {
	r := clitest.RunWith(t, clitest.Options{Stdin: secret}, created, createArgs("--secret-value", "-", "--debug")...)
	if r.Code != 0 || strings.Contains(r.Err, "s3cr3t") || strings.Contains(r.Err, "afsk_unit_test_key") {
		t.Fatalf("leak: %+v", r)
	}
	if !strings.Contains(r.Err, `"value": "[redacted]"`) && !strings.Contains(r.Err, `"value":"[redacted]"`) {
		t.Errorf("value not shown as redacted:\n%s", r.Err)
	}
	if !strings.Contains(strings.ToLower(r.Err), "authorization: [redacted]") && !strings.Contains(r.Err, "Authorization") {
		t.Errorf("no redacted Authorization line:\n%s", r.Err)
	}
}

// CLI-US-3.1: put-secret-value answers 204 and prints nothing.
func TestPutSecretValue204PrintsNothing(t *testing.T) {
	r := clitest.RunWith(t, clitest.Options{Stdin: secret}, func(w http.ResponseWriter, req *http.Request) {
		if req.Method == "GET" {
			clitest.Reply(w, 200, `{"id":"sec_1","version":2}`)
			return
		}
		w.WriteHeader(204)
	}, "secrets", "put-secret-value", "--team-id", "team_1", "--secret-id", "sec_1", "--secret-value", "-")
	if r.Code != 0 || r.Out != "" || r.Req.Method != "PUT" || r.Req.URL.Path != "/v1/teams/team_1/secrets/sec_1/value" {
		t.Fatalf("%+v", r)
	}
	assertSent := strings.Contains(string(r.ReqBody), `"value":"s3cr3t-Value\n"`)
	if !assertSent {
		t.Errorf("body %s", r.ReqBody)
	}
}

// CLI-1: a delete blocked by dependents is the API's Conflict, with the dependents listed (254).
func TestDeleteSecret409PrintsDependents(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, req *http.Request) {
		if req.Method == "GET" {
			clitest.Reply(w, 200, `{"id":"sec_1","version":2}`)
			return
		}
		clitest.Reply(w, 409, `{"code":"Conflict","message":"the secret is bound","requestId":"r","details":{"dependents":["agt_1","agt_2"]}}`)
	}, "secrets", "delete-secret", "--team-id", "team_1", "--secret-id", "sec_1")
	if r.Code != 254 || !strings.Contains(r.Err, "(Conflict) when calling the DeleteSecret operation") || !strings.Contains(r.Err, "Dependents: agt_1, agt_2") {
		t.Fatalf("%+v", r)
	}
}
