// Package secretinput reads secret values only from a file, standard input or (M2) a no-echo
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
