package output

import (
	"fmt"
	"io"
)

// RenderOptions describe one rendering of a response body.
type RenderOptions struct {
	Format string // json | yaml | table | text
	Query  string
	Options
}

// Render writes a JSON response body in the requested format. json without --query re-indents
// the raw body; every other combination decodes into an ordered tree first (D-C5).
func Render(w io.Writer, raw []byte, o RenderOptions) error {
	var tree any
	var err error
	if o.Query != "" {
		if tree, err = Query(o.Query, raw); err != nil {
			return err
		}
		if o.Format == "json" {
			return writeTreeJSON(w, tree)
		}
		// after a query the shape is the query's: do not guess an items key from the operation
		o.ItemsKey = ""
	} else {
		if o.Format == "json" {
			return WriteJSON(w, raw)
		}
		if tree, err = ParseTree(raw); err != nil {
			return err
		}
	}
	switch o.Format {
	case "yaml":
		return WriteYAML(w, tree)
	case "table":
		return WriteTable(w, tree, o.Options)
	case "text":
		return WriteText(w, tree, o.Options)
	}
	return fmt.Errorf("cannot render format %q", o.Format)
}
