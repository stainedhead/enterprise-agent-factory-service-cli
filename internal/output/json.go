package output

import (
	"bytes"
	"encoding/json"
	"io"
)

// WriteJSON re-indents raw JSON without touching member order, nulls, numbers or unknown
// fields (D-C5, CLI-6).
func WriteJSON(w io.Writer, raw []byte) error {
	var b bytes.Buffer
	if err := json.Indent(&b, bytes.TrimSpace(raw), "", "  "); err != nil {
		return err
	}
	b.WriteByte('\n')
	_, err := w.Write(b.Bytes())
	return err
}

// writeTreeJSON writes a tree (a query result) indented.
func writeTreeJSON(w io.Writer, v any) error {
	raw, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return WriteJSON(w, raw)
}
