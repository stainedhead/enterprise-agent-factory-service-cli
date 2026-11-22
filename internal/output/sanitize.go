package output

import (
	"fmt"
	"io"
	"strings"
	"unicode/utf8"
)

// Sanitize makes untrusted text safe to print on a terminal: C0 control characters other than newline
// and tab, DEL, C1 controls (U+0080 to U+009F) and invalid UTF-8 bytes are replaced by a visible
// escape (\x1b, \u0085, \xff), so a value stays diagnosable and cannot move the cursor, clear the
// screen, set the title or write to the clipboard (FR-003).
func Sanitize(s string) string {
	clean := true
	for i := 0; i < len(s); i++ {
		if c := s[i]; c < 0x20 && c != '\n' && c != '\t' || c >= 0x7f {
			clean = false
			break
		}
	}
	if clean {
		return s
	}
	var b strings.Builder
	b.Grow(len(s) + 8)
	for i := 0; i < len(s); {
		r, n := utf8.DecodeRuneInString(s[i:])
		switch {
		case r == utf8.RuneError && n == 1:
			fmt.Fprintf(&b, `\x%02x`, s[i])
		case r < 0x20 && r != '\n' && r != '\t', r == 0x7f:
			fmt.Fprintf(&b, `\x%02x`, r)
		case r >= 0x80 && r <= 0x9f:
			fmt.Fprintf(&b, `\u%04x`, r)
		default:
			b.WriteString(s[i : i+n])
		}
		i += n
	}
	return b.String()
}

type sanitizingWriter struct{ w io.Writer }

// NewSanitizingWriter wraps w so everything written through it is passed through Sanitize. Each Write
// must hold whole characters (fmt.Fprintf calls do).
func NewSanitizingWriter(w io.Writer) io.Writer { return sanitizingWriter{w} }

func (s sanitizingWriter) Write(p []byte) (int, error) {
	if _, err := io.WriteString(s.w, Sanitize(string(p))); err != nil {
		return 0, err
	}
	return len(p), nil
}
