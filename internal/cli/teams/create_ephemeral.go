package teams

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
)

func init() { registrations = append(registrations, registerEphemeral) }

// ephemeralOptions are --budget-usd and --budget-thresholds: the lifetime budget of an ephemeral
// team, sent as budget {period: lifetime, amountUsd, thresholds} (CLI-17, D-C19).
var ephemeralOptions = []cli.Option{
	{Flag: "--budget-usd", In: "local", Kind: cli.String, Desc: "Lifetime budget of an ephemeral team in USD, for example 500.00 (required with --kind ephemeral)."},
	{Flag: "--budget-thresholds", In: "local", Kind: cli.List, Elem: "integer", Desc: "Percent-of-budget alert thresholds of an ephemeral team, for example 50 80 100."},
}

func text(raw json.RawMessage) string {
	var s string
	if json.Unmarshal(raw, &s) == nil {
		return s
	}
	return string(raw)
}

// buildEphemeral checks the three required members of an ephemeral team before any call and
// folds the budget flags into the body's budget object. The members may come from options or
// from --cli-input-json (CLI-US-8.1).
func buildEphemeral(p *cli.Parsed, parts *cli.RequestParts) error {
	if parts.Body == nil {
		return nil
	}
	kind := ""
	if raw, ok := parts.Body.Get("kind"); ok {
		kind = text(raw)
	}
	flags := p.Has("budget-usd") || p.Has("budget-thresholds")
	if kind != "ephemeral" {
		if flags {
			return fmt.Errorf("--budget-usd and --budget-thresholds apply to --kind ephemeral")
		}
		return nil
	}
	var missing []string
	if _, ok := parts.Body.Get("expiresAt"); !ok {
		missing = append(missing, "--expires-at")
	}
	if _, ok := parts.Body.Get("teardownRule"); !ok {
		missing = append(missing, "--teardown-rule")
	}
	budget := map[string]any{}
	if raw, ok := parts.Body.Get("budget"); ok {
		_ = json.Unmarshal(raw, &budget)
	}
	if !p.Has("budget-usd") && budget["amountUsd"] == nil {
		missing = append(missing, "--budget-usd")
	}
	if len(missing) > 0 {
		return fmt.Errorf("an ephemeral team needs --expires-at, --teardown-rule and --budget-usd; missing %s", strings.Join(missing, ", "))
	}
	if !flags {
		return nil
	}
	budget["period"] = "lifetime"
	if p.Has("budget-usd") {
		budget["amountUsd"] = p.Value("budget-usd")
	}
	if p.Has("budget-thresholds") {
		var ts []json.Number
		for _, v := range p.Values["budget-thresholds"] {
			ts = append(ts, json.Number(v))
		}
		budget["thresholds"] = ts
	}
	raw, err := json.Marshal(budget)
	if err != nil {
		return err
	}
	parts.Body.Set("budget", raw)
	return nil
}

func registerEphemeral(r *cli.Registry) {
	r.Override(cli.Override{Op: "CreateTeam", Extra: ephemeralOptions, Build: buildEphemeral, Examples: []cli.Example{
		{Title: "Create an ephemeral team with a lifetime budget", Command: `afs teams create-team --slug cve-sprint --display-name "CVE sprint" --kind ephemeral --max-budget-usd 500.00 --first-admin '{"email":"ada@example.com"}' --expires-at 2026-12-31T00:00:00Z --teardown-rule at_cap_or_expiry --budget-usd 500.00 --budget-thresholds 50 80 100`},
	}})
}
