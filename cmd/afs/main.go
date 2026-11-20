// Command afs is the Agent Factory CLI. It follows AWS CLI conventions:
// "afs <service> <operation> [--options]".
package main

import (
	"os"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
)

func main() {
	os.Exit(cli.Run(os.Args[1:], os.Stdout, os.Stderr))
}
