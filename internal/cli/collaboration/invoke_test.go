package collaboration_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/clitest"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/sso"
)

const (
	pcKey      = "pckey-SENTINEL-7f3a9c"
	afsKey     = "afsk_unit_test_key" // the clitest default AFS_API_KEY
	queryToken = "qtok-SENTINEL-51d2"
)

// fakeEdge is a Paperclip edge on a loopback address; it counts what reaches it.
type fakeEdge struct {
	*httptest.Server
	n    atomic.Int32
	last atomic.Pointer[http.Request]
	body atomic.Pointer[[]byte]
}

func newEdge(t *testing.T, h http.HandlerFunc) *fakeEdge {
	t.Helper()
	e := &fakeEdge{}
	e.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		e.n.Add(1)
		b := readAll(r)
		e.last.Store(r)
		e.body.Store(&b)
		h(w, r)
	}))
	t.Cleanup(e.Close)
	return e
}

func readAll(r *http.Request) []byte {
	var b []byte
	buf := make([]byte, 4096)
	for {
		n, err := r.Body.Read(buf)
		b = append(b, buf[:n]...)
		if err != nil {
			return b
		}
	}
}

// api answers GetTeamCollaboration with restApiUrl set to url (null when url is empty).
func api(url string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := "null"
		if url != "" {
			u = `"` + url + `"`
		}
		clitest.Reply(w, 200, strings.Replace(collab("t.collab.example", true, true, "[]", 1, 1, 1, "Ready"), `"restApiUrl":null`, `"restApiUrl":`+u, 1))
	}
}

func keyFile(t *testing.T, content string, mode os.FileMode) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "pc.key")
	if err := os.WriteFile(p, []byte(content), mode); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(p, mode); err != nil {
		t.Fatal(err)
	}
	return p
}

func ok(body string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) { clitest.Reply(w, 200, body) }
}

func invokeArgs(extra ...string) []string {
	return append([]string{"collaboration", "invoke-api", "--team-id", team}, extra...)
}

// CLI-US-25.1
func TestInvokeAPIRequestShape(t *testing.T) {
	e := newEdge(t, ok(`{"id":"iss_1"}`))
	kf := keyFile(t, pcKey+"\n", 0o600)
	bf := filepath.Join(t.TempDir(), "issue.json")
	_ = os.WriteFile(bf, []byte(`{"title":"x"}`), 0o600)
	r := clitest.Run(t, api(e.URL+"/api"), invokeArgs("--method", "post", "--path", "/api/issues", "--body", "file://"+bf, "--paperclip-key-file", kf)...)
	if r.Code != 0 || len(r.Requests) != 1 || r.Req.URL.Path != "/v1/teams/"+team+"/collaboration" || e.n.Load() != 1 {
		t.Fatalf("%+v edge=%d", r, e.n.Load())
	}
	q, body := e.last.Load(), *e.body.Load()
	h := q.Header
	if q.Method != "POST" || q.URL.Path != "/api/issues" || string(body) != `{"title":"x"}` ||
		h.Get("X-Afs-Authorization") != "Bearer "+afsKey || h.Get("Authorization") != "Bearer "+pcKey ||
		h.Get("Content-Type") != "application/json" || h.Get("Accept") != "application/json" ||
		!strings.HasPrefix(h.Get("User-Agent"), "afs-cli/") || h.Get("Cookie") != "" {
		t.Errorf("%s %s %v", q.Method, q.URL, h)
	}
}

func TestInvokeAPIResponseBodyUnchanged(t *testing.T) {
	const doc = "{ \"a\" :  [1,2],\n \"b\":\"é\" }\n"
	e := newEdge(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write([]byte(doc))
	})
	r := clitest.Run(t, api(e.URL), invokeArgs("--path", "/api/x", "--paperclip-key-file", keyFile(t, pcKey, 0o600))...)
	if r.Code != 0 || r.Out != doc {
		t.Errorf("%d %q", r.Code, r.Out)
	}
	r = clitest.Run(t, api(e.URL), invokeArgs("--path", "/api/x", "--include", "--paperclip-key-file", keyFile(t, pcKey, 0o600))...)
	if !strings.HasPrefix(r.Out, "HTTP/1.1 200 OK\n") || !strings.HasSuffix(r.Out, "\n\n"+doc) {
		t.Errorf("--include: %q", r.Out)
	}
}

