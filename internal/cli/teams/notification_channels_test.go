package teams_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/clitest"
)

func TestCommandsDoneNotificationChannels(t *testing.T) {
	clitest.CheckOperations(t, "notification-channels", "ListNotificationChannels", "GetNotificationChannel", "CreateNotificationChannel", "UpdateNotificationChannel", "DeleteNotificationChannel", "TestNotificationChannel")
}

// CLI-US-19.2: describe prints the stored configuration and delivery status, with the secret's id
// only; the path parameter is --channel-id.
func TestDescribeChannelShowsConfigAndDelivery(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) {
		clitest.Reply(w, 200, `{"id":"nch_1","name":"payments-slack","kind":"slack","config":{"channel":"#alerts"},"secretId":"sec_9","delivery":{"healthy":false,"lastError":"channel_not_found"},"version":2}`)
	}, "teams", "describe-notification-channel", "--team-id", "team_1", "--channel-id", "nch_1")
	if r.Code != 0 || r.Req.URL.Path != "/v1/teams/team_1/notification-channels/nch_1" {
		t.Fatalf("%+v", r)
	}
	for _, want := range []string{`"channel": "#alerts"`, `"lastError": "channel_not_found"`, `"secretId": "sec_9"`} {
		if !strings.Contains(r.Out, want) {
			t.Errorf("missing %s in %s", want, r.Out)
		}
	}
	tbl := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) {
		clitest.Reply(w, 200, `{"id":"nch_1","name":"n","kind":"slack","delivery":{"healthy":false,"lastError":"channel_not_found"},"version":2}`)
	}, "teams", "describe-notification-channel", "--team-id", "team_1", "--channel-id", "nch_1", "--output", "table")
	if !strings.Contains(tbl.Out, "channel_not_found") {
		t.Errorf("table lacks the delivery error:\n%s", tbl.Out)
	}
}

// CLI-12: --wait on create prints the final channel once Ready.
func TestCreateChannelWait(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, req *http.Request) {
		if req.Method == "POST" {
			clitest.Reply(w, 201, `{"id":"nch_1","generation":1}`)
			return
		}
		clitest.Reply(w, 200, `{"id":"nch_1","name":"payments-slack","generation":1,"status":{"phase":"Ready","observedGeneration":1}}`)
	}, "teams", "create-notification-channel", "--team-id", "team_1", "--name", "payments-slack", "--kind", "slack",
		"--config", `{"channel":"#a"}`, "--event-types", "budget.*", "--wait")
	if r.Code != 0 || !strings.Contains(r.Out, "payments-slack") || !strings.Contains(string(r.RequestBody[0]), `"eventTypes":["budget.*"]`) {
		t.Fatalf("%+v body=%s", r, r.ReqBody)
	}
}

// CLI-10: update and delete read the version from the channel itself.
func TestUpdateAndDeleteChannelReadVersion(t *testing.T) {
	h := func(w http.ResponseWriter, req *http.Request) {
		if req.Method == "GET" {
			clitest.Reply(w, 200, `{"id":"nch_1","version":6}`)
			return
		}
		clitest.Reply(w, 202, `{"id":"nch_1","version":7}`)
	}
	r := clitest.Run(t, h, "teams", "update-notification-channel", "--team-id", "team_1", "--channel-id", "nch_1", "--no-enabled")
	if r.Code != 0 || r.Req.Header.Get("If-Match") != `"6"` || string(r.ReqBody) != `{"enabled":false}` {
		t.Fatalf("update: %+v body=%s", r, r.ReqBody)
	}
	r = clitest.Run(t, h, "teams", "delete-notification-channel", "--team-id", "team_1", "--channel-id", "nch_1")
	if r.Code != 0 || r.Req.Method != "DELETE" || r.Req.Header.Get("If-Match") != `"6"` {
		t.Fatalf("delete: %+v", r)
	}
}

func testChannel(status int, body string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) { clitest.Reply(w, status, body) }
}

var testArgs = []string{"teams", "test-notification-channel", "--team-id", "team_1", "--channel-id", "nch_1"}

// CLI-US-19.1, D-C20: delivered is exit 0 and prints the result.
func TestNotificationTestDelivered0(t *testing.T) {
	r := clitest.Run(t, testChannel(200, `{"delivered":true,"pending":false,"testId":"ntest_1"}`), testArgs...)
	if r.Code != 0 || !strings.Contains(r.Out, `"delivered": true`) || r.Err != "" || r.Req.URL.Path != "/v1/teams/team_1/notification-channels/nch_1:test" {
		t.Fatalf("%+v", r)
	}
}

// CLI-US-19.1, D-C20: a failed delivery prints the result, then the AWS-style DeliveryFailed
// error of the TestNotificationChannel operation, and exits 254.
func TestNotificationTestFailed254(t *testing.T) {
	r := clitest.Run(t, testChannel(200, `{"delivered":false,"error":"channel_not_found","pending":false,"testId":"ntest_1"}`), testArgs...)
	if r.Code != 254 || !strings.Contains(r.Out, "channel_not_found") {
		t.Fatalf("%+v", r)
	}
	if !strings.Contains(r.Err, "An error occurred (DeliveryFailed) when calling the TestNotificationChannel operation: channel_not_found") || !strings.Contains(r.Err, "Request ID: req_test") {
		t.Errorf("stderr %q", r.Err)
	}
}

// D-C20: 202 pending is exit 0 with a pointer to describe-notification-channel.
func TestNotificationTestPending202Exit0(t *testing.T) {
	r := clitest.Run(t, testChannel(202, `{"pending":true,"testId":"ntest_1"}`), testArgs...)
	if r.Code != 0 || !strings.Contains(r.Err, "The test is still running; run afs teams describe-notification-channel to see the result.") || !strings.Contains(r.Out, "ntest_1") {
		t.Fatalf("%+v", r)
	}
}
