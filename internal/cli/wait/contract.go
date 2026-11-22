package wait

import "github.com/stainedhead/enterprise-agent-factory-service-cli/internal/client/gen"

// Every member of a status response that Evaluate reads is referenced here, so a rename in the contract
// fails the build instead of making --wait run until its timeout (FR-017, ADR-014).
var _ = func(s gen.Status, p gen.PoolStatus) {
	_, _, _ = s.Phase, s.ObservedGeneration, s.Conditions
	_, _, _ = p.Generation, p.Workers.Desired, p.Workers.Reporting
}