// CLI-US-25.2
func TestInvokeAPICredentialChain(t *testing.T) {
	e := newEdge(t, ok(`{}`))
	kf := keyFile(t, pcKey, 0o600)
	r := clitest.RunWith(t, clitest.Options{Vars: map[string]string{"AFS_API_KEY": "afsk_other"}}, api(e.URL), invokeArgs("--path", "/api/x", "--paperclip-key-file", kf, "--debug")...)
	if got := e.last.Load().Header.Get("X-Afs-Authorization"); got != "Bearer afsk_other" {
		t.Errorf("X-Afs-Authorization %q", got)
	}
	if !strings.Contains(r.Err, "credential source: env AFS_API_KEY") || strings.Contains(r.Err, "afsk_other") {
		t.Errorf("debug must name the source and not the key: %s", r.Err)
	}
	// an SSO profile: the cached access token is sent
	home := t.TempDir()
	if err := sso.NewCache(home).Save(&sso.Token{Issuer: "https://idp.example.com", ClientID: "c1", AccessToken: "sso-token-1", ExpiresAt: time.Now().Add(time.Hour), TokenType: "Bearer"}); err != nil {
		t.Fatal(err)
	}
	cfg := filepath.Join(home, ".afs")
	_ = os.MkdirAll(cfg, 0o700)
	_ = os.WriteFile(filepath.Join(cfg, "config"), []byte("[default]\nsso_issuer = https://idp.example.com\nsso_client_id = c1\n"), 0o600)
	r = clitest.RunWith(t, clitest.Options{Home: home, Vars: map[string]string{"AFS_API_KEY": ""}}, api(e.URL), invokeArgs("--path", "/api/x", "--paperclip-key-file", kf)...)
	if got := e.last.Load().Header.Get("X-Afs-Authorization"); r.Code != 0 || got != "Bearer sso-token-1" {
		t.Errorf("%d %q %s", r.Code, got, r.Err)
	}
}

func TestInvokeAPIWorkloadTokenRefused255(t *testing.T) {
	e := newEdge(t, ok(`{}`))
	tok := filepath.Join(t.TempDir(), "token")
	_ = os.WriteFile(tok, []byte("wl-token-SENTINEL"), 0o600)
	r := clitest.RunWith(t, clitest.Options{Vars: map[string]string{"AFS_API_KEY": "", "AFS_API_TOKEN_FILE": tok}}, api(e.URL),
		invokeArgs("--path", "/api/x", "--paperclip-key-file", keyFile(t, pcKey, 0o600))...)
	if r.Code != 255 || e.n.Load() != 0 || len(r.Requests) != 0 || !strings.Contains(r.Err, "does not admit workload credentials; use an API key or afs sso login.") || strings.Contains(r.Err, "wl-token-SENTINEL") {
		t.Errorf("%d %q edge=%d api=%d", r.Code, r.Err, e.n.Load(), len(r.Requests))
	}
}

func TestInvokeAPIExpiredToken255(t *testing.T) {
	e := newEdge(t, ok(`{}`))
	home := t.TempDir()
	_ = sso.NewCache(home).Save(&sso.Token{Issuer: "https://idp.example.com", ClientID: "c1", AccessToken: "old", ExpiresAt: time.Now().Add(-time.Hour), TokenType: "Bearer"})
	cfg := filepath.Join(home, ".afs")
	_ = os.MkdirAll(cfg, 0o700)
	_ = os.WriteFile(filepath.Join(cfg, "config"), []byte("[default]\nsso_issuer = https://idp.example.com\nsso_client_id = c1\n"), 0o600)
	r := clitest.RunWith(t, clitest.Options{Home: home, Vars: map[string]string{"AFS_API_KEY": ""}}, api(e.URL), invokeArgs("--path", "/api/x", "--paperclip-key-file", keyFile(t, pcKey, 0o600))...)
	if r.Code != 255 || !strings.Contains(r.Err, "afs sso login") || e.n.Load() != 0 || len(r.Requests) != 0 {
		t.Errorf("%d %q", r.Code, r.Err)
	}
}

// CLI-US-25.3
func TestInvokeAPIKeyFromStdin(t *testing.T) {
	e := newEdge(t, ok(`{}`))
	r := clitest.RunWith(t, clitest.Options{Stdin: pcKey + " \n\n"}, api(e.URL), invokeArgs("--path", "/api/x", "--paperclip-key-file", "-")...)
	if r.Code != 0 || e.last.Load().Header.Get("Authorization") != "Bearer "+pcKey {
		t.Errorf("%d %q", r.Code, r.Err)
	}
	r = clitest.RunWith(t, clitest.Options{Stdin: "x"}, api(e.URL), invokeArgs("--path", "/api/x", "--method", "POST", "--body", "-", "--paperclip-key-file", "-")...)
	if r.Code != 252 {
		t.Errorf("--body - must exit 252, got %d", r.Code)
	}
}

func TestInvokeAPIKeyFileMode255(t *testing.T) {
	e := newEdge(t, ok(`{}`))
	kf := keyFile(t, pcKey, 0o644)
	r := clitest.Run(t, api(e.URL), invokeArgs("--path", "/api/x", "--paperclip-key-file", kf)...)
	if r.Code != 255 || !strings.Contains(r.Err, "chmod 600 "+kf) || e.n.Load() != 0 || strings.Contains(r.Err, pcKey) {
		t.Errorf("%d %q", r.Code, r.Err)
	}
	r = clitest.Run(t, api(e.URL), invokeArgs("--path", "/api/x", "--paperclip-key-file", filepath.Join(t.TempDir(), "missing"))...)
	if r.Code != 255 {
		t.Errorf("missing file: %d", r.Code)
	}
}

