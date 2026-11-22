// Package secretinput reads secret values only from a file, standard input or a no-echo
// prompt, never from a plain argument (CLI-9, D-C16).
package secretinput

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"
)

// MaxSecretBytes is the largest value the API accepts (SecretValue maxLength).
const MaxSecretBytes = 65536

// ErrPlainArgument refuses a secret given directly on the command line.
var ErrPlainArgument = errors.New("a secret value must not be passed as an argument: use file://PATH, fileb://PATH or - (standard input)")

// Reader supplies the files and stdin to read from.
type Reader struct {
	Stdin    io.Reader
	ReadFile func(string) ([]byte, error)
	// ReadPassword reads one line from the terminal without echo; nil when there is no terminal.
	ReadPassword func() ([]byte, error)
}

// ErrNoValue is returned by Prompt when no terminal is available to ask on.
var ErrNoValue = errors.New("no secret value given and no terminal to prompt on: use file://PATH, fileb://PATH or - (standard input)")

// Prompt asks for the value on the terminal without echo (CLI-9); the same limits apply as for files.
func (r Reader) Prompt() ([]byte, error) {
	if r.ReadPassword == nil {
		return nil, ErrNoValue
	}
	b, err := r.ReadPassword()
	if err != nil {
		return nil, fmt.Errorf("cannot read the secret value: %w", err)
	}
	return check(b)
}

// Read returns the secret named by arg: file://PATH or fileb://PATH read the file, "-" reads stdin.
// The bytes are returned unchanged (no newline trimming); they must be valid UTF-8 of 1 to 65536
// bytes. Anything else is refused.
func (r Reader) Read(arg string) ([]byte, error) {
	var b []byte
	var err error
	switch {
	case arg == "-":
		b, err = io.ReadAll(io.LimitReader(r.Stdin, MaxSecretBytes+1))
	case strings.HasPrefix(arg, "file://"), strings.HasPrefix(arg, "fileb://"):
		path := arg[strings.Index(arg, "://")+3:]
		rf := r.ReadFile
		if rf == nil {
			rf = os.ReadFile
		}
		b, err = rf(path)
	default:
		return nil, ErrPlainArgument
	}
	if err != nil {
		return nil, fmt.Errorf("cannot read the secret value: %w", err)
	}
	return check(b)
}

// check applies the limits of the SecretValue schema.
func check(b []byte) ([]byte, error) {
	switch {
	case len(b) == 0:
		return nil, errors.New("the secret value is empty")
	case len(b) > MaxSecretBytes:
		return nil, fmt.Errorf("the secret value is larger than %d bytes", MaxSecretBytes)
	case !utf8.Valid(b):
		return nil, errors.New("the secret value is not valid UTF-8 text")
	}
	return b, nil
}
