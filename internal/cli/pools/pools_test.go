package pools_test

import (
	"net/http"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/clitest"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/golden"
)

func TestCommandsDoneStages(t *testing.T) {
	clitest.CheckOperations(t, "stages", "ListStages", "CreateStage", "UpdateStage", "DeleteStage")
}

func TestCommandsDonePools(t *testing.T) {
	clitest.CheckOperations(t, "pools", "ListPools", "CreatePool", "GetPool", "UpdatePool", "DeletePool")
}

func TestCommandsDoneScale(t *testing.T) {
	clitest.CheckOperations(t, "scale", "ScalePool", "PausePool", "ResumePool", "GetPoolStatus")
}

func TestCommandsDoneBindings(t *testing.T) {
	clitest.CheckOperations(t, "bindings", "ListPoolHandoffs", "PutPoolHandoffs", "ListPoolSecretBindings", "PutPoolSecretBindings",
		"ListPoolConfigValues", "PutPoolConfigValue", "DeletePoolConfigValue")
}

func TestCommandsDoneLinks(t *testing.T) {
	clitest.CheckOperations(t, "links", "ListPoolExternalSystems", "PutPoolExternalSystem", "ListPoolRoleLinks", "CreatePoolRoleLink", "VerifyPoolRoleLink", "DeletePoolRoleLink")
}

const pool = `{"id":"pool_01JD2Q3R4S5T6V7W8X9Y0Z0001","name":"fixer","stageId":"stg_01JD2Q3R4S5T6V7W8X9Y0Z0001","template":"claude-code","workers":10,"state":"enabled","generation":2,"version":3,"status":{"phase":"Ready","observedGeneration":2}}`

// CLI-US-8.4: create-pool on a persistent team is the API's Conflict, 254.
func TestCreatePoolPersistentConflict254(t *testing.T) {
	r := clitest.Run(t, clitest.Fail(409, "Conflict", "pools exist only on ephemeral teams"),
		"pools", "create-pool", "--team-id", "team_1", "--name", "fixer", "--template", "claude-code", "--image", "registry.example.com/a@sha256:"+strings.Repeat("a", 64),
		"--model-id", "m1", "--workers", "10")
	if r.Code != 254 || !strings.Contains(r.Err, "(Conflict) when calling the CreatePool operation") {
		t.Fatalf("%+v", r)
	}
}

// CLI-US-8.3 (D-C31): describe-pool shows desired workers, state and phase.
func TestDescribePoolTableGolden(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) { clitest.Reply(w, 200, pool) },
		"pools", "describe-pool", "--team-id", "team_1", "--pool-id", "pool_1", "--output", "table")
	if r.Code != 0 || !strings.Contains(r.Out, "Workers") || !strings.Contains(r.Out, "enabled") || !strings.Contains(r.Out, "Ready") {
		t.Fatalf("%+v", r)
	}
	golden.Check(t, "pools/GetPool.desired.table.golden", r.Out)
}

// CLI-1: the scratch size option reads well (scratch-size-gib), and pool names are in the path only for ids.
func TestCreatePoolScratchSizeOption(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) { clitest.Reply(w, 201, pool) },
		"pools", "create-pool", "--team-id", "team_1", "--name", "fixer", "--template", "claude-code", "--image", "registry.example.com/a@sha256:"+strings.Repeat("a", 64),
		"--model-id", "m1", "--workers", "10", "--scratch-size-gib", "20")
	if r.Code != 0 || !strings.Contains(string(r.ReqBody), `"scratchSizeGiB":20`) {
		t.Fatalf("%+v body=%s", r, r.ReqBody)
	}
}

// CLI-US-8.2: --workers 250 is a usage error before any call (the range is 0 to 200).
func TestScalePoolRange252(t *testing.T) {
	r := clitest.Run(t, nil, "pools", "scale-pool", "--team-id", "team_1", "--pool-id", "pool_1", "--workers", "250")
	if r.Code != 252 || len(r.Requests) != 0 || !strings.Contains(r.Err, "250") {
		t.Fatalf("%+v", r)
	}
}

// CLI-US-8.2: --workers 20 sends ScalePool with an Idempotency-Key.
func TestScalePoolIdempotencyKey(t *testing.T) {
	h := func(w http.ResponseWriter, _ *http.Request) { clitest.Reply(w, 200, pool) }
	r := clitest.Run(t, h, "pools", "scale-pool", "--team-id", "team_1", "--pool-id", "pool_1", "--workers", "20")
	if r.Code != 0 || r.Req.Method != "POST" || r.Req.URL.Path != "/v1/teams/team_1/pools/pool_1:scale" || string(r.ReqBody) != `{"workers":20}` || r.Req.Header.Get("Idempotency-Key") == "" {
		t.Fatalf("%+v body=%s", r, r.ReqBody)
	}
	r = clitest.Run(t, h, "pools", "scale-pool", "--team-id", "team_1", "--pool-id", "pool_1", "--workers", "20", "--idempotency-key", "scale-1")
	if r.Req.Header.Get("Idempotency-Key") != "scale-1" {
		t.Errorf("key %q", r.Req.Header.Get("Idempotency-Key"))
	}
}

