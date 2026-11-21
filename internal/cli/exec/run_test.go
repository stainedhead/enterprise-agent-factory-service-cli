package exec_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli/services"
)

type result struct {
	code        int
	out, errOut string
}

// run executes afs against a handler with an API key in the environment.
func run(t *testing.T, h http.HandlerFunc, args ...string) result {
	t.Helper()
	srv := httptest.NewServer(h)
	t.Cleanup(srv.Close)
	return runAt(t, srv.URL, args...)
}

func runAt(t *testing.T, endpoint string, args ...string) result {
	t.Helper()
	reg, err := services.NewRegistry()
	if err != nil {
		t.Fatal(err)
	}
	var out, errb bytes.Buffer
	env := cli.NewEnv(strings.NewReader(""), &out, &errb)
	env.Vars = map[string]string{"AFS_ENDPOINT_URL": endpoint, "AFS_API_KEY": "afsk_unit_test_key"}
	env.Home = t.TempDir()
	env.Sleep = func(context.Context, time.Duration) error { return nil }
	code := cli.Run(context.Background(), args, env, reg)
	return result{code, out.String(), errb.String()}
}

func jsonReply(w http.ResponseWriter, status int, body string) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Request-Id", "req_1")
	w.WriteHeader(status)
	_, _ = w.Write([]byte(body))
}

// teamsPages serves 25 teams; nextToken is "o<offset>", default page size 10.
func teamsPages(mu *sync.Mutex, queries *[]string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		*queries = append(*queries, r.URL.RawQuery)
		mu.Unlock()
		start := 0
		if tk := r.URL.Query().Get("nextToken"); tk != "" {
			_, _ = fmt.Sscanf(tk, "o%d", &start)
		}
		n := 10
		if v := r.URL.Query().Get("maxResults"); v != "" {
			_, _ = fmt.Sscan(v, &n)
		}
		end := min(start+n, 25)
		var items []string
		for i := start; i < end; i++ {
			items = append(items, fmt.Sprintf(`{"id":"team_%02d"}`, i))
		}
		next := ""
		if end < 25 {
			next = fmt.Sprintf(`,"nextToken":"o%d"`, end)
		}
		jsonReply(w, 200, `{"items":[`+strings.Join(items, ",")+`]`+next+`}`)
	}
}

func count(out string) int { return strings.Count(out, `"id"`) }

func TestAutoPaginateMergesItems(t *testing.T) {
	var mu sync.Mutex
	var q []string
	r := run(t, teamsPages(&mu, &q), "teams", "list-teams")
	if r.code != 0 || count(r.out) != 25 || strings.Contains(r.out, "nextToken") {
		t.Fatalf("code=%d items=%d out=%s err=%s", r.code, count(r.out), r.out, r.errOut)
	}
}

func TestMaxItemsPrintsNextTokenAndStartingTokenResumes(t *testing.T) {
	var mu sync.Mutex
	var q []string
	r := run(t, teamsPages(&mu, &q), "teams", "list-teams", "--max-items", "10")
	if r.code != 0 || count(r.out) != 10 || !strings.Contains(r.out, "NextToken") {
		t.Fatalf("%+v", r)
	}
	var doc struct{ NextToken string }
	if err := json.Unmarshal([]byte(r.out), &doc); err != nil {
		t.Fatal(err)
	}
	r2 := run(t, teamsPages(&mu, &q), "teams", "list-teams", "--max-items", "10", "--starting-token", doc.NextToken)
	if r2.code != 0 || !strings.Contains(r2.out, "team_10") || strings.Contains(r2.out, "team_09") {
		t.Fatalf("%+v", r2)
	}
}

func TestMaxItemsMidPageResumesAfterItem(t *testing.T) {
	var mu sync.Mutex
	var q []string
	srv := httptest.NewServer(teamsPages(&mu, &q))
	defer srv.Close()
	r := runAt(t, srv.URL, "teams", "list-teams", "--max-items", "4", "--page-size", "10")
	var doc struct{ NextToken string }
	_ = json.Unmarshal([]byte(r.out), &doc)
	r2 := runAt(t, srv.URL, "teams", "list-teams", "--max-items", "2", "--page-size", "10", "--starting-token", doc.NextToken)
	if !strings.Contains(r2.out, "team_04") || strings.Contains(r2.out, "team_03") {
		t.Fatalf("%s", r2.out)
	}
}

