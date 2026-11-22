package models_test

import (
	"testing"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/clitest"
)

func TestCommandsDone(t *testing.T) {
	clitest.CheckOperations(t, "models", "ListModels")
}
