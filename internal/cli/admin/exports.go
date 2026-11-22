package admin

import (
	"fmt"
	"os"
	"strings"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/output"
)

func init() { registrations = append(registrations, registerExports) }

// maxBundleBytes is the API's limit for an import bundle (D-C18, CLI-US-18.3).
const maxBundleBytes = 100 << 20

var exportColumns = []output.Column{
	{Header: "Id", Path: "id"}, {Header: "Status", Path: "status"}, {Header: "TeamIds", Path: "teamIds"},
	{Header: "SizeBytes", Path: "sizeBytes"}, {Header: "Sha256", Path: "sha256"}, {Header: "FinishedAt", Path: "finishedAt"},
	{Header: "DownloadUrl", Path: "downloadUrl"}, {Header: "Error", Path: "error"},
}

var exportListColumns = []output.Column{
	{Header: "Id", Path: "id"}, {Header: "Status", Path: "status"}, {Header: "TeamIds", Path: "teamIds"},
	{Header: "SizeBytes", Path: "sizeBytes"}, {Header: "FinishedAt", Path: "finishedAt"}, {Header: "Error", Path: "error"},
}

var importColumns = []output.Column{
	{Header: "Id", Path: "id"}, {Header: "DryRun", Path: "dryRun"}, {Header: "Teams", Path: "teams"},
}

// checkBundleSize fails before the bundle is read when its file is over the limit (CLI-US-18.3).
func checkBundleSize(p *cli.Parsed) error {
	path, ok := strings.CutPrefix(p.Value("bundle"), "file://")
	if !ok {
		return nil
	}
	st, err := os.Stat(path)
	if err != nil {
		return nil // the read reports a missing file
	}
	if st.Size() > maxBundleBytes {
		return fmt.Errorf("argument --bundle: %s is %d bytes; the limit is 100 MB (%d bytes)", path, st.Size(), maxBundleBytes)
	}
	return nil
}

func registerExports(r *cli.Registry) {
	r.Override(cli.Override{Op: "ListExports", Columns: exportListColumns, Examples: []cli.Example{
		{Title: "List exports", Command: "afs admin list-exports"},
	}})
	r.Override(cli.Override{Op: "CreateExport", Columns: exportColumns, Examples: []cli.Example{
		{Title: "Export two teams", Command: "afs admin create-export --team-ids team_01J9ZQ3ABCDEFGHJKMNPQRSTVW team_01J9ZQ3ABCDEFGHJKMNPQRSTVX"},
		{Title: "Export every team", Command: "afs admin create-export"},
	}})
	r.Override(cli.Override{Op: "GetExport", Columns: exportColumns, Examples: []cli.Example{
		{Title: "Show an export (sha256, size and download URL when it has finished)", Command: "afs admin describe-export --export-id exp_01J9ZQ3ABCDEFGHJKMNPQRSTVW"},
	}})
	r.Override(cli.Override{Op: "CreateImport", Columns: importColumns, Precheck: checkBundleSize,
		Build: func(p *cli.Parsed, parts *cli.RequestParts) error {
			if !p.Has("dry-run") {
				parts.Query.Set("dryRun", "true") // an import is a dry run unless --no-dry-run (CLI-US-18.2)
			}
			return nil
		}, Examples: []cli.Example{
			{Title: "Dry-run an import (nothing is changed)", Command: "afs admin create-import --bundle file://bundle.json --account-map file://accounts.json"},
			{Title: "Apply the import", Command: "afs admin create-import --bundle file://bundle.json --account-map file://accounts.json --no-dry-run"},
		}})
}
