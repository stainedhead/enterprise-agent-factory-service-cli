package exec

import (
	"encoding/json"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/spec"
)

func placeholder(s spec.Schema) json.RawMessage {
	switch s.Type {
	case "integer", "number":
		return json.RawMessage("0")
	case "boolean":
		return json.RawMessage("false")
	case "array":
		return json.RawMessage("[]")
	case "object":
		return json.RawMessage("{}")
	}
	return json.RawMessage(`""`)
}

// Skeleton returns the --generate-cli-skeleton document of an operation: every path and query
// parameter, every writable body field, "version" for If-Match operations and "Idempotency-Key"
// for idempotent ones. No call is made.
func Skeleton(op *spec.Operation) *cli.Object {
	o := cli.NewObject()
	for _, p := range op.Params {
		o.Set(p.Name, placeholder(p.Schema))
	}
	if op.Body != nil {
		for _, f := range op.Body.Fields {
			if f.ReadOnly {
				continue
			}
			o.Set(f.Name, placeholder(f.Schema))
		}
	}
	if op.IfMatch != spec.IfMatchNone {
		o.Set("version", json.RawMessage("0"))
	}
	if op.Idempotent {
		o.Set("Idempotency-Key", json.RawMessage(`""`))
	}
	return o
}
