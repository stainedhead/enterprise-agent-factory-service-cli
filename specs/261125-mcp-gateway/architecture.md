# Architecture

No new layer. `internal/cli/mcp` registers overrides on the registry like every service package; the generic executor
builds, sends and renders. Two core changes: `options.go` marks `clientSecret` sensitive, and `safeupdate` learns a
wrapper member and a versionless list for the MCP writes.
