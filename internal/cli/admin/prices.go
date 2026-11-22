package admin

import (
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/output"
)

func init() { registrations = append(registrations, registerPrices) }

var priceColumns = []output.Column{
	{Header: "Model", Path: "modelId"}, {Header: "EffectiveFrom", Path: "effectiveFrom"}, {Header: "InputUsdPerMtok", Path: "inputUsdPerMtok"},
	{Header: "OutputUsdPerMtok", Path: "outputUsdPerMtok"}, {Header: "CacheRead", Path: "cacheReadUsdPerMtok"},
	{Header: "CacheWrite", Path: "cacheWriteUsdPerMtok"}, {Header: "Source", Path: "source"},
}

func registerPrices(r *cli.Registry) {
	r.Override(cli.Override{Op: "ListModelPrices", Columns: priceColumns, Examples: []cli.Example{
		{Title: "Show a model's price history", Command: "afs admin list-model-prices --model-id anthropic.claude-sonnet-4-5"},
	}})
	r.Override(cli.Override{Op: "CreateModelPrice", Columns: priceColumns, Examples: []cli.Example{
		{Title: "Record a new price (USD per million tokens, as strings)", Command: "afs admin create-model-price --model-id anthropic.claude-sonnet-4-5 --effective-from 2026-12-01T00:00:00Z --input-usd-per-mtok 3.00 --output-usd-per-mtok 15.00 --cache-read-usd-per-mtok 0.30 --cache-write-usd-per-mtok 3.75"},
	}})
}
