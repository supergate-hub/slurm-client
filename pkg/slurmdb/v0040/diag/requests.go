// SPDX-FileCopyrightText: Copyright (C) SchedMD LLC.
// SPDX-License-Identifier: Apache-2.0

package diag

import (
	"context"

	api "github.com/supergate-hub/slurm-client/api/v0040"
	"github.com/supergate-hub/slurm-client/pkg/slurmclient"
)

// Get retrieves diagnostic information from the Slurm database daemon.
func Get(ctx context.Context, client *slurmclient.ServiceClient) (*DiagInfo, error) {
	resp, err := client.Get(ctx, getURL(client))
	if err != nil {
		return nil, err
	}

	var result api.V0040OpenapiSlurmdbdStatsResp
	if err := slurmclient.DecodeResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result.Statistics, nil
}
