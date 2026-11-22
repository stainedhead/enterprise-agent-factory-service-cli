// Package edge is the client of a team's Paperclip edge (API contract section 4.4c) used by
// "afs collaboration invoke-api": origin validation, the Paperclip key reader and one HTTP request with the
// narrow retry rules of CLI-US-25.8. It holds no business rule: the edge and Paperclip decide everything.
package edge

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"strings"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/config"
)

// OriginError is a refused origin; the message names the rule and never a credential.
type OriginError struct{ Rule string }

func (e *OriginError) Error() string { return "the collaboration origin is refused: " + e.Rule }

// ErrNotConfigured is a null restApiUrl: the service has no collaboration base domain.
var ErrNotConfigured = errors.New("Public collaboration is not configured for this service (collab.base_domain is unset).") //nolint:staticcheck // user-facing sentence from the PRD

func refuse(format string, a ...any) error { return &OriginError{Rule: fmt.Sprintf(format, a...)} }

// Origin validates a restApiUrl or --collab-url and returns "scheme://host": only the scheme and the host
// are used (PC-CLI-D4). Rules (CLI-US-25.5): https (http only for a loopback host, as endpoint_url), no
// userinfo, query or fragment, no port other than 443, a DNS name rather than an IP literal, and, when
// baseDomain is set, exactly one label under it (PC-CLI-D5). A loopback http origin is exempt from the
// port, address and domain rules so tests and local stubs can be reached.
func Origin(raw, baseDomain string) (string, error) {
	if strings.TrimSpace(raw) == "" {
		return "", refuse("the URL is empty")
	}
	if strings.ContainsAny(raw, "?#") {
		return "", refuse("the URL must not contain a query or a fragment")
	}
	u, err := url.Parse(raw)
	if err != nil || u.Hostname() == "" {
		return "", refuse("the URL has no host")
	}
	if u.User != nil {
		return "", refuse("the URL must not contain credentials")
	}
	host := strings.ToLower(u.Hostname())
	loopback := u.Scheme == "http" && config.IsLoopbackHost(host)
	switch {
	case u.Scheme == "https", loopback:
	case u.Scheme == "http":
		return "", refuse("the scheme must be https")
	default:
		return "", refuse("the scheme must be https")
	}
	if loopback {
		return "http://" + u.Host, nil
	}
	if p := u.Port(); p != "" && p != "443" {
		return "", refuse("the port must be 443")
	}
	if net.ParseIP(host) != nil || strings.HasPrefix(u.Host, "[") {
		return "", refuse("the host must be a DNS name, not an IP address")
	}
	if baseDomain != "" {
		d := strings.ToLower(strings.Trim(baseDomain, "."))
		label, ok := strings.CutSuffix(host, "."+d)
		if !ok || label == "" || strings.Contains(label, ".") {
			return "", refuse("the host must be one label under the profile's collab_base_domain")
		}
	}
	return "https://" + host, nil
}
