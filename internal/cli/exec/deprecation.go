package exec

import (
	"fmt"
	"io"
	"net/http"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/output"
)

// removalVersion is the release that removes what 1.1.0 deprecated (ORC-27, ORC-30).
const removalVersion = "2.0.0"

// warnDeprecated writes one stderr line for a deprecated operation and one for each deprecated option the
// caller gave, by flag or in --cli-input-json (ORC-30, OR-CLI-27, OR-CLI-28). It reports whether any line was
// written. Nothing here changes the exit code or the output on stdout.
func warnDeprecated(w io.Writer, cmd *cli.Command, p *cli.Parsed, b *Built) bool {
	printed := false
	if cmd.Deprecated != "" {
		_, _ = fmt.Fprintf(w, "Warning: afs %s is deprecated and will be removed in %s; use %s instead.\n", cmd.Key(), removalVersion, cmd.Deprecated)
		printed = true
	}
	for _, o := range cmd.Options {
		if o.Deprecated == "" {
			continue
		}
		given := p.Has(o.Name())
		if !given && o.In == "body" && b != nil && b.Body != nil {
			_, given = b.Body.Get(o.Wire)
		}
		if given {
			_, _ = fmt.Fprintf(w, "Warning: %s is deprecated and will be removed in %s; %s\n", o.Flag, removalVersion, o.Deprecated)
			printed = true
		}
	}
	return printed
}

// warnServiceDeprecation covers what the CLI has no static knowledge of: a response that carries the
// Deprecation header (RFC 9745) prints one line with the Sunset date, when the service sent one. The text
// that is printed is sanitized because it comes from the service (FR-003).
func warnServiceDeprecation(w io.Writer, h http.Header) {
	if h == nil || h.Get("Deprecation") == "" {
		return
	}
	line := "Warning: the service reports that this request uses a deprecated feature"
	if s := h.Get("Sunset"); s != "" {
		line += "; it is planned for removal on " + s
	}
	_, _ = fmt.Fprintln(output.NewSanitizingWriter(w), line+".")
}
