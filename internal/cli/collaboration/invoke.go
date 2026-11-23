package collaboration

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/auth"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/buildinfo"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli/exec"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli/orchestration"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli/runtime"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/client"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/edge"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/output"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/redact"
)

func init() { registrations = append(registrations, registerInvoke) }

// invokeOp names invoke-api in error messages; it is not an operation of the contract (API contract 4.4b).
const invokeOp = "InvokeApi"

// maxInvokeBody is the largest body invoke-api sends: the maximum of the service setting collab.max_body_mb.
// The edge applies the real limit (default 25 MB) and answers 400 body_too_large (CLI-US-25.4).
const maxInvokeBody = 100 << 20

// maxErrorPeek bounds how much of an error response is held in memory to tell the edge's envelope from
// Paperclip's own body.
const maxErrorPeek = 1 << 20

var invokeMethods = map[string]bool{"GET": true, "HEAD": true, "POST": true, "PUT": true, "PATCH": true, "DELETE": true}

// refusedGlobals are the global options that mean nothing for a command that prints the response unchanged
// (PC-CLI-D7).
var refusedGlobals = []string{"query", "output", "no-paginate", "max-items", "page-size", "starting-token", "cli-input-json", "generate-cli-skeleton"}

func registerInvoke(r *cli.Registry) {
	r.AddLocal(&cli.Command{
		Service: "collaboration", Name: "invoke-api",
		Summary: "Call a team's Paperclip REST API (the orchestration agent's Web API) with your AFS credential and a Paperclip key.",
		Options: []cli.Option{
			{Flag: "--team-id", In: "local", Kind: cli.String, Desc: "The team (or the profile's team_id). Not needed with --collab-url."},
			{Flag: "--method", In: "local", Kind: cli.String, Desc: "GET, HEAD, POST, PUT, PATCH or DELETE (default GET)."},
			{Flag: "--path", In: "local", Kind: cli.String, Required: true, Desc: "The Paperclip path, starting with / (for example /api/issues)."},
			{Flag: "--body", In: "local", Kind: cli.String, Desc: "The request body: file://PATH or fileb://PATH."},
			{Flag: "--content-type", In: "local", Kind: cli.String, Desc: "The body's media type (default application/json)."},
			{Flag: "--paperclip-key-file", In: "local", Kind: cli.String, Required: true, Desc: "A file (mode 0600) holding the Paperclip key, or - for standard input. The key is never an argument."},
			{Flag: "--collab-url", In: "local", Kind: cli.String, Desc: "Use this https origin instead of the Web API address of the team's orchestration agent."},
			{Flag: "--include", In: "local", Kind: cli.Bool, Desc: "Print the status line and response headers before the body."},
			{Flag: "--wait", In: "local", Kind: cli.Bool, Desc: "Not supported: invoke-api prints the response unchanged."},
		},
		Local: invoke,
		Override: cli.Override{Examples: []cli.Example{
			{Title: "Read the team's companies through the Web API at its effective level (private or public)", Command: "afs collaboration invoke-api --team-id " + exTeam + " --path /api/companies --paperclip-key-file ./pc.key"},
			{Title: "Create an issue from a pipeline", Command: "afs collaboration invoke-api --team-id " + exTeam + " --method POST --path /api/issues --body file://issue.json --paperclip-key-file ./pc.key"},
		}},
	})
}

//nolint:staticcheck // ST1005: the sentence is fixed by CLI-US-25.2
var errWorkload = errors.New("The collaboration endpoint does not admit workload credentials; use an API key or afs sso login.")

type invokeArgs struct {
	teamID, method, path, ctype, keyFile, collabURL string
	body                                            []byte
	include                                         bool
}

