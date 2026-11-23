# Architecture: orchestration-and-model-access

The CLI keeps its layering (`cmd/afs` -> `internal/cli` -> `internal/client`, `internal/config`, `internal/output`);
command handlers parse flags, make one API call and render. This change adds no layer.

## New and changed packages

| Package | Change |
|---|---|
| `internal/cli/orchestration` | New service: overrides of `GetTeamOrchestration`, `PutTeamOrchestration`, `ChangeCollaborationTool`; exposure options and their Build hook; notices (`AfterDescribe`); `Columns` |
| `internal/cli/modelaccess` | New: the model options shared by agents and pools (`--models`, `--default-model-id`, deprecated `--model-id`) as one override |
| `internal/cli` (`override.go`, `options.go`) | Override gains `Deprecated`, `DeprecatedOptions`, `Hidden`, `Desc`, `Enum`, `Notes`; Option gains `Hidden`, `Deprecated`; a `Kinds` override to a list defaults its elements to strings |
| `internal/cli/help` | Prints deprecation notes and notes, skips hidden options, service summary for `orchestration` |
| `internal/cli/exec` | `deprecation.go` (stderr lines for deprecated operations, options and the `Deprecation` header); hints table for the 1.4.0 reasons |
| `internal/output` | `Column.NAUnless` (N/A for a surface the tool lacks) |
| `internal/cli/wait` | `GetTeamOrchestration` poll target for put and change-tool |
| `internal/cli/admin`, `teams`, `agents`, `pools`, `models`, `collaboration`, `swarmtemplates` | overrides, columns, examples, hooks as in the FR table |

## Flow of a deprecation

```
parse -> Build (options, hooks) -> warnDeprecated (op and options, stderr) -> send -> warnServiceDeprecation
                                                                                  (only if no static line)
```
Nothing on stdout changes and the exit code is never affected.

## Flow of `invoke-api`

`GetTeamOrchestration` (one call) -> `webApi.supported` false: exit 254 "Web API is off (N/A for <Tool>)";
`webApi.url` null: exit 254 "Web API is off for this team (<reasons>)" + hint; otherwise origin check (unchanged) and the edge call.