func TestInvokeAPIKeyShapeRefused252(t *testing.T) {
	e := newEdge(t, ok(`{}`))
	for name, k := range map[string]string{"empty": "\n", "spaces": "abc def", "control": "abc\x01def", "long": strings.Repeat("k", 4097), "tab": "a\tb"} {
		r := clitest.Run(t, api(e.URL), invokeArgs("--path", "/api/x", "--paperclip-key-file", keyFile(t, k, 0o600))...)
		if r.Code != 252 || e.n.Load() != 0 || len(r.Requests) != 0 || strings.Contains(r.Err, "abc") && name != "empty" {
			t.Errorf("%s: %d %q", name, r.Code, r.Err)
		}
	}
}

func TestInvokeAPIKeyAsArgumentUnknown252(t *testing.T) {
	e := newEdge(t, ok(`{}`))
	for _, a := range [][]string{{"--paperclip-key", pcKey}, {"--paperclip-key=" + pcKey}} {
		r := clitest.Run(t, api(e.URL), invokeArgs(append([]string{"--path", "/api/x"}, a...)...)...)
		if r.Code != 252 || !strings.Contains(r.Err, "--paperclip-key-file") || strings.Contains(r.Err, pcKey) || e.n.Load() != 0 {
			t.Errorf("%v: %d %q", a, r.Code, r.Err)
		}
	}
	r := clitest.Run(t, api(e.URL), invokeArgs("--path", "/api/x")...)
	if r.Code != 252 || !strings.Contains(r.Err, "--paperclip-key-file") {
		t.Errorf("missing key option: %d %q", r.Code, r.Err)
	}
}

// CLI-US-25.4
func TestInvokeAPIValidation252(t *testing.T) {
	e := newEdge(t, ok(`{}`))
	kf := keyFile(t, pcKey, 0o600)
	bf := filepath.Join(t.TempDir(), "b.json")
	_ = os.WriteFile(bf, []byte(`{}`), 0o600)
	big := filepath.Join(t.TempDir(), "big")
	f, _ := os.Create(big) //nolint:gosec // test file in a temp dir
	_ = f.Truncate(100<<20 + 1)
	_ = f.Close()
	cases := map[string][]string{
		"method":      {"--method", "TRACE", "--path", "/api/x"},
		"no path":     {},
		"relative":    {"--path", "api/x"},
		"scheme":      {"--path", "/api/x?u=https://evil.example"},
		"host":        {"--path", "//evil.example/x"},
		"dotdot":      {"--path", "/api/../x"},
		"encoded dot": {"--path", "/api/%2e%2E/x"},
		"control":     {"--path", "/api/x\n"},
		"space":       {"--path", "/api/x y"},
		"fragment":    {"--path", "/api/x#f"},
		"body get":    {"--path", "/api/x", "--body", "file://" + bf},
		"body head":   {"--method", "head", "--path", "/api/x", "--body", "file://" + bf},
		"body scheme": {"--method", "POST", "--path", "/api/x", "--body", bf},
		"body big":    {"--method", "POST", "--path", "/api/x", "--body", "file://" + big},
		"ctype":       {"--method", "POST", "--path", "/api/x", "--body", "file://" + bf, "--content-type", "a\r\nX: y"},
		"positional":  {"--path", "/api/x", "stray"},
	}
	for name, a := range cases {
		r := clitest.Run(t, api(e.URL), invokeArgs(append(a, "--paperclip-key-file", kf)...)...)
		if r.Code != 252 || e.n.Load() != 0 || len(r.Requests) != 0 {
			t.Errorf("%s: code %d edge=%d api=%d %q", name, r.Code, e.n.Load(), len(r.Requests), r.Err)
		}
	}
	r := clitest.Run(t, api(e.URL), invokeArgs("--method", "POST", "--path", "/api/x", "--body", "file://"+filepath.Join(t.TempDir(), "gone"), "--paperclip-key-file", kf)...)
	if r.Code != 255 {
		t.Errorf("unreadable body file is a local problem: %d", r.Code)
	}
}

func TestInvokeAPIReservedPathRefused252(t *testing.T) {
	e := newEdge(t, ok(`{}`))
	for _, p := range []string{"/_afs/health", "/_afs", "/_AFS/x", "/_afs/a?b=c"} {
		r := clitest.Run(t, api(e.URL), invokeArgs("--path", p, "--paperclip-key-file", keyFile(t, pcKey, 0o600))...)
		if r.Code != 252 || e.n.Load() != 0 {
			t.Errorf("%s: %d", p, r.Code)
		}
	}
	r := clitest.Run(t, api(e.URL), invokeArgs("--path", "/api/_afs/x", "--paperclip-key-file", keyFile(t, pcKey, 0o600))...)
	if r.Code != 0 {
		t.Errorf("a path that only contains _afs is Paperclip's: %d", r.Code)
	}
}

