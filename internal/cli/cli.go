package cli

import (
	"fmt"
	"io"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/buildinfo"
)

// Exit codes follow the AWS CLI: 0 success, 252 usage error, 254 service error, 255 other.
const (
	ExitOK      = 0
	ExitUsage   = 252
	ExitService = 254
	ExitOther   = 255
)

// Services lists the top-level afs services (solution design section 11).
var Services = []string{
	"configure", "sso", "teams", "agents", "members", "identities", "secrets", "budgets",
	"usage", "logs", "notifications", "api-keys", "audit", "work-orders", "models", "admin",
}

const usage = `usage: afs [options] <service> <operation> [parameters]

To see help text, you can run:
  afs help
  afs <service> help
  afs <service> <operation> help
`

// Run executes the CLI with args (excluding the program name) and returns the process exit code.
func Run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		printf(stderr, "%s", usage)
		return ExitUsage
	}
	switch args[0] {
	case "--version", "version":
		printf(stdout, "%s\n", buildinfo.String("afs"))
		return ExitOK
	case "help", "--help", "-h":
		printf(stdout, "%s", usage)
		printf(stdout, "\nAvailable services:\n")
		for _, s := range Services {
			printf(stdout, "  - %s\n", s)
		}
		return ExitOK
	}
	printf(stderr, "%s\nafs: error: argument service: Invalid choice: %q\n", usage, args[0])
	return ExitUsage
}

// printf writes to a terminal stream; a failed write to stdout/stderr has no useful recovery.
func printf(w io.Writer, format string, a ...any) {
	_, _ = fmt.Fprintf(w, format, a...)
}
