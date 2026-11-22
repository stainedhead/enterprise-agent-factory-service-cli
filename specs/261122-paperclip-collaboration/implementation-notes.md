# Implementation notes: paperclip-collaboration

**Date:** 2026-11-22. Updated as work proceeds; records decisions, edge cases and deviations.

## Technical decisions

Inherited decisions PC-CLI-D1..D12 are in the PRD section 7.1. Decisions taken while implementing are listed below.

(none yet)

## Edge cases and solutions

(none yet)

## Deviations from plan

(none yet)

## Root contract follow-ups

(none yet)

## Parked

| Item | Why parked | Needs |
|---|---|---|
| End-to-end proofs against a dev environment (create collaboration agent, connect two agents; pipeline creates a Paperclip issue; service admin sets base domain) | No real cluster, AWS account or deployed `afs-edge` | A dev environment with `afs-edge`, `collab.base_domain`, a team provider with a collaboration client |
| Release signing, Homebrew tap | Credentials and tap repo not in this repository (as in v1.0.0) | CI secrets, tap repository |

## Lessons learned

(at completion)
