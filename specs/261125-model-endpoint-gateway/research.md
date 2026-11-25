# Research: model-endpoint-gateway

## Sources

- Root `model-endpoint-gateway-PRD.md` Appendix D.3.4 (operations) and D.4 (MEG-CLI-1 to 4), Appendix C.7 (TC-CLI-01 to 05).
- Root `docs/contracts/api-contract.md` section 4.13 (v1.5.0).
- `-restapi` `api/openapi.yaml` at `7d48ea70d2f9cb77ef4308b0225ca9d4168dea60` (contract 1.5.0, 221 operations). At that
  commit the 27 new operations answer `501 not_implemented` (`h_model_endpoints_pending.go`); the CLI builds from the
  contract, and the handlers land in the `-restapi` 1.2.0 implementation change.

## Findings

1. **27 operations, not 26.** Section 4.13 has 26 rows, but the `GetInferenceReport` row covers two operations
   (`GetInferenceReport` and `GetTeamInferenceReport`). The snapshot has 27 new operationIds; the CLI has 27 new
   commands (219 commands over 221 operations; `GetHealth` and `GetReadiness` have none).
2. **Service names come from `x-afs-cli`.** 23 operations are in `model-endpoints`, two in `teams`
   (`put-team-gateway-settings`, `put-team-model-endpoints`) and two in `reports`. The CLI uses them unchanged.
3. **`ImportProviderStatement` has `period` and `source` both as query parameters (for a `text/csv` body) and as
   body members (for a JSON body).** One option each; the override puts them where the chosen body needs them.
4. **Credential versions are not readable.** `PutModelEndpointCredential` and `PutTeamModelEndpointCredential` take an
   optional `If-Match`, but `ModelEndpointCredential` has no `version` and no operation reads one. The CLI sends
   `If-Match` only when `--if-match` is given; to replace a service credential, `rotate-model-endpoint-credential`
   needs no version. Raised as a contract note.
5. **`UpdateInferenceProvider` needs `If-Match` but there is no single-provider GET.** The version is read from
   `ListInferenceProviders` (not paginated), matching the item's `provider`.
6. **`PutTeamModelEndpoints` would read the team's version** through the generic rule (same path prefix as
   `GetTeam`), but the selections response carries no version. The CLI treats it as versionless.
7. **`PutTeamGatewaySettings`, `PutEndpointFederation`** have no reader in the snapshot; If-Match is optional and sent
   only with `--if-match` (the response prints the version to use).
8. **Permission of `ImportProviderStatement`.** PRD D.3.4 lists `admin:model-endpoints`; the contract and the snapshot
   say `admin:pricing`. The CLI follows the snapshot (it only names the permission in tests and help).
9. Agents and pools gain `models.routeAccess`, `models.routes[]` and `models.defaultRoute` inside the existing `models`
   member (no new top-level field); usage and cost reports gain `route` and `groupBy=route`.
