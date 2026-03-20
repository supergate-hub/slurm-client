# Slurm Client

## Overview

This Go SDK provides a proxy-style client for the [Slurm REST API][rest-api],
inspired by [openstacksdk]. It communicates with [slurmrestd] using JWT-based
authentication.

## Architecture

The SDK is organized in three layers:

```
┌─────────────────────────────────────────────────┐
│  Client                                          │
│  ├── Slurm   (SlurmProxy)                        │
│  │   ├── Jobs()       → JobsService              │
│  │   ├── Nodes()      → NodesService             │
│  │   ├── Partitions() → PartitionsService         │
│  │   └── ... (9 resources)                        │
│  └── Slurmdb (SlurmdbProxy)                       │
│      ├── Accounts()   → AccountsService           │
│      ├── Users()      → UsersService              │
│      └── ... (11 resources)                       │
├─────────────────────────────────────────────────┤
│  Transport                                        │
│  ├── HTTP request/response                        │
│  ├── Bearer token auth (masked in logs)           │
│  ├── Retry with exponential backoff               │
│  ├── Structured logging (slog)                    │
│  ├── URL construction: base/domain/version/path   │
│  └── Unix socket support                          │
├─────────────────────────────────────────────────┤
│  api/v0040 (OpenAPI generated types)              │
│  └── V0040JobInfo, V0040Node, V0040Account, ...   │
└─────────────────────────────────────────────────┘
```

### Request Flow

```
User code
  │
  ▼
client.Slurm.Jobs().List(ctx)
  │
  ▼
JobsService.List()
  │  builds resource path: "/jobs"
  ▼
Transport.Get(ctx, "/jobs", &resp)
  │  constructs URL: baseURL + "/slurm/v0.0.40" + "/jobs"
  │  injects: Authorization: Bearer <token>
  │  retry on: 503, 429, connection reset
  ▼
slurmrestd HTTP API
  │
  ▼
JSON response → unmarshal into v0040.V0040OpenapiJobInfoResp
  │
  ▼
return []v0040.V0040JobInfo
```

### Version Compatibility

Response parsing uses v0.0.40 type definitions. All supported versions
(v0.0.39 through v0.0.44) share the same JSON field names, so v0.0.40
types can safely unmarshal responses from any version. The Transport
layer handles the version-specific URL prefix.

### Testing

The `slurmtest` package provides an in-memory mock slurmrestd server
for unit testing without a real Slurm cluster.

<!-- Links -->

[rest-api]: https://slurm.schedmd.com/rest_api.html
[openstacksdk]: https://docs.openstack.org/openstacksdk/latest/
[slurmrestd]: https://slurm.schedmd.com/slurmrestd.html
