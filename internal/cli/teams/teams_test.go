package teams_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/clitest"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/spec"
)

func TestCommandsDoneA(t *testing.T) {
	clitest.CheckOperations(t, "teams", "ListTeams", "GetTeam", "GetTeamStatus", "CreateTeam")
}

// CLI-1, CLI-US-4.6: --state and --kind filters reach the query; --query shapes the output.
func TestListTeamsFiltersAndQuery(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) {
		clitest.Reply(w, 200, `{"items":[{"id":"team_1","displayName":"One"},{"id":"team_2","displayName":"Two"}]}`)
	}, "teams", "list-teams", "--state", "active", "--kind", "persistent", "--query", "items[].id", "--output", "text")
	if r.Code != 0 || r.Out != "team_1\tteam_2\n" {
		t.Fatalf("%+v", r)
	}
	q := r.Req.URL.Query()
	if q.Get("state") != "active" || q.Get("kind") != "persistent" {
		t.Fatalf("query %v", q)
	}
}

// CLI-US-2.1: the skeleton lists every CreateTeam field.
func TestCreateTeamSkeleton(t *testing.T) {
	r := clitest.Run(t, nil, "teams", "create-team", "--generate-cli-skeleton")
	if r.Code != 0 {
		t.Fatalf("%+v", r)
	}
	var got map[string]any
	if err := json.Unmarshal([]byte(r.Out), &got); err != nil {
		t.Fatal(err)
	}
	for _, f := range spec.ByID("CreateTeam").Body.Fields {
		if f.ReadOnly {
			continue
		}
		if _, ok := got[f.Name]; !ok {
			t.Errorf("skeleton lacks %s", f.Name)
		}
	}
}

// CLI-US-2.1: file://team.json is sent unchanged (members, order and values).
func TestCreateTeamInputFileUnchanged(t *testing.T) {
	doc := `{"slug":"payments","displayName":"Payments","kind":"persistent","maxBudgetUsd":"500.00","firstAdmin":{"email":"ada@example.com"},"tags":{"project":"checkout"}}`
	f := filepath.Join(t.TempDir(), "team.json")
	if err := os.WriteFile(f, []byte(doc), 0o600); err != nil {
		t.Fatal(err)
	}
	r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) {
		clitest.Reply(w, 201, `{"id":"team_1"}`)
	}, "teams", "create-team", "--cli-input-json", "file://"+f)
	if r.Code != 0 || string(r.ReqBody) != doc || r.Req.Method != "POST" || r.Req.URL.Path != "/v1/teams" {
		t.Fatalf("%+v body=%s", r, r.ReqBody)
	}
	if r.Req.Header.Get("Idempotency-Key") == "" {
		t.Error("no Idempotency-Key")
	}
}

// CLI-US-2.2: a missing required tag is the API's ValidationError: 254, naming the tag.
func TestCreateTeamMissingTag254(t *testing.T) {
	r := clitest.Run(t, clitest.Fail(422, "ValidationError", "required tag costCenter is missing"),
		"teams", "create-team", "--slug", "p", "--display-name", "P", "--kind", "persistent", "--max-budget-usd", "1.00", "--first-admin", `{"email":"a@example.com"}`)
	if r.Code != 254 || !strings.Contains(r.Err, "(ValidationError) when calling the CreateTeam operation: required tag costCenter is missing") {
		t.Fatalf("%+v", r)
	}
}

func waitHandler(phase string, polls *int) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		switch {
		case req.Method == "POST":
			clitest.Reply(w, 201, `{"id":"team_1","generation":1}`)
		case strings.HasSuffix(req.URL.Path, "/status"):
			*polls++
			clitest.Reply(w, 200, fmt.Sprintf(`{"id":"team_1","generation":1,"status":{"phase":%q,"observedGeneration":1,"conditions":[{"type":"Reconciled","status":"False","reason":"QuotaFailed"}]}}`, phase))
		default:
			clitest.Reply(w, 200, `{"id":"team_1","slug":"payments","status":{"phase":"Ready"}}`)
		}
	}
}

var createArgs = []string{"teams", "create-team", "--slug", "payments", "--display-name", "P", "--kind", "persistent", "--max-budget-usd", "1.00", "--first-admin", `{"email":"a@example.com"}`, "--wait"}

// CLI-US-2.3: --wait prints the final team once Ready.
func TestCreateTeamWaitReady(t *testing.T) {
	polls := 0
	r := clitest.Run(t, waitHandler("Ready", &polls), createArgs...)
	if r.Code != 0 || !strings.Contains(r.Out, `"slug": "payments"`) || polls < 1 {
		t.Fatalf("%+v polls=%d", r, polls)
	}
}

// CLI-US-2.3: Failed exits 254 with the conditions on stderr.
func TestCreateTeamWaitFailed(t *testing.T) {
	polls := 0
	r := clitest.Run(t, waitHandler("Failed", &polls), createArgs...)
	if r.Code != 254 || !strings.Contains(r.Err, "Failed") || !strings.Contains(r.Err, "QuotaFailed") {
		t.Fatalf("%+v", r)
	}
}
