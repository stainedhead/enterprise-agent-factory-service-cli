// Package config reads and writes ~/.afs/config and ~/.afs/credentials (AWS-style INI) and
// resolves settings with the precedence option, environment, profile, default (CLI-3).
package config

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// File is a line-preserving INI file: updating one section leaves every other line untouched.
type File struct {
	Path  string
	lines []string
}

// Load reads an INI file; a missing file is an empty one.
func Load(path string) (*File, error) {
	f := &File{Path: path}
	b, err := os.ReadFile(path) //nolint:gosec // user-selected config path
	if errors.Is(err, fs.ErrNotExist) {
		return f, nil
	}
	if err != nil {
		return nil, err
	}
	s := strings.ReplaceAll(string(b), "\r\n", "\n")
	s = strings.TrimSuffix(s, "\n")
	if s != "" {
		f.lines = strings.Split(s, "\n")
	}
	return f, nil
}

func sectionName(line string) (string, bool) {
	t := strings.TrimSpace(line)
	if strings.HasPrefix(t, "[") && strings.HasSuffix(t, "]") {
		return strings.TrimSpace(t[1 : len(t)-1]), true
	}
	return "", false
}

func keyValue(line string) (string, string, bool) {
	t := strings.TrimSpace(line)
	if t == "" || strings.HasPrefix(t, "#") || strings.HasPrefix(t, ";") || strings.HasPrefix(t, "[") {
		return "", "", false
	}
	k, v, ok := strings.Cut(t, "=")
	if !ok {
		return "", "", false
	}
	return strings.TrimSpace(k), strings.TrimSpace(v), true
}

// span returns the line range [start, end) of a section's body, or ok=false.
func (f *File) span(section string) (int, int, bool) {
	start := -1
	for i, l := range f.lines {
		if name, ok := sectionName(l); ok {
			if start >= 0 {
				return start, i, true
			}
			if name == section {
				start = i + 1
			}
		}
	}
	if start >= 0 {
		return start, len(f.lines), true
	}
	return 0, 0, false
}

// HasSection reports whether a section exists.
func (f *File) HasSection(section string) bool { _, _, ok := f.span(section); return ok }

// Get returns a key of a section.
func (f *File) Get(section, key string) (string, bool) {
	s, e, ok := f.span(section)
	if !ok {
		return "", false
	}
	for _, l := range f.lines[s:e] {
		if k, v, ok := keyValue(l); ok && k == key {
			return v, true
		}
	}
	return "", false
}

// Set sets a key, creating the section when needed.
func (f *File) Set(section, key, value string) {
	s, e, ok := f.span(section)
	if !ok {
		if len(f.lines) > 0 && strings.TrimSpace(f.lines[len(f.lines)-1]) != "" {
			f.lines = append(f.lines, "")
		}
		f.lines = append(f.lines, "["+section+"]", key+" = "+value)
		return
	}
	for i := s; i < e; i++ {
		if k, _, ok := keyValue(f.lines[i]); ok && k == key {
			f.lines[i] = key + " = " + value
			return
		}
	}
	// insert after the last non-blank line of the section
	at := e
	for at > s && strings.TrimSpace(f.lines[at-1]) == "" {
		at--
	}
	f.lines = append(f.lines[:at], append([]string{key + " = " + value}, f.lines[at:]...)...)
}

// String renders the file.
func (f *File) String() string {
	if len(f.lines) == 0 {
		return ""
	}
	return strings.Join(f.lines, "\n") + "\n"
}

// Save writes the file atomically with mode 0600, creating its directory with 0700.
func (f *File) Save() error {
	dir := filepath.Dir(f.Path)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(dir, ".afs-*")
	if err != nil {
		return err
	}
	name := tmp.Name()
	if err := tmp.Chmod(0o600); err != nil {
		_ = tmp.Close()
		_ = os.Remove(name)
		return err
	}
	if _, err := tmp.WriteString(f.String()); err != nil {
		_ = tmp.Close()
		_ = os.Remove(name)
		return err
	}
	if err := tmp.Close(); err != nil {
		_ = os.Remove(name)
		return err
	}
	return os.Rename(name, f.Path)
}

// CheckPrivate refuses a file readable or writable by group or others (CLI-5): the error
// carries the chmod hint. Windows has no such modes; the check is skipped there.
func CheckPrivate(path string) error {
	if runtime.GOOS == "windows" {
		return nil
	}
	st, err := os.Stat(path)
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	if st.Mode().Perm()&^0o600 != 0 {
		return fmt.Errorf("%s has mode %04o, wider than 0600; run: chmod 600 %s", path, st.Mode().Perm(), path)
	}
	return nil
}
