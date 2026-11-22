# Research: cli-feature-complete-auto-review

Created 2026-11-22. Source PRD: cli-feature-complete-auto-review-PRD.md.

## Research questions
1. OQ-1 FR-003 escape style: `\x1b` form (readable). Decision D-R1.
2. OQ-2 FR-011 `Retry-After` cap: 20 s or the read timeout when smaller. Decision D-R2.
3. OQ-3 FR-014 kubeconfig prefix: read from the pinned snapshot's `GetTeamKubeconfig` example. Decision D-R3.
4. OQ-4 FR-016 align code to the PRD (notice after the write). Decision D-R4.
5. OQ-5 FR-017 use or remove generated types. Decision D-R5.
6. OQ-6 FR-015 does exiting 252 break completion? Verify with completion tests.

## Industry standards
OIDC Discovery 1.0 section 4.3 (issuer equality); RFC 8628 device flow; RFC 9110 `Retry-After`; AWS CLI exit codes.

## Existing implementations
AWS CLI v2 for usage errors (252) and `NextToken` handling.

## API documentation, best practices, references
`api/openapi.yaml` (pinned snapshot, read only). Findings and evidence: the source PRD.

## Open questions
Resolved in implementation-notes.md as decisions D-R1..D-R6.
