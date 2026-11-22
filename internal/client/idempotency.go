package client

import (
	"math/rand/v2"

	"github.com/oklog/ulid/v2"
)

// NewIdempotencyKey returns a ULID, one per invocation (CLI-11).
func NewIdempotencyKey() string { return ulid.Make().String() }

func defaultRand() float64 { return rand.Float64() } //nolint:gosec // jitter, not a secret
