package teams_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/clitest"
)

func TestCommandsDoneLogDestinations(t *testing.T) {
	clitest.CheckOperations(t, "log-destinations", "ListLogDestinations", "CreateLogDestination", "GetLogDestination", "UpdateLogDestination", "DeleteLogDestination")
}

// CLI-US-19.2: describe prints the configuration and delivery status, never a credential.
func TestDescribeLogDestinationNoCredential(t *testing.T) {
	body := `{"id":"log_1","name":"team-splunk","kind":"splunk_hec","enabled":true,"config":{"endpoint":"https://splunk.example.com:8088"},"secretId":"sec_1","version":2,"generation":2,"status":{"phase":"Degraded","observedGeneration":2,"conditions":[{"type":"Delivered","status":"False","reason":"Rejected","message":"HTTP 403 from the collector"}]}}`
	r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) { clitest.Reply(w, 200, body) },
		"teams", "describe-log-destination", "--team-id", "team_1", "--destination-id", "log_1", "--output", "table")
	if r.Code != 0 || !strings.Contains(r.Out, "splunk.example.com:8088") || !strings.Contains(r.Out, "Degraded") || !strings.Contains(r.Out, "sec_1") {
		t.Fatalf("%+v", r)
	}
	if r.Req.URL.Path != "/v1/teams/team_1/log-destinations/log_1" {
		t.Errorf("path %s", r.Req.URL.Path)
	}
}

// CLI-12: create-log-destination --wait polls the destination until Ready.
func TestCreateLogDestinationWait(t *testing.T) {
	polls := 0
	r := clitest.Run(t, func(w http.ResponseWriter, req *http.Request) {
		if req.Method == "POST" {
			clitest.Reply(w, 201, `{"id":"log_1","generation":1,"status":{"phase":"Pending","observedGeneration":0}}`)
			return
		}
		polls++
		phase := "Pending"
		if polls >= 2 {
			phase = "Ready"
		}
		clitest.Reply(w, 200, `{"id":"log_1","name":"n","kind":"cloudwatch","generation":1,"version":1,"status":{"phase":"`+phase+`","observedGeneration":1}}`)
	}, "teams", "create-log-destination", "--team-id", "team_1", "--name", "n", "--kind", "cloudwatch", "--config", `{"logGroup":"/afs/x","region":"us-east-1"}`, "--wait")
	if r.Code != 0 || polls < 2 || !strings.Contains(r.Out, "Ready") {
		t.Fatalf("%+v polls=%d", r, polls)
	}
}
