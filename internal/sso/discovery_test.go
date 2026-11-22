package sso

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"
)

// hostileIdP serves a discovery document built by doc and records every request, so tests can prove
// that no form carrying a device code or refresh token was sent.
type hostileIdP struct {
	srv  *httptest.Server
	mu   sync.Mutex
	reqs []string
}

func newHostileIdP(t *testing.T, doc func(base string) string) *hostileIdP {
	t.Helper()
	h := &hostileIdP{}
	h.srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		h.mu.Lock()
		h.reqs = append(h.reqs, r.Method+" "+r.URL.Path+" "+r.Form.Encode())
		h.mu.Unlock()
		if r.URL.Path == "/.well-known/openid-configuration" {
			_, _ = w.Write([]byte(doc(h.srv.URL)))
			return
		}
		_, _ = w.Write([]byte(`{"access_token":"a"}`))
	}))
	t.Cleanup(h.srv.Close)
	return h
}

func (h *hostileIdP) sentSecrets() bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	for _, r := range h.reqs {
		if strings.Contains(r, "device_code") || strings.Contains(r, "refresh_token") {
			return true
		}
	}
	return false
}

// FR-002: discovery must be https (loopback http only for a loopback issuer) and name the issuer asked for.
func TestDiscoveryRejectsUnsafeDocuments(t *testing.T) {
	cases := []struct {
		name string
		doc  func(base string) string
		want string
	}{
		{"issuer differs", func(b string) string {
			return fmt.Sprintf(`{"issuer":"https://other.example","device_authorization_endpoint":"%[1]s/device","token_endpoint":"%[1]s/token"}`, b)
		}, "issuer"},
		{"issuer missing", func(b string) string {
			return fmt.Sprintf(`{"device_authorization_endpoint":"%[1]s/device","token_endpoint":"%[1]s/token"}`, b)
		}, "issuer"},
		{"plain http token endpoint on another host", func(b string) string {
			return fmt.Sprintf(`{"issuer":%q,"device_authorization_endpoint":"%s/device","token_endpoint":"http://idp.example.com/token"}`, b, b)
		}, "token_endpoint"},
		{"plain http device endpoint on another host", func(b string) string {
			return fmt.Sprintf(`{"issuer":%q,"device_authorization_endpoint":"http://idp.example.com/device","token_endpoint":"%s/token"}`, b, b)
		}, "device_authorization_endpoint"},
		{"endpoint is not a URL", func(b string) string {
			return fmt.Sprintf(`{"issuer":%q,"token_endpoint":"/token"}`, b)
		}, "token_endpoint"},
		{"endpoint has userinfo", func(b string) string {
			return fmt.Sprintf(`{"issuer":%q,"token_endpoint":"https://u:p@idp.example.com/token"}`, b)
		}, "token_endpoint"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			h := newHostileIdP(t, c.doc)
			if _, err := Discover(context.Background(), http.DefaultClient, h.srv.URL); err == nil || !strings.Contains(err.Error(), c.want) {
				t.Fatalf("Discover: %v", err)
			}
			clk := &clock{t: time.Now()}
			if _, err := Login(context.Background(), loginOpts(clk, &bytes.Buffer{}), h.srv.URL, "c"); err == nil {
				t.Error("Login succeeded")
			}
			tok := &Token{Issuer: h.srv.URL, ClientID: "c", RefreshToken: "rt-secret", AccessToken: "x", ExpiresAt: time.Now().Add(-time.Hour)}
			if _, err := Refresh(context.Background(), http.DefaultClient, time.Now(), tok); err == nil {
				t.Error("Refresh succeeded")
			}
			if h.sentSecrets() {
				t.Errorf("a form with a secret was sent: %v", h.reqs)
			}
		})
	}
}

func TestDiscoveryRejectsUnsafeIssuers(t *testing.T) {
	for _, iss := range []string{"http://idp.example.com", "https://u:p@idp.example.com", "https://idp.example.com?x=1", "idp.example.com", ""} {
		h := newHostileIdP(t, func(string) string { return `{}` })
		if _, err := Discover(context.Background(), http.DefaultClient, iss); err == nil || !strings.Contains(err.Error(), "sso_issuer") {
			t.Errorf("%q: %v", iss, err)
		}
		if len(h.reqs) != 0 {
			t.Errorf("%q: request sent", iss)
		}
	}
}

func TestDiscoveryAcceptsTrailingSlashIssuer(t *testing.T) {
	h := newHostileIdP(t, func(b string) string {
		return fmt.Sprintf(`{"issuer":"%[1]s/","device_authorization_endpoint":"%[1]s/device","token_endpoint":"%[1]s/token"}`, b)
	})
	if _, err := Discover(context.Background(), http.DefaultClient, h.srv.URL); err != nil {
		t.Fatal(err)
	}
}

// FR-003: the verification URL and user code come from the identity provider and are printed raw otherwise.
func TestDevicePromptIsSanitized(t *testing.T) {
	var srv *httptest.Server
	srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/.well-known/openid-configuration":
			_, _ = fmt.Fprintf(w, `{"issuer":%q,"device_authorization_endpoint":"%s/device","token_endpoint":"%s/token"}`, srv.URL, srv.URL, srv.URL)
		case "/device":
			_, _ = w.Write([]byte(`{"device_code":"d","user_code":"AB\u001b[2JCD","verification_uri":"https://v\u001b]0;x\u0007","verification_uri_complete":"https://v?c=\u009b2J"}`))
		default:
			_, _ = w.Write([]byte(`{"access_token":"a"}`))
		}
	}))
	defer srv.Close()
	clk := &clock{t: time.Now()}
	var out bytes.Buffer
	if _, err := Login(context.Background(), loginOpts(clk, &out), srv.URL, "c"); err != nil {
		t.Fatal(err)
	}
	for _, r := range out.String() {
		if (r < 0x20 && r != '\n') || (r >= 0x7f && r <= 0x9f) {
			t.Fatalf("control character %U in %q", r, out.String())
		}
	}
}
