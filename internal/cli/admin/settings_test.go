package admin_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/clitest"
)

func TestCommandsDoneSettings(t *testing.T) {
	clitest.CheckOperations(t, "settings", "GetSettings", "PutSetting")
}

func putSetting(t *testing.T, h http.HandlerFunc, args ...string) clitest.Result {
	t.Helper()
	return clitest.Run(t, h, append([]string{"admin", "put-setting"}, args...)...)
}

func ok(w http.ResponseWriter, _ *http.Request) {
	clitest.Reply(w, 200, `{"key":"k","value":1,"version":3}`)
}

// CLI-US-15.3: --value is JSON when it parses as JSON.
func TestPutSettingJSONValue(t *testing.T) {
	r := putSetting(t, ok, "--key", "api.cors_origins", "--value", `["https://console.example.com"]`, "--if-match", "2")
	if r.Code != 0 || string(r.ReqBody) != `{"value":["https://console.example.com"]}` || r.Req.URL.Path != "/v1/admin/settings/api.cors_origins" {
		t.Fatalf("%+v body=%s", r, r.ReqBody)
	}
}

// A plain word that is not JSON is sent as a string.
func TestPutSettingPlainStringValue(t *testing.T) {
	r := putSetting(t, ok, "--key", "secrets.path_prefix", "--value", "afs/prod", "--if-match", "2")
	if r.Code != 0 || string(r.ReqBody) != `{"value":"afs/prod"}` {
		t.Fatalf("%+v body=%s", r, r.ReqBody)
	}
	r = putSetting(t, ok, "--key", "api.idempotency_key_ttl_seconds", "--value", "5000", "--if-match", "2")
	if string(r.ReqBody) != `{"value":5000}` {
		t.Fatalf("number: %s", r.ReqBody)
	}
}

// CLI-US-21.1: an out-of-range value is the API's ValidationError naming the range (254).
func TestPutSettingValidationError(t *testing.T) {
	r := putSetting(t, clitest.Fail(422, "ValidationError", "value 99999 is outside 1 to 5000"), "--key", "api.cors_origins", "--value", "99999", "--if-match", "1")
	if r.Code != 254 || !strings.Contains(r.Err, "(ValidationError) when calling the PutSetting operation: value 99999 is outside 1 to 5000") {
		t.Fatalf("%+v", r)
	}
}

// CLI-US-21.2: a team-admin gets AccessDenied.
func TestPutSettingAccessDenied(t *testing.T) {
	r := putSetting(t, clitest.Denied("admin:settings"), "--key", "api.cors_origins", "--value", "1", "--if-match", "1")
	if r.Code != 254 || !strings.Contains(r.Err, "(AccessDenied) when calling the PutSetting operation") || !strings.Contains(r.Err, "admin:settings") {
		t.Fatalf("%+v", r)
	}
}

// CLI-US-15.3: get-settings --output json prints every setting as the API returned it.
func TestGetSettingsJSONPrintsEverySetting(t *testing.T) {
	body := `{"items":[{"key":"a","value":1,"version":1},{"key":"b","value":["x"],"version":2}]}`
	r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) { clitest.Reply(w, 200, body) }, "admin", "get-settings", "--output", "json")
	if r.Code != 0 || !strings.Contains(r.Out, `"key": "a"`) || !strings.Contains(r.Out, `"key": "b"`) {
		t.Fatalf("%+v", r)
	}
}
