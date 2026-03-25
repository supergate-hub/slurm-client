package main

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerSlurmdbTools(s *server.MCPServer, cr *clusterResolver) {
	clusterDesc := ""
	if cr.isMultiCluster() {
		clusterDesc = " Use the 'cluster' parameter to target a specific cluster."
	}

	// --- Accounts ---
	s.AddTool(mcp.NewTool("slurmdb_list_accounts",
		mcp.WithDescription("List all accounts in the Slurm accounting database."+clusterDesc),
		mcp.WithString("cluster", mcp.Description("Target cluster name (required in multi-cluster mode)")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		client, err := cr.resolve(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		accts, err := client.Slurmdb.Accounts().List(ctx)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return jsonResult(accts)
	})

	s.AddTool(mcp.NewTool("slurmdb_get_account",
		mcp.WithDescription("Get detailed information about a specific account."+clusterDesc),
		mcp.WithString("account_name", mcp.Required(), mcp.Description("The account name")),
		mcp.WithString("cluster", mcp.Description("Target cluster name (required in multi-cluster mode)")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		client, err := cr.resolve(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		name, err := req.RequireString("account_name")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		acct, err := client.Slurmdb.Accounts().Get(ctx, name)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return jsonResult(acct)
	})

	// --- Users ---
	s.AddTool(mcp.NewTool("slurmdb_list_users",
		mcp.WithDescription("List all users in the Slurm accounting database."+clusterDesc),
		mcp.WithString("cluster", mcp.Description("Target cluster name (required in multi-cluster mode)")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		client, err := cr.resolve(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		users, err := client.Slurmdb.Users().List(ctx)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return jsonResult(users)
	})

	s.AddTool(mcp.NewTool("slurmdb_get_user",
		mcp.WithDescription("Get detailed information about a specific user."+clusterDesc),
		mcp.WithString("user_name", mcp.Required(), mcp.Description("The user name")),
		mcp.WithString("cluster", mcp.Description("Target cluster name (required in multi-cluster mode)")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		client, err := cr.resolve(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		name, err := req.RequireString("user_name")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		user, err := client.Slurmdb.Users().Get(ctx, name)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return jsonResult(user)
	})

	// --- Associations ---
	s.AddTool(mcp.NewTool("slurmdb_list_associations",
		mcp.WithDescription("List all associations in the Slurm accounting database."+clusterDesc),
		mcp.WithString("cluster", mcp.Description("Target cluster name (required in multi-cluster mode)")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		client, err := cr.resolve(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		assocs, err := client.Slurmdb.Associations().List(ctx)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return jsonResult(assocs)
	})

	// --- Clusters ---
	s.AddTool(mcp.NewTool("slurmdb_list_clusters",
		mcp.WithDescription("List all clusters in the Slurm accounting database."+clusterDesc),
		mcp.WithString("cluster", mcp.Description("Target cluster name (required in multi-cluster mode)")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		client, err := cr.resolve(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		clusters, err := client.Slurmdb.Clusters().List(ctx)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return jsonResult(clusters)
	})

	// --- QOS ---
	s.AddTool(mcp.NewTool("slurmdb_list_qos",
		mcp.WithDescription("List all Quality of Service entries."+clusterDesc),
		mcp.WithString("cluster", mcp.Description("Target cluster name (required in multi-cluster mode)")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		client, err := cr.resolve(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		qos, err := client.Slurmdb.QOS().List(ctx)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return jsonResult(qos)
	})

	// --- SlurmDB Jobs ---
	s.AddTool(mcp.NewTool("slurmdb_list_jobs",
		mcp.WithDescription("List job accounting records from the Slurm database."+clusterDesc),
		mcp.WithString("cluster", mcp.Description("Target cluster name (required in multi-cluster mode)")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		client, err := cr.resolve(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		jobs, err := client.Slurmdb.Jobs().List(ctx)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return jsonResult(jobs)
	})

	s.AddTool(mcp.NewTool("slurmdb_get_job",
		mcp.WithDescription("Get accounting record for a specific job."+clusterDesc),
		mcp.WithString("job_id", mcp.Required(), mcp.Description("The job ID")),
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
		job, err := client.Slurmdb.Jobs().Get(ctx, jobID)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return jsonResult(job)
	})

	// --- TRES ---
	s.AddTool(mcp.NewTool("slurmdb_list_tres",
		mcp.WithDescription("List all Trackable Resources (TRES)."+clusterDesc),
		mcp.WithString("cluster", mcp.Description("Target cluster name (required in multi-cluster mode)")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		client, err := cr.resolve(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		tres, err := client.Slurmdb.TRES().List(ctx)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return jsonResult(tres)
	})

	// --- WCKeys ---
	s.AddTool(mcp.NewTool("slurmdb_list_wckeys",
		mcp.WithDescription("List all Workload Characterization Keys."+clusterDesc),
		mcp.WithString("cluster", mcp.Description("Target cluster name (required in multi-cluster mode)")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		client, err := cr.resolve(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		wckeys, err := client.Slurmdb.WCKeys().List(ctx)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return jsonResult(wckeys)
	})

	// --- SlurmDB Diag ---
	s.AddTool(mcp.NewTool("slurmdb_diag",
		mcp.WithDescription("Get SlurmDB (slurmdbd) diagnostics and statistics."+clusterDesc),
		mcp.WithString("cluster", mcp.Description("Target cluster name (required in multi-cluster mode)")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		client, err := cr.resolve(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		diag, err := client.Slurmdb.SlurmdbDiag().Get(ctx)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return jsonResult(diag)
	})
}
