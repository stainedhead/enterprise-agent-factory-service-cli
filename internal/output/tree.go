package output

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"sort"
)

// Object is a JSON object that keeps member order. Values are nil, bool, json.Number, string,
// []any or *Object.
type Object struct {
	Keys []string
	Vals map[string]any
}

// Get returns a member.
func (o *Object) Get(k string) (any, bool) { v, ok := o.Vals[k]; return v, ok }

// ParseTree decodes JSON into an ordered tree; numbers keep their text (json.Number).
func ParseTree(raw []byte) (any, error) {
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	v, err := readValue(dec)
	if err != nil {
		return nil, err
	}
	if _, err := dec.Token(); err != io.EOF {
		return nil, fmt.Errorf("unexpected data after the JSON value")
	}
	return v, nil
}

func readValue(dec *json.Decoder) (any, error) {
	tok, err := dec.Token()
	if err != nil {
		return nil, err
	}
	d, ok := tok.(json.Delim)
	if !ok {
		return tok, nil
	}
	switch d {
	case '{':
		o := &Object{Vals: map[string]any{}}
		for dec.More() {
			kt, err := dec.Token()
			if err != nil {
				return nil, err
			}
			k, _ := kt.(string)
			v, err := readValue(dec)
			if err != nil {
				return nil, err
			}
			if _, dup := o.Vals[k]; !dup {
				o.Keys = append(o.Keys, k)
			}
			o.Vals[k] = v
		}
		_, err = dec.Token()
		return o, err
	case '[':
		arr := []any{}
		for dec.More() {
			v, err := readValue(dec)
			if err != nil {
				return nil, err
			}
			arr = append(arr, v)
		}
		_, err = dec.Token()
		return arr, err
	}
	return nil, fmt.Errorf("unexpected delimiter %v", d)
}

// MarshalJSON writes the object in member order.
func (o *Object) MarshalJSON() ([]byte, error) {
	var b bytes.Buffer
	b.WriteByte('{')
	for i, k := range o.Keys {
		if i > 0 {
			b.WriteByte(',')
		}
		kb, _ := json.Marshal(k)
		b.Write(kb)
		b.WriteByte(':')
		vb, err := json.Marshal(o.Vals[k])
		if err != nil {
			return nil, err
		}
		b.Write(vb)
	}
	b.WriteByte('}')
	return b.Bytes(), nil
}

// fromGeneric converts decoded generic JSON (maps, slices) into an ordered tree with sorted keys.
func fromGeneric(v any) any {
	switch x := v.(type) {
	case map[string]any:
		o := &Object{Vals: make(map[string]any, len(x))}
		for k := range x {
			o.Keys = append(o.Keys, k)
		}
		sort.Strings(o.Keys)
		for _, k := range o.Keys {
			o.Vals[k] = fromGeneric(x[k])
		}
		return o
	case []any:
		out := make([]any, len(x))
		for i := range x {
			out[i] = fromGeneric(x[i])
		}
		return out
	case float64:
		return json.Number(fmtFloat(x))
	}
	return v
}

func fmtFloat(f float64) string {
	b, _ := json.Marshal(f)
	return string(b)
}

// IsScalar reports whether a tree value is not an object or an array.
func IsScalar(v any) bool {
	switch v.(type) {
	case *Object, []any:
		return false
	}
	return true
}

// Scalar renders a scalar for table and text output; null is empty.
func Scalar(v any) string {
	switch x := v.(type) {
	case nil:
		return ""
	case string:
		return x
	case bool:
		if x {
			return "true"
		}
		return "false"
	case json.Number:
		return x.String()
	}
	b, _ := json.Marshal(v)
	return string(b)
}
