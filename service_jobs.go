package slurmclient

import (
	"context"
	"fmt"

	v0040 "github.com/supergate-hub/slurm-client/api/v0040"
)

type jobsServiceImpl struct{ t *Transport }

func (s *jobsServiceImpl) List(ctx context.Context, opts ...JobListOpts) ([]v0040.V0040JobInfo, error) {
	path := "/jobs"
	if len(opts) > 0 && opts[0].UpdateTime != nil {
		path += fmt.Sprintf("?update_time=%d", *opts[0].UpdateTime)
	}
	var resp v0040.V0040OpenapiJobInfoResp
	if err := s.t.Get(ctx, path, &resp); err != nil {
		return nil, err
	}
	return resp.Jobs, nil
}

func (s *jobsServiceImpl) Get(ctx context.Context, jobID string) (*v0040.V0040JobInfo, error) {
	var resp v0040.V0040OpenapiJobInfoResp
	if err := s.t.Get(ctx, "/job/"+jobID, &resp); err != nil {
		return nil, err
	}
	if len(resp.Jobs) == 0 {
		return nil, ErrNotFound
	}
	return &resp.Jobs[0], nil
}

func (s *jobsServiceImpl) Submit(ctx context.Context, opts JobSubmitOpts) (*JobSubmitResult, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}
	body := v0040.V0040JobSubmitReq{
		Script: &opts.Script,
	}
	if opts.Job != nil {
		body.Job = opts.Job
	}
	if len(opts.Jobs) > 0 {
		jobs := v0040.V0040JobDescMsgList(opts.Jobs)
		body.Jobs = &jobs
	}
	var resp v0040.V0040OpenapiJobSubmitResponse
	if err := s.t.Post(ctx, "/job/submit", body, &resp); err != nil {
		return nil, err
	}
	result := &JobSubmitResult{}
	if resp.JobId != nil {
		result.JobID = int64(*resp.JobId)
	}
	if resp.StepId != nil {
		result.StepID = *resp.StepId
	}
	return result, nil
}

func (s *jobsServiceImpl) Update(ctx context.Context, jobID string, opts JobUpdateOpts) error {
	body := opts.Job
	return s.t.Post(ctx, "/job/"+jobID, body, nil)
}

func (s *jobsServiceImpl) Delete(ctx context.Context, jobID string) error {
	return s.t.Delete(ctx, "/job/"+jobID, nil)
}
