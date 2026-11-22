package exec_test

import (
	"net/http"
	"strings"
	"testing"
	"unicode/utf8"
)

func assertNoControl(t *testing.T, what, s string) {
	t.Helper()
	for _, r := range s {
		if (r < 0x20 && r != '\n' && r != '\t') || (r >= 0x7f && r <= 0x9f) {
			t.Fatalf("%s: control character %U in %q", what, r, s)
		}
	}
	if !utf8.ValidString(s) {
		t.Fatalf("%s: invalid UTF-8", what)
	}
}

// FR-003: an API error message, its detail lines and listed names reach the terminal without escapes.
func TestServiceTextIsSanitizedOnTheTerminal(t *testing.T) {
	r := run(t, func(w http.ResponseWriter, _ *http.Request) {
		jsonReply(w, 403, `{"code":"AccessDenied","message":"no \u001b]0;pwned\u0007 \r\u009b2J","requestId":"req\u001b[2J","details":{"requiredPermission":"p\u001b[31m","dependents":["a\u001b[2J"],"missingTags":["t\u0007"]}}`)
	}, "teams", "list-teams")
	if r.code != 254 {
		t.Fatalf("%+v", r)
	}
	assertNoControl(t, "stderr", r.errOut)
	if !strings.Contains(r.errOut, `\x1b]0;pwned\x07`) {
		t.Errorf("escape not visible: %q", r.errOut)
	}

	for _, format := range []string{"table", "text"} {
		r = run(t, func(w http.ResponseWriter, _ *http.Request) {
			jsonReply(w, 200, `{"items":[{"id":"team_1","slug":"\u001b[2J\u001b]52;c;eA==\u0007\r","state":"active","version":1}]}`)
		}, "teams", "list-teams", "--output", format)
		if r.code != 0 {
			t.Fatalf("%+v", r)
		}
		assertNoControl(t, format, r.out)
	}
}