func TestInvokeAPIBodyTooLargeFromEdge254(t *testing.T) {
	e := newEdge(t, func(w http.ResponseWriter, _ *http.Request) {
		clitest.Reply(w, 400, `{"code":"ValidationError","message":"too big","requestId":"req_b","details":{"reason":"body_too_large"}}`)
	})
	r := clitest.Run(t, api(e.URL), invokeArgs("--path", "/api/x", "--paperclip-key-file", keyFile(t, pcKey, 0o600))...)
	if r.Code != 254 || !strings.Contains(r.Err, "(ValidationError) when calling the InvokeApi operation") || !strings.Contains(r.Err, "collab.max_body_mb") || r.Out != "" {
		t.Errorf("%d %q", r.Code, r.Err)
	}
}

// CLI-US-25.5
func TestInvokeAPIOriginFromRestApiURL(t *testing.T) {
	e := newEdge(t, ok(`{}`))
	// only the scheme and the host are used: the URL's own path (/api) is not
	r := clitest.Run(t, api(e.URL+"/api"), invokeArgs("--path", "/api/issues?status=open", "--paperclip-key-file", keyFile(t, pcKey, 0o600))...)
	q := e.last.Load()
	if r.Code != 0 || q.URL.Path != "/api/issues" || q.URL.RawQuery != "status=open" {
		t.Errorf("%d %s", r.Code, q.URL)
	}
}

func TestInvokeAPICollabURLSkipsRead(t *testing.T) {
	e := newEdge(t, ok(`{}`))
	r := clitest.RunWith(t, clitest.Options{Vars: map[string]string{"AFS_ENDPOINT_URL": ""}}, clitest.Fail(500, "InternalError", "no"),
		"collaboration", "invoke-api", "--collab-url", e.URL, "--path", "/api/x", "--paperclip-key-file", keyFile(t, pcKey, 0o600))
	if r.Code != 0 || len(r.Requests) != 0 || e.n.Load() != 1 {
		t.Errorf("%d api=%d edge=%d %q", r.Code, len(r.Requests), e.n.Load(), r.Err)
	}
}

func TestInvokeAPIOriginRefused255(t *testing.T) {
	e := newEdge(t, ok(`{}`))
	kf := keyFile(t, pcKey, 0o600)
	cases := map[string]string{ //nolint:gosec // G101: a URL with a made-up password is the case under test
		"http":       "http://t.collab.example/api",
		"userinfo":   "https://user:pw@t.collab.example/api",
		"port":       "https://t.collab.example:8443/api",
		"ipv4":       "https://10.0.0.1/api",
		"ipv6":       "https://[::1]/api",
		"query":      "https://t.collab.example/?q=1",
		"fragment":   "https://t.collab.example/#f",
		"empty host": "https:///api",
		"ftp":        "ftp://t.collab.example/",
	}
	for name, u := range cases {
		r := clitest.Run(t, api(u), invokeArgs("--path", "/api/x", "--paperclip-key-file", kf)...)
		if r.Code != 255 || e.n.Load() != 0 || strings.Contains(r.Err, pcKey) || strings.Contains(r.Err, "pw@") {
			t.Errorf("%s: %d %q", name, r.Code, r.Err)
		}
		if name != "empty host" && !strings.Contains(r.Err, "refused") {
			t.Errorf("%s: the message must name the rule: %q", name, r.Err)
		}
	}
	// collab_base_domain pins the host to one label under the domain
	home := t.TempDir()
	_ = os.MkdirAll(filepath.Join(home, ".afs"), 0o700)
	_ = os.WriteFile(filepath.Join(home, ".afs", "config"), []byte("[default]\ncollab_base_domain = collab.example\n"), 0o600)
	for u, want := range map[string]int{"https://t.other.example/api": 255, "https://a.b.collab.example/api": 255, "https://collab.example/api": 255, "https://evilcollab.example/api": 255} {
		r := clitest.RunWith(t, clitest.Options{Home: home}, api(u), invokeArgs("--path", "/api/x", "--paperclip-key-file", kf)...)
		if r.Code != want || e.n.Load() != 0 {
			t.Errorf("%s: %d %q", u, r.Code, r.Err)
		}
	}
}

func TestInvokeAPINullRestApiURL255(t *testing.T) {
	e := newEdge(t, ok(`{}`))
	r := clitest.Run(t, api(""), invokeArgs("--path", "/api/x", "--paperclip-key-file", keyFile(t, pcKey, 0o600))...)
	if r.Code != 255 || !strings.Contains(r.Err, "Public collaboration is not configured for this service (collab.base_domain is unset).") || e.n.Load() != 0 {
		t.Errorf("%d %q", r.Code, r.Err)
	}
}

