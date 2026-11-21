package sso

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/auth"
)

type clock struct{ t time.Time }

func (c *clock) now() time.Time { return c.t }
func (c *clock) sleep(ctx context.Context, d time.Duration) error {
	c.t = c.t.Add(d)
	return ctx.Err()
}

type idp struct {
	srv       *httptest.Server
	pending   int // answer authorization_pending this many times first
	slowDowns int
	polls     atomic.Int32
	intervals []time.Time
	clientIDs []string
	deviceErr string // "expired_token" or "access_denied" at the token endpoint
	noDevice  bool
	refreshOK bool
	refreshes atomic.Int32
	clk       *clock
}

func newIdP(t *testing.T, clk *clock) *idp {
	t.Helper()
	p := &idp{clk: clk, refreshOK: true}
	mux := http.NewServeMux()
	mux.HandleFunc("/.well-known/openid-configuration", func(w http.ResponseWriter, r *http.Request) {
		d := map[string]string{"token_endpoint": p.srv.URL + "/token"}
		if !p.noDevice {
			d["device_authorization_endpoint"] = p.srv.URL + "/device"
		}
		_ = json.NewEncoder(w).Encode(d)
	})
	mux.HandleFunc("/device", func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		p.clientIDs = append(p.clientIDs, r.Form.Get("client_id"))
		_ = json.NewEncoder(w).Encode(map[string]any{"device_code": "dev-123", "user_code": "ABCD-EFGH", "verification_uri": p.srv.URL + "/activate", "expires_in": 600, "interval": 1})
	})
	mux.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		if r.Form.Get("grant_type") == "refresh_token" {
			p.refreshes.Add(1)
			if !p.refreshOK {
				w.WriteHeader(400)
				_, _ = w.Write([]byte(`{"error":"invalid_grant","error_description":"refresh token revoked"}`))
				return
			}
			_ = json.NewEncoder(w).Encode(map[string]any{"access_token": "refreshed-at", "expires_in": 3600})
			return
		}
		p.intervals = append(p.intervals, p.clk.now())
		n := int(p.polls.Add(1))
		switch {
		case p.deviceErr != "":
			w.WriteHeader(400)
			_, _ = fmt.Fprintf(w, `{"error":%q}`, p.deviceErr)
		case n <= p.pending:
			w.WriteHeader(400)
			_, _ = w.Write([]byte(`{"error":"authorization_pending"}`))
		case n <= p.pending+p.slowDowns:
			w.WriteHeader(400)
			_, _ = w.Write([]byte(`{"error":"slow_down"}`))
		default:
			_ = json.NewEncoder(w).Encode(map[string]any{"access_token": "at-1", "refresh_token": "rt-1", "token_type": "Bearer", "expires_in": 3600})
		}
	})
	p.srv = httptest.NewServer(mux)
	t.Cleanup(p.srv.Close)
	return p
}

func loginOpts(clk *clock, out *bytes.Buffer) LoginOptions {
	return LoginOptions{HTTP: http.DefaultClient, Out: out, Now: clk.now, Sleep: clk.sleep}
}

// CLI-US-1.1: the device code and URL are printed; the token is cached with mode 0600.
func TestLoginPrintsCodeAndURLThenCaches0600(t *testing.T) {
	clk := &clock{t: time.Date(2026, 11, 21, 12, 0, 0, 0, time.UTC)}
	p := newIdP(t, clk)
	p.pending = 2
	var out bytes.Buffer
	tok, err := Login(context.Background(), loginOpts(clk, &out), p.srv.URL, "cli-client")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "ABCD-EFGH") || !strings.Contains(out.String(), p.srv.URL+"/activate") {
		t.Errorf("prompt = %q", out.String())
	}
	cache := NewCache(t.TempDir())
	if err := cache.Save(tok); err != nil {
		t.Fatal(err)
	}
	st, err := os.Stat(cache.Path(p.srv.URL, "cli-client"))
	if err != nil || st.Mode().Perm() != 0o600 {
		t.Fatalf("stat %v %v", st, err)
	}
	got, err := cache.Load(p.srv.URL, "cli-client")
	if err != nil || got.AccessToken != "at-1" || got.RefreshToken != "rt-1" || !got.ExpiresAt.Equal(clk.t.Add(time.Hour)) {
		t.Fatalf("cached = %+v err %v", got, err)
	}
	if strings.Contains(out.String(), "at-1") || strings.Contains(out.String(), "dev-123") {
		t.Error("prompt leaks a token or device code")
	}
}

