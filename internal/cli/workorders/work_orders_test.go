package workorders_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/clitest"
)

const (
	tmpl  = "swt_01J9ZQ3ABCDEFGHJKMNPQRSTVW"
	team1 = "team_01J9ZQ3ABCDEFGHJKMNPQRSTVW"
	wo    = "wo_01J9ZQ3ABCDEFGHJKMNPQRSTVW"
)

func TestCommandsDoneWorkOrders(t *testing.T) {
	clitest.CheckOperations(t, "workorders", "CreateWorkOrder", "ListWorkOrders", "GetWorkOrder", "UpdateWorkOrder",
		"CancelWorkOrder", "ListWorkOrderRepos", "UpdateWorkOrderRepo", "ListTeamWorkOrders")
}

func create(extra ...string) []string {
	return append([]string{"work-orders", "create-work-order", "--template-id", tmpl, "--title", "Patch", "--severity", "critical",
		"--sla-due-at", "2026-12-01T00:00:00Z"}, extra...)
}

const created = `{"id":"` + wo + `","title":"Patch","severity":"critical","status":"open","teamId":null,"version":1,"repos":[` +
	`{"host":"github.com","repo":"acme/ledger","status":"pending","ownerSource":"team_config","owningTeamId":"` + team1 + `"},` +
	`{"host":"github.com","repo":"acme/billing","status":"pending","ownerSource":"explicit","owningContact":"a@example.com"}]}`

// CLI-US-9.1 and CLI-18: one CreateWorkOrder with the owners as classified; the key is a security-operator key.
func TestCreateWorkOrderReposBody(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) { clitest.Reply(w, 201, created) },
		create("--repos", "github.com/acme/ledger", "github.com/acme/billing=a@example.com", "--idempotency-key", "k1")...)
	want := `{"templateId":"` + tmpl + `","title":"Patch","severity":"critical","slaDueAt":"2026-12-01T00:00:00Z","repos":[{"host":"github.com","repo":"acme/ledger"},{"host":"github.com","repo":"acme/billing","owningContact":"a@example.com"}]}`
	if r.Code != 0 || len(r.Requests) != 1 || r.Req.Method != "POST" || r.Req.URL.Path != "/v1/work-orders" || !strings.Contains(string(r.ReqBody), `"repos":[{"host":"github.com","repo":"acme/ledger"},{"host":"github.com","repo":"acme/billing","owningContact":"a@example.com"}]`) {
		t.Fatalf("%+v body=%s want %s", r, r.ReqBody, want)
	}
	if r.Req.Header.Get("Idempotency-Key") != "k1" || r.Req.Header.Get("Authorization") == "" {
		t.Errorf("headers %v", r.Req.Header)
	}
}

// CLI-US-9.1: the same --idempotency-key replays the first answer.
func TestCreateWorkOrderReplay(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Idempotent-Replayed", "true")
		clitest.Reply(w, 201, created)
	}, create("--repos", "github.com/acme/ledger", "--idempotency-key", "k1")...)
	if r.Code != 0 || r.Req.Header.Get("Idempotency-Key") != "k1" || !strings.Contains(r.Out, wo) {
		t.Fatalf("%+v", r)
	}
}

// CLI-US-9.3: no '=' sends no owner fields; the response's ownerSource is printed; explicit =team_ is sent.
func TestCreateWorkOrderTeamConfigOwner(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) { clitest.Reply(w, 201, created) },
		create("--repos", "github.com/acme/ledger", "github.com/acme/billing="+team1)...)
	if r.Code != 0 || !strings.Contains(string(r.ReqBody), `{"host":"github.com","repo":"acme/ledger"},`) ||
		!strings.Contains(string(r.ReqBody), `"owningTeamId":"`+team1+`"`) || !strings.Contains(r.Out, `"ownerSource": "team_config"`) {
		t.Fatalf("%+v body=%s", r, r.ReqBody)
	}
}

// CLI-US-9.2: a team slug is a usage error (252) and nothing is sent.
func TestCreateWorkOrderSlugOwner252(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) { clitest.Reply(w, 201, created) },
		create("--repos", "github.com/acme/ledger=payments")...)
	if r.Code != 252 || len(r.Requests) != 0 || !strings.Contains(r.Err, "use a team id (team_...) or an email address") {
		t.Fatalf("%+v", r)
	}
}