func TestPageSizeSentAndRange(t *testing.T) {
	var mu sync.Mutex
	var q []string
	run(t, teamsPages(&mu, &q), "teams", "list-teams", "--page-size", "20", "--no-paginate")
	if len(q) != 1 || !strings.Contains(q[0], "maxResults=20") {
		t.Fatalf("%v", q)
	}
	r := run(t, teamsPages(&mu, &q), "teams", "list-teams", "--page-size", "101")
	if r.code != 252 {
		t.Fatalf("%+v", r)
	}
}

func TestNoPaginateReturnsFirstPageWithToken(t *testing.T) {
	var mu sync.Mutex
	var q []string
	r := run(t, teamsPages(&mu, &q), "teams", "list-teams", "--no-paginate")
	if count(r.out) != 10 || !strings.Contains(r.out, "NextToken") {
		t.Fatalf("%s", r.out)
	}
}

func TestBadStartingToken252(t *testing.T) {
	r := run(t, func(w http.ResponseWriter, _ *http.Request) { jsonReply(w, 200, `{"items":[]}`) }, "teams", "list-teams", "--starting-token", "!!")
	if r.code != 252 {
		t.Fatalf("%+v", r)
	}
}

func TestCSVFollowsNextTokenHeader(t *testing.T) {
	calls := 0
	r := run(t, func(w http.ResponseWriter, req *http.Request) {
		calls++
		if req.Header.Get("Accept") != "text/csv" {
			t.Errorf("accept %q", req.Header.Get("Accept"))
		}
		w.Header().Set("Content-Type", "text/csv")
		if req.URL.Query().Get("nextToken") == "" {
			w.Header().Set("X-Afs-Next-Token", "n2")
			_, _ = w.Write([]byte("a,b\n1,2\n"))
			return
		}
		_, _ = w.Write([]byte("a,b\n3,4\n"))
	}, "budgets", "list-costs", "--team-id", "team_1", "--output", "csv")
	if r.code != 0 || r.out != "a,b\n1,2\n3,4\n" || calls != 2 {
		t.Fatalf("%+v calls=%d", r, calls)
	}
}

func TestCSVRules252(t *testing.T) {
	h := func(w http.ResponseWriter, _ *http.Request) { jsonReply(w, 200, `{"items":[]}`) }
	for _, args := range [][]string{
		{"teams", "list-teams", "--output", "csv"},
		{"budgets", "list-costs", "--team-id", "team_1", "--output", "csv", "--query", "items"},
		{"budgets", "list-costs", "--team-id", "team_1", "--output", "csv", "--max-items", "3"},
	} {
		if r := run(t, h, args...); r.code != 252 {
			t.Errorf("%v: %+v", args, r)
		}
	}
}

func TestAccessDeniedExit254(t *testing.T) {
	r := run(t, func(w http.ResponseWriter, _ *http.Request) {
		jsonReply(w, 403, `{"code":"AccessDenied","message":"not allowed","requestId":"req_9","details":{"requiredPermission":"teams:read"}}`)
	}, "teams", "describe-team", "--team-id", "team_1")
	if r.code != 254 || !strings.Contains(r.errOut, "An error occurred (AccessDenied) when calling the GetTeam operation: not allowed") ||
		!strings.Contains(r.errOut, "teams:read") || !strings.Contains(r.errOut, "Request ID: req_9") {
		t.Fatalf("%+v", r)
	}
}

func TestUnauthorizedHint(t *testing.T) {
	r := run(t, func(w http.ResponseWriter, _ *http.Request) {
		jsonReply(w, 401, `{"code":"Unauthenticated","message":"bad token","requestId":"r"}`)
	}, "teams", "describe-team", "--team-id", "team_1")
	if r.code != 254 || !strings.Contains(r.errOut, "afs sso login --profile default") {
		t.Fatalf("%+v", r)
	}
}

