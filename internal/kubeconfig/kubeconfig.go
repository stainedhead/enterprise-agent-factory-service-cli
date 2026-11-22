// Package kubeconfig merges a kubeconfig returned by the API into the user's kubeconfig file
// (CLI-16): only entries with the same names are replaced, the write is atomic, a new file is
// 0600 and an existing file keeps its mode.
package kubeconfig

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// sections are the named lists of a kubeconfig that are merged by entry name.
var sections = []string{"clusters", "contexts", "users"}

// Options tune a merge.
type Options struct {
	// AWSProfile, when set, adds AWS_PROFILE=<AWSProfile> to the env block of every exec credential plugin.
	AWSProfile string
}

// Result says what a write did.
type Result struct {
	Path    string
	Created bool        // the file did not exist
	Mode    fs.FileMode // the mode of the written file
	Wide    bool        // an existing file is readable or writable by others (wider than 0600); its mode was kept
	Context string      // the current-context of the merged-in configuration
}

func mapping(n *yaml.Node) *yaml.Node {
	if n != nil && n.Kind == yaml.DocumentNode && len(n.Content) == 1 {
		n = n.Content[0]
	}
	if n == nil || n.Kind != yaml.MappingNode {
		return nil
	}
	return n
}

// get returns the value node of a key of a mapping node.
func get(m *yaml.Node, key string) *yaml.Node {
	for i := 0; i+1 < len(m.Content); i += 2 {
		if m.Content[i].Value == key {
			return m.Content[i+1]
		}
	}
	return nil
}

func set(m *yaml.Node, key string, v *yaml.Node) {
	for i := 0; i+1 < len(m.Content); i += 2 {
		if m.Content[i].Value == key {
			m.Content[i+1] = v
			return
		}
	}
	m.Content = append(m.Content, &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: key}, v)
}

func scalar(s string) *yaml.Node { return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: s} }

func entryName(item *yaml.Node) string {
	if item.Kind != yaml.MappingNode {
		return ""
	}
	if n := get(item, "name"); n != nil {
		return n.Value
	}
	return ""
}

func parse(b []byte) (*yaml.Node, error) {
	var doc yaml.Node
	if err := yaml.Unmarshal(b, &doc); err != nil {
		return nil, err
	}
	if doc.Kind == 0 { // an empty document
		doc = yaml.Node{Kind: yaml.DocumentNode, Content: []*yaml.Node{{Kind: yaml.MappingNode, Tag: "!!map"}}}
	}
	if mapping(&doc) == nil {
		return nil, errors.New("a kubeconfig must be a YAML mapping")
	}
	return &doc, nil
}

// addAWSProfile sets AWS_PROFILE in the env block of every user's exec plugin.
func addAWSProfile(in *yaml.Node, profile string) {
	users := get(mapping(in), "users")
	if users == nil {
		return
	}
	for _, u := range users.Content {
		user := get(u, "user")
		if user == nil || user.Kind != yaml.MappingNode {
			continue
		}
		ex := get(user, "exec")
		if ex == nil || ex.Kind != yaml.MappingNode {
			continue
		}
		env := get(ex, "env")
		if env == nil || env.Kind != yaml.SequenceNode {
			env = &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}
			set(ex, "env", env)
		}
		entry := &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map", Content: []*yaml.Node{scalar("name"), scalar("AWS_PROFILE"), scalar("value"), scalar(profile)}}
		replaced := false
		for i, e := range env.Content {
			if entryName(e) == "AWS_PROFILE" {
				env.Content[i], replaced = entry, true
			}
		}
		if !replaced {
			env.Content = append(env.Content, entry)
		}
	}
}

// Merge folds the incoming kubeconfig into the existing one (nil or empty for none): clusters,
// contexts and users are replaced by name or appended, everything else in the existing file is
// kept as it is, and current-context becomes the incoming one. It returns the merged YAML and
// the incoming current-context.
func Merge(existing, incoming []byte, o Options) ([]byte, string, error) {
	in, err := parse(incoming)
	if err != nil {
		return nil, "", fmt.Errorf("the service returned an invalid kubeconfig: %w", err)
	}
	if o.AWSProfile != "" {
		addAWSProfile(in, o.AWSProfile)
	}
	base, err := parse(existing)
	if err != nil {
		return nil, "", fmt.Errorf("the existing kubeconfig is not valid YAML: %w", err)
	}
	bm, im := mapping(base), mapping(in)
	if get(bm, "apiVersion") == nil {
		set(bm, "apiVersion", scalar("v1"))
	}
	if get(bm, "kind") == nil {
		set(bm, "kind", scalar("Config"))
	}
	for _, sec := range sections {
		add := get(im, sec)
		if add == nil || add.Kind != yaml.SequenceNode {
			continue
		}
		list := get(bm, sec)
		if list == nil || list.Kind != yaml.SequenceNode {
			list = &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}
			set(bm, sec, list)
		}
		for _, item := range add.Content {
			name, replaced := entryName(item), false
			for i, cur := range list.Content {
				if name != "" && entryName(cur) == name {
					list.Content[i], replaced = item, true
					break
				}
			}
			if !replaced {
				list.Content = append(list.Content, item)
			}
		}
	}
	cur := ""
	if c := get(im, "current-context"); c != nil {
		cur = c.Value
		set(bm, "current-context", scalar(cur))
	}
	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	if err := enc.Encode(base); err != nil {
		return nil, "", err
	}
	if err := enc.Close(); err != nil {
		return nil, "", err
	}
	return buf.Bytes(), cur, nil
}

// rename is replaced by tests to fail the final step of an atomic write.
var rename = os.Rename

// WriteFile writes data to path atomically (a temporary file in the same directory, then a
// rename). A new file is 0600 and a new directory 0700; an existing file keeps its mode.
func WriteFile(path string, data []byte) (Result, error) {
	res := Result{Path: path, Mode: 0o600, Created: true}
	if target, err := filepath.EvalSymlinks(path); err == nil {
		path = target
	}
	if st, err := os.Stat(path); err == nil {
		res.Created, res.Mode = false, st.Mode().Perm()
		res.Wide = res.Mode&0o077 != 0
	} else if !errors.Is(err, fs.ErrNotExist) {
		return res, err
	}
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return res, err
	}
	tmp, err := os.CreateTemp(dir, ".kubeconfig-*")
	if err != nil {
		return res, err
	}
	name := tmp.Name()
	cleanup := func() { _ = os.Remove(name) }
	if err := tmp.Chmod(res.Mode); err != nil {
		_ = tmp.Close()
		cleanup()
		return res, err
	}
	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		cleanup()
		return res, err
	}
	if err := tmp.Sync(); err != nil {
		_ = tmp.Close()
		cleanup()
		return res, err
	}
	if err := tmp.Close(); err != nil {
		cleanup()
		return res, err
	}
	if err := rename(name, path); err != nil {
		cleanup()
		return res, err
	}
	return res, nil
}

// Apply merges incoming into the file at path (created when missing) and writes it.
func Apply(path string, incoming []byte, o Options) (Result, error) {
	existing, err := os.ReadFile(path) //nolint:gosec // the user's kubeconfig path
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return Result{Path: path}, err
	}
	merged, cur, err := Merge(existing, incoming, o)
	if err != nil {
		return Result{Path: path}, err
	}
	res, err := WriteFile(path, merged)
	res.Context = cur
	return res, err
}
