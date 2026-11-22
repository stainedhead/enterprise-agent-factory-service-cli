package client

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

type sleeps struct{ d []time.Duration }

func (s *sleeps) sleep(ctx context.Context, d time.Duration) error {
	s.d = append(s.d, d)
	return ctx.Err()
}

func newClient(t *testing.T, srv *httptest.Server, mod func(*Config)) (*Client, *sleeps, *bytes.Buffer) {
	t.Helper()
	sl := &sleeps{}
	var errb bytes.Buffer
	cfg := Config{Endpoint: srv.URL, UserAgent: "afs-cli/test (linux/amd64)", Sleep: sl.sleep, Rand: func() float64 { return 0.5 },
		Stderr: &errb, ContractVersion: "1.1.0",
		Credentials: func(context.Context) (Credential, error) {
			return Credential{Header: "Bearer tok", Source: "env AFS_API_KEY"}, nil
		}}
	if mod != nil {
		mod(&cfg)
	}
	return New(cfg), sl, &errb
}

func jsonReply(w http.ResponseWriter, status int, body string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = io.WriteString(w, body)
}

func errBody(code, reason string) string {
	return fmt.Sprintf(`{"code":%q,"message":"m","requestId":"req_x","details":{"reason":%q}}`, code, reason)
}

// CLI-1 / D-C11: each path value is escaped as one segment.
func TestPathEscapingPerSegment(t *testing.T) {
	got := BuildPath("/v1/work-orders/{workOrderId}/repos/{host}/{owner}/{repo}",
		map[string]string{"workOrderId": "wo_1", "host": "github.com", "owner": "pay ments", "repo": "a/b"})
	if want := "/v1/work-orders/wo_1/repos/github.com/pay%20ments/a%2Fb"; got != want {
		t.Errorf("path = %s, want %s", got, want)
	}
	if BuildPath("/x/{a", nil) != "/x/{a" {
		t.Error("unterminated placeholder changed")
	}
}

func TestUserAgentAcceptHeadersAndIfMatch(t *testing.T) {
	var got http.Header
	var gotURL string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got, gotURL = r.Header.Clone(), r.URL.String()
		if r.Header.Get("Accept") == "text/csv" {
			w.Header().Set("Content-Type", "text/csv")
			_, _ = io.WriteString(w, "a,b\r\n")
			return
		}
		if r.Header.Get("Accept") == "application/yaml" {
			w.Header().Set("Content-Type", "application/yaml; charset=utf-8")
			_, _ = io.WriteString(w, "a: b\n")
			return
		}
		jsonReply(w, 200, `{}`)
	}))
	defer srv.Close()
	c, _, _ := newClient(t, srv, nil)
	// CLI-US-4: Accept headers.
	for _, accept := range []string{"", "text/csv", "application/yaml"} {
		if _, err := c.Do(context.Background(), Request{Method: "PATCH", Path: "/v1/x", Query: url.Values{"a": {"b c"}}, Accept: accept, IfMatch: "7", Body: []byte(`{"a":1}`)}); err != nil {
			t.Fatal(err)
		}
		want := accept
		if want == "" {
			want = "application/json"
		}
		if got.Get("Accept") != want {
			t.Errorf("Accept = %q, want %q", got.Get("Accept"), want)
		}
	}
	if got.Get("User-Agent") != "afs-cli/test (linux/amd64)" {
		t.Errorf("User-Agent = %q", got.Get("User-Agent"))
	}
	if got.Get("If-Match") != `"7"` || got.Get("Authorization") != "Bearer tok" || got.Get("Content-Type") != "application/json" {
		t.Errorf("headers = %v", got)
	}
	if gotURL != "/v1/x?a=b+c" {
		t.Errorf("url = %s", gotURL)
	}
	if _, err := c.Do(context.Background(), Request{Method: "GET", Path: "/v1/x", NoAuth: true}); err != nil {
		t.Fatal(err)
	}
	if got.Get("Authorization") != "" {
		t.Error("public request carried a credential")
	}
}

