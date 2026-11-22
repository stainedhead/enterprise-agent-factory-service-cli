# Plan: cli-feature-complete-auto-review

Sequential streams (no sub-agents available), each fix Red -> Green -> Refactor, a short review, a status
update and a commit. Order:

1. FR-001 alone (config injection).
2. P1: FR-002, FR-003, FR-004.
3. Stream A rest: FR-006, FR-007, FR-008.
4. Stream C: FR-005, FR-009, FR-010, FR-011.
5. Stream D: FR-012, FR-013, FR-016.
6. Stream B rest: FR-018, FR-020.
7. Stream E: FR-014, FR-015, FR-017, FR-019.
8. Docs: errors-and-exit-codes, pagination, SSO, troubleshooting, ADR entries, README.

Gate after each group: `make gate`.
