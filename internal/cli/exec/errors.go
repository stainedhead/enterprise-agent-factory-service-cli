package exec

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/auth"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli/safeupdate"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli/wait"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/client"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/output"
)

// ReportError prints err in the AWS style on w and returns the exit code (CLI-8):
// 254 for an answer from the service, 130 for Ctrl-C, 255 for anything else.
func ReportError(w io.Writer, opID, profile string, err error) int {
	w = output.NewSanitizingWriter(w) // service-supplied text must not carry terminal escapes (FR-003)
	if errors.Is(err, context.Canceled) {
		return cli.ExitInterrupted
	}
	var ae *client.APIError
	if errors.As(err, &ae) {
		_, _ = fmt.Fprintf(w, "\nAn error occurred (%s) when calling the %s operation: %s\n", ae.Code, opID, ae.Message)
		if p := ae.RequiredPermission(); p != "" {
			_, _ = fmt.Fprintf(w, "Required permission: %s\n", p)
		}
		if d := ae.Dependents(); len(d) > 0 {
			_, _ = fmt.Fprintf(w, "Dependents: %s\n", strings.Join(d, ", "))
		}
		if m := ae.MissingTags(); len(m) > 0 {
			_, _ = fmt.Fprintf(w, "Missing tags: %s\n", strings.Join(m, ", "))
		}
		for _, a := range ae.AmbiguousRepos() {
			_, _ = fmt.Fprintf(w, "Ambiguous repository: %s is owned by more than one team (%s)\n", a.Repo, strings.Join(a.TeamIDs, ", "))
		}
		if u := ae.UnownedRepos(); len(u) > 0 {
			_, _ = fmt.Fprintf(w, "Unowned repositories: %s\n", strings.Join(u, ", "))
		}
		if len(ae.AmbiguousRepos()) > 0 || len(ae.UnownedRepos()) > 0 {
			_, _ = fmt.Fprintln(w, "Name the owner of each in --repos as host/owner/repo=<team_id or email> (for example github.com/acme/ledger=team_01J9ZQ3ABCDEFGHJKMNPQRSTVW).")
		}
		var rd *safeupdate.ReadDeniedError
		if errors.As(err, &rd) {
			_, _ = fmt.Fprintln(w, safeupdate.ReadDeniedHint)
		}
		if ae.Status == 401 {
			_, _ = fmt.Fprintf(w, "%s.\n", auth.LoginHint(profile))
		}
		if ae.Status == 412 {
			if v := strings.Trim(ae.Header.Get("ETag"), `"`); v != "" {
				_, _ = fmt.Fprintf(w, "The resource changed (version %s). Re-run with --if-match %s to apply your change to the latest version.\n", v, v)
			}
		}
		if ae.Status == 428 {
			_, _ = fmt.Fprintln(w, "The service requires a version: pass --if-match <version> (see the resource's current version with its describe or list operation).")
		}
		if ae.RequestID != "" {
			_, _ = fmt.Fprintf(w, "Request ID: %s\n", ae.RequestID)
		}
		return cli.ExitService
	}
	var fe *wait.FailedError
	if errors.As(err, &fe) {
		_, _ = fmt.Fprintf(w, "\nAn error occurred (ResourceFailed) when calling the %s operation: %s\n", opID, fe.Error())
		return cli.ExitService
	}
	var lr *auth.LoginRequiredError
	if errors.As(err, &lr) {
		_, _ = fmt.Fprintf(w, "\n%s\n", lr.Error())
		return cli.ExitOther
	}
	_, _ = fmt.Fprintf(w, "\nafs: error: %v\n", err)
	return cli.ExitOther
}
