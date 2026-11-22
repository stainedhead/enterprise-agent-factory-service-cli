package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// APIError is a non-2xx answer from afs-api, after retries.
type APIError struct {
	Status    int
	Code      string
	Message   string
	RequestID string
	Details   json.RawMessage
	Header    http.Header
	Body      []byte
}

func (e *APIError) Error() string { return fmt.Sprintf("%s: %s", e.Code, e.Message) }

// detail reads one string member of details.
func (e *APIError) detail(key string) string {
	var d map[string]json.RawMessage
	if json.Unmarshal(e.Details, &d) != nil {
		return ""
	}
	var s string
	if json.Unmarshal(d[key], &s) != nil {
		return ""
	}
	return s
}

// detailList reads one string-array member of details.
func (e *APIError) detailList(key string) []string {
	var d map[string]json.RawMessage
	if json.Unmarshal(e.Details, &d) != nil {
		return nil
	}
	var l []string
	if json.Unmarshal(d[key], &l) != nil {
		return nil
	}
	return l
}

// Dependents is details.dependents: ids of the objects that block a delete.
func (e *APIError) Dependents() []string { return e.detailList("dependents") }

// MissingTags is details.missingTags of a ValidationError.
func (e *APIError) MissingTags() []string { return e.detailList("missingTags") }

// Reason is details.reason (for example idempotency_in_progress).
func (e *APIError) Reason() string { return e.detail("reason") }

// RequiredPermission is details.requiredPermission of an AccessDenied.
func (e *APIError) RequiredPermission() string { return e.detail("requiredPermission") }

// ErrUnexpectedMediaType is returned for a 2xx answer whose media type is not the one requested.
var ErrUnexpectedMediaType = errors.New("unexpected response media type")

// ErrReadTimeout is the cause when the server does not answer within the read timeout.
var ErrReadTimeout = errors.New("read timeout: the server did not respond in time")

func decodeError(status int, h http.Header, body []byte) *APIError {
	e := &APIError{Status: status, Header: h, Body: body, RequestID: h.Get("X-Request-Id")}
	var env struct {
		Code      string          `json:"code"`
		Message   string          `json:"message"`
		RequestID string          `json:"requestId"`
		Details   json.RawMessage `json:"details"`
	}
	if mt := mediaType(h.Get("Content-Type")); strings.Contains(mt, "json") && json.Unmarshal(body, &env) == nil && env.Code != "" {
		e.Code, e.Message, e.Details = env.Code, env.Message, env.Details
		if env.RequestID != "" {
			e.RequestID = env.RequestID
		}
		return e
	}
	e.Code = "ServiceError"
	if status >= 500 {
		e.Code = "InternalError"
	}
	e.Message = fmt.Sprintf("the server answered HTTP %d without an error document", status)
	return e
}

func mediaType(ct string) string {
	mt, _, _ := strings.Cut(ct, ";")
	return strings.ToLower(strings.TrimSpace(mt))
}