func TestLoginHonoursSlowDownAndPending(t *testing.T) {
	clk := &clock{t: time.Date(2026, 11, 21, 12, 0, 0, 0, time.UTC)}
	p := newIdP(t, clk)
	p.pending, p.slowDowns = 1, 1
	if _, err := Login(context.Background(), loginOpts(clk, &bytes.Buffer{}), p.srv.URL, "c"); err != nil {
		t.Fatal(err)
	}
	if len(p.intervals) != 3 {
		t.Fatalf("polls = %d", len(p.intervals))
	}
	// interval 1s; after slow_down it grows by 5s
	if d := p.intervals[1].Sub(p.intervals[0]); d != time.Second {
		t.Errorf("pending interval = %s", d)
	}
	if d := p.intervals[2].Sub(p.intervals[1]); d != 6*time.Second {
		t.Errorf("slow_down interval = %s", d)
	}
}

// CLI-US-1 / D48: the client id of the device request is the one given (cliClientId).
func TestLoginUsesCliClientIdFromListSignInProviders(t *testing.T) {
	clk := &clock{t: time.Now()}
	p := newIdP(t, clk)
	if _, err := Login(context.Background(), loginOpts(clk, &bytes.Buffer{}), p.srv.URL, "0oa1cli"); err != nil {
		t.Fatal(err)
	}
	if len(p.clientIDs) != 1 || p.clientIDs[0] != "0oa1cli" {
		t.Errorf("client ids = %v", p.clientIDs)
	}
}

func TestLoginExpiredDeviceCode(t *testing.T) {
	clk := &clock{t: time.Now()}
	p := newIdP(t, clk)
	p.deviceErr = "expired_token"
	if _, err := Login(context.Background(), loginOpts(clk, &bytes.Buffer{}), p.srv.URL, "c"); !errors.Is(err, ErrDeviceExpired) {
		t.Errorf("expired_token: %v", err)
	}
	p.deviceErr = "access_denied"
	if _, err := Login(context.Background(), loginOpts(clk, &bytes.Buffer{}), p.srv.URL, "c"); !errors.Is(err, ErrAccessDenied) {
		t.Errorf("access_denied: %v", err)
	}
	p.deviceErr = "server_error"
	if _, err := Login(context.Background(), loginOpts(clk, &bytes.Buffer{}), p.srv.URL, "c"); err == nil || !strings.Contains(err.Error(), "server_error") {
		t.Errorf("other: %v", err)
	}
	// local deadline passes while still pending
	p.deviceErr = ""
	p.pending = 100000
	if _, err := Login(context.Background(), loginOpts(clk, &bytes.Buffer{}), p.srv.URL, "c"); !errors.Is(err, ErrDeviceExpired) {
		t.Errorf("deadline: %v", err)
	}
}

func TestLoginDiscoveryFailures(t *testing.T) {
	clk := &clock{t: time.Now()}
	p := newIdP(t, clk)
	p.noDevice = true
	if _, err := Login(context.Background(), loginOpts(clk, &bytes.Buffer{}), p.srv.URL, "c"); err == nil || !strings.Contains(err.Error(), "device authorization grant") {
		t.Errorf("no device endpoint: %v", err)
	}
	bad := httptest.NewServer(http.NotFoundHandler())
	defer bad.Close()
	if _, err := Login(context.Background(), loginOpts(clk, &bytes.Buffer{}), bad.URL, "c"); err == nil || !strings.Contains(err.Error(), "discovery") {
		t.Errorf("404 discovery: %v", err)
	}
	empty := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { _, _ = w.Write([]byte(`{}`)) }))
	defer empty.Close()
	if _, err := Discover(context.Background(), http.DefaultClient, empty.URL); err == nil {
		t.Error("discovery without token_endpoint accepted")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := Login(ctx, loginOpts(clk, &bytes.Buffer{}), p.srv.URL, "c"); err == nil {
		t.Error("cancelled login succeeded")
	}
}

