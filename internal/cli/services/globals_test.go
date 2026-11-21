package services_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli/services"
)

type nopWC struct{ io.Writer }

func (nopWC) Close() error { return nil }

func meServer(t *testing.T, delay time.Duration) (*httptest.Server, *[]http.Header) {
	t.Helper()
	var hdrs []http.Header
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hdrs = append(hdrs, r.Header.Clone())
		time.Sleep(delay)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"principal":{"id":"usr_1","type":"user"},"serviceAdmin":false,"securityOperator":false,"teams":[]}`))
	}))
	t.Cleanup(srv.Close)
	return srv, &hdrs
}

// One case per global option of 5.1.
func TestGlobalOptions(t *testing.T) {
	srv, hdrs := meServer(t, 0)
	home := t.TempDir()
	if err := os.MkdirAll(filepath.Join(home, ".afs"), 0o700); err != nil {
		t.Fatal(err)
	}
	cfg := "[profile p2]\nendpoint_url = " + srv.URL + "\noutput = yaml\n"
	if err := os.WriteFile(filepath.Join(home, ".afs", "config"), []byte(cfg), 0o600); err != nil {
		t.Fatal(err)
	}
	key := map[string]string{"AFS_API_KEY": "afsk_secret_value_123"}
	me := []string{"sso", "get-caller-identity"}
	with := func(extra ...string) []string { return append(append([]string{}, me...), extra...) }

	t.Run("profile", func(t *testing.T) {
		r := afs(t, home, key, "", with("--profile", "p2")...)
		if r.code != 0 || !strings.Contains(r.out, "serviceAdmin: false") {
			t.Fatalf("%+v", r)
		}
	})
	t.Run("endpoint-url", func(t *testing.T) {
		r := afs(t, home, key, "", with("--endpoint-url", srv.URL)...)
		if r.code != 0 || !strings.Contains(r.out, `"serviceAdmin"`) {
			t.Fatalf("%+v", r)
		}
	})
	t.Run("output and query", func(t *testing.T) {
		r := afs(t, home, map[string]string{"AFS_API_KEY": "k", "AFS_ENDPOINT_URL": srv.URL}, "", with("--output", "text", "--query", "principal.id")...)
		if strings.TrimSpace(r.out) != "usr_1" {
			t.Fatalf("%+v", r)
		}
	})
	t.Run("missing profile", func(t *testing.T) {
		if r := afs(t, home, key, "", with("--profile", "nope")...); r.code != 255 {
			t.Fatalf("%+v", r)
		}
	})
	t.Run("debug redacts the credential", func(t *testing.T) {
		r := afs(t, home, key, "", with("--profile", "p2", "--debug")...)
		if r.code != 0 || r.err == "" || strings.Contains(r.err, "afsk_secret_value_123") || strings.Contains(r.out, "afsk_secret_value_123") {
			t.Fatalf("%+v", r)
		}
	})
	t.Run("color on and off", func(t *testing.T) {
		r := afs(t, home, key, "", with("--profile", "p2", "--output", "table", "--color", "on")...)
		if !strings.Contains(r.out, "\x1b[") {
			t.Fatalf("no color: %q", r.out)
		}
		r = afs(t, home, key, "", with("--profile", "p2", "--output", "table", "--color", "off")...)
		if strings.Contains(r.out, "\x1b[") {
			t.Fatalf("color: %q", r.out)
		}
	})
	t.Run("no-cli-pager", func(t *testing.T) {
		for _, tc := range []struct {
			extra   []string
			started bool
		}{{nil, true}, {[]string{"--no-cli-pager"}, false}} {
			reg, _ := services.NewRegistry()
			var out, eb bytes.Buffer
			env := cli.NewEnv(strings.NewReader(""), &out, &eb)
			env.Vars = map[string]string{"AFS_API_KEY": "k", "AFS_ENDPOINT_URL": srv.URL, "AFS_PAGER": "fakepager"}
			env.Home = home
			env.StdoutTTY = true
			started := false
			env.StartPager = func(cmd string, o io.Writer) (io.WriteCloser, func() error, error) {
				started = true
				return nopWC{o}, func() error { return nil }, nil
			}
			if code := cli.Run(context.Background(), with(tc.extra...), env, reg); code != 0 || started != tc.started {
				t.Fatalf("code=%d started=%v want %v err=%s", code, started, tc.started, eb.String())
			}
		}
	})
	t.Run("generate-cli-skeleton makes no call", func(t *testing.T) {
		n := len(*hdrs)
		r := afs(t, home, key, "", "teams", "create-team", "--generate-cli-skeleton")
		if r.code != 0 || len(*hdrs) != n {
			t.Fatalf("%+v", r)
		}
	})
	t.Run("version", func(t *testing.T) {
		r := afs(t, home, nil, "", "--version")
		if r.code != 0 || !strings.HasPrefix(r.out, "afs-cli/") {
			t.Fatalf("%+v", r)
		}
	})
	t.Run("user agent and authorization", func(t *testing.T) {
		h := (*hdrs)[0]
		if !strings.HasPrefix(h.Get("User-Agent"), "afs-cli/") || h.Get("Authorization") != "Bearer afsk_secret_value_123" {
			t.Fatalf("%v", h)
		}
	})
}

func TestReadTimeoutExit255(t *testing.T) {
	srv, _ := meServer(t, 1500*time.Millisecond)
	r := afs(t, t.TempDir(), map[string]string{"AFS_API_KEY": "k", "AFS_ENDPOINT_URL": srv.URL}, "", "sso", "get-caller-identity", "--cli-read-timeout", "1", "--cli-connect-timeout", "2")
	if r.code != 255 || !strings.Contains(r.err, "timeout") {
		t.Fatalf("%+v", r)
	}
}
