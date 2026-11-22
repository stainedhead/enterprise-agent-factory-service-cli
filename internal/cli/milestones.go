package cli

import "github.com/stainedhead/enterprise-agent-factory-service-cli/internal/spec"

// CurrentMilestone is the highest milestone whose operations are complete (example, table
// columns, golden files). It is bumped only by milestone exit tasks (D-C4).
const CurrentMilestone = 5

// MilestoneOf returns the milestone (0..6) in which an operation is completed, or -1 when unknown.
func MilestoneOf(opID string) int {
	if m, ok := milestoneOf[opID]; ok {
		return m
	}
	return -1
}

// StreamOf returns the stream that owns an operation (tasks.md Appendix A).
func StreamOf(opID string) string { return streamOf[opID] }

// Pending lists operations at or below CurrentMilestone that lack an example or table columns.
func (r *Registry) Pending() []string {
	var out []string
	for _, c := range r.byOp {
		if MilestoneOf(c.Op.ID) > CurrentMilestone {
			continue
		}
		if len(c.Examples) == 0 || (len(c.Columns) == 0 && !c.Dynamic && !noBody(c.Op)) {
			out = append(out, c.Op.ID)
		}
	}
	return out
}

// noBody is true when every success response is 204: there is nothing to put in columns.
func noBody(op *spec.Operation) bool {
	for _, s := range op.Success {
		if s != 204 {
			return false
		}
	}
	return len(op.Success) > 0
}

var _ = spec.ContractVersion