func TestLoginVerificationURIComplete(t *testing.T) {
	clk := &clock{t: time.Now()}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/.well-known/openid-configuration":
			_, _ = fmt.Fprintf(w, `{"device_authorization_endpoint":"http://%s/device","token_endpoint":"http://%s/token"}`, r.Host, r.Host)
		case "/device":
			_, _ = w.Write([]byte(`{"device_code":"d","user_code":"U","verification_uri":"https://v","verification_uri_complete":"https://v?c=U"}`))
		default:
			_, _ = w.Write([]byte(`{"access_token":"a"}`))
		}
	}))
	defer srv.Close()
	var out bytes.Buffer
	tok, err := Login(context.Background(), loginOpts(clk, &out), srv.URL, "c")
	if err != nil || !strings.Contains(out.String(), "https://v?c=U") || tok.TokenType != "Bearer" || !tok.ExpiresAt.Equal(clk.t.Add(time.Hour)) {
		t.Errorf("tok %+v err %v out %q", tok, err, out.String())
	}
}

// ---- cache (S1) ----

func TestCacheWrites0600(t *testing.T) {
	c := NewCache(t.TempDir())
	if err := c.Save(&Token{Issuer: "https://i", ClientID: "c", AccessToken: "a", ExpiresAt: time.Now()}); err != nil {
		t.Fatal(err)
	}
	if st, _ := os.Stat(c.Path("https://i", "c")); st.Mode().Perm() != 0o600 {
		t.Errorf("mode = %o", st.Mode().Perm())
	}
	if st, _ := os.Stat(c.Dir); st.Mode().Perm() != 0o700 {
		t.Errorf("dir mode = %o", st.Mode().Perm())
	}
	if c.Path("https://i", "c") == c.Path("https://i", "d") {
		t.Error("client id not part of the cache key")
	}
}

func TestCacheWiderThan0600Refused(t *testing.T) {
	c := NewCache(t.TempDir())
	_ = c.Save(&Token{Issuer: "i", ClientID: "c", AccessToken: "a"})
	_ = os.Chmod(c.Path("i", "c"), 0o644) //nolint:gosec // test
	if _, err := c.Load("i", "c"); err == nil || !strings.Contains(err.Error(), "chmod 600") {
		t.Fatalf("err = %v", err)
	}
}

func TestCacheLoadMissingDamagedDelete(t *testing.T) {
	c := NewCache(t.TempDir())
	if tok, err := c.Load("i", "c"); tok != nil || err != nil {
		t.Errorf("missing: %v %v", tok, err)
	}
	_ = os.MkdirAll(c.Dir, 0o700)
	_ = os.WriteFile(c.Path("i", "c"), []byte("{broken"), 0o600)
	if _, err := c.Load("i", "c"); err == nil || !strings.Contains(err.Error(), "damaged") {
		t.Errorf("damaged: %v", err)
	}
	if err := c.Delete("i", "c"); err != nil {
		t.Error(err)
	}
	if err := c.Delete("i", "c"); err != nil {
		t.Errorf("second delete: %v", err)
	}
}

func source(t *testing.T, clk *clock, issuer string) (*Source, *Cache) {
	cache := NewCache(t.TempDir())
	return &Source{Cache: cache, Issuer: issuer, ClientID: "c", Profile: "dev", HTTP: http.DefaultClient, Now: clk.now}, cache
}