func TestInvokeAPINoRedirectFollowed(t *testing.T) {
	var target atomic.Int32
	t2 := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { target.Add(1) }))
	defer t2.Close()
	e := newEdge(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Location", t2.URL+"/x")
		w.WriteHeader(302)
	})
	r := clitest.Run(t, api(e.URL), invokeArgs("--path", "/api/x", "--include", "--paperclip-key-file", keyFile(t, pcKey, 0o600))...)
	if r.Code != 0 || target.Load() != 0 || !strings.Contains(r.Out, "302 Found") || !strings.Contains(r.Out, "Location: "+t2.URL+"/x") {
		t.Errorf("%d %q target=%d", r.Code, r.Out, target.Load())
	}
}

func TestInvokeAPINoCookieJar(t *testing.T) {
	e := newEdge(t, func(w http.ResponseWriter, _ *http.Request) {
		http.SetCookie(w, &http.Cookie{Name: "__Host-afs-collab", Value: "cookie-SENTINEL", Path: "/", Secure: true, HttpOnly: true, SameSite: http.SameSiteLaxMode})
		clitest.Reply(w, 200, `{}`)
	})
	r := clitest.Run(t, api(e.URL), invokeArgs("--path", "/api/x", "--include", "--paperclip-key-file", keyFile(t, pcKey, 0o600))...)
	if strings.Contains(r.Out, "cookie-SENTINEL") || !strings.Contains(r.Out, "Set-Cookie: [redacted]") {
		t.Errorf("Set-Cookie must be masked: %q", r.Out)
	}
	// a second run sends no Cookie (nothing is stored between runs, nor within one)
	_ = clitest.Run(t, api(e.URL), invokeArgs("--path", "/api/x", "--paperclip-key-file", keyFile(t, pcKey, 0o600))...)
	if e.last.Load().Header.Get("Cookie") != "" {
		t.Error("a cookie was sent")
	}
}

// CLI-US-25.6
func TestInvokeAPIEdgeErrors254(t *testing.T) {
	type tc struct {
		status       int
		code, reason string
		extra        string
		wantAll      []string
	}
	cases := []tc{
		{403, "AccessDenied", "surface_disabled", "", []string{"The REST API is switched off for this team. A team-admin can turn it on: afs collaboration put-collaboration --team-id <id> --rest-api-enabled"}},
		{403, "AccessDenied", "permission_required", `,"requiredPermission":"collaboration:invoke"`, []string{"Required permission: collaboration:invoke"}},
		{403, "AccessDenied", "no_team_role", "", []string{"Your credential has no role on this team"}},
		{401, "Unauthenticated", "token_expired", "", []string{"Run afs sso login to sign in again."}},
		{401, "Unauthenticated", "token_invalid", "", []string{"Run afs sso login"}},
		{401, "Unauthenticated", "key_expired", "", []string{"afs api-keys create-api-key"}},
		{401, "Unauthenticated", "key_revoked", "", []string{"afs api-keys create-api-key"}},
		{400, "ValidationError", "ambiguous_credentials", "", []string{"Send one AFS credential only."}},
		{400, "ValidationError", "malformed_credential", "", []string{"malformed"}},
		{400, "ValidationError", "unsupported_encoding", "", []string{"send the body uncompressed"}},
		{404, "NotFound", "", "", []string{"The team's Paperclip is not reachable through the public endpoint; run afs collaboration describe-collaboration to see why."}},
		{429, "Throttled", "", "", []string{"(Throttled)"}},
		{503, "Unavailable", "paused", "", []string{"is paused"}},
		{503, "Unavailable", "team_suspended", "", []string{"The team is suspended"}},
		{503, "Unavailable", "starting", "", []string{"is starting"}},
		{503, "Unavailable", "upstream_unavailable", "", []string{"did not answer"}},
	}
	for _, c := range cases {
		details := ""
		if c.reason != "" || c.extra != "" {
			details = `,"details":{"reason":"` + c.reason + `"` + c.extra + `}`
		}
		body := fmt.Sprintf(`{"code":%q,"message":"edge says no","requestId":"req_edge1"%s}`, c.code, details)
		e := newEdge(t, func(w http.ResponseWriter, _ *http.Request) { clitest.Reply(w, c.status, body) })
		method := "POST" // no retry of an answered request
		r := clitest.Run(t, api(e.URL), invokeArgs("--method", method, "--path", "/api/x", "--body", "file://"+keyFile(t, "{}", 0o600), "--paperclip-key-file", keyFile(t, pcKey, 0o600))...)
		if r.Code != 254 || r.Out != "" || !strings.Contains(r.Err, "An error occurred ("+c.code+") when calling the InvokeApi operation: edge says no") || !strings.Contains(r.Err, "Request ID: req_edge1") {
			t.Errorf("%d/%s: %d %q", c.status, c.reason, r.Code, r.Err)
		}
		for _, w := range c.wantAll {
			if !strings.Contains(r.Err, w) {
				t.Errorf("%d/%s: want %q in %q", c.status, c.reason, w, r.Err)
			}
		}
	}
}