// CLI-US-9.4: ambiguous repositories name the repositories and teams and say how to choose; exit 254.
func TestCreateWorkOrderAmbiguous254(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) {
		clitest.Reply(w, 422, `{"code":"ValidationError","message":"repositories need an owner","requestId":"req_a","details":{"ambiguousRepos":[{"repo":"github.com/acme/ledger","teamIds":["team_A","team_B"]}]}}`)
	}, create("--repos", "github.com/acme/ledger")...)
	for _, want := range []string{"(ValidationError) when calling the CreateWorkOrder operation", "github.com/acme/ledger", "team_A, team_B", "host/owner/repo=<team_id or email>"} {
		if !strings.Contains(r.Err, want) {
			t.Errorf("stderr lacks %q: %s", want, r.Err)
		}
	}
	if r.Code != 254 {
		t.Errorf("exit %d", r.Code)
	}
}

// CLI-US-9.4: unowned repositories are named too.
func TestCreateWorkOrderUnowned254(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) {
		clitest.Reply(w, 422, `{"code":"ValidationError","message":"repositories need an owner","requestId":"req_a","details":{"unownedRepos":["github.com/acme/ledger","github.com/acme/billing"]}}`)
	}, create("--repos", "github.com/acme/ledger", "github.com/acme/billing")...)
	if r.Code != 254 || !strings.Contains(r.Err, "Unowned repositories: github.com/acme/ledger, github.com/acme/billing") || !strings.Contains(r.Err, "=<team_id or email>") {
		t.Fatalf("%+v", r)
	}
}

// CLI-US-9.5: a second open work order for the same CVE and repository is a Conflict.
func TestCreateWorkOrderDuplicateCVE409(t *testing.T) {
	r := clitest.Run(t, clitest.Fail(409, "Conflict", "an open work order exists for CVE-2026-1234 on github.com/acme/ledger"),
		create("--cve-id", "CVE-2026-1234", "--kind", "cve_remediation", "--repos", "github.com/acme/ledger")...)
	if r.Code != 254 || !strings.Contains(r.Err, "(Conflict) when calling the CreateWorkOrder operation") {
		t.Fatalf("%+v", r)
	}
}

// CLI-US-9.6: the service key on a team command is refused by the API (254).
func TestServiceKeyOnTeamCommandDenied(t *testing.T) {
	r := clitest.Run(t, clitest.Denied("agents:read"), "agents", "list-agents", "--team-id", team1)
	if r.Code != 254 || !strings.Contains(r.Err, "(AccessDenied)") {
		t.Fatalf("%+v", r)
	}
}

// CLI-US-22.2: a disabled template is a ValidationError (254).
func TestDisabledTemplate254(t *testing.T) {
	r := clitest.Run(t, clitest.Fail(422, "ValidationError", "swarm template is disabled"), create("--repos", "github.com/acme/ledger")...)
	if r.Code != 254 || !strings.Contains(r.Err, "(ValidationError) when calling the CreateWorkOrder operation: swarm template is disabled") {
		t.Fatalf("%+v", r)
	}
}

// CLI-6 (CSV): list-work-orders --output csv asks for text/csv and prints it unchanged.
func TestListWorkOrdersCSV(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/csv")
		w.WriteHeader(200)
		_, _ = w.Write([]byte("id,title\n" + wo + ",Patch\n"))
	}, "work-orders", "list-work-orders", "--output", "csv")
	if r.Code != 0 || r.Req.Header.Get("Accept") != "text/csv" || !strings.HasPrefix(r.Out, "id,title\n"+wo) {
		t.Fatalf("%+v", r)
	}
}

// update-work-order sends add-repos in the same classified form.
func TestUpdateWorkOrderAddRepos(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, req *http.Request) {
		clitest.Reply(w, 200, `{"id":"`+wo+`","version":4}`)
	}, "work-orders", "update-work-order", "--work-order-id", wo, "--if-match", "3", "--add-repos", "github.com/acme/gateway="+team1)
	if r.Code != 0 || r.Req.Header.Get("If-Match") != `"3"` || !strings.Contains(string(r.ReqBody), `"addRepos":[{"host":"github.com","repo":"acme/gateway","owningTeamId":"`+team1+`"}]`) {
		t.Fatalf("%+v body=%s", r, r.ReqBody)
	}
}

// CLI-US-10.2: a workload token on cancel-work-order is refused by the API (254).
func TestCancelWorkOrderWorkloadDenied(t *testing.T) {
	r := clitest.Run(t, clitest.Denied("work-orders:write"), "work-orders", "cancel-work-order", "--work-order-id", wo, "--reason", "no longer needed")
	if r.Code != 254 || !strings.Contains(r.Err, "(AccessDenied) when calling the CancelWorkOrder operation") || !strings.Contains(r.Err, "work-orders:write") {
		t.Fatalf("%+v", r)
	}
}

