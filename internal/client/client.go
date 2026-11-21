package client

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/redact"
)

// Credential is the Authorization header value and where it came from (for --debug).
type Credential struct {
	Header string // "Bearer <token>"
	Source string // for example "env AFS_API_KEY"
}

// CredentialFunc returns the credential for one request; it is called once per Do, so a
// projected token file is re-read on every command.
type CredentialFunc func(ctx context.Context) (Credential, error)

// Config configures a Client.
type Config struct {
	Endpoint    string // host root, https://afs.example.com
	HTTP        http.RoundTripper
	Credentials CredentialFunc
	ReadTimeout time.Duration // 0 = none
	Stderr      io.Writer     // notices
	Debug       io.Writer     // nil = off
	UserAgent   string
	// ContractVersion is the contract this build was generated for (info.version).
	ContractVersion string
	Now             func() time.Time
	Sleep           func(ctx context.Context, d time.Duration) error
	Rand            func() float64
}

// Client sends requests to afs-api.
type Client struct {
	cfg   Config
	dbg   *debugLog
	mu    sync.Mutex
	noted bool
}

// New builds a client; missing hooks get real ones.
func New(cfg Config) *Client {
	if cfg.HTTP == nil {
		cfg.HTTP = NewTransport(TransportOptions{ConnectTimeout: DefaultConnectTimeout})
	}
	if cfg.Now == nil {
		cfg.Now = time.Now
	}
	if cfg.Sleep == nil {
		cfg.Sleep = func(ctx context.Context, d time.Duration) error {
			t := time.NewTimer(d)
			defer t.Stop()
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-t.C:
				return nil
			}
		}
	}
	if cfg.Rand == nil {
		cfg.Rand = defaultRand
	}
	c := &Client{cfg: cfg}
	if cfg.Debug != nil {
		c.dbg = &debugLog{w: cfg.Debug, now: cfg.Now}
	}
	return c
}

// Do sends the request, retrying per CLI-11. A 2xx answer is returned; any other status is an
// *APIError; transport failures are returned as they are (context.Canceled is a Ctrl-C).
func (c *Client) Do(ctx context.Context, req Request) (*Response, error) {
	key := req.IdempotencyKey
	if req.Idempotent && key == "" {
		key = NewIdempotencyKey()
	}
	var cred Credential
	if !req.NoAuth && c.cfg.Credentials != nil {
		var err error
		if cred, err = c.cfg.Credentials(ctx); err != nil {
			return nil, err
		}
	}
	if cred.Source != "" {
		c.dbg.printf("auth", "credential source: %s", cred.Source)
	}
	var last error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		resp, fail := c.once(ctx, &req, key, cred, attempt)
		if fail.apiErr == nil && fail.err == nil {
			resp.Attempts = attempt
			return resp, nil
		}
		last = fail.err
		if fail.apiErr != nil {
			last = fail.apiErr
		}
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		if attempt == maxAttempts || !retryable(&req, key, fail) {
			return nil, last
		}
		var ra time.Duration
		if fail.apiErr != nil {
			ra = parseRetryAfter(fail.apiErr.Header, c.cfg.Now())
		}
		d := backoff(attempt, c.cfg.Rand, ra)
		c.dbg.printf("client", "retrying in %s (attempt %d of %d): %v", d, attempt+1, maxAttempts, last)
		if err := c.cfg.Sleep(ctx, d); err != nil {
			return nil, err
		}
	}
	return nil, last
}

func (c *Client) buildURL(req *Request) (string, error) {
	base, err := url.Parse(strings.TrimRight(c.cfg.Endpoint, "/"))
	if err != nil || base.Host == "" {
		return "", fmt.Errorf("invalid endpoint URL %q", c.cfg.Endpoint)
	}
	u := base.String() + req.Path
	if len(req.Query) > 0 {
		u += "?" + req.Query.Encode()
	}
	return u, nil
}