func TestInvokeAPISurfaceDisabledHint(t *testing.T) {
	e := newEdge(t, func(w http.ResponseWriter, _ *http.Request) {
		clitest.Reply(w, 403, `{"code":"AccessDenied","message":"off","requestId":"r","details":{"reason":"surface_disabled"}}`)
	})
	r := clitest.Run(t, api(e.URL), invokeArgs("--path", "/api/x", "--paperclip-key-file", keyFile(t, pcKey, 0o600))...)
	if !strings.Contains(r.Err, "afs collaboration put-collaboration --team-id <id> --rest-api-enabled") || strings.Count(r.Err, "switched off") != 1 {
		t.Errorf("%q", r.Err)
	}
}

func TestInvokeAPISecretDetectedPatternsOnly(t *testing.T) {
	e := newEdge(t, func(w http.ResponseWriter, _ *http.Request) {
		clitest.Reply(w, 422, `{"code":"SecretDetected","message":"The body contains a credential","requestId":"r","details":{"patterns":["aws_access_key","github_token"],"matched":"AKIA-SENTINEL-MATCH"}}`)
	})
	r := clitest.Run(t, api(e.URL), invokeArgs("--method", "POST", "--path", "/api/x", "--body", "file://"+keyFile(t, "{}", 0o600), "--paperclip-key-file", keyFile(t, pcKey, 0o600))...)
	if r.Code != 254 || !strings.Contains(r.Err, "Patterns: aws_access_key, github_token\n") ||
		!strings.Contains(r.Err, "Store the credential as an AFS secret instead (afs secrets create-secret).") || strings.Contains(r.Err, "SENTINEL") {
		t.Errorf("%d %q", r.Code, r.Err)
	}
}

// CLI-US-25.7
func TestInvokeAPIPaperclipErrorBodyOnStdout254(t *testing.T) {
	const pbody = `{"error":"Unprocessable","message":"title is required, token=PCBODY-SENTINEL"}`
	e := newEdge(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(422)
		_, _ = w.Write([]byte(pbody))
	})
	r := clitest.Run(t, api(e.URL), invokeArgs("--path", "/api/x", "--paperclip-key-file", keyFile(t, pcKey, 0o600))...)
	if r.Code != 254 || r.Out != pbody || !strings.Contains(r.Err, "An error occurred (HttpError) when calling the InvokeApi operation: HTTP 422 Unprocessable Entity") || strings.Contains(r.Err, "PCBODY") {
		t.Errorf("%d out=%q err=%q", r.Code, r.Out, r.Err)
	}
	// a Paperclip body that has a code member but is not the envelope is still Paperclip's
	e2 := newEdge(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(409)
		_, _ = w.Write([]byte(`{"code":"dup"}`))
	})
	r = clitest.Run(t, api(e2.URL), invokeArgs("--path", "/api/x", "--paperclip-key-file", keyFile(t, pcKey, 0o600))...)
	if r.Code != 254 || r.Out != `{"code":"dup"}` || !strings.Contains(r.Err, "(HttpError)") {
		t.Errorf("%d out=%q err=%q", r.Code, r.Out, r.Err)
	}
}

func TestInvokeAPIWAFBlock254(t *testing.T) {
	e := newEdge(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(403)
		_, _ = w.Write([]byte("<html>Blocked by \x1b]0;evil\x07 WAF</html>"))
	})
	r := clitest.Run(t, api(e.URL), invokeArgs("--path", "/api/x", "--paperclip-key-file", keyFile(t, pcKey, 0o600))...)
	if r.Code != 254 || !strings.Contains(r.Out, "Blocked by") || !strings.Contains(r.Err, "(HttpError)") || !strings.Contains(r.Err, "HTTP 403 Forbidden") || strings.Contains(r.Err, "WAF") {
		t.Errorf("%d out=%q err=%q", r.Code, r.Out, r.Err)
	}
}

func TestInvokeAPI204AndHead(t *testing.T) {
	e := newEdge(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "HEAD" {
			w.Header().Set("Content-Length", "12")
			return
		}
		w.WriteHeader(204)
	})
	for _, m := range []string{"DELETE", "head"} {
		r := clitest.Run(t, api(e.URL), invokeArgs("--method", m, "--path", "/api/x", "--paperclip-key-file", keyFile(t, pcKey, 0o600))...)
		if r.Code != 0 || r.Out != "" {
			t.Errorf("%s: %d %q", m, r.Code, r.Out)
		}
	}
}

