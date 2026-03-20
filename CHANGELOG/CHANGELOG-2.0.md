## v2.0.0

### Added

- New openstacksdk-style proxy client: `client.Slurm.Jobs().List(ctx)`.
- Full v0.0.40 resource coverage: 9 Slurm + 11 SlurmDB resources.
- Automatic API version detection via slurmrestd `/openapi` endpoint.
- Retry with exponential backoff on transient failures (503, 429).
- Structured logging support via `*slog.Logger` with token masking.
- Unix socket support for local slurmrestd connections.
- Client-side validation for job submission (empty script rejection).
- `slurmtest` package: in-memory mock slurmrestd for testing.
- `hack/generate-version.sh`: scaffold OpenAPI client for new Slurm versions.
- `hack/check-compat.sh`: verify SDK compatibility with new OpenAPI specs.

### Changed

- SDK entry point changed from `pkg/client.NewClient()` to `slurmclient.NewClient()`.
- Response types use v0.0.40 definitions for all API versions (JSON-compatible v0.0.39-v0.0.44).
- Single import: `import slurm "github.com/supergate-hub/slurm-client"`.

### Removed

- Removed `pkg/client` (k8s-style Reader/Writer/Informer).
- Removed `pkg/slurmclient` (legacy HTTP transport).
- Removed `pkg/slurm/v0040` and `pkg/slurmdb/v0040` (old domain packages).
- Removed `pkg/types` (manual per-version type wrappers).
- Removed `pkg/object`, `pkg/event`, `pkg/utils`.

### Breaking Changes

- All `pkg/` imports are removed. Migrate to `slurmclient.NewClient()`.
- `UnifiedClient` replaced by proxy pattern (`client.Slurm.Jobs()`).
- k8s-style `Reader`/`Writer`/`Informer` interfaces removed.
