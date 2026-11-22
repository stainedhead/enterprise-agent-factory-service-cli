package services_test

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/clitest"
)

// CLI-US-10.1: a worker with only the projected token is a workload principal; get-caller-identity
// sends it as the bearer token and prints the API's answer.
func TestGetCallerIdentityWorkload(t *testing.T) {
	tok := filepath.Join(t.TempDir(), "token")
	if err := os.WriteFile(tok, []byte("projected-sa-token\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	r := clitest.RunWith(t, clitest.Options{Vars: map[string]string{"AFS_API_KEY": "", "AFS_API_TOKEN_FILE": tok}}, func(w http.ResponseWriter, req *http.Request) {
		clitest.Reply(w, 200, `{"type":"workload","subject":"system:serviceaccount:afs-pool:worker","roles":["workload"]}`)
	}, "sso", "get-caller-identity")
	if r.Code != 0 || r.Req.Header.Get("Authorization") != "Bearer projected-sa-token" || !strings.Contains(r.Out, `"type": "workload"`) {
		t.Fatalf("%+v auth=%q", r, r.Req.Header.Get("Authorization"))
	}
	if strings.Contains(r.Err, "projected-sa-token") {
		t.Error("the token reached stderr")
	}
}
