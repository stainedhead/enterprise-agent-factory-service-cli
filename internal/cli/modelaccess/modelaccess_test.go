package modelaccess

import (
	"strings"
	"testing"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
)

func parsed(values map[string][]string) *cli.Parsed { return &cli.Parsed{Values: values} }

// FR-012 (v1.5.0): --routes all, --routes <slugs> and --default-route fold into models; bad words are refused.
func TestRoutesFoldIntoModels(t *testing.T) {
	for name, tc := range map[string]struct {
		values map[string][]string
		want   string
		err    string
	}{
		"all":      {values: map[string][]string{"routes": {"all"}}, want: `{"models":{"routeAccess":"all"}}`},
		"selected": {values: map[string][]string{"routes": {"a", "b"}, "default-route": {"a"}}, want: `{"models":{"routeAccess":"selected","routes":["a","b"],"defaultRoute":"a"}}`},
		"mixed":    {values: map[string][]string{"routes": {"all", "a"}}, err: "cannot be combined"},
		"empty":    {values: map[string][]string{"routes": {}}, err: "at least one route"},
		"blank":    {values: map[string][]string{"routes": {"a", ""}}, err: "cannot be empty"},
	} {
		parts := &cli.RequestParts{Body: cli.NewObject()}
		err := Build(parsed(tc.values), parts)
		if tc.err != "" {
			if err == nil || !strings.Contains(err.Error(), tc.err) {
				t.Errorf("%s: err %v, want %q", name, err, tc.err)
			}
			continue
		}
		b, _ := parts.Body.MarshalJSON()
		if err != nil || string(b) != tc.want {
			t.Errorf("%s: %s %v, want %s", name, b, err, tc.want)
		}
	}
}
