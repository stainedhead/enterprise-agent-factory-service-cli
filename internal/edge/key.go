package edge

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/config"
)

// MaxKeyBytes is the longest Paperclip key accepted.
const MaxKeyBytes = 4096

// KeyShapeError is a key that is empty, too long or has whitespace or control characters (exit 252).
type KeyShapeError struct{ Reason string }

func (e *KeyShapeError) Error() string { return "the Paperclip key " + e.Reason }

// ReadKey reads the Paperclip key from a file or, for "-", from stdin. A file wider than 0600 is refused
// with the chmod hint (an ordinary error, exit 255); a key of the wrong shape is a *KeyShapeError. The key
// is returned trimmed of trailing whitespace and is never part of any error message.
func ReadKey(src string, stdin io.Reader) (string, error) {
	var raw []byte
	var err error
	if src == "-" {
		if stdin == nil {
			return "", errors.New("no standard input to read the Paperclip key from")
		}
		raw, err = io.ReadAll(io.LimitReader(stdin, MaxKeyBytes+64))
	} else {
		if err = config.CheckPrivate(src); err != nil {
			return "", err
		}
		var f *os.File
		if f, err = os.Open(src); err == nil { //nolint:gosec // user-selected key file
			defer func() { _ = f.Close() }()
			raw, err = io.ReadAll(io.LimitReader(f, MaxKeyBytes+64))
		}
	}
	if err != nil {
		return "", fmt.Errorf("cannot read the Paperclip key from %s: %w", describe(src), err)
	}
	key := strings.TrimRightFunc(string(raw), unicode.IsSpace)
	switch {
	case key == "":
		return "", &KeyShapeError{Reason: "is empty"}
	case len(key) > MaxKeyBytes:
		return "", &KeyShapeError{Reason: fmt.Sprintf("is longer than %d bytes", MaxKeyBytes)}
	}
	for _, r := range key {
		if unicode.IsSpace(r) || unicode.IsControl(r) {
			return "", &KeyShapeError{Reason: "must not contain whitespace or control characters"}
		}
	}
	return key, nil
}

func describe(src string) string {
	if src == "-" {
		return "standard input"
	}
	return src
}
