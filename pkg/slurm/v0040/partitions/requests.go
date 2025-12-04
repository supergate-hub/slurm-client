package partitions

import (
	"context"
	"fmt"
	"net/url"

	api "github.com/supergate-hub/slurm-client/api/v0040"
	"github.com/supergate-hub/slurm-client/pkg/slurmclient"
)

// List retrieves all partitions from the Slurm cluster.
func List(ctx context.Context, client *slurmclient.ServiceClient, opts ListOpts) ([]Partition, error) {
	query := url.Values{}
	if opts.UpdateTime != nil {
		query.Set("update_time", fmt.Sprintf("%d", *opts.UpdateTime))
	}

	resp, err := client.GetWithQuery(ctx, listURL(client), query)
	if err != nil {
		return nil, err
	}

	var result api.V0040OpenapiPartitionResp
	if err := slurmclient.DecodeResponse(resp, &result); err != nil {
		return nil, err
	}

	if len(result.Partitions) == 0 {
		return []Partition{}, nil
	}
	return result.Partitions, nil
}

// Get retrieves a single partition by name.
func Get(ctx context.Context, client *slurmclient.ServiceClient, partitionName string, opts ...GetOpts) (*Partition, error) {
	query := url.Values{}
	if len(opts) > 0 && opts[0].UpdateTime != nil {
		query.Set("update_time", fmt.Sprintf("%d", *opts[0].UpdateTime))
	}

	resp, err := client.GetWithQuery(ctx, getURL(client, partitionName), query)
	if err != nil {
		return nil, err
	}

	var result api.V0040OpenapiPartitionResp
	if err := slurmclient.DecodeResponse(resp, &result); err != nil {
		return nil, err
	}

	if len(result.Partitions) == 0 {
		return nil, fmt.Errorf("partition %s not found", partitionName)
	}

	return &result.Partitions[0], nil
}
