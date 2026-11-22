// Package output renders API results: json, yaml, table, text and csv, with JMESPath queries.
package output

// Column is one table or text column: a header and the JMESPath-lite path of the value
// (dotted keys, for example "status.phase").
type Column struct {
	Header string
	Path   string
	// Null is what a table or text cell shows when the value is absent or null ("-" for a hostname
	// that is not set yet); empty keeps the cell empty.
	Null string
}

// cellOf renders the cell of a row for a column.
func cellOf(row any, c Column) string {
	v := lookup(row, c.Path)
	if v == nil && c.Null != "" {
		return c.Null
	}
	return cell(v)
}
