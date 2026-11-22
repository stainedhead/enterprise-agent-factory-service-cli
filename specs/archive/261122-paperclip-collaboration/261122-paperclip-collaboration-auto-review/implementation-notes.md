# Implementation notes
- FR-R1: `client.RetryAfterExceedsCap`; `invoke-api` returns the 429/503 answer at once when Retry-After is above the cap. The earlier edge test that expected a clamped 20 s wait was replaced.
- FR-R2, FR-R3, FR-R4: tests added; the behaviour already held (the review found gaps in proof, not in code).
- No fix deferred.
