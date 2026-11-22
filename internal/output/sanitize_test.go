package output

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"unicode/utf8"
)

const hostile = "\x1b[2J\x1b]52;c;ZXZpbA==\x07\r\x7f\u0085\u009b2J\x00end"

func TestSanitize(t *testing.T) {
	cases := []struct{ in, want string }{
		{"plain text é 日本", "plain text é 日本"},
		{"a\nb\tc", "a\nb\tc"},
		{"\x1b[2J", `\x1b[2J`},
		{"x\x07y", `x\x07y`},
		{"a\rb", `a\x0db`},
		{"d\x7f", `d\x7f`},
		{"c1\u0085\u009b", `c1\u0085\u009b`},
		{"bad\xffutf8", `bad\xffutf8`},
		{"", ""},
	}
	for _, c := range cases {
		if got := Sanitize(c.in); got != c.want {
			t.Errorf("Sanitize(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func cleanOutput(t *testing.T, s string) {
	t.Helper()
	for _, r := range s {
		if (r < 0x20 && r != '\n' && r != '\t') || (r >= 0x7f && r <= 0x9f) {
			t.Fatalf("control character %U in %q", r, s)
		}
	}
	if !utf8.ValidString(s) {
		t.Fatalf("invalid UTF-8 in %q", s)
	}
}

func rawJSON(t *testing.T, v any) string {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func TestTableAndTextNeverEmitControlCharacters(t *testing.T) {
	cols := []Column{{Header: "Id", Path: "id"}, {Header: "Name", Path: "name"}}
	docs := map[string]string{
		"list": rawJSON(t, map[string]any{"items": []any{
			map[string]any{"id": "team_1", "name": hostile},
			map[string]any{"id": "team_2", "name": "ok"},
		}}),
		"object":  rawJSON(t, map[string]any{"id": "x", "name": hostile}),
		"scalar":  rawJSON(t, hostile),
		"scalars": rawJSON(t, []any{hostile, "x"}),
	}
	for doc, raw := range docs {
		for _, mode := range []string{"table", "text"} {
			out := render(t, raw, RenderOptions{Format: mode, Options: Options{Columns: cols, ItemsKey: "items", Width: 120}})
			cleanOutput(t, out)
			if doc == "list" && !strings.Contains(out, `\x1b[2J`) {
				t.Errorf("%s/%s lost the visible escape: %q", doc, mode, out)
			}
		}
	}
}

// ADR-005: json and yaml already escape control characters, so they stay as they are.
func TestJSONAndYAMLAreUnchangedByTheSanitizer(t *testing.T) {
	raw := rawJSON(t, map[string]any{"name": hostile})
	j := render(t, raw, RenderOptions{Format: "json"})
	y := render(t, raw, RenderOptions{Format: "yaml"})
	cleanOutput(t, y)
	// json passes the service's bytes through (ADR-005): C0 controls arrive escaped by the service's
	// encoder; DEL and C1 are not escaped by JSON encoders and stay as sent (D-R7).
	for _, r := range j {
		if r < 0x20 && r != '\n' {
			t.Fatalf("raw C0 character %U in json output", r)
		}
	}
	var back map[string]any
	if err := json.Unmarshal([]byte(j), &back); err != nil || back["name"] != hostile {
		t.Errorf("json round trip: %v", err)
	}
}

func TestSanitizingWriter(t *testing.T) {
	var b bytes.Buffer
	w := NewSanitizingWriter(&b)
	n, err := w.Write([]byte("a\x1bb\n"))
	if err != nil || n != 4 {
		t.Fatalf("n=%d err=%v", n, err)
	}
	if b.String() != "a\\x1bb\n" {
		t.Errorf("got %q", b.String())
	}
}

func BenchmarkTable10kRows(b *testing.B) {
	rows := make([]any, 10000)
	for i := range rows {
		rows[i] = map[string]any{"id": "team_x", "name": "a reasonably long team name"}
	}
	raw, _ := json.Marshal(map[string]any{"items": rows})
	o := RenderOptions{Format: "table", Options: Options{ItemsKey: "items", Width: 120, Columns: []Column{{Header: "Id", Path: "id"}, {Header: "Name", Path: "name"}}}}
	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		_ = Render(&buf, raw, o)
	}
}
