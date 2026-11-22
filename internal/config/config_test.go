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
	_ = f.Set("a", "k", "1")
	_ = f.Set("b", "k", "2")
	_ = f.Set("a", "j", "3")
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

func TestSetRejectsInjection(t *testing.T) {
	cases := []struct{ name, section, key, value string }{
		{"newline in value", "profile p", "k", "a\n[profile prod]\nendpoint_url = https://attacker.example"},
		{"carriage return", "profile p", "k", "a\rb"},
		{"tab is a control char", "profile p", "k", "a\tb"},
		{"NUL", "profile p", "k", "a\x00b"},
		{"DEL", "profile p", "k", "a\x7fb"},
		{"newline in key", "profile p", "k\nx", "v"},
		{"equals in key", "profile p", "k=x", "v"},
		{"bracket in key", "profile p", "[k", "v"},
		{"empty key", "profile p", "", "v"},
		{"newline in section", "profile p\n[x", "k", "v"},
		{"bracket in section", "profile ]p", "k", "v"},
		{"empty section", "", "k", "v"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			f := &File{Path: filepath.Join(t.TempDir(), "config")}
			if err := f.Set("default", "endpoint_url", "https://ok.example"); err != nil {
				t.Fatal(err)
			}
			before := f.String()
			if err := f.Set(c.section, c.key, c.value); err == nil {
				t.Fatalf("accepted %+v", c)
			}
			if f.String() != before {
				t.Errorf("file changed:\n%s", f.String())
			}
		})
	}
}

func TestValidateIssuer(t *testing.T) {
	ok := []string{"https://idp.example.com", "https://idp.example.com/realms/a", "http://127.0.0.1:8080", "http://localhost/x", "https://idp.example.com/"}
	bad := []string{"", "http://idp.example.com", "https://u:p@idp.example.com", "https://idp.example.com?x=1", "https://idp.example.com#f", "idp.example.com", "https://", "ftp://idp.example.com", "https://idp.example.com\nx"}
	for _, s := range ok {
		if err := ValidateIssuer("sso_issuer", s); err != nil {
			t.Errorf("%q: %v", s, err)
		}
	}
	for _, s := range bad {
		err := ValidateIssuer("sso_issuer", s)
		if err == nil {
			t.Errorf("%q accepted", s)
		} else if !strings.Contains(err.Error(), "sso_issuer") {
			t.Errorf("%q: message does not name the key: %v", s, err)
		}
	}
}

func TestValidateIdentifier(t *testing.T) {
	for _, s := range []string{"idp_1", "cli1", "0oa1b2c3d4e5", "abc.def-ghi_jk@x", "a:b"} {
		if err := ValidateIdentifier("cliClientId", s); err != nil {
			t.Errorf("%q: %v", s, err)
		}
	}
	for _, s := range []string{"", "a b", "a\nb", "a=b", "[x]", "a/b", strings.Repeat("a", 300), "-lead"} {
		if err := ValidateIdentifier("cliClientId", s); err == nil {
			t.Errorf("%q accepted", s)
		}
	}
}

func TestWriteProfileRefusesUnsafeIssuer(t *testing.T) {
	home := t.TempDir()
	s := Source{Getenv: func(string) string { return "" }, Home: home}
	if _, err := WriteProfile(s, "p", map[string]string{"sso_issuer": "http://idp.example.com"}, []string{"sso_issuer"}); err == nil || !strings.Contains(err.Error(), "sso_issuer") {
		t.Fatalf("%v", err)
	}
	if _, err := os.Stat(configPath(s)); err == nil {
		t.Error("config written")
	}
}

// FR-005: an endpoint with credentials, a query or a fragment is refused, and the message never echoes them.
func TestValidateEndpointRefusesUserinfoQueryFragment(t *testing.T) {
	for _, raw := range []string{"https://user:hunter2@afs.example.com", "https://hunter2@afs.example.com", "http://user:hunter2@127.0.0.1:8080", "https://afs.example.com?token=hunter2", "https://afs.example.com#hunter2"} {
		err := ValidateEndpoint(raw)
		if err == nil {
			t.Errorf("%q accepted", raw)
			continue
		}
		if strings.Contains(err.Error(), "hunter2") {
			t.Errorf("%q: message echoes the secret: %v", raw, err)
		}
		if !strings.Contains(err.Error(), "endpoint_url must not contain credentials, a query or a fragment") {
			t.Errorf("%q: %v", raw, err)
		}
	}
	// other rejections do not echo userinfo either
	if err := ValidateEndpoint("http://user:hunter2@afs.example.com"); err == nil || strings.Contains(err.Error(), "hunter2") {
		t.Errorf("%v", err)
	}
}

func dirEntries(t *testing.T, dir string) []string {
	t.Helper()
	es, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	for _, e := range es {
		names = append(names, e.Name())
	}
	return names
}

// FR-019: a failed save leaves the original untouched and no temporary file behind.
func TestSaveFailuresLeaveOriginalAndNoTempFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config")
	if err := os.WriteFile(path, []byte("[default]\noutput = json\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	// the directory cannot be created: its parent is a file
	f := &File{Path: filepath.Join(path, "sub", "config")}
	_ = f.Set("default", "output", "table")
	if err := f.Save(); err == nil {
		t.Fatal("saved under a file")
	}
	if b, _ := os.ReadFile(path); string(b) != "[default]\noutput = json\n" { //nolint:gosec // test file
		t.Errorf("original changed: %q", b)
	}
	// the final rename fails: the target is a non-empty directory
	target := filepath.Join(dir, "target")
	if err := os.MkdirAll(filepath.Join(target, "keep"), 0o700); err != nil {
		t.Fatal(err)
	}
	g := &File{Path: target}
	_ = g.Set("default", "output", "table")
	if err := g.Save(); err == nil {
		t.Fatal("renamed over a directory")
	}
	for _, n := range dirEntries(t, dir) {
		if strings.HasPrefix(n, ".afs-") {
			t.Errorf("temporary file left behind: %s", n)
		}
	}
	if _, err := os.Stat(filepath.Join(target, "keep")); err != nil {
		t.Errorf("target damaged: %v", err)
	}
}
