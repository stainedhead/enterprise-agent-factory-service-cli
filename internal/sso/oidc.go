package sso

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/config"
)

// HTTPDoer is the part of http.Client the flow uses.
type HTTPDoer interface {
	Do(*http.Request) (*http.Response, error)
}

// Discovery is the part of the OIDC discovery document the flow needs.
type Discovery struct {
	Issuer                      string `json:"issuer"`
	DeviceAuthorizationEndpoint string `json:"device_authorization_endpoint"`
	TokenEndpoint               string `json:"token_endpoint"`
}

// Discover fetches {issuer}/.well-known/openid-configuration.
//
// The issuer must be https (http only for a loopback host), the document must name that same issuer
// (OIDC Discovery 4.3), and the endpoints it names must be https, or loopback http for a loopback
// issuer. Nothing carrying a device code or token is sent otherwise.
func Discover(ctx context.Context, hc HTTPDoer, issuer string) (*Discovery, error) {
	if err := config.ValidateIssuer("sso_issuer", issuer); err != nil {
		return nil, err
	}
	u := strings.TrimRight(issuer, "/") + "/.well-known/openid-configuration"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	var d Discovery
	if err := doJSON(hc, req, &d); err != nil {
		return nil, fmt.Errorf("OIDC discovery for %s failed: %w", issuer, err)
	}
	if strings.TrimRight(d.Issuer, "/") != strings.TrimRight(issuer, "/") {
		return nil, fmt.Errorf("OIDC discovery for %s names a different issuer; refusing to continue", issuer)
	}
	if d.TokenEndpoint == "" {
		return nil, fmt.Errorf("OIDC discovery for %s has no token_endpoint", issuer)
	}
	loopbackIssuer := false
	if u, err := url.Parse(issuer); err == nil {
		loopbackIssuer = config.IsLoopbackHost(u.Hostname())
	}
	if err := checkEndpoint("token_endpoint", d.TokenEndpoint, loopbackIssuer); err != nil {
		return nil, err
	}
	if d.DeviceAuthorizationEndpoint != "" {
		if err := checkEndpoint("device_authorization_endpoint", d.DeviceAuthorizationEndpoint, loopbackIssuer); err != nil {
			return nil, err
		}
	}
	return &d, nil
}

// checkEndpoint requires an absolute https URL without credentials; http is allowed only when the
// issuer and the endpoint are both loopback.
func checkEndpoint(key, raw string, loopbackIssuer bool) error {
	u, err := url.Parse(raw)
	if err != nil || u.Host == "" || u.User != nil {
		return fmt.Errorf("the identity provider's %s is not a valid https URL", key)
	}
	switch {
	case u.Scheme == "https":
		return nil
	case u.Scheme == "http" && loopbackIssuer && config.IsLoopbackHost(u.Hostname()):
		return nil
	}
	return fmt.Errorf("the identity provider's %s must use https", key)
}

// oauthError is an RFC 6749 error response.
type oauthError struct {
	Code        string `json:"error"`
	Description string `json:"error_description"`
}

func (e *oauthError) Error() string {
	if e.Description != "" {
		return e.Code + ": " + e.Description
	}
	return e.Code
}

func doJSON(hc HTTPDoer, req *http.Request, out any) error {
	resp, err := hc.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	b, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return err
	}
	if resp.StatusCode >= 300 {
		var oe oauthError
		if json.Unmarshal(b, &oe) == nil && oe.Code != "" {
			return &oe
		}
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	return json.Unmarshal(b, out)
}

func postForm(ctx context.Context, hc HTTPDoer, endpoint string, form url.Values, out any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	return doJSON(hc, req, out)
}

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
}

// ErrNoAccessToken is returned when the identity provider answers 200 without a usable access token.
var ErrNoAccessToken = errors.New("the identity provider returned no access token")

func (r tokenResponse) token(issuer, clientID string, now time.Time, prevRefresh string) (*Token, error) {
	if r.AccessToken == "" {
		return nil, ErrNoAccessToken
	}
	if r.TokenType != "" && !strings.EqualFold(r.TokenType, "Bearer") {
		return nil, fmt.Errorf("the identity provider returned an unsupported token_type; afs needs Bearer")
	}
	t := &Token{Issuer: issuer, ClientID: clientID, AccessToken: r.AccessToken, RefreshToken: r.RefreshToken, TokenType: r.TokenType}
	if t.TokenType == "" {
		t.TokenType = "Bearer"
	}
	if t.RefreshToken == "" {
		t.RefreshToken = prevRefresh // the IdP may not rotate refresh tokens
	}
	exp := r.ExpiresIn
	if exp <= 0 {
		exp = 3600
	}
	t.ExpiresAt = now.Add(time.Duration(exp) * time.Second)
	return t, nil
}
