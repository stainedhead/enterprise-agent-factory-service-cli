package output

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/golden"
)

const fixture = `{"items":[
 {"id":"team_01JD2Q3R4S5T6V7W8X9Y0Z0001","slug":"payments","state":"active","version":3,"quotas":{"cpu":"16"},"tags":["a","b"],"note":null,"ceilingOverrideUsd":null,"status":{"phase":"Ready"}},
 {"id":"team_01JD2Q3R4S5T6V7W8X9Y0Z0002","slug":"billing","state":"suspended","version":12,"quotas":{"cpu":"4"},"tags":[],"note":"a very long note that should be truncated when the table is narrow enough to need it","status":{"phase":"Pending"}}
],"nextToken":"abc"}`

var teamCols = []Column{{"Id", "id"}, {"Slug", "slug"}, {"State", "state"}, {"Phase", "status.phase"}, {"Version", "version"}}

func render(t *testing.T, raw string, o RenderOptions) string {
	t.Helper()
	var b bytes.Buffer
	if err := Render(&b, []byte(raw), o); err != nil {
		t.Fatal(err)
	}
	return b.String()
}

// CLI-US-4.1 / D-C5: the raw body is re-indented, order, nulls, unknown fields and numbers intact.
func TestJSONPreservesOrderNullsUnknownFields(t *testing.T) {
	in := `{"z":1,"a":null,"future":{"x":[1,2.50,12345678901234567890]},"m":"é"}`
	got := render(t, in, RenderOptions{Format: "json"})
	want := "{\n  \"z\": 1,\n  \"a\": null,\n  \"future\": {\n    \"x\": [\n      1,\n      2.50,\n      12345678901234567890\n    ]\n  },\n  \"m\": \"é\"\n}\n"
	if got != want {
		t.Errorf("got:\n%s", got)
	}
	// byte-stable: rendering the output again changes nothing
	if again := render(t, got, RenderOptions{Format: "json"}); again != got {
		t.Error("re-indent is not stable")
	}
	if err := Render(io.Discard, []byte("{bad"), RenderOptions{Format: "json"}); err == nil {
		t.Error("invalid JSON accepted")
	}
}

func TestYAMLFromOrderedTree(t *testing.T) {
	got := render(t, `{"b":1,"a":{"c":[true,null,"x"],"s":"true","e":""},"f":1.5}`, RenderOptions{Format: "yaml"})
	want := "b: 1\na:\n  c:\n    - true\n    - null\n    - x\n  s: \"true\"\n  e: \"\"\nf: 1.5\n"
	if got != want {
		t.Errorf("yaml:\n%s\nwant:\n%s", got, want)
	}
}

func TestNumbersKeepPrecision(t *testing.T) {
	got := render(t, `{"n":12345678901234567890,"f":0.1000000000000000055511}`, RenderOptions{Format: "yaml"})
	if !strings.Contains(got, "12345678901234567890") || !strings.Contains(got, "0.1000000000000000055511") {
		t.Errorf("precision lost: %s", got)
	}
	q := render(t, `{"items":[{"n":12345678901234567890}]}`, RenderOptions{Format: "json", Query: "items[0].n"})
	if strings.TrimSpace(q) != `"12345678901234567890"` && strings.TrimSpace(q) != "12345678901234567890" {
		t.Errorf("query number = %s", q)
	}
}

// CLI-US-4.6: --query applies before rendering, for json, yaml, table and text.
func TestQueryBeforeRenderJSONYAMLTableText(t *testing.T) {
	q := "items[?state==`active`].slug"
	for format, want := range map[string]string{
		"json":  "[\n  \"payments\"\n]\n",
		"yaml":  "- payments\n",
		"table": "Value\npayments\n",
		"text":  "payments\n",
	} {
		if got := render(t, fixture, RenderOptions{Format: format, Query: q, Options: Options{Columns: teamCols, ItemsKey: "items"}}); got != want {
			t.Errorf("%s: got %q, want %q", format, got, want)
		}
	}
	// a query that keeps item objects renders with the operation's columns
	got := render(t, fixture, RenderOptions{Format: "table", Query: "items[?state==`active`]", Options: Options{Columns: teamCols}})
	if !strings.HasPrefix(got, "Id") || !strings.Contains(got, "payments") || strings.Contains(got, "billing") {
		t.Errorf("object query table:\n%s", got)
	}
	// numbers in the query engine
	if got := render(t, fixture, RenderOptions{Format: "json", Query: "items[?version>`5`].slug"}); got != "[\n  \"billing\"\n]\n" {
		t.Errorf("numeric filter: %s", got)
	}
}

