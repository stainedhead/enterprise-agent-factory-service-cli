# Plan: orchestration-and-model-access

Milestones continue the numbering of the 1.0.1 work (M0 to M10). `CurrentMilestone` ends at 12.

| Milestone | Scope | Exit |
|---|---|---|
| M11 | Re-pin; core mechanisms (deprecation, hidden, desc, notes); `orchestration` service; `admin list-orchestration-endpoints`; hints for the 1.4.0 reasons; `invoke-api` on `GetTeamOrchestration`; deprecations of the collaboration commands | gate green, coverage >= 90%, goldens for the 4 operations |
| M12 | Model access: `put-service-model-access`, team/agent/pool flags and columns, `--model-id` deprecation, settings keys, swarm skeleton | gate green, goldens for `PutServiceModelAccess` |
| Docs | README (with a diagram), user-docs, CHANGELOG 1.1.0, AGENTS.md, docs/, the root skills entry | `make docs` clean, links checked |
| Review | Code and design review against the AGENTS.md checklist; findings fixed | no P0/P1 left |
| Release | PR into `main`, CI green, merge (no tag: J tags) | merged |

Order inside a milestone is Red -> Green -> Refactor per task; `status.md` is updated after each task.
