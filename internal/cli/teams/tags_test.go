package teams_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/clitest"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/golden"
)

const teamTags = `{"version":3,"requestMetadataDroppedKeys":[],"tags":[{"key":"afs:team","source":"system","value":"payments"},{"key":"billing-code","source":"service","value":"B-77"},{"key":"cost-center","source":"team","value":"CC-1001"},{"key":"owner","source":"service_default","value":"unassigned"},{"key":"project","source":"missing"}]}`

func TestCommandsDoneTags(t *testing.T) {
	clitest.CheckOperations(t, "tags", "GetTeamTags", "PutTeamTags")
}

// CLI-US-16.1: --tags k=v pairs send the team-owned keys as {key, value} objects.
func TestPutTeamTagsShorthand(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) { clitest.Reply(w, 200, teamTags) },
		"teams", "put-team-tags", "--team-id", "team_1", "--tags", "cost-center=CC-2002", "project=checkout", "--if-match", "3")
	want := `{"tags":[{"key":"cost-center","value":"CC-2002"},{"key":"project","value":"checkout"}]}`
	if r.Code != 0 || string(r.ReqBody) != want || r.Req.Method != "PUT" {
		t.Fatalf("%+v body=%s", r, r.ReqBody)
	}
	// JSON is accepted too, and a pair without "=" is a usage error before any call
	r = clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) { clitest.Reply(w, 200, teamTags) },
		"teams", "put-team-tags", "--team-id", "team_1", "--tags", `[{"key":"a","value":"b"}]`, "--if-match", "3")
	if r.Code != 0 || string(r.ReqBody) != `{"tags":[{"key":"a","value":"b"}]}` {
		t.Fatalf("json: %+v", r)
	}
	r = clitest.Run(t, nil, "teams", "put-team-tags", "--team-id", "team_1", "--tags", "oops")
	if r.Code != 252 || len(r.Requests) != 0 || !strings.Contains(r.Err, "key=value") {
		t.Fatalf("no =: %+v", r)
	}
}

// CLI-US-16.2: the table shows every tag with its source, including system, service and missing.
func TestGetTeamTagsTableGolden(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) { clitest.Reply(w, 200, teamTags) },
		"teams", "get-team-tags", "--team-id", "team_1", "--output", "table")
	if r.Code != 0 {
		t.Fatalf("%+v", r)
	}
	for _, src := range []string{"system", "service", "team", "service_default", "missing"} {
		if !strings.Contains(r.Out, src) {
			t.Errorf("source %s missing:\n%s", src, r.Out)
		}
	}
	golden.Check(t, "tags/GetTeamTags.sources.table.golden", r.Out)
}

// CLI-US-16.3: a reserved prefix or a service-owned key is refused by the API (254).
func TestPutTeamTagsReservedPrefix254(t *testing.T) {
	r := clitest.Run(t, clitest.Fail(400, "ValidationError", "tag key afs:team uses a reserved prefix"),
		"teams", "put-team-tags", "--team-id", "team_1", "--tags", "afs:team=x", "--if-match", "3")
	if r.Code != 254 || !strings.Contains(r.Err, "(ValidationError) when calling the PutTeamTags operation") {
		t.Fatalf("%+v", r)
	}
}

// CLI-US-6.3: without --if-match the CLI reads the version and sends it; with it, the given one.
func TestPutTeamTagsIfMatch(t *testing.T) {
	h := func(w http.ResponseWriter, _ *http.Request) { clitest.Reply(w, 200, teamTags) }
	r := clitest.Run(t, h, "teams", "put-team-tags", "--team-id", "team_1", "--tags", "project=x")
	if r.Code != 0 || r.Req.Method != "PUT" || r.Req.Header.Get("If-Match") != `"3"` || len(r.Requests) != 2 {
		t.Fatalf("read first: %+v", r)
	}
	r = clitest.Run(t, h, "teams", "put-team-tags", "--team-id", "team_1", "--tags", "project=x", "--if-match", "7")
	if r.Code != 0 || len(r.Requests) != 1 || r.Req.Header.Get("If-Match") != `"7"` {
		t.Fatalf("explicit: %+v", r)
	}
}
