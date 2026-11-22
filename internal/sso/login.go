package sso

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"time"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/output"
)

// Errors of the device flow; all are exit 255 (not API errors).
var (
	ErrDeviceExpired = errors.New("the sign-in code expired before it was approved; run afs sso login again")
	ErrAccessDenied  = errors.New("the sign-in was denied at the identity provider")
)

// LoginOptions inject the pieces of the flow that tests replace.
type LoginOptions struct {
	HTTP  HTTPDoer
	Out   io.Writer // verification URL and code (stderr)
	Now   func() time.Time
	Sleep func(ctx context.Context, d time.Duration) error
	Scope string // default "openid offline_access"
}

type deviceResponse struct {
	DeviceCode              string `json:"device_code"`
	UserCode                string `json:"user_code"`
	VerificationURI         string `json:"verification_uri"`
	VerificationURIComplete string `json:"verification_uri_complete"`
	ExpiresIn               int    `json:"expires_in"`
	Interval                int    `json:"interval"`
}

const grantDevice = "urn:ietf:params:oauth:grant-type:device_code"

// Login runs the device authorization grant against the issuer and returns the token.
func Login(ctx context.Context, o LoginOptions, issuer, clientID string) (*Token, error) {
	disc, err := Discover(ctx, o.HTTP, issuer)
	if err != nil {
		return nil, err
	}
	if disc.DeviceAuthorizationEndpoint == "" {
		return nil, fmt.Errorf("the identity provider %s does not offer the device authorization grant; enable it for client %s", issuer, clientID)
	}
	scope := o.Scope
	if scope == "" {
		scope = "openid offline_access"
	}
	var dev deviceResponse
	if err := postForm(ctx, o.HTTP, disc.DeviceAuthorizationEndpoint, url.Values{"client_id": {clientID}, "scope": {scope}}, &dev); err != nil {
		return nil, fmt.Errorf("device authorization failed: %w", err)
	}
	uri := dev.VerificationURIComplete
	if uri == "" {
		uri = dev.VerificationURI
	}
	_, _ = fmt.Fprintf(o.Out, "To sign in, open this URL in a browser:\n  %s\nand enter the code: %s\nWaiting for approval...\n", output.Sanitize(dev.VerificationURI), output.Sanitize(dev.UserCode))
	if dev.VerificationURIComplete != "" {
		_, _ = fmt.Fprintf(o.Out, "Or open this URL directly:\n  %s\n", output.Sanitize(uri))
	}
	interval := time.Duration(dev.Interval) * time.Second
	if interval <= 0 {
		interval = 5 * time.Second
	}
	deadline := o.Now().Add(time.Duration(dev.ExpiresIn) * time.Second)
	if dev.ExpiresIn <= 0 {
		deadline = o.Now().Add(15 * time.Minute)
	}
	for {
		if err := o.Sleep(ctx, interval); err != nil {
			return nil, err
		}
		if !o.Now().Before(deadline) {
			return nil, ErrDeviceExpired
		}
		var tr tokenResponse
		err := postForm(ctx, o.HTTP, disc.TokenEndpoint, url.Values{"grant_type": {grantDevice}, "device_code": {dev.DeviceCode}, "client_id": {clientID}}, &tr)
		var oe *oauthError
		switch {
		case err == nil:
			return tr.token(issuer, clientID, o.Now(), "")
		case errors.As(err, &oe) && oe.Code == "authorization_pending":
			continue
		case errors.As(err, &oe) && oe.Code == "slow_down":
			interval += 5 * time.Second
		case errors.As(err, &oe) && oe.Code == "expired_token":
			return nil, ErrDeviceExpired
		case errors.As(err, &oe) && oe.Code == "access_denied":
			return nil, ErrAccessDenied
		default:
			return nil, fmt.Errorf("token request failed: %w", err)
		}
	}
}
