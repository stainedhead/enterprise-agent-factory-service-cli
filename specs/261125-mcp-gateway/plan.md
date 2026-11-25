# Plan

1. Re-pin and regenerate (T1); fix the two code generator collisions with the existing overlay, not by editing the contract.
2. Core changes only where the generic path needs them (T2): a second sensitive member, readers.
3. One package `internal/cli/mcp` holds every MCP override, whichever service the contract puts the command in (T3).
4. Milestones, hints, tests, docs (T4 to T6).
