// SPDX-FileCopyrightText: Copyright (C) SchedMD LLC.
// SPDX-License-Identifier: Apache-2.0

package jobs

import (
	"context"
	"fmt"
	"net/url"

	api "github.com/supergate-hub/slurm-client/api/v0040"
	"github.com/supergate-hub/slurm-client/pkg/slurmclient"
)

// List retrieves jobs from the Slurm database.
func List(ctx context.Context, client *slurmclient.ServiceClient, opts ListOpts) ([]Job, error) {
	query := url.Values{}
	for _, user := range opts.Users {
		query.Add("users", user)
	}
	for _, account := range opts.Accounts {
		query.Add("accounts", account)
	}
	for _, cluster := range opts.Clusters {
		query.Add("cluster", cluster)
	}
	if opts.StartTime != nil {
		query.Set("start_time", *opts.StartTime)
	}
	if opts.EndTime != nil {
		query.Set("end_time", *opts.EndTime)
	}
	for _, flag := range opts.Flags {
		query.Add("flags", flag)
	}

	resp, err := client.GetWithQuery(ctx, listURL(client), query)
	if err != nil {
		return nil, err
	}

	var result api.V0040OpenapiSlurmdbdJobsResp
	if err := slurmclient.DecodeResponse(resp, &result); err != nil {
		return nil, err
	}

	if len(result.Jobs) == 0 {
		return []Job{}, nil
	}
	return result.Jobs, nil
}

// Get retrieves a single job by ID from the Slurm database.
func Get(ctx context.Context, client *slurmclient.ServiceClient, jobID string) (*Job, error) {
	resp, err := client.Get(ctx, getURL(client, jobID))
	if err != nil {
		return nil, err
	}

	var result api.V0040OpenapiSlurmdbdJobsResp
	if err := slurmclient.DecodeResponse(resp, &result); err != nil {
		return nil, err
	}

	if len(result.Jobs) == 0 {
		return nil, fmt.Errorf("job %s not found", jobID)
	}

	return &result.Jobs[0], nil
}
