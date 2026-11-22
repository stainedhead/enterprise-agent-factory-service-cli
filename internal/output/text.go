package output

import (
	"io"
	"strings"
)

// WriteText renders a result as AWS-style text: tab-separated values, one line per item (a
// single object is one line, a list of scalars is one line). An empty list prints nothing.
func WriteText(w io.Writer, v any, o Options) error {
	rows, list := split(v, o.ItemsKey)
	if IsScalar(v) {
		_, err := io.WriteString(w, Sanitize(Scalar(v))+"\n")
		return err
	}
	if len(rows) > 0 && list && allScalar(rows) {
		parts := make([]string, len(rows))
		for i, r := range rows {
			parts[i] = Sanitize(Scalar(r))
		}
		_, err := io.WriteString(w, strings.Join(parts, "\t")+"\n")
		return err
	}
	cols := pickColumns(rows, o.Columns)
	var b strings.Builder
	for _, row := range rows {
		vals := make([]string, len(cols))
		for i, c := range cols {
			vals[i] = Sanitize(strings.NewReplacer("\t", " ", "\n", " ").Replace(cell(lookup(row, c.Path))))
		}
		b.WriteString(strings.Join(vals, "\t") + "\n")
	}
	_, err := io.WriteString(w, b.String())
	return err
}

func allScalar(rows []any) bool {
	for _, r := range rows {
		if !IsScalar(r) {
			return false
		}
	}
	return true
}
