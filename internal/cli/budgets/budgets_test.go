package budgets_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/clitest"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/golden"
)

func TestCommandsDoneBudgets(t *testing.T) {
	clitest.CheckOperations(t, "budgets", "GetBudget", "PutBudget")
}

var put = []string{"budgets", "put-budget", "--team-id", "team_1", "--period", "monthly", "--amount-usd", "100.50", "--thresholds", "50", "80", "100"}

func budgetHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		clitest.Reply(w, 200, `{"period":"monthly","amountUsd":"100.50","version":4}`)
		return
	}
	clitest.Reply(w, 200, `{"period":"monthly","amountUsd":"100.50","version":5}`)
}

// CLI-US-14.1: the amount is a decimal string and the thresholds are numbers.
func TestPutBudgetBody(t *testing.T) {
	r := clitest.Run(t, budgetHandler, put...)
	if r.Code != 0 || r.Req.Method != "PUT" || string(r.ReqBody) != `{"period":"monthly","amountUsd":"100.50","thresholds":[50,80,100]}` && !strings.Contains(string(r.ReqBody), `"amountUsd":"100.50"`) {
		t.Fatalf("%+v body=%s", r, r.ReqBody)
	}
	if !strings.Contains(string(r.ReqBody), `"thresholds":[50,80,100]`) {
		t.Errorf("thresholds: %s", r.ReqBody)
	}
}

// CLI-US-14.2: an amount above the team's ceiling is the API's ValidationError (254).
func TestPutBudgetAboveCeiling254(t *testing.T) {
	r := clitest.Run(t, clitest.Fail(422, "ValidationError", "amountUsd 99999.00 is above the team ceiling 5000.00"), put...)
	if r.Code != 254 || !strings.Contains(r.Err, "(ValidationError) when calling the PutBudget operation: amountUsd 99999.00 is above the team ceiling 5000.00") {
		t.Fatalf("%+v", r)
	}
}

// CLI-US-14.3: the table shows amount, spent, reserved, percent used and forecast.
func TestGetBudgetTableGolden(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) {
		clitest.Reply(w, 200, `{"id":"bud_1","period":"monthly","amountUsd":"2000.00","spentUsd":"812.40","reservedUsd":"3.20","percentUsed":40.6,"forecastUsd":"1950.00","onExhaust":"block","thresholds":[50,80,100],"version":1}`)
	}, "budgets", "get-budget", "--team-id", "team_1", "--output", "table")
	if r.Code != 0 {
		t.Fatalf("%+v", r)
	}
	for _, want := range []string{"Amount", "2000.00", "Spent", "812.40", "Reserved", "3.20", "PercentUsed", "40.6", "Forecast", "1950.00"} {
		if !strings.Contains(r.Out, want) {
			t.Errorf("table lacks %q:\n%s", want, r.Out)
		}
	}
	golden.Check(t, "budgets/get-budget-fixture.table.golden", r.Out)
}

// CLI-US-6.3: put-budget sends If-Match: the given version, or the one it read.
func TestPutBudgetIfMatch(t *testing.T) {
	r := clitest.Run(t, budgetHandler, append(append([]string{}, put...), "--if-match", "3")...)
	if r.Code != 0 || len(r.Requests) != 1 || r.Req.Header.Get("If-Match") != `"3"` {
		t.Fatalf("explicit: %+v", r)
	}
	r = clitest.Run(t, budgetHandler, put...)
	if r.Code != 0 || r.Req.Header.Get("If-Match") != `"4"` || !strings.Contains(r.Err, "version 4") {
		t.Fatalf("read: %+v", r)
	}
}
