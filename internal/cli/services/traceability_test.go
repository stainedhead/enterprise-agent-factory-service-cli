package services_test

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// TestEveryStoryCriterionHasATest (M6-X1): every criterion of every CLI-US-n in the spec's
// traceability table is named by a "CLI-US-n.k" comment on at least one test, and so is every
// "5.1" row of the table by a "5.1-<name>" comment.
func TestEveryStoryCriterionHasATest(t *testing.T) {
	tasks, err := os.ReadFile(filepath.Join("..", "..", "..", "specs", "261121-cli-feature-complete", "tasks.md"))
	if err != nil {
		t.Skip("the spec directory is not present (archived or not checked out)")
	}
	want := map[string]bool{}
	story := regexp.MustCompile(`^\| (CLI-US-\d+) \| (.*) \|$`)
	crit := regexp.MustCompile(`\.(\d+)`)
	sec := regexp.MustCompile(`^\| 5\.1 ([a-z ]+) \|`)
	for _, line := range strings.Split(string(tasks), "\n") {
		if m := story.FindStringSubmatch(line); m != nil {
			for _, c := range crit.FindAllStringSubmatch(m[2], -1) {
				want[m[1]+"."+c[1]] = false
			}
		}
		if m := sec.FindStringSubmatch(line); m != nil {
			want["5.1-"+strings.ReplaceAll(m[1], " ", "-")] = false
		}
	}
	if len(want) < 60 {
		t.Fatalf("the traceability table yielded only %d criteria", len(want))
	}
	ref := regexp.MustCompile(`CLI-US-\d+\.\d+|5\.1-[a-z-]+`)
	root := filepath.Join("..", "..", "..")
	_ = filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(p, "_test.go") {
			return nil
		}
		b, _ := os.ReadFile(p) //nolint:gosec // test sources of this repository
		for _, id := range ref.FindAllString(string(b), -1) {
			if _, ok := want[id]; ok {
				want[id] = true
			}
		}
		return nil
	})
	for id, seen := range want {
		if !seen {
			t.Errorf("no test names criterion %s", id)
		}
	}
}
