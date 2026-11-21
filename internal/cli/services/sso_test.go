package services_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli/services"
	clisso "github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli/sso"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/sso"
)

type res struct {
	code     int
	out, err string
}

func afs(t *testing.T, home string, vars map[string]string, stdin string, args ...string) res {
	t.Helper()
	reg, err := services.NewRegistry()
	if err != nil {
		t.Fatal(err)
	}
	var out, eb bytes.Buffer
	env := cli.NewEnv(strings.NewReader(stdin), &out, &eb)
	env.Vars = vars
	env.Home = home
	env.Sleep = func(context.Context, time.Duration) error { return nil }
	code := cli.Run(context.Background(), args, env, reg)
	return res{code, out.String(), eb.String()}
}

func fakeAPI(t *testing.T, providers string) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/v1/identity-providers":
			_, _ = w.Write([]byte(providers))
		case "/v1/me":
			_, _ = w.Write([]byte(`{"principal":{"id":"usr_1","type":"user","email":"a@b.c"},"serviceAdmin":false,"securityOperator":false,"teams":[]}`))
		default:
			w.WriteHeader(404)
			_, _ = w.Write([]byte(`{"code":"NotFound","message":"x","requestId":"r"}`))
		}
	}))
	t.Cleanup(srv.Close)
	return srv
}

const oneProvider = `{"items":[{"id":"idp_1","kind":"okta","name":"Okta","issuerUrl":"https://idp.example.com","cliClientId":"cli1","consoleClientId":"c","isDefault":true}]}`

func TestGetCallerIdentityJSONUnchanged(t *testing.T) {
	srv := fakeAPI(t, oneProvider)
	r := afs(t, t.TempDir(), map[string]string{"AFS_ENDPOINT_URL": srv.URL, "AFS_API_KEY": "k"}, "", "sso", "get-caller-identity")
	if r.code != 0 || !strings.Contains(r.out, `"email": "a@b.c"`) {
		t.Fatalf("%+v", r)
	}
	r = afs(t, t.TempDir(), map[string]string{"AFS_ENDPOINT_URL": srv.URL, "AFS_API_KEY": "k"}, "", "sso", "get-caller-identity", "--output", "table")
	if !strings.Contains(r.out, "usr_1") {
		t.Fatalf("%+v", r)
	}
}

func TestListProvidersNeedsNoCredential(t *testing.T) {
	srv := fakeAPI(t, oneProvider)
	r := afs(t, t.TempDir(), map[string]string{"AFS_ENDPOINT_URL": srv.URL}, "", "sso", "list-providers", "--output", "table")
	if r.code != 0 || !strings.Contains(r.out, "Okta") {
		t.Fatalf("%+v", r)
	}
}

func TestConfigureWritesProfile0600(t *testing.T) {
	home := t.TempDir()
	r := afs(t, home, map[string]string{}, "https://afs.example.com\ntable\nteam_1\n", "configure", "--profile", "prod")
	if r.code != 0 {
		t.Fatalf("%+v", r)
	}
	p := filepath.Join(home, ".afs", "config")
	b, err := os.ReadFile(p) //nolint:gosec // test file
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"[profile prod]", "endpoint_url = https://afs.example.com", "output = table", "team_id = team_1"} {
		if !strings.Contains(string(b), want) {
			t.Errorf("missing %q in\n%s", want, b)
		}
	}
	if fi, _ := os.Stat(p); fi.Mode().Perm() != 0o600 {
		t.Errorf("mode %v", fi.Mode())
	}
	if r := afs(t, home, map[string]string{}, "http://insecure.example.com\n\n\n", "configure"); r.code != 255 {
		t.Errorf("http endpoint: %+v", r)
	}
	if r := afs(t, home, map[string]string{}, "\nbogus\n\n", "configure"); r.code != 252 {
		t.Errorf("bad output: %+v", r)
	}
	if r := afs(t, home, map[string]string{}, "", "configure", "nope"); r.code != 252 {
		t.Errorf("bad arg: %+v", r)
	}
}