func TestInvalidQueryExit252(t *testing.T) {
	err := Render(io.Discard, []byte(`{}`), RenderOptions{Format: "json", Query: "items[?"})
	var qe *QueryError
	if !errors.As(err, &qe) || !strings.Contains(qe.Error(), "invalid --query") || qe.Unwrap() == nil {
		t.Fatalf("err = %v", err)
	}
	if Validate("a.b") != nil || Validate("[") == nil {
		t.Error("Validate")
	}
	if err := Render(io.Discard, []byte(`not json`), RenderOptions{Format: "json", Query: "a"}); err == nil {
		t.Error("query on invalid JSON accepted")
	}
}

func TestGoldenTableAndText(t *testing.T) {
	o := RenderOptions{Options: Options{Columns: teamCols, ItemsKey: "items", Width: 100}}
	o.Format = "table"
	golden.Check(t, "table_list.golden", render(t, fixture, o))
	o.Format = "text"
	golden.Check(t, "text_list.golden", render(t, fixture, o))
	o.Format = "table"
	golden.Check(t, "table_object.golden", render(t, `{"id":"team_1","slug":"payments","state":"active","status":{"phase":"Ready"},"version":3}`, o))
	o.Format = "text"
	golden.Check(t, "text_object.golden", render(t, `{"id":"team_1","slug":"payments","state":"active","status":{"phase":"Ready"},"version":3}`, o))
	o.Columns = nil
	o.Format = "table"
	golden.Check(t, "table_generic.golden", render(t, fixture, o))
}

// NFR accessibility: table readable at 80 columns, with a truncation mark.
func TestTableFitsIn80Columns(t *testing.T) {
	cols := []Column{{"Id", "id"}, {"Slug", "slug"}, {"State", "state"}, {"Note", "note"}, {"Version", "version"}}
	got := render(t, fixture, RenderOptions{Format: "table", Options: Options{Columns: cols}})
	for _, l := range strings.Split(strings.TrimRight(got, "\n"), "\n") {
		if n := len([]rune(l)); n > 80 {
			t.Errorf("line of %d runes: %q", n, l)
		}
	}
	if !strings.Contains(got, "…") {
		t.Errorf("no truncation mark:\n%s", got)
	}
	vert := render(t, `{"note":"`+strings.Repeat("x", 200)+`"}`, RenderOptions{Format: "table", Options: Options{Columns: []Column{{"Note", "note"}}}})
	for _, l := range strings.Split(strings.TrimRight(vert, "\n"), "\n") {
		if len([]rune(l)) > 80 {
			t.Errorf("vertical line too long: %d", len([]rune(l)))
		}
	}
}

func TestGenericColumnsFallback(t *testing.T) {
	got := render(t, fixture, RenderOptions{Format: "table", Options: Options{Width: 200}})
	first := strings.Split(got, "\n")[0]
	for _, h := range []string{"id", "slug", "state", "version"} {
		if !strings.Contains(first, h) {
			t.Errorf("header lacks %s: %q", h, first)
		}
	}
	if strings.Contains(first, "quotas") || strings.Contains(first, "tags") {
		t.Errorf("non-scalar columns in %q", first)
	}
}

// CLI-6: an empty list is a header only in a table and nothing in text.
func TestEmptyListRendering(t *testing.T) {
	o := Options{Columns: teamCols, ItemsKey: "items"}
	if got := render(t, `{"items":[]}`, RenderOptions{Format: "table", Options: o}); got != "Id  Slug  State  Phase  Version\n" {
		t.Errorf("table = %q", got)
	}
	if got := render(t, `{"items":[]}`, RenderOptions{Format: "text", Options: o}); got != "" {
		t.Errorf("text = %q", got)
	}
	if got := render(t, `{"items":[]}`, RenderOptions{Format: "json"}); got != "{\n  \"items\": []\n}\n" {
		t.Errorf("json = %q", got)
	}
	if got := render(t, `{"items":[]}`, RenderOptions{Format: "table"}); got != "" {
		t.Errorf("generic empty table = %q", got)
	}
}

func TestTextScalarsAndLists(t *testing.T) {
	for in, want := range map[string]string{`"abc"`: "abc\n", `5`: "5\n", `["a","b"]`: "a\tb\n", `null`: "\n", `{"a":1,"b":"x"}`: "1\tx\n"} {
		if got := render(t, in, RenderOptions{Format: "text"}); got != want {
			t.Errorf("%s: %q, want %q", in, got, want)
		}
	}
}

