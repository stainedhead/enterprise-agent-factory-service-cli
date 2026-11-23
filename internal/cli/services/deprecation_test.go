package services_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/clitest"
)

// ORC-27, CLI-US-27: a response that carries the Deprecation header prints one line with the Sunset date, for a
// request the CLI had no static knowledge of; an operation the CLI already says is deprecated prints no second line.
func TestServiceDeprecationHeaderLine(t *testing.T) {
	h := func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Deprecation", "@1795392000")
		w.Header().Set("Sunset", "Tue, 23 Nov 2027 00:00:00 GMT")
		clitest.Reply(w, 200, `{"items":[]}`)
	}
	r := clitest.Run(t, h, "agents", "list-agents", "--team-id", "team_1")
	want := "Warning: the service reports that this request uses a deprecated feature; it is planned for removal on Tue, 23 Nov 2027 00:00:00 GMT.\n"
	if r.Code != 0 || r.Err != want {
		t.Errorf("exit %d err %q", r.Code, r.Err)
	}
	// no Sunset date: the line has none
	h2 := func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Deprecation", "true")
		clitest.Reply(w, 200, `{"items":[]}`)
	}
	if r := clitest.Run(t, h2, "agents", "list-agents", "--team-id", "team_1"); r.Err != "Warning: the service reports that this request uses a deprecated feature.\n" {
		t.Errorf("err %q", r.Err)
	}
	// a static deprecation covers the header: one line only
	d := clitest.Run(t, h, "admin", "list-collaboration-endpoints")
	if strings.Count(d.Err, "Warning:") != 1 || !strings.Contains(d.Err, "afs admin list-collaboration-endpoints is deprecated") {
		t.Errorf("err %q", d.Err)
	}
	// terminal escapes in a service-supplied date are not passed through
	h3 := func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Deprecation", "true")
		w.Header().Set("Sunset", "soon\x1b[31m")
		clitest.Reply(w, 200, `{"items":[]}`)
	}
	if r := clitest.Run(t, h3, "agents", "list-agents", "--team-id", "team_1"); strings.Contains(r.Err, "\x1b") {
		t.Errorf("escape in %q", r.Err)
	}
	// no header, no line; the skeleton prints no warning either
	if r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) { clitest.Reply(w, 200, `{"items":[]}`) }, "agents", "list-agents", "--team-id", "team_1"); r.Err != "" {
		t.Errorf("err %q", r.Err)
	}
	if s := clitest.Run(t, nil, "teams", "update-team", "--generate-cli-skeleton"); s.Err != "" || s.Code != 0 {
		t.Errorf("skeleton: %d %q", s.Code, s.Err)
	}
}