func invoke(ctx context.Context, inv *cli.Invocation) int {
	env := inv.Env
	args, err := parseInvoke(inv)
	var fe *bodyFileError
	if errors.As(err, &fe) {
		return localFail(env, err)
	}
	if err != nil {
		return exec.UsageExit(env, &cli.UsageError{Msg: err.Error(), Usage: fmt.Sprintf("usage: afs %s [options]\n\nTo see help text, run: afs %s help\n", inv.Cmd.Key(), inv.Cmd.Key())})
	}
	st, credFn, err := runtime.Credentials(env, inv.P.G)
	if err != nil {
		return exec.ReportError(env.Stderr, invokeOp, "default", err)
	}
	if args.teamID == "" {
		args.teamID = st.TeamID
	}
	if args.teamID == "" && args.collabURL == "" {
		return exec.UsageExit(env, &cli.UsageError{Msg: "--team-id is required (or set team_id in the profile, or pass --collab-url)", Usage: "usage: afs collaboration invoke-api [options]\n"})
	}
	key, err := edge.ReadKey(args.keyFile, env.Stdin)
	if err != nil {
		var ks *edge.KeyShapeError
		if errors.As(err, &ks) {
			return exec.UsageExit(env, &cli.UsageError{Msg: err.Error(), Usage: "usage: afs collaboration invoke-api [options]\n"})
		}
		return localFail(env, err)
	}
	cred, err := credFn(ctx)
	if err != nil {
		return exec.ReportError(env.Stderr, invokeOp, st.Profile, err)
	}
	switch {
	case cred.Source == auth.WorkloadSource:
		return localFail(env, errWorkload)
	case cred.Header == "":
		return localFail(env, errors.New("no credentials are configured: set AFS_API_KEY or run "+auth.LoginHint(st.Profile)))
	}
	ec := &edge.Client{HTTP: runtime.Transport(env, inv.P.G), ReadTimeout: runtime.ReadTimeout(inv.P.G), Now: env.Now, Sleep: env.Sleep, Rand: env.Rand}
	if inv.P.G.Debug {
		ec.Debug = env.Stderr
	}
	ec.Debugf("credential source: %s", cred.Source)

	rawOrigin := args.collabURL
	if rawOrigin == "" {
		var rerr error
		if rawOrigin, rerr = webAPIURL(ctx, inv, args.teamID); rerr != nil {
			var off *surfaceOffError
			if errors.As(rerr, &off) {
				return exec.ReportError(env.Stderr, invokeOp, st.Profile, off.apiError())
			}
			return exec.ReportError(env.Stderr, "GetTeamOrchestration", st.Profile, rerr)
		}
	}
	origin, err := edge.Origin(rawOrigin, st.CollabBaseDomain)
	if err != nil {
		return localFail(env, err)
	}
	ec.Debugf("origin: %s", strings.TrimPrefix(strings.TrimPrefix(origin, "https://"), "http://"))

	resp, err := ec.Do(ctx, edge.Request{Origin: origin, Method: args.method, Path: args.path, ContentType: args.ctype, Body: args.body,
		AFSAuth: cred.Header, AFSSource: cred.Source, Key: key, UserAgent: buildinfo.UserAgent()})
	if err != nil {
		return exec.ReportError(env.Stderr, invokeOp, st.Profile, err)
	}
	defer resp.Close()
	return render(env, st.Profile, resp, args)
}

// render writes the response: bodies unchanged to stdout, the edge's envelope as an AWS-style error.
func render(env *cli.Env, profile string, resp *edge.Response, args invokeArgs) int {
	if resp.Status >= 400 {
		peek, err := io.ReadAll(io.LimitReader(resp, maxErrorPeek+1))
		if err != nil {
			return exec.ReportError(env.Stderr, invokeOp, profile, err)
		}
		if len(peek) <= maxErrorPeek && client.IsEnvelope(resp.Header, peek) {
			return exec.ReportError(env.Stderr, invokeOp, profile, client.DecodeError(resp.Status, resp.Header, peek))
		}
		printHead(env, resp, args.include)
		_, _ = env.Stdout.Write(peek)
		if _, err := io.Copy(env.Stdout, resp); err != nil {
			return exec.ReportError(env.Stderr, invokeOp, profile, err)
		}
		w := output.NewSanitizingWriter(env.Stderr)
		_, _ = fmt.Fprintf(w, "\nAn error occurred (HttpError) when calling the %s operation: HTTP %d %s\n", invokeOp, resp.Status, resp.Reason)
		if id := resp.Header.Get("X-Request-Id"); id != "" {
			_, _ = fmt.Fprintf(w, "Request ID: %s\n", id)
		}
		return cli.ExitService
	}
	printHead(env, resp, args.include)
	if _, err := io.Copy(env.Stdout, resp); err != nil {
		return exec.ReportError(env.Stderr, invokeOp, profile, err)
	}
	return cli.ExitOK
}

