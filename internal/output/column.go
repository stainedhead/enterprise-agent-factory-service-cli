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
	// NAUnless names a boolean member of the row (a surface's "supported"); when it is false the cell shows
	// "N/A" whatever Path holds, so a tool without a surface does not print a level it cannot have.
	NAUnless string
}

// cellOf renders the cell of a row for a column.
func cellOf(row any, c Column) string {
	if c.NAUnless != "" {
		if ok, isBool := lookup(row, c.NAUnless).(bool); isBool && !ok {
			return "N/A"
		}
	}
	v := lookup(row, c.Path)
	if v == nil && c.Null != "" {
		return c.Null
	}
	return cell(v)
}
