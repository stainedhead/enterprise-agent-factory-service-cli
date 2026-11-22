package safeupdate

import "github.com/stainedhead/enterprise-agent-factory-service-cli/internal/client/gen"

// The members the version readers match on are referenced here, so a rename in the contract fails the
// build instead of making a safe update read no version (FR-017, ADR-014).
var _ = func(t gen.Team, r gen.WorkOrderRepo) {
	_, _ = t.Version, t.Id
	_, _, _ = r.Host, r.Repo, r.Version
}