// CLI-US-25.8
func TestInvokeAPIRetriesGetOn503(t *testing.T) {
	var n atomic.Int32
	e := newEdge(t, func(w http.ResponseWriter, _ *http.Request) {
		if n.Add(1) < 3 {
			w.Header().Set("Retry-After", "2")
			w.WriteHeader(503)
			return
		}
		clitest.Reply(w, 200, `{"ok":true}`)
	})
	var slept []time.Duration
	r := clitest.RunWith(t, clitest.Options{Sleep: func(_ context.Context, d time.Duration) error { slept = append(slept, d); return nil }}, api(e.URL),
		invokeArgs("--path", "/api/x", "--paperclip-key-file", keyFile(t, pcKey, 0o600))...)
	if r.Code != 0 || r.Out != `{"ok":true}` || e.n.Load() != 3 || len(slept) != 2 || slept[0] < 2*time.Second {
		t.Errorf("%d %q n=%d slept=%v", r.Code, r.Out, e.n.Load(), slept)
	}
	// three 503s in a row: the last answer is the result
	e2 := newEdge(t, func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(503) })
	r = clitest.Run(t, api(e2.URL), invokeArgs("--path", "/api/x", "--paperclip-key-file", keyFile(t, pcKey, 0o600))...)
	if r.Code != 254 || e2.n.Load() != 3 {
		t.Errorf("%d n=%d", r.Code, e2.n.Load())
	}
}

func TestInvokeAPIPostNotRetriedAfterSend(t *testing.T) {
	for _, st := range []int{429, 503} {
		e := newEdge(t, func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(st) })
		r := clitest.Run(t, api(e.URL), invokeArgs("--method", "POST", "--path", "/api/x", "--body", "file://"+keyFile(t, "{}", 0o600), "--paperclip-key-file", keyFile(t, pcKey, 0o600))...)
		if r.Code != 254 || e.n.Load() != 1 {
			t.Errorf("%d: code %d n=%d", st, r.Code, e.n.Load())
		}
	}
	// a connect failure is retried for any method, at most 3 times
	closed := httptest.NewServer(http.NotFoundHandler())
	url := closed.URL
	closed.Close()
	var slept int
	r := clitest.RunWith(t, clitest.Options{Sleep: func(context.Context, time.Duration) error { slept++; return nil }}, api(url),
		invokeArgs("--method", "POST", "--path", "/api/x", "--body", "file://"+keyFile(t, "{}", 0o600), "--paperclip-key-file", keyFile(t, pcKey, 0o600))...)
	if r.Code != 255 || slept != 2 {
		t.Errorf("%d slept=%d %q", r.Code, slept, r.Err)
	}
}

func TestInvokeAPINoIdempotencyKey(t *testing.T) {
	e := newEdge(t, ok(`{}`))
	_ = clitest.Run(t, api(e.URL), invokeArgs("--method", "POST", "--path", "/api/x", "--body", "file://"+keyFile(t, "{}", 0o600), "--paperclip-key-file", keyFile(t, pcKey, 0o600))...)
	if k := e.last.Load().Header.Get("Idempotency-Key"); k != "" {
		t.Errorf("Idempotency-Key %q", k)
	}
}

// CLI-US-25.9
func TestInvokeAPIGlobalOptionsRefused252(t *testing.T) {
	e := newEdge(t, ok(`{}`))
	kf := keyFile(t, pcKey, 0o600)
	for _, g := range [][]string{{"--query", "a"}, {"--output", "table"}, {"--no-paginate"}, {"--max-items", "1"}, {"--page-size", "5"}, {"--starting-token", "x"},
		{"--cli-input-json", "{}"}, {"--generate-cli-skeleton"}, {"--wait"}} {
		r := clitest.Run(t, api(e.URL), invokeArgs(append([]string{"--path", "/api/x", "--paperclip-key-file", kf}, g...)...)...)
		if r.Code != 252 || !strings.Contains(r.Err, "invoke-api prints the response unchanged; pipe it to a JSON tool instead.") || e.n.Load() != 0 || len(r.Requests) != 0 {
			t.Errorf("%v: %d %q", g, r.Code, r.Err)
		}
	}
	r := clitest.Run(t, api(e.URL), invokeArgs("--path", "/api/x", "--paperclip-key-file", kf, "--no-cli-pager", "--cli-read-timeout", "5", "--cli-connect-timeout", "5", "--color", "off")...)
	if r.Code != 0 {
		t.Errorf("usual globals: %d %q", r.Code, r.Err)
	}
}

func TestInvokeAPIReadTimeout255(t *testing.T) {
	e := newEdge(t, func(w http.ResponseWriter, r *http.Request) { <-r.Context().Done() })
	r := clitest.Run(t, api(e.URL), invokeArgs("--path", "/api/x", "--cli-read-timeout", "1", "--paperclip-key-file", keyFile(t, pcKey, 0o600))...)
	if r.Code != 255 || !strings.Contains(r.Err, "read timeout") {
		t.Errorf("%d %q", r.Code, r.Err)
	}
}

