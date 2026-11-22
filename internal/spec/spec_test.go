package spec_test

import (
	"crypto/sha256"
	"encoding/hex"
	"regexp"
	"strings"
	"testing"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/api"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/spec"
)

func TestSnapshotCounts(t *testing.T) {
	var all, cmds, pag, csv, im, idem int
	re := regexp.MustCompile(`^[a-z]+(-[a-z]+)*$`)
	for _, o := range spec.All() {
		all++
		if o.Service != "" {
			cmds++
			if !re.MatchString(o.Command) || !re.MatchString(o.Service) {
				t.Errorf("%s: bad command %q %q", o.ID, o.Service, o.Command)
			}
		}
		if o.Paginated {
			pag++
			if o.ResultKey == "" {
				t.Errorf("%s paginated without ResultKey", o.ID)
			}
		}
		if o.CSV {
			csv++
		}
		if o.IfMatch != spec.IfMatchNone {
			im++
			switch o.Method {
			case "PATCH", "DELETE":
				if o.IfMatch != spec.IfMatchRequired {
					t.Errorf("%s: want required", o.ID)
				}
			case "PUT":
				if o.IfMatch != spec.IfMatchOptional {
					t.Errorf("%s: want optional", o.ID)
				}
			}
		}
		if o.Idempotent {
			idem++
			if o.Method != "POST" {
				t.Errorf("%s: idempotent non-POST", o.ID)
			}
		}
	}
	for name, got := range map[string][2]int{"ops": {all, 182}, "cmds": {cmds, 180}, "paginated": {pag, 52}, "csv": {csv, 8}, "ifmatch": {im, 59}, "idem": {idem, 44}} {
		if got[0] != got[1] {
			t.Errorf("%s = %d, want %d", name, got[0], got[1])
		}
	}
}

func TestContractVersionAndChecksum(t *testing.T) {
	if spec.ContractVersion != "1.1.0" {
		t.Errorf("ContractVersion = %s", spec.ContractVersion)
	}
	h := sha256.Sum256(api.OpenAPI)
	if spec.SnapshotSHA256 != hex.EncodeToString(h[:]) {
		t.Error("SnapshotSHA256 does not match the embedded snapshot")
	}
	for _, l := range strings.Split(string(api.Source), "\n") {
		if v, ok := strings.CutPrefix(l, "sha256: "); ok && v != spec.SnapshotSHA256 {
			t.Error("SnapshotSHA256 does not match api/SOURCE")
		}
	}
}

func TestLookups(t *testing.T) {
	o := spec.ByID("UpdateAgent")
	if o == nil || o.Service != "agents" || o.Command != "update-agent" || o.Body == nil {
		t.Fatalf("UpdateAgent = %+v", o)
	}
	if spec.ByID("Nope") != nil {
		t.Error("unknown id found")
	}
	if p := spec.ByID("GetTeamReport"); p.ResultKey != "rows" {
		t.Errorf("report ResultKey = %q", p.ResultKey)
	}
	if p := spec.ByID("DeleteAgent"); len(p.Success) != 1 || p.Success[0] != 202 {
		t.Errorf("DeleteAgent success = %v", p.Success)
	}
	cs := spec.Commands()
	if len(cs) != 180 {
		t.Errorf("commands = %d", len(cs))
	}
	if spec.ByID("UpdateWorkOrderRepo").AnyOf == nil {
		t.Error("AnyOf missing")
	}
}
