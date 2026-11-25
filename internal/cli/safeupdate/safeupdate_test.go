package safeupdate

import (
	"context"
	"errors"
	"net/url"
	"strings"
	"testing"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/client"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/spec"
)

type fake struct {
	byPath map[string][]any
	calls  []string
}

func (f *fake) Do(_ context.Context, r client.Request) (*client.Response, error) {
	f.calls = append(f.calls, r.Path+"?"+r.Query.Encode())
	q := r.Query.Get("nextToken")
	rs := f.byPath[r.Path+"|"+q]
	if rs == nil {
		return nil, &client.APIError{Status: 404, Code: "NotFound"}
	}
	if e, ok := rs[0].(error); ok {
		return nil, e
	}
	return &client.Response{Status: 200, Body: []byte(rs[0].(string))}, nil
}

func TestReaderTableCoversAllIfMatchOperations(t *testing.T) {
	var withReader, without []string
	for i := range spec.All() {
		op := &spec.All()[i]
		if op.IfMatch == spec.IfMatchNone {
			continue
		}
		if r := ReaderFor(op); r != nil {
			withReader = append(withReader, op.ID)
			if r.Get == nil || r.Get.Method != "GET" {
				t.Errorf("%s: reader is not a GET", op.ID)
			}
		} else {
			without = append(without, op.ID)
		}
	}
	if len(withReader)+len(without) != 74 {
		t.Fatalf("%d If-Match operations, want 74", len(withReader)+len(without))
	}
	t.Logf("no reader (send no If-Match): %s", strings.Join(without, ", "))
	// Every operation without a reader takes an optional If-Match (CD-2), so the write is sent without one
	// unless --if-match is given; v1.5.0 adds five (no readable version: credentials, federation, gateway
	// settings, selections).
	want := "PutEndpointFederation, PutModelEndpointCredential, PutServiceModelAccess, PutSetting, PutTeamGatewaySettings, PutTeamModelEndpointCredential, PutTeamModelEndpoints, PutTeamModels"
	if got := strings.Join(without, ", "); got != want {
		t.Errorf("operations without a reader:\n got %s\nwant %s", got, want)
	}
	for _, id := range without {
		if spec.ByID(id).IfMatch != spec.IfMatchOptional {
			t.Errorf("%s has no reader but requires If-Match", id)
		}
	}
}

func TestDirectRead(t *testing.T) {
	f := &fake{byPath: map[string][]any{"/v1/teams/team_1|": {`{"id":"team_1","version":7}`}}}
	v, ok, err := Version(context.Background(), f, spec.ByID("DeleteTeam"), map[string]string{"teamId": "team_1"})
	if err != nil || !ok || v != "7" {
		t.Fatalf("%q %v %v", v, ok, err)
	}
	if _, ok, err := Version(context.Background(), f, spec.ByID("DeleteTeam"), map[string]string{"teamId": "gone"}); ok || err != nil {
		t.Fatalf("404 must send no If-Match: %v %v", ok, err)
	}
}

func TestListReaderPaginatesToFindItem(t *testing.T) {
	f := &fake{byPath: map[string][]any{
		"/v1/teams/team_1/stages|":   {`{"items":[{"id":"s0","version":1}],"nextToken":"n2"}`},
		"/v1/teams/team_1/stages|n2": {`{"items":[{"id":"s1","version":9}]}`},
	}}
	v, ok, err := Version(context.Background(), f, spec.ByID("DeleteStage"), map[string]string{"teamId": "team_1", "stageId": "s1"})
	if err != nil || !ok || v != "9" || len(f.calls) != 2 {
		t.Fatalf("%q %v %v %v", v, ok, err, f.calls)
	}
	if _, ok, _ := Version(context.Background(), f, spec.ByID("DeleteStage"), map[string]string{"teamId": "team_1", "stageId": "zzz"}); ok {
		t.Fatal("found a missing item")
	}
}

