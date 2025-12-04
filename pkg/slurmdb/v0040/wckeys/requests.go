// SPDX-FileCopyrightText: Copyright (C) SchedMD LLC.
// SPDX-License-Identifier: Apache-2.0

package wckeys

import (
	"context"
	"fmt"
	"net/url"

	api "github.com/supergate-hub/slurm-client/api/v0040"
	"github.com/supergate-hub/slurm-client/pkg/slurmclient"
)

// List retrieves all WCKeys from the Slurm database.
func List(ctx context.Context, client *slurmclient.ServiceClient, opts ListOpts) ([]WCKey, error) {
	query := url.Values{}
	if opts.WithDeleted != nil && *opts.WithDeleted {
		query.Set("with_deleted", "true")
	}

	resp, err := client.GetWithQuery(ctx, listURL(client), query)
	if err != nil {
		return nil, err
	}

	var result api.V0040OpenapiWckeyResp
	if err := slurmclient.DecodeResponse(resp, &result); err != nil {
		return nil, err
	}

	if len(result.Wckeys) == 0 {
		return []WCKey{}, nil
	}
	return result.Wckeys, nil
}

// Get retrieves a single WCKey by ID.
func Get(ctx context.Context, client *slurmclient.ServiceClient, wckeyID string) (*WCKey, error) {
	resp, err := client.Get(ctx, getURL(client, wckeyID))
	if err != nil {
		return nil, err
	}

	var result api.V0040OpenapiWckeyResp
	if err := slurmclient.DecodeResponse(resp, &result); err != nil {
		return nil, err
	}

	if len(result.Wckeys) == 0 {
		return nil, fmt.Errorf("wckey %s not found", wckeyID)
	}

	return &result.Wckeys[0], nil
}

// Create creates one or more WCKeys in the Slurm database.
func Create(ctx context.Context, client *slurmclient.ServiceClient, opts CreateOpts) error {
	reqBody := struct {
		WCKeys []WCKey `json:"wckeys"`
	}{
		WCKeys: opts.WCKeys,
	}

	resp, err := client.Post(ctx, createURL(client), reqBody)
	if err != nil {
		return err
	}

	return slurmclient.CheckResponse(resp)
}

// Delete removes a WCKey from the Slurm database.
func Delete(ctx context.Context, client *slurmclient.ServiceClient, wckeyID string) error {
	resp, err := client.Delete(ctx, deleteURL(client, wckeyID))
	if err != nil {
		return err
	}

	return slurmclient.CheckResponse(resp)
}
