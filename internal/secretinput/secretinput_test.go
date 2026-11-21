package secretinput

import (
	"errors"
	"strings"
	"testing"
)

func TestReadSources(t *testing.T) {
	r := Reader{Stdin: strings.NewReader("from-stdin\n"), ReadFile: func(p string) ([]byte, error) {
		if p == "/tmp/tok" {
			return []byte("from-file\n"), nil
		}
		return nil, errors.New("nope")
	}}
	for arg, want := range map[string]string{"-": "from-stdin\n", "file:///tmp/tok": "from-file\n", "fileb:///tmp/tok": "from-file\n"} {
		r.Stdin = strings.NewReader("from-stdin\n")
		got, err := r.Read(arg)
		if err != nil || string(got) != want {
			t.Errorf("%s: %q %v", arg, got, err)
		}
	}
}

// CLI-US-3: anything that is not file://, fileb:// or - is refused.
func TestPlainArgumentRefused(t *testing.T) {
	if _, err := (Reader{}).Read("hunter2"); !errors.Is(err, ErrPlainArgument) {
		t.Fatalf("err = %v", err)
	}
}

func TestReadRejectsBadValues(t *testing.T) {
	cases := map[string][]byte{"empty": {}, "binary": {0xff, 0xfe}, "big": []byte(strings.Repeat("x", MaxSecretBytes+1))}
	for name, content := range cases {
		r := Reader{ReadFile: func(string) ([]byte, error) { return content, nil }}
		if _, err := r.Read("file://x"); err == nil {
			t.Errorf("%s accepted", name)
		}
	}
	if _, err := (Reader{ReadFile: func(string) ([]byte, error) { return nil, errors.New("missing") }}).Read("file://x"); err == nil {
		t.Error("read error swallowed")
	}
	if _, err := (Reader{Stdin: strings.NewReader("")}).Read("-"); err == nil {
		t.Error("empty stdin accepted")
	}
	if _, err := (Reader{}).Read("file://" + t.TempDir() + "/nope"); err == nil {
		t.Error("missing real file accepted")
	}
}
