package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func src(t *testing.T, env map[string]string) (Source, string) {
	t.Helper()
	home := t.TempDir()
	return Source{Getenv: func(k string) string { return env[k] }, Home: home}, home
}

func writeConfig(t *testing.T, home, body string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Join(home, ".afs"), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(home, ".afs", "config"), []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
}

// 5.1-configuration: option > AFS_ENDPOINT_URL > profile endpoint_url > default.
func TestPrecedenceEndpointURL(t *testing.T) {
	env := map[string]string{}
	s, home := src(t, env)
	writeConfig(t, home, "[default]\nendpoint_url = https://profile.example\n")
	steps := []struct {
		ov   string
		env  string
		want string
		from string
	}{
		{"https://option.example", "https://env.example", "https://option.example", "option"},
		{"", "https://env.example", "https://env.example", "environment"},
		{"", "", "https://profile.example", "profile"},
	}
	for _, st := range steps {
		env["AFS_ENDPOINT_URL"] = st.env
		got, err := Resolve(Overrides{EndpointURL: st.ov}, s)
		if err != nil || got.EndpointURL != st.want || got.Sources["endpoint_url"] != st.from {
			t.Errorf("%+v: got %q from %q err %v", st, got.EndpointURL, got.Sources["endpoint_url"], err)
		}
	}
	s2, _ := src(t, map[string]string{})
	got, err := Resolve(Overrides{}, s2)
	if err != nil || got.EndpointURL != "" || got.Sources["endpoint_url"] != "default" {
		t.Errorf("default: %+v %v", got, err)
	}
}

func TestAFSProfileSelectsProfile(t *testing.T) {
	env := map[string]string{"AFS_PROFILE": "p2"}
	s, home := src(t, env)
	writeConfig(t, home, "[default]\nendpoint_url = https://d.example\n\n[profile p2]\nendpoint_url = https://p2.example\nteam_id = team_1\nsso_provider_id = idp_1\nsso_issuer = https://iss\nsso_client_id = cid\n")
	got, err := Resolve(Overrides{}, s)
	if err != nil || got.Profile != "p2" || got.EndpointURL != "https://p2.example" || got.TeamID != "team_1" || got.SSOIssuer != "https://iss" || got.SSOClientID != "cid" || got.SSOProviderID != "idp_1" {
		t.Fatalf("got %+v err %v", got, err)
	}
	got, _ = Resolve(Overrides{Profile: "default"}, s)
	if got.Profile != "default" || got.EndpointURL != "https://d.example" {
		t.Errorf("--profile default: %+v", got)
	}
}

// CLI-3: output defaults to json; AFS_DEFAULT_OUTPUT and the profile are honoured.
func TestDefaultOutputJSON(t *testing.T) {
	env := map[string]string{}
	s, home := src(t, env)
	if got, _ := Resolve(Overrides{}, s); got.Output != "json" || got.Sources["output"] != "default" {
		t.Errorf("default output = %+v", got)
	}
	writeConfig(t, home, "[default]\noutput = table\n")
	if got, _ := Resolve(Overrides{}, s); got.Output != "table" || got.Sources["output"] != "profile" {
		t.Errorf("profile output = %+v", got)
	}
	env["AFS_DEFAULT_OUTPUT"] = "yaml"
	if got, _ := Resolve(Overrides{}, s); got.Output != "yaml" {
		t.Errorf("env output = %+v", got)
	}
	if got, _ := Resolve(Overrides{Output: "text"}, s); got.Output != "text" || got.Sources["output"] != "option" {
		t.Errorf("option output = %+v", got)
	}
}

// D-C35: an invalid value from the environment or a profile is a configuration error (255).
func TestInvalidOutputFromEnvOrProfile255(t *testing.T) {
	env := map[string]string{"AFS_DEFAULT_OUTPUT": "xml"}
	s, home := src(t, env)
	if _, err := Resolve(Overrides{}, s); err == nil || !strings.Contains(err.Error(), "AFS_DEFAULT_OUTPUT") {
		t.Errorf("env: %v", err)
	}
	delete(env, "AFS_DEFAULT_OUTPUT")
	writeConfig(t, home, "[default]\noutput = bad\n")
	if _, err := Resolve(Overrides{}, s); err == nil || !strings.Contains(err.Error(), "profile default") {
		t.Errorf("profile: %v", err)
	}
}

func TestUnknownProfile255(t *testing.T) {
	s, home := src(t, map[string]string{})
	writeConfig(t, home, "[default]\n")
	_, err := Resolve(Overrides{Profile: "nope"}, s)
	if err == nil || err.Error() != "The config profile (nope) could not be found" {
		t.Errorf("err = %v", err)
	}
	s2, _ := src(t, map[string]string{"AFS_PROFILE": "nope"})
	if _, err := Resolve(Overrides{}, s2); err == nil {
		t.Error("AFS_PROFILE=nope accepted")
	}
}

// CLI-3 / D-C10: team_id comes from the profile; the resolver never applies it itself.
func TestTeamIDDefaultsPathOnly(t *testing.T) {
	s, home := src(t, map[string]string{})
	writeConfig(t, home, "[default]\nteam_id = team_01\n")
	if got, _ := Resolve(Overrides{}, s); got.TeamID != "team_01" {
		t.Errorf("team id = %q", got.TeamID)
	}
}