func TestErrorsAndOddBodies(t *testing.T) {
	boom := errors.New("boom")
	f := &fake{byPath: map[string][]any{"/v1/teams/t|": {boom}, "/v1/teams/team_1/stages|": {boom}}}
	if _, _, err := Version(context.Background(), f, spec.ByID("DeleteTeam"), map[string]string{"teamId": "t"}); err == nil {
		t.Error("read error swallowed")
	}
	if _, _, err := Version(context.Background(), f, spec.ByID("DeleteStage"), map[string]string{"teamId": "team_1", "stageId": "s"}); err == nil {
		t.Error("list error swallowed")
	}
	for body, want := range map[string]string{`{"version":"3"}`: "3", `{"version":4}`: "4", `{}`: "", `[]`: ""} {
		v, ok, _ := versionOf([]byte(body))
		if v != want || ok != (want != "") {
			t.Errorf("%s: %q %v", body, v, ok)
		}
	}
	if Notice("5") != "Applied to version 5; use --if-match to guard against concurrent changes." {
		t.Error("notice")
	}
	_ = url.Values{}
}

func TestOperationWithoutReader(t *testing.T) {
	op := spec.ByID("UpdateWorkOrderRepo")
	if ReaderFor(op) != nil {
		t.Skip("has a reader")
	}
	v, ok, err := Version(context.Background(), &fake{}, op, nil)
	if v != "" || ok || err != nil {
		t.Fatal("no reader must mean no version")
	}
}

func repoPath() map[string]string {
	return map[string]string{"workOrderId": "wo_1", "host": "github.com", "owner": "acme", "repo": "ledger"}
}

// FR-019: the composite reader of a work order's repositories matches host and owner/repo across pages.
func TestCompositeReaderFindsRepositoryAcrossPages(t *testing.T) {
	op := spec.ByID("UpdateWorkOrderRepo")
	f := &fake{byPath: map[string][]any{
		"/v1/work-orders/wo_1/repos|":   {`{"items":[{"host":"github.com","repo":"acme/other","version":1}],"nextToken":"n2"}`},
		"/v1/work-orders/wo_1/repos|n2": {`{"items":[{"host":"gitlab.com","repo":"acme/ledger","version":2},{"host":"github.com","repo":"acme/ledger","version":8}]}`},
	}}
	v, ok, err := Version(context.Background(), f, op, repoPath())
	if err != nil || !ok || v != "8" || len(f.calls) != 2 {
		t.Fatalf("%q %v %v %v", v, ok, err, f.calls)
	}
	// not in the list, an unusable page, and a gone work order send no If-Match
	for name, f := range map[string]*fake{
		"absent":   {byPath: map[string][]any{"/v1/work-orders/wo_1/repos|": {`{"items":[]}`}}},
		"not json": {byPath: map[string][]any{"/v1/work-orders/wo_1/repos|": {`<html>`}}},
		"404":      {},
	} {
		if v, ok, err := Version(context.Background(), f, op, repoPath()); ok || err != nil {
			t.Errorf("%s: %q %v %v", name, v, ok, err)
		}
	}
}

func TestCompositeReaderDeniedAndFailed(t *testing.T) {
	op := spec.ByID("UpdateWorkOrderRepo")
	denied := &client.APIError{Status: 403, Code: "AccessDenied", Message: "not allowed"}
	f := &fake{byPath: map[string][]any{"/v1/work-orders/wo_1/repos|": {denied}}}
	_, _, err := Version(context.Background(), f, op, repoPath())
	var rd *ReadDeniedError
	if !errors.As(err, &rd) {
		t.Fatalf("%v", err)
	}
	if rd.Error() != denied.Error() || rd.Unwrap() != error(denied) { //nolint:errorlint // comparing the unwrapped value itself
		t.Errorf("Error/Unwrap: %q", rd.Error())
	}
	var ae *client.APIError
	if !errors.As(err, &ae) || ae.Status != 403 {
		t.Errorf("does not unwrap to the API error: %v", err)
	}
	f = &fake{byPath: map[string][]any{"/v1/work-orders/wo_1/repos|": {&client.APIError{Status: 500, Code: "InternalError"}}}}
	if _, _, err := Version(context.Background(), f, op, repoPath()); err == nil || errors.As(err, &rd) {
		t.Errorf("a 500 must be a plain error: %v", err)
	}
}