func TestColor(t *testing.T) {
	env := map[string]string{}
	get := func(k string) string { return env[k] }
	if ColorEnabled("auto", false, get) {
		t.Error("auto on without a terminal")
	}
	if !ColorEnabled("auto", true, get) {
		t.Error("auto off on a terminal")
	}
	env["NO_COLOR"] = "1"
	if ColorEnabled("auto", true, get) {
		t.Error("NO_COLOR ignored")
	}
	if !ColorEnabled("on", false, get) || ColorEnabled("off", true, get) {
		t.Error("explicit on/off")
	}
	got := render(t, fixture, RenderOptions{Format: "table", Options: Options{Columns: teamCols, Color: true}})
	if !strings.Contains(got, "\x1b[1mId\x1b[0m") {
		t.Errorf("no bold header: %q", got)
	}
	if plain := render(t, fixture, RenderOptions{Format: "table", Options: Options{Columns: teamCols}}); strings.Contains(plain, "\x1b") {
		t.Error("escape codes without color")
	}
	vert := render(t, `{"a":1}`, RenderOptions{Format: "table", Options: Options{Columns: []Column{{"A", "a"}}, Color: true}})
	if !strings.Contains(vert, "\x1b[1mValue") {
		t.Errorf("vertical header: %q", vert)
	}
}

func TestRenderUnknownFormat(t *testing.T) {
	if err := Render(io.Discard, []byte(`{}`), RenderOptions{Format: "xml"}); err == nil {
		t.Error("xml rendered")
	}
}

// ---- csv and pager (O4) ----

// CLI-6: pages concatenate under one header; bytes are unchanged.
func TestCSVPagesKeepOneHeader(t *testing.T) {
	var b bytes.Buffer
	w := &CSVWriter{W: &b}
	pages := []string{"h1,h2\r\na,1\r\n", "h1,h2\r\nb,2\r\n", "h1,h2\r\nc,3\r\n", "h1,h2"}
	for _, p := range pages {
		if err := w.WritePage([]byte(p)); err != nil {
			t.Fatal(err)
		}
	}
	if want := "h1,h2\r\na,1\r\nb,2\r\nc,3\r\n"; b.String() != want {
		t.Errorf("csv = %q, want %q", b.String(), want)
	}
}

func TestCSVBytesUnchanged(t *testing.T) {
	var b bytes.Buffer
	body := "x,\"y\nz\"\r\n\xff,é\r\n"
	if err := (&CSVWriter{W: &b}).WritePage([]byte(body)); err != nil || b.String() != body {
		t.Errorf("bytes changed: %q", b.String())
	}
}

type countingWriter struct{ writes int }

func (c *countingWriter) Write(p []byte) (int, error) { c.writes++; return len(p), nil }

func TestCSVStreamsPageByPage(t *testing.T) {
	cw := &countingWriter{}
	w := &CSVWriter{W: cw}
	_ = w.WritePage([]byte("h\r\na\r\n"))
	if cw.writes != 1 {
		t.Error("first page not written immediately")
	}
	_ = w.WritePage([]byte("h\r\nb\r\n"))
	if cw.writes != 2 {
		t.Error("second page not written immediately")
	}
}

type fakePager struct {
	cmd     string
	started bool
	buf     bytes.Buffer
	closed  bool
}

func (f *fakePager) start(cmd string, out io.Writer) (io.WriteCloser, func() error, error) {
	f.cmd, f.started = cmd, true
	return &pagerWriter{f: f}, func() error { return nil }, nil
}

type pagerWriter struct{ f *fakePager }

func (p *pagerWriter) Write(b []byte) (int, error) { return p.f.buf.Write(b) }
func (p *pagerWriter) Close() error                { p.f.closed = true; return nil }

func pagerEnv(m map[string]string) func(string) string { return func(k string) string { return m[k] } }

func TestPagerOnlyOnTTY(t *testing.T) {
	var out bytes.Buffer
	fp := &fakePager{}
	w, done := OpenPager(PagerOptions{Out: &out, TTY: false, Getenv: pagerEnv(nil), Start: fp.start}, false)
	_, _ = w.Write([]byte("x"))
	_ = done()
	if fp.started || out.String() != "x" {
		t.Error("pager used without a terminal")
	}
	out.Reset()
	w, done = OpenPager(PagerOptions{Out: &out, TTY: true, Getenv: pagerEnv(nil), Start: fp.start}, false)
	_, _ = w.Write([]byte("y"))
	_ = done()
	if !fp.started || fp.cmd != "less -R" || fp.buf.String() != "y" || !fp.closed {
		t.Errorf("pager: %+v", fp)
	}
}

