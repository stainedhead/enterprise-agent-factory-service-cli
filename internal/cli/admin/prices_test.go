package admin_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/clitest"
)

func TestCommandsDonePrices(t *testing.T) {
	clitest.CheckOperations(t, "prices", "ListModelPrices", "CreateModelPrice")
}

// CLI-1: money options are strings, sent as given (no float conversion).
func TestCreateModelPriceMoneyAsStrings(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) { clitest.Reply(w, 201, `{"modelId":"m1"}`) },
		"admin", "create-model-price", "--model-id", "m1", "--effective-from", "2026-12-01T00:00:00Z",
		"--input-usd-per-mtok", "3.00", "--output-usd-per-mtok", "15.00", "--cache-read-usd-per-mtok", "0.30")
	body := string(r.ReqBody)
	if r.Code != 0 || r.Req.URL.Path != "/v1/admin/models/m1/prices" || !strings.Contains(body, `"inputUsdPerMtok":"3.00"`) || !strings.Contains(body, `"outputUsdPerMtok":"15.00"`) || !strings.Contains(body, `"cacheReadUsdPerMtok":"0.30"`) {
		t.Fatalf("%+v body=%s", r, body)
	}
}
