package config

import (
	"fmt"
	"net"
	"net/url"
	"path/filepath"
	"strings"
)

// Source is what Resolve reads from the process: environment and home directory.
type Source struct {
	Getenv func(string) string
	Home   string
}

// Overrides are the command-line values (already validated by the parser); "" means not given.
type Overrides struct {
	Profile     string
	EndpointURL string
	Output      string
	// AllowMissing accepts an explicit profile that has no section yet (afs configure creates it).
	AllowMissing bool
}

// Settings are the resolved settings of one invocation.
type Settings struct {
	Profile         string
	EndpointURL     string
	Output          string
	SSOProviderID   string
	SSOIssuer       string
	SSOClientID     string
	TeamID          string
	ConfigPath      string
	CredentialsPath string
	// Sources names where Output and EndpointURL came from ("option", "environment", "profile", "default").
	Sources map[string]string
}

// Valid output formats.
var outputs = []string{"json", "yaml", "table", "text", "csv"}

// ConfigPath and CredentialsPath honour AFS_CONFIG_FILE and AFS_SHARED_CREDENTIALS_FILE (D-C15).
func configPath(s Source) string {
	if p := s.Getenv("AFS_CONFIG_FILE"); p != "" {
		return p
	}
	return filepath.Join(s.Home, ".afs", "config")
}

func credentialsPath(s Source) string {
	if p := s.Getenv("AFS_SHARED_CREDENTIALS_FILE"); p != "" {
		return p
	}
	return filepath.Join(s.Home, ".afs", "credentials")
}

// ConfigSection is the section header of a profile in ~/.afs/config.
func ConfigSection(profile string) string {
	if profile == "default" {
		return "default"
	}
	return "profile " + profile
}

// ProfileName resolves the profile: --profile, then AFS_PROFILE, then "default".
func ProfileName(ov Overrides, s Source) (name string, explicit bool) {
	if ov.Profile != "" {
		return ov.Profile, true
	}
	if p := s.Getenv("AFS_PROFILE"); p != "" {
		return p, true
	}
	return "default", false
}

// Resolve applies the precedence option, environment, profile, default (CLI-3, D-C35).
// Errors are configuration errors (exit 255).
func Resolve(ov Overrides, s Source) (Settings, error) {
	name, explicit := ProfileName(ov, s)
	st := Settings{Profile: name, ConfigPath: configPath(s), CredentialsPath: credentialsPath(s), Sources: map[string]string{}}
	f, err := Load(st.ConfigPath)
	if err != nil {
		return st, fmt.Errorf("reading %s: %w", st.ConfigPath, err)
	}
	sec := ConfigSection(name)
	if explicit && name != "default" && !ov.AllowMissing && !f.HasSection(sec) {
		return st, fmt.Errorf("The config profile (%s) could not be found", name) //nolint:staticcheck // AWS CLI wording
	}
	get := func(k string) string { v, _ := f.Get(sec, k); return v }

	st.SSOProviderID, st.SSOIssuer, st.SSOClientID, st.TeamID = get("sso_provider_id"), get("sso_issuer"), get("sso_client_id"), get("team_id")

	switch {
	case ov.EndpointURL != "":
		st.EndpointURL, st.Sources["endpoint_url"] = ov.EndpointURL, "option"
	case s.Getenv("AFS_ENDPOINT_URL") != "":
		st.EndpointURL, st.Sources["endpoint_url"] = s.Getenv("AFS_ENDPOINT_URL"), "environment"
	case get("endpoint_url") != "":
		st.EndpointURL, st.Sources["endpoint_url"] = get("endpoint_url"), "profile"
	default:
		st.Sources["endpoint_url"] = "default"
	}

	type cand struct{ val, src, where string }
	for _, c := range []cand{
		{ov.Output, "option", "--output"},
		{s.Getenv("AFS_DEFAULT_OUTPUT"), "environment", "AFS_DEFAULT_OUTPUT"},
		{get("output"), "profile", fmt.Sprintf("output in profile %s", name)},
	} {
		if c.val == "" {
			continue
		}
		if !contains(outputs, c.val) {
			return st, fmt.Errorf("invalid output format %q from %s (choose from %s)", c.val, c.where, strings.Join(outputs, ", "))
		}
		st.Output, st.Sources["output"] = c.val, c.src
		break
	}
	if st.Output == "" {
		st.Output, st.Sources["output"] = "json", "default"
	}
	return st, nil
}

func contains(l []string, s string) bool {
	for _, x := range l {
		if x == s {
			return true
		}
	}
	return false
}

// ValidateEndpoint requires https; http is accepted only for loopback hosts (D-C33).
func ValidateEndpoint(raw string) error {
	if raw == "" {
		return fmt.Errorf("no endpoint configured; run afs configure, set AFS_ENDPOINT_URL or pass --endpoint-url")
	}
	u, err := url.Parse(raw)
	if err != nil || u.Host == "" {
		return fmt.Errorf("endpoint_url %q is not a valid URL", raw)
	}
	switch u.Scheme {
	case "https":
		return nil
	case "http":
		h := u.Hostname()
		if h == "localhost" {
			return nil
		}
		if ip := net.ParseIP(h); ip != nil && ip.IsLoopback() {
			return nil
		}
		return fmt.Errorf("endpoint_url must use https (http is allowed only for loopback hosts): %s", raw)
	}
	return fmt.Errorf("endpoint_url must use https: %s", raw)
}

// WriteProfile sets keys of a profile in ~/.afs/config, leaving every other profile untouched.
func WriteProfile(s Source, profile string, kv map[string]string, order []string) (string, error) {
	path := configPath(s)
	f, err := Load(path)
	if err != nil {
		return path, err
	}
	for _, k := range order {
		if v, ok := kv[k]; ok && v != "" {
			f.Set(ConfigSection(profile), k, v)
		}
	}
	return path, f.Save()
}

// CredentialsFile loads ~/.afs/credentials after checking its mode.
func CredentialsFile(s Source) (*File, error) {
	p := credentialsPath(s)
	if err := CheckPrivate(p); err != nil {
		return nil, err
	}
	return Load(p)
}

// ConfigFilePath is the resolved path of the config file.
func ConfigFilePath(s Source) string { return configPath(s) }
