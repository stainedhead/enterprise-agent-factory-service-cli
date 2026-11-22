package client

import (
	"context"
	"errors"
	"net"
	"net/http"
	"strconv"
	"time"
)

// maxAttempts is the total number of tries (CLI-11).
const maxAttempts = 3

const baseDelay = 200 * time.Millisecond

// maxRetryAfter caps a server-requested delay (FR-011, D-R2); the read timeout lowers it further.
const maxRetryAfter = 20 * time.Second

// retryAfterCap is the longest Retry-After the client waits for.
func retryAfterCap(readTimeout time.Duration) time.Duration {
	if readTimeout > 0 && readTimeout < maxRetryAfter {
		return readTimeout
	}
	return maxRetryAfter
}

// backoff is full jitter: rand(0, 200ms * 2^(attempt-1)); a longer Retry-After wins (D-C24).
func backoff(attempt int, rnd func() float64, retryAfter time.Duration) time.Duration {
	ceiling := baseDelay << (attempt - 1)
	d := time.Duration(rnd() * float64(ceiling))
	if retryAfter > d {
		return retryAfter
	}
	return d
}

func parseRetryAfter(h http.Header, now time.Time) time.Duration {
	v := h.Get("Retry-After")
	if v == "" {
		return 0
	}
	if n, err := strconv.Atoi(v); err == nil && n > 0 {
		return time.Duration(n) * time.Second
	}
	if t, err := http.ParseTime(v); err == nil && t.After(now) {
		return t.Sub(now)
	}
	return 0
}

// attemptFailure describes why one try failed.
type attemptFailure struct {
	apiErr *APIError
	err    error
	// sent is true when the request may have reached the server.
	sent bool
}

// retryable applies the CLI-11 rules.
func retryable(req *Request, key string, f attemptFailure) bool {
	safe := req.Method == http.MethodGet || key != ""
	if f.apiErr != nil {
		switch f.apiErr.Status {
		case http.StatusTooManyRequests:
			return true
		case http.StatusServiceUnavailable:
			return safe
		case http.StatusConflict:
			return f.apiErr.Reason() == "idempotency_in_progress"
		}
		return false
	}
	if f.err == nil || errors.Is(f.err, context.Canceled) || errors.Is(f.err, ErrReadTimeout) || errors.Is(f.err, ErrUnexpectedMediaType) {
		return false
	}
	if !f.sent {
		return true // connect failures are retried on any request
	}
	return safe
}

// isDialError reports a failure to establish the connection.
func isDialError(err error) bool {
	var op *net.OpError
	return errors.As(err, &op) && op.Op == "dial"
}

// RetryAfter returns the delay a response asks for with Retry-After, capped at 20 s or the read timeout
// when that is smaller. Zero means no usable header.
func RetryAfter(h http.Header, now time.Time, readTimeout time.Duration) time.Duration {
	d := parseRetryAfter(h, now)
	if limit := retryAfterCap(readTimeout); d > limit {
		return limit
	}
	return d
}
