// Package release_test checks the GoReleaser configuration locally, without GoReleaser: the
// platform matrix, the version stamping and the checksum file are read from .goreleaser.yaml,
// and a binary is built with the same link flags (CLI-15, D-C14).
package release_test

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

type config struct {
	Builds []struct {
		Main    string   `yaml:"main"`
		Goos    []string `yaml:"goos"`
		Goarch  []string `yaml:"goarch"`
		Ldflags []string `yaml:"ldflags"`
		Ignore  []struct {
			Goos   string `yaml:"goos"`
			Goarch string `yaml:"goarch"`
		} `yaml:"ignore"`
	} `yaml:"builds"`
	Checksum struct {
		Name string `yaml:"name_template"`
	} `yaml:"checksum"`
	Signs []struct {
		Cmd       string   `yaml:"cmd"`
		Artifacts string   `yaml:"artifacts"`
		Args      []string `yaml:"args"`
	} `yaml:"signs"`
}

func load(t *testing.T) config {
	t.Helper()
	b, err := os.ReadFile("../../.goreleaser.yaml")
	if err != nil {
		t.Fatal(err)
	}
	var c config
	if err := yaml.Unmarshal(b, &c); err != nil {
		t.Fatal(err)
	}
	return c
}

// CLI-15: archives for darwin/arm64, darwin/amd64, linux/arm64, linux/amd64 and windows/amd64.
// 5.1-distribution: the release produces the five platform archives (see also the signing tests).
func TestReleaseTargets(t *testing.T) {
	c := load(t)
	if len(c.Builds) != 1 || c.Builds[0].Main != "./cmd/afs" {
		t.Fatalf("builds: %+v", c.Builds)
	}
	b := c.Builds[0]
	var got []string
	for _, os := range b.Goos {
		for _, arch := range b.Goarch {
			skip := false
			for _, ig := range b.Ignore {
				skip = skip || (ig.Goos == os && ig.Goarch == arch)
			}
			if !skip {
				got = append(got, os+"/"+arch)
			}
		}
	}
	sort.Strings(got)
	want := "darwin/amd64 darwin/arm64 linux/amd64 linux/arm64 windows/amd64"
	if strings.Join(got, " ") != want {
		t.Errorf("targets %v, want %s", got, want)
	}
	if c.Checksum.Name != "checksums.txt" {
		t.Errorf("checksum file %q", c.Checksum.Name)
	}
}

// D-C14: checksums.txt is signed with cosign; the signing step can be skipped locally.
func TestChecksumsSignedWithCosign(t *testing.T) {
	c := load(t)
	if len(c.Signs) != 1 || c.Signs[0].Cmd != "cosign" || c.Signs[0].Artifacts != "checksum" || !strings.Contains(strings.Join(c.Signs[0].Args, " "), "sign-blob") {
		t.Fatalf("signs: %+v", c.Signs)
	}
}

// CLI-15, M4-G1: the binary built with the configuration's link flags reports the tag.
func TestSnapshotBinaryReportsVersion(t *testing.T) {
	c := load(t)
	const snapshot = "0.1.0-SNAPSHOT-abc1234"
	var flags []string
	for _, f := range c.Builds[0].Ldflags {
		f = strings.NewReplacer("{{ .Version }}", snapshot, "{{ .ShortCommit }}", "abc1234", "{{ .Date }}", "2026-11-22T00:00:00Z").Replace(f)
		flags = append(flags, f)
	}
	bin := filepath.Join(t.TempDir(), "afs")
	build := exec.CommandContext(context.Background(), "go", "build", "-trimpath", "-ldflags", strings.Join(flags, " "), "-o", bin, "../../cmd/afs") //nolint:gosec // fixed arguments
	build.Env = append(os.Environ(), "CGO_ENABLED=0")
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build: %v\n%s", err, out)
	}
	out, err := exec.CommandContext(context.Background(), bin, "--version").Output() //nolint:gosec // the binary just built
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(string(out), "afs-cli/"+snapshot+" contract/1.6 go") {
		t.Errorf("--version: %q", out)
	}
}
