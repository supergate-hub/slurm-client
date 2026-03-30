package main

import (
	"context"
	"encoding/json"
	"fmt"

	slurmclient "github.com/supergate-hub/slurm-client"
)

// RestBackend implements MCPBackend using the Slurm REST API via the Go SDK.
type RestBackend struct {
	client *slurmclient.Client
}

// NewRestBackend creates a RestBackend wrapping the given SDK client.
func NewRestBackend(client *slurmclient.Client) MCPBackend {
	return &RestBackend{client: client}
}

func marshalJSON(v any) (string, error) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshal error: %w", err)
	}
	return string(data), nil
}

// --- Slurm operations ---

func (b *RestBackend) ListJobs(ctx context.Context) (string, error) {
	jobs, err := b.client.Slurm.Jobs().List(ctx)
	if err != nil {
		return "", err
	}
	return marshalJSON(jobs)
}

func (b *RestBackend) GetJob(ctx context.Context, jobID string) (string, error) {
	job, err := b.client.Slurm.Jobs().Get(ctx, jobID)
	if err != nil {
		return "", err
	}
	return marshalJSON(job)
}

func (b *RestBackend) SubmitJob(ctx context.Context, script string) (string, error) {
	result, err := b.client.Slurm.Jobs().Submit(ctx, slurmclient.JobSubmitOpts{
		Script: script,
	})
	if err != nil {
		return "", err
	}
	return marshalJSON(result)
}

func (b *RestBackend) CancelJob(ctx context.Context, jobID string) error {
	return b.client.Slurm.Jobs().Delete(ctx, jobID)
}

func (b *RestBackend) ListNodes(ctx context.Context) (string, error) {
	nodes, err := b.client.Slurm.Nodes().List(ctx)
	if err != nil {
		return "", err
	}
	return marshalJSON(nodes)
}

func (b *RestBackend) GetNode(ctx context.Context, nodeName string) (string, error) {
	node, err := b.client.Slurm.Nodes().Get(ctx, nodeName)
	if err != nil {
		return "", err
	}
	return marshalJSON(node)
}

func (b *RestBackend) ListPartitions(ctx context.Context) (string, error) {
	parts, err := b.client.Slurm.Partitions().List(ctx)
	if err != nil {
		return "", err
	}
	return marshalJSON(parts)
}

func (b *RestBackend) GetPartition(ctx context.Context, partitionName string) (string, error) {
	part, err := b.client.Slurm.Partitions().Get(ctx, partitionName)
	if err != nil {
		return "", err
	}
	return marshalJSON(part)
}

func (b *RestBackend) ListReservations(ctx context.Context) (string, error) {
	res, err := b.client.Slurm.Reservations().List(ctx)
	if err != nil {
		return "", err
	}
	return marshalJSON(res)
}

func (b *RestBackend) ListLicenses(ctx context.Context) (string, error) {
	lic, err := b.client.Slurm.Licenses().List(ctx)
	if err != nil {
		return "", err
	}
	return marshalJSON(lic)
}

func (b *RestBackend) GetDiag(ctx context.Context) (string, error) {
	diag, err := b.client.Slurm.Diag().Get(ctx)
	if err != nil {
		return "", err
	}
	return marshalJSON(diag)
}

func (b *RestBackend) Ping(ctx context.Context) (string, error) {
	pings, err := b.client.Slurm.Ping().Get(ctx)
	if err != nil {
		return "", err
	}
	return marshalJSON(pings)
}

// --- SlurmDB operations ---

func (b *RestBackend) ListAccounts(ctx context.Context) (string, error) {
	accts, err := b.client.Slurmdb.Accounts().List(ctx)
	if err != nil {
		return "", err
	}
	return marshalJSON(accts)
}

func (b *RestBackend) GetAccount(ctx context.Context, name string) (string, error) {
	acct, err := b.client.Slurmdb.Accounts().Get(ctx, name)
	if err != nil {
		return "", err
	}
	return marshalJSON(acct)
}

func (b *RestBackend) ListUsers(ctx context.Context) (string, error) {
	users, err := b.client.Slurmdb.Users().List(ctx)
	if err != nil {
		return "", err
	}
	return marshalJSON(users)
}

func (b *RestBackend) GetUser(ctx context.Context, name string) (string, error) {
	user, err := b.client.Slurmdb.Users().Get(ctx, name)
	if err != nil {
		return "", err
	}
	return marshalJSON(user)
}

func (b *RestBackend) ListSlurmdbJobs(ctx context.Context) (string, error) {
	jobs, err := b.client.Slurmdb.Jobs().List(ctx)
	if err != nil {
		return "", err
	}
	return marshalJSON(jobs)
}

func (b *RestBackend) GetSlurmdbJob(ctx context.Context, jobID string) (string, error) {
	job, err := b.client.Slurmdb.Jobs().Get(ctx, jobID)
	if err != nil {
		return "", err
	}
	return marshalJSON(job)
}

func (b *RestBackend) ListClusters(ctx context.Context) (string, error) {
	clusters, err := b.client.Slurmdb.Clusters().List(ctx)
	if err != nil {
		return "", err
	}
	return marshalJSON(clusters)
}

func (b *RestBackend) ListQOS(ctx context.Context) (string, error) {
	qos, err := b.client.Slurmdb.QOS().List(ctx)
	if err != nil {
		return "", err
	}
	return marshalJSON(qos)
}

func (b *RestBackend) ListAssociations(ctx context.Context) (string, error) {
	assocs, err := b.client.Slurmdb.Associations().List(ctx)
	if err != nil {
		return "", err
	}
	return marshalJSON(assocs)
}

func (b *RestBackend) ListTRES(ctx context.Context) (string, error) {
	tres, err := b.client.Slurmdb.TRES().List(ctx)
	if err != nil {
		return "", err
	}
	return marshalJSON(tres)
}

func (b *RestBackend) ListWCKeys(ctx context.Context) (string, error) {
	wckeys, err := b.client.Slurmdb.WCKeys().List(ctx)
	if err != nil {
		return "", err
	}
	return marshalJSON(wckeys)
}

func (b *RestBackend) GetSlurmdbDiag(ctx context.Context) (string, error) {
	diag, err := b.client.Slurmdb.SlurmdbDiag().Get(ctx)
	if err != nil {
		return "", err
	}
	return marshalJSON(diag)
}
