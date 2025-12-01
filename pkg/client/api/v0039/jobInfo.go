// SPDX-FileCopyrightText: Copyright (C) SchedMD LLC.
// SPDX-License-Identifier: Apache-2.0

package v0039

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	utilerrors "k8s.io/apimachinery/pkg/util/errors"

	api "github.com/supergate-hub/slurm-client/api/v0039"
	"github.com/supergate-hub/slurm-client/pkg/utils"
)

type JobInfoInterface interface {
	CreateJobInfo(ctx context.Context, req any) (*int32, error)
	// CreateJobInfoAlloc(ctx context.Context, req any) (*int32, error) // Not supported in v0039?
	DeleteJobInfo(ctx context.Context, jobId string) error
	UpdateJobInfo(ctx context.Context, jobId string, req any) error
	GetJobInfo(ctx context.Context, jobId string) (*api.V0039JobInfo, error)
	ListJobInfo(ctx context.Context) ([]api.V0039JobInfo, error)
}

// CreateJobInfo implements JobInfoInterface
func (c *SlurmClient) CreateJobInfo(ctx context.Context, req any) (*int32, error) {
	// Support V0039JobDescMsg for better UX
	if r, ok := req.(api.V0039JobDescMsg); ok {
		payload := map[string]any{
			"job": r,
		}
		if r.Script != nil {
			payload["script"] = *r.Script
		}
		
		buf, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}

		res, err := c.SlurmV0039SubmitJobWithBodyWithResponse(ctx, "application/json", bytes.NewReader(buf))
		if err != nil {
			return nil, err
		}
		if res.StatusCode() != 200 {
			errs := []error{errors.New(http.StatusText(res.StatusCode()))}
			if res.JSON200 != nil && res.JSON200.Errors != nil {
				for _, e := range *res.JSON200.Errors {
					if e.Error != nil {
						errs = append(errs, errors.New(*e.Error))
					}
				}
			}
			return nil, utilerrors.NewAggregate(errs)
		}
		if res.JSON200 == nil || res.JSON200.JobId == nil {
			return nil, errors.New("empty response or missing job id")
		}
		jobId := int32(*res.JSON200.JobId)
		return &jobId, nil
	}

	// Fallback to raw JSON body if user really wants to use it
	if r, ok := req.(api.SlurmV0039SubmitJobJSONBody); ok {
		body := api.SlurmV0039SubmitJobJSONRequestBody(r)
		res, err := c.SlurmV0039SubmitJobWithResponse(ctx, body)
		if err != nil {
			return nil, err
		}
		if res.StatusCode() != 200 {
			errs := []error{errors.New(http.StatusText(res.StatusCode()))}
			if res.JSON200 != nil && res.JSON200.Errors != nil {
				for _, e := range *res.JSON200.Errors {
					if e.Error != nil {
						errs = append(errs, errors.New(*e.Error))
					}
				}
			}
			return nil, utilerrors.NewAggregate(errs)
		}
		if res.JSON200 == nil || res.JSON200.JobId == nil {
			return nil, errors.New("empty response or missing job id")
		}
		jobId := int32(*res.JSON200.JobId)
		return &jobId, nil
	}

	return nil, errors.New("expected req to be api.V0039JobDescMsg or api.SlurmV0039SubmitJobJSONBody")
}

// CreateJobInfoAlloc implements ClientInterface (Stub)
func (c *SlurmClient) CreateJobInfoAlloc(ctx context.Context, req any) (*int32, error) {
	return nil, errors.New("CreateJobInfoAlloc not implemented for v0039")
}

// DeleteJobInfo implements JobInfoInterface
func (c *SlurmClient) DeleteJobInfo(ctx context.Context, jobId string) error {
	// CancelJob in v0039
	res, err := c.SlurmV0039CancelJobWithResponse(ctx, jobId, &api.SlurmV0039CancelJobParams{})
	if err != nil {
		return err
	}
	if res.StatusCode() != 200 {
		errs := []error{errors.New(http.StatusText(res.StatusCode()))}
		if res.JSON200 != nil && res.JSON200.Errors != nil {
			for _, e := range *res.JSON200.Errors {
				if e.Error != nil {
					errs = append(errs, errors.New(*e.Error))
				}
			}
		}
		return utilerrors.NewAggregate(errs)
	}
	return nil
}

// UpdateJobInfo implements JobInfoInterface
func (c *SlurmClient) UpdateJobInfo(ctx context.Context, jobId string, req any) error {
	r, ok := req.(api.SlurmV0039UpdateJobJSONBody)
	if !ok {
		return errors.New("expected req to be SlurmV0039UpdateJobJSONBody")
	}
	body := api.SlurmV0039UpdateJobJSONRequestBody(r)
	res, err := c.SlurmV0039UpdateJobWithResponse(ctx, jobId, body)
	if err != nil {
		return err
	}
	if res.StatusCode() != 200 {
		errs := []error{errors.New(http.StatusText(res.StatusCode()))}
		if res.JSON200 != nil && res.JSON200.Errors != nil {
			for _, e := range *res.JSON200.Errors {
				if e.Error != nil {
					errs = append(errs, errors.New(*e.Error))
				}
			}
		}
		return utilerrors.NewAggregate(errs)
	}
	return nil
}

// GetJobInfo implements JobInfoInterface
func (c *SlurmClient) GetJobInfo(ctx context.Context, jobId string) (*api.V0039JobInfo, error) {
	res, err := c.SlurmV0039GetJobWithResponse(ctx, jobId)
	if err != nil {
		return nil, err
	}
	if res.StatusCode() != 200 {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode())
	}
	if res.JSON200 == nil {
		return nil, fmt.Errorf("empty response body")
	}
	if res.JSON200.Errors != nil && len(*res.JSON200.Errors) > 0 {
		errs := []error{}
		for _, e := range *res.JSON200.Errors {
			if e.Error != nil {
				errs = append(errs, errors.New(*e.Error))
			}
		}
		return nil, errors.Join(errs...)
	}
	if res.JSON200.Jobs == nil || len(*res.JSON200.Jobs) == 0 {
		return nil, nil
	}
	
	out := &api.V0039JobInfo{}
	utils.RemarshalOrDie((*res.JSON200.Jobs)[0], out)
	return out, nil
}

// ListJobInfo implements JobInfoInterface
func (c *SlurmClient) ListJobInfo(ctx context.Context) ([]api.V0039JobInfo, error) {
	res, err := c.SlurmV0039GetJobsWithResponse(ctx, &api.SlurmV0039GetJobsParams{})
	if err != nil {
		return nil, err
	}
	if res.StatusCode() != 200 {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode())
	}
	if res.JSON200 == nil {
		return nil, fmt.Errorf("empty response body")
	}
	if res.JSON200.Errors != nil && len(*res.JSON200.Errors) > 0 {
		errs := []error{}
		for _, e := range *res.JSON200.Errors {
			if e.Error != nil {
				errs = append(errs, errors.New(*e.Error))
			}
		}
		return nil, errors.Join(errs...)
	}
	
	jobs := []api.V0039JobInfo{}
	if res.JSON200.Jobs != nil {
		for _, job := range *res.JSON200.Jobs {
			j := api.V0039JobInfo{}
			utils.RemarshalOrDie(job, &j)
			jobs = append(jobs, j)
		}
	}
	return jobs, nil
}