// printHead writes the status line and headers for --include; Set-Cookie and credentials are masked.
func printHead(env *cli.Env, resp *edge.Response, include bool) {
	if !include {
		return
	}
	_, _ = fmt.Fprintf(env.Stdout, "%s %d %s\n", resp.Proto, resp.Status, resp.Reason)
	for _, l := range redact.Headers(resp.Header) {
		_, _ = fmt.Fprintln(env.Stdout, l)
	}
	_, _ = fmt.Fprintln(env.Stdout)
}

func localFail(env *cli.Env, err error) int {
	_, _ = fmt.Fprintf(output.NewSanitizingWriter(env.Stderr), "\nafs: error: %v\n", err)
	return cli.ExitOther
}

// surfaceOffError is a Web API that has no address to call: the tool has none (Buzz), or the team's exposure
// (or the service's switches) leaves it off. It is reported as an answer of the service would be, exit 254.
type surfaceOffError struct {
	message string
	reason  string // details.reason for the hint ("" when the tool has no Web API)
}

func (e *surfaceOffError) Error() string { return e.message }

func (e *surfaceOffError) apiError() *client.APIError {
	ae := &client.APIError{Status: http.StatusConflict, Code: "SurfaceOff", Message: e.message}
	if e.reason != "" {
		ae.Details, _ = json.Marshal(map[string]string{"reason": e.reason})
	}
	return ae
}

// webAPIURL reads the address of the Web API's effective level with one GetTeamOrchestration (ORC-30): a tool
// without a Web API (Buzz) and a surface that is off are refused with "Web API is off" before any key is sent.
func webAPIURL(ctx context.Context, inv *cli.Invocation, teamID string) (string, error) {
	sess, err := runtime.New(inv.Env, inv.P.G, false)
	if err != nil {
		return "", err
	}
	resp, err := sess.Client.Do(ctx, client.Request{Method: http.MethodGet, Path: "/v1/teams/" + url.PathEscape(teamID) + "/orchestration", Operation: "GetTeamOrchestration"})
	if err != nil {
		return "", err
	}
	var doc struct {
		Tool   string `json:"tool"`
		WebAPI struct {
			Supported *bool    `json:"supported"`
			Effective string   `json:"effective"`
			Reasons   []string `json:"reasons"`
			URL       *string  `json:"url"`
		} `json:"webApi"`
	}
	if err := json.Unmarshal(resp.Body, &doc); err != nil {
		return "", fmt.Errorf("the orchestration answer is not valid JSON: %w", err)
	}
	w := doc.WebAPI
	switch {
	case w.Supported != nil && !*w.Supported:
		return "", &surfaceOffError{message: "Web API is off (N/A for " + orchestration.ToolName(doc.Tool) + ")"}
	case w.URL == nil || *w.URL == "":
		msg := "Web API is off for this team"
		if len(w.Reasons) > 0 {
			msg += " (" + strings.Join(w.Reasons, ", ") + ")"
		}
		return "", &surfaceOffError{message: msg, reason: "surface_disabled"}
	}
	return *w.URL, nil
}

