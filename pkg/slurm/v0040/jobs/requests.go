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

// List retrieves all jobs from the Slurm controller.
func List(ctx context.Context, client *slurmclient.ServiceClient, opts ListOpts) ([]Job, error) {
	query := url.Values{}
	if opts.UpdateTime != nil {
		query.Set("update_time", fmt.Sprintf("%d", *opts.UpdateTime))
	}

	resp, err := client.GetWithQuery(ctx, listURL(client), query)
	if err != nil {
		return nil, err
	}

	var result api.V0040OpenapiJobInfoResp
	if err := slurmclient.DecodeResponse(resp, &result); err != nil {
		return nil, err
	}

	if len(result.Jobs) == 0 {
		return []Job{}, nil
	}
	return result.Jobs, nil
}

// Get retrieves a single job by ID.
func Get(ctx context.Context, client *slurmclient.ServiceClient, jobID string, opts ...GetOpts) (*Job, error) {
	query := url.Values{}
	if len(opts) > 0 && opts[0].UpdateTime != nil {
		query.Set("update_time", fmt.Sprintf("%d", *opts[0].UpdateTime))
	}

	resp, err := client.GetWithQuery(ctx, getURL(client, jobID), query)
	if err != nil {
		return nil, err
	}

	var result api.V0040OpenapiJobInfoResp
	if err := slurmclient.DecodeResponse(resp, &result); err != nil {
		return nil, err
	}

	if len(result.Jobs) == 0 {
		return nil, fmt.Errorf("job %s not found", jobID)
	}

	return &result.Jobs[0], nil
}

// Submit submits a new job to the Slurm controller.
func Submit(ctx context.Context, client *slurmclient.ServiceClient, opts SubmitOpts) (*JobSubmitResponse, error) {
	reqBody := api.V0040JobSubmitReq{}

	if opts.Script != "" {
		reqBody.Script = &opts.Script
	}
	if opts.Job != nil {
		reqBody.Job = opts.Job
	}
	if len(opts.Jobs) > 0 {
		jobsList := api.V0040JobDescMsgList(opts.Jobs)
		reqBody.Jobs = &jobsList
	}

	resp, err := client.Post(ctx, submitURL(client), reqBody)
	if err != nil {
		return nil, err
	}

	var result api.V0040OpenapiJobSubmitResponse
	if err := slurmclient.DecodeResponse(resp, &result); err != nil {
		return nil, err
	}

	return result.Result, nil
}

// Update modifies an existing job.
func Update(ctx context.Context, client *slurmclient.ServiceClient, jobID string, opts UpdateOpts) ([]JobArrayResponse, error) {
	if opts.Job == nil {
		return nil, fmt.Errorf("job description is required for update")
	}

	reqBody := api.V0040JobDescMsg{}
	if opts.Job != nil {
		reqBody = *opts.Job
	}

	resp, err := client.Post(ctx, updateURL(client, jobID), reqBody)
	if err != nil {
		return nil, err
	}

	var result api.V0040OpenapiJobPostResponse
	if err := slurmclient.DecodeResponse(resp, &result); err != nil {
		return nil, err
	}

	if result.Results == nil {
		return []JobArrayResponse{}, nil
	}
	return *result.Results, nil
}

// Delete cancels/removes a job.
func Delete(ctx context.Context, client *slurmclient.ServiceClient, jobID string, opts ...DeleteOpts) error {
	query := url.Values{}
	if len(opts) > 0 {
		if opts[0].Signal != nil {
			query.Set("signal", *opts[0].Signal)
		}
		if opts[0].Flags != nil {
			for _, flag := range *opts[0].Flags {
				query.Add("flags", flag)
			}
		}
	}

	resp, err := client.DeleteWithQuery(ctx, deleteURL(client, jobID), query)
	if err != nil {
		return err
	}

	return slurmclient.CheckResponse(resp)
}

