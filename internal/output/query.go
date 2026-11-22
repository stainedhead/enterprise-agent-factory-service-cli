package output

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"strings"

	"github.com/jmespath-community/go-jmespath"
)

// QueryError is an invalid or failing --query (usage error, 252).
type QueryError struct{ Err error }

func (e *QueryError) Error() string { return fmt.Sprintf("invalid --query: %v", e.Err) }
func (e *QueryError) Unwrap() error { return e.Err }

// Query applies a JMESPath expression to raw JSON and returns the result as an ordered tree
// (keys sorted). Numbers are read with UseNumber, so large integers survive (D-C5).
func Query(expr string, raw []byte) (any, error) {
	compiled, err := jmespath.Compile(normalizeLiterals(expr))
	if err != nil {
		return nil, &QueryError{err}
	}
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	var data any
	if err := dec.Decode(&data); err != nil {
		return nil, err
	}
	res, err := compiled.Search(numbersToFloat(data))
	if err != nil {
		return nil, &QueryError{err}
	}
	return fromGeneric(res), nil
}

// Validate compiles an expression without data.
func Validate(expr string) error {
	if _, err := jmespath.Compile(normalizeLiterals(expr)); err != nil {
		return &QueryError{err}
	}
	return nil
}

// numbersToFloat converts json.Number to float64 for the JMESPath engine; numbers that float64
// cannot hold exactly stay strings so no digit is lost.
func numbersToFloat(v any) any {
	switch x := v.(type) {
	case json.Number:
		if i, err := x.Int64(); err == nil {
			if i > 1<<53 || i < -(1<<53) {
				return x.String()
			}
			return float64(i)
		}
		if !strings.ContainsAny(x.String(), ".eE") {
			return x.String() // an integer too big for int64
		}
		f, err := x.Float64()
		if err != nil || math.IsInf(f, 0) {
			return x.String()
		}
		return f
	case map[string]any:
		for k, val := range x {
			x[k] = numbersToFloat(val)
		}
	case []any:
		for i := range x {
			x[i] = numbersToFloat(x[i])
		}
	}
	return v
}

// normalizeLiterals accepts the AWS CLI habit of unquoted literals, such as `active`, which
// JMESPath proper requires to be JSON (`"active"`): a backtick literal that is not valid JSON
// becomes a JSON string.
func normalizeLiterals(expr string) string {
	var b strings.Builder
	for i := 0; i < len(expr); i++ {
		c := expr[i]
		switch c {
		case '\'':
			// a raw string literal: copy through the closing quote
			j := i + 1
			for j < len(expr) && (expr[j] != '\'' || expr[j-1] == '\\') {
				j++
			}
			if j >= len(expr) {
				j = len(expr) - 1
			}
			b.WriteString(expr[i : j+1])
			i = j
		case '"':
			j := i + 1
			for j < len(expr) && (expr[j] != '"' || expr[j-1] == '\\') {
				j++
			}
			if j >= len(expr) {
				j = len(expr) - 1
			}
			b.WriteString(expr[i : j+1])
			i = j
		case '`':
			j := i + 1
			for j < len(expr) && (expr[j] != '`' || expr[j-1] == '\\') {
				j++
			}
			if j >= len(expr) {
				b.WriteString(expr[i:])
				return b.String()
			}
			lit := expr[i+1 : j]
			if json.Valid([]byte(lit)) {
				b.WriteString(expr[i : j+1])
			} else {
				q, _ := json.Marshal(lit)
				b.WriteString("`" + string(q) + "`")
			}
			i = j
		default:
			b.WriteByte(c)
		}
	}
	return b.String()
}
