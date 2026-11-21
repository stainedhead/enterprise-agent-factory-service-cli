package cli

import (
	"bytes"
	"encoding/json"
)

// Object is a JSON object that keeps member order and raw values, so a body supplied with
// --cli-input-json is sent with its order and values unchanged (D-C6).
type Object struct {
	keys []string
	vals map[string]json.RawMessage
}

// NewObject returns an empty object.
func NewObject() *Object { return &Object{vals: map[string]json.RawMessage{}} }

// ParseObject decodes a JSON object keeping member order.
func ParseObject(raw []byte) (*Object, error) {
	dec := json.NewDecoder(bytes.NewReader(raw))
	tok, err := dec.Token()
	if err != nil {
		return nil, err
	}
	if d, ok := tok.(json.Delim); !ok || d != '{' {
		return nil, &json.UnmarshalTypeError{Value: "non-object", Offset: 0}
	}
	o := NewObject()
	for dec.More() {
		kt, err := dec.Token()
		if err != nil {
			return nil, err
		}
		var v json.RawMessage
		if err := dec.Decode(&v); err != nil {
			return nil, err
		}
		o.Set(kt.(string), v)
	}
	if _, err := dec.Token(); err != nil {
		return nil, err
	}
	return o, nil
}

// Set sets a member, appending it when new.
func (o *Object) Set(k string, v json.RawMessage) {
	if _, ok := o.vals[k]; !ok {
		o.keys = append(o.keys, k)
	}
	o.vals[k] = v
}

// Get returns a member.
func (o *Object) Get(k string) (json.RawMessage, bool) {
	v, ok := o.vals[k]
	return v, ok
}

// Delete removes a member.
func (o *Object) Delete(k string) {
	if _, ok := o.vals[k]; !ok {
		return
	}
	delete(o.vals, k)
	for i, x := range o.keys {
		if x == k {
			o.keys = append(o.keys[:i], o.keys[i+1:]...)
			break
		}
	}
}

// Keys returns the member names in order.
func (o *Object) Keys() []string { return append([]string(nil), o.keys...) }

// Len is the number of members.
func (o *Object) Len() int { return len(o.keys) }

// MarshalJSON writes the members in order.
func (o *Object) MarshalJSON() ([]byte, error) {
	var b bytes.Buffer
	b.WriteByte('{')
	for i, k := range o.keys {
		if i > 0 {
			b.WriteByte(',')
		}
		kb, _ := json.Marshal(k)
		b.Write(kb)
		b.WriteByte(':')
		b.Write(o.vals[k])
	}
	b.WriteByte('}')
	return b.Bytes(), nil
}
