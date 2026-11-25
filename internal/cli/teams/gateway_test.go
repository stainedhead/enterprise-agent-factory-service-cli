package teams_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/clitest"
)

// FR-010, MEG-CLI-2: the two team commands of v1.5.0.
func TestCommandsDoneGateway(t *testing.T) {
	clitest.CheckOperations(t, "gateway", "PutTeamGatewaySettings", "PutTeamModelEndpoints")
}

// FR-010: put-team-model-endpoints sends the list, reads no version, and the gateway settings send their booleans.
func TestPutTeamModelEndpointsBody(t *testing.T) {
	ok := func(w http.ResponseWriter, _ *http.Request) {
		clitest.Reply(w, 200, `{"teamId":"team_1","endpointIds":["mep_1"],"selections":[{"endpointId":"mep_1","slug":"openai-prod","selectedBy":"team_admin"}]}`)
	}
	r := clitest.Run(t, ok, "teams", "put-team-model-endpoints", "--team-id", "team_1", "--endpoint-ids", "mep_1", "mep_2")
	if r.Code != 0 || len(r.Requests) != 1 || r.Req.Method != "PUT" || r.Req.URL.Path != "/v1/teams/team_1/model-endpoint-selections" || r.Req.Header.Get("If-Match") != "" {
		t.Fatalf("%+v", r)
	}
	var b struct {
		EndpointIDs []string `json:"endpointIds"`
	}
	if err := json.Unmarshal(r.ReqBody, &b); err != nil || len(b.EndpointIDs) != 2 || b.EndpointIDs[1] != "mep_2" {
		t.Errorf("body %s", r.ReqBody)
	}
	g := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) {
		clitest.Reply(w, 200, `{"teamId":"team_1","deployment":"per_namespace","mayAddModelEndpoints":false,"gatewayState":"provisioning","gatewayUrl":null,"version":7}`)
	}, "teams", "put-team-gateway-settings", "--team-id", "team_1", "--deployment", "per_namespace", "--no-may-add-model-endpoints")
	var gb map[string]any
	_ = json.Unmarshal(g.ReqBody, &gb)
	if g.Code != 0 || g.Req.URL.Path != "/v1/admin/teams/team_1/gateway-settings" || gb["deployment"] != "per_namespace" || gb["mayAddModelEndpoints"] != false {
		t.Errorf("%+v %v", g, gb)
	}
	if r := clitest.Run(t, nil, "teams", "put-team-gateway-settings", "--team-id", "team_1", "--deployment", "dedicated"); r.Code != 252 {
		t.Errorf("bad enum: %+v", r)
	}
}
