# Plan: model-endpoint-gateway

Milestones continue the numbering (M0 to M12 done). `CurrentMilestone` ends at 14.

| Milestone | Scope | Exit |
|---|---|---|
| M13 | Re-pin; core changes (shared query/body names, raw body, safe-update readers); `model-endpoints` service for the service-admin operations (catalog, models, prices, proposals, teams, credentials, federation, statements, providers) | gate green, goldens for the 18 operations |
| M14 | Team side: team endpoints, selections, gateway settings, inference reports, route options on agents and pools, report `--route`, hints | gate green, goldens for the 9 operations |
| Docs | `make docs`, user-docs examples, CHANGELOG 1.2.0, README, AGENTS.md, CLI PRD section 4.2 | `docs.TestReferenceUpToDate` green |
| Release | PR into `main`, CI green, merge (no tag: J tags) | merged |

Red -> Green -> Refactor per task; `status.md` is updated after each task.
