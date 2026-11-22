package output

import (
	"encoding/json"
	"io"

	"gopkg.in/yaml.v3"
)

func toYAMLNode(v any) *yaml.Node {
	switch x := v.(type) {
	case nil:
		return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!null", Value: "null"}
	case bool:
		s := "false"
		if x {
			s = "true"
		}
		return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!bool", Value: s}
	case json.Number:
		tag := "!!float"
		if _, err := x.Int64(); err == nil {
			tag = "!!int"
		}
		return &yaml.Node{Kind: yaml.ScalarNode, Tag: tag, Value: x.String()}
	case string:
		n := &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: x}
		if x == "" || needsQuote(x) {
			n.Style = yaml.DoubleQuotedStyle
		}
		return n
	case []any:
		n := &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}
		for _, e := range x {
			n.Content = append(n.Content, toYAMLNode(e))
		}
		return n
	case *Object:
		n := &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}
		for _, k := range x.Keys {
			n.Content = append(n.Content, &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: k}, toYAMLNode(x.Vals[k]))
		}
		return n
	}
	return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: Scalar(v)}
}

// needsQuote reports strings YAML would read as another type (true, 12, null, ...).
func needsQuote(s string) bool {
	var v any
	if err := yaml.Unmarshal([]byte(s), &v); err != nil {
		return true
	}
	_, isStr := v.(string)
	return !isStr
}

// WriteYAML renders a tree as YAML in member order.
func WriteYAML(w io.Writer, v any) error {
	enc := yaml.NewEncoder(w)
	enc.SetIndent(2)
	if err := enc.Encode(toYAMLNode(v)); err != nil {
		return err
	}
	return enc.Close()
}