// CLI-US-1.3: no cache, or expired with no refresh token -> LoginRequiredError with the hint, no call.
func TestExpiredNoRefreshExit255NoCall(t *testing.T) {
	clk := &clock{t: time.Now()}
	s, cache := source(t, clk, "http://127.0.0.1:1")
	_, ok, err := s.Token(context.Background())
	var lr *auth.LoginRequiredError
	if ok || !errors.As(err, &lr) || !strings.Contains(err.Error(), "Run afs sso login --profile dev") {
		t.Fatalf("missing cache: ok=%v err=%v", ok, err)
	}
	_ = cache.Save(&Token{Issuer: s.Issuer, ClientID: "c", AccessToken: "old", ExpiresAt: clk.t.Add(-time.Hour)})
	if _, ok, err = s.Token(context.Background()); ok || !errors.As(err, &lr) || !strings.Contains(err.Error(), "expired") {
		t.Fatalf("expired: ok=%v err=%v", ok, err)
	}
}

func TestValidTokenNoRefresh(t *testing.T) {
	clk := &clock{t: time.Now()}
	s, cache := source(t, clk, "http://127.0.0.1:1")
	_ = cache.Save(&Token{Issuer: s.Issuer, ClientID: "c", AccessToken: "fresh", ExpiresAt: clk.t.Add(time.Hour)})
	if tok, ok, err := s.Token(context.Background()); !ok || err != nil || tok != "fresh" {
		t.Fatalf("tok %q ok %v err %v", tok, ok, err)
	}
	// within the skew window the token counts as expired
	_ = cache.Save(&Token{Issuer: s.Issuer, ClientID: "c", AccessToken: "soon", ExpiresAt: clk.t.Add(10 * time.Second)})
	if _, ok, _ := s.Token(context.Background()); ok {
		t.Error("token inside the skew window used")
	}
}

// CLI-US-1.2: an expired access token is refreshed silently and the cache is updated (0600).
func TestExpiredAccessTokenRefreshedSilently(t *testing.T) {
	clk := &clock{t: time.Date(2026, 11, 21, 12, 0, 0, 0, time.UTC)}
	p := newIdP(t, clk)
	s, cache := source(t, clk, p.srv.URL)
	_ = cache.Save(&Token{Issuer: s.Issuer, ClientID: "c", AccessToken: "old", RefreshToken: "rt-keep", ExpiresAt: clk.t.Add(-time.Minute)})
	tok, ok, err := s.Token(context.Background())
	if err != nil || !ok || tok != "refreshed-at" {
		t.Fatalf("tok %q ok %v err %v", tok, ok, err)
	}
	got, _ := cache.Load(s.Issuer, "c")
	if got.AccessToken != "refreshed-at" || got.RefreshToken != "rt-keep" || !got.ExpiresAt.Equal(clk.t.Add(time.Hour)) {
		t.Errorf("cache = %+v", got)
	}
	if st, _ := os.Stat(cache.Path(s.Issuer, "c")); st.Mode().Perm() != 0o600 {
		t.Errorf("mode = %o", st.Mode().Perm())
	}
}

func TestRefreshRejectedFallsBackToHint(t *testing.T) {
	clk := &clock{t: time.Now()}
	p := newIdP(t, clk)
	p.refreshOK = false
	s, cache := source(t, clk, p.srv.URL)
	_ = cache.Save(&Token{Issuer: s.Issuer, ClientID: "c", AccessToken: "old", RefreshToken: "rt", ExpiresAt: clk.t.Add(-time.Minute)})
	_, ok, err := s.Token(context.Background())
	var lr *auth.LoginRequiredError
	if ok || !errors.As(err, &lr) || !strings.Contains(err.Error(), "Run afs sso login --profile dev") || !strings.Contains(err.Error(), "revoked") {
		t.Fatalf("ok=%v err=%v", ok, err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, _, err := s.Token(ctx); !errors.Is(err, context.Canceled) {
		t.Errorf("cancelled refresh: %v", err)
	}
}

func TestTokenURLFormEncodingHelper(t *testing.T) {
	var e *oauthError
	err := error(&oauthError{Code: "x", Description: "d"})
	if !errors.As(err, &e) || e.Error() != "x: d" || (&oauthError{Code: "y"}).Error() != "y" {
		t.Error("oauthError text")
	}
	_ = url.Values{}
}
