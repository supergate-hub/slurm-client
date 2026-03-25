package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	slurmclient "github.com/supergate-hub/slurm-client"
)

func registerSlurmTools(s *server.MCPServer, cr *clusterResolver) {
	clusterDesc := ""
	if cr.isMultiCluster() {
		clusterDesc = " Use the 'cluster' parameter to target a specific cluster."
	}

	// --- Jobs ---
	s.AddTool(mcp.NewTool("slurm_list_jobs",
		mcp.WithDescription("List all jobs in the Slurm cluster."+clusterDesc),
		mcp.WithString("cluster", mcp.Description("Target cluster name (required in multi-cluster mode)")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		client, err := cr.resolve(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		jobs, err := client.Slurm.Jobs().List(ctx)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return jsonResult(jobs)
	})

	s.AddTool(mcp.NewTool("slurm_get_job",
		mcp.WithDescription("Get detailed information about a specific job by its ID."+clusterDesc),
		mcp.WithString("job_id", mcp.Required(), mcp.Description("The Slurm job ID")),
		mcp.WithString("cluster", mcp.Description("Target cluster name (required in multi-cluster mode)")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		client, err := cr.resolve(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		jobID, err := req.RequireString("job_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		job, err := client.Slurm.Jobs().Get(ctx, jobID)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return jsonResult(job)
	})

	s.AddTool(mcp.NewTool("slurm_submit_job",
		mcp.WithDescription("Submit a new batch job to the Slurm cluster."+clusterDesc),
		mcp.WithString("script", mcp.Required(), mcp.Description("The batch script content")),
		mcp.WithString("cluster", mcp.Description("Target cluster name (required in multi-cluster mode)")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		client, err := cr.resolve(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		script, err := req.RequireString("script")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		result, err := client.Slurm.Jobs().Submit(ctx, slurmclient.JobSubmitOpts{
			Script: script,
		})
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return jsonResult(result)
	})

	s.AddTool(mcp.NewTool("slurm_cancel_job",
		mcp.WithDescription("Cancel a running or pending job by its ID."+clusterDesc),
		mcp.WithString("job_id", mcp.Required(), mcp.Description("The Slurm job ID to cancel")),
		mcp.WithString("cluster", mcp.Description("Target cluster name (required in multi-cluster mode)")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		client, err := cr.resolve(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		jobID, err := req.RequireString("job_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		if err := client.Slurm.Jobs().Delete(ctx, jobID); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("Job %s cancelled", jobID)), nil
	})

	// --- Nodes ---
	s.AddTool(mcp.NewTool("slurm_list_nodes",
		mcp.WithDescription("List all nodes in the Slurm cluster with their states and resources."+clusterDesc),
		mcp.WithString("cluster", mcp.Description("Target cluster name (required in multi-cluster mode)")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		client, err := cr.resolve(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		nodes, err := client.Slurm.Nodes().List(ctx)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return jsonResult(nodes)
	})

	s.AddTool(mcp.NewTool("slurm_get_node",
		mcp.WithDescription("Get detailed information about a specific node by name."+clusterDesc),
		mcp.WithString("node_name", mcp.Required(), mcp.Description("The node name")),
		mcp.WithString("cluster", mcp.Description("Target cluster name (required in multi-cluster mode)")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		client, err := cr.resolve(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		name, err := req.RequireString("node_name")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		node, err := client.Slurm.Nodes().Get(ctx, name)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return jsonResult(node)
	})

	// --- Partitions ---
	s.AddTool(mcp.NewTool("slurm_list_partitions",
		mcp.WithDescription("List all partitions in the Slurm cluster."+clusterDesc),
		mcp.WithString("cluster", mcp.Description("Target cluster name (required in multi-cluster mode)")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		client, err := cr.resolve(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		parts, err := client.Slurm.Partitions().List(ctx)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return jsonResult(parts)
	})

	s.AddTool(mcp.NewTool("slurm_get_partition",
		mcp.WithDescription("Get detailed information about a specific partition."+clusterDesc),
		mcp.WithString("partition_name", mcp.Required(), mcp.Description("The partition name")),
		mcp.WithString("cluster", mcp.Description("Target cluster name (required in multi-cluster mode)")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		client, err := cr.resolve(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		name, err := req.RequireString("partition_name")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		part, err := client.Slurm.Partitions().Get(ctx, name)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return jsonResult(part)
	})

	// --- Reservations ---
	s.AddTool(mcp.NewTool("slurm_list_reservations",
		mcp.WithDescription("List all reservations in the Slurm cluster."+clusterDesc),
		mcp.WithString("cluster", mcp.Description("Target cluster name (required in multi-cluster mode)")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		client, err := cr.resolve(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		res, err := client.Slurm.Reservations().List(ctx)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return jsonResult(res)
	})

	// --- Licenses ---
	s.AddTool(mcp.NewTool("slurm_list_licenses",
		mcp.WithDescription("List all licenses tracked by Slurm."+clusterDesc),
		mcp.WithString("cluster", mcp.Description("Target cluster name (required in multi-cluster mode)")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		client, err := cr.resolve(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		lic, err := client.Slurm.Licenses().List(ctx)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return jsonResult(lic)
	})

	// --- Diag ---
	s.AddTool(mcp.NewTool("slurm_diag",
		mcp.WithDescription("Get Slurm controller diagnostics and statistics."+clusterDesc),
		mcp.WithString("cluster", mcp.Description("Target cluster name (required in multi-cluster mode)")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		client, err := cr.resolve(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		diag, err := client.Slurm.Diag().Get(ctx)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return jsonResult(diag)
	})

	// --- Ping ---
	s.AddTool(mcp.NewTool("slurm_ping",
		mcp.WithDescription("Ping Slurm controllers to check if they are responding."+clusterDesc),
		mcp.WithString("cluster", mcp.Description("Target cluster name (required in multi-cluster mode)")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		client, err := cr.resolve(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		pings, err := client.Slurm.Ping().Get(ctx)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return jsonResult(pings)
	})
}

func jsonResult(v any) (*mcp.CallToolResult, error) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("marshal error: %v", err)), nil
	}
	return mcp.NewToolResultText(string(data)), nil
}