// parseInvoke checks everything that can be checked before any call (CLI-US-25.3, 25.4, 25.9).
func parseInvoke(inv *cli.Invocation) (invokeArgs, error) {
	p := inv.P
	var a invokeArgs
	for _, g := range refusedGlobals {
		if p.G.Set[g] {
			return a, errors.New("invoke-api prints the response unchanged; pipe it to a JSON tool instead. (--" + g + " is not supported)")
		}
	}
	if p.Has("wait") {
		return a, errors.New("invoke-api prints the response unchanged; pipe it to a JSON tool instead. (--wait is not supported)")
	}
	if len(p.Positional) > 0 {
		return a, errors.New("unexpected argument; invoke-api takes options only")
	}
	a.teamID, a.collabURL, a.keyFile, a.include = p.Value("team-id"), p.Value("collab-url"), p.Value("paperclip-key-file"), p.Value("include") == "true"
	if a.keyFile == "" {
		return a, errors.New("the following arguments are required: --paperclip-key-file (a file or - for standard input; the key is never an argument)")
	}
	a.method = strings.ToUpper(p.Value("method"))
	if a.method == "" {
		a.method = http.MethodGet
	}
	if !invokeMethods[a.method] {
		return a, fmt.Errorf("argument --method: invalid choice %q (choose from GET, HEAD, POST, PUT, PATCH, DELETE)", p.Value("method"))
	}
	a.path = p.Value("path")
	if err := checkPath(a.path); err != nil {
		return a, err
	}
	if b := p.Value("body"); p.Has("body") {
		if a.method == http.MethodGet || a.method == http.MethodHead {
			return a, fmt.Errorf("--body cannot be sent with %s", a.method)
		}
		data, err := readBody(b)
		if err != nil {
			return a, err
		}
		a.body = data
		a.ctype = "application/json"
		if ct := p.Value("content-type"); ct != "" {
			if strings.ContainsAny(ct, "\r\n\x00") {
				return a, errors.New("argument --content-type: invalid value")
			}
			a.ctype = ct
		}
	}
	return a, nil
}

// checkPath refuses a path that could leave the team's Paperclip or reach the edge's own paths.
func checkPath(path string) error {
	if path == "" {
		return errors.New("the following arguments are required: --path")
	}
	if !strings.HasPrefix(path, "/") || strings.HasPrefix(path, "//") {
		return errors.New("argument --path: must start with a single /")
	}
	for _, r := range path {
		if r < 0x21 || r == 0x7f || r == '#' || r == '\\' {
			return errors.New("argument --path: must not contain a space, control character, # or backslash")
		}
	}
	if strings.Contains(path, "://") {
		return errors.New("argument --path: must be a path, not a URL")
	}
	bare, _, _ := strings.Cut(path, "?")
	for _, seg := range strings.Split(strings.ToLower(bare), "/") {
		if seg == ".." || seg == "%2e%2e" || seg == ".%2e" || seg == "%2e." {
			return errors.New("argument --path: must not contain .. segments")
		}
	}
	if l := strings.ToLower(bare); l == "/_afs" || strings.HasPrefix(l, "/_afs/") {
		return errors.New("argument --path: /_afs/ paths belong to the AFS edge and are not forwarded")
	}
	return nil
}

// bodyFileError is a --body file that cannot be read: a local file problem (exit 255), not a usage error.
type bodyFileError struct{ err error }

func (e *bodyFileError) Error() string { return "cannot read the --body file: " + e.err.Error() }
func (e *bodyFileError) Unwrap() error { return e.err }

func readBody(v string) ([]byte, error) {
	p, ok := strings.CutPrefix(v, "file://")
	if !ok {
		if p, ok = strings.CutPrefix(v, "fileb://"); !ok {
			return nil, errors.New("argument --body: must be file://PATH or fileb://PATH (standard input is for the key)")
		}
	}
	st, err := os.Stat(p)
	if err != nil {
		return nil, &bodyFileError{err}
	}
	if st.Size() > maxInvokeBody {
		return nil, errors.New("argument --body: the file is over 100 MB, the largest body the service accepts")
	}
	b, err := os.ReadFile(p) //nolint:gosec // user-selected input file
	if err != nil {
		return nil, &bodyFileError{err}
	}
	return b, nil
}
