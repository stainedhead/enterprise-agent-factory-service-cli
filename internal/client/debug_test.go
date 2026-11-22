package client

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/golden"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/redact"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/spec"
)

// CLI-US-3.4: the debug log names the credential source, redacts Authorization and body fields.
func TestDebugLogGolden(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Request-Id", "req_42")
		jsonReply(w, 201, `{"id":"key_1","secretKey":"afsk_topsecret","name":"ci"}`)
	}))
	defer srv.Close()
	var dbg bytes.Buffer
	tick := time.Date(2026, 11, 21, 12, 0, 0, 0, time.UTC)
	c, _, _ := newClient(t, srv, func(cfg *Config) {
		cfg.Debug = &dbg
		cfg.Now = func() time.Time { tick = tick.Add(10 * time.Millisecond); return tick }
	})
	_, err := c.Do(context.Background(), Request{Method: "POST", Path: "/v1/teams/t/api-keys", Operation: "CreateApiKey", Idempotent: true,
		IdempotencyKey: "idem-1", Body: []byte(`{"name":"ci","value":"hunter2"}`)})
	if err != nil {
		t.Fatal(err)
	}
	got := strings.ReplaceAll(dbg.String(), srv.URL, "http://HOST")
	for _, leak := range []string{"Bearer tok", "hunter2", "afsk_topsecret"} {
		if strings.Contains(got, leak) {
			t.Fatalf("debug log leaks %q:\n%s", leak, got)
		}
	}
	golden.Check(t, "debug.golden", got)
}

func TestDebugOffWritesNothingAndTruncates(t *testing.T) {
	var d *debugLog
	d.printf("x", "y") // nil-safe
	var buf bytes.Buffer
	d = &debugLog{w: &buf, now: time.Now}
	d.body("c", "big", []byte(`{"a":"`+strings.Repeat("x", 5000)+`"}`))
	if !strings.Contains(buf.String(), "truncated") {
		t.Error("long body not truncated")
	}
}

// No query parameter of any operation carries a secret (so URLs are safe to log).
func TestNoQueryParamCarriesSecret(t *testing.T) {
	secret := map[string]bool{}
	for _, f := range redact.Fields() {
		secret[strings.ToLower(f)] = true
	}
	for _, op := range spec.All() {
		for _, p := range op.Params {
			if p.In == "query" && secret[strings.ToLower(p.Name)] {
				t.Errorf("%s: query parameter %s may carry a secret", op.ID, p.Name)
			}
		}
	}
}

// FR-005 (defence in depth): a request URL with credentials is logged redacted.
func TestDebugRequestLineRedactsUserinfo(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { jsonReply(w, 200, `{}`) }))
	defer srv.Close()
	var dbg bytes.Buffer
	c, _, _ := newClient(t, srv, func(cfg *Config) {
		cfg.Debug = &dbg
		cfg.Endpoint = strings.Replace(srv.URL, "http://", "http://user:hunter2@", 1)
	})
	if _, err := c.Do(context.Background(), Request{Method: "GET", Path: "/v1/teams", Operation: "ListTeams"}); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(dbg.String(), "hunter2") || !strings.Contains(dbg.String(), "xxxxx@") {
		t.Errorf("%s", dbg.String())
	}
}
