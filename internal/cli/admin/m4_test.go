package admin_test

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/clitest"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/golden"
)

func TestCommandsDoneServiceChannels(t *testing.T) {
	clitest.CheckOperations(t, "service-channels", "ListServiceNotificationChannels", "CreateServiceNotificationChannel", "UpdateServiceNotificationChannel", "DeleteServiceNotificationChannel")
}

func TestCommandsDoneServiceTags(t *testing.T) {
	clitest.CheckOperations(t, "service-tags", "ListServiceTags", "PutServiceTag", "DeleteServiceTag")
}

func TestCommandsDonePriceProposals(t *testing.T) {
	clitest.CheckOperations(t, "price-proposals", "ListModelPriceProposals", "ApproveModelPriceProposal", "RejectModelPriceProposal")
}

func TestCommandsDoneReconciliations(t *testing.T) {
	clitest.CheckOperations(t, "reconciliations", "ListCostReconciliations", "ListOrgAccounts")
}

func TestCommandsDoneExports(t *testing.T) {
	clitest.CheckOperations(t, "exports", "ListExports", "CreateExport", "GetExport", "CreateImport")
}

// CLI-1: deleting a channel that a setting names is the API's 409, printed and exit 254.
func TestDeleteServiceChannelNamedInSetting409(t *testing.T) {
	r := clitest.Run(t, clitest.Fail(409, "Conflict", "the channel is used by setting notifications.default_channel_id"),
		"admin", "delete-notification-channel", "--channel-id", "nch_1", "--if-match", "2")
	if r.Code != 254 || !strings.Contains(r.Err, "(Conflict) when calling the DeleteServiceNotificationChannel operation") || !strings.Contains(r.Err, "default_channel_id") {
		t.Fatalf("%+v", r)
	}
}

// CLI-1: --mode fixed|required and --allowed-values as a list.
func TestPutServiceTagModeAndAllowedValues(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) {
		clitest.Reply(w, 200, `{"key":"cost-center","mode":"required","version":2}`)
	},
		"admin", "put-service-tag", "--key", "cost-center", "--mode", "required", "--allowed-values", "CC-1", "CC-2", "--if-match", "1")
	if r.Code != 0 || string(r.ReqBody) != `{"allowedValues":["CC-1","CC-2"],"mode":"required"}` && string(r.ReqBody) != `{"mode":"required","allowedValues":["CC-1","CC-2"]}` {
		t.Fatalf("%+v body=%s", r, r.ReqBody)
	}
	if r.Req.URL.Path != "/v1/admin/service-tags/cost-center" {
		t.Errorf("path %s", r.Req.URL.Path)
	}
	r = clitest.Run(t, nil, "admin", "put-service-tag", "--key", "k", "--mode", "sometimes")
	if r.Code != 252 || len(r.Requests) != 0 {
		t.Fatalf("enum: %+v", r)
	}
}

const proposals = `{"items":[{"id":"1042","modelId":"anthropic.claude-sonnet-4-5","effectiveFrom":"2026-12-01T00:00:00Z","inputUsdPerMtok":"3.00","outputUsdPerMtok":"15.00","status":"pending","version":1},{"id":"1043","modelId":"anthropic.claude-haiku-4-5","effectiveFrom":"2026-12-01T00:00:00Z","inputUsdPerMtok":"1.00","outputUsdPerMtok":"5.00","status":"pending","version":1}]}`

// CLI-US-15.1: one table row per pending proposal.
func TestListPriceProposalsTableGolden(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) { clitest.Reply(w, 200, proposals) },
		"admin", "list-price-proposals", "--status", "pending", "--output", "table")
	if r.Code != 0 || r.Req.URL.Query().Get("status") != "pending" || len(strings.Split(strings.TrimSpace(r.Out), "\n")) != 3 {
		t.Fatalf("%+v", r)
	}
	golden.Check(t, "price-proposals/ListModelPriceProposals.two.table.golden", r.Out)
}

// CLI-US-15.2: approve and reject send the note and print the decided proposal.
func TestApprovePriceProposal(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) {
		clitest.Reply(w, 200, `{"id":"1042","status":"approved","note":"ok","decidedBy":"usr_1","version":2}`)
	}, "admin", "approve-price-proposal", "--proposal-id", "1042", "--note", "ok", "--output", "table")
	if r.Code != 0 || r.Req.URL.Path != "/v1/admin/model-price-proposals/1042:approve" || string(r.ReqBody) != `{"note":"ok"}` || !strings.Contains(r.Out, "approved") {
		t.Fatalf("%+v body=%s", r, r.ReqBody)
	}
}

func TestRejectPriceProposal(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) {
		clitest.Reply(w, 200, `{"id":"1042","status":"rejected","note":"too high","version":2}`)
	}, "admin", "reject-price-proposal", "--proposal-id", "1042", "--note", "too high", "--output", "table")
	if r.Code != 0 || r.Req.URL.Path != "/v1/admin/model-price-proposals/1042:reject" || !strings.Contains(r.Out, "rejected") {
		t.Fatalf("%+v", r)
	}
}