func TestInvokeAPISwitchingProtocols255(t *testing.T) {
	e := newEdge(t, func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(101) })
	r := clitest.Run(t, api(e.URL), invokeArgs("--path", "/api/x", "--paperclip-key-file", keyFile(t, pcKey, 0o600))...)
	if r.Code != 255 {
		t.Errorf("%d %q", r.Code, r.Err)
	}
}

// CLI-US-25.10
func TestInvokeAPIDebugRedactsEverything(t *testing.T) {
	e := newEdge(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Set-Cookie", "c=cookie-SENTINEL")
		clitest.Reply(w, 200, `{"x":1}`)
	})
	bf := keyFile(t, `{"body":"BODY-SENTINEL"}`, 0o600)
	r := clitest.Run(t, api(e.URL+"/api"), invokeArgs("--method", "POST", "--path", "/api/x?token="+queryToken, "--body", "file://"+bf, "--paperclip-key-file", keyFile(t, pcKey, 0o600), "--debug")...)
	if r.Code != 0 {
		t.Fatalf("%d %q", r.Code, r.Err)
	}
	for _, s := range []string{afsKey, pcKey, queryToken, "BODY-SENTINEL", "cookie-SENTINEL"} {
		if strings.Contains(r.Err, s) {
			t.Errorf("debug output has %q:\n%s", s, r.Err)
		}
	}
	for _, want := range []string{"DEBUG collab request: POST ", "/api/x attempt=1", "DEBUG collab response: 200", "request-id=req_test", "Authorization: [redacted]", "X-Afs-Authorization: [redacted]", "Set-Cookie: [redacted]", "credential source: env AFS_API_KEY"} {
		if !strings.Contains(r.Err, want) {
			t.Errorf("debug output lacks %q:\n%s", want, r.Err)
		}
	}
	if strings.Contains(r.Out, afsKey) || strings.Contains(r.Out, pcKey) {
		t.Errorf("stdout: %q", r.Out)
	}
}

func TestInvokeAPISecretsNeverInArgsOrStderr(t *testing.T) {
	kf := keyFile(t, pcKey, 0o600)
	// every failure path: no secret in stderr, and the arguments never carry one (the key is a file)
	e := newEdge(t, func(w http.ResponseWriter, _ *http.Request) { clitest.Reply(w, 500, `not json`) })
	runs := [][]string{
		{"--path", "/api/x?token=" + queryToken, "--paperclip-key-file", kf},
		{"--path", "/api/x?token=" + queryToken, "--paperclip-key-file", keyFile(t, pcKey, 0o644)},
		{"--path", "/api/../x?token=" + queryToken, "--paperclip-key-file", kf},
		{"--path", "/api/x", "--paperclip-key=" + pcKey},
		{"--collab-url", "http://evil.example/?token=" + queryToken, "--path", "/api/x", "--paperclip-key-file", kf},
	}
	for _, a := range runs {
		args := invokeArgs(a...)
		r := clitest.Run(t, api(e.URL), args...)
		for _, s := range []string{afsKey, queryToken} {
			if strings.Contains(r.Err, s) && !strings.Contains(strings.Join(a, " "), "--paperclip-key=") {
				t.Errorf("%v: stderr leaks %q: %s", a, s, r.Err)
			}
		}
		if strings.Contains(r.Err, pcKey) {
			t.Errorf("%v: stderr leaks the Paperclip key: %s", a, r.Err)
		}
		if strings.Contains(strings.Join(args, " "), pcKey) && !strings.Contains(strings.Join(a, " "), "--paperclip-key=") {
			t.Errorf("the key is in the process arguments: %v", args)
		}
	}
}

// the sentinels never reach stdout (other than a body the server returned) or stderr on any path of the
// command, including the generated operations of this service.
func TestNoSecretValueInAnyCollaborationOutput(t *testing.T) {
	e := newEdge(t, ok(`{"ok":true}`))
	kf := keyFile(t, pcKey, 0o600)
	paths := [][]string{
		invokeArgs("--path", "/api/x", "--paperclip-key-file", kf, "--debug"),
		invokeArgs("--path", "/api/x", "--paperclip-key-file", kf, "--include"),
		{"collaboration", "describe-collaboration", "--team-id", team, "--debug"},
		{"collaboration", "rotate-collaboration-keys", "--team-id", team, "--debug"},
	}
	for _, a := range paths {
		r := clitest.Run(t, func(w http.ResponseWriter, rq *http.Request) {
			if strings.HasSuffix(rq.URL.Path, ":rotate-keys") {
				clitest.Reply(w, 200, rotated)
				return
			}
			api(e.URL)(w, rq)
		}, a...)
		for _, s := range []string{afsKey, pcKey} {
			if strings.Contains(r.Out+r.Err, s) {
				t.Errorf("%v leaks %q", a, s)
			}
		}
	}
}
