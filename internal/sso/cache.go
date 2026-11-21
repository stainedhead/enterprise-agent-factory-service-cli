// Package sso implements the OIDC device authorization grant (RFC 8628), the token cache under
// ~/.afs/sso/cache and silent refresh (CLI-4).
package sso

import (
	"crypto/sha1" //nolint:gosec // cache file name, not security
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/config"
)

// Token is a cached SSO session.
type Token struct {
	Issuer       string    `json:"issuer"`
	ClientID     string    `json:"clientId"`
	AccessToken  string    `json:"accessToken"`
	RefreshToken string    `json:"refreshToken,omitempty"`
	ExpiresAt    time.Time `json:"expiresAt"`
	TokenType    string    `json:"tokenType"`
}

// expirySkew treats a token as expired slightly early so it does not lapse in flight.
const expirySkew = 30 * time.Second

// Expired reports whether the access token is expired at now.
func (t *Token) Expired(now time.Time) bool { return !now.Add(expirySkew).Before(t.ExpiresAt) }

// Cache is the directory of cached tokens, one file per issuer and client id.
type Cache struct{ Dir string }

// NewCache returns the cache under home (~/.afs/sso/cache).
func NewCache(home string) *Cache { return &Cache{Dir: filepath.Join(home, ".afs", "sso", "cache")} }

// Path is the cache file of an issuer and client id: sha1(issuer + "\n" + client id).json.
func (c *Cache) Path(issuer, clientID string) string {
	sum := sha1.Sum([]byte(issuer + "\n" + clientID)) //nolint:gosec // file name only
	return filepath.Join(c.Dir, hex.EncodeToString(sum[:])+".json")
}

// Load reads a token; it returns nil with no error when none is cached, and refuses a file
// with a mode wider than 0600.
func (c *Cache) Load(issuer, clientID string) (*Token, error) {
	p := c.Path(issuer, clientID)
	if err := config.CheckPrivate(p); err != nil {
		return nil, err
	}
	b, err := os.ReadFile(p) //nolint:gosec // path derived from a hash
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var t Token
	if err := json.Unmarshal(b, &t); err != nil {
		return nil, fmt.Errorf("the SSO cache file %s is damaged; delete it and run afs sso login: %w", p, err)
	}
	return &t, nil
}

// Save writes the token with mode 0600, atomically.
func (c *Cache) Save(t *Token) error {
	b, err := json.MarshalIndent(t, "", "  ") //nolint:gosec // the cache file holds the token by design and is written 0600
	if err != nil {
		return err
	}
	f := &config.File{Path: c.Path(t.Issuer, t.ClientID)}
	if err := os.MkdirAll(c.Dir, 0o700); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(c.Dir, ".tok-*")
	if err != nil {
		return err
	}
	name := tmp.Name()
	fail := func(err error) error { _ = tmp.Close(); _ = os.Remove(name); return err }
	if err := tmp.Chmod(0o600); err != nil {
		return fail(err)
	}
	if _, err := tmp.Write(b); err != nil {
		return fail(err)
	}
	if err := tmp.Close(); err != nil {
		_ = os.Remove(name)
		return err
	}
	return os.Rename(name, f.Path)
}

// Delete removes a cached token; a missing one is not an error.
func (c *Cache) Delete(issuer, clientID string) error {
	err := os.Remove(c.Path(issuer, clientID))
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	return err
}
