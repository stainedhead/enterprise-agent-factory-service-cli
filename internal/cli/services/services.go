// Package services is the composition root of the command table: it imports every service
// package, wires the executor and help, and returns the finished registry.
package services

import (
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli/admin"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli/agents"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli/apikeys"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli/audit"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli/budgets"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli/collaboration"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli/completion"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli/configure"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli/exec"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli/help"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli/members"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli/models"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli/pools"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli/reports"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli/secrets"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli/sso"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli/swarmtemplates"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli/teams"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli/wait"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli/workorders"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/spec"
)

// Register applies every service's registrations to r.
func Register(r *cli.Registry) {
	for _, f := range []func(*cli.Registry){
		teams.Register, members.Register, apikeys.Register, agents.Register, secrets.Register,
		budgets.Register, reports.Register, models.Register, audit.Register, admin.Register,
		pools.Register, swarmtemplates.Register, workorders.Register, sso.Register,
		configure.Register, completion.Register, collaboration.Register,
	} {
		f(r)
	}
}

// NewRegistry builds the complete command table for the pinned snapshot.
func NewRegistry() (*cli.Registry, error) {
	r := cli.NewRegistry(spec.All())
	Register(r)
	r.Exec, r.Help = exec.Run, help.Run
	for id := range wait.Table {
		r.AddOptions(id, exec.WaitOptions()...)
	}
	if err := r.Finalize(); err != nil {
		return nil, err
	}
	return r, nil
}
