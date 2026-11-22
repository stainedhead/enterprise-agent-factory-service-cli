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
	if len(withReader)+len(without) != 59 {
		t.Fatalf("%d If-Match operations, want 59", len(withReader)+len(without))
	}
	t.Logf("no reader (send no If-Match): %s", strings.Join(without, ", "))
	if len(without) > 6 {
		t.Errorf("too many operations without a reader: %v", without)
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
