package edge

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestOrigin(t *testing.T) {
	for _, c := range []struct {
		raw, domain, want string
		bad               bool
	}{
		{"https://t.collab.example/api", "", "https://t.collab.example", false},
		{"https://T.Collab.Example:443/api/x", "collab.example", "https://t.collab.example", false},
		{"http://127.0.0.1:8080/api", "collab.example", "http://127.0.0.1:8080", false},
		{"http://localhost:9/", "", "http://localhost:9", false},
		{"", "", "", true},
		{"https:///x", "", "", true},
		{"http://t.collab.example", "", "", true},
		{"ftp://t.collab.example", "", "", true},
		{"https://u:p@t.collab.example", "", "", true},
		{"https://t.collab.example:8443", "", "", true},
		{"https://10.0.0.1", "", "", true},
		{"https://[::1]", "", "", true},
		{"https://[2001:db8::1]:443", "", "", true},
		{"https://t.collab.example/?a=1", "", "", true},
		{"https://t.collab.example/#a", "", "", true},
		{"https://t.other.example", "collab.example", "", true},
		{"https://a.b.collab.example", "collab.example", "", true},
		{"https://collab.example", "collab.example", "", true},
		{"https://xcollab.example", "collab.example", "", true},
		{"https://localhost", "", "https://localhost", false},
	} {
		got, err := Origin(c.raw, c.domain)
		if (err != nil) != c.bad || got != c.want {
			t.Errorf("%q (%q): got %q, %v", c.raw, c.domain, got, err)
		}
		var oe *OriginError
		if err != nil && (!errors.As(err, &oe) || !strings.Contains(err.Error(), "refused")) {
			t.Errorf("%q: error %v", c.raw, err)
		}
	}
}

func TestReadKey(t *testing.T) {
	dir := t.TempDir()
	write := func(content string, mode os.FileMode) string {
		p := filepath.Join(dir, "k"+content[:1]+mode.String())
		_ = os.WriteFile(p, []byte(content), mode)
		_ = os.Chmod(p, mode)
		return p
	}
	if k, err := ReadKey(write("abc123\r\n", 0o600), nil); err != nil || k != "abc123" {
		t.Errorf("%q %v", k, err)
	}
	if k, err := ReadKey("-", strings.NewReader("stdin-key\n")); err != nil || k != "stdin-key" {
		t.Errorf("%q %v", k, err)
	}
	if _, err := ReadKey("-", nil); err == nil {
		t.Error("nil stdin")
	}
	if _, err := ReadKey(write("abc", 0o640), nil); err == nil || !strings.Contains(err.Error(), "chmod 600") {
		t.Errorf("mode: %v", err)
	}
	if _, err := ReadKey(filepath.Join(dir, "nope"), nil); err == nil {
		t.Error("missing file")
	}
	for _, in := range []string{"", " \n", "a b", "a\x00b", strings.Repeat("x", MaxKeyBytes+1), "a\tb"} {
		_, err := ReadKey("-", strings.NewReader(in))
		var ks *KeyShapeError
		if !errors.As(err, &ks) || strings.Contains(err.Error(), "a b") {
			t.Errorf("%q: %v", in, err)
		}
	}
}

func newClient(rt http.RoundTripper, sleeps *[]time.Duration) *Client {
	return &Client{HTTP: rt, Now: time.Now, Rand: func() float64 { return 1 }, Sleep: func(_ context.Context, d time.Duration) error { *sleeps = append(*sleeps, d); return nil }}
}

func TestDoHeadersAndNoCompression(t *testing.T) {
	var got http.Header
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { got = r.Header.Clone(); _, _ = w.Write([]byte("hi")) }))
	defer srv.Close()
	var s []time.Duration
	c := newClient(http.DefaultTransport, &s)
	resp, err := c.Do(context.Background(), Request{Origin: srv.URL, Method: "GET", Path: "/a", AFSAuth: "Bearer a", Key: "k", UserAgent: "ua"})
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Close()
	b, _ := io.ReadAll(resp)
	if string(b) != "hi" || got.Get("Accept-Encoding") != "identity" || got.Get("Authorization") != "Bearer k" || got.Get("Content-Type") != "" {
		t.Errorf("%q %v", b, got)
	}
}

type flaky struct {
	n    int
	fail int
	rt   http.RoundTripper
}

func (f *flaky) RoundTrip(r *http.Request) (*http.Response, error) {
	f.n++
	if f.n <= f.fail {
		return nil, &net.OpError{Op: "dial", Err: errors.New("refused")}
	}
	return f.rt.RoundTrip(r)
}

