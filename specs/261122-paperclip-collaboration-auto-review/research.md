# Research
`client.Do` stops retrying when `Retry-After` is above `min(20 s, read timeout)`; `internal/edge` used `client.RetryAfter`, which clamps. `os.ReadFile` of a directory returns a read error that names the path.
