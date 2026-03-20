package slurmclient

import (
	"context"

	v0040 "github.com/supergate-hub/slurm-client/api/v0040"
)

type partitionsServiceImpl struct{ t *Transport }

func (s *partitionsServiceImpl) List(ctx context.Context) ([]v0040.V0040PartitionInfo, error) {
	var resp v0040.V0040OpenapiPartitionResp
	if err := s.t.Get(ctx, "/partitions", &resp); err != nil {
		return nil, err
	}
	return resp.Partitions, nil
}

func (s *partitionsServiceImpl) Get(ctx context.Context, partitionName string) (*v0040.V0040PartitionInfo, error) {
	var resp v0040.V0040OpenapiPartitionResp
	if err := s.t.Get(ctx, "/partition/"+partitionName, &resp); err != nil {
		return nil, err
	}
	if len(resp.Partitions) == 0 {
		return nil, ErrNotFound
	}
	return &resp.Partitions[0], nil
}
