package auth

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func tokenFile(t *testing.T, content string) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "token")
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	return p
}

func workloadInputs(t *testing.T, env map[string]string, creds string) Inputs {
	t.Helper()
	in := setup(t, env, creds, 0o600)
	in.Workload = WorkloadToken{Getenv: func(k string) string { return env[k] }, DefaultPath: filepath.Join(t.TempDir(), "absent")}
	return in
}

// CLI-5, CLI-US-10.1, D-C29: with no other credential the workload token file is the credential.
func TestWorkloadTokenUsedWhenNoOtherCredential(t *testing.T) {
	path := tokenFile(t, "wl-token-1\n")
	in := workloadInputs(t, map[string]string{"AFS_API_TOKEN_FILE": path}, "")
	c, err := Chain(in)(context.Background())
	if err != nil || c.Header != "Bearer wl-token-1" || c.Source != "workload token file" {
		t.Fatalf("%+v %v", c, err)
	}
	// the default path is used when the file exists
	def := tokenFile(t, "wl-default")
	in.Workload = WorkloadToken{Getenv: func(string) string { return "" }, DefaultPath: def}
	if c, err = Chain(in)(context.Background()); err != nil || c.Header != "Bearer wl-default" {
		t.Fatalf("default path: %+v %v", c, err)
	}
	// and none is offered when it does not
	in.Workload = WorkloadToken{Getenv: func(string) string { return "" }, DefaultPath: filepath.Join(t.TempDir(), "absent")}
	if c, err = Chain(in)(context.Background()); err != nil || c.Header != "" {
		t.Fatalf("no file: %+v %v", c, err)
	}
	// a file named by AFS_API_TOKEN_FILE that is missing is an error: the user asked for it
	in = workloadInputs(t, map[string]string{"AFS_API_TOKEN_FILE": filepath.Join(t.TempDir(), "missing")}, "")
	if _, err = Chain(in)(context.Background()); err == nil || !strings.Contains(err.Error(), "workload token file") {
		t.Fatalf("explicit missing file: %v", err)
	}
	// an empty file is no credential
	in = workloadInputs(t, map[string]string{"AFS_API_TOKEN_FILE": tokenFile(t, "  \n")}, "")
	if c, err = Chain(in)(context.Background()); err != nil || c.Header != "" {
		t.Fatalf("empty file: %+v %v", c, err)
	}
}

// CLI-5: the kubelet rotates the token, so it is read again for every request.
func TestWorkloadTokenReReadPerRequest(t *testing.T) {
	path := tokenFile(t, "first")
	chain := Chain(workloadInputs(t, map[string]string{"AFS_API_TOKEN_FILE": path}, ""))
	if c, _ := chain(context.Background()); c.Header != "Bearer first" {
		t.Fatalf("%+v", c)
	}
	if err := os.WriteFile(path, []byte("second"), 0o600); err != nil {
		t.Fatal(err)
	}
	if c, _ := chain(context.Background()); c.Header != "Bearer second" {
		t.Fatalf("not re-read: %+v", c)
	}
}

// CLI-5: the workload token is last: the env key, the credentials file and a cached SSO token all win.
func TestWorkloadTokenLastInChain(t *testing.T) {
	path := tokenFile(t, "wl")
	env := map[string]string{"AFS_API_TOKEN_FILE": path}

	in := workloadInputs(t, map[string]string{"AFS_API_TOKEN_FILE": path, "AFS_API_KEY": "afsk_env"}, "")
	if c, _ := Chain(in)(context.Background()); c.Header != "Bearer afsk_env" {
		t.Errorf("env key: %+v", c)
	}
	in = workloadInputs(t, env, "[default]\napi_key = afsk_file\n")
	if c, _ := Chain(in)(context.Background()); c.Header != "Bearer afsk_file" {
		t.Errorf("credentials file: %+v", c)
	}
	in = workloadInputs(t, env, "")
	in.SSO = fakeSSO{tok: "sso-token", ok: true}
	if c, _ := Chain(in)(context.Background()); c.Header != "Bearer sso-token" {
		t.Errorf("sso: %+v", c)
	}
	// an expired SSO session does not hide the workload token
	in.SSO = fakeSSO{err: &LoginRequiredError{Profile: "default"}}
	if c, err := Chain(in)(context.Background()); err != nil || c.Header != "Bearer wl" {
		t.Errorf("expired sso: %+v %v", c, err)
	}
}