// CLI-8: envelope decoding; a non-JSON 5xx becomes InternalError with the status.
func TestAPIErrorDecode(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/denied":
			w.Header().Set("X-Request-Id", "hdr_id")
			jsonReply(w, 403, `{"code":"AccessDenied","message":"no","requestId":"req_9","details":{"requiredPermission":"teams:write"}}`)
		case "/html":
			w.Header().Set("X-Request-Id", "hdr_id")
			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(502)
			_, _ = io.WriteString(w, "<html>bad gateway</html>")
		default:
			w.WriteHeader(418)
		}
	}))
	defer srv.Close()
	c, _, _ := newClient(t, srv, nil)
	_, err := c.Do(context.Background(), Request{Method: "GET", Path: "/denied"})
	var ae *APIError
	if !errors.As(err, &ae) || ae.Code != "AccessDenied" || ae.Message != "no" || ae.RequestID != "req_9" || ae.RequiredPermission() != "teams:write" || ae.Status != 403 {
		t.Fatalf("err = %#v", err)
	}
	if ae.Error() != "AccessDenied: no" {
		t.Error(ae.Error())
	}
	_, err = c.Do(context.Background(), Request{Method: "GET", Path: "/html"})
	if !errors.As(err, &ae) || ae.Code != "InternalError" || ae.Status != 502 || ae.RequestID != "hdr_id" {
		t.Fatalf("html err = %#v", err)
	}
	_, err = c.Do(context.Background(), Request{Method: "GET", Path: "/other"})
	if !errors.As(err, &ae) || ae.Code != "ServiceError" || ae.Reason() != "" {
		t.Fatalf("418 err = %#v", err)
	}
}

// NFR compatibility: one notice per process for a newer minor or a different major.
func TestContractVersionNoticeOncePerProcess(t *testing.T) {
	cases := []struct {
		server string
		notice bool
	}{{"1.2.0", true}, {"1.3.1", true}, {"1.1.0", false}, {"1.0.5", false}, {"2.0.0", true}, {"", false}, {"junk", false}}
	for _, tc := range cases {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if tc.server != "" {
				w.Header().Set("X-Afs-Contract-Version", tc.server)
			}
			jsonReply(w, 200, `{}`)
		}))
		c, _, errb := newClient(t, srv, nil)
		for i := 0; i < 3; i++ {
			if _, err := c.Do(context.Background(), Request{Method: "GET", Path: "/x"}); err != nil {
				t.Fatal(err)
			}
		}
		if n := strings.Count(errb.String(), "Notice:"); (n == 1) != tc.notice || n > 1 {
			t.Errorf("server %q: %d notices (%q)", tc.server, n, errb.String())
		}
		srv.Close()
	}
}

func TestUnexpectedMediaType255(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		_, _ = io.WriteString(w, "<html/>")
	}))
	defer srv.Close()
	c, _, _ := newClient(t, srv, nil)
	if _, err := c.Do(context.Background(), Request{Method: "GET", Path: "/x"}); !errors.Is(err, ErrUnexpectedMediaType) {
		t.Fatalf("err = %v", err)
	}
}

func TestInvalidEndpoint(t *testing.T) {
	c := New(Config{Endpoint: "not a url"})
	if _, err := c.Do(context.Background(), Request{Method: "GET", Path: "/x"}); err == nil {
		t.Fatal("want error")
	}
	c = New(Config{Endpoint: "http://x", Credentials: func(context.Context) (Credential, error) { return Credential{}, errors.New("no cred") }})
	if _, err := c.Do(context.Background(), Request{Method: "GET", Path: "/x"}); err == nil || err.Error() != "no cred" {
		t.Fatalf("err = %v", err)
	}
}

// ---- retries (CLI-11, 5.1-retries) ----

func TestIdempotencyInProgressTwiceThenSuccess(t *testing.T) {
	var n atomic.Int32
	var keys []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		keys = append(keys, r.Header.Get("Idempotency-Key"))
		if n.Add(1) <= 2 {
			jsonReply(w, 409, errBody("Conflict", "idempotency_in_progress"))
			return
		}
		jsonReply(w, 201, `{"id":"x"}`)
	}))
	defer srv.Close()
	c, sl, _ := newClient(t, srv, nil)
	resp, err := c.Do(context.Background(), Request{Method: "POST", Path: "/v1/teams", Idempotent: true, Body: []byte(`{}`)})
	if err != nil || resp.Attempts != 3 || n.Load() != 3 {
		t.Fatalf("resp %+v err %v attempts %d", resp, err, n.Load())
	}
	if keys[0] == "" || keys[0] != keys[1] || keys[1] != keys[2] {
		t.Errorf("keys = %v", keys)
	}
	if len(sl.d) != 2 {
		t.Errorf("sleeps = %v", sl.d)
	}
}

