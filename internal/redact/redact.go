// Package redact removes secrets from anything the CLI writes to a log: header values and
// JSON body fields (data-dictionary section 8).
package redact

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
)

// Mask replaces a redacted value.
const Mask = "[redacted]"

// sensitiveHeaders are masked wherever headers are logged. x-afs-authorization carries the AFS
// credential to the team's collaboration edge (CLI-25); the edge cookie __Host-afs-collab travels
// in Cookie and Set-Cookie, so it is masked with them.
var sensitiveHeaders = map[string]bool{
	"authorization": true, "cookie": true, "set-cookie": true, "proxy-authorization": true,
	"x-afs-authorization": true,
}

// sensitiveFields are JSON member names whose values are never logged, at any depth.
var sensitiveFields = map[string]bool{
	"value": true, "secretKey": true, "token": true, "refreshToken": true, "accessToken": true,
	"clientSecret": true, "password": true, "idToken": true, "deviceCode": true,
}

// Fields returns the redacted JSON member names (for tests and docs).
func Fields() []string {
	out := make([]string, 0, len(sensitiveFields))
	for k := range sensitiveFields {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// HeaderValue returns the value to log for a header.
func HeaderValue(name, value string) string {
	if sensitiveHeaders[strings.ToLower(name)] {
		return Mask
	}
	return value
}

// Headers renders headers as sorted "Name: value" lines with sensitive values masked.
func Headers(h http.Header) []string {
	names := make([]string, 0, len(h))
	for k := range h {
		names = append(names, k)
	}
	sort.Strings(names)
	var out []string
	for _, k := range names {
		for _, v := range h[k] {
			out = append(out, fmt.Sprintf("%s: %s", k, HeaderValue(k, v)))
		}
	}
	return out
}

// JSON masks sensitive members of a JSON document at any depth, inside arrays too. A body that
// is not JSON is never logged: a placeholder with its size is returned.
func JSON(body []byte) string {
	if len(bytes.TrimSpace(body)) == 0 {
		return ""
	}
	dec := json.NewDecoder(bytes.NewReader(body))
	dec.UseNumber()
	var v any
	if err := dec.Decode(&v); err != nil {
		return fmt.Sprintf("<%d bytes, not JSON>", len(body))
	}
	out, err := json.Marshal(mask(v))
	if err != nil {
		return fmt.Sprintf("<%d bytes, not JSON>", len(body))
	}
	return string(out)
}

// secretValue is true for a member whose name ends in "Secret" and whose value is not a reference:
// collabClientSecretArn names a Secrets Manager secret and is safe to log, a "...Secret" that holds a
// value is not (CLI-US-31.2). Members whose name ends in "Arn" are references by definition.
func secretValue(name string, v any) bool {
	lower := strings.ToLower(name)
	if strings.HasSuffix(lower, "arn") || !strings.HasSuffix(lower, "secret") {
		return false
	}
	s, ok := v.(string)
	if ok && strings.HasPrefix(s, "arn:") {
		return false
	}
	return true
}

func mask(v any) any {
	switch x := v.(type) {
	case map[string]any:
		for k, val := range x {
			if sensitiveFields[k] || secretValue(k, val) {
				x[k] = Mask
				continue
			}
			x[k] = mask(val)
		}
		return x
	case []any:
		for i := range x {
			x[i] = mask(x[i])
		}
		return x
	}
	return v
}
