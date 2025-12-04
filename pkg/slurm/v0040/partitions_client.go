// SPDX-FileCopyrightText: Copyright (C) SchedMD LLC.
// SPDX-License-Identifier: Apache-2.0

package v0040

import (
	"context"

	"github.com/supergate-hub/slurm-client/pkg/slurm/v0040/partitions"
	"github.com/supergate-hub/slurm-client/pkg/slurmclient"
)

// partitionsClient implements PartitionsInterface.
type partitionsClient struct {
	client *slurmclient.ServiceClient
}

var _ PartitionsInterface = (*partitionsClient)(nil)

// List retrieves all partitions.
func (c *partitionsClient) List(ctx context.Context, opts partitions.ListOpts) ([]Partition, error) {
	return partitions.List(ctx, c.client, opts)
}

// Get retrieves a single partition by name.
func (c *partitionsClient) Get(ctx context.Context, partitionName string, opts ...partitions.GetOpts) (*Partition, error) {
	return partitions.Get(ctx, c.client, partitionName, opts...)
}



