// SPDX-FileCopyrightText: Copyright (C) SchedMD LLC.
// SPDX-License-Identifier: Apache-2.0

package instances

import (
	"context"
	"net/url"

	api "github.com/supergate-hub/slurm-client/api/v0040"
	"github.com/supergate-hub/slurm-client/pkg/slurmclient"
)

// List retrieves all instances from the Slurm database.
func List(ctx context.Context, client *slurmclient.ServiceClient, opts ListOpts) ([]Instance, error) {
	query := url.Values{}
	if opts.Cluster != nil {
		query.Set("cluster", *opts.Cluster)
	}
	if opts.Extra != nil {
		query.Set("extra", *opts.Extra)
	}
	if opts.Format != nil {
		query.Set("format", *opts.Format)
	}
	if opts.InstanceID != nil {
		query.Set("instance_id", *opts.InstanceID)
	}
	if opts.InstanceType != nil {
		query.Set("instance_type", *opts.InstanceType)
	}
	if opts.NodeList != nil {
		query.Set("node_list", *opts.NodeList)
	}

	resp, err := client.GetWithQuery(ctx, listURL(client), query)
	if err != nil {
		return nil, err
	}

	var result api.V0040OpenapiInstancesResp
	if err := slurmclient.DecodeResponse(resp, &result); err != nil {
		return nil, err
	}

	if len(result.Instances) == 0 {
		return []Instance{}, nil
	}
	return result.Instances, nil
}

// Get retrieves a single instance based on the provided parameters.
func Get(ctx context.Context, client *slurmclient.ServiceClient, opts GetOpts) (*Instance, error) {
	query := url.Values{}
	if opts.Cluster != nil {
		query.Set("cluster", *opts.Cluster)
	}
	if opts.InstanceID != nil {
		query.Set("instance_id", *opts.InstanceID)
	}
	if opts.InstanceType != nil {
		query.Set("instance_type", *opts.InstanceType)
	}
	if opts.NodeList != nil {
		query.Set("node_list", *opts.NodeList)
	}

	resp, err := client.GetWithQuery(ctx, getURL(client), query)
	if err != nil {
		return nil, err
	}

	var result api.V0040OpenapiInstancesResp
	if err := slurmclient.DecodeResponse(resp, &result); err != nil {
		return nil, err
	}

	if len(result.Instances) == 0 {
		return nil, nil
	}

	return &result.Instances[0], nil
}
