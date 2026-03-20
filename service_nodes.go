package slurmclient

import (
	"context"
	"fmt"

	v0040 "github.com/supergate-hub/slurm-client/api/v0040"
)

type nodesServiceImpl struct{ t *Transport }

func (s *nodesServiceImpl) List(ctx context.Context, opts ...NodeListOpts) ([]v0040.V0040Node, error) {
	path := "/nodes"
	if len(opts) > 0 && opts[0].UpdateTime != nil {
		path += fmt.Sprintf("?update_time=%d", *opts[0].UpdateTime)
	}
	var resp v0040.V0040OpenapiNodesResp
	if err := s.t.Get(ctx, path, &resp); err != nil {
		return nil, err
	}
	return resp.Nodes, nil
}

func (s *nodesServiceImpl) Get(ctx context.Context, nodeName string) (*v0040.V0040Node, error) {
	var resp v0040.V0040OpenapiNodesResp
	if err := s.t.Get(ctx, "/node/"+nodeName, &resp); err != nil {
		return nil, err
	}
	if len(resp.Nodes) == 0 {
		return nil, ErrNotFound
	}
	return &resp.Nodes[0], nil
}

func (s *nodesServiceImpl) Delete(ctx context.Context, nodeName string) error {
	return s.t.Delete(ctx, "/node/"+nodeName, nil)
}
