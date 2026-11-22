package workorders

import (
	"strings"
	"testing"
)

const team1 = "team_01J9ZQ3ABCDEFGHJKMNPQRSTVW"

// CLI-18, CLI-US-9.2: how --repos items are classified into owner fields.
func TestReposClassification(t *testing.T) {
	for _, tc := range []struct {
		name, in, want, fail string
	}{
		{"team", "github.com/acme/ledger=" + team1, `{"host":"github.com","repo":"acme/ledger","owningTeamId":"` + team1 + `"}`, ""},
		{"email", "github.com/acme/ledger=a@example.com", `{"host":"github.com","repo":"acme/ledger","owningContact":"a@example.com"}`, ""},
		{"both", "github.com/acme/ledger=" + team1 + ",a@example.com", `{"host":"github.com","repo":"acme/ledger","owningTeamId":"` + team1 + `","owningContact":"a@example.com"}`, ""},
		{"none", "github.com/acme/ledger", `{"host":"github.com","repo":"acme/ledger"}`, ""},
		{"slug", "github.com/acme/ledger=payments", "", "use a team id (team_...) or an email address"},
		{"empty owner", "github.com/acme/ledger=", "", "names no owner"},
		{"two teams", "github.com/acme/ledger=" + team1 + "," + team1, "", "two owning teams"},
		{"two emails", "github.com/acme/ledger=a@example.com,b@example.com", "", "two owning contacts"},
		{"no host", "acme/ledger", "", "must look like host/owner/repo"},
		{"too deep", "github.com/acme/ledger/extra", "", "must look like host/owner/repo"},
		{"space", "github.com/acme/led ger", "", "must look like host/owner/repo"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ParseRepo(tc.in)
			if tc.fail != "" {
				if err == nil || !strings.Contains(err.Error(), tc.fail) {
					t.Fatalf("err = %v, want %q", err, tc.fail)
				}
				return
			}
			if err != nil || string(got) != tc.want {
				t.Fatalf("got %s, %v; want %s", got, err, tc.want)
			}
		})
	}
}
