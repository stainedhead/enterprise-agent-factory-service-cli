package cli

import (
	"strings"
	"testing"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/client"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/output"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/spec"
)

func TestKebab(t *testing.T) {
	for in, want := range map[string]string{"teamId": "team-id", "amountUsd": "amount-usd", "iamRoleArn": "iam-role-arn", "name": "name", "sizeGiB": "size-gi-b"} {
		if got := Kebab(in); got != want {
			t.Errorf("Kebab(%s) = %s, want %s", in, got, want)
		}
	}
}

// CLI-1: commands take ids; every parameter and writable body member becomes an option.
func TestDeriveOptions(t *testing.T) {
	find := func(opts []Option, flag string) *Option {
		for i := range opts {
			if opts[i].Flag == flag {
				return &opts[i]
			}
		}
		return nil
	}
	op := spec.ByID("UpdateAgent")
	opts := DeriveOptions(op, nil)
	if o := find(opts, "--team-id"); o == nil || o.In != "path" || !o.Required || o.Kind != String {
		t.Errorf("--team-id = %+v", o)
	}
	if o := find(opts, "--if-match"); o == nil || o.In != "header" {
		t.Error("--if-match missing")
	}
	if find(opts, "--idempotency-key") != nil {
		t.Error("PATCH has no idempotency key")
	}
	if find(opts, "--replicas") != nil {
		t.Error("--replicas exists")
	}
	if o := find(opts, "--volume"); o == nil || o.Kind != JSON {
		t.Errorf("--volume = %+v", o)
	}
	b := DeriveOptions(spec.ByID("PutBudget"), nil)
	if o := find(b, "--thresholds"); o == nil || o.Kind != List || o.Elem != "integer" {
		t.Errorf("--thresholds = %+v", o)
	}
	if o := find(b, "--amount-usd"); o == nil || o.Kind != String {
		t.Errorf("--amount-usd = %+v", o)
	}
	s := DeriveOptions(spec.ByID("CreateSecret"), nil)
	if o := find(s, "--secret-value"); o == nil || !o.Sensitive {
		t.Errorf("--secret-value = %+v", o)
	}
	if find(s, "--idempotency-key") == nil {
		t.Error("POST has idempotency key")
	}
	pv := DeriveOptions(spec.ByID("PutSetting"), nil)
	if o := find(pv, "--value"); o == nil || o.Kind != Any {
		t.Errorf("--value = %+v", o)
	}
	ov := &Override{Rename: map[string]string{"id": "channel-id"}, Extra: []Option{{Flag: "--repos", In: "local", Kind: List}}}
	n := DeriveOptions(spec.ByID("GetNotificationChannel"), ov)
	if find(n, "--channel-id") == nil || find(n, "--repos") == nil || find(n, "--id") != nil {
		t.Error("rename or extra not applied")
	}
	if o := find(DeriveOptions(spec.ByID("ListTeams"), nil), "--state"); o == nil || len(o.Enum) == 0 {
		t.Error("enum missing")
	}
}

// Every derived option list is free of collisions with itself and with the globals.
func TestAllOperationsDeriveCleanly(t *testing.T) {
	r := NewRegistry(spec.All())
	if err := r.Finalize(); err != nil {
		t.Fatal(err)
	}
}

