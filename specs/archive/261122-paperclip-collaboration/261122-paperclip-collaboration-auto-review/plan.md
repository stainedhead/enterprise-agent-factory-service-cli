# Plan
Sequence: FR-R1 (edge), FR-R4 (edge key), FR-R3 (cli parse), FR-R2 (runtime). Each is red-green-refactor, then `make gate`. Layers: adapter-level only (`internal/edge`, `internal/cli`, `internal/cli/runtime`); no domain change.
