// SPDX-FileCopyrightText: Copyright (C) SchedMD LLC.
// SPDX-License-Identifier: Apache-2.0

package associations

import (
	"context"
	"net/url"

	api "github.com/supergate-hub/slurm-client/api/v0040"
	"github.com/supergate-hub/slurm-client/pkg/slurmclient"
)

// List retrieves all associations from the Slurm database.
func List(ctx context.Context, client *slurmclient.ServiceClient, opts ListOpts) ([]Association, error) {
	query := url.Values{}
	if opts.Account != nil {
		query.Set("account", *opts.Account)
	}
	if opts.Cluster != nil {
		query.Set("cluster", *opts.Cluster)
	}
	if opts.Partition != nil {
		query.Set("partition", *opts.Partition)
	}
	if opts.User != nil {
		query.Set("user", *opts.User)
	}
	if opts.WithDeleted != nil && *opts.WithDeleted {
		query.Set("with_deleted", "true")
	}

	resp, err := client.GetWithQuery(ctx, listURL(client), query)
	if err != nil {
		return nil, err
	}

	var result api.V0040OpenapiAssocsResp
	if err := slurmclient.DecodeResponse(resp, &result); err != nil {
		return nil, err
	}

	if len(result.Associations) == 0 {
		return []Association{}, nil
	}
	return result.Associations, nil
}

// Get retrieves a single association based on the provided parameters.
func Get(ctx context.Context, client *slurmclient.ServiceClient, opts GetOpts) (*Association, error) {
	query := url.Values{}
	if opts.Account != nil {
		query.Set("account", *opts.Account)
	}
	if opts.Cluster != nil {
		query.Set("cluster", *opts.Cluster)
	}
	if opts.Partition != nil {
		query.Set("partition", *opts.Partition)
	}
	if opts.User != nil {
		query.Set("user", *opts.User)
	}

	resp, err := client.GetWithQuery(ctx, getURL(client), query)
	if err != nil {
		return nil, err
	}

	var result api.V0040OpenapiAssocsResp
	if err := slurmclient.DecodeResponse(resp, &result); err != nil {
		return nil, err
	}

	if len(result.Associations) == 0 {
		return nil, nil
	}

	return &result.Associations[0], nil
}

// Create creates one or more associations in the Slurm database.
func Create(ctx context.Context, client *slurmclient.ServiceClient, opts CreateOpts) error {
	reqBody := struct {
		Associations []Association `json:"associations"`
	}{
		Associations: opts.Associations,
	}

	resp, err := client.Post(ctx, createURL(client), reqBody)
	if err != nil {
		return err
	}

	return slurmclient.CheckResponse(resp)
}

// Delete removes an association from the Slurm database.
func Delete(ctx context.Context, client *slurmclient.ServiceClient, opts DeleteOpts) error {
	query := url.Values{}
	if opts.Account != nil {
		query.Set("account", *opts.Account)
	}
	if opts.Cluster != nil {
		query.Set("cluster", *opts.Cluster)
	}
	if opts.Partition != nil {
		query.Set("partition", *opts.Partition)
	}
	if opts.User != nil {
		query.Set("user", *opts.User)
	}

	resp, err := client.DeleteWithQuery(ctx, deleteURL(client), query)
	if err != nil {
		return err
	}

	return slurmclient.CheckResponse(resp)
}
