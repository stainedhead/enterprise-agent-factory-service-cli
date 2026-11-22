package edge

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/client"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/redact"
)

// Request is one call to the edge. Credentials are header values supplied by the caller.
type Request struct {
	Origin      string // validated by Origin
	Method      string
	Path        string // starts with "/", may carry a query string
	ContentType string // used when Body is not empty
	Body        []byte
	AFSAuth     string // value of X-Afs-Authorization ("Bearer ...")
	AFSSource   string // where the AFS credential came from (debug only)
	Key         string // the Paperclip key, sent as "Authorization: Bearer <key>"
	UserAgent   string
}

// Client sends requests. HTTP is used with RoundTrip directly, so no redirect is followed and no cookie is
// stored or sent.
type Client struct {
	HTTP        http.RoundTripper
	ReadTimeout time.Duration
	Now         func() time.Time
	Sleep       func(ctx context.Context, d time.Duration) error
	Rand        func() float64
	Debug       io.Writer
}

// ErrReadTimeout is returned when the edge sends nothing for the read timeout.
var ErrReadTimeout = client.ErrReadTimeout

// ErrSwitchingProtocols is a 101 answer: WebSocket is not supported (CLI-US-25.9).
var ErrSwitchingProtocols = errors.New("the server answered 101 Switching Protocols; WebSocket and streaming upgrades are not supported by invoke-api")

// Response is an answer whose body is read lazily; Close it.
type Response struct {
	Status   int
	Reason   string
	Proto    string
	Header   http.Header
	Attempts int
	body     io.ReadCloser
	rd       *idleReader
	cancel   context.CancelCauseFunc
}

// Read streams the body; every read re-arms the idle timeout.
func (r *Response) Read(p []byte) (int, error) { return r.rd.Read(p) }

// Close releases the connection.
func (r *Response) Close() {
	if r.rd != nil {
		r.rd.stop()
	}
	if r.body != nil {
		_ = r.body.Close()
	}
	r.cancel(nil)
}

type idleReader struct {
	r     io.Reader
	timer *time.Timer
	d     time.Duration
	ctx   context.Context
}

func (i *idleReader) Read(p []byte) (int, error) {
	n, err := i.r.Read(p)
	if i.timer != nil {
		i.timer.Reset(i.d)
	}
	if err != nil && !errors.Is(err, io.EOF) && errors.Is(context.Cause(i.ctx), ErrReadTimeout) {
		err = ErrReadTimeout
	}
	return n, err
}

func (i *idleReader) stop() {
	if i.timer != nil {
		i.timer.Stop()
	}
}

// Do sends the request. A response of any status is returned (the caller classifies it); GET and HEAD retry
// 429 and 503, honouring Retry-After, and every method retries a failure to connect; never more than
// client.MaxAttempts tries, and no Idempotency-Key is sent (CLI-US-25.8).
func (c *Client) Do(ctx context.Context, req Request) (*Response, error) {
	if c.Now == nil {
		c.Now = time.Now
	}
	if c.Rand == nil {
		c.Rand = func() float64 { return 0.5 }
	}
	idempotent := req.Method == http.MethodGet || req.Method == http.MethodHead
	var last error
	for attempt := 1; attempt <= client.MaxAttempts; attempt++ {
		resp, err := c.once(ctx, req, attempt)
		var delay time.Duration
		switch {
		case err != nil:
			if ctx.Err() != nil {
				return nil, ctx.Err()
			}
			last = err
			if !client.IsDialError(err) || attempt == client.MaxAttempts {
				return nil, err
			}
		case idempotent && (resp.Status == http.StatusTooManyRequests || resp.Status == http.StatusServiceUnavailable) && attempt < client.MaxAttempts:
			delay = client.RetryAfter(resp.Header, c.Now(), c.ReadTimeout) // capped at 20 s or the read timeout
			resp.Close()
			last = fmt.Errorf("HTTP %d", resp.Status)
		default:
			resp.Attempts = attempt
			return resp, nil
		}
		d := client.Backoff(attempt, c.Rand, delay)
		c.debugf("retrying in %s (attempt %d of %d): %v", d, attempt+1, client.MaxAttempts, last)
		if c.Sleep != nil {
			if err := c.Sleep(ctx, d); err != nil {
				return nil, err
			}
		}
	}
	return nil, last
}

