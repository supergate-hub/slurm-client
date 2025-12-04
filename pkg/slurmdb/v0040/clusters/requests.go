// SPDX-FileCopyrightText: Copyright (C) SchedMD LLC.
// SPDX-License-Identifier: Apache-2.0

package clusters

import (
	"context"
	"fmt"
	"net/url"

	api "github.com/supergate-hub/slurm-client/api/v0040"
	"github.com/supergate-hub/slurm-client/pkg/slurmclient"
)

// List retrieves all clusters from the Slurm database.
func List(ctx context.Context, client *slurmclient.ServiceClient, opts ListOpts) ([]Cluster, error) {
	query := url.Values{}
	if opts.WithDeleted != nil && *opts.WithDeleted {
		query.Set("with_deleted", "true")
	}

	resp, err := client.GetWithQuery(ctx, listURL(client), query)
	if err != nil {
		return nil, err
	}

	var result api.V0040OpenapiClustersResp
	if err := slurmclient.DecodeResponse(resp, &result); err != nil {
		return nil, err
	}

	if len(result.Clusters) == 0 {
		return []Cluster{}, nil
	}
	return result.Clusters, nil
}

// Get retrieves a single cluster by name.
func Get(ctx context.Context, client *slurmclient.ServiceClient, clusterName string, opts ...GetOpts) (*Cluster, error) {
	query := url.Values{}
	if len(opts) > 0 && opts[0].WithDeleted != nil && *opts[0].WithDeleted {
		query.Set("with_deleted", "true")
	}

	resp, err := client.GetWithQuery(ctx, getURL(client, clusterName), query)
	if err != nil {
		return nil, err
	}

	var result api.V0040OpenapiClustersResp
	if err := slurmclient.DecodeResponse(resp, &result); err != nil {
		return nil, err
	}

	if len(result.Clusters) == 0 {
		return nil, fmt.Errorf("cluster %s not found", clusterName)
	}

	return &result.Clusters[0], nil
}

// Create creates one or more clusters in the Slurm database.
func Create(ctx context.Context, client *slurmclient.ServiceClient, opts CreateOpts) error {
	reqBody := struct {
		Clusters []Cluster `json:"clusters"`
	}{
		Clusters: opts.Clusters,
	}

	resp, err := client.Post(ctx, createURL(client), reqBody)
	if err != nil {
		return err
	}

	return slurmclient.CheckResponse(resp)
}

// Delete removes a cluster from the Slurm database.
func Delete(ctx context.Context, client *slurmclient.ServiceClient, clusterName string, opts ...DeleteOpts) error {
	query := url.Values{}
	if len(opts) > 0 && opts[0].Classification != nil {
		query.Set("classification", *opts[0].Classification)
	}

	resp, err := client.DeleteWithQuery(ctx, deleteURL(client, clusterName), query)
	if err != nil {
		return err
	}

	return slurmclient.CheckResponse(resp)
}
