package sso

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

// debugDoer logs every identity-provider request for --debug in the API client's format:
// "<RFC 3339 time> DEBUG sso <METHOD> <URL> -> <status> (<duration>)". The URL is printed without
// credentials. Request forms and response bodies are never logged: they carry the device code, user
// code, refresh token and tokens (FR-008).
type debugDoer struct {
	inner HTTPDoer
	w     io.Writer
	now   func() time.Time
}

// WithDebug wraps hc so each request is logged to w. A nil w returns hc unchanged.
func WithDebug(hc HTTPDoer, w io.Writer, now func() time.Time) HTTPDoer {
	if w == nil {
		return hc
	}
	if now == nil {
		now = time.Now
	}
	return &debugDoer{inner: hc, w: w, now: now}
}

func (d *debugDoer) Do(req *http.Request) (*http.Response, error) {
	start := d.now()
	resp, err := d.inner.Do(req)
	elapsed := d.now().Sub(start).Round(time.Millisecond)
	status := "error"
	switch {
	case err == nil:
		status = fmt.Sprint(resp.StatusCode)
	case req.Context().Err() != nil:
		status = "cancelled"
	}
	_, _ = fmt.Fprintf(d.w, "%s DEBUG sso %s %s -> %s (%s)\n", d.now().UTC().Format(time.RFC3339), req.Method, req.URL.Redacted(), status, elapsed)
	return resp, err
}