func TestNetworkErrorExit255(t *testing.T) {
	r := runAt(t, "http://127.0.0.1:1", "teams", "describe-team", "--team-id", "team_1")
	if r.code != 255 || !strings.Contains(r.errOut, "afs: error") {
		t.Fatalf("%+v", r)
	}
}

func TestNoEndpointExit255(t *testing.T) {
	r := runAt(t, "", "teams", "describe-team", "--team-id", "team_1")
	if r.code != 255 || !strings.Contains(r.errOut, "no endpoint") {
		t.Fatalf("%+v", r)
	}
}

func TestDelete204PrintsNothingAndReadsVersion(t *testing.T) {
	var ifm string
	r := run(t, func(w http.ResponseWriter, req *http.Request) {
		if req.Method == "GET" {
			jsonReply(w, 200, `{"items":[{"id":"stage_0","version":2},{"id":"stage_1","version":7}]}`)
			return
		}
		ifm = req.Header.Get("If-Match")
		w.WriteHeader(204)
	}, "pools", "delete-stage", "--team-id", "team_1", "--stage-id", "stage_1")
	if r.code != 0 || r.out != "" || ifm != `"7"` || !strings.Contains(r.errOut, "version 7") {
		t.Fatalf("%+v ifm=%q", r, ifm)
	}
}

func TestIfMatchFlagSentAsGiven(t *testing.T) {
	var ifm string
	gets := 0
	r := run(t, func(w http.ResponseWriter, req *http.Request) {
		if req.Method == "GET" {
			gets++
		}
		ifm = req.Header.Get("If-Match")
		w.WriteHeader(204)
	}, "pools", "delete-stage", "--team-id", "team_1", "--stage-id", "stage_1", "--if-match", "3")
	if r.code != 0 || ifm != `"3"` || gets != 0 {
		t.Fatalf("%+v ifm=%q gets=%d", r, ifm, gets)
	}
}

func TestPrecondition412And428(t *testing.T) {
	for _, status := range []int{412, 428} {
		calls := 0
		r := run(t, func(w http.ResponseWriter, _ *http.Request) {
			calls++
			w.Header().Set("ETag", `"9"`)
			jsonReply(w, status, `{"code":"PreconditionFailed","message":"stale","requestId":"r"}`)
		}, "pools", "delete-stage", "--team-id", "team_1", "--stage-id", "s", "--if-match", "1")
		if r.code != 254 || calls != 1 {
			t.Fatalf("%d: %+v calls=%d", status, r, calls)
		}
	}
}

func TestSkeletonMakesNoCall(t *testing.T) {
	r := runAt(t, "http://127.0.0.1:1", "teams", "create-team", "--generate-cli-skeleton")
	var m map[string]any
	if r.code != 0 || json.Unmarshal([]byte(r.out), &m) != nil || m["slug"] == nil || m["displayName"] == nil {
		t.Fatalf("%+v", r)
	}
}

func TestCreateSendsIdempotencyKeyAndBody(t *testing.T) {
	var key, body string
	r := run(t, func(w http.ResponseWriter, req *http.Request) {
		key = req.Header.Get("Idempotency-Key")
		b := new(bytes.Buffer)
		_, _ = b.ReadFrom(req.Body)
		body = b.String()
		jsonReply(w, 201, `{"id":"team_1","name":"x"}`)
	}, "teams", "create-team", "--cli-input-json", `{"slug":"a","displayName":"A","firstAdmin":{"email":"a@b.c"},"maxBudgetUsd":"10"}`, "--kind", "persistent", "--idempotency-key", "k-1")
	if r.code != 0 || key != "k-1" || !strings.Contains(body, `"slug":"a"`) || !strings.Contains(body, `"kind":"persistent"`) {
		t.Fatalf("%+v key=%q body=%s", r, key, body)
	}
}

