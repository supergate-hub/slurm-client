// SPDX-FileCopyrightText: Copyright (C) SchedMD LLC.
// SPDX-License-Identifier: Apache-2.0

package v0040

import (
	"context"

	"github.com/supergate-hub/slurm-client/pkg/slurm/v0040/nodes"
	"github.com/supergate-hub/slurm-client/pkg/slurmclient"
)

// nodesClient implements NodesInterface.
type nodesClient struct {
	client *slurmclient.ServiceClient
}

var _ NodesInterface = (*nodesClient)(nil)

// List retrieves all nodes.
func (c *nodesClient) List(ctx context.Context, opts nodes.ListOpts) ([]Node, error) {
	return nodes.List(ctx, c.client, opts)
}

// Get retrieves a single node by name.
func (c *nodesClient) Get(ctx context.Context, nodeName string, opts ...nodes.GetOpts) (*Node, error) {
	return nodes.Get(ctx, c.client, nodeName, opts...)
}

// Update modifies a node's properties.
func (c *nodesClient) Update(ctx context.Context, nodeName string, opts nodes.UpdateOpts) error {
	return nodes.Update(ctx, c.client, nodeName, opts)
}

// Delete removes a node from Slurm's configuration.
func (c *nodesClient) Delete(ctx context.Context, nodeName string) error {
	return nodes.Delete(ctx, c.client, nodeName)
}



