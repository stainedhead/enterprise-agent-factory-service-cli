package gen_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/client/gen"
)

func TestCreateTeamRequestRoundTrip(t *testing.T) {
	in := `{"displayName":"Payments","slug":"payments","ownerEmail":"lead@example.com","kind":"persistent"}`
	var r gen.CreateTeamRequest
	if err := json.Unmarshal([]byte(in), &r); err != nil {
		t.Fatal(err)
	}
	if r.DisplayName != "Payments" {
		t.Fatalf("displayName = %q", r.DisplayName)
	}
	out, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), `"displayName":"Payments"`) || !strings.Contains(string(out), `"slug":"payments"`) {
		t.Fatalf("round trip lost fields: %s", out)
	}
}

func TestPutBudgetMoneyIsString(t *testing.T) {
	r := gen.PutBudgetRequest{AmountUsd: "100.50", Period: "monthly"}
	out, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), `"amountUsd":"100.50"`) {
		t.Fatalf("money is not a string: %s", out)
	}
}