func TestEnumAndRequired252BeforeCall(t *testing.T) {
	called := false
	h := func(w http.ResponseWriter, _ *http.Request) { called = true; jsonReply(w, 200, `{}`) }
	if r := run(t, h, "teams", "create-team", "--kind", "bogus"); r.code != 252 {
		t.Fatalf("%+v", r)
	}
	if r := run(t, h, "teams", "describe-team"); r.code != 252 || !strings.Contains(r.errOut, "required") {
		t.Fatalf("%+v", r)
	}
	if called {
		t.Fatal("a call was made")
	}
}

func TestWaitReadyAndFailed(t *testing.T) {
	polls := 0
	phase := "Ready"
	h := func(w http.ResponseWriter, req *http.Request) {
		switch {
		case req.Method == "DELETE" || req.Method == "POST":
			jsonReply(w, 202, `{"id":"team_1","generation":3}`)
		case strings.HasSuffix(req.URL.Path, "/status"):
			polls++
			g := 2
			if polls > 1 {
				g = 3
			}
			jsonReply(w, 200, fmt.Sprintf(`{"status":{"phase":%q,"observedGeneration":%d}}`, phase, g))
		default:
			jsonReply(w, 200, `{"id":"team_1","slug":"final"}`)
		}
	}
	r := run(t, h, "teams", "resume-team", "--team-id", "team_1", "--wait")
	if r.code != 0 || !strings.Contains(r.out, "final") || polls != 2 {
		t.Fatalf("%+v polls=%d", r, polls)
	}
	phase, polls = "Failed", 1
	r = run(t, h, "teams", "resume-team", "--team-id", "team_1", "--wait")
	if r.code != 254 || !strings.Contains(r.errOut, "Failed") {
		t.Fatalf("%+v", r)
	}
}

func TestWaitDeleteUntil404AndTimeout(t *testing.T) {
	gets := 0
	r := run(t, func(w http.ResponseWriter, req *http.Request) {
		if req.Method == "DELETE" {
			jsonReply(w, 202, `{}`)
			return
		}
		gets++
		if gets < 3 {
			jsonReply(w, 200, `{"id":"team_1"}`)
			return
		}
		jsonReply(w, 404, `{"code":"NotFound","message":"gone","requestId":"r"}`)
	}, "teams", "delete-team", "--team-id", "team_1", "--wait", "--if-match", "1")
	if r.code != 0 || gets != 3 {
		t.Fatalf("%+v gets=%d", r, gets)
	}
	r = run(t, func(w http.ResponseWriter, req *http.Request) {
		if req.Method == "DELETE" {
			jsonReply(w, 202, `{}`)
			return
		}
		jsonReply(w, 200, `{"id":"team_1"}`)
	}, "teams", "delete-team", "--team-id", "team_1", "--wait", "--wait-timeout", "1s", "--if-match", "1")
	if r.code != 255 || !strings.Contains(r.errOut, "timed out") {
		t.Fatalf("%+v", r)
	}
}

func TestTableAndQueryOutput(t *testing.T) {
	var mu sync.Mutex
	var q []string
	r := run(t, teamsPages(&mu, &q), "teams", "list-teams", "--no-paginate", "--query", "items[0].id", "--output", "text")
	if r.code != 0 || strings.TrimSpace(r.out) != "team_00" {
		t.Fatalf("%+v", r)
	}
	r = run(t, teamsPages(&mu, &q), "teams", "list-teams", "--no-paginate", "--query", "[[", "--output", "text")
	if r.code != 252 {
		t.Fatalf("%+v", r)
	}
}

func TestReplayedNotice(t *testing.T) {
	r := run(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Idempotent-Replayed", "true")
		jsonReply(w, 201, `{"id":"team_1"}`)
	}, "teams", "create-team", "--cli-input-json", `{"slug":"a","displayName":"A","kind":"persistent","firstAdmin":{"email":"a@b.c"},"maxBudgetUsd":"10"}`)
	if !strings.Contains(r.out, "team_1") || !strings.Contains(r.errOut, "stored response") {
		t.Fatalf("%+v", r)
	}
}

func newServer(t *testing.T, h http.HandlerFunc) string {
	t.Helper()
	srv := httptest.NewServer(h)
	t.Cleanup(srv.Close)
	return srv.URL
}
