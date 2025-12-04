// SPDX-FileCopyrightText: Copyright (C) SchedMD LLC.
// SPDX-License-Identifier: Apache-2.0

package qos

import (
	"context"
	"fmt"
	"net/url"

	api "github.com/supergate-hub/slurm-client/api/v0040"
	"github.com/supergate-hub/slurm-client/pkg/slurmclient"
)

// List retrieves all QOS from the Slurm database.
func List(ctx context.Context, client *slurmclient.ServiceClient, opts ListOpts) ([]QOS, error) {
	query := url.Values{}
	if opts.Description != nil {
		query.Set("description", *opts.Description)
	}
	if opts.WithDeleted != nil && *opts.WithDeleted {
		query.Set("with_deleted", "true")
	}

	resp, err := client.GetWithQuery(ctx, listURL(client), query)
	if err != nil {
		return nil, err
	}

	var result api.V0040OpenapiSlurmdbdQosResp
	if err := slurmclient.DecodeResponse(resp, &result); err != nil {
		return nil, err
	}

	if len(result.Qos) == 0 {
		return []QOS{}, nil
	}
	return result.Qos, nil
}

// Get retrieves a single QOS by name.
func Get(ctx context.Context, client *slurmclient.ServiceClient, qosName string, opts ...GetOpts) (*QOS, error) {
	query := url.Values{}
	if len(opts) > 0 && opts[0].WithDeleted != nil && *opts[0].WithDeleted {
		query.Set("with_deleted", "true")
	}

	resp, err := client.GetWithQuery(ctx, getURL(client, qosName), query)
	if err != nil {
		return nil, err
	}

	var result api.V0040OpenapiSlurmdbdQosResp
	if err := slurmclient.DecodeResponse(resp, &result); err != nil {
		return nil, err
	}

	if len(result.Qos) == 0 {
		return nil, fmt.Errorf("qos %s not found", qosName)
	}

	return &result.Qos[0], nil
}

// Create creates one or more QOS in the Slurm database.
func Create(ctx context.Context, client *slurmclient.ServiceClient, opts CreateOpts) error {
	reqBody := struct {
		QOS []QOS `json:"qos"`
	}{
		QOS: opts.QOS,
	}

	resp, err := client.Post(ctx, createURL(client), reqBody)
	if err != nil {
		return err
	}

	return slurmclient.CheckResponse(resp)
}

// Delete removes a QOS from the Slurm database.
func Delete(ctx context.Context, client *slurmclient.ServiceClient, qosName string) error {
	resp, err := client.Delete(ctx, deleteURL(client, qosName))
	if err != nil {
		return err
	}

	return slurmclient.CheckResponse(resp)
}
