package client

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/redact"
)

// debugLog writes the --debug lines: "<RFC 3339 time> DEBUG <component> <message>".
type debugLog struct {
	w   io.Writer
	now func() time.Time
}

func (d *debugLog) printf(component, format string, a ...any) {
	if d == nil || d.w == nil {
		return
	}
	_, _ = fmt.Fprintf(d.w, "%s DEBUG %s %s\n", d.now().UTC().Format(time.RFC3339), component, fmt.Sprintf(format, a...))
}

// headers writes one line per request header with sensitive values masked (CLI-25).
func (d *debugLog) headers(component string, h http.Header) {
	for _, l := range redact.Headers(h) {
		d.printf(component, "request header: %s", l)
	}
}

const maxLoggedBody = 4096

func (d *debugLog) body(component, label string, b []byte) {
	if d == nil || d.w == nil || len(b) == 0 {
		return
	}
	s := redact.JSON(b)
	if len(s) > maxLoggedBody {
		s = s[:maxLoggedBody] + "...(truncated)"
	}
	d.printf(component, "%s: %s", label, s)
}