func status(reporting int, phase string) string {
	return `{"id":"pool_1","generation":2,"status":{"phase":"` + phase + `","observedGeneration":2},"workers":{"desired":20,"reporting":` + strconv.Itoa(reporting) + `,"working":0,"idle":0,"impaired":0}}`
}

// CLI-US-8.2, CLI-12: --wait returns 0 when 20 of 20 workers report.
func TestScalePoolWaitReporting(t *testing.T) {
	var mu sync.Mutex
	polls := 0
	r := clitest.Run(t, func(w http.ResponseWriter, req *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		switch {
		case req.Method == "POST":
			clitest.Reply(w, 200, `{"id":"pool_1","workers":20,"generation":2,"status":{"phase":"Pending","observedGeneration":1}}`)
		case strings.HasSuffix(req.URL.Path, "/status"):
			polls++
			counts := []int{5, 12, 20}
			clitest.Reply(w, 200, status(counts[min(polls-1, 2)], "Ready"))
		default:
			clitest.Reply(w, 200, pool)
		}
	}, "pools", "scale-pool", "--team-id", "team_1", "--pool-id", "pool_1", "--workers", "20", "--wait")
	if r.Code != 0 || polls != 3 {
		t.Fatalf("%+v polls=%d", r, polls)
	}
}

// CLI-US-8.2: --wait-timeout before the workers report is exit 255 and says the change stands.
func TestScalePoolWaitTimeout255(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, req *http.Request) {
		if req.Method == "POST" {
			clitest.Reply(w, 200, `{"id":"pool_1","workers":20,"generation":2}`)
			return
		}
		clitest.Reply(w, 200, status(5, "Ready"))
	}, "pools", "scale-pool", "--team-id", "team_1", "--pool-id", "pool_1", "--workers", "20", "--wait", "--wait-timeout", "1s")
	if r.Code != 255 || !strings.Contains(r.Err, "timed out") || !strings.Contains(r.Err, "the change stands") {
		t.Fatalf("%+v", r)
	}
}

// CLI-US-8.3: get-pool-status shows desired, reporting, working, idle and impaired.
func TestGetPoolStatusTableGolden(t *testing.T) {
	body := `{"id":"pool_01JD2Q3R4S5T6V7W8X9Y0Z0001","generation":3,"status":{"phase":"Ready","observedGeneration":3},"workers":{"desired":20,"reporting":18,"working":11,"idle":7,"impaired":2}}`
	r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) { clitest.Reply(w, 200, body) },
		"pools", "get-pool-status", "--team-id", "team_1", "--pool-id", "pool_1", "--output", "table")
	if r.Code != 0 {
		t.Fatalf("%+v", r)
	}
	for _, h := range []string{"Desired", "Reporting", "Working", "Idle", "Impaired"} {
		if !strings.Contains(r.Out, h) {
			t.Errorf("no %s column:\n%s", h, r.Out)
		}
	}
	golden.Check(t, "scale/GetPoolStatus.counts.table.golden", r.Out)
}

// CLI-17: --to-pool-ids is a list; CLI-10: the PUT reads the version and sends If-Match.
func TestPutHandoffsListAndIfMatch(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) {
		clitest.Reply(w, 200, `{"toPoolIds":["pool_2","pool_3"],"version":4}`)
	},
		"pools", "put-handoffs", "--team-id", "team_1", "--pool-id", "pool_1", "--to-pool-ids", "pool_2", "pool_3")
	if r.Code != 0 || len(r.Requests) != 2 || r.Req.Method != "PUT" || string(r.ReqBody) != `{"toPoolIds":["pool_2","pool_3"]}` || r.Req.Header.Get("If-Match") == "" {
		t.Fatalf("%+v body=%s", r, r.ReqBody)
	}
}

// CLI-17: a role link is created by role name and account; the trust is pending until verified.
func TestCreatePoolRoleLink(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) {
		clitest.Reply(w, 201, `{"id":"prl_1","accountId":"111122223333","roleName":"afs-cve-reader","profileName":"target","status":"pending_trust","externalId":"afs-7Q2M9K"}`)
	}, "pools", "create-role-link", "--team-id", "team_1", "--pool-id", "pool_1", "--account-id", "111122223333", "--role-name", "afs-cve-reader", "--profile-name", "target", "--output", "table")
	if r.Code != 0 || !strings.Contains(r.Out, "pending_trust") || !strings.Contains(r.Out, "afs-7Q2M9K") {
		t.Fatalf("%+v", r)
	}
}