func (c *Client) once(ctx context.Context, req Request, attempt int) (*Response, error) {
	actx, cancel := context.WithCancelCause(ctx)
	var timer *time.Timer
	if c.ReadTimeout > 0 {
		timer = time.AfterFunc(c.ReadTimeout, func() { cancel(ErrReadTimeout) })
	}
	fail := func(err error) (*Response, error) {
		if timer != nil {
			timer.Stop()
		}
		if errors.Is(context.Cause(actx), ErrReadTimeout) {
			err = ErrReadTimeout
		}
		cancel(nil)
		return nil, err
	}
	var body io.Reader
	if len(req.Body) > 0 {
		body = bytes.NewReader(req.Body)
	}
	hr, err := http.NewRequestWithContext(actx, req.Method, req.Origin+req.Path, body)
	if err != nil {
		return fail(err)
	}
	hr.Header.Set("Accept", "application/json")
	hr.Header.Set("Accept-Encoding", "identity")
	hr.Header.Set("User-Agent", req.UserAgent)
	hr.Header.Set("X-Afs-Authorization", req.AFSAuth)
	hr.Header.Set("Authorization", "Bearer "+req.Key)
	if len(req.Body) > 0 {
		hr.Header.Set("Content-Type", req.ContentType)
	}
	start := c.Now()
	c.debugf("request: %s %s%s attempt=%d", req.Method, hostOf(req.Origin), stripQuery(req.Path), attempt)
	if c.Debug != nil {
		for _, l := range redact.Headers(hr.Header) {
			c.debugf("request header: %s", l)
		}
	}
	resp, err := c.HTTP.RoundTrip(hr)
	if err != nil {
		c.debugf("request failed: %v attempt=%d", err, attempt)
		return fail(err)
	}
	c.debugf("response: %d in %s request-id=%s attempt=%d", resp.StatusCode, c.Now().Sub(start), resp.Header.Get("X-Request-Id"), attempt)
	if c.Debug != nil {
		for _, l := range redact.Headers(resp.Header) {
			c.debugf("response header: %s", l)
		}
	}
	if resp.StatusCode == http.StatusSwitchingProtocols {
		_ = resp.Body.Close()
		return fail(ErrSwitchingProtocols)
	}
	reason := strings.TrimSpace(strings.TrimPrefix(resp.Status, fmt.Sprintf("%d", resp.StatusCode)))
	if reason == "" {
		reason = http.StatusText(resp.StatusCode)
	}
	if timer != nil {
		timer.Reset(c.ReadTimeout)
	}
	rd := &idleReader{r: resp.Body, timer: timer, d: c.ReadTimeout, ctx: actx}
	return &Response{Status: resp.StatusCode, Reason: reason, Proto: resp.Proto, Header: resp.Header, body: resp.Body, rd: rd, cancel: cancel}, nil
}

func hostOf(origin string) string {
	_, h, _ := strings.Cut(origin, "://")
	return h
}

// stripQuery drops the query string: it can carry a token and is never logged.
func stripQuery(p string) string {
	if i := strings.IndexByte(p, '?'); i >= 0 {
		return p[:i]
	}
	return p
}

func (c *Client) debugf(format string, a ...any) {
	if c.Debug == nil {
		return
	}
	_, _ = fmt.Fprintf(c.Debug, "%s DEBUG collab %s\n", c.Now().UTC().Format(time.RFC3339), fmt.Sprintf(format, a...))
}

// Debugf lets the caller write a line in the same format (origin read, credential source).
func (c *Client) Debugf(format string, a ...any) {
	if c.Now == nil {
		c.Now = time.Now
	}
	c.debugf(format, a...)
}
