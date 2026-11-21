// Package output renders API results: json, yaml, table, text and csv, with JMESPath queries.
package output

// Column is one table or text column: a header and the JMESPath-lite path of the value
// (dotted keys, for example "status.phase").
type Column struct {
	Header string
	Path   string
}
