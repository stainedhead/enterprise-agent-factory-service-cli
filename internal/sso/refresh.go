package sso

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/auth"
)

// Refresh exchanges a refresh token for a new access token.
func Refresh(ctx context.Context, hc HTTPDoer, now time.Time, t *Token) (*Token, error) {
	disc, err := Discover(ctx, hc, t.Issuer)
	if err != nil {
		return nil, err
	}
	var tr tokenResponse
	err = postForm(ctx, hc, disc.TokenEndpoint, url.Values{"grant_type": {"refresh_token"}, "refresh_token": {t.RefreshToken}, "client_id": {t.ClientID}}, &tr)
	if err != nil {
		return nil, fmt.Errorf("refresh rejected: %w", err)
	}
	return tr.token(t.Issuer, t.ClientID, now, t.RefreshToken), nil
}

// Source is the cached SSO session of a profile as a credential source (auth.TokenSource).
type Source struct {
	Cache    *Cache
	Issuer   string
	ClientID string
	Profile  string
	HTTP     HTTPDoer
	Now      func() time.Time
}

// Token returns the access token, refreshing silently when it has expired (CLI-US-1.2). A missing
// cache, or an expired one with no usable refresh token, is an auth.LoginRequiredError.
func (s *Source) Token(ctx context.Context) (string, bool, error) {
	t, err := s.Cache.Load(s.Issuer, s.ClientID)
	if err != nil {
		return "", false, err
	}
	if t == nil {
		return "", false, &auth.LoginRequiredError{Profile: s.Profile}
	}
	if !t.Expired(s.Now()) {
		return t.AccessToken, true, nil
	}
	if t.RefreshToken == "" {
		return "", false, &auth.LoginRequiredError{Profile: s.Profile, Reason: "Your session has expired"}
	}
	nt, err := Refresh(ctx, s.HTTP, s.Now(), t)
	if err != nil {
		if ctx.Err() != nil {
			return "", false, ctx.Err()
		}
		return "", false, &auth.LoginRequiredError{Profile: s.Profile, Reason: "Your session has expired and could not be refreshed (" + err.Error() + ")"}
	}
	if err := s.Cache.Save(nt); err != nil {
		return "", false, err
	}
	return nt.AccessToken, true, nil
}
