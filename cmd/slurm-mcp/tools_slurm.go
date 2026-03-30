package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
)

// toolAdder is the interface for registering MCP tools.
// Both *server.MCPServer and *toolRegistrar satisfy this.
type toolAdder interface {
	AddTool(tool mcp.Tool, handler func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error))
}

func registerSlurmTools(s toolAdder, cr *clusterResolver, rbac *RBAC) {
	clusterDesc := ""
	if cr.isMultiCluster() {
		clusterDesc = " Use the 'cluster' parameter to target a specific cluster."
	}

	// --- Jobs ---
	s.AddTool(mcp.NewTool("slurm_list_jobs",
		mcp.WithDescription("List all jobs in the Slurm cluster."+clusterDesc),
		mcp.WithString("cluster", mcp.Description("Target cluster name (required in multi-cluster mode)")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if err := checkRBAC(rbac, "slurm_list_jobs"); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		backend, err := cr.resolveBackend(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		result, err := backend.ListJobs(ctx)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(result), nil
	})

	s.AddTool(mcp.NewTool("slurm_get_job",
		mcp.WithDescription("Get detailed information about a specific job by its ID."+clusterDesc),
		mcp.WithString("job_id", mcp.Required(), mcp.Description("The Slurm job ID")),
		mcp.WithString("cluster", mcp.Description("Target cluster name (required in multi-cluster mode)")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if err := checkRBAC(rbac, "slurm_get_job"); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		backend, err := cr.resolveBackend(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		jobID, err := req.RequireString("job_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		result, err := backend.GetJob(ctx, jobID)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(result), nil
	})

	s.AddTool(mcp.NewTool("slurm_submit_job",
		mcp.WithDescription("Submit a new batch job to the Slurm cluster."+clusterDesc),
		mcp.WithString("script", mcp.Required(), mcp.Description("The batch script content")),
		mcp.WithString("cluster", mcp.Description("Target cluster name (required in multi-cluster mode)")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if err := checkRBAC(rbac, "slurm_submit_job"); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		backend, err := cr.resolveBackend(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		script, err := req.RequireString("script")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		result, err := backend.SubmitJob(ctx, script)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(result), nil
	})

	s.AddTool(mcp.NewTool("slurm_cancel_job",
		mcp.WithDescription("Cancel a running or pending job by its ID."+clusterDesc),
		mcp.WithString("job_id", mcp.Required(), mcp.Description("The Slurm job ID to cancel")),
		mcp.WithString("cluster", mcp.Description("Target cluster name (required in multi-cluster mode)")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if err := checkRBAC(rbac, "slurm_cancel_job"); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		backend, err := cr.resolveBackend(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		jobID, err := req.RequireString("job_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		if err := backend.CancelJob(ctx, jobID); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("Job %s cancelled", jobID)), nil
	})

	// --- Nodes ---
	s.AddTool(mcp.NewTool("slurm_list_nodes",
		mcp.WithDescription("List all nodes in the Slurm cluster with their states and resources."+clusterDesc),
		mcp.WithString("cluster", mcp.Description("Target cluster name (required in multi-cluster mode)")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if err := checkRBAC(rbac, "slurm_list_nodes"); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		backend, err := cr.resolveBackend(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		result, err := backend.ListNodes(ctx)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(result), nil
	})

	s.AddTool(mcp.NewTool("slurm_get_node",
		mcp.WithDescription("Get detailed information about a specific node by name."+clusterDesc),
		mcp.WithString("node_name", mcp.Required(), mcp.Description("The node name")),
		mcp.WithString("cluster", mcp.Description("Target cluster name (required in multi-cluster mode)")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if err := checkRBAC(rbac, "slurm_get_node"); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		backend, err := cr.resolveBackend(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		name, err := req.RequireString("node_name")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		result, err := backend.GetNode(ctx, name)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(result), nil
	})

	// --- Partitions ---
	s.AddTool(mcp.NewTool("slurm_list_partitions",
		mcp.WithDescription("List all partitions in the Slurm cluster."+clusterDesc),
		mcp.WithString("cluster", mcp.Description("Target cluster name (required in multi-cluster mode)")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if err := checkRBAC(rbac, "slurm_list_partitions"); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		backend, err := cr.resolveBackend(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		result, err := backend.ListPartitions(ctx)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(result), nil
	})

	s.AddTool(mcp.NewTool("slurm_get_partition",
		mcp.WithDescription("Get detailed information about a specific partition."+clusterDesc),
		mcp.WithString("partition_name", mcp.Required(), mcp.Description("The partition name")),
		mcp.WithString("cluster", mcp.Description("Target cluster name (required in multi-cluster mode)")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if err := checkRBAC(rbac, "slurm_get_partition"); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		backend, err := cr.resolveBackend(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		name, err := req.RequireString("partition_name")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		result, err := backend.GetPartition(ctx, name)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(result), nil
	})

	// --- Reservations ---
	s.AddTool(mcp.NewTool("slurm_list_reservations",
		mcp.WithDescription("List all reservations in the Slurm cluster."+clusterDesc),
		mcp.WithString("cluster", mcp.Description("Target cluster name (required in multi-cluster mode)")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if err := checkRBAC(rbac, "slurm_list_reservations"); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		backend, err := cr.resolveBackend(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		result, err := backend.ListReservations(ctx)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(result), nil
	})

	// --- Licenses ---
	s.AddTool(mcp.NewTool("slurm_list_licenses",
		mcp.WithDescription("List all licenses tracked by Slurm."+clusterDesc),
		mcp.WithString("cluster", mcp.Description("Target cluster name (required in multi-cluster mode)")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if err := checkRBAC(rbac, "slurm_list_licenses"); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		backend, err := cr.resolveBackend(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		result, err := backend.ListLicenses(ctx)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(result), nil
	})

	// --- Diag ---
	s.AddTool(mcp.NewTool("slurm_diag",
		mcp.WithDescription("Get Slurm controller diagnostics and statistics."+clusterDesc),
		mcp.WithString("cluster", mcp.Description("Target cluster name (required in multi-cluster mode)")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if err := checkRBAC(rbac, "slurm_diag"); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		backend, err := cr.resolveBackend(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		result, err := backend.GetDiag(ctx)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(result), nil
	})

	// --- Ping ---
	s.AddTool(mcp.NewTool("slurm_ping",
		mcp.WithDescription("Ping Slurm controllers to check if they are responding."+clusterDesc),
		mcp.WithString("cluster", mcp.Description("Target cluster name (required in multi-cluster mode)")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if err := checkRBAC(rbac, "slurm_ping"); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		backend, err := cr.resolveBackend(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		result, err := backend.Ping(ctx)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(result), nil
	})
}

func jsonResult(v any) (*mcp.CallToolResult, error) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("marshal error: %v", err)), nil
	}
	return mcp.NewToolResultText(string(data)), nil
}

func jsonMarshal(v any) ([]byte, error)     { return json.Marshal(v) }
func jsonUnmarshal(data []byte, v any) error { return json.Unmarshal(data, v) }
