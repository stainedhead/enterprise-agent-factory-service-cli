package admin

import (
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/output"
)

func init() { registrations = append(registrations, registerPriceProposals) }

var proposalColumns = []output.Column{
	{Header: "Id", Path: "id"}, {Header: "Model", Path: "modelId"}, {Header: "EffectiveFrom", Path: "effectiveFrom"},
	{Header: "InputUsdPerMtok", Path: "inputUsdPerMtok"}, {Header: "OutputUsdPerMtok", Path: "outputUsdPerMtok"},
	{Header: "Status", Path: "status"}, {Header: "DecidedBy", Path: "decidedBy"}, {Header: "Note", Path: "note"},
}

func registerPriceProposals(r *cli.Registry) {
	r.Override(cli.Override{Op: "ListModelPriceProposals", Columns: proposalColumns, Examples: []cli.Example{
		{Title: "List the pending price proposals", Command: "afs admin list-price-proposals --status pending --output table"},
	}})
	r.Override(cli.Override{Op: "ApproveModelPriceProposal", Columns: proposalColumns, Examples: []cli.Example{
		{Title: "Approve a proposal", Command: `afs admin approve-price-proposal --proposal-id 1042 --note "Matches the AWS price list"`},
	}})
	r.Override(cli.Override{Op: "RejectModelPriceProposal", Columns: proposalColumns, Examples: []cli.Example{
		{Title: "Reject a proposal", Command: `afs admin reject-price-proposal --proposal-id 1042 --note "Waiting for the contract rate"`},
	}})
}
