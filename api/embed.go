// Package api embeds the pinned snapshot of the AFS REST API's OpenAPI document.
package api

import _ "embed"

// OpenAPI is the verbatim copy of -restapi's api/openapi.yaml named by SOURCE.
//
//go:embed openapi.yaml
var OpenAPI []byte

// Source is the provenance record (repo, ref, commit, sha256, synced_at).
//
//go:embed SOURCE
var Source []byte
