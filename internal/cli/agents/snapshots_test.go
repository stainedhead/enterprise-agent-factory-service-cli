package agents_test

import (
	"testing"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/clitest"
)

func TestCommandsDoneSnapshotsAndCatalog(t *testing.T) {
	clitest.CheckOperations(t, "snapshots", "ListAgentSnapshots", "CreateAgentSnapshot")
	clitest.CheckOperations(t, "catalog", "ListStorageClasses", "ListAgentTemplates")
}