func TestOverrideConflicts(t *testing.T) {
	cols := []output.Column{{Header: "Id", Path: "id"}}
	r := NewRegistry(spec.All())
	r.Override(Override{Op: "GetTeam", Columns: cols})
	r.Override(Override{Op: "GetTeam", Columns: cols})
	if err := r.Finalize(); err == nil || !strings.Contains(err.Error(), "Columns") {
		t.Errorf("two Columns: %v", err)
	}
	r = NewRegistry(spec.All())
	r.Override(Override{Op: "GetTeam", After: func(*Invocation, *client.Response) int { return -1 }})
	r.Override(Override{Op: "GetTeam", After: func(*Invocation, *client.Response) int { return -1 }})
	if err := r.Finalize(); err == nil {
		t.Error("two After accepted")
	}
	r = NewRegistry(spec.All())
	r.Override(Override{Op: "GetTeam", Rename: map[string]string{"teamId": "x"}})
	r.Override(Override{Op: "GetTeam", Rename: map[string]string{"teamId": "y"}})
	if err := r.Finalize(); err == nil {
		t.Error("two Rename accepted")
	}
	r = NewRegistry(spec.All())
	r.Override(Override{Op: "Nope"})
	if err := r.Finalize(); err == nil {
		t.Error("unknown op accepted")
	}
	r = NewRegistry(spec.All())
	r.AddOptions("GetTeam", Option{Flag: "--team-id", Kind: String})
	if err := r.Finalize(); err == nil {
		t.Error("duplicate option accepted")
	}
	r = NewRegistry(spec.All())
	r.AddOptions("GetTeam", Option{Flag: "--output", Kind: String})
	if err := r.Finalize(); err == nil {
		t.Error("global collision accepted")
	}
}

// D-C32: several registrations merge; examples append, validators and builders chain.
func TestOverrideMerge(t *testing.T) {
	r := NewRegistry(spec.All())
	var order []string
	r.Override(Override{Op: "CreateTeam", Examples: []Example{{Title: "a", Command: "afs teams create-team"}},
		Validate: func(*Parsed) error { order = append(order, "v1"); return nil },
		Build:    func(*Parsed, *RequestParts) error { order = append(order, "b1"); return nil },
		Rename:   map[string]string{"slug": "team-slug"}})
	r.Override(Override{Op: "CreateTeam", Examples: []Example{{Title: "b", Command: "afs teams create-team --kind ephemeral"}},
		Validate: func(*Parsed) error { order = append(order, "v2"); return nil },
		Build:    func(*Parsed, *RequestParts) error { order = append(order, "b2"); return nil },
		Columns:  []output.Column{{Header: "Id", Path: "id"}}})
	if err := r.Finalize(); err != nil {
		t.Fatal(err)
	}
	c := r.ByOp("CreateTeam")
	if len(c.Examples) != 2 || c.Columns == nil {
		t.Fatal("merge lost data")
	}
	_ = c.Validate(nil)
	_ = c.Build(nil, nil)
	if strings.Join(order, ",") != "v1,v2,b1,b2" {
		t.Errorf("order = %v", order)
	}
}

func TestRegistryQueries(t *testing.T) {
	r := testRegistry(t)
	if len(r.Services()) < 10 || len(r.Commands("teams")) == 0 || len(r.All()) != 180 {
		t.Errorf("services %d, teams %d, all %d", len(r.Services()), len(r.Commands("teams")), len(r.All()))
	}
	if _, ok := r.Lookup("teams", "nope"); ok {
		t.Error("lookup found nothing")
	}
}

// D-C4: every operation has a milestone (Appendix A).
func TestEveryOperationHasMilestone(t *testing.T) {
	for _, o := range spec.Commands() {
		m := MilestoneOf(o.ID)
		if m < 0 || m > 6 {
			t.Errorf("%s: milestone %d", o.ID, m)
		}
		if StreamOf(o.ID) == "" {
			t.Errorf("%s: no stream", o.ID)
		}
	}
	if len(milestoneOf) != 180 {
		t.Errorf("milestone map has %d entries", len(milestoneOf))
	}
	if MilestoneOf("Nope") != -1 {
		t.Error("unknown op has a milestone")
	}
}

func TestPendingListsOperationsWithoutExamples(t *testing.T) {
	r := testRegistry(t) // no service overrides: every operation up to CurrentMilestone is pending
	want := 0
	for _, c := range r.All() {
		if c.Op != nil && MilestoneOf(c.Op.ID) <= CurrentMilestone {
			want++
		}
	}
	if p := r.Pending(); len(p) != want {
		t.Errorf("pending = %d operations %v, want %d", len(p), p, want)
	}
}