func TestIdempotencyInProgressThreeTimes(t *testing.T) {
	var n atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n.Add(1)
		jsonReply(w, 409, errBody("Conflict", "idempotency_in_progress"))
	}))
	defer srv.Close()
	c, _, _ := newClient(t, srv, nil)
	_, err := c.Do(context.Background(), Request{Method: "POST", Path: "/x", Idempotent: true})
	var ae *APIError
	if !errors.As(err, &ae) || n.Load() != 3 {
		t.Fatalf("err %v attempts %d", err, n.Load())
	}
}

func TestIdempotencyMismatchNotRetried(t *testing.T) {
	var n atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n.Add(1)
		jsonReply(w, 409, errBody("Conflict", "idempotency_mismatch"))
	}))
	defer srv.Close()
	c, _, _ := newClient(t, srv, nil)
	if _, err := c.Do(context.Background(), Request{Method: "POST", Path: "/x", Idempotent: true}); err == nil || n.Load() != 1 {
		t.Fatalf("err %v attempts %d", err, n.Load())
	}
}

func Test429RetryAfter2WaitsAtLeast2s(t *testing.T) {
	var n atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if n.Add(1) == 1 {
			w.Header().Set("Retry-After", "2")
			jsonReply(w, 429, errBody("Throttled", ""))
			return
		}
		jsonReply(w, 200, `{}`)
	}))
	defer srv.Close()
	c, sl, _ := newClient(t, srv, nil)
	if _, err := c.Do(context.Background(), Request{Method: "PATCH", Path: "/x"}); err != nil {
		t.Fatal(err)
	}
	if len(sl.d) != 1 || sl.d[0] < 2*time.Second {
		t.Errorf("sleeps = %v", sl.d)
	}
}

func Test503RetriedOnGetAndIdempotentPostOnly(t *testing.T) {
	for _, tc := range []struct {
		method    string
		idem      bool
		wantTries int32
	}{{"GET", false, 3}, {"POST", true, 3}, {"PATCH", false, 1}, {"POST", false, 1}} {
		var n atomic.Int32
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			n.Add(1)
			jsonReply(w, 503, errBody("Unavailable", ""))
		}))
		c, _, _ := newClient(t, srv, nil)
		_, _ = c.Do(context.Background(), Request{Method: tc.method, Path: "/x", Idempotent: tc.idem})
		if n.Load() != tc.wantTries {
			t.Errorf("%s idem=%v: %d tries, want %d", tc.method, tc.idem, n.Load(), tc.wantTries)
		}
		srv.Close()
	}
}

func Test4xxNotRetried(t *testing.T) {
	var n atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n.Add(1)
		jsonReply(w, 400, errBody("ValidationError", ""))
	}))
	defer srv.Close()
	c, _, _ := newClient(t, srv, nil)
	_, _ = c.Do(context.Background(), Request{Method: "GET", Path: "/x"})
	if n.Load() != 1 {
		t.Errorf("tries = %d", n.Load())
	}
}

// A PATCH whose connection drops after the request was sent is not retried; a GET is.
func TestPatchDroppedAfterSendNotRetried(t *testing.T) {
	var n atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n.Add(1)
		conn, _, _ := w.(http.Hijacker).Hijack()
		_ = conn.Close()
	}))
	defer srv.Close()
	c, _, _ := newClient(t, srv, nil)
	if _, err := c.Do(context.Background(), Request{Method: "PATCH", Path: "/x", Body: []byte(`{}`)}); err == nil || n.Load() != 1 {
		t.Fatalf("err %v tries %d", err, n.Load())
	}
	n.Store(0)
	_, _ = c.Do(context.Background(), Request{Method: "GET", Path: "/x"})
	if n.Load() != 3 {
		t.Errorf("GET tries = %d, want 3", n.Load())
	}
}

// A connect failure is retried on any request (and the injected dialer is used).
func TestPatchConnectFailureRetried(t *testing.T) {
	var dials atomic.Int32
	tr := NewTransport(TransportOptions{DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
		dials.Add(1)
		return nil, &net.OpError{Op: "dial", Net: network, Err: errors.New("connection refused")}
	}})
	c := New(Config{Endpoint: "http://127.0.0.1:9", HTTP: tr, Sleep: (&sleeps{}).sleep, Rand: func() float64 { return 0 }})
	if _, err := c.Do(context.Background(), Request{Method: "PATCH", Path: "/x"}); err == nil || dials.Load() != 3 {
		t.Fatalf("err %v dials %d", err, dials.Load())
	}
}

