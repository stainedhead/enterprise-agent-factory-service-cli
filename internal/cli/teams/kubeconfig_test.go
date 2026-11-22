package teams_test

import (
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/clitest"
)

const apiKubeconfig = "apiVersion: v1\nkind: Config\nclusters:\n  - name: afs\n    cluster:\n      server: https://EXAMPLE.gr7.us-east-1.eks.amazonaws.com\ncontexts:\n  - name: afs-payments\n    context: {cluster: afs, namespace: afs-payments, user: afs}\ncurrent-context: afs-payments\nusers:\n  - name: afs\n    user:\n      exec:\n        apiVersion: client.authentication.k8s.io/v1beta1\n        command: aws\n        args: [eks, get-token, --cluster-name, afs]\n"

func yamlReply(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/yaml")
	w.Header().Set("X-Request-Id", "req_k")
	_, _ = w.Write([]byte(apiKubeconfig))
}

func haveAWS(string) (string, error) { return "/usr/local/bin/aws", nil }

func TestCommandsDoneKubeconfig(t *testing.T) {
	clitest.CheckOperations(t, "kubeconfig", "GetTeamKubeconfig")
}

// CLI-US-11.1: with no --kubeconfig and no $KUBECONFIG the context is merged into ~/.kube/config,
// 0600, with no token and the aws exec plugin.
func TestGetKubeconfigMergeDefaultPath(t *testing.T) {
	home := t.TempDir()
	r := clitest.RunWith(t, clitest.Options{Home: home, LookPath: haveAWS}, yamlReply, "teams", "get-kubeconfig", "--team-id", "team_1")
	if r.Code != 0 || r.Req.Header.Get("Accept") != "application/yaml" {
		t.Fatalf("%+v", r)
	}
	path := filepath.Join(home, ".kube", "config")
	b, err := os.ReadFile(path) //nolint:gosec // test file
	if err != nil {
		t.Fatal(err)
	}
	got := string(b)
	if !strings.Contains(got, "afs-payments") || !strings.Contains(got, "command: aws") || strings.Contains(got, "token:") {
		t.Errorf("kubeconfig:\n%s", got)
	}
	if st, _ := os.Stat(path); st.Mode().Perm() != 0o600 {
		t.Errorf("mode %v", st.Mode().Perm())
	}
	if !strings.Contains(r.Out, "afs-payments") || r.Err != "" {
		t.Errorf("out %q err %q", r.Out, r.Err)
	}
	// $KUBECONFIG: the first path wins over ~/.kube/config
	first := filepath.Join(t.TempDir(), "first")
	r = clitest.RunWith(t, clitest.Options{Home: home, LookPath: haveAWS, Vars: map[string]string{"KUBECONFIG": first + string(os.PathListSeparator) + "/nonexistent/second"}},
		yamlReply, "teams", "get-kubeconfig", "--team-id", "team_1")
	if _, err := os.Stat(first); r.Code != 0 || err != nil {
		t.Fatalf("$KUBECONFIG: %+v %v", r, err)
	}
}

// CLI-US-11.2: --kubeconfig PATH writes there only, --aws-profile adds the env, a second run keeps one context.
func TestGetKubeconfigExplicitPath(t *testing.T) {
	home := t.TempDir()
	kc := filepath.Join(t.TempDir(), "kc.yaml")
	for i := 0; i < 2; i++ {
		r := clitest.RunWith(t, clitest.Options{Home: home, LookPath: haveAWS}, yamlReply,
			"teams", "get-kubeconfig", "--team-id", "team_1", "--kubeconfig", kc, "--aws-profile", "ops")
		if r.Code != 0 {
			t.Fatalf("run %d: %+v", i, r)
		}
	}
	b, _ := os.ReadFile(kc) //nolint:gosec // test file
	if strings.Count(string(b), "name: afs-payments") != 1 || !strings.Contains(string(b), "AWS_PROFILE") || !strings.Contains(string(b), "value: ops") {
		t.Errorf("kubeconfig:\n%s", b)
	}
	if _, err := os.Stat(filepath.Join(home, ".kube", "config")); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("~/.kube/config was touched: %v", err)
	}
	// a wide existing file keeps its mode and warns
	wide := filepath.Join(t.TempDir(), "wide")
	if err := os.WriteFile(wide, []byte("apiVersion: v1\nkind: Config\n"), 0o644); err != nil { //nolint:gosec // the test needs a wide file
		t.Fatal(err)
	}
	if err := os.Chmod(wide, 0o644); err != nil { //nolint:gosec // the test needs a wide file
		t.Fatal(err)
	}
	r := clitest.RunWith(t, clitest.Options{LookPath: haveAWS}, yamlReply, "teams", "get-kubeconfig", "--team-id", "team_1", "--kubeconfig", wide)
	if st, _ := os.Stat(wide); r.Code != 0 || st.Mode().Perm() != 0o644 || !strings.Contains(r.Err, "wider than 0600") {
		t.Fatalf("%+v mode %v", r, st.Mode().Perm())
	}
}

// CLI-US-11.2: --kubeconfig - prints the API's YAML unchanged and touches no file.
func TestGetKubeconfigStdout(t *testing.T) {
	home := t.TempDir()
	r := clitest.RunWith(t, clitest.Options{Home: home}, yamlReply, "teams", "get-kubeconfig", "--team-id", "team_1", "--kubeconfig", "-")
	if r.Code != 0 || r.Out != apiKubeconfig {
		t.Fatalf("%+v", r)
	}
	if _, err := os.Stat(filepath.Join(home, ".kube")); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("a file was written: %v", err)
	}
}

// CLI-US-11.3: a viewer with kubectl switched off gets the API's AccessDenied, 254.
func TestGetKubeconfigAccessDenied(t *testing.T) {
	r := clitest.Run(t, clitest.Denied("teams:get-kubeconfig"), "teams", "get-kubeconfig", "--team-id", "team_1", "--kubeconfig", "-")
	if r.Code != 254 || !strings.Contains(r.Err, "(AccessDenied) when calling the GetTeamKubeconfig operation") {
		t.Fatalf("%+v", r)
	}
}

// CLI-16: a missing aws CLI is a warning, not an error.
func TestMissingAWSWarns(t *testing.T) {
	r := clitest.RunWith(t, clitest.Options{LookPath: func(string) (string, error) { return "", errors.New("not found") }}, yamlReply,
		"teams", "get-kubeconfig", "--team-id", "team_1", "--kubeconfig", filepath.Join(t.TempDir(), "kc"))
	if r.Code != 0 || !strings.Contains(r.Err, "aws CLI") {
		t.Fatalf("%+v", r)
	}
}
