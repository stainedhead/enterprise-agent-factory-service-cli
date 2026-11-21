package exec

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/client"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/output"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/spec"
)

type token struct {
	T string `json:"t"` // the server's nextToken that starts the page
	S int    `json:"s"` // items of that page already printed
}

func encodeToken(t token) string {
	b, _ := json.Marshal(t)
	return base64.RawURLEncoding.EncodeToString(b)
}

func decodeToken(s string) (token, error) {
	b, err := base64.RawURLEncoding.DecodeString(s)
	var t token
	if err != nil || json.Unmarshal(b, &t) != nil {
		return t, fmt.Errorf("--starting-token is not a token printed by afs")
	}
	return t, nil
}

// checkPaging validates the paging options (usage errors, 252).
func checkPaging(cmd *cli.Command, g cli.Globals, format string) error {
	op := cmd.Op
	if g.Set["page-size"] && (g.PageSize < 1 || g.PageSize > 100) {
		return fmt.Errorf("argument --page-size: must be between 1 and 100, got %d", g.PageSize)
	}
	if g.Set["max-items"] && g.MaxItems < 1 {
		return fmt.Errorf("argument --max-items: must be at least 1, got %d", g.MaxItems)
	}
	if format == "csv" {
		switch {
		case !op.CSV:
			return fmt.Errorf("--output csv is not available for %s", cmd.Key())
		case g.Query != "":
			return fmt.Errorf("--output csv cannot be combined with --query")
		case g.Set["max-items"]:
			return fmt.Errorf("--output csv cannot be combined with --max-items")
		}
	}
	if (g.Set["starting-token"] || g.Set["max-items"] || g.Set["page-size"] || g.Set["no-paginate"]) && !op.Paginated && !op.CSV {
		return fmt.Errorf("%s does not paginate", cmd.Key())
	}
	return nil
}

// pages is how a paginated JSON listing is fetched.
type pages struct {
	do   func(ctx context.Context, q url.Values) (*client.Response, error)
	op   *spec.Operation
	g    cli.Globals
	base url.Values
}

// fetchJSON returns the merged listing body and the last response.
func (p *pages) fetchJSON(ctx context.Context) ([]byte, *client.Response, error) {
	key := p.op.ResultKey
	if key == "" {
		key = "items"
	}
	var start token
	if p.g.StartingToken != "" {
		var err error
		if start, err = decodeToken(p.g.StartingToken); err != nil {
			return nil, nil, &usageFail{err}
		}
	}
	var all []json.RawMessage
	var first *cli.Object
	next := start.T
	skip := start.S
	var last *client.Response
	remaining := -1
	if p.g.MaxItems > 0 {
		remaining = p.g.MaxItems
	}
	outToken := ""
	for {
		q := url.Values{}
		for k, v := range p.base {
			q[k] = v
		}
		if p.g.PageSize > 0 {
			q.Set("maxResults", strconv.Itoa(p.g.PageSize))
		}
		if remaining >= 0 && skip == 0 && (p.g.PageSize == 0 || remaining < p.g.PageSize) {
			q.Set("maxResults", strconv.Itoa(min(remaining, 100)))
		}
		if next != "" {
			q.Set("nextToken", next)
		}
		resp, err := p.do(ctx, q)
		if err != nil {
			return nil, last, err
		}
		last = resp
		obj, err := cli.ParseObject(resp.Body)
		if err != nil {
			return nil, last, fmt.Errorf("the service returned a page that is not a JSON object: %w", err)
		}
		if first == nil {
			first = obj
		}
		var items []json.RawMessage
		if raw, ok := obj.Get(key); ok {
			_ = json.Unmarshal(raw, &items)
		}
		pageToken := next
		if skip > 0 {
			if skip > len(items) {
				skip = len(items)
			}
			items = items[skip:]
		}
		printed := skip
		skip = 0
		if remaining >= 0 && len(items) > remaining {
			items = items[:remaining]
			outToken = encodeToken(token{T: pageToken, S: printed + len(items)})
			all = append(all, items...)
			remaining = 0
			break
		}
		all = append(all, items...)
		if remaining >= 0 {
			remaining -= len(items)
		}
		var nt string
		if raw, ok := obj.Get("nextToken"); ok {
			_ = json.Unmarshal(raw, &nt)
		}
		next = nt
		if nt == "" {
			break
		}
		if p.g.NoPaginate {
			outToken = encodeToken(token{T: nt})
			break
		}
		if remaining == 0 {
			outToken = encodeToken(token{T: nt})
			break
		}
	}
	if all == nil {
		all = []json.RawMessage{}
	}
	merged, _ := json.Marshal(all)
	first.Set(key, merged)
	first.Delete("nextToken")
	if outToken != "" {
		t, _ := json.Marshal(outToken)
		first.Set("NextToken", t)
	}
	body, err := json.Marshal(first)
	return body, last, err
}

// fetchCSV streams every page to w (first page whole, later pages without their header).
func (p *pages) fetchCSV(ctx context.Context, w *output.CSVWriter) (*client.Response, error) {
	var last *client.Response
	next := ""
	if p.g.StartingToken != "" {
		next = p.g.StartingToken
	}
	for {
		q := url.Values{}
		for k, v := range p.base {
			q[k] = v
		}
		if p.g.PageSize > 0 {
			q.Set("maxResults", strconv.Itoa(p.g.PageSize))
		}
		if next != "" {
			q.Set("nextToken", next)
		}
		resp, err := p.do(ctx, q)
		if err != nil {
			return last, err
		}
		last = resp
		if err := w.WritePage(resp.Body); err != nil {
			return last, err
		}
		next = resp.Header.Get("X-Afs-Next-Token")
		if next == "" || p.g.NoPaginate {
			return last, nil
		}
	}
}

// usageFail marks an error that is a usage error (exit 252).
type usageFail struct{ error }
