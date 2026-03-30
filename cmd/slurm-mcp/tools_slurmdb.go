package main

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
)

func registerSlurmdbTools(s toolAdder, cr *clusterResolver, rbac *RBAC) {
	clusterDesc := ""
	if cr.isMultiCluster() {
		clusterDesc = " Use the 'cluster' parameter to target a specific cluster."
	}

	// --- Accounts ---
	s.AddTool(mcp.NewTool("slurmdb_list_accounts",
		mcp.WithDescription("List all accounts in the Slurm accounting database."+clusterDesc),
		mcp.WithString("cluster", mcp.Description("Target cluster name (required in multi-cluster mode)")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if err := checkRBAC(rbac, "slurmdb_list_accounts"); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		backend, err := cr.resolveBackend(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		result, err := backend.ListAccounts(ctx)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(result), nil
	})

	s.AddTool(mcp.NewTool("slurmdb_get_account",
		mcp.WithDescription("Get detailed information about a specific account."+clusterDesc),
		mcp.WithString("account_name", mcp.Required(), mcp.Description("The account name")),
		mcp.WithString("cluster", mcp.Description("Target cluster name (required in multi-cluster mode)")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if err := checkRBAC(rbac, "slurmdb_get_account"); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		backend, err := cr.resolveBackend(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		name, err := req.RequireString("account_name")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		result, err := backend.GetAccount(ctx, name)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(result), nil
	})

	// --- Users ---
	s.AddTool(mcp.NewTool("slurmdb_list_users",
		mcp.WithDescription("List all users in the Slurm accounting database."+clusterDesc),
		mcp.WithString("cluster", mcp.Description("Target cluster name (required in multi-cluster mode)")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if err := checkRBAC(rbac, "slurmdb_list_users"); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		backend, err := cr.resolveBackend(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		result, err := backend.ListUsers(ctx)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(result), nil
	})

	s.AddTool(mcp.NewTool("slurmdb_get_user",
		mcp.WithDescription("Get detailed information about a specific user."+clusterDesc),
		mcp.WithString("user_name", mcp.Required(), mcp.Description("The user name")),
		mcp.WithString("cluster", mcp.Description("Target cluster name (required in multi-cluster mode)")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if err := checkRBAC(rbac, "slurmdb_get_user"); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		backend, err := cr.resolveBackend(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		name, err := req.RequireString("user_name")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		result, err := backend.GetUser(ctx, name)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(result), nil
	})

	// --- Associations ---
	s.AddTool(mcp.NewTool("slurmdb_list_associations",
		mcp.WithDescription("List all associations in the Slurm accounting database."+clusterDesc),
		mcp.WithString("cluster", mcp.Description("Target cluster name (required in multi-cluster mode)")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if err := checkRBAC(rbac, "slurmdb_list_associations"); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		backend, err := cr.resolveBackend(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		result, err := backend.ListAssociations(ctx)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(result), nil
	})

	// --- Clusters ---
	s.AddTool(mcp.NewTool("slurmdb_list_clusters",
		mcp.WithDescription("List all clusters in the Slurm accounting database."+clusterDesc),
		mcp.WithString("cluster", mcp.Description("Target cluster name (required in multi-cluster mode)")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if err := checkRBAC(rbac, "slurmdb_list_clusters"); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		backend, err := cr.resolveBackend(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		result, err := backend.ListClusters(ctx)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(result), nil
	})

	// --- QOS ---
	s.AddTool(mcp.NewTool("slurmdb_list_qos",
		mcp.WithDescription("List all Quality of Service entries."+clusterDesc),
		mcp.WithString("cluster", mcp.Description("Target cluster name (required in multi-cluster mode)")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if err := checkRBAC(rbac, "slurmdb_list_qos"); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		backend, err := cr.resolveBackend(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		result, err := backend.ListQOS(ctx)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(result), nil
	})

	// --- SlurmDB Jobs ---
	s.AddTool(mcp.NewTool("slurmdb_list_jobs",
		mcp.WithDescription("List job accounting records from the Slurm database."+clusterDesc),
		mcp.WithString("cluster", mcp.Description("Target cluster name (required in multi-cluster mode)")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if err := checkRBAC(rbac, "slurmdb_list_jobs"); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		backend, err := cr.resolveBackend(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		result, err := backend.ListSlurmdbJobs(ctx)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(result), nil
	})

	s.AddTool(mcp.NewTool("slurmdb_get_job",
		mcp.WithDescription("Get accounting record for a specific job."+clusterDesc),
		mcp.WithString("job_id", mcp.Required(), mcp.Description("The job ID")),
		mcp.WithString("cluster", mcp.Description("Target cluster name (required in multi-cluster mode)")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if err := checkRBAC(rbac, "slurmdb_get_job"); err != nil {
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
		result, err := backend.GetSlurmdbJob(ctx, jobID)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(result), nil
	})

	// --- TRES ---
	s.AddTool(mcp.NewTool("slurmdb_list_tres",
		mcp.WithDescription("List all Trackable Resources (TRES)."+clusterDesc),
		mcp.WithString("cluster", mcp.Description("Target cluster name (required in multi-cluster mode)")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if err := checkRBAC(rbac, "slurmdb_list_tres"); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		backend, err := cr.resolveBackend(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		result, err := backend.ListTRES(ctx)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(result), nil
	})

	// --- WCKeys ---
	s.AddTool(mcp.NewTool("slurmdb_list_wckeys",
		mcp.WithDescription("List all Workload Characterization Keys."+clusterDesc),
		mcp.WithString("cluster", mcp.Description("Target cluster name (required in multi-cluster mode)")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if err := checkRBAC(rbac, "slurmdb_list_wckeys"); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		backend, err := cr.resolveBackend(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		result, err := backend.ListWCKeys(ctx)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(result), nil
	})

	// --- SlurmDB Diag ---
	s.AddTool(mcp.NewTool("slurmdb_diag",
		mcp.WithDescription("Get SlurmDB (slurmdbd) diagnostics and statistics."+clusterDesc),
		mcp.WithString("cluster", mcp.Description("Target cluster name (required in multi-cluster mode)")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if err := checkRBAC(rbac, "slurmdb_diag"); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		backend, err := cr.resolveBackend(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		result, err := backend.GetSlurmdbDiag(ctx)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(result), nil
	})
}
