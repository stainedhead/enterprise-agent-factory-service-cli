package modelendpoints

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/output"
)

func init() { registrations = append(registrations, registerStatements) }

// maxStatementBytes bounds the statement file the CLI reads (10 000 rows of the contract fit well inside it).
const maxStatementBytes = 16 << 20

var statementColumns = []output.Column{
	{Header: "Id", Path: "id"}, {Header: "Endpoint", Path: "endpointId"}, {Header: "Period", Path: "period"},
	{Header: "Source", Path: "source"}, {Header: "Rows", Path: "rowCount"}, {Header: "TotalUsd", Path: "totalUsd"},
	{Header: "ImportedAt", Path: "importedAt"},
}

// statementOption is the file of the statement: the provider's CSV (columns teamId,model,amountUsd,tokens) or JSON.
var statementOption = cli.Option{Flag: "--statement", In: "local", Kind: cli.String,
	Desc: "The provider's statement, file://PATH only: a CSV file (columns teamId,model,amountUsd,tokens) is sent as text/csv; a JSON file is either the whole body {period, source, rows} or the array of rows."}

// checkStatement refuses, before anything is read, a statement that is not a file or is too large.
func checkStatement(p *cli.Parsed) error {
	if !p.Has("statement") {
		return nil
	}
	path, ok := strings.CutPrefix(p.Value("statement"), "file://")
	if !ok || path == "" {
		return errors.New("argument --statement: give the statement as file://PATH (a CSV or JSON file)")
	}
	if st, err := os.Stat(path); err == nil && st.Size() > maxStatementBytes {
		return fmt.Errorf("argument --statement: %s is %d bytes; the limit is %d bytes", path, st.Size(), maxStatementBytes)
	}
	return nil
}

// buildStatement reads --statement and chooses the body: a CSV file is sent unchanged as text/csv with period and
// source in the query; JSON (a file, --rows or --cli-input-json) is the body, with period and source moved into it.
// Options given on the command line win over the file's members.
func buildStatement(p *cli.Parsed, parts *cli.RequestParts) error {
	if p.Has("statement") {
		path := strings.TrimPrefix(p.Value("statement"), "file://")
		b, err := os.ReadFile(path) //nolint:gosec // user-selected input file
		if err != nil {
			return fmt.Errorf("argument --statement: %w", err)
		}
		trimmed := bytes.TrimSpace(b)
		switch {
		case len(trimmed) == 0:
			return errors.New("argument --statement: the file is empty")
		case trimmed[0] == '{':
			obj, err := cli.ParseObject(trimmed)
			if err != nil {
				return fmt.Errorf("argument --statement: %w", err)
			}
			for _, k := range obj.Keys() {
				v, _ := obj.Get(k)
				if (k == "period" || k == "source") && parts.Query.Has(k) {
					continue // the command line wins
				}
				if k == "period" || k == "source" {
					var s string
					if json.Unmarshal(v, &s) != nil {
						return fmt.Errorf("argument --statement: %s must be a string", k)
					}
					parts.Query.Set(k, s)
					continue
				}
				parts.Body.Set(k, v)
			}
		case trimmed[0] == '[':
			if !json.Valid(trimmed) {
				return errors.New("argument --statement: the file is not valid JSON")
			}
			parts.Body.Set("rows", json.RawMessage(trimmed))
		default:
			if !parts.Query.Has("period") || !parts.Query.Has("source") {
				return errors.New("the following arguments are required with a CSV statement: --period, --source")
			}
			parts.Raw, parts.ContentType = b, "text/csv"
			return nil
		}
	}
	for _, k := range []string{"period", "source"} {
		if parts.Query.Has(k) {
			v, _ := json.Marshal(parts.Query.Get(k))
			parts.Body.Set(k, v)
			parts.Query.Del(k)
		}
	}
	var missing []string
	for _, m := range []struct{ wire, flag string }{{"period", "--period"}, {"source", "--source"}, {"rows", "--statement"}} {
		if _, ok := parts.Body.Get(m.wire); !ok {
			missing = append(missing, m.flag)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("the following arguments are required: %s", strings.Join(missing, ", "))
	}
	return nil
}

func registerStatements(r *cli.Registry) {
	r.Override(cli.Override{Op: "ImportProviderStatement", Columns: statementColumns,
		Extra: []cli.Option{statementOption}, Precheck: checkStatement, Build: buildStatement,
		Optional: []string{"rows"}, Hidden: []string{"rows"},
		Desc: map[string]string{
			"period": "The statement's month, YYYY-MM.",
			"source": "The statement's format: nous_statement or opencode_csv.",
		},
		Notes: []string{"Imports a provider's monthly statement so spend can be reconciled per team (variance over 5% is alerted). A route shows Not reconciled until its first import."},
		Examples: []cli.Example{
			{Title: "Import OpenCode's usage CSV for November", Command: "afs model-endpoints import-provider-statement --endpoint-id " + exEndpoint + " --period 2026-11 --source opencode_csv --statement file://opencode-2026-11.csv"},
			{Title: "Import a statement prepared as JSON", Command: "afs model-endpoints import-provider-statement --endpoint-id " + exEndpoint + " --statement file://statement-2026-11.json"},
		}})
}
