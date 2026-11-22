package auth

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"
)

// DefaultWorkloadTokenPath is where a pod's projected ServiceAccount token for the AFS API is mounted.
const DefaultWorkloadTokenPath = "/var/run/secrets/afs/api/token" //nolint:gosec // a file path, not a credential

// WorkloadToken yields the projected ServiceAccount token of a worker pod (CLI-5, D-C29). The path
// is AFS_API_TOKEN_FILE, else DefaultWorkloadTokenPath when that file exists. The file is read on
// every call because the kubelet rotates the token.
type WorkloadToken struct {
	Getenv func(string) string
	// DefaultPath replaces DefaultWorkloadTokenPath (tests).
	DefaultPath string
	// ReadFile replaces os.ReadFile (tests).
	ReadFile func(string) ([]byte, error)
}

// Token reads the token file. No file at the default path means no workload credential; a file named
// by AFS_API_TOKEN_FILE that cannot be read is an error, because the user asked for it.
func (w WorkloadToken) Token(_ context.Context) (string, bool, error) {
	path, explicit := "", false
	if w.Getenv != nil {
		path = w.Getenv("AFS_API_TOKEN_FILE")
		explicit = path != ""
	}
	if path == "" {
		path = w.DefaultPath
		if path == "" {
			path = DefaultWorkloadTokenPath
		}
	}
	read := w.ReadFile
	if read == nil {
		read = os.ReadFile
	}
	b, err := read(path) //nolint:gosec // the configured token path
	switch {
	case errors.Is(err, fs.ErrNotExist) && !explicit:
		return "", false, nil
	case err != nil:
		return "", false, fmt.Errorf("cannot read the workload token file %s: %w", path, err)
	}
	tok := strings.TrimSpace(string(b))
	return tok, tok != "", nil
}
