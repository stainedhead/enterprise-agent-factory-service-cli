package teams_test

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/clitest"
)

func ephemeralArgs(extra ...string) []string {
	base := []string{"teams", "create-team", "--slug", "cve-sprint", "--display-name", "CVE sprint", "--kind", "ephemeral", "--max-budget-usd", "500.00",
		"--first-admin", `{"email":"ada@example.com"}`}
	return append(base, extra...)
}

var okTeam = func(w http.ResponseWriter, _ *http.Request) {
	clitest.Reply(w, 201, `{"id":"team_1","slug":"cve-sprint","kind":"ephemeral"}`)
}

// CLI-US-8.1: leaving out one of --expires-at, --teardown-rule, --budget-usd is a usage error before any call.
func TestEphemeralRequiresAllThree252(t *testing.T) {
	all := map[string][]string{
		"--expires-at":    {"--expires-at", "2026-12-31T00:00:00Z"},
		"--teardown-rule": {"--teardown-rule", "at_cap_or_expiry"},
		"--budget-usd":    {"--budget-usd", "500.00"},
	}
	for missing := range all {
		var args []string
		for flag, v := range all {
			if flag != missing {
				args = append(args, v...)
			}
		}
		r := clitest.Run(t, okTeam, ephemeralArgs(args...)...)
		if r.Code != 252 || len(r.Requests) != 0 || !strings.Contains(r.Err, missing) {
			t.Errorf("without %s: %+v", missing, r)
		}
	}
	// a persistent team needs none of them, and the budget flags are refused on it
	r := clitest.Run(t, okTeam, "teams", "create-team", "--slug", "p", "--display-name", "P", "--kind", "persistent", "--max-budget-usd", "1.00", "--first-admin", `{"email":"a@example.com"}`)
	if r.Code != 0 {
		t.Fatalf("persistent: %+v", r)
	}
	r = clitest.Run(t, okTeam, "teams", "create-team", "--slug", "p", "--display-name", "P", "--kind", "persistent", "--max-budget-usd", "1.00", "--first-admin", `{"email":"a@example.com"}`, "--budget-usd", "5.00")
	if r.Code != 252 {
		t.Fatalf("budget flag on a persistent team: %+v", r)
	}
}

// CLI-17, D-C19: the budget is sent as {period: lifetime, amountUsd, thresholds}.
func TestEphemeralBudgetBody(t *testing.T) {
	r := clitest.Run(t, okTeam, ephemeralArgs("--expires-at", "2026-12-31T00:00:00Z", "--teardown-rule", "at_cap_or_expiry", "--budget-usd", "500.00", "--budget-thresholds", "50", "80", "100")...)
	if r.Code != 0 {
		t.Fatalf("%+v", r)
	}
	var body struct {
		Budget map[string]any `json:"budget"`
	}
	if err := json.Unmarshal(r.ReqBody, &body); err != nil {
		t.Fatal(err)
	}
	b := body.Budget
	th, _ := b["thresholds"].([]any)
	if b["period"] != "lifetime" || b["amountUsd"] != "500.00" || len(th) != 3 || th[0] != float64(50) || th[2] != float64(100) || len(b) != 3 {
		t.Errorf("budget %v (body %s)", b, r.ReqBody)
	}
}

// CLI-17: the full object through --cli-input-json is accepted and checked the same way.
func TestEphemeralViaInputJSON(t *testing.T) {
	dir := t.TempDir()
	full := filepath.Join(dir, "team.json")
	doc := `{"slug":"cve-sprint","displayName":"CVE sprint","kind":"ephemeral","maxBudgetUsd":"500.00","firstAdmin":{"email":"ada@example.com"},"expiresAt":"2026-12-31T00:00:00Z","teardownRule":"at_work_complete","budget":{"period":"lifetime","amountUsd":"500.00"}}`
	if err := os.WriteFile(full, []byte(doc), 0o600); err != nil {
		t.Fatal(err)
	}
	r := clitest.Run(t, okTeam, "teams", "create-team", "--cli-input-json", "file://"+full)
	if r.Code != 0 || !strings.Contains(string(r.ReqBody), `"budget":{"period":"lifetime","amountUsd":"500.00"}`) {
		t.Fatalf("%+v body=%s", r, r.ReqBody)
	}
	// without the budget object the same document is refused
	part := filepath.Join(dir, "part.json")
	if err := os.WriteFile(part, []byte(strings.Replace(doc, `,"budget":{"period":"lifetime","amountUsd":"500.00"}`, "", 1)), 0o600); err != nil {
		t.Fatal(err)
	}
	r = clitest.Run(t, okTeam, "teams", "create-team", "--cli-input-json", "file://"+part)
	if r.Code != 252 || len(r.Requests) != 0 || !strings.Contains(r.Err, "--budget-usd") {
		t.Fatalf("%+v", r)
	}
	// and --budget-usd completes it
	r = clitest.Run(t, okTeam, "teams", "create-team", "--cli-input-json", "file://"+part, "--budget-usd", "250.00")
	if r.Code != 0 || !strings.Contains(string(r.ReqBody), `"amountUsd":"250.00"`) {
		t.Fatalf("%+v body=%s", r, r.ReqBody)
	}
}
