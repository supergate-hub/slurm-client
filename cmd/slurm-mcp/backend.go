package main

import "context"

// MCPBackend abstracts data fetching for MCP tools.
// REST backend uses the Go SDK. SSH backend (Lane 2) will implement this with CLI commands.
type MCPBackend interface {
	// Slurm operations
	ListJobs(ctx context.Context) (string, error)
	GetJob(ctx context.Context, jobID string) (string, error)
	SubmitJob(ctx context.Context, script string) (string, error)
	CancelJob(ctx context.Context, jobID string) error
	ListNodes(ctx context.Context) (string, error)
	GetNode(ctx context.Context, nodeName string) (string, error)
	ListPartitions(ctx context.Context) (string, error)
	GetPartition(ctx context.Context, partitionName string) (string, error)
	ListReservations(ctx context.Context) (string, error)
	ListLicenses(ctx context.Context) (string, error)
	GetDiag(ctx context.Context) (string, error)
	Ping(ctx context.Context) (string, error)

	// SlurmDB operations
	ListAccounts(ctx context.Context) (string, error)
	GetAccount(ctx context.Context, name string) (string, error)
	ListUsers(ctx context.Context) (string, error)
	GetUser(ctx context.Context, name string) (string, error)
	ListSlurmdbJobs(ctx context.Context) (string, error)
	GetSlurmdbJob(ctx context.Context, jobID string) (string, error)
	ListClusters(ctx context.Context) (string, error)
	ListQOS(ctx context.Context) (string, error)
	ListAssociations(ctx context.Context) (string, error)
	ListTRES(ctx context.Context) (string, error)
	ListWCKeys(ctx context.Context) (string, error)
	GetSlurmdbDiag(ctx context.Context) (string, error)
}
