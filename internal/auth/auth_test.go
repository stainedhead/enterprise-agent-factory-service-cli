package auth

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/config"
)

type fakeSSO struct {
	tok string
	ok  bool
	err error
}

func (f fakeSSO) Token(context.Context) (string, bool, error) { return f.tok, f.ok, f.err }

func setup(t *testing.T, env map[string]string, creds string, mode os.FileMode) Inputs {
	t.Helper()
	home := t.TempDir()
	if creds != "" {
		if err := os.MkdirAll(filepath.Join(home, ".afs"), 0o700); err != nil {
			t.Fatal(err)
		}
		p := filepath.Join(home, ".afs", "credentials")
		if err := os.WriteFile(p, []byte(creds), 0o600); err != nil {
			t.Fatal(err)
		}
		if err := os.Chmod(p, mode); err != nil { //nolint:gosec // test
			t.Fatal(err)
		}
	}
	src := config.Source{Getenv: func(k string) string { return env[k] }, Home: home}
	st, err := config.Resolve(config.Overrides{}, src)
	if err != nil {
		t.Fatal(err)
	}
	return Inputs{Source: src, Settings: st}
}

// CLI-5: AFS_API_KEY wins even when an SSO token is cached.
func TestEnvKeyWinsOverProfileSSO(t *testing.T) {
	in := setup(t, map[string]string{"AFS_API_KEY": "afsk_env"}, "[default]\napi_key = afsk_file\n", 0o600)
	in.SSO = fakeSSO{tok: "sso-token", ok: true}
	c, err := Chain(in)(context.Background())
	if err != nil || c.Header != "Bearer afsk_env" || c.Source != "env AFS_API_KEY" {
		t.Fatalf("cred %+v err %v", c, err)
	}
}

// 5.1-credentials: a credentials file wider than 0600 is refused with a chmod hint.
func TestCredentialsFile0644Refused(t *testing.T) {
	in := setup(t, map[string]string{}, "[default]\napi_key = afsk_file\n", 0o644)
	_, err := Chain(in)(context.Background())
	if err == nil || !strings.Contains(err.Error(), "chmod 600") {
		t.Fatalf("err = %v", err)
	}
}

func TestChainOrderAndSourceNames(t *testing.T) {
	in := setup(t, map[string]string{}, "[default]\napi_key = afsk_file\n", 0o600)
	in.SSO = fakeSSO{tok: "sso-token", ok: true}
	c, err := Chain(in)(context.Background())
	if err != nil || c.Header != "Bearer afsk_file" || !strings.Contains(c.Source, "credentials file") || !strings.Contains(c.Source, "profile default") {
		t.Fatalf("file cred %+v err %v", c, err)
	}
	if strings.Contains(c.Source, "afsk_file") {
		t.Error("source names the key")
	}
	in = setup(t, map[string]string{}, "", 0o600)
	in.SSO = fakeSSO{tok: "sso-token", ok: true}
	in.Workload = fakeSSO{tok: "wl", ok: true}
	if c, _ = Chain(in)(context.Background()); c.Header != "Bearer sso-token" || c.Source != "sso cache profile default" {
		t.Errorf("sso cred %+v", c)
	}
	in.SSO = fakeSSO{}
	if c, _ = Chain(in)(context.Background()); c.Header != "Bearer wl" || c.Source != "workload token file" {
		t.Errorf("workload cred %+v", c)
	}
	in.Workload = nil
	if c, err = Chain(in)(context.Background()); err != nil || c.Header != "" {
		t.Errorf("empty chain: %+v %v", c, err)
	}
}

func TestSSOLoginRequiredAndOtherErrors(t *testing.T) {
	in := setup(t, map[string]string{}, "", 0o600)
	in.SSO = fakeSSO{err: &LoginRequiredError{Profile: "default", Reason: "Your session expired"}}
	_, err := Chain(in)(context.Background())
	var lr *LoginRequiredError
	if !errors.As(err, &lr) || !strings.Contains(err.Error(), "Run afs sso login --profile default") {
		t.Fatalf("err = %v", err)
	}
	// a workload token still works after a missing SSO session
	in.Workload = fakeSSO{tok: "wl", ok: true}
	if c, err := Chain(in)(context.Background()); err != nil || c.Header != "Bearer wl" {
		t.Errorf("fallback: %+v %v", c, err)
	}
	in.SSO = fakeSSO{err: errors.New("disk on fire")}
	if _, err := Chain(in)(context.Background()); err == nil || err.Error() != "disk on fire" {
		t.Errorf("other error: %v", err)
	}
	in.SSO = nil
	in.Workload = fakeSSO{err: errors.New("bad token file")}
	if _, err := Chain(in)(context.Background()); err == nil {
		t.Error("workload error swallowed")
	}
	if (&LoginRequiredError{Profile: "p"}).Error() != "You are not signed in. Run afs sso login --profile p" {
		t.Error("default reason text")
	}
}