func TestDoRetriesDialErrorsOnly(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(204) }))
	defer srv.Close()
	var s []time.Duration
	f := &flaky{fail: 2, rt: http.DefaultTransport}
	resp, err := newClient(f, &s).Do(context.Background(), Request{Origin: srv.URL, Method: "POST", Path: "/", Body: []byte("{}"), ContentType: "application/json"})
	if err != nil || resp.Attempts != 3 || len(s) != 2 {
		t.Fatalf("%v %v %v", resp, err, s)
	}
	resp.Close()
	s = nil
	f = &flaky{fail: 9, rt: http.DefaultTransport}
	if _, err = newClient(f, &s).Do(context.Background(), Request{Origin: srv.URL, Method: "GET", Path: "/"}); err == nil || f.n != 3 {
		t.Errorf("%v %d", err, f.n)
	}
	// an error after the request was sent is not retried
	boom := roundTripFunc(func(*http.Request) (*http.Response, error) { return nil, errors.New("reset") })
	if _, err = newClient(boom, &s).Do(context.Background(), Request{Origin: srv.URL, Method: "GET", Path: "/"}); err == nil {
		t.Error("expected the error")
	}
	// cancel
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err = newClient(&flaky{fail: 9, rt: http.DefaultTransport}, &s).Do(ctx, Request{Origin: srv.URL, Method: "GET", Path: "/"}); !errors.Is(err, context.Canceled) {
		t.Errorf("%v", err)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestDoRetryAfterHonoured(t *testing.T) {
	n := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		n++
		if n == 1 {
			w.Header().Set("Retry-After", "3")
			w.WriteHeader(429)
			return
		}
		if n == 2 {
			w.Header().Set("Retry-After", "9999")
			w.WriteHeader(503)
			return
		}
		w.WriteHeader(200)
	}))
	defer srv.Close()
	var s []time.Duration
	resp, err := newClient(http.DefaultTransport, &s).Do(context.Background(), Request{Origin: srv.URL, Method: "HEAD", Path: "/"})
	// 429 with Retry-After 3 is waited for; the 503 that asks for 9999 s is returned as it is (FR-R1)
	if err != nil || resp.Status != 503 || len(s) != 1 || s[0] < 3*time.Second || n != 2 {
		t.Fatalf("%v %v %v", resp, err, s)
	}
	resp.Close()
}

func TestDoIdleTimeoutAndUpgrade(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.(http.Flusher).Flush()
		<-r.Context().Done()
	}))
	defer srv.Close()
	var s []time.Duration
	c := newClient(http.DefaultTransport, &s)
	c.ReadTimeout = 200 * time.Millisecond
	resp, err := c.Do(context.Background(), Request{Origin: srv.URL, Method: "GET", Path: "/"})
	if err != nil {
		t.Fatal(err)
	}
	_, err = io.ReadAll(resp)
	resp.Close()
	if !errors.Is(err, ErrReadTimeout) {
		t.Errorf("%v", err)
	}
	up := roundTripFunc(func(*http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: 101, Status: "101 Switching Protocols", Header: http.Header{}, Body: io.NopCloser(bytes.NewReader(nil))}, nil
	})
	if _, err = newClient(up, &s).Do(context.Background(), Request{Origin: "https://t.example", Method: "GET", Path: "/"}); !errors.Is(err, ErrSwitchingProtocols) {
		t.Errorf("%v", err)
	}
}

func TestDebugNeverLogsQueryOrBody(t *testing.T) {
	var dbg bytes.Buffer
	rt := roundTripFunc(func(*http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: 200, Status: "200 OK", Header: http.Header{"Set-Cookie": {"c=SECRET"}}, Body: io.NopCloser(strings.NewReader("BODYSECRET"))}, nil
	})
	var s []time.Duration
	c := newClient(rt, &s)
	c.Debug = &dbg
	resp, err := c.Do(context.Background(), Request{Origin: "https://t.example", Method: "POST", Path: "/p?token=QSECRET", Body: []byte("REQSECRET"), ContentType: "text/plain", AFSAuth: "Bearer AFSSECRET", Key: "KEYSECRET"})
	if err != nil {
		t.Fatal(err)
	}
	resp.Close()
	c.Debugf("origin: %s", "t.example")
	for _, bad := range []string{"QSECRET", "BODYSECRET", "REQSECRET", "AFSSECRET", "KEYSECRET", "c=SECRET"} {
		if strings.Contains(dbg.String(), bad) {
			t.Errorf("debug has %s: %s", bad, dbg.String())
		}
	}
	if !strings.Contains(dbg.String(), "DEBUG collab request: POST t.example/p attempt=1") {
		t.Error(dbg.String())
	}
}

// FR-R1
func TestDoRetryAfterOverCapNotRetried(t *testing.T) {
	n := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		n++
		w.Header().Set("Retry-After", "9999")
		w.WriteHeader(503)
	}))
	defer srv.Close()
	var s []time.Duration
	resp, err := newClient(http.DefaultTransport, &s).Do(context.Background(), Request{Origin: srv.URL, Method: "GET", Path: "/"})
	if err != nil || resp.Status != 503 || n != 1 || len(s) != 0 {
		t.Fatalf("%v %v n=%d slept=%v", resp, err, n, s)
	}
	resp.Close()
}

// FR-R4
func TestReadKeyDirectory(t *testing.T) {
	d := t.TempDir()
	_, err := ReadKey(d, nil)
	if err == nil || !strings.Contains(err.Error(), d) {
		t.Errorf("%v", err)
	}
}