func TestApproveAccessDeniedTeamAdmin(t *testing.T) {
	r := clitest.Run(t, clitest.Denied("admin:pricing"), "admin", "approve-price-proposal", "--proposal-id", "1042")
	if r.Code != 254 || !strings.Contains(r.Err, "(AccessDenied) when calling the ApproveModelPriceProposal operation") {
		t.Fatalf("%+v", r)
	}
}

// CLI-US-18.1: create-export prints the id; describe-export prints sha256, size and the download URL.
func TestExportCreateAndDescribe(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) {
		clitest.Reply(w, 202, `{"id":"exp_1","status":"pending","teamIds":["team_1"]}`)
	}, "admin", "create-export", "--team-ids", "team_1")
	if r.Code != 0 || !strings.Contains(r.Out, "exp_1") || string(r.ReqBody) != `{"teamIds":["team_1"]}` {
		t.Fatalf("%+v", r)
	}
	done := `{"id":"exp_1","status":"succeeded","sizeBytes":48213,"sha256":"9f9f","downloadUrl":"https://afs-exports.example.com/exports/exp_1.json?X-Amz-Signature=abc","teamIds":[]}`
	r = clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) { clitest.Reply(w, 200, done) },
		"admin", "describe-export", "--export-id", "exp_1", "--output", "table")
	for _, want := range []string{"9f9f", "48213", "https://afs-exports.example.com/exports/exp_1.json?X-Amz-Signature=abc"} {
		if r.Code != 0 || !strings.Contains(r.Out, want) {
			t.Fatalf("describe lacks %s:\n%+v", want, r)
		}
	}
	if r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) { clitest.Reply(w, 200, `{"items":[]}`) }, "admin", "list-exports"); r.Code != 0 {
		t.Fatalf("list: %+v", r)
	}
}

func importFiles(t *testing.T) (bundle, accounts string) {
	t.Helper()
	dir := t.TempDir()
	bundle, accounts = filepath.Join(dir, "bundle.json"), filepath.Join(dir, "accounts.json")
	if err := os.WriteFile(bundle, []byte(`{"schemaVersion":14,"sha256":"ab","teams":[]}`), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(accounts, []byte(`{"111122223333":"444455556666"}`), 0o600); err != nil {
		t.Fatal(err)
	}
	return
}

var importReply = func(w http.ResponseWriter, _ *http.Request) {
	clitest.Reply(w, 200, `{"dryRun":true,"id":"exp_2","teams":[{"slug":"payments","result":"updated","changes":[],"error":null}]}`)
}

// CLI-US-18.2: an import is a dry run unless --no-dry-run.
func TestCreateImportDryRunDefault(t *testing.T) {
	bundle, accounts := importFiles(t)
	r := clitest.Run(t, importReply, "admin", "create-import", "--bundle", "file://"+bundle, "--account-map", "file://"+accounts)
	if r.Code != 0 || r.Req.URL.Query().Get("dryRun") != "true" || r.Req.URL.Path != "/v1/imports" {
		t.Fatalf("%+v", r)
	}
}

func TestCreateImportNoDryRun(t *testing.T) {
	bundle, accounts := importFiles(t)
	r := clitest.Run(t, importReply, "admin", "create-import", "--bundle", "file://"+bundle, "--account-map", "file://"+accounts, "--no-dry-run")
	if r.Code != 0 || r.Req.URL.Query().Get("dryRun") != "false" {
		t.Fatalf("%+v", r)
	}
}

// D-C18: the body wraps the files as {bundle, accountMap}.
func TestCreateImportBodyWrapsFiles(t *testing.T) {
	bundle, accounts := importFiles(t)
	r := clitest.Run(t, importReply, "admin", "create-import", "--bundle", "file://"+bundle, "--account-map", "file://"+accounts)
	var got struct {
		Bundle     map[string]any    `json:"bundle"`
		AccountMap map[string]string `json:"accountMap"`
	}
	if err := json.Unmarshal(r.ReqBody, &got); err != nil || got.Bundle["schemaVersion"] != float64(14) || got.AccountMap["111122223333"] != "444455556666" || len(got.Bundle) != 3 {
		t.Fatalf("body %s (%v)", r.ReqBody, err)
	}
}

// CLI-US-18.3: a bundle over 100 MB is a usage error before any call and before it is read.
func TestCreateImportOver100MB252(t *testing.T) {
	f := filepath.Join(t.TempDir(), "big.json")
	fh, err := os.Create(f) //nolint:gosec // test fixture
	if err != nil {
		t.Fatal(err)
	}
	if err := fh.Truncate(101 << 20); err != nil { // sparse
		t.Skipf("no sparse files: %v", err)
	}
	_ = fh.Close()
	r := clitest.Run(t, importReply, "admin", "create-import", "--bundle", "file://"+f)
	if r.Code != 252 || len(r.Requests) != 0 || !strings.Contains(r.Err, "100 MB") {
		t.Fatalf("%+v", r)
	}
}