func TestConfigureSSOListsProvidersAndWritesKeys(t *testing.T) {
	srv := fakeAPI(t, `{"items":[{"id":"idp_1","kind":"okta","name":"A","issuerUrl":"https://a","cliClientId":"ca"},{"id":"idp_2","kind":"entra","name":"B","issuerUrl":"https://b","cliClientId":"cb","isDefault":true}]}`)
	home := t.TempDir()
	r := afs(t, home, map[string]string{"AFS_ENDPOINT_URL": srv.URL}, "2\n", "configure", "sso", "--profile", "p")
	if r.code != 0 || !strings.Contains(r.err, "1) A") {
		t.Fatalf("%+v", r)
	}
	b, _ := os.ReadFile(filepath.Join(home, ".afs", "config")) //nolint:gosec // test file
	for _, want := range []string{"sso_provider_id = idp_2", "sso_issuer = https://b", "sso_client_id = cb"} {
		if !strings.Contains(string(b), want) {
			t.Errorf("missing %q in\n%s", want, b)
		}
	}
	if r := afs(t, home, map[string]string{"AFS_ENDPOINT_URL": srv.URL}, "9\n", "configure", "sso"); r.code != 252 {
		t.Errorf("bad choice: %+v", r)
	}
	if r := afs(t, home, map[string]string{"AFS_ENDPOINT_URL": srv.URL}, "\n", "configure", "sso", "--profile", "q"); r.code != 0 {
		t.Errorf("default choice: %+v", r)
	}
	one := fakeAPI(t, oneProvider)
	if r := afs(t, home, map[string]string{"AFS_ENDPOINT_URL": one.URL}, "", "configure", "sso", "--profile", "r"); r.code != 0 {
		t.Errorf("single provider: %+v", r)
	}
	none := fakeAPI(t, `{"items":[]}`)
	if r := afs(t, home, map[string]string{"AFS_ENDPOINT_URL": none.URL}, "", "configure", "sso"); r.code != 255 {
		t.Errorf("no providers: %+v", r)
	}
}

func TestLogoutDeletesCacheThenNextCommandExit255(t *testing.T) {
	srv := fakeAPI(t, oneProvider)
	home := t.TempDir()
	cfg := "[default]\nendpoint_url = " + srv.URL + "\nsso_issuer = https://idp.example.com\nsso_client_id = cli1\n"
	if err := os.MkdirAll(filepath.Join(home, ".afs"), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(home, ".afs", "config"), []byte(cfg), 0o600); err != nil {
		t.Fatal(err)
	}
	cache := sso.NewCache(home)
	if err := cache.Save(&sso.Token{Issuer: "https://idp.example.com", ClientID: "cli1", AccessToken: "tok", ExpiresAt: time.Now().Add(time.Hour)}); err != nil {
		t.Fatal(err)
	}
	vars := map[string]string{}
	if r := afs(t, home, vars, "", "sso", "get-caller-identity"); r.code != 0 {
		t.Fatalf("with session: %+v", r)
	}
	if r := afs(t, home, vars, "", "sso", "logout"); r.code != 0 || !strings.Contains(r.out, "Signed out") {
		t.Fatalf("logout: %+v", r)
	}
	r := afs(t, home, vars, "", "sso", "get-caller-identity")
	if r.code != 255 || !strings.Contains(r.err, "afs sso login --profile default") {
		t.Fatalf("after logout: %+v", r)
	}
	if r := afs(t, t.TempDir(), vars, "", "sso", "logout"); r.code != 0 || !strings.Contains(r.out, "no single sign-on") {
		t.Fatalf("no sso: %+v", r)
	}
}