// FR-019: when two items collide on different fields, the id beats the name, and two items with the
// same id are an error rather than a guess.
func TestListMatchPrefersIDAndRefusesAmbiguity(t *testing.T) {
	path := map[string]string{"teamId": "team_1", "stageId": "s1"}
	op := spec.ByID("DeleteStage")
	f := &fake{byPath: map[string][]any{"/v1/teams/team_1/stages|": {`{"items":[{"id":"a","name":"s1","version":1},{"id":"s1","version":9}]}`}}}
	if v, ok, err := Version(context.Background(), f, op, path); err != nil || !ok || v != "9" {
		t.Fatalf("id must beat name: %q %v %v", v, ok, err)
	}
	f = &fake{byPath: map[string][]any{"/v1/teams/team_1/stages|": {`{"items":[{"id":"s1","version":1},{"id":"s1","version":2}]}`}}}
	if _, ok, err := Version(context.Background(), f, op, path); err == nil || ok || !strings.Contains(err.Error(), "--if-match") {
		t.Fatalf("ambiguous items must be an error naming --if-match: %v %v", ok, err)
	}
	// same name on two items, no id match: ambiguous as well
	f = &fake{byPath: map[string][]any{"/v1/teams/team_1/stages|": {`{"items":[{"id":"a","name":"s1","version":1},{"id":"b","name":"s1","version":2}]}`}}}
	if _, ok, err := Version(context.Background(), f, op, path); err == nil || ok {
		t.Fatalf("ambiguous names: %v %v", ok, err)
	}
	// a single name match still works
	f = &fake{byPath: map[string][]any{"/v1/teams/team_1/stages|": {`{"items":[{"id":"a","name":"s1","version":4},{"id":"b","name":"x","version":2}]}`}}}
	if v, ok, err := Version(context.Background(), f, op, path); err != nil || !ok || v != "4" {
		t.Fatalf("name match: %q %v %v", v, ok, err)
	}
}

func TestListWithOddPagesSendsNoVersion(t *testing.T) {
	path := map[string]string{"teamId": "team_1", "stageId": "s1"}
	op := spec.ByID("DeleteStage")
	for name, body := range map[string]string{"not json": `<html>`, "no items": `{}`, "item not an object": `{"items":["x",3]}`} {
		f := &fake{byPath: map[string][]any{"/v1/teams/team_1/stages|": {body}}}
		if v, ok, err := Version(context.Background(), f, op, path); ok || err != nil {
			t.Errorf("%s: %q %v %v", name, v, ok, err)
		}
	}
}

// CLI-US-34.3 (ADR-008): PutTeamCollaboration reads GetTeamCollaboration at the same path; the two
// new POST actions have no reader and no If-Match, and no operation of the collaboration service
// outside the Put has one.
func TestReaderTableMatchesSnapshot(t *testing.T) {
	put := spec.ByID("PutTeamCollaboration")
	r := ReaderFor(put)
	if r == nil || r.Get == nil || r.Get.ID != "GetTeamCollaboration" || r.Path != put.Path || r.Match != "" {
		t.Fatalf("reader of PutTeamCollaboration = %+v", r)
	}
	for _, id := range []string{"ConnectAgentToCollaboration", "DisconnectAgentFromCollaboration", "RotateCollaborationKeys"} {
		op := spec.ByID(id)
		if op.IfMatch != spec.IfMatchNone { // the executor reads a version only when the operation takes If-Match
			t.Errorf("%s must take no If-Match, so no version is read", id)
		}
	}
}

// FR-008 (v1.5.0): an inference provider has no GET of its own; its version is read from the provider list (not
// paginated) by the item's provider.
func TestInferenceProviderVersionFromList(t *testing.T) {
	op := spec.ByID("UpdateInferenceProvider")
	r := ReaderFor(op)
	if r == nil || r.Get.ID != "ListInferenceProviders" || r.Match != "provider" {
		t.Fatalf("reader = %+v", r)
	}
	f := &fake{byPath: map[string][]any{"/v1/admin/inference-providers|": {`{"items":[{"provider":"bedrock","version":3},{"provider":"azure","version":8}]}`}}}
	v, ok, err := Version(context.Background(), f, op, map[string]string{"provider": "azure"})
	if err != nil || !ok || v != "8" {
		t.Fatalf("v=%q ok=%v err=%v", v, ok, err)
	}
	if len(f.calls) != 1 || f.calls[0] != "/v1/admin/inference-providers?" {
		t.Errorf("requests %v", f.calls)
	}
	for _, id := range []string{"PutTeamModelEndpoints", "PutTeamModelEndpointCredential", "PutModelEndpointCredential"} {
		if ReaderFor(spec.ByID(id)) != nil {
			t.Errorf("%s must read no version", id)
		}
	}
}
