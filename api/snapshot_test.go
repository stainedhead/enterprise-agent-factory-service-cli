package api

import (
	"crypto/sha256"
	"encoding/hex"
	"regexp"
	"strings"
	"testing"
)

func sourceField(t *testing.T, raw []byte, key string) string {
	t.Helper()
	for _, l := range strings.Split(string(raw), "\n") {
		if k, v, ok := strings.Cut(l, ": "); ok && k == key {
			return strings.TrimSpace(v)
		}
	}
	t.Fatalf("SOURCE has no %q", key)
	return ""
}

func checksum(b []byte) string {
	h := sha256.Sum256(b)
	return hex.EncodeToString(h[:])
}

func TestSnapshotChecksumMatchesSource(t *testing.T) {
	if got, want := checksum(OpenAPI), sourceField(t, Source, "sha256"); got != want {
		t.Fatalf("openapi.yaml sha256 %s != SOURCE %s", got, want)
	}
}

func TestTamperedCopyFailsCheck(t *testing.T) {
	tampered := append([]byte("# edited\n"), OpenAPI...)
	if checksum(tampered) == sourceField(t, Source, "sha256") {
		t.Fatal("a tampered snapshot must not match SOURCE")
	}
}

func TestSnapshotIsOpenAPI31(t *testing.T) {
	s := string(OpenAPI)
	if !regexp.MustCompile(`(?m)^openapi: 3\.1\.\d+$`).MatchString(s) {
		t.Fatal("snapshot is not OpenAPI 3.1")
	}
	if !regexp.MustCompile(`(?m)^  version: \d+\.\d+\.\d+$`).MatchString(s) {
		t.Fatal("info.version is not semver")
	}
}

func TestSourceHasAllFields(t *testing.T) {
	for _, k := range []string{"repo", "ref", "commit", "sha256", "synced_at"} {
		if sourceField(t, Source, k) == "" {
			t.Errorf("empty %s", k)
		}
	}
}
