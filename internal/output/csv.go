package output

import (
	"bytes"
	"io"
)

// CSVWriter concatenates CSV pages: the first page whole, later pages without their header row.
// Bytes are written unchanged, so the CSV is the API's, byte for byte (CLI-6).
type CSVWriter struct {
	W io.Writer
	// Resume drops the header of the first page: the output continues an earlier one (--starting-token).
	Resume bool
	wrote  bool
}

// WritePage writes one page as it arrives (pages are never buffered whole).
func (c *CSVWriter) WritePage(page []byte) error {
	if !c.wrote && !c.Resume {
		c.wrote = true
		_, err := c.W.Write(page)
		return err
	}
	c.wrote = true
	i := bytes.IndexByte(page, '\n')
	if i < 0 {
		return nil // header only, no newline
	}
	_, err := c.W.Write(page[i+1:])
	return err
}