// list-team-work-orders: a security operator is refused by the API.
func TestListTeamWorkOrdersDenied(t *testing.T) {
	r := clitest.Run(t, clitest.Denied("work-orders:read-team"), "work-orders", "list-team-work-orders", "--team-id", team1)
	if r.Code != 254 || r.Req.URL.Path != "/v1/teams/"+team1+"/work-orders" {
		t.Fatalf("%+v", r)
	}
}

// ---- update-repo ----

const repoList = `{"items":[{"host":"github.com","repo":"acme/other","status":"pending","version":9},{"host":"github.com","repo":"acme/ledger","status":"in_progress","version":3}]}`

// CLI-US-10.1: --repo host/owner/repo becomes three path segments.
func TestUpdateRepoPathSegments(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, req *http.Request) {
		if req.Method == "GET" {
			clitest.Reply(w, 200, repoList)
			return
		}
		clitest.Reply(w, 200, `{"host":"github.com","repo":"acme/ledger","status":"pr_opened","version":4}`)
	}, "work-orders", "update-repo", "--work-order-id", wo, "--repo", "github.com/acme/ledger", "--status", "pr_opened",
		"--pr-url", "https://github.com/acme/ledger/pull/42", "--pr-number", "42")
	patch := r.Requests[len(r.Requests)-1]
	if r.Code != 0 || patch.Method != "PATCH" || patch.URL.Path != "/v1/work-orders/"+wo+"/repos/github.com/acme/ledger" {
		t.Fatalf("%+v", r)
	}
	if patch.Header.Get("If-Match") != `"3"` {
		t.Errorf("If-Match %q: the composite host and repo match must pick version 3", patch.Header.Get("If-Match"))
	}
	if string(r.ReqBody) != `{"status":"pr_opened","prUrl":"https://github.com/acme/ledger/pull/42","prNumber":42}` && !strings.Contains(string(r.ReqBody), `"status":"pr_opened"`) {
		t.Errorf("body %s", r.ReqBody)
	}
}

// CLI-US-10.1, D-C17: with --if-match nothing is read; the workload token is the bearer.
func TestUpdateRepoWithIfMatch(t *testing.T) {
	r := clitest.RunWith(t, clitest.Options{Vars: map[string]string{"AFS_API_KEY": "", "AFS_API_TOKEN_FILE": tokenFile(t)}},
		func(w http.ResponseWriter, _ *http.Request) {
			clitest.Reply(w, 200, `{"host":"github.com","repo":"acme/ledger","status":"merged","version":4}`)
		},
		"work-orders", "update-repo", "--work-order-id", wo, "--repo", "github.com/acme/ledger", "--status", "merged", "--if-match", "3")
	if r.Code != 0 || len(r.Requests) != 1 || r.Req.Method != "PATCH" || r.Req.Header.Get("If-Match") != `"3"` || r.Req.Header.Get("Authorization") != "Bearer workload-token" {
		t.Fatalf("%+v", r)
	}
}

// D-C17: when the read is denied the API's answer is followed by the --if-match hint (254), and no PATCH is sent.
func TestUpdateRepoReadDeniedHint254(t *testing.T) {
	r := clitest.Run(t, clitest.Denied("work-orders:read"), "work-orders", "update-repo", "--work-order-id", wo,
		"--repo", "github.com/acme/ledger", "--status", "pr_opened")
	if r.Code != 254 || len(r.Requests) != 1 || r.Req.Method != "GET" ||
		!strings.Contains(r.Err, "(AccessDenied)") || !strings.Contains(r.Err, "Pass --if-match <version> (the version printed by the previous update-repo).") {
		t.Fatalf("%+v", r)
	}
}

// update-repo takes --host --owner --repo separately too, and refuses a mix.
func TestUpdateRepoSeparateParts(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) {
		clitest.Reply(w, 200, `{"host":"github.com","repo":"acme/ledger","status":"skipped","version":5}`)
	},
		"work-orders", "update-repo", "--work-order-id", wo, "--host", "github.com", "--owner", "acme", "--repo", "ledger", "--status", "skipped", "--if-match", "4")
	if r.Code != 0 || r.Req.URL.Path != "/v1/work-orders/"+wo+"/repos/github.com/acme/ledger" {
		t.Fatalf("%+v", r)
	}
	m := clitest.Run(t, nil, "work-orders", "update-repo", "--work-order-id", wo, "--host", "github.com", "--repo", "github.com/acme/ledger", "--status", "skipped")
	if m.Code != 252 {
		t.Errorf("mixed forms: exit %d", m.Code)
	}
	n := clitest.Run(t, nil, "work-orders", "update-repo", "--work-order-id", wo, "--status", "skipped")
	if n.Code != 252 || !strings.Contains(n.Err, "--repo host/owner/repo") {
		t.Errorf("no repo: %+v", n)
	}
}
