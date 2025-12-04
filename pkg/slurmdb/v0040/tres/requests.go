// SPDX-FileCopyrightText: Copyright (C) SchedMD LLC.
// SPDX-License-Identifier: Apache-2.0

package tres

import (
	"context"

	api "github.com/supergate-hub/slurm-client/api/v0040"
	"github.com/supergate-hub/slurm-client/pkg/slurmclient"
)

// List retrieves all TRES from the Slurm database.
func List(ctx context.Context, client *slurmclient.ServiceClient, opts ListOpts) ([]TRES, error) {
	resp, err := client.Get(ctx, listURL(client))
	if err != nil {
		return nil, err
	}

	var result api.V0040OpenapiTresResp
	if err := slurmclient.DecodeResponse(resp, &result); err != nil {
		return nil, err
	}

	if len(result.TRES) == 0 {
		return []TRES{}, nil
	}
	return result.TRES, nil
}

// Create creates one or more TRES in the Slurm database.
func Create(ctx context.Context, client *slurmclient.ServiceClient, opts CreateOpts) error {
	reqBody := struct {
		TRES []TRES `json:"tres"`
	}{
		TRES: opts.TRES,
	}

	resp, err := client.Post(ctx, createURL(client), reqBody)
	if err != nil {
		return err
	}

	return slurmclient.CheckResponse(resp)
}