func (c *Client) once(ctx context.Context, req *Request, key string, cred Credential, attempt int) (*Response, attemptFailure) {
	u, err := c.buildURL(req)
	if err != nil {
		return nil, attemptFailure{err: err}
	}
	actx, cancel := context.WithCancelCause(ctx)
	defer cancel(nil)
	var timer *time.Timer
	arm := func() {}
	stop := func() {}
	if c.cfg.ReadTimeout > 0 {
		timer = time.AfterFunc(c.cfg.ReadTimeout, func() { cancel(ErrReadTimeout) })
		arm = func() { timer.Reset(c.cfg.ReadTimeout) }
		stop = func() { timer.Stop() }
	}
	defer stop()

	var wrote, hooked bool
	trace := &httptrace.ClientTrace{
		GotConn:      func(httptrace.GotConnInfo) { hooked = true },
		WroteHeaders: func() { wrote = true },
	}
	actx = httptrace.WithClientTrace(actx, trace)
	var body io.Reader
	if len(req.Body) > 0 {
		body = bytes.NewReader(req.Body)
	}
	hr, err := http.NewRequestWithContext(actx, req.Method, u, body)
	if err != nil {
		return nil, attemptFailure{err: err}
	}
	accept := req.Accept
	if accept == "" {
		accept = "application/json"
	}
	hr.Header.Set("Accept", accept)
	hr.Header.Set("User-Agent", c.cfg.UserAgent)
	if len(req.Body) > 0 {
		hr.Header.Set("Content-Type", "application/json")
	}
	if cred.Header != "" && !req.NoAuth {
		hr.Header.Set("Authorization", cred.Header)
	}
	if req.IfMatch != "" {
		hr.Header.Set("If-Match", `"`+strings.Trim(req.IfMatch, `"`)+`"`)
	}
	if key != "" {
		hr.Header.Set("Idempotency-Key", key)
	}
	start := c.cfg.Now()
	c.dbg.printf("client", "request: %s %s operation=%s attempt=%d", req.Method, u, req.Operation, attempt)
	if c.dbg != nil {
		for _, l := range redact.Headers(hr.Header) {
			c.dbg.printf("client", "request header: %s", l)
		}
		c.dbg.body("client", "request body", req.Body)
	}

	resp, err := c.cfg.HTTP.RoundTrip(hr)
	if err != nil {
		if cause := context.Cause(actx); cause != nil && errors.Is(cause, ErrReadTimeout) {
			err = ErrReadTimeout
		}
		sent := wrote || (!hooked && !isDialError(err))
		c.dbg.printf("client", "request failed: %v attempt=%d", err, attempt)
		return nil, attemptFailure{err: err, sent: sent}
	}
	defer func() { _ = resp.Body.Close() }()
	arm()
	data, rerr := readAll(resp.Body, arm)
	if rerr != nil {
		if cause := context.Cause(actx); errors.Is(cause, ErrReadTimeout) {
			rerr = ErrReadTimeout
		}
		return nil, attemptFailure{err: rerr, sent: true}
	}
	c.dbg.printf("client", "response: %d in %s request-id=%s attempt=%d", resp.StatusCode, c.cfg.Now().Sub(start), resp.Header.Get("X-Request-Id"), attempt)
	c.dbg.body("client", "response body", data)
	c.noteVersion(resp.Header.Get("X-Afs-Contract-Version"))

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, attemptFailure{apiErr: decodeError(resp.StatusCode, resp.Header, data), sent: true}
	}
	if len(data) > 0 {
		want, got := mediaType(accept), mediaType(resp.Header.Get("Content-Type"))
		if got != want && (want != "application/json" || !strings.HasSuffix(got, "+json")) {
			return nil, attemptFailure{err: fmt.Errorf("%w: wanted %s, got %q", ErrUnexpectedMediaType, want, got), sent: true}
		}
	}
	return &Response{Status: resp.StatusCode, Header: resp.Header, Body: data, Replayed: resp.Header.Get("Idempotent-Replayed") == "true"}, attemptFailure{}
}

// readAll reads the body, calling tick after every read (the read timeout is per read).
func readAll(r io.Reader, tick func()) ([]byte, error) {
	var buf bytes.Buffer
	chunk := make([]byte, 32*1024)
	for {
		n, err := r.Read(chunk)
		buf.Write(chunk[:n])
		tick()
		if errors.Is(err, io.EOF) {
			return buf.Bytes(), nil
		}
		if err != nil {
			return nil, err
		}
	}
}

func (c *Client) noteVersion(server string) {
	if server == "" || c.cfg.Stderr == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.noted {
		return
	}
	if msg := contractNotice(c.cfg.ContractVersion, server); msg != "" {
		c.noted = true
		_, _ = fmt.Fprintln(c.cfg.Stderr, msg)
	}
}
