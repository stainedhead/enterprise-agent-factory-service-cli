package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestRunNoArgsIsUsageError(t *testing.T) {
	var out, errOut bytes.Buffer
	if code := Run(nil, &out, &errOut); code != ExitUsage {
		t.Fatalf("exit = %d, want %d", code, ExitUsage)
	}
}

func TestRunHelpListsServices(t *testing.T) {
	var out, errOut bytes.Buffer
	if code := Run([]string{"help"}, &out, &errOut); code != ExitOK {
		t.Fatalf("exit = %d, want 0", code)
	}
	for _, s := range []string{"teams", "agents", "secrets", "budgets"} {
		if !strings.Contains(out.String(), "- "+s) {
			t.Errorf("help output missing service %q", s)
		}
	}
}

func TestRunUnknownServiceIsUsageError(t *testing.T) {
	var out, errOut bytes.Buffer
	if code := Run([]string{"bogus"}, &out, &errOut); code != ExitUsage {
		t.Fatalf("exit = %d, want %d", code, ExitUsage)
	}
}
