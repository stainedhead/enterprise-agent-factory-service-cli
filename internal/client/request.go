// Package client is the HTTP transport to afs-api: request building, credentials, retries,
// idempotency keys, timeouts, debug logging with redaction and contract-version notices.
package client

import (
	"net/http"
	"net/url"
	"strings"
)

// Request is one API call.
type Request struct {
	Method string
	// Path is the escaped path, built with BuildPath.
	Path  string
	Query url.Values
	Body  []byte
	// ContentType is the media type of Body; "" means application/json.
	ContentType string
	Accept      string // media type wanted; "" means application/json
	// IfMatch is the version sent as If-Match: "<version>" ("" sends no header).
	IfMatch string
	// Idempotent marks a request that carries an Idempotency-Key (every POST).
	Idempotent     bool
	IdempotencyKey string // "" generates a ULID when Idempotent
	// NoAuth sends no credential (public operations).
	NoAuth bool
	// Operation names the operationId for debug logs.
	Operation string
}

// Response is a successful (2xx) reply.
type Response struct {
	Status int
	Header http.Header
	Body   []byte
	// Replayed is true when the server answered from a stored idempotent response.
	Replayed bool
	// Attempts is how many requests were made.
	Attempts int
}

// RequestID returns the X-Request-Id of the reply.
func (r *Response) RequestID() string { return r.Header.Get("X-Request-Id") }

// BuildPath fills {name} placeholders with each value escaped as one path segment.
func BuildPath(template string, params map[string]string) string {
	var b strings.Builder
	for i := 0; i < len(template); {
		if template[i] == '{' {
			if j := strings.IndexByte(template[i:], '}'); j > 0 {
				name := template[i+1 : i+j]
				b.WriteString(url.PathEscape(params[name]))
				i += j + 1
				continue
			}
		}
		b.WriteByte(template[i])
		i++
	}
	return b.String()
}