// A fake server that is both the API and the identity provider (device flow).
func idpServer(t *testing.T) *httptest.Server {
	t.Helper()
	var srv *httptest.Server
	srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/v1/identity-providers":
			_, _ = w.Write([]byte(`{"items":[{"id":"idp_1","kind":"okta","name":"Okta","issuerUrl":"` + srv.URL + `","cliClientId":"cli1","isDefault":true}]}`))
		case "/.well-known/openid-configuration":
			_, _ = w.Write([]byte(`{"device_authorization_endpoint":"` + srv.URL + `/device","token_endpoint":"` + srv.URL + `/token"}`))
		case "/device":
			_, _ = w.Write([]byte(`{"device_code":"dc","user_code":"ABCD-1234","verification_uri":"https://v.example.com","expires_in":600,"interval":1}`))
		case "/token":
			_, _ = w.Write([]byte(`{"access_token":"access-xyz","refresh_token":"refresh-xyz","expires_in":3600,"token_type":"Bearer"}`))
		case "/v1/me":
			if r.Header.Get("Authorization") != "Bearer access-xyz" {
				w.WriteHeader(401)
				_, _ = w.Write([]byte(`{"code":"Unauthenticated","message":"no","requestId":"r"}`))
				return
			}
			_, _ = w.Write([]byte(`{"principal":{"id":"usr_1","type":"user"},"serviceAdmin":false,"securityOperator":false,"teams":[]}`))
		}
	}))
	t.Cleanup(srv.Close)
	return srv
}

func TestSSOLoginThenCommandUsesCachedToken(t *testing.T) {
	srv := idpServer(t)
	home := t.TempDir()
	vars := map[string]string{"AFS_ENDPOINT_URL": srv.URL}
	r := afs(t, home, vars, "", "sso", "login", "--profile", "dev")
	if r.code != 0 || !strings.Contains(r.err, "ABCD-1234") || !strings.Contains(r.out, "Signed in") {
		t.Fatalf("login: %+v", r)
	}
	if strings.Contains(r.out+r.err, "access-xyz") || strings.Contains(r.out+r.err, "refresh-xyz") {
		t.Fatal("a token reached the terminal")
	}
	cfg, _ := os.ReadFile(filepath.Join(home, ".afs", "config")) //nolint:gosec // test file
	if !strings.Contains(string(cfg), "sso_client_id = cli1") {
		t.Fatalf("login did not write the profile keys:\n%s", cfg)
	}
	r = afs(t, home, vars, "", "sso", "get-caller-identity", "--profile", "dev")
	if r.code != 0 || !strings.Contains(r.out, "usr_1") {
		t.Fatalf("after login: %+v", r)
	}
	r = afs(t, home, vars, "", "sso", "login", "--profile", "dev") // second login: keys already present
	if r.code != 0 {
		t.Fatalf("relogin: %+v", r)
	}
}

func TestSSOLoginFailures(t *testing.T) {
	none := fakeAPI(t, `{"items":[]}`)
	if r := afs(t, t.TempDir(), map[string]string{"AFS_ENDPOINT_URL": none.URL}, "", "sso", "login"); r.code != 255 {
		t.Fatalf("no providers: %+v", r)
	}
	two := fakeAPI(t, `{"items":[{"id":"a","issuerUrl":"https://a"},{"id":"b","issuerUrl":"https://b"}]}`)
	if r := afs(t, t.TempDir(), map[string]string{"AFS_ENDPOINT_URL": two.URL}, "", "sso", "login"); r.code != 255 || !strings.Contains(r.err, "configure sso") {
		t.Fatalf("ambiguous: %+v", r)
	}
	if r := afs(t, t.TempDir(), map[string]string{}, "", "sso", "login"); r.code != 255 {
		t.Fatalf("no endpoint: %+v", r)
	}
}

func TestPick(t *testing.T) {
	ps := []clisso.Provider{{ID: "a"}, {ID: "b", IsDefault: true}}
	if p, err := clisso.Pick(ps, "a"); err != nil || p.ID != "a" {
		t.Error("by id")
	}
	if p, err := clisso.Pick(ps, ""); err != nil || p.ID != "b" {
		t.Error("default")
	}
	if _, err := clisso.Pick(ps, "zzz"); err == nil {
		t.Error("unknown id")
	}
	if _, err := clisso.Pick(nil, ""); err == nil {
		t.Error("none")
	}
}
