# Implementation notes: cli-feature-complete-auto-review

Decisions are recorded as D-R<n>; J authorized the whole flow, so open questions were decided here.

- D-R1 (OQ-1): control characters render as `\xNN` (`\u0085` style for C1), keeping values diagnosable.
- D-R2 (OQ-2): `Retry-After` cap is 20 s, or the read timeout when smaller.
- D-R4 (OQ-4): the safe-update notice follows the source PRD (after the write).
