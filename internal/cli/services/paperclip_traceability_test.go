package services_test

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"testing"
)

// prdCandidates are where the Paperclip change PRD lives while the spec is active and after it is archived
// (the fixes spec is archived inside it).
var prdCandidates = []string{
	"specs/261122-paperclip-collaboration/paperclip-collaboration-PRD.md",
	"specs/archive/261122-paperclip-collaboration/paperclip-collaboration-PRD.md",
}

// supersededIn110 lists the tests that proved a 1.0.1 criterion which release 1.1.0 changed on purpose (the
// orchestration agent replaces the collaboration agent, API contract section 4.4d). The criterion stays in the
// archived PRD as history; the replacement test is named here and proved by the 1.1.0 traceability test.
var supersededIn110 = map[string]string{
	"collaboration.TestInvokeAPINullRestApiURL255":             "collaboration.TestInvokeAPIWebAPIOff254",
	"agents.TestCreateAgentRoleEnum252":                        "agents.TestCreateAgentRoleSentUnchanged",
	"agents.TestCreateAgentCollaborationRoleWithConnect252":    "agents.TestCreateAgentRoleHiddenAndRefusedByAPI",
	"agents.TestAgentTableRoleAndPaperclipColumnsGolden":       "agents.TestAgentTableRoleAndLinkColumnsGolden",
	"agents.TestListAgentsRoleEnum252":                         "agents.TestListAgentsRoleOpenEnumPassesAnyValue",
	"swarmtemplates.TestSwarmTemplateSkeletonHasCollaboration": "swarmtemplates.TestSwarmTemplateSkeletonHasOrchestration",
}

var (
	storyRe     = regexp.MustCompile(`(?m)^\*\*CLI-US-(\d+) `)
	criterionRe = regexp.MustCompile(`(?m)^- (\d+\.\d+) `)
	provenRe    = regexp.MustCompile("`([a-z]+)\\.(Test[A-Za-z0-9_]+)`")
)

// testFuncs maps a package directory name to the names of the Test functions in its _test.go files.
func testFuncs(t *testing.T, root string) map[string]map[string]bool {
	t.Helper()
	out := map[string]map[string]bool{}
	fn := regexp.MustCompile(`(?m)^func (Test[A-Za-z0-9_]+)\(`)
	err := filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && (d.Name() == "specs" || d.Name() == ".git" || d.Name() == "bin") {
			return fs.SkipDir
		}
		if d.IsDir() || !strings.HasSuffix(p, "_test.go") {
			return nil
		}
		b, err := os.ReadFile(p) //nolint:gosec // test reads the repository's own files
		if err != nil {
			return err
		}
		dir := filepath.Base(filepath.Dir(p))
		if dir == "docgen" {
			dir = "docs"
		}
		if out[dir] == nil {
			out[dir] = map[string]bool{}
		}
		for _, m := range fn.FindAllStringSubmatch(string(b), -1) {
			out[dir][m[1]] = true
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// CLI-US-24 to CLI-US-35 (Paperclip): every acceptance criterion names the test that proves it, and that
// test exists. A criterion without one, or a name that no test file defines, fails here.
func TestPaperclipCriteriaTraceToTests(t *testing.T) {
	root := filepath.Join("..", "..", "..")
	var prd []byte
	for _, c := range prdCandidates {
		if b, err := os.ReadFile(filepath.Join(root, c)); err == nil { //nolint:gosec // fixed repository paths
			prd = b
			break
		}
	}
	if prd == nil {
		t.Fatalf("the Paperclip change PRD is in none of %v", prdCandidates)
	}
	funcs := testFuncs(t, root)
	text := string(prd)
	stories := storyRe.FindAllStringSubmatchIndex(text, -1)
	seen := map[string]bool{}
	var missing []string
	for i, s := range stories {
		end := len(text)
		if i+1 < len(stories) {
			end = stories[i+1][0]
		}
		n := text[s[2]:s[3]]
		if n < "24" || n > "35" || len(n) != 2 {
			continue
		}
		seen[n] = true
		body := text[s[0]:end]
		if j := strings.Index(body, "\n## "); j >= 0 {
			body = body[:j]
		}
		crit := criterionRe.FindAllStringSubmatchIndex(body, -1)
		if len(crit) == 0 {
			t.Errorf("CLI-US-%s has no numbered criteria", n)
		}
		for k, c := range crit {
			cend := len(body)
			if k+1 < len(crit) {
				cend = crit[k+1][0]
			}
			id := body[c[2]:c[3]]
			tail := body[c[0]:cend]
			p := strings.Index(tail, "*Proven by*")
			if p < 0 {
				t.Errorf("criterion %s has no *Proven by*", id)
				continue
			}
			names := provenRe.FindAllStringSubmatch(tail[p:], -1)
			if len(names) == 0 {
				t.Errorf("criterion %s names no test", id)
			}
			for _, m := range names {
				if repl, ok := supersededIn110[m[1]+"."+m[2]]; ok {
					if !funcs[strings.SplitN(repl, ".", 2)[0]][strings.SplitN(repl, ".", 2)[1]] {
						missing = append(missing, "CLI-US-"+id+": replacement "+repl)
					}
					continue
				}
				if !funcs[m[1]][m[2]] {
					missing = append(missing, "CLI-US-"+id+": "+m[1]+"."+m[2])
				}
			}
		}
	}
	for n := 24; n <= 35; n++ {
		if !seen[itoa(n)] {
			t.Errorf("CLI-US-%d not found in the PRD", n)
		}
	}
	sort.Strings(missing)
	for _, m := range missing {
		t.Errorf("no such test: %s", m)
	}
}

func itoa(n int) string { return strconv.Itoa(n) }
