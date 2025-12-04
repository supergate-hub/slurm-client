// SPDX-FileCopyrightText: Copyright (C) SchedMD LLC.
// SPDX-License-Identifier: Apache-2.0

package v0040

import (
	"context"

	"github.com/supergate-hub/slurm-client/pkg/slurm/v0040/jobs"
	"github.com/supergate-hub/slurm-client/pkg/slurmclient"
)

// jobsClient implements JobsInterface.
type jobsClient struct {
	client *slurmclient.ServiceClient
}

var _ JobsInterface = (*jobsClient)(nil)

// List retrieves all jobs.
func (c *jobsClient) List(ctx context.Context, opts jobs.ListOpts) ([]Job, error) {
	return jobs.List(ctx, c.client, opts)
}

// Get retrieves a single job by ID.
func (c *jobsClient) Get(ctx context.Context, jobID string, opts ...jobs.GetOpts) (*Job, error) {
	return jobs.Get(ctx, c.client, jobID, opts...)
}

// Submit submits a new job.
func (c *jobsClient) Submit(ctx context.Context, opts jobs.SubmitOpts) (*JobSubmitResponse, error) {
	return jobs.Submit(ctx, c.client, opts)
}

// Update modifies an existing job.
func (c *jobsClient) Update(ctx context.Context, jobID string, opts jobs.UpdateOpts) ([]JobArrayResponse, error) {
	return jobs.Update(ctx, c.client, jobID, opts)
}

// Delete cancels/removes a job.
func (c *jobsClient) Delete(ctx context.Context, jobID string, opts ...jobs.DeleteOpts) error {
	return jobs.Delete(ctx, c.client, jobID, opts...)
}

