package exec

import (
	"strings"
	"testing"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
)

// FR-009: integer and number values must be JSON numbers; anything else is a usage error before a request.
func TestConvertIntegerAndNumberStrictness(t *testing.T) {
	cases := []struct {
		kind cli.Kind
		in   string
		ok   bool
	}{
		{cli.Int, "3", true}, {cli.Int, "-3", true}, {cli.Int, "0", true},
		{cli.Int, "+3", false}, {cli.Int, "03", false}, {cli.Int, "0x10", false}, {cli.Int, "1e3", false},
		{cli.Int, "1.0", false}, {cli.Int, "", false}, {cli.Int, " 3", false}, {cli.Int, "99999999999999999999", false},
		{cli.Number, "1.5", true}, {cli.Number, "-0", true}, {cli.Number, "2e3", true}, {cli.Number, "0.25", true},
		{cli.Number, "NaN", false}, {cli.Number, "Inf", false}, {cli.Number, "+Inf", false}, {cli.Number, "0x1p3", false},
		{cli.Number, "+1", false}, {cli.Number, "01", false}, {cli.Number, ".5", false}, {cli.Number, "1.", false}, {cli.Number, "1e999", false},
	}
	for _, c := range cases {
		o := cli.Option{Flag: "--workers", Wire: "workers", Kind: c.kind}
		raw, text, err := convert(o, []string{c.in}, Sources{}, nil)
		if c.ok {
			if err != nil || string(raw) != c.in || text != c.in {
				t.Errorf("%v %q: raw=%q text=%q err=%v", c.kind, c.in, raw, text, err)
			}
			continue
		}
		if err == nil {
			t.Errorf("%v %q accepted as %q", c.kind, c.in, raw)
		} else if !strings.Contains(err.Error(), "expected") || !strings.Contains(err.Error(), c.in) && c.in != "" {
			t.Errorf("%v %q: message %v", c.kind, c.in, err)
		}
	}
}

func TestConvertListsOfNumbersAreStrict(t *testing.T) {
	for _, elem := range []string{"integer", "number"} {
		o := cli.Option{Flag: "--thresholds", Wire: "thresholds", Kind: cli.List, Elem: elem}
		if raw, _, err := convert(o, []string{"50", "80", "100"}, Sources{}, nil); err != nil || string(raw) != "[50,80,100]" {
			t.Errorf("%s ok list: %s %v", elem, raw, err)
		}
		for _, bad := range []string{"050", "+80", "NaN", "0x10", "1e", ""} {
			if _, _, err := convert(o, []string{"50", bad}, Sources{}, nil); err == nil {
				t.Errorf("%s list accepted %q", elem, bad)
			}
		}
	}
}

func TestParameterNumbersAreStrict(t *testing.T) {
	for _, kind := range []cli.Kind{cli.Int, cli.Number} {
		for _, bad := range []string{"+3", "03", "0x10", "NaN"} {
			b := &Built{Path: map[string]string{}, Query: map[string][]string{}}
			if err := setParam(b, cli.Option{Flag: "--max-results", Wire: "maxResults", In: "query", Kind: kind}, bad); err == nil {
				t.Errorf("%v %q accepted", kind, bad)
			}
		}
	}
}
