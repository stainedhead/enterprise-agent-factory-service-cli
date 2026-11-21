package output

import (
	"encoding/json"
	"io"
	"strings"
	"unicode/utf8"
)

// Options control table and text rendering.
type Options struct {
	Columns  []Column
	ItemsKey string // "items" or "rows": the list member of a list response
	Color    bool
	Width    int // table width limit; 0 means 80
}

const minColWidth = 6

// lookup follows a dotted path through objects.
func lookup(v any, path string) any {
	if path == "" {
		return v
	}
	for _, part := range strings.Split(path, ".") {
		o, ok := v.(*Object)
		if !ok {
			return nil
		}
		v, ok = o.Vals[part]
		if !ok {
			return nil
		}
	}
	return v
}

func cell(v any) string {
	switch x := v.(type) {
	case []any:
		allScalar := true
		parts := make([]string, len(x))
		for i, e := range x {
			if !IsScalar(e) {
				allScalar = false
				break
			}
			parts[i] = Scalar(e)
		}
		if allScalar {
			return strings.Join(parts, ",")
		}
	}
	if IsScalar(v) {
		return Scalar(v)
	}
	b, _ := json.Marshal(v)
	return string(b)
}

// split returns the rows of a result: the items of a list response, a bare array, or the single
// object. list is false for a single object.
func split(v any, itemsKey string) (rows []any, list bool) {
	switch x := v.(type) {
	case []any:
		return x, true
	case *Object:
		for _, key := range []string{itemsKey, "items", "rows"} {
			if key == "" {
				continue
			}
			if arr, ok := x.Vals[key].([]any); ok {
				return arr, true
			}
		}
		return []any{x}, false
	}
	return []any{v}, false
}

// pickColumns uses the operation's columns when they fit the rows, else the generic ones, so a
// --query that reshapes the result still renders.
func pickColumns(rows []any, want []Column) []Column {
	if len(rows) > 0 && allScalar(rows) {
		return []Column{{Header: "Value", Path: ""}}
	}
	if want != nil && len(rows) > 0 {
		for _, c := range want {
			if lookup(rows[0], c.Path) != nil {
				return want
			}
		}
		return genericColumns(rows)
	}
	if want != nil {
		return want // empty list: header only
	}
	return genericColumns(rows)
}

// genericColumns are the scalar top-level members of the first row (TestGenericColumnsFallback).
func genericColumns(rows []any) []Column {
	if len(rows) == 0 {
		return nil
	}
	o, ok := rows[0].(*Object)
	if !ok {
		return nil
	}
	var cols []Column
	for _, k := range o.Keys {
		if IsScalar(o.Vals[k]) {
			cols = append(cols, Column{Header: k, Path: k})
		}
	}
	return cols
}

func pad(s string, w int) string {
	if n := utf8.RuneCountInString(s); n < w {
		return s + strings.Repeat(" ", w-n)
	}
	return s
}

func trunc(s string, w int) string {
	if utf8.RuneCountInString(s) <= w {
		return s
	}
	if w <= 1 {
		return "…"
	}
	r := []rune(s)
	return string(r[:w-1]) + "…"
}

// WriteTable renders a result as aligned columns: a header and one line per item, or two
// columns (Field, Value) for a single object. Lines fit the width (default 80) with a … mark.
func WriteTable(w io.Writer, v any, o Options) error {
	width := o.Width
	if width <= 0 {
		width = 80
	}
	rows, list := split(v, o.ItemsKey)
	cols := pickColumns(rows, o.Columns)
	if len(cols) == 0 {
		return nil
	}
	if !list {
		return writeVertical(w, rows[0], cols, o, width)
	}
	cells := make([][]string, len(rows))
	widths := make([]int, len(cols))
	for i, c := range cols {
		widths[i] = utf8.RuneCountInString(c.Header)
	}
	for r, row := range rows {
		cells[r] = make([]string, len(cols))
		for i, c := range cols {
			cells[r][i] = strings.NewReplacer("\n", " ", "\t", " ").Replace(cell(lookup(row, c.Path)))
			if n := utf8.RuneCountInString(cells[r][i]); n > widths[i] {
				widths[i] = n
			}
		}
	}
	fit(widths, width)
	var b strings.Builder
	line := func(vals []string, header bool) {
		for i, s := range vals {
			s = trunc(s, widths[i])
			if i == len(vals)-1 {
				s = strings.TrimRight(s, " ")
			} else {
				s = pad(s, widths[i])
			}
			if header {
				s = boldKeepPad(s, o.Color)
			}
			b.WriteString(s)
			if i < len(vals)-1 {
				b.WriteString("  ")
			}
		}
		b.WriteByte('\n')
	}
	headers := make([]string, len(cols))
	for i, c := range cols {
		headers[i] = c.Header
	}
	line(headers, true)
	for _, r := range cells {
		line(r, false)
	}
	_, err := io.WriteString(w, b.String())
	return err
}

func boldKeepPad(s string, on bool) string {
	if !on {
		return s
	}
	trimmed := strings.TrimRight(s, " ")
	return bold(trimmed, true) + s[len(trimmed):]
}

// fit shrinks the widest columns until the line (with two-space gaps) fits.
func fit(widths []int, limit int) {
	total := func() int {
		t := 2 * (len(widths) - 1)
		for _, w := range widths {
			t += w
		}
		return t
	}
	for total() > limit {
		wi := 0
		for i, w := range widths {
			if w > widths[wi] {
				wi = i
			}
		}
		if widths[wi] <= minColWidth {
			return
		}
		widths[wi]--
	}
}

func writeVertical(w io.Writer, row any, cols []Column, o Options, width int) error {
	kw := len("Field")
	for _, c := range cols {
		if n := utf8.RuneCountInString(c.Header); n > kw {
			kw = n
		}
	}
	vw := width - kw - 2
	if vw < minColWidth {
		vw = minColWidth
	}
	var b strings.Builder
	b.WriteString(boldKeepPad(pad("Field", kw), o.Color) + "  " + bold("Value", o.Color) + "\n")
	for _, c := range cols {
		b.WriteString(pad(c.Header, kw) + "  " + strings.TrimRight(trunc(strings.NewReplacer("\n", " ", "\t", " ").Replace(cell(lookup(row, c.Path))), vw), " ") + "\n")
	}
	_, err := io.WriteString(w, b.String())
	return err
}