func TestNoCLIPager(t *testing.T) {
	var out bytes.Buffer
	fp := &fakePager{}
	w, _ := OpenPager(PagerOptions{Out: &out, TTY: true, NoPager: true, Getenv: pagerEnv(nil), Start: fp.start}, false)
	_, _ = w.Write([]byte("z"))
	if fp.started || out.String() != "z" {
		t.Error("--no-cli-pager ignored")
	}
}

func TestPagerCommandPrecedenceAndEmptyDisables(t *testing.T) {
	if c, ok := PagerCommand(pagerEnv(map[string]string{"AFS_PAGER": "bat", "PAGER": "more"}), true); c != "bat" || !ok {
		t.Errorf("AFS_PAGER: %q %v", c, ok)
	}
	if c, _ := PagerCommand(pagerEnv(map[string]string{"PAGER": "more"}), false); c != "more" {
		t.Errorf("PAGER: %q", c)
	}
	if _, ok := PagerCommand(pagerEnv(nil), true); ok {
		t.Error("empty AFS_PAGER should disable")
	}
	var out bytes.Buffer
	fp := &fakePager{}
	w, _ := OpenPager(PagerOptions{Out: &out, TTY: true, Getenv: pagerEnv(nil), Start: fp.start}, true)
	_, _ = w.Write([]byte("q"))
	if fp.started {
		t.Error("empty AFS_PAGER did not disable the pager")
	}
}

func TestEmptyAFSPagerDisables(t *testing.T) { TestPagerCommandPrecedenceAndEmptyDisables(t) }

func TestPagerStartFailureWritesDirect(t *testing.T) {
	var out bytes.Buffer
	w, _ := OpenPager(PagerOptions{Out: &out, TTY: true, Getenv: pagerEnv(nil), Start: func(string, io.Writer) (io.WriteCloser, func() error, error) {
		return nil, nil, errors.New("no such pager")
	}}, false)
	_, _ = w.Write([]byte("d"))
	if out.String() != "d" {
		t.Error("fallback to direct output failed")
	}
}

func TestStartProcessRealPager(t *testing.T) {
	var out bytes.Buffer
	w, wait, err := startProcess("cat", &out)
	if err != nil {
		t.Skip("cat not available")
	}
	_, _ = w.Write([]byte("through cat"))
	_ = w.Close()
	if err := wait(); err != nil || out.String() != "through cat" {
		t.Errorf("out %q err %v", out.String(), err)
	}
	if _, _, err := startProcess("", &out); err == nil {
		t.Error("empty command accepted")
	}
	if _, _, err := startProcess("definitely-not-a-program-xyz", &out); err == nil {
		t.Error("missing program started")
	}
}

func TestNormalizeLiterals(t *testing.T) {
	for in, want := range map[string]string{
		"a==`active`":      "a==`\"active\"`",
		"a==`5`":           "a==`5`",
		"a==`\"x\"`":       "a==`\"x\"`",
		"a=='it`s'":        "a=='it`s'",
		"\"a`b\".c":        "\"a`b\".c",
		"a==`unterminated": "a==`unterminated",
		"a==`x y`":         "a==`\"x y\"`",
		"'unterminated":    "'unterminated",
		"\"unterminated":   "\"unterminated",
	} {
		if got := normalizeLiterals(in); got != want {
			t.Errorf("%s -> %s, want %s", in, got, want)
		}
	}
}

func TestNumbersToFloatEdges(t *testing.T) {
	if got := numbersToFloat(json.Number("12345678901234567890")); got != "12345678901234567890" {
		t.Errorf("big int = %v", got)
	}
	if got := numbersToFloat(json.Number("9007199254740993")); got != "9007199254740993" {
		t.Errorf("2^53+1 = %v", got)
	}
	if got := numbersToFloat(json.Number("2.5")); got != 2.5 {
		t.Errorf("float = %v", got)
	}
	if got := numbersToFloat(json.Number("1e999")); got != "1e999" {
		t.Errorf("inf = %v", got)
	}
}

// CLI-13: a usage row's group key is spread into key.<member> columns.
func TestGenericColumnsSpreadUsageKey(t *testing.T) {
	tree, err := ParseTree([]byte(`{"rows":[{"calls":2,"key":{"agent":"agt_1","model":"m1"},"other":{"x":1}}]}`))
	if err != nil {
		t.Fatal(err)
	}
	var b strings.Builder
	if err := WriteTable(&b, tree, Options{}); err != nil {
		t.Fatal(err)
	}
	head := strings.SplitN(b.String(), "\n", 2)[0]
	if !strings.Contains(head, "key.agent") || !strings.Contains(head, "key.model") || strings.Contains(head, "other") {
		t.Errorf("header %q", head)
	}
}
