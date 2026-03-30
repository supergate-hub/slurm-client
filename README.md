# Slurm Client

<div align="center">

[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg?style=for-the-badge)](./LICENSES/Apache-2.0.txt)
[![Tag](https://img.shields.io/github/v/tag/supergate-hub/slurm-client?style=for-the-badge)](https://github.com/supergate-hub/slurm-client/tags/)
[![Go-Version](https://img.shields.io/github/go-mod/go-version/supergate-hub/slurm-client?style=for-the-badge)](./go.mod)
[![E2E](https://img.shields.io/github/actions/workflow/status/supergate-hub/slurm-client/e2e.yaml?style=for-the-badge&label=E2E)](https://github.com/supergate-hub/slurm-client/actions/workflows/e2e.yaml)

**MCP server + Go SDK for [Slurm][slurm] HPC clusters.**

Let AI agents manage your Slurm cluster via natural language,
or use the Go SDK for programmatic control.

</div>

[slurm]: https://slurm.schedmd.com/
[rest-api]: https://slurm.schedmd.com/rest_api.html
[openstacksdk]: https://docs.openstack.org/openstacksdk/latest/
[mcp]: https://modelcontextprotocol.io/

## Architecture

```mermaid
graph TD
    Agent["🤖 AI Agent / Claude"]

    subgraph MCP["slurm-mcp server"]
        direction TB
        Tools["Tool Gating\n8 core · 31 with --all-tools"]
        RBAC["RBAC + Audit Logging"]
        Backend["MCPBackend Interface"]
        Tools --> RBAC --> Backend
    end

    Agent -- "MCP (stdio / SSE)" --> MCP

    subgraph REST["REST Backend"]
        SDK["Go SDK"]
        SDK -- "HTTP / Unix Socket" --> slurmrestd["slurmrestd"]
    end

    subgraph SSH["SSH Backend"]
        CLI["CLI Executor"]
        CLI -- "SSH" --> cmds["squeue · sbatch\nscontrol · sacctmgr"]
    end

    Backend --> SDK
    Backend --> CLI

    subgraph Clusters["Slurm Clusters"]
        direction LR
        C1["GPU Cluster"]
        C2["CPU Cluster"]
        C3["Training Cluster"]
    end

    slurmrestd --> Clusters
    cmds --> Clusters

    style MCP fill:#dbeafe,stroke:#3b82f6,stroke-width:2px
    style REST fill:#dcfce7,stroke:#22c55e,stroke-width:2px
    style SSH fill:#ffedd5,stroke:#f97316,stroke-width:2px
    style Clusters fill:#f3f4f6,stroke:#6b7280,stroke-width:1px
    style Agent fill:#ede9fe,stroke:#8b5cf6,stroke-width:2px
```

## MCP Server

The MCP server enables AI agents like Claude to manage Slurm clusters via natural language.
Single Go binary, no Python, no pip, no SSH configuration required for REST mode.

### Quick Start (2 minutes)

**1. Download**

```bash
go install github.com/supergate-hub/slurm-client/cmd/slurm-mcp@latest
# or download from GitHub Releases
```

**2. Configure Claude Desktop**

Add to `~/Library/Application Support/Claude/claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "slurm": {
      "command": "slurm-mcp",
      "env": {
        "SLURM_ENDPOINT": "http://slurmctld:6820",
        "SLURM_TOKEN": "your-jwt-token"
      }
    }
  }
}
```

**3. Restart Claude Desktop.** Ask Claude: "Show me all running jobs on the cluster."

### Backend Modes

```bash
# REST backend (default) - requires slurmrestd
SLURM_ENDPOINT=http://slurmctld:6820 SLURM_TOKEN=<jwt> slurm-mcp

# SSH backend - for clusters without slurmrestd
slurm-mcp --backend=ssh  # (coming soon)

# Multi-cluster
slurm-mcp --config clusters.yaml

# SSE transport (remote / Docker)
slurm-mcp --config clusters.yaml --transport=sse --port=8080 --bearer-token=secret

# All tools (default: 8 core tools)
slurm-mcp --all-tools
```

### Core Tools (default)

These 8 tools are exposed by default. Use `--all-tools` to enable all 31.

| Tool | Description |
|------|------------|
| `slurm_ping` | Cluster health check |
| `slurm_list_jobs` | List running/pending jobs |
| `slurm_get_job` | Job details |
| `slurm_submit_job` | Submit a job script |
| `slurm_cancel_job` | Cancel a job |
| `slurm_list_nodes` | Node status overview |
| `slurm_list_partitions` | Partition info |
| `slurm_queue_depth` | Queue depth analysis |

### All Tools (31)

<details>
<summary>Expand full tool list</summary>

**Slurm API** (all tools accept optional `cluster` parameter):

`slurm_ping`, `slurm_list_jobs`, `slurm_get_job`, `slurm_submit_job`,
`slurm_cancel_job`, `slurm_list_nodes`, `slurm_get_node`,
`slurm_list_partitions`, `slurm_get_partition`, `slurm_list_reservations`,
`slurm_list_licenses`, `slurm_diag`

**SlurmDB API:**

`slurmdb_list_accounts`, `slurmdb_get_account`, `slurmdb_list_users`,
`slurmdb_get_user`, `slurmdb_list_jobs`, `slurmdb_get_job`,
`slurmdb_list_clusters`, `slurmdb_list_qos`, `slurmdb_list_associations`,
`slurmdb_list_tres`, `slurmdb_list_wckeys`, `slurmdb_diag`

**Multi-Cluster Tools** (multi-cluster mode only):

`slurm_list_clusters`, `slurm_cluster_status`, `slurm_cross_cluster_jobs`,
`slurm_node_health_summary`, `slurm_resource_availability`

**Analysis Tools:**

`slurm_gpu_allocation`, `slurm_queue_depth`

</details>

### MCP Resources

Read-only data automatically available as context to AI agents:

`slurm://{cluster}/nodes`, `slurm://{cluster}/jobs`, `slurm://{cluster}/partitions`

### MCP Prompts

Pre-built prompt templates: `cluster-status-summary`, `job-queue-analysis`, `node-troubleshooting`

### RBAC & Audit Logging

Control which operations AI agents can perform:

```yaml
# In clusters.yaml
rbac:
  default_access: operator  # read-only | operator | admin
  audit_log: /var/log/slurm-mcp-audit.jsonl
```

```bash
# Or via environment variables
MCP_ACCESS_LEVEL=read-only MCP_AUDIT_LOG=/tmp/audit.log slurm-mcp
```

| Access Level | Allowed Operations |
|-------------|-------------------|
| `read-only` | List, Get, Ping, Diag (all read operations) |
| `operator` | Read + Submit job, Cancel job |
| `admin` | All operations (default) |

### Docker

```bash
docker build -t slurm-mcp .
docker run -p 8080:8080 \
  -v ./clusters.yaml:/etc/slurm-mcp/clusters.yaml \
  -e MCP_BEARER_TOKEN=secret \
  slurm-mcp --config /etc/slurm-mcp/clusters.yaml
```

### Environment Variables

| Variable | Description | Required |
|----------|------------|----------|
| `SLURM_ENDPOINT` | slurmrestd URL | Yes (single-cluster) |
| `SLURM_TOKEN` | JWT token | Yes (single-cluster) |
| `SLURM_VERSION` | API version (e.g., `v0.0.40`) | No (auto-detect) |
| `SLURM_UNIX_SOCKET` | Unix socket path | No |
| `MCP_BEARER_TOKEN` | Bearer token for SSE transport | Yes (SSE mode) |
| `MCP_ACCESS_LEVEL` | RBAC access level | No (default: admin) |
| `MCP_AUDIT_LOG` | Audit log file path | No |

### Flags

| Flag | Description | Default |
|------|------------|---------|
| `--config` | Path to clusters.yaml | — |
| `--transport` | Transport: `stdio` or `sse` | `stdio` |
| `--port` | Port for SSE transport | `8080` |
| `--bearer-token` | Auth token for SSE | — |
| `--all-tools` | Expose all 31 tools | `false` (8 core) |
| `--backend` | Backend: `rest` or `ssh` | `rest` |

---

## Go SDK

Go SDK for the [Slurm REST API][rest-api], inspired by [openstacksdk].
Proxy-style client with IDE-friendly auto-completion chains.

### Install

```bash
go get github.com/supergate-hub/slurm-client
```

### Quick Start

```go
package main

import (
    "context"
    "fmt"

    slurm "github.com/supergate-hub/slurm-client"
)

func main() {
    ctx := context.Background()

    client, err := slurm.NewClient(ctx, slurm.AuthOpts{
        Endpoint:  "http://slurmctld:6820",
        AuthToken: "your-jwt-token",
        Version:   slurm.V0040, // or omit for auto-detection
    })
    if err != nil {
        panic(err)
    }

    // List jobs
    jobs, _ := client.Slurm.Jobs().List(ctx)
    for _, j := range jobs {
        fmt.Printf("Job %d: %s\n", *j.JobId, *j.Name)
    }

    // Get a specific node
    node, _ := client.Slurm.Nodes().Get(ctx, "node01")
    fmt.Printf("Node: %s\n", *node.Name)

    // Submit a job
    result, _ := client.Slurm.Jobs().Submit(ctx, slurm.JobSubmitOpts{
        Script: "#!/bin/bash\nsrun hostname",
    })
    fmt.Printf("Submitted job: %d\n", result.JobID)

    // SlurmDB accounts
    accounts, _ := client.Slurmdb.Accounts().List(ctx)
    for _, a := range accounts {
        fmt.Printf("Account: %s\n", a.Name)
    }
}
```

### Multi-Cluster

```go
manager, err := slurm.NewManager(ctx, slurm.ManagerOpts{
    Clusters: map[string]slurm.AuthOpts{
        "gpu": {Endpoint: "http://gpu-cluster:6820", AuthToken: "token1"},
        "cpu": {Endpoint: "http://cpu-cluster:6820", AuthToken: "token2"},
    },
})

client, _ := manager.Cluster("gpu")
jobs, _ := client.Slurm.Jobs().List(ctx)

// Query all clusters in parallel
allJobs, _ := manager.AllJobs(ctx)
for _, r := range allJobs {
    fmt.Printf("[%s] %d jobs (err=%v)\n", r.Cluster, len(r.Data), r.Err)
}
```

Or load from YAML:

```yaml
# clusters.yaml
clusters:
  gpu:
    endpoint: http://gpu-cluster:6820
    token: ${GPU_SLURM_TOKEN}
  cpu:
    endpoint: http://cpu-cluster:6820
    token: ${CPU_SLURM_TOKEN}
```

```go
opts, _ := slurm.ParseConfig("clusters.yaml")
manager, _ := slurm.NewManager(ctx, *opts)
```

### SDK Features

- **Proxy pattern**: `client.Slurm.Jobs().List(ctx)` with IDE auto-completion
- **Multi-cluster Manager**: route to specific clusters or query all in parallel
- **YAML config**: `clusters.yaml` with `${ENV_VAR}` expansion
- **Version auto-detection**: queries slurmrestd's `/openapi` endpoint
- **Retry with backoff**: automatic retry on 503, 429, and connection errors
- **Structured logging**: `*slog.Logger` via `AuthOpts.Logger` (tokens always masked)
- **Unix socket support**: `AuthOpts.UnixSocket` for local slurmrestd
- **Full API coverage**: 9 Slurm + 11 SlurmDB resources (v0.0.39 ~ v0.0.44)
- **Mock server**: `slurmtest` package for unit testing

### Supported Resources

<details>
<summary>Slurm API</summary>

| Resource | Operations |
|----------|-----------|
| Jobs | List, Get, Submit, Update, Delete |
| Nodes | List, Get, Delete |
| Partitions | List, Get |
| Reservations | List, Get |
| Licenses | List |
| Shares | Get |
| Diag | Get |
| Ping | Get |
| Reconfigure | Reconfigure |

</details>

<details>
<summary>SlurmDB API</summary>

| Resource | Operations |
|----------|-----------|
| Accounts | List, Get |
| Users | List, Get |
| Associations | List |
| Clusters | List, Get |
| QOS | List, Get |
| Jobs | List, Get |
| WCKeys | List |
| TRES | List |
| Instances | List |
| Config | Get |
| Diag | Get |

</details>

### Testing with Mock Server

```go
import (
    slurm "github.com/supergate-hub/slurm-client"
    "github.com/supergate-hub/slurm-client/slurmtest"
)

func TestMyCode(t *testing.T) {
    mock := slurmtest.NewServer()
    defer mock.Close()

    mock.AddJob(slurmtest.MockJob{ID: 123, Name: "test-job"})
    mock.AddNode(slurmtest.MockNode{Name: "node01"})

    client, _ := slurm.NewClient(ctx, slurm.AuthOpts{
        Endpoint:  mock.URL(),
        AuthToken: "test",
        Version:   slurm.V0040,
    })

    jobs, _ := client.Slurm.Jobs().List(ctx)
    // assert len(jobs) == 1
}
```

### SDK Configuration

| Option | Description | Default |
|--------|------------|---------|
| `Endpoint` | slurmrestd HTTP URL | required |
| `AuthToken` | JWT token | required |
| `Version` | API version (e.g., `V0040`) | auto-detect |
| `UnixSocket` | Unix socket path | — |
| `HTTPClient` | Custom `*http.Client` | 30s timeout |
| `Logger` | `*slog.Logger` for structured logging | disabled |
| `MaxRetries` | Max retry attempts | 3 |
| `RetryBaseDelay` | Base delay for exponential backoff | 500ms |

### Supported Slurm Versions

| Constant | Slurm API Version |
|----------|-------------------|
| `V0039` | v0.0.39 |
| `V0040` | v0.0.40 |
| `V0041` | v0.0.41 |
| `V0042` | v0.0.42 |
| `V0043` | v0.0.43 |
| `V0044` | v0.0.44 |

## Adding New Versions

```bash
./hack/generate-version.sh v0.0.45 path/to/openapi-spec.json
./hack/check-compat.sh path/to/new-spec.json
```

## License

[Apache 2.0](./LICENSES/Apache-2.0.txt)
