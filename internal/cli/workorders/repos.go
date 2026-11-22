package workorders

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
)

// repoInput is one WorkOrderRepoInput; the member order is the wire order.
type repoInput struct {
	Host          string `json:"host"`
	Repo          string `json:"repo"`
	OwningTeamID  string `json:"owningTeamId,omitempty"`
	OwningContact string `json:"owningContact,omitempty"`
}

// ParseRepo reads one --repos item, host/owner/repo[=<owners>] (CLI-18). Owners are comma separated:
// a value with "@" is an owning contact (email), one starting with "team_" an owning team id, and
// anything else (a team slug, for instance) is refused. Without "=" the item carries no owner and the
// service applies the team configuration.
func ParseRepo(item string) (json.RawMessage, error) {
	spec, owners, hasOwners := strings.Cut(item, "=")
	parts := strings.Split(spec, "/")
	if len(parts) != 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" || strings.ContainsAny(spec, " \t") {
		return nil, fmt.Errorf("--repos item %q must look like host/owner/repo or host/owner/repo=<team id or email>", item)
	}
	in := repoInput{Host: parts[0], Repo: parts[1] + "/" + parts[2]}
	if hasOwners {
		if owners == "" {
			return nil, fmt.Errorf("--repos item %q names no owner after '='; give a team id (team_...) or an email address, or leave out the '='", item)
		}
		for _, o := range strings.Split(owners, ",") {
			switch {
			case strings.Contains(o, "@"):
				if in.OwningContact != "" {
					return nil, fmt.Errorf("--repos item %q names two owning contacts; give at most one email address", item)
				}
				in.OwningContact = o
			case strings.HasPrefix(o, "team_"):
				if in.OwningTeamID != "" {
					return nil, fmt.Errorf("--repos item %q names two owning teams; give at most one team id", item)
				}
				in.OwningTeamID = o
			default:
				return nil, fmt.Errorf("--repos item %q: owner %q is neither a team id nor an email; use a team id (team_...) or an email address, not a team name", item, o)
			}
		}
	}
	return json.Marshal(in)
}

// expandRepos turns the string items of a body member into WorkOrderRepoInput objects; items that
// are already JSON objects (from --cli-input-json) pass through.
func expandRepos(member string) cli.BuildFunc {
	return func(_ *cli.Parsed, parts *cli.RequestParts) error {
		raw, ok := parts.Body.Get(member)
		if !ok {
			return nil
		}
		var items []json.RawMessage
		if err := json.Unmarshal(raw, &items); err != nil {
			return fmt.Errorf("--%s must be a list of host/owner/repo[=owners] items", cli.Kebab(member))
		}
		for i, it := range items {
			var s string
			if json.Unmarshal(it, &s) != nil {
				continue
			}
			obj, err := ParseRepo(s)
			if err != nil {
				return err
			}
			items[i] = obj
		}
		out, _ := json.Marshal(items)
		parts.Body.Set(member, out)
		return nil
	}
}

// splitRepoPath lets update-repo take --repo host/owner/repo as well as --host, --owner and --repo.
func splitRepoPath(_ *cli.Parsed, parts *cli.RequestParts) error {
	if s := parts.Path["repo"]; strings.Count(s, "/") == 2 {
		seg := strings.Split(s, "/")
		if parts.Path["host"] != "" || parts.Path["owner"] != "" {
			return fmt.Errorf("give --repo as host/owner/repo, or --host, --owner and --repo separately, not both")
		}
		parts.Path["host"], parts.Path["owner"], parts.Path["repo"] = seg[0], seg[1], seg[2]
	}
	for _, k := range []string{"host", "owner", "repo"} {
		if parts.Path[k] == "" {
			return fmt.Errorf("the repository is required: --repo host/owner/repo")
		}
	}
	return nil
}
