package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	slurmclient "github.com/supergate-hub/slurm-client"
)

func registerSlurmTools(s *server.MCPServer, client *slurmclient.Client) {
	// --- Jobs ---
	s.AddTool(mcp.NewTool("slurm_list_jobs",
		mcp.WithDescription("List all jobs in the Slurm cluster. Returns job IDs, names, states, and resource allocations."),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		jobs, err := client.Slurm.Jobs().List(ctx)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return jsonResult(jobs)
	})

	s.AddTool(mcp.NewTool("slurm_get_job",
		mcp.WithDescription("Get detailed information about a specific job by its ID."),
		mcp.WithString("job_id", mcp.Required(), mcp.Description("The Slurm job ID")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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
		mcp.WithDescription("Submit a new batch job to the Slurm cluster. Requires a bash script."),
		mcp.WithString("script", mcp.Required(), mcp.Description("The batch script content (e.g., '#!/bin/bash\\nsrun hostname')")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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
		mcp.WithDescription("Cancel a running or pending job by its ID."),
		mcp.WithString("job_id", mcp.Required(), mcp.Description("The Slurm job ID to cancel")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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
		mcp.WithDescription("List all nodes in the Slurm cluster with their states and resources."),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		nodes, err := client.Slurm.Nodes().List(ctx)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return jsonResult(nodes)
	})

	s.AddTool(mcp.NewTool("slurm_get_node",
		mcp.WithDescription("Get detailed information about a specific node by name."),
		mcp.WithString("node_name", mcp.Required(), mcp.Description("The node name")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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
		mcp.WithDescription("List all partitions in the Slurm cluster with their configurations."),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		parts, err := client.Slurm.Partitions().List(ctx)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return jsonResult(parts)
	})

	s.AddTool(mcp.NewTool("slurm_get_partition",
		mcp.WithDescription("Get detailed information about a specific partition by name."),
		mcp.WithString("partition_name", mcp.Required(), mcp.Description("The partition name")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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
		mcp.WithDescription("List all reservations in the Slurm cluster."),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		res, err := client.Slurm.Reservations().List(ctx)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return jsonResult(res)
	})

	// --- Licenses ---
	s.AddTool(mcp.NewTool("slurm_list_licenses",
		mcp.WithDescription("List all licenses tracked by Slurm."),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		lic, err := client.Slurm.Licenses().List(ctx)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return jsonResult(lic)
	})

	// --- Diag ---
	s.AddTool(mcp.NewTool("slurm_diag",
		mcp.WithDescription("Get Slurm controller diagnostics and statistics."),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		diag, err := client.Slurm.Diag().Get(ctx)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return jsonResult(diag)
	})

	// --- Ping ---
	s.AddTool(mcp.NewTool("slurm_ping",
		mcp.WithDescription("Ping Slurm controllers to check if they are responding."),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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
