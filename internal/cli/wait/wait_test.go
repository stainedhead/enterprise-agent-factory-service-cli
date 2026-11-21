package wait

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/client"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/spec"
)

type fake struct {
	replies []any // *client.Response, or error
	calls   []string
}

func (f *fake) Do(_ context.Context, r client.Request) (*client.Response, error) {
	f.calls = append(f.calls, r.Method+" "+r.Path)
	if len(f.replies) == 0 {
		return nil, errors.New("no reply")
	}
	x := f.replies[0]
	if len(f.replies) > 1 {
		f.replies = f.replies[1:]
	}
	if e, ok := x.(error); ok {
		return nil, e
	}
	return x.(*client.Response), nil
}

func ok(body string) *client.Response { return &client.Response{Status: 200, Body: []byte(body)} }

func opts() Options {
	now := time.Unix(0, 0)
	return Options{Timeout: time.Minute, Interval: time.Second,
		Now:   func() time.Time { return now },
		Sleep: func(_ context.Context, d time.Duration) error { now = now.Add(d); return nil }}
}

func TestEvaluate(t *testing.T) {
	cases := []struct {
		body string
		gen  int64
		done bool
		err  string
	}{
		{`{"status":{"phase":"Ready","observedGeneration":3}}`, 3, true, ""},
		{`{"status":{"phase":"Ready","observedGeneration":2}}`, 3, false, ""},
		{`{"phase":"Ready"}`, 0, true, ""},
		{`{"status":{"phase":"Pending"}}`, 0, false, ""},
		{`{"status":{"phase":"Failed","conditions":[{"type":"X"}]}}`, 0, false, "Failed"},
		{`{"status":{"phase":"Degraded"}}`, 0, false, "Degraded"},
		{`not json`, 0, false, ""},
	}
	for _, c := range cases {
		done, err := Evaluate([]byte(c.body), c.gen)
		if done != c.done || (err != nil) != (c.err != "") || (err != nil && !strings.Contains(err.Error(), c.err)) {
			t.Errorf("%s: done=%v err=%v", c.body, done, err)
		}
	}
	if Generation([]byte(`{"generation":4}`)) != 4 || Generation([]byte(`{}`)) != 0 || Generation([]byte(`x`)) != 0 {
		t.Error("Generation")
	}
}

func TestTableOpsExistInSnapshot(t *testing.T) {
	for id, tg := range Table {
		for _, o := range []string{id, tg.PollOp, tg.GetOp} {
			if spec.ByID(o) == nil {
				t.Errorf("%s: operation %s is not in the snapshot", id, o)
			}
		}
	}
}

func TestIgnoresStaleReadyByGeneration(t *testing.T) {
	f := &fake{replies: []any{ok(`{"status":{"phase":"Ready","observedGeneration":1}}`), ok(`{"status":{"phase":"Ready","observedGeneration":2}}`), ok(`{"final":true}`)}}
	tg := Table["UpdateTeam"]
	out, err := Until(context.Background(), f, tg, []byte(`{"generation":2}`), "/poll", "/get", opts())
	if err != nil || !strings.Contains(string(out), "final") || len(f.calls) != 3 {
		t.Fatalf("%v %s %v", err, out, f.calls)
	}
}

func TestTimeoutAndCancel(t *testing.T) {
	f := &fake{replies: []any{ok(`{"status":{"phase":"Pending"}}`)}}
	_, err := Until(context.Background(), f, Table["UpdateTeam"], nil, "/p", "/g", opts())
	if !errors.Is(err, ErrTimeout) {
		t.Fatalf("%v", err)
	}
	o := opts()
	o.Sleep = func(context.Context, time.Duration) error { return context.Canceled }
	_, err = Until(context.Background(), f, Table["UpdateTeam"], nil, "/p", "/g", o)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("%v", err)
	}
	f = &fake{replies: []any{errors.New("boom")}}
	if _, err = Until(context.Background(), f, Table["UpdateTeam"], nil, "/p", "/g", opts()); err == nil {
		t.Fatal("poll error swallowed")
	}
}

func TestDeleteWaitErrors(t *testing.T) {
	f := &fake{replies: []any{&client.APIError{Status: 500, Code: "InternalError"}}}
	if _, err := Until(context.Background(), f, Table["DeleteTeam"], nil, "/p", "/g", opts()); err == nil {
		t.Fatal("a 500 ended the wait")
	}
	f = &fake{replies: []any{&client.APIError{Status: 404, Code: "NotFound"}}}
	if out, err := Until(context.Background(), f, Table["DeleteTeam"], nil, "/p", "/g", opts()); err != nil || out != nil {
		t.Fatalf("%v %s", err, out)
	}
}

func TestFailedError(t *testing.T) {
	f := &fake{replies: []any{ok(`{"status":{"phase":"Failed","conditions":[{"type":"Ready"}]}}`)}}
	_, err := Until(context.Background(), f, Table["CreateTeam"], nil, "/p", "/g", opts())
	var fe *FailedError
	if !errors.As(err, &fe) || !strings.Contains(err.Error(), "conditions") {
		t.Fatalf("%v", err)
	}
}

func TestPathsFromRequestThenResponseID(t *testing.T) {
	poll, get := Paths(Table["CreateAgent"], map[string]string{"teamId": "team_1"}, []byte(`{"id":"agt_9"}`))
	if poll != "/v1/teams/team_1/agents/agt_9/status" || get != "/v1/teams/team_1/agents/agt_9" {
		t.Fatalf("%s %s", poll, get)
	}
	if p, _ := Paths(Target{PollOp: "Nope"}, nil, nil); p != "" {
		t.Fatal("unknown op")
	}
}