// 5.1-configuration bullet 4: writing keeps other profiles and comments, mode 0600.
func TestWriteKeepsOtherProfilesAnd0600(t *testing.T) {
	s, home := src(t, map[string]string{})
	writeConfig(t, home, "# my comment\n[default]\nendpoint_url = https://d.example\n\n[profile other]\noutput = table\n")
	path, err := WriteProfile(s, "new", map[string]string{"endpoint_url": "https://n.example", "output": "yaml", "team_id": ""}, []string{"endpoint_url", "output", "team_id"})
	if err != nil {
		t.Fatal(err)
	}
	b, _ := os.ReadFile(path) //nolint:gosec // test file
	got := string(b)
	for _, want := range []string{"# my comment", "[default]\nendpoint_url = https://d.example", "[profile other]\noutput = table", "[profile new]\nendpoint_url = https://n.example\noutput = yaml"} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q in:\n%s", want, got)
		}
	}
	if strings.Contains(got, "team_id") {
		t.Error("empty value written")
	}
	if st, _ := os.Stat(path); st.Mode().Perm() != 0o600 {
		t.Errorf("mode = %o", st.Mode().Perm())
	}
	// update an existing key in place
	if _, err := WriteProfile(s, "default", map[string]string{"endpoint_url": "https://changed"}, []string{"endpoint_url"}); err != nil {
		t.Fatal(err)
	}
	b, _ = os.ReadFile(path) //nolint:gosec // test file
	if !strings.Contains(string(b), "[default]\nendpoint_url = https://changed\n") || !strings.Contains(string(b), "[profile other]") {
		t.Errorf("update lost data:\n%s", b)
	}
}

func TestConfigAndCredentialsFileEnvOverrides(t *testing.T) {
	dir := t.TempDir()
	cfg := filepath.Join(dir, "c.ini")
	cred := filepath.Join(dir, "k.ini")
	if err := os.WriteFile(cfg, []byte("[default]\nendpoint_url = https://via-env\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(cred, []byte("[default]\napi_key = afsk_x\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	vars := map[string]string{"AFS_CONFIG_FILE": cfg, "AFS_SHARED_CREDENTIALS_FILE": cred}
	s := Source{Getenv: func(k string) string { return vars[k] }, Home: "/nonexistent"}
	got, err := Resolve(Overrides{}, s)
	if err != nil || got.EndpointURL != "https://via-env" || got.CredentialsPath != cred || ConfigFilePath(s) != cfg {
		t.Fatalf("got %+v err %v", got, err)
	}
	f, err := CredentialsFile(s)
	if err != nil {
		t.Fatal(err)
	}
	if v, ok := f.Get("default", "api_key"); !ok || v != "afsk_x" {
		t.Errorf("api_key = %q", v)
	}
}

// D-C33: https required; http only for loopback.
func TestHTTPEndpointOnlyLoopback(t *testing.T) {
	for url, ok := range map[string]bool{
		"https://afs.example.com": true, "http://localhost:8080": true, "http://127.0.0.1:9": true, "http://[::1]:1": true,
		"http://afs.example.com": false, "ftp://x": false, "": false, "https://": false, "http://10.0.0.1": false,
	} {
		if err := ValidateEndpoint(url); (err == nil) != ok {
			t.Errorf("%q: err = %v, want ok=%v", url, err, ok)
		}
	}
}

func TestCheckPrivate(t *testing.T) {
	p := filepath.Join(t.TempDir(), "f")
	if err := CheckPrivate(p); err != nil {
		t.Errorf("missing file: %v", err)
	}
	if err := os.WriteFile(p, []byte("x"), 0o644); err != nil { //nolint:gosec // test
		t.Fatal(err)
	}
	if err := os.Chmod(p, 0o644); err != nil { //nolint:gosec // test
		t.Fatal(err)
	}
	if err := CheckPrivate(p); err == nil || !strings.Contains(err.Error(), "chmod 600 "+p) {
		t.Errorf("0644 accepted: %v", err)
	}
	for _, m := range []os.FileMode{0o600, 0o400} {
		_ = os.Chmod(p, m)
		if err := CheckPrivate(p); err != nil {
			t.Errorf("%o refused: %v", m, err)
		}
	}
}

func TestINIEdgeCases(t *testing.T) {
	f := &File{Path: filepath.Join(t.TempDir(), "sub", "x")}
	f.Set("a", "k", "1")
	f.Set("b", "k", "2")
	f.Set("a", "j", "3")
	if f.String() != "[a]\nk = 1\nj = 3\n\n[b]\nk = 2\n" {
		t.Errorf("render:\n%s", f.String())
	}
	if _, ok := f.Get("a", "nope"); ok {
		t.Error("phantom key")
	}
	if f.HasSection("zz") {
		t.Error("phantom section")
	}
	if err := f.Save(); err != nil {
		t.Fatal(err)
	}
	g, err := Load(f.Path)
	if err != nil || g.String() != f.String() {
		t.Errorf("reload: %v %q", err, g.String())
	}
	if (&File{}).String() != "" {
		t.Error("empty render")
	}
	if _, err := Load(t.TempDir()); err == nil {
		t.Error("loading a directory succeeded")
	}
	crlf := filepath.Join(t.TempDir(), "w")
	_ = os.WriteFile(crlf, []byte("[x]\r\nk=v\r\n; c\r\n"), 0o600)
	h, _ := Load(crlf)
	if v, _ := h.Get("x", "k"); v != "v" {
		t.Errorf("crlf value = %q", v)
	}
	if ConfigSection("default") != "default" || ConfigSection("p") != "profile p" {
		t.Error("section names")
	}
}