func TestULIDPerInvocationOrFlag(t *testing.T) {
	var keys []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		keys = append(keys, r.Header.Get("Idempotency-Key"))
		jsonReply(w, 200, `{}`)
	}))
	defer srv.Close()
	c, _, _ := newClient(t, srv, nil)
	ulid := regexp.MustCompile(`^[0-9A-HJKMNP-TV-Z]{26}$`)
	for i := 0; i < 2; i++ {
		_, _ = c.Do(context.Background(), Request{Method: "POST", Path: "/x", Idempotent: true})
	}
	_, _ = c.Do(context.Background(), Request{Method: "POST", Path: "/x", Idempotent: true, IdempotencyKey: "my-key"})
	_, _ = c.Do(context.Background(), Request{Method: "PATCH", Path: "/x"})
	if !ulid.MatchString(keys[0]) || keys[0] == keys[1] || keys[2] != "my-key" || keys[3] != "" {
		t.Errorf("keys = %v", keys)
	}
}

func TestReplayHeaderSurfaced(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Idempotent-Replayed", "true")
		w.Header().Set("X-Request-Id", "req_7")
		jsonReply(w, 201, `{"id":"x"}`)
	}))
	defer srv.Close()
	c, _, _ := newClient(t, srv, nil)
	resp, err := c.Do(context.Background(), Request{Method: "POST", Path: "/x", Idempotent: true})
	if err != nil || !resp.Replayed || resp.RequestID() != "req_7" {
		t.Fatalf("resp %+v err %v", resp, err)
	}
}

func TestBackoffAndRetryAfterParsing(t *testing.T) {
	if d := backoff(1, func() float64 { return 1 }, 0); d != 200*time.Millisecond {
		t.Errorf("attempt 1 ceiling = %s", d)
	}
	if d := backoff(2, func() float64 { return 1 }, 0); d != 400*time.Millisecond {
		t.Errorf("attempt 2 ceiling = %s", d)
	}
	if d := backoff(1, func() float64 { return 0 }, 3*time.Second); d != 3*time.Second {
		t.Errorf("retry-after = %s", d)
	}
	h := http.Header{}
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	h.Set("Retry-After", now.Add(5*time.Second).UTC().Format(http.TimeFormat))
	if d := parseRetryAfter(h, now); d != 5*time.Second {
		t.Errorf("http-date = %s", d)
	}
	h.Set("Retry-After", "junk")
	if parseRetryAfter(h, now) != 0 || parseRetryAfter(http.Header{}, now) != 0 {
		t.Error("junk parsed")
	}
}

// ---- timeouts and cancellation (T3) ----

func TestConnectTimeoutDefault30s(t *testing.T) {
	for _, tc := range []struct {
		timeout  time.Duration
		deadline bool
	}{{DefaultConnectTimeout, true}, {0, false}} {
		var dl time.Time
		var has bool
		tr := NewTransport(TransportOptions{ConnectTimeout: tc.timeout, DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			dl, has = ctx.Deadline()
			return nil, errors.New("stop")
		}})
		c := New(Config{Endpoint: "http://127.0.0.1:9", HTTP: tr, Sleep: (&sleeps{}).sleep})
		_, _ = c.Do(context.Background(), Request{Method: "POST", Path: "/x"})
		if has != tc.deadline {
			t.Errorf("timeout %s: deadline=%v", tc.timeout, has)
		}
		if has && (time.Until(dl) > 30*time.Second || time.Until(dl) < 29*time.Second) {
			t.Errorf("deadline in %s", time.Until(dl))
		}
	}
}

func TestReadTimeoutHeadersAndBody(t *testing.T) {
	release := make(chan struct{})
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/headers":
			select {
			case <-release:
			case <-r.Context().Done():
			}
		case "/body":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			_, _ = io.WriteString(w, `{"a":`)
			w.(http.Flusher).Flush()
			select {
			case <-release:
			case <-r.Context().Done():
			}
		}
	}))
	defer srv.Close()
	defer close(release)
	c, _, _ := newClient(t, srv, func(cfg *Config) { cfg.ReadTimeout = 100 * time.Millisecond })
	for _, p := range []string{"/headers", "/body"} {
		_, err := c.Do(context.Background(), Request{Method: "PATCH", Path: p})
		if !errors.Is(err, ErrReadTimeout) {
			t.Errorf("%s: err = %v", p, err)
		}
	}
	// 0 disables the timeout: a quick reply works.
	srv2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { jsonReply(w, 200, `{}`) }))
	defer srv2.Close()
	c2, _, _ := newClient(t, srv2, nil)
	if _, err := c2.Do(context.Background(), Request{Method: "GET", Path: "/"}); err != nil {
		t.Fatal(err)
	}
}

// CLI-US-5.3: cancelling the context aborts the in-flight request.
func TestContextCancelAbortsInFlight(t *testing.T) {
	started := make(chan struct{})
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		close(started)
		<-r.Context().Done()
	}))
	defer srv.Close()
	c, _, _ := newClient(t, srv, nil)
	ctx, cancel := context.WithCancel(context.Background())
	go func() { <-started; cancel() }()
	if _, err := c.Do(ctx, Request{Method: "GET", Path: "/x"}); !errors.Is(err, context.Canceled) {
		t.Fatalf("err = %v", err)
	}
}

func TestCancelDuringBackoff(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { jsonReply(w, 429, errBody("Throttled", "")) }))
	defer srv.Close()
	ctx, cancel := context.WithCancel(context.Background())
	c, _, _ := newClient(t, srv, func(cfg *Config) {
		cfg.Sleep = func(ctx context.Context, d time.Duration) error { cancel(); return ctx.Err() }
	})
	if _, err := c.Do(ctx, Request{Method: "GET", Path: "/x"}); !errors.Is(err, context.Canceled) {
		t.Fatalf("err = %v", err)
	}
}

func TestDefaultsAreWired(t *testing.T) {
	c := New(Config{Endpoint: "http://x"})
	if c.cfg.HTTP == nil || c.cfg.Now == nil || c.cfg.Sleep == nil || c.cfg.Rand == nil {
		t.Fatal("defaults missing")
	}
	if err := c.cfg.Sleep(context.Background(), time.Millisecond); err != nil {
		t.Error(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := c.cfg.Sleep(ctx, time.Hour); err == nil {
		t.Error("sleep ignored cancel")
	}
	if r := c.cfg.Rand(); r < 0 || r >= 1 {
		t.Errorf("rand = %v", r)
	}
	if NewIdempotencyKey() == NewIdempotencyKey() { //nolint:staticcheck // two calls must differ
		t.Error("keys equal")
	}
}

// FR-011: Retry-After is capped (20 s, or the read timeout when smaller); over the cap the client gives up.
func TestRetryAfterOverCapGivesUpWithTheAPIError(t *testing.T) {
	future := time.Now().Add(26 * time.Hour).UTC().Format(http.TimeFormat)
	for name, c := range map[string]struct {
		header      string
		readTimeout time.Duration
	}{
		"delta seconds":       {"86400", 0},
		"HTTP date":           {future, 0},
		"just over 20s":       {"21", 0},
		"over the read limit": {"5", 3 * time.Second},
	} {
		t.Run(name, func(t *testing.T) {
			var n atomic.Int32
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				n.Add(1)
				w.Header().Set("Retry-After", c.header)
				jsonReply(w, 429, errBody("Throttled", ""))
			}))
			defer srv.Close()
			var dbg bytes.Buffer
			cl, sl, _ := newClient(t, srv, func(cfg *Config) { cfg.ReadTimeout = c.readTimeout; cfg.Debug = &dbg; cfg.Now = time.Now })
			_, err := cl.Do(context.Background(), Request{Method: "GET", Path: "/x"})
			var ae *APIError
			if !errors.As(err, &ae) || ae.Status != 429 {
				t.Fatalf("err = %v", err)
			}
			if n.Load() != 1 || len(sl.d) != 0 {
				t.Errorf("attempts=%d sleeps=%v", n.Load(), sl.d)
			}
			if !strings.Contains(dbg.String(), "Retry-After") || !strings.Contains(dbg.String(), "not retrying") {
				t.Errorf("debug does not say why:\n%s", dbg.String())
			}
		})
	}
}

func TestRetryAfterWithinCapIsHonouredAndLogged(t *testing.T) {
	var n atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if n.Add(1) == 1 {
			w.Header().Set("Retry-After", "20")
			jsonReply(w, 429, errBody("Throttled", ""))
			return
		}
		jsonReply(w, 200, `{}`)
	}))
	defer srv.Close()
	var dbg bytes.Buffer
	cl, sl, _ := newClient(t, srv, func(cfg *Config) { cfg.Debug = &dbg; cfg.Now = time.Now })
	if _, err := cl.Do(context.Background(), Request{Method: "GET", Path: "/x"}); err != nil {
		t.Fatal(err)
	}
	if len(sl.d) != 1 || sl.d[0] != 20*time.Second || !strings.Contains(dbg.String(), "Retry-After 20s") {
		t.Errorf("sleeps=%v\n%s", sl.d, dbg.String())
	}
}
